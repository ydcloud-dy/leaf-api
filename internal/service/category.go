package service

import (
	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/leaf-api/internal/biz"
	"github.com/ydcloud-dy/leaf-api/internal/model/dto"
	"github.com/ydcloud-dy/leaf-api/pkg/response"
)

// CategoryService 分类服务
type CategoryService struct {
	categoryUseCase biz.CategoryUseCase
}

// NewCategoryService 创建分类服务
func NewCategoryService(categoryUseCase biz.CategoryUseCase) *CategoryService {
	return &CategoryService{
		categoryUseCase: categoryUseCase,
	}
}

// List 查询分类列表
func (s *CategoryService) List(c *gin.Context) {
	categories, err := s.categoryUseCase.List()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, categories)
}

// Create 创建分类
func (s *CategoryService) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required,max=50"`
		Description string `json:"description" binding:"max=200"`
		Sort        int    `json:"sort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := s.categoryUseCase.Create(req.Name, req.Description, req.Sort); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除分类
func (s *CategoryService) Delete(c *gin.Context) {
	var req dto.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := s.categoryUseCase.Delete(req.ID); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
