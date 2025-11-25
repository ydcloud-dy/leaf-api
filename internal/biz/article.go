package biz

import (
	"errors"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/ydcloud-dy/leaf-api/internal/data"
	"github.com/ydcloud-dy/leaf-api/internal/model/dto"
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
)

// ArticleUseCase 文章业务用例接口
type ArticleUseCase interface {
	// Create 创建文章
	Create(req *dto.CreateArticleRequest, authorID uint) (*dto.ArticleResponse, error)
	// Update 更新文章
	Update(id uint, req *dto.UpdateArticleRequest) (*dto.ArticleResponse, error)
	// Delete 删除文章
	Delete(id uint) error
	// GetByID 根据 ID 查询文章
	GetByID(id uint) (*dto.ArticleResponse, error)
	// List 查询文章列表
	List(req *dto.ArticleListRequest) (*dto.PageResponse, error)
	// UpdateStatus 更新文章状态
	UpdateStatus(id uint, status int) error
	// Search 搜索文章
	Search(keyword string, page, limit int) (*dto.PageResponse, error)
	// Archive 获取归档文章（按月份分组）
	Archive(page, limit int) (*dto.PageResponse, error)
	// GetDefaultCategoryID 获取默认分类ID
	GetDefaultCategoryID() (uint, error)
}

// articleUseCase 文章业务用例实现
type articleUseCase struct {
	data *data.Data
}

// NewArticleUseCase 创建文章业务用例
func NewArticleUseCase(d *data.Data) ArticleUseCase {
	return &articleUseCase{data: d}
}

// Create 创建文章
func (uc *articleUseCase) Create(req *dto.CreateArticleRequest, authorID uint) (*dto.ArticleResponse, error) {
	// 验证分类是否存在
	if _, err := uc.data.CategoryRepo.FindByID(req.CategoryID); err != nil {
		return nil, errors.New("分类不存在")
	}

	// 如果没有提供 HTML，则自动从 Markdown 转换
	contentHTML := req.ContentHTML
	if contentHTML == "" {
		contentHTML = markdownToHTML(req.ContentMarkdown)
	}

	// 创建文章
	article := &po.Article{
		Title:           req.Title,
		ContentMarkdown: req.ContentMarkdown,
		ContentHTML:     contentHTML,
		Summary:         req.Summary,
		Cover:           req.Cover,
		AuthorID:        authorID,
		CategoryID:      req.CategoryID,
		ChapterID:       req.ChapterID,
		Status:          req.Status,
	}

	if err := uc.data.ArticleRepo.Create(article); err != nil {
		return nil, errors.New("创建文章失败: " + err.Error())
	}

	// 关联标签
	if len(req.TagIDs) > 0 {
		if err := uc.data.ArticleRepo.AssociateTags(article.ID, req.TagIDs); err != nil {
			return nil, errors.New("关联标签失败: " + err.Error())
		}
	}

	// 重新查询文章（包含关联数据）
	return uc.GetByID(article.ID)
}

// Update 更新文章
func (uc *articleUseCase) Update(id uint, req *dto.UpdateArticleRequest) (*dto.ArticleResponse, error) {
	// 查询文章
	article, err := uc.data.ArticleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	// 更新字段
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.ContentMarkdown != "" {
		article.ContentMarkdown = req.ContentMarkdown
		// 如果提供了 Markdown，自动转换为 HTML（除非明确提供了 HTML）
		if req.ContentHTML != "" {
			article.ContentHTML = req.ContentHTML
		} else {
			article.ContentHTML = markdownToHTML(req.ContentMarkdown)
		}
	}
	if req.Summary != "" {
		article.Summary = req.Summary
	}
	if req.Cover != "" {
		article.Cover = req.Cover
	}
	if req.CategoryID > 0 {
		// 验证分类是否存在
		if _, err := uc.data.CategoryRepo.FindByID(req.CategoryID); err != nil {
			return nil, errors.New("分类不存在")
		}
		article.CategoryID = req.CategoryID
	}
	// 设置章节ID（可为空）
	article.ChapterID = req.ChapterID
	if req.Status >= 0 {
		article.Status = req.Status
	}

	if err := uc.data.ArticleRepo.Update(article); err != nil {
		return nil, errors.New("更新文章失败")
	}

	// 更新标签关联
	if len(req.TagIDs) > 0 {
		if err := uc.data.ArticleRepo.AssociateTags(article.ID, req.TagIDs); err != nil {
			return nil, errors.New("更新标签失败")
		}
	}

	// 重新查询文章
	return uc.GetByID(id)
}

// Delete 删除文章
func (uc *articleUseCase) Delete(id uint) error {
	// 检查文章是否存在
	if _, err := uc.data.ArticleRepo.FindByID(id); err != nil {
		return errors.New("文章不存在")
	}

	if err := uc.data.ArticleRepo.Delete(id); err != nil {
		return errors.New("删除文章失败")
	}

	return nil
}

// GetByID 根据 ID 查询文章
func (uc *articleUseCase) GetByID(id uint) (*dto.ArticleResponse, error) {
	article, err := uc.data.ArticleRepo.FindByIDWithRelations(id)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	return uc.convertToArticleResponse(article), nil
}

// List 查询文章列表
func (uc *articleUseCase) List(req *dto.ArticleListRequest) (*dto.PageResponse, error) {
	// 解析查询参数
	var categoryID, tagID uint
	if req.Category != "" {
		category, err := uc.data.CategoryRepo.FindByName(req.Category)
		if err == nil {
			categoryID = category.ID
		}
	}
	if req.Tag != "" {
		tag, err := uc.data.TagRepo.FindByName(req.Tag)
		if err == nil {
			tagID = tag.ID
		}
	}

	// 查询文章列表
	articles, total, err := uc.data.ArticleRepo.List(
		req.Page, req.Limit,
		categoryID, tagID,
		req.Status, req.Keyword,
	)
	if err != nil {
		return nil, errors.New("查询文章列表失败")
	}

	// 转换为 DTO
	items := make([]dto.ArticleListItem, 0, len(articles))
	for _, article := range articles {
		items = append(items, uc.convertToArticleListItem(article))
	}

	return &dto.PageResponse{
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
		Data:  items,
	}, nil
}

// UpdateStatus 更新文章状态
func (uc *articleUseCase) UpdateStatus(id uint, status int) error {
	// 检查文章是否存在
	if _, err := uc.data.ArticleRepo.FindByID(id); err != nil {
		return errors.New("文章不存在")
	}

	if err := uc.data.ArticleRepo.UpdateStatus(id, status); err != nil {
		return errors.New("更新状态失败")
	}

	return nil
}

// convertToArticleResponse 转换为文章响应
func (uc *articleUseCase) convertToArticleResponse(article *po.Article) *dto.ArticleResponse {
	resp := &dto.ArticleResponse{
		ID:              article.ID,
		Title:           article.Title,
		ContentMarkdown: article.ContentMarkdown,
		ContentHTML:     article.ContentHTML,
		Summary:         article.Summary,
		Cover:           article.Cover,
		AuthorID:        article.AuthorID,
		CategoryID:      article.CategoryID,
		ChapterID:       article.ChapterID,
		Status:          article.Status,
		ViewCount:       article.ViewCount,
		LikeCount:       article.LikeCount,
		FavoriteCount:   article.FavoriteCount,
		CommentCount:    article.CommentCount,
		CreatedAt:       article.CreatedAt,
		UpdatedAt:       article.UpdatedAt,
	}

	// 作者信息
	if article.Author.ID > 0 {
		resp.Author = &dto.AuthorInfo{
			ID:       article.Author.ID,
			Username: article.Author.Username,
			Avatar:   article.Author.Avatar,
		}
	}

	// 分类信息
	if article.Category.ID > 0 {
		resp.Category = &dto.CategoryInfo{
			ID:          article.Category.ID,
			Name:        article.Category.Name,
			Description: article.Category.Description,
		}
	}

	// 标签信息
	if len(article.Tags) > 0 {
		tags := make([]dto.TagInfo, 0, len(article.Tags))
		for _, tag := range article.Tags {
			tags = append(tags, dto.TagInfo{
				ID:    tag.ID,
				Name:  tag.Name,
				Color: tag.Color,
			})
		}
		resp.Tags = tags
	}

	return resp
}

// convertToArticleListItem 转换为文章列表项
func (uc *articleUseCase) convertToArticleListItem(article *po.Article) dto.ArticleListItem {
	item := dto.ArticleListItem{
		ID:            article.ID,
		Title:         article.Title,
		Summary:       article.Summary,
		Cover:         article.Cover,
		Status:        article.Status,
		ViewCount:     article.ViewCount,
		LikeCount:     article.LikeCount,
		FavoriteCount: article.FavoriteCount,
		CommentCount:  article.CommentCount,
		CreatedAt:     article.CreatedAt,
	}

	// 作者信息
	if article.Author.ID > 0 {
		item.Author = &dto.AuthorInfo{
			ID:       article.Author.ID,
			Username: article.Author.Username,
			Avatar:   article.Author.Avatar,
		}
	}

	// 分类信息
	if article.Category.ID > 0 {
		item.Category = &dto.CategoryInfo{
			ID:          article.Category.ID,
			Name:        article.Category.Name,
			Description: article.Category.Description,
		}
	}

	// 标签信息
	if len(article.Tags) > 0 {
		tags := make([]dto.TagInfo, 0, len(article.Tags))
		for _, tag := range article.Tags {
			tags = append(tags, dto.TagInfo{
				ID:    tag.ID,
				Name:  tag.Name,
				Color: tag.Color,
			})
		}
		item.Tags = tags
	}

	return item
}

// markdownToHTML 将 Markdown 转换为 HTML
func markdownToHTML(md string) string {
	// 创建 Markdown 解析器
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	// 创建 HTML 渲染器
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// 渲染为 HTML
	return string(markdown.Render(doc, renderer))
}

// Search 搜索文章
func (uc *articleUseCase) Search(keyword string, page, limit int) (*dto.PageResponse, error) {
	// 使用文章列表请求结构，设置搜索关键词
	req := &dto.ArticleListRequest{
		PageRequest: dto.PageRequest{
			Page:  page,
			Limit: limit,
		},
		Keyword: keyword,
		Status:  "1", // 只搜索已发布的文章
	}
	return uc.List(req)
}

// Archive 获取归档文章（返回所有已发布的文章，前端按月份分组）
func (uc *articleUseCase) Archive(page, limit int) (*dto.PageResponse, error) {
	req := &dto.ArticleListRequest{
		PageRequest: dto.PageRequest{
			Page:  page,
			Limit: limit,
		},
		Status: "1", // 只返回已发布的文章
	}
	return uc.List(req)
}

// GetDefaultCategoryID 获取默认分类ID
func (uc *articleUseCase) GetDefaultCategoryID() (uint, error) {
	categories, err := uc.data.CategoryRepo.List()
	if err != nil {
		return 0, errors.New("查询分类列表失败")
	}
	if len(categories) == 0 {
		return 0, errors.New("系统中没有可用的分类")
	}
	return categories[0].ID, nil
}

