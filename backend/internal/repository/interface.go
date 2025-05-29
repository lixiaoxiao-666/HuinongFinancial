package repository

import (
	"context"
	"huinong-backend/internal/model"
	"time"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	// 基本CRUD
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint64) (*model.User, error)
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	GetByUUID(ctx context.Context, uuid string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint64) error

	// 查询方法
	List(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error)
	GetByUserType(ctx context.Context, userType string, limit, offset int) ([]*model.User, error)
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]*model.User, error)

	// 登录信息更新
	UpdateLoginInfo(ctx context.Context, userID uint64, loginIP string) error

	// 认证相关
	GetUserAuth(ctx context.Context, userID uint64, authType string) (*model.UserAuth, error)
	CreateUserAuth(ctx context.Context, auth *model.UserAuth) error
	UpdateUserAuth(ctx context.Context, auth *model.UserAuth) error

	// 会话管理
	CreateSession(ctx context.Context, session *model.UserSession) error
	GetSession(ctx context.Context, sessionID string) (*model.UserSession, error)
	UpdateSession(ctx context.Context, session *model.UserSession) error
	DeleteSession(ctx context.Context, sessionID string) error
	GetUserSessions(ctx context.Context, userID uint64) ([]*model.UserSession, error)

	// 标签管理
	AddUserTag(ctx context.Context, tag *model.UserTag) error
	GetUserTags(ctx context.Context, userID uint64, tagType string) ([]*model.UserTag, error)
	RemoveUserTag(ctx context.Context, userID uint64, tagKey string) error

	// 统计方法
	GetUserCount(ctx context.Context) (int64, error)
	GetUserCountByType(ctx context.Context, userType string) (int64, error)
}

// LoanRepository 贷款相关数据访问接口
type LoanRepository interface {
	// 贷款产品相关
	GetProductByID(ctx context.Context, id uint) (*model.LoanProduct, error)
	GetProductByCode(ctx context.Context, code string) (*model.LoanProduct, error)
	GetProductsByUserType(ctx context.Context, userType string, page, limit int) ([]*model.LoanProduct, int64, error)
	GetAllProducts(ctx context.Context) ([]*model.LoanProduct, error)
	GetActiveProducts(ctx context.Context, userType string) ([]*model.LoanProduct, error)
	CreateProduct(ctx context.Context, product *model.LoanProduct) error
	UpdateProduct(ctx context.Context, product *model.LoanProduct) error
	DeleteProduct(ctx context.Context, id uint) error
	ListProducts(ctx context.Context, req *ListProductsRequest) (*ListProductsResponse, error)

	// 贷款申请相关
	CreateApplication(ctx context.Context, application *model.LoanApplication) error
	GetApplicationByID(ctx context.Context, id uint) (*model.LoanApplication, error)
	GetApplicationByNo(ctx context.Context, applicationNo string) (*model.LoanApplication, error)
	GetUserApplications(ctx context.Context, userID uint, page, limit int, status string) ([]*model.LoanApplication, int64, error)
	GetPendingApplications(ctx context.Context, limit, offset int) ([]*model.LoanApplication, error)
	UpdateApplication(ctx context.Context, application *model.LoanApplication) error
	UpdateApplicationStatus(ctx context.Context, id uint, status string) error
	GetApplicationsForAdmin(ctx context.Context, page, limit int, status string) ([]*model.LoanApplication, int64, error)
	GetApplicationStatistics(ctx context.Context) (map[string]interface{}, error)

	// 审批日志相关
	CreateApprovalLog(ctx context.Context, log *model.ApprovalLog) error
	GetApprovalLogs(ctx context.Context, applicationID uint) ([]*model.ApprovalLog, error)

	// Dify工作流日志相关
	CreateDifyLog(ctx context.Context, log *model.DifyWorkflowLog) error
	GetDifyLogs(ctx context.Context, applicationID uint) ([]*model.DifyWorkflowLog, error)
}

// MachineRepository 农机数据访问接口
type MachineRepository interface {
	// 设备管理
	Create(ctx context.Context, machine *model.Machine) error
	GetByID(ctx context.Context, id uint64) (*model.Machine, error)
	GetByCode(ctx context.Context, code string) (*model.Machine, error)
	Update(ctx context.Context, machine *model.Machine) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, req *ListMachinesRequest) (*ListMachinesResponse, error)

	// 地理位置搜索
	SearchNearby(ctx context.Context, longitude, latitude, radius float64, machineType string) ([]*model.Machine, error)
	GetByLocation(ctx context.Context, province, city, county string) ([]*model.Machine, error)

	// 租赁订单
	CreateOrder(ctx context.Context, order *model.RentalOrder) error
	GetOrderByID(ctx context.Context, id uint64) (*model.RentalOrder, error)
	GetOrderByNo(ctx context.Context, orderNo string) (*model.RentalOrder, error)
	UpdateOrder(ctx context.Context, order *model.RentalOrder) error
	ListOrders(ctx context.Context, req *ListOrdersRequest) (*ListOrdersResponse, error)
	GetUserOrders(ctx context.Context, userID uint64, userType string, limit, offset int) ([]*model.RentalOrder, error)

	// 统计方法
	GetMachineCount(ctx context.Context) (int64, error)
	GetAvailableMachineCount(ctx context.Context) (int64, error)
	GetOrderCount(ctx context.Context) (int64, error)
}

// ArticleRepository 文章数据访问接口
type ArticleRepository interface {
	// 基本CRUD
	Create(ctx context.Context, article *model.Article) error
	GetByID(ctx context.Context, id uint64) (*model.Article, error)
	Update(ctx context.Context, article *model.Article) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, req *ListArticlesRequest) (*ListArticlesResponse, error)

	// 分类管理
	CreateCategory(ctx context.Context, category *model.Category) error
	GetCategoryByID(ctx context.Context, id uint64) (*model.Category, error)
	GetCategoryByName(ctx context.Context, name string) (*model.Category, error)
	UpdateCategory(ctx context.Context, category *model.Category) error
	DeleteCategory(ctx context.Context, id uint64) error
	ListCategories(ctx context.Context) ([]*model.Category, error)

	// 内容查询
	GetByCategory(ctx context.Context, category string, limit, offset int) ([]*model.Article, error)
	GetFeatured(ctx context.Context, limit int) ([]*model.Article, error)
	GetTopArticles(ctx context.Context, limit int) ([]*model.Article, error)
	Search(ctx context.Context, keyword string, limit, offset int) ([]*model.Article, error)

	// 统计更新
	IncrementViewCount(ctx context.Context, id uint64) error
	IncrementLikeCount(ctx context.Context, id uint64) error
	IncrementShareCount(ctx context.Context, id uint64) error
}

// ExpertRepository 专家数据访问接口
type ExpertRepository interface {
	Create(ctx context.Context, expert *model.Expert) error
	GetByID(ctx context.Context, id uint64) (*model.Expert, error)
	Update(ctx context.Context, expert *model.Expert) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, req *ListExpertsRequest) (*ListExpertsResponse, error)
	SearchBySpecialty(ctx context.Context, specialties []string, serviceArea string, limit, offset int) ([]*model.Expert, error)
	GetVerifiedExperts(ctx context.Context, limit, offset int) ([]*model.Expert, error)
}

// SystemConfigRepository 系统配置数据访问接口
type SystemConfigRepository interface {
	Get(ctx context.Context, key string) (*model.SystemConfig, error)
	Set(ctx context.Context, config *model.SystemConfig) error
	GetByGroup(ctx context.Context, group string) ([]*model.SystemConfig, error)
	List(ctx context.Context) ([]*model.SystemConfig, error)
	Delete(ctx context.Context, key string) error
}

// FileRepository 文件数据访问接口
type FileRepository interface {
	Create(ctx context.Context, file *model.FileUpload) error
	GetByID(ctx context.Context, id uint64) (*model.FileUpload, error)
	GetByHash(ctx context.Context, hash string) (*model.FileUpload, error)
	GetByBusiness(ctx context.Context, businessType string, businessID uint64) ([]*model.FileUpload, error)
	Update(ctx context.Context, file *model.FileUpload) error
	Delete(ctx context.Context, id uint64) error
	IncrementAccessCount(ctx context.Context, id uint64) error
}

// OfflineQueueRepository 离线队列数据访问接口
type OfflineQueueRepository interface {
	Add(ctx context.Context, action *model.OfflineQueue) error
	GetPending(ctx context.Context, limit int) ([]*model.OfflineQueue, error)
	GetUserActions(ctx context.Context, userID uint64, limit, offset int) ([]*model.OfflineQueue, error)
	Update(ctx context.Context, action *model.OfflineQueue) error
	Delete(ctx context.Context, id uint64) error
	GetRetryableActions(ctx context.Context) ([]*model.OfflineQueue, error)
}

// APILogRepository API日志数据访问接口
type APILogRepository interface {
	Create(ctx context.Context, log *model.APILog) error
	GetByRequestID(ctx context.Context, requestID string) (*model.APILog, error)
	List(ctx context.Context, req *ListAPILogsRequest) (*ListAPILogsResponse, error)
	DeleteOldLogs(ctx context.Context, before time.Time) error
	GetStatistics(ctx context.Context, startTime, endTime time.Time) (*APIStatistics, error)
}

// OARepository OA后台数据访问接口
type OARepository interface {
	// OA用户管理
	CreateOAUser(ctx context.Context, user *model.OAUser) error
	GetOAUserByID(ctx context.Context, id uint64) (*model.OAUser, error)
	GetOAUserByUsername(ctx context.Context, username string) (*model.OAUser, error)
	UpdateOAUser(ctx context.Context, user *model.OAUser) error
	DeleteOAUser(ctx context.Context, id uint64) error
	ListOAUsers(ctx context.Context, req *ListOAUsersRequest) (*ListOAUsersResponse, error)

	// 角色管理
	CreateRole(ctx context.Context, role *model.OARole) error
	GetRoleByID(ctx context.Context, id uint64) (*model.OARole, error)
	GetRoleByName(ctx context.Context, name string) (*model.OARole, error)
	UpdateRole(ctx context.Context, role *model.OARole) error
	DeleteRole(ctx context.Context, id uint64) error
	ListRoles(ctx context.Context) ([]*model.OARole, error)
}

// 请求和响应结构体
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

type ListProductsRequest struct {
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
	ProductType string `json:"product_type"`
	Status      string `json:"status"`
	Keyword     string `json:"keyword"`
}

type ListProductsResponse struct {
	Products []*model.LoanProduct `json:"products"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
}

type ListApplicationsRequest struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	UserID    uint64 `json:"user_id"`
	ProductID uint64 `json:"product_id"`
	Status    string `json:"status"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type ListApplicationsResponse struct {
	Applications []*model.LoanApplication `json:"applications"`
	Total        int64                    `json:"total"`
	Page         int                      `json:"page"`
	Limit        int                      `json:"limit"`
}

type ListMachinesRequest struct {
	Page        int     `json:"page"`
	Limit       int     `json:"limit"`
	MachineType string  `json:"machine_type"`
	Status      string  `json:"status"`
	Province    string  `json:"province"`
	City        string  `json:"city"`
	County      string  `json:"county"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	Radius      float64 `json:"radius"`
	Keyword     string  `json:"keyword"`
}

type ListMachinesResponse struct {
	Machines []*model.Machine `json:"machines"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
}

type ListOrdersRequest struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	MachineID uint64 `json:"machine_id"`
	RenterID  uint64 `json:"renter_id"`
	OwnerID   uint64 `json:"owner_id"`
	Status    string `json:"status"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type ListOrdersResponse struct {
	Orders []*model.RentalOrder `json:"orders"`
	Total  int64                `json:"total"`
	Page   int                  `json:"page"`
	Limit  int                  `json:"limit"`
}

type ListArticlesRequest struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Category   string `json:"category"`
	Status     string `json:"status"`
	AuthorID   uint64 `json:"author_id"`
	IsTop      *bool  `json:"is_top"`
	IsFeatured *bool  `json:"is_featured"`
	Keyword    string `json:"keyword"`
}

type ListArticlesResponse struct {
	Articles []*model.Article `json:"articles"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
}

type ListExpertsRequest struct {
	Page        int      `json:"page"`
	Limit       int      `json:"limit"`
	Specialties []string `json:"specialties"`
	ServiceArea string   `json:"service_area"`
	IsVerified  *bool    `json:"is_verified"`
	Status      string   `json:"status"`
	Keyword     string   `json:"keyword"`
}

type ListExpertsResponse struct {
	Experts []*model.Expert `json:"experts"`
	Total   int64           `json:"total"`
	Page    int             `json:"page"`
	Limit   int             `json:"limit"`
}

type ListAPILogsRequest struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Method     string `json:"method"`
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	UserID     uint64 `json:"user_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
}

type ListAPILogsResponse struct {
	Logs  []*model.APILog `json:"logs"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}

type ListOAUsersRequest struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	RoleID     uint64 `json:"role_id"`
	Department string `json:"department"`
	Status     string `json:"status"`
	Keyword    string `json:"keyword"`
}

type ListOAUsersResponse struct {
	Users []*model.OAUser `json:"users"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}

type APIStatistics struct {
	TotalRequests   int64          `json:"total_requests"`
	SuccessRequests int64          `json:"success_requests"`
	ErrorRequests   int64          `json:"error_requests"`
	AvgResponseTime float64        `json:"avg_response_time"`
	TopEndpoints    []EndpointStat `json:"top_endpoints"`
}

type EndpointStat struct {
	Endpoint string `json:"endpoint"`
	Count    int64  `json:"count"`
}
