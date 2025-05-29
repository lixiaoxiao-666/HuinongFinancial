package repository

import (
	"context"
	"fmt"

	"huinong-backend/internal/model"

	"gorm.io/gorm"
)

// expertRepository 专家Repository实现
type expertRepository struct {
	db *gorm.DB
}

// NewExpertRepository 创建专家Repository实例
func NewExpertRepository(db *gorm.DB) ExpertRepository {
	return &expertRepository{
		db: db,
	}
}

// Create 创建专家
func (r *expertRepository) Create(ctx context.Context, expert *model.Expert) error {
	return r.db.WithContext(ctx).Create(expert).Error
}

// GetByID 根据ID获取专家
func (r *expertRepository) GetByID(ctx context.Context, id uint64) (*model.Expert, error) {
	var expert model.Expert
	err := r.db.WithContext(ctx).First(&expert, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("专家不存在")
		}
		return nil, err
	}
	return &expert, nil
}

// Update 更新专家
func (r *expertRepository) Update(ctx context.Context, expert *model.Expert) error {
	return r.db.WithContext(ctx).Save(expert).Error
}

// Delete 删除专家
func (r *expertRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Expert{}, id).Error
}

// List 专家列表查询
func (r *expertRepository) List(ctx context.Context, req *ListExpertsRequest) (*ListExpertsResponse, error) {
	var experts []*model.Expert
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Expert{})

	// 条件筛选
	if len(req.Specialties) > 0 {
		query = query.Where("specialty IN ?", req.Specialties)
	}
	if req.ServiceArea != "" {
		query = query.Where("service_area LIKE ?", "%"+req.ServiceArea+"%")
	}
	if req.IsVerified != nil {
		query = query.Where("is_verified = ?", *req.IsVerified)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		query = query.Where("expert_name LIKE ? OR description LIKE ?",
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
		Find(&experts).Error
	if err != nil {
		return nil, err
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
	query := r.db.WithContext(ctx).Where("status = ?", "active")

	if len(specialties) > 0 {
		query = query.Where("specialty IN ?", specialties)
	}
	if serviceArea != "" {
		query = query.Where("service_area LIKE ?", "%"+serviceArea+"%")
	}

	err := query.Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&experts).Error
	return experts, err
}

// GetVerifiedExperts 获取已认证专家
func (r *expertRepository) GetVerifiedExperts(ctx context.Context, limit, offset int) ([]*model.Expert, error) {
	var experts []*model.Expert
	err := r.db.WithContext(ctx).
		Where("is_verified = ? AND status = ?", true, "active").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&experts).Error
	return experts, err
}
