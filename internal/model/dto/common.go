package dto

// PageRequest 分页请求
type PageRequest struct {
	Page  int `form:"page" json:"page" binding:"min=1"`
	Limit int `form:"limit" json:"limit" binding:"min=1,max=100"`
}

// PageResponse 分页响应
type PageResponse struct {
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Data  interface{} `json:"data"`
}

// IDRequest ID 请求
type IDRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}
