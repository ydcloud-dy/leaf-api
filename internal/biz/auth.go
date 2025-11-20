package biz

import (
	"errors"

	"github.com/ydcloud-dy/leaf-api/internal/data"
	"github.com/ydcloud-dy/leaf-api/internal/model/dto"
	"github.com/ydcloud-dy/leaf-api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// AuthUseCase 认证业务用例接口
type AuthUseCase interface {
	// Login 管理员登录
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	// GetProfile 获取管理员信息
	GetProfile(adminID uint) (*dto.AdminInfo, error)
}

// authUseCase 认证业务用例实现
type authUseCase struct {
	data *data.Data
}

// NewAuthUseCase 创建认证业务用例
func NewAuthUseCase(d *data.Data) AuthUseCase {
	return &authUseCase{data: d}
}

// Login 管理员登录
func (uc *authUseCase) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 查询用户（统一使用users表）
	user, err := uc.data.UserRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查状态
	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	// 检查是否是管理员角色
	if user.Role != "admin" && user.Role != "super_admin" {
		return nil, errors.New("无权限访问管理后台")
	}

	// 生成 Token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("生成 Token 失败")
	}

	// 返回登录结果
	return &dto.LoginResponse{
		Token: token,
		Admin: &dto.AdminInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Avatar:    user.Avatar,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// GetProfile 获取管理员信息
func (uc *authUseCase) GetProfile(adminID uint) (*dto.AdminInfo, error) {
	user, err := uc.data.UserRepo.FindByID(adminID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return &dto.AdminInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}, nil
}
