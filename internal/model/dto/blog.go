package dto

import "time"

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	ArticleID     *uint  `json:"article_id"` // 可为空，NULL表示留言板消息
	UserID        uint   `json:"user_id"`
	ParentID      *uint  `json:"parent_id"`
	ReplyToUserID *uint  `json:"reply_to_user_id"` // 被回复的用户ID
	Content       string `json:"content" binding:"required,min=1,max=1000"`
}

// CreateGuestbookMessageRequest 创建留言板消息请求
type CreateGuestbookMessageRequest struct {
	UserID        uint   `json:"user_id"`
	ParentID      *uint  `json:"parent_id"`
	ReplyToUserID *uint  `json:"reply_to_user_id"` // 被回复的用户ID
	Content       string `json:"content" binding:"required,min=1,max=1000"`
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID           uint               `json:"id"`
	ArticleID    *uint              `json:"article_id"` // 可为空
	UserID       uint               `json:"user_id"`
	ParentID     *uint              `json:"parent_id"`
	Content      string             `json:"content"`
	LikeCount    int                `json:"like_count"`
	IsLiked      bool               `json:"is_liked"`
	Status       int                `json:"status"`
	CreatedAt    time.Time          `json:"created_at"`
	User         *UserInfo          `json:"user,omitempty"`
	ReplyToUser  *UserInfo          `json:"reply_to_user,omitempty"`
	Replies      []CommentResponse  `json:"replies,omitempty"`
}

// CommentListResponse 评论列表响应
type CommentListResponse struct {
	List  []CommentResponse `json:"list"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}

// ArticleDetailResponse 文章详情响应（包含用户状态）
type ArticleDetailResponse struct {
	ArticleResponse
	IsLiked     bool `json:"is_liked"`
	IsFavorited bool `json:"is_favorited"`
}

// LikeInfo 点赞信息
type LikeInfo struct {
	ID        uint             `json:"id"`
	ArticleID uint             `json:"article_id"`
	UserID    uint             `json:"user_id"`
	CreatedAt time.Time        `json:"created_at"`
	Article   *ArticleResponse `json:"article,omitempty"`
}

// LikeListResponse 点赞列表响应
type LikeListResponse struct {
	List  []LikeInfo `json:"list"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
}

// FavoriteInfo 收藏信息
type FavoriteInfo struct {
	ID        uint             `json:"id"`
	ArticleID uint             `json:"article_id"`
	UserID    uint             `json:"user_id"`
	CreatedAt time.Time        `json:"created_at"`
	Article   *ArticleResponse `json:"article,omitempty"`
}

// FavoriteListResponse 收藏列表响应
type FavoriteListResponse struct {
	List  []FavoriteInfo `json:"list"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// UserStatsResponse 用户统计响应
type UserStatsResponse struct {
	ArticlesCount  int64 `json:"articles_count"`
	LikesCount     int64 `json:"likes_count"`
	FavoritesCount int64 `json:"favorites_count"`
	CommentsCount  int64 `json:"comments_count"`
}

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Email    string `json:"email"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
