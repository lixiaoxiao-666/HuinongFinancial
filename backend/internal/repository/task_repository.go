package repository

import (
	"context"
	"fmt"
	"huinong-backend/internal/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

// taskRepositoryImpl 任务仓库实现
type taskRepositoryImpl struct {
	db *gorm.DB
}

// NewTaskRepository 创建任务仓库实例
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepositoryImpl{db: db}
}

// CreateTask 创建任务
func (r *taskRepositoryImpl) CreateTask(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

// GetTaskByID 根据ID获取任务
func (r *taskRepositoryImpl) GetTaskByID(ctx context.Context, id uint64) (*model.Task, error) {
	var task model.Task
	err := r.db.WithContext(ctx).
		Preload("AssignedUser").
		Preload("CreatedUser").
		First(&task, id).Error

	if err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateTask 更新任务
func (r *taskRepositoryImpl) UpdateTask(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.Task{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteTask 删除任务
func (r *taskRepositoryImpl) DeleteTask(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Task{}, id).Error
}

// ListTasks 获取任务列表
func (r *taskRepositoryImpl) ListTasks(ctx context.Context, conditions map[string]interface{}, page, limit int, sortBy, sortOrder, keyword string) ([]*model.Task, int64, error) {
	var tasks []*model.Task
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Task{})

	// 添加查询条件
	for key, value := range conditions {
		if value != nil && value != "" {
			switch key {
			case "type", "status", "priority", "business_type":
				query = query.Where(fmt.Sprintf("%s = ?", key), value)
			case "assigned_to", "created_by", "business_id":
				query = query.Where(fmt.Sprintf("%s = ?", key), value)
			}
		}
	}

	// 关键词搜索
	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	orderBy := "created_at desc"
	if sortBy != "" {
		if sortOrder == "" {
			sortOrder = "desc"
		}
		orderBy = fmt.Sprintf("%s %s", sortBy, sortOrder)
	}

	// 分页查询
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).
		Order(orderBy).
		Preload("AssignedUser").
		Preload("CreatedUser").
		Find(&tasks).Error

	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// CreateTaskAction 创建任务操作记录
func (r *taskRepositoryImpl) CreateTaskAction(ctx context.Context, action *model.TaskAction) error {
	return r.db.WithContext(ctx).Create(action).Error
}

// GetTaskActions 获取任务操作记录
func (r *taskRepositoryImpl) GetTaskActions(ctx context.Context, taskID uint64) ([]*model.TaskAction, error) {
	var actions []*model.TaskAction
	err := r.db.WithContext(ctx).
		Where("task_id = ?", taskID).
		Preload("Operator").
		Order("created_at desc").
		Find(&actions).Error

	if err != nil {
		return nil, err
	}
	return actions, nil
}

// GetTaskActionsByUser 获取用户的任务操作记录
func (r *taskRepositoryImpl) GetTaskActionsByUser(ctx context.Context, userID uint64, limit int) ([]*model.TaskAction, error) {
	var actions []*model.TaskAction
	err := r.db.WithContext(ctx).
		Where("operator_id = ?", userID).
		Preload("Task").
		Preload("Operator").
		Order("created_at desc").
		Limit(limit).
		Find(&actions).Error

	if err != nil {
		return nil, err
	}
	return actions, nil
}

// GetTaskStatistics 获取任务统计
func (r *taskRepositoryImpl) GetTaskStatistics(ctx context.Context, conditions map[string]interface{}) (*model.TaskStatistics, error) {
	statistics := &model.TaskStatistics{}

	// 总任务数
	query := r.db.WithContext(ctx).Model(&model.Task{})
	for key, value := range conditions {
		if value != nil && value != "" {
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	if err := query.Count(&statistics.TotalTasks).Error; err != nil {
		return nil, err
	}

	// 按状态统计
	statusCounts := []struct {
		Status string
		Count  int64
	}{}

	baseQuery := r.db.WithContext(ctx).Model(&model.Task{})
	for key, value := range conditions {
		if key != "status" && value != nil && value != "" {
			baseQuery = baseQuery.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	if err := baseQuery.Select("status, COUNT(*) as count").Group("status").Scan(&statusCounts).Error; err != nil {
		return nil, err
	}

	for _, sc := range statusCounts {
		switch sc.Status {
		case model.TaskStatusPending:
			statistics.PendingTasks = sc.Count
		case model.TaskStatusProcessing:
			statistics.ProcessingTasks = sc.Count
		case model.TaskStatusCompleted:
			statistics.CompletedTasks = sc.Count
		case model.TaskStatusCancelled:
			statistics.CancelledTasks = sc.Count
		}
	}

	// 按优先级统计
	if err := baseQuery.Where("priority IN ?", []string{model.TaskPriorityHigh, model.TaskPriorityUrgent}).Count(&statistics.HighPriorityTasks).Error; err != nil {
		return nil, err
	}

	if err := baseQuery.Where("priority = ?", model.TaskPriorityUrgent).Count(&statistics.UrgentTasks).Error; err != nil {
		return nil, err
	}

	// 未分配任务
	if err := baseQuery.Where("assigned_to IS NULL").Count(&statistics.UnassignedTasks).Error; err != nil {
		return nil, err
	}

	return statistics, nil
}

// GetTaskCountByStatus 根据状态获取任务数量
func (r *taskRepositoryImpl) GetTaskCountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Task{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// GetTaskCountByType 根据类型获取任务数量
func (r *taskRepositoryImpl) GetTaskCountByType(ctx context.Context, taskType string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Task{}).Where("type = ?", taskType).Count(&count).Error
	return count, err
}

// GetTaskCountByPriority 根据优先级获取任务数量
func (r *taskRepositoryImpl) GetTaskCountByPriority(ctx context.Context, priority string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Task{}).Where("priority = ?", priority).Count(&count).Error
	return count, err
}

// GetOverdueTasks 获取超时任务
func (r *taskRepositoryImpl) GetOverdueTasks(ctx context.Context, limit int) ([]*model.Task, error) {
	var tasks []*model.Task

	// 计算超时时间点
	now := time.Now()

	// 紧急任务2小时，高优先级8小时，中优先级24小时，低优先级72小时
	overdueConditions := []string{
		fmt.Sprintf("(priority = '%s' AND status IN ('%s', '%s') AND created_at < '%s')",
			model.TaskPriorityUrgent, model.TaskStatusPending, model.TaskStatusProcessing,
			now.Add(-2*time.Hour).Format("2006-01-02 15:04:05")),
		fmt.Sprintf("(priority = '%s' AND status IN ('%s', '%s') AND created_at < '%s')",
			model.TaskPriorityHigh, model.TaskStatusPending, model.TaskStatusProcessing,
			now.Add(-8*time.Hour).Format("2006-01-02 15:04:05")),
		fmt.Sprintf("(priority = '%s' AND status IN ('%s', '%s') AND created_at < '%s')",
			model.TaskPriorityMedium, model.TaskStatusPending, model.TaskStatusProcessing,
			now.Add(-24*time.Hour).Format("2006-01-02 15:04:05")),
		fmt.Sprintf("(priority = '%s' AND status IN ('%s', '%s') AND created_at < '%s')",
			model.TaskPriorityLow, model.TaskStatusPending, model.TaskStatusProcessing,
			now.Add(-72*time.Hour).Format("2006-01-02 15:04:05")),
	}

	whereClause := strings.Join(overdueConditions, " OR ")

	err := r.db.WithContext(ctx).
		Where(whereClause).
		Preload("AssignedUser").
		Preload("CreatedUser").
		Order("priority desc, created_at asc").
		Limit(limit).
		Find(&tasks).Error

	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTasksByAssignee 根据分配人获取任务
func (r *taskRepositoryImpl) GetTasksByAssignee(ctx context.Context, assigneeID uint64, page, limit int) ([]*model.Task, int64, error) {
	var tasks []*model.Task
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Task{}).Where("assigned_to = ?", assigneeID)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).
		Order("created_at desc").
		Preload("AssignedUser").
		Preload("CreatedUser").
		Find(&tasks).Error

	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// GetUnassignedTasks 获取未分配任务
func (r *taskRepositoryImpl) GetUnassignedTasks(ctx context.Context, page, limit int) ([]*model.Task, int64, error) {
	var tasks []*model.Task
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Task{}).Where("assigned_to IS NULL")

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).
		Order("priority desc, created_at asc").
		Preload("CreatedUser").
		Find(&tasks).Error

	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// GetTasksByCreator 根据创建人获取任务
func (r *taskRepositoryImpl) GetTasksByCreator(ctx context.Context, creatorID uint64, page, limit int) ([]*model.Task, int64, error) {
	var tasks []*model.Task
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Task{}).Where("created_by = ?", creatorID)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).
		Order("created_at desc").
		Preload("AssignedUser").
		Preload("CreatedUser").
		Find(&tasks).Error

	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// GetTasksByBusiness 根据业务获取任务
func (r *taskRepositoryImpl) GetTasksByBusiness(ctx context.Context, businessType string, businessID uint64) ([]*model.Task, error) {
	var tasks []*model.Task

	err := r.db.WithContext(ctx).
		Where("business_type = ? AND business_id = ?", businessType, businessID).
		Preload("AssignedUser").
		Preload("CreatedUser").
		Order("created_at desc").
		Find(&tasks).Error

	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// CreateTaskFromBusiness 从业务创建任务
func (r *taskRepositoryImpl) CreateTaskFromBusiness(ctx context.Context, businessType string, businessID uint64, taskType string, title, description string, priority string, assigneeID *uint64) (*model.Task, error) {
	task := &model.Task{
		Title:        title,
		Description:  description,
		Type:         taskType,
		Priority:     priority,
		Status:       model.TaskStatusPending,
		BusinessID:   businessID,
		BusinessType: businessType,
		CreatedBy:    1, // TODO: 从上下文获取
		AssignedTo:   assigneeID,
	}

	if err := r.db.WithContext(ctx).Create(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// ==================== 辅助方法 ====================

// buildWhereClause 构建查询条件
func (r *taskRepositoryImpl) buildWhereClause(conditions map[string]interface{}) (string, []interface{}) {
	var whereParts []string
	var args []interface{}

	for key, value := range conditions {
		if value != nil && value != "" {
			whereParts = append(whereParts, fmt.Sprintf("%s = ?", key))
			args = append(args, value)
		}
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = strings.Join(whereParts, " AND ")
	}

	return whereClause, args
}

// validateSortField 验证排序字段
func (r *taskRepositoryImpl) validateSortField(sortBy string) string {
	validFields := map[string]bool{
		"id":           true,
		"title":        true,
		"type":         true,
		"priority":     true,
		"status":       true,
		"created_at":   true,
		"updated_at":   true,
		"completed_at": true,
	}

	if validFields[sortBy] {
		return sortBy
	}
	return "created_at" // 默认排序字段
}
