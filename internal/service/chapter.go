package service

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/leaf-api/internal/data"
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"github.com/ydcloud-dy/leaf-api/pkg/response"
	"gorm.io/gorm"
)

type ChapterService struct {
	data *data.Data
}

func NewChapterService(d *data.Data) *ChapterService {
	return &ChapterService{data: d}
}

// GetChapters 获取章节列表(按标签)
func (s *ChapterService) GetChapters(c *gin.Context) {
	tagID := c.Query("tag_id")
	
	var chapters []po.Chapter
	query := s.data.GetDB().Model(&po.Chapter{}).Preload("Tag")
	
	if tagID != "" {
		query = query.Where("tag_id = ?", tagID)
	}
	
	query.Order("sort ASC, id ASC").Find(&chapters)
	response.Success(c, chapters)
}

// GetChapter 获取章节详情
func (s *ChapterService) GetChapter(c *gin.Context) {
	id := c.Param("id")
	
	var chapter po.Chapter
	if err := s.data.GetDB().Preload("Tag").First(&chapter, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "章节不存在")
			return
		}
		response.Error(c, 500, "获取章节失败")
		return
	}
	
	response.Success(c, chapter)
}

// CreateChapter 创建章节
func (s *ChapterService) CreateChapter(c *gin.Context) {
	var req struct {
		TagID    uint   `json:"tag_id" binding:"required"`
		ParentID *uint  `json:"parent_id"` // 父章节ID，可选
		Name     string `json:"name" binding:"required"`
		Sort     int    `json:"sort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	chapter := po.Chapter{
		TagID:    req.TagID,
		ParentID: req.ParentID,
		Name:     req.Name,
		Sort:     req.Sort,
	}

	if err := s.data.GetDB().Create(&chapter).Error; err != nil {
		response.Error(c, 500, "创建章节失败")
		return
	}

	response.Success(c, chapter)
}

// UpdateChapter 更新章节
func (s *ChapterService) UpdateChapter(c *gin.Context) {
	id := c.Param("id")

	var chapter po.Chapter
	if err := s.data.GetDB().First(&chapter, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "章节不存在")
			return
		}
		response.Error(c, 500, "获取章节失败")
		return
	}

	var req struct {
		TagID    *uint   `json:"tag_id"`
		ParentID *uint   `json:"parent_id"` // 支持更新父章节ID
		Name     *string `json:"name"`
		Sort     *int    `json:"sort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.TagID != nil {
		updates["tag_id"] = *req.TagID
	}
	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}

	if err := s.data.GetDB().Model(&chapter).Updates(updates).Error; err != nil {
		response.Error(c, 500, "更新章节失败")
		return
	}

	response.Success(c, chapter)
}

// DeleteChapter 删除章节
func (s *ChapterService) DeleteChapter(c *gin.Context) {
	id := c.Param("id")
	
	var chapter po.Chapter
	if err := s.data.GetDB().First(&chapter, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "章节不存在")
			return
		}
		response.Error(c, 500, "获取章节失败")
		return
	}
	
	// 检查是否有文章关联
	var count int64
	s.data.GetDB().Model(&po.Article{}).Where("chapter_id = ?", id).Count(&count)
	if count > 0 {
		response.Error(c, 400, "该章节下还有文章,无法删除")
		return
	}
	
	if err := s.data.GetDB().Delete(&chapter).Error; err != nil {
		response.Error(c, 500, "删除章节失败")
		return
	}
	
	response.Success(c, nil)
}

// GetChaptersByTag 获取标签下的章节及文章(用于前端笔记页面)
// 支持多级目录结构
func (s *ChapterService) GetChaptersByTag(c *gin.Context) {
	tagName := c.Param("tag")

	// 先查找标签
	var tag po.Tag
	if err := s.data.GetDB().Where("name = ?", tagName).First(&tag).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, 404, "标签不存在")
			return
		}
		response.Error(c, 500, "查询失败")
		return
	}

	// 查询该标签下的所有章节
	var chapters []po.Chapter
	s.data.GetDB().Where("tag_id = ?", tag.ID).Order("sort ASC, id ASC").Find(&chapters)

	// 如果没有章节，直接返回
	if len(chapters) == 0 {
		response.Success(c, []interface{}{})
		return
	}

	// 收集所有章节ID
	chapterIDs := make([]uint, len(chapters))
	for i, chapter := range chapters {
		chapterIDs[i] = chapter.ID
	}

	// 一次性查询所有章节的文章 (优化N+1查询问题)
	var allArticles []po.Article
	s.data.GetDB().Where("chapter_id IN ? AND status = 1", chapterIDs).
		Select("id, title, chapter_id, view_count, created_at").
		Find(&allArticles)

	// 按章节ID分组文章
	articlesByChapter := make(map[uint][]po.Article)
	for _, article := range allArticles {
		if article.ChapterID != nil {
			articlesByChapter[*article.ChapterID] = append(articlesByChapter[*article.ChapterID], article)
		}
	}

	// 为每个章节的文章排序
	for chapterID := range articlesByChapter {
		sortArticlesByTitleNumber(articlesByChapter[chapterID])
	}

	// 组装结果 - 构建树形结构
	type ChapterWithArticles struct {
		po.Chapter
		Articles     []po.Article           `json:"articles"`
		SubChapters  []ChapterWithArticles  `json:"sub_chapters"`
	}

	// 先构建一个章节ID到章节的映射
	chapterMap := make(map[uint]*ChapterWithArticles)
	for _, chapter := range chapters {
		articles := articlesByChapter[chapter.ID]
		if articles == nil {
			articles = []po.Article{}
		}

		chapterMap[chapter.ID] = &ChapterWithArticles{
			Chapter:     chapter,
			Articles:    articles,
			SubChapters: []ChapterWithArticles{},
		}
	}

	// 构建树形结构 - 先找出所有子章节并添加到父章节的 SubChapters 中
	for _, chapter := range chapters {
		if chapter.ParentID != nil {
			// 如果是子章节，添加到父章节的 SubChapters 中
			if parentChapter, ok := chapterMap[*chapter.ParentID]; ok {
				chapterItem := chapterMap[chapter.ID]
				parentChapter.SubChapters = append(parentChapter.SubChapters, *chapterItem)
			}
		}
	}

	// 然后构建结果列表,只包含一级章节
	var result []ChapterWithArticles
	for _, chapter := range chapters {
		if chapter.ParentID == nil {
			result = append(result, *chapterMap[chapter.ID])
		}
	}

	response.Success(c, result)
}

// sortArticlesByTitleNumber 按标题中的序号对文章进行排序
func sortArticlesByTitleNumber(articles []po.Article) {
	// 中文数字映射表
	chineseNumbers := map[rune]int{
		'一': 1, '二': 2, '三': 3, '四': 4, '五': 5,
		'六': 6, '七': 7, '八': 8, '九': 9, '十': 10,
		'壹': 1, '贰': 2, '叁': 3, '肆': 4, '伍': 5,
		'陆': 6, '柒': 7, '捌': 8, '玖': 9, '拾': 10,
	}

	// 提取标题中的序号
	extractNumber := func(title string) int {
		if title == "" {
			return 999999 // 没有标题的排在最后
		}

		// 去除空格
		title = strings.TrimSpace(title)

		// 尝试匹配阿拉伯数字: "1."、"1、"、"1 "、"1-" 等
		arabicRegex := regexp.MustCompile(`^(\d+)[.\s、,，-]`)
		if matches := arabicRegex.FindStringSubmatch(title); len(matches) > 1 {
			if num, err := strconv.Atoi(matches[1]); err == nil {
				return num
			}
		}

		// 尝试匹配中文数字: "一、"、"一."、"一 " 等
		if len(title) > 0 {
			firstChar := []rune(title)[0]
			if num, ok := chineseNumbers[firstChar]; ok {
				return num
			}
		}

		// 如果都没有匹配到,返回一个很大的数,排在后面
		return 999999
	}

	// 使用 sort.Slice 进行排序
	sort.Slice(articles, func(i, j int) bool {
		numI := extractNumber(articles[i].Title)
		numJ := extractNumber(articles[j].Title)

		// 如果序号相同,按创建时间排序
		if numI == numJ {
			return articles[i].CreatedAt.Before(articles[j].CreatedAt)
		}

		return numI < numJ
	})
}
