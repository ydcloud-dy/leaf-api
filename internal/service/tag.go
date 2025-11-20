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
func (s *TagService) List(c *gin.Context) {
	tags, err := s.tagUseCase.List()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, tags)
}

// Create 创建标签
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
