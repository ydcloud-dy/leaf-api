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
// @Summary 获取系统设置
// @Description 获取所有系统配置项
// @Tags 系统设置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /settings [get]
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
// @Summary 更新系统设置
// @Description 批量更新系统配置项
// @Tags 系统设置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "配置项键值对"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /settings [put]
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
