package service

import (
	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/leaf-api/internal/data"
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"github.com/ydcloud-dy/leaf-api/pkg/response"
)

// SettingsService 设置服务
type SettingsService struct {
	data *data.Data
}

// NewSettingsService 创建设置服务
func NewSettingsService(d *data.Data) *SettingsService {
	return &SettingsService{
		data: d,
	}
}

// Get 获取所有设置
func (s *SettingsService) Get(c *gin.Context) {
	settings, err := s.data.SettingRepo.List()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// 转换为 map 格式
	settingsMap := make(map[string]string)
	for _, setting := range settings {
		settingsMap[setting.Key] = setting.Value
	}

	response.Success(c, settingsMap)
}

// Update 更新设置
func (s *SettingsService) Update(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 批量更新设置
	var settings []*po.Setting
	for key, value := range req {
		// 先查询是否存在
		existing, err := s.data.SettingRepo.FindByKey(key)
		if err != nil {
			// 不存在则创建
			settings = append(settings, &po.Setting{
				Key:   key,
				Value: value,
			})
		} else {
			// 存在则更新
			existing.Value = value
			settings = append(settings, existing)
		}
	}

	if err := s.data.SettingRepo.BatchUpdate(settings); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
