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

// NewLoanHandler 创建贷款处理器实例
func NewLoanHandler(loanService service.LoanService) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
	}
}

// GetProducts 获取贷款产品列表
// @Summary 获取贷款产品列表
// @Description 获取可用的贷款产品
// @Tags 贷款管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Security BearerAuth
// @Success 200 {object} service.GetProductsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /loan/products [get]
func (h *LoanHandler) GetProducts(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	req := &service.GetProductsRequest{
		Page:  page,
		Limit: limit,
	}

	// 调用服务
	resp, err := h.loanService.GetProducts(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取贷款产品失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    http.StatusOK,
		Message: "获取贷款产品成功",
		Data:    resp,
	})
}

// CreateApplication 提交贷款申请
// @Summary 提交贷款申请
// @Description 创建新的贷款申请
// @Tags 贷款管理
// @Accept json
// @Produce json
// @Param request body service.CreateApplicationRequest true "申请信息"
// @Security BearerAuth
// @Success 200 {object} service.CreateApplicationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /loan/applications [post]
func (h *LoanHandler) CreateApplication(c *gin.Context) {
	var req service.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	// 调用服务
	resp, err := h.loanService.CreateApplication(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "创建贷款申请失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    http.StatusOK,
		Message: "贷款申请提交成功",
		Data:    resp,
	})
}

// GetApplicationDetails 获取申请详情
// @Summary 获取申请详情
// @Description 获取贷款申请的详细信息
// @Tags 贷款管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Security BearerAuth
// @Success 200 {object} service.GetApplicationDetailsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /loan/applications/{id} [get]
func (h *LoanHandler) GetApplicationDetails(c *gin.Context) {
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

	req := &service.GetApplicationDetailsRequest{
		ID: uint(id),
	}

	// 调用服务
	resp, err := h.loanService.GetApplicationDetails(c.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "申请不存在" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "无权限查看此申请" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, ErrorResponse{
			Code:    statusCode,
			Message: "获取申请详情失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    http.StatusOK,
		Message: "获取申请详情成功",
		Data:    resp,
	})
}

// GetUserApplications 获取用户申请列表
// @Summary 获取用户申请列表
// @Description 获取当前用户的贷款申请列表
// @Tags 贷款管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Param status query string false "申请状态"
// @Security BearerAuth
// @Success 200 {object} service.GetUserApplicationsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /loan/applications [get]
func (h *LoanHandler) GetUserApplications(c *gin.Context) {
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

	req := &service.GetUserApplicationsRequest{
		Page:   page,
		Limit:  limit,
		Status: status,
	}

	// 调用服务
	resp, err := h.loanService.GetUserApplications(c.Request.Context(), req)
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
