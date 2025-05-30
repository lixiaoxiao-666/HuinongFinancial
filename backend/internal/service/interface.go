package service

import (
	"context"
	"huinong-backend/internal/model"
	"huinong-backend/internal/repository"
	"io"
	"time"
)

// UserService 用户服务接口
type UserService interface {
	// 用户注册和登录
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context, sessionID string) error
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)

	// 用户管理
	GetProfile(ctx context.Context, userID uint64) (*UserProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uint64, req *UpdateProfileRequest) error
	ChangePassword(ctx context.Context, userID uint64, req *ChangePasswordRequest) error
	FreezeUser(ctx context.Context, userID uint64, reason string) error
	UnfreezeUser(ctx context.Context, userID uint64) error

	// 获取用户信息
	GetUserByID(userID uint) (*model.User, error)

	// 用户认证
	SubmitRealNameAuth(ctx context.Context, userID uint64, req *RealNameAuthRequest) error
	SubmitBankCardAuth(ctx context.Context, userID uint64, req *BankCardAuthRequest) error
	ReviewAuth(ctx context.Context, authID uint64, req *ReviewAuthRequest) error

	// 用户标签
	AddUserTag(ctx context.Context, userID uint64, req *AddTagRequest) error
	GetUserTags(ctx context.Context, userID uint64, tagType string) ([]*model.UserTag, error)
	RemoveUserTag(ctx context.Context, userID uint64, tagKey string) error

	// 列表查询
	ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error)
	GetUserStatistics(ctx context.Context) (*UserStatistics, error)

	// 头像上传
	UploadAvatar(ctx context.Context, userID uint64, file io.Reader, filename string) (string, error)
}

// LoanService 贷款服务接口
type LoanService interface {
	// 贷款产品查询
	GetProducts(ctx context.Context, req *GetProductsRequest) (*GetProductsResponse, error)
	CreateApplication(ctx context.Context, req *CreateApplicationRequest) (*CreateApplicationResponse, error)
	GetApplicationDetails(ctx context.Context, req *GetApplicationDetailsRequest) (*GetApplicationDetailsResponse, error)
	GetUserApplications(ctx context.Context, req *GetUserApplicationsRequest) (*GetUserApplicationsResponse, error)
	ApproveApplication(ctx context.Context, req *ApproveApplicationRequest) error
	RejectApplication(ctx context.Context, req *RejectApplicationRequest) error
	GetAdminApplications(ctx context.Context, req *GetAdminApplicationsRequest) (*GetAdminApplicationsResponse, error)
	GetStatistics(ctx context.Context) (*LoanStatisticsResponse, error)

	// 获取贷款申请信息
	GetLoanApplicationByID(applicationID uint) (*model.LoanApplication, error)
	UpdateLoanApplication(application *model.LoanApplication) error
	CreateApprovalLog(log *model.ApprovalLog) error

	// 贷款产品管理
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*model.LoanProduct, error)
	UpdateProduct(ctx context.Context, productID uint64, req *UpdateProductRequest) error
	DeleteProduct(ctx context.Context, productID uint64) error
	GetProduct(ctx context.Context, productID uint64) (*model.LoanProduct, error)
	ListProducts(ctx context.Context, req *repository.ListProductsRequest) (*repository.ListProductsResponse, error)
	GetActiveProducts(ctx context.Context, userType string) ([]*model.LoanProduct, error)

	// 审批流程
	StartReview(ctx context.Context, applicationID uint64, reviewerID uint64) error
	ReturnApplication(ctx context.Context, applicationID uint64, req *ReturnRequest) error

	// AI辅助审批
	TriggerAIAssessment(ctx context.Context, applicationID uint64, workflowType string) error
	ProcessDifyCallback(ctx context.Context, req *DifyCallbackRequest) error

	// 统计和报表
	GetLoanStatistics(ctx context.Context, req *StatisticsRequest) (*LoanStatistics, error)
	GenerateApprovalReport(ctx context.Context, req *ReportRequest) (*ApprovalReport, error)
}

// MachineService 农机租赁服务接口
type MachineService interface {
	// 农机管理
	RegisterMachine(ctx context.Context, req *RegisterMachineRequest) (*RegisterMachineResponse, error)
	GetMachine(ctx context.Context, machineID uint64) (*MachineDetailResponse, error)
	UpdateMachine(ctx context.Context, machineID uint64, req *UpdateMachineRequest) error
	DeleteMachine(ctx context.Context, machineID uint64) error
	GetUserMachines(ctx context.Context, userID uint64, req *GetUserMachinesRequest) (*GetUserMachinesResponse, error)

	// 获取农机和订单信息
	GetMachineByID(machineID uint) (*model.Machine, error)
	GetRentalOrderByID(orderID uint) (*model.RentalOrder, error)
	CheckRentalTimeConflict(machineID uint, startTime, endTime time.Time, excludeOrderID uint) (bool, error)

	// 农机搜索
	SearchMachines(ctx context.Context, req *SearchMachinesRequest) (*SearchMachinesResponse, error)
	GetAvailableMachines(ctx context.Context, req *GetAvailableMachinesRequest) (*GetAvailableMachinesResponse, error)

	// 租赁订单
	CreateRentalOrder(ctx context.Context, req *CreateRentalOrderRequest) (*CreateRentalOrderResponse, error)
	GetRentalOrder(ctx context.Context, orderID uint64) (*RentalOrderDetailResponse, error)
	ConfirmOrder(ctx context.Context, orderID uint64, req *ConfirmOrderRequest) error
	CancelOrder(ctx context.Context, orderID uint64, req *CancelOrderRequest) error
	GetUserOrders(ctx context.Context, userID uint64, req *GetUserOrdersRequest) (*GetUserOrdersResponse, error)

	// 订单支付
	PayOrder(ctx context.Context, orderID uint64, req *PayOrderRequest) (*PayOrderResponse, error)
	CompleteOrder(ctx context.Context, orderID uint64, req *CompleteOrderRequest) error

	// 评价系统
	SubmitRating(ctx context.Context, orderID uint64, req *SubmitRatingRequest) error
	GetMachineRatings(ctx context.Context, machineID uint64) (*MachineRatingsResponse, error)
}

// ContentService 内容管理服务接口
type ContentService interface {
	// 文章管理
	CreateArticle(ctx context.Context, req *CreateArticleRequest) (*CreateArticleResponse, error)
	GetArticle(ctx context.Context, articleID uint64) (*ArticleDetailResponse, error)
	UpdateArticle(ctx context.Context, articleID uint64, req *UpdateArticleRequest) error
	DeleteArticle(ctx context.Context, articleID uint64) error
	ListArticles(ctx context.Context, req *ListArticlesRequest) (*ListArticlesResponse, error)
	PublishArticle(ctx context.Context, articleID uint64) error
	GetFeaturedArticles(ctx context.Context, limit int) (*FeaturedArticlesResponse, error)

	// 分类管理
	CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*model.Category, error)
	GetCategories(ctx context.Context) ([]*model.Category, error)
	UpdateCategory(ctx context.Context, categoryID uint64, req *UpdateCategoryRequest) error
	DeleteCategory(ctx context.Context, categoryID uint64) error

	// 专家管理
	CreateExpert(ctx context.Context, req *CreateExpertRequest) (*model.Expert, error)
	GetExpert(ctx context.Context, expertID uint64) (*ExpertDetailResponse, error)
	ListExperts(ctx context.Context, req *ListExpertsRequest) (*ListExpertsResponse, error)
	UpdateExpert(ctx context.Context, expertID uint64, req *UpdateExpertRequest) error
	DeleteExpert(ctx context.Context, expertID uint64) error

	// 专家咨询
	SubmitConsultation(ctx context.Context, req *SubmitConsultationRequest) (*SubmitConsultationResponse, error)
	GetConsultations(ctx context.Context, userID uint64) (*ConsultationsResponse, error)
}

// SystemService 系统服务接口
type SystemService interface {
	// 配置管理
	GetConfig(ctx context.Context, configKey string) (string, error)
	SetConfig(ctx context.Context, configKey, configValue string) error
	GetConfigs(ctx context.Context, configGroup string) (map[string]string, error)

	// 文件管理
	UploadFile(ctx context.Context, req *UploadFileRequest) (*UploadFileResponse, error)
	GetFile(ctx context.Context, fileID uint64) (*model.FileUpload, error)
	DeleteFile(ctx context.Context, fileID uint64) error

	// 健康检查
	HealthCheck(ctx context.Context) (*HealthCheckResponse, error)

	// 系统统计
	GetSystemStats(ctx context.Context) (*SystemStatsResponse, error)

	// 系统版本
	GetSystemVersion(ctx context.Context) (*SystemVersionResponse, error)

	// 获取公共配置
	GetPublicConfigs(ctx context.Context) (map[string]string, error)
}

// OAService OA后台服务接口
type OAService interface {
	// OA用户管理
	CreateOAUser(ctx context.Context, req *CreateOAUserRequest) (*model.OAUser, error)
	GetOAUser(ctx context.Context, userID uint64) (*model.OAUser, error)
	UpdateOAUser(ctx context.Context, userID uint64, req *UpdateOAUserRequest) error
	DeleteOAUser(ctx context.Context, userID uint64) error
	ListOAUsers(ctx context.Context, req *ListOAUsersRequest) (*ListOAUsersResponse, error)

	// OA角色管理
	CreateRole(ctx context.Context, req *CreateRoleRequest) (*model.OARole, error)
	GetRole(ctx context.Context, roleID uint64) (*model.OARole, error)
	UpdateRole(ctx context.Context, roleID uint64, req *UpdateRoleRequest) error
	DeleteRole(ctx context.Context, roleID uint64) error
	ListRoles(ctx context.Context) ([]*model.OARole, error)

	// OA认证
	OALogin(ctx context.Context, req *OALoginRequest) (*OALoginResponse, error)
	OALogout(ctx context.Context, sessionID string) error

	// 工作台
	GetDashboard(ctx context.Context, userID uint64) (*DashboardResponse, error)
	GetPendingTasks(ctx context.Context, userID uint64) (*PendingTasksResponse, error)
}

// TaskService 任务服务接口
type TaskService interface {
	// 任务管理
	CreateTask(ctx context.Context, req *CreateTaskRequest) (*CreateTaskResponse, error)
	GetTask(ctx context.Context, taskID uint64) (*TaskDetailResponse, error)
	UpdateTask(ctx context.Context, taskID uint64, req *UpdateTaskRequest) error
	DeleteTask(ctx context.Context, taskID uint64) error
	ListTasks(ctx context.Context, req *ListTasksRequest) (*ListTasksResponse, error)

	// 任务分配
	AssignTask(ctx context.Context, taskID uint64, assigneeID uint64) error
	UnassignTask(ctx context.Context, taskID uint64) error

	// 任务进度
	GetTaskProgress(ctx context.Context, taskID uint64) (*TaskProgressResponse, error)
	UpdateTaskProgress(ctx context.Context, taskID uint64, progress float64) error

	// 任务完成
	CompleteTask(ctx context.Context, taskID uint64) error
	ReassignTask(ctx context.Context, taskID uint64, newAssigneeID uint64) error
}

// ==================== 基础结构体定义 ====================

// 用户相关结构体
type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username"`
	UserType string `json:"user_type" binding:"required"`
	RealName string `json:"real_name"`
	Email    string `json:"email"`
	Province string `json:"province"`
	City     string `json:"city"`
	County   string `json:"county"`
	Address  string `json:"address"`
	SmsCode  string `json:"sms_code" binding:"required"`
}

type RegisterResponse struct {
	User         *model.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int         `json:"expires_in"`
}

type LoginRequest struct {
	Phone      string `json:"phone" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Platform   string `json:"platform" binding:"required"`
	DeviceID   string `json:"device_id"`
	DeviceType string `json:"device_type"`
	AppVersion string `json:"app_version"`
}

type LoginResponse struct {
	User         *model.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int         `json:"expires_in"`
	SessionID    string      `json:"session_id"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type UserProfileResponse struct {
	User     *model.User                `json:"user"`
	AuthInfo map[string]*model.UserAuth `json:"auth_info"`
	Tags     []*model.UserTag           `json:"tags"`
	Sessions []*model.UserSession       `json:"sessions,omitempty"`
}

type UpdateProfileRequest struct {
	Username string `json:"username"`
	RealName string `json:"real_name"`
	Email    string `json:"email"`
	Province string `json:"province"`
	City     string `json:"city"`
	County   string `json:"county"`
	Address  string `json:"address"`
	Avatar   string `json:"avatar"`
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type RealNameAuthRequest struct {
	RealName    string `json:"real_name" binding:"required"`
	IDCard      string `json:"id_card" binding:"required"`
	IDCardFront string `json:"id_card_front" binding:"required"`
	IDCardBack  string `json:"id_card_back" binding:"required"`
}

type BankCardAuthRequest struct {
	BankCard    string `json:"bank_card" binding:"required"`
	BankName    string `json:"bank_name" binding:"required"`
	AccountName string `json:"account_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type ReviewAuthRequest struct {
	Status     string `json:"status" binding:"required"` // approved, rejected
	ReviewNote string `json:"review_note"`
}

type AddTagRequest struct {
	TagKey   string `json:"tag_key" binding:"required"`
	TagValue string `json:"tag_value" binding:"required"`
	TagType  string `json:"tag_type" binding:"required"`
}

type ListUsersRequest struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	UserType string `json:"user_type"`
	Status   string `json:"status"`
	Keyword  string `json:"keyword"`
}

type ListUsersResponse struct {
	Users []*model.User `json:"users"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
}

type UserStatistics struct {
	TotalUsers      int64            `json:"total_users"`
	ActiveUsers     int64            `json:"active_users"`
	UsersByType     map[string]int64 `json:"users_by_type"`
	NewUsersToday   int64            `json:"new_users_today"`
	VerifiedUsers   int64            `json:"verified_users"`
	AuthCompletions map[string]int64 `json:"auth_completions"`
}

// 贷款相关结构体
type GetProductsRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type ProductResponse struct {
	ID               uint    `json:"id"`
	ProductName      string  `json:"product_name"`
	ProductCode      string  `json:"product_code"`
	ProductType      string  `json:"product_type"`
	Description      string  `json:"description"`
	MinAmount        int64   `json:"min_amount"`
	MaxAmount        int64   `json:"max_amount"`
	InterestRate     float64 `json:"interest_rate"`
	TermMonths       int     `json:"term_months"`
	RequiredAuth     string  `json:"required_auth"`
	IsActive         bool    `json:"is_active"`
	EligibleUserType string  `json:"eligible_user_type"`
}

type GetProductsResponse struct {
	Products []*ProductResponse `json:"products"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	Limit    int                `json:"limit"`
}

type CreateApplicationRequest struct {
	ProductID     uint   `json:"product_id" binding:"required"`
	LoanAmount    int64  `json:"loan_amount" binding:"required"`
	LoanPurpose   string `json:"loan_purpose" binding:"required"`
	TermMonths    int    `json:"term_months" binding:"required"`
	ContactPhone  string `json:"contact_phone" binding:"required"`
	ContactEmail  string `json:"contact_email"`
	MaterialsJSON string `json:"materials_json"`
	Remarks       string `json:"remarks"`
}

type CreateApplicationResponse struct {
	ID            uint      `json:"id"`
	ApplicationNo string    `json:"application_no"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type GetApplicationDetailsRequest struct {
	ID uint `json:"id"`
}

type ApplicationDetailsResponse struct {
	ID            uint      `json:"id"`
	ApplicationNo string    `json:"application_no"`
	ProductID     uint      `json:"product_id"`
	LoanAmount    int64     `json:"loan_amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type ApprovalLogResponse struct {
	ID        uint      `json:"id"`
	Step      string    `json:"step"`
	Status    string    `json:"status"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

type DifyLogResponse struct {
	ID           uint      `json:"id"`
	WorkflowType string    `json:"workflow_type"`
	Status       string    `json:"status"`
	Result       string    `json:"result"`
	CreatedAt    time.Time `json:"created_at"`
}

type GetApplicationDetailsResponse struct {
	Application  *ApplicationDetailsResponse `json:"application"`
	Product      *ProductResponse            `json:"product"`
	ApprovalLogs []*ApprovalLogResponse      `json:"approval_logs"`
	DifyLogs     []*DifyLogResponse          `json:"dify_logs"`
}

type GetUserApplicationsRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Status string `json:"status"`
}

type UserApplicationResponse struct {
	ID            uint      `json:"id"`
	ApplicationNo string    `json:"application_no"`
	ProductName   string    `json:"product_name"`
	LoanAmount    int64     `json:"loan_amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type GetUserApplicationsResponse struct {
	Applications []*UserApplicationResponse `json:"applications"`
	Total        int64                      `json:"total"`
	Page         int                        `json:"page"`
	Limit        int                        `json:"limit"`
}

type ApproveApplicationRequest struct {
	ID           uint   `json:"id"`
	ApprovalNote string `json:"approval_note"`
}

type RejectApplicationRequest struct {
	ID            uint   `json:"id"`
	RejectionNote string `json:"rejection_note"`
}

type GetAdminApplicationsRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Status string `json:"status"`
}

type AdminApplicationResponse struct {
	ID            uint      `json:"id"`
	ApplicationNo string    `json:"application_no"`
	UserName      string    `json:"user_name"`
	ProductName   string    `json:"product_name"`
	LoanAmount    int64     `json:"loan_amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type GetAdminApplicationsResponse struct {
	Applications []*AdminApplicationResponse `json:"applications"`
	Total        int64                       `json:"total"`
	Page         int                         `json:"page"`
	Limit        int                         `json:"limit"`
}

type LoanStatisticsResponse struct {
	TotalApplications   int64       `json:"total_applications"`
	MonthlyApplications int64       `json:"monthly_applications"`
	StatusStatistics    interface{} `json:"status_statistics"`
}

type CreateProductRequest struct {
	ProductName      string  `json:"product_name" binding:"required"`
	ProductCode      string  `json:"product_code" binding:"required"`
	ProductType      string  `json:"product_type" binding:"required"`
	Description      string  `json:"description"`
	MinAmount        int64   `json:"min_amount" binding:"required"`
	MaxAmount        int64   `json:"max_amount" binding:"required"`
	InterestRate     float64 `json:"interest_rate" binding:"required"`
	TermMonths       int     `json:"term_months" binding:"required"`
	RequiredAuth     string  `json:"required_auth"`
	IsActive         bool    `json:"is_active"`
	SortOrder        int     `json:"sort_order"`
	DifyWorkflowID   string  `json:"dify_workflow_id"`
	EligibleUserType string  `json:"eligible_user_type"`
}

type UpdateProductRequest struct {
	ProductName      string  `json:"product_name"`
	Description      string  `json:"description"`
	MinAmount        int64   `json:"min_amount"`
	MaxAmount        int64   `json:"max_amount"`
	InterestRate     float64 `json:"interest_rate"`
	TermMonths       int     `json:"term_months"`
	RequiredAuth     string  `json:"required_auth"`
	IsActive         bool    `json:"is_active"`
	SortOrder        int     `json:"sort_order"`
	DifyWorkflowID   string  `json:"dify_workflow_id"`
	EligibleUserType string  `json:"eligible_user_type"`
}

type ReturnRequest struct {
	ReturnReason string `json:"return_reason" binding:"required"`
	ReturnNote   string `json:"return_note"`
}

type DifyCallbackRequest struct {
	ApplicationID uint64                 `json:"application_id"`
	WorkflowType  string                 `json:"workflow_type"`
	ExecutionID   string                 `json:"execution_id"`
	Status        string                 `json:"status"`
	Result        map[string]interface{} `json:"result"`
	Error         string                 `json:"error"`
}

type StatisticsRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	GroupBy   string `json:"group_by"`
}

type LoanStatistics struct {
	TotalApplications     int64            `json:"total_applications"`
	PendingApplications   int64            `json:"pending_applications"`
	ApprovedApplications  int64            `json:"approved_applications"`
	RejectedApplications  int64            `json:"rejected_applications"`
	TotalLoanAmount       int64            `json:"total_loan_amount"`
	ApprovedLoanAmount    int64            `json:"approved_loan_amount"`
	ApplicationsByProduct map[string]int64 `json:"applications_by_product"`
	ApprovalRate          float64          `json:"approval_rate"`
}

type ReportRequest struct {
	ReportType string `json:"report_type"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Format     string `json:"format"`
}

type ApprovalReport struct {
	TotalApplications    int64            `json:"total_applications"`
	ApprovedApplications int64            `json:"approved_applications"`
	RejectedApplications int64            `json:"rejected_applications"`
	PendingApplications  int64            `json:"pending_applications"`
	ReportData           []ReportDataItem `json:"report_data"`
	GeneratedAt          time.Time        `json:"generated_at"`
}

type ReportDataItem struct {
	Date             string  `json:"date"`
	ApplicationCount int64   `json:"application_count"`
	ApprovalCount    int64   `json:"approval_count"`
	ApprovalRate     float64 `json:"approval_rate"`
}

// ==================== 农机租赁相关结构体 ====================

type RegisterMachineRequest struct {
	MachineName       string                 `json:"machine_name" binding:"required"`
	MachineType       string                 `json:"machine_type" binding:"required"`
	Brand             string                 `json:"brand" binding:"required"`
	Model             string                 `json:"model" binding:"required"`
	Specifications    map[string]interface{} `json:"specifications"`
	Description       string                 `json:"description"`
	Images            []string               `json:"images"`
	Province          string                 `json:"province" binding:"required"`
	City              string                 `json:"city" binding:"required"`
	County            string                 `json:"county" binding:"required"`
	Address           string                 `json:"address"`
	Longitude         float64                `json:"longitude"`
	Latitude          float64                `json:"latitude"`
	HourlyRate        int64                  `json:"hourly_rate"`
	DailyRate         int64                  `json:"daily_rate"`
	PerAcreRate       int64                  `json:"per_acre_rate"`
	DepositAmount     int64                  `json:"deposit_amount" binding:"required"`
	AvailableSchedule map[string]interface{} `json:"available_schedule"`
	MinRentalHours    int                    `json:"min_rental_hours"`
	MaxAdvanceDays    int                    `json:"max_advance_days"`
}

type RegisterMachineResponse struct {
	ID          uint64 `json:"id"`
	MachineCode string `json:"machine_code"`
	Status      string `json:"status"`
}

type MachineDetailResponse struct {
	Machine *model.Machine  `json:"machine"`
	Owner   *model.User     `json:"owner"`
	Ratings []MachineRating `json:"ratings"`
}

type MachineRating struct {
	ID        uint64    `json:"id"`
	Rating    float32   `json:"rating"`
	Comment   string    `json:"comment"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateMachineRequest struct {
	MachineName       string                 `json:"machine_name"`
	Description       string                 `json:"description"`
	Images            []string               `json:"images"`
	HourlyRate        int64                  `json:"hourly_rate"`
	DailyRate         int64                  `json:"daily_rate"`
	PerAcreRate       int64                  `json:"per_acre_rate"`
	DepositAmount     int64                  `json:"deposit_amount"`
	AvailableSchedule map[string]interface{} `json:"available_schedule"`
	Status            string                 `json:"status"`
}

type GetUserMachinesRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Status string `json:"status"`
}

type GetUserMachinesResponse struct {
	Machines []*model.Machine `json:"machines"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
}

type SearchMachinesRequest struct {
	MachineType string  `json:"machine_type"`
	Province    string  `json:"province"`
	City        string  `json:"city"`
	County      string  `json:"county"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	Radius      int     `json:"radius"`
	Page        int     `json:"page"`
	Limit       int     `json:"limit"`
}

type SearchMachinesResponse struct {
	Machines []*model.Machine `json:"machines"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
}

type GetAvailableMachinesRequest struct {
	MachineType string `json:"machine_type"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
}

type GetAvailableMachinesResponse struct {
	Machines []*model.Machine `json:"machines"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
}

type CreateRentalOrderRequest struct {
	MachineID      uint64  `json:"machine_id" binding:"required"`
	StartTime      string  `json:"start_time" binding:"required"`
	EndTime        string  `json:"end_time" binding:"required"`
	RentalLocation string  `json:"rental_location" binding:"required"`
	ContactPerson  string  `json:"contact_person" binding:"required"`
	ContactPhone   string  `json:"contact_phone" binding:"required"`
	BillingMethod  string  `json:"billing_method" binding:"required"`
	Quantity       float64 `json:"quantity" binding:"required"`
	Remarks        string  `json:"remarks"`
}

type CreateRentalOrderResponse struct {
	ID            uint64 `json:"id"`
	OrderNo       string `json:"order_no"`
	TotalAmount   int64  `json:"total_amount"`
	DepositAmount int64  `json:"deposit_amount"`
}

type RentalOrderDetailResponse struct {
	Order   *model.RentalOrder `json:"order"`
	Machine *model.Machine     `json:"machine"`
	Renter  *model.User        `json:"renter"`
	Owner   *model.User        `json:"owner"`
}

type ConfirmOrderRequest struct {
	ConfirmType string `json:"confirm_type" binding:"required"` // owner, renter
}

type CancelOrderRequest struct {
	CancelReason string `json:"cancel_reason" binding:"required"`
}

type GetUserOrdersRequest struct {
	OrderType string `json:"order_type"` // renter, owner
	Status    string `json:"status"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
}

type GetUserOrdersResponse struct {
	Orders []*model.RentalOrder `json:"orders"`
	Total  int64                `json:"total"`
	Page   int                  `json:"page"`
	Limit  int                  `json:"limit"`
}

type PayOrderRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required"`
}

type PayOrderResponse struct {
	PaymentID     string `json:"payment_id"`
	PaymentAmount int64  `json:"payment_amount"`
	PaymentURL    string `json:"payment_url,omitempty"`
}

type CompleteOrderRequest struct {
	ActualEndTime string  `json:"actual_end_time"`
	WorkHours     float64 `json:"work_hours"`
	WorkAcres     float64 `json:"work_acres"`
	Notes         string  `json:"notes"`
}

type SubmitRatingRequest struct {
	RatingType string  `json:"rating_type" binding:"required"` // renter, owner
	Rating     float32 `json:"rating" binding:"required"`
	Comment    string  `json:"comment"`
}

type MachineRatingsResponse struct {
	Ratings       []MachineRating `json:"ratings"`
	AverageRating float32         `json:"average_rating"`
	TotalCount    int64           `json:"total_count"`
}

// ==================== 内容管理相关结构体 ====================

type CreateArticleRequest struct {
	Title          string   `json:"title" binding:"required"`
	Subtitle       string   `json:"subtitle"`
	Content        string   `json:"content" binding:"required"`
	Summary        string   `json:"summary"`
	Category       string   `json:"category" binding:"required"`
	Tags           []string `json:"tags"`
	CoverImage     string   `json:"cover_image"`
	IsTop          bool     `json:"is_top"`
	IsFeatured     bool     `json:"is_featured"`
	SEOTitle       string   `json:"seo_title"`
	SEODescription string   `json:"seo_description"`
	SEOKeywords    string   `json:"seo_keywords"`
}

type CreateArticleResponse struct {
	ID     uint64 `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type ArticleDetailResponse struct {
	Article *model.Article `json:"article"`
	Author  *model.OAUser  `json:"author"`
}

type UpdateArticleRequest struct {
	Title          string   `json:"title"`
	Subtitle       string   `json:"subtitle"`
	Content        string   `json:"content"`
	Summary        string   `json:"summary"`
	Category       string   `json:"category"`
	Tags           []string `json:"tags"`
	CoverImage     string   `json:"cover_image"`
	IsTop          bool     `json:"is_top"`
	IsFeatured     bool     `json:"is_featured"`
	SEOTitle       string   `json:"seo_title"`
	SEODescription string   `json:"seo_description"`
	SEOKeywords    string   `json:"seo_keywords"`
}

type ListArticlesRequest struct {
	Category   string `json:"category"`
	Tag        string `json:"tag"`
	Keyword    string `json:"keyword"`
	Status     string `json:"status"`
	IsTop      *bool  `json:"is_top"`
	IsFeatured *bool  `json:"is_featured"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

type ListArticlesResponse struct {
	Articles []*model.Article `json:"articles"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
}

type FeaturedArticlesResponse struct {
	Articles []*model.Article `json:"articles"`
}

type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required"`
	DisplayName string  `json:"display_name" binding:"required"`
	Description string  `json:"description"`
	ParentID    *uint64 `json:"parent_id"`
	Icon        string  `json:"icon"`
	SortOrder   int     `json:"sort_order"`
}

type UpdateCategoryRequest struct {
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
	Status      string `json:"status"`
}

type CreateExpertRequest struct {
	Name            string   `json:"name" binding:"required"`
	Title           string   `json:"title"`
	Organization    string   `json:"organization"`
	Specialties     []string `json:"specialties"`
	Phone           string   `json:"phone"`
	Email           string   `json:"email"`
	WeChat          string   `json:"wechat"`
	Avatar          string   `json:"avatar"`
	Biography       string   `json:"biography"`
	ExperienceYears int      `json:"experience_years"`
	ServiceAreas    []string `json:"service_areas"`
}

type ExpertDetailResponse struct {
	Expert *model.Expert `json:"expert"`
}

type ListExpertsRequest struct {
	Specialty string `json:"specialty"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Status    string `json:"status"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
}

type ListExpertsResponse struct {
	Experts []*model.Expert `json:"experts"`
	Total   int64           `json:"total"`
	Page    int             `json:"page"`
	Limit   int             `json:"limit"`
}

type UpdateExpertRequest struct {
	Name            string   `json:"name"`
	Title           string   `json:"title"`
	Organization    string   `json:"organization"`
	Specialties     []string `json:"specialties"`
	Phone           string   `json:"phone"`
	Email           string   `json:"email"`
	WeChat          string   `json:"wechat"`
	Avatar          string   `json:"avatar"`
	Biography       string   `json:"biography"`
	ExperienceYears int      `json:"experience_years"`
	ServiceAreas    []string `json:"service_areas"`
	Status          string   `json:"status"`
}

type SubmitConsultationRequest struct {
	ExpertID    uint64   `json:"expert_id" binding:"required"`
	Question    string   `json:"question" binding:"required"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	ContactType string   `json:"contact_type" binding:"required"` // phone, wechat, email
}

type SubmitConsultationResponse struct {
	ID     uint64 `json:"id"`
	Status string `json:"status"`
}

type ConsultationsResponse struct {
	Consultations []ConsultationDetail `json:"consultations"`
	Total         int64                `json:"total"`
}

type ConsultationDetail struct {
	ID         uint64        `json:"id"`
	Expert     *model.Expert `json:"expert"`
	Question   string        `json:"question"`
	Answer     string        `json:"answer"`
	Status     string        `json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
	AnsweredAt *time.Time    `json:"answered_at"`
}

// ==================== 系统管理相关结构体 ====================

type UploadFileRequest struct {
	File         io.Reader `json:"-"`
	FileName     string    `json:"file_name" binding:"required"`
	BusinessType string    `json:"business_type"`
	BusinessID   uint64    `json:"business_id"`
	IsPublic     bool      `json:"is_public"`
}

type UploadFileResponse struct {
	ID       uint64 `json:"id"`
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
	FileSize int64  `json:"file_size"`
}

type HealthCheckResponse struct {
	Status    string                 `json:"status"`
	Database  map[string]interface{} `json:"database"`
	Redis     map[string]interface{} `json:"redis"`
	Services  map[string]interface{} `json:"services"`
	Timestamp int64                  `json:"timestamp"`
}

type SystemStatsResponse struct {
	UserCount        int64 `json:"user_count"`
	ApplicationCount int64 `json:"application_count"`
	MachineCount     int64 `json:"machine_count"`
	OrderCount       int64 `json:"order_count"`
	ArticleCount     int64 `json:"article_count"`
}

type SystemVersionResponse struct {
	Version     string `json:"version"`
	BuildTime   string `json:"build_time"`
	GitCommit   string `json:"git_commit"`
	GoVersion   string `json:"go_version"`
	Environment string `json:"environment"`
}

// ==================== OA管理相关结构体 ====================

type CreateOAUserRequest struct {
	Username   string `json:"username" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Phone      string `json:"phone"`
	Password   string `json:"password" binding:"required,min=6"`
	RealName   string `json:"real_name" binding:"required"`
	RoleID     uint64 `json:"role_id" binding:"required"`
	Department string `json:"department"`
	Position   string `json:"position"`
}

type UpdateOAUserRequest struct {
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	RealName   string `json:"real_name"`
	RoleID     uint64 `json:"role_id"`
	Department string `json:"department"`
	Position   string `json:"position"`
	Status     string `json:"status"`
}

type ListOAUsersRequest struct {
	RoleID     uint64 `json:"role_id"`
	Department string `json:"department"`
	Status     string `json:"status"`
	Keyword    string `json:"keyword"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

type ListOAUsersResponse struct {
	Users []*model.OAUser `json:"users"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}

type CreateRoleRequest struct {
	Name        string                 `json:"name" binding:"required"`
	DisplayName string                 `json:"display_name" binding:"required"`
	Description string                 `json:"description"`
	Permissions map[string]interface{} `json:"permissions"`
	IsSuper     bool                   `json:"is_super"`
}

type UpdateRoleRequest struct {
	DisplayName string                 `json:"display_name"`
	Description string                 `json:"description"`
	Permissions map[string]interface{} `json:"permissions"`
	Status      string                 `json:"status"`
}

type OALoginRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Platform   string `json:"platform" binding:"required"`
	CaptchaID  string `json:"captcha_id"`
	Captcha    string `json:"captcha"`
	DeviceType string `json:"device_type"`
	DeviceName string `json:"device_name"`
	DeviceID   string `json:"device_id"`
	AppVersion string `json:"app_version"`
	IPAddress  string `json:"ip_address"`
	Location   string `json:"location"`
}

type OALoginResponse struct {
	User         *model.OAUser `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int           `json:"expires_in"`
	SessionID    string        `json:"session_id"`
}

type DashboardResponse struct {
	UserStats        map[string]int64 `json:"user_stats"`
	LoanStats        map[string]int64 `json:"loan_stats"`
	MachineStats     map[string]int64 `json:"machine_stats"`
	RecentActivities []ActivityLog    `json:"recent_activities"`
}

type ActivityLog struct {
	ID          uint64    `json:"id"`
	UserName    string    `json:"user_name"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type PendingTasksResponse struct {
	LoanApprovals       []PendingLoanApproval `json:"loan_approvals"`
	UserAuthentications []PendingUserAuth     `json:"user_authentications"`
	TotalCount          int64                 `json:"total_count"`
}

type PendingLoanApproval struct {
	ID            uint64    `json:"id"`
	ApplicationNo string    `json:"application_no"`
	UserName      string    `json:"user_name"`
	ProductName   string    `json:"product_name"`
	Amount        int64     `json:"amount"`
	SubmittedAt   time.Time `json:"submitted_at"`
}

type PendingUserAuth struct {
	ID          uint64    `json:"id"`
	UserName    string    `json:"user_name"`
	AuthType    string    `json:"auth_type"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// ==================== 服务类型别名 ====================

// 为了兼容router中的引用，提供类型别名
type ArticleService = ContentService
type ExpertService = ContentService
type FileService = SystemService

// ==================== 任务管理相关结构体 ====================

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	Type         string  `json:"type" binding:"required"`
	Priority     string  `json:"priority"`
	BusinessID   uint64  `json:"business_id" binding:"required"`
	BusinessType string  `json:"business_type" binding:"required"`
	AssignedTo   *uint64 `json:"assigned_to"`
	Data         string  `json:"data"`
}

// CreateTaskResponse 创建任务响应
type CreateTaskResponse struct {
	ID     uint64 `json:"id"`
	TaskNo string `json:"task_no"`
	Status string `json:"status"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Priority    string  `json:"priority"`
	Status      string  `json:"status"`
	AssignedTo  *uint64 `json:"assigned_to"`
	Data        string  `json:"data"`
}

// ListTasksRequest 任务列表请求
type ListTasksRequest struct {
	Page         int     `json:"page"`
	Limit        int     `json:"limit"`
	Type         string  `json:"type"`
	Status       string  `json:"status"`
	Priority     string  `json:"priority"`
	AssignedTo   *uint64 `json:"assigned_to"`
	CreatedBy    *uint64 `json:"created_by"`
	BusinessType string  `json:"business_type"`
	BusinessID   *uint64 `json:"business_id"`
	IsOverdue    *bool   `json:"is_overdue"`
	Keyword      string  `json:"keyword"`
	SortBy       string  `json:"sort_by"`
	SortOrder    string  `json:"sort_order"`
}

// ListTasksResponse 任务列表响应
type ListTasksResponse struct {
	Tasks      []*TaskDetailResponse `json:"tasks"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	Statistics *TaskStatisticsInfo   `json:"statistics"`
}

// TaskDetailResponse 任务详情响应
type TaskDetailResponse struct {
	ID           uint64     `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Type         string     `json:"type"`
	TypeName     string     `json:"type_name"`
	Priority     string     `json:"priority"`
	PriorityName string     `json:"priority_name"`
	Status       string     `json:"status"`
	StatusName   string     `json:"status_name"`
	Progress     float64    `json:"progress"`
	IsOverdue    bool       `json:"is_overdue"`
	BusinessID   uint64     `json:"business_id"`
	BusinessType string     `json:"business_type"`
	Data         string     `json:"data"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CompletedAt  *time.Time `json:"completed_at"`

	// 关联信息
	AssignedUser *TaskUserInfo     `json:"assigned_user,omitempty"`
	CreatedUser  *TaskUserInfo     `json:"created_user,omitempty"`
	Actions      []*TaskActionInfo `json:"actions,omitempty"`

	// 业务相关信息
	BusinessInfo map[string]interface{} `json:"business_info,omitempty"`
}

// TaskUserInfo 任务用户信息
type TaskUserInfo struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
	Avatar   string `json:"avatar,omitempty"`
}

// TaskActionInfo 任务操作信息
type TaskActionInfo struct {
	ID         uint64        `json:"id"`
	Action     string        `json:"action"`
	ActionName string        `json:"action_name"`
	Comment    string        `json:"comment"`
	CreatedAt  time.Time     `json:"created_at"`
	Operator   *TaskUserInfo `json:"operator"`
}

// TaskStatisticsInfo 任务统计信息
type TaskStatisticsInfo struct {
	TotalTasks        int64 `json:"total_tasks"`
	PendingTasks      int64 `json:"pending_tasks"`
	ProcessingTasks   int64 `json:"processing_tasks"`
	CompletedTasks    int64 `json:"completed_tasks"`
	CancelledTasks    int64 `json:"cancelled_tasks"`
	OverdueTasks      int64 `json:"overdue_tasks"`
	HighPriorityTasks int64 `json:"high_priority_tasks"`
	UrgentTasks       int64 `json:"urgent_tasks"`
	MyTasks           int64 `json:"my_tasks"`
	UnassignedTasks   int64 `json:"unassigned_tasks"`
}

// TaskProgressResponse 任务进度响应
type TaskProgressResponse struct {
	TaskID    uint64        `json:"task_id"`
	Progress  float64       `json:"progress"`
	UpdatedAt time.Time     `json:"updated_at"`
	UpdatedBy *TaskUserInfo `json:"updated_by"`
}

// ProcessTaskRequest 处理任务请求
type ProcessTaskRequest struct {
	Action   string  `json:"action" binding:"required"`
	Comment  string  `json:"comment"`
	Progress float64 `json:"progress"`
}

// ProcessTaskResponse 处理任务响应
type ProcessTaskResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	TaskID    uint64    `json:"task_id"`
	NewStatus string    `json:"new_status"`
	UpdatedAt time.Time `json:"updated_at"`
}
