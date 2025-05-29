package repository

import (
	"context"

	"huinong-backend/internal/model"

	"gorm.io/gorm"
)

// offlineQueueRepository 离线队列Repository实现
type offlineQueueRepository struct {
	db *gorm.DB
}

// NewOfflineQueueRepository 创建离线队列Repository实例
func NewOfflineQueueRepository(db *gorm.DB) OfflineQueueRepository {
	return &offlineQueueRepository{
		db: db,
	}
}

// Add 添加离线操作
func (r *offlineQueueRepository) Add(ctx context.Context, action *model.OfflineQueue) error {
	return r.db.WithContext(ctx).Create(action).Error
}

// GetPending 获取待处理的操作
func (r *offlineQueueRepository) GetPending(ctx context.Context, limit int) ([]*model.OfflineQueue, error) {
	var actions []*model.OfflineQueue
	err := r.db.WithContext(ctx).
		Where("status = ?", "pending").
		Order("created_at ASC").
		Limit(limit).
		Find(&actions).Error
	return actions, err
}

// GetUserActions 获取用户操作列表
func (r *offlineQueueRepository) GetUserActions(ctx context.Context, userID uint64, limit, offset int) ([]*model.OfflineQueue, error) {
	var actions []*model.OfflineQueue
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&actions).Error
	return actions, err
}

// Update 更新操作状态
func (r *offlineQueueRepository) Update(ctx context.Context, action *model.OfflineQueue) error {
	return r.db.WithContext(ctx).Save(action).Error
}

// Delete 删除操作记录
func (r *offlineQueueRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.OfflineQueue{}, id).Error
}

// GetRetryableActions 获取可重试的操作
func (r *offlineQueueRepository) GetRetryableActions(ctx context.Context) ([]*model.OfflineQueue, error) {
	var actions []*model.OfflineQueue
	err := r.db.WithContext(ctx).
		Where("status = ? AND retry_count < max_retry", "failed").
		Order("created_at ASC").
		Find(&actions).Error
	return actions, err
}
