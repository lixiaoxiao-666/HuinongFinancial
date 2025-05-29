package service

import (
	"context"
	"fmt"
	"time"

	"huinong-backend/internal/model"
	"huinong-backend/internal/repository"
)

// loanService 贷款服务实现
type loanService struct {
	loanRepo repository.LoanRepository
	userRepo repository.UserRepository
}

// NewLoanService 创建贷款服务实例
func NewLoanService(
	loanRepo repository.LoanRepository,
	userRepo repository.UserRepository,
) LoanService {
	return &loanService{
		loanRepo: loanRepo,
		userRepo: userRepo,
	}
}

// ==================== 贷款产品查询 ====================

// GetProducts 获取贷款产品列表
func (s *loanService) GetProducts(ctx context.Context, req *GetProductsRequest) (*GetProductsResponse, error) {
	// TODO: 从上下文获取用户类型
	userType := "farmer" // 临时默认值

	products, err := s.loanRepo.GetActiveProducts(ctx, userType)
	if err != nil {
		return nil, fmt.Errorf("获取贷款产品失败: %v", err)
	}

	// 转换为响应格式
	var productResponses []*ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, &ProductResponse{
			ID:               uint(product.ID),
			ProductName:      product.ProductName,
			ProductCode:      product.ProductCode,
			ProductType:      product.ProductType,
			Description:      product.Description,
			MinAmount:        product.MinAmount,
			MaxAmount:        product.MaxAmount,
			InterestRate:     product.InterestRate,
			TermMonths:       product.TermMonths,
			RequiredAuth:     product.RequiredAuth,
			IsActive:         product.IsActive,
			EligibleUserType: product.EligibleUserType,
		})
	}

	return &GetProductsResponse{
		Products: productResponses,
		Total:    int64(len(productResponses)),
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}

// ==================== 贷款申请处理 ====================

// CreateApplication 创建贷款申请
func (s *loanService) CreateApplication(ctx context.Context, req *CreateApplicationRequest) (*CreateApplicationResponse, error) {
	// TODO: 从上下文获取用户ID
	userID := uint64(1) // 临时默认值

	// 验证产品是否存在
	product, err := s.loanRepo.GetProductByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("产品不存在: %v", err)
	}

	// 验证申请金额
	if req.LoanAmount < product.MinAmount || req.LoanAmount > product.MaxAmount {
		return nil, fmt.Errorf("申请金额超出产品限制范围")
	}

	// 生成申请编号
	applicationNo := generateApplicationNo()

	// 创建申请记录
	application := &model.LoanApplication{
		ApplicationNo:        applicationNo,
		UserID:               userID,
		ProductID:            uint64(req.ProductID),
		ApplyAmount:          req.LoanAmount,
		LoanAmount:           req.LoanAmount,
		ApplyTermMonths:      req.TermMonths,
		TermMonths:           req.TermMonths,
		LoanPurpose:          req.LoanPurpose,
		ContactPhone:         req.ContactPhone,
		ApplicantPhone:       req.ContactPhone,
		ContactEmail:         req.ContactEmail,
		MaterialsJSON:        req.MaterialsJSON,
		ApplicationDocuments: req.MaterialsJSON,
		Status:               "pending",
		SubmittedAt:          time.Now(),
		Remarks:              req.Remarks,
	}

	err = s.loanRepo.CreateApplication(ctx, application)
	if err != nil {
		return nil, fmt.Errorf("创建贷款申请失败: %v", err)
	}

	// TODO: 触发AI审批流程
	// go s.TriggerAIAssessment(context.Background(), application.ID, "loan_approval")

	return &CreateApplicationResponse{
		ID:            uint(application.ID),
		ApplicationNo: application.ApplicationNo,
		Status:        application.Status,
		CreatedAt:     application.CreatedAt,
	}, nil
}

// GetApplicationDetails 获取申请详情
func (s *loanService) GetApplicationDetails(ctx context.Context, req *GetApplicationDetailsRequest) (*GetApplicationDetailsResponse, error) {
	application, err := s.loanRepo.GetApplicationByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("获取申请详情失败: %v", err)
	}

	// 获取产品信息
	product, err := s.loanRepo.GetProductByID(ctx, uint(application.ProductID))
	if err != nil {
		return nil, fmt.Errorf("获取产品信息失败: %v", err)
	}

	// 获取审批日志
	approvalLogs, _ := s.loanRepo.GetApprovalLogs(ctx, uint(application.ID))

	// 获取Dify工作流日志
	difyLogs, _ := s.loanRepo.GetDifyLogs(ctx, uint(application.ID))

	// 转换审批日志
	var approvalLogResponses []*ApprovalLogResponse
	for _, log := range approvalLogs {
		approvalLogResponses = append(approvalLogResponses, &ApprovalLogResponse{
			ID:        uint(log.ID),
			Step:      log.Step,
			Status:    log.Status,
			Note:      log.Note,
			CreatedAt: log.CreatedAt,
		})
	}

	// 转换Dify日志
	var difyLogResponses []*DifyLogResponse
	for _, log := range difyLogs {
		difyLogResponses = append(difyLogResponses, &DifyLogResponse{
			ID:           uint(log.ID),
			WorkflowType: log.WorkflowType,
			Status:       log.Status,
			Result:       log.Result,
			CreatedAt:    log.CreatedAt,
		})
	}

	return &GetApplicationDetailsResponse{
		Application: &ApplicationDetailsResponse{
			ID:            uint(application.ID),
			ApplicationNo: application.ApplicationNo,
			ProductID:     uint(application.ProductID),
			LoanAmount:    application.LoanAmount,
			Status:        application.Status,
			CreatedAt:     application.CreatedAt,
		},
		Product: &ProductResponse{
			ID:           uint(product.ID),
			ProductName:  product.ProductName,
			ProductCode:  product.ProductCode,
			ProductType:  product.ProductType,
			Description:  product.Description,
			MinAmount:    product.MinAmount,
			MaxAmount:    product.MaxAmount,
			InterestRate: product.InterestRate,
			TermMonths:   product.TermMonths,
			RequiredAuth: product.RequiredAuth,
			IsActive:     product.IsActive,
		},
		ApprovalLogs: approvalLogResponses,
		DifyLogs:     difyLogResponses,
	}, nil
}

// GetUserApplications 获取用户申请列表
func (s *loanService) GetUserApplications(ctx context.Context, req *GetUserApplicationsRequest) (*GetUserApplicationsResponse, error) {
	// TODO: 从上下文获取用户ID
	userID := uint(1) // 临时默认值

	applications, total, err := s.loanRepo.GetUserApplications(ctx, userID, req.Page, req.Limit, req.Status)
	if err != nil {
		return nil, fmt.Errorf("获取用户申请列表失败: %v", err)
	}

	var responses []*UserApplicationResponse
	for _, app := range applications {
		// 获取产品名称
		product, _ := s.loanRepo.GetProductByID(ctx, uint(app.ProductID))
		productName := "未知产品"
		if product != nil {
			productName = product.ProductName
		}

		responses = append(responses, &UserApplicationResponse{
			ID:            uint(app.ID),
			ApplicationNo: app.ApplicationNo,
			ProductName:   productName,
			LoanAmount:    app.LoanAmount,
			Status:        app.Status,
			CreatedAt:     app.CreatedAt,
		})
	}

	return &GetUserApplicationsResponse{
		Applications: responses,
		Total:        total,
		Page:         req.Page,
		Limit:        req.Limit,
	}, nil
}

// ==================== 审批流程 ====================

// ApproveApplication 批准申请
func (s *loanService) ApproveApplication(ctx context.Context, req *ApproveApplicationRequest) error {
	application, err := s.loanRepo.GetApplicationByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("获取申请信息失败: %v", err)
	}

	// 更新申请状态
	application.Status = "approved"
	if application.FinalApprovedAt == nil {
		now := time.Now()
		application.FinalApprovedAt = &now
	}

	err = s.loanRepo.UpdateApplication(ctx, application)
	if err != nil {
		return fmt.Errorf("更新申请状态失败: %v", err)
	}

	// TODO: 记录审批日志
	// TODO: 发送通知

	return nil
}

// RejectApplication 拒绝申请
func (s *loanService) RejectApplication(ctx context.Context, req *RejectApplicationRequest) error {
	application, err := s.loanRepo.GetApplicationByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("获取申请信息失败: %v", err)
	}

	// 更新申请状态
	application.Status = "rejected"
	application.RejectionReason = req.RejectionNote

	err = s.loanRepo.UpdateApplication(ctx, application)
	if err != nil {
		return fmt.Errorf("更新申请状态失败: %v", err)
	}

	// TODO: 记录审批日志
	// TODO: 发送通知

	return nil
}

// GetAdminApplications 获取管理员申请列表
func (s *loanService) GetAdminApplications(ctx context.Context, req *GetAdminApplicationsRequest) (*GetAdminApplicationsResponse, error) {
	applications, total, err := s.loanRepo.GetApplicationsForAdmin(ctx, req.Page, req.Limit, req.Status)
	if err != nil {
		return nil, fmt.Errorf("获取申请列表失败: %v", err)
	}

	var responses []*AdminApplicationResponse
	for _, app := range applications {
		// 获取用户名和产品名
		user, _ := s.userRepo.GetByID(ctx, app.UserID)
		product, _ := s.loanRepo.GetProductByID(ctx, uint(app.ProductID))

		userName := "未知用户"
		if user != nil {
			userName = user.RealName
			if userName == "" {
				userName = user.Username
			}
		}

		productName := "未知产品"
		if product != nil {
			productName = product.ProductName
		}

		responses = append(responses, &AdminApplicationResponse{
			ID:            uint(app.ID),
			ApplicationNo: app.ApplicationNo,
			UserName:      userName,
			ProductName:   productName,
			LoanAmount:    app.LoanAmount,
			Status:        app.Status,
			CreatedAt:     app.CreatedAt,
		})
	}

	return &GetAdminApplicationsResponse{
		Applications: responses,
		Total:        total,
		Page:         req.Page,
		Limit:        req.Limit,
	}, nil
}

// GetStatistics 获取贷款统计
func (s *loanService) GetStatistics(ctx context.Context) (*LoanStatisticsResponse, error) {
	// TODO: 实现统计查询
	stats, err := s.loanRepo.GetApplicationStatistics(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取统计数据失败: %v", err)
	}

	return &LoanStatisticsResponse{
		TotalApplications:   stats["total_applications"].(int64),
		MonthlyApplications: stats["monthly_applications"].(int64),
		StatusStatistics:    stats["status_statistics"],
	}, nil
}

// ==================== 产品管理 ====================

// CreateProduct 创建贷款产品
func (s *loanService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*model.LoanProduct, error) {
	product := &model.LoanProduct{
		ProductCode:      req.ProductCode,
		ProductName:      req.ProductName,
		ProductType:      req.ProductType,
		Description:      req.Description,
		MinAmount:        req.MinAmount,
		MaxAmount:        req.MaxAmount,
		InterestRate:     req.InterestRate,
		TermMonths:       req.TermMonths,
		RequiredAuth:     req.RequiredAuth,
		IsActive:         req.IsActive,
		SortOrder:        req.SortOrder,
		DifyWorkflowID:   req.DifyWorkflowID,
		EligibleUserType: req.EligibleUserType,
	}

	err := s.loanRepo.CreateProduct(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("创建产品失败: %v", err)
	}

	return product, nil
}

// UpdateProduct 更新产品
func (s *loanService) UpdateProduct(ctx context.Context, productID uint64, req *UpdateProductRequest) error {
	product, err := s.loanRepo.GetProductByID(ctx, uint(productID))
	if err != nil {
		return fmt.Errorf("获取产品失败: %v", err)
	}

	// 更新字段
	if req.ProductName != "" {
		product.ProductName = req.ProductName
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.MinAmount > 0 {
		product.MinAmount = req.MinAmount
	}
	if req.MaxAmount > 0 {
		product.MaxAmount = req.MaxAmount
	}
	if req.InterestRate > 0 {
		product.InterestRate = req.InterestRate
	}
	if req.TermMonths > 0 {
		product.TermMonths = req.TermMonths
	}
	if req.RequiredAuth != "" {
		product.RequiredAuth = req.RequiredAuth
	}
	product.IsActive = req.IsActive
	if req.SortOrder > 0 {
		product.SortOrder = req.SortOrder
	}
	if req.DifyWorkflowID != "" {
		product.DifyWorkflowID = req.DifyWorkflowID
	}
	if req.EligibleUserType != "" {
		product.EligibleUserType = req.EligibleUserType
	}

	return s.loanRepo.UpdateProduct(ctx, product)
}

// DeleteProduct 删除产品
func (s *loanService) DeleteProduct(ctx context.Context, productID uint64) error {
	return s.loanRepo.DeleteProduct(ctx, uint(productID))
}

// GetProduct 获取产品详情
func (s *loanService) GetProduct(ctx context.Context, productID uint64) (*model.LoanProduct, error) {
	return s.loanRepo.GetProductByID(ctx, uint(productID))
}

// ListProducts 获取产品列表
func (s *loanService) ListProducts(ctx context.Context, req *repository.ListProductsRequest) (*repository.ListProductsResponse, error) {
	return s.loanRepo.ListProducts(ctx, req)
}

// GetActiveProducts 获取活跃产品
func (s *loanService) GetActiveProducts(ctx context.Context, userType string) ([]*model.LoanProduct, error) {
	return s.loanRepo.GetActiveProducts(ctx, userType)
}

// ==================== 高级审批功能 ====================

// StartReview 开始人工审核
func (s *loanService) StartReview(ctx context.Context, applicationID uint64, reviewerID uint64) error {
	application, err := s.loanRepo.GetApplicationByID(ctx, uint(applicationID))
	if err != nil {
		return fmt.Errorf("获取申请失败: %v", err)
	}

	application.Status = "reviewing"
	application.CurrentApprover = &reviewerID

	return s.loanRepo.UpdateApplication(ctx, application)
}

// ReturnApplication 退回申请
func (s *loanService) ReturnApplication(ctx context.Context, applicationID uint64, req *ReturnRequest) error {
	application, err := s.loanRepo.GetApplicationByID(ctx, uint(applicationID))
	if err != nil {
		return fmt.Errorf("获取申请失败: %v", err)
	}

	application.Status = "returned"

	return s.loanRepo.UpdateApplication(ctx, application)
}

// ==================== AI集成功能 ====================

// TriggerAIAssessment 触发AI评估
func (s *loanService) TriggerAIAssessment(ctx context.Context, applicationID uint64, workflowType string) error {
	// TODO: 调用Dify工作流API
	// 1. 准备申请数据
	// 2. 调用Dify API
	// 3. 记录调用日志

	return fmt.Errorf("AI评估功能尚未实现")
}

// ProcessDifyCallback 处理Dify回调
func (s *loanService) ProcessDifyCallback(ctx context.Context, req *DifyCallbackRequest) error {
	// TODO: 处理Dify工作流回调
	// 1. 验证回调签名
	// 2. 更新申请状态
	// 3. 记录结果

	return fmt.Errorf("Dify回调处理功能尚未实现")
}

// ==================== 统计报表功能 ====================

// GetLoanStatistics 获取贷款统计
func (s *loanService) GetLoanStatistics(ctx context.Context, req *StatisticsRequest) (*LoanStatistics, error) {
	// TODO: 实现详细统计
	return &LoanStatistics{
		TotalApplications:     0,
		PendingApplications:   0,
		ApprovedApplications:  0,
		RejectedApplications:  0,
		TotalLoanAmount:       0,
		ApprovedLoanAmount:    0,
		ApplicationsByProduct: make(map[string]int64),
		ApprovalRate:          0.0,
	}, nil
}

// GenerateApprovalReport 生成审批报告
func (s *loanService) GenerateApprovalReport(ctx context.Context, req *ReportRequest) (*ApprovalReport, error) {
	// TODO: 实现报告生成
	return &ApprovalReport{
		TotalApplications:    0,
		ApprovedApplications: 0,
		RejectedApplications: 0,
		PendingApplications:  0,
		ReportData:           []ReportDataItem{},
		GeneratedAt:          time.Now(),
	}, nil
}

// ==================== 辅助方法 ====================

// CreateApprovalLog 创建审批日志
func (s *loanService) CreateApprovalLog(log *model.ApprovalLog) error {
	ctx := context.Background()
	return s.loanRepo.CreateApprovalLog(ctx, log)
}

// GetLoanApplicationByID 根据ID获取贷款申请
func (s *loanService) GetLoanApplicationByID(applicationID uint) (*model.LoanApplication, error) {
	ctx := context.Background()
	return s.loanRepo.GetApplicationByID(ctx, uint(applicationID))
}

// UpdateLoanApplication 更新贷款申请
func (s *loanService) UpdateLoanApplication(application *model.LoanApplication) error {
	ctx := context.Background()
	return s.loanRepo.UpdateApplication(ctx, application)
}

func generateApplicationNo() string {
	return fmt.Sprintf("LA%d%06d", time.Now().Unix(), time.Now().Nanosecond()%1000000)
}
