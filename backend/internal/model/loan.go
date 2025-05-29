package model

import (
	"time"

	"gorm.io/gorm"
)

// LoanProduct 贷款产品表
type LoanProduct struct {
	ID           uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductName  string  `gorm:"type:varchar(100);not null" json:"product_name"`
	ProductCode  string  `gorm:"type:varchar(50);unique;not null" json:"product_code"`
	ProductType  string  `gorm:"type:varchar(50);not null" json:"product_type"` // micro_loan(小额贷款), farm_loan(农场贷款), cooperative_loan(合作社贷款)
	Description  string  `gorm:"type:text" json:"description"`
	MinAmount    int64   `gorm:"not null" json:"min_amount"`                      // 最小贷款金额（分）
	MaxAmount    int64   `gorm:"not null" json:"max_amount"`                      // 最大贷款金额（分）
	InterestRate float64 `gorm:"type:decimal(5,4);not null" json:"interest_rate"` // 年利率
	TermMonths   int     `gorm:"not null" json:"term_months"`                     // 贷款期限（月）
	RequiredAuth string  `gorm:"type:text" json:"required_auth"`                  // 所需认证类型，逗号分隔
	IsActive     bool    `gorm:"default:true" json:"is_active"`                   // 是否激活
	SortOrder    int     `gorm:"default:0" json:"sort_order"`                     // 排序权重

	// Dify工作流配置
	DifyWorkflowID   string `gorm:"type:varchar(100)" json:"dify_workflow_id"`  // Dify工作流ID
	EligibleUserType string `gorm:"type:varchar(50)" json:"eligible_user_type"` // 适用用户类型

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	LoanApplications []LoanApplication `gorm:"foreignKey:ProductID" json:"loan_applications,omitempty"`
}

// LoanApplication 贷款申请表
type LoanApplication struct {
	ID            uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ApplicationNo string `gorm:"type:varchar(30);uniqueIndex;not null" json:"application_no"`

	// 关联信息
	UserID    uint64 `gorm:"not null;index" json:"user_id"`
	ProductID uint64 `gorm:"not null;index" json:"product_id"`

	// 申请基本信息
	ApplyAmount     int64      `gorm:"not null" json:"apply_amount"`                   // 申请金额(分)
	LoanAmount      int64      `gorm:"not null" json:"loan_amount"`                    // 申请金额(分) - 兼容service层
	ApplyTermMonths int        `gorm:"not null" json:"apply_term_months"`              // 申请期限(月)
	TermMonths      int        `gorm:"not null" json:"term_months"`                    // 申请期限(月) - 兼容service层
	LoanPurpose     string     `gorm:"type:varchar(100);not null" json:"loan_purpose"` // 贷款用途
	ExpectedUseDate *time.Time `json:"expected_use_date"`                              // 预期用款日期

	// 申请人信息
	ApplicantName   string `gorm:"type:varchar(50);not null" json:"applicant_name"`
	ApplicantIDCard string `gorm:"type:varchar(18);not null" json:"applicant_id_card"`
	ApplicantPhone  string `gorm:"type:varchar(20);not null" json:"applicant_phone"`
	ContactPhone    string `gorm:"type:varchar(20);not null" json:"contact_phone"` // 兼容service层
	ContactEmail    string `gorm:"type:varchar(100)" json:"contact_email"`         // 兼容service层

	// 收入信息
	MonthlyIncome int64  `json:"monthly_income"`                         // 月收入(分)
	YearlyIncome  int64  `json:"yearly_income"`                          // 年收入(分)
	IncomeSource  string `gorm:"type:varchar(100)" json:"income_source"` // 收入来源
	OtherDebts    int64  `json:"other_debts"`                            // 其他负债(分)

	// 经营信息(农业相关)
	FarmArea          float64 `gorm:"type:decimal(10,2)" json:"farm_area"`       // 经营面积(亩)
	CropTypes         string  `gorm:"type:json" json:"crop_types"`               // 种植作物(JSON数组)
	YearsOfExperience int     `json:"years_of_experience"`                       // 从业年限
	LandCertificate   string  `gorm:"type:varchar(100)" json:"land_certificate"` // 土地证明

	// 申请材料(JSON格式存储文件路径)
	ApplicationDocuments string `gorm:"type:json" json:"application_documents"`
	MaterialsJSON        string `gorm:"type:json" json:"materials_json"` // 兼容service层

	// 审批流程
	Status             string     `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	CurrentApprover    *uint64    `json:"current_approver"`                          // 当前审批人
	ApprovalLevel      int        `gorm:"default:1" json:"approval_level"`           // 当前审批层级
	AutoApprovalPassed bool       `gorm:"default:false" json:"auto_approval_passed"` // 是否通过自动审批
	FinalApprovedAt    *time.Time `json:"final_approved_at"`                         // 最终审批时间
	FinalApprover      *uint64    `json:"final_approver"`                            // 最终审批人
	SubmittedAt        time.Time  `json:"submitted_at"`                              // 提交时间 - 兼容service层

	// 审批结果
	ApprovedAmount     *int64   `json:"approved_amount"`                        // 批准金额(分)
	ApprovedTermMonths *int     `json:"approved_term_months"`                   // 批准期限(月)
	ApprovedRate       *float64 `gorm:"type:decimal(5,4)" json:"approved_rate"` // 批准利率
	RejectionReason    string   `gorm:"type:text" json:"rejection_reason"`      // 拒绝原因

	// 风控评估
	CreditScore     int     `json:"credit_score"`                               // 信用评分
	RiskLevel       string  `gorm:"type:varchar(20)" json:"risk_level"`         // 风险等级
	DebtIncomeRatio float64 `gorm:"type:decimal(5,4)" json:"debt_income_ratio"` // 负债收入比
	RiskAssessment  string  `gorm:"type:text" json:"risk_assessment"`           // 风险评估报告

	// Dify AI工作流
	DifyConversationID string `gorm:"type:varchar(100)" json:"dify_conversation_id"` // 对话ID
	AIRecommendation   string `gorm:"type:text" json:"ai_recommendation"`            // AI审批建议

	// 放款信息
	DisbursementMethod  string     `gorm:"type:varchar(20)" json:"disbursement_method"`  // 放款方式
	DisbursementAccount string     `gorm:"type:varchar(50)" json:"disbursement_account"` // 放款账户
	DisbursedAt         *time.Time `json:"disbursed_at"`                                 // 放款时间
	DisbursedAmount     *int64     `json:"disbursed_amount"`                             // 实际放款金额(分)

	// 备注
	Remarks string `gorm:"type:text" json:"remarks"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	User         User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product      LoanProduct       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ApprovalLogs []ApprovalLog     `gorm:"foreignKey:ApplicationID" json:"approval_logs,omitempty"`
	DifyLogs     []DifyWorkflowLog `gorm:"foreignKey:ApplicationID" json:"dify_logs,omitempty"`
}

// ApprovalLog 审批日志表
type ApprovalLog struct {
	ID            uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ApplicationID uint64 `gorm:"not null;index" json:"application_id"`
	ApproverID    uint64 `gorm:"not null;index" json:"approver_id"`

	// 审批动作：submit(提交)、approve(通过)、reject(拒绝)、return(退回)、withdraw(撤回)
	Action string `gorm:"type:varchar(20);not null" json:"action"`
	Step   string `gorm:"type:varchar(20);not null" json:"step"` // 兼容service层

	// 审批层级
	ApprovalLevel int `gorm:"not null" json:"approval_level"`

	// 审批结果：pending(待处理)、approved(通过)、rejected(拒绝)、returned(退回)
	Result string `gorm:"type:varchar(20);not null" json:"result"`
	Status string `gorm:"type:varchar(20);not null" json:"status"` // 兼容service层

	// 审批意见
	Comment string `gorm:"type:text" json:"comment"`
	Note    string `gorm:"type:text" json:"note"` // 兼容service层

	// 审批前后状态
	PreviousStatus string `gorm:"type:varchar(20)" json:"previous_status"`
	NewStatus      string `gorm:"type:varchar(20)" json:"new_status"`

	// 审批时间
	ActionTime time.Time  `json:"action_time"`
	ApprovedAt *time.Time `json:"approved_at"` // 兼容service层

	// 审批参数调整
	AmountAdjustment *int64   `json:"amount_adjustment"`                        // 金额调整(分)
	TermAdjustment   *int     `json:"term_adjustment"`                          // 期限调整(月)
	RateAdjustment   *float64 `gorm:"type:decimal(5,4)" json:"rate_adjustment"` // 利率调整

	// 附件
	AttachmentFiles string `gorm:"type:json" json:"attachment_files"`

	CreatedAt time.Time `json:"created_at"`

	// 关联
	Application LoanApplication `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
	Approver    OAUser          `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
}

// DifyWorkflowLog Dify工作流调用日志表
type DifyWorkflowLog struct {
	ID            uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ApplicationID uint64 `gorm:"not null;index" json:"application_id"`

	// Dify相关信息
	WorkflowID     string `gorm:"type:varchar(100);not null" json:"workflow_id"`
	ConversationID string `gorm:"type:varchar(100)" json:"conversation_id"`
	MessageID      string `gorm:"type:varchar(100)" json:"message_id"`

	// 工作流类型：risk_assessment(风险评估)、document_review(文档审核)、decision_support(决策支持)
	WorkflowType string `gorm:"type:varchar(30);not null" json:"workflow_type"`

	// 调用信息
	RequestData  string `gorm:"type:json;not null" json:"request_data"` // 请求数据
	ResponseData string `gorm:"type:json" json:"response_data"`         // 响应数据

	// 执行状态：pending(处理中)、completed(完成)、failed(失败)、timeout(超时)
	Status string `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`

	// 执行结果
	Result          string  `gorm:"type:text" json:"result"`                   // 工作流执行结果
	Recommendation  string  `gorm:"type:text" json:"recommendation"`           // AI建议
	ConfidenceScore float64 `gorm:"type:decimal(3,2)" json:"confidence_score"` // 置信度评分

	// 性能指标
	StartTime  time.Time  `json:"start_time"`
	EndTime    *time.Time `json:"end_time"`
	Duration   *int       `json:"duration"`    // 执行时长(毫秒)
	TokenUsage int        `json:"token_usage"` // Token使用量
	CostAmount int64      `json:"cost_amount"` // 成本(分)

	// 错误信息
	ErrorCode    string `gorm:"type:varchar(50)" json:"error_code"`
	ErrorMessage string `gorm:"type:text" json:"error_message"`

	// 重试信息
	RetryCount int `gorm:"default:0" json:"retry_count"`
	MaxRetries int `gorm:"default:3" json:"max_retries"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	Application LoanApplication `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
}

// LoanApplyConditions 贷款申请条件结构
type LoanApplyConditions struct {
	MinAge               int      `json:"min_age"`               // 最小年龄
	MaxAge               int      `json:"max_age"`               // 最大年龄
	MinYearsExperience   int      `json:"min_years_experience"`  // 最小从业年限
	MinFarmArea          float64  `json:"min_farm_area"`         // 最小经营面积
	RequiredRegions      []string `json:"required_regions"`      // 适用地区
	RequiredCertificates []string `json:"required_certificates"` // 必需证件
	BlacklistIndustries  []string `json:"blacklist_industries"`  // 禁入行业
}

// LoanFeatures 贷款产品特色
type LoanFeatures struct {
	FastApproval      bool   `json:"fast_approval"`      // 快速审批
	OnlineApplication bool   `json:"online_application"` // 在线申请
	FlexibleRepayment bool   `json:"flexible_repayment"` // 灵活还款
	NoCollateral      bool   `json:"no_collateral"`      // 免抵押
	LowInterestRate   bool   `json:"low_interest_rate"`  // 低利率
	SeasonalSupport   bool   `json:"seasonal_support"`   // 季节性支持
	Description       string `json:"description"`        // 特色描述
}

// ApplicationDocuments 申请材料结构
type ApplicationDocuments struct {
	IDCardFront     string   `json:"id_card_front"`    // 身份证正面
	IDCardBack      string   `json:"id_card_back"`     // 身份证背面
	BankStatement   []string `json:"bank_statement"`   // 银行流水
	IncomeProof     []string `json:"income_proof"`     // 收入证明
	LandCertificate []string `json:"land_certificate"` // 土地证明
	BusinessLicense string   `json:"business_license"` // 营业执照
	CropRecords     []string `json:"crop_records"`     // 种植记录
	Other           []string `json:"other"`            // 其他材料
}

// DifyRequestData Dify请求数据结构
type DifyRequestData struct {
	Query  string                 `json:"query"`  // 查询内容
	Inputs map[string]interface{} `json:"inputs"` // 输入变量
	UserID string                 `json:"user"`   // 用户标识
	Files  []DifyFile             `json:"files"`  // 文件列表
}

// DifyFile Dify文件结构
type DifyFile struct {
	Type         string `json:"type"`            // 文件类型
	TransferType string `json:"transfer_method"` // 传输方式
	URL          string `json:"url"`             // 文件URL
	UploadFileID string `json:"upload_file_id"`  // 上传文件ID
}

// DifyResponseData Dify响应数据结构
type DifyResponseData struct {
	ConversationID string                 `json:"conversation_id"` // 对话ID
	MessageID      string                 `json:"message_id"`      // 消息ID
	Answer         string                 `json:"answer"`          // 回答内容
	Metadata       map[string]interface{} `json:"metadata"`        // 元数据
	Usage          DifyUsage              `json:"usage"`           // 使用统计
}

// DifyUsage Dify使用统计
type DifyUsage struct {
	PromptTokens     int     `json:"prompt_tokens"`     // 提示词Token数
	CompletionTokens int     `json:"completion_tokens"` // 完成Token数
	TotalTokens      int     `json:"total_tokens"`      // 总Token数
	TotalPrice       string  `json:"total_price"`       // 总价格
	Currency         string  `json:"currency"`          // 货币单位
	Latency          float64 `json:"latency"`           // 延迟(秒)
}

// TableName 设置表名
func (LoanProduct) TableName() string {
	return "loan_products"
}

func (LoanApplication) TableName() string {
	return "loan_applications"
}

func (ApprovalLog) TableName() string {
	return "approval_logs"
}

func (DifyWorkflowLog) TableName() string {
	return "dify_workflow_logs"
}
