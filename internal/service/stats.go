package service

import (
	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/leaf-api/internal/data"
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"github.com/ydcloud-dy/leaf-api/pkg/response"
)

// StatsService 统计服务
type StatsService struct {
	data *data.Data
}

// NewStatsService 创建统计服务
func NewStatsService(d *data.Data) *StatsService {
	return &StatsService{
		data: d,
	}
}

// GetStats 获取统计数据
func (s *StatsService) GetStats(c *gin.Context) {
	var stats struct {
		ArticleCount  int64 `json:"article_count"`
		UserCount     int64 `json:"user_count"`
		CommentCount  int64 `json:"comment_count"`
		TotalViews    int64 `json:"total_views"`
		TodayViews    int64 `json:"today_views"`
	}

	// 统计文章数
	s.data.GetDB().Model(&po.Article{}).Count(&stats.ArticleCount)

	// 统计用户数
	s.data.GetDB().Model(&po.User{}).Count(&stats.UserCount)

	// 统计评论数
	s.data.GetDB().Model(&po.Comment{}).Where("status = ?", 1).Count(&stats.CommentCount)

	// 统计总浏览量
	s.data.GetDB().Model(&po.View{}).Count(&stats.TotalViews)

	// 统计今日浏览量
	todayViews, _ := s.data.ViewRepo.CountToday()
	stats.TodayViews = todayViews

	response.Success(c, stats)
}

// GetHotArticles 获取热门文章
func (s *StatsService) GetHotArticles(c *gin.Context) {
	var articles []po.Article
	s.data.GetDB().Model(&po.Article{}).
		Preload("Category").
		Where("status = ?", 1).
		Order("view_count DESC").
		Limit(10).
		Find(&articles)

	response.Success(c, articles)
}
