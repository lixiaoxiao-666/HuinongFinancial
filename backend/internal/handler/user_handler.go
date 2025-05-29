package handler

import (
	"context"
	"huinong-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册请求参数"
// @Success 200 {object} SuccessResponse{data=service.RegisterResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	ctx := SetClientInfoToContext(c)
	resp, err := h.userService.Register(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "注册失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("注册成功", resp))
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录请求参数"
// @Success 200 {object} SuccessResponse{data=service.LoginResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	ctx := SetClientInfoToContext(c)
	resp, err := h.userService.Login(ctx, &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "登录失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("登录成功", resp))
}

// RefreshToken 刷新访问令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer refresh_token"
// @Success 200 {object} SuccessResponse{data=service.TokenResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/refresh-token [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "缺少刷新令牌", ""))
		return
	}

	// 去掉 Bearer 前缀
	if len(refreshToken) > 7 && refreshToken[:7] == "Bearer " {
		refreshToken = refreshToken[7:]
	}

	ctx := SetClientInfoToContext(c)
	resp, err := h.userService.RefreshToken(ctx, refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "刷新令牌失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("刷新令牌成功", resp))
}

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取当前用户的详细资料
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} SuccessResponse{data=service.UserProfileResponse}
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := GetUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", ""))
		return
	}

	ctx := SetClientInfoToContext(c)
	resp, err := h.userService.GetProfile(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取用户资料失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取用户资料成功", resp))
}

// UpdateProfile 更新用户资料
// @Summary 更新用户资料
// @Description 更新当前用户的资料信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body service.UpdateProfileRequest true "更新资料请求参数"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := GetUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", ""))
		return
	}

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	ctx := SetClientInfoToContext(c)
	err := h.userService.UpdateProfile(ctx, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "更新用户资料失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("更新用户资料成功", nil))
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户的登录密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body service.ChangePasswordRequest true "修改密码请求参数"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := GetUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", ""))
		return
	}

	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	ctx := SetClientInfoToContext(c)
	err := h.userService.ChangePassword(ctx, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "修改密码失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("修改密码成功", nil))
}

// Logout 用户退出登录
// @Summary 用户退出登录
// @Description 用户退出登录，清除会话
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	userID := GetUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", ""))
		return
	}

	// 从请求头获取session_id或从JWT中解析
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		// 可以从JWT token中解析session_id，这里简单处理
		sessionID = "current_session"
	}

	ctx := SetClientInfoToContext(c)
	err := h.userService.Logout(ctx, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "退出登录失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("退出登录成功", nil))
}

// ListUsers 获取用户列表（管理员接口）
// @Summary 获取用户列表
// @Description 管理员获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Param user_type query string false "用户类型"
// @Param status query string false "用户状态"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} ListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	userType := c.Query("user_type")
	status := c.Query("status")
	keyword := c.Query("keyword")

	req := &service.ListUsersRequest{
		Page:     page,
		Limit:    limit,
		UserType: userType,
		Status:   status,
		Keyword:  keyword,
	}

	ctx := SetClientInfoToContext(c)
	resp, err := h.userService.ListUsers(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取用户列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewListResponse("获取用户列表成功", resp.Users, resp.Total, resp.Page, resp.Limit))
}

// GetUserStatistics 获取用户统计信息（管理员接口）
// @Summary 获取用户统计信息
// @Description 管理员获取用户统计数据
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} SuccessResponse{data=service.UserStatistics}
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/users/statistics [get]
func (h *UserHandler) GetUserStatistics(c *gin.Context) {
	ctx := SetClientInfoToContext(c)
	resp, err := h.userService.GetUserStatistics(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取统计信息失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取统计信息成功", resp))
}

// GetUserIDFromContext 从上下文获取用户ID
func GetUserIDFromContext(c *gin.Context) uint64 {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint64); ok {
			return id
		}
	}
	return 0
}

// SetClientInfoToContext 设置客户端信息到上下文
func SetClientInfoToContext(c *gin.Context) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "client_ip", c.ClientIP())
	ctx = context.WithValue(ctx, "user_agent", c.GetHeader("User-Agent"))
	return ctx
}

// SubmitRealNameAuth 提交实名认证
func (h *UserHandler) SubmitRealNameAuth(c *gin.Context) {
	userID := GetUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", ""))
		return
	}

	var req service.RealNameAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	ctx := SetClientInfoToContext(c)
	err := h.userService.SubmitRealNameAuth(ctx, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "提交实名认证失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("提交实名认证成功", nil))
}

// SubmitBankCardAuth 提交银行卡认证
func (h *UserHandler) SubmitBankCardAuth(c *gin.Context) {
	userID := GetUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", ""))
		return
	}

	var req service.BankCardAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	ctx := SetClientInfoToContext(c)
	err := h.userService.SubmitBankCardAuth(ctx, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "提交银行卡认证失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("提交银行卡认证成功", nil))
}

// GetUserTags 获取用户标签
func (h *UserHandler) GetUserTags(c *gin.Context) {
	userID := GetUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", ""))
		return
	}

	tagType := c.Query("tag_type")

	ctx := SetClientInfoToContext(c)
	tags, err := h.userService.GetUserTags(ctx, userID, tagType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取用户标签失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取用户标签成功", tags))
}

// AddUserTag 添加用户标签
func (h *UserHandler) AddUserTag(c *gin.Context) {
	userID := GetUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", ""))
		return
	}

	var req service.AddTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	ctx := SetClientInfoToContext(c)
	err := h.userService.AddUserTag(ctx, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "添加用户标签失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("添加用户标签成功", nil))
}

// RemoveUserTag 删除用户标签
func (h *UserHandler) RemoveUserTag(c *gin.Context) {
	userID := GetUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", ""))
		return
	}

	tagKey := c.Param("tag_key")
	if tagKey == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "标签键不能为空", ""))
		return
	}

	ctx := SetClientInfoToContext(c)
	err := h.userService.RemoveUserTag(ctx, userID, tagKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "删除用户标签失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("删除用户标签成功", nil))
}
