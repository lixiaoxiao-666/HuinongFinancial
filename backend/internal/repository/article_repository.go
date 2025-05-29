package repository

import (
	"context"
	"fmt"

	"huinong-backend/internal/model"

	"gorm.io/gorm"
)

// articleRepository 文章数据访问层实现
type articleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章数据访问层实例
func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

// ==================== 文章管理 ====================

// Create 创建文章
func (r *articleRepository) Create(ctx context.Context, article *model.Article) error {
	if err := r.db.WithContext(ctx).Create(article).Error; err != nil {
		return fmt.Errorf("创建文章失败: %w", err)
	}
	return nil
}

// GetByID 根据ID获取文章
func (r *articleRepository) GetByID(ctx context.Context, id uint64) (*model.Article, error) {
	var article model.Article
	if err := r.db.WithContext(ctx).
		Preload("Author").
		First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("文章不存在")
		}
		return nil, fmt.Errorf("获取文章失败: %w", err)
	}
	return &article, nil
}

// Update 更新文章
func (r *articleRepository) Update(ctx context.Context, article *model.Article) error {
	if err := r.db.WithContext(ctx).Save(article).Error; err != nil {
		return fmt.Errorf("更新文章失败: %w", err)
	}
	return nil
}

// Delete 删除文章
func (r *articleRepository) Delete(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Delete(&model.Article{}, id).Error; err != nil {
		return fmt.Errorf("删除文章失败: %w", err)
	}
	return nil
}

// List 获取文章列表
func (r *articleRepository) List(ctx context.Context, req *ListArticlesRequest) (*ListArticlesResponse, error) {
	var articles []*model.Article
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Article{})

	// 添加查询条件
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
		keyword := "%" + req.Keyword + "%"
		query = query.Where("title LIKE ? OR content LIKE ? OR summary LIKE ?", keyword, keyword, keyword)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取文章总数失败: %w", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.Limit
	if err := query.
		Preload("Author").
		Order("is_top DESC, created_at DESC").
		Limit(req.Limit).
		Offset(offset).
		Find(&articles).Error; err != nil {
		return nil, fmt.Errorf("获取文章列表失败: %w", err)
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
func (r *articleRepository) CreateCategory(ctx context.Context, category *model.Category) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return fmt.Errorf("创建分类失败: %w", err)
	}
	return nil
}

// GetCategoryByID 根据ID获取分类
func (r *articleRepository) GetCategoryByID(ctx context.Context, id uint64) (*model.Category, error) {
	var category model.Category
	if err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("分类不存在")
		}
		return nil, fmt.Errorf("获取分类失败: %w", err)
	}
	return &category, nil
}

// GetCategoryByName 根据名称获取分类
func (r *articleRepository) GetCategoryByName(ctx context.Context, name string) (*model.Category, error) {
	var category model.Category
	if err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("分类不存在")
		}
		return nil, fmt.Errorf("获取分类失败: %w", err)
	}
	return &category, nil
}

// UpdateCategory 更新分类
func (r *articleRepository) UpdateCategory(ctx context.Context, category *model.Category) error {
	if err := r.db.WithContext(ctx).Save(category).Error; err != nil {
		return fmt.Errorf("更新分类失败: %w", err)
	}
	return nil
}

// DeleteCategory 删除分类
func (r *articleRepository) DeleteCategory(ctx context.Context, id uint64) error {
	// 检查是否有子分类
	var childCount int64
	if err := r.db.WithContext(ctx).
		Model(&model.Category{}).
		Where("parent_id = ?", id).
		Count(&childCount).Error; err != nil {
		return fmt.Errorf("检查子分类失败: %w", err)
	}
	if childCount > 0 {
		return fmt.Errorf("存在子分类，无法删除")
	}

	// 检查是否有文章使用该分类
	var articleCount int64
	if err := r.db.WithContext(ctx).
		Model(&model.Article{}).
		Where("category = (SELECT name FROM categories WHERE id = ?)", id).
		Count(&articleCount).Error; err != nil {
		return fmt.Errorf("检查关联文章失败: %w", err)
	}
	if articleCount > 0 {
		return fmt.Errorf("存在关联文章，无法删除")
	}

	if err := r.db.WithContext(ctx).Delete(&model.Category{}, id).Error; err != nil {
		return fmt.Errorf("删除分类失败: %w", err)
	}
	return nil
}

// ListCategories 获取分类列表
func (r *articleRepository) ListCategories(ctx context.Context) ([]*model.Category, error) {
	var categories []*model.Category
	if err := r.db.WithContext(ctx).
		Where("status = ?", "active").
		Order("sort_order ASC, created_at ASC").
		Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("获取分类列表失败: %w", err)
	}
	return categories, nil
}

// ==================== 内容查询 ====================

// GetByCategory 根据分类获取文章
func (r *articleRepository) GetByCategory(ctx context.Context, category string, limit, offset int) ([]*model.Article, error) {
	var articles []*model.Article
	if err := r.db.WithContext(ctx).
		Where("category = ? AND status = ?", category, "published").
		Order("is_top DESC, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&articles).Error; err != nil {
		return nil, fmt.Errorf("获取分类文章失败: %w", err)
	}
	return articles, nil
}

// GetFeatured 获取推荐文章
func (r *articleRepository) GetFeatured(ctx context.Context, limit int) ([]*model.Article, error) {
	var articles []*model.Article
	if err := r.db.WithContext(ctx).
		Where("is_featured = ? AND status = ?", true, "published").
		Order("created_at DESC").
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, fmt.Errorf("获取推荐文章失败: %w", err)
	}
	return articles, nil
}

// GetTopArticles 获取置顶文章
func (r *articleRepository) GetTopArticles(ctx context.Context, limit int) ([]*model.Article, error) {
	var articles []*model.Article
	if err := r.db.WithContext(ctx).
		Where("is_top = ? AND status = ?", true, "published").
		Order("created_at DESC").
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, fmt.Errorf("获取置顶文章失败: %w", err)
	}
	return articles, nil
}

// Search 搜索文章
func (r *articleRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*model.Article, error) {
	var articles []*model.Article
	searchPattern := "%" + keyword + "%"
	if err := r.db.WithContext(ctx).
		Where("(title LIKE ? OR content LIKE ? OR summary LIKE ?) AND status = ?",
			searchPattern, searchPattern, searchPattern, "published").
		Order("view_count DESC, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&articles).Error; err != nil {
		return nil, fmt.Errorf("搜索文章失败: %w", err)
	}
	return articles, nil
}

// ==================== 统计更新 ====================

// IncrementViewCount 增加阅读量
func (r *articleRepository) IncrementViewCount(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).
		Model(&model.Article{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error; err != nil {
		return fmt.Errorf("更新阅读量失败: %w", err)
	}
	return nil
}

// IncrementLikeCount 增加点赞量
func (r *articleRepository) IncrementLikeCount(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).
		Model(&model.Article{}).
		Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
		return fmt.Errorf("更新点赞量失败: %w", err)
	}
	return nil
}

// IncrementShareCount 增加分享量
func (r *articleRepository) IncrementShareCount(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).
		Model(&model.Article{}).
		Where("id = ?", id).
		UpdateColumn("share_count", gorm.Expr("share_count + 1")).Error; err != nil {
		return fmt.Errorf("更新分享量失败: %w", err)
	}
	return nil
}
