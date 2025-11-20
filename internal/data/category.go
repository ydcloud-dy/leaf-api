package data

import (
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"gorm.io/gorm"
)

// CategoryRepo 分类仓储接口
type CategoryRepo interface {
	// Create 创建分类
	Create(category *po.Category) error
	// Update 更新分类
	Update(category *po.Category) error
	// Delete 删除分类
	Delete(id uint) error
	// FindByID 根据 ID 查询分类
	FindByID(id uint) (*po.Category, error)
	// FindByName 根据名称查询分类
	FindByName(name string) (*po.Category, error)
	// List 查询分类列表
	List() ([]*po.Category, error)
	// HasArticles 检查分类下是否有文章
	HasArticles(id uint) (bool, error)
}

// categoryRepo 分类仓储实现
type categoryRepo struct {
	db *gorm.DB
}

// NewCategoryRepo 创建分类仓储
func NewCategoryRepo(db *gorm.DB) CategoryRepo {
	return &categoryRepo{db: db}
}

// Create 创建分类
func (r *categoryRepo) Create(category *po.Category) error {
	return r.db.Create(category).Error
}

// Update 更新分类
func (r *categoryRepo) Update(category *po.Category) error {
	return r.db.Save(category).Error
}

// Delete 删除分类
func (r *categoryRepo) Delete(id uint) error {
	return r.db.Delete(&po.Category{}, id).Error
}

// FindByID 根据 ID 查询分类
func (r *categoryRepo) FindByID(id uint) (*po.Category, error) {
	var category po.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// FindByName 根据名称查询分类
func (r *categoryRepo) FindByName(name string) (*po.Category, error) {
	var category po.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// List 查询分类列表
func (r *categoryRepo) List() ([]*po.Category, error) {
	var categories []*po.Category
	err := r.db.Order("sort ASC, created_at DESC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// HasArticles 检查分类下是否有文章
func (r *categoryRepo) HasArticles(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&po.Article{}).Where("category_id = ?", id).Count(&count).Error
	return count > 0, err
}
