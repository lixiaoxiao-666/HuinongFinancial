package service

import (
	"context"
	"fmt"

	"huinong-backend/internal/model"
	"huinong-backend/internal/repository"
)

// contentService 内容管理服务实现
type contentService struct {
	articleRepo repository.ArticleRepository
	expertRepo  repository.ExpertRepository
}

// NewContentService 创建内容管理服务实例
func NewContentService(articleRepo repository.ArticleRepository, expertRepo repository.ExpertRepository) ContentService {
	return &contentService{
		articleRepo: articleRepo,
		expertRepo:  expertRepo,
	}
}

// ==================== 文章管理 ====================

// CreateArticle 创建文章
func (s *contentService) CreateArticle(ctx context.Context, req *CreateArticleRequest) (*CreateArticleResponse, error) {
	article := &model.Article{
		Title:          req.Title,
		Subtitle:       req.Subtitle,
		Content:        req.Content,
		Summary:        req.Summary,
		Category:       req.Category,
		CoverImage:     req.CoverImage,
		IsTop:          req.IsTop,
		IsFeatured:     req.IsFeatured,
		SEOTitle:       req.SEOTitle,
		SEODescription: req.SEODescription,
		SEOKeywords:    req.SEOKeywords,
		Status:         "draft", // 默认草稿状态
		// TODO: 设置作者ID从上下文中获取
		// AuthorID: getUserIDFromContext(ctx),
	}

	// TODO: 处理标签
	// if len(req.Tags) > 0 {
	//     article.Tags = strings.Join(req.Tags, ",")
	// }

	err := s.articleRepo.Create(ctx, article)
	if err != nil {
		return nil, fmt.Errorf("创建文章失败: %v", err)
	}

	return &CreateArticleResponse{
		ID:     article.ID,
		Title:  article.Title,
		Status: article.Status,
	}, nil
}

// GetArticle 获取文章详情
func (s *contentService) GetArticle(ctx context.Context, articleID uint64) (*ArticleDetailResponse, error) {
	article, err := s.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return nil, fmt.Errorf("获取文章失败: %v", err)
	}

	// 增加阅读量
	s.articleRepo.IncrementViewCount(ctx, articleID)

	// TODO: 获取作者信息
	// author, _ := s.oaRepo.GetOAUserByID(ctx, article.AuthorID)

	return &ArticleDetailResponse{
		Article: article,
		// Author:  author,
	}, nil
}

// UpdateArticle 更新文章
func (s *contentService) UpdateArticle(ctx context.Context, articleID uint64, req *UpdateArticleRequest) error {
	article, err := s.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return fmt.Errorf("获取文章失败: %v", err)
	}

	// 更新字段
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Subtitle != "" {
		article.Subtitle = req.Subtitle
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.Summary != "" {
		article.Summary = req.Summary
	}
	if req.Category != "" {
		article.Category = req.Category
	}
	if req.CoverImage != "" {
		article.CoverImage = req.CoverImage
	}
	article.IsTop = req.IsTop
	article.IsFeatured = req.IsFeatured
	if req.SEOTitle != "" {
		article.SEOTitle = req.SEOTitle
	}
	if req.SEODescription != "" {
		article.SEODescription = req.SEODescription
	}
	if req.SEOKeywords != "" {
		article.SEOKeywords = req.SEOKeywords
	}

	// TODO: 处理标签更新
	// if len(req.Tags) > 0 {
	//     article.Tags = strings.Join(req.Tags, ",")
	// }

	return s.articleRepo.Update(ctx, article)
}

// DeleteArticle 删除文章
func (s *contentService) DeleteArticle(ctx context.Context, articleID uint64) error {
	return s.articleRepo.Delete(ctx, articleID)
}

// ListArticles 获取文章列表
func (s *contentService) ListArticles(ctx context.Context, req *ListArticlesRequest) (*ListArticlesResponse, error) {
	// 转换请求参数
	repoReq := &repository.ListArticlesRequest{
		Page:       req.Page,
		Limit:      req.Limit,
		Category:   req.Category,
		Status:     req.Status,
		IsTop:      req.IsTop,
		IsFeatured: req.IsFeatured,
		Keyword:    req.Keyword,
	}

	result, err := s.articleRepo.List(ctx, repoReq)
	if err != nil {
		return nil, fmt.Errorf("获取文章列表失败: %v", err)
	}

	return &ListArticlesResponse{
		Articles: result.Articles,
		Total:    result.Total,
		Page:     result.Page,
		Limit:    result.Limit,
	}, nil
}

// PublishArticle 发布文章
func (s *contentService) PublishArticle(ctx context.Context, articleID uint64) error {
	article, err := s.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return fmt.Errorf("获取文章失败: %v", err)
	}

	article.Status = "published"
	return s.articleRepo.Update(ctx, article)
}

// GetFeaturedArticles 获取推荐文章
func (s *contentService) GetFeaturedArticles(ctx context.Context, limit int) (*FeaturedArticlesResponse, error) {
	articles, err := s.articleRepo.GetFeatured(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("获取推荐文章失败: %v", err)
	}

	return &FeaturedArticlesResponse{
		Articles: articles,
	}, nil
}

// ==================== 分类管理 ====================

// CreateCategory 创建分类
func (s *contentService) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*model.Category, error) {
	category := &model.Category{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		ParentID:    req.ParentID,
		Icon:        req.Icon,
		SortOrder:   req.SortOrder,
		Status:      "active",
	}

	err := s.articleRepo.CreateCategory(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("创建分类失败: %v", err)
	}

	return category, nil
}

// GetCategories 获取分类列表
func (s *contentService) GetCategories(ctx context.Context) ([]*model.Category, error) {
	return s.articleRepo.ListCategories(ctx)
}

// UpdateCategory 更新分类
func (s *contentService) UpdateCategory(ctx context.Context, categoryID uint64, req *UpdateCategoryRequest) error {
	category, err := s.articleRepo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("获取分类失败: %v", err)
	}

	if req.DisplayName != "" {
		category.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	if req.SortOrder > 0 {
		category.SortOrder = req.SortOrder
	}
	if req.Status != "" {
		category.Status = req.Status
	}

	return s.articleRepo.UpdateCategory(ctx, category)
}

// DeleteCategory 删除分类
func (s *contentService) DeleteCategory(ctx context.Context, categoryID uint64) error {
	return s.articleRepo.DeleteCategory(ctx, categoryID)
}

// ==================== 专家管理 ====================

// CreateExpert 创建专家
func (s *contentService) CreateExpert(ctx context.Context, req *CreateExpertRequest) (*model.Expert, error) {
	expert := &model.Expert{
		Name:            req.Name,
		Title:           req.Title,
		Organization:    req.Organization,
		Phone:           req.Phone,
		Email:           req.Email,
		WeChat:          req.WeChat,
		Avatar:          req.Avatar,
		Biography:       req.Biography,
		ExperienceYears: req.ExperienceYears,
		Status:          "active",
		IsVerified:      false, // 默认未认证
	}

	// TODO: 处理专业领域和服务地区
	// expert.Specialties = strings.Join(req.Specialties, ",")
	// expert.ServiceAreas = strings.Join(req.ServiceAreas, ",")

	err := s.expertRepo.Create(ctx, expert)
	if err != nil {
		return nil, fmt.Errorf("创建专家失败: %v", err)
	}

	return expert, nil
}

// GetExpert 获取专家详情
func (s *contentService) GetExpert(ctx context.Context, expertID uint64) (*ExpertDetailResponse, error) {
	expert, err := s.expertRepo.GetByID(ctx, expertID)
	if err != nil {
		return nil, fmt.Errorf("获取专家信息失败: %v", err)
	}

	return &ExpertDetailResponse{
		Expert: expert,
	}, nil
}

// ListExperts 获取专家列表
func (s *contentService) ListExperts(ctx context.Context, req *ListExpertsRequest) (*ListExpertsResponse, error) {
	// 转换请求参数
	repoReq := &repository.ListExpertsRequest{
		Page:        req.Page,
		Limit:       req.Limit,
		Specialties: []string{req.Specialty}, // 将单个specialty转为数组
		ServiceArea: req.Province,            // 使用Province作为ServiceArea
		IsVerified:  nil,                     // service层没有此字段
		Status:      req.Status,
		Keyword:     "", // service层没有此字段
	}

	result, err := s.expertRepo.List(ctx, repoReq)
	if err != nil {
		return nil, fmt.Errorf("获取专家列表失败: %v", err)
	}

	return &ListExpertsResponse{
		Experts: result.Experts,
		Total:   result.Total,
		Page:    result.Page,
		Limit:   result.Limit,
	}, nil
}

// UpdateExpert 更新专家信息
func (s *contentService) UpdateExpert(ctx context.Context, expertID uint64, req *UpdateExpertRequest) error {
	expert, err := s.expertRepo.GetByID(ctx, expertID)
	if err != nil {
		return fmt.Errorf("获取专家信息失败: %v", err)
	}

	// 更新字段
	if req.Name != "" {
		expert.Name = req.Name
	}
	if req.Title != "" {
		expert.Title = req.Title
	}
	if req.Organization != "" {
		expert.Organization = req.Organization
	}
	if req.Phone != "" {
		expert.Phone = req.Phone
	}
	if req.Email != "" {
		expert.Email = req.Email
	}
	if req.WeChat != "" {
		expert.WeChat = req.WeChat
	}
	if req.Avatar != "" {
		expert.Avatar = req.Avatar
	}
	if req.Biography != "" {
		expert.Biography = req.Biography
	}
	if req.ExperienceYears > 0 {
		expert.ExperienceYears = req.ExperienceYears
	}
	if req.Status != "" {
		expert.Status = req.Status
	}

	// TODO: 更新专业领域和服务地区
	// if len(req.Specialties) > 0 {
	//     expert.Specialties = strings.Join(req.Specialties, ",")
	// }
	// if len(req.ServiceAreas) > 0 {
	//     expert.ServiceAreas = strings.Join(req.ServiceAreas, ",")
	// }

	return s.expertRepo.Update(ctx, expert)
}

// DeleteExpert 删除专家
func (s *contentService) DeleteExpert(ctx context.Context, expertID uint64) error {
	return s.expertRepo.Delete(ctx, expertID)
}

// ==================== 专家咨询 ====================

// SubmitConsultation 提交专家咨询
func (s *contentService) SubmitConsultation(ctx context.Context, req *SubmitConsultationRequest) (*SubmitConsultationResponse, error) {
	// TODO: 实现专家咨询逻辑
	// 1. 验证专家是否存在且可用
	// 2. 创建咨询记录
	// 3. 发送通知给专家

	return &SubmitConsultationResponse{
		ID:     1, // 临时返回
		Status: "pending",
	}, nil
}

// GetConsultations 获取用户咨询记录
func (s *contentService) GetConsultations(ctx context.Context, userID uint64) (*ConsultationsResponse, error) {
	// TODO: 实现获取咨询记录逻辑
	return &ConsultationsResponse{
		Consultations: []ConsultationDetail{},
		Total:         0,
	}, nil
}
