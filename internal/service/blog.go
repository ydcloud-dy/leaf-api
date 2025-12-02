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
// @Summary 用户注册
// @Description 博客前台用户注册
// @Tags 博客前台
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response "注册成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /blog/auth/register [post]
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
// @Summary 用户登录
// @Description 博客前台用户登录
// @Tags 博客前台
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录信息"
// @Success 200 {object} response.Response "登录成功"
// @Failure 401 {object} response.Response "认证失败"
// @Router /blog/auth/login [post]
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
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 博客前台
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "用户不存在"
// @Router /blog/auth/me [get]
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
// @Summary 获取文章详情
// @Description 获取文章详细内容，包含用户点赞收藏状态（需登录）
// @Tags 博客前台
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "文章不存在"
// @Router /blog/articles/{id} [get]
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
// @Summary 点赞文章
// @Description 用户点赞文章
// @Tags 博客前台
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response "点赞成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /blog/articles/{id}/like [post]
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
// @Summary 取消点赞文章
// @Description 用户取消点赞文章
// @Tags 博客前台
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response "取消成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /blog/articles/{id}/like [delete]
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
// @Summary 收藏文章
// @Description 用户收藏文章
// @Tags 博客前台
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response "收藏成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /blog/articles/{id}/favorite [post]
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
// @Summary 取消收藏文章
// @Description 用户取消收藏文章
// @Tags 博客前台
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response "取消成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /blog/articles/{id}/favorite [delete]
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
// @Summary 获取用户点赞列表
// @Description 获取当前用户点赞的文章列表
// @Tags 博客前台
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /blog/user/likes [get]
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
// @Summary 获取用户收藏列表
// @Description 获取当前用户收藏的文章列表
// @Tags 博客前台
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /blog/user/favorites [get]
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
// @Summary 创建评论
// @Description 用户发表评论
// @Tags 博客前台
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCommentRequest true "评论信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /blog/comments [post]
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
// @Summary 获取文章评论
// @Description 获取指定文章的评论列表（包含用户点赞状态）
// @Tags 博客前台
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /blog/articles/{id}/comments [get]
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
// @Summary 点赞评论
// @Description 用户点赞评论
// @Tags 博客前台
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response "点赞成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /blog/comments/{id}/like [post]
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
// @Summary 取消点赞评论
// @Description 用户取消点赞评论
// @Tags 博客前台
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response "取消成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /blog/comments/{id}/like [delete]
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
// @Summary 删除评论
// @Description 用户删除自己的评论
// @Tags 博客前台
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /blog/comments/{id} [delete]
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
// @Summary 获取留言板消息
// @Description 获取留言板消息列表（包含用户点赞状态）
// @Tags 博客前台
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} response.Response "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /blog/guestbook [get]
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
// @Summary 创建留言板消息
// @Description 用户在留言板发表留言
// @Tags 博客前台
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateGuestbookMessageRequest true "留言信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /blog/guestbook [post]
func (s *BlogService) CreateGuestbookMessage(c *gin.Context) {
	userID := c.GetUint("user_id")

	var guestbookReq dto.CreateGuestbookMessageRequest
	if err := c.ShouldBindJSON(&guestbookReq); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 转换为评论请求，ArticleID为nil表示留言板消息
	req := dto.CreateCommentRequest{
		ArticleID:     nil, // nil表示留言板消息
		UserID:        userID,
		ParentID:      guestbookReq.ParentID,
		ReplyToUserID: guestbookReq.ReplyToUserID,
		Content:       guestbookReq.Content,
	}

	resp, err := s.blogUseCase.CreateComment(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetUserStats 获取用户统计信息
// @Summary 获取用户统计信息
// @Description 获取当前用户的统计信息（点赞数、收藏数等）
// @Tags 博客前台
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /blog/user/stats [get]
func (s *BlogService) GetUserStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	resp, err := s.blogUseCase.GetUserStats(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// UpdateProfile 更新用户资料
// @Summary 更新用户资料
// @Description 用户更新个人资料
// @Tags 博客前台
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "用户资料"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /blog/auth/profile [put]
func (s *BlogService) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := s.blogUseCase.UpdateProfile(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 返回更新后的用户数据（不包含密码）
	response.Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"nickname":   user.Nickname,
		"avatar":     user.Avatar,
		"bio":        user.Bio,
		"skills":     user.Skills,
		"contacts":   user.Contacts,
		"role":       user.Role,
		"is_blogger": user.IsBlogger,
	})
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 用户修改登录密码
// @Tags 博客前台
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "密码信息"
// @Success 200 {object} response.Response "修改成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /blog/auth/password [put]
func (s *BlogService) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := s.blogUseCase.ChangePassword(userID, &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetBloggerInfo 获取博主信息
// @Summary 获取博主信息
// @Description 获取博主的个人信息（用于关于页面）
// @Tags 博客前台
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /blog/blogger [get]
func (s *BlogService) GetBloggerInfo(c *gin.Context) {
	resp, err := s.blogUseCase.GetBloggerInfo()
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, resp)
}
