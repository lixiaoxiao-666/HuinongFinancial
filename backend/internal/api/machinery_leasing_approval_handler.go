package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"backend/internal/service"
	"backend/pkg"
)

type MachineryLeasingApprovalHandler struct {
	leasingApprovalService *service.MachineryLeasingApprovalService
	log                    *zap.Logger
}

func NewMachineryLeasingApprovalHandler(leasingApprovalService *service.MachineryLeasingApprovalService, log *zap.Logger) *MachineryLeasingApprovalHandler {
	return &MachineryLeasingApprovalHandler{
		leasingApprovalService: leasingApprovalService,
		log:                    log,
	}
}

// GetLeasingApplicationInfo 获取农机租赁申请信息
// @Summary 获取农机租赁申请详细信息
// @Description Dify工作流调用此接口获取农机租赁申请的完整信息，包括承租方、出租方、农机、申请详情等
// @Tags AI智能体-农机租赁
// @Accept json
// @Produce json
// @Param application_id path string true "申请ID"
// @Success 200 {object} service.MachineryLeasingApplicationInfo "农机租赁申请信息"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "申请不存在"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /ai-agent/machinery-leasing/applications/{application_id} [get]
func (h *MachineryLeasingApprovalHandler) GetLeasingApplicationInfo(c *gin.Context) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		h.log.Warn("Missing application_id parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "申请ID不能为空",
		})
		return
	}

	h.log.Info("Getting leasing application info",
		zap.String("applicationId", applicationID),
		zap.String("userAgent", c.GetHeader("User-Agent")),
		zap.String("clientIP", c.ClientIP()))

	info, err := h.leasingApprovalService.GetLeasingApplicationInfoWithLog(applicationID, c.Request)
	if err != nil {
		h.log.Error("Failed to get leasing application info",
			zap.String("applicationId", applicationID),
			zap.Error(err))

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "获取农机租赁申请信息失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, info)
}

// SubmitLeasingAIDecision 提交农机租赁AI审批决策
// @Summary 提交农机租赁AI审批决策
// @Description Dify工作流调用此接口提交对农机租赁申请的AI分析结果和审批决策
// @Tags AI智能体-农机租赁
// @Accept json
// @Produce json
// @Param application_id path string true "申请ID"
// @Param decision body service.LeasingAIDecisionRequest true "AI决策请求"
// @Success 200 {object} map[string]interface{} "提交成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "申请不存在"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /ai-agent/machinery-leasing/applications/{application_id}/decision [post]
func (h *MachineryLeasingApprovalHandler) SubmitLeasingAIDecision(c *gin.Context) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		h.log.Warn("Missing application_id parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "申请ID不能为空",
		})
		return
	}

	var request service.LeasingAIDecisionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "请求参数格式错误",
			"details": err.Error(),
		})
		return
	}

	// 验证必要的字段
	if request.AIDecision.Decision == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_required_field",
			"message": "AI决策结果不能为空",
		})
		return
	}

	h.log.Info("Submitting leasing AI decision",
		zap.String("applicationId", applicationID),
		zap.String("decision", request.AIDecision.Decision),
		zap.Float64("riskScore", request.AIAnalysis.RiskScore),
		zap.String("riskLevel", request.AIAnalysis.RiskLevel))

	err := h.leasingApprovalService.SubmitLeasingAIDecisionWithLog(applicationID, &request, c.Request)
	if err != nil {
		h.log.Error("Failed to submit leasing AI decision",
			zap.String("applicationId", applicationID),
			zap.Error(err))

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "提交AI决策失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "AI决策提交成功",
		"data": gin.H{
			"application_id": applicationID,
			"decision":       request.AIDecision.Decision,
			"processed_at":   request.ProcessingInfo.ProcessedAt,
		},
	})
}

// SubmitLeasingAIDecisionQuery 查询参数方式提交农机租赁AI决策
// @Summary 查询参数方式提交农机租赁AI决策
// @Description 通过URL查询参数的方式提交农机租赁AI审批决策，适用于简化的工作流调用
// @Tags AI智能体-农机租赁
// @Accept json
// @Produce json
// @Param application_id query string true "申请ID"
// @Param decision query string true "AI决策 (AUTO_APPROVE/AUTO_REJECT/REQUIRE_HUMAN_REVIEW/REQUIRE_DEPOSIT_ADJUSTMENT)"
// @Param risk_score query number true "风险评分 (0-1)"
// @Param confidence_score query number false "置信度评分 (0-1)"
// @Param suggested_deposit query number false "建议押金"
// @Param risk_level query string false "风险等级 (LOW/MEDIUM/HIGH)"
// @Param analysis_summary query string false "分析摘要"
// @Param ai_model_version query string false "AI模型版本"
// @Param workflow_id query string false "工作流ID"
// @Success 200 {object} service.LeasingAIDecisionResult "决策结果"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /ai-agent/machinery-leasing/decision-query [post]
func (h *MachineryLeasingApprovalHandler) SubmitLeasingAIDecisionQuery(c *gin.Context) {
	// 获取必要参数
	applicationID := c.Query("application_id")
	decision := c.Query("decision")
	riskScoreStr := c.Query("risk_score")

	if applicationID == "" || decision == "" || riskScoreStr == "" {
		h.log.Warn("Missing required parameters",
			zap.String("applicationId", applicationID),
			zap.String("decision", decision),
			zap.String("riskScore", riskScoreStr))

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_required_parameters",
			"message": "申请ID、决策和风险评分是必要参数",
		})
		return
	}

	// 解析风险评分
	riskScore, err := strconv.ParseFloat(riskScoreStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_risk_score",
			"message": "风险评分格式错误，应为0-1之间的数字",
		})
		return
	}

	// 验证决策值
	validDecisions := []string{"AUTO_APPROVE", "AUTO_REJECT", "REQUIRE_HUMAN_REVIEW", "REQUIRE_DEPOSIT_ADJUSTMENT"}
	if !pkg.Contains(validDecisions, decision) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":        "invalid_decision",
			"message":      "无效的决策值",
			"valid_values": validDecisions,
		})
		return
	}

	// 获取可选参数
	confidenceScoreStr := c.DefaultQuery("confidence_score", "0.8")
	confidenceScore, _ := strconv.ParseFloat(confidenceScoreStr, 64)

	suggestedDepositStr := c.DefaultQuery("suggested_deposit", "0")
	suggestedDeposit, _ := strconv.ParseFloat(suggestedDepositStr, 64)

	riskLevel := c.DefaultQuery("risk_level", "MEDIUM")
	analysisSummary := c.DefaultQuery("analysis_summary", "AI农机租赁风险分析")
	aiModelVersion := c.DefaultQuery("ai_model_version", "v1.0.0")
	workflowID := c.DefaultQuery("workflow_id", "dify_leasing_ai_workflow")

	// 构建参数
	params := &service.LeasingAIDecisionParams{
		ApplicationID:       applicationID,
		Decision:            decision,
		RiskScore:           riskScore,
		ConfidenceScore:     confidenceScore,
		SuggestedDeposit:    suggestedDeposit,
		SuggestedConditions: []string{"请关注后续通知"},
		RiskLevel:           riskLevel,
		AnalysisSummary:     analysisSummary,
		DetailedAnalysis: map[string]interface{}{
			"analysis_summary":  analysisSummary,
			"suggested_deposit": suggestedDeposit,
			"conditions":        []string{"需要审核"},
			"confidence_score":  confidenceScore,
			"decision":          decision,
			"detailed_analysis": map[string]interface{}{
				"lessee_analysis":    "承租方分析",
				"lessor_analysis":    "出租方分析",
				"machinery_analysis": "农机状况分析",
				"risk_factors":       []string{"待评估"},
				"strengths":          []string{"待评估"},
			},
			"recommendations": []string{"建议审核"},
			"risk_level":      riskLevel,
			"risk_score":      riskScore,
		},
		Recommendations: []string{"请关注后续通知"},
		AIModelVersion:  aiModelVersion,
		WorkflowID:      workflowID,
	}

	h.log.Info("Submitting leasing AI decision via query",
		zap.String("applicationId", applicationID),
		zap.String("decision", decision),
		zap.Float64("riskScore", riskScore),
		zap.String("riskLevel", riskLevel))

	result, err := h.leasingApprovalService.SubmitLeasingAIDecisionQuery(c.Request.Context(), params)
	if err != nil {
		h.log.Error("Failed to submit leasing AI decision via query",
			zap.String("applicationId", applicationID),
			zap.Error(err))

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "提交AI决策失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateLeasingApplicationStatus 更新农机租赁申请状态
// @Summary 更新农机租赁申请状态
// @Description 人工审批员更新农机租赁申请的审批状态
// @Tags 农机租赁审批
// @Accept json
// @Produce json
// @Param application_id path string true "申请ID"
// @Param request body map[string]interface{} true "状态更新请求"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /machinery-leasing/applications/{application_id}/status [put]
func (h *MachineryLeasingApprovalHandler) UpdateLeasingApplicationStatus(c *gin.Context) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_parameter",
			"message": "申请ID不能为空",
		})
		return
	}

	var request struct {
		Status   string                 `json:"status" binding:"required"`
		Operator string                 `json:"operator"`
		Remarks  string                 `json:"remarks"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "请求参数格式错误",
			"details": err.Error(),
		})
		return
	}

	h.log.Info("Updating leasing application status",
		zap.String("applicationId", applicationID),
		zap.String("status", request.Status),
		zap.String("operator", request.Operator))

	err := h.leasingApprovalService.UpdateLeasingApplicationStatus(
		applicationID,
		request.Status,
		request.Operator,
		request.Remarks,
		request.Metadata,
	)

	if err != nil {
		h.log.Error("Failed to update leasing application status",
			zap.String("applicationId", applicationID),
			zap.Error(err))

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "更新申请状态失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "申请状态更新成功",
		"data": gin.H{
			"application_id": applicationID,
			"new_status":     request.Status,
			"operator":       request.Operator,
		},
	})
}

// RegisterMachineryLeasingApprovalRoutes 注册农机租赁审批路由
func RegisterMachineryLeasingApprovalRoutes(router *gin.RouterGroup, handler *MachineryLeasingApprovalHandler, authMiddleware gin.HandlerFunc) {
	// AI智能体路由组 (无需认证，由Dify工作流调用)
	aiAgentGroup := router.Group("/ai-agent/machinery-leasing")
	{
		aiAgentGroup.GET("/applications/:application_id", handler.GetLeasingApplicationInfo)
		aiAgentGroup.POST("/applications/:application_id/decision", handler.SubmitLeasingAIDecision)
		aiAgentGroup.POST("/decision-query", handler.SubmitLeasingAIDecisionQuery)
	}

	// 人工审批路由组 (需要认证)
	approvalGroup := router.Group("/machinery-leasing", authMiddleware)
	{
		approvalGroup.PUT("/applications/:application_id/status", handler.UpdateLeasingApplicationStatus)
	}
}
