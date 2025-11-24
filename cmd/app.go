package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ydcloud-dy/leaf-api/config"
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"github.com/ydcloud-dy/leaf-api/pkg/logger"
	"github.com/ydcloud-dy/leaf-api/pkg/oss"
	"github.com/ydcloud-dy/leaf-api/pkg/redis"
	"golang.org/x/crypto/bcrypt"
)

// Run 运行应用
func Run(configPath string) error {
	// 加载配置
	if err := config.LoadConfig(configPath); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// 初始化日志
	logger.Init()
	logger.Info("Starting Blog Admin API...")

	// 初始化数据库
	if err := config.InitDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// 自动迁移数据库
	if err := po.AutoMigrate(config.DB); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// 初始化 Redis
	if err := redis.InitRedis(); err != nil {
		logger.Warn("Failed to initialize Redis: ", err)
		logger.Warn("Online user tracking and visit duration recording will be disabled")
	} else {
		logger.Info("Redis connected successfully")
	}

	// 初始化 OSS
	if err := oss.Init(); err != nil {
		logger.Warn("Failed to initialize OSS: ", err)
	}

	// 创建默认管理员
	initDefaultAdmin()

	// 初始化应用（依赖注入）
	app, err := InitApp(config.DB)
	if err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}

	// 创建 HTTP 服务器
	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.GetEngine(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 在 goroutine 中启动服务器
	go func() {
		logger.Info("Server starting on ", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: ", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", err)
	}

	// 关闭数据库连接
	if sqlDB, err := config.DB.DB(); err == nil {
		sqlDB.Close()
	}

	// 关闭 Redis 连接
	if err := redis.Close(); err != nil {
		logger.Error("Failed to close Redis: ", err)
	}

	logger.Info("Server exited gracefully")
	return nil
}

// initDefaultAdmin 创建默认管理员
func initDefaultAdmin() {
	var count int64
	// 检查users表中是否已有管理员
	config.DB.Model(&po.User{}).Where("role IN ?", []string{"admin", "super_admin"}).Count(&count)
	if count > 0 {
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := po.User{
		Username: "admin",
		Password: string(password),
		Email:    "admin@example.com",
		Nickname: "管理员",
		Role:     "admin",
		Status:   1,
	}

	if err := config.DB.Create(&admin).Error; err != nil {
		logger.Error("Failed to create default admin: ", err)
		return
	}

	logger.Info("Default admin created: admin / admin123")
}
