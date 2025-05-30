package service

import (
	"context"
	"errors"
	"fmt"
	"huinong-backend/internal/model"
	"huinong-backend/internal/repository"
	"strings"
	"time"

	"gorm.io/gorm"
)

// taskServiceImpl 任务服务实现
type taskServiceImpl struct {
	taskRepo repository.TaskRepository
	userRepo repository.UserRepository
	loanRepo repository.LoanRepository
	db       *gorm.DB
}

// NewTaskService 创建任务服务实例
func NewTaskService(
	taskRepo repository.TaskRepository,
	userRepo repository.UserRepository,
	loanRepo repository.LoanRepository,
	db *gorm.DB,
) TaskService {
	return &taskServiceImpl{
		taskRepo: taskRepo,
		userRepo: userRepo,
		loanRepo: loanRepo,
		db:       db,
	}
}

// CreateTask 创建任务
func (s *taskServiceImpl) CreateTask(ctx context.Context, req *CreateTaskRequest) (*CreateTaskResponse, error) {
	if req.Title == "" {
		return nil, errors.New("任务标题不能为空")
	}

	if req.Type == "" {
		return nil, errors.New("任务类型不能为空")
	}

	// 设置默认优先级
	if req.Priority == "" {
		req.Priority = model.TaskPriorityMedium
	}

	// 验证任务类型
	if !isValidTaskType(req.Type) {
		return nil, errors.New("无效的任务类型")
	}

	// 验证优先级
	if !isValidTaskPriority(req.Priority) {
		return nil, errors.New("无效的任务优先级")
	}

	// 从上下文获取创建人ID
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("无法获取用户信息")
	}

	// 创建任务模型
	task := &model.Task{
		Title:        req.Title,
		Description:  req.Description,
		Type:         req.Type,
		Priority:     req.Priority,
		Status:       model.TaskStatusPending,
		CreatedBy:    userID,
		BusinessID:   req.BusinessID,
		BusinessType: req.BusinessType,
		Data:         req.Data,
		AssignedTo:   req.AssignedTo,
	}

	// 保存任务
	if err := s.taskRepo.CreateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("创建任务失败: %v", err)
	}

	// 记录任务创建操作
	if err := s.recordTaskAction(ctx, task.ID, model.TaskActionCreate, "任务已创建", userID); err != nil {
		// 记录操作失败不影响任务创建
		fmt.Printf("记录任务操作失败: %v\n", err)
	}

	// 如果指定了分配人，记录分配操作
	if req.AssignedTo != nil {
		if err := s.recordTaskAction(ctx, task.ID, model.TaskActionAssign,
			fmt.Sprintf("任务已分配给用户ID: %d", *req.AssignedTo), userID); err != nil {
			fmt.Printf("记录任务分配操作失败: %v\n", err)
		}
	}

	// 生成任务编号
	taskNo := s.generateTaskNo(task.Type, task.ID)

	return &CreateTaskResponse{
		ID:     task.ID,
		TaskNo: taskNo,
		Status: task.Status,
	}, nil
}

// GetTask 获取任务详情
func (s *taskServiceImpl) GetTask(ctx context.Context, taskID uint64) (*TaskDetailResponse, error) {
	task, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("任务不存在")
		}
		return nil, fmt.Errorf("获取任务失败: %v", err)
	}

	// 获取任务操作记录
	actions, err := s.taskRepo.GetTaskActions(ctx, taskID)
	if err != nil {
		fmt.Printf("获取任务操作记录失败: %v\n", err)
	}

	// 转换为响应格式
	response := s.convertTaskToResponse(task, actions)

	// 获取业务相关信息
	businessInfo, err := s.getBusinessInfo(ctx, task.BusinessType, task.BusinessID)
	if err != nil {
		fmt.Printf("获取业务信息失败: %v\n", err)
	} else {
		response.BusinessInfo = businessInfo
	}

	return response, nil
}

// UpdateTask 更新任务
func (s *taskServiceImpl) UpdateTask(ctx context.Context, taskID uint64, req *UpdateTaskRequest) error {
	task, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("任务不存在")
		}
		return fmt.Errorf("获取任务失败: %v", err)
	}

	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		return errors.New("无法获取用户信息")
	}

	// 更新字段
	updates := make(map[string]interface{})
	var actions []string

	if req.Title != "" && req.Title != task.Title {
		updates["title"] = req.Title
		actions = append(actions, fmt.Sprintf("标题从 '%s' 修改为 '%s'", task.Title, req.Title))
	}

	if req.Description != task.Description {
		updates["description"] = req.Description
		actions = append(actions, "描述已更新")
	}

	if req.Priority != "" && isValidTaskPriority(req.Priority) && req.Priority != task.Priority {
		updates["priority"] = req.Priority
		actions = append(actions, fmt.Sprintf("优先级从 '%s' 修改为 '%s'",
			getTaskPriorityName(task.Priority), getTaskPriorityName(req.Priority)))
	}

	if req.Status != "" && isValidTaskStatus(req.Status) && req.Status != task.Status {
		updates["status"] = req.Status
		actions = append(actions, fmt.Sprintf("状态从 '%s' 修改为 '%s'",
			getTaskStatusName(task.Status), getTaskStatusName(req.Status)))

		// 如果状态变为完成，设置完成时间
		if req.Status == model.TaskStatusCompleted {
			now := time.Now()
			updates["completed_at"] = &now
		}
	}

	if req.AssignedTo != nil {
		if task.AssignedTo == nil || *req.AssignedTo != *task.AssignedTo {
			updates["assigned_to"] = req.AssignedTo
			if *req.AssignedTo == 0 {
				actions = append(actions, "任务已取消分配")
			} else {
				actions = append(actions, fmt.Sprintf("任务已重新分配给用户ID: %d", *req.AssignedTo))
			}
		}
	}

	if req.Data != "" && req.Data != task.Data {
		updates["data"] = req.Data
		actions = append(actions, "任务数据已更新")
	}

	// 如果有更新内容，执行更新
	if len(updates) > 0 {
		if err := s.taskRepo.UpdateTask(ctx, taskID, updates); err != nil {
			return fmt.Errorf("更新任务失败: %v", err)
		}

		// 记录更新操作
		actionDesc := strings.Join(actions, "；")
		if err := s.recordTaskAction(ctx, taskID, "update", actionDesc, userID); err != nil {
			fmt.Printf("记录任务更新操作失败: %v\n", err)
		}
	}

	return nil
}

// DeleteTask 删除任务
func (s *taskServiceImpl) DeleteTask(ctx context.Context, taskID uint64) error {
	task, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("任务不存在")
		}
		return fmt.Errorf("获取任务失败: %v", err)
	}

	// 只能删除待处理或已取消的任务
	if task.Status != model.TaskStatusPending && task.Status != model.TaskStatusCancelled {
		return errors.New("只能删除待处理或已取消的任务")
	}

	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		return errors.New("无法获取用户信息")
	}

	// 记录删除操作
	if err := s.recordTaskAction(ctx, taskID, "delete", "任务已删除", userID); err != nil {
		fmt.Printf("记录任务删除操作失败: %v\n", err)
	}

	// 删除任务
	if err := s.taskRepo.DeleteTask(ctx, taskID); err != nil {
		return fmt.Errorf("删除任务失败: %v", err)
	}

	return nil
}

// ListTasks 获取任务列表
func (s *taskServiceImpl) ListTasks(ctx context.Context, req *ListTasksRequest) (*ListTasksResponse, error) {
	// 设置默认分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// 设置默认排序
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	// 构建查询条件
	conditions := make(map[string]interface{})

	if req.Type != "" {
		conditions["type"] = req.Type
	}
	if req.Status != "" {
		conditions["status"] = req.Status
	}
	if req.Priority != "" {
		conditions["priority"] = req.Priority
	}
	if req.AssignedTo != nil {
		conditions["assigned_to"] = *req.AssignedTo
	}
	if req.CreatedBy != nil {
		conditions["created_by"] = *req.CreatedBy
	}
	if req.BusinessType != "" {
		conditions["business_type"] = req.BusinessType
	}
	if req.BusinessID != nil {
		conditions["business_id"] = *req.BusinessID
	}

	// 获取任务列表
	tasks, total, err := s.taskRepo.ListTasks(ctx, conditions, req.Page, req.Limit, req.SortBy, req.SortOrder, req.Keyword)
	if err != nil {
		return nil, fmt.Errorf("获取任务列表失败: %v", err)
	}

	// 转换为响应格式
	taskResponses := make([]*TaskDetailResponse, 0, len(tasks))
	for _, task := range tasks {
		taskResponse := s.convertTaskToResponse(task, nil) // 列表不包含操作记录

		// 处理超时任务过滤
		if req.IsOverdue != nil {
			if *req.IsOverdue != taskResponse.IsOverdue {
				continue
			}
		}

		taskResponses = append(taskResponses, taskResponse)
	}

	// 获取统计信息
	statistics, err := s.getTaskStatistics(ctx, conditions)
	if err != nil {
		fmt.Printf("获取任务统计失败: %v\n", err)
	}

	return &ListTasksResponse{
		Tasks:      taskResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		Statistics: statistics,
	}, nil
}

// AssignTask 分配任务
func (s *taskServiceImpl) AssignTask(ctx context.Context, taskID uint64, assigneeID uint64) error {
	task, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("任务不存在")
		}
		return fmt.Errorf("获取任务失败: %v", err)
	}

	// 验证任务状态
	if task.Status != model.TaskStatusPending {
		return errors.New("只能分配待处理的任务")
	}

	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		return errors.New("无法获取用户信息")
	}

	// 验证分配目标用户是否存在
	if _, err := s.userRepo.GetByID(ctx, assigneeID); err != nil {
		return errors.New("分配目标用户不存在")
	}

	// 更新任务分配
	updates := map[string]interface{}{
		"assigned_to": assigneeID,
	}

	if err := s.taskRepo.UpdateTask(ctx, taskID, updates); err != nil {
		return fmt.Errorf("分配任务失败: %v", err)
	}

	// 记录分配操作
	actionDesc := fmt.Sprintf("任务已分配给用户ID: %d", assigneeID)
	if err := s.recordTaskAction(ctx, taskID, model.TaskActionAssign, actionDesc, userID); err != nil {
		fmt.Printf("记录任务分配操作失败: %v\n", err)
	}

	return nil
}

// UnassignTask 取消分配任务
func (s *taskServiceImpl) UnassignTask(ctx context.Context, taskID uint64) error {
	task, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("任务不存在")
		}
		return fmt.Errorf("获取任务失败: %v", err)
	}

	if task.AssignedTo == nil {
		return errors.New("任务未分配，无需取消")
	}

	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		return errors.New("无法获取用户信息")
	}

	// 取消分配
	updates := map[string]interface{}{
		"assigned_to": nil,
	}

	if err := s.taskRepo.UpdateTask(ctx, taskID, updates); err != nil {
		return fmt.Errorf("取消分配失败: %v", err)
	}

	// 记录取消分配操作
	if err := s.recordTaskAction(ctx, taskID, "unassign", "任务分配已取消", userID); err != nil {
		fmt.Printf("记录任务取消分配操作失败: %v\n", err)
	}

	return nil
}

// GetTaskProgress 获取任务进度
func (s *taskServiceImpl) GetTaskProgress(ctx context.Context, taskID uint64) (*TaskProgressResponse, error) {
	task, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("任务不存在")
		}
		return nil, fmt.Errorf("获取任务失败: %v", err)
	}

	// 根据任务状态计算进度
	progress := calculateTaskProgress(task.Status)

	return &TaskProgressResponse{
		TaskID:    task.ID,
		Progress:  progress,
		UpdatedAt: task.UpdatedAt,
		UpdatedBy: nil, // TODO: 添加更新人信息
	}, nil
}

// UpdateTaskProgress 更新任务进度
func (s *taskServiceImpl) UpdateTaskProgress(ctx context.Context, taskID uint64, progress float64) error {
	if progress < 0 || progress > 100 {
		return errors.New("进度值必须在0-100之间")
	}

	_, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("任务不存在")
		}
		return fmt.Errorf("获取任务失败: %v", err)
	}

	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		return errors.New("无法获取用户信息")
	}

	// 根据进度更新状态
	var newStatus string
	if progress == 0 {
		newStatus = model.TaskStatusPending
	} else if progress == 100 {
		newStatus = model.TaskStatusCompleted
	} else {
		newStatus = model.TaskStatusProcessing
	}

	updates := map[string]interface{}{
		"status": newStatus,
	}

	// 如果任务完成，设置完成时间
	if newStatus == model.TaskStatusCompleted {
		now := time.Now()
		updates["completed_at"] = &now
	}

	if err := s.taskRepo.UpdateTask(ctx, taskID, updates); err != nil {
		return fmt.Errorf("更新任务进度失败: %v", err)
	}

	// 记录进度更新操作
	actionDesc := fmt.Sprintf("任务进度更新为 %.1f%%", progress)
	if err := s.recordTaskAction(ctx, taskID, "progress", actionDesc, userID); err != nil {
		fmt.Printf("记录任务进度操作失败: %v\n", err)
	}

	return nil
}

// CompleteTask 完成任务
func (s *taskServiceImpl) CompleteTask(ctx context.Context, taskID uint64) error {
	task, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("任务不存在")
		}
		return fmt.Errorf("获取任务失败: %v", err)
	}

	if task.Status == model.TaskStatusCompleted {
		return errors.New("任务已完成")
	}

	if task.Status == model.TaskStatusCancelled {
		return errors.New("已取消的任务无法完成")
	}

	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		return errors.New("无法获取用户信息")
	}

	// 更新任务状态
	now := time.Now()
	updates := map[string]interface{}{
		"status":       model.TaskStatusCompleted,
		"completed_at": &now,
	}

	if err := s.taskRepo.UpdateTask(ctx, taskID, updates); err != nil {
		return fmt.Errorf("完成任务失败: %v", err)
	}

	// 记录完成操作
	if err := s.recordTaskAction(ctx, taskID, model.TaskActionComplete, "任务已完成", userID); err != nil {
		fmt.Printf("记录任务完成操作失败: %v\n", err)
	}

	return nil
}

// ReassignTask 重新分配任务
func (s *taskServiceImpl) ReassignTask(ctx context.Context, taskID uint64, newAssigneeID uint64) error {
	task, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("任务不存在")
		}
		return fmt.Errorf("获取任务失败: %v", err)
	}

	if task.Status == model.TaskStatusCompleted {
		return errors.New("已完成的任务无法重新分配")
	}

	if task.Status == model.TaskStatusCancelled {
		return errors.New("已取消的任务无法重新分配")
	}

	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		return errors.New("无法获取用户信息")
	}

	// 验证新的分配目标用户是否存在
	if _, err := s.userRepo.GetByID(ctx, newAssigneeID); err != nil {
		return errors.New("分配目标用户不存在")
	}

	// 更新任务分配
	updates := map[string]interface{}{
		"assigned_to": newAssigneeID,
	}

	if err := s.taskRepo.UpdateTask(ctx, taskID, updates); err != nil {
		return fmt.Errorf("重新分配任务失败: %v", err)
	}

	// 记录重新分配操作
	actionDesc := fmt.Sprintf("任务已重新分配给用户ID: %d", newAssigneeID)
	if err := s.recordTaskAction(ctx, taskID, model.TaskActionReassign, actionDesc, userID); err != nil {
		fmt.Printf("记录任务重新分配操作失败: %v\n", err)
	}

	return nil
}

// ==================== 辅助方法 ====================

// recordTaskAction 记录任务操作
func (s *taskServiceImpl) recordTaskAction(ctx context.Context, taskID uint64, action, comment string, operatorID uint64) error {
	taskAction := &model.TaskAction{
		TaskID:     taskID,
		Action:     action,
		Comment:    comment,
		OperatorID: operatorID,
	}

	return s.taskRepo.CreateTaskAction(ctx, taskAction)
}

// convertTaskToResponse 转换任务为响应格式
func (s *taskServiceImpl) convertTaskToResponse(task *model.Task, actions []*model.TaskAction) *TaskDetailResponse {
	response := &TaskDetailResponse{
		ID:           task.ID,
		Title:        task.Title,
		Description:  task.Description,
		Type:         task.Type,
		TypeName:     getTaskTypeName(task.Type),
		Priority:     task.Priority,
		PriorityName: getTaskPriorityName(task.Priority),
		Status:       task.Status,
		StatusName:   getTaskStatusName(task.Status),
		Progress:     calculateTaskProgress(task.Status),
		IsOverdue:    task.IsOverdue(),
		BusinessID:   task.BusinessID,
		BusinessType: task.BusinessType,
		Data:         task.Data,
		CreatedAt:    task.CreatedAt,
		UpdatedAt:    task.UpdatedAt,
		CompletedAt:  task.CompletedAt,
	}

	// 添加用户信息
	if task.AssignedUser != nil {
		response.AssignedUser = &TaskUserInfo{
			ID:       task.AssignedUser.ID,
			Username: task.AssignedUser.Username,
			RealName: task.AssignedUser.RealName,
		}
	}

	if task.CreatedUser != nil {
		response.CreatedUser = &TaskUserInfo{
			ID:       task.CreatedUser.ID,
			Username: task.CreatedUser.Username,
			RealName: task.CreatedUser.RealName,
		}
	}

	// 添加操作记录
	if actions != nil {
		response.Actions = make([]*TaskActionInfo, 0, len(actions))
		for _, action := range actions {
			actionInfo := &TaskActionInfo{
				ID:         action.ID,
				Action:     action.Action,
				ActionName: getTaskActionName(action.Action),
				Comment:    action.Comment,
				CreatedAt:  action.CreatedAt,
			}

			if action.Operator != nil {
				actionInfo.Operator = &TaskUserInfo{
					ID:       action.Operator.ID,
					Username: action.Operator.Username,
					RealName: action.Operator.RealName,
				}
			}

			response.Actions = append(response.Actions, actionInfo)
		}
	}

	return response
}

// getBusinessInfo 获取业务相关信息
func (s *taskServiceImpl) getBusinessInfo(ctx context.Context, businessType string, businessID uint64) (map[string]interface{}, error) {
	businessInfo := make(map[string]interface{})

	switch businessType {
	case "loan_application":
		// 获取贷款申请信息
		if application, err := s.loanRepo.GetApplicationByID(ctx, uint(businessID)); err == nil {
			businessInfo["application_no"] = application.ApplicationNo
			businessInfo["loan_amount"] = application.LoanAmount
			businessInfo["user_name"] = "" // TODO: 获取用户名
		}
	case "user_auth":
		// 获取用户认证信息
		// TODO: 实现用户认证信息获取
	default:
		businessInfo["type"] = businessType
		businessInfo["id"] = businessID
	}

	return businessInfo, nil
}

// getTaskStatistics 获取任务统计信息
func (s *taskServiceImpl) getTaskStatistics(ctx context.Context, conditions map[string]interface{}) (*TaskStatisticsInfo, error) {
	statistics, err := s.taskRepo.GetTaskStatistics(ctx, conditions)
	if err != nil {
		return nil, err
	}

	return &TaskStatisticsInfo{
		TotalTasks:        statistics.TotalTasks,
		PendingTasks:      statistics.PendingTasks,
		ProcessingTasks:   statistics.ProcessingTasks,
		CompletedTasks:    statistics.CompletedTasks,
		CancelledTasks:    statistics.CancelledTasks,
		OverdueTasks:      0, // TODO: 计算超时任务
		HighPriorityTasks: statistics.HighPriorityTasks,
		UrgentTasks:       statistics.UrgentTasks,
		MyTasks:           statistics.MyTasks,
		UnassignedTasks:   statistics.UnassignedTasks,
	}, nil
}

// generateTaskNo 生成任务编号
func (s *taskServiceImpl) generateTaskNo(taskType string, taskID uint64) string {
	prefix := ""
	switch taskType {
	case model.TaskTypeLoanApproval:
		prefix = "LA"
	case model.TaskTypeUserAuth:
		prefix = "UA"
	case model.TaskTypeMachineVerify:
		prefix = "MV"
	case model.TaskTypeContentReview:
		prefix = "CR"
	case model.TaskTypeSystemMaintenance:
		prefix = "SM"
	default:
		prefix = "TK"
	}

	now := time.Now()
	return fmt.Sprintf("%s%s%06d", prefix, now.Format("20060102"), taskID)
}

// getUserIDFromContext 从上下文获取用户ID
func getUserIDFromContext(ctx context.Context) uint64 {
	// TODO: 实现从上下文获取用户ID的逻辑
	// 这里暂时返回固定值
	return 1
}

// calculateTaskProgress 根据状态计算任务进度
func calculateTaskProgress(status string) float64 {
	switch status {
	case model.TaskStatusPending:
		return 0
	case model.TaskStatusProcessing:
		return 50
	case model.TaskStatusCompleted:
		return 100
	case model.TaskStatusCancelled:
		return 0
	default:
		return 0
	}
}

// 验证和名称获取函数
func isValidTaskType(taskType string) bool {
	validTypes := []string{
		model.TaskTypeLoanApproval,
		model.TaskTypeUserAuth,
		model.TaskTypeMachineVerify,
		model.TaskTypeContentReview,
		model.TaskTypeSystemMaintenance,
		model.TaskTypeCustom,
	}

	for _, validType := range validTypes {
		if taskType == validType {
			return true
		}
	}
	return false
}

func isValidTaskPriority(priority string) bool {
	validPriorities := []string{
		model.TaskPriorityLow,
		model.TaskPriorityMedium,
		model.TaskPriorityHigh,
		model.TaskPriorityUrgent,
	}

	for _, validPriority := range validPriorities {
		if priority == validPriority {
			return true
		}
	}
	return false
}

func isValidTaskStatus(status string) bool {
	validStatuses := []string{
		model.TaskStatusPending,
		model.TaskStatusProcessing,
		model.TaskStatusCompleted,
		model.TaskStatusCancelled,
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func getTaskTypeName(taskType string) string {
	names := map[string]string{
		model.TaskTypeLoanApproval:      "贷款审批",
		model.TaskTypeUserAuth:          "用户认证审核",
		model.TaskTypeMachineVerify:     "农机设备验证",
		model.TaskTypeContentReview:     "内容审核",
		model.TaskTypeSystemMaintenance: "系统维护",
		model.TaskTypeCustom:            "自定义任务",
	}

	if name, exists := names[taskType]; exists {
		return name
	}
	return taskType
}

func getTaskPriorityName(priority string) string {
	names := map[string]string{
		model.TaskPriorityLow:    "低",
		model.TaskPriorityMedium: "中",
		model.TaskPriorityHigh:   "高",
		model.TaskPriorityUrgent: "紧急",
	}

	if name, exists := names[priority]; exists {
		return name
	}
	return priority
}

func getTaskStatusName(status string) string {
	names := map[string]string{
		model.TaskStatusPending:    "待处理",
		model.TaskStatusProcessing: "进行中",
		model.TaskStatusCompleted:  "已完成",
		model.TaskStatusCancelled:  "已取消",
	}

	if name, exists := names[status]; exists {
		return name
	}
	return status
}

func getTaskActionName(action string) string {
	names := map[string]string{
		model.TaskActionCreate:   "创建",
		model.TaskActionAssign:   "分配",
		model.TaskActionStart:    "开始处理",
		model.TaskActionComplete: "完成",
		model.TaskActionCancel:   "取消",
		model.TaskActionComment:  "评论",
		model.TaskActionReassign: "重新分配",
		"update":                 "更新",
		"progress":               "进度更新",
		"unassign":               "取消分配",
		"delete":                 "删除",
	}

	if name, exists := names[action]; exists {
		return name
	}
	return action
}
