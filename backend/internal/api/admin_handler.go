package api

import (
	"backend/internal/data"
	"backend/internal/service"
	"backend/pkg"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AdminHandler OA管理员处理器
type AdminHandler struct {
	adminService *service.AdminService
	loanService  *service.LoanService
	log          *zap.Logger
}

// NewAdminHandler 创建OA管理员处理器
func NewAdminHandler(adminService *service.AdminService, loanService *service.LoanService, log *zap.Logger) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		loanService:  loanService,
		log:          log,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login OA用户登录
func (h *AdminHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ValidationError(c, err.Error())
		return
	}

	// 验证用户名和密码
	oaUser, token, err := h.adminService.Login(req.Username, req.Password)
	if err != nil {
		h.log.Error("OA用户登录失败", zap.String("username", req.Username), zap.Error(err))
		pkg.Unauthorized(c, "用户名或密码错误")
		return
	}

	pkg.Success(c, map[string]interface{}{
		"admin_user_id": oaUser.OAUserID,
		"username":      oaUser.Username,
		"role":          oaUser.Role,
		"display_name":  oaUser.DisplayName,
		"token":         token,
		"expires_in":    3600,
	})
}

// GetPendingApplicationsRequest 获取待审批申请列表请求
type GetPendingApplicationsRequest struct {
	StatusFilter  string `form:"status_filter"`
	ApplicantName string `form:"applicant_name"`
	ApplicationID string `form:"application_id"`
	Page          int    `form:"page,default=1"`
	Limit         int    `form:"limit,default=10"`
}

// GetPendingApplications 获取待审批贷款申请列表
func (h *AdminHandler) GetPendingApplications(c *gin.Context) {
	var req GetPendingApplicationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkg.ValidationError(c, err.Error())
		return
	}

	// 构建查询条件
	filters := map[string]interface{}{}

	if req.StatusFilter != "" {
		filters["status"] = req.StatusFilter
	}

	if req.ApplicationID != "" {
		filters["application_id"] = req.ApplicationID
	}

	if req.ApplicantName != "" {
		filters["applicant_name"] = req.ApplicantName
	}

	// 获取待审批申请列表
	applications, total, err := h.adminService.GetPendingApplications(filters, req.Page, req.Limit)
	if err != nil {
		h.log.Error("获取待审批申请列表失败", zap.Error(err))
		pkg.InternalError(c, "获取申请列表失败")
		return
	}

	// 构建响应数据
	var responseData []map[string]interface{}
	for _, app := range applications {
		// 解析申请人信息
		var applicantInfo map[string]interface{}
		if app.ApplicantSnapshot != nil {
			json.Unmarshal(app.ApplicantSnapshot, &applicantInfo)
		}

		applicantName := ""
		if realName, exists := applicantInfo["real_name"]; exists {
			applicantName = realName.(string)
		}

		// 计算等待时间
		waitingDays := int(time.Since(app.SubmittedAt).Hours() / 24)

		item := map[string]interface{}{
			"application_id": app.ApplicationID,
			"applicant_name": applicantName,
			"product_name":   h.getProductName(app.ProductID),
			"amount_applied": app.AmountApplied,
			"status":         app.Status,
			"priority":       h.calculatePriority(app),
			"submitted_at":   app.SubmittedAt,
			"waiting_days":   waitingDays,
			"ai_risk_score":  app.AIRiskScore,
			"ai_suggestion":  app.AISuggestion,
		}
		responseData = append(responseData, item)
	}

	// 构建分页响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Success",
		"data":    responseData,
		"total":   total,
		"page":    req.Page,
		"limit":   req.Limit,
	})
}

// GetApplicationDetail 获取贷款申请详情 (管理员视角)
func (h *AdminHandler) GetApplicationDetail(c *gin.Context) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		pkg.BadRequest(c, "申请ID不能为空")
		return
	}

	// 获取申请详情
	application, err := h.adminService.GetApplicationDetail(applicationID)
	if err != nil {
		h.log.Error("获取申请详情失败", zap.String("application_id", applicationID), zap.Error(err))
		pkg.NotFound(c, "申请不存在")
		return
	}

	// 获取申请历史
	history, err := h.adminService.GetApplicationHistory(applicationID)
	if err != nil {
		h.log.Error("获取申请历史失败", zap.String("application_id", applicationID), zap.Error(err))
		// 不阻止返回，历史记录获取失败不影响主要信息
		history = []data.LoanApplicationHistory{}
	}

	// 获取相关文件
	files, err := h.adminService.GetApplicationFiles(applicationID)
	if err != nil {
		h.log.Error("获取申请文件失败", zap.String("application_id", applicationID), zap.Error(err))
		files = []data.UploadedFile{}
	}

	// 解析申请人信息
	var applicantInfo map[string]interface{}
	if application.ApplicantSnapshot != nil {
		json.Unmarshal(application.ApplicantSnapshot, &applicantInfo)
	}

	// 构建AI分析结果
	aiAnalysis := h.buildAIAnalysisResult(*application)

	// 构建响应数据
	responseData := map[string]interface{}{
		"application_id":       application.ApplicationID,
		"user_id":              application.UserID,
		"product_id":           application.ProductID,
		"product_name":         h.getProductName(application.ProductID),
		"amount_applied":       application.AmountApplied,
		"term_months_applied":  application.TermMonthsApplied,
		"purpose":              application.Purpose,
		"status":               application.Status,
		"submitted_at":         application.SubmittedAt,
		"updated_at":           application.UpdatedAt,
		"applicant_info":       applicantInfo,
		"ai_analysis":          aiAnalysis,
		"approved_amount":      application.ApprovedAmount,
		"approved_term_months": application.ApprovedTermMonths,
		"final_decision":       application.FinalDecision,
		"decision_reason":      application.DecisionReason,
		"processed_by":         application.ProcessedBy,
		"processed_at":         application.ProcessedAt,
		"history":              h.formatHistory(history),
		"uploaded_files":       h.formatFiles(files),
	}

	pkg.Success(c, responseData)
}

// ReviewApplicationRequest 审批申请请求
type ReviewApplicationRequest struct {
	Decision            string   `json:"decision" binding:"required,oneof=approved rejected request_more_info"`
	ApprovedAmount      *float64 `json:"approved_amount"`
	ApprovedTermMonths  *int     `json:"approved_term_months"`
	Comments            string   `json:"comments"`
	RequiredInfoDetails *string  `json:"required_info_details"`
}

// ReviewApplication 提交审批决策
func (h *AdminHandler) ReviewApplication(c *gin.Context) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		pkg.BadRequest(c, "申请ID不能为空")
		return
	}

	var req ReviewApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ValidationError(c, err.Error())
		return
	}

	// 获取当前审批员信息
	adminUserID, exists := c.Get("admin_user_id")
	if !exists {
		pkg.Unauthorized(c, "未获取到审批员信息")
		return
	}

	// 验证审批决策
	if req.Decision == "approved" && req.ApprovedAmount == nil {
		pkg.BadRequest(c, "批准申请时必须指定批准金额")
		return
	}

	// 执行审批
	err := h.adminService.ReviewApplication(applicationID, adminUserID.(string), req.Decision, &service.ReviewDetails{
		ApprovedAmount:      req.ApprovedAmount,
		ApprovedTermMonths:  req.ApprovedTermMonths,
		Comments:            req.Comments,
		RequiredInfoDetails: req.RequiredInfoDetails,
	})

	if err != nil {
		h.log.Error("提交审批决策失败",
			zap.String("application_id", applicationID),
			zap.String("decision", req.Decision),
			zap.Error(err))
		pkg.InternalError(c, "审批决策提交失败")
		return
	}

	h.log.Info("审批决策提交成功",
		zap.String("application_id", applicationID),
		zap.String("decision", req.Decision),
		zap.String("admin_user_id", adminUserID.(string)))

	pkg.SuccessWithMessage(c, "审批决策提交成功", nil)
}

// ToggleAIApprovalRequest AI审批开关请求
type ToggleAIApprovalRequest struct {
	Enabled bool `json:"enabled"`
}

// ToggleAIApproval 控制AI审批流程开关
func (h *AdminHandler) ToggleAIApproval(c *gin.Context) {
	var req ToggleAIApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ValidationError(c, err.Error())
		return
	}

	// 获取当前管理员信息
	adminUserID, exists := c.Get("admin_user_id")
	if !exists {
		pkg.Unauthorized(c, "未获取到管理员信息")
		return
	}

	// 更新AI审批开关状态
	err := h.adminService.ToggleAIApproval(req.Enabled, adminUserID.(string))
	if err != nil {
		h.log.Error("更新AI审批开关失败",
			zap.Bool("enabled", req.Enabled),
			zap.String("admin_user_id", adminUserID.(string)),
			zap.Error(err))
		pkg.InternalError(c, "AI审批流程状态更新失败")
		return
	}

	h.log.Info("AI审批开关状态已更新",
		zap.Bool("enabled", req.Enabled),
		zap.String("admin_user_id", adminUserID.(string)))

	pkg.SuccessWithMessage(c, "AI审批流程状态更新成功", map[string]interface{}{
		"enabled": req.Enabled,
	})
}

// GetSystemStats 获取系统统计信息
func (h *AdminHandler) GetSystemStats(c *gin.Context) {
	stats, err := h.adminService.GetSystemStats()
	if err != nil {
		h.log.Error("获取系统统计信息失败", zap.Error(err))
		pkg.InternalError(c, "获取统计信息失败")
		return
	}

	pkg.Success(c, stats)
}

// GetDashboard 获取OA首页/工作台信息
func (h *AdminHandler) GetDashboard(c *gin.Context) {
	// 获取当前管理员信息
	adminUserID, exists := c.Get("admin_user_id")
	if !exists {
		pkg.Unauthorized(c, "未获取到管理员信息")
		return
	}

	// 获取系统统计信息
	stats, err := h.adminService.GetSystemStats()
	if err != nil {
		h.log.Error("获取系统统计信息失败", zap.Error(err))
		pkg.InternalError(c, "获取统计信息失败")
		return
	}

	// 获取待办事项
	pendingTasks, err := h.adminService.GetPendingTasks(adminUserID.(string))
	if err != nil {
		h.log.Error("获取待办事项失败", zap.Error(err))
		// 不阻止返回，待办事项获取失败不影响主要信息
		pendingTasks = []map[string]interface{}{}
	}

	// 构建响应数据
	dashboardData := map[string]interface{}{
		"system_stats":  stats,
		"pending_tasks": pendingTasks,
		"quick_actions": []map[string]interface{}{
			{"name": "审批看板", "path": "/admin/approval", "icon": "approval"},
			{"name": "系统配置", "path": "/admin/system", "icon": "setting"},
			{"name": "操作日志", "path": "/admin/logs", "icon": "log"},
			{"name": "用户管理", "path": "/admin/users", "icon": "user"},
		},
		"recent_activities": h.getRecentActivities(),
	}

	pkg.Success(c, dashboardData)
}

// GetOAUsers 获取OA用户列表
func (h *AdminHandler) GetOAUsers(c *gin.Context) {
	page := pkg.GetIntParam(c, "page", 1)
	limit := pkg.GetIntParam(c, "limit", 10)
	role := c.Query("role")

	users, total, err := h.adminService.GetOAUsers(page, limit, role)
	if err != nil {
		h.log.Error("获取OA用户列表失败", zap.Error(err))
		pkg.InternalError(c, "获取用户列表失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Success",
		"data":    users,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// CreateOAUserRequest 创建OA用户请求
type CreateOAUserRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	Role        string `json:"role" binding:"required,oneof=审批员 系统管理员"`
	DisplayName string `json:"display_name" binding:"required"`
	Email       string `json:"email" binding:"email"`
}

// CreateOAUser 创建OA用户
func (h *AdminHandler) CreateOAUser(c *gin.Context) {
	var req CreateOAUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ValidationError(c, err.Error())
		return
	}

	// 获取当前管理员信息
	adminUserID, exists := c.Get("admin_user_id")
	if !exists {
		pkg.Unauthorized(c, "未获取到管理员信息")
		return
	}

	err := h.adminService.CreateOAUser(&req, adminUserID.(string))
	if err != nil {
		h.log.Error("创建OA用户失败",
			zap.String("username", req.Username),
			zap.String("creator", adminUserID.(string)),
			zap.Error(err))
		if err.Error() == "用户名已存在" {
			pkg.BadRequest(c, "用户名已存在")
			return
		}
		pkg.InternalError(c, "创建用户失败")
		return
	}

	h.log.Info("OA用户创建成功",
		zap.String("username", req.Username),
		zap.String("role", req.Role),
		zap.String("creator", adminUserID.(string)))

	pkg.SuccessWithMessage(c, "用户创建成功", nil)
}

// UpdateOAUserStatusRequest 更新OA用户状态请求
type UpdateOAUserStatusRequest struct {
	Status int8 `json:"status" binding:"oneof=0 1"` // 0:启用, 1:禁用
}

// UpdateOAUserStatus 更新OA用户状态
func (h *AdminHandler) UpdateOAUserStatus(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		pkg.BadRequest(c, "用户ID不能为空")
		return
	}

	var req UpdateOAUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ValidationError(c, err.Error())
		return
	}

	// 获取当前管理员信息
	adminUserID, exists := c.Get("admin_user_id")
	if !exists {
		pkg.Unauthorized(c, "未获取到管理员信息")
		return
	}

	err := h.adminService.UpdateOAUserStatus(userID, req.Status, adminUserID.(string))
	if err != nil {
		h.log.Error("更新OA用户状态失败",
			zap.String("user_id", userID),
			zap.Int8("status", req.Status),
			zap.String("operator", adminUserID.(string)),
			zap.Error(err))
		pkg.InternalError(c, "更新用户状态失败")
		return
	}

	h.log.Info("OA用户状态更新成功",
		zap.String("user_id", userID),
		zap.Int8("status", req.Status),
		zap.String("operator", adminUserID.(string)))

	statusText := "启用"
	if req.Status == 1 {
		statusText = "禁用"
	}

	pkg.SuccessWithMessage(c, "用户状态已"+statusText, nil)
}

// GetOperationLogsRequest 获取操作日志请求
type GetOperationLogsRequest struct {
	OperatorID string `form:"operator_id"`
	Action     string `form:"action"`
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=10"`
}

// GetOperationLogs 获取操作日志
func (h *AdminHandler) GetOperationLogs(c *gin.Context) {
	var req GetOperationLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkg.ValidationError(c, err.Error())
		return
	}

	logs, total, err := h.adminService.GetOperationLogs(&req)
	if err != nil {
		h.log.Error("获取操作日志失败", zap.Error(err))
		pkg.InternalError(c, "获取操作日志失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Success",
		"data":    logs,
		"total":   total,
		"page":    req.Page,
		"limit":   req.Limit,
	})
}

// GetSystemConfigurations 获取系统配置
func (h *AdminHandler) GetSystemConfigurations(c *gin.Context) {
	configs, err := h.adminService.GetSystemConfigurations()
	if err != nil {
		h.log.Error("获取系统配置失败", zap.Error(err))
		pkg.InternalError(c, "获取系统配置失败")
		return
	}

	pkg.Success(c, configs)
}

// UpdateSystemConfigurationRequest 更新系统配置请求
type UpdateSystemConfigurationRequest struct {
	ConfigValue string `json:"config_value" binding:"required"`
}

// UpdateSystemConfiguration 更新系统配置
func (h *AdminHandler) UpdateSystemConfiguration(c *gin.Context) {
	configKey := c.Param("config_key")
	if configKey == "" {
		pkg.BadRequest(c, "配置键不能为空")
		return
	}

	var req UpdateSystemConfigurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ValidationError(c, err.Error())
		return
	}

	// 获取当前管理员信息
	adminUserID, exists := c.Get("admin_user_id")
	if !exists {
		pkg.Unauthorized(c, "未获取到管理员信息")
		return
	}

	err := h.adminService.UpdateSystemConfiguration(configKey, req.ConfigValue, adminUserID.(string))
	if err != nil {
		h.log.Error("更新系统配置失败",
			zap.String("config_key", configKey),
			zap.String("operator", adminUserID.(string)),
			zap.Error(err))
		pkg.InternalError(c, "更新系统配置失败")
		return
	}

	h.log.Info("系统配置更新成功",
		zap.String("config_key", configKey),
		zap.String("config_value", req.ConfigValue),
		zap.String("operator", adminUserID.(string)))

	pkg.SuccessWithMessage(c, "系统配置更新成功", nil)
}

// 辅助方法

// getProductName 获取产品名称
func (h *AdminHandler) getProductName(productID string) string {
	product, err := h.loanService.GetProductByID(productID)
	if err != nil {
		return "未知产品"
	}
	return product.Name
}

// calculatePriority 计算申请优先级
func (h *AdminHandler) calculatePriority(app data.LoanApplication) string {
	// 根据等待时间和金额计算优先级
	waitingDays := int(time.Since(app.SubmittedAt).Hours() / 24)

	if waitingDays > 3 {
		return "HIGH"
	} else if app.AmountApplied > 100000 {
		return "MEDIUM"
	}
	return "NORMAL"
}

// buildAIAnalysisResult 构建AI分析结果
func (h *AdminHandler) buildAIAnalysisResult(app data.LoanApplication) map[string]interface{} {
	result := map[string]interface{}{
		"risk_score":      app.AIRiskScore,
		"suggestion":      app.AISuggestion,
		"analysis_status": "completed",
	}

	// 模拟详细的AI分析结果
	if app.AIRiskScore != nil {
		result["risk_level"] = h.getRiskLevel(*app.AIRiskScore)
		result["confidence"] = h.getConfidenceLevel(*app.AIRiskScore)

		// 模拟各项分析结果
		result["data_verification"] = map[string]interface{}{
			"identity_verified":     true,
			"phone_verified":        true,
			"income_consistency":    85,
			"document_authenticity": "pass",
		}

		result["risk_factors"] = []string{
			"申请人年龄适中",
			"收入来源稳定",
			"无不良信用记录",
		}

		if *app.AIRiskScore > 70 {
			result["risk_factors"] = append(result["risk_factors"].([]string), "历史还款记录良好")
		}
	}

	return result
}

// getRiskLevel 获取风险等级
func (h *AdminHandler) getRiskLevel(score int) string {
	if score >= 80 {
		return "LOW"
	} else if score >= 60 {
		return "MEDIUM"
	}
	return "HIGH"
}

// getConfidenceLevel 获取置信度等级
func (h *AdminHandler) getConfidenceLevel(score int) string {
	if score >= 85 {
		return "HIGH"
	} else if score >= 70 {
		return "MEDIUM"
	}
	return "LOW"
}

// formatHistory 格式化历史记录
func (h *AdminHandler) formatHistory(history []data.LoanApplicationHistory) []map[string]interface{} {
	var result []map[string]interface{}

	for _, h := range history {
		item := map[string]interface{}{
			"status_from":   h.StatusFrom,
			"status_to":     h.StatusTo,
			"operator_type": h.OperatorType,
			"operator_id":   h.OperatorID,
			"comments":      h.Comments,
			"occurred_at":   h.OccurredAt,
		}
		result = append(result, item)
	}

	return result
}

// formatFiles 格式化文件列表
func (h *AdminHandler) formatFiles(files []data.UploadedFile) []map[string]interface{} {
	var result []map[string]interface{}

	for _, f := range files {
		item := map[string]interface{}{
			"file_id":      f.FileID,
			"file_name":    f.FileName,
			"file_type":    f.FileType,
			"file_size":    f.FileSize,
			"storage_path": f.StoragePath,
			"purpose":      f.Purpose,
			"uploaded_at":  f.UploadedAt,
		}
		result = append(result, item)
	}

	return result
}

// getRecentActivities 获取最近活动记录
func (h *AdminHandler) getRecentActivities() []map[string]interface{} {
	// 模拟最近活动数据
	return []map[string]interface{}{
		{
			"time":     time.Now().Add(-30 * time.Minute),
			"action":   "审批申请",
			"target":   "贷款申请 APP20241210001",
			"operator": "审批员张三",
			"result":   "已批准",
		},
		{
			"time":     time.Now().Add(-2 * time.Hour),
			"action":   "AI审批配置",
			"target":   "启用自动审批",
			"operator": "管理员李四",
			"result":   "配置成功",
		},
		{
			"time":     time.Now().Add(-4 * time.Hour),
			"action":   "用户管理",
			"target":   "创建审批员账号",
			"operator": "管理员李四",
			"result":   "创建成功",
		},
	}
}

// RegisterAdminRoutes 注册管理员路由
func RegisterAdminRoutes(r *gin.RouterGroup, handler *AdminHandler, adminAuthMiddleware gin.HandlerFunc) {
	// OA用户登录（不需要认证）
	r.POST("/login", handler.Login)

	// 需要管理员认证的路由
	auth := r.Group("")
	auth.Use(adminAuthMiddleware)
	{
		// 审批相关
		auth.GET("/loans/applications/pending", handler.GetPendingApplications)
		auth.GET("/loans/applications/:application_id", handler.GetApplicationDetail)
		auth.POST("/loans/applications/:application_id/review", handler.ReviewApplication)

		// 系统管理
		auth.POST("/system/ai-approval/toggle", handler.ToggleAIApproval)
		auth.GET("/system/stats", handler.GetSystemStats)

		// 新增功能
		auth.GET("/dashboard", handler.GetDashboard)
		auth.GET("/users", handler.GetOAUsers)
		auth.POST("/users", handler.CreateOAUser)
		auth.PUT("/users/:user_id/status", handler.UpdateOAUserStatus)
		auth.GET("/logs", handler.GetOperationLogs)
		auth.GET("/configs", handler.GetSystemConfigurations)
		auth.PUT("/configs/:config_key", handler.UpdateSystemConfiguration)
	}
}
