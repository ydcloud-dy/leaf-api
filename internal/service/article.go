package service

import (
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/leaf-api/internal/biz"
	"github.com/ydcloud-dy/leaf-api/internal/model/dto"
	"github.com/ydcloud-dy/leaf-api/pkg/response"
)

// ArticleService 文章服务
type ArticleService struct {
	articleUseCase biz.ArticleUseCase
}

// NewArticleService 创建文章服务
func NewArticleService(articleUseCase biz.ArticleUseCase) *ArticleService {
	return &ArticleService{
		articleUseCase: articleUseCase,
	}
}

// Create 创建文章
func (s *ArticleService) Create(c *gin.Context) {
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 获取作者 ID
	adminID, _ := c.Get("admin_id")

	resp, err := s.articleUseCase.Create(&req, adminID.(uint))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// Update 更新文章
func (s *ArticleService) Update(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := s.articleUseCase.Update(idReq.ID, &req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// Delete 删除文章
func (s *ArticleService) Delete(c *gin.Context) {
	var req dto.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := s.articleUseCase.Delete(req.ID); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetByID 根据 ID 查询文章
func (s *ArticleService) GetByID(c *gin.Context) {
	var req dto.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := s.articleUseCase.GetByID(req.ID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// List 查询文章列表
func (s *ArticleService) List(c *gin.Context) {
	var req dto.ArticleListRequest

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	req.Page = page
	req.Limit = limit

	// 解析过滤参数
	req.Category = c.Query("category")
	req.Tag = c.Query("tag")
	req.Status = c.Query("status")
	req.Keyword = c.Query("keyword")

	resp, err := s.articleUseCase.List(&req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithPage(c, resp.Data, resp.Total, resp.Page, resp.Limit)
}

// UpdateStatus 更新文章状态
func (s *ArticleService) UpdateStatus(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var req dto.UpdateArticleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := s.articleUseCase.UpdateStatus(idReq.ID, req.Status); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// Search 搜索文章
func (s *ArticleService) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := s.articleUseCase.Search(keyword, page, limit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithPage(c, resp.Data, resp.Total, resp.Page, resp.Limit)
}

// Archive 获取归档文章
func (s *ArticleService) Archive(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	resp, err := s.articleUseCase.Archive(page, limit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithPage(c, resp.Data, resp.Total, resp.Page, resp.Limit)
}

// ImportMarkdown 批量导入 Markdown 文件
func (s *ArticleService) ImportMarkdown(c *gin.Context) {
	// 获取作者 ID
	adminID, exists := c.Get("admin_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	// 解析 multipart form
	form, err := c.MultipartForm()
	if err != nil {
		response.BadRequest(c, "解析表单失败: "+err.Error())
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.BadRequest(c, "没有上传文件")
		return
	}

	successCount := 0
	failedFiles := []string{}

	// 遍历所有文件
	for _, file := range files {
		// 检查文件扩展名
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".md" && ext != ".markdown" {
			failedFiles = append(failedFiles, file.Filename+": 不支持的文件格式")
			continue
		}

		// 打开文件
		f, err := file.Open()
		if err != nil {
			failedFiles = append(failedFiles, file.Filename+": 打开文件失败")
			continue
		}

		// 读取文件内容
		content, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			failedFiles = append(failedFiles, file.Filename+": 读取文件失败")
			continue
		}

		// 提取文件名作为标题（去掉扩展名）
		title := strings.TrimSuffix(file.Filename, ext)

		// 创建文章
		req := &dto.CreateArticleRequest{
			Title:           title,
			ContentMarkdown: string(content),
			Summary:         generateSummary(string(content), 200),
			Status:          0, // 默认为草稿
			CategoryID:      1, // 默认分类，可以根据需要修改
			TagIDs:          []uint{},
		}

		_, err = s.articleUseCase.Create(req, adminID.(uint))
		if err != nil {
			failedFiles = append(failedFiles, file.Filename+": 创建文章失败")
			continue
		}

		successCount++
	}

	result := map[string]interface{}{
		"total":   len(files),
		"success": successCount,
		"failed":  len(failedFiles),
	}

	if len(failedFiles) > 0 {
		result["failed_files"] = failedFiles
	}

	response.Success(c, result)
}

// generateSummary 从内容中生成摘要
func generateSummary(content string, maxLen int) string {
	// 移除 Markdown 标记
	content = strings.ReplaceAll(content, "#", "")
	content = strings.ReplaceAll(content, "*", "")
	content = strings.ReplaceAll(content, "_", "")
	content = strings.ReplaceAll(content, "`", "")
	content = strings.ReplaceAll(content, "\n", " ")
	content = strings.TrimSpace(content)

	// 截取指定长度
	runes := []rune(content)
	if len(runes) > maxLen {
		return string(runes[:maxLen]) + "..."
	}
	return content
}

