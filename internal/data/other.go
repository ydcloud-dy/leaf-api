package data

import (
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"gorm.io/gorm"
)

// LikeRepo 点赞仓储接口
type LikeRepo interface {
	// Create 创建点赞
	Create(like *po.Like) error
	// Delete 删除点赞
	Delete(articleID, userID uint) error
	// Exists 检查是否已点赞
	Exists(articleID, userID uint) (bool, error)
	// List 查询点赞列表
	List(articleID uint, page, limit int) ([]*po.Like, int64, error)
	// ListByUser 根据用户ID查询点赞列表
	ListByUser(userID uint, page, limit int) ([]*po.Like, int64, error)
	// CountByArticle 统计文章点赞数
	CountByArticle(articleID uint) (int64, error)
	// CountByUser 统计用户点赞数
	CountByUser(userID uint) (int64, error)
}

// likeRepo 点赞仓储实现
type likeRepo struct {
	db *gorm.DB
}

// NewLikeRepo 创建点赞仓储
func NewLikeRepo(db *gorm.DB) LikeRepo {
	return &likeRepo{db: db}
}

// Create 创建点赞
func (r *likeRepo) Create(like *po.Like) error {
	return r.db.Create(like).Error
}

// Delete 删除点赞
func (r *likeRepo) Delete(articleID, userID uint) error {
	return r.db.Where("article_id = ? AND user_id = ?", articleID, userID).Delete(&po.Like{}).Error
}

// Exists 检查是否已点赞
func (r *likeRepo) Exists(articleID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&po.Like{}).Where("article_id = ? AND user_id = ?", articleID, userID).Count(&count).Error
	return count > 0, err
}

// List 查询点赞列表
func (r *likeRepo) List(articleID uint, page, limit int) ([]*po.Like, int64, error) {
	var likes []*po.Like
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&po.Like{}).Preload("User").Where("article_id = ?", articleID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&likes).Error; err != nil {
		return nil, 0, err
	}

	return likes, total, nil
}

// ListByUser 根据用户ID查询点赞列表
func (r *likeRepo) ListByUser(userID uint, page, limit int) ([]*po.Like, int64, error) {
	var likes []*po.Like
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&po.Like{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Article").Preload("Article.Category").Offset(offset).Limit(limit).Order("created_at DESC").Find(&likes).Error; err != nil {
		return nil, 0, err
	}

	return likes, total, nil
}

// CountByArticle 统计文章点赞数
func (r *likeRepo) CountByArticle(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&po.Like{}).Where("article_id = ?", articleID).Count(&count).Error
	return count, err
}

// CountByUser 统计用户点赞数
func (r *likeRepo) CountByUser(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&po.Like{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// FavoriteRepo 收藏仓储接口
type FavoriteRepo interface {
	// Create 创建收藏
	Create(favorite *po.Favorite) error
	// Delete 删除收藏
	Delete(articleID, userID uint) error
	// Exists 检查是否已收藏
	Exists(articleID, userID uint) (bool, error)
	// List 查询收藏列表
	List(articleID uint, page, limit int) ([]*po.Favorite, int64, error)
	// ListByUser 根据用户ID查询收藏列表
	ListByUser(userID uint, page, limit int) ([]*po.Favorite, int64, error)
	// CountByArticle 统计文章收藏数
	CountByArticle(articleID uint) (int64, error)
	// CountByUser 统计用户收藏数
	CountByUser(userID uint) (int64, error)
}

// favoriteRepo 收藏仓储实现
type favoriteRepo struct {
	db *gorm.DB
}

// NewFavoriteRepo 创建收藏仓储
func NewFavoriteRepo(db *gorm.DB) FavoriteRepo {
	return &favoriteRepo{db: db}
}

// Create 创建收藏
func (r *favoriteRepo) Create(favorite *po.Favorite) error {
	return r.db.Create(favorite).Error
}

// Delete 删除收藏
func (r *favoriteRepo) Delete(articleID, userID uint) error {
	return r.db.Where("article_id = ? AND user_id = ?", articleID, userID).Delete(&po.Favorite{}).Error
}

// Exists 检查是否已收藏
func (r *favoriteRepo) Exists(articleID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&po.Favorite{}).Where("article_id = ? AND user_id = ?", articleID, userID).Count(&count).Error
	return count > 0, err
}

// List 查询收藏列表
func (r *favoriteRepo) List(articleID uint, page, limit int) ([]*po.Favorite, int64, error) {
	var favorites []*po.Favorite
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&po.Favorite{}).Preload("User").Where("article_id = ?", articleID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&favorites).Error; err != nil {
		return nil, 0, err
	}

	return favorites, total, nil
}

// ListByUser 根据用户ID查询收藏列表
func (r *favoriteRepo) ListByUser(userID uint, page, limit int) ([]*po.Favorite, int64, error) {
	var favorites []*po.Favorite
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&po.Favorite{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Article").Preload("Article.Category").Offset(offset).Limit(limit).Order("created_at DESC").Find(&favorites).Error; err != nil {
		return nil, 0, err
	}

	return favorites, total, nil
}

// CountByArticle 统计文章收藏数
func (r *favoriteRepo) CountByArticle(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&po.Favorite{}).Where("article_id = ?", articleID).Count(&count).Error
	return count, err
}

// CountByUser 统计用户收藏数
func (r *favoriteRepo) CountByUser(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&po.Favorite{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CommentLikeRepo 评论点赞仓储接口
type CommentLikeRepo interface {
	// Create 创建评论点赞
	Create(like *po.CommentLike) error
	// Delete 删除评论点赞
	Delete(commentID, userID uint) error
	// Exists 检查是否已点赞评论
	Exists(commentID, userID uint) (bool, error)
	// CountByComment 统计评论点赞数
	CountByComment(commentID uint) (int64, error)
}

// commentLikeRepo 评论点赞仓储实现
type commentLikeRepo struct {
	db *gorm.DB
}

// NewCommentLikeRepo 创建评论点赞仓储
func NewCommentLikeRepo(db *gorm.DB) CommentLikeRepo {
	return &commentLikeRepo{db: db}
}

// Create 创建评论点赞
func (r *commentLikeRepo) Create(like *po.CommentLike) error {
	return r.db.Create(like).Error
}

// Delete 删除评论点赞
func (r *commentLikeRepo) Delete(commentID, userID uint) error {
	return r.db.Where("comment_id = ? AND user_id = ?", commentID, userID).Delete(&po.CommentLike{}).Error
}

// Exists 检查是否已点赞评论
func (r *commentLikeRepo) Exists(commentID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&po.CommentLike{}).Where("comment_id = ? AND user_id = ?", commentID, userID).Count(&count).Error
	return count > 0, err
}

// CountByComment 统计评论点赞数
func (r *commentLikeRepo) CountByComment(commentID uint) (int64, error) {
	var count int64
	err := r.db.Model(&po.CommentLike{}).Where("comment_id = ?", commentID).Count(&count).Error
	return count, err
}

// ViewRepo 浏览记录仓储接口
type ViewRepo interface {
	// Create 创建浏览记录
	Create(view *po.View) error
	// CountByArticle 统计文章浏览量
	CountByArticle(articleID uint) (int64, error)
	// CountToday 统计今日浏览量
	CountToday() (int64, error)
}

// viewRepo 浏览记录仓储实现
type viewRepo struct {
	db *gorm.DB
}

// NewViewRepo 创建浏览记录仓储
func NewViewRepo(db *gorm.DB) ViewRepo {
	return &viewRepo{db: db}
}

// Create 创建浏览记录
func (r *viewRepo) Create(view *po.View) error {
	return r.db.Create(view).Error
}

// CountByArticle 统计文章浏览量
func (r *viewRepo) CountByArticle(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&po.View{}).Where("article_id = ?", articleID).Count(&count).Error
	return count, err
}

// CountToday 统计今日浏览量
func (r *viewRepo) CountToday() (int64, error) {
	var count int64
	err := r.db.Model(&po.View{}).Where("DATE(created_at) = CURDATE()").Count(&count).Error
	return count, err
}

// FileRepo 文件仓储接口
type FileRepo interface {
	// Create 创建文件
	Create(file *po.File) error
	// Delete 删除文件
	Delete(id uint) error
	// FindByID 根据 ID 查询文件
	FindByID(id uint) (*po.File, error)
	// List 查询文件列表
	List(page, limit int) ([]*po.File, int64, error)
}

// fileRepo 文件仓储实现
type fileRepo struct {
	db *gorm.DB
}

// NewFileRepo 创建文件仓储
func NewFileRepo(db *gorm.DB) FileRepo {
	return &fileRepo{db: db}
}

// Create 创建文件
func (r *fileRepo) Create(file *po.File) error {
	return r.db.Create(file).Error
}

// Delete 删除文件
func (r *fileRepo) Delete(id uint) error {
	return r.db.Delete(&po.File{}, id).Error
}

// FindByID 根据 ID 查询文件
func (r *fileRepo) FindByID(id uint) (*po.File, error) {
	var file po.File
	err := r.db.First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// List 查询文件列表
func (r *fileRepo) List(page, limit int) ([]*po.File, int64, error) {
	var files []*po.File
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&po.File{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// SettingRepo 设置仓储接口
type SettingRepo interface {
	// Create 创建设置
	Create(setting *po.Setting) error
	// Update 更新设置
	Update(setting *po.Setting) error
	// Delete 删除设置
	Delete(key string) error
	// FindByKey 根据 Key 查询设置
	FindByKey(key string) (*po.Setting, error)
	// List 查询所有设置
	List() ([]*po.Setting, error)
	// BatchUpdate 批量更新设置
	BatchUpdate(settings []*po.Setting) error
}

// settingRepo 设置仓储实现
type settingRepo struct {
	db *gorm.DB
}

// NewSettingRepo 创建设置仓储
func NewSettingRepo(db *gorm.DB) SettingRepo {
	return &settingRepo{db: db}
}

// Create 创建设置
func (r *settingRepo) Create(setting *po.Setting) error {
	return r.db.Create(setting).Error
}

// Update 更新设置
func (r *settingRepo) Update(setting *po.Setting) error {
	return r.db.Save(setting).Error
}

// Delete 删除设置
func (r *settingRepo) Delete(key string) error {
	return r.db.Where("key = ?", key).Delete(&po.Setting{}).Error
}

// FindByKey 根据 Key 查询设置
func (r *settingRepo) FindByKey(key string) (*po.Setting, error) {
	var setting po.Setting
	err := r.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// List 查询所有设置
func (r *settingRepo) List() ([]*po.Setting, error) {
	var settings []*po.Setting
	err := r.db.Find(&settings).Error
	if err != nil {
		return nil, err
	}
	return settings, nil
}

// BatchUpdate 批量更新设置
func (r *settingRepo) BatchUpdate(settings []*po.Setting) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, setting := range settings {
			if err := tx.Save(setting).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
