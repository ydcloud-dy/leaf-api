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
// @Summary 创建文章
// @Description 创建一篇新文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateArticleRequest true "文章信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /articles [post]
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
// @Summary 更新文章
// @Description 更新文章信息
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Param request body dto.UpdateArticleRequest true "文章信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /articles/{id} [put]
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
// @Summary 删除文章
// @Description 根据ID删除文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /articles/{id} [delete]
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
// @Summary 获取文章详情
// @Description 根据ID获取文章详细信息
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /articles/{id} [get]
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
// @Summary 获取文章列表
// @Description 分页获取文章列表，支持筛选和搜索
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param category query string false "分类"
// @Param tag query string false "标签"
// @Param status query string false "状态"
// @Param keyword query string false "搜索关键词"
// @Param sort query string false "排序方式" default(latest)
// @Success 200 {object} response.Response "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /blog/articles [get]
func (s *ArticleService) List(c *gin.Context) {
	var req dto.ArticleListRequest

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	// 如果传了page_size，优先使用page_size
	if pageSize > 0 {
		limit = pageSize
	}
	req.Page = page
	req.Limit = limit

	// 解析过滤参数
	req.Category = c.Query("category")
	req.Tag = c.Query("tag")
	req.Status = c.Query("status")
	req.Keyword = c.Query("keyword")
	req.Sort = c.DefaultQuery("sort", "latest") // 默认按最新排序

	resp, err := s.articleUseCase.List(&req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithPage(c, resp.Data, resp.Total, resp.Page, resp.Limit)
}

// UpdateStatus 更新文章状态
// @Summary 更新文章状态
// @Description 更新文章的发布状态（草稿/已发布）
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Param request body dto.UpdateArticleStatusRequest true "状态信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /articles/{id}/status [patch]
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
// @Summary 搜索文章
// @Description 根据关键词搜索文章
// @Tags 博客前台
// @Accept json
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param sort query string false "排序方式" default(latest)
// @Success 200 {object} response.Response "搜索成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /blog/articles/search [get]
func (s *ArticleService) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	// 如果传了page_size，优先使用page_size
	if pageSize > 0 {
		limit = pageSize
	}
	sort := c.DefaultQuery("sort", "latest") // 默认按最新排序

	resp, err := s.articleUseCase.Search(keyword, page, limit, sort)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithPage(c, resp.Data, resp.Total, resp.Page, resp.Limit)
}

// Archive 获取归档文章
// @Summary 获取归档文章
// @Description 获取文章归档列表，按时间归档
// @Tags 博客前台
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(100)
// @Success 200 {object} response.Response "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /blog/articles/archive [get]
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
// @Summary 批量导入Markdown文件
// @Description 批量导入Markdown文件为文章
// @Tags 文章管理
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param files formData file true "Markdown文件（可多个）"
// @Success 200 {object} response.Response "导入成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /articles/import [post]
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

	// 获取默认分类ID（使用第一个可用分类）
	defaultCategoryID, err := s.articleUseCase.GetDefaultCategoryID()
	if err != nil {
		response.BadRequest(c, "获取默认分类失败: "+err.Error())
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
			CategoryID:      defaultCategoryID,
			TagIDs:          []uint{},
		}

		_, err = s.articleUseCase.Create(req, adminID.(uint))
		if err != nil {
			failedFiles = append(failedFiles, file.Filename+": 创建文章失败 - "+err.Error())
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

// BatchUpdateCover 批量更新封面
// @Summary 批量更新文章封面
// @Description 批量更新多篇文章的封面图
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.BatchUpdateCoverRequest true "封面更新信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /articles/batch-update-cover [post]
func (s *ArticleService) BatchUpdateCover(c *gin.Context) {
	var req dto.BatchUpdateCoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := s.articleUseCase.BatchUpdateCover(req.ArticleIDs, req.Cover); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"updated": len(req.ArticleIDs),
	})
}

// BatchUpdateFields 批量更新字段
// @Summary 批量更新文章字段
// @Description 批量更新多篇文章的指定字段（状态、分类、标签等）
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.BatchUpdateFieldsRequest true "字段更新信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /articles/batch-update-fields [post]
func (s *ArticleService) BatchUpdateFields(c *gin.Context) {
	var req dto.BatchUpdateFieldsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := s.articleUseCase.BatchUpdateFields(&req); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"updated": len(req.ArticleIDs),
	})
}

// BatchDelete 批量删除
// @Summary 批量删除文章
// @Description 批量删除多篇文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.BatchDeleteRequest true "删除文章ID列表"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /articles/batch-delete [post]
func (s *ArticleService) BatchDelete(c *gin.Context) {
	var req dto.BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := s.articleUseCase.BatchDelete(req.ArticleIDs); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"deleted": len(req.ArticleIDs),
	})
}

