package service

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"backend/internal/data"
)

// UnifiedApplicationProcessor 统一申请处理器
type UnifiedApplicationProcessor struct {
	data           *data.Data
	loanService    *LoanService
	leasingService *MachineryLeasingApprovalService
	aiService      *AIAgentService
	log            *zap.Logger
}

// NewUnifiedApplicationProcessor 创建统一申请处理器
func NewUnifiedApplicationProcessor(
	data *data.Data,
	loanService *LoanService,
	leasingService *MachineryLeasingApprovalService,
	aiService *AIAgentService,
	log *zap.Logger,
) *UnifiedApplicationProcessor {
	return &UnifiedApplicationProcessor{
		data:           data,
		loanService:    loanService,
		leasingService: leasingService,
		aiService:      aiService,
		log:            log,
	}
}

// ApplicationTypeDetectionResult 申请类型识别结果
type ApplicationTypeDetectionResult struct {
	ApplicationType string  `json:"application_type"`
	Confidence      float64 `json:"confidence"`
	DetectionMethod string  `json:"detection_method"`
	Reason          string  `json:"reason"`
}

// UnifiedApplicationInfo 统一申请信息响应
type UnifiedApplicationInfo struct {
	ApplicationType string                 `json:"application_type"`
	ApplicationID   string                 `json:"application_id"`
	UserID          string                 `json:"user_id"`
	Status          string                 `json:"status"`
	SubmittedAt     time.Time              `json:"submitted_at"`
	BasicInfo       map[string]interface{} `json:"basic_info"`
	BusinessInfo    map[string]interface{} `json:"business_info"`
	ApplicantInfo   map[string]interface{} `json:"applicant_info"`
	FinancialInfo   map[string]interface{} `json:"financial_info"`
	RiskAssessment  map[string]interface{} `json:"risk_assessment"`
	Documents       []interface{}          `json:"documents"`
}

// UnifiedExternalDataResponse 统一外部数据响应
type UnifiedExternalDataResponse struct {
	UserID          string                 `json:"user_id"`
	ApplicationType string                 `json:"application_type"`
	DataTypes       []string               `json:"data_types"`
	CreditData      map[string]interface{} `json:"credit_data,omitempty"`
	BankData        map[string]interface{} `json:"bank_data,omitempty"`
	BlacklistData   map[string]interface{} `json:"blacklist_data,omitempty"`
	GovernmentData  map[string]interface{} `json:"government_data,omitempty"`
	FarmingData     map[string]interface{} `json:"farming_data,omitempty"`
	RetrievedAt     time.Time              `json:"retrieved_at"`
}

// UnifiedDecisionParams 统一决策参数
type UnifiedDecisionParams struct {
	ApplicationID         string                 `json:"application_id"`
	Decision              string                 `json:"decision"`
	RiskScore             float64                `json:"risk_score"`
	RiskLevel             string                 `json:"risk_level"`
	ConfidenceScore       float64                `json:"confidence_score"`
	AnalysisSummary       string                 `json:"analysis_summary"`
	ApprovedAmount        *float64               `json:"approved_amount,omitempty"`
	ApprovedTermMonths    *int                   `json:"approved_term_months,omitempty"`
	SuggestedInterestRate *string                `json:"suggested_interest_rate,omitempty"`
	SuggestedDeposit      *float64               `json:"suggested_deposit,omitempty"`
	DetailedAnalysis      map[string]interface{} `json:"detailed_analysis"`
	Recommendations       []string               `json:"recommendations"`
	Conditions            []string               `json:"conditions"`
	AIModelVersion        string                 `json:"ai_model_version"`
	WorkflowID            string                 `json:"workflow_id"`
}

// UnifiedDecisionResponse 统一决策响应
type UnifiedDecisionResponse struct {
	ApplicationID     string                 `json:"application_id"`
	ApplicationType   string                 `json:"application_type"`
	Decision          string                 `json:"decision"`
	NewStatus         string                 `json:"new_status"`
	NextStep          string                 `json:"next_step"`
	DecisionID        string                 `json:"decision_id"`
	AIOperationID     string                 `json:"ai_operation_id"`
	ProcessingSummary map[string]interface{} `json:"processing_summary"`
}

// DetectApplicationType 智能识别申请类型
func (p *UnifiedApplicationProcessor) DetectApplicationType(applicationID string) (*ApplicationTypeDetectionResult, error) {
	p.log.Info("开始检测申请类型", zap.String("application_id", applicationID))

	// 1. 基于ID模式识别（快速路径）
	if typeResult := p.detectByIDPattern(applicationID); typeResult.Confidence >= 0.9 {
		p.log.Info("基于ID模式识别成功",
			zap.String("application_id", applicationID),
			zap.String("type", typeResult.ApplicationType),
			zap.Float64("confidence", typeResult.Confidence))
		return typeResult, nil
	}

	// 2. 基于数据库查询识别（准确路径）
	if typeResult := p.detectByDatabaseQuery(applicationID); typeResult.Confidence >= 0.95 {
		p.log.Info("基于数据库查询识别成功",
			zap.String("application_id", applicationID),
			zap.String("type", typeResult.ApplicationType),
			zap.Float64("confidence", typeResult.Confidence))
		return typeResult, nil
	}

	// 3. 如果都无法识别，返回未知
	return &ApplicationTypeDetectionResult{
		ApplicationType: "UNKNOWN",
		Confidence:      0.0,
		DetectionMethod: "FALLBACK",
		Reason:          "无法通过ID模式或数据库查询识别申请类型",
	}, fmt.Errorf("无法识别申请类型: %s", applicationID)
}

// detectByIDPattern 基于ID模式识别
func (p *UnifiedApplicationProcessor) detectByIDPattern(applicationID string) *ApplicationTypeDetectionResult {
	// 农机租赁申请模式
	if strings.HasPrefix(applicationID, "ml_") || strings.HasPrefix(applicationID, "leasing_") {
		return &ApplicationTypeDetectionResult{
			ApplicationType: "MACHINERY_LEASING",
			Confidence:      0.95,
			DetectionMethod: "ID_PATTERN",
			Reason:          fmt.Sprintf("申请ID '%s' 匹配农机租赁模式 (ml_* 或 leasing_*)", applicationID),
		}
	}

	// 贷款申请模式
	if strings.HasPrefix(applicationID, "test_app_") ||
		strings.HasPrefix(applicationID, "app_") ||
		strings.HasPrefix(applicationID, "loan_") {
		return &ApplicationTypeDetectionResult{
			ApplicationType: "LOAN_APPLICATION",
			Confidence:      0.95,
			DetectionMethod: "ID_PATTERN",
			Reason:          fmt.Sprintf("申请ID '%s' 匹配贷款申请模式 (test_app_*, app_*, loan_*)", applicationID),
		}
	}

	return &ApplicationTypeDetectionResult{
		ApplicationType: "UNKNOWN",
		Confidence:      0.0,
		DetectionMethod: "ID_PATTERN",
		Reason:          fmt.Sprintf("申请ID '%s' 不匹配任何已知模式", applicationID),
	}
}

// detectByDatabaseQuery 基于数据库查询识别
func (p *UnifiedApplicationProcessor) detectByDatabaseQuery(applicationID string) *ApplicationTypeDetectionResult {
	// 检查贷款申请表
	var loanApp data.LoanApplication
	if err := p.data.DB.Where("application_id = ?", applicationID).First(&loanApp).Error; err == nil {
		return &ApplicationTypeDetectionResult{
			ApplicationType: "LOAN_APPLICATION",
			Confidence:      0.99,
			DetectionMethod: "DATABASE_QUERY",
			Reason:          fmt.Sprintf("在贷款申请表中找到记录: %s", applicationID),
		}
	}

	// 检查农机租赁申请表
	var leasingApp data.MachineryLeasingApplication
	if err := p.data.DB.Where("application_id = ?", applicationID).First(&leasingApp).Error; err == nil {
		return &ApplicationTypeDetectionResult{
			ApplicationType: "MACHINERY_LEASING",
			Confidence:      0.99,
			DetectionMethod: "DATABASE_QUERY",
			Reason:          fmt.Sprintf("在农机租赁申请表中找到记录: %s", applicationID),
		}
	}

	return &ApplicationTypeDetectionResult{
		ApplicationType: "UNKNOWN",
		Confidence:      0.0,
		DetectionMethod: "DATABASE_QUERY",
		Reason:          fmt.Sprintf("在任何申请表中都未找到记录: %s", applicationID),
	}
}

// GetApplicationInfoUnified 统一获取申请信息
func (p *UnifiedApplicationProcessor) GetApplicationInfoUnified(applicationID string) (*UnifiedApplicationInfo, error) {
	// 1. 识别申请类型
	typeResult, err := p.DetectApplicationType(applicationID)
	if err != nil {
		return nil, fmt.Errorf("识别申请类型失败: %w", err)
	}

	if typeResult.Confidence < 0.9 {
		return nil, fmt.Errorf("申请类型识别置信度不足: %.2f", typeResult.Confidence)
	}

	// 2. 根据类型获取具体信息
	switch typeResult.ApplicationType {
	case "LOAN_APPLICATION":
		return p.getLoanApplicationInfo(applicationID)
	case "MACHINERY_LEASING":
		return p.getMachineryLeasingInfo(applicationID)
	default:
		return nil, fmt.Errorf("不支持的申请类型: %s", typeResult.ApplicationType)
	}
}

// getLoanApplicationInfo 获取贷款申请信息
func (p *UnifiedApplicationProcessor) getLoanApplicationInfo(applicationID string) (*UnifiedApplicationInfo, error) {
	// 获取贷款申请详情
	var loanApp data.LoanApplication
	if err := p.data.DB.Where("application_id = ?", applicationID).First(&loanApp).Error; err != nil {
		return nil, fmt.Errorf("贷款申请不存在: %w", err)
	}

	// 获取用户信息
	var user data.User
	if err := p.data.DB.Where("user_id = ?", loanApp.UserID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("申请人不存在: %w", err)
	}

	// 获取用户画像
	var userProfile data.UserProfile
	p.data.DB.Where("user_id = ?", loanApp.UserID).First(&userProfile)

	// 获取产品信息
	var product data.LoanProduct
	if err := p.data.DB.Where("product_id = ?", loanApp.ProductID).First(&product).Error; err != nil {
		return nil, fmt.Errorf("贷款产品不存在: %w", err)
	}

	// 构建统一响应
	return &UnifiedApplicationInfo{
		ApplicationType: "LOAN_APPLICATION",
		ApplicationID:   loanApp.ApplicationID,
		UserID:          loanApp.UserID,
		Status:          loanApp.Status,
		SubmittedAt:     loanApp.SubmittedAt,
		BasicInfo: map[string]interface{}{
			"amount":      loanApp.AmountApplied,
			"term_months": loanApp.TermMonthsApplied,
			"purpose":     loanApp.Purpose,
		},
		BusinessInfo: map[string]interface{}{
			"product_id":           product.ProductID,
			"product_name":         product.Name,
			"category":             product.Category,
			"interest_rate_yearly": product.InterestRateYearly,
			"max_amount":           product.MaxAmount,
		},
		ApplicantInfo: map[string]interface{}{
			"user_id":        user.UserID,
			"phone":          user.Phone,
			"real_name":      userProfile.RealName,
			"id_card_number": userProfile.IDCardNumber,
			"address":        userProfile.Address,
		},
		FinancialInfo: map[string]interface{}{
			"annual_income": userProfile.AnnualIncome,
			"occupation":    userProfile.Occupation,
		},
		RiskAssessment: map[string]interface{}{
			"ai_risk_score": loanApp.AIRiskScore,
			"ai_suggestion": loanApp.AISuggestion,
		},
		Documents: []interface{}{}, // TODO: 获取相关文档
	}, nil
}

// getMachineryLeasingInfo 获取农机租赁申请信息
func (p *UnifiedApplicationProcessor) getMachineryLeasingInfo(applicationID string) (*UnifiedApplicationInfo, error) {
	// 获取农机租赁申请详情
	var leasingApp data.MachineryLeasingApplication
	if err := p.data.DB.Where("application_id = ?", applicationID).First(&leasingApp).Error; err != nil {
		return nil, fmt.Errorf("农机租赁申请不存在: %w", err)
	}

	// 获取承租方信息
	var lesseeUser data.User
	if err := p.data.DB.Where("user_id = ?", leasingApp.LesseeUserID).First(&lesseeUser).Error; err != nil {
		return nil, fmt.Errorf("承租方不存在: %w", err)
	}

	// 获取出租方信息
	var lessorUser data.User
	if err := p.data.DB.Where("user_id = ?", leasingApp.LessorUserID).First(&lessorUser).Error; err != nil {
		return nil, fmt.Errorf("出租方不存在: %w", err)
	}

	// 获取农机信息
	var machinery data.FarmMachinery
	if err := p.data.DB.Where("machinery_id = ?", leasingApp.MachineryID).First(&machinery).Error; err != nil {
		return nil, fmt.Errorf("农机不存在: %w", err)
	}

	// 构建统一响应
	return &UnifiedApplicationInfo{
		ApplicationType: "MACHINERY_LEASING",
		ApplicationID:   leasingApp.ApplicationID,
		UserID:          leasingApp.LesseeUserID, // 主要申请人为承租方
		Status:          leasingApp.ApplicationStatus,
		SubmittedAt:     leasingApp.SubmittedAt,
		BasicInfo: map[string]interface{}{
			"requested_start_date": leasingApp.RequestedStartDate,
			"requested_end_date":   leasingApp.RequestedEndDate,
			"rental_days":          leasingApp.RentalDays,
			"total_amount":         leasingApp.TotalAmount,
			"deposit_amount":       leasingApp.DepositAmount,
			"usage_purpose":        leasingApp.UsagePurpose,
		},
		BusinessInfo: map[string]interface{}{
			"machinery_id":   machinery.MachineryID,
			"machinery_type": machinery.Type,
			"brand_model":    machinery.BrandModel,
			"daily_rent":     machinery.DailyRent,
			"location":       machinery.LocationText,
		},
		ApplicantInfo: map[string]interface{}{
			"lessee_info": map[string]interface{}{
				"user_id": lesseeUser.UserID,
				"phone":   lesseeUser.Phone,
			},
			"lessor_info": map[string]interface{}{
				"user_id": lessorUser.UserID,
				"phone":   lessorUser.Phone,
			},
		},
		FinancialInfo: map[string]interface{}{
			"total_amount":   leasingApp.TotalAmount,
			"deposit_amount": leasingApp.DepositAmount,
		},
		RiskAssessment: map[string]interface{}{
			"ai_risk_score": leasingApp.AIRiskScore,
			"ai_suggestion": leasingApp.AISuggestion,
			"risk_level":    leasingApp.RiskLevel,
		},
		Documents: []interface{}{}, // TODO: 获取相关文档
	}, nil
}

// GetExternalDataUnified 统一获取外部数据
func (p *UnifiedApplicationProcessor) GetExternalDataUnified(userID string, dataTypes []string, applicationID string) (*UnifiedExternalDataResponse, error) {
	// 1. 如果提供了申请ID，识别申请类型以优化数据获取
	var applicationType string
	if applicationID != "" {
		if typeResult, err := p.DetectApplicationType(applicationID); err == nil && typeResult.Confidence >= 0.9 {
			applicationType = typeResult.ApplicationType
		}
	}

	// 2. 根据申请类型过滤和优化数据类型
	filteredDataTypes := p.filterDataTypesByApplicationType(dataTypes, applicationType)

	// 3. 获取外部数据
	response := &UnifiedExternalDataResponse{
		UserID:          userID,
		ApplicationType: applicationType,
		DataTypes:       filteredDataTypes,
		RetrievedAt:     time.Now(),
	}

	// 4. 根据数据类型获取相应数据
	for _, dataType := range filteredDataTypes {
		switch dataType {
		case "credit_report":
			response.CreditData = p.getMockCreditData(userID)
		case "bank_flow":
			response.BankData = p.getMockBankData(userID)
		case "blacklist_check":
			response.BlacklistData = p.getMockBlacklistData(userID)
		case "government_subsidy":
			if applicationType == "LOAN_APPLICATION" || applicationType == "MACHINERY_LEASING" {
				response.GovernmentData = p.getMockGovernmentData(userID)
			}
		case "farming_qualification":
			if applicationType == "MACHINERY_LEASING" {
				response.FarmingData = p.getMockFarmingData(userID)
			}
		}
	}

	return response, nil
}

// filterDataTypesByApplicationType 根据申请类型过滤数据类型
func (p *UnifiedApplicationProcessor) filterDataTypesByApplicationType(dataTypes []string, applicationType string) []string {
	if applicationType == "" {
		return dataTypes // 如果无法识别类型，返回全部数据类型
	}

	var filtered []string
	for _, dataType := range dataTypes {
		switch dataType {
		case "credit_report", "bank_flow", "blacklist_check":
			// 这些数据对所有申请类型都有用
			filtered = append(filtered, dataType)
		case "government_subsidy":
			// 政府补贴对农业相关申请有用
			if applicationType == "LOAN_APPLICATION" || applicationType == "MACHINERY_LEASING" {
				filtered = append(filtered, dataType)
			}
		case "farming_qualification":
			// 农业资质主要用于农机租赁
			if applicationType == "MACHINERY_LEASING" {
				filtered = append(filtered, dataType)
			}
		}
	}

	return filtered
}

// SubmitAIDecisionUnified 统一提交AI决策
func (p *UnifiedApplicationProcessor) SubmitAIDecisionUnified(params *UnifiedDecisionParams) (*UnifiedDecisionResponse, error) {
	// 1. 识别申请类型
	typeResult, err := p.DetectApplicationType(params.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("识别申请类型失败: %w", err)
	}

	if typeResult.Confidence < 0.9 {
		return nil, fmt.Errorf("申请类型识别置信度不足: %.2f", typeResult.Confidence)
	}

	// 2. 验证决策是否适用于申请类型
	if err := p.validateDecisionForApplicationType(params.Decision, typeResult.ApplicationType); err != nil {
		return nil, fmt.Errorf("决策验证失败: %w", err)
	}

	// 3. 根据申请类型路由到相应处理器
	switch typeResult.ApplicationType {
	case "LOAN_APPLICATION":
		return p.submitLoanDecision(params, typeResult.ApplicationType)
	case "MACHINERY_LEASING":
		return p.submitLeasingDecision(params, typeResult.ApplicationType)
	default:
		return nil, fmt.Errorf("不支持的申请类型: %s", typeResult.ApplicationType)
	}
}

// validateDecisionForApplicationType 验证决策是否适用于申请类型
func (p *UnifiedApplicationProcessor) validateDecisionForApplicationType(decision, applicationType string) error {
	switch applicationType {
	case "LOAN_APPLICATION":
		validDecisions := []string{"AUTO_APPROVED", "AUTO_REJECTED", "REQUIRE_HUMAN_REVIEW"}
		for _, valid := range validDecisions {
			if decision == valid {
				return nil
			}
		}
		return fmt.Errorf("贷款申请不支持决策: %s", decision)
	case "MACHINERY_LEASING":
		validDecisions := []string{"AUTO_APPROVE", "AUTO_REJECT", "REQUIRE_HUMAN_REVIEW", "REQUIRE_DEPOSIT_ADJUSTMENT"}
		for _, valid := range validDecisions {
			if decision == valid {
				return nil
			}
		}
		return fmt.Errorf("农机租赁申请不支持决策: %s", decision)
	default:
		return fmt.Errorf("未知申请类型: %s", applicationType)
	}
}

// submitLoanDecision 提交贷款决策
func (p *UnifiedApplicationProcessor) submitLoanDecision(params *UnifiedDecisionParams, applicationType string) (*UnifiedDecisionResponse, error) {
	// 更新贷款申请状态
	newStatus := p.getNewStatusFromDecision(params.Decision, applicationType)

	updateData := map[string]interface{}{
		"status":        newStatus,
		"ai_risk_score": int(params.RiskScore * 100), // 转换为0-100分数
		"ai_suggestion": params.AnalysisSummary,
		"processed_at":  time.Now(),
	}

	if params.ApprovedAmount != nil {
		updateData["approved_amount"] = *params.ApprovedAmount
	}
	if params.ApprovedTermMonths != nil {
		updateData["approved_term_months"] = *params.ApprovedTermMonths
	}

	if err := p.data.DB.Model(&data.LoanApplication{}).
		Where("application_id = ?", params.ApplicationID).
		Updates(updateData).Error; err != nil {
		return nil, fmt.Errorf("更新贷款申请失败: %w", err)
	}

	// 记录AI操作日志
	aiOperationID := p.recordAIOperation(params, applicationType)

	return &UnifiedDecisionResponse{
		ApplicationID:   params.ApplicationID,
		ApplicationType: applicationType,
		Decision:        params.Decision,
		NewStatus:       newStatus,
		NextStep:        p.getNextStepFromDecision(params.Decision, applicationType),
		DecisionID:      generateDecisionID(),
		AIOperationID:   aiOperationID,
		ProcessingSummary: map[string]interface{}{
			"risk_score":   params.RiskScore,
			"risk_level":   params.RiskLevel,
			"confidence":   params.ConfidenceScore,
			"processed_at": time.Now(),
		},
	}, nil
}

// submitLeasingDecision 提交农机租赁决策
func (p *UnifiedApplicationProcessor) submitLeasingDecision(params *UnifiedDecisionParams, applicationType string) (*UnifiedDecisionResponse, error) {
	// 更新农机租赁申请状态
	newStatus := p.getNewStatusFromDecision(params.Decision, applicationType)

	updateData := map[string]interface{}{
		"application_status": newStatus,
		"ai_risk_score":      params.RiskScore,
		"ai_suggestion":      params.AnalysisSummary,
		"risk_level":         params.RiskLevel,
		"approved_at":        time.Now(),
	}

	if err := p.data.DB.Model(&data.MachineryLeasingApplication{}).
		Where("application_id = ?", params.ApplicationID).
		Updates(updateData).Error; err != nil {
		return nil, fmt.Errorf("更新农机租赁申请失败: %w", err)
	}

	// 记录AI操作日志
	aiOperationID := p.recordAIOperation(params, applicationType)

	return &UnifiedDecisionResponse{
		ApplicationID:   params.ApplicationID,
		ApplicationType: applicationType,
		Decision:        params.Decision,
		NewStatus:       newStatus,
		NextStep:        p.getNextStepFromDecision(params.Decision, applicationType),
		DecisionID:      generateDecisionID(),
		AIOperationID:   aiOperationID,
		ProcessingSummary: map[string]interface{}{
			"risk_score":        params.RiskScore,
			"risk_level":        params.RiskLevel,
			"confidence":        params.ConfidenceScore,
			"suggested_deposit": params.SuggestedDeposit,
			"processed_at":      time.Now(),
		},
	}, nil
}

// getNewStatusFromDecision 根据决策获取新状态
func (p *UnifiedApplicationProcessor) getNewStatusFromDecision(decision, applicationType string) string {
	switch applicationType {
	case "LOAN_APPLICATION":
		switch decision {
		case "AUTO_APPROVED":
			return "AI_APPROVED"
		case "AUTO_REJECTED":
			return "AI_REJECTED"
		case "REQUIRE_HUMAN_REVIEW":
			return "MANUAL_REVIEW_REQUIRED"
		default:
			return "PROCESSING"
		}
	case "MACHINERY_LEASING":
		switch decision {
		case "AUTO_APPROVE":
			return "AI_APPROVED"
		case "AUTO_REJECT":
			return "AI_REJECTED"
		case "REQUIRE_HUMAN_REVIEW":
			return "MANUAL_REVIEW_REQUIRED"
		case "REQUIRE_DEPOSIT_ADJUSTMENT":
			return "DEPOSIT_ADJUSTMENT_REQUIRED"
		default:
			return "PROCESSING"
		}
	default:
		return "PROCESSING"
	}
}

// getNextStepFromDecision 根据决策获取下一步操作
func (p *UnifiedApplicationProcessor) getNextStepFromDecision(decision, applicationType string) string {
	switch decision {
	case "AUTO_APPROVED", "AUTO_APPROVE":
		return "GENERATE_CONTRACT"
	case "AUTO_REJECTED", "AUTO_REJECT":
		return "SEND_REJECTION_NOTICE"
	case "REQUIRE_HUMAN_REVIEW":
		return "ASSIGN_TO_REVIEWER"
	case "REQUIRE_DEPOSIT_ADJUSTMENT":
		return "ADJUST_DEPOSIT_TERMS"
	default:
		return "CONTINUE_PROCESSING"
	}
}

// recordAIOperation 记录AI操作
func (p *UnifiedApplicationProcessor) recordAIOperation(params *UnifiedDecisionParams, applicationType string) string {
	aiOperationID := generateAIOperationID()

	aiLog := &data.AIOperationLog{
		OperationID:     aiOperationID,
		ApplicationID:   params.ApplicationID,
		ApplicationType: applicationType,
		Operation:       "AI_DECISION",
		Result:          params.Decision,
		Details:         params.AnalysisSummary,
		Timestamp:       time.Now(),
	}

	p.data.DB.Create(aiLog)

	return aiOperationID
}

// Mock data generation methods
func (p *UnifiedApplicationProcessor) getMockCreditData(userID string) map[string]interface{} {
	return map[string]interface{}{
		"credit_score":  750,
		"credit_grade":  "A",
		"overdue_count": 0,
		"total_loans":   2,
		"last_update":   time.Now().AddDate(0, -1, 0),
	}
}

func (p *UnifiedApplicationProcessor) getMockBankData(userID string) map[string]interface{} {
	return map[string]interface{}{
		"monthly_income":   15000,
		"monthly_expense":  8000,
		"account_balance":  50000,
		"income_stability": "STABLE",
		"last_6_months_flow": []map[string]interface{}{
			{"month": "2024-11", "income": 15500, "expense": 8200},
			{"month": "2024-10", "income": 14800, "expense": 7900},
		},
	}
}

func (p *UnifiedApplicationProcessor) getMockBlacklistData(userID string) map[string]interface{} {
	return map[string]interface{}{
		"is_blacklisted": false,
		"risk_level":     "LOW",
		"check_time":     time.Now(),
	}
}

func (p *UnifiedApplicationProcessor) getMockGovernmentData(userID string) map[string]interface{} {
	return map[string]interface{}{
		"received_subsidies": []map[string]interface{}{
			{"year": 2024, "type": "农业生产补贴", "amount": 5000},
			{"year": 2023, "type": "土地流转补贴", "amount": 8000},
		},
		"total_subsidies": 13000,
	}
}

func (p *UnifiedApplicationProcessor) getMockFarmingData(userID string) map[string]interface{} {
	return map[string]interface{}{
		"certification_level":        "高级农机操作员",
		"experience_years":           8,
		"machinery_operation_skills": []string{"拖拉机", "收割机", "播种机"},
		"safety_training_completed":  true,
		"last_training_date":         "2024-03-15",
	}
}

// Utility functions
func generateDecisionID() string {
	return fmt.Sprintf("decision_%d", time.Now().UnixNano())
}

func generateAIOperationID() string {
	return fmt.Sprintf("ai_op_%d", time.Now().UnixNano())
}
