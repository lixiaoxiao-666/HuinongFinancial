package repository

import (
	"context"
	"fmt"

	"huinong-backend/internal/model"

	"gorm.io/gorm"
)

// contentRepository 内容Repository实现
type contentRepository struct {
	db *gorm.DB
}

// NewContentRepository 创建内容Repository实例
func NewContentRepository(db *gorm.DB) ArticleRepository {
	return &contentRepository{
		db: db,
	}
}

// ==================== 文章管理 ====================

// Create 创建文章
func (r *contentRepository) Create(ctx context.Context, article *model.Article) error {
	return r.db.WithContext(ctx).Create(article).Error
}

// GetByID 根据ID获取文章
func (r *contentRepository) GetByID(ctx context.Context, id uint64) (*model.Article, error) {
	var article model.Article
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Author").
		First(&article, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("文章不存在")
		}
		return nil, err
	}
	return &article, nil
}

// Update 更新文章
func (r *contentRepository) Update(ctx context.Context, article *model.Article) error {
	return r.db.WithContext(ctx).Save(article).Error
}

// Delete 删除文章
func (r *contentRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Article{}, id).Error
}

// List 文章列表查询
func (r *contentRepository) List(ctx context.Context, req *ListArticlesRequest) (*ListArticlesResponse, error) {
	var articles []*model.Article
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Article{}).
		Preload("Category").
		Preload("Author")

	// 条件筛选
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.AuthorID > 0 {
		query = query.Where("author_id = ?", req.AuthorID)
	}
	if req.IsTop != nil {
		query = query.Where("is_top = ?", *req.IsTop)
	}
	if req.IsFeatured != nil {
		query = query.Where("is_featured = ?", *req.IsFeatured)
	}
	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR summary LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.Limit
	err = query.Order("created_at DESC").
		Offset(offset).Limit(req.Limit).
		Find(&articles).Error
	if err != nil {
		return nil, err
	}

	return &ListArticlesResponse{
		Articles: articles,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}

// ==================== 分类管理 ====================

// CreateCategory 创建分类
func (r *contentRepository) CreateCategory(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// GetCategoryByID 根据ID获取分类
func (r *contentRepository) GetCategoryByID(ctx context.Context, id uint64) (*model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("分类不存在")
		}
		return nil, err
	}
	return &category, nil
}

// GetCategoryByName 根据名称获取分类
func (r *contentRepository) GetCategoryByName(ctx context.Context, name string) (*model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("分类不存在")
		}
		return nil, err
	}
	return &category, nil
}

// UpdateCategory 更新分类
func (r *contentRepository) UpdateCategory(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// DeleteCategory 删除分类
func (r *contentRepository) DeleteCategory(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Category{}, id).Error
}

// ListCategories 获取分类列表
func (r *contentRepository) ListCategories(ctx context.Context) ([]*model.Category, error) {
	var categories []*model.Category
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("sort_order ASC, created_at DESC").
		Find(&categories).Error
	return categories, err
}

// ==================== 内容查询 ====================

// GetByCategory 根据分类获取文章
func (r *contentRepository) GetByCategory(ctx context.Context, category string, limit, offset int) ([]*model.Article, error) {
	var articles []*model.Article
	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("category = ? AND status = ?", category, "published").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&articles).Error
	return articles, err
}

// GetFeatured 获取推荐文章
func (r *contentRepository) GetFeatured(ctx context.Context, limit int) ([]*model.Article, error) {
	var articles []*model.Article
	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("is_featured = ? AND status = ?", true, "published").
		Order("created_at DESC").
		Limit(limit).
		Find(&articles).Error
	return articles, err
}

// GetTopArticles 获取置顶文章
func (r *contentRepository) GetTopArticles(ctx context.Context, limit int) ([]*model.Article, error) {
	var articles []*model.Article
	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("is_top = ? AND status = ?", true, "published").
		Order("created_at DESC").
		Limit(limit).
		Find(&articles).Error
	return articles, err
}

// Search 搜索文章
func (r *contentRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*model.Article, error) {
	var articles []*model.Article
	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("(title LIKE ? OR content LIKE ?) AND status = ?",
			"%"+keyword+"%", "%"+keyword+"%", "published").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&articles).Error
	return articles, err
}

// ==================== 统计更新 ====================

// IncrementViewCount 增加浏览次数
func (r *contentRepository) IncrementViewCount(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.Article{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// IncrementLikeCount 增加点赞次数
func (r *contentRepository) IncrementLikeCount(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.Article{}).
		Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

// IncrementShareCount 增加分享次数
func (r *contentRepository) IncrementShareCount(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.Article{}).
		Where("id = ?", id).
		UpdateColumn("share_count", gorm.Expr("share_count + 1")).Error
}
