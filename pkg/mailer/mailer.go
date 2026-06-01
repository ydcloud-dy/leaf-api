package mailer

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

const defaultRecipientChunkSize = 100

// Config 邮件发送配置
type Config struct {
	Host                    string
	Port                    int
	Username                string
	Password                string
	From                    string
	FromName                string
	UseSSL                  bool
	Timeout                 time.Duration
	MaxRecipientsPerMessage int
}

// Client SMTP 邮件客户端
type Client struct {
	cfg Config
}

// New 创建邮件客户端
func New(cfg Config) *Client {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 10 * time.Second
	}
	if cfg.MaxRecipientsPerMessage <= 0 {
		cfg.MaxRecipientsPerMessage = defaultRecipientChunkSize
	}
	return &Client{cfg: cfg}
}

// Configured 判断是否具备发送条件
func (c *Client) Configured() bool {
	return strings.TrimSpace(c.cfg.Host) != "" &&
		strings.TrimSpace(c.cfg.From) != "" &&
		c.cfg.Port > 0
}

// SendHTML 发送 HTML 邮件
func (c *Client) SendHTML(ctx context.Context, recipients []string, subject, htmlBody string) error {
	if !c.Configured() {
		return errors.New("mailer is not configured")
	}

	addresses := normalizeRecipients(recipients)
	if len(addresses) == 0 {
		return errors.New("no valid recipients")
	}

	for _, chunk := range chunkStrings(addresses, c.cfg.MaxRecipientsPerMessage) {
		if err := c.sendChunk(ctx, chunk, subject, htmlBody); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) sendChunk(ctx context.Context, recipients []string, subject, htmlBody string) error {
	addr := net.JoinHostPort(c.cfg.Host, fmt.Sprintf("%d", c.cfg.Port))
	serverName := c.cfg.Host

	dialer := &net.Dialer{}
	if deadline, ok := ctx.Deadline(); ok {
		dialer.Timeout = time.Until(deadline)
	} else {
		dialer.Timeout = c.cfg.Timeout
	}

	var (
		client *smtp.Client
		conn   net.Conn
		err    error
	)

	if c.cfg.UseSSL {
		tlsConn, dialErr := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
			ServerName:         serverName,
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: false,
		})
		if dialErr != nil {
			return dialErr
		}
		conn = tlsConn
		client, err = smtp.NewClient(conn, serverName)
		if err != nil {
			_ = conn.Close()
			return err
		}
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
		if err != nil {
			return err
		}
		client, err = smtp.NewClient(conn, serverName)
		if err != nil {
			_ = conn.Close()
			return err
		}

		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(&tls.Config{
				ServerName:         serverName,
				MinVersion:         tls.VersionTLS12,
				InsecureSkipVerify: false,
			}); err != nil {
				_ = client.Close()
				return err
			}
		}
	}
	defer client.Close()

	if c.cfg.Username != "" || c.cfg.Password != "" {
		if err := client.Auth(smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, serverName)); err != nil {
			return err
		}
	}

	if err := client.Mail(c.cfg.From); err != nil {
		return err
	}
	for _, recipient := range recipients {
		if err := client.Rcpt(recipient); err != nil {
			return err
		}
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	if _, err := writer.Write(buildMessage(c.cfg.From, c.cfg.FromName, recipients, subject, htmlBody)); err != nil {
		_ = writer.Close()
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func buildMessage(from, fromName string, recipients []string, subject, htmlBody string) []byte {
	var buf bytes.Buffer

	headers := map[string]string{
		"From":                      formatAddress(fromName, from),
		"To":                        "undisclosed-recipients:;",
		"Subject":                   mime.QEncoding.Encode("UTF-8", subject),
		"MIME-Version":              "1.0",
		"Content-Type":              "text/html; charset=UTF-8",
		"Content-Transfer-Encoding": "8bit",
		"Date":                      time.Now().Format(time.RFC1123Z),
		"Message-ID":                fmt.Sprintf("<%d@%s>", time.Now().UnixNano(), messageDomain(from)),
	}

	for k, v := range headers {
		buf.WriteString(k)
		buf.WriteString(": ")
		buf.WriteString(v)
		buf.WriteString("\r\n")
	}
	buf.WriteString("\r\n")
	buf.WriteString(htmlBody)

	return buf.Bytes()
}

func formatAddress(name, address string) string {
	if strings.TrimSpace(name) == "" {
		return address
	}
	return (&mail.Address{Name: name, Address: address}).String()
}

func messageDomain(from string) string {
	parts := strings.SplitN(from, "@", 2)
	if len(parts) == 2 && parts[1] != "" {
		return parts[1]
	}
	return "localhost"
}

func normalizeRecipients(recipients []string) []string {
	seen := make(map[string]struct{}, len(recipients))
	result := make([]string, 0, len(recipients))

	for _, recipient := range recipients {
		recipient = strings.TrimSpace(recipient)
		if recipient == "" {
			continue
		}

		parsed, err := mail.ParseAddress(recipient)
		if err != nil {
			continue
		}

		key := strings.ToLower(parsed.Address)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, parsed.Address)
	}

	return result
}

func chunkStrings(items []string, size int) [][]string {
	if size <= 0 {
		size = defaultRecipientChunkSize
	}

	var chunks [][]string
	for start := 0; start < len(items); start += size {
		end := start + size
		if end > len(items) {
			end = len(items)
		}
		chunk := make([]string, end-start)
		copy(chunk, items[start:end])
		chunks = append(chunks, chunk)
	}
	return chunks
}
