package service

import (
	"time"
)

// ============= 农机租赁审批相关类型 =============

// RateOrderRequest 评价订单请求
type RateOrderRequest struct {
	OrderID    uint64  `json:"order_id" validate:"required"`
	UserID     uint64  `json:"user_id" validate:"required"`
	RatingType string  `json:"rating_type" validate:"required,oneof=renter owner"`
	Rating     float32 `json:"rating" validate:"required,min=1,max=5"`
	Comment    string  `json:"comment" validate:"max=500"`
}

// GetRentalApplicationsRequest 获取租赁申请请求
type GetRentalApplicationsRequest struct {
	Status      string `json:"status,omitempty"`
	MachineType string `json:"machine_type,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
	RiskLevel   string `json:"risk_level,omitempty"`
	Page        int    `json:"page" validate:"min=1"`
	Limit       int    `json:"limit" validate:"min=1,max=100"`
	SortBy      string `json:"sort_by,omitempty"`
	SortOrder   string `json:"sort_order,omitempty"`
}

// GetRentalApplicationsResponse 获取租赁申请响应
type GetRentalApplicationsResponse struct {
	Applications []*RentalApplicationItem     `json:"applications"`
	Total        int64                        `json:"total"`
	Page         int                          `json:"page"`
	Limit        int                          `json:"limit"`
	Statistics   *RentalApplicationStatistics `json:"statistics"`
}

// RentalApplicationItem 租赁申请条目
type RentalApplicationItem struct {
	ID            uint      `json:"id"`
	ApplicationNo string    `json:"application_no"`
	UserID        uint64    `json:"user_id"`
	UserName      string    `json:"user_name"`
	MachineID     uint64    `json:"machine_id"`
	MachineName   string    `json:"machine_name"`
	MachineType   string    `json:"machine_type"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	TotalAmount   int64     `json:"total_amount"`
	Status        string    `json:"status"`
	RiskLevel     string    `json:"risk_level"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// RentalApplicationStatistics 租赁申请统计
type RentalApplicationStatistics struct {
	TotalApplications    int64 `json:"total_applications"`
	PendingApplications  int64 `json:"pending_applications"`
	ApprovedApplications int64 `json:"approved_applications"`
	RejectedApplications int64 `json:"rejected_applications"`
}

// GetRentalApplicationDetailRequest 获取租赁申请详情请求
type GetRentalApplicationDetailRequest struct {
	ID uint `json:"id" validate:"required"`
}

// GetRentalApplicationDetailResponse 获取租赁申请详情响应
type GetRentalApplicationDetailResponse struct {
	Application *RentalApplicationDetail `json:"application"`
	User        interface{}              `json:"user"`
	Machine     interface{}              `json:"machine"`
	ReviewLogs  interface{}              `json:"review_logs"`
}

// RentalApplicationDetail 租赁申请详情
type RentalApplicationDetail struct {
	ID             uint      `json:"id"`
	ApplicationNo  string    `json:"application_no"`
	UserID         uint      `json:"user_id"`
	UserName       string    `json:"user_name"`
	UserRealName   string    `json:"user_real_name"`
	UserPhone      string    `json:"user_phone"`
	MachineID      uint      `json:"machine_id"`
	MachineName    string    `json:"machine_name"`
	MachineType    string    `json:"machine_type"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	RentalLocation string    `json:"rental_location"`
	ContactPerson  string    `json:"contact_person"`
	ContactPhone   string    `json:"contact_phone"`
	BillingMethod  string    `json:"billing_method"`
	Quantity       float64   `json:"quantity"`
	TotalAmount    int64     `json:"total_amount"`
	DepositAmount  int64     `json:"deposit_amount"`
	Status         string    `json:"status"`
	RiskLevel      string    `json:"risk_level"`
	RiskFactors    []string  `json:"risk_factors"`
	AIAssessment   string    `json:"ai_assessment"`
	Remarks        string    `json:"remarks"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// MachineDetailSummary 农机详情摘要
type MachineDetailSummary struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	Model           string    `json:"model"`
	Brand           string    `json:"brand"`
	Year            int       `json:"year"`
	Status          string    `json:"status"`
	DailyRate       float64   `json:"daily_rate"`
	OwnerID         uint      `json:"owner_id"`
	OwnerName       string    `json:"owner_name"`
	OwnerPhone      string    `json:"owner_phone"`
	Location        string    `json:"location"`
	Condition       string    `json:"condition"`
	LastMaintenance time.Time `json:"last_maintenance"`
	TotalRentals    int       `json:"total_rentals"`
	AverageRating   float64   `json:"average_rating"`
}

// ApproveRentalRequest 批准租赁请求
type ApproveRentalRequest struct {
	ID           uint   `json:"id" validate:"required"`
	ReviewerID   uint64 `json:"reviewer_id" validate:"required"`
	ApprovalNote string `json:"approval_note" validate:"required"`
}

// RejectRentalRequest 拒绝租赁请求
type RejectRentalRequest struct {
	ID              uint   `json:"id" validate:"required"`
	ReviewerID      uint64 `json:"reviewer_id" validate:"required"`
	RejectionReason string `json:"rejection_reason" validate:"required"`
}

// GetRentalStatisticsRequest 获取租赁统计请求
type GetRentalStatisticsRequest struct {
	Period    string `json:"period,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

// GetRentalStatisticsResponse 获取租赁统计响应
type GetRentalStatisticsResponse struct {
	Overview         *RentalStatisticsOverview `json:"overview"`
	Trends           interface{}               `json:"trends"`
	MachineTypeStats interface{}               `json:"machine_type_stats"`
	RegionStats      interface{}               `json:"region_stats"`
}

// RentalStatisticsOverview 租赁统计概览
type RentalStatisticsOverview struct {
	TotalApplications    int64   `json:"total_applications"`
	PendingApplications  int64   `json:"pending_applications"`
	ApprovedApplications int64   `json:"approved_applications"`
	RejectedApplications int64   `json:"rejected_applications"`
	TotalRentalAmount    int64   `json:"total_rental_amount"`
	ApprovedRentalAmount int64   `json:"approved_rental_amount"`
	ApprovalRate         float64 `json:"approval_rate"`
}

// BatchApproveRentalsRequest 批量审批租赁请求
type BatchApproveRentalsRequest struct {
	ApplicationIDs []uint64 `json:"application_ids" validate:"required"`
	ReviewerID     uint64   `json:"reviewer_id" validate:"required"`
	Decision       string   `json:"decision" validate:"required,oneof=approve reject"`
	Note           string   `json:"note" validate:"required"`
}

// BatchApprovalResult 批量审批结果
type BatchApprovalResult struct {
	ApplicationID uint64 `json:"application_id"`
	Success       bool   `json:"success"`
	Error         string `json:"error,omitempty"`
}

// BatchApproveRentalsResponse 批量审批租赁响应
type BatchApproveRentalsResponse struct {
	Results      []*BatchApprovalResult `json:"results"`
	SuccessCount int                    `json:"success_count"`
	FailureCount int                    `json:"failure_count"`
	TotalCount   int                    `json:"total_count"`
}

// ============= 贷款审批增强类型 =============

// RetryAIAssessmentRequest 重试AI评估请求
type RetryAIAssessmentRequest struct {
	ID uint `json:"id" validate:"required"`
}

// BatchApproveLoanRequest 批量审批贷款请求
type BatchApproveLoanRequest struct {
	ApplicationIDs []uint   `json:"application_ids" validate:"required"`
	ReviewerID     uint64   `json:"reviewer_id" validate:"required"`
	Decision       string   `json:"decision" validate:"required,oneof=approve reject"`
	ReviewComments string   `json:"review_comments" validate:"required"`
	ApprovedAmount *float64 `json:"approved_amount,omitempty"`
	ApprovedTerm   *int     `json:"approved_term,omitempty"`
	InterestRate   *float64 `json:"interest_rate,omitempty"`
	AutoNotify     bool     `json:"auto_notify"`
}

// BatchApproveLoanResponse 批量审批贷款响应
type BatchApproveLoanResponse struct {
	TotalCount   int                    `json:"total_count"`
	SuccessCount int                    `json:"success_count"`
	FailureCount int                    `json:"failure_count"`
	Results      []BatchOperationResult `json:"results"`
	BatchID      string                 `json:"batch_id"`
	ProcessedBy  uint64                 `json:"processed_by"`
	ProcessedAt  time.Time              `json:"processed_at"`
}

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
	ID      uint   `json:"id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// EnableAutoApprovalRequest 启用自动审批请求
type EnableAutoApprovalRequest struct {
	OperatorID    uint64                 `json:"operator_id" validate:"required"`
	BusinessType  string                 `json:"business_type" validate:"required,oneof=loan machine_rental"`
	Conditions    AutoApprovalConditions `json:"conditions" validate:"required"`
	EffectiveDate time.Time              `json:"effective_date"`
	ExpiryDate    *time.Time             `json:"expiry_date,omitempty"`
	Description   string                 `json:"description" validate:"required"`
	ApprovalLimit *float64               `json:"approval_limit,omitempty"`
}

// DisableAutoApprovalRequest 禁用自动审批请求
type DisableAutoApprovalRequest struct {
	OperatorID   uint64 `json:"operator_id" validate:"required"`
	BusinessType string `json:"business_type" validate:"required,oneof=loan machine_rental"`
	Reason       string `json:"reason" validate:"required"`
}

// AutoApprovalConditions 自动审批条件
type AutoApprovalConditions struct {
	MaxAmount             float64  `json:"max_amount"`
	MinCreditScore        float64  `json:"min_credit_score"`
	MaxRiskLevel          string   `json:"max_risk_level"`
	MinConfidence         float64  `json:"min_confidence"`
	UserTypes             []string `json:"user_types,omitempty"`
	RequiredVerifications []string `json:"required_verifications,omitempty"`
}

// GetAutoApprovalConfigResponse 获取自动审批配置响应
type GetAutoApprovalConfigResponse struct {
	LoanConfig     *AutoApprovalConfig        `json:"loan_config"`
	RentalConfig   *AutoApprovalConfig        `json:"rental_config"`
	GlobalSettings GlobalAutoApprovalSettings `json:"global_settings"`
}

// AutoApprovalConfig 自动审批配置
type AutoApprovalConfig struct {
	Enabled       bool                   `json:"enabled"`
	Conditions    AutoApprovalConditions `json:"conditions"`
	CreatedBy     uint64                 `json:"created_by"`
	CreatedAt     time.Time              `json:"created_at"`
	EffectiveDate time.Time              `json:"effective_date"`
	ExpiryDate    *time.Time             `json:"expiry_date"`
	Description   string                 `json:"description"`
}

// GlobalAutoApprovalSettings 全局自动审批设置
type GlobalAutoApprovalSettings struct {
	MaxDailyAutoApprovals int             `json:"max_daily_auto_approvals"`
	MaxAmountPerDay       float64         `json:"max_amount_per_day"`
	RequireSecondApproval bool            `json:"require_second_approval"`
	MonitoringEnabled     bool            `json:"monitoring_enabled"`
	AlertThresholds       AlertThresholds `json:"alert_thresholds"`
}

// AlertThresholds 告警阈值
type AlertThresholds struct {
	HighRiskApprovals      int     `json:"high_risk_approvals"`
	LargeAmountApprovals   int     `json:"large_amount_approvals"`
	RejectionRateThreshold float64 `json:"rejection_rate_threshold"`
}

// GetApplicationsByRiskLevelRequest 按风险等级获取申请请求
type GetApplicationsByRiskLevelRequest struct {
	RiskLevel string `json:"risk_level" validate:"required,oneof=low medium high"`
	Status    string `json:"status,omitempty"`
	Page      int    `json:"page" validate:"min=1"`
	Limit     int    `json:"limit" validate:"min=1,max=100"`
}

// GetApplicationsByRiskLevelResponse 按风险等级获取申请响应
type GetApplicationsByRiskLevelResponse struct {
	Applications []LoanApplicationSummary `json:"applications"`
	Pagination   PaginationInfo           `json:"pagination"`
	RiskAnalysis RiskLevelAnalysis        `json:"risk_analysis"`
}

// RiskLevelAnalysis 风险等级分析
type RiskLevelAnalysis struct {
	RiskLevel      string  `json:"risk_level"`
	TotalCount     int     `json:"total_count"`
	ApprovalRate   float64 `json:"approval_rate"`
	AverageAmount  float64 `json:"average_amount"`
	AverageScore   float64 `json:"average_score"`
	TrendDirection string  `json:"trend_direction"`
}

// GetAIAssessmentHistoryRequest 获取AI评估历史请求
type GetAIAssessmentHistoryRequest struct {
	ApplicationID uint `json:"application_id" validate:"required"`
}

// GetAIAssessmentHistoryResponse 获取AI评估历史响应
type GetAIAssessmentHistoryResponse struct {
	ApplicationID uint                 `json:"application_id"`
	Assessments   []AIAssessmentRecord `json:"assessments"`
	Summary       AIAssessmentSummary  `json:"summary"`
}

// AIAssessmentRecord AI评估记录
type AIAssessmentRecord struct {
	ID                uint      `json:"id"`
	Version           int       `json:"version"`
	RiskLevel         string    `json:"risk_level"`
	RiskScore         float64   `json:"risk_score"`
	Decision          string    `json:"decision"`
	Confidence        float64   `json:"confidence"`
	RiskFactors       []string  `json:"risk_factors"`
	RecommendedAmount *float64  `json:"recommended_amount"`
	RecommendedTerm   *int      `json:"recommended_term"`
	RecommendedRate   *float64  `json:"recommended_rate"`
	ProcessingTime    float64   `json:"processing_time_ms"`
	ModelVersion      string    `json:"model_version"`
	CreatedAt         time.Time `json:"created_at"`
	Comments          string    `json:"comments"`
}

// AIAssessmentSummary AI评估摘要
type AIAssessmentSummary struct {
	TotalAssessments  int     `json:"total_assessments"`
	LatestVersion     int     `json:"latest_version"`
	CurrentRiskLevel  string  `json:"current_risk_level"`
	RiskTrend         string  `json:"risk_trend"`
	AverageConfidence float64 `json:"average_confidence"`
	ModelConsistency  float64 `json:"model_consistency"`
}

// CreateApplicationTaskRequest 创建申请任务请求
type CreateApplicationTaskRequest struct {
	ApplicationID uint                   `json:"application_id" validate:"required"`
	CreatorID     uint64                 `json:"creator_id" validate:"required"`
	TaskType      string                 `json:"task_type" validate:"required"`
	Priority      string                 `json:"priority" validate:"required"`
	Title         string                 `json:"title" validate:"required"`
	Description   string                 `json:"description" validate:"required"`
	AssignedTo    *uint64                `json:"assigned_to,omitempty"`
	DueDate       *time.Time             `json:"due_date,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// CreateApplicationTaskResponse 创建申请任务响应
type CreateApplicationTaskResponse struct {
	TaskID                  uint64     `json:"task_id"`
	ApplicationID           uint       `json:"application_id"`
	Status                  string     `json:"status"`
	CreatedAt               time.Time  `json:"created_at"`
	EstimatedCompletionTime *time.Time `json:"estimated_completion_time"`
}

// GetApplicationTasksRequest 获取申请任务请求
type GetApplicationTasksRequest struct {
	ApplicationID uint `json:"application_id" validate:"required"`
}

// GetApplicationTasksResponse 获取申请任务响应
type GetApplicationTasksResponse struct {
	ApplicationID uint        `json:"application_id"`
	Tasks         []TaskInfo  `json:"tasks"`
	Summary       TaskSummary `json:"summary"`
}

// TaskInfo 任务信息
type TaskInfo struct {
	ID             uint64     `json:"id"`
	Type           string     `json:"type"`
	Title          string     `json:"title"`
	Status         string     `json:"status"`
	Priority       string     `json:"priority"`
	AssignedTo     *uint64    `json:"assigned_to"`
	AssignedToName *string    `json:"assigned_to_name"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DueDate        *time.Time `json:"due_date"`
	Progress       int        `json:"progress"`
}

// TaskSummary 任务摘要
type TaskSummary struct {
	TotalTasks      int `json:"total_tasks"`
	PendingTasks    int `json:"pending_tasks"`
	InProgressTasks int `json:"in_progress_tasks"`
	CompletedTasks  int `json:"completed_tasks"`
	OverdueTasks    int `json:"overdue_tasks"`
}

// GetAdvancedStatisticsRequest 获取高级统计请求
type GetAdvancedStatisticsRequest struct {
	Period        string `json:"period,omitempty"`
	StartDate     string `json:"start_date,omitempty"`
	EndDate       string `json:"end_date,omitempty"`
	IncludeTrends bool   `json:"include_trends"`
}

// GetAdvancedStatisticsResponse 获取高级统计响应
type GetAdvancedStatisticsResponse struct {
	Overview      StatisticsOverview `json:"overview"`
	TrendAnalysis *TrendAnalysis     `json:"trend_analysis,omitempty"`
	RiskAnalysis  RiskAnalysis       `json:"risk_analysis"`
	Performance   PerformanceMetrics `json:"performance"`
	Predictions   *PredictionData    `json:"predictions,omitempty"`
}

// StatisticsOverview 统计概览
type StatisticsOverview struct {
	TotalApplications    int     `json:"total_applications"`
	ApprovedApplications int     `json:"approved_applications"`
	RejectedApplications int     `json:"rejected_applications"`
	PendingApplications  int     `json:"pending_applications"`
	TotalAmount          float64 `json:"total_amount"`
	ApprovedAmount       float64 `json:"approved_amount"`
	ApprovalRate         float64 `json:"approval_rate"`
	AverageProcessTime   float64 `json:"average_process_time_hours"`
}

// TrendAnalysis 趋势分析
type TrendAnalysis struct {
	ApplicationTrend  []TrendPoint `json:"application_trend"`
	AmountTrend       []TrendPoint `json:"amount_trend"`
	ApprovalRateTrend []TrendPoint `json:"approval_rate_trend"`
}

// TrendPoint 趋势点
type TrendPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
	Count int     `json:"count,omitempty"`
}

// RiskAnalysis 风险分析
type RiskAnalysis struct {
	RiskDistribution map[string]int    `json:"risk_distribution"`
	HighRiskFactors  []RiskFactorStats `json:"high_risk_factors"`
	RiskTrends       []RiskTrendPoint  `json:"risk_trends"`
}

// RiskFactorStats 风险因子统计
type RiskFactorStats struct {
	Factor      string  `json:"factor"`
	Frequency   int     `json:"frequency"`
	ImpactScore float64 `json:"impact_score"`
}

// RiskTrendPoint 风险趋势点
type RiskTrendPoint struct {
	Date         string         `json:"date"`
	Distribution map[string]int `json:"distribution"`
	AverageScore float64        `json:"average_score"`
}

// PerformanceMetrics 性能指标
type PerformanceMetrics struct {
	AverageProcessingTime float64 `json:"average_processing_time_hours"`
	AutoApprovalRate      float64 `json:"auto_approval_rate"`
	ManualReviewRate      float64 `json:"manual_review_rate"`
	ReviewerEfficiency    float64 `json:"reviewer_efficiency"`
	SystemUptime          float64 `json:"system_uptime"`
	ErrorRate             float64 `json:"error_rate"`
}

// PredictionData 预测数据
type PredictionData struct {
	NextWeekApplications int                `json:"next_week_applications"`
	NextMonthAmount      float64            `json:"next_month_amount"`
	ExpectedApprovalRate float64            `json:"expected_approval_rate"`
	ResourceRequirement  ResourcePrediction `json:"resource_requirement"`
}

// ResourcePrediction 资源预测
type ResourcePrediction struct {
	RequiredReviewers int     `json:"required_reviewers"`
	EstimatedWorkload float64 `json:"estimated_workload_hours"`
	PeakLoadTime      string  `json:"peak_load_time"`
}

// ============= 通用类型 =============

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

// UserProfileSummary 用户资料摘要
type UserProfileSummary struct {
	UserID               uint       `json:"user_id"`
	Username             string     `json:"username"`
	RealName             string     `json:"real_name"`
	Phone                string     `json:"phone"`
	Email                string     `json:"email"`
	UserType             string     `json:"user_type"`
	IsVerified           bool       `json:"is_verified"`
	CreditScore          float64    `json:"credit_score"`
	RegistrationDate     time.Time  `json:"registration_date"`
	LastLoginDate        *time.Time `json:"last_login_date"`
	TotalApplications    int        `json:"total_applications"`
	ApprovedApplications int        `json:"approved_applications"`
	UserRating           float64    `json:"user_rating"`
}

// RiskAssessmentResult 风险评估结果
type RiskAssessmentResult struct {
	RiskLevel    string    `json:"risk_level"`
	RiskScore    float64   `json:"risk_score"`
	RiskFactors  []string  `json:"risk_factors"`
	Confidence   float64   `json:"confidence"`
	AssessedAt   time.Time `json:"assessed_at"`
	AssessedBy   string    `json:"assessed_by"`
	ModelVersion string    `json:"model_version"`
}

// ReviewRecord 审核记录
type ReviewRecord struct {
	ID           uint      `json:"id"`
	ReviewerID   uint64    `json:"reviewer_id"`
	ReviewerName string    `json:"reviewer_name"`
	Action       string    `json:"action"`
	Comments     string    `json:"comments"`
	CreatedAt    time.Time `json:"created_at"`
	Duration     float64   `json:"duration_minutes"`
}

// RentalStatisticsSummary 租赁统计摘要
type RentalStatisticsSummary struct {
	TotalApplications   int     `json:"total_applications"`
	PendingApplications int     `json:"pending_applications"`
	ApprovalRate        float64 `json:"approval_rate"`
	AverageProcessTime  float64 `json:"average_process_time_hours"`
}

// LoanApplicationSummary 贷款申请摘要
type LoanApplicationSummary struct {
	ID              uint      `json:"id"`
	ApplicationNo   string    `json:"application_no"`
	UserID          uint      `json:"user_id"`
	ApplicantName   string    `json:"applicant_name"`
	ApplyAmount     int64     `json:"apply_amount"`
	ApplyTermMonths int       `json:"apply_term_months"`
	Status          string    `json:"status"`
	RiskLevel       string    `json:"risk_level"`
	CreatedAt       time.Time `json:"created_at"`
	AssignedTo      *uint64   `json:"assigned_to"`
}

// RequestMoreInfoRentalRequest 请求农机租赁补充信息请求
type RequestMoreInfoRentalRequest struct {
	ID           uint   `json:"id" binding:"required"`
	ReviewerID   uint64 `json:"reviewer_id" binding:"required"`
	RequiredInfo string `json:"required_info" binding:"required"` // 需要补充的信息描述
	RequestNote  string `json:"request_note"`                     // 审批员的备注
}

// BatchApproveRentalsRequest 批量批准农机租赁请求
// ... existing code ...
