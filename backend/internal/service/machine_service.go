package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"huinong-backend/internal/model"
	"huinong-backend/internal/repository"
)

// machineService 农机服务实现
type machineService struct {
	machineRepo repository.MachineRepository
	userRepo    repository.UserRepository
}

// NewMachineService 创建农机服务实例
func NewMachineService(machineRepo repository.MachineRepository, userRepo repository.UserRepository) MachineService {
	return &machineService{
		machineRepo: machineRepo,
		userRepo:    userRepo,
	}
}

// ==================== 农机管理 ====================

// RegisterMachine 注册农机
func (s *machineService) RegisterMachine(ctx context.Context, req *RegisterMachineRequest) (*RegisterMachineResponse, error) {
	// TODO: 从上下文中获取用户ID
	// userID := getUserIDFromContext(ctx)

	// 生成农机编码
	machineCode := generateMachineCode()

	machine := &model.Machine{
		MachineCode:    machineCode,
		MachineName:    req.MachineName,
		MachineType:    req.MachineType,
		Brand:          req.Brand,
		Model:          req.Model,
		Description:    req.Description,
		Province:       req.Province,
		City:           req.City,
		County:         req.County,
		Address:        req.Address,
		Longitude:      req.Longitude,
		Latitude:       req.Latitude,
		HourlyRate:     req.HourlyRate,
		DailyRate:      req.DailyRate,
		PerAcreRate:    req.PerAcreRate,
		DepositAmount:  req.DepositAmount,
		MinRentalHours: req.MinRentalHours,
		MaxAdvanceDays: req.MaxAdvanceDays,
		Status:         "pending", // 待审核
		// OwnerID:        userID,
		// Images:         strings.Join(req.Images, ","),
		// Specifications: marshalToJSON(req.Specifications),
		// AvailableSchedule: marshalToJSON(req.AvailableSchedule),
	}

	err := s.machineRepo.Create(ctx, machine)
	if err != nil {
		return nil, fmt.Errorf("注册农机失败: %v", err)
	}

	return &RegisterMachineResponse{
		ID:          machine.ID,
		MachineCode: machine.MachineCode,
		Status:      machine.Status,
	}, nil
}

// GetMachine 获取农机详情
func (s *machineService) GetMachine(ctx context.Context, machineID uint64) (*MachineDetailResponse, error) {
	machine, err := s.machineRepo.GetByID(ctx, machineID)
	if err != nil {
		return nil, fmt.Errorf("获取农机信息失败: %v", err)
	}

	// 获取机主信息
	// owner, _ := s.userRepo.GetByID(ctx, machine.OwnerID)

	// TODO: 获取评价信息
	ratings := []MachineRating{}

	return &MachineDetailResponse{
		Machine: machine,
		// Owner:   owner,
		Ratings: ratings,
	}, nil
}

// UpdateMachine 更新农机信息
func (s *machineService) UpdateMachine(ctx context.Context, machineID uint64, req *UpdateMachineRequest) error {
	machine, err := s.machineRepo.GetByID(ctx, machineID)
	if err != nil {
		return fmt.Errorf("获取农机信息失败: %v", err)
	}

	// 更新字段
	if req.MachineName != "" {
		machine.MachineName = req.MachineName
	}
	if req.Description != "" {
		machine.Description = req.Description
	}
	if req.HourlyRate > 0 {
		machine.HourlyRate = req.HourlyRate
	}
	if req.DailyRate > 0 {
		machine.DailyRate = req.DailyRate
	}
	if req.PerAcreRate > 0 {
		machine.PerAcreRate = req.PerAcreRate
	}
	if req.DepositAmount > 0 {
		machine.DepositAmount = req.DepositAmount
	}
	if req.Status != "" {
		machine.Status = req.Status
	}

	// TODO: 更新图片和其他JSON字段
	// if len(req.Images) > 0 {
	//     machine.Images = strings.Join(req.Images, ",")
	// }

	return s.machineRepo.Update(ctx, machine)
}

// DeleteMachine 删除农机
func (s *machineService) DeleteMachine(ctx context.Context, machineID uint64) error {
	return s.machineRepo.Delete(ctx, machineID)
}

// GetUserMachines 获取用户农机列表
func (s *machineService) GetUserMachines(ctx context.Context, userID uint64, req *GetUserMachinesRequest) (*GetUserMachinesResponse, error) {
	// TODO: 实现获取用户农机列表逻辑
	return &GetUserMachinesResponse{
		Machines: []*model.Machine{},
		Total:    0,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}

// ==================== 农机搜索 ====================

// SearchMachines 搜索农机
func (s *machineService) SearchMachines(ctx context.Context, req *SearchMachinesRequest) (*SearchMachinesResponse, error) {
	// 转换为Repository请求
	repoReq := &repository.ListMachinesRequest{
		Page:        req.Page,
		Limit:       req.Limit,
		MachineType: req.MachineType,
		Province:    req.Province,
		City:        req.City,
		County:      req.County,
		Longitude:   req.Longitude,
		Latitude:    req.Latitude,
		Radius:      float64(req.Radius),
		Status:      "active", // 只搜索可用的农机
	}

	result, err := s.machineRepo.List(ctx, repoReq)
	if err != nil {
		return nil, fmt.Errorf("搜索农机失败: %v", err)
	}

	return &SearchMachinesResponse{
		Machines: result.Machines,
		Total:    result.Total,
		Page:     result.Page,
		Limit:    result.Limit,
	}, nil
}

// GetAvailableMachines 获取可用农机
func (s *machineService) GetAvailableMachines(ctx context.Context, req *GetAvailableMachinesRequest) (*GetAvailableMachinesResponse, error) {
	// TODO: 根据时间段查询可用农机
	repoReq := &repository.ListMachinesRequest{
		Page:        req.Page,
		Limit:       req.Limit,
		MachineType: req.MachineType,
		Status:      "active",
	}

	result, err := s.machineRepo.List(ctx, repoReq)
	if err != nil {
		return nil, fmt.Errorf("获取可用农机失败: %v", err)
	}

	return &GetAvailableMachinesResponse{
		Machines: result.Machines,
		Total:    result.Total,
		Page:     result.Page,
		Limit:    result.Limit,
	}, nil
}

// ==================== 租赁订单 ====================

// CreateRentalOrder 创建租赁订单
func (s *machineService) CreateRentalOrder(ctx context.Context, req *CreateRentalOrderRequest) (*CreateRentalOrderResponse, error) {
	// TODO: 从上下文中获取用户ID
	// renterID := getUserIDFromContext(ctx)

	// 获取农机信息
	machine, err := s.machineRepo.GetByID(ctx, req.MachineID)
	if err != nil {
		return nil, fmt.Errorf("获取农机信息失败: %v", err)
	}

	// 生成订单号
	orderNo := generateOrderNo()

	// 计算订单金额
	var totalAmount int64
	switch req.BillingMethod {
	case "hourly":
		totalAmount = machine.HourlyRate * int64(req.Quantity)
	case "daily":
		totalAmount = machine.DailyRate * int64(req.Quantity)
	case "per_acre":
		totalAmount = machine.PerAcreRate * int64(req.Quantity)
	}

	order := &model.RentalOrder{
		OrderNo:   orderNo,
		MachineID: req.MachineID,
		// RenterID:       renterID,
		// OwnerID:        machine.OwnerID,
		RentalLocation: req.RentalLocation,
		ContactPerson:  req.ContactPerson,
		ContactPhone:   req.ContactPhone,
		BillingMethod:  req.BillingMethod,
		Quantity:       req.Quantity,
		TotalAmount:    totalAmount,
		DepositAmount:  machine.DepositAmount,
		Remarks:        req.Remarks,
		Status:         "pending", // 待确认
		// StartTime:      parseTime(req.StartTime),
		// EndTime:        parseTime(req.EndTime),
	}

	err = s.machineRepo.CreateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("创建订单失败: %v", err)
	}

	return &CreateRentalOrderResponse{
		ID:            order.ID,
		OrderNo:       order.OrderNo,
		TotalAmount:   order.TotalAmount,
		DepositAmount: order.DepositAmount,
	}, nil
}

// GetRentalOrder 获取订单详情
func (s *machineService) GetRentalOrder(ctx context.Context, orderID uint64) (*RentalOrderDetailResponse, error) {
	order, err := s.machineRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("获取订单信息失败: %v", err)
	}

	// 获取关联信息
	machine, _ := s.machineRepo.GetByID(ctx, order.MachineID)
	// renter, _ := s.userRepo.GetByID(ctx, order.RenterID)
	// owner, _ := s.userRepo.GetByID(ctx, order.OwnerID)

	return &RentalOrderDetailResponse{
		Order:   order,
		Machine: machine,
		// Renter:  renter,
		// Owner:   owner,
	}, nil
}

// ConfirmOrder 确认订单
func (s *machineService) ConfirmOrder(ctx context.Context, orderID uint64, req *ConfirmOrderRequest) error {
	order, err := s.machineRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("获取订单信息失败: %v", err)
	}

	// 根据确认类型更新状态
	switch req.ConfirmType {
	case "owner":
		order.Status = "confirmed"
	case "renter":
		order.Status = "paid" // 租户确认后转为待支付
	}

	return s.machineRepo.UpdateOrder(ctx, order)
}

// CancelOrder 取消订单
func (s *machineService) CancelOrder(ctx context.Context, orderID uint64, req *CancelOrderRequest) error {
	order, err := s.machineRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("获取订单信息失败: %v", err)
	}

	order.Status = "cancelled"
	// TODO: 需要在RentalOrder模型中添加CancelReason字段
	// order.CancelReason = req.CancelReason

	return s.machineRepo.UpdateOrder(ctx, order)
}

// GetUserOrders 获取用户订单
func (s *machineService) GetUserOrders(ctx context.Context, userID uint64, req *GetUserOrdersRequest) (*GetUserOrdersResponse, error) {
	orders, err := s.machineRepo.GetUserOrders(ctx, userID, req.OrderType, req.Limit, (req.Page-1)*req.Limit)
	if err != nil {
		return nil, fmt.Errorf("获取用户订单失败: %v", err)
	}

	return &GetUserOrdersResponse{
		Orders: orders,
		Total:  int64(len(orders)), // TODO: 需要实际的总数
		Page:   req.Page,
		Limit:  req.Limit,
	}, nil
}

// ==================== 订单支付 ====================

// PayOrder 支付订单
func (s *machineService) PayOrder(ctx context.Context, orderID uint64, req *PayOrderRequest) (*PayOrderResponse, error) {
	order, err := s.machineRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("获取订单信息失败: %v", err)
	}

	// TODO: 集成支付服务
	// 1. 调用支付接口
	// 2. 生成支付订单
	// 3. 返回支付参数

	return &PayOrderResponse{
		PaymentID:     "pay_" + strconv.FormatUint(orderID, 10),
		PaymentAmount: order.TotalAmount + order.DepositAmount,
		PaymentURL:    "", // 支付链接（如果需要）
	}, nil
}

// CompleteOrder 完成订单
func (s *machineService) CompleteOrder(ctx context.Context, orderID uint64, req *CompleteOrderRequest) error {
	order, err := s.machineRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("获取订单信息失败: %v", err)
	}

	order.Status = "completed"
	// TODO: 更新实际结束时间等字段
	// order.ActualEndTime = parseTime(req.ActualEndTime)
	// order.WorkHours = req.WorkHours
	// order.WorkAcres = req.WorkAcres
	// order.Notes = req.Notes

	return s.machineRepo.UpdateOrder(ctx, order)
}

// ==================== 评价系统 ====================

// SubmitRating 提交评价
func (s *machineService) SubmitRating(ctx context.Context, orderID uint64, req *SubmitRatingRequest) error {
	// TODO: 实现评价逻辑
	// 1. 验证订单是否已完成
	// 2. 创建评价记录
	// 3. 更新农机平均评分

	return nil
}

// GetMachineRatings 获取农机评价
func (s *machineService) GetMachineRatings(ctx context.Context, machineID uint64) (*MachineRatingsResponse, error) {
	// TODO: 实现获取评价逻辑
	return &MachineRatingsResponse{
		Ratings:       []MachineRating{},
		AverageRating: 0.0,
		TotalCount:    0,
	}, nil
}

// ==================== 工具函数 ====================

// generateMachineCode 生成农机编码
func generateMachineCode() string {
	return "M" + strconv.FormatInt(time.Now().Unix(), 10)
}

// generateOrderNo 生成订单号
func generateOrderNo() string {
	return "O" + strconv.FormatInt(time.Now().Unix(), 10)
}
