package biz

import (
	"errors"

	"github.com/ydcloud-dy/leaf-api/internal/data"
	"github.com/ydcloud-dy/leaf-api/internal/model/dto"
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"golang.org/x/crypto/bcrypt"
)

// UserUseCase 用户业务用例接口（实际管理管理员账号）
type UserUseCase interface {
	// Create 创建管理员
	Create(req *dto.CreateUserRequest) (*dto.UserResponse, error)
	// Update 更新管理员
	Update(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	// Delete 删除管理员
	Delete(id uint) error
	// GetByID 根据 ID 查询管理员
	GetByID(id uint) (*dto.UserResponse, error)
	// List 查询管理员列表
	List(req *dto.UserListRequest) (*dto.PageResponse, error)
}

// userUseCase 用户业务用例实现
type userUseCase struct {
	data *data.Data
}

// NewUserUseCase 创建用户业务用例
func NewUserUseCase(d *data.Data) UserUseCase {
	return &userUseCase{data: d}
}

// Create 创建管理员
func (uc *userUseCase) Create(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// 检查用户名是否已存在
	if _, err := uc.data.UserRepo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if _, err := uc.data.UserRepo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户（管理员）
	user := &po.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Avatar:   req.Avatar,
		Role:     "admin", // 管理端创建的用户默认是admin角色
		Status:   req.Status,
	}

	if err := uc.data.UserRepo.Create(user); err != nil {
		return nil, errors.New("创建用户失败")
	}

	return uc.convertToUserResponse(user), nil
}

// Update 更新管理员
func (uc *userUseCase) Update(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// 查询用户
	user, err := uc.data.UserRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 更新字段
	if req.Username != "" && req.Username != user.Username {
		// 检查用户名是否已存在
		if _, err := uc.data.UserRepo.FindByUsername(req.Username); err == nil {
			return nil, errors.New("用户名已存在")
		}
		user.Username = req.Username
	}
	if req.Email != "" && req.Email != user.Email {
		// 检查邮箱是否已存在
		if _, err := uc.data.UserRepo.FindByEmail(req.Email); err == nil {
			return nil, errors.New("邮箱已存在")
		}
		user.Email = req.Email
	}
	if req.Password != "" {
		// 加密密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("密码加密失败")
		}
		user.Password = string(hashedPassword)
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	// 更新技术栈和联系方式（允许空字符串清空）
	user.Skills = req.Skills
	user.Contacts = req.Contacts
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := uc.data.UserRepo.Update(user); err != nil {
		return nil, errors.New("更新用户失败")
	}

	return uc.convertToUserResponse(user), nil
}

// Delete 删除用户
func (uc *userUseCase) Delete(id uint) error {
	// 检查用户是否存在
	user, err := uc.data.UserRepo.FindByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 只有当删除的是管理员时，才检查是否是最后一个管理员
	if user.Role == "admin" || user.Role == "super_admin" {
		// 统计管理员数量
		var adminCount int64
		uc.data.GetDB().Model(&po.User{}).Where("role IN ?", []string{"admin", "super_admin"}).Count(&adminCount)
		if adminCount <= 1 {
			return errors.New("不能删除最后一个管理员")
		}
	}

	if err := uc.data.UserRepo.Delete(id); err != nil {
		return errors.New("删除用户失败")
	}

	return nil
}

// GetByID 根据 ID 查询用户
func (uc *userUseCase) GetByID(id uint) (*dto.UserResponse, error) {
	user, err := uc.data.UserRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return uc.convertToUserResponse(user), nil
}

// List 查询用户列表
func (uc *userUseCase) List(req *dto.UserListRequest) (*dto.PageResponse, error) {
	users, total, err := uc.data.UserRepo.List(req.Page, req.Limit, req.Keyword, req.Status)
	if err != nil {
		return nil, errors.New("查询管理员列表失败")
	}

	// 转换为 DTO
	items := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		items = append(items, *uc.convertToUserResponse(user))
	}

	return &dto.PageResponse{
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
		Data:  items,
	}, nil
}

// convertToUserResponse 转换为用户响应
func (uc *userUseCase) convertToUserResponse(user *po.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		Skills:    user.Skills,
		Contacts:  user.Contacts,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// CategoryUseCase 分类业务用例接口
type CategoryUseCase interface {
	// Create 创建分类
	Create(name, description string, sort int) error
	// Delete 删除分类
	Delete(id uint) error
	// List 查询分类列表
	List() ([]po.Category, error)
}

// categoryUseCase 分类业务用例实现
type categoryUseCase struct {
	data *data.Data
}

// NewCategoryUseCase 创建分类业务用例
func NewCategoryUseCase(d *data.Data) CategoryUseCase {
	return &categoryUseCase{data: d}
}

// Create 创建分类
func (uc *categoryUseCase) Create(name, description string, sort int) error {
	// 检查分类名称是否已存在
	if _, err := uc.data.CategoryRepo.FindByName(name); err == nil {
		return errors.New("分类名称已存在")
	}

	category := &po.Category{
		Name:        name,
		Description: description,
		Sort:        sort,
	}

	if err := uc.data.CategoryRepo.Create(category); err != nil {
		return errors.New("创建分类失败")
	}

	return nil
}

// Delete 删除分类
func (uc *categoryUseCase) Delete(id uint) error {
	// 检查分类是否存在
	if _, err := uc.data.CategoryRepo.FindByID(id); err != nil {
		return errors.New("分类不存在")
	}

	// 检查分类下是否有文章
	hasArticles, err := uc.data.CategoryRepo.HasArticles(id)
	if err != nil {
		return errors.New("查询失败")
	}
	if hasArticles {
		return errors.New("该分类下存在文章，无法删除")
	}

	if err := uc.data.CategoryRepo.Delete(id); err != nil {
		return errors.New("删除分类失败")
	}

	return nil
}

// List 查询分类列表
func (uc *categoryUseCase) List() ([]po.Category, error) {
	categories, err := uc.data.CategoryRepo.List()
	if err != nil {
		return nil, errors.New("查询分类列表失败")
	}

	result := make([]po.Category, 0, len(categories))
	for _, category := range categories {
		result = append(result, *category)
	}

	return result, nil
}

// TagUseCase 标签业务用例接口
type TagUseCase interface {
	// Create 创建标签
	Create(name, color string) error
	// Delete 删除标签
	Delete(id uint) error
	// List 查询标签列表
	List() ([]po.Tag, error)
}

// tagUseCase 标签业务用例实现
type tagUseCase struct {
	data *data.Data
}

// NewTagUseCase 创建标签业务用例
func NewTagUseCase(d *data.Data) TagUseCase {
	return &tagUseCase{data: d}
}

// Create 创建标签
func (uc *tagUseCase) Create(name, color string) error {
	// 检查标签名称是否已存在
	if _, err := uc.data.TagRepo.FindByName(name); err == nil {
		return errors.New("标签名称已存在")
	}

	tag := &po.Tag{
		Name:  name,
		Color: color,
	}

	if err := uc.data.TagRepo.Create(tag); err != nil {
		return errors.New("创建标签失败")
	}

	return nil
}

// Delete 删除标签
func (uc *tagUseCase) Delete(id uint) error {
	// 检查标签是否存在
	if _, err := uc.data.TagRepo.FindByID(id); err != nil {
		return errors.New("标签不存在")
	}

	if err := uc.data.TagRepo.Delete(id); err != nil {
		return errors.New("删除标签失败")
	}

	return nil
}

// List 查询标签列表
func (uc *tagUseCase) List() ([]po.Tag, error) {
	tags, err := uc.data.TagRepo.List()
	if err != nil {
		return nil, errors.New("查询标签列表失败")
	}

	result := make([]po.Tag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, *tag)
	}

	return result, nil
}

// CommentUseCase 评论业务用例接口
type CommentUseCase interface {
	// Delete 删除评论
	Delete(id uint) error
	// UpdateStatus 更新评论状态
	UpdateStatus(id uint, status int) error
	// List 查询评论列表
	List(page, limit int, articleID uint, status string) ([]*po.Comment, int64, error)
}

// commentUseCase 评论业务用例实现
type commentUseCase struct {
	data *data.Data
}

// NewCommentUseCase 创建评论业务用例
func NewCommentUseCase(d *data.Data) CommentUseCase {
	return &commentUseCase{data: d}
}

// Delete 删除评论
func (uc *commentUseCase) Delete(id uint) error {
	// 检查评论是否存在
	if _, err := uc.data.CommentRepo.FindByID(id); err != nil {
		return errors.New("评论不存在")
	}

	if err := uc.data.CommentRepo.Delete(id); err != nil {
		return errors.New("删除评论失败")
	}

	return nil
}

// UpdateStatus 更新评论状态
func (uc *commentUseCase) UpdateStatus(id uint, status int) error {
	// 检查评论是否存在
	if _, err := uc.data.CommentRepo.FindByID(id); err != nil {
		return errors.New("评论不存在")
	}

	if err := uc.data.CommentRepo.UpdateStatus(id, status); err != nil {
		return errors.New("更新状态失败")
	}

	return nil
}

// List 查询评论列表
func (uc *commentUseCase) List(page, limit int, articleID uint, status string) ([]*po.Comment, int64, error) {
	comments, total, err := uc.data.CommentRepo.List(page, limit, articleID, status)
	if err != nil {
		return nil, 0, errors.New("查询评论列表失败")
	}

	return comments, total, nil
}
