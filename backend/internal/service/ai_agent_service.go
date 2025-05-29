package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"backend/internal/data"
)

// AIAgentService AI智能体服务
type AIAgentService struct {
	data *data.Data
	log  *zap.Logger
}

// NewAIAgentService 创建AI智能体服务
func NewAIAgentService(data *data.Data, log *zap.Logger) *AIAgentService {
	return &AIAgentService{
		data: data,
		log:  log,
	}
}

// ApplicationInfo 申请信息响应结构
type ApplicationInfo struct {
	ApplicationID     string           `json:"application_id"`
	ProductInfo       ProductInfo      `json:"product_info"`
	ApplicationInfo   AppInfo          `json:"application_info"`
	ApplicantInfo     ApplicantInfo    `json:"applicant_info"`
	FinancialInfo     FinancialInfo    `json:"financial_info"`
	UploadedDocuments []AIDocumentInfo `json:"uploaded_documents"`
	ExternalData      ExternalDataInfo `json:"external_data"`
}

type ProductInfo struct {
	ProductID          string  `json:"product_id"`
	Name               string  `json:"name"`
	Category           string  `json:"category"`
	MaxAmount          float64 `json:"max_amount"`
	InterestRateYearly string  `json:"interest_rate_yearly"`
}

type AppInfo struct {
	Amount      float64   `json:"amount"`
	TermMonths  int       `json:"term_months"`
	Purpose     string    `json:"purpose"`
	SubmittedAt time.Time `json:"submitted_at"`
	Status      string    `json:"status"`
}

type ApplicantInfo struct {
	UserID       string `json:"user_id"`
	RealName     string `json:"real_name"`
	IDCardNumber string `json:"id_card_number"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	Age          int    `json:"age"`
	IsVerified   bool   `json:"is_verified"`
}

type FinancialInfo struct {
	AnnualIncome      float64 `json:"annual_income"`
	ExistingLoans     int     `json:"existing_loans"`
	CreditScore       int     `json:"credit_score"`
	AccountBalance    float64 `json:"account_balance"`
	LandArea          string  `json:"land_area"`
	FarmingExperience string  `json:"farming_experience"`
}

type AIDocumentInfo struct {
	DocType       string            `json:"doc_type"`
	FileID        string            `json:"file_id"`
	FileURL       string            `json:"file_url"`
	OCRResult     map[string]string `json:"ocr_result,omitempty"`
	ExtractedInfo map[string]string `json:"extracted_info,omitempty"`
}

type ExternalDataInfo struct {
	CreditBureauScore     int           `json:"credit_bureau_score"`
	BlacklistCheck        bool          `json:"blacklist_check"`
	PreviousLoanHistory   []interface{} `json:"previous_loan_history"`
	LandOwnershipVerified bool          `json:"land_ownership_verified"`
}

// AIDecisionRequest AI决策请求
type AIDecisionRequest struct {
	ApplicationID         string   `json:"application_id"`
	ApplicationType       string   `json:"application_type"`
	Decision              string   `json:"decision"`
	RiskScore             float64  `json:"risk_score"`
	RiskLevel             string   `json:"risk_level"`
	ConfidenceScore       float64  `json:"confidence_score"`
	AnalysisSummary       string   `json:"analysis_summary"`
	ApprovedAmount        *float64 `json:"approved_amount,omitempty"`
	ApprovedTermMonths    *int     `json:"approved_term_months,omitempty"`
	SuggestedInterestRate string   `json:"suggested_interest_rate,omitempty"`
	SuggestedDeposit      *float64 `json:"suggested_deposit,omitempty"`
	DetailedAnalysis      string   `json:"detailed_analysis,omitempty"`
	Recommendations       []string `json:"recommendations,omitempty"`
	Conditions            []string `json:"conditions,omitempty"`
	AIModelVersion        string   `json:"ai_model_version,omitempty"`
	WorkflowID            string   `json:"workflow_id,omitempty"`
}

// AIDecisionResponse AI决策响应
type AIDecisionResponse struct {
	Success         bool   `json:"success"`
	AIOperationID   string `json:"ai_operation_id"`
	ProcessedAt     string `json:"processed_at"`
	NextStepMessage string `json:"next_step_message"`
}

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

// GetApplicationInfoUnified 获取申请信息（统一接口）
func (s *AIAgentService) GetApplicationInfoUnified(ctx context.Context, applicationID string) (*AIApplicationInfo, error) {
	// 根据申请ID格式判断申请类型
	applicationType := s.detectApplicationType(applicationID)

	switch applicationType {
	case "LOAN_APPLICATION":
		return s.getLoanApplicationInfo(ctx, applicationID)
	case "MACHINERY_LEASING":
		return s.getMachineryLeasingInfo(ctx, applicationID)
	default:
		return nil, fmt.Errorf("申请不存在")
	}
}

// SubmitAIDecisionUnified 提交AI决策结果（统一接口）
func (s *AIAgentService) SubmitAIDecisionUnified(ctx context.Context, request *AIDecisionRequest) (*AIDecisionResponse, error) {
	// 验证申请是否存在
	exists, err := s.validateApplicationExists(request.ApplicationID, request.ApplicationType)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("申请不存在")
	}

	// 生成操作ID
	operationID := fmt.Sprintf("ai_op_%d", time.Now().Unix())

	// 记录AI决策日志
	aiLog := data.AIOperationLog{
		OperationID:     operationID,
		ApplicationID:   request.ApplicationID,
		ApplicationType: request.ApplicationType,
		Operation:       "AI_DECISION",
		Result:          request.Decision,
		Details:         request.AnalysisSummary,
		Timestamp:       time.Now(),
	}

	if err := s.data.DB.Create(&aiLog).Error; err != nil {
		s.log.Error("保存AI决策日志失败", zap.Error(err))
		return nil, fmt.Errorf("AI服务暂时不可用")
	}

	// 根据申请类型更新申请状态
	switch request.ApplicationType {
	case "LOAN_APPLICATION":
		err = s.updateLoanApplicationDecision(request)
	case "MACHINERY_LEASING":
		err = s.updateMachineryLeasingDecision(request)
	default:
		return nil, fmt.Errorf("不支持的申请类型")
	}

	if err != nil {
		return nil, err
	}

	return &AIDecisionResponse{
		Success:         true,
		AIOperationID:   operationID,
		ProcessedAt:     time.Now().Format(time.RFC3339),
		NextStepMessage: "AI决策已成功处理，系统将继续后续流程",
	}, nil
}

// GetExternalDataUnified 获取外部数据（多类型支持）
func (s *AIAgentService) GetExternalDataUnified(ctx context.Context, userID string, dataTypes []string) (*ExternalDataResponse, error) {
	// 验证用户是否存在
	var user data.User
	if err := s.data.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	response := &ExternalDataResponse{
		UserID:      userID,
		DataTypes:   dataTypes,
		RetrievedAt: time.Now(),
	}

	// 根据请求的数据类型获取相应数据
	for _, dataType := range dataTypes {
		switch dataType {
		case "credit_report":
			response.CreditData = s.getCreditReportData(userID)
		case "bank_flow":
			response.BankData = s.getBankFlowData(userID)
		case "blacklist_check":
			response.BlacklistData = s.getBlacklistData(userID)
		case "government_subsidy":
			response.GovernmentData = s.getGovernmentSubsidyData(userID)
		case "farming_qualification":
			response.FarmingData = s.getFarmingQualificationData(userID)
		default:
			return nil, fmt.Errorf("不支持的数据类型")
		}
	}

	return response, nil
}

// GetAIModelConfigUnified 获取AI模型配置（多类型支持）
func (s *AIAgentService) GetAIModelConfigUnified(ctx context.Context) (*AIModelConfigResponse, error) {
	return &AIModelConfigResponse{
		Models: []AIModelConfig{
			{
				ModelName: "claude-3.5-sonnet",
				Version:   "v1.0",
				Type:      "LOAN_APPLICATION",
				Enabled:   true,
			},
			{
				ModelName: "gpt-4o",
				Version:   "v1.0",
				Type:      "MACHINERY_LEASING",
				Enabled:   true,
			},
		},
		RiskThresholds: map[string]float64{
			"low_risk_threshold":     0.3,
			"medium_risk_threshold":  0.7,
			"auto_approve_threshold": 0.2,
			"auto_reject_threshold":  0.8,
		},
		BusinessRules: map[string]interface{}{
			"max_loan_amount":       1000000,
			"min_credit_score":      600,
			"max_leasing_term_days": 365,
			"required_documents":    []string{"id_card", "income_proof"},
		},
		UpdatedAt: time.Now(),
	}, nil
}

// GetMachineryLeasingApplicationInfo 获取农机租赁申请信息（专用接口）
func (s *AIAgentService) GetMachineryLeasingApplicationInfo(ctx context.Context, applicationID string) (*MachineryLeasingApplicationInfo, error) {
	// 这里应该查询农机租赁申请数据
	// 目前返回模拟数据
	return &MachineryLeasingApplicationInfo{
		ApplicationID: applicationID,
		LesseeInfo: map[string]interface{}{
			"name":    "张三",
			"id_card": "123456789012345678",
			"phone":   "13800138000",
		},
		LessorInfo: map[string]interface{}{
			"company_name": "农机租赁公司",
			"contact":      "李四",
		},
		MachineryInfo: map[string]interface{}{
			"type":  "拖拉机",
			"model": "东方红LX1000",
			"value": 150000,
		},
		LeasingDetails: map[string]interface{}{
			"rental_amount": 5000,
			"term_days":     30,
			"deposit":       10000,
		},
		Status:      "SUBMITTED",
		SubmittedAt: time.Now(),
	}, nil
}

// GetAIOperationLogs 获取AI操作日志
func (s *AIAgentService) GetAIOperationLogs(ctx context.Context, applicationID, applicationType string, page, limit int) ([]AIOperationLog, int64, error) {
	query := s.data.DB.Model(&data.AIOperationLog{})

	if applicationID != "" {
		query = query.Where("application_id = ?", applicationID)
	}
	if applicationType != "" {
		query = query.Where("application_type = ?", applicationType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []data.AIOperationLog
	offset := (page - 1) * limit
	if err := query.Order("timestamp DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	var result []AIOperationLog
	for _, log := range logs {
		result = append(result, AIOperationLog{
			OperationID:     log.OperationID,
			ApplicationID:   log.ApplicationID,
			ApplicationType: log.ApplicationType,
			Operation:       log.Operation,
			Result:          log.Result,
			Details:         log.Details,
			Timestamp:       log.Timestamp,
		})
	}

	return result, total, nil
}

// 辅助方法

// detectApplicationType 根据申请ID检测申请类型
func (s *AIAgentService) detectApplicationType(applicationID string) string {
	if strings.Contains(applicationID, "ml_") || strings.Contains(applicationID, "leasing_") {
		return "MACHINERY_LEASING"
	}
	if strings.Contains(applicationID, "app_") || strings.Contains(applicationID, "loan_") || strings.Contains(applicationID, "test_app_") {
		return "LOAN_APPLICATION"
	}
	return "UNKNOWN"
}

// getLoanApplicationInfo 获取贷款申请信息
func (s *AIAgentService) getLoanApplicationInfo(ctx context.Context, applicationID string) (*AIApplicationInfo, error) {
	var application data.LoanApplication
	if err := s.data.DB.Where("application_id = ?", applicationID).First(&application).Error; err != nil {
		return nil, fmt.Errorf("申请不存在")
	}

	// 解析申请人信息
	var applicantInfo map[string]interface{}
	if application.ApplicantSnapshot != nil {
		json.Unmarshal(application.ApplicantSnapshot, &applicantInfo)
	}

	// 获取产品信息
	var product data.LoanProduct
	s.data.DB.Where("product_id = ?", application.ProductID).First(&product)

	return &AIApplicationInfo{
		ApplicationType: "LOAN_APPLICATION",
		ApplicationID:   application.ApplicationID,
		UserID:          application.UserID,
		Status:          application.Status,
		BasicInfo: map[string]interface{}{
			"amount":      application.AmountApplied,
			"term_months": application.TermMonthsApplied,
			"purpose":     application.Purpose,
		},
		FinancialInfo: applicantInfo,
		ProductInfo: map[string]interface{}{
			"product_id":   product.ProductID,
			"product_name": product.Name,
			"category":     product.Category,
		},
		SubmittedAt: application.SubmittedAt,
	}, nil
}

// getMachineryLeasingInfo 获取农机租赁申请信息
func (s *AIAgentService) getMachineryLeasingInfo(ctx context.Context, applicationID string) (*AIApplicationInfo, error) {
	// 这里应该查询真实的农机租赁数据
	// 目前返回模拟数据
	return &AIApplicationInfo{
		ApplicationType: "MACHINERY_LEASING",
		ApplicationID:   applicationID,
		UserID:          "user_123",
		Status:          "SUBMITTED",
		BasicInfo: map[string]interface{}{
			"machinery_type": "拖拉机",
			"rental_amount":  5000,
			"term_days":      30,
		},
		SubmittedAt: time.Now(),
	}, nil
}

// validateApplicationExists 验证申请是否存在
func (s *AIAgentService) validateApplicationExists(applicationID, applicationType string) (bool, error) {
	switch applicationType {
	case "LOAN_APPLICATION":
		var count int64
		err := s.data.DB.Model(&data.LoanApplication{}).Where("application_id = ?", applicationID).Count(&count).Error
		return count > 0, err
	case "MACHINERY_LEASING":
		// 这里应该查询农机租赁表
		// 目前直接返回true
		return true, nil
	}
	return false, nil
}

// updateLoanApplicationDecision 更新贷款申请决策
func (s *AIAgentService) updateLoanApplicationDecision(request *AIDecisionRequest) error {
	updates := map[string]interface{}{
		"status":        "AI_PROCESSED",
		"ai_risk_score": int(request.RiskScore * 100),
		"ai_suggestion": request.AnalysisSummary,
	}

	if request.ApprovedAmount != nil {
		updates["approved_amount"] = *request.ApprovedAmount
	}
	if request.ApprovedTermMonths != nil {
		updates["approved_term_months"] = *request.ApprovedTermMonths
	}

	return s.data.DB.Model(&data.LoanApplication{}).
		Where("application_id = ?", request.ApplicationID).
		Updates(updates).Error
}

// updateMachineryLeasingDecision 更新农机租赁决策
func (s *AIAgentService) updateMachineryLeasingDecision(request *AIDecisionRequest) error {
	// 这里应该更新农机租赁申请
	// 目前直接返回成功
	return nil
}

// 外部数据获取方法

func (s *AIAgentService) getCreditReportData(userID string) map[string]interface{} {
	return map[string]interface{}{
		"credit_score":  750,
		"credit_grade":  "A",
		"overdue_count": 0,
	}
}

func (s *AIAgentService) getBankFlowData(userID string) map[string]interface{} {
	return map[string]interface{}{
		"monthly_income":  15000,
		"monthly_expense": 8000,
		"balance":         50000,
	}
}

func (s *AIAgentService) getBlacklistData(userID string) map[string]interface{} {
	return map[string]interface{}{
		"is_blacklisted": false,
		"risk_level":     "LOW",
	}
}

func (s *AIAgentService) getGovernmentSubsidyData(userID string) map[string]interface{} {
	return map[string]interface{}{
		"has_subsidy":    true,
		"subsidy_amount": 5000,
		"subsidy_type":   "农业补贴",
	}
}

func (s *AIAgentService) getFarmingQualificationData(userID string) map[string]interface{} {
	return map[string]interface{}{
		"is_qualified":        true,
		"qualification_level": "专业农户",
		"farming_years":       10,
	}
}
