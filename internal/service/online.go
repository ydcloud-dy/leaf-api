package service

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/leaf-api/internal/data"
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"github.com/ydcloud-dy/leaf-api/pkg/redis"
	"github.com/ydcloud-dy/leaf-api/pkg/response"
)

const (
	// 在线用户 Redis Key 前缀
	onlineUserPrefix = "online:user:"
	// 在线游客 Redis Key 前缀
	onlineGuestPrefix = "online:guest:"
	// 在线用户过期时间（60秒）
	onlineUserExpire = 60 * time.Second
)

// OnlineService 在线用户追踪服务
type OnlineService struct {
	data *data.Data
}

// NewOnlineService 创建在线用户追踪服务
func NewOnlineService(d *data.Data) *OnlineService {
	return &OnlineService{
		data: d,
	}
}

// RecordHeartbeat 记录用户心跳（保持在线状态）
// @Summary 记录用户心跳
// @Description 记录用户在线状态，登录用户按UserID追踪，未登录按IP追踪
// @Tags 在线追踪
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "记录成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /blog/heartbeat [post]
func (s *OnlineService) RecordHeartbeat(c *gin.Context) {
	// 获取用户ID（如果已登录）
	userIDValue, exists := c.Get("user_id")

	var key string
	if exists && userIDValue != nil {
		userID := userIDValue.(uint)
		key = fmt.Sprintf("%s%d", onlineUserPrefix, userID)
	} else {
		// 未登录用户，使用 IP 作为标识
		ip := c.ClientIP()
		key = fmt.Sprintf("%s%s", onlineGuestPrefix, ip)
	}

	// 设置带过期时间的键
	err := redis.SetWithExpire(key, time.Now().Unix(), onlineUserExpire)
	if err != nil {
		response.Error(c, 500, "记录在线状态失败")
		return
	}

	response.Success(c, gin.H{"status": "ok"})
}

// GetOnlineCount 获取在线人数
func (s *OnlineService) GetOnlineCount() (int64, error) {
	// 查询所有在线用户键
	userKeys, err := redis.Keys(onlineUserPrefix + "*")
	if err != nil {
		return 0, err
	}

	// 查询所有在线游客键
	guestKeys, err := redis.Keys(onlineGuestPrefix + "*")
	if err != nil {
		return 0, err
	}

	return int64(len(userKeys) + len(guestKeys)), nil
}

// VisitService 页面访问时长记录服务
type VisitService struct {
	data *data.Data
}

// NewVisitService 创建访问时长记录服务
func NewVisitService(d *data.Data) *VisitService {
	return &VisitService{
		data: d,
	}
}

// RecordVisitDuration 记录页面访问时长
// @Summary 记录页面访问时长
// @Description 记录用户访问页面的时长信息，用于统计分析
// @Tags 在线追踪
// @Accept json
// @Produce json
// @Param request body object{path=string,duration=int} true "访问信息 path:页面路径 duration:停留时长(秒)"
// @Success 200 {object} response.Response "记录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /blog/visit [post]
func (s *VisitService) RecordVisitDuration(c *gin.Context) {
	var req struct {
		Path     string `json:"path" binding:"required"`
		Duration int    `json:"duration" binding:"min=0"` // 秒，0表示刚进入页面
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	// 获取用户ID（如果已登录）
	var userID *uint
	if userIDValue, exists := c.Get("user_id"); exists && userIDValue != nil {
		uid := userIDValue.(uint)
		userID = &uid
	}

	// 创建访问记录
	visit := &po.PageVisit{
		UserID:    userID,
		IP:        c.ClientIP(),
		Path:      req.Path,
		Duration:  req.Duration,
		UserAgent: c.GetHeader("User-Agent"),
		Referrer:  c.GetHeader("Referer"),
		CreatedAt: time.Now(),
	}

	if err := s.data.GetDB().Create(visit).Error; err != nil {
		response.Error(c, 500, "记录访问时长失败")
		return
	}

	response.Success(c, gin.H{"status": "ok"})
}

// GetAverageVisitDuration 获取平均访问时长（秒）
// 只统计 duration > 0 的记录（排除刚进入页面的记录）
func (s *VisitService) GetAverageVisitDuration() (float64, error) {
	var avgDuration float64
	err := s.data.GetDB().Model(&po.PageVisit{}).
		Select("COALESCE(AVG(duration), 0)").
		Where("created_at >= ? AND duration > 0", time.Now().Add(-24*time.Hour)).
		Row().
		Scan(&avgDuration)

	return avgDuration, err
}

// GetAverageVisitDurationByPath 按页面路径获取平均访问时长
func (s *VisitService) GetAverageVisitDurationByPath(path string) (float64, error) {
	var avgDuration float64
	err := s.data.GetDB().Model(&po.PageVisit{}).
		Select("COALESCE(AVG(duration), 0)").
		Where("path = ? AND created_at >= ?", path, time.Now().Add(-24*time.Hour)).
		Row().
		Scan(&avgDuration)

	return avgDuration, err
}

// Get24HourPageViews 获取最近24小时的页面访问量（PV）
func (s *VisitService) Get24HourPageViews() (int64, error) {
	var count int64
	err := s.data.GetDB().Model(&po.PageVisit{}).
		Where("created_at >= ?", time.Now().Add(-24*time.Hour)).
		Count(&count).Error

	return count, err
}
