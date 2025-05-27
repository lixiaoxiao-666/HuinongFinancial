package api

import (
	"backend/internal/service"
	"backend/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
	log         *zap.Logger
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService *service.UserService, log *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		log:         log,
	}
}

// SendVerificationCode 发送验证码
func (h *UserHandler) SendVerificationCode(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required,len=11"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ValidationError(c, err.Error())
		return
	}

	if err := h.userService.SendVerificationCode(c.Request.Context(), req.Phone); err != nil {
		h.log.Error("发送验证码失败", zap.Error(err))
		pkg.BadRequest(c, err.Error())
		return
	}

	pkg.SuccessWithMessage(c, "验证码发送成功", nil)
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ValidationError(c, err.Error())
		return
	}

	result, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		h.log.Error("用户注册失败", zap.Error(err))
		pkg.BadRequest(c, err.Error())
		return
	}

	pkg.Created(c, result)
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ValidationError(c, err.Error())
		return
	}

	result, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		h.log.Error("用户登录失败", zap.Error(err))
		pkg.BadRequest(c, err.Error())
		return
	}

	pkg.Success(c, result)
}

// GetUserInfo 获取用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		pkg.Unauthorized(c, "User not found in context")
		return
	}

	result, err := h.userService.GetUserInfo(c.Request.Context(), userID.(string))
	if err != nil {
		h.log.Error("获取用户信息失败", zap.Error(err))
		pkg.BadRequest(c, err.Error())
		return
	}

	pkg.Success(c, result)
}

// UpdateUserInfo 更新用户信息
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		pkg.Unauthorized(c, "User not found in context")
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ValidationError(c, err.Error())
		return
	}

	if err := h.userService.UpdateUserInfo(c.Request.Context(), userID.(string), &req); err != nil {
		h.log.Error("更新用户信息失败", zap.Error(err))
		pkg.BadRequest(c, err.Error())
		return
	}

	pkg.SuccessWithMessage(c, "用户信息更新成功", nil)
}

// RegisterUserRoutes 注册用户路由
func RegisterUserRoutes(r *gin.RouterGroup, handler *UserHandler, authMiddleware gin.HandlerFunc) {
	// 无需认证的路由
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.POST("/send-verification-code", handler.SendVerificationCode)

	// 需要认证的路由
	auth := r.Group("", authMiddleware)
	{
		auth.GET("/me", handler.GetUserInfo)
		auth.PUT("/me", handler.UpdateUserInfo)
	}
}
