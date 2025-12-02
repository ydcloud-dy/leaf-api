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
// @Summary 获取评论列表
// @Description 分页获取评论列表，支持按文章ID和状态筛选
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Param article_id query int false "文章ID"
// @Param status query string false "评论状态"
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /comments [get]
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
// @Summary 删除评论
// @Description 根据ID删除评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /comments/{id} [delete]
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
// @Summary 更新评论状态
// @Description 更新评论的审核状态（待审核/已通过/已拒绝）
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Param request body object{status=int} true "状态信息 0:待审核 1:已通过 2:已拒绝"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /comments/{id}/status [patch]
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
