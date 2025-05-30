package repository

import (
	"context"
	"time"

	"huinong-backend/internal/model"

	"gorm.io/gorm"
)

// SessionRepository 会话存储库接口
type SessionRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, session *model.UserSession) error
	GetByID(ctx context.Context, id uint64) (*model.UserSession, error)
	GetBySessionID(ctx context.Context, sessionID string) (*model.UserSession, error)
	Update(ctx context.Context, session *model.UserSession) error
	Delete(ctx context.Context, id uint64) error

	// 批量操作
	GetByUserID(ctx context.Context, userID uint64, status string) ([]*model.UserSession, error)
	GetActiveSessionsByUserID(ctx context.Context, userID uint64) ([]*model.UserSession, error)
	DeleteExpiredSessions(ctx context.Context, expiredBefore time.Time) error

	// 状态管理
	UpdateStatus(ctx context.Context, sessionID string, status string, logoutTime time.Time) error
	UpdateLastActive(ctx context.Context, sessionID string, lastActiveAt time.Time) error
	UpdateTokens(ctx context.Context, sessionID string, accessTokenHash, refreshTokenHash string, accessExpiresAt, refreshExpiresAt time.Time) error

	// 统计查询
	CountActiveSessionsByUserID(ctx context.Context, userID uint64) (int64, error)
	GetRecentLoginsByUserID(ctx context.Context, userID uint64, limit int) ([]*model.UserSession, error)

	// 安全相关
	GetSessionsByIPAddress(ctx context.Context, ipAddress string, limit int) ([]*model.UserSession, error)
	GetSessionsByDeviceID(ctx context.Context, deviceID string) ([]*model.UserSession, error)
}

// sessionRepository 会话存储库实现
type sessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository 创建会话存储库实例
func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

// Create 创建会话记录
func (r *sessionRepository) Create(ctx context.Context, session *model.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetByID 根据ID获取会话
func (r *sessionRepository) GetByID(ctx context.Context, id uint64) (*model.UserSession, error) {
	var session model.UserSession
	err := r.db.WithContext(ctx).First(&session, id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetBySessionID 根据会话ID获取会话
func (r *sessionRepository) GetBySessionID(ctx context.Context, sessionID string) (*model.UserSession, error) {
	var session model.UserSession
	err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Update 更新会话
func (r *sessionRepository) Update(ctx context.Context, session *model.UserSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// Delete 删除会话
func (r *sessionRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.UserSession{}, id).Error
}

// GetByUserID 获取用户的会话列表
func (r *sessionRepository) GetByUserID(ctx context.Context, userID uint64, status string) ([]*model.UserSession, error) {
	var sessions []*model.UserSession
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("last_active_at DESC").Find(&sessions).Error
	return sessions, err
}

// GetActiveSessionsByUserID 获取用户活跃会话
func (r *sessionRepository) GetActiveSessionsByUserID(ctx context.Context, userID uint64) ([]*model.UserSession, error) {
	return r.GetByUserID(ctx, userID, "active")
}

// DeleteExpiredSessions 删除过期会话
func (r *sessionRepository) DeleteExpiredSessions(ctx context.Context, expiredBefore time.Time) error {
	return r.db.WithContext(ctx).
		Where("refresh_expires_at < ?", expiredBefore).
		Delete(&model.UserSession{}).Error
}

// UpdateStatus 更新会话状态
func (r *sessionRepository) UpdateStatus(ctx context.Context, sessionID string, status string, logoutTime time.Time) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if !logoutTime.IsZero() {
		updates["logout_time"] = logoutTime
	}
	return r.db.WithContext(ctx).
		Model(&model.UserSession{}).
		Where("session_id = ?", sessionID).
		Updates(updates).Error
}

// UpdateLastActive 更新最后活跃时间
func (r *sessionRepository) UpdateLastActive(ctx context.Context, sessionID string, lastActiveAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.UserSession{}).
		Where("session_id = ?", sessionID).
		Updates(map[string]interface{}{
			"last_active_at": lastActiveAt,
			"updated_at":     time.Now(),
		}).Error
}

// UpdateTokens 更新Token信息
func (r *sessionRepository) UpdateTokens(ctx context.Context, sessionID string, accessTokenHash, refreshTokenHash string, accessExpiresAt, refreshExpiresAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.UserSession{}).
		Where("session_id = ?", sessionID).
		Updates(map[string]interface{}{
			"access_token_hash":  accessTokenHash,
			"refresh_token_hash": refreshTokenHash,
			"token_expires_at":   accessExpiresAt,
			"refresh_expires_at": refreshExpiresAt,
			"updated_at":         time.Now(),
		}).Error
}

// CountActiveSessionsByUserID 统计用户活跃会话数
func (r *sessionRepository) CountActiveSessionsByUserID(ctx context.Context, userID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.UserSession{}).
		Where("user_id = ? AND status = ?", userID, "active").
		Count(&count).Error
	return count, err
}

// GetRecentLoginsByUserID 获取用户最近登录记录
func (r *sessionRepository) GetRecentLoginsByUserID(ctx context.Context, userID uint64, limit int) ([]*model.UserSession, error) {
	var sessions []*model.UserSession
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("login_time DESC").
		Limit(limit).
		Find(&sessions).Error
	return sessions, err
}

// GetSessionsByIPAddress 根据IP地址获取会话
func (r *sessionRepository) GetSessionsByIPAddress(ctx context.Context, ipAddress string, limit int) ([]*model.UserSession, error) {
	var sessions []*model.UserSession
	err := r.db.WithContext(ctx).
		Where("ip_address = ?", ipAddress).
		Order("login_time DESC").
		Limit(limit).
		Find(&sessions).Error
	return sessions, err
}

// GetSessionsByDeviceID 根据设备ID获取会话
func (r *sessionRepository) GetSessionsByDeviceID(ctx context.Context, deviceID string) ([]*model.UserSession, error) {
	var sessions []*model.UserSession
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("login_time DESC").
		Find(&sessions).Error
	return sessions, err
}
