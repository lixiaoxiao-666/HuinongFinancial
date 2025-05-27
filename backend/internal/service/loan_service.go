package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"backend/internal/data"
	"backend/pkg"

	"go.uber.org/zap"
)

// LoanService 贷款服务
type LoanService struct {
	data *data.Data
	log  *zap.Logger
}

// NewLoanService 创建贷款服务
func NewLoanService(data *data.Data, log *zap.Logger) *LoanService {
	return &LoanService{
		data: data,
		log:  log,
	}
}

// LoanProductResponse 贷款产品响应
type LoanProductResponse struct {
	ProductID             string      `json:"product_id"`
	Name                  string      `json:"name"`
	Description           string      `json:"description"`
	Category              string      `json:"category"`
	MinAmount             float64     `json:"min_amount"`
	MaxAmount             float64     `json:"max_amount"`
	MinTermMonths         int         `json:"min_term_months"`
	MaxTermMonths         int         `json:"max_term_months"`
	InterestRateYearly    string      `json:"interest_rate_yearly"`
	RepaymentMethods      interface{} `json:"repayment_methods"`
	ApplicationConditions string      `json:"application_conditions,omitempty"`
	RequiredDocuments     interface{} `json:"required_documents,omitempty"`
}

// SubmitLoanApplicationRequest 提交贷款申请请求
type SubmitLoanApplicationRequest struct {
	ProductID         string                 `json:"product_id" binding:"required"`
	Amount            float64                `json:"amount" binding:"required,gt=0"`
	TermMonths        int                    `json:"term_months" binding:"required,gt=0"`
	Purpose           string                 `json:"purpose" binding:"required"`
	ApplicantInfo     map[string]interface{} `json:"applicant_info" binding:"required"`
	UploadedDocuments []DocumentInfo         `json:"uploaded_documents"`
}

// DocumentInfo 文档信息
type DocumentInfo struct {
	DocType string `json:"doc_type" binding:"required"`
	FileID  string `json:"file_id" binding:"required"`
}

// LoanApplicationResponse 贷款申请响应
type LoanApplicationResponse struct {
	ApplicationID       string                 `json:"application_id"`
	UserID              string                 `json:"user_id"`
	ProductID           string                 `json:"product_id"`
	AmountApplied       float64                `json:"amount_applied"`
	TermMonthsApplied   int                    `json:"term_months_applied"`
	Purpose             string                 `json:"purpose"`
	Status              string                 `json:"status"`
	ApplicantSnapshot   map[string]interface{} `json:"applicant_snapshot,omitempty"`
	SubmittedAt         time.Time              `json:"submitted_at"`
	AIRiskScore         *int                   `json:"ai_risk_score,omitempty"`
	AISuggestion        string                 `json:"ai_suggestion,omitempty"`
	ApprovedAmount      *float64               `json:"approved_amount,omitempty"`
	ApprovedTermMonths  *int                   `json:"approved_term_months,omitempty"`
	FinalDecision       string                 `json:"final_decision,omitempty"`
	DecisionReason      string                 `json:"decision_reason,omitempty"`
	ProcessedAt         *time.Time             `json:"processed_at,omitempty"`
	History             []HistoryItem          `json:"history,omitempty"`
}

// HistoryItem 历史记录项
type HistoryItem struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Operator  string    `json:"operator"`
	Comments  string    `json:"comments,omitempty"`
}

// GetLoanProducts 获取贷款产品列表
func (s *LoanService) GetLoanProducts(ctx context.Context, category string) ([]LoanProductResponse, error) {
	var products []data.LoanProduct
	query := s.data.DB.Where("status = ?", 0) // 只查询有效产品
	
	if category != "" {
		query = query.Where("category = ?", category)
	}
	
	if err := query.Find(&products).Error; err != nil {
		s.log.Error("查询贷款产品失败", zap.Error(err))
		return nil, errors.New("查询产品失败")
	}

	var result []LoanProductResponse
	for _, product := range products {
		var repaymentMethods interface{}
		if product.RepaymentMethods != nil {
			json.Unmarshal(product.RepaymentMethods, &repaymentMethods)
		}

		result = append(result, LoanProductResponse{
			ProductID:          product.ProductID,
			Name:               product.Name,
			Description:        product.Description,
			Category:           product.Category,
			MinAmount:          product.MinAmount,
			MaxAmount:          product.MaxAmount,
			MinTermMonths:      product.MinTermMonths,
			MaxTermMonths:      product.MaxTermMonths,
			InterestRateYearly: product.InterestRateYearly,
			RepaymentMethods:   repaymentMethods,
		})
	}

	return result, nil
}

// GetLoanProduct 获取贷款产品详情
func (s *LoanService) GetLoanProduct(ctx context.Context, productID string) (*LoanProductResponse, error) {
	var product data.LoanProduct
	if err := s.data.DB.Where("product_id = ? AND status = ?", productID, 0).First(&product).Error; err != nil {
		return nil, errors.New("产品不存在")
	}

	var repaymentMethods, requiredDocuments interface{}
	if product.RepaymentMethods != nil {
		json.Unmarshal(product.RepaymentMethods, &repaymentMethods)
	}
	if product.RequiredDocuments != nil {
		json.Unmarshal(product.RequiredDocuments, &requiredDocuments)
	}

	return &LoanProductResponse{
		ProductID:             product.ProductID,
		Name:                  product.Name,
		Description:           product.Description,
		Category:              product.Category,
		MinAmount:             product.MinAmount,
		MaxAmount:             product.MaxAmount,
		MinTermMonths:         product.MinTermMonths,
		MaxTermMonths:         product.MaxTermMonths,
		InterestRateYearly:    product.InterestRateYearly,
		RepaymentMethods:      repaymentMethods,
		ApplicationConditions: product.ApplicationConditions,
		RequiredDocuments:     requiredDocuments,
	}, nil
}

// SubmitLoanApplication 提交贷款申请
func (s *LoanService) SubmitLoanApplication(ctx context.Context, userID string, req *SubmitLoanApplicationRequest) (*LoanApplicationResponse, error) {
	// 验证产品是否存在
	var product data.LoanProduct
	if err := s.data.DB.Where("product_id = ? AND status = ?", req.ProductID, 0).First(&product).Error; err != nil {
		return nil, errors.New("贷款产品不存在")
	}

	// 验证申请金额和期限
	if req.Amount < product.MinAmount || req.Amount > product.MaxAmount {
		return nil, errors.New("申请金额超出产品范围")
	}
	if req.TermMonths < product.MinTermMonths || req.TermMonths > product.MaxTermMonths {
		return nil, errors.New("申请期限超出产品范围")
	}

	// 序列化申请人信息
	applicantSnapshot, _ := json.Marshal(req.ApplicantInfo)

	// 创建贷款申请
	application := data.LoanApplication{
		ApplicationID:     pkg.GenerateApplicationID(),
		UserID:            userID,
		ProductID:         req.ProductID,
		AmountApplied:     req.Amount,
		TermMonthsApplied: req.TermMonths,
		Purpose:           req.Purpose,
		Status:            "SUBMITTED",
		ApplicantSnapshot: applicantSnapshot,
		SubmittedAt:       time.Now(),
	}

	if err := s.data.DB.Create(&application).Error; err != nil {
		s.log.Error("创建贷款申请失败", zap.Error(err))
		return nil, errors.New("提交申请失败")
	}

	// 记录申请历史
	history := data.LoanApplicationHistory{
		ApplicationID: application.ApplicationID,
		StatusTo:      "SUBMITTED",
		OperatorType:  "USER",
		OperatorID:    userID,
		Comments:      "用户提交申请",
		OccurredAt:    time.Now(),
	}
	s.data.DB.Create(&history)

	// 处理文档关联
	for _, doc := range req.UploadedDocuments {
		s.data.DB.Model(&data.UploadedFile{}).
			Where("file_id = ? AND user_id = ?", doc.FileID, userID).
			Update("related_id", application.ApplicationID)
	}

	// 异步触发AI审批流程
	go s.triggerAIProcessing(application.ApplicationID)

	return &LoanApplicationResponse{
		ApplicationID:     application.ApplicationID,
		UserID:            application.UserID,
		ProductID:         application.ProductID,
		AmountApplied:     application.AmountApplied,
		TermMonthsApplied: application.TermMonthsApplied,
		Purpose:           application.Purpose,
		Status:            application.Status,
		SubmittedAt:       application.SubmittedAt,
	}, nil
}

// GetLoanApplication 获取贷款申请详情
func (s *LoanService) GetLoanApplication(ctx context.Context, applicationID string, userID string) (*LoanApplicationResponse, error) {
	var application data.LoanApplication
	query := s.data.DB.Where("application_id = ?", applicationID)
	
	// 如果指定了用户ID，则只能查看自己的申请
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	
	if err := query.First(&application).Error; err != nil {
		return nil, errors.New("申请不存在")
	}

	// 获取申请历史
	var historyRecords []data.LoanApplicationHistory
	s.data.DB.Where("application_id = ?", applicationID).
		Order("occurred_at ASC").
		Find(&historyRecords)

	var history []HistoryItem
	for _, record := range historyRecords {
		history = append(history, HistoryItem{
			Status:    record.StatusTo,
			Timestamp: record.OccurredAt,
			Operator:  record.OperatorType,
			Comments:  record.Comments,
		})
	}

	var applicantSnapshot map[string]interface{}
	if application.ApplicantSnapshot != nil {
		json.Unmarshal(application.ApplicantSnapshot, &applicantSnapshot)
	}

	return &LoanApplicationResponse{
		ApplicationID:      application.ApplicationID,
		UserID:             application.UserID,
		ProductID:          application.ProductID,
		AmountApplied:      application.AmountApplied,
		TermMonthsApplied:  application.TermMonthsApplied,
		Purpose:            application.Purpose,
		Status:             application.Status,
		ApplicantSnapshot:  applicantSnapshot,
		SubmittedAt:        application.SubmittedAt,
		AIRiskScore:        application.AIRiskScore,
		AISuggestion:       application.AISuggestion,
		ApprovedAmount:     application.ApprovedAmount,
		ApprovedTermMonths: application.ApprovedTermMonths,
		FinalDecision:      application.FinalDecision,
		DecisionReason:     application.DecisionReason,
		ProcessedAt:        application.ProcessedAt,
		History:            history,
	}, nil
}

// GetMyLoanApplications 获取我的贷款申请列表
func (s *LoanService) GetMyLoanApplications(ctx context.Context, userID string, status string, page, limit int) ([]LoanApplicationResponse, int64, error) {
	offset, validLimit := pkg.GetPagination(page, limit)
	
	query := s.data.DB.Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&data.LoanApplication{}).Count(&total)

	var applications []data.LoanApplication
	if err := query.Offset(offset).Limit(validLimit).
		Order("submitted_at DESC").
		Find(&applications).Error; err != nil {
		s.log.Error("查询贷款申请失败", zap.Error(err))
		return nil, 0, errors.New("查询失败")
	}

	var result []LoanApplicationResponse
	for _, app := range applications {
		result = append(result, LoanApplicationResponse{
			ApplicationID:     app.ApplicationID,
			UserID:            app.UserID,
			ProductID:         app.ProductID,
			AmountApplied:     app.AmountApplied,
			TermMonthsApplied: app.TermMonthsApplied,
			Purpose:           app.Purpose,
			Status:            app.Status,
			SubmittedAt:       app.SubmittedAt,
			AIRiskScore:       app.AIRiskScore,
			AISuggestion:      app.AISuggestion,
			FinalDecision:     app.FinalDecision,
			ProcessedAt:       app.ProcessedAt,
		})
	}

	return result, total, nil
}

// triggerAIProcessing 触发AI处理流程（模拟实现）
func (s *LoanService) triggerAIProcessing(applicationID string) {
	// 模拟AI处理延迟
	time.Sleep(2 * time.Second)

	// 更新申请状态为AI处理中
	s.data.DB.Model(&data.LoanApplication{}).
		Where("application_id = ?", applicationID).
		Updates(map[string]interface{}{
			"status": "AI_PROCESSING",
		})

	// 记录状态变更历史
	history := data.LoanApplicationHistory{
		ApplicationID: applicationID,
		StatusFrom:    "SUBMITTED",
		StatusTo:      "AI_PROCESSING",
		OperatorType:  "SYSTEM",
		OperatorID:    "SYSTEM_AI",
		Comments:      "AI系统开始处理申请",
		OccurredAt:    time.Now(),
	}
	s.data.DB.Create(&history)

	s.log.Info("AI处理流程已触发", zap.String("applicationId", applicationID))
} 