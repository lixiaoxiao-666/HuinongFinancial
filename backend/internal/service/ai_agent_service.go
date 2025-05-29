package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"backend/internal/data"
	"backend/pkg"
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

// LogAIAgentAction 记录AI智能体操作日志
func (s *AIAgentService) LogAIAgentAction(actionType, agentType, applicationID string, requestData, responseData interface{}, status string, errorMessage string, duration int, req *http.Request) {
	logID := pkg.GenerateUUID()

	var reqBytes, respBytes []byte
	if requestData != nil {
		reqBytes, _ = json.Marshal(requestData)
	}
	if responseData != nil {
		respBytes, _ = json.Marshal(responseData)
	}

	ipAddress := ""
	userAgent := ""
	if req != nil {
		ipAddress = getClientIP(req)
		userAgent = req.UserAgent()
	}

	aiLog := &data.AIAgentLog{
		LogID:         logID,
		ApplicationID: applicationID,
		ActionType:    actionType,
		AgentType:     agentType,
		RequestData:   reqBytes,
		ResponseData:  respBytes,
		Status:        status,
		ErrorMessage:  errorMessage,
		Duration:      duration,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		OccurredAt:    time.Now(),
		CreatedAt:     time.Now(),
	}

	if err := s.data.DB.Create(aiLog).Error; err != nil {
		s.log.Error("记录AI Agent日志失败",
			zap.Error(err),
			zap.String("logId", logID),
			zap.String("actionType", actionType),
			zap.String("applicationId", applicationID),
		)
	}
}

// getClientIP 获取客户端IP地址
func getClientIP(req *http.Request) string {
	// 尝试从X-Forwarded-For头获取
	if ip := req.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	// 尝试从X-Real-IP头获取
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	// 返回RemoteAddr
	return req.RemoteAddr
}

// GetApplicationInfo 获取申请详细信息供AI分析
func (s *AIAgentService) GetApplicationInfoWithLog(applicationID string, req *http.Request) (*ApplicationInfo, error) {
	startTime := time.Now()
	s.log.Info("获取申请信息", zap.String("applicationId", applicationID))

	// 记录请求日志
	requestData := map[string]interface{}{
		"application_id": applicationID,
	}

	info, err := s.getApplicationInfoInternal(applicationID)

	// 计算处理时间
	duration := int(time.Since(startTime).Milliseconds())

	// 记录操作日志
	status := "SUCCESS"
	errorMessage := ""
	if err != nil {
		status = "ERROR"
		errorMessage = err.Error()
	}

	s.LogAIAgentAction(
		"GET_APPLICATION_INFO",
		"DIFY_WORKFLOW",
		applicationID,
		requestData,
		info,
		status,
		errorMessage,
		duration,
		req,
	)

	return info, err
}

// GetApplicationInfo 保持原方法签名用于向后兼容
func (s *AIAgentService) GetApplicationInfo(applicationID string) (*ApplicationInfo, error) {
	return s.getApplicationInfoInternal(applicationID)
}

// GetApplicationInfoWithSupport 支持多种申请类型的统一获取方法
func (s *AIAgentService) GetApplicationInfoWithSupport(applicationID string, req *http.Request) (interface{}, error) {
	startTime := time.Now()
	s.log.Info("获取申请信息（多类型支持）", zap.String("applicationId", applicationID))

	var result interface{}
	var err error

	// 检查是否为农机租赁申请
	if s.isLeasingApplication(applicationID) {
		// 获取农机租赁申请信息
		var leasingService *MachineryLeasingApprovalService
		if s.data != nil {
			leasingService = NewMachineryLeasingApprovalService(s.data, s.log)
		}

		if leasingService != nil {
			leasingInfo, leasingErr := leasingService.GetLeasingApplicationInfoWithLog(applicationID, req)
			if leasingErr == nil {
				result = leasingInfo
			} else {
				err = leasingErr
			}
		} else {
			err = fmt.Errorf("农机租赁服务不可用")
		}
	} else {
		// 获取贷款申请信息
		loanInfo, loanErr := s.getApplicationInfoInternal(applicationID)
		if loanErr == nil {
			result = loanInfo
		} else {
			err = loanErr
		}
	}

	// 计算处理时间
	duration := int(time.Since(startTime).Milliseconds())

	// 记录操作日志
	status := "SUCCESS"
	errorMessage := ""
	if err != nil {
		status = "ERROR"
		errorMessage = err.Error()
	}

	applicationType := "LOAN_APPLICATION"
	if s.isLeasingApplication(applicationID) {
		applicationType = "LEASING_APPLICATION"
	}

	s.LogAIAgentAction(
		"GET_APPLICATION_INFO_UNIFIED",
		"DIFY_WORKFLOW",
		applicationID,
		map[string]interface{}{
			"application_id":   applicationID,
			"application_type": applicationType,
		},
		result,
		status,
		errorMessage,
		duration,
		req,
	)

	return result, err
}

// isLeasingApplication 检查是否为农机租赁申请
func (s *AIAgentService) isLeasingApplication(applicationID string) bool {
	// 检查农机租赁申请表
	var count int64
	s.data.DB.Model(&data.MachineryLeasingApplication{}).Where("application_id = ?", applicationID).Count(&count)
	return count > 0
}

// getApplicationInfoInternal 内部方法，实际的业务逻辑
func (s *AIAgentService) getApplicationInfoInternal(applicationID string) (*ApplicationInfo, error) {
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
			RealName:     maskName(profile.RealName),
			IDCardNumber: maskIDCard(profile.IDCardNumber),
			Phone:        maskPhone(user.Phone),
			Address:      profile.Address,
			Age:          calculateAge(profile.BirthDate),
			IsVerified:   profile.IDCardNumber != "",
		},
		FinancialInfo: FinancialInfo{
			AnnualIncome:      profile.AnnualIncome,
			ExistingLoans:     0,     // 需要实际计算
			CreditScore:       750,   // 模拟数据
			AccountBalance:    15000, // 模拟数据
			LandArea:          "10亩", // 模拟数据
			FarmingExperience: "10年", // 模拟数据
		},
		UploadedDocuments: make([]AIDocumentInfo, 0),
		ExternalData: ExternalDataInfo{
			CreditBureauScore:     750,
			BlacklistCheck:        false,
			PreviousLoanHistory:   []interface{}{},
			LandOwnershipVerified: true,
		},
	}

	// 处理文件信息
	for _, file := range files {
		docInfo := AIDocumentInfo{
			DocType: file.Purpose,
			FileID:  file.FileID,
			FileURL: file.StoragePath,
		}

		// 模拟OCR结果
		if file.Purpose == "id_card_front" {
			docInfo.OCRResult = map[string]string{
				"name":      profile.RealName,
				"id_number": maskIDCard(profile.IDCardNumber),
				"address":   profile.Address,
			}
		}

		info.UploadedDocuments = append(info.UploadedDocuments, docInfo)
	}

	if !profileExists {
		// 如果用户画像不存在，提供默认值
		info.ApplicantInfo.RealName = "未完善"
		info.ApplicantInfo.IDCardNumber = ""
		info.ApplicantInfo.Address = ""
		info.ApplicantInfo.Age = 0
		info.ApplicantInfo.IsVerified = false
	}

	return info, nil
}

// SubmitAIDecision 提交AI审批决策
func (s *AIAgentService) SubmitAIDecision(applicationID string, request *AIDecisionRequest) error {
	return s.SubmitAIDecisionWithLog(applicationID, request, nil)
}

// SubmitAIDecisionWithLog 提交AI审批决策（包含日志记录）
func (s *AIAgentService) SubmitAIDecisionWithLog(applicationID string, request *AIDecisionRequest, req *http.Request) error {
	startTime := time.Now()
	s.log.Info("提交AI审批决策", zap.String("applicationId", applicationID))

	err := s.submitAIDecisionInternal(applicationID, request)

	// 计算处理时间
	duration := int(time.Since(startTime).Milliseconds())

	// 记录操作日志
	status := "SUCCESS"
	errorMessage := ""
	if err != nil {
		status = "ERROR"
		errorMessage = err.Error()
	}

	responseData := &AIDecisionResult{
		ApplicationID: applicationID,
		NewStatus:     getNewStatusFromDecision(request.AIDecision.Decision),
		NextStep:      request.AIDecision.NextAction,
	}

	s.LogAIAgentAction(
		"SUBMIT_AI_DECISION",
		"DIFY_WORKFLOW",
		applicationID,
		request,
		responseData,
		status,
		errorMessage,
		duration,
		req,
	)

	return err
}

// submitAIDecisionInternal 内部方法，实际的业务逻辑
func (s *AIAgentService) submitAIDecisionInternal(applicationID string, request *AIDecisionRequest) error {
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
	newStatus := getNewStatusFromDecision(request.AIDecision.Decision)

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
		Comments:      fmt.Sprintf("AI决策: %s, 风险评分: %.2f", request.AIDecision.Decision, request.AIAnalysis.RiskScore),
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

// getNewStatusFromDecision 根据AI决策获取新状态
func getNewStatusFromDecision(decision string) string {
	switch decision {
	case "AUTO_APPROVED":
		return "AI_APPROVED"
	case "REQUIRE_HUMAN_REVIEW":
		return "MANUAL_REVIEW_REQUIRED"
	case "AUTO_REJECTED":
		return "AI_REJECTED"
	default:
		return "AI_PROCESSED"
	}
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
	return s.GetAIModelConfigWithLog(nil)
}

// GetAIModelConfigWithLog 获取AI模型配置（包含日志记录）
func (s *AIAgentService) GetAIModelConfigWithLog(req *http.Request) (map[string]interface{}, error) {
	startTime := time.Now()
	s.log.Info("获取AI模型配置")

	// 记录请求日志
	requestData := map[string]interface{}{
		"action": "get_model_config",
	}

	config, err := s.getAIModelConfigInternal()

	// 计算处理时间
	duration := int(time.Since(startTime).Milliseconds())

	// 记录操作日志
	status := "SUCCESS"
	errorMessage := ""
	if err != nil {
		status = "ERROR"
		errorMessage = err.Error()
	}

	s.LogAIAgentAction(
		"GET_AI_MODEL_CONFIG",
		"DIFY_WORKFLOW",
		"", // 这里没有application_id
		requestData,
		config,
		status,
		errorMessage,
		duration,
		req,
	)

	return config, err
}

// getAIModelConfigInternal 内部方法，实际的业务逻辑
func (s *AIAgentService) getAIModelConfigInternal() (map[string]interface{}, error) {
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
	return s.GetExternalDataWithLog(userID, dataTypes, nil)
}

// GetExternalDataWithLog 获取外部数据（包含日志记录）
func (s *AIAgentService) GetExternalDataWithLog(userID, dataTypes string, req *http.Request) (map[string]interface{}, error) {
	startTime := time.Now()
	s.log.Info("获取外部数据",
		zap.String("userId", userID),
		zap.String("dataTypes", dataTypes),
	)

	// 记录请求日志
	requestData := map[string]interface{}{
		"user_id":    userID,
		"data_types": dataTypes,
	}

	result, err := s.getExternalDataInternal(userID, dataTypes)

	// 计算处理时间
	duration := int(time.Since(startTime).Milliseconds())

	// 记录操作日志
	status := "SUCCESS"
	errorMessage := ""
	if err != nil {
		status = "ERROR"
		errorMessage = err.Error()
	}

	s.LogAIAgentAction(
		"GET_EXTERNAL_DATA",
		"DIFY_WORKFLOW",
		"", // 这里没有application_id，使用空字符串
		requestData,
		result,
		status,
		errorMessage,
		duration,
		req,
	)

	return result, err
}

// getExternalDataInternal 内部方法，实际的业务逻辑
func (s *AIAgentService) getExternalDataInternal(userID, dataTypes string) (map[string]interface{}, error) {
	if dataTypes == "" {
		dataTypes = "credit_report,bank_flow,blacklist_check,government_subsidy"
	}

	result := map[string]interface{}{
		"user_id": userID,
	}

	// 模拟征信报告数据
	if pkg.Contains(strings.Split(dataTypes, ","), "credit_report") || strings.Contains(dataTypes, "credit") {
		result["credit_report"] = map[string]interface{}{
			"score":           750,
			"grade":           "优秀",
			"report_date":     "2024-03-01",
			"loan_history":    []interface{}{},
			"overdue_records": 0,
		}
	}

	// 模拟银行流水数据
	if pkg.Contains(strings.Split(dataTypes, ","), "bank_flow") || strings.Contains(dataTypes, "bank") {
		result["bank_flow"] = map[string]interface{}{
			"average_monthly_income": 6500,
			"account_stability":      "稳定",
			"last_6_months_flow": []map[string]interface{}{
				{"month": "2024-02", "income": 7200, "expense": 4800},
				{"month": "2024-01", "income": 6800, "expense": 5100},
			},
		}
	}

	// 模拟黑名单检查
	if pkg.Contains(strings.Split(dataTypes, ","), "blacklist_check") || strings.Contains(dataTypes, "blacklist") {
		result["blacklist_check"] = map[string]interface{}{
			"is_blacklisted": false,
			"check_time":     time.Now().Format(time.RFC3339),
		}
	}

	// 模拟政府补贴信息
	if pkg.Contains(strings.Split(dataTypes, ","), "government_subsidy") || strings.Contains(dataTypes, "subsidy") {
		result["government_subsidy"] = map[string]interface{}{
			"received_subsidies": []map[string]interface{}{
				{"year": 2023, "type": "种粮补贴", "amount": 1200},
				{"year": 2022, "type": "农机购置补贴", "amount": 3000},
			},
		}
	}

	// 记录查询结果到数据库
	queryID := pkg.GenerateUUID()
	queryResult, _ := json.Marshal(result)

	externalQuery := &data.ExternalDataQuery{
		QueryID:       queryID,
		UserID:        userID,
		ApplicationID: "", // 这里可能需要传入application_id
		DataTypes:     dataTypes,
		QueryResult:   queryResult,
		Status:        "SUCCESS",
		ErrorMessage:  "",
		QueryDuration: int(time.Since(time.Now()).Milliseconds()),
		QueriedAt:     time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.data.DB.Create(externalQuery).Error; err != nil {
		s.log.Error("记录外部数据查询失败", zap.Error(err))
	}

	return result, nil
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

// SubmitAIDecisionUnified 统一的AI决策提交方法，支持贷款和农机租赁
func (s *AIAgentService) SubmitAIDecisionUnified(ctx context.Context, applicationID string, params *AIDecisionParams) (*AIDecisionResult, error) {
	startTime := time.Now()
	s.log.Info("提交AI决策（统一接口）",
		zap.String("applicationId", applicationID),
		zap.String("decision", params.Decision))

	var result *AIDecisionResult
	var err error

	// 检查申请类型并调用相应的处理方法
	if s.isLeasingApplication(applicationID) {
		// 处理农机租赁申请
		var leasingService *MachineryLeasingApprovalService
		if s.data != nil {
			leasingService = NewMachineryLeasingApprovalService(s.data, s.log)
		}

		if leasingService != nil {
			// 转换参数为农机租赁参数格式
			leasingParams := &LeasingAIDecisionParams{
				ApplicationID:       params.ApplicationID,
				Decision:            params.Decision,
				RiskScore:           params.RiskScore,
				ConfidenceScore:     params.ConfidenceScore,
				SuggestedDeposit:    params.ApprovedAmount, // 在农机租赁中，approved_amount 映射为建议押金
				SuggestedConditions: params.Conditions,
				RiskLevel:           params.RiskLevel,
				AnalysisSummary:     params.AnalysisSummary,
				DetailedAnalysis:    params.DetailedAnalysis,
				Recommendations:     params.Recommendations,
				AIModelVersion:      params.AIModelVersion,
				WorkflowID:          params.WorkflowID,
			}

			leasingResult, leasingErr := leasingService.SubmitLeasingAIDecisionQuery(ctx, leasingParams)
			if leasingErr == nil {
				result = &AIDecisionResult{
					ApplicationID: leasingResult.ApplicationID,
					NewStatus:     leasingResult.NewStatus,
					NextStep:      leasingResult.NextStep,
				}
			} else {
				err = leasingErr
			}
		} else {
			err = fmt.Errorf("农机租赁服务不可用")
		}
	} else {
		// 处理贷款申请
		result, err = s.submitAIDecisionQueryInternal(ctx, params)
	}

	// 计算处理时间
	duration := int(time.Since(startTime).Milliseconds())

	// 记录操作日志
	status := "SUCCESS"
	errorMessage := ""
	if err != nil {
		status = "ERROR"
		errorMessage = err.Error()
	}

	applicationType := "LOAN_APPLICATION"
	if s.isLeasingApplication(applicationID) {
		applicationType = "LEASING_APPLICATION"
	}

	// 从context中获取request（如果有的话）
	var req *http.Request
	if ctx != nil {
		if r, ok := ctx.Value("request").(*http.Request); ok {
			req = r
		}
	}

	s.LogAIAgentAction(
		"SUBMIT_AI_DECISION_UNIFIED",
		"DIFY_WORKFLOW",
		applicationID,
		map[string]interface{}{
			"application_type": applicationType,
			"decision":         params.Decision,
			"risk_score":       params.RiskScore,
		},
		result,
		status,
		errorMessage,
		duration,
		req,
	)

	return result, err
}

// SubmitAIDecisionQuery 处理查询参数方式的AI决策提交
func (s *AIAgentService) SubmitAIDecisionQuery(ctx context.Context, params *AIDecisionParams) (*AIDecisionResult, error) {
	startTime := time.Now()
	s.log.Info("Processing AI decision with query parameters",
		zap.String("application_id", params.ApplicationID),
		zap.String("decision", params.Decision),
		zap.Float64("risk_score", params.RiskScore))

	result, err := s.submitAIDecisionQueryInternal(ctx, params)

	// 计算处理时间
	duration := int(time.Since(startTime).Milliseconds())

	// 记录操作日志
	status := "SUCCESS"
	errorMessage := ""
	if err != nil {
		status = "ERROR"
		errorMessage = err.Error()
	}

	// 从context中获取request（如果有的话）
	var req *http.Request
	if ctx != nil {
		if r, ok := ctx.Value("request").(*http.Request); ok {
			req = r
		}
	}

	s.LogAIAgentAction(
		"SUBMIT_AI_DECISION_QUERY",
		"DIFY_WORKFLOW",
		params.ApplicationID,
		params,
		result,
		status,
		errorMessage,
		duration,
		req,
	)

	return result, err
}

// submitAIDecisionQueryInternal 内部方法，实际的业务逻辑
func (s *AIAgentService) submitAIDecisionQueryInternal(ctx context.Context, params *AIDecisionParams) (*AIDecisionResult, error) {
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
	if !pkg.Contains(allowedStatuses, application.Status) {
		return nil, fmt.Errorf("申请状态不允许AI决策: %s", application.Status)
	}

	// 3. 构建AI分析结果数据
	detailedAnalysisJSON, _ := json.Marshal(params.DetailedAnalysis)
	recommendationsJSON, _ := json.Marshal(params.Recommendations)
	riskFactorsJSON, _ := json.Marshal([]map[string]interface{}{
		{
			"factor":      "risk_score",
			"value":       params.RiskScore,
			"description": params.AnalysisSummary,
		},
	})

	// 4. 保存AI分析结果
	aiResult := &data.AIAnalysisResult{
		ApplicationID:         params.ApplicationID,
		WorkflowExecutionID:   params.WorkflowID,
		RiskLevel:             params.RiskLevel,
		RiskScore:             params.RiskScore,
		ConfidenceScore:       params.ConfidenceScore,
		AnalysisSummary:       params.AnalysisSummary,
		DetailedAnalysis:      detailedAnalysisJSON,
		RiskFactors:           riskFactorsJSON,
		Recommendations:       recommendationsJSON,
		AIDecision:            params.Decision,
		ApprovedAmount:        &params.ApprovedAmount,
		ApprovedTermMonths:    &params.ApprovedTermMonths,
		SuggestedInterestRate: params.SuggestedInterestRate,
		NextAction:            getNextAction(params.Decision),
		AIModelVersion:        params.AIModelVersion,
		ProcessingTimeMs:      2000, // 模拟处理时间
		ProcessedAt:           time.Now(),
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	// 5. 开始数据库事务
	tx := s.data.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 保存AI分析结果
	if err := tx.Create(aiResult).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("保存AI分析结果失败: %w", err)
	}

	// 6. 确定新状态
	newStatus := getNewStatusFromDecision(params.Decision)

	// 7. 更新申请状态
	updateData := map[string]interface{}{
		"status":        newStatus,
		"ai_risk_score": int(params.RiskScore * 1000), // 转换为千分制
		"ai_suggestion": params.AnalysisSummary,
		"updated_at":    time.Now(),
	}

	// 根据决策类型更新其他字段
	if params.Decision == "AUTO_APPROVED" {
		updateData["approved_amount"] = params.ApprovedAmount
		updateData["approved_term_months"] = params.ApprovedTermMonths
		updateData["final_decision"] = "APPROVED"
		updateData["decision_reason"] = "AI自动审批通过"
		updateData["processed_by"] = "AI_SYSTEM"
		updateData["processed_at"] = time.Now()
	} else if params.Decision == "AUTO_REJECTED" {
		updateData["final_decision"] = "REJECTED"
		updateData["decision_reason"] = params.AnalysisSummary
		updateData["processed_by"] = "AI_SYSTEM"
		updateData["processed_at"] = time.Now()
	}

	if err := tx.Model(&application).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新申请状态失败: %w", err)
	}

	// 8. 记录状态变更历史
	history := &data.LoanApplicationHistory{
		ApplicationID: params.ApplicationID,
		StatusFrom:    application.Status,
		StatusTo:      newStatus,
		OperatorType:  "AI_SYSTEM",
		OperatorID:    "ai_agent",
		Comments:      fmt.Sprintf("AI决策: %s, 风险评分: %.2f", params.Decision, params.RiskScore),
		OccurredAt:    time.Now(),
	}

	if err := tx.Create(history).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("记录状态历史失败: %w", err)
	}

	// 9. 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	// 10. 构建响应
	result := &AIDecisionResult{
		ApplicationID: params.ApplicationID,
		NewStatus:     newStatus,
		NextStep:      getNextAction(params.Decision),
	}

	s.log.Info("AI decision processed successfully",
		zap.String("application_id", params.ApplicationID),
		zap.String("decision", params.Decision),
		zap.String("new_status", newStatus))

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

// maskName 脱敏姓名
func maskName(name string) string {
	if len(name) == 0 {
		return ""
	}
	if len(name) <= 2 {
		return "*" + name[len(name)-1:]
	}
	return name[:1] + "*" + name[len(name)-1:]
}

// calculateAge 计算年龄
func calculateAge(birthDate *time.Time) int {
	if birthDate == nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}
