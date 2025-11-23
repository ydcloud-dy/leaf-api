package server

import (
	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/leaf-api/internal/server/middleware"
	"github.com/ydcloud-dy/leaf-api/internal/service"
)

// registerRoutes 注册路由
func registerRoutes(
	r *gin.Engine,
	authService *service.AuthService,
	articleService *service.ArticleService,
	userService *service.UserService,
	categoryService *service.CategoryService,
	tagService *service.TagService,
	commentService *service.CommentService,
	chapterService *service.ChapterService,
	statsService *service.StatsService,
	settingsService *service.SettingsService,
	fileService *service.FileService,
	blogService *service.BlogService,
) {
	// 管理后台认证路由（不需要 JWT 验证）
	auth := r.Group("/auth")
	{
		auth.POST("/login", authService.Login)
		auth.POST("/logout", authService.Logout)
		auth.GET("/profile", middleware.JWTAuth(), authService.GetProfile)
	}

	// 博客前台认证路由（不需要 JWT 验证）
	blogAuth := r.Group("/blog/auth")
	{
		blogAuth.POST("/register", blogService.Register)
		blogAuth.POST("/login", blogService.Login)
		blogAuth.GET("/me", middleware.JWTAuth(), blogService.GetUserInfo)
	}

	// 博客公开路由（不需要认证）
	blog := r.Group("/blog")
	{
		// 文章相关
		blog.GET("/articles", articleService.List)           // 文章列表
		blog.GET("/articles/search", articleService.Search)  // 搜索文章
		blog.GET("/articles/archive", articleService.Archive) // 归档文章

		// 分类和标签
		blog.GET("/categories", categoryService.List) // 分类列表
		blog.GET("/tags", tagService.List)            // 标签列表

		// 章节
		blog.GET("/chapters/:tag", chapterService.GetChaptersByTag) // 获取标签下的章节及文章

		// 统计
		blog.GET("/stats", statsService.GetStats) // 站点统计
		blog.GET("/stats/hot-articles", statsService.GetHotArticles) // 热门文章
	}

	// 博客可选认证路由（支持登录和未登录状态）
	blogOptionalAuth := r.Group("/blog")
	blogOptionalAuth.Use(middleware.OptionalJWTAuth())
	{
		// 文章详情（登录用户可查看点赞收藏状态）
		blogOptionalAuth.GET("/articles/:id", blogService.GetArticleDetail)
		// 文章评论（登录用户可查看点赞状态）
		blogOptionalAuth.GET("/articles/:id/comments", blogService.GetArticleComments)
		// 留言板（登录用户可查看点赞状态）
		blogOptionalAuth.GET("/guestbook", blogService.GetGuestbookMessages)
	}

	// 博客需要认证的路由
	blogAuthed := r.Group("/blog")
	blogAuthed.Use(middleware.JWTAuth())
	{
		// 点赞
		blogAuthed.POST("/articles/:id/like", blogService.LikeArticle)
		blogAuthed.DELETE("/articles/:id/like", blogService.UnlikeArticle)

		// 收藏
		blogAuthed.POST("/articles/:id/favorite", blogService.FavoriteArticle)
		blogAuthed.DELETE("/articles/:id/favorite", blogService.UnfavoriteArticle)

		// 用户点赞和收藏列表
		blogAuthed.GET("/user/likes", blogService.GetUserLikes)
		blogAuthed.GET("/user/favorites", blogService.GetUserFavorites)
		blogAuthed.GET("/user/stats", blogService.GetUserStats)

		// 评论
		blogAuthed.POST("/comments", blogService.CreateComment)
		blogAuthed.POST("/comments/:id/like", blogService.LikeComment)
		blogAuthed.DELETE("/comments/:id/like", blogService.UnlikeComment)
		blogAuthed.DELETE("/comments/:id", blogService.DeleteComment)

		// 留言板
		blogAuthed.POST("/guestbook", blogService.CreateGuestbookMessage)
	}

	// 管理后台 API 路由（需要 JWT 验证）
	api := r.Group("/")
	api.Use(middleware.JWTAuth())
	{
		// 用户管理
		users := api.Group("/users")
		{
			users.GET("", userService.List)
			users.GET("/:id", userService.GetByID)
			users.POST("", userService.Create)
			users.PUT("/:id", userService.Update)
			users.DELETE("/:id", userService.Delete)
		}

		// 文章管理
		articles := api.Group("/articles")
		{
			articles.GET("", articleService.List)
			articles.GET("/:id", articleService.GetByID)
			articles.POST("", articleService.Create)
			articles.POST("/import", articleService.ImportMarkdown)
			articles.PUT("/:id", articleService.Update)
			articles.PATCH("/:id/status", articleService.UpdateStatus)
			articles.DELETE("/:id", articleService.Delete)
		}

		// 评论管理
		comments := api.Group("/comments")
		{
			comments.GET("", commentService.List)
			comments.DELETE("/:id", commentService.Delete)
			comments.PATCH("/:id/status", commentService.UpdateStatus)
		}

		// 标签管理
		tags := api.Group("/tags")
		{
			tags.GET("", tagService.List)
			tags.POST("", tagService.Create)
			tags.DELETE("/:id", tagService.Delete)
		}

		// 分类管理
		categories := api.Group("/categories")
		{
			categories.GET("", categoryService.List)
			categories.POST("", categoryService.Create)
			categories.DELETE("/:id", categoryService.Delete)
		}

		// 章节管理
		chapters := api.Group("/chapters")
		{
			chapters.GET("", chapterService.GetChapters)
			chapters.GET("/:id", chapterService.GetChapter)
			chapters.POST("", chapterService.CreateChapter)
			chapters.PUT("/:id", chapterService.UpdateChapter)
			chapters.DELETE("/:id", chapterService.DeleteChapter)
		}

		// 统计
		stats := api.Group("/stats")
		{
			stats.GET("", statsService.GetStats)
			stats.GET("/hot-articles", statsService.GetHotArticles)
		}

		// 设置
		settings := api.Group("/settings")
		{
			settings.GET("", settingsService.Get)
			settings.PUT("", settingsService.Update)
		}

		// 文件上传
		files := api.Group("/files")
		{
			files.POST("/upload", fileService.Upload)
			files.GET("", fileService.List)
			files.DELETE("/:id", fileService.Delete)
		}
	}
}
