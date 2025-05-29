package api

import (
	"backend/internal/service"
	"backend/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
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

// RegisterUserRoutes 注册用户路由
func RegisterUserRoutes(group *gin.RouterGroup, handler *UserHandler, authMiddleware gin.HandlerFunc) {
	// 无需认证的路由
	group.POST("/register", handler.Register)
	group.POST("/login", handler.Login)

	// 需要认证的路由
	authenticated := group.Group("")
	authenticated.Use(authMiddleware)
	{
		authenticated.GET("/me", handler.GetUserInfo)
		authenticated.PUT("/me", handler.UpdateUserInfo)
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 用户服务
// @Accept json
// @Produce json
// @Param request body UserRegisterRequest true "注册请求"
// @Success 201 {object} CommonResponse{data=UserRegisterResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("用户注册参数绑定失败", zap.Error(err))
		pkg.BadRequestWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	// 验证密码确认
	if req.Password != req.ConfirmPassword {
		pkg.BadRequestWithMessage(c, "两次输入的密码不一致")
		return
	}

	// 构建服务层请求
	serviceReq := &service.RegisterRequest{
		Phone:    req.Phone,
		Password: req.Password,
	}

	// 调用服务层进行注册
	loginResp, err := h.userService.Register(c.Request.Context(), serviceReq)
	if err != nil {
		h.log.Error("用户注册失败", zap.Error(err))

		// 根据错误类型返回不同的响应
		switch err.Error() {
		case "手机号已注册":
			pkg.BadRequestWithMessage(c, "该手机号已被注册")
		default:
			pkg.InternalError(c, "注册失败，请稍后重试")
		}
		return
	}

	// 记录成功日志
	h.log.Info("用户注册成功", zap.String("user_id", loginResp.UserID), zap.String("phone", req.Phone))

	// 返回成功响应
	c.JSON(201, CommonResponse{
		Code:    CodeSuccess,
		Message: "Success",
		Data: UserRegisterResponse{
			UserID: loginResp.UserID,
		},
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 用户服务
// @Accept json
// @Produce json
// @Param request body UserLoginRequest true "登录请求"
// @Success 200 {object} CommonResponse{data=UserLoginResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("用户登录参数绑定失败", zap.Error(err))
		pkg.BadRequestWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	// 构建服务层请求
	serviceReq := &service.LoginRequest{
		Phone:    req.Username, // 使用username字段作为手机号
		Password: req.Password,
	}

	// 调用服务层进行登录
	loginResp, err := h.userService.Login(c.Request.Context(), serviceReq)
	if err != nil {
		h.log.Warn("用户登录失败", zap.Error(err), zap.String("username", req.Username))

		// 根据错误类型返回不同的响应
		switch err.Error() {
		case "用户不存在":
			pkg.UnauthorizedWithMessage(c, "手机号或密码错误")
		case "密码错误":
			pkg.UnauthorizedWithMessage(c, "手机号或密码错误")
		case "用户已被禁用":
			pkg.UnauthorizedWithMessage(c, "账户已被禁用，请联系客服")
		default:
			pkg.InternalError(c, "登录失败，请稍后重试")
		}
		return
	}

	// 记录成功日志
	h.log.Info("用户登录成功", zap.String("user_id", loginResp.UserID), zap.String("username", req.Username))

	// 返回成功响应
	pkg.Success(c, UserLoginResponse{
		Token: loginResp.Token,
		UserInfo: User{
			UserID:     loginResp.UserID,
			Username:   req.Username,
			RealName:   "", // 需要从profile获取
			Phone:      "",
			Email:      "",
			IDCard:     "",
			Address:    "",
			IsVerified: false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	})
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的信息
// @Tags 用户服务
// @Produce json
// @Success 200 {object} CommonResponse{data=UserInfoResponse}
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/me [get]
// @Security BearerAuth
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		h.log.Warn("获取用户信息时未找到用户ID")
		pkg.UnauthorizedWithMessage(c, "用户未登录")
		return
	}

	// 调用服务层获取用户信息
	userInfo, err := h.userService.GetUserInfo(c.Request.Context(), userID.(string))
	if err != nil {
		h.log.Error("获取用户信息失败", zap.Error(err), zap.String("user_id", userID.(string)))

		// 根据错误类型返回不同的响应
		switch err.Error() {
		case "用户不存在":
			pkg.NotFoundWithMessage(c, "用户不存在")
		default:
			pkg.InternalError(c, "获取用户信息失败，请稍后重试")
		}
		return
	}

	// 返回成功响应
	pkg.Success(c, userInfo)
}

// UpdateUserInfo 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前登录用户的信息
// @Tags 用户服务
// @Accept json
// @Produce json
// @Param request body UserUpdateRequest true "更新请求"
// @Success 200 {object} CommonResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/me [put]
// @Security BearerAuth
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		h.log.Warn("更新用户信息时未找到用户ID")
		pkg.UnauthorizedWithMessage(c, "用户未登录")
		return
	}

	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("更新用户信息参数绑定失败", zap.Error(err))
		pkg.BadRequestWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	// 构建服务层请求
	serviceReq := &service.UpdateUserRequest{
		Nickname:  req.RealName, // 暂时用RealName作为Nickname
		AvatarURL: "",           // 头像URL暂时留空
	}

	// 调用服务层更新用户信息
	err := h.userService.UpdateUserInfo(c.Request.Context(), userID.(string), serviceReq)
	if err != nil {
		h.log.Error("更新用户信息失败", zap.Error(err), zap.String("user_id", userID.(string)))

		// 根据错误类型返回不同的响应
		switch err.Error() {
		case "用户不存在":
			pkg.NotFoundWithMessage(c, "用户不存在")
		default:
			pkg.InternalError(c, "更新用户信息失败，请稍后重试")
		}
		return
	}

	// 记录成功日志
	h.log.Info("用户信息更新成功", zap.String("user_id", userID.(string)))

	// 返回成功响应
	pkg.Success(c, gin.H{"message": "用户信息更新成功"})
}

// 辅助函数：脱敏手机号
func maskPhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}

// 辅助函数：脱敏身份证号
func maskIDCard(idCard string) string {
	if len(idCard) < 8 {
		return idCard
	}
	if len(idCard) == 18 {
		return idCard[:3] + "***********" + idCard[14:]
	}
	// 15位身份证号
	return idCard[:3] + "***********" + idCard[11:]
}
