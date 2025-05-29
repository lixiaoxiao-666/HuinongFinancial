package repository

import (
	"context"
	"fmt"

	"huinong-backend/internal/model"

	"gorm.io/gorm"
)

// oaRepository OARepository实现
type oaRepository struct {
	db *gorm.DB
}

// NewOARepository 创建OARepository实例
func NewOARepository(db *gorm.DB) OARepository {
	return &oaRepository{
		db: db,
	}
}

// ==================== OA用户管理 ====================

// CreateOAUser 创建OA用户
func (r *oaRepository) CreateOAUser(ctx context.Context, user *model.OAUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetOAUserByID 根据ID获取OA用户
func (r *oaRepository) GetOAUserByID(ctx context.Context, id uint64) (*model.OAUser, error) {
	var user model.OAUser
	err := r.db.WithContext(ctx).
		Preload("Role").
		First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("OA用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// GetOAUserByUsername 根据用户名获取OA用户
func (r *oaRepository) GetOAUserByUsername(ctx context.Context, username string) (*model.OAUser, error) {
	var user model.OAUser
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("OA用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateOAUser 更新OA用户
func (r *oaRepository) UpdateOAUser(ctx context.Context, user *model.OAUser) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// DeleteOAUser 删除OA用户
func (r *oaRepository) DeleteOAUser(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.OAUser{}, id).Error
}

// ListOAUsers OA用户列表查询
func (r *oaRepository) ListOAUsers(ctx context.Context, req *ListOAUsersRequest) (*ListOAUsersResponse, error) {
	var users []*model.OAUser
	var total int64

	query := r.db.WithContext(ctx).Model(&model.OAUser{}).
		Preload("Role")

	// 条件筛选
	if req.RoleID > 0 {
		query = query.Where("role_id = ?", req.RoleID)
	}
	if req.Department != "" {
		query = query.Where("department = ?", req.Department)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ? OR email LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
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
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	return &ListOAUsersResponse{
		Users: users,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

// ==================== 角色管理 ====================

// CreateRole 创建角色
func (r *oaRepository) CreateRole(ctx context.Context, role *model.OARole) error {
	return r.db.WithContext(ctx).Create(role).Error
}

// GetRoleByID 根据ID获取角色
func (r *oaRepository) GetRoleByID(ctx context.Context, id uint64) (*model.OARole, error) {
	var role model.OARole
	err := r.db.WithContext(ctx).First(&role, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("角色不存在")
		}
		return nil, err
	}
	return &role, nil
}

// GetRoleByName 根据名称获取角色
func (r *oaRepository) GetRoleByName(ctx context.Context, name string) (*model.OARole, error) {
	var role model.OARole
	err := r.db.WithContext(ctx).Where("role_name = ?", name).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("角色不存在")
		}
		return nil, err
	}
	return &role, nil
}

// UpdateRole 更新角色
func (r *oaRepository) UpdateRole(ctx context.Context, role *model.OARole) error {
	return r.db.WithContext(ctx).Save(role).Error
}

// DeleteRole 删除角色
func (r *oaRepository) DeleteRole(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.OARole{}, id).Error
}

// ListRoles 获取角色列表
func (r *oaRepository) ListRoles(ctx context.Context) ([]*model.OARole, error) {
	var roles []*model.OARole
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("sort_order ASC, created_at ASC").
		Find(&roles).Error
	return roles, err
}
