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
	loanRepo    repository.LoanRepository
	userRepo    repository.UserRepository
	difyService DifyService
	taskService TaskService
}

// NewLoanService 创建贷款服务实例
func NewLoanService(
	loanRepo repository.LoanRepository,
	userRepo repository.UserRepository,
	difyService DifyService,
	taskService TaskService,
) LoanService {
	return &loanService{
		loanRepo:    loanRepo,
		userRepo:    userRepo,
		difyService: difyService,
		taskService: taskService,
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

	// 异步触发AI审批流程
	go func() {
		// 使用新的context避免原context被取消
		assessmentCtx := context.Background()
		if triggerErr := s.TriggerAIAssessment(assessmentCtx, application.ID, "loan_approval"); triggerErr != nil {
			// 记录错误日志，但不影响申请创建的响应
			fmt.Printf("触发AI评估失败 - 申请ID: %d, 错误: %v\n", application.ID, triggerErr)
		}
	}()

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
	previousStatus := application.Status
	application.Status = "approved"
	if application.FinalApprovedAt == nil {
		now := time.Now()
		application.FinalApprovedAt = &now
	}

	err = s.loanRepo.UpdateApplication(ctx, application)
	if err != nil {
		return fmt.Errorf("更新申请状态失败: %v", err)
	}

	// 记录审批日志
	approvalLog := &model.ApprovalLog{
		ApplicationID:  application.ID,
		ApproverID:     0, // TODO: 从上下文中获取当前操作员ID
		Action:         "approve",
		Step:           "manual_review", // 假设这是人工审核步骤
		Result:         "approved",
		Status:         "approved",
		Comment:        req.ApprovalNote,
		Note:           "人工批准贷款申请",
		PreviousStatus: previousStatus,
		NewStatus:      application.Status,
		ActionTime:     time.Now(),
	}
	if err := s.loanRepo.CreateApprovalLog(ctx, approvalLog); err != nil {
		// 记录日志创建失败，但不影响主流程
		fmt.Printf("创建贷款审批日志失败 - 申请ID %d: %v\n", application.ID, err)
	}

	// 如果存在关联的待处理任务，则将其标记为完成
	// 这里假设我们通过某种方式（例如，在AI评估后创建任务时存储关联）找到任务ID
	// 简化处理：假设从 application 模型的某个字段获取 taskID
	/*
		if application.RelatedTaskID != nil && *application.RelatedTaskID > 0 {
			err = s.taskService.CompleteTask(ctx, *application.RelatedTaskID)
			if err != nil {
				fmt.Printf("完成贷款审批任务失败 - 任务ID %d: %v\n", *application.RelatedTaskID, err)
				// 即使任务完成失败，也不应阻塞主流程，但需要记录
			}
		}
	*/

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
	previousStatus := application.Status
	application.Status = "rejected"
	application.RejectionReason = req.RejectionNote

	err = s.loanRepo.UpdateApplication(ctx, application)
	if err != nil {
		return fmt.Errorf("更新申请状态失败: %v", err)
	}

	// 记录审批日志
	approvalLog := &model.ApprovalLog{
		ApplicationID:  application.ID,
		ApproverID:     0, // TODO: 从上下文中获取当前操作员ID
		Action:         "reject",
		Step:           "manual_review", // 假设这是人工审核步骤
		Result:         "rejected",
		Status:         "rejected",
		Comment:        req.RejectionNote,
		Note:           "人工拒绝贷款申请",
		PreviousStatus: previousStatus,
		NewStatus:      application.Status,
		ActionTime:     time.Now(),
	}
	if err := s.loanRepo.CreateApprovalLog(ctx, approvalLog); err != nil {
		fmt.Printf("创建贷款审批日志失败 - 申请ID %d: %v\n", application.ID, err)
	}

	// 如果存在关联的待处理任务，则将其标记为完成（或取消）
	/*
		if application.RelatedTaskID != nil && *application.RelatedTaskID > 0 {
			// 拒绝也意味着任务结束
			err = s.taskService.CompleteTask(ctx, *application.RelatedTaskID) // 或者一个特定的取消方法
			if err != nil {
				fmt.Printf("完成/取消贷款审批任务失败 - 任务ID %d: %v\n", *application.RelatedTaskID, err)
			}
		}
	*/

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

	previousStatus := application.Status
	application.Status = "returned"
	// 可以考虑将退回原因保存到某个字段，例如 application.ReturnReason = req.ReturnReason

	if err := s.loanRepo.UpdateApplication(ctx, application); err != nil {
		return fmt.Errorf("更新申请状态为退回失败: %v", err)
	}

	// 记录审批日志
	approverID := uint64(0) // TODO: 从上下文中获取当前操作员ID
	approvalLog := &model.ApprovalLog{
		ApplicationID:  application.ID,
		ApproverID:     approverID,
		Action:         "return",
		Step:           "manual_review",
		Result:         "returned",
		Status:         "returned",
		Comment:        req.ReturnNote,
		Note:           fmt.Sprintf("贷款申请被退回，原因: %s", req.ReturnReason),
		PreviousStatus: previousStatus,
		NewStatus:      application.Status,
		ActionTime:     time.Now(),
	}
	if err := s.loanRepo.CreateApprovalLog(ctx, approvalLog); err != nil {
		fmt.Printf("创建贷款退回审批日志失败 - 申请ID %d: %v\n", application.ID, err)
	}

	// 创建任务：通知用户修改申请或进行其他处理
	// 假设退回的任务需要指派给申请提交用户或特定处理人
	// var assignedTo *uint64 // 可以根据业务逻辑设置指派人
	// if application.UserID > 0 { assignedTo = &application.UserID }

	taskReq := &CreateTaskRequest{
		Title:        fmt.Sprintf("贷款申请被退回: %s", application.ApplicationNo),
		Description:  fmt.Sprintf("您的贷款申请 %s 已被退回，原因: %s。请根据审批意见修改后重新提交。退回备注: %s", application.ApplicationNo, req.ReturnReason, req.ReturnNote),
		Type:         "loan_application_returned",
		Priority:     "high",
		BusinessID:   application.ID,
		BusinessType: "loan_application",
		// AssignedTo:   assignedTo, // 指派给申请人或特定处理角色
		Data: fmt.Sprintf(`{"application_id": %d, "application_no": "%s", "return_reason": "%s", "return_note": "%s"}`, application.ID, application.ApplicationNo, req.ReturnReason, req.ReturnNote),
	}

	taskResp, err := s.taskService.CreateTask(ctx, taskReq)
	if err != nil {
		fmt.Printf("为退回的贷款申请创建任务失败 - 申请ID %d: %v\n", application.ID, err)
		// 即使创建任务失败，也不应阻塞主流程，但需要记录
	} else {
		fmt.Printf("为退回的贷款申请创建任务成功 - 申请ID %d, 任务ID %d\n", application.ID, taskResp.ID)
		// 可以考虑将taskID关联到application模型
		// application.RelatedTaskID = &taskResp.ID
		// s.loanRepo.UpdateApplication(ctx, application) // 再次更新
	}

	return nil
}

// ==================== AI集成功能 ====================

// TriggerAIAssessment 触发AI评估
func (s *loanService) TriggerAIAssessment(ctx context.Context, applicationID uint64, workflowType string) error {
	// 获取申请详情
	application, err := s.loanRepo.GetApplicationByID(ctx, uint(applicationID))
	if err != nil {
		return fmt.Errorf("获取申请详情失败: %v", err)
	}

	// 更新申请状态为AI评估中
	application.Status = "ai_processing"
	if updateErr := s.loanRepo.UpdateApplication(ctx, application); updateErr != nil {
		fmt.Printf("更新申请状态失败: %v\n", updateErr)
	}

	// 创建审批日志记录AI评估开始
	startLog := &model.ApprovalLog{
		ApplicationID:  application.ID,
		ApproverID:     0, // AI系统
		Action:         "submit",
		Step:           "ai_assessment",
		ApprovalLevel:  0,
		Result:         "pending",
		Status:         "pending",
		Comment:        "开始AI智能评估",
		Note:           "系统自动触发AI工作流进行贷款申请评估",
		PreviousStatus: "pending",
		NewStatus:      "ai_processing",
		ActionTime:     time.Now(),
	}
	if logErr := s.loanRepo.CreateApprovalLog(ctx, startLog); logErr != nil {
		fmt.Printf("创建审批日志失败: %v\n", logErr)
	}

	// 调用Dify工作流
	response, err := s.difyService.CallLoanApprovalWorkflow(uint(applicationID), uint(application.UserID))
	if err != nil {
		// 更新申请状态为AI评估失败
		application.Status = "ai_failed"
		application.RejectionReason = fmt.Sprintf("AI评估失败: %v", err)
		if updateErr := s.loanRepo.UpdateApplication(ctx, application); updateErr != nil {
			fmt.Printf("更新失败状态失败: %v\n", updateErr)
		}

		// 记录失败日志
		failLog := &model.ApprovalLog{
			ApplicationID:  application.ID,
			ApproverID:     0,
			Action:         "reject",
			Step:           "ai_assessment",
			ApprovalLevel:  0,
			Result:         "rejected",
			Status:         "rejected",
			Comment:        fmt.Sprintf("AI评估失败: %v", err),
			Note:           "AI工作流调用失败",
			PreviousStatus: "ai_processing",
			NewStatus:      "ai_failed",
			ActionTime:     time.Now(),
		}
		if logErr := s.loanRepo.CreateApprovalLog(ctx, failLog); logErr != nil {
			fmt.Printf("创建失败日志失败: %v\n", logErr)
		}

		return fmt.Errorf("调用Dify工作流失败: %v", err)
	}

	// 处理AI评估结果
	if err := s.processAIAssessmentResult(ctx, application, response); err != nil {
		return fmt.Errorf("处理AI评估结果失败: %v", err)
	}

	return nil
}

// processAIAssessmentResult 处理AI评估结果
func (s *loanService) processAIAssessmentResult(ctx context.Context, application *model.LoanApplication, response *DifyWorkflowResponse) error {
	// 更新申请的Dify对话信息
	application.DifyConversationID = response.ConversationID

	// 根据AI响应状态处理结果
	if response.Status == "succeeded" && response.Data != nil {
		// 解析AI评估结果
		result, ok := response.Data["result"]
		if !ok {
			return fmt.Errorf("AI响应中缺少result字段")
		}

		resultMap, ok := result.(map[string]interface{})
		if !ok {
			return fmt.Errorf("AI结果格式不正确")
		}

		// 获取AI建议
		recommendation := ""
		if rec, exists := resultMap["recommendation"]; exists {
			recommendation = fmt.Sprintf("%v", rec)
		}

		// 获取AI决策
		decision := ""
		if dec, exists := resultMap["decision"]; exists {
			decision = fmt.Sprintf("%v", dec)
		}

		// 获取信用评分
		if score, exists := resultMap["credit_score"]; exists {
			if scoreFloat, ok := score.(float64); ok {
				application.CreditScore = int(scoreFloat)
			}
		}

		// 获取风险等级
		if risk, exists := resultMap["risk_level"]; exists {
			application.RiskLevel = fmt.Sprintf("%v", risk)
		}

		// 更新AI建议
		application.AIRecommendation = recommendation

		// 根据AI决策更新申请状态
		var newStatus string
		var actionResult string
		var actionComment string

		switch decision {
		case "approve", "approved":
			newStatus = "ai_approved"
			actionResult = "approved"
			actionComment = "AI智能评估通过，建议批准申请"
			application.AutoApprovalPassed = true

			// 如果有推荐金额，设置批准金额
			if approvedAmount, exists := resultMap["approved_amount"]; exists {
				if amountFloat, ok := approvedAmount.(float64); ok {
					amount := int64(amountFloat)
					application.ApprovedAmount = &amount
				}
			}

		case "reject", "rejected":
			newStatus = "ai_rejected"
			actionResult = "rejected"
			actionComment = "AI智能评估不通过，建议拒绝申请"

			// 设置拒绝原因
			if reason, exists := resultMap["rejection_reason"]; exists {
				application.RejectionReason = fmt.Sprintf("%v", reason)
			} else {
				application.RejectionReason = "AI评估认为不符合贷款条件"
			}

		case "manual_review", "manual":
			newStatus = "manual_review"
			actionResult = "returned"
			actionComment = "AI智能评估建议人工审核"
			application.ApprovalLevel = 1

			// 创建人工审核任务
			taskTitle := fmt.Sprintf("贷款申请需人工审核: %s", application.ApplicationNo)
			taskDescription := fmt.Sprintf("贷款申请 %s (ID: %d) AI评估完成，建议进行人工审核。AI建议: %s", application.ApplicationNo, application.ID, recommendation)
			taskData := fmt.Sprintf(`{"application_id": %d, "application_no": "%s", "ai_recommendation": "%s", "risk_level": "%s"}`, application.ID, application.ApplicationNo, recommendation, application.RiskLevel)

			taskReq := &CreateTaskRequest{
				Title:        taskTitle,
				Description:  taskDescription,
				Type:         "loan_manual_review",
				Priority:     "medium",
				BusinessID:   application.ID,
				BusinessType: "loan_application",
				// AssignedTo: nil, // 可以后续通过OA系统指派给特定审核员或审核组
				Data: taskData,
			}
			taskResp, taskErr := s.taskService.CreateTask(ctx, taskReq)
			if taskErr != nil {
				fmt.Printf("为贷款申请创建人工审核任务失败 - 申请ID %d: %v\n", application.ID, taskErr)
				// 创建任务失败也记录，但不中断流程
				actionComment = actionComment + " (创建审核任务失败)"
			} else {
				fmt.Printf("为贷款申请创建人工审核任务成功 - 申请ID %d, 任务ID %d\n", application.ID, taskResp.ID)
				// TODO: 考虑将taskID关联到application，用于后续跟踪
				// application.RelatedTaskID = &taskResp.ID
			}

		default:
			newStatus = "manual_review"
			actionResult = "returned"
			actionComment = "AI评估完成，建议人工审核 (未知决策)"
			application.ApprovalLevel = 1

			// 也为未知决策创建人工审核任务
			taskTitle := fmt.Sprintf("贷款申请需人工审核 (未知AI决策): %s", application.ApplicationNo)
			taskDescription := fmt.Sprintf("贷款申请 %s (ID: %d) AI评估决策未知，建议进行人工审核。AI建议: %s", application.ApplicationNo, application.ID, recommendation)
			taskData := fmt.Sprintf(`{"application_id": %d, "application_no": "%s", "ai_recommendation": "%s", "risk_level": "%s", "raw_decision": "%s"}`, application.ID, application.ApplicationNo, recommendation, application.RiskLevel, decision)

			taskReq := &CreateTaskRequest{
				Title:        taskTitle,
				Description:  taskDescription,
				Type:         "loan_manual_review_unknown",
				Priority:     "medium",
				BusinessID:   application.ID,
				BusinessType: "loan_application",
				Data:         taskData,
			}
			taskResp, taskErr := s.taskService.CreateTask(ctx, taskReq)
			if taskErr != nil {
				fmt.Printf("为贷款申请创建人工审核任务失败 (未知AI决策) - 申请ID %d: %v\n", application.ID, taskErr)
			} else {
				fmt.Printf("为贷款申请创建人工审核任务成功 (未知AI决策) - 申请ID %d, 任务ID %d\n", application.ID, taskResp.ID)
			}
		}

		application.Status = newStatus

		// 更新申请记录
		if err := s.loanRepo.UpdateApplication(ctx, application); err != nil {
			return fmt.Errorf("更新申请记录失败: %v", err)
		}

		// 创建AI评估完成日志
		resultLog := &model.ApprovalLog{
			ApplicationID:  application.ID,
			ApproverID:     0, // AI系统
			Action:         actionResult,
			Step:           "ai_assessment",
			ApprovalLevel:  0,
			Result:         actionResult,
			Status:         actionResult,
			Comment:        actionComment,
			Note:           fmt.Sprintf("AI评估建议: %s", recommendation),
			PreviousStatus: "ai_processing",
			NewStatus:      newStatus,
			ActionTime:     time.Now(),
		}

		if err := s.loanRepo.CreateApprovalLog(ctx, resultLog); err != nil {
			fmt.Printf("创建AI评估结果日志失败: %v\n", err)
		}

	} else {
		// AI评估失败
		application.Status = "ai_failed"
		application.RejectionReason = "AI评估过程中出现错误"

		if err := s.loanRepo.UpdateApplication(ctx, application); err != nil {
			return fmt.Errorf("更新失败状态失败: %v", err)
		}

		// 记录失败日志
		failLog := &model.ApprovalLog{
			ApplicationID:  application.ID,
			ApproverID:     0,
			Action:         "reject",
			Step:           "ai_assessment",
			ApprovalLevel:  0,
			Result:         "rejected",
			Status:         "rejected",
			Comment:        "AI评估失败",
			Note:           fmt.Sprintf("AI工作流状态: %s", response.Status),
			PreviousStatus: "ai_processing",
			NewStatus:      "ai_failed",
			ActionTime:     time.Now(),
		}

		if err := s.loanRepo.CreateApprovalLog(ctx, failLog); err != nil {
			fmt.Printf("创建失败日志失败: %v\n", err)
		}

		return fmt.Errorf("AI评估未成功完成，状态: %s", response.Status)
	}

	return nil
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

// generateApplicationNo 生成申请编号
func generateApplicationNo() string {
	return fmt.Sprintf("LA%d", time.Now().UnixNano()/1000000)
}

// ==================== 增强审批功能 ====================

// BatchApproveLoanApplications 批量审批贷款申请
func (s *loanService) BatchApproveLoanApplications(ctx context.Context, req *BatchApproveLoanRequest) (*BatchApproveLoanResponse, error) {
	results := make([]BatchOperationResult, 0, len(req.ApplicationIDs))
	successCount := 0
	failureCount := 0
	batchID := fmt.Sprintf("BATCH_%d", time.Now().UnixNano())

	for _, appID := range req.ApplicationIDs {
		result := BatchOperationResult{
			ID:      appID,
			Success: false,
		}

		// 获取申请
		_, err := s.loanRepo.GetApplicationByID(ctx, uint(appID))
		if err != nil {
			result.Error = fmt.Sprintf("获取申请失败: %v", err)
			failureCount++
		} else {
			// 执行审批
			if req.Decision == "approve" {
				err = s.ApproveApplication(ctx, &ApproveApplicationRequest{
					ID:           uint(appID),
					ApprovalNote: req.ReviewComments,
				})
			} else {
				err = s.RejectApplication(ctx, &RejectApplicationRequest{
					ID:            uint(appID),
					RejectionNote: req.ReviewComments,
				})
			}

			if err != nil {
				result.Error = fmt.Sprintf("审批失败: %v", err)
				failureCount++
			} else {
				result.Success = true
				result.Message = "审批成功"
				successCount++
			}
		}

		results = append(results, result)
	}

	return &BatchApproveLoanResponse{
		TotalCount:   len(req.ApplicationIDs),
		SuccessCount: successCount,
		FailureCount: failureCount,
		Results:      results,
		BatchID:      batchID,
		ProcessedBy:  req.ReviewerID,
		ProcessedAt:  time.Now(),
	}, nil
}

// RetryAIAssessment 重试AI评估
func (s *loanService) RetryAIAssessment(ctx context.Context, req *RetryAIAssessmentRequest) error {
	return s.TriggerAIAssessment(ctx, uint64(req.ID), "loan_approval")
}

// EnableAutoApproval 启用自动审批
func (s *loanService) EnableAutoApproval(ctx context.Context, req *EnableAutoApprovalRequest) error {
	// TODO: 实现自动审批配置的创建
	return fmt.Errorf("暂未实现自动审批功能")
}

// DisableAutoApproval 禁用自动审批
func (s *loanService) DisableAutoApproval(ctx context.Context, req *DisableAutoApprovalRequest) error {
	// TODO: 实现自动审批配置的禁用
	return fmt.Errorf("暂未实现自动审批功能")
}

// GetAutoApprovalConfig 获取自动审批配置
func (s *loanService) GetAutoApprovalConfig(ctx context.Context) (*GetAutoApprovalConfigResponse, error) {
	// TODO: 实现获取自动审批配置
	return &GetAutoApprovalConfigResponse{
		LoanConfig:   nil,
		RentalConfig: nil,
		GlobalSettings: GlobalAutoApprovalSettings{
			MaxDailyAutoApprovals: 0,
			MaxAmountPerDay:       0,
			RequireSecondApproval: true,
			MonitoringEnabled:     true,
		},
	}, nil
}

// GetApplicationsByRiskLevel 按风险等级获取申请
func (s *loanService) GetApplicationsByRiskLevel(ctx context.Context, req *GetApplicationsByRiskLevelRequest) (*GetApplicationsByRiskLevelResponse, error) {
	// TODO: 实现按风险等级查询申请
	return &GetApplicationsByRiskLevelResponse{
		Applications: []LoanApplicationSummary{},
		Pagination: PaginationInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Total: 0,
			Pages: 0,
		},
		RiskAnalysis: RiskLevelAnalysis{
			RiskLevel:      req.RiskLevel,
			TotalCount:     0,
			ApprovalRate:   0,
			AverageAmount:  0,
			AverageScore:   0,
			TrendDirection: "stable",
		},
	}, nil
}

// GetAIAssessmentHistory 获取AI评估历史
func (s *loanService) GetAIAssessmentHistory(ctx context.Context, req *GetAIAssessmentHistoryRequest) (*GetAIAssessmentHistoryResponse, error) {
	// TODO: 实现获取AI评估历史
	return &GetAIAssessmentHistoryResponse{
		ApplicationID: req.ApplicationID,
		Assessments:   []AIAssessmentRecord{},
		Summary: AIAssessmentSummary{
			TotalAssessments:  0,
			LatestVersion:     1,
			CurrentRiskLevel:  "medium",
			RiskTrend:         "stable",
			AverageConfidence: 0.8,
			ModelConsistency:  0.9,
		},
	}, nil
}

// CreateApplicationTask 创建申请任务
func (s *loanService) CreateApplicationTask(ctx context.Context, req *CreateApplicationTaskRequest) (*CreateApplicationTaskResponse, error) {
	// TODO: 实现创建申请任务
	return &CreateApplicationTaskResponse{
		TaskID:        uint64(time.Now().UnixNano()),
		ApplicationID: req.ApplicationID,
		Status:        "created",
		CreatedAt:     time.Now(),
	}, nil
}

// GetApplicationTasks 获取申请任务
func (s *loanService) GetApplicationTasks(ctx context.Context, req *GetApplicationTasksRequest) (*GetApplicationTasksResponse, error) {
	// TODO: 实现获取申请任务
	return &GetApplicationTasksResponse{
		ApplicationID: req.ApplicationID,
		Tasks:         []TaskInfo{},
		Summary: TaskSummary{
			TotalTasks:      0,
			PendingTasks:    0,
			InProgressTasks: 0,
			CompletedTasks:  0,
			OverdueTasks:    0,
		},
	}, nil
}

// GetAdvancedStatistics 获取高级统计
func (s *loanService) GetAdvancedStatistics(ctx context.Context, req *GetAdvancedStatisticsRequest) (*GetAdvancedStatisticsResponse, error) {
	// TODO: 实现获取高级统计
	return &GetAdvancedStatisticsResponse{
		Overview: StatisticsOverview{
			TotalApplications:    0,
			ApprovedApplications: 0,
			RejectedApplications: 0,
			PendingApplications:  0,
			TotalAmount:          0,
			ApprovedAmount:       0,
			ApprovalRate:         0,
			AverageProcessTime:   0,
		},
		RiskAnalysis: RiskAnalysis{
			RiskDistribution: map[string]int{},
			HighRiskFactors:  []RiskFactorStats{},
			RiskTrends:       []RiskTrendPoint{},
		},
		Performance: PerformanceMetrics{
			AverageProcessingTime: 0,
			AutoApprovalRate:      0,
			ManualReviewRate:      0,
			ReviewerEfficiency:    0,
			SystemUptime:          0.99,
			ErrorRate:             0.01,
		},
	}, nil
}
