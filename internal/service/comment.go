package service

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/leaf-api/internal/biz"
	"github.com/ydcloud-dy/leaf-api/internal/model/dto"
	"github.com/ydcloud-dy/leaf-api/pkg/response"
)

// CommentService 评论服务
type CommentService struct {
	commentUseCase biz.CommentUseCase
}

// NewCommentService 创建评论服务
func NewCommentService(commentUseCase biz.CommentUseCase) *CommentService {
	return &CommentService{
		commentUseCase: commentUseCase,
	}
}

// List 查询评论列表
func (s *CommentService) List(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 解析过滤参数
	articleIDStr := c.Query("article_id")
	var articleID uint
	if articleIDStr != "" {
		id, _ := strconv.ParseUint(articleIDStr, 10, 32)
		articleID = uint(id)
	}
	status := c.Query("status")

	comments, total, err := s.commentUseCase.List(page, limit, articleID, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithPage(c, comments, total, page, limit)
}

// Delete 删除评论
func (s *CommentService) Delete(c *gin.Context) {
	var req dto.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := s.commentUseCase.Delete(req.ID); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateStatus 更新评论状态
func (s *CommentService) UpdateStatus(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1 2"` // 0: pending, 1: approved, 2: rejected
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := s.commentUseCase.UpdateStatus(idReq.ID, req.Status); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
