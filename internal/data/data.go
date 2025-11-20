package data

import (
	"gorm.io/gorm"
)

// Data 数据层结构，包含所有 Repository
type Data struct {
	db              *gorm.DB
	AdminRepo       AdminRepo
	UserRepo        UserRepo
	ArticleRepo     ArticleRepo
	CategoryRepo    CategoryRepo
	TagRepo         TagRepo
	CommentRepo     CommentRepo
	LikeRepo        LikeRepo
	FavoriteRepo    FavoriteRepo
	CommentLikeRepo CommentLikeRepo
	ViewRepo        ViewRepo
	FileRepo        FileRepo
	SettingRepo     SettingRepo
}

// NewData 创建数据层实例
func NewData(db *gorm.DB) (*Data, error) {
	return &Data{
		db:              db,
		AdminRepo:       NewAdminRepo(db),
		UserRepo:        NewUserRepo(db),
		ArticleRepo:     NewArticleRepo(db),
		CategoryRepo:    NewCategoryRepo(db),
		TagRepo:         NewTagRepo(db),
		CommentRepo:     NewCommentRepo(db),
		LikeRepo:        NewLikeRepo(db),
		FavoriteRepo:    NewFavoriteRepo(db),
		CommentLikeRepo: NewCommentLikeRepo(db),
		ViewRepo:        NewViewRepo(db),
		FileRepo:        NewFileRepo(db),
		SettingRepo:     NewSettingRepo(db),
	}, nil
}

// GetDB 获取数据库实例（用于事务）
func (d *Data) GetDB() *gorm.DB {
	return d.db
}
