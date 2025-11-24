package dto

import "time"

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar" binding:"max=500"`
	Bio      string `json:"bio" binding:"max=500"`
	Status   int    `json:"status" binding:"oneof=0 1"` // 0: banned, 1: active
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6,max=50"`
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar" binding:"max=500"`
	Bio      string `json:"bio" binding:"max=500"`
	Skills   string `json:"skills" binding:"max=500"`
	Contacts string `json:"contacts" binding:"max=1000"`
	Status   *int   `json:"status" binding:"omitempty,oneof=0 1"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	PageRequest
	Keyword string `form:"keyword"`
	Status  string `form:"status"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Skills    string    `json:"skills"`
	Contacts  string    `json:"contacts"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
