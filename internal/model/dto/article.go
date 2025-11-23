package dto

import "time"

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Title           string `json:"title" binding:"required,max=200"`
	ContentMarkdown string `json:"content_markdown" binding:"required"`
	ContentHTML     string `json:"content_html"` // 可选，如果不传则自动从 Markdown 转换
	Summary         string `json:"summary" binding:"max=500"`
	Cover           string `json:"cover" binding:"max=500"`
	CategoryID      uint   `json:"category_id" binding:"required"`
	ChapterID       *uint  `json:"chapter_id"` // 章节ID，可为空
	TagIDs          []uint `json:"tag_ids"`
	Status          int    `json:"status" binding:"oneof=0 1 2"` // 0: draft, 1: published, 2: offline
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	Title           string `json:"title" binding:"omitempty,max=200"`
	ContentMarkdown string `json:"content_markdown"`
	ContentHTML     string `json:"content_html"` // 可选
	Summary         string `json:"summary" binding:"max=500"`
	Cover           string `json:"cover" binding:"max=500"`
	CategoryID      uint   `json:"category_id"`
	ChapterID       *uint  `json:"chapter_id"` // 章节ID，可为空
	TagIDs          []uint `json:"tag_ids"`
	Status          int    `json:"status" binding:"omitempty,oneof=0 1 2"`
}

// UpdateArticleStatusRequest 更新文章状态请求
type UpdateArticleStatusRequest struct {
	Status int `json:"status" binding:"required,oneof=0 1 2"`
}

// ArticleListRequest 文章列表请求
type ArticleListRequest struct {
	PageRequest
	Category string `form:"category"`
	Tag      string `form:"tag"`
	Status   string `form:"status"`
	Keyword  string `form:"keyword"`
}

// ArticleResponse 文章响应
type ArticleResponse struct {
	ID              uint             `json:"id"`
	Title           string           `json:"title"`
	ContentMarkdown string           `json:"content_markdown"`
	ContentHTML     string           `json:"content_html"`
	Summary         string           `json:"summary"`
	Cover           string           `json:"cover"`
	AuthorID        uint             `json:"author_id"`
	CategoryID      uint             `json:"category_id"`
	ChapterID       *uint            `json:"chapter_id"`
	Status          int              `json:"status"`
	ViewCount       int              `json:"view_count"`
	LikeCount       int              `json:"like_count"`
	FavoriteCount   int              `json:"favorite_count"`
	CommentCount    int              `json:"comment_count"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	Author          *AuthorInfo      `json:"author,omitempty"`
	Category        *CategoryInfo    `json:"category,omitempty"`
	Tags            []TagInfo        `json:"tags,omitempty"`
}

// ArticleListItem 文章列表项
type ArticleListItem struct {
	ID            uint          `json:"id"`
	Title         string        `json:"title"`
	Summary       string        `json:"summary"`
	Cover         string        `json:"cover"`
	Status        int           `json:"status"`
	ViewCount     int           `json:"view_count"`
	LikeCount     int           `json:"like_count"`
	FavoriteCount int           `json:"favorite_count"`
	CommentCount  int           `json:"comment_count"`
	CreatedAt     time.Time     `json:"created_at"`
	Author        *AuthorInfo   `json:"author,omitempty"`
	Category      *CategoryInfo `json:"category,omitempty"`
	Tags          []TagInfo     `json:"tags,omitempty"`
}

// CategoryInfo 分类信息
type CategoryInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TagInfo 标签信息
type TagInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// AuthorInfo 作者信息
type AuthorInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
