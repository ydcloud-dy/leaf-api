package data

import (
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"gorm.io/gorm"
)

// UserRepo 用户仓储接口
type UserRepo interface {
	// Create 创建用户
	Create(user *po.User) error
	// Update 更新用户
	Update(user *po.User) error
	// Delete 删除用户
	Delete(id uint) error
	// FindByID 根据 ID 查询用户
	FindByID(id uint) (*po.User, error)
	// FindByUsername 根据用户名查询用户
	FindByUsername(username string) (*po.User, error)
	// FindByEmail 根据邮箱查询用户
	FindByEmail(email string) (*po.User, error)
	// List 查询用户列表
	List(page, limit int, keyword, status string) ([]*po.User, int64, error)
}

// userRepo 用户仓储实现
type userRepo struct {
	db *gorm.DB
}

// NewUserRepo 创建用户仓储
func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

// Create 创建用户
func (r *userRepo) Create(user *po.User) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *userRepo) Update(user *po.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *userRepo) Delete(id uint) error {
	return r.db.Delete(&po.User{}, id).Error
}

// FindByID 根据 ID 查询用户
func (r *userRepo) FindByID(id uint) (*po.User, error) {
	var user po.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查询用户
func (r *userRepo) FindByUsername(username string) (*po.User, error) {
	var user po.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查询用户
func (r *userRepo) FindByEmail(email string) (*po.User, error) {
	var user po.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List 查询用户列表
func (r *userRepo) List(page, limit int, keyword, status string) ([]*po.User, int64, error) {
	var users []*po.User
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&po.User{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
