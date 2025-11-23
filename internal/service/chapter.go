package service

import (
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
		TagID uint   `json:"tag_id" binding:"required"`
		Name  string `json:"name" binding:"required"`
		Sort  int    `json:"sort"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	
	chapter := po.Chapter{
		TagID: req.TagID,
		Name:  req.Name,
		Sort:  req.Sort,
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
		TagID *uint   `json:"tag_id"`
		Name  *string `json:"name"`
		Sort  *int    `json:"sort"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	
	updates := make(map[string]interface{})
	if req.TagID != nil {
		updates["tag_id"] = *req.TagID
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
	
	// 为每个章节查询文章
	type ChapterWithArticles struct {
		po.Chapter
		Articles []po.Article `json:"articles"`
	}
	
	var result []ChapterWithArticles
	for _, chapter := range chapters {
		var articles []po.Article
		s.data.GetDB().Where("chapter_id = ? AND status = 1", chapter.ID).
			Order("created_at DESC").
			Find(&articles)
		
		result = append(result, ChapterWithArticles{
			Chapter:  chapter,
			Articles: articles,
		})
	}
	
	response.Success(c, result)
}
