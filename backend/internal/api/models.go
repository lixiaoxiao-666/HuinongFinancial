package api

import (
	"time"
)

// CommonResponse 通用响应格式
type CommonResponse struct {
	Code    int         `json:"code"`    // 0表示成功，其他表示错误码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 业务数据
}

// PaginationResponse 分页响应格式
type PaginationResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`  // 数据列表
	Total   int64       `json:"total"` // 总数量
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

// ================================
// 用户相关模型
// ================================

// User 用户信息
type User struct {
	UserID     string    `json:"user_id"`
	Username   string    `json:"username"`
	RealName   string    `json:"real_name"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	IDCard     string    `json:"id_card"`
	Address    string    `json:"address"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=30"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	RealName        string `json:"real_name" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
}

// UserRegisterResponse 用户注册响应
type UserRegisterResponse struct {
	UserID string `json:"user_id"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	Token    string `json:"token"`
	UserInfo User   `json:"user_info"`
}

// UserUpdateRequest 用户信息更新请求
type UserUpdateRequest struct {
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	IDCard   string `json:"id_card"`
	Address  string `json:"address"`
}

// ================================
// 贷款相关模型
// ================================

// LoanProduct 贷款产品
type LoanProduct struct {
	ProductID             string             `json:"product_id"`
	Name                  string             `json:"name"`
	Description           string             `json:"description"`
	Category              string             `json:"category"`
	MinAmount             int64              `json:"min_amount"`
	MaxAmount             int64              `json:"max_amount"`
	MinTermMonths         int                `json:"min_term_months"`
	MaxTermMonths         int                `json:"max_term_months"`
	InterestRateYearly    string             `json:"interest_rate_yearly"`
	RepaymentMethods      []string           `json:"repayment_methods"`
	ApplicationConditions string             `json:"application_conditions,omitempty"`
	RequiredDocuments     []RequiredDocument `json:"required_documents,omitempty"`
	Status                int                `json:"status"`
}

// RequiredDocument 必需文档
type RequiredDocument struct {
	Type string `json:"type"`
	Desc string `json:"desc"`
}

// LoanApplicationRequest 贷款申请请求
type LoanApplicationRequest struct {
	ProductID         string         `json:"product_id" binding:"required"`
	Amount            int64          `json:"amount" binding:"required,gt=0"`
	TermMonths        int            `json:"term_months" binding:"required,gt=0"`
	Purpose           string         `json:"purpose" binding:"required"`
	ApplicantInfo     ApplicantInfo  `json:"applicant_info" binding:"required"`
	UploadedDocuments []DocumentInfo `json:"uploaded_documents"`
}

// ApplicantInfo 申请人信息
type ApplicantInfo struct {
	RealName     string `json:"real_name" binding:"required"`
	IDCardNumber string `json:"id_card_number" binding:"required"`
	Address      string `json:"address" binding:"required"`
}

// DocumentInfo 文档信息
type DocumentInfo struct {
	DocType string `json:"doc_type" binding:"required"`
	FileID  string `json:"file_id" binding:"required"`
}

// LoanApplicationResponse 贷款申请响应
type LoanApplicationResponse struct {
	ApplicationID string `json:"application_id"`
}

// LoanApplicationDetail 贷款申请详情
type LoanApplicationDetail struct {
	ApplicationID  string               `json:"application_id"`
	ProductID      string               `json:"product_id"`
	UserID         string               `json:"user_id"`
	Amount         int64                `json:"amount"`
	TermMonths     int                  `json:"term_months"`
	Purpose        string               `json:"purpose"`
	Status         string               `json:"status"`
	SubmittedAt    time.Time            `json:"submitted_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	ApprovedAmount *int64               `json:"approved_amount,omitempty"`
	Remarks        string               `json:"remarks,omitempty"`
	History        []ApplicationHistory `json:"history,omitempty"`
}

// ApplicationHistory 申请历史
type ApplicationHistory struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Operator  string    `json:"operator"`
}

// MyLoanApplication 我的贷款申请
type MyLoanApplication struct {
	ApplicationID string    `json:"application_id"`
	ProductName   string    `json:"product_name"`
	Amount        int64     `json:"amount"`
	Status        string    `json:"status"`
	SubmittedAt   time.Time `json:"submitted_at"`
}

// ================================
// 文件服务相关模型
// ================================

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	FileID     string    `json:"file_id"`
	FileName   string    `json:"file_name"`
	FileSize   int64     `json:"file_size"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// ================================
// AI智能体相关模型
// ================================

// AIApplicationInfo AI申请信息
type AIApplicationInfo struct {
	ApplicationType string                 `json:"application_type"`
	ApplicationID   string                 `json:"application_id"`
	UserID          string                 `json:"user_id"`
	Status          string                 `json:"status"`
	BasicInfo       map[string]interface{} `json:"basic_info"`
	FinancialInfo   map[string]interface{} `json:"financial_info"`
	ProductInfo     map[string]interface{} `json:"product_info"`
	Documents       []DocumentInfo         `json:"documents"`
	SubmittedAt     time.Time              `json:"submitted_at"`
}

// AIDecisionResponse AI决策响应
type AIDecisionResponse struct {
	Success         bool   `json:"success"`
	AIOperationID   string `json:"ai_operation_id"`
	ProcessedAt     string `json:"processed_at"`
	NextStepMessage string `json:"next_step_message"`
}

// ExternalDataResponse 外部数据响应
type ExternalDataResponse struct {
	UserID         string                 `json:"user_id"`
	DataTypes      []string               `json:"data_types"`
	CreditData     map[string]interface{} `json:"credit_data,omitempty"`
	BankData       map[string]interface{} `json:"bank_data,omitempty"`
	BlacklistData  map[string]interface{} `json:"blacklist_data,omitempty"`
	GovernmentData map[string]interface{} `json:"government_data,omitempty"`
	FarmingData    map[string]interface{} `json:"farming_data,omitempty"`
	RetrievedAt    time.Time              `json:"retrieved_at"`
}

// AIModelConfigResponse AI模型配置响应
type AIModelConfigResponse struct {
	Models         []AIModelConfig        `json:"models"`
	RiskThresholds map[string]float64     `json:"risk_thresholds"`
	BusinessRules  map[string]interface{} `json:"business_rules"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// AIModelConfig AI模型配置
type AIModelConfig struct {
	ModelName string `json:"model_name"`
	Version   string `json:"version"`
	Type      string `json:"type"`
	Enabled   bool   `json:"enabled"`
}

// MachineryLeasingApplicationInfo 农机租赁申请信息
type MachineryLeasingApplicationInfo struct {
	ApplicationID  string                 `json:"application_id"`
	LesseeInfo     map[string]interface{} `json:"lessee_info"`
	LessorInfo     map[string]interface{} `json:"lessor_info"`
	MachineryInfo  map[string]interface{} `json:"machinery_info"`
	LeasingDetails map[string]interface{} `json:"leasing_details"`
	Status         string                 `json:"status"`
	SubmittedAt    time.Time              `json:"submitted_at"`
}

// AIOperationLog AI操作日志
type AIOperationLog struct {
	OperationID     string    `json:"operation_id"`
	ApplicationID   string    `json:"application_id"`
	ApplicationType string    `json:"application_type"`
	Operation       string    `json:"operation"`
	Result          string    `json:"result"`
	Details         string    `json:"details"`
	Timestamp       time.Time `json:"timestamp"`
}

// AIOperationLogsResponse AI操作日志响应
type AIOperationLogsResponse struct {
	Logs       []AIOperationLog `json:"logs"`
	Pagination Pagination       `json:"pagination"`
}

// Pagination 分页信息
type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalCount  int `json:"total_count"`
	Limit       int `json:"limit"`
}

// ================================
// 管理员相关模型
// ================================

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminLoginResponse 管理员登录响应
type AdminLoginResponse struct {
	Token     string `json:"token"`
	AdminInfo Admin  `json:"admin_info"`
}

// Admin 管理员信息
type Admin struct {
	AdminID   string    `json:"admin_id"`
	Username  string    `json:"username"`
	RealName  string    `json:"real_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// LoanApplicationForAdmin 管理员视角的贷款申请
type LoanApplicationForAdmin struct {
	ApplicationID  string    `json:"application_id"`
	UserID         string    `json:"user_id"`
	UserName       string    `json:"user_name"`
	ProductName    string    `json:"product_name"`
	Amount         int64     `json:"amount"`
	TermMonths     int       `json:"term_months"`
	Status         string    `json:"status"`
	SubmittedAt    time.Time `json:"submitted_at"`
	AIRiskScore    *int      `json:"ai_risk_score,omitempty"`
	AISuggestion   string    `json:"ai_suggestion,omitempty"`
	RequiresReview bool      `json:"requires_review"`
}

// ApprovalRequest 审批请求
type ApprovalRequest struct {
	Action         string `json:"action" binding:"required,oneof=approve reject"`
	ApprovedAmount *int64 `json:"approved_amount"`
	Comments       string `json:"comments"`
}

// ApprovalResponse 审批响应
type ApprovalResponse struct {
	ApplicationID string `json:"application_id"`
	Action        string `json:"action"`
	ProcessedAt   string `json:"processed_at"`
}

// ================================
// 常量定义
// ================================

// 错误码常量
const (
	CodeSuccess               = 0
	CodeBadRequest            = 1001
	CodeUnauthorized          = 2001
	CodeForbidden             = 3001
	CodeNotFound              = 4001
	CodeInternalError         = 5001
	CodeAIApplicationNotFound = 4002
	CodeAIStatusConflict      = 4003
	CodeAIAnalysisError       = 4004
	CodeAIServiceUnavailable  = 5002
)

// 申请状态常量
const (
	StatusSubmitted    = "SUBMITTED"              // 已提交
	StatusAIReviewing  = "AI_REVIEWING"           // AI审核中
	StatusAIApproved   = "AI_APPROVED"            // AI通过
	StatusAIRejected   = "AI_REJECTED"            // AI拒绝
	StatusManualReview = "MANUAL_REVIEW_REQUIRED" // 需要人工审核
	StatusApproved     = "APPROVED"               // 已批准
	StatusRejected     = "REJECTED"               // 已拒绝
	StatusRequireInfo  = "REQUIRE_INFO"           // 需要补充信息
)

// 用户类型常量
const (
	UserTypeUser  = "user"  // 普通用户
	UserTypeAdmin = "admin" // 管理员
)

// 审批决策常量
const (
	DecisionApproved    = "approved"     // 批准
	DecisionRejected    = "rejected"     // 拒绝
	DecisionRequireInfo = "require_info" // 需要补充信息
)

// AI决策常量（贷款申请）
const (
	AIDecisionAutoApproved       = "AUTO_APPROVED"        // 自动批准
	AIDecisionAutoRejected       = "AUTO_REJECTED"        // 自动拒绝
	AIDecisionRequireHumanReview = "REQUIRE_HUMAN_REVIEW" // 需要人工审核
)

// AI决策常量（农机租赁）
const (
	AIDecisionAutoApprove              = "AUTO_APPROVE"               // 自动通过
	AIDecisionAutoReject               = "AUTO_REJECT"                // 自动拒绝
	AIDecisionRequireDepositAdjustment = "REQUIRE_DEPOSIT_ADJUSTMENT" // 需要调整押金
)

// 申请类型常量
const (
	ApplicationTypeLoan             = "LOAN_APPLICATION"  // 贷款申请
	ApplicationTypeMachineryLeasing = "MACHINERY_LEASING" // 农机租赁申请
)
