package service

import (
	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/leaf-api/internal/biz"
	"github.com/ydcloud-dy/leaf-api/internal/model/dto"
	"github.com/ydcloud-dy/leaf-api/pkg/response"
)

// AuthService 认证服务
type AuthService struct {
	authUseCase biz.AuthUseCase
}

// NewAuthService 创建认证服务
func NewAuthService(authUseCase biz.AuthUseCase) *AuthService {
	return &AuthService{
		authUseCase: authUseCase,
	}
}

// Login 登录
// @Summary 管理员登录
// @Description 使用用户名和密码登录，返回 JWT Token
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录信息"
// @Success 200 {object} response.Response "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "认证失败"
// @Router /auth/login [post]
func (s *AuthService) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := s.authUseCase.Login(&req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// Logout 登出
// @Summary 管理员登出
// @Description 退出登录
// @Tags 认证管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "登出成功"
// @Router /auth/logout [post]
func (s *AuthService) Logout(c *gin.Context) {
	response.Success(c, nil)
}

// GetProfile 获取当前用户信息
// @Summary 获取当前管理员信息
// @Description 获取当前登录管理员的详细信息
// @Tags 认证管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /auth/profile [get]
func (s *AuthService) GetProfile(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	resp, err := s.authUseCase.GetProfile(adminID.(uint))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// UpdateProfile 更新当前用户信息
// @Summary 更新当前管理员信息
// @Description 更新当前登录管理员的个人信息
// @Tags 认证管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "更新信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /auth/profile [put]
func (s *AuthService) UpdateProfile(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := s.authUseCase.UpdateProfile(adminID.(uint), &req)
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
