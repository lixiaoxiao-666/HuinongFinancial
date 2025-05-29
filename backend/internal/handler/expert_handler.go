package handler

import (
	"net/http"
	"strconv"

	"huinong-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// ExpertHandler 专家处理器
type ExpertHandler struct {
	contentService service.ContentService
}

// NewExpertHandler 创建专家处理器
func NewExpertHandler(contentService service.ContentService) *ExpertHandler {
	return &ExpertHandler{
		contentService: contentService,
	}
}

// CreateExpert 创建专家
// @Summary 创建专家
// @Description 管理员创建专家信息
// @Tags 专家管理
// @Accept json
// @Produce json
// @Param request body service.CreateExpertRequest true "专家信息"
// @Success 200 {object} StandardResponse{data=model.Expert}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/content/experts [post]
func (h *ExpertHandler) CreateExpert(c *gin.Context) {
	var req service.CreateExpertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	// 从上下文获取操作员ID
	_, exists := c.Get("oaUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "未登录", "OA用户认证信息缺失"))
		return
	}

	expert, err := h.contentService.CreateExpert(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "创建专家失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("专家创建成功", expert))
}

// GetExpert 获取专家详情
// @Summary 获取专家详情
// @Description 获取指定专家的详细信息
// @Tags 专家管理
// @Accept json
// @Produce json
// @Param id path int true "专家ID"
// @Success 200 {object} StandardResponse{data=service.ExpertDetailResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/content/experts/{id} [get]
func (h *ExpertHandler) GetExpert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的专家ID", err.Error()))
		return
	}

	response, err := h.contentService.GetExpert(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "专家不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "专家不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取专家详情失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// ListExperts 获取专家列表
// @Summary 获取专家列表
// @Description 获取专家列表，支持专业领域、地区等筛选
// @Tags 专家管理
// @Accept json
// @Produce json
// @Param specialty query string false "专业领域"
// @Param province query string false "省份"
// @Param city query string false "城市"
// @Param status query string false "状态"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Success 200 {object} StandardResponse{data=service.ListExpertsResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/content/experts [get]
func (h *ExpertHandler) ListExperts(c *gin.Context) {
	specialty := c.Query("specialty")
	province := c.Query("province")
	city := c.Query("city")
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

	req := &service.ListExpertsRequest{
		Specialty: specialty,
		Province:  province,
		City:      city,
		Status:    status,
		Page:      page,
		Limit:     limit,
	}

	response, err := h.contentService.ListExperts(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取专家列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// UpdateExpert 更新专家信息
// @Summary 更新专家信息
// @Description 管理员更新专家信息
// @Tags 专家管理
// @Accept json
// @Produce json
// @Param id path int true "专家ID"
// @Param request body service.UpdateExpertRequest true "更新信息"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/content/experts/{id} [put]
func (h *ExpertHandler) UpdateExpert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的专家ID", err.Error()))
		return
	}

	var req service.UpdateExpertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	err = h.contentService.UpdateExpert(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "专家不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "专家不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "更新专家信息失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("专家信息更新成功", nil))
}

// DeleteExpert 删除专家
// @Summary 删除专家
// @Description 管理员删除专家
// @Tags 专家管理
// @Accept json
// @Produce json
// @Param id path int true "专家ID"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/content/experts/{id} [delete]
func (h *ExpertHandler) DeleteExpert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的专家ID", err.Error()))
		return
	}

	err = h.contentService.DeleteExpert(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "专家不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "专家不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "删除专家失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("专家删除成功", nil))
}

// SubmitConsultation 提交咨询
// @Summary 提交咨询
// @Description 用户向专家提交咨询问题
// @Tags 专家咨询
// @Accept json
// @Produce json
// @Param request body service.SubmitConsultationRequest true "咨询信息"
// @Success 200 {object} StandardResponse{data=service.SubmitConsultationResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/consultations [post]
func (h *ExpertHandler) SubmitConsultation(c *gin.Context) {
	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	var req service.SubmitConsultationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	response, err := h.contentService.SubmitConsultation(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "专家不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "专家不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "提交咨询失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("咨询提交成功", response))
}

// GetConsultations 获取我的咨询
// @Summary 获取我的咨询
// @Description 获取当前用户的咨询记录
// @Tags 专家咨询
// @Accept json
// @Produce json
// @Success 200 {object} StandardResponse{data=service.ConsultationsResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/consultations [get]
func (h *ExpertHandler) GetConsultations(c *gin.Context) {
	// 从上下文获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	userID, ok := userIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "用户ID格式错误", "userID type assertion failed"))
		return
	}

	response, err := h.contentService.GetConsultations(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取咨询记录失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}
