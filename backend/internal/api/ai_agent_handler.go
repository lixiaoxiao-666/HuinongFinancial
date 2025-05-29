package api

import (
	"backend/internal/service"
	"backend/pkg"

	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AIAgentHandler AI智能体处理器
type AIAgentHandler struct {
	aiAgentService *service.AIAgentService
	log            *zap.Logger
}

// NewAIAgentHandler 创建AI智能体处理器
func NewAIAgentHandler(aiAgentService *service.AIAgentService, log *zap.Logger) *AIAgentHandler {
	return &AIAgentHandler{
		aiAgentService: aiAgentService,
		log:            log,
	}
}

// GetApplicationInfo 获取申请信息供Dify工作流调用
func (h *AIAgentHandler) GetApplicationInfo(c *gin.Context) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		pkg.BadRequest(c, "申请ID不能为空")
		return
	}

	h.log.Info("处理获取申请信息请求", zap.String("applicationId", applicationID))

	// 使用带日志记录的方法
	info, err := h.aiAgentService.GetApplicationInfoWithLog(applicationID, c.Request)
	if err != nil {
		h.log.Error("获取申请信息失败", zap.Error(err), zap.String("applicationId", applicationID))
		pkg.InternalError(c, err.Error())
		return
	}

	pkg.Success(c, info)
}

// SubmitAIDecision 接收AI决策结果（修改为查询参数方式）
func (h *AIAgentHandler) SubmitAIDecision(c *gin.Context) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		h.log.Warn("Missing application_id parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "申请ID不能为空",
			"data":    nil,
		})
		return
	}

	// 从查询参数获取必需的AI决策数据
	decision := c.Query("decision")
	riskScoreStr := c.Query("risk_score")
	confidenceScoreStr := c.Query("confidence_score")
	approvedAmountStr := c.Query("approved_amount")
	approvedTermMonthsStr := c.Query("approved_term_months")
	suggestedInterestRate := c.Query("suggested_interest_rate")
	riskLevel := c.Query("risk_level")
	analysisSummary := c.Query("analysis_summary")

	// 验证必需参数
	if decision == "" || riskScoreStr == "" || confidenceScoreStr == "" ||
		approvedAmountStr == "" || approvedTermMonthsStr == "" {
		h.log.Warn("Missing required parameters",
			zap.String("decision", decision),
			zap.String("risk_score", riskScoreStr),
			zap.String("confidence_score", confidenceScoreStr))
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "缺少必需的AI决策参数",
			"data":    nil,
		})
		return
	}

	// 类型转换和验证
	riskScore, err := strconv.ParseFloat(riskScoreStr, 64)
	if err != nil || riskScore < 0 || riskScore > 1 {
		h.log.Warn("Invalid risk_score parameter", zap.String("risk_score", riskScoreStr))
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "风险分数格式错误，应为0-1之间的数值",
			"data":    nil,
		})
		return
	}

	confidenceScore, err := strconv.ParseFloat(confidenceScoreStr, 64)
	if err != nil || confidenceScore < 0 || confidenceScore > 1 {
		h.log.Warn("Invalid confidence_score parameter", zap.String("confidence_score", confidenceScoreStr))
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "置信度分数格式错误，应为0-1之间的数值",
			"data":    nil,
		})
		return
	}

	approvedAmount, err := strconv.ParseFloat(approvedAmountStr, 64)
	if err != nil || approvedAmount < 0 {
		h.log.Warn("Invalid approved_amount parameter", zap.String("approved_amount", approvedAmountStr))
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "批准金额格式错误",
			"data":    nil,
		})
		return
	}

	approvedTermMonths, err := strconv.Atoi(approvedTermMonthsStr)
	if err != nil || approvedTermMonths < 1 {
		h.log.Warn("Invalid approved_term_months parameter", zap.String("approved_term_months", approvedTermMonthsStr))
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "批准期限格式错误",
			"data":    nil,
		})
		return
	}

	// 验证决策枚举值
	validDecisions := []string{"AUTO_APPROVED", "AUTO_REJECTED", "REQUIRE_HUMAN_REVIEW"}
	if !contains(validDecisions, decision) {
		h.log.Warn("Invalid decision parameter", zap.String("decision", decision))
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "无效的决策类型",
			"data":    nil,
		})
		return
	}

	// 验证风险等级枚举值
	validRiskLevels := []string{"LOW", "MEDIUM", "HIGH"}
	if riskLevel != "" && !contains(validRiskLevels, riskLevel) {
		h.log.Warn("Invalid risk_level parameter", zap.String("risk_level", riskLevel))
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "无效的风险等级",
			"data":    nil,
		})
		return
	}

	// 获取可选参数
	detailedAnalysis := c.Query("detailed_analysis")
	recommendations := c.Query("recommendations") // 逗号分隔的字符串
	conditions := c.Query("conditions")           // 逗号分隔的字符串
	aiModelVersion := c.Query("ai_model_version")
	workflowID := c.Query("workflow_id")

	// 设置默认值
	if riskLevel == "" {
		if riskScore < 0.3 {
			riskLevel = "LOW"
		} else if riskScore < 0.7 {
			riskLevel = "MEDIUM"
		} else {
			riskLevel = "HIGH"
		}
	}

	if suggestedInterestRate == "" {
		suggestedInterestRate = "5.0%"
	}

	if analysisSummary == "" {
		analysisSummary = "AI智能分析完成"
	}

	if aiModelVersion == "" {
		aiModelVersion = "v1.0.0"
	}

	if workflowID == "" {
		workflowID = "dify_ai_workflow"
	}

	// 解析详细分析JSON（如果提供）
	var detailedAnalysisObj map[string]interface{}
	if detailedAnalysis != "" {
		// 清理可能包含的schema信息
		detailedAnalysis = cleanDetailedAnalysis(detailedAnalysis)
		if err := json.Unmarshal([]byte(detailedAnalysis), &detailedAnalysisObj); err != nil {
			h.log.Warn("Failed to parse detailed_analysis JSON", zap.Error(err))
			// 不返回错误，使用默认值
			detailedAnalysisObj = map[string]interface{}{
				"credit_analysis":    "信用分析",
				"financial_analysis": "财务分析",
				"risk_factors":       []string{"待评估"},
				"strengths":          []string{"待评估"},
			}
		}
	}

	// 解析推荐和条件（逗号分隔字符串转数组）
	var recommendationsList []string
	var conditionsList []string

	if recommendations != "" {
		recommendationsList = strings.Split(recommendations, ",")
		// 清理空格
		for i, rec := range recommendationsList {
			recommendationsList[i] = strings.TrimSpace(rec)
		}
	} else {
		recommendationsList = []string{"请关注后续通知"}
	}

	if conditions != "" {
		conditionsList = strings.Split(conditions, ",")
		// 清理空格
		for i, cond := range conditionsList {
			conditionsList[i] = strings.TrimSpace(cond)
		}
	} else {
		conditionsList = []string{"无特殊条件"}
	}

	// 构建AI决策请求对象
	request := &service.AIDecisionParams{
		ApplicationID:         applicationID,
		Decision:              decision,
		RiskScore:             riskScore,
		ConfidenceScore:       confidenceScore,
		ApprovedAmount:        approvedAmount,
		ApprovedTermMonths:    approvedTermMonths,
		SuggestedInterestRate: suggestedInterestRate,
		RiskLevel:             riskLevel,
		AnalysisSummary:       analysisSummary,
		DetailedAnalysis:      detailedAnalysisObj,
		Recommendations:       recommendationsList,
		Conditions:            conditionsList,
		AIModelVersion:        aiModelVersion,
		WorkflowID:            workflowID,
	}

	h.log.Info("Processing AI decision submission",
		zap.String("application_id", applicationID),
		zap.String("decision", decision),
		zap.Float64("risk_score", riskScore),
		zap.String("risk_level", riskLevel))

	// 设置request到context中，以便service层记录日志
	ctx := context.WithValue(c.Request.Context(), "request", c.Request)

	// 调用服务层处理
	result, err := h.aiAgentService.SubmitAIDecisionQuery(ctx, request)
	if err != nil {
		h.log.Error("Failed to submit AI decision", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "处理AI决策失败",
			"data":    nil,
		})
		return
	}

	h.log.Info("AI decision submitted successfully", zap.String("application_id", applicationID))

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "AI审批结果已成功处理",
		"data":    result,
	})
}

// 辅助函数：检查字符串是否在数组中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// 辅助函数：清理详细分析JSON中的schema信息
func cleanDetailedAnalysis(detailedAnalysis string) string {
	// 尝试解析JSON
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(detailedAnalysis), &data); err != nil {
		return detailedAnalysis
	}

	// 移除schema相关字段
	delete(data, "type")
	delete(data, "properties")
	delete(data, "required")

	// 重新序列化
	cleaned, err := json.Marshal(data)
	if err != nil {
		return detailedAnalysis
	}

	return string(cleaned)
}

// TriggerWorkflow 触发AI审批工作流
func (h *AIAgentHandler) TriggerWorkflow(c *gin.Context) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		pkg.BadRequest(c, "申请ID不能为空")
		return
	}

	var request struct {
		WorkflowType string `json:"workflow_type" binding:"required"`
		Priority     string `json:"priority"`
		CallbackURL  string `json:"callback_url"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.Error("解析工作流触发请求失败", zap.Error(err))
		pkg.BadRequest(c, "请求参数格式错误")
		return
	}

	// 设置默认优先级
	if request.Priority == "" {
		request.Priority = "NORMAL"
	}

	h.log.Info("处理工作流触发请求",
		zap.String("applicationId", applicationID),
		zap.String("workflowType", request.WorkflowType),
		zap.String("priority", request.Priority),
	)

	execution, err := h.aiAgentService.TriggerWorkflow(
		applicationID,
		request.WorkflowType,
		request.Priority,
		request.CallbackURL,
	)
	if err != nil {
		h.log.Error("触发工作流失败", zap.Error(err), zap.String("applicationId", applicationID))
		pkg.InternalError(c, err.Error())
		return
	}

	pkg.Success(c, map[string]interface{}{
		"workflow_execution_id":     execution.ExecutionID,
		"estimated_completion_time": execution.EstimatedCompletion,
		"status":                    execution.Status,
	})
}

// GetAIModelConfig 获取AI模型配置信息
func (h *AIAgentHandler) GetAIModelConfig(c *gin.Context) {
	h.log.Info("处理获取AI模型配置请求")

	// 使用带日志记录的方法
	config, err := h.aiAgentService.GetAIModelConfigWithLog(c.Request)
	if err != nil {
		h.log.Error("获取AI模型配置失败", zap.Error(err))
		pkg.InternalError(c, err.Error())
		return
	}

	pkg.Success(c, config)
}

// GetExternalData 获取外部数据
func (h *AIAgentHandler) GetExternalData(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		pkg.BadRequest(c, "用户ID不能为空")
		return
	}

	dataTypes := c.Query("data_types")
	if dataTypes == "" {
		dataTypes = "credit_report,bank_flow,blacklist_check" // 默认查询类型
	}

	h.log.Info("处理获取外部数据请求",
		zap.String("userId", userID),
		zap.String("dataTypes", dataTypes),
	)

	// 使用带日志记录的方法
	data, err := h.aiAgentService.GetExternalDataWithLog(userID, dataTypes, c.Request)
	if err != nil {
		h.log.Error("获取外部数据失败", zap.Error(err), zap.String("userId", userID))
		pkg.InternalError(c, err.Error())
		return
	}

	pkg.Success(c, data)
}

// UpdateApplicationStatus 更新申请状态
func (h *AIAgentHandler) UpdateApplicationStatus(c *gin.Context) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		pkg.BadRequest(c, "申请ID不能为空")
		return
	}

	var request struct {
		Status   string                 `json:"status" binding:"required"`
		Operator string                 `json:"operator" binding:"required"`
		Remarks  string                 `json:"remarks"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.Error("解析状态更新请求失败", zap.Error(err))
		pkg.BadRequest(c, "请求参数格式错误")
		return
	}

	h.log.Info("处理申请状态更新",
		zap.String("applicationId", applicationID),
		zap.String("status", request.Status),
		zap.String("operator", request.Operator),
	)

	if err := h.aiAgentService.UpdateApplicationStatus(
		applicationID,
		request.Status,
		request.Operator,
		request.Remarks,
		request.Metadata,
	); err != nil {
		h.log.Error("更新申请状态失败", zap.Error(err), zap.String("applicationId", applicationID))
		pkg.InternalError(c, err.Error())
		return
	}

	pkg.Success(c, nil)
}

// RegisterAIAgentRoutes 注册AI智能体路由
func RegisterAIAgentRoutes(router *gin.RouterGroup, handler *AIAgentHandler, authMiddleware gin.HandlerFunc) {
	// AI智能体路由组
	aiAgent := router.Group("/ai-agent")

	// 使用AI Agent Token认证
	aiAgent.Use(authMiddleware)

	{
		// 获取申请信息供Dify调用
		aiAgent.GET("/applications/:application_id/info", handler.GetApplicationInfo)

		// 提交AI决策结果
		aiAgent.POST("/applications/:application_id/ai-decision", handler.SubmitAIDecision)

		// 获取AI模型配置
		aiAgent.GET("/config/models", handler.GetAIModelConfig)

		// 获取外部数据
		aiAgent.GET("/external-data", handler.GetExternalData)

		// 更新申请状态
		aiAgent.PUT("/applications/:application_id/status", handler.UpdateApplicationStatus)
	}

	// 系统内部调用的路由（需要系统Token）
	system := router.Group("/ai-agent")
	system.Use(authMiddleware) // 这里应该使用系统Token认证，暂时复用
	{
		// 触发工作流
		system.POST("/applications/:application_id/trigger-workflow", handler.TriggerWorkflow)
	}
}
