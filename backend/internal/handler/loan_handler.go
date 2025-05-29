package handler

import (
	"net/http"
	"strconv"

	"huinong-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// LoanHandler 贷款处理器
type LoanHandler struct {
	loanService service.LoanService
}

// NewLoanHandler 创建贷款处理器
func NewLoanHandler(loanService service.LoanService) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
	}
}

// GetProducts 获取贷款产品列表
// @Summary 获取贷款产品列表
// @Description 获取可用的贷款产品列表，支持按用户类型筛选
// @Tags 贷款产品
// @Accept json
// @Produce json
// @Param user_type query string false "用户类型"
// @Param status query string false "产品状态"
// @Success 200 {object} StandardResponse{data=service.GetProductsResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/loan/products [get]
func (h *LoanHandler) GetProducts(c *gin.Context) {
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

	req := &service.GetProductsRequest{
		Page:  page,
		Limit: limit,
	}

	response, err := h.loanService.GetProducts(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取产品列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// GetProductDetail 获取产品详情
// @Summary 获取贷款产品详情
// @Description 根据产品ID获取详细信息
// @Tags 贷款产品
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} StandardResponse{data=model.LoanProduct}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/loan/products/{id} [get]
func (h *LoanHandler) GetProductDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的产品ID", err.Error()))
		return
	}

	product, err := h.loanService.GetProduct(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "产品不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "产品不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取产品详情失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", product))
}

// SubmitApplication 提交贷款申请
// @Summary 提交贷款申请
// @Description 用户提交贷款申请，系统自动触发AI评估
// @Tags 贷款申请
// @Accept json
// @Produce json
// @Param request body service.CreateApplicationRequest true "申请信息"
// @Success 200 {object} StandardResponse{data=service.CreateApplicationResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/loan/applications [post]
func (h *LoanHandler) SubmitApplication(c *gin.Context) {
	var req service.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	// 从上下文获取用户ID (注意：此处需要在service层内部处理用户ID，不在request中)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	// 这里我们需要创建一个临时的request结构，包含用户ID
	// 由于service.CreateApplicationRequest没有UserID字段，我们需要在service层内部处理
	// 暂时通过context传递用户ID
	ctx := c.Request.Context()

	response, err := h.loanService.CreateApplication(ctx, &req)
	if err != nil {
		if err.Error() == "产品不存在" || err.Error() == "申请金额超出产品限制" {
			c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(http.StatusUnprocessableEntity, err.Error(), err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "提交申请失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("申请提交成功，AI评估正在进行中", response))
}

// GetUserApplications 获取我的申请列表
// @Summary 获取我的申请列表
// @Description 获取当前用户的贷款申请列表
// @Tags 贷款申请
// @Accept json
// @Produce json
// @Param status query string false "申请状态"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Success 200 {object} StandardResponse{data=service.GetUserApplicationsResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/loan/applications [get]
func (h *LoanHandler) GetUserApplications(c *gin.Context) {
	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
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

	req := &service.GetUserApplicationsRequest{
		Page:   page,
		Limit:  limit,
		Status: status,
	}

	response, err := h.loanService.GetUserApplications(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取申请列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// GetApplicationDetail 获取申请详情
// @Summary 获取申请详情
// @Description 获取指定申请的详细信息，包括AI评估结果和审批日志
// @Tags 贷款申请
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Success 200 {object} StandardResponse{data=service.GetApplicationDetailsResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/loan/applications/{id} [get]
func (h *LoanHandler) GetApplicationDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的申请ID", err.Error()))
		return
	}

	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	req := &service.GetApplicationDetailsRequest{
		ID: uint(id),
	}

	response, err := h.loanService.GetApplicationDetails(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "申请不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "申请不存在", err.Error()))
			return
		}
		if err.Error() == "无权访问该申请" {
			c.JSON(http.StatusForbidden, NewErrorResponse(http.StatusForbidden, "无权访问该申请", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取申请详情失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// CancelApplication 取消申请
// @Summary 取消申请
// @Description 取消指定的贷款申请（仅限pending和待审核状态）
// @Tags 贷款申请
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/loan/applications/{id} [delete]
func (h *LoanHandler) CancelApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的申请ID", err.Error()))
		return
	}

	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	// 获取申请详情检查权限和状态
	req := &service.GetApplicationDetailsRequest{
		ID: uint(id),
	}

	application, err := h.loanService.GetApplicationDetails(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "申请不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "申请不存在", err.Error()))
			return
		}
		if err.Error() == "无权访问该申请" {
			c.JSON(http.StatusForbidden, NewErrorResponse(http.StatusForbidden, "无权访问该申请", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取申请信息失败", err.Error()))
		return
	}

	// 检查是否可以取消
	if application.Application.Status != "pending" && application.Application.Status != "ai_processing" && application.Application.Status != "manual_review" {
		c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(http.StatusUnprocessableEntity, "当前状态不允许取消", "申请状态不支持取消操作"))
		return
	}

	// 获取申请对象并更新状态
	app, err := h.loanService.GetLoanApplicationByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取申请信息失败", err.Error()))
		return
	}

	// 更新申请状态为已取消
	app.Status = "cancelled"
	if err := h.loanService.UpdateLoanApplication(app); err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "取消申请失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("申请已取消", nil))
}

// ApproveApplication 审批通过申请
// @Summary 审批通过申请
// @Description 管理员审批通过贷款申请
// @Tags 贷款管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Param request body service.ApproveApplicationRequest true "审批信息"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/loan/applications/{id}/approve [put]
func (h *LoanHandler) ApproveApplication(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "无效的申请ID",
			Error:   err.Error(),
		})
		return
	}

	var reqBody struct {
		ApprovalNote   string  `json:"approval_note"`
		ApprovedAmount int64   `json:"approved_amount"`
		ApprovedTerms  int     `json:"approved_terms"`
		ApprovedRate   float64 `json:"approved_rate"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	req := &service.ApproveApplicationRequest{
		ID:           uint(id),
		ApprovalNote: reqBody.ApprovalNote,
	}

	// 调用服务
	err = h.loanService.ApproveApplication(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "审批申请失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    http.StatusOK,
		Message: "申请审批通过",
	})
}

// RejectApplication 拒绝申请
// @Summary 拒绝申请
// @Description 管理员拒绝贷款申请
// @Tags 贷款管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Param request body service.RejectApplicationRequest true "拒绝信息"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/loan/applications/{id}/reject [put]
func (h *LoanHandler) RejectApplication(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "无效的申请ID",
			Error:   err.Error(),
		})
		return
	}

	var reqBody struct {
		RejectionReason string `json:"rejection_reason"`
		RejectionNote   string `json:"rejection_note"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	req := &service.RejectApplicationRequest{
		ID:            uint(id),
		RejectionNote: reqBody.RejectionNote,
	}

	// 调用服务
	err = h.loanService.RejectApplication(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "拒绝申请失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    http.StatusOK,
		Message: "申请已拒绝",
	})
}

// GetAdminApplications 管理员获取申请列表
// @Summary 管理员获取申请列表
// @Description 管理员获取所有贷款申请列表
// @Tags 贷款管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Param status query string false "申请状态"
// @Security BearerAuth
// @Success 200 {object} service.GetAdminApplicationsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/loan/applications [get]
func (h *LoanHandler) GetAdminApplications(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	req := &service.GetAdminApplicationsRequest{
		Page:   page,
		Limit:  limit,
		Status: status,
	}

	// 调用服务
	resp, err := h.loanService.GetAdminApplications(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取申请列表失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    http.StatusOK,
		Message: "获取申请列表成功",
		Data:    resp,
	})
}

// GetStatistics 获取贷款统计数据
// @Summary 获取贷款统计数据
// @Description 获取贷款申请的统计信息
// @Tags 贷款管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} service.LoanStatisticsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/loan/statistics [get]
func (h *LoanHandler) GetStatistics(c *gin.Context) {
	// 调用服务
	resp, err := h.loanService.GetStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取统计数据失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    http.StatusOK,
		Message: "获取统计数据成功",
		Data:    resp,
	})
}
