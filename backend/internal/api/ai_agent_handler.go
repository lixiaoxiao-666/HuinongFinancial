package api

import (
	"backend/internal/service"
	"backend/pkg"

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

	info, err := h.aiAgentService.GetApplicationInfo(applicationID)
	if err != nil {
		h.log.Error("获取申请信息失败", zap.Error(err), zap.String("applicationId", applicationID))
		pkg.InternalError(c, err.Error())
		return
	}

	pkg.Success(c, info)
}

// SubmitAIDecision 接收Dify工作流的AI决策结果
func (h *AIAgentHandler) SubmitAIDecision(c *gin.Context) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		pkg.BadRequest(c, "申请ID不能为空")
		return
	}

	var request service.AIDecisionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.Error("解析AI决策请求失败", zap.Error(err))
		pkg.BadRequest(c, "请求参数格式错误")
		return
	}

	h.log.Info("处理AI决策提交",
		zap.String("applicationId", applicationID),
		zap.String("decision", request.AIDecision.Decision),
		zap.Float64("riskScore", request.AIAnalysis.RiskScore),
	)

	if err := h.aiAgentService.SubmitAIDecision(applicationID, &request); err != nil {
		h.log.Error("处理AI决策失败", zap.Error(err), zap.String("applicationId", applicationID))
		pkg.InternalError(c, err.Error())
		return
	}

	// 确定下一步状态
	var newStatus, nextStep string
	switch request.AIDecision.Decision {
	case "AUTO_APPROVED":
		newStatus = "AI_APPROVED"
		nextStep = "AWAIT_FINAL_CONFIRMATION"
	case "REQUIRE_HUMAN_REVIEW":
		newStatus = "MANUAL_REVIEW_REQUIRED"
		nextStep = "AWAIT_HUMAN_REVIEW"
	case "AUTO_REJECTED":
		newStatus = "AI_REJECTED"
		nextStep = "PROCESS_COMPLETED"
	default:
		newStatus = "AI_PROCESSED"
		nextStep = "AWAIT_NEXT_ACTION"
	}

	pkg.Success(c, map[string]interface{}{
		"application_id": applicationID,
		"new_status":     newStatus,
		"next_step":      nextStep,
	})
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

	config, err := h.aiAgentService.GetAIModelConfig()
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

	data, err := h.aiAgentService.GetExternalData(userID, dataTypes)
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
