package logger

import (
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/ydcloud-dy/leaf-api/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

func Init() {
	Log = logrus.New()

	// Set log level
	switch config.AppConfig.Log.Level {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}

	// Set log format
	if config.AppConfig.Log.Format == "json" {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
	}

	// Set output
	if config.AppConfig.Log.Output == "file" && config.AppConfig.Log.FilePath != "" {
		writer := &lumberjack.Logger{
			Filename:   config.AppConfig.Log.FilePath,
			MaxSize:    config.AppConfig.Log.MaxSize,
			MaxBackups: config.AppConfig.Log.MaxBackups,
			MaxAge:     config.AppConfig.Log.MaxAge,
			Compress:   true,
		}
		Log.SetOutput(io.MultiWriter(os.Stdout, writer))
	} else {
		Log.SetOutput(os.Stdout)
	}
}

// GinLogger returns a gin middleware for logging requests
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		entry := Log.WithFields(logrus.Fields{
			"status":     statusCode,
			"latency":    latency.String(),
			"client_ip":  clientIP,
			"method":     method,
			"path":       path,
			"query":      query,
			"user_agent": c.Request.UserAgent(),
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else if statusCode >= 500 {
			entry.Error("Server error")
		} else if statusCode >= 400 {
			entry.Warn("Client error")
		} else {
			entry.Info("Request completed")
		}
	}
}

// GinRecovery returns a gin middleware for recovering from panics
func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				Log.WithFields(logrus.Fields{
					"error":  err,
					"path":   c.Request.URL.Path,
					"method": c.Request.Method,
				}).Error("Panic recovered")

				c.AbortWithStatusJSON(500, gin.H{
					"code":    500,
					"message": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}

func Debug(args ...interface{}) {
	Log.Debug(args...)
}

func Info(args ...interface{}) {
	Log.Info(args...)
}

func Warn(args ...interface{}) {
	Log.Warn(args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}
