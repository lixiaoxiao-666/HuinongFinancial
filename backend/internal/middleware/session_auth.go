package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"huinong-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SessionAuthMiddleware 基于Redis会话的认证中间件
type SessionAuthMiddleware struct {
	sessionService service.SessionService
}

// NewSessionAuthMiddleware 创建会话认证中间件实例
func NewSessionAuthMiddleware(sessionService service.SessionService) *SessionAuthMiddleware {
	return &SessionAuthMiddleware{
		sessionService: sessionService,
	}
}

// RequireAuth 要求认证的中间件
func (m *SessionAuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
				"error":   "Authorization header required",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
				"error":   "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		accessToken := parts[1]

		// 验证Token并获取会话信息
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		sessionInfo, err := m.sessionService.ValidateToken(ctx, accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户和会话信息存储到上下文
		c.Set("user_id", sessionInfo.UserID)
		c.Set("session_id", sessionInfo.SessionID)
		c.Set("platform", sessionInfo.Platform)
		c.Set("device_info", sessionInfo.DeviceInfo)
		c.Set("network_info", sessionInfo.NetworkInfo)

		c.Next()
	}
}

// OptionalAuth 可选认证中间件
func (m *SessionAuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		accessToken := parts[1]

		// 验证Token
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		sessionInfo, err := m.sessionService.ValidateToken(ctx, accessToken)
		if err == nil && sessionInfo != nil {
			// 验证成功，存储用户信息
			c.Set("user_id", sessionInfo.UserID)
			c.Set("session_id", sessionInfo.SessionID)
			c.Set("platform", sessionInfo.Platform)
			c.Set("device_info", sessionInfo.DeviceInfo)
			c.Set("network_info", sessionInfo.NetworkInfo)
		}

		c.Next()
	}
}

// RequireRole 角色权限中间件（仅用于OA系统）
func (m *SessionAuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 首先检查是否已经认证
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
				"error":   "User not authenticated",
			})
			c.Abort()
			return
		}

		// 检查是否为OA平台
		platform, exists := c.Get("platform")
		if !exists || platform != "oa" {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "访问被拒绝",
				"error":   "Role check only available for OA platform",
			})
			c.Abort()
			return
		}

		// TODO: 实现OA用户角色检查逻辑
		// 这里需要从数据库查询OAUser和对应的角色信息
		// 暂时简单处理，后续需要注入UserService来查询角色
		userIDVal := userID.(uint64)
		_ = userIDVal

		// 临时处理：如果需要admin角色，检查用户ID是否为管理员
		// 实际应该查询oa_users表和oa_roles表
		for _, requiredRole := range roles {
			if requiredRole == "admin" {
				// TODO: 查询数据库验证用户角色
				// 现在暂时允许通过，后续需要实现具体的角色查询逻辑
				c.Set("user_role", "admin") // 临时设置
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "权限不足",
			"error":   "Required role not found",
		})
		c.Abort()
	}
}

// CheckPlatform 检查平台权限的中间件
func (m *SessionAuthMiddleware) CheckPlatform(requiredPlatform string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户平台信息（在RequireAuth中间件中设置）
		platform, exists := c.Get("platform")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "访问被拒绝",
				"error":   "Platform information not found",
			})
			c.Abort()
			return
		}

		// 检查平台权限
		if platform != requiredPlatform {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "访问被拒绝",
				"error":   "Platform access not authorized",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RefreshToken 刷新Token中间件
func (m *SessionAuthMiddleware) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var refreshToken string

		// 首先尝试从form-data获取
		refreshToken = c.PostForm("refresh_token")

		// 如果form-data中没有，尝试从JSON body获取
		if refreshToken == "" {
			var requestBody struct {
				RefreshToken string `json:"refresh_token"`
			}
			if err := c.ShouldBindJSON(&requestBody); err == nil {
				refreshToken = requestBody.RefreshToken
			}
		}

		if refreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
				"error":   "refresh_token required",
			})
			c.Abort()
			return
		}

		// 刷新Token
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		tokenPair, err := m.sessionService.RefreshSession(ctx, refreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "刷新失败",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "刷新成功",
			"data":    tokenPair,
		})
	}
}

// Logout 登出中间件
func (m *SessionAuthMiddleware) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, exists := c.Get("session_id")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
				"error":   "No active session",
			})
			c.Abort()
			return
		}

		// 注销会话
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		err := m.sessionService.RevokeSession(ctx, sessionID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "注销失败",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "注销成功",
		})
	}
}

// SessionInfo 获取会话信息中间件
func (m *SessionAuthMiddleware) SessionInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
			})
			c.Abort()
			return
		}

		// 获取用户所有会话
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		sessions, err := m.sessionService.GetUserSessions(ctx, userID.(uint64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取会话信息失败",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取成功",
			"data":    sessions,
		})
	}
}

// RevokeOtherSessions 注销其他会话中间件
func (m *SessionAuthMiddleware) RevokeOtherSessions() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
			})
			c.Abort()
			return
		}

		sessionID, exists := c.Get("session_id")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			c.Abort()
			return
		}

		// 注销用户的其他会话
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		err := m.sessionService.RevokeUserSessions(ctx, userID.(uint64), sessionID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "注销其他会话失败",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "注销其他会话成功",
		})
	}
}
