package repository

import (
	"context"
	"fmt"

	"huinong-backend/internal/model"

	"gorm.io/gorm"
)

// systemConfigRepository 系统配置数据访问层实现
type systemConfigRepository struct {
	db *gorm.DB
}

// NewSystemConfigRepository 创建系统配置数据访问层实例
func NewSystemConfigRepository(db *gorm.DB) SystemConfigRepository {
	return &systemConfigRepository{db: db}
}

// Get 根据键获取配置
func (r *systemConfigRepository) Get(ctx context.Context, key string) (*model.SystemConfig, error) {
	var config model.SystemConfig
	if err := r.db.WithContext(ctx).
		Where("config_key = ?", key).
		First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("配置项不存在: %s", key)
		}
		return nil, fmt.Errorf("获取配置失败: %w", err)
	}
	return &config, nil
}

// Set 设置配置
func (r *systemConfigRepository) Set(ctx context.Context, config *model.SystemConfig) error {
	// 使用Upsert操作，如果存在则更新，不存在则创建
	if err := r.db.WithContext(ctx).
		Where("config_key = ?", config.ConfigKey).
		Assign(config).
		FirstOrCreate(config).Error; err != nil {
		return fmt.Errorf("设置配置失败: %w", err)
	}
	return nil
}

// GetByGroup 根据分组获取配置列表
func (r *systemConfigRepository) GetByGroup(ctx context.Context, group string) ([]*model.SystemConfig, error) {
	var configs []*model.SystemConfig
	if err := r.db.WithContext(ctx).
		Where("config_group = ?", group).
		Order("config_key ASC").
		Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("获取配置组失败: %w", err)
	}
	return configs, nil
}

// List 获取所有配置
func (r *systemConfigRepository) List(ctx context.Context) ([]*model.SystemConfig, error) {
	var configs []*model.SystemConfig
	if err := r.db.WithContext(ctx).
		Order("config_group ASC, config_key ASC").
		Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("获取配置列表失败: %w", err)
	}
	return configs, nil
}

// Delete 删除配置
func (r *systemConfigRepository) Delete(ctx context.Context, key string) error {
	result := r.db.WithContext(ctx).
		Where("config_key = ?", key).
		Delete(&model.SystemConfig{})

	if result.Error != nil {
		return fmt.Errorf("删除配置失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("配置项不存在: %s", key)
	}

	return nil
}
