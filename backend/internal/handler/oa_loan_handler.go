package handler

import (
	"net/http"
	"strconv"

	"huinong-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// OALoanHandler OA贷款审批处理器
type OALoanHandler struct {
	loanService service.LoanService
	oaService   service.OAService
}

// NewOALoanHandler 创建OA贷款审批处理器
func NewOALoanHandler(loanService service.LoanService, oaService service.OAService) *OALoanHandler {
	return &OALoanHandler{
		loanService: loanService,
		oaService:   oaService,
	}
}

// GetApplications 获取申请列表(管理员视图)
// @Summary 获取贷款申请列表
// @Description 获取所有贷款申请列表，支持多种筛选条件
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Param status query string false "申请状态"
// @Param product_id query int false "产品ID"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Success 200 {object} StandardResponse{data=service.GetAdminApplicationsResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/loans/applications [get]
func (h *OALoanHandler) GetApplications(c *gin.Context) {
	status := c.Query("status")
	// productIDStr := c.Query("product_id")  // TODO: 后续扩展产品筛选功能
	// startDate := c.Query("start_date")     // TODO: 后续扩展日期筛选功能
	// endDate := c.Query("end_date")         // TODO: 后续扩展日期筛选功能
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

	// var productID uint64 = 0  // TODO: 后续使用产品筛选
	// if productIDStr != "" {
	//     productID, _ = strconv.ParseUint(productIDStr, 10, 64)
	// }

	req := &service.GetAdminApplicationsRequest{
		Status: status,
		Page:   page,
		Limit:  limit,
	}

	response, err := h.loanService.GetAdminApplications(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取申请列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// GetApplicationDetail 获取申请详情(管理员视图)
// @Summary 获取申请详情
// @Description 获取指定申请的详细信息，包括用户信息、AI评估结果等
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Success 200 {object} StandardResponse{data=service.GetApplicationDetailsResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/loans/applications/{id} [get]
func (h *OALoanHandler) GetApplicationDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的申请ID", err.Error()))
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
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取申请详情失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// ApproveApplication 批准申请
// @Summary 人工批准申请
// @Description 审批员人工批准贷款申请
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Param request body service.ApproveApplicationRequest true "批准信息"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/loans/applications/{id}/approve [post]
func (h *OALoanHandler) ApproveApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的申请ID", err.Error()))
		return
	}

	var req service.ApproveApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	// 设置申请ID
	req.ID = uint(id)

	// 从上下文获取审批员ID
	_, exists := c.Get("oaUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "审批员未登录", "OA用户认证信息缺失"))
		return
	}

	// 这里需要扩展ApproveApplicationRequest来包含审批员ID
	// 暂时通过context传递

	err = h.loanService.ApproveApplication(c.Request.Context(), &req)
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

// RejectApplication 拒绝申请
// @Summary 人工拒绝申请
// @Description 审批员人工拒绝贷款申请
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Param request body service.RejectApplicationRequest true "拒绝信息"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/loans/applications/{id}/reject [post]
func (h *OALoanHandler) RejectApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的申请ID", err.Error()))
		return
	}

	var req service.RejectApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	// 设置申请ID
	req.ID = uint(id)

	// 从上下文获取审批员ID
	_, exists := c.Get("oaUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "审批员未登录", "OA用户认证信息缺失"))
		return
	}

	err = h.loanService.RejectApplication(c.Request.Context(), &req)
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

// ReturnApplication 退回申请
// @Summary 退回申请补充材料
// @Description 审批员要求申请人补充材料
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Param request body service.ReturnRequest true "退回信息"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/loans/applications/{id}/return [post]
func (h *OALoanHandler) ReturnApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的申请ID", err.Error()))
		return
	}

	var req service.ReturnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	// 从上下文获取审批员ID
	_, exists := c.Get("oaUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "审批员未登录", "OA用户认证信息缺失"))
		return
	}

	err = h.loanService.ReturnApplication(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "申请不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "申请不存在", err.Error()))
			return
		}
		if err.Error() == "申请状态不允许退回" {
			c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(http.StatusUnprocessableEntity, "申请状态不允许退回", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "退回申请失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("申请已退回", nil))
}

// StartReview 开始审核
// @Summary 开始审核申请
// @Description 审批员开始审核指定申请，更新审核状态
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/loans/applications/{id}/start-review [post]
func (h *OALoanHandler) StartReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的申请ID", err.Error()))
		return
	}

	// 从上下文获取审批员ID
	oaUserID, exists := c.Get("oaUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "审批员未登录", "OA用户认证信息缺失"))
		return
	}

	err = h.loanService.StartReview(c.Request.Context(), id, oaUserID.(uint64))
	if err != nil {
		if err.Error() == "申请不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "申请不存在", err.Error()))
			return
		}
		if err.Error() == "申请状态不允许审核" {
			c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(http.StatusUnprocessableEntity, "申请状态不允许审核", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "开始审核失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("已开始审核", nil))
}

// GetStatistics 获取统计数据
// @Summary 获取贷款申请统计
// @Description 获取贷款申请的各种统计数据
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Success 200 {object} StandardResponse{data=service.LoanStatisticsResponse}
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/loans/statistics [get]
func (h *OALoanHandler) GetStatistics(c *gin.Context) {
	response, err := h.loanService.GetStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取统计数据失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// RetryAIAssessment 重试AI评估
// @Summary 重试AI评估
// @Description 对AI评估失败的申请重新触发评估
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/loans/applications/{id}/retry-ai [post]
func (h *OALoanHandler) RetryAIAssessment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的申请ID", err.Error()))
		return
	}

	// 从上下文获取操作员ID
	_, exists := c.Get("oaUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "操作员未登录", "OA用户认证信息缺失"))
		return
	}

	err = h.loanService.TriggerAIAssessment(c.Request.Context(), id, "loan_approval")
	if err != nil {
		if err.Error() == "申请不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "申请不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "重试AI评估失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("AI评估已重新启动", nil))
}
