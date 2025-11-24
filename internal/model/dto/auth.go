package dto

import "time"

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar" binding:"max=500"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string     `json:"token"`
	Admin *AdminInfo `json:"admin,omitempty"`
	User  *UserInfo  `json:"user,omitempty"`
}

// AdminInfo 管理员信息
type AdminInfo struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Skills    string    `json:"skills"`
	Contacts  string    `json:"contacts"`
	Role      string    `json:"role"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Skills    string    `json:"skills"`
	Contacts  string    `json:"contacts"`
	Role      string    `json:"role"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
