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
func (s *AuthService) Logout(c *gin.Context) {
	response.Success(c, nil)
}

// GetProfile 获取当前用户信息
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
