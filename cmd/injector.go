//go:build !wireinject
// +build !wireinject

package cmd

import (
	"github.com/ydcloud-dy/leaf-api/internal/biz"
	"github.com/ydcloud-dy/leaf-api/internal/data"
	"github.com/ydcloud-dy/leaf-api/internal/server"
	"gorm.io/gorm"
)

// InitApp 初始化应用（手动依赖注入）
func InitApp(db *gorm.DB) (*server.HTTPServer, error) {
	// 初始化 Data 层
	dataLayer, err := data.NewData(db)
	if err != nil {
		return nil, err
	}

	// 初始化 Biz 层
	bizLayer := biz.NewBiz(dataLayer)

	// 初始化 Server 层
	httpServer := server.NewHTTPServer(bizLayer, dataLayer)

	return httpServer, nil
}
