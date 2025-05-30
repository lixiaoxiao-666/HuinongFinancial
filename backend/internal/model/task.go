package model

import (
	"time"
)

// Task 待处理任务基础模型
type Task struct {
	ID           uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Title        string     `json:"title" gorm:"type:varchar(255);not null;comment:任务标题"`
	Description  string     `json:"description" gorm:"type:text;comment:任务描述"`
	Type         string     `json:"type" gorm:"type:varchar(50);not null;comment:任务类型"`
	Priority     string     `json:"priority" gorm:"type:varchar(20);default:medium;comment:任务优先级(low,medium,high,urgent)"`
	Status       string     `json:"status" gorm:"type:varchar(20);default:pending;comment:任务状态(pending,processing,completed,cancelled)"`
	AssignedTo   *uint64    `json:"assigned_to" gorm:"comment:分配给的用户ID"`
	CreatedBy    uint64     `json:"created_by" gorm:"not null;comment:创建人ID"`
	BusinessID   uint64     `json:"business_id" gorm:"not null;comment:关联业务ID"`
	BusinessType string     `json:"business_type" gorm:"type:varchar(50);not null;comment:业务类型"`
	Data         string     `json:"data" gorm:"type:json;comment:任务相关数据"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	CompletedAt  *time.Time `json:"completed_at" gorm:"comment:完成时间"`

	// 关联信息
	AssignedUser *OAUser `json:"assigned_user,omitempty" gorm:"foreignKey:AssignedTo"`
	CreatedUser  *OAUser `json:"created_user,omitempty" gorm:"foreignKey:CreatedBy"`
}

// TaskAction 任务操作记录
type TaskAction struct {
	ID         uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	TaskID     uint64    `json:"task_id" gorm:"not null;comment:任务ID"`
	Action     string    `json:"action" gorm:"type:varchar(50);not null;comment:操作类型"`
	Comment    string    `json:"comment" gorm:"type:text;comment:操作备注"`
	OperatorID uint64    `json:"operator_id" gorm:"not null;comment:操作人ID"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`

	// 关联信息
	Task     *Task   `json:"task,omitempty" gorm:"foreignKey:TaskID"`
	Operator *OAUser `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
}

// TaskStatistics 任务统计
type TaskStatistics struct {
	TotalTasks        int64 `json:"total_tasks"`
	PendingTasks      int64 `json:"pending_tasks"`
	ProcessingTasks   int64 `json:"processing_tasks"`
	CompletedTasks    int64 `json:"completed_tasks"`
	CancelledTasks    int64 `json:"cancelled_tasks"`
	HighPriorityTasks int64 `json:"high_priority_tasks"`
	UrgentTasks       int64 `json:"urgent_tasks"`
	MyTasks           int64 `json:"my_tasks"`
	UnassignedTasks   int64 `json:"unassigned_tasks"`
}

// TaskType 任务类型常量
const (
	TaskTypeLoanApproval      = "loan_approval"      // 贷款审批
	TaskTypeUserAuth          = "user_auth"          // 用户认证审核
	TaskTypeMachineVerify     = "machine_verify"     // 农机设备验证
	TaskTypeContentReview     = "content_review"     // 内容审核
	TaskTypeSystemMaintenance = "system_maintenance" // 系统维护
	TaskTypeCustom            = "custom"             // 自定义任务
)

// TaskPriority 任务优先级常量
const (
	TaskPriorityLow    = "low"
	TaskPriorityMedium = "medium"
	TaskPriorityHigh   = "high"
	TaskPriorityUrgent = "urgent"
)

// TaskStatus 任务状态常量
const (
	TaskStatusPending    = "pending"
	TaskStatusProcessing = "processing"
	TaskStatusCompleted  = "completed"
	TaskStatusCancelled  = "cancelled"
)

// TaskAction 操作类型常量
const (
	TaskActionCreate   = "create"   // 创建任务
	TaskActionAssign   = "assign"   // 分配任务
	TaskActionStart    = "start"    // 开始处理
	TaskActionComplete = "complete" // 完成任务
	TaskActionCancel   = "cancel"   // 取消任务
	TaskActionComment  = "comment"  // 添加评论
	TaskActionReassign = "reassign" // 重新分配
)

// GetPriorityLevel 获取优先级等级数值
func (t *Task) GetPriorityLevel() int {
	switch t.Priority {
	case TaskPriorityUrgent:
		return 4
	case TaskPriorityHigh:
		return 3
	case TaskPriorityMedium:
		return 2
	case TaskPriorityLow:
		return 1
	default:
		return 2
	}
}

// IsOverdue 判断任务是否超时
func (t *Task) IsOverdue() bool {
	// 根据任务类型和优先级判断超时时间
	var overdueHours int
	switch t.Priority {
	case TaskPriorityUrgent:
		overdueHours = 2 // 紧急任务2小时
	case TaskPriorityHigh:
		overdueHours = 8 // 高优先级8小时
	case TaskPriorityMedium:
		overdueHours = 24 // 中等优先级1天
	case TaskPriorityLow:
		overdueHours = 72 // 低优先级3天
	default:
		overdueHours = 24
	}

	if t.Status == TaskStatusPending || t.Status == TaskStatusProcessing {
		overdueTime := t.CreatedAt.Add(time.Duration(overdueHours) * time.Hour)
		return time.Now().After(overdueTime)
	}

	return false
}

// TableName 指定表名
func (Task) TableName() string {
	return "tasks"
}

// TableName 指定表名
func (TaskAction) TableName() string {
	return "task_actions"
}
