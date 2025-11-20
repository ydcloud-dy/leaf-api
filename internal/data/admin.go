package data

import (
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"gorm.io/gorm"
)

// AdminRepo 管理员仓储接口
type AdminRepo interface {
	// Create 创建管理员
	Create(admin *po.Admin) error
	// Update 更新管理员
	Update(admin *po.Admin) error
	// Delete 删除管理员
	Delete(id uint) error
	// FindByID 根据 ID 查询管理员
	FindByID(id uint) (*po.Admin, error)
	// FindByUsername 根据用户名查询管理员
	FindByUsername(username string) (*po.Admin, error)
	// FindByEmail 根据邮箱查询管理员
	FindByEmail(email string) (*po.Admin, error)
	// List 查询管理员列表
	List(page, limit int, keyword, status string) ([]*po.Admin, int64, error)
}

// adminRepo 管理员仓储实现
type adminRepo struct {
	db *gorm.DB
}

// NewAdminRepo 创建管理员仓储
func NewAdminRepo(db *gorm.DB) AdminRepo {
	return &adminRepo{db: db}
}

// Create 创建管理员
func (r *adminRepo) Create(admin *po.Admin) error {
	return r.db.Create(admin).Error
}

// Update 更新管理员
func (r *adminRepo) Update(admin *po.Admin) error {
	return r.db.Save(admin).Error
}

// Delete 删除管理员
func (r *adminRepo) Delete(id uint) error {
	return r.db.Delete(&po.Admin{}, id).Error
}

// FindByID 根据 ID 查询管理员
func (r *adminRepo) FindByID(id uint) (*po.Admin, error) {
	var admin po.Admin
	err := r.db.First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// FindByUsername 根据用户名查询管理员
func (r *adminRepo) FindByUsername(username string) (*po.Admin, error) {
	var admin po.Admin
	err := r.db.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// FindByEmail 根据邮箱查询管理员
func (r *adminRepo) FindByEmail(email string) (*po.Admin, error) {
	var admin po.Admin
	err := r.db.Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// List 查询管理员列表
func (r *adminRepo) List(page, limit int, keyword, status string) ([]*po.Admin, int64, error) {
	var admins []*po.Admin
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&po.Admin{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&admins).Error; err != nil {
		return nil, 0, err
	}

	return admins, total, nil
}
