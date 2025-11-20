package service

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/leaf-api/internal/biz"
	"github.com/ydcloud-dy/leaf-api/internal/model/dto"
	"github.com/ydcloud-dy/leaf-api/pkg/response"
)

// ArticleService 文章服务
type ArticleService struct {
	articleUseCase biz.ArticleUseCase
}

// NewArticleService 创建文章服务
func NewArticleService(articleUseCase biz.ArticleUseCase) *ArticleService {
	return &ArticleService{
		articleUseCase: articleUseCase,
	}
}

// Create 创建文章
func (s *ArticleService) Create(c *gin.Context) {
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 获取作者 ID
	adminID, _ := c.Get("admin_id")

	resp, err := s.articleUseCase.Create(&req, adminID.(uint))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// Update 更新文章
func (s *ArticleService) Update(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := s.articleUseCase.Update(idReq.ID, &req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// Delete 删除文章
func (s *ArticleService) Delete(c *gin.Context) {
	var req dto.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := s.articleUseCase.Delete(req.ID); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetByID 根据 ID 查询文章
func (s *ArticleService) GetByID(c *gin.Context) {
	var req dto.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := s.articleUseCase.GetByID(req.ID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// List 查询文章列表
func (s *ArticleService) List(c *gin.Context) {
	var req dto.ArticleListRequest

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	req.Page = page
	req.Limit = limit

	// 解析过滤参数
	req.Category = c.Query("category")
	req.Tag = c.Query("tag")
	req.Status = c.Query("status")
	req.Keyword = c.Query("keyword")

	resp, err := s.articleUseCase.List(&req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithPage(c, resp.Data, resp.Total, resp.Page, resp.Limit)
}

// UpdateStatus 更新文章状态
func (s *ArticleService) UpdateStatus(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var req dto.UpdateArticleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := s.articleUseCase.UpdateStatus(idReq.ID, req.Status); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// Search 搜索文章
func (s *ArticleService) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := s.articleUseCase.Search(keyword, page, limit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithPage(c, resp.Data, resp.Total, resp.Page, resp.Limit)
}

// Archive 获取归档文章
func (s *ArticleService) Archive(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	resp, err := s.articleUseCase.Archive(page, limit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithPage(c, resp.Data, resp.Total, resp.Page, resp.Limit)
}
