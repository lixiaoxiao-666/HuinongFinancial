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
// @Description 重新触发AI风险评估流程
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/loans/applications/{id}/retry-ai [post]
func (h *OALoanHandler) RetryAIAssessment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的申请ID", err.Error()))
		return
	}

	req := &service.RetryAIAssessmentRequest{
		ID: uint(id),
	}

	err = h.loanService.RetryAIAssessment(c.Request.Context(), req)
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

// ============= 新增的增强审批功能 =============

// BatchApproveLoanApplications 批量审批贷款申请
// @Summary 批量审批贷款申请
// @Description 批量审批多个贷款申请，支持统一条件设置
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Param request body service.BatchApproveLoanRequest true "批量审批信息"
// @Success 200 {object} StandardResponse{data=service.BatchApproveLoanResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/loans/applications/batch-approve [post]
func (h *OALoanHandler) BatchApproveLoanApplications(c *gin.Context) {
	var req service.BatchApproveLoanRequest
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

	response, err := h.loanService.BatchApproveLoanApplications(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "批量审批失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("批量审批完成", response))
}

// EnableAutoApproval 启用自动审批
// @Summary 启用自动审批
// @Description 为指定条件的申请启用自动审批功能
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Param request body service.EnableAutoApprovalRequest true "自动审批配置"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/loans/auto-approval/enable [post]
func (h *OALoanHandler) EnableAutoApproval(c *gin.Context) {
	var req service.EnableAutoApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	// 从上下文获取操作员ID
	operatorIDInterface, exists := c.Get("oaUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "操作员未登录", "OA用户认证信息缺失"))
		return
	}

	operatorID, ok := operatorIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "操作员ID格式错误", "operatorID type assertion failed"))
		return
	}

	req.OperatorID = operatorID

	err := h.loanService.EnableAutoApproval(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "启用自动审批失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("自动审批已启用", nil))
}

// DisableAutoApproval 禁用自动审批
// @Summary 禁用自动审批
// @Description 禁用自动审批功能
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Param request body service.DisableAutoApprovalRequest true "禁用配置"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/loans/auto-approval/disable [post]
func (h *OALoanHandler) DisableAutoApproval(c *gin.Context) {
	var req service.DisableAutoApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	// 从上下文获取操作员ID
	operatorIDInterface, exists := c.Get("oaUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "操作员未登录", "OA用户认证信息缺失"))
		return
	}

	operatorID, ok := operatorIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "操作员ID格式错误", "operatorID type assertion failed"))
		return
	}

	req.OperatorID = operatorID

	err := h.loanService.DisableAutoApproval(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "禁用自动审批失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("自动审批已禁用", nil))
}

// GetAutoApprovalConfig 获取自动审批配置
// @Summary 获取自动审批配置
// @Description 获取当前自动审批配置信息
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Success 200 {object} StandardResponse{data=service.GetAutoApprovalConfigResponse}
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/loans/auto-approval/config [get]
func (h *OALoanHandler) GetAutoApprovalConfig(c *gin.Context) {
	response, err := h.loanService.GetAutoApprovalConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取自动审批配置失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// GetApplicationsByRiskLevel 按风险等级获取申请
// @Summary 按风险等级获取申请
// @Description 根据AI评估的风险等级筛选申请
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Param risk_level query string true "风险等级 (low/medium/high)"
// @Param status query string false "申请状态"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Success 200 {object} StandardResponse{data=service.GetApplicationsByRiskLevelResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/loans/applications/by-risk-level [get]
func (h *OALoanHandler) GetApplicationsByRiskLevel(c *gin.Context) {
	riskLevel := c.Query("risk_level")
	if riskLevel == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "风险等级不能为空", "risk_level is required"))
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

	req := &service.GetApplicationsByRiskLevelRequest{
		RiskLevel: riskLevel,
		Status:    status,
		Page:      page,
		Limit:     limit,
	}

	response, err := h.loanService.GetApplicationsByRiskLevel(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取申请列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// GetAIAssessmentHistory 获取AI评估历史
// @Summary 获取AI评估历史
// @Description 获取指定申请的AI评估历史记录
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Success 200 {object} StandardResponse{data=service.GetAIAssessmentHistoryResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/loans/applications/{id}/ai-assessment-history [get]
func (h *OALoanHandler) GetAIAssessmentHistory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的申请ID", err.Error()))
		return
	}

	req := &service.GetAIAssessmentHistoryRequest{
		ApplicationID: uint(id),
	}

	response, err := h.loanService.GetAIAssessmentHistory(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "申请不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "申请不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取AI评估历史失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// CreateApplicationTask 创建申请任务
// @Summary 创建申请任务
// @Description 为贷款申请创建审批任务
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Param request body service.CreateApplicationTaskRequest true "任务信息"
// @Success 200 {object} StandardResponse{data=service.CreateApplicationTaskResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/loans/applications/create-task [post]
func (h *OALoanHandler) CreateApplicationTask(c *gin.Context) {
	var req service.CreateApplicationTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	// 从上下文获取创建者ID
	creatorIDInterface, exists := c.Get("oaUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "OA用户认证信息缺失"))
		return
	}

	creatorID, ok := creatorIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "用户ID格式错误", "creatorID type assertion failed"))
		return
	}

	req.CreatorID = creatorID

	response, err := h.loanService.CreateApplicationTask(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "创建任务失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("任务创建成功", response))
}

// GetApplicationTasks 获取申请相关任务
// @Summary 获取申请相关任务
// @Description 获取指定申请的所有相关任务
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Success 200 {object} StandardResponse{data=service.GetApplicationTasksResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/loans/applications/{id}/tasks [get]
func (h *OALoanHandler) GetApplicationTasks(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的申请ID", err.Error()))
		return
	}

	req := &service.GetApplicationTasksRequest{
		ApplicationID: uint(id),
	}

	response, err := h.loanService.GetApplicationTasks(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "申请不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "申请不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取任务列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// GetAdvancedStatistics 获取高级统计数据
// @Summary 获取高级统计数据
// @Description 获取贷款业务的高级统计分析数据
// @Tags OA贷款管理
// @Accept json
// @Produce json
// @Param period query string false "统计周期 (day/week/month/quarter/year)" default(month)
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param include_trends query bool false "是否包含趋势分析" default(true)
// @Success 200 {object} StandardResponse{data=service.GetAdvancedStatisticsResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/loans/advanced-statistics [get]
func (h *OALoanHandler) GetAdvancedStatistics(c *gin.Context) {
	period := c.DefaultQuery("period", "month")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	includeTrendsStr := c.DefaultQuery("include_trends", "true")

	includeTrends, err := strconv.ParseBool(includeTrendsStr)
	if err != nil {
		includeTrends = true
	}

	req := &service.GetAdvancedStatisticsRequest{
		Period:        period,
		StartDate:     startDate,
		EndDate:       endDate,
		IncludeTrends: includeTrends,
	}

	response, err := h.loanService.GetAdvancedStatistics(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取统计数据失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}
