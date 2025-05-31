package handler

import (
	"net/http"
	"strconv"

	"huinong-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// MachineHandler 农机处理器
type MachineHandler struct {
	machineService service.MachineService
}

// NewMachineHandler 创建农机处理器
func NewMachineHandler(machineService service.MachineService) *MachineHandler {
	return &MachineHandler{
		machineService: machineService,
	}
}

// RegisterMachine 注册农机
// @Summary 注册农机
// @Description 农机所有者注册新的农机设备
// @Tags 农机管理
// @Accept json
// @Produce json
// @Param request body service.RegisterMachineRequest true "农机信息"
// @Success 200 {object} StandardResponse{data=service.RegisterMachineResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/machines [post]
func (h *MachineHandler) RegisterMachine(c *gin.Context) {
	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	var req service.RegisterMachineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	response, err := h.machineService.RegisterMachine(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "注册农机失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("农机注册成功", response))
}

// GetUserMachines 获取我的农机列表
// @Summary 获取我的农机列表
// @Description 获取当前用户的所有农机设备
// @Tags 农机管理
// @Accept json
// @Produce json
// @Param status query string false "设备状态"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Success 200 {object} StandardResponse{data=service.GetUserMachinesResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/machines [get]
func (h *MachineHandler) GetUserMachines(c *gin.Context) {
	// 从上下文获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	userID, ok := userIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "用户ID格式错误", "userID type assertion failed"))
		return
	}

	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	req := &service.GetUserMachinesRequest{
		Status: status,
		Page:   page,
		Limit:  limit,
	}

	response, err := h.machineService.GetUserMachines(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取农机列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// SearchMachines 搜索附近农机
// @Summary 搜索附近农机
// @Description 根据位置和条件搜索可租赁的农机
// @Tags 农机租赁
// @Accept json
// @Produce json
// @Param longitude query float64 true "经度"
// @Param latitude query float64 true "纬度"
// @Param radius query int false "搜索半径(公里)" default(10)
// @Param machine_type query string false "设备类型"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Success 200 {object} StandardResponse{data=service.SearchMachinesResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/machines/search [get]
func (h *MachineHandler) SearchMachines(c *gin.Context) {
	longitudeStr := c.Query("longitude")
	latitudeStr := c.Query("latitude")
	radiusStr := c.DefaultQuery("radius", "10")
	machineType := c.Query("machine_type")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	if longitudeStr == "" || latitudeStr == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "经纬度不能为空", "longitude and latitude are required"))
		return
	}

	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的经度", err.Error()))
		return
	}

	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的纬度", err.Error()))
		return
	}

	radius, err := strconv.Atoi(radiusStr)
	if err != nil || radius < 1 || radius > 100 {
		radius = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	req := &service.SearchMachinesRequest{
		Longitude:   longitude,
		Latitude:    latitude,
		Radius:      radius,
		MachineType: machineType,
		StartTime:   startTime,
		EndTime:     endTime,
		Page:        page,
		Limit:       limit,
	}

	response, err := h.machineService.SearchMachines(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "搜索农机失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("搜索成功", response))
}

// GetMachine 获取农机详情
// @Summary 获取农机详情
// @Description 获取指定农机的详细信息
// @Tags 农机租赁
// @Accept json
// @Produce json
// @Param id path int true "农机ID"
// @Success 200 {object} StandardResponse{data=service.MachineDetailResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/machines/{id} [get]
func (h *MachineHandler) GetMachine(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的农机ID", err.Error()))
		return
	}

	response, err := h.machineService.GetMachine(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "农机不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "农机不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取农机详情失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// CreateOrder 创建租赁订单
// @Summary 创建租赁订单
// @Description 为指定农机创建租赁订单
// @Tags 农机租赁
// @Accept json
// @Produce json
// @Param id path int true "农机ID"
// @Param request body service.CreateRentalOrderRequest true "订单信息"
// @Success 200 {object} StandardResponse{data=service.CreateRentalOrderResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/machines/{id}/orders [post]
func (h *MachineHandler) CreateOrder(c *gin.Context) {
	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	idStr := c.Param("id")
	machineID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的农机ID", err.Error()))
		return
	}

	var req service.CreateRentalOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	// 设置农机ID
	req.MachineID = machineID

	response, err := h.machineService.CreateRentalOrder(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "农机不存在" || err.Error() == "农机不可用" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, err.Error(), err.Error()))
			return
		}
		if err.Error() == "时间段冲突" || err.Error() == "租赁时间无效" {
			c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(http.StatusUnprocessableEntity, err.Error(), err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "创建订单失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("订单创建成功", response))
}

// GetUserOrders 获取我的订单
// @Summary 获取我的订单
// @Description 获取当前用户的租赁订单列表
// @Tags 农机租赁
// @Accept json
// @Produce json
// @Param order_type query string false "订单类型" Enums(renter,owner)
// @Param status query string false "订单状态"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Success 200 {object} StandardResponse{data=service.GetUserOrdersResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/orders [get]
func (h *MachineHandler) GetUserOrders(c *gin.Context) {
	// 从上下文获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	userID, ok := userIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "用户ID格式错误", "userID type assertion failed"))
		return
	}

	orderType := c.Query("order_type")
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	req := &service.GetUserOrdersRequest{
		OrderType: orderType,
		Status:    status,
		Page:      page,
		Limit:     limit,
	}

	response, err := h.machineService.GetUserOrders(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取订单列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// ConfirmOrder 确认订单
// @Summary 确认订单
// @Description 农机所有者确认租赁订单
// @Tags 农机租赁
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param request body service.ConfirmOrderRequest true "确认信息"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/orders/{id}/confirm [put]
func (h *MachineHandler) ConfirmOrder(c *gin.Context) {
	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的订单ID", err.Error()))
		return
	}

	var req service.ConfirmOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	err = h.machineService.ConfirmOrder(c.Request.Context(), orderID, &req)
	if err != nil {
		if err.Error() == "订单不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "订单不存在", err.Error()))
			return
		}
		if err.Error() == "无权操作该订单" {
			c.JSON(http.StatusForbidden, NewErrorResponse(http.StatusForbidden, "无权操作该订单", err.Error()))
			return
		}
		if err.Error() == "订单状态不允许确认" {
			c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(http.StatusUnprocessableEntity, "订单状态不允许确认", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "确认订单失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("订单已确认", nil))
}

// PayOrder 支付订单
// @Summary 支付订单
// @Description 租赁者支付订单费用
// @Tags 农机租赁
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param request body service.PayOrderRequest true "支付信息"
// @Success 200 {object} StandardResponse{data=service.PayOrderResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/orders/{id}/pay [post]
func (h *MachineHandler) PayOrder(c *gin.Context) {
	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的订单ID", err.Error()))
		return
	}

	var req service.PayOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	response, err := h.machineService.PayOrder(c.Request.Context(), orderID, &req)
	if err != nil {
		if err.Error() == "订单不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "订单不存在", err.Error()))
			return
		}
		if err.Error() == "无权操作该订单" {
			c.JSON(http.StatusForbidden, NewErrorResponse(http.StatusForbidden, "无权操作该订单", err.Error()))
			return
		}
		if err.Error() == "订单状态不允许支付" {
			c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(http.StatusUnprocessableEntity, "订单状态不允许支付", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "订单支付失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("支付成功", response))
}

// CompleteOrder 完成订单
// @Summary 完成订单
// @Description 标记订单为已完成
// @Tags 农机租赁
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param request body service.CompleteOrderRequest true "完成信息"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/orders/{id}/complete [put]
func (h *MachineHandler) CompleteOrder(c *gin.Context) {
	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的订单ID", err.Error()))
		return
	}

	var req service.CompleteOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	err = h.machineService.CompleteOrder(c.Request.Context(), orderID, &req)
	if err != nil {
		if err.Error() == "订单不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "订单不存在", err.Error()))
			return
		}
		if err.Error() == "订单状态不允许完成" {
			c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(http.StatusUnprocessableEntity, "订单状态不允许完成", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "完成订单失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("订单已完成", nil))
}

// CancelOrder 取消订单
// @Summary 取消订单
// @Description 取消租赁订单
// @Tags 农机租赁
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param request body service.CancelOrderRequest true "取消原因"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/orders/{id}/cancel [put]
func (h *MachineHandler) CancelOrder(c *gin.Context) {
	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的订单ID", err.Error()))
		return
	}

	var req service.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	err = h.machineService.CancelOrder(c.Request.Context(), orderID, &req)
	if err != nil {
		if err.Error() == "订单不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "订单不存在", err.Error()))
			return
		}
		if err.Error() == "无权操作该订单" {
			c.JSON(http.StatusForbidden, NewErrorResponse(http.StatusForbidden, "无权操作该订单", err.Error()))
			return
		}
		if err.Error() == "订单状态不允许取消" {
			c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(http.StatusUnprocessableEntity, "订单状态不允许取消", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "取消订单失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("订单已取消", nil))
}

// RateOrder 评价订单
// @Summary 评价订单
// @Description 用户对完成的订单进行评价
// @Tags 农机租赁
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param request body service.RateOrderRequest true "评价信息"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/orders/{id}/rate [post]
func (h *MachineHandler) RateOrder(c *gin.Context) {
	// 从上下文获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	userID, ok := userIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "用户ID格式错误", "userID type assertion failed"))
		return
	}

	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的订单ID", err.Error()))
		return
	}

	var req service.RateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	req.OrderID = orderID
	req.UserID = userID

	err = h.machineService.RateOrder(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "订单不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "订单不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "评价失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("评价成功", nil))
}

// ============= OA管理员农机租赁审批功能 =============

// GetRentalApplications 获取农机租赁申请列表(OA管理员)
// @Summary 获取农机租赁申请列表
// @Description OA管理员获取所有农机租赁申请，支持筛选和分页
// @Tags OA农机管理
// @Accept json
// @Produce json
// @Param status query string false "申请状态 (pending/approved/rejected/cancelled)"
// @Param machine_type query string false "设备类型筛选"
// @Param start_date query string false "申请开始日期"
// @Param end_date query string false "申请结束日期"
// @Param risk_level query string false "风险等级 (low/medium/high)"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Param sort_by query string false "排序字段" default(created_at)
// @Param sort_order query string false "排序方向" default(desc)
// @Success 200 {object} StandardResponse{data=service.GetRentalApplicationsResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/machines/rental-applications [get]
func (h *MachineHandler) GetRentalApplications(c *gin.Context) {
	status := c.Query("status")
	machineType := c.Query("machine_type")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	riskLevel := c.Query("risk_level")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	req := &service.GetRentalApplicationsRequest{
		Status:      status,
		MachineType: machineType,
		StartDate:   startDate,
		EndDate:     endDate,
		RiskLevel:   riskLevel,
		Page:        page,
		Limit:       limit,
		SortBy:      sortBy,
		SortOrder:   sortOrder,
	}

	response, err := h.machineService.GetRentalApplications(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取租赁申请列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// GetRentalApplicationDetail 获取农机租赁申请详情(OA管理员)
// @Summary 获取农机租赁申请详情
// @Description OA管理员获取指定租赁申请的详细信息
// @Tags OA农机管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Success 200 {object} StandardResponse{data=service.GetRentalApplicationDetailResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/machines/rental-applications/{id} [get]
func (h *MachineHandler) GetRentalApplicationDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的申请ID", err.Error()))
		return
	}

	req := &service.GetRentalApplicationDetailRequest{
		ID: uint(id),
	}

	response, err := h.machineService.GetRentalApplicationDetail(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "申请不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "申请不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取申请详情失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// ApproveRentalApplication 批准农机租赁申请(OA管理员)
// @Summary 批准农机租赁申请
// @Description OA管理员批准农机租赁申请
// @Tags OA农机管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Param request body service.ApproveRentalRequest true "批准信息"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/machines/rental-applications/{id}/approve [post]
func (h *MachineHandler) ApproveRentalApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的申请ID", err.Error()))
		return
	}

	var req service.ApproveRentalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	// 从上下文获取审批员ID
	reviewerIDInterface, exists := c.Get("oaUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "审批员未登录", "OA用户认证信息缺失"))
		return
	}

	reviewerID, ok := reviewerIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "审批员ID格式错误", "reviewerID type assertion failed"))
		return
	}

	req.ID = uint(id)
	req.ReviewerID = reviewerID

	err = h.machineService.ApproveRentalApplication(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "申请不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "申请不存在", err.Error()))
			return
		}
		if err.Error() == "申请状态不允许审批" {
			c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(http.StatusUnprocessableEntity, "申请状态不允许审批", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "批准申请失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("申请已批准", nil))
}

// RejectRentalApplication 拒绝农机租赁申请(OA管理员)
// @Summary 拒绝农机租赁申请
// @Description OA管理员拒绝农机租赁申请
// @Tags OA农机管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Param request body service.RejectRentalRequest true "拒绝信息"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/machines/rental-applications/{id}/reject [post]
func (h *MachineHandler) RejectRentalApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的申请ID", err.Error()))
		return
	}

	var req service.RejectRentalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	// 从上下文获取审批员ID
	reviewerIDInterface, exists := c.Get("oaUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "审批员未登录", "OA用户认证信息缺失"))
		return
	}

	reviewerID, ok := reviewerIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "审批员ID格式错误", "reviewerID type assertion failed"))
		return
	}

	req.ID = uint(id)
	req.ReviewerID = reviewerID

	err = h.machineService.RejectRentalApplication(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "申请不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "申请不存在", err.Error()))
			return
		}
		if err.Error() == "申请状态不允许审批" {
			c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(http.StatusUnprocessableEntity, "申请状态不允许审批", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "拒绝申请失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("申请已拒绝", nil))
}

// GetRentalStatistics 获取农机租赁统计(OA管理员)
// @Summary 获取农机租赁统计
// @Description OA管理员获取农机租赁业务统计数据
// @Tags OA农机管理
// @Accept json
// @Produce json
// @Param period query string false "统计周期 (day/week/month/year)" default(month)
// @Param start_date query string false "统计开始日期"
// @Param end_date query string false "统计结束日期"
// @Success 200 {object} StandardResponse{data=service.GetRentalStatisticsResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/machines/rental-statistics [get]
func (h *MachineHandler) GetRentalStatistics(c *gin.Context) {
	period := c.DefaultQuery("period", "month")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	req := &service.GetRentalStatisticsRequest{
		Period:    period,
		StartDate: startDate,
		EndDate:   endDate,
	}

	response, err := h.machineService.GetRentalStatistics(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取统计数据失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// BatchApproveRentals 批量审批农机租赁申请(OA管理员)
// @Summary 批量审批农机租赁申请
// @Description OA管理员批量审批多个农机租赁申请
// @Tags OA农机管理
// @Accept json
// @Produce json
// @Param request body service.BatchApproveRentalsRequest true "批量审批信息"
// @Success 200 {object} StandardResponse{data=service.BatchApproveRentalsResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/machines/rental-applications/batch-approve [post]
func (h *MachineHandler) BatchApproveRentals(c *gin.Context) {
	var req service.BatchApproveRentalsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	// 从上下文获取审批员ID
	reviewerIDInterface, exists := c.Get("oaUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "审批员未登录", "OA用户认证信息缺失"))
		return
	}

	reviewerID, ok := reviewerIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "审批员ID格式错误", "reviewerID type assertion failed"))
		return
	}

	req.ReviewerID = reviewerID

	response, err := h.machineService.BatchApproveRentals(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "批量审批失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("批量审批完成", response))
}
