package api

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"backend/internal/service"
	"backend/pkg"
)

// AIAgentHandler AI智能体处理器（统一版本）
type AIAgentHandler struct {
	aiAgentService   *service.AIAgentService
	unifiedProcessor *service.UnifiedApplicationProcessor
	log              *zap.Logger
}

// NewAIAgentHandler 创建AI智能体处理器
func NewAIAgentHandler(
	aiAgentService *service.AIAgentService,
	unifiedProcessor *service.UnifiedApplicationProcessor,
	log *zap.Logger,
) *AIAgentHandler {
	return &AIAgentHandler{
		aiAgentService:   aiAgentService,
		unifiedProcessor: unifiedProcessor,
		log:              log,
	}
}

// RegisterAIAgentRoutes 注册AI智能体路由（统一版本）
func RegisterAIAgentRoutes(group *gin.RouterGroup, handler *AIAgentHandler, aiAgentAuthMiddleware gin.HandlerFunc) {
	// 使用统一处理架构的路由
	aiGroup := group.Group("/ai-agent")
	aiGroup.Use(aiAgentAuthMiddleware)
	{
		// 统一申请信息获取接口
		aiGroup.GET("/applications/:application_id/info", handler.GetApplicationInfoUnified)

		// 统一外部数据获取接口（智能适配）
		aiGroup.GET("/external-data/:user_id", handler.GetExternalDataUnified)

		// 统一AI决策提交接口（智能路由）
		aiGroup.POST("/applications/:application_id/decisions", handler.SubmitAIDecisionUnified)

		// AI模型配置获取接口（动态适配）
		aiGroup.GET("/config/models", handler.GetAIModelConfigUnified)

		// AI操作日志查询接口（统一查询）
		aiGroup.GET("/logs", handler.GetAIOperationLogs)

		// 保留专用接口作为兼容性支持（标记为deprecated）
		aiGroup.GET("/machinery-leasing/applications/:application_id", handler.GetMachineryLeasingApplicationInfo)
	}
}

// GetApplicationInfoUnified 获取申请信息（统一处理）
// @Summary 获取申请信息（统一处理）
// @Description 智能识别申请类型并返回对应的完整申请信息，支持贷款申请和农机租赁申请
// @Tags 统一AI智能体
// @Param application_id path string true "申请ID"
// @Success 200 {object} CommonResponse{data=UnifiedApplicationInfo}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/ai-agent/applications/{application_id}/info [get]
func (h *AIAgentHandler) GetApplicationInfoUnified(c *gin.Context) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		h.log.Warn("申请ID参数为空")
		pkg.BadRequestWithMessage(c, "申请ID不能为空")
		return
	}

	h.log.Info("开始获取统一申请信息",
		zap.String("application_id", applicationID),
		zap.String("client_ip", c.ClientIP()))

	// 使用统一处理器获取申请信息
	applicationInfo, err := h.unifiedProcessor.GetApplicationInfoUnified(applicationID)
	if err != nil {
		h.log.Error("获取统一申请信息失败",
			zap.String("application_id", applicationID),
			zap.Error(err))

		if strings.Contains(err.Error(), "不存在") {
			pkg.NotFoundWithMessage(c, "申请不存在")
			return
		}
		if strings.Contains(err.Error(), "置信度不足") {
			pkg.BadRequestWithMessage(c, "无法识别申请类型: "+err.Error())
			return
		}

		pkg.InternalErrorWithMessage(c, "获取申请信息失败")
		return
	}

	h.log.Info("统一申请信息获取成功",
		zap.String("application_id", applicationID),
		zap.String("application_type", applicationInfo.ApplicationType))

	pkg.Success(c, applicationInfo)
}

// GetExternalDataUnified 获取外部数据（智能适配）
// @Summary 获取外部数据（智能适配）
// @Description 根据用户类型和申请上下文智能获取相关外部数据
// @Tags 统一AI智能体
// @Param user_id path string true "用户ID"
// @Param data_types query string true "数据类型，逗号分隔"
// @Param application_id query string false "申请ID，用于上下文识别"
// @Success 200 {object} CommonResponse{data=UnifiedExternalDataResponse}
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/ai-agent/external-data/{user_id} [get]
func (h *AIAgentHandler) GetExternalDataUnified(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		h.log.Warn("用户ID参数为空")
		pkg.BadRequestWithMessage(c, "用户ID不能为空")
		return
	}

	dataTypesStr := c.Query("data_types")
	if dataTypesStr == "" {
		h.log.Warn("数据类型参数为空", zap.String("user_id", userID))
		pkg.BadRequestWithMessage(c, "数据类型不能为空")
		return
	}

	applicationID := c.Query("application_id") // 可选参数，用于上下文识别

	// 解析数据类型
	dataTypes := strings.Split(dataTypesStr, ",")
	for i, dataType := range dataTypes {
		dataTypes[i] = strings.TrimSpace(dataType)
	}

	h.log.Info("开始获取统一外部数据",
		zap.String("user_id", userID),
		zap.Strings("data_types", dataTypes),
		zap.String("application_id", applicationID))

	// 使用统一处理器获取外部数据
	externalData, err := h.unifiedProcessor.GetExternalDataUnified(userID, dataTypes, applicationID)
	if err != nil {
		h.log.Error("获取统一外部数据失败",
			zap.String("user_id", userID),
			zap.Error(err))
		pkg.InternalErrorWithMessage(c, "获取外部数据失败")
		return
	}

	h.log.Info("统一外部数据获取成功",
		zap.String("user_id", userID),
		zap.String("application_type", externalData.ApplicationType),
		zap.Strings("filtered_data_types", externalData.DataTypes))

	pkg.Success(c, externalData)
}

// SubmitAIDecisionUnified 提交AI决策（智能路由）
// @Summary 提交AI决策（智能路由）
// @Description 接收LLM分析后的决策结果，系统自动识别申请类型并路由到对应的业务处理逻辑
// @Tags 统一AI智能体
// @Param application_id path string true "申请ID"
// @Param decision query string true "AI决策结果"
// @Param risk_score query number true "风险分数(0-1)"
// @Param risk_level query string true "风险等级"
// @Param confidence_score query number true "置信度(0-1)"
// @Param analysis_summary query string true "分析摘要"
// @Param approved_amount query number false "批准金额"
// @Param approved_term_months query integer false "批准期限（月）"
// @Param suggested_interest_rate query string false "建议利率"
// @Param suggested_deposit query number false "建议押金"
// @Param detailed_analysis query string false "详细分析JSON字符串"
// @Param recommendations query string false "建议列表，逗号分隔"
// @Param conditions query string false "条件列表，逗号分隔"
// @Param ai_model_version query string false "AI模型版本"
// @Param workflow_id query string false "工作流ID"
// @Success 200 {object} CommonResponse{data=UnifiedDecisionResponse}
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/ai-agent/applications/{application_id}/decisions [post]
func (h *AIAgentHandler) SubmitAIDecisionUnified(c *gin.Context) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		h.log.Warn("申请ID参数为空")
		pkg.BadRequestWithMessage(c, "申请ID不能为空")
		return
	}

	// 解析必需参数
	decision := c.Query("decision")
	riskScoreStr := c.Query("risk_score")
	riskLevel := c.Query("risk_level")
	confidenceScoreStr := c.Query("confidence_score")
	analysisSummary := c.Query("analysis_summary")

	if decision == "" || riskScoreStr == "" || riskLevel == "" || confidenceScoreStr == "" || analysisSummary == "" {
		h.log.Warn("必需参数缺失",
			zap.String("application_id", applicationID),
			zap.String("decision", decision),
			zap.String("risk_score", riskScoreStr),
			zap.String("risk_level", riskLevel))
		pkg.BadRequestWithMessage(c, "缺少必需参数: decision, risk_score, risk_level, confidence_score, analysis_summary")
		return
	}

	// 转换数值参数
	riskScore, err := strconv.ParseFloat(riskScoreStr, 64)
	if err != nil || riskScore < 0 || riskScore > 1 {
		pkg.BadRequestWithMessage(c, "风险分数必须是0-1之间的数字")
		return
	}

	confidenceScore, err := strconv.ParseFloat(confidenceScoreStr, 64)
	if err != nil || confidenceScore < 0 || confidenceScore > 1 {
		pkg.BadRequestWithMessage(c, "置信度必须是0-1之间的数字")
		return
	}

	// 构建统一决策参数
	params := &service.UnifiedDecisionParams{
		ApplicationID:   applicationID,
		Decision:        decision,
		RiskScore:       riskScore,
		RiskLevel:       riskLevel,
		ConfidenceScore: confidenceScore,
		AnalysisSummary: analysisSummary,
		AIModelVersion:  c.Query("ai_model_version"),
		WorkflowID:      c.Query("workflow_id"),
	}

	// 解析可选参数
	if approvedAmountStr := c.Query("approved_amount"); approvedAmountStr != "" {
		if approvedAmount, err := strconv.ParseFloat(approvedAmountStr, 64); err == nil && approvedAmount >= 0 {
			params.ApprovedAmount = &approvedAmount
		}
	}

	if approvedTermMonthsStr := c.Query("approved_term_months"); approvedTermMonthsStr != "" {
		if approvedTermMonths, err := strconv.Atoi(approvedTermMonthsStr); err == nil && approvedTermMonths > 0 {
			params.ApprovedTermMonths = &approvedTermMonths
		}
	}

	if suggestedInterestRate := c.Query("suggested_interest_rate"); suggestedInterestRate != "" {
		params.SuggestedInterestRate = &suggestedInterestRate
	}

	if suggestedDepositStr := c.Query("suggested_deposit"); suggestedDepositStr != "" {
		if suggestedDeposit, err := strconv.ParseFloat(suggestedDepositStr, 64); err == nil && suggestedDeposit >= 0 {
			params.SuggestedDeposit = &suggestedDeposit
		}
	}

	// 解析详细分析JSON
	if detailedAnalysisStr := c.Query("detailed_analysis"); detailedAnalysisStr != "" {
		params.DetailedAnalysis = make(map[string]interface{})
		// 简单处理，实际使用中可能需要JSON解析
		params.DetailedAnalysis["raw"] = detailedAnalysisStr
	}

	// 解析建议和条件列表
	if recommendationsStr := c.Query("recommendations"); recommendationsStr != "" {
		params.Recommendations = strings.Split(recommendationsStr, ",")
		for i, rec := range params.Recommendations {
			params.Recommendations[i] = strings.TrimSpace(rec)
		}
	}

	if conditionsStr := c.Query("conditions"); conditionsStr != "" {
		params.Conditions = strings.Split(conditionsStr, ",")
		for i, cond := range params.Conditions {
			params.Conditions[i] = strings.TrimSpace(cond)
		}
	}

	h.log.Info("开始提交统一AI决策",
		zap.String("application_id", applicationID),
		zap.String("decision", decision),
		zap.Float64("risk_score", riskScore))

	// 使用统一处理器提交决策
	decisionResponse, err := h.unifiedProcessor.SubmitAIDecisionUnified(params)
	if err != nil {
		h.log.Error("提交统一AI决策失败",
			zap.String("application_id", applicationID),
			zap.Error(err))

		if strings.Contains(err.Error(), "不支持决策") {
			pkg.BadRequestWithMessage(c, err.Error())
			return
		}
		if strings.Contains(err.Error(), "置信度不足") {
			pkg.BadRequestWithMessage(c, "申请类型识别失败: "+err.Error())
			return
		}

		pkg.InternalErrorWithMessage(c, "提交AI决策失败")
		return
	}

	h.log.Info("统一AI决策提交成功",
		zap.String("application_id", applicationID),
		zap.String("application_type", decisionResponse.ApplicationType),
		zap.String("new_status", decisionResponse.NewStatus))

	pkg.Success(c, decisionResponse)
}

// GetAIModelConfigUnified 获取AI模型配置（动态适配）
// @Summary 获取AI模型配置（动态适配）
// @Description 获取当前可用的AI模型配置，根据申请类型动态调整阈值和规则
// @Tags 统一AI智能体
// @Param application_type query string false "申请类型"
// @Success 200 {object} CommonResponse{data=AIModelConfigResponse}
// @Router /api/v1/ai-agent/config/models [get]
func (h *AIAgentHandler) GetAIModelConfigUnified(c *gin.Context) {
	applicationType := c.Query("application_type")

	h.log.Info("开始获取AI模型配置",
		zap.String("application_type", applicationType))

	// 使用原有的AI服务获取模型配置（保持兼容性）
	modelConfig, err := h.aiAgentService.GetAIModelConfigUnified(c.Request.Context())
	if err != nil {
		h.log.Error("获取AI模型配置失败", zap.Error(err))
		pkg.InternalErrorWithMessage(c, "获取AI模型配置失败")
		return
	}

	// 根据申请类型过滤配置（如果指定了类型）
	if applicationType != "" {
		h.log.Info("根据申请类型过滤模型配置", zap.String("type", applicationType))
		// 可以在这里添加类型特定的过滤逻辑
	}

	h.log.Info("AI模型配置获取成功",
		zap.String("application_type", applicationType))

	pkg.Success(c, modelConfig)
}

// GetAIOperationLogs 获取AI操作日志（统一查询）
// @Summary 获取AI操作日志（统一查询）
// @Description 查询AI操作的详细日志，支持多种申请类型的统一查询和过滤
// @Tags 统一AI智能体
// @Param application_id query string false "申请ID"
// @Param application_type query string false "申请类型"
// @Param operation_type query string false "操作类型"
// @Param page query integer false "页码"
// @Param limit query integer false "每页数量"
// @Success 200 {object} CommonResponse{data=AIOperationLogsResponse}
// @Router /api/v1/ai-agent/logs [get]
func (h *AIAgentHandler) GetAIOperationLogs(c *gin.Context) {
	// 解析查询参数
	applicationID := c.Query("application_id")
	applicationType := c.Query("application_type")
	operationType := c.Query("operation_type")

	page := pkg.GetIntParam(c, "page", 1)
	limit := pkg.GetIntParam(c, "limit", 20)

	if limit > 100 {
		limit = 100
	}

	h.log.Info("开始获取AI操作日志",
		zap.String("application_id", applicationID),
		zap.String("application_type", applicationType),
		zap.String("operation_type", operationType),
		zap.Int("page", page),
		zap.Int("limit", limit))

	// 使用原有的AI服务获取操作日志（保持兼容性）
	logs, total, err := h.aiAgentService.GetAIOperationLogs(c.Request.Context(), applicationID, applicationType, page, limit)
	if err != nil {
		h.log.Error("获取AI操作日志失败", zap.Error(err))
		pkg.InternalErrorWithMessage(c, "获取AI操作日志失败")
		return
	}

	// 构建响应
	var responseLogs []AIOperationLog
	for _, log := range logs {
		responseLogs = append(responseLogs, AIOperationLog{
			OperationID:     log.OperationID,
			ApplicationID:   log.ApplicationID,
			ApplicationType: log.ApplicationType,
			Operation:       log.Operation,
			Result:          log.Result,
			Details:         log.Details,
			Timestamp:       log.Timestamp,
		})
	}

	response := AIOperationLogsResponse{
		Logs: responseLogs,
		Pagination: Pagination{
			CurrentPage: page,
			TotalPages:  int((total + int64(limit) - 1) / int64(limit)),
			TotalCount:  int(total),
			Limit:       limit,
		},
	}

	h.log.Info("AI操作日志获取成功",
		zap.Int("logs_count", len(responseLogs)),
		zap.Int("total_count", int(total)))

	pkg.Success(c, response)
}

// GetMachineryLeasingApplicationInfo 获取农机租赁申请信息（专用接口，兼容性保留）
// @Summary 获取农机租赁申请信息（已废弃）
// @Description 专门用于农机租赁申请的信息获取，建议使用统一接口 /applications/{id}/info
// @Tags 农机租赁专用（已废弃）
// @Param application_id path string true "农机租赁申请ID"
// @Success 200 {object} CommonResponse{data=MachineryLeasingApplicationInfo}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/ai-agent/machinery-leasing/applications/{application_id} [get]
// @deprecated
func (h *AIAgentHandler) GetMachineryLeasingApplicationInfo(c *gin.Context) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		h.log.Warn("申请ID参数为空")
		pkg.BadRequestWithMessage(c, "申请ID不能为空")
		return
	}

	h.log.Warn("使用了已废弃的专用接口",
		zap.String("application_id", applicationID),
		zap.String("deprecated_endpoint", "/machinery-leasing/applications/:id"),
		zap.String("recommended_endpoint", "/applications/:id/info"))

	// 转发到统一接口
	h.GetApplicationInfoUnified(c)
}
