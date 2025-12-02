package service

import (
	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/leaf-api/internal/biz"
	"github.com/ydcloud-dy/leaf-api/internal/model/dto"
	"github.com/ydcloud-dy/leaf-api/pkg/response"
)

// TagService 标签服务
type TagService struct {
	tagUseCase biz.TagUseCase
}

// NewTagService 创建标签服务
func NewTagService(tagUseCase biz.TagUseCase) *TagService {
	return &TagService{
		tagUseCase: tagUseCase,
	}
}

// List 查询标签列表
// @Summary 获取标签列表
// @Description 获取所有文章标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /blog/tags [get]
func (s *TagService) List(c *gin.Context) {
	tags, err := s.tagUseCase.List()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, tags)
}

// Create 创建标签
// @Summary 创建标签
// @Description 创建新的文章标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{name=string,color=string} true "标签信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /tags [post]
func (s *TagService) Create(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required,max=50"`
		Color string `json:"color" binding:"max=20"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := s.tagUseCase.Create(req.Name, req.Color); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除标签
// @Summary 删除标签
// @Description 根据ID删除文章标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "标签ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /tags/{id} [delete]
func (s *TagService) Delete(c *gin.Context) {
	var req dto.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := s.tagUseCase.Delete(req.ID); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
