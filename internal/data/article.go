package data

import (
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"gorm.io/gorm"
)

// ArticleRepo 文章仓储接口
type ArticleRepo interface {
	// Create 创建文章
	Create(article *po.Article) error
	// Update 更新文章
	Update(article *po.Article) error
	// Delete 删除文章
	Delete(id uint) error
	// FindByID 根据 ID 查询文章
	FindByID(id uint) (*po.Article, error)
	// FindByIDWithRelations 根据 ID 查询文章（包含关联数据）
	FindByIDWithRelations(id uint) (*po.Article, error)
	// List 查询文章列表
	List(page, limit int, categoryID, tagID uint, status, keyword string) ([]*po.Article, int64, error)
	// UpdateStatus 更新文章状态
	UpdateStatus(id uint, status int) error
	// IncrementViewCount 增加浏览量
	IncrementViewCount(id uint) error
	// IncrementLikeCount 增加点赞数
	IncrementLikeCount(id uint) error
	// DecrementLikeCount 减少点赞数
	DecrementLikeCount(id uint) error
	// IncrementFavoriteCount 增加收藏数
	IncrementFavoriteCount(id uint) error
	// DecrementFavoriteCount 减少收藏数
	DecrementFavoriteCount(id uint) error
	// IncrementCommentCount 增加评论数
	IncrementCommentCount(id uint) error
	// DecrementCommentCount 减少评论数
	DecrementCommentCount(id uint) error
	// AssociateTags 关联标签
	AssociateTags(articleID uint, tagIDs []uint) error
}

// articleRepo 文章仓储实现
type articleRepo struct {
	db *gorm.DB
}

// NewArticleRepo 创建文章仓储
func NewArticleRepo(db *gorm.DB) ArticleRepo {
	return &articleRepo{db: db}
}

// Create 创建文章
func (r *articleRepo) Create(article *po.Article) error {
	return r.db.Create(article).Error
}

// Update 更新文章
func (r *articleRepo) Update(article *po.Article) error {
	return r.db.Save(article).Error
}

// Delete 删除文章
func (r *articleRepo) Delete(id uint) error {
	return r.db.Select("Tags").Delete(&po.Article{ID: id}).Error
}

// FindByID 根据 ID 查询文章
func (r *articleRepo) FindByID(id uint) (*po.Article, error) {
	var article po.Article
	err := r.db.First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// FindByIDWithRelations 根据 ID 查询文章（包含关联数据）
func (r *articleRepo) FindByIDWithRelations(id uint) (*po.Article, error) {
	var article po.Article
	err := r.db.Preload("Author").Preload("Category").Preload("Tags").First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// List 查询文章列表
func (r *articleRepo) List(page, limit int, categoryID, tagID uint, status, keyword string) ([]*po.Article, int64, error) {
	var articles []*po.Article
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&po.Article{}).Preload("Author").Preload("Category").Preload("Tags")

	// 分类过滤
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	// 标签过滤
	if tagID > 0 {
		query = query.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ?", tagID)
	}

	// 状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 关键词搜索
	if keyword != "" {
		query = query.Where("title LIKE ? OR summary LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// UpdateStatus 更新文章状态
func (r *articleRepo) UpdateStatus(id uint, status int) error {
	return r.db.Model(&po.Article{}).Where("id = ?", id).Update("status", status).Error
}

// IncrementViewCount 增加浏览量
func (r *articleRepo) IncrementViewCount(id uint) error {
	return r.db.Model(&po.Article{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// IncrementLikeCount 增加点赞数
func (r *articleRepo) IncrementLikeCount(id uint) error {
	return r.db.Model(&po.Article{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
}

// DecrementLikeCount 减少点赞数
func (r *articleRepo) DecrementLikeCount(id uint) error {
	return r.db.Model(&po.Article{}).Where("id = ? AND like_count > 0", id).
		UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
}

// IncrementFavoriteCount 增加收藏数
func (r *articleRepo) IncrementFavoriteCount(id uint) error {
	return r.db.Model(&po.Article{}).Where("id = ?", id).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
}

// DecrementFavoriteCount 减少收藏数
func (r *articleRepo) DecrementFavoriteCount(id uint) error {
	return r.db.Model(&po.Article{}).Where("id = ? AND favorite_count > 0", id).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
}

// IncrementCommentCount 增加评论数
func (r *articleRepo) IncrementCommentCount(id uint) error {
	return r.db.Model(&po.Article{}).Where("id = ?", id).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error
}

// DecrementCommentCount 减少评论数
func (r *articleRepo) DecrementCommentCount(id uint) error {
	return r.db.Model(&po.Article{}).Where("id = ? AND comment_count > 0", id).
		UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error
}

// AssociateTags 关联标签
func (r *articleRepo) AssociateTags(articleID uint, tagIDs []uint) error {
	var article po.Article
	if err := r.db.First(&article, articleID).Error; err != nil {
		return err
	}

	var tags []po.Tag
	if err := r.db.Find(&tags, tagIDs).Error; err != nil {
		return err
	}

	return r.db.Model(&article).Association("Tags").Replace(tags)
}
