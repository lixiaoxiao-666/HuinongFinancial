package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"backend/internal/data"
	"backend/pkg"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AIAgentService AI智能体服务
type AIAgentService struct {
	data *data.Data
	log  *zap.Logger
}

// NewAIAgentService 创建AI智能体服务实例
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

// AIDecisionRequest AI决策请求结构
type AIDecisionRequest struct {
	AIAnalysis     AIAnalysisData     `json:"ai_analysis"`
	AIDecision     AIDecisionData     `json:"ai_decision"`
	ProcessingInfo ProcessingInfoData `json:"processing_info"`
}

type AIAnalysisData struct {
	RiskLevel        string                   `json:"risk_level"`
	RiskScore        float64                  `json:"risk_score"`
	ConfidenceScore  float64                  `json:"confidence_score"`
	AnalysisSummary  string                   `json:"analysis_summary"`
	DetailedAnalysis map[string]interface{}   `json:"detailed_analysis"`
	RiskFactors      []map[string]interface{} `json:"risk_factors"`
	Recommendations  []string                 `json:"recommendations"`
}

type AIDecisionData struct {
	Decision              string   `json:"decision"`
	ApprovedAmount        float64  `json:"approved_amount"`
	ApprovedTermMonths    int      `json:"approved_term_months"`
	SuggestedInterestRate string   `json:"suggested_interest_rate"`
	Conditions            []string `json:"conditions"`
	NextAction            string   `json:"next_action"`
}

type ProcessingInfoData struct {
	AIModelVersion   string    `json:"ai_model_version"`
	ProcessingTimeMs int       `json:"processing_time_ms"`
	WorkflowID       string    `json:"workflow_id"`
	ProcessedAt      time.Time `json:"processed_at"`
}

// AIDecisionParams 查询参数方式的AI决策参数
type AIDecisionParams struct {
	ApplicationID         string                 `json:"application_id"`
	Decision              string                 `json:"decision"`
	RiskScore             float64                `json:"risk_score"`
	ConfidenceScore       float64                `json:"confidence_score"`
	ApprovedAmount        float64                `json:"approved_amount"`
	ApprovedTermMonths    int                    `json:"approved_term_months"`
	SuggestedInterestRate string                 `json:"suggested_interest_rate"`
	RiskLevel             string                 `json:"risk_level"`
	AnalysisSummary       string                 `json:"analysis_summary"`
	DetailedAnalysis      map[string]interface{} `json:"detailed_analysis"`
	Recommendations       []string               `json:"recommendations"`
	Conditions            []string               `json:"conditions"`
	AIModelVersion        string                 `json:"ai_model_version"`
	WorkflowID            string                 `json:"workflow_id"`
}

// AIDecisionResult AI决策处理结果
type AIDecisionResult struct {
	ApplicationID string `json:"application_id"`
	NewStatus     string `json:"new_status"`
	NextStep      string `json:"next_step"`
}

// GetApplicationInfo 获取申请详细信息供AI分析
func (s *AIAgentService) GetApplicationInfo(applicationID string) (*ApplicationInfo, error) {
	s.log.Info("获取申请信息", zap.String("applicationId", applicationID))

	// 查询申请基本信息
	var application data.LoanApplication
	if err := s.data.DB.Where("application_id = ?", applicationID).First(&application).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("申请不存在")
		}
		return nil, fmt.Errorf("查询申请失败: %w", err)
	}

	// 查询产品信息
	var product data.LoanProduct
	if err := s.data.DB.Where("product_id = ?", application.ProductID).First(&product).Error; err != nil {
		return nil, fmt.Errorf("查询产品失败: %w", err)
	}

	// 查询用户信息
	var user data.User
	if err := s.data.DB.Where("user_id = ?", application.UserID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 查询用户画像信息
	var profile data.UserProfile
	profileExists := true
	if err := s.data.DB.Where("user_id = ?", application.UserID).First(&profile).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("查询用户画像失败: %w", err)
		}
		profileExists = false
	}

	// 查询上传文件
	var files []data.UploadedFile
	s.data.DB.Where("user_id = ? AND related_id = ?", application.UserID, applicationID).Find(&files)

	// 构建响应
	info := &ApplicationInfo{
		ApplicationID: applicationID,
		ProductInfo: ProductInfo{
			ProductID:          product.ProductID,
			Name:               product.Name,
			Category:           product.Category,
			MaxAmount:          product.MaxAmount,
			InterestRateYearly: product.InterestRateYearly,
		},
		ApplicationInfo: AppInfo{
			Amount:      application.AmountApplied,
			TermMonths:  application.TermMonthsApplied,
			Purpose:     application.Purpose,
			SubmittedAt: application.SubmittedAt,
			Status:      application.Status,
		},
		ApplicantInfo: ApplicantInfo{
			UserID:       user.UserID,
			RealName:     "张三", // 模拟数据
			IDCardNumber: "310***********1234",
			Phone:        maskPhone(user.Phone),
			Address:      "XX省XX市XX村",
			Age:          35,
			IsVerified:   true,
		},
		FinancialInfo: FinancialInfo{
			AnnualIncome:      80000,
			ExistingLoans:     0,
			CreditScore:       750,
			AccountBalance:    15000,
			LandArea:          "10亩",
			FarmingExperience: "10年",
		},
		ExternalData: ExternalDataInfo{
			CreditBureauScore:     750,
			BlacklistCheck:        false,
			PreviousLoanHistory:   []interface{}{},
			LandOwnershipVerified: true,
		},
	}

	// 填充用户画像信息（如果存在）
	if profileExists {
		info.ApplicantInfo.RealName = profile.RealName
		info.ApplicantInfo.Address = profile.Address
		info.FinancialInfo.AnnualIncome = profile.AnnualIncome
		if profile.IDCardNumber != "" {
			info.ApplicantInfo.IDCardNumber = maskIDCard(profile.IDCardNumber)
		}
	}

	// 填充文件信息
	for _, file := range files {
		docInfo := AIDocumentInfo{
			DocType: file.Purpose,
			FileID:  file.FileID,
			FileURL: file.StoragePath,
		}

		// 模拟OCR结果
		if file.Purpose == "id_card_front" {
			docInfo.OCRResult = map[string]string{
				"name":      info.ApplicantInfo.RealName,
				"id_number": info.ApplicantInfo.IDCardNumber,
				"address":   info.ApplicantInfo.Address,
			}
		} else if file.Purpose == "land_contract" {
			docInfo.ExtractedInfo = map[string]string{
				"land_area":       "10亩",
				"contract_period": "30年",
				"location":        info.ApplicantInfo.Address,
			}
		}

		info.UploadedDocuments = append(info.UploadedDocuments, docInfo)
	}

	s.log.Info("申请信息获取成功", zap.String("applicationId", applicationID))
	return info, nil
}

// SubmitAIDecision 提交AI审批决策
func (s *AIAgentService) SubmitAIDecision(applicationID string, request *AIDecisionRequest) error {
	s.log.Info("提交AI审批决策", zap.String("applicationId", applicationID))

	// 开始事务
	tx := s.data.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查申请是否存在
	var application data.LoanApplication
	if err := tx.Where("application_id = ?", applicationID).First(&application).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("申请不存在")
		}
		return fmt.Errorf("查询申请失败: %w", err)
	}

	// 保存AI分析结果
	detailedAnalysisBytes, _ := json.Marshal(request.AIAnalysis.DetailedAnalysis)
	riskFactorsBytes, _ := json.Marshal(request.AIAnalysis.RiskFactors)
	recommendationsBytes, _ := json.Marshal(request.AIAnalysis.Recommendations)

	aiResult := &data.AIAnalysisResult{
		ApplicationID:         applicationID,
		WorkflowExecutionID:   request.ProcessingInfo.WorkflowID,
		RiskLevel:             request.AIAnalysis.RiskLevel,
		RiskScore:             request.AIAnalysis.RiskScore,
		ConfidenceScore:       request.AIAnalysis.ConfidenceScore,
		AnalysisSummary:       request.AIAnalysis.AnalysisSummary,
		DetailedAnalysis:      detailedAnalysisBytes,
		RiskFactors:           riskFactorsBytes,
		Recommendations:       recommendationsBytes,
		AIDecision:            request.AIDecision.Decision,
		ApprovedAmount:        &request.AIDecision.ApprovedAmount,
		ApprovedTermMonths:    &request.AIDecision.ApprovedTermMonths,
		SuggestedInterestRate: request.AIDecision.SuggestedInterestRate,
		NextAction:            request.AIDecision.NextAction,
		AIModelVersion:        request.ProcessingInfo.AIModelVersion,
		ProcessingTimeMs:      request.ProcessingInfo.ProcessingTimeMs,
		ProcessedAt:           request.ProcessingInfo.ProcessedAt,
	}

	if err := tx.Create(aiResult).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("保存AI分析结果失败: %w", err)
	}

	// 更新申请状态
	newStatus := ""
	switch request.AIDecision.Decision {
	case "AUTO_APPROVED":
		newStatus = "AI_APPROVED"
	case "REQUIRE_HUMAN_REVIEW":
		newStatus = "MANUAL_REVIEW_REQUIRED"
	case "AUTO_REJECTED":
		newStatus = "AI_REJECTED"
	default:
		newStatus = "AI_PROCESSED"
	}

	// 更新申请记录
	updateData := map[string]interface{}{
		"status":        newStatus,
		"ai_risk_score": int(request.AIAnalysis.RiskScore * 1000), // 转换为千分制
		"ai_suggestion": request.AIAnalysis.AnalysisSummary,
		"updated_at":    time.Now(),
	}

	if request.AIDecision.Decision == "AUTO_APPROVED" {
		updateData["approved_amount"] = request.AIDecision.ApprovedAmount
		updateData["approved_term_months"] = request.AIDecision.ApprovedTermMonths
		updateData["final_decision"] = "APPROVED"
		updateData["decision_reason"] = "AI自动审批通过"
		updateData["processed_by"] = "AI_SYSTEM"
		updateData["processed_at"] = time.Now()
	} else if request.AIDecision.Decision == "AUTO_REJECTED" {
		updateData["final_decision"] = "REJECTED"
		updateData["decision_reason"] = request.AIAnalysis.AnalysisSummary
		updateData["processed_by"] = "AI_SYSTEM"
		updateData["processed_at"] = time.Now()
	}

	if err := tx.Model(&application).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新申请状态失败: %w", err)
	}

	// 记录状态变更历史
	history := &data.LoanApplicationHistory{
		ApplicationID: applicationID,
		StatusFrom:    application.Status,
		StatusTo:      newStatus,
		OperatorType:  "AI_SYSTEM",
		OperatorID:    "ai_agent",
		Comments:      request.AIAnalysis.AnalysisSummary,
		OccurredAt:    time.Now(),
	}

	if err := tx.Create(history).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("记录状态历史失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	s.log.Info("AI审批决策处理完成",
		zap.String("applicationId", applicationID),
		zap.String("decision", request.AIDecision.Decision),
		zap.String("newStatus", newStatus),
	)

	return nil
}

// TriggerWorkflow 触发AI审批工作流
func (s *AIAgentService) TriggerWorkflow(applicationID, workflowType, priority, callbackURL string) (*data.WorkflowExecution, error) {
	s.log.Info("触发AI工作流",
		zap.String("applicationId", applicationID),
		zap.String("workflowType", workflowType),
	)

	// 检查申请是否存在
	var application data.LoanApplication
	if err := s.data.DB.Where("application_id = ?", applicationID).First(&application).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("申请不存在")
		}
		return nil, fmt.Errorf("查询申请失败: %w", err)
	}

	// 创建工作流执行记录
	execution := &data.WorkflowExecution{
		ExecutionID:         pkg.GenerateUUID(),
		ApplicationID:       applicationID,
		WorkflowType:        workflowType,
		Status:              "RUNNING",
		Priority:            priority,
		StartedAt:           time.Now(),
		EstimatedCompletion: timePtr(time.Now().Add(5 * time.Minute)),
		CurrentStage:        "DATA_COLLECTION",
		Progress:            0,
		CallbackURL:         callbackURL,
		RetryCount:          0,
	}

	if err := s.data.DB.Create(execution).Error; err != nil {
		return nil, fmt.Errorf("创建工作流记录失败: %w", err)
	}

	// 更新申请状态为AI处理中
	if err := s.data.DB.Model(&application).Updates(map[string]interface{}{
		"status":     "AI_TRIGGERED",
		"updated_at": time.Now(),
	}).Error; err != nil {
		s.log.Error("更新申请状态失败", zap.Error(err))
	}

	// 记录状态变更历史
	history := &data.LoanApplicationHistory{
		ApplicationID: applicationID,
		StatusFrom:    application.Status,
		StatusTo:      "AI_TRIGGERED",
		OperatorType:  "SYSTEM",
		OperatorID:    "workflow_trigger",
		Comments:      fmt.Sprintf("触发AI工作流: %s", workflowType),
		OccurredAt:    time.Now(),
	}
	s.data.DB.Create(history)

	s.log.Info("AI工作流触发成功",
		zap.String("executionId", execution.ExecutionID),
		zap.String("applicationId", applicationID),
	)

	return execution, nil
}

// GetAIModelConfig 获取AI模型配置
func (s *AIAgentService) GetAIModelConfig() (map[string]interface{}, error) {
	s.log.Info("获取AI模型配置")

	// 模拟配置数据，实际项目中从数据库或配置文件读取
	config := map[string]interface{}{
		"active_models": []map[string]interface{}{
			{
				"model_id":   "risk_assessment_v2",
				"model_type": "RISK_EVALUATION",
				"version":    "2.1.0",
				"status":     "ACTIVE",
				"thresholds": map[string]float64{
					"low_risk":    0.3,
					"medium_risk": 0.7,
					"high_risk":   0.9,
				},
			},
			{
				"model_id":    "fraud_detection_v1",
				"model_type":  "FRAUD_DETECTION",
				"version":     "1.5.2",
				"status":      "ACTIVE",
				"sensitivity": 0.85,
			},
		},
		"approval_rules": map[string]interface{}{
			"auto_approval_threshold":  0.3,
			"auto_rejection_threshold": 0.8,
			"max_auto_approval_amount": 50000,
			"required_human_review_conditions": []string{
				"申请金额超过5万元",
				"信用评分低于700分",
				"存在潜在欺诈风险",
			},
		},
		"business_parameters": map[string]interface{}{
			"max_debt_to_income_ratio": 0.5,
			"min_credit_score":         600,
			"max_loan_amount_by_category": map[string]float64{
				"种植贷": 50000,
				"设备贷": 200000,
				"其他":  30000,
			},
		},
	}

	return config, nil
}

// GetExternalData 获取外部数据
func (s *AIAgentService) GetExternalData(userID, dataTypes string) (map[string]interface{}, error) {
	s.log.Info("获取外部数据",
		zap.String("userId", userID),
		zap.String("dataTypes", dataTypes),
	)

	// 记录查询请求
	query := &data.ExternalDataQuery{
		QueryID:       pkg.GenerateUUID(),
		UserID:        userID,
		DataTypes:     dataTypes,
		Status:        "SUCCESS",
		QueryDuration: 1500, // 模拟查询耗时1.5秒
		QueriedAt:     time.Now(),
	}

	// 模拟外部数据
	externalData := map[string]interface{}{
		"user_id": userID,
		"credit_report": map[string]interface{}{
			"score":           750,
			"grade":           "优秀",
			"report_date":     "2024-03-01",
			"loan_history":    []interface{}{},
			"overdue_records": 0,
		},
		"bank_flow": map[string]interface{}{
			"average_monthly_income": 6500,
			"account_stability":      "稳定",
			"last_6_months_flow": []map[string]interface{}{
				{"month": "2024-02", "income": 7200, "expense": 4800},
				{"month": "2024-01", "income": 6800, "expense": 5100},
			},
		},
		"blacklist_check": map[string]interface{}{
			"is_blacklisted": false,
			"check_time":     time.Now().Format(time.RFC3339),
		},
		"government_subsidy": map[string]interface{}{
			"received_subsidies": []map[string]interface{}{
				{"year": 2023, "type": "种粮补贴", "amount": 1200},
				{"year": 2022, "type": "农机购置补贴", "amount": 3000},
			},
		},
	}

	// 保存查询结果
	queryResultBytes, _ := json.Marshal(externalData)
	query.QueryResult = queryResultBytes

	if err := s.data.DB.Create(query).Error; err != nil {
		s.log.Error("保存外部数据查询记录失败", zap.Error(err))
	}

	return externalData, nil
}

// UpdateApplicationStatus 更新申请状态
func (s *AIAgentService) UpdateApplicationStatus(applicationID, status, operator, remarks string, metadata map[string]interface{}) error {
	s.log.Info("更新申请状态",
		zap.String("applicationId", applicationID),
		zap.String("status", status),
		zap.String("operator", operator),
	)

	// 查询当前申请
	var application data.LoanApplication
	if err := s.data.DB.Where("application_id = ?", applicationID).First(&application).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("申请不存在")
		}
		return fmt.Errorf("查询申请失败: %w", err)
	}

	oldStatus := application.Status

	// 更新申请状态
	if err := s.data.DB.Model(&application).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("更新申请状态失败: %w", err)
	}

	// 记录状态变更历史
	history := &data.LoanApplicationHistory{
		ApplicationID: applicationID,
		StatusFrom:    oldStatus,
		StatusTo:      status,
		OperatorType:  "AI_SYSTEM",
		OperatorID:    operator,
		Comments:      remarks,
		OccurredAt:    time.Now(),
	}

	if err := s.data.DB.Create(history).Error; err != nil {
		s.log.Error("记录状态历史失败", zap.Error(err))
	}

	return nil
}

// SubmitAIDecisionQuery 处理查询参数方式的AI决策提交
func (s *AIAgentService) SubmitAIDecisionQuery(ctx context.Context, params *AIDecisionParams) (*AIDecisionResult, error) {
	s.log.Info("Processing AI decision with query parameters",
		zap.String("application_id", params.ApplicationID),
		zap.String("decision", params.Decision),
		zap.Float64("risk_score", params.RiskScore))

	// 1. 验证申请是否存在
	var application data.LoanApplication
	if err := s.data.DB.Where("application_id = ?", params.ApplicationID).First(&application).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("申请不存在: %s", params.ApplicationID)
		}
		return nil, fmt.Errorf("获取申请信息失败: %w", err)
	}

	// 2. 检查申请状态是否允许AI决策
	allowedStatuses := []string{"SUBMITTED", "AI_TRIGGERED", "AI_REVIEWING", "pending_review"}
	if !contains(allowedStatuses, application.Status) {
		return nil, fmt.Errorf("申请状态不允许AI决策: %s", application.Status)
	}

	// 3. 构建AI分析结果数据
	detailedAnalysisJSON, _ := json.Marshal(params.DetailedAnalysis)
	recommendationsJSON, _ := json.Marshal(params.Recommendations)

	// 4. 开始数据库事务
	tx := s.data.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("开始事务失败: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 5. 创建AI分析结果记录
	aiAnalysisResult := &data.AIAnalysisResult{
		ApplicationID:         params.ApplicationID,
		WorkflowExecutionID:   params.WorkflowID + "_" + fmt.Sprintf("%d", time.Now().Unix()),
		RiskLevel:             params.RiskLevel,
		RiskScore:             params.RiskScore,
		ConfidenceScore:       params.ConfidenceScore,
		AnalysisSummary:       params.AnalysisSummary,
		DetailedAnalysis:      detailedAnalysisJSON,
		RiskFactors:           recommendationsJSON, // 临时使用recommendations作为risk_factors
		Recommendations:       recommendationsJSON,
		AIDecision:            params.Decision,
		ApprovedAmount:        &params.ApprovedAmount,
		ApprovedTermMonths:    &params.ApprovedTermMonths,
		SuggestedInterestRate: params.SuggestedInterestRate,
		NextAction:            getNextAction(params.Decision),
		AIModelVersion:        params.AIModelVersion,
		ProcessingTimeMs:      2000, // 默认处理时间
		ProcessedAt:           time.Now(),
	}

	if err := tx.Create(aiAnalysisResult).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("保存AI分析结果失败: %w", err)
	}

	// 6. 确定新的申请状态
	var newStatus, nextStep string
	oldStatus := application.Status

	switch params.Decision {
	case "AUTO_APPROVED":
		newStatus = "AI_APPROVED"
		nextStep = "AWAIT_FINAL_CONFIRMATION"
		// 更新批准金额和期限
		application.ApprovedAmount = &params.ApprovedAmount
		application.ApprovedTermMonths = &params.ApprovedTermMonths
	case "REQUIRE_HUMAN_REVIEW":
		newStatus = "MANUAL_REVIEW_REQUIRED"
		nextStep = "AWAIT_HUMAN_REVIEW"
	case "AUTO_REJECTED":
		newStatus = "AI_REJECTED"
		nextStep = "PROCESS_COMPLETED"
	default:
		newStatus = "AI_PROCESSED"
		nextStep = "AWAIT_NEXT_ACTION"
	}

	// 7. 更新申请状态
	application.Status = newStatus
	application.AISuggestion = params.AnalysisSummary
	application.UpdatedAt = time.Now()

	if err := tx.Save(&application).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新申请状态失败: %w", err)
	}

	// 8. 记录状态变更历史
	historyRecord := &data.LoanApplicationHistory{
		ApplicationID: params.ApplicationID,
		StatusFrom:    oldStatus,
		StatusTo:      newStatus,
		OperatorType:  "AI_SYSTEM",
		OperatorID:    "ai_agent",
		Comments:      fmt.Sprintf("AI决策: %s, 风险评分: %.2f", params.Decision, params.RiskScore),
		OccurredAt:    time.Now(),
	}

	if err := tx.Create(historyRecord).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("记录状态历史失败: %w", err)
	}

	// 9. 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	s.log.Info("AI decision processed successfully",
		zap.String("application_id", params.ApplicationID),
		zap.String("new_status", newStatus),
		zap.String("decision", params.Decision))

	// 10. 返回处理结果
	result := &AIDecisionResult{
		ApplicationID: params.ApplicationID,
		NewStatus:     newStatus,
		NextStep:      nextStep,
	}

	return result, nil
}

// getNextAction 根据决策确定下一步行动
func getNextAction(decision string) string {
	actionMap := map[string]string{
		"AUTO_APPROVED":        "GENERATE_CONTRACT",
		"AUTO_REJECTED":        "SEND_REJECTION_NOTICE",
		"REQUIRE_HUMAN_REVIEW": "ASSIGN_TO_REVIEWER",
	}
	if action, exists := actionMap[decision]; exists {
		return action
	}
	return "MANUAL_REVIEW"
}

// 辅助函数：检查字符串是否在切片中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// 辅助函数
func maskPhone(phone string) string {
	if len(phone) < 7 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

func maskIDCard(idCard string) string {
	if len(idCard) < 10 {
		return idCard
	}
	return idCard[:3] + "***********" + idCard[len(idCard)-4:]
}

func timePtr(t time.Time) *time.Time {
	return &t
}
