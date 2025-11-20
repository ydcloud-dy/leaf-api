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
func (s *FileService) List(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	files, total, err := s.data.FileRepo.List(page, limit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithPage(c, files, total, page, limit)
}

// Delete 删除文件
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
