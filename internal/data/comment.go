package data

import (
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"gorm.io/gorm"
)

// CommentRepo 评论仓储接口
type CommentRepo interface {
	// Create 创建评论
	Create(comment *po.Comment) error
	// Update 更新评论
	Update(comment *po.Comment) error
	// Delete 删除评论
	Delete(id uint) error
	// FindByID 根据 ID 查询评论
	FindByID(id uint) (*po.Comment, error)
	// List 查询评论列表
	List(page, limit int, articleID uint, status string) ([]*po.Comment, int64, error)
	// UpdateStatus 更新评论状态
	UpdateStatus(id uint, status int) error
	// CountByArticle 统计文章评论数
	CountByArticle(articleID uint) (int64, error)
	// CountByUser 统计用户评论数
	CountByUser(userID uint) (int64, error)
}

// commentRepo 评论仓储实现
type commentRepo struct {
	db *gorm.DB
}

// NewCommentRepo 创建评论仓储
func NewCommentRepo(db *gorm.DB) CommentRepo {
	return &commentRepo{db: db}
}

// Create 创建评论
func (r *commentRepo) Create(comment *po.Comment) error {
	return r.db.Create(comment).Error
}

// Update 更新评论
func (r *commentRepo) Update(comment *po.Comment) error {
	return r.db.Save(comment).Error
}

// Delete 删除评论
func (r *commentRepo) Delete(id uint) error {
	return r.db.Delete(&po.Comment{}, id).Error
}

// FindByID 根据 ID 查询评论
func (r *commentRepo) FindByID(id uint) (*po.Comment, error) {
	var comment po.Comment
	err := r.db.Preload("User").Preload("ReplyToUser").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// List 查询评论列表
func (r *commentRepo) List(page, limit int, articleID uint, status string) ([]*po.Comment, int64, error) {
	var comments []*po.Comment
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&po.Comment{}).Preload("User").Preload("ReplyToUser").Preload("Article")

	// 文章过滤
	if articleID > 0 {
		query = query.Where("article_id = ?", articleID)
	}

	// 状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// UpdateStatus 更新评论状态
func (r *commentRepo) UpdateStatus(id uint, status int) error {
	return r.db.Model(&po.Comment{}).Where("id = ?", id).Update("status", status).Error
}

// CountByArticle 统计文章评论数
func (r *commentRepo) CountByArticle(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&po.Comment{}).Where("article_id = ? AND status = ?", articleID, 1).Count(&count).Error
	return count, err
}

// CountByUser 统计用户评论数
func (r *commentRepo) CountByUser(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&po.Comment{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
