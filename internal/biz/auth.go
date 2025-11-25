package biz

import (
	"errors"

	"github.com/ydcloud-dy/leaf-api/internal/data"
	"github.com/ydcloud-dy/leaf-api/internal/model/dto"
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"github.com/ydcloud-dy/leaf-api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// AuthUseCase 认证业务用例接口
type AuthUseCase interface {
	// Login 管理员登录
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	// GetProfile 获取管理员信息
	GetProfile(adminID uint) (*dto.AdminInfo, error)
	// UpdateProfile 更新管理员信息
	UpdateProfile(adminID uint, req *dto.UpdateProfileRequest) (*po.User, error)
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
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			Bio:       user.Bio,
			Skills:    user.Skills,
			Contacts:  user.Contacts,
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
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		Skills:    user.Skills,
		Contacts:  user.Contacts,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}, nil
}

// UpdateProfile 更新管理员信息
func (uc *authUseCase) UpdateProfile(adminID uint, req *dto.UpdateProfileRequest) (*po.User, error) {
	user, err := uc.data.UserRepo.FindByID(adminID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.Email != "" && req.Email != user.Email {
		// 检查邮箱是否已被使用
		existingUser, err := uc.data.UserRepo.FindByEmail(req.Email)
		if err == nil && existingUser.ID != adminID {
			return nil, errors.New("邮箱已被其他用户使用")
		}
		user.Email = req.Email
	}
	// 更新技术栈和联系方式（允许空字符串清空）
	user.Skills = req.Skills
	user.Contacts = req.Contacts

	// 更新博主标识（仅管理员可以设置）
	if req.IsBlogger != nil {
		// 如果要设置为博主，先取消其他用户的博主标识
		if *req.IsBlogger {
			uc.data.GetDB().Model(&po.User{}).Where("is_blogger = ?", true).Update("is_blogger", false)
		}
		user.IsBlogger = *req.IsBlogger
	}

	if err := uc.data.UserRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
