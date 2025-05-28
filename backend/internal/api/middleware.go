package api

import (
	"net/http"
	"strings"
	"time"

	"backend/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(jwtManager *pkg.JWTManager) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			pkg.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			pkg.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := jwtManager.ParseToken(tokenString)
		if err != nil {
			pkg.Unauthorized(c, "Invalid token", err.Error())
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_type", claims.UserType)
		c.Set("role", claims.Role)

		c.Next()
	})
}

// AdminAuthMiddleware OA后台认证中间件
func AdminAuthMiddleware(jwtManager *pkg.JWTManager) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			pkg.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			pkg.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := jwtManager.ParseToken(tokenString)
		if err != nil {
			pkg.Unauthorized(c, "Invalid token", err.Error())
			c.Abort()
			return
		}

		// 验证是否为OA用户
		if claims.UserType != "oa_user" {
			pkg.Forbidden(c, "Access denied")
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("admin_user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	})
}

// RoleMiddleware 角色权限中间件
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			pkg.Forbidden(c, "Role information not found")
			c.Abort()
			return
		}

		userRole := role.(string)
		for _, requiredRole := range requiredRoles {
			if userRole == requiredRole {
				c.Next()
				return
			}
		}

		pkg.Forbidden(c, "Insufficient permissions")
		c.Abort()
	})
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}

// RequestLoggerMiddleware 请求日志中间件
func RequestLoggerMiddleware(log *zap.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(start)

		// 获取客户端IP
		clientIP := c.ClientIP()

		// 获取请求方法
		method := c.Request.Method

		// 获取状态码
		statusCode := c.Writer.Status()

		// 记录日志
		if raw != "" {
			path = path + "?" + raw
		}

		log.Info("Request processed",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.String("client_ip", clientIP),
			zap.Duration("latency", latency),
		)
	})
}

// ErrorHandlerMiddleware 错误处理中间件
func ErrorHandlerMiddleware(log *zap.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)
				pkg.InternalError(c, "Internal server error")
			}
		}()

		c.Next()
	})
}

// RateLimitMiddleware 简单的请求限制中间件（基于IP）
func RateLimitMiddleware() gin.HandlerFunc {
	// 这里可以实现基于Redis的限流逻辑
	return gin.HandlerFunc(func(c *gin.Context) {
		// 简化实现，实际应该使用Redis计数器
		c.Next()
	})
}

// AIAgentAuthMiddleware AI智能体Token认证中间件
func AIAgentAuthMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			pkg.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		// 检查是否为AI Agent Token格式
		if strings.HasPrefix(authHeader, "AI-Agent-Token ") {
			token := strings.TrimPrefix(authHeader, "AI-Agent-Token ")

			// 这里应该验证AI Agent Token的有效性
			// 实际项目中应该从配置文件或数据库读取有效的Token列表
			validTokens := []string{
				"ai_agent_secure_token_123456",
				"dify_agent_token_789",
			}

			isValid := false
			for _, validToken := range validTokens {
				if token == validToken {
					isValid = true
					break
				}
			}

			if !isValid {
				pkg.Unauthorized(c, "Invalid AI Agent token")
				c.Abort()
				return
			}

			// 设置AI Agent上下文信息
			c.Set("agent_type", "ai_agent")
			c.Set("agent_token", token)
			c.Next()
			return
		}

		// 检查是否为系统Token格式
		if strings.HasPrefix(authHeader, "System-Token ") {
			token := strings.TrimPrefix(authHeader, "System-Token ")

			// 验证系统Token
			validSystemTokens := []string{
				"system_secure_token",
				"internal_api_token_456",
			}

			isValid := false
			for _, validToken := range validSystemTokens {
				if token == validToken {
					isValid = true
					break
				}
			}

			if !isValid {
				pkg.Unauthorized(c, "Invalid system token")
				c.Abort()
				return
			}

			// 设置系统调用上下文信息
			c.Set("agent_type", "system")
			c.Set("system_token", token)
			c.Next()
			return
		}

		pkg.Unauthorized(c, "Invalid authorization header format")
		c.Abort()
	})
}
