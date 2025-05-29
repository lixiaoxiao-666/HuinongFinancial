package repository

import (
	"context"
	"huinong-backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户Repository实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByPhone 根据手机号获取用户
func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUUID 根据UUID获取用户
func (r *userRepository) GetByUUID(ctx context.Context, uuid string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户
func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// List 获取用户列表
func (r *userRepository) List(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})

	// 添加查询条件
	if req.UserType != "" {
		query = query.Where("user_type = ?", req.UserType)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("username LIKE ? OR real_name LIKE ? OR phone LIKE ?", keyword, keyword, keyword)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, err
	}

	return &ListUsersResponse{
		Users: users,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

// GetByUserType 根据用户类型获取用户
func (r *userRepository) GetByUserType(ctx context.Context, userType string, limit, offset int) ([]*model.User, error) {
	var users []*model.User
	err := r.db.WithContext(ctx).Where("user_type = ?", userType).
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	return users, err
}

// GetByStatus 根据状态获取用户
func (r *userRepository) GetByStatus(ctx context.Context, status string, limit, offset int) ([]*model.User, error) {
	var users []*model.User
	err := r.db.WithContext(ctx).Where("status = ?", status).
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	return users, err
}

// GetUserAuth 获取用户认证信息
func (r *userRepository) GetUserAuth(ctx context.Context, userID uint64, authType string) (*model.UserAuth, error) {
	var auth model.UserAuth
	err := r.db.WithContext(ctx).Where("user_id = ? AND auth_type = ?", userID, authType).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// CreateUserAuth 创建用户认证信息
func (r *userRepository) CreateUserAuth(ctx context.Context, auth *model.UserAuth) error {
	return r.db.WithContext(ctx).Create(auth).Error
}

// UpdateUserAuth 更新用户认证信息
func (r *userRepository) UpdateUserAuth(ctx context.Context, auth *model.UserAuth) error {
	return r.db.WithContext(ctx).Save(auth).Error
}

// CreateSession 创建用户会话
func (r *userRepository) CreateSession(ctx context.Context, session *model.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetSession 获取用户会话
func (r *userRepository) GetSession(ctx context.Context, sessionID string) (*model.UserSession, error) {
	var session model.UserSession
	err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// UpdateSession 更新用户会话
func (r *userRepository) UpdateSession(ctx context.Context, session *model.UserSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// DeleteSession 删除用户会话
func (r *userRepository) DeleteSession(ctx context.Context, sessionID string) error {
	return r.db.WithContext(ctx).Where("session_id = ?", sessionID).Delete(&model.UserSession{}).Error
}

// GetUserSessions 获取用户所有会话
func (r *userRepository) GetUserSessions(ctx context.Context, userID uint64) ([]*model.UserSession, error) {
	var sessions []*model.UserSession
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&sessions).Error
	return sessions, err
}

// AddUserTag 添加用户标签
func (r *userRepository) AddUserTag(ctx context.Context, tag *model.UserTag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

// GetUserTags 获取用户标签
func (r *userRepository) GetUserTags(ctx context.Context, userID uint64, tagType string) ([]*model.UserTag, error) {
	var tags []*model.UserTag
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if tagType != "" {
		query = query.Where("tag_type = ?", tagType)
	}
	err := query.Find(&tags).Error
	return tags, err
}

// RemoveUserTag 移除用户标签
func (r *userRepository) RemoveUserTag(ctx context.Context, userID uint64, tagKey string) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND tag_key = ?", userID, tagKey).
		Delete(&model.UserTag{}).Error
}

// GetUserCount 获取用户总数
func (r *userRepository) GetUserCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error
	return count, err
}

// GetUserCountByType 根据类型获取用户数量
func (r *userRepository) GetUserCountByType(ctx context.Context, userType string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("user_type = ?", userType).Count(&count).Error
	return count, err
}

// UpdateLoginInfo 更新用户登录信息
func (r *userRepository) UpdateLoginInfo(ctx context.Context, userID uint64, loginIP string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"last_login_time": time.Now(),
		"last_login_ip":   loginIP,
		"login_count":     gorm.Expr("login_count + 1"),
	}).Error
}
