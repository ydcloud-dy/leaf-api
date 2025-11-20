package service

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/leaf-api/internal/biz"
	"github.com/ydcloud-dy/leaf-api/internal/model/dto"
	"github.com/ydcloud-dy/leaf-api/pkg/response"
)

// UserService 用户服务
type UserService struct {
	userUseCase biz.UserUseCase
}

// NewUserService 创建用户服务
func NewUserService(userUseCase biz.UserUseCase) *UserService {
	return &UserService{
		userUseCase: userUseCase,
	}
}

// Create 创建用户
func (s *UserService) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := s.userUseCase.Create(&req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// Update 更新用户
func (s *UserService) Update(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := s.userUseCase.Update(idReq.ID, &req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// Delete 删除用户
func (s *UserService) Delete(c *gin.Context) {
	var req dto.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := s.userUseCase.Delete(req.ID); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetByID 根据 ID 查询用户
func (s *UserService) GetByID(c *gin.Context) {
	var req dto.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := s.userUseCase.GetByID(req.ID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// List 查询用户列表
func (s *UserService) List(c *gin.Context) {
	var req dto.UserListRequest

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	req.Page = page
	req.Limit = limit

	// 解析过滤参数
	req.Keyword = c.Query("keyword")
	req.Status = c.Query("status")

	resp, err := s.userUseCase.List(&req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithPage(c, resp.Data, resp.Total, resp.Page, resp.Limit)
}
