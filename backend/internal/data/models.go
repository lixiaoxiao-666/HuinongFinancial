package data

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User 用户表
type User struct {
	ID           uint       `gorm:"primarykey;autoIncrement" json:"id"`
	UserID       string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"userId"`
	Phone        string     `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`
	Nickname     string     `gorm:"type:varchar(100)" json:"nickname"`
	AvatarURL    string     `gorm:"type:varchar(512)" json:"avatarUrl"`
	Status       int8       `gorm:"type:tinyint;not null;default:0" json:"status"`
	RegisteredAt time.Time  `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"registeredAt"`
	LastLoginAt  *time.Time `gorm:"type:datetime(3)" json:"lastLoginAt"`
	CreatedAt    time.Time  `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// UserProfile 用户画像/详情表
type UserProfile struct {
	UserID           string     `gorm:"type:varchar(64);primaryKey" json:"userId"`
	RealName         string     `gorm:"type:varchar(100)" json:"realName"`
	IDCardNumber     string     `gorm:"type:varchar(30);index" json:"idCardNumber"`
	IDCardFrontURL   string     `gorm:"type:varchar(512)" json:"idCardFrontUrl"`
	IDCardBackURL    string     `gorm:"type:varchar(512)" json:"idCardBackUrl"`
	Address          string     `gorm:"type:varchar(255)" json:"address"`
	Gender           int8       `gorm:"type:tinyint" json:"gender"`
	BirthDate        *time.Time `gorm:"type:date" json:"birthDate"`
	Occupation       string     `gorm:"type:varchar(100)" json:"occupation"`
	AnnualIncome     float64    `gorm:"type:decimal(15,2)" json:"annualIncome"`
	CreditAuthAgreed bool       `gorm:"type:boolean;not null;default:0" json:"creditAuthAgreed"`
	CreatedAt        time.Time  `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// LoanProduct 贷款产品表
type LoanProduct struct {
	ID                    uint            `gorm:"primarykey;autoIncrement" json:"id"`
	ProductID             string          `gorm:"type:varchar(64);uniqueIndex;not null" json:"productId"`
	Name                  string          `gorm:"type:varchar(255);not null" json:"name"`
	Description           string          `gorm:"type:text" json:"description"`
	Category              string          `gorm:"type:varchar(50);index" json:"category"`
	MinAmount             float64         `gorm:"type:decimal(15,2);not null" json:"minAmount"`
	MaxAmount             float64         `gorm:"type:decimal(15,2);not null" json:"maxAmount"`
	MinTermMonths         int             `gorm:"type:int;not null" json:"minTermMonths"`
	MaxTermMonths         int             `gorm:"type:int;not null" json:"maxTermMonths"`
	InterestRateYearly    string          `gorm:"type:varchar(50);not null" json:"interestRateYearly"`
	RepaymentMethods      json.RawMessage `gorm:"type:json" json:"repaymentMethods"`
	ApplicationConditions string          `gorm:"type:text" json:"applicationConditions"`
	RequiredDocuments     json.RawMessage `gorm:"type:json" json:"requiredDocuments"`
	Status                int8            `gorm:"type:tinyint;not null;default:0;index" json:"status"`
	CreatedAt             time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt             time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// LoanApplication 贷款申请表
type LoanApplication struct {
	ID                 uint            `gorm:"primarykey;autoIncrement" json:"id"`
	ApplicationID      string          `gorm:"type:varchar(64);uniqueIndex;not null" json:"applicationId"`
	UserID             string          `gorm:"type:varchar(64);not null;index" json:"userId"`
	ProductID          string          `gorm:"type:varchar(64);not null;index" json:"productId"`
	AmountApplied      float64         `gorm:"type:decimal(15,2);not null" json:"amountApplied"`
	TermMonthsApplied  int             `gorm:"type:int;not null" json:"termMonthsApplied"`
	Purpose            string          `gorm:"type:varchar(500)" json:"purpose"`
	Status             string          `gorm:"type:varchar(50);not null;index" json:"status"`
	ApplicantSnapshot  json.RawMessage `gorm:"type:json" json:"applicantSnapshot"`
	SubmittedAt        time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"submittedAt"`
	AIRiskScore        *int            `gorm:"type:int" json:"aiRiskScore"`
	AISuggestion       string          `gorm:"type:text" json:"aiSuggestion"`
	ApprovedAmount     *float64        `gorm:"type:decimal(15,2)" json:"approvedAmount"`
	ApprovedTermMonths *int            `gorm:"type:int" json:"approvedTermMonths"`
	FinalDecision      string          `gorm:"type:varchar(50)" json:"finalDecision"`
	DecisionReason     string          `gorm:"type:text" json:"decisionReason"`
	ProcessedBy        string          `gorm:"type:varchar(64)" json:"processedBy"`
	ProcessedAt        *time.Time      `gorm:"type:datetime(3)" json:"processedAt"`
	CreatedAt          time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt          time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// LoanApplicationHistory 贷款申请审批历史表
type LoanApplicationHistory struct {
	ID            uint      `gorm:"primarykey;autoIncrement" json:"id"`
	ApplicationID string    `gorm:"type:varchar(64);not null;index" json:"applicationId"`
	StatusFrom    string    `gorm:"type:varchar(50)" json:"statusFrom"`
	StatusTo      string    `gorm:"type:varchar(50);not null" json:"statusTo"`
	OperatorType  string    `gorm:"type:varchar(20);not null" json:"operatorType"`
	OperatorID    string    `gorm:"type:varchar(64)" json:"operatorId"`
	Comments      string    `gorm:"type:text" json:"comments"`
	OccurredAt    time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"occurredAt"`
}

// UploadedFile 上传文件记录表
type UploadedFile struct {
	ID          uint      `gorm:"primarykey;autoIncrement" json:"id"`
	FileID      string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"fileId"`
	UserID      string    `gorm:"type:varchar(64);not null;index" json:"userId"`
	FileName    string    `gorm:"type:varchar(255);not null" json:"fileName"`
	FileType    string    `gorm:"type:varchar(50)" json:"fileType"`
	FileSize    int64     `gorm:"type:bigint" json:"fileSize"`
	StoragePath string    `gorm:"type:varchar(512);not null" json:"storagePath"`
	Purpose     string    `gorm:"type:varchar(100);index" json:"purpose"`
	RelatedID   string    `gorm:"type:varchar(64);index" json:"relatedId"`
	UploadedAt  time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"uploadedAt"`
	CreatedAt   time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// FarmMachinery 农机信息表
type FarmMachinery struct {
	ID           uint            `gorm:"primarykey;autoIncrement" json:"id"`
	MachineryID  string          `gorm:"type:varchar(64);uniqueIndex;not null" json:"machineryId"`
	OwnerUserID  string          `gorm:"type:varchar(64);not null;index" json:"ownerUserId"`
	Type         string          `gorm:"type:varchar(100);not null;index" json:"type"`
	BrandModel   string          `gorm:"type:varchar(255)" json:"brandModel"`
	Description  string          `gorm:"type:text" json:"description"`
	Images       json.RawMessage `gorm:"type:json" json:"images"`
	DailyRent    float64         `gorm:"type:decimal(10,2);not null" json:"dailyRent"`
	Deposit      *float64        `gorm:"type:decimal(10,2)" json:"deposit"`
	LocationText string          `gorm:"type:varchar(255)" json:"locationText"`
	LocationGeo  string          `gorm:"type:varchar(100)" json:"locationGeo"`
	Status       string          `gorm:"type:varchar(50);not null;default:'AVAILABLE'" json:"status"`
	PublishedAt  *time.Time      `gorm:"type:datetime(3)" json:"publishedAt"`
	CreatedAt    time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt    time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// MachineryLeasingOrder 农机租赁订单表
type MachineryLeasingOrder struct {
	ID            uint       `gorm:"primarykey;autoIncrement" json:"id"`
	OrderID       string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"orderId"`
	MachineryID   string     `gorm:"type:varchar(64);not null;index" json:"machineryId"`
	LesseeUserID  string     `gorm:"type:varchar(64);not null;index" json:"lesseeUserId"`
	LessorUserID  string     `gorm:"type:varchar(64);not null;index" json:"lessorUserId"`
	StartDate     time.Time  `gorm:"type:date;not null" json:"startDate"`
	EndDate       time.Time  `gorm:"type:date;not null" json:"endDate"`
	TotalRent     float64    `gorm:"type:decimal(10,2);not null" json:"totalRent"`
	DepositAmount *float64   `gorm:"type:decimal(10,2)" json:"depositAmount"`
	Status        string     `gorm:"type:varchar(50);not null;index" json:"status"`
	LesseeNotes   string     `gorm:"type:text" json:"lesseeNotes"`
	LessorNotes   string     `gorm:"type:text" json:"lessorNotes"`
	ConfirmedAt   *time.Time `gorm:"type:datetime(3)" json:"confirmedAt"`
	CompletedAt   *time.Time `gorm:"type:datetime(3)" json:"completedAt"`
	CreatedAt     time.Time  `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// OAUser OA后台用户表
type OAUser struct {
	ID           uint      `gorm:"primarykey;autoIncrement" json:"id"`
	OAUserID     string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"oaUserId"`
	Username     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         string    `gorm:"type:varchar(50);not null;index" json:"role"`
	DisplayName  string    `gorm:"type:varchar(100)" json:"displayName"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	Status       int8      `gorm:"type:tinyint;not null;default:0" json:"status"`
	CreatedAt    time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// SystemConfiguration 系统配置表
type SystemConfiguration struct {
	ConfigKey   string    `gorm:"type:varchar(100);primaryKey" json:"configKey"`
	ConfigValue string    `gorm:"type:text;not null" json:"configValue"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	UpdatedAt   time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// AIAnalysisResult AI分析结果表
type AIAnalysisResult struct {
	ID                    uint            `gorm:"primarykey;autoIncrement" json:"id"`
	ApplicationID         string          `gorm:"type:varchar(64);not null;index" json:"applicationId"`
	WorkflowExecutionID   string          `gorm:"type:varchar(64);uniqueIndex" json:"workflowExecutionId"`
	RiskLevel             string          `gorm:"type:varchar(20);not null" json:"riskLevel"`
	RiskScore             float64         `gorm:"type:decimal(5,4);not null" json:"riskScore"`
	ConfidenceScore       float64         `gorm:"type:decimal(5,4);not null" json:"confidenceScore"`
	AnalysisSummary       string          `gorm:"type:text" json:"analysisSummary"`
	DetailedAnalysis      json.RawMessage `gorm:"type:json" json:"detailedAnalysis"`
	RiskFactors           json.RawMessage `gorm:"type:json" json:"riskFactors"`
	Recommendations       json.RawMessage `gorm:"type:json" json:"recommendations"`
	AIDecision            string          `gorm:"type:varchar(50);not null" json:"aiDecision"`
	ApprovedAmount        *float64        `gorm:"type:decimal(15,2)" json:"approvedAmount"`
	ApprovedTermMonths    *int            `gorm:"type:int" json:"approvedTermMonths"`
	SuggestedInterestRate string          `gorm:"type:varchar(20)" json:"suggestedInterestRate"`
	NextAction            string          `gorm:"type:varchar(50);not null" json:"nextAction"`
	AIModelVersion        string          `gorm:"type:varchar(50);not null" json:"aiModelVersion"`
	ProcessingTimeMs      int             `gorm:"type:int;not null" json:"processingTimeMs"`
	ProcessedAt           time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"processedAt"`
	CreatedAt             time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt             time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// WorkflowExecution 工作流执行记录表
type WorkflowExecution struct {
	ID                  uint            `gorm:"primarykey;autoIncrement" json:"id"`
	ExecutionID         string          `gorm:"type:varchar(64);uniqueIndex;not null" json:"executionId"`
	ApplicationID       string          `gorm:"type:varchar(64);not null;index" json:"applicationId"`
	WorkflowType        string          `gorm:"type:varchar(50);not null" json:"workflowType"`
	Status              string          `gorm:"type:varchar(50);not null;index" json:"status"`
	Priority            string          `gorm:"type:varchar(20);not null;default:'NORMAL'" json:"priority"`
	StartedAt           time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"startedAt"`
	CompletedAt         *time.Time      `gorm:"type:datetime(3)" json:"completedAt"`
	EstimatedCompletion *time.Time      `gorm:"type:datetime(3)" json:"estimatedCompletion"`
	CurrentStage        string          `gorm:"type:varchar(100)" json:"currentStage"`
	Progress            int             `gorm:"type:int;not null;default:0" json:"progress"`
	ErrorMessage        string          `gorm:"type:text" json:"errorMessage"`
	Metadata            json.RawMessage `gorm:"type:json" json:"metadata"`
	CallbackURL         string          `gorm:"type:varchar(512)" json:"callbackUrl"`
	RetryCount          int             `gorm:"type:int;not null;default:0" json:"retryCount"`
	CreatedAt           time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt           time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// AIModelConfig AI模型配置表
type AIModelConfig struct {
	ID            uint            `gorm:"primarykey;autoIncrement" json:"id"`
	ModelID       string          `gorm:"type:varchar(100);uniqueIndex;not null" json:"modelId"`
	ModelType     string          `gorm:"type:varchar(50);not null;index" json:"modelType"`
	Version       string          `gorm:"type:varchar(50);not null" json:"version"`
	Status        string          `gorm:"type:varchar(20);not null;default:'ACTIVE'" json:"status"`
	Configuration json.RawMessage `gorm:"type:json" json:"configuration"`
	Thresholds    json.RawMessage `gorm:"type:json" json:"thresholds"`
	Description   string          `gorm:"type:text" json:"description"`
	CreatedBy     string          `gorm:"type:varchar(64);not null" json:"createdBy"`
	ActivatedAt   *time.Time      `gorm:"type:datetime(3)" json:"activatedAt"`
	DeactivatedAt *time.Time      `gorm:"type:datetime(3)" json:"deactivatedAt"`
	CreatedAt     time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt     time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// ExternalDataQuery 外部数据查询记录表
type ExternalDataQuery struct {
	ID            uint            `gorm:"primarykey;autoIncrement" json:"id"`
	QueryID       string          `gorm:"type:varchar(64);uniqueIndex;not null" json:"queryId"`
	UserID        string          `gorm:"type:varchar(64);not null;index" json:"userId"`
	ApplicationID string          `gorm:"type:varchar(64);index" json:"applicationId"`
	DataTypes     string          `gorm:"type:varchar(200);not null" json:"dataTypes"`
	QueryResult   json.RawMessage `gorm:"type:json" json:"queryResult"`
	Status        string          `gorm:"type:varchar(20);not null" json:"status"`
	ErrorMessage  string          `gorm:"type:text" json:"errorMessage"`
	QueryDuration int             `gorm:"type:int" json:"queryDuration"`
	QueriedAt     time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"queriedAt"`
	CreatedAt     time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt     time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// AIAgentLog AI智能体操作日志表
type AIAgentLog struct {
	ID            uint            `gorm:"primarykey;autoIncrement" json:"id"`
	LogID         string          `gorm:"type:varchar(64);uniqueIndex;not null" json:"logId"`
	ApplicationID string          `gorm:"type:varchar(64);index" json:"applicationId"`
	ActionType    string          `gorm:"type:varchar(50);not null;index" json:"actionType"`
	AgentType     string          `gorm:"type:varchar(50);not null" json:"agentType"`
	RequestData   json.RawMessage `gorm:"type:json" json:"requestData"`
	ResponseData  json.RawMessage `gorm:"type:json" json:"responseData"`
	Status        string          `gorm:"type:varchar(20);not null" json:"status"`
	ErrorMessage  string          `gorm:"type:text" json:"errorMessage"`
	Duration      int             `gorm:"type:int" json:"duration"`
	IPAddress     string          `gorm:"type:varchar(45)" json:"ipAddress"`
	UserAgent     string          `gorm:"type:text" json:"userAgent"`
	OccurredAt    time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"occurredAt"`
	CreatedAt     time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
}

// MachineryLeasingApplication 农机租赁审批申请表
type MachineryLeasingApplication struct {
	ID                 uint       `gorm:"primarykey;autoIncrement" json:"id"`
	ApplicationID      string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"applicationId"`
	LesseeUserID       string     `gorm:"type:varchar(64);not null;index" json:"lesseeUserId"`
	LessorUserID       string     `gorm:"type:varchar(64);not null;index" json:"lessorUserId"`
	MachineryID        string     `gorm:"type:varchar(64);not null;index" json:"machineryId"`
	RequestedStartDate time.Time  `gorm:"type:date;not null" json:"requestedStartDate"`
	RequestedEndDate   time.Time  `gorm:"type:date;not null" json:"requestedEndDate"`
	RentalDays         int        `gorm:"type:int;not null" json:"rentalDays"`
	TotalAmount        float64    `gorm:"type:decimal(10,2);not null" json:"totalAmount"`
	DepositAmount      *float64   `gorm:"type:decimal(10,2)" json:"depositAmount"`
	UsagePurpose       string     `gorm:"type:varchar(500)" json:"usagePurpose"`
	LesseeNotes        string     `gorm:"type:text" json:"lesseeNotes"`
	ApplicationStatus  string     `gorm:"type:varchar(50);not null;default:'PENDING_REVIEW';index" json:"applicationStatus"`
	RiskLevel          string     `gorm:"type:varchar(20)" json:"riskLevel"`
	AIRiskScore        *float64   `gorm:"type:decimal(5,4)" json:"aiRiskScore"`
	AISuggestion       string     `gorm:"type:text" json:"aiSuggestion"`
	ApprovalResult     string     `gorm:"type:varchar(50)" json:"approvalResult"`
	ApprovalComments   string     `gorm:"type:text" json:"approvalComments"`
	ApprovedBy         string     `gorm:"type:varchar(64)" json:"approvedBy"`
	ApprovedAt         *time.Time `gorm:"type:datetime(3)" json:"approvedAt"`
	SubmittedAt        time.Time  `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"submittedAt"`
	CreatedAt          time.Time  `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt          time.Time  `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// MachineryLeasingApprovalHistory 农机租赁审批历史记录表
type MachineryLeasingApprovalHistory struct {
	ID            uint      `gorm:"primarykey;autoIncrement" json:"id"`
	ApplicationID string    `gorm:"type:varchar(64);not null;index" json:"applicationId"`
	StatusFrom    string    `gorm:"type:varchar(50)" json:"statusFrom"`
	StatusTo      string    `gorm:"type:varchar(50);not null" json:"statusTo"`
	OperatorType  string    `gorm:"type:varchar(20);not null" json:"operatorType"`
	OperatorID    string    `gorm:"type:varchar(64)" json:"operatorId"`
	Comments      string    `gorm:"type:text" json:"comments"`
	OccurredAt    time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"occurredAt"`
}

// LessorQualification 农机主资质审核表
type LessorQualification struct {
	ID                     uint       `gorm:"primarykey;autoIncrement" json:"id"`
	UserID                 string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"userId"`
	RealName               string     `gorm:"type:varchar(100);not null" json:"realName"`
	IDCardNumber           string     `gorm:"type:varchar(30);not null" json:"idCardNumber"`
	Phone                  string     `gorm:"type:varchar(20);not null" json:"phone"`
	Address                string     `gorm:"type:varchar(255)" json:"address"`
	BusinessLicenseNumber  string     `gorm:"type:varchar(100)" json:"businessLicenseNumber"`
	FarmScale              string     `gorm:"type:varchar(100)" json:"farmScale"`
	VerificationStatus     string     `gorm:"type:varchar(20);not null;default:'PENDING';index" json:"verificationStatus"`
	CreditRating           string     `gorm:"type:varchar(10)" json:"creditRating"`
	TotalMachineryCount    int        `gorm:"type:int;not null;default:0" json:"totalMachineryCount"`
	SuccessfulLeasingCount int        `gorm:"type:int;not null;default:0" json:"successfulLeasingCount"`
	AverageRating          *float64   `gorm:"type:decimal(3,2)" json:"averageRating"`
	VerifiedAt             *time.Time `gorm:"type:datetime(3)" json:"verifiedAt"`
	CreatedAt              time.Time  `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt              time.Time  `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// MachineryLeasingAIResult 农机租赁AI分析结果表
type MachineryLeasingAIResult struct {
	ID                   uint            `gorm:"primarykey;autoIncrement" json:"id"`
	ApplicationID        string          `gorm:"type:varchar(64);not null;index" json:"applicationId"`
	WorkflowExecutionID  string          `gorm:"type:varchar(64);uniqueIndex" json:"workflowExecutionId"`
	RiskLevel            string          `gorm:"type:varchar(20);not null" json:"riskLevel"`
	RiskScore            float64         `gorm:"type:decimal(5,4);not null" json:"riskScore"`
	ConfidenceScore      float64         `gorm:"type:decimal(5,4);not null" json:"confidenceScore"`
	AnalysisSummary      string          `gorm:"type:text" json:"analysisSummary"`
	DetailedAnalysis     json.RawMessage `gorm:"type:json" json:"detailedAnalysis"`
	RiskFactors          json.RawMessage `gorm:"type:json" json:"riskFactors"`
	Recommendations      json.RawMessage `gorm:"type:json" json:"recommendations"`
	AIDecision           string          `gorm:"type:varchar(50);not null" json:"aiDecision"`
	SuggestedDepositRate *float64        `gorm:"type:decimal(5,4)" json:"suggestedDepositRate"`
	SuggestedConditions  json.RawMessage `gorm:"type:json" json:"suggestedConditions"`
	NextAction           string          `gorm:"type:varchar(50);not null" json:"nextAction"`
	AIModelVersion       string          `gorm:"type:varchar(50);not null" json:"aiModelVersion"`
	ProcessingTimeMs     int             `gorm:"type:int;not null" json:"processingTimeMs"`
	ProcessedAt          time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"processedAt"`
	CreatedAt            time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt            time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// TableName methods for machinery leasing approval models
func (MachineryLeasingApplication) TableName() string { return "machinery_leasing_applications" }
func (MachineryLeasingApprovalHistory) TableName() string {
	return "machinery_leasing_approval_history"
}
func (LessorQualification) TableName() string      { return "lessor_qualifications" }
func (MachineryLeasingAIResult) TableName() string { return "machinery_leasing_ai_results" }

// TableName methods for existing models
func (User) TableName() string                   { return "users" }
func (UserProfile) TableName() string            { return "user_profiles" }
func (LoanProduct) TableName() string            { return "loan_products" }
func (LoanApplication) TableName() string        { return "loan_applications" }
func (LoanApplicationHistory) TableName() string { return "loan_application_history" }
func (UploadedFile) TableName() string           { return "uploaded_files" }
func (FarmMachinery) TableName() string          { return "farm_machinery" }
func (MachineryLeasingOrder) TableName() string  { return "machinery_leasing_orders" }
func (OAUser) TableName() string                 { return "oa_users" }
func (SystemConfiguration) TableName() string    { return "system_configurations" }
func (AIAnalysisResult) TableName() string       { return "ai_analysis_results" }
func (WorkflowExecution) TableName() string      { return "workflow_executions" }
func (AIModelConfig) TableName() string          { return "ai_model_configs" }
func (ExternalDataQuery) TableName() string      { return "external_data_queries" }
func (AIAgentLog) TableName() string             { return "ai_agent_logs" }

// FileUpload 文件上传记录表
type FileUpload struct {
	ID           uint      `gorm:"primarykey;autoIncrement" json:"id"`
	FileID       string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"fileId"`
	UserID       string    `gorm:"type:varchar(64);not null;index" json:"userId"`
	FileName     string    `gorm:"type:varchar(255);not null" json:"fileName"`
	FileSize     int64     `gorm:"type:bigint" json:"fileSize"`
	FileType     string    `gorm:"type:varchar(50)" json:"fileType"`
	BusinessType string    `gorm:"type:varchar(50)" json:"businessType"`
	StoragePath  string    `gorm:"type:varchar(512);not null" json:"storagePath"`
	Status       int8      `gorm:"type:tinyint;not null;default:1" json:"status"`
	UploadedAt   time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"uploadedAt"`
	CreatedAt    time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

// AIOperationLog AI操作日志表
type AIOperationLog struct {
	ID              uint      `gorm:"primarykey;autoIncrement" json:"id"`
	OperationID     string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"operationId"`
	ApplicationID   string    `gorm:"type:varchar(64);index" json:"applicationId"`
	ApplicationType string    `gorm:"type:varchar(50);not null;index" json:"applicationType"`
	Operation       string    `gorm:"type:varchar(50);not null;index" json:"operation"`
	Result          string    `gorm:"type:varchar(50);not null" json:"result"`
	Details         string    `gorm:"type:text" json:"details"`
	Timestamp       time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"timestamp"`
	CreatedAt       time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

func (FileUpload) TableName() string     { return "file_uploads" }
func (AIOperationLog) TableName() string { return "ai_operation_logs" }
