package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"backend/internal/data"
)

type MachineryLeasingApprovalService struct {
	data *data.Data
	log  *zap.Logger
}

func NewMachineryLeasingApprovalService(data *data.Data, log *zap.Logger) *MachineryLeasingApprovalService {
	return &MachineryLeasingApprovalService{
		data: data,
		log:  log,
	}
}

// MachineryLeasingApplicationInfo 农机租赁申请完整信息
type MachineryLeasingApplicationInfo struct {
	ApplicationID   string                  `json:"application_id"`
	LesseeInfo      LesseeInfo              `json:"lessee_info"`
	LessorInfo      LessorInfo              `json:"lessor_info"`
	MachineryInfo   MachineryInfo           `json:"machinery_info"`
	ApplicationInfo LeasingApplicationInfo  `json:"application_info"`
	ExternalData    LeasingExternalDataInfo `json:"external_data"`
}

type LesseeInfo struct {
	UserID       string `json:"user_id"`
	RealName     string `json:"real_name"`
	IDCardNumber string `json:"id_card_number"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	CreditRating string `json:"credit_rating"`
	IsVerified   bool   `json:"is_verified"`
}

type LessorInfo struct {
	UserID                 string  `json:"user_id"`
	RealName               string  `json:"real_name"`
	Phone                  string  `json:"phone"`
	VerificationStatus     string  `json:"verification_status"`
	CreditRating           string  `json:"credit_rating"`
	TotalMachineryCount    int     `json:"total_machinery_count"`
	SuccessfulLeasingCount int     `json:"successful_leasing_count"`
	AverageRating          float64 `json:"average_rating"`
}

type MachineryInfo struct {
	MachineryID  string  `json:"machinery_id"`
	Type         string  `json:"type"`
	BrandModel   string  `json:"brand_model"`
	DailyRent    float64 `json:"daily_rent"`
	Deposit      float64 `json:"deposit"`
	Status       string  `json:"status"`
	LocationText string  `json:"location_text"`
}

type LeasingApplicationInfo struct {
	RequestedStartDate time.Time `json:"requested_start_date"`
	RequestedEndDate   time.Time `json:"requested_end_date"`
	RentalDays         int       `json:"rental_days"`
	TotalAmount        float64   `json:"total_amount"`
	DepositAmount      float64   `json:"deposit_amount"`
	UsagePurpose       string    `json:"usage_purpose"`
	LesseeNotes        string    `json:"lessee_notes"`
	SubmittedAt        time.Time `json:"submitted_at"`
	Status             string    `json:"status"`
}

type LeasingExternalDataInfo struct {
	LesseeHistoryCount int     `json:"lessee_history_count"`
	LessorReliability  float64 `json:"lessor_reliability"`
	MachineryCondition string  `json:"machinery_condition"`
	SeasonalDemand     string  `json:"seasonal_demand"`
}

// LeasingAIDecisionRequest AI决策请求
type LeasingAIDecisionRequest struct {
	AIAnalysis     LeasingAIAnalysisData     `json:"ai_analysis"`
	AIDecision     LeasingAIDecisionData     `json:"ai_decision"`
	ProcessingInfo LeasingProcessingInfoData `json:"processing_info"`
}

type LeasingAIAnalysisData struct {
	RiskLevel        string                   `json:"risk_level"`
	RiskScore        float64                  `json:"risk_score"`
	ConfidenceScore  float64                  `json:"confidence_score"`
	AnalysisSummary  string                   `json:"analysis_summary"`
	DetailedAnalysis map[string]interface{}   `json:"detailed_analysis"`
	RiskFactors      []map[string]interface{} `json:"risk_factors"`
	Recommendations  []string                 `json:"recommendations"`
}

type LeasingAIDecisionData struct {
	Decision            string   `json:"decision"`
	SuggestedDeposit    float64  `json:"suggested_deposit"`
	SuggestedConditions []string `json:"suggested_conditions"`
	NextAction          string   `json:"next_action"`
}

type LeasingProcessingInfoData struct {
	AIModelVersion   string    `json:"ai_model_version"`
	ProcessingTimeMs int       `json:"processing_time_ms"`
	WorkflowID       string    `json:"workflow_id"`
	ProcessedAt      time.Time `json:"processed_at"`
}

// LeasingAIDecisionParams AI决策参数
type LeasingAIDecisionParams struct {
	ApplicationID       string                 `json:"application_id"`
	Decision            string                 `json:"decision"`
	RiskScore           float64                `json:"risk_score"`
	ConfidenceScore     float64                `json:"confidence_score"`
	SuggestedDeposit    float64                `json:"suggested_deposit"`
	SuggestedConditions []string               `json:"suggested_conditions"`
	RiskLevel           string                 `json:"risk_level"`
	AnalysisSummary     string                 `json:"analysis_summary"`
	DetailedAnalysis    map[string]interface{} `json:"detailed_analysis"`
	Recommendations     []string               `json:"recommendations"`
	AIModelVersion      string                 `json:"ai_model_version"`
	WorkflowID          string                 `json:"workflow_id"`
}

// LeasingAIDecisionResult AI决策结果
type LeasingAIDecisionResult struct {
	ApplicationID string `json:"application_id"`
	NewStatus     string `json:"new_status"`
	NextStep      string `json:"next_step"`
}

// LogLeasingAIAgentAction 记录农机租赁AI操作日志
func (s *MachineryLeasingApprovalService) LogLeasingAIAgentAction(actionType, agentType, applicationID string, requestData, responseData interface{}, status string, errorMessage string, duration int, req *http.Request) {
	logID := uuid.New().String()

	requestJSON, _ := json.Marshal(requestData)
	responseJSON, _ := json.Marshal(responseData)

	clientIP := getClientIP(req)
	userAgent := req.Header.Get("User-Agent")

	aiLog := &data.AIAgentLog{
		LogID:         logID,
		ApplicationID: applicationID,
		ActionType:    actionType,
		AgentType:     agentType,
		RequestData:   requestJSON,
		ResponseData:  responseJSON,
		Status:        status,
		ErrorMessage:  errorMessage,
		Duration:      duration,
		IPAddress:     clientIP,
		UserAgent:     userAgent,
		OccurredAt:    time.Now(),
		CreatedAt:     time.Now(),
	}

	if err := s.data.DB.Create(aiLog).Error; err != nil {
		s.log.Error("Failed to save leasing AI agent log",
			zap.String("logId", logID),
			zap.String("actionType", actionType),
			zap.Error(err))
	} else {
		s.log.Info("Leasing AI agent action logged successfully",
			zap.String("logId", logID),
			zap.String("actionType", actionType),
			zap.String("applicationId", applicationID))
	}
}

// GetLeasingApplicationInfoWithLog 获取农机租赁申请信息（带日志）
func (s *MachineryLeasingApprovalService) GetLeasingApplicationInfoWithLog(applicationID string, req *http.Request) (*MachineryLeasingApplicationInfo, error) {
	startTime := time.Now()

	requestData := map[string]string{"application_id": applicationID}

	info, err := s.getLeasingApplicationInfoInternal(applicationID)

	duration := int(time.Since(startTime).Milliseconds())

	if err != nil {
		s.LogLeasingAIAgentAction("GET_LEASING_APPLICATION_INFO", "DIFY_WORKFLOW", applicationID,
			requestData, nil, "ERROR", err.Error(), duration, req)
		return nil, err
	}

	s.LogLeasingAIAgentAction("GET_LEASING_APPLICATION_INFO", "DIFY_WORKFLOW", applicationID,
		requestData, info, "SUCCESS", "", duration, req)

	return info, nil
}

// getLeasingApplicationInfoInternal 内部获取农机租赁申请信息方法
func (s *MachineryLeasingApprovalService) getLeasingApplicationInfoInternal(applicationID string) (*MachineryLeasingApplicationInfo, error) {
	// 1. 获取租赁申请基本信息
	var application data.MachineryLeasingApplication
	if err := s.data.DB.Where("application_id = ?", applicationID).First(&application).Error; err != nil {
		return nil, fmt.Errorf("农机租赁申请不存在: %w", err)
	}

	// 2. 获取承租方信息
	var lesseeUser data.User
	var lesseeProfile data.UserProfile
	if err := s.data.DB.Where("user_id = ?", application.LesseeUserID).First(&lesseeUser).Error; err != nil {
		return nil, fmt.Errorf("承租方用户不存在: %w", err)
	}
	s.data.DB.Where("user_id = ?", application.LesseeUserID).First(&lesseeProfile)

	// 3. 获取出租方信息
	var lessorUser data.User
	var lessorQualification data.LessorQualification
	if err := s.data.DB.Where("user_id = ?", application.LessorUserID).First(&lessorUser).Error; err != nil {
		return nil, fmt.Errorf("出租方用户不存在: %w", err)
	}
	s.data.DB.Where("user_id = ?", application.LessorUserID).First(&lessorQualification)

	// 4. 获取农机信息
	var machinery data.FarmMachinery
	if err := s.data.DB.Where("machinery_id = ?", application.MachineryID).First(&machinery).Error; err != nil {
		return nil, fmt.Errorf("农机信息不存在: %w", err)
	}

	// 5. 构建响应数据
	info := &MachineryLeasingApplicationInfo{
		ApplicationID: applicationID,
		LesseeInfo: LesseeInfo{
			UserID:       lesseeUser.UserID,
			RealName:     maskName(lesseeProfile.RealName),
			IDCardNumber: maskIDCard(lesseeProfile.IDCardNumber),
			Phone:        maskPhone(lesseeUser.Phone),
			Address:      lesseeProfile.Address,
			CreditRating: "良好", // 模拟信用等级
			IsVerified:   true, // 模拟认证状态
		},
		LessorInfo: LessorInfo{
			UserID:                 lessorUser.UserID,
			RealName:               maskName(lessorQualification.RealName),
			Phone:                  maskPhone(lessorUser.Phone),
			VerificationStatus:     lessorQualification.VerificationStatus,
			CreditRating:           lessorQualification.CreditRating,
			TotalMachineryCount:    lessorQualification.TotalMachineryCount,
			SuccessfulLeasingCount: lessorQualification.SuccessfulLeasingCount,
			AverageRating:          *lessorQualification.AverageRating,
		},
		MachineryInfo: MachineryInfo{
			MachineryID:  machinery.MachineryID,
			Type:         machinery.Type,
			BrandModel:   machinery.BrandModel,
			DailyRent:    machinery.DailyRent,
			Deposit:      *machinery.Deposit,
			Status:       machinery.Status,
			LocationText: machinery.LocationText,
		},
		ApplicationInfo: LeasingApplicationInfo{
			RequestedStartDate: application.RequestedStartDate,
			RequestedEndDate:   application.RequestedEndDate,
			RentalDays:         application.RentalDays,
			TotalAmount:        application.TotalAmount,
			DepositAmount:      *application.DepositAmount,
			UsagePurpose:       application.UsagePurpose,
			LesseeNotes:        application.LesseeNotes,
			SubmittedAt:        application.SubmittedAt,
			Status:             application.ApplicationStatus,
		},
		ExternalData: LeasingExternalDataInfo{
			LesseeHistoryCount: 3,    // 模拟承租历史次数
			LessorReliability:  0.92, // 模拟出租方可靠性
			MachineryCondition: "良好", // 模拟农机状况
			SeasonalDemand:     "高",  // 模拟季节性需求
		},
	}

	return info, nil
}

// SubmitLeasingAIDecisionWithLog 提交农机租赁AI决策（带日志）
func (s *MachineryLeasingApprovalService) SubmitLeasingAIDecisionWithLog(applicationID string, request *LeasingAIDecisionRequest, req *http.Request) error {
	startTime := time.Now()

	err := s.submitLeasingAIDecisionInternal(applicationID, request)

	duration := int(time.Since(startTime).Milliseconds())

	if err != nil {
		s.LogLeasingAIAgentAction("SUBMIT_LEASING_AI_DECISION", "DIFY_WORKFLOW", applicationID,
			request, nil, "ERROR", err.Error(), duration, req)
		return err
	}

	result := LeasingAIDecisionResult{
		ApplicationID: applicationID,
		NewStatus:     getLeasingNewStatusFromDecision(request.AIDecision.Decision),
		NextStep:      getLeasingNextAction(request.AIDecision.Decision),
	}

	s.LogLeasingAIAgentAction("SUBMIT_LEASING_AI_DECISION", "DIFY_WORKFLOW", applicationID,
		request, result, "SUCCESS", "", duration, req)

	return nil
}

// submitLeasingAIDecisionInternal 内部提交农机租赁AI决策方法
func (s *MachineryLeasingApprovalService) submitLeasingAIDecisionInternal(applicationID string, request *LeasingAIDecisionRequest) error {
	tx := s.data.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 保存AI分析结果
	detailedAnalysisJSON, _ := json.Marshal(request.AIAnalysis.DetailedAnalysis)
	riskFactorsJSON, _ := json.Marshal(request.AIAnalysis.RiskFactors)
	recommendationsJSON, _ := json.Marshal(request.AIAnalysis.Recommendations)
	suggestedConditionsJSON, _ := json.Marshal(request.AIDecision.SuggestedConditions)

	aiResult := &data.MachineryLeasingAIResult{
		ApplicationID:        applicationID,
		WorkflowExecutionID:  request.ProcessingInfo.WorkflowID,
		RiskLevel:            request.AIAnalysis.RiskLevel,
		RiskScore:            request.AIAnalysis.RiskScore,
		ConfidenceScore:      request.AIAnalysis.ConfidenceScore,
		AnalysisSummary:      request.AIAnalysis.AnalysisSummary,
		DetailedAnalysis:     detailedAnalysisJSON,
		RiskFactors:          riskFactorsJSON,
		Recommendations:      recommendationsJSON,
		AIDecision:           request.AIDecision.Decision,
		SuggestedDepositRate: &request.AIDecision.SuggestedDeposit,
		SuggestedConditions:  suggestedConditionsJSON,
		NextAction:           request.AIDecision.NextAction,
		AIModelVersion:       request.ProcessingInfo.AIModelVersion,
		ProcessingTimeMs:     request.ProcessingInfo.ProcessingTimeMs,
		ProcessedAt:          request.ProcessingInfo.ProcessedAt,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := tx.Create(aiResult).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("保存AI分析结果失败: %w", err)
	}

	// 2. 更新申请状态
	newStatus := getLeasingNewStatusFromDecision(request.AIDecision.Decision)

	var application data.MachineryLeasingApplication
	if err := tx.Where("application_id = ?", applicationID).First(&application).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("农机租赁申请不存在: %w", err)
	}

	oldStatus := application.ApplicationStatus

	updates := map[string]interface{}{
		"application_status": newStatus,
		"risk_level":         request.AIAnalysis.RiskLevel,
		"ai_risk_score":      request.AIAnalysis.RiskScore,
		"ai_suggestion":      request.AIAnalysis.AnalysisSummary,
		"updated_at":         time.Now(),
	}

	if err := tx.Model(&application).Updates(updates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新申请状态失败: %w", err)
	}

	// 3. 记录状态变更历史
	history := &data.MachineryLeasingApprovalHistory{
		ApplicationID: applicationID,
		StatusFrom:    oldStatus,
		StatusTo:      newStatus,
		OperatorType:  "AI_SYSTEM",
		OperatorID:    "ai_agent",
		Comments:      fmt.Sprintf("AI决策: %s, 风险评分: %.2f", request.AIDecision.Decision, request.AIAnalysis.RiskScore),
		OccurredAt:    time.Now(),
	}

	if err := tx.Create(history).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("记录状态变更历史失败: %w", err)
	}

	return tx.Commit().Error
}

// SubmitLeasingAIDecisionQuery 查询参数方式提交农机租赁AI决策
func (s *MachineryLeasingApprovalService) SubmitLeasingAIDecisionQuery(ctx context.Context, params *LeasingAIDecisionParams) (*LeasingAIDecisionResult, error) {
	startTime := time.Now()

	result, err := s.submitLeasingAIDecisionQueryInternal(ctx, params)

	duration := int(time.Since(startTime).Milliseconds())

	// 创建模拟的HTTP请求用于日志记录
	req, _ := http.NewRequest("POST", "/api/ai-agent/machinery-leasing/decision-query", nil)
	req.Header.Set("User-Agent", "dify-workflow/1.0")
	req.RemoteAddr = "172.20.0.9:0"

	if err != nil {
		s.LogLeasingAIAgentAction("SUBMIT_LEASING_AI_DECISION_QUERY", "DIFY_WORKFLOW", params.ApplicationID,
			params, nil, "ERROR", err.Error(), duration, req)
		return nil, err
	}

	s.LogLeasingAIAgentAction("SUBMIT_LEASING_AI_DECISION_QUERY", "DIFY_WORKFLOW", params.ApplicationID,
		params, result, "SUCCESS", "", duration, req)

	return result, nil
}

// submitLeasingAIDecisionQueryInternal 内部查询参数方式提交农机租赁AI决策
func (s *MachineryLeasingApprovalService) submitLeasingAIDecisionQueryInternal(ctx context.Context, params *LeasingAIDecisionParams) (*LeasingAIDecisionResult, error) {
	tx := s.data.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 保存AI分析结果
	detailedAnalysisJSON, _ := json.Marshal(params.DetailedAnalysis)
	riskFactorsJSON, _ := json.Marshal([]map[string]interface{}{
		{
			"factor":      "risk_score",
			"value":       params.RiskScore,
			"description": params.AnalysisSummary,
		},
	})
	recommendationsJSON, _ := json.Marshal(params.Recommendations)
	suggestedConditionsJSON, _ := json.Marshal(params.SuggestedConditions)

	aiResult := &data.MachineryLeasingAIResult{
		ApplicationID:        params.ApplicationID,
		WorkflowExecutionID:  params.WorkflowID,
		RiskLevel:            params.RiskLevel,
		RiskScore:            params.RiskScore,
		ConfidenceScore:      params.ConfidenceScore,
		AnalysisSummary:      params.AnalysisSummary,
		DetailedAnalysis:     detailedAnalysisJSON,
		RiskFactors:          riskFactorsJSON,
		Recommendations:      recommendationsJSON,
		AIDecision:           params.Decision,
		SuggestedDepositRate: &params.SuggestedDeposit,
		SuggestedConditions:  suggestedConditionsJSON,
		NextAction:           getLeasingNextAction(params.Decision),
		AIModelVersion:       params.AIModelVersion,
		ProcessingTimeMs:     2000, // 模拟处理时间
		ProcessedAt:          time.Now(),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := tx.Create(aiResult).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("保存AI分析结果失败: %w", err)
	}

	// 2. 更新申请状态
	newStatus := getLeasingNewStatusFromDecision(params.Decision)

	var application data.MachineryLeasingApplication
	if err := tx.Where("application_id = ?", params.ApplicationID).First(&application).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("农机租赁申请不存在: %w", err)
	}

	oldStatus := application.ApplicationStatus

	updates := map[string]interface{}{
		"application_status": newStatus,
		"risk_level":         params.RiskLevel,
		"ai_risk_score":      params.RiskScore,
		"ai_suggestion":      params.AnalysisSummary,
		"updated_at":         time.Now(),
	}

	if err := tx.Model(&application).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新申请状态失败: %w", err)
	}

	// 3. 记录状态变更历史
	history := &data.MachineryLeasingApprovalHistory{
		ApplicationID: params.ApplicationID,
		StatusFrom:    oldStatus,
		StatusTo:      newStatus,
		OperatorType:  "AI_SYSTEM",
		OperatorID:    "ai_agent",
		Comments:      fmt.Sprintf("AI决策: %s, 风险评分: %.2f", params.Decision, params.RiskScore),
		OccurredAt:    time.Now(),
	}

	if err := tx.Create(history).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("记录状态变更历史失败: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	result := &LeasingAIDecisionResult{
		ApplicationID: params.ApplicationID,
		NewStatus:     newStatus,
		NextStep:      getLeasingNextAction(params.Decision),
	}

	return result, nil
}

// UpdateLeasingApplicationStatus 更新农机租赁申请状态
func (s *MachineryLeasingApprovalService) UpdateLeasingApplicationStatus(applicationID, status, operator, remarks string, metadata map[string]interface{}) error {
	var application data.MachineryLeasingApplication
	if err := s.data.DB.Where("application_id = ?", applicationID).First(&application).Error; err != nil {
		return fmt.Errorf("农机租赁申请不存在: %w", err)
	}

	oldStatus := application.ApplicationStatus

	tx := s.data.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新申请状态
	updates := map[string]interface{}{
		"application_status": status,
		"updated_at":         time.Now(),
	}

	if status == "APPROVED" || status == "REJECTED" {
		updates["approved_by"] = operator
		updates["approved_at"] = time.Now()
		if remarks != "" {
			updates["approval_comments"] = remarks
		}
	}

	if err := tx.Model(&application).Updates(updates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新申请状态失败: %w", err)
	}

	// 记录状态变更历史
	history := &data.MachineryLeasingApprovalHistory{
		ApplicationID: applicationID,
		StatusFrom:    oldStatus,
		StatusTo:      status,
		OperatorType:  "HUMAN",
		OperatorID:    operator,
		Comments:      remarks,
		OccurredAt:    time.Now(),
	}

	if err := tx.Create(history).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("记录状态变更历史失败: %w", err)
	}

	return tx.Commit().Error
}

// 获取农机租赁新状态
func getLeasingNewStatusFromDecision(decision string) string {
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
		return "MANUAL_REVIEW_REQUIRED"
	}
}

// 获取农机租赁下一步动作
func getLeasingNextAction(decision string) string {
	switch decision {
	case "AUTO_APPROVE":
		return "CREATE_LEASING_ORDER"
	case "AUTO_REJECT":
		return "NOTIFY_REJECTION"
	case "REQUIRE_HUMAN_REVIEW":
		return "ASSIGN_TO_REVIEWER"
	case "REQUIRE_DEPOSIT_ADJUSTMENT":
		return "ADJUST_DEPOSIT_TERMS"
	default:
		return "ASSIGN_TO_REVIEWER"
	}
}
