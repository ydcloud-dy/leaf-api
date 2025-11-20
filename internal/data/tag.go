package data

import (
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"gorm.io/gorm"
)

// TagRepo 标签仓储接口
type TagRepo interface {
	// Create 创建标签
	Create(tag *po.Tag) error
	// Update 更新标签
	Update(tag *po.Tag) error
	// Delete 删除标签
	Delete(id uint) error
	// FindByID 根据 ID 查询标签
	FindByID(id uint) (*po.Tag, error)
	// FindByName 根据名称查询标签
	FindByName(name string) (*po.Tag, error)
	// List 查询标签列表
	List() ([]*po.Tag, error)
	// FindByIDs 根据 ID 列表查询标签
	FindByIDs(ids []uint) ([]*po.Tag, error)
}

// tagRepo 标签仓储实现
type tagRepo struct {
	db *gorm.DB
}

// NewTagRepo 创建标签仓储
func NewTagRepo(db *gorm.DB) TagRepo {
	return &tagRepo{db: db}
}

// Create 创建标签
func (r *tagRepo) Create(tag *po.Tag) error {
	return r.db.Create(tag).Error
}

// Update 更新标签
func (r *tagRepo) Update(tag *po.Tag) error {
	return r.db.Save(tag).Error
}

// Delete 删除标签
func (r *tagRepo) Delete(id uint) error {
	return r.db.Select("Articles").Delete(&po.Tag{ID: id}).Error
}

// FindByID 根据 ID 查询标签
func (r *tagRepo) FindByID(id uint) (*po.Tag, error) {
	var tag po.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// FindByName 根据名称查询标签
func (r *tagRepo) FindByName(name string) (*po.Tag, error) {
	var tag po.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// List 查询标签列表
func (r *tagRepo) List() ([]*po.Tag, error) {
	var tags []*po.Tag
	err := r.db.Order("created_at DESC").Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// FindByIDs 根据 ID 列表查询标签
func (r *tagRepo) FindByIDs(ids []uint) ([]*po.Tag, error) {
	var tags []*po.Tag
	err := r.db.Find(&tags, ids).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}
