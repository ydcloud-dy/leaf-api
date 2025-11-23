package biz

import (
	"errors"
	"time"

	"github.com/ydcloud-dy/leaf-api/internal/data"
	"github.com/ydcloud-dy/leaf-api/internal/model/dto"
	"github.com/ydcloud-dy/leaf-api/internal/model/po"
	"github.com/ydcloud-dy/leaf-api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// BlogUseCase 博客用户业务用例接口
type BlogUseCase interface {
	// Register 用户注册
	Register(req *dto.RegisterRequest) (*dto.LoginResponse, error)
	// Login 用户登录
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	// GetUserInfo 获取用户信息
	GetUserInfo(userID uint) (*dto.UserInfo, error)

	// GetArticleDetail 获取文章详情（包含用户点赞收藏状态）
	GetArticleDetail(articleID, userID uint) (*dto.ArticleDetailResponse, error)

	// LikeArticle 点赞文章
	LikeArticle(userID, articleID uint) error
	// UnlikeArticle 取消点赞
	UnlikeArticle(userID, articleID uint) error
	// IsLiked 检查是否已点赞
	IsLiked(userID, articleID uint) (bool, error)
	// GetUserLikes 获取用户点赞列表
	GetUserLikes(userID uint, page, limit int) (*dto.LikeListResponse, error)

	// FavoriteArticle 收藏文章
	FavoriteArticle(userID, articleID uint) error
	// UnfavoriteArticle 取消收藏
	UnfavoriteArticle(userID, articleID uint) error
	// IsFavorited 检查是否已收藏
	IsFavorited(userID, articleID uint) (bool, error)
	// GetUserFavorites 获取用户收藏列表
	GetUserFavorites(userID uint, page, limit int) (*dto.FavoriteListResponse, error)

	// CreateComment 创建评论
	CreateComment(req *dto.CreateCommentRequest) (*dto.CommentResponse, error)
	// GetArticleComments 获取文章评论列表
	GetArticleComments(articleID, userID uint, page, limit int) (*dto.CommentListResponse, error)
	// LikeComment 点赞评论
	LikeComment(userID, commentID uint) error
	// UnlikeComment 取消点赞评论
	UnlikeComment(userID, commentID uint) error
	// DeleteComment 删除评论
	DeleteComment(commentID, userID uint) error
	// GetUserStats 获取用户统计信息
	GetUserStats(userID uint) (*dto.UserStatsResponse, error)
}

// blogUseCase 博客用户业务用例实现
type blogUseCase struct {
	data *data.Data
}

// NewBlogUseCase 创建博客用户业务用例
func NewBlogUseCase(d *data.Data) BlogUseCase {
	return &blogUseCase{data: d}
}

// Register 用户注册
func (uc *blogUseCase) Register(req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	// 检查用户名是否已存在
	if _, err := uc.data.UserRepo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if _, err := uc.data.UserRepo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := &po.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Status:   1, // 默认启用
	}

	if err := uc.data.UserRepo.Create(user); err != nil {
		return nil, errors.New("创建用户失败")
	}

	// 生成 Token
	token, err := jwt.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		return nil, errors.New("生成 Token 失败")
	}

	return &dto.LoginResponse{
		Token: token,
		User: &dto.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// Login 用户登录
func (uc *blogUseCase) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 查询用户
	user, err := uc.data.UserRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查状态
	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	// 生成 Token
	token, err := jwt.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		return nil, errors.New("生成 Token 失败")
	}

	return &dto.LoginResponse{
		Token: token,
		User: &dto.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// GetUserInfo 获取用户信息
func (uc *blogUseCase) GetUserInfo(userID uint) (*dto.UserInfo, error) {
	user, err := uc.data.UserRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return &dto.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}, nil
}

// GetArticleDetail 获取文章详情（包含用户点赞收藏状态）
func (uc *blogUseCase) GetArticleDetail(articleID, userID uint) (*dto.ArticleDetailResponse, error) {
	// 获取文章基本信息
	article, err := uc.data.ArticleRepo.FindByIDWithRelations(articleID)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	// 博客前台只能查看已发布的文章（status = 1）
	if article.Status != 1 {
		return nil, errors.New("文章不存在或未发布")
	}

	// 增加浏览量（异步更新，不影响返回）
	go func() {
		_ = uc.data.ArticleRepo.IncrementViewCount(articleID)
	}()

	// 转换为响应结构
	articleResp := &dto.ArticleResponse{
		ID:              article.ID,
		Title:           article.Title,
		ContentMarkdown: article.ContentMarkdown,
		ContentHTML:     article.ContentHTML,
		Summary:         article.Summary,
		Cover:           article.Cover,
		AuthorID:        article.AuthorID,
		CategoryID:      article.CategoryID,
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
		articleResp.Author = &dto.AuthorInfo{
			ID:       article.Author.ID,
			Username: article.Author.Username,
			Avatar:   article.Author.Avatar,
		}
	}

	// 分类信息
	if article.Category.ID > 0 {
		articleResp.Category = &dto.CategoryInfo{
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
		articleResp.Tags = tags
	}

	// 检查用户点赞和收藏状态
	var isLiked, isFavorited bool
	if userID > 0 {
		isLiked, _ = uc.data.LikeRepo.Exists(articleID, userID)
		isFavorited, _ = uc.data.FavoriteRepo.Exists(articleID, userID)
	}

	return &dto.ArticleDetailResponse{
		ArticleResponse: *articleResp,
		IsLiked:         isLiked,
		IsFavorited:     isFavorited,
	}, nil
}

// LikeArticle 点赞文章
func (uc *blogUseCase) LikeArticle(userID, articleID uint) error {
	// 检查是否已点赞
	exists, err := uc.data.LikeRepo.Exists(articleID, userID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("已经点赞过了")
	}

	// 创建点赞记录
	like := &po.Like{
		UserID:    userID,
		ArticleID: articleID,
		CreatedAt: time.Now(),
	}
	if err := uc.data.LikeRepo.Create(like); err != nil {
		return err
	}

	// 更新文章点赞数
	return uc.data.ArticleRepo.IncrementLikeCount(articleID)
}

// UnlikeArticle 取消点赞
func (uc *blogUseCase) UnlikeArticle(userID, articleID uint) error {
	if err := uc.data.LikeRepo.Delete(articleID, userID); err != nil {
		return err
	}
	// 更新文章点赞数
	return uc.data.ArticleRepo.DecrementLikeCount(articleID)
}

// IsLiked 检查是否已点赞
func (uc *blogUseCase) IsLiked(userID, articleID uint) (bool, error) {
	return uc.data.LikeRepo.Exists(articleID, userID)
}

// GetUserLikes 获取用户点赞列表
func (uc *blogUseCase) GetUserLikes(userID uint, page, limit int) (*dto.LikeListResponse, error) {
	likes, total, err := uc.data.LikeRepo.ListByUser(userID, page, limit)
	if err != nil {
		return nil, err
	}

	likeList := make([]dto.LikeInfo, 0, len(likes))
	for _, like := range likes {
		likeList = append(likeList, dto.LikeInfo{
			ID:        like.ID,
			ArticleID: like.ArticleID,
			UserID:    like.UserID,
			CreatedAt: like.CreatedAt,
			Article: &dto.ArticleResponse{
				ID:      like.Article.ID,
				Title:   like.Article.Title,
				Summary: like.Article.Summary,
				Cover:   like.Article.Cover,
			},
		})
	}

	return &dto.LikeListResponse{
		List:  likeList,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// FavoriteArticle 收藏文章
func (uc *blogUseCase) FavoriteArticle(userID, articleID uint) error {
	// 检查是否已收藏
	exists, err := uc.data.FavoriteRepo.Exists(articleID, userID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("已经收藏过了")
	}

	// 创建收藏记录
	favorite := &po.Favorite{
		UserID:    userID,
		ArticleID: articleID,
		CreatedAt: time.Now(),
	}
	if err := uc.data.FavoriteRepo.Create(favorite); err != nil {
		return err
	}

	// 更新文章收藏数
	return uc.data.ArticleRepo.IncrementFavoriteCount(articleID)
}

// UnfavoriteArticle 取消收藏
func (uc *blogUseCase) UnfavoriteArticle(userID, articleID uint) error {
	if err := uc.data.FavoriteRepo.Delete(articleID, userID); err != nil {
		return err
	}
	// 更新文章收藏数
	return uc.data.ArticleRepo.DecrementFavoriteCount(articleID)
}

// IsFavorited 检查是否已收藏
func (uc *blogUseCase) IsFavorited(userID, articleID uint) (bool, error) {
	return uc.data.FavoriteRepo.Exists(articleID, userID)
}

// GetUserFavorites 获取用户收藏列表
func (uc *blogUseCase) GetUserFavorites(userID uint, page, limit int) (*dto.FavoriteListResponse, error) {
	favorites, total, err := uc.data.FavoriteRepo.ListByUser(userID, page, limit)
	if err != nil {
		return nil, err
	}

	favoriteList := make([]dto.FavoriteInfo, 0, len(favorites))
	for _, favorite := range favorites {
		favoriteList = append(favoriteList, dto.FavoriteInfo{
			ID:        favorite.ID,
			ArticleID: favorite.ArticleID,
			UserID:    favorite.UserID,
			CreatedAt: favorite.CreatedAt,
			Article: &dto.ArticleResponse{
				ID:      favorite.Article.ID,
				Title:   favorite.Article.Title,
				Summary: favorite.Article.Summary,
				Cover:   favorite.Article.Cover,
			},
		})
	}

	return &dto.FavoriteListResponse{
		List:  favoriteList,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// CreateComment 创建评论
func (uc *blogUseCase) CreateComment(req *dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	comment := &po.Comment{
		ArticleID:     req.ArticleID,
		UserID:        req.UserID,
		ParentID:      req.ParentID,
		ReplyToUserID: req.ReplyToUserID,
		Content:       req.Content,
		Status:        1, // 默认审核通过
		CreatedAt:     time.Now(),
	}

	if err := uc.data.CommentRepo.Create(comment); err != nil {
		return nil, err
	}

	// 更新文章评论数
	_ = uc.data.ArticleRepo.IncrementCommentCount(req.ArticleID)

	// 查询创建的评论（带用户信息）
	createdComment, err := uc.data.CommentRepo.FindByID(comment.ID)
	if err != nil {
		return nil, err
	}

	response := &dto.CommentResponse{
		ID:        createdComment.ID,
		ArticleID: createdComment.ArticleID,
		UserID:    createdComment.UserID,
		ParentID:  createdComment.ParentID,
		Content:   createdComment.Content,
		LikeCount: createdComment.LikeCount,
		Status:    createdComment.Status,
		CreatedAt: createdComment.CreatedAt,
		User: &dto.UserInfo{
			ID:       createdComment.User.ID,
			Username: createdComment.User.Username,
			Nickname: createdComment.User.Nickname,
			Avatar:   createdComment.User.Avatar,
		},
	}

	// 如果有回复目标用户，添加用户信息
	if createdComment.ReplyToUser != nil {
		response.ReplyToUser = &dto.UserInfo{
			ID:       createdComment.ReplyToUser.ID,
			Username: createdComment.ReplyToUser.Username,
			Nickname: createdComment.ReplyToUser.Nickname,
			Avatar:   createdComment.ReplyToUser.Avatar,
		}
	}

	return response, nil
}

// GetArticleComments 获取文章评论列表
func (uc *blogUseCase) GetArticleComments(articleID, userID uint, page, limit int) (*dto.CommentListResponse, error) {
	// 获取所有评论（不分页，为了构建完整的树形结构）
	comments, _, err := uc.data.CommentRepo.List(1, 1000, articleID, "1") // 只返回审核通过的
	if err != nil {
		return nil, err
	}

	// 将评论转换为map，以便快速查找
	commentMap := make(map[uint]*dto.CommentResponse)
	for _, comment := range comments {
		commentResp := &dto.CommentResponse{
			ID:        comment.ID,
			ArticleID: comment.ArticleID,
			UserID:    comment.UserID,
			ParentID:  comment.ParentID,
			Content:   comment.Content,
			LikeCount: comment.LikeCount,
			Status:    comment.Status,
			CreatedAt: comment.CreatedAt,
			User: &dto.UserInfo{
				ID:       comment.User.ID,
				Username: comment.User.Username,
				Nickname: comment.User.Nickname,
				Avatar:   comment.User.Avatar,
			},
			Replies: make([]dto.CommentResponse, 0),
		}

		// 添加回复目标用户信息
		if comment.ReplyToUser != nil {
			commentResp.ReplyToUser = &dto.UserInfo{
				ID:       comment.ReplyToUser.ID,
				Username: comment.ReplyToUser.Username,
				Nickname: comment.ReplyToUser.Nickname,
				Avatar:   comment.ReplyToUser.Avatar,
			}
		}

		// 检查当前用户是否已点赞该评论
		if userID > 0 {
			isLiked, _ := uc.data.CommentLikeRepo.Exists(comment.ID, userID)
			commentResp.IsLiked = isLiked
		}

		commentMap[comment.ID] = commentResp
	}

	// 构建树形结构
	topLevelComments := make([]dto.CommentResponse, 0)
	for _, comment := range comments {
		commentResp := commentMap[comment.ID]
		if comment.ParentID == nil {
			// 顶级评论
			topLevelComments = append(topLevelComments, *commentResp)
		} else {
			// 回复评论，添加到父评论的replies中
			if parent, ok := commentMap[*comment.ParentID]; ok {
				parent.Replies = append(parent.Replies, *commentResp)
			}
		}
	}

	// 应用分页到顶级评论
	start := (page - 1) * limit
	end := start + limit
	if start > len(topLevelComments) {
		start = len(topLevelComments)
	}
	if end > len(topLevelComments) {
		end = len(topLevelComments)
	}
	pagedComments := topLevelComments[start:end]

	// 重新构建完整的树（包括子评论）
	result := make([]dto.CommentResponse, 0, len(pagedComments))
	for _, topComment := range pagedComments {
		result = append(result, *commentMap[topComment.ID])
	}

	return &dto.CommentListResponse{
		List:  result,
		Total: int64(len(topLevelComments)), // 只统计顶级评论数
		Page:  page,
		Limit: limit,
	}, nil
}

// LikeComment 点赞评论
func (uc *blogUseCase) LikeComment(userID, commentID uint) error {
	// 检查是否已点赞
	exists, err := uc.data.CommentLikeRepo.Exists(commentID, userID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("已经点赞过了")
	}

	// 创建点赞记录
	like := &po.CommentLike{
		UserID:    userID,
		CommentID: commentID,
		CreatedAt: time.Now(),
	}

	// 使用事务：创建点赞记录 + 更新评论点赞数
	return uc.data.GetDB().Transaction(func(tx *gorm.DB) error {
		// 创建点赞记录
		if err := uc.data.CommentLikeRepo.Create(like); err != nil {
			return err
		}
		// 增加评论点赞数
		return tx.Model(&po.Comment{}).Where("id = ?", commentID).
			UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
	})
}

// UnlikeComment 取消点赞评论
func (uc *blogUseCase) UnlikeComment(userID, commentID uint) error {
	// 使用事务：删除点赞记录 + 更新评论点赞数
	return uc.data.GetDB().Transaction(func(tx *gorm.DB) error {
		// 删除点赞记录
		if err := uc.data.CommentLikeRepo.Delete(commentID, userID); err != nil {
			return err
		}
		// 减少评论点赞数
		return tx.Model(&po.Comment{}).Where("id = ?", commentID).
			UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
	})
}

// DeleteComment 删除评论（含权限检查）
func (uc *blogUseCase) DeleteComment(commentID, userID uint) error {
	// 查询评论
	comment, err := uc.data.CommentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}

	// 权限检查：只有评论作者本人可以删除自己的评论
	// 如果是子评论，父评论作者也可以删除
	canDelete := comment.UserID == userID

	if !canDelete && comment.ParentID != nil {
		// 检查是否为父评论作者
		parentComment, err := uc.data.CommentRepo.FindByID(*comment.ParentID)
		if err == nil && parentComment.UserID == userID {
			canDelete = true
		}
	}

	if !canDelete {
		return errors.New("无权删除该评论")
	}

	// 删除评论（如果是父评论，子评论会被级联删除）
	if err := uc.data.CommentRepo.Delete(commentID); err != nil {
		return err
	}

	// 更新文章评论数
	_ = uc.data.ArticleRepo.DecrementCommentCount(comment.ArticleID)

	return nil
}

// GetUserStats 获取用户统计信息
func (uc *blogUseCase) GetUserStats(userID uint) (*dto.UserStatsResponse, error) {
	// 获取点赞数
	likesCount, err := uc.data.LikeRepo.CountByUser(userID)
	if err != nil {
		likesCount = 0
	}

	// 获取收藏数
	favoritesCount, err := uc.data.FavoriteRepo.CountByUser(userID)
	if err != nil {
		favoritesCount = 0
	}

	// 获取评论数
	commentsCount, err := uc.data.CommentRepo.CountByUser(userID)
	if err != nil {
		commentsCount = 0
	}

	// 暂时将文章数设为0（可以后续根据需求添加）
	var articlesCount int64 = 0

	return &dto.UserStatsResponse{
		ArticlesCount:  articlesCount,
		LikesCount:     likesCount,
		FavoritesCount: favoritesCount,
		CommentsCount:  commentsCount,
	}, nil
}
