package po

import (
	"time"

	"gorm.io/gorm"
)

// Admin 管理员模型
type Admin struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Email     string         `gorm:"size:100" json:"email"`
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Avatar    string         `gorm:"size:500" json:"avatar"`
	Bio       string         `gorm:"size:500" json:"bio"`
	Skills    string         `gorm:"type:text" json:"skills"`   // 技术栈，逗号分隔
	Contacts  string         `gorm:"type:text" json:"contacts"` // 联系方式，JSON格式
	Role      string         `gorm:"size:20;default:admin" json:"role"`
	Status    int            `gorm:"default:1" json:"status"` // 1: active, 0: inactive
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User 前台用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"size:100;uniqueIndex" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Avatar    string         `gorm:"size:500" json:"avatar"`
	Bio       string         `gorm:"size:500" json:"bio"`
	Skills    string         `gorm:"type:text" json:"skills"`     // JSON数组格式的技术栈
	Contacts  string         `gorm:"type:text" json:"contacts"`   // JSON对象格式的联系方式
	Role      string         `gorm:"size:20;default:'user'" json:"role"` // user, admin, super_admin
	IsBlogger bool           `gorm:"default:false" json:"is_blogger"`    // 是否为博主（用于关于页面展示）
	Status    int            `gorm:"default:1" json:"status"`            // 1: active, 0: banned
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Article 文章模型
type Article struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	Title           string         `gorm:"size:200;not null" json:"title"`
	ContentMarkdown string         `gorm:"type:longtext" json:"content_markdown"`
	ContentHTML     string         `gorm:"type:longtext" json:"content_html"`
	Summary         string         `gorm:"size:500" json:"summary"`
	Cover           string         `gorm:"size:500" json:"cover"`
	AuthorID        uint           `gorm:"index" json:"author_id"`
	CategoryID      uint           `gorm:"index" json:"category_id"`
	ChapterID       *uint          `gorm:"index" json:"chapter_id"` // 所属章节ID,可为空
	Status          int            `gorm:"default:0" json:"status"` // 0: draft, 1: published, 2: offline
	ViewCount       int            `gorm:"default:0" json:"view_count"`
	LikeCount       int            `gorm:"default:0" json:"like_count"`
	FavoriteCount   int            `gorm:"default:0" json:"favorite_count"`
	CommentCount    int            `gorm:"default:0" json:"comment_count"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	Author   User      `gorm:"foreignKey:AuthorID;references:ID;constraint:OnDelete:SET NULL" json:"author,omitempty"`
	Category Category  `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
	Chapter  *Chapter  `gorm:"foreignKey:ChapterID;references:ID" json:"chapter,omitempty"`
	Tags     []Tag     `gorm:"many2many:article_tags" json:"tags,omitempty"`
}

// Category 分类模型
type Category struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Description string         `gorm:"size:200" json:"description"`
	Sort        int            `gorm:"default:0" json:"sort"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Tag 标签模型
type Tag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Color     string         `gorm:"size:20" json:"color"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Comment 评论模型
type Comment struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	ArticleID     *uint          `gorm:"index" json:"article_id"` // 可为空，NULL表示留言板消息
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	ParentID      *uint          `gorm:"index" json:"parent_id"`
	ReplyToUserID *uint          `gorm:"index" json:"reply_to_user_id"` // 被回复的用户ID
	Content       string         `gorm:"type:text;not null" json:"content"`
	LikeCount     int            `gorm:"default:0" json:"like_count"`
	Status        int            `gorm:"default:0" json:"status"` // 0: pending, 1: approved, 2: rejected
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ReplyToUser *User     `gorm:"foreignKey:ReplyToUserID" json:"reply_to_user,omitempty"`
	Article     *Article  `gorm:"foreignKey:ArticleID;constraint:OnDelete:SET NULL;" json:"article,omitempty"`
	Replies     []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// Like 点赞记录
type Like struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	ArticleID uint      `gorm:"index;not null" json:"article_id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Article Article `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
}

// Favorite 收藏记录
type Favorite struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	ArticleID uint      `gorm:"index;not null" json:"article_id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Article Article `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
}

// CommentLike 评论点赞记录
type CommentLike struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CommentID uint      `gorm:"index;not null" json:"comment_id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comment Comment `gorm:"foreignKey:CommentID" json:"comment,omitempty"`
}

// View 浏览记录
type View struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	ArticleID uint      `gorm:"index;not null" json:"article_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	IP        string    `gorm:"size:50" json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

// PageVisit 页面访问时长记录
type PageVisit struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    *uint     `gorm:"index" json:"user_id"` // 可为空，游客访问
	IP        string    `gorm:"size:50;index" json:"ip"`
	Path      string    `gorm:"size:500" json:"path"`         // 访问路径
	Duration  int       `gorm:"not null" json:"duration"`      // 停留时长（秒）
	UserAgent string    `gorm:"size:500" json:"user_agent"`    // 用户代理
	Referrer  string    `gorm:"size:500" json:"referrer"`      // 来源页面
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

// File 文件模型
type File struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:200;not null" json:"name"`
	URL       string         `gorm:"size:500;not null" json:"url"`
	Size      int64          `json:"size"`
	Type      string         `gorm:"size:50" json:"type"`
	MimeType  string         `gorm:"size:100" json:"mime_type"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Setting 系统设置
type Setting struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Key       string    `gorm:"size:100;uniqueIndex;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Admin{},
		&User{},
		&Article{},
		&Category{},
		&Tag{},
		&Chapter{},
		&Comment{},
		&Like{},
		&Favorite{},
		&CommentLike{},
		&View{},
		&PageVisit{},
		&File{},
		&Setting{},
	)
}
