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
// @Summary 获取分类列表
// @Description 获取所有文章分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /blog/categories [get]
func (s *CategoryService) List(c *gin.Context) {
	categories, err := s.categoryUseCase.List()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, categories)
}

// Create 创建分类
// @Summary 创建分类
// @Description 创建新的文章分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{name=string,description=string,sort=int} true "分类信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /categories [post]
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
// @Summary 删除分类
// @Description 根据ID删除文章分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /categories/{id} [delete]
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
