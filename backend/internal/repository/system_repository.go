package repository

import (
	"context"
	"fmt"
	"time"

	"huinong-backend/internal/model"

	"gorm.io/gorm"
)

// systemConfigRepository 系统配置Repository实现
type systemConfigRepository struct {
	db *gorm.DB
}

// NewSystemConfigRepository 创建系统配置Repository实例
func NewSystemConfigRepository(db *gorm.DB) SystemConfigRepository {
	return &systemConfigRepository{
		db: db,
	}
}

// Get 根据键获取配置
func (r *systemConfigRepository) Get(ctx context.Context, key string) (*model.SystemConfig, error) {
	var config model.SystemConfig
	err := r.db.WithContext(ctx).Where("config_key = ?", key).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("配置不存在")
		}
		return nil, err
	}
	return &config, nil
}

// Set 设置配置
func (r *systemConfigRepository) Set(ctx context.Context, config *model.SystemConfig) error {
	var existing model.SystemConfig
	err := r.db.WithContext(ctx).Where("config_key = ?", config.ConfigKey).First(&existing).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if err == gorm.ErrRecordNotFound {
		// 创建新配置
		return r.db.WithContext(ctx).Create(config).Error
	} else {
		// 更新现有配置
		config.ID = existing.ID
		return r.db.WithContext(ctx).Save(config).Error
	}
}

// GetByGroup 根据分组获取配置
func (r *systemConfigRepository) GetByGroup(ctx context.Context, group string) ([]*model.SystemConfig, error) {
	var configs []*model.SystemConfig
	err := r.db.WithContext(ctx).
		Where("config_group = ?", group).
		Order("sort_order ASC, created_at ASC").
		Find(&configs).Error
	return configs, err
}

// List 获取所有配置
func (r *systemConfigRepository) List(ctx context.Context) ([]*model.SystemConfig, error) {
	var configs []*model.SystemConfig
	err := r.db.WithContext(ctx).
		Order("config_group ASC, sort_order ASC, created_at ASC").
		Find(&configs).Error
	return configs, err
}

// Delete 删除配置
func (r *systemConfigRepository) Delete(ctx context.Context, key string) error {
	return r.db.WithContext(ctx).Where("config_key = ?", key).Delete(&model.SystemConfig{}).Error
}

// ==================== 文件管理Repository ====================

// fileRepository 文件Repository实现
type fileRepository struct {
	db *gorm.DB
}

// NewFileRepository 创建文件Repository实例
func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{
		db: db,
	}
}

// Create 创建文件记录
func (r *fileRepository) Create(ctx context.Context, file *model.FileUpload) error {
	return r.db.WithContext(ctx).Create(file).Error
}

// GetByID 根据ID获取文件
func (r *fileRepository) GetByID(ctx context.Context, id uint64) (*model.FileUpload, error) {
	var file model.FileUpload
	err := r.db.WithContext(ctx).First(&file, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("文件不存在")
		}
		return nil, err
	}
	return &file, nil
}

// GetByHash 根据文件哈希获取文件
func (r *fileRepository) GetByHash(ctx context.Context, hash string) (*model.FileUpload, error) {
	var file model.FileUpload
	err := r.db.WithContext(ctx).Where("file_hash = ?", hash).First(&file).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("文件不存在")
		}
		return nil, err
	}
	return &file, nil
}

// GetByBusiness 根据业务类型和业务ID获取文件列表
func (r *fileRepository) GetByBusiness(ctx context.Context, businessType string, businessID uint64) ([]*model.FileUpload, error) {
	var files []*model.FileUpload
	err := r.db.WithContext(ctx).
		Where("business_type = ? AND business_id = ?", businessType, businessID).
		Order("created_at DESC").
		Find(&files).Error
	return files, err
}

// Update 更新文件记录
func (r *fileRepository) Update(ctx context.Context, file *model.FileUpload) error {
	return r.db.WithContext(ctx).Save(file).Error
}

// Delete 删除文件记录
func (r *fileRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.FileUpload{}, id).Error
}

// IncrementAccessCount 增加访问次数
func (r *fileRepository) IncrementAccessCount(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.FileUpload{}).
		Where("id = ?", id).
		UpdateColumn("access_count", gorm.Expr("access_count + 1")).Error
}

// ==================== API日志Repository ====================

// apiLogRepository API日志Repository实现
type apiLogRepository struct {
	db *gorm.DB
}

// NewAPILogRepository 创建API日志Repository实例
func NewAPILogRepository(db *gorm.DB) APILogRepository {
	return &apiLogRepository{
		db: db,
	}
}

// Create 创建API日志
func (r *apiLogRepository) Create(ctx context.Context, log *model.APILog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetByRequestID 根据请求ID获取日志
func (r *apiLogRepository) GetByRequestID(ctx context.Context, requestID string) (*model.APILog, error) {
	var log model.APILog
	err := r.db.WithContext(ctx).Where("request_id = ?", requestID).First(&log).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("日志不存在")
		}
		return nil, err
	}
	return &log, nil
}

// List API日志列表查询
func (r *apiLogRepository) List(ctx context.Context, req *ListAPILogsRequest) (*ListAPILogsResponse, error) {
	var logs []*model.APILog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.APILog{})

	// 条件筛选
	if req.Method != "" {
		query = query.Where("method = ?", req.Method)
	}
	if req.URL != "" {
		query = query.Where("url LIKE ?", "%"+req.URL+"%")
	}
	if req.StatusCode > 0 {
		query = query.Where("status_code = ?", req.StatusCode)
	}
	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.StartDate != "" && req.EndDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)
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
		Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return &ListAPILogsResponse{
		Logs:  logs,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

// DeleteOldLogs 删除旧日志
func (r *apiLogRepository) DeleteOldLogs(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).Where("created_at < ?", before).Delete(&model.APILog{}).Error
}

// GetStatistics 获取API统计
func (r *apiLogRepository) GetStatistics(ctx context.Context, startTime, endTime time.Time) (*APIStatistics, error) {
	stats := &APIStatistics{}

	// 总请求数
	err := r.db.WithContext(ctx).Model(&model.APILog{}).
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Count(&stats.TotalRequests).Error
	if err != nil {
		return nil, err
	}

	// 成功请求数
	err = r.db.WithContext(ctx).Model(&model.APILog{}).
		Where("created_at BETWEEN ? AND ? AND status_code >= 200 AND status_code < 300", startTime, endTime).
		Count(&stats.SuccessRequests).Error
	if err != nil {
		return nil, err
	}

	// 错误请求数
	stats.ErrorRequests = stats.TotalRequests - stats.SuccessRequests

	// 平均响应时间
	var avgTime *float64
	err = r.db.WithContext(ctx).Model(&model.APILog{}).
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Select("AVG(response_time)").
		Row().Scan(&avgTime)
	if err != nil {
		return nil, err
	}
	if avgTime != nil {
		stats.AvgResponseTime = *avgTime
	}

	// 热门端点
	var endpoints []EndpointStat
	err = r.db.WithContext(ctx).Model(&model.APILog{}).
		Select("url as endpoint, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Group("url").
		Order("count DESC").
		Limit(10).
		Find(&endpoints).Error
	if err != nil {
		return nil, err
	}
	stats.TopEndpoints = endpoints

	return stats, nil
}
