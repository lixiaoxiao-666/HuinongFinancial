package repository

import (
	"context"
	"fmt"

	"huinong-backend/internal/model"

	"gorm.io/gorm"
)

// expertRepository 专家数据访问层实现
type expertRepository struct {
	db *gorm.DB
}

// NewExpertRepository 创建专家数据访问层实例
func NewExpertRepository(db *gorm.DB) ExpertRepository {
	return &expertRepository{db: db}
}

// Create 创建专家
func (r *expertRepository) Create(ctx context.Context, expert *model.Expert) error {
	if err := r.db.WithContext(ctx).Create(expert).Error; err != nil {
		return fmt.Errorf("创建专家失败: %w", err)
	}
	return nil
}

// GetByID 根据ID获取专家
func (r *expertRepository) GetByID(ctx context.Context, id uint64) (*model.Expert, error) {
	var expert model.Expert
	if err := r.db.WithContext(ctx).First(&expert, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("专家不存在")
		}
		return nil, fmt.Errorf("获取专家失败: %w", err)
	}
	return &expert, nil
}

// Update 更新专家
func (r *expertRepository) Update(ctx context.Context, expert *model.Expert) error {
	if err := r.db.WithContext(ctx).Save(expert).Error; err != nil {
		return fmt.Errorf("更新专家失败: %w", err)
	}
	return nil
}

// Delete 删除专家
func (r *expertRepository) Delete(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Delete(&model.Expert{}, id).Error; err != nil {
		return fmt.Errorf("删除专家失败: %w", err)
	}
	return nil
}

// List 获取专家列表
func (r *expertRepository) List(ctx context.Context, req *ListExpertsRequest) (*ListExpertsResponse, error) {
	var experts []*model.Expert
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Expert{})

	// 添加查询条件
	if len(req.Specialties) > 0 {
		// 专业领域匹配（JSON字段查询）
		query = query.Where("JSON_CONTAINS(specialties, ?) OR JSON_OVERLAPS(specialties, ?)",
			fmt.Sprintf(`["%s"]`, req.Specialties[0]), fmt.Sprintf(`["%s"]`, req.Specialties[0]))
	}
	if req.ServiceArea != "" {
		// 服务地区匹配（JSON字段查询）
		query = query.Where("JSON_CONTAINS(service_areas, ?)", fmt.Sprintf(`["%s"]`, req.ServiceArea))
	}
	if req.IsVerified != nil {
		query = query.Where("is_verified = ?", *req.IsVerified)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("name LIKE ? OR title LIKE ? OR organization LIKE ? OR biography LIKE ?",
			keyword, keyword, keyword, keyword)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取专家总数失败: %w", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.Limit
	if err := query.
		Order("is_verified DESC, rating DESC, created_at DESC").
		Limit(req.Limit).
		Offset(offset).
		Find(&experts).Error; err != nil {
		return nil, fmt.Errorf("获取专家列表失败: %w", err)
	}

	return &ListExpertsResponse{
		Experts: experts,
		Total:   total,
		Page:    req.Page,
		Limit:   req.Limit,
	}, nil
}

// SearchBySpecialty 根据专业领域搜索专家
func (r *expertRepository) SearchBySpecialty(ctx context.Context, specialties []string, serviceArea string, limit, offset int) ([]*model.Expert, error) {
	var experts []*model.Expert
	query := r.db.WithContext(ctx).Model(&model.Expert{}).
		Where("status = ? AND is_verified = ?", "active", true)

	// 专业领域匹配
	if len(specialties) > 0 {
		specialtyConditions := make([]string, len(specialties))
		args := make([]interface{}, len(specialties))
		for i, specialty := range specialties {
			specialtyConditions[i] = "JSON_CONTAINS(specialties, ?)"
			args[i] = fmt.Sprintf(`["%s"]`, specialty)
		}
		query = query.Where(fmt.Sprintf("(%s)",
			fmt.Sprintf("(%s)", fmt.Sprintf("%s", specialtyConditions[0]))), args[0])
	}

	// 服务地区匹配
	if serviceArea != "" {
		query = query.Where("JSON_CONTAINS(service_areas, ?)", fmt.Sprintf(`["%s"]`, serviceArea))
	}

	if err := query.
		Order("rating DESC, experience_years DESC, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&experts).Error; err != nil {
		return nil, fmt.Errorf("搜索专家失败: %w", err)
	}

	return experts, nil
}

// GetVerifiedExperts 获取已认证专家
func (r *expertRepository) GetVerifiedExperts(ctx context.Context, limit, offset int) ([]*model.Expert, error) {
	var experts []*model.Expert
	if err := r.db.WithContext(ctx).
		Where("is_verified = ? AND status = ?", true, "active").
		Order("rating DESC, experience_years DESC").
		Limit(limit).
		Offset(offset).
		Find(&experts).Error; err != nil {
		return nil, fmt.Errorf("获取认证专家失败: %w", err)
	}
	return experts, nil
}
