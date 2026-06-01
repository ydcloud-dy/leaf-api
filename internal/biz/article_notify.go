package biz

import (
	"context"
	"fmt"
	"html"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"github.com/ydcloud-dy/leaf-api/pkg/logger"
	"github.com/ydcloud-dy/leaf-api/pkg/mailer"
)

const (
	defaultSiteName = "运维工程师的技术笔记"
	defaultSiteURL  = "https://dycloud.fun"
)

type publishNoticeArticle struct {
	ID          uint
	Title       string
	Summary     string
	Category    string
	Tags        []string
	PublishedAt time.Time
}

type publishNoticeConfig struct {
	Enabled  bool
	SiteName string
	SiteURL  string
	Mailer   mailer.Config
}

func (uc *articleUseCase) notifyPublishedArticle(articleID uint) {
	uc.notifyPublishedArticlesByIDs([]uint{articleID})
}

func (uc *articleUseCase) notifyPublishedArticlesByIDs(articleIDs []uint) {
	if len(articleIDs) == 0 {
		return
	}

	articles, err := uc.data.ArticleRepo.FindByIDs(uniqueUintIDs(articleIDs))
	if err != nil || len(articles) == 0 {
		if err != nil {
			logger.WithFields(logrus.Fields{
				"article_ids": articleIDs,
				"error":       err.Error(),
			}).Warn("加载发布通知文章失败")
		}
		return
	}

	snapshots := make([]publishNoticeArticle, 0, len(articles))
	publishedAt := time.Now()
	for _, article := range articles {
		snapshots = append(snapshots, snapshotPublishedArticle(article, publishedAt))
	}

	go func(items []publishNoticeArticle) {
		if err := uc.sendPublishNotification(items); err != nil {
			logger.WithFields(logrus.Fields{
				"article_count": len(items),
				"error":         err.Error(),
			}).Warn("文章发布邮件通知发送失败")
		}
	}(snapshots)
}

func (uc *articleUseCase) sendPublishNotification(items []publishNoticeArticle) error {
	if len(items) == 0 {
		return nil
	}

	cfg, err := uc.loadPublishNoticeConfig()
	if err != nil {
		return err
	}
	if !cfg.Enabled {
		return nil
	}

	recipients, err := uc.loadNotificationRecipients(cfg.Mailer.From, cfg.Mailer.Username)
	if err != nil {
		return err
	}
	if len(recipients) == 0 {
		return nil
	}

	client := mailer.New(cfg.Mailer)
	if !client.Configured() {
		return nil
	}

	subject, body := renderPublishNotification(cfg.SiteName, cfg.SiteURL, items)
	if err := client.SendHTML(context.Background(), recipients, subject, body); err != nil {
		return err
	}

	logger.WithFields(logrus.Fields{
		"article_count":   len(items),
		"recipient_count": len(recipients),
	}).Info("文章发布邮件通知发送成功")

	return nil
}

func (uc *articleUseCase) loadPublishNoticeConfig() (*publishNoticeConfig, error) {
	settings, err := uc.data.SettingRepo.List()
	if err != nil {
		return nil, err
	}

	values := make(map[string]string, len(settings))
	for _, setting := range settings {
		values[strings.ToLower(strings.TrimSpace(setting.Key))] = strings.TrimSpace(setting.Value)
	}

	siteName := firstNonEmpty(
		settingValue(values, "site_name", "site.title", "blog_name"),
		envValue("SITE_NAME", "BLOG_NAME"),
		defaultSiteName,
	)
	siteURL := normalizeBaseURL(firstNonEmpty(
		settingValue(values, "site_url", "site.url", "blog_url"),
		envValue("SITE_URL", "BLOG_URL"),
		defaultSiteURL,
	))

	host := firstNonEmpty(
		settingValue(values, "mail_smtp_host", "smtp_host", "mail.smtp.host"),
		envValue("MAIL_SMTP_HOST", "SMTP_HOST"),
	)
	username := firstNonEmpty(
		settingValue(values, "mail_smtp_username", "smtp_username", "mail.smtp.username"),
		envValue("MAIL_SMTP_USERNAME", "SMTP_USERNAME"),
	)
	password := firstNonEmpty(
		settingValue(values, "mail_smtp_password", "smtp_password", "mail.smtp.password"),
		envValue("MAIL_SMTP_PASSWORD", "SMTP_PASSWORD"),
	)
	from := firstNonEmpty(
		settingValue(values, "mail_smtp_from", "mail_from", "smtp_from"),
		envValue("MAIL_SMTP_FROM", "MAIL_FROM", "SMTP_FROM"),
	)
	if from == "" {
		from = username
	}
	fromName := firstNonEmpty(
		settingValue(values, "mail_smtp_from_name", "mail_from_name", "smtp_from_name"),
		envValue("MAIL_SMTP_FROM_NAME", "MAIL_FROM_NAME"),
		siteName,
	)

	portValue := firstNonEmpty(
		settingValue(values, "mail_smtp_port", "smtp_port", "mail.smtp.port"),
		envValue("MAIL_SMTP_PORT", "SMTP_PORT"),
	)
	port, _ := strconv.Atoi(portValue)

	useSSL, useSSLSet := parseOptionalBool(firstNonEmpty(
		settingValue(values, "mail_smtp_use_ssl", "smtp_use_ssl", "mail.smtp.use_ssl"),
		envValue("MAIL_SMTP_USE_SSL", "SMTP_USE_SSL"),
	))
	if !useSSLSet {
		useSSL = port == 465
	}
	if port <= 0 {
		if useSSL {
			port = 465
		} else {
			port = 587
		}
	}

	timeoutSeconds, _ := strconv.Atoi(firstNonEmpty(
		settingValue(values, "mail_smtp_timeout_seconds", "smtp_timeout_seconds"),
		envValue("MAIL_SMTP_TIMEOUT_SECONDS", "SMTP_TIMEOUT_SECONDS"),
	))
	maxRecipients, _ := strconv.Atoi(firstNonEmpty(
		settingValue(values, "mail_smtp_max_recipients", "smtp_max_recipients"),
		envValue("MAIL_SMTP_MAX_RECIPIENTS", "SMTP_MAX_RECIPIENTS"),
	))

	explicitEnabled, enabledSet := parseOptionalBool(firstNonEmpty(
		settingValue(values, "mail_enabled", "mail.enable", "smtp_enabled"),
		envValue("MAIL_ENABLED", "SMTP_ENABLED"),
	))
	enabled := explicitEnabled
	if !enabledSet {
		enabled = host != "" && from != ""
	}

	mailCfg := mailer.Config{
		Host:                    host,
		Port:                    port,
		Username:                username,
		Password:                password,
		From:                    from,
		FromName:                fromName,
		UseSSL:                  useSSL,
		MaxRecipientsPerMessage: maxRecipients,
	}
	if timeoutSeconds > 0 {
		mailCfg.Timeout = time.Duration(timeoutSeconds) * time.Second
	}

	return &publishNoticeConfig{
		Enabled:  enabled,
		SiteName: siteName,
		SiteURL:  siteURL,
		Mailer:   mailCfg,
	}, nil
}

func (uc *articleUseCase) loadNotificationRecipients(excludedEmails ...string) ([]string, error) {
	users, err := uc.data.UserRepo.ListActiveRegisteredWithEmail()
	if err != nil {
		return nil, err
	}

	recipients := make([]string, 0, len(users))
	seen := make(map[string]struct{}, len(users))
	excluded := make(map[string]struct{}, len(excludedEmails)+1)
	excluded["admin@example.com"] = struct{}{}
	for _, email := range excludedEmails {
		email = strings.TrimSpace(email)
		if email == "" {
			continue
		}
		excluded[strings.ToLower(email)] = struct{}{}
	}

	for _, user := range users {
		email := strings.TrimSpace(user.Email)
		if email == "" {
			continue
		}
		key := strings.ToLower(email)
		if _, skip := excluded[key]; skip {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		recipients = append(recipients, email)
	}

	return recipients, nil
}

func renderPublishNotification(siteName, siteURL string, items []publishNoticeArticle) (string, string) {
	if len(items) == 1 {
		subject := fmt.Sprintf("【%s】新文章发布：%s", siteName, items[0].Title)
		return subject, renderSingleArticleNotification(siteName, siteURL, items[0])
	}

	subject := fmt.Sprintf("【%s】有 %d 篇新文章发布", siteName, len(items))
	return subject, renderBatchArticleNotification(siteName, siteURL, items)
}

func renderSingleArticleNotification(siteName, siteURL string, item publishNoticeArticle) string {
	link := articleLink(siteURL, item.ID)
	meta := renderArticleMeta(item)

	var builder strings.Builder
	builder.WriteString(`<!doctype html><html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">`)
	builder.WriteString(`<title>`)
	builder.WriteString(html.EscapeString(siteName))
	builder.WriteString(`</title></head><body style="margin:0;padding:0;background:#f5f7fb;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,'PingFang SC','Hiragino Sans GB','Microsoft YaHei',sans-serif;color:#1f2937;">`)
	builder.WriteString(`<div style="max-width:720px;margin:0 auto;padding:32px 16px;">`)
	builder.WriteString(`<div style="background:#fff;border:1px solid #e5e7eb;border-radius:16px;overflow:hidden;box-shadow:0 12px 32px rgba(15,23,42,.08);">`)
	builder.WriteString(`<div style="padding:28px 28px 20px;background:linear-gradient(135deg,#0f172a 0%,#1d4ed8 100%);color:#fff;">`)
	builder.WriteString(`<div style="font-size:13px;opacity:.85;letter-spacing:.08em;text-transform:uppercase;">NEW ARTICLE</div>`)
	builder.WriteString(`<h1 style="margin:10px 0 0;font-size:28px;line-height:1.35;">`)
	builder.WriteString(html.EscapeString(item.Title))
	builder.WriteString(`</h1></div>`)
	builder.WriteString(`<div style="padding:24px 28px 28px;">`)
	if item.Summary != "" {
		builder.WriteString(`<p style="margin:0 0 18px;font-size:15px;line-height:1.8;color:#374151;">`)
		builder.WriteString(html.EscapeString(item.Summary))
		builder.WriteString(`</p>`)
	}
	builder.WriteString(`<div style="padding:14px 16px;border-radius:12px;background:#f9fafb;border:1px solid #e5e7eb;margin-bottom:20px;">`)
	builder.WriteString(meta)
	builder.WriteString(`</div>`)
	builder.WriteString(`<a href="`)
	builder.WriteString(html.EscapeString(link))
	builder.WriteString(`" style="display:inline-block;padding:12px 22px;border-radius:999px;background:#2563eb;color:#fff;text-decoration:none;font-weight:700;">阅读全文</a>`)
	builder.WriteString(`</div></div>`)
	builder.WriteString(`<div style="padding:16px 8px 0;font-size:12px;color:#6b7280;line-height:1.8;">`)
	builder.WriteString(`你收到这封邮件，是因为你是 `)
	builder.WriteString(html.EscapeString(siteName))
	builder.WriteString(` 的注册用户。`)
	builder.WriteString(`</div></div></body></html>`)
	return builder.String()
}

func renderBatchArticleNotification(siteName, siteURL string, items []publishNoticeArticle) string {
	maxVisible := 10
	visible := items
	hiddenCount := 0
	if len(items) > maxVisible {
		visible = items[:maxVisible]
		hiddenCount = len(items) - maxVisible
	}

	var builder strings.Builder
	builder.WriteString(`<!doctype html><html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">`)
	builder.WriteString(`<title>`)
	builder.WriteString(html.EscapeString(siteName))
	builder.WriteString(`</title></head><body style="margin:0;padding:0;background:#f5f7fb;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,'PingFang SC','Hiragino Sans GB','Microsoft YaHei',sans-serif;color:#1f2937;">`)
	builder.WriteString(`<div style="max-width:760px;margin:0 auto;padding:32px 16px;">`)
	builder.WriteString(`<div style="background:#fff;border:1px solid #e5e7eb;border-radius:16px;overflow:hidden;box-shadow:0 12px 32px rgba(15,23,42,.08);">`)
	builder.WriteString(`<div style="padding:28px;background:linear-gradient(135deg,#111827 0%,#2563eb 100%);color:#fff;">`)
	builder.WriteString(`<div style="font-size:13px;opacity:.85;letter-spacing:.08em;text-transform:uppercase;">PUBLISHED ARTICLES</div>`)
	builder.WriteString(`<h1 style="margin:10px 0 0;font-size:28px;line-height:1.35;">`)
	builder.WriteString(html.EscapeString(fmt.Sprintf("新发布了 %d 篇文章", len(items))))
	builder.WriteString(`</h1>`)
	builder.WriteString(`<p style="margin:12px 0 0;font-size:14px;line-height:1.8;opacity:.92;">`)
	builder.WriteString(html.EscapeString(siteName))
	builder.WriteString(` 刚刚更新了一批内容，下面是本次发布摘要。`)
	builder.WriteString(`</p></div>`)
	builder.WriteString(`<div style="padding:22px 28px 28px;">`)
	builder.WriteString(`<div style="display:grid;gap:14px;">`)
	for _, item := range visible {
		link := articleLink(siteURL, item.ID)
		builder.WriteString(`<div style="padding:16px;border-radius:14px;border:1px solid #e5e7eb;background:#f9fafb;">`)
		builder.WriteString(`<div style="display:flex;justify-content:space-between;gap:12px;align-items:flex-start;">`)
		builder.WriteString(`<div style="min-width:0;flex:1;">`)
		builder.WriteString(`<a href="`)
		builder.WriteString(html.EscapeString(link))
		builder.WriteString(`" style="color:#111827;text-decoration:none;font-size:16px;font-weight:700;line-height:1.6;">`)
		builder.WriteString(html.EscapeString(item.Title))
		builder.WriteString(`</a>`)
		builder.WriteString(`<div style="margin-top:8px;font-size:13px;line-height:1.8;color:#6b7280;">`)
		builder.WriteString(renderArticleMeta(item))
		builder.WriteString(`</div>`)
		if item.Summary != "" {
			builder.WriteString(`<p style="margin:10px 0 0;font-size:14px;line-height:1.8;color:#374151;">`)
			builder.WriteString(html.EscapeString(limitString(item.Summary, 180)))
			builder.WriteString(`</p>`)
		}
		builder.WriteString(`</div>`)
		builder.WriteString(`<a href="`)
		builder.WriteString(html.EscapeString(link))
		builder.WriteString(`" style="white-space:nowrap;align-self:center;padding:9px 14px;border-radius:999px;background:#2563eb;color:#fff;text-decoration:none;font-size:13px;font-weight:700;">阅读</a>`)
		builder.WriteString(`</div></div>`)
	}
	builder.WriteString(`</div>`)
	if hiddenCount > 0 {
		builder.WriteString(`<div style="margin-top:18px;padding:14px 16px;border-radius:12px;background:#eff6ff;color:#1d4ed8;font-size:14px;line-height:1.8;">`)
		builder.WriteString(html.EscapeString(fmt.Sprintf("还有 %d 篇文章未展开显示，访问站点可以查看完整列表。", hiddenCount)))
		builder.WriteString(`</div>`)
	}
	builder.WriteString(`<div style="margin-top:22px;">`)
	builder.WriteString(`<a href="`)
	builder.WriteString(html.EscapeString(strings.TrimRight(siteURL, "/")))
	builder.WriteString(`" style="display:inline-block;padding:12px 22px;border-radius:999px;background:#111827;color:#fff;text-decoration:none;font-weight:700;">访问站点</a>`)
	builder.WriteString(`</div></div></div>`)
	builder.WriteString(`<div style="padding:16px 8px 0;font-size:12px;color:#6b7280;line-height:1.8;">`)
	builder.WriteString(`你收到这封邮件，是因为你是 `)
	builder.WriteString(html.EscapeString(siteName))
	builder.WriteString(` 的注册用户。`)
	builder.WriteString(`</div></div></body></html>`)
	return builder.String()
}

func renderArticleMeta(item publishNoticeArticle) string {
	parts := make([]string, 0, 3)
	if item.Category != "" {
		parts = append(parts, "分类："+html.EscapeString(item.Category))
	}
	if len(item.Tags) > 0 {
		parts = append(parts, "标签："+html.EscapeString(strings.Join(item.Tags, " / ")))
	}
	if !item.PublishedAt.IsZero() {
		parts = append(parts, "发布时间："+item.PublishedAt.Format("2006-01-02 15:04"))
	}
	if len(parts) == 0 {
		return "暂无补充信息"
	}
	return strings.Join(parts, " · ")
}

func articleLink(siteURL string, id uint) string {
	return strings.TrimRight(siteURL, "/") + "/articles/" + strconv.FormatUint(uint64(id), 10)
}

func snapshotPublishedArticle(article *po.Article, publishedAt time.Time) publishNoticeArticle {
	if article == nil {
		return publishNoticeArticle{PublishedAt: publishedAt}
	}

	tags := make([]string, 0, len(article.Tags))
	for _, tag := range article.Tags {
		if strings.TrimSpace(tag.Name) != "" {
			tags = append(tags, tag.Name)
		}
	}

	category := ""
	if strings.TrimSpace(article.Category.Name) != "" {
		category = article.Category.Name
	}

	return publishNoticeArticle{
		ID:          article.ID,
		Title:       article.Title,
		Summary:     article.Summary,
		Category:    category,
		Tags:        tags,
		PublishedAt: publishedAt,
	}
}

func uniqueUintIDs(ids []uint) []uint {
	seen := make(map[uint]struct{}, len(ids))
	result := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

func settingValue(values map[string]string, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(values[strings.ToLower(key)]); value != "" {
			return value
		}
	}
	return ""
}

func envValue(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func normalizeBaseURL(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return defaultSiteURL
	}
	return strings.TrimRight(value, "/")
}

func parseOptionalBool(value string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "y", "on", "enabled":
		return true, true
	case "0", "false", "no", "n", "off", "disabled":
		return false, true
	default:
		return false, false
	}
}

func limitString(value string, limit int) string {
	runes := []rune(strings.TrimSpace(value))
	if len(runes) <= limit {
		return strings.TrimSpace(value)
	}
	return string(runes[:limit]) + "..."
}
