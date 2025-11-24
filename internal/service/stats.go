package service

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/leaf-api/internal/data"
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"github.com/ydcloud-dy/leaf-api/pkg/response"
)

// StatsService 统计服务
type StatsService struct {
	data          *data.Data
	onlineService *OnlineService
	visitService  *VisitService
}

// NewStatsService 创建统计服务
func NewStatsService(d *data.Data) *StatsService {
	return &StatsService{
		data:          d,
		onlineService: NewOnlineService(d),
		visitService:  NewVisitService(d),
	}
}

// GetStats 获取统计数据
func (s *StatsService) GetStats(c *gin.Context) {
	var stats struct {
		ArticleCount       int64   `json:"article_count"`        // 文章总数
		ChapterCount       int64   `json:"chapter_count"`        // 章节/笔记总数
		CategoryCount      int64   `json:"category_count"`       // 文章分类数
		TagCount           int64   `json:"tag_count"`            // 标签数
		UserCount          int64   `json:"user_count"`           // 用户总数
		CommentCount       int64   `json:"comment_count"`        // 评论总数
		TotalViews         int64   `json:"total_views"`          // 总浏览量
		TodayViews         int64   `json:"today_views"`          // 今日浏览量
		OnlineCount        int64   `json:"online_count"`         // 当前在线人数
		AvgVisitDuration   float64 `json:"avg_visit_duration"`   // 平均访问时长（秒）
		SiteRuntime        int64   `json:"site_runtime"`         // 网站运行天数
	}

	// 统计文章数
	s.data.GetDB().Model(&po.Article{}).Where("status = ?", 1).Count(&stats.ArticleCount)

	// 统计章节数（笔记数）
	s.data.GetDB().Model(&po.Chapter{}).Count(&stats.ChapterCount)

	// 统计分类数
	s.data.GetDB().Model(&po.Category{}).Count(&stats.CategoryCount)

	// 统计标签数
	s.data.GetDB().Model(&po.Tag{}).Count(&stats.TagCount)

	// 统计用户数
	s.data.GetDB().Model(&po.User{}).Count(&stats.UserCount)

	// 统计评论数
	s.data.GetDB().Model(&po.Comment{}).Where("status = ?", 1).Count(&stats.CommentCount)

	// 统计总浏览量（所有文章的浏览量之和）
	s.data.GetDB().Model(&po.Article{}).Select("COALESCE(SUM(view_count), 0)").Row().Scan(&stats.TotalViews)

	// 统计24小时访问量（PV）- 从 page_visits 表统计
	pv24h, _ := s.visitService.Get24HourPageViews()
	stats.TodayViews = pv24h

	// 获取在线人数
	onlineCount, _ := s.onlineService.GetOnlineCount()
	stats.OnlineCount = onlineCount

	// 获取平均访问时长
	avgDuration, _ := s.visitService.GetAverageVisitDuration()
	stats.AvgVisitDuration = avgDuration

	// 计算网站运行天数（从 settings 表读取网站启动时间）
	var setting po.Setting
	err := s.data.GetDB().Where("`key` = ?", "site_start_date").First(&setting).Error
	if err == nil {
		// 解析启动时间
		if startTime, err := time.Parse("2006-01-02", setting.Value); err == nil {
			stats.SiteRuntime = int64(time.Since(startTime).Hours() / 24)
		}
	}

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
