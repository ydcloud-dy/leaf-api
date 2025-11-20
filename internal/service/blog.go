package service

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/leaf-api/internal/biz"
	"github.com/ydcloud-dy/leaf-api/internal/model/dto"
	"github.com/ydcloud-dy/leaf-api/pkg/response"
)

// BlogService 博客服务
type BlogService struct {
	blogUseCase biz.BlogUseCase
}

// NewBlogService 创建博客服务
func NewBlogService(blogUseCase biz.BlogUseCase) *BlogService {
	return &BlogService{
		blogUseCase: blogUseCase,
	}
}

// Register 用户注册
func (s *BlogService) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := s.blogUseCase.Register(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// Login 用户登录
func (s *BlogService) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := s.blogUseCase.Login(&req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetUserInfo 获取用户信息
func (s *BlogService) GetUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := s.blogUseCase.GetUserInfo(userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetArticleDetail 获取文章详情（包含用户状态）
func (s *BlogService) GetArticleDetail(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	// 获取用户ID（如果已登录）
	userID := uint(0)
	if id, exists := c.Get("user_id"); exists {
		userID = id.(uint)
	}

	resp, err := s.blogUseCase.GetArticleDetail(uint(articleID), userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// LikeArticle 点赞文章
func (s *BlogService) LikeArticle(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	if err := s.blogUseCase.LikeArticle(userID, uint(articleID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// UnlikeArticle 取消点赞
func (s *BlogService) UnlikeArticle(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	if err := s.blogUseCase.UnlikeArticle(userID, uint(articleID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// FavoriteArticle 收藏文章
func (s *BlogService) FavoriteArticle(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	if err := s.blogUseCase.FavoriteArticle(userID, uint(articleID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// UnfavoriteArticle 取消收藏
func (s *BlogService) UnfavoriteArticle(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	if err := s.blogUseCase.UnfavoriteArticle(userID, uint(articleID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetUserLikes 获取用户点赞列表
func (s *BlogService) GetUserLikes(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := s.blogUseCase.GetUserLikes(userID, page, limit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetUserFavorites 获取用户收藏列表
func (s *BlogService) GetUserFavorites(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := s.blogUseCase.GetUserFavorites(userID, page, limit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// CreateComment 创建评论
func (s *BlogService) CreateComment(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	req.UserID = userID

	resp, err := s.blogUseCase.CreateComment(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetArticleComments 获取文章评论列表
func (s *BlogService) GetArticleComments(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	// 获取用户ID（如果已登录）
	userID := uint(0)
	if id, exists := c.Get("user_id"); exists {
		userID = id.(uint)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	resp, err := s.blogUseCase.GetArticleComments(uint(articleID), userID, page, limit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// LikeComment 点赞评论
func (s *BlogService) LikeComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的评论ID")
		return
	}

	if err := s.blogUseCase.LikeComment(userID, uint(commentID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// UnlikeComment 取消点赞评论
func (s *BlogService) UnlikeComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的评论ID")
		return
	}

	if err := s.blogUseCase.UnlikeComment(userID, uint(commentID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteComment 删除评论
func (s *BlogService) DeleteComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的评论ID")
		return
	}

	if err := s.blogUseCase.DeleteComment(uint(commentID), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetGuestbookMessages 获取留言板消息列表
func (s *BlogService) GetGuestbookMessages(c *gin.Context) {
	// 获取用户ID（如果已登录）
	userID := uint(0)
	if id, exists := c.Get("user_id"); exists {
		userID = id.(uint)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 使用 article_id = 0 表示留言板消息
	resp, err := s.blogUseCase.GetArticleComments(0, userID, page, limit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// CreateGuestbookMessage 创建留言板消息
func (s *BlogService) CreateGuestbookMessage(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 设置 article_id = 0 表示留言板消息
	req.ArticleID = 0
	req.UserID = userID

	resp, err := s.blogUseCase.CreateComment(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetUserStats 获取用户统计信息
func (s *BlogService) GetUserStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	resp, err := s.blogUseCase.GetUserStats(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}
