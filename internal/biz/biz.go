package biz

import (
	"github.com/ydcloud-dy/leaf-api/internal/data"
)

// Biz 业务逻辑层结构
type Biz struct {
	AuthUseCase     AuthUseCase
	ArticleUseCase  ArticleUseCase
	UserUseCase     UserUseCase
	CategoryUseCase CategoryUseCase
	TagUseCase      TagUseCase
	CommentUseCase  CommentUseCase
	BlogUseCase     BlogUseCase
}

// NewBiz 创建业务逻辑层实例
func NewBiz(d *data.Data) *Biz {
	return &Biz{
		AuthUseCase:     NewAuthUseCase(d),
		ArticleUseCase:  NewArticleUseCase(d),
		UserUseCase:     NewUserUseCase(d),
		CategoryUseCase: NewCategoryUseCase(d),
		TagUseCase:      NewTagUseCase(d),
		CommentUseCase:  NewCommentUseCase(d),
		BlogUseCase:     NewBlogUseCase(d),
	}
}
