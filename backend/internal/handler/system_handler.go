package handler

import (
	"net/http"

	"huinong-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SystemHandler 系统处理器
type SystemHandler struct {
	systemService service.SystemService
}

// NewSystemHandler 创建系统处理器
func NewSystemHandler(systemService service.SystemService) *SystemHandler {
	return &SystemHandler{
		systemService: systemService,
	}
}

// GetConfig 获取配置
// @Summary 获取配置
// @Description 获取指定配置项的值
// @Tags 系统管理
// @Accept json
// @Produce json
// @Param key query string true "配置键"
// @Success 200 {object} StandardResponse{data=string}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/system/config [get]
func (h *SystemHandler) GetConfig(c *gin.Context) {
	configKey := c.Query("key")
	if configKey == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "配置键不能为空", "config key is required"))
		return
	}

	value, err := h.systemService.GetConfig(c.Request.Context(), configKey)
	if err != nil {
		if err.Error() == "配置不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "配置不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取配置失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", value))
}

// SetConfig 设置配置
// @Summary 设置配置
// @Description 设置指定配置项的值
// @Tags 系统管理
// @Accept json
// @Produce json
// @Param request body map[string]string true "配置键值对"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/system/config [put]
func (h *SystemHandler) SetConfig(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	configKey, exists := req["key"]
	if !exists || configKey == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "配置键不能为空", "config key is required"))
		return
	}

	configValue, exists := req["value"]
	if !exists {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "配置值不能为空", "config value is required"))
		return
	}

	err := h.systemService.SetConfig(c.Request.Context(), configKey, configValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "设置配置失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("配置设置成功", nil))
}

// GetConfigs 获取配置组
// @Summary 获取配置组
// @Description 获取指定配置组的所有配置
// @Tags 系统管理
// @Accept json
// @Produce json
// @Param group query string false "配置组"
// @Success 200 {object} StandardResponse{data=map[string]string}
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/system/configs [get]
func (h *SystemHandler) GetConfigs(c *gin.Context) {
	configGroup := c.Query("group")

	configs, err := h.systemService.GetConfigs(c.Request.Context(), configGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取配置组失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", configs))
}

// GetPublicConfigs 获取公开配置
// @Summary 获取公开配置
// @Description 获取APP端需要的公开配置信息
// @Tags 系统管理
// @Accept json
// @Produce json
// @Success 200 {object} StandardResponse{data=map[string]string}
// @Failure 500 {object} ErrorResponse
// @Router /api/public/configs [get]
func (h *SystemHandler) GetPublicConfigs(c *gin.Context) {
	// 获取公开配置组
	configs, err := h.systemService.GetConfigs(c.Request.Context(), "public")
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取公开配置失败", err.Error()))
		return
	}

	// 确保只返回公开的配置项
	publicConfigs := map[string]string{
		"app_version":            "1.0.0",
		"customer_service_phone": "400-123-4567",
		"customer_service_hours": "9:00-18:00",
		"about_us":               "数字惠农是专业的农业金融服务平台",
		"privacy_policy_url":     "https://example.com/privacy",
		"terms_of_service_url":   "https://example.com/terms",
		"file_upload_max_size":   "10MB",
		"supported_file_types":   "jpg,jpeg,png,pdf,doc,docx",
	}

	// 合并从数据库获取的配置
	for k, v := range configs {
		publicConfigs[k] = v
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", publicConfigs))
}

// HealthCheck 健康检查
// @Summary 健康检查
// @Description 检查系统健康状态
// @Tags 系统管理
// @Accept json
// @Produce json
// @Success 200 {object} StandardResponse{data=service.HealthCheckResponse}
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/system/health [get]
func (h *SystemHandler) HealthCheck(c *gin.Context) {
	response, err := h.systemService.HealthCheck(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "健康检查失败", err.Error()))
		return
	}

	// 根据健康状态返回相应的HTTP状态码
	httpStatus := http.StatusOK
	if response.Status != "healthy" {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, NewSuccessResponse("健康检查完成", response))
}

// GetSystemStats 获取系统统计
// @Summary 获取系统统计
// @Description 获取系统总体统计数据
// @Tags 系统管理
// @Accept json
// @Produce json
// @Success 200 {object} StandardResponse{data=service.SystemStatsResponse}
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/system/statistics [get]
func (h *SystemHandler) GetSystemStats(c *gin.Context) {
	response, err := h.systemService.GetSystemStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取系统统计失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// GetSystemVersion 获取系统版本信息
// @Summary 获取系统版本信息
// @Description 获取API系统版本和基本信息
// @Tags 系统管理
// @Accept json
// @Produce json
// @Success 200 {object} StandardResponse
// @Router /api/public/version [get]
func (h *SystemHandler) GetSystemVersion(c *gin.Context) {
	version := map[string]interface{}{
		"version":     "1.0.0",
		"name":        "数字惠农API",
		"description": "数字惠农APP及OA后台管理系统API接口",
		"build_time":  "2024-01-01T00:00:00Z",
		"go_version":  "1.21",
		"environment": "production", // 可以从环境变量获取
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", version))
}
