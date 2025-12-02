package service

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/leaf-api/internal/data"
	"github.com/ydcloud-dy/leaf-api/internal/model/dto"
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"github.com/ydcloud-dy/leaf-api/pkg/oss"
	"github.com/ydcloud-dy/leaf-api/pkg/response"
)

// FileService 文件服务
type FileService struct {
	data *data.Data
}

// NewFileService 创建文件服务
func NewFileService(d *data.Data) *FileService {
	return &FileService{
		data: d,
	}
}

// Upload 上传文件
// @Summary 上传文件
// @Description 上传文件到OSS，支持图片、视频等多种文件类型
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "文件"
// @Param folder formData string false "文件夹名称" default(uploads)
// @Success 200 {object} response.Response "上传成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /files/upload [post]
func (s *FileService) Upload(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "文件上传失败")
		return
	}

	// 获取文件夹参数
	folder := c.DefaultPostForm("folder", "uploads")

	// 上传到 OSS
	url, err := oss.UploadFile(file, folder)
	if err != nil {
		response.ServerError(c, "上传文件失败: "+err.Error())
		return
	}

	// 保存文件记录
	fileRecord := &po.File{
		Name:     file.Filename,
		URL:      url,
		Size:     file.Size,
		Type:     folder,
		MimeType: file.Header.Get("Content-Type"),
	}

	if err := s.data.FileRepo.Create(fileRecord); err != nil {
		response.ServerError(c, "保存文件记录失败")
		return
	}

	response.Success(c, gin.H{
		"url":  url,
		"name": file.Filename,
		"size": file.Size,
		"id":   fileRecord.ID,
	})
}

// List 查询文件列表
// @Summary 获取文件列表
// @Description 分页获取已上传的文件列表
// @Tags 文件管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /files [get]
func (s *FileService) List(c *gin.Context) {
	// 解析分页参数 (兼容 page_size 和 limit 两种参数名)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := c.Query("page_size")
	if pageSize == "" {
		pageSize = c.DefaultQuery("limit", "20")
	}
	limit, _ := strconv.Atoi(pageSize)

	files, total, err := s.data.FileRepo.List(page, limit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithPage(c, files, total, page, limit)
}

// Delete 删除文件
// @Summary 删除文件
// @Description 删除OSS上的文件及数据库记录
// @Tags 文件管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文件ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /files/{id} [delete]
func (s *FileService) Delete(c *gin.Context) {
	var req dto.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 查询文件记录
	file, err := s.data.FileRepo.FindByID(req.ID)
	if err != nil {
		response.ServerError(c, "文件不存在")
		return
	}

	// 从 OSS 删除文件
	objectKey := oss.GetObjectKeyFromURL(file.URL)
	if err := oss.DeleteFile(objectKey); err != nil {
		response.ServerError(c, "删除文件失败: "+err.Error())
		return
	}

	// 删除数据库记录
	if err := s.data.FileRepo.Delete(req.ID); err != nil {
		response.ServerError(c, "删除文件记录失败")
		return
	}

	response.Success(c, nil)
}
