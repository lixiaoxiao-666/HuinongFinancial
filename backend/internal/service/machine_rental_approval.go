package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"huinong-backend/internal/model"
	"huinong-backend/internal/repository"
)

// RateOrder 评价订单
func (s *machineService) RateOrder(ctx context.Context, req *RateOrderRequest) error {
	// 验证订单存在且已完成
	order, err := s.machineRepo.GetRentalOrderByID(ctx, req.OrderID)
	if err != nil {
		return fmt.Errorf("订单不存在")
	}

	if order.Status != "completed" {
		return fmt.Errorf("只能对已完成的订单进行评价")
	}

	// 验证用户权限（只有租赁者或机主可以评价）
	if order.RenterID != req.UserID && order.OwnerID != req.UserID {
		return fmt.Errorf("无权评价该订单")
	}

	// 创建评价记录
	rating := &model.RentalRating{
		OrderID:    req.OrderID,
		RaterID:    req.UserID,
		RatingType: req.RatingType, // "renter" 或 "owner"
		Rating:     req.Rating,
		Comment:    req.Comment,
		CreatedAt:  time.Now(),
	}

	return s.machineRepo.CreateRentalRating(ctx, rating)
}

// GetRentalApplications 获取农机租赁申请列表(OA管理员)
func (s *machineService) GetRentalApplications(ctx context.Context, req *GetRentalApplicationsRequest) (*GetRentalApplicationsResponse, error) {
	// 转换为Repository请求
	repoReq := &repository.ListRentalApplicationsRequest{
		Status:      req.Status,
		MachineType: req.MachineType,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		RiskLevel:   req.RiskLevel,
		Page:        req.Page,
		Limit:       req.Limit,
		SortBy:      req.SortBy,
		SortOrder:   req.SortOrder,
	}

	result, err := s.machineRepo.ListRentalApplications(ctx, repoReq)
	if err != nil {
		return nil, fmt.Errorf("获取租赁申请列表失败: %v", err)
	}

	// 转换为响应格式
	applications := make([]*RentalApplicationItem, 0, len(result.Applications))
	for _, app := range result.Applications {
		// 获取用户信息以补充用户名称
		user, err := s.userRepo.GetByID(ctx, app.UserID)
		if err != nil {
			continue // 跳过无法获取用户信息的记录
		}

		// 获取农机信息以补充农机名称和类型
		machine, err := s.machineRepo.GetByID(ctx, app.MachineID)
		if err != nil {
			continue // 跳过无法获取农机信息的记录
		}

		item := &RentalApplicationItem{
			ID:            uint(app.ID),
			ApplicationNo: app.ApplicationNo,
			UserID:        app.UserID,
			UserName:      user.Username,
			MachineID:     app.MachineID,
			MachineName:   machine.MachineName,
			MachineType:   machine.MachineType,
			StartTime:     app.StartTime,
			EndTime:       app.EndTime,
			TotalAmount:   app.TotalAmount,
			Status:        app.Status,
			RiskLevel:     app.RiskLevel,
			CreatedAt:     app.CreatedAt,
			UpdatedAt:     app.UpdatedAt,
		}
		applications = append(applications, item)
	}

	return &GetRentalApplicationsResponse{
		Applications: applications,
		Total:        result.Total,
		Page:         result.Page,
		Limit:        result.Limit,
		Statistics: &RentalApplicationStatistics{
			TotalApplications:    result.Statistics.TotalApplications,
			PendingApplications:  result.Statistics.PendingApplications,
			ApprovedApplications: result.Statistics.ApprovedApplications,
			RejectedApplications: result.Statistics.RejectedApplications,
		},
	}, nil
}

// GetRentalApplicationDetail 获取农机租赁申请详情(OA管理员)
func (s *machineService) GetRentalApplicationDetail(ctx context.Context, req *GetRentalApplicationDetailRequest) (*GetRentalApplicationDetailResponse, error) {
	// 获取申请详情
	application, err := s.machineRepo.GetRentalApplicationByID(ctx, uint64(req.ID))
	if err != nil {
		return nil, fmt.Errorf("申请不存在")
	}

	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, application.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 获取农机信息
	machine, err := s.machineRepo.GetByID(ctx, application.MachineID)
	if err != nil {
		return nil, fmt.Errorf("获取农机信息失败: %v", err)
	}

	// 获取审核记录
	reviewLogs, err := s.machineRepo.GetRentalReviewLogs(ctx, uint64(req.ID))
	if err != nil {
		reviewLogs = []*model.RentalReviewLog{} // 如果没有审核记录，返回空列表
	}

	// 解析风险因素
	var riskFactors []string
	if application.RiskFactors != "" {
		// 假设RiskFactors是JSON字符串，需要解析
		// 这里简化处理，直接转为字符串数组
		riskFactors = []string{application.RiskFactors}
	}

	return &GetRentalApplicationDetailResponse{
		Application: &RentalApplicationDetail{
			ID:             uint(application.ID),
			ApplicationNo:  application.ApplicationNo,
			UserID:         uint(application.UserID),
			UserName:       user.Username,
			UserRealName:   user.RealName,
			UserPhone:      user.Phone,
			MachineID:      uint(application.MachineID),
			MachineName:    machine.MachineName,
			MachineType:    machine.MachineType,
			StartTime:      application.StartTime,
			EndTime:        application.EndTime,
			RentalLocation: application.RentalLocation,
			ContactPerson:  application.ContactPerson,
			ContactPhone:   application.ContactPhone,
			BillingMethod:  application.BillingMethod,
			Quantity:       application.Quantity,
			TotalAmount:    application.TotalAmount,
			DepositAmount:  application.DepositAmount,
			Status:         application.Status,
			RiskLevel:      application.RiskLevel,
			RiskFactors:    riskFactors,
			AIAssessment:   application.AIAssessment,
			Remarks:        application.Remarks,
			CreatedAt:      application.CreatedAt,
			UpdatedAt:      application.UpdatedAt,
		},
		User:       user,
		Machine:    machine,
		ReviewLogs: reviewLogs,
	}, nil
}

// ApproveRentalApplication 批准农机租赁申请(OA管理员)
func (s *machineService) ApproveRentalApplication(ctx context.Context, req *ApproveRentalRequest) error {
	// 获取申请信息
	application, err := s.machineRepo.GetRentalApplicationByID(ctx, uint64(req.ID))
	if err != nil {
		return fmt.Errorf("申请不存在")
	}

	if application.Status != "pending" {
		return fmt.Errorf("申请状态不允许审批")
	}

	// 更新申请状态
	application.Status = "approved"
	application.ReviewerID = &req.ReviewerID
	application.ReviewedAt = &time.Time{}
	*application.ReviewedAt = time.Now()
	application.ReviewNote = req.ApprovalNote
	application.UpdatedAt = time.Now()

	err = s.machineRepo.UpdateRentalApplication(ctx, application)
	if err != nil {
		return fmt.Errorf("更新申请状态失败: %v", err)
	}

	// 创建审核记录
	reviewLog := &model.RentalReviewLog{
		ApplicationID: uint64(req.ID),
		ReviewerID:    req.ReviewerID,
		Action:        "approve",
		Note:          req.ApprovalNote,
		CreatedAt:     time.Now(),
	}

	err = s.machineRepo.CreateRentalReviewLog(ctx, reviewLog)
	if err != nil {
		return fmt.Errorf("创建审核记录失败: %v", err)
	}

	// TODO: 发送通知给申请人
	// s.notificationService.SendApprovalNotification(application.UserID, application.ID, "approved")

	return nil
}

// RejectRentalApplication 拒绝农机租赁申请(OA管理员)
func (s *machineService) RejectRentalApplication(ctx context.Context, req *RejectRentalRequest) error {
	// 获取申请信息
	application, err := s.machineRepo.GetRentalApplicationByID(ctx, uint64(req.ID))
	if err != nil {
		return fmt.Errorf("申请不存在")
	}

	if application.Status != "pending" {
		return fmt.Errorf("申请状态不允许审批")
	}

	// 更新申请状态
	application.Status = "rejected"
	application.ReviewerID = &req.ReviewerID
	application.ReviewedAt = &time.Time{}
	*application.ReviewedAt = time.Now()
	application.ReviewNote = req.RejectionReason
	application.UpdatedAt = time.Now()

	err = s.machineRepo.UpdateRentalApplication(ctx, application)
	if err != nil {
		return fmt.Errorf("更新申请状态失败: %v", err)
	}

	// 创建审核记录
	reviewLog := &model.RentalReviewLog{
		ApplicationID: uint64(req.ID),
		ReviewerID:    req.ReviewerID,
		Action:        "reject",
		Note:          req.RejectionReason,
		CreatedAt:     time.Now(),
	}

	err = s.machineRepo.CreateRentalReviewLog(ctx, reviewLog)
	if err != nil {
		return fmt.Errorf("创建审核记录失败: %v", err)
	}

	// TODO: 发送通知给申请人
	// s.notificationService.SendApprovalNotification(application.UserID, application.ID, "rejected")

	return nil
}

// GetRentalStatistics 获取农机租赁统计(OA管理员)
func (s *machineService) GetRentalStatistics(ctx context.Context, req *GetRentalStatisticsRequest) (*GetRentalStatisticsResponse, error) {
	// 由于MachineRepository接口中没有GetRentalStatistics方法，我们先实现一个简化版本
	// TODO: 在repository接口和实现中添加GetRentalStatistics方法

	// 这里先返回一个基础的统计响应
	return &GetRentalStatisticsResponse{
		Overview: &RentalStatisticsOverview{
			TotalApplications:    0,
			PendingApplications:  0,
			ApprovedApplications: 0,
			RejectedApplications: 0,
			TotalRentalAmount:    0,
			ApprovedRentalAmount: 0,
			ApprovalRate:         0.0,
		},
		Trends:           nil,
		MachineTypeStats: nil,
		RegionStats:      nil,
	}, nil
}

// BatchApproveRentals 批量审批农机租赁申请(OA管理员)
func (s *machineService) BatchApproveRentals(ctx context.Context, req *BatchApproveRentalsRequest) (*BatchApproveRentalsResponse, error) {
	results := make([]*BatchApprovalResult, 0, len(req.ApplicationIDs))
	successCount := 0
	failureCount := 0

	for _, appID := range req.ApplicationIDs {
		result := &BatchApprovalResult{
			ApplicationID: appID,
			Success:       false,
		}

		// 根据决策类型处理
		switch req.Decision {
		case "approve":
			approveReq := &ApproveRentalRequest{
				ID:           uint(appID),
				ReviewerID:   req.ReviewerID,
				ApprovalNote: req.Note,
			}
			err := s.ApproveRentalApplication(ctx, approveReq)
			if err != nil {
				result.Error = err.Error()
				failureCount++
			} else {
				result.Success = true
				successCount++
			}

		case "reject":
			rejectReq := &RejectRentalRequest{
				ID:              uint(appID),
				ReviewerID:      req.ReviewerID,
				RejectionReason: req.Note,
			}
			err := s.RejectRentalApplication(ctx, rejectReq)
			if err != nil {
				result.Error = err.Error()
				failureCount++
			} else {
				result.Success = true
				successCount++
			}

		default:
			result.Error = "无效的审批决策"
			failureCount++
		}

		results = append(results, result)
	}

	return &BatchApproveRentalsResponse{
		Results:      results,
		SuccessCount: successCount,
		FailureCount: failureCount,
		TotalCount:   len(req.ApplicationIDs),
	}, nil
}

// 辅助函数：生成申请编号
func generateRentalApplicationNo() string {
	return "RA" + strconv.FormatInt(time.Now().UnixNano(), 10)[10:]
}

// 辅助函数：计算风险等级
func calculateRiskLevel(user *model.User, machine *model.Machine, amount int64) string {
	// 简单的风险评估逻辑
	if amount > 50000 {
		return "high"
	}
	if amount > 20000 {
		return "medium"
	}
	return "low"
}

// 辅助函数：评估风险因素
func assessRiskFactors(user *model.User, machine *model.Machine, amount int64) []string {
	factors := []string{}

	// 根据用户信息评估风险
	// 由于User模型中没有CreditScore字段，我们跳过这个检查
	// if user.CreditScore < 60 {
	//     factors = append(factors, "用户信用分较低")
	// }

	// 根据金额评估风险
	if amount > 30000 {
		factors = append(factors, "租赁金额较高")
	}

	// 根据农机状态评估风险
	if machine.Status != "available" {
		factors = append(factors, "农机状态异常")
	}

	return factors
}
