package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
	jwtSecret string
}

// NewAuthMiddleware 创建认证中间件实例
func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}
}

// RequireAuth JWT认证中间件
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
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

		tokenString := parts[1]

		// 验证Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
				"error":   "Invalid token: " + err.Error(),
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 检查Token类型
			if tokenType, ok := claims["type"].(string); ok && tokenType != "access" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "未授权",
					"error":   "Invalid token type",
				})
				c.Abort()
				return
			}

			// 提取用户ID
			if userID, ok := claims["user_id"].(float64); ok {
				c.Set("user_id", uint64(userID))
			}

			// 提取平台信息
			if platform, ok := claims["platform"].(string); ok {
				c.Set("platform", platform)
			}

			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
				"error":   "Invalid token claims",
			})
			c.Abort()
			return
		}
	}
}

// OptionalAuth 可选认证中间件（不强制要求认证）
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
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

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return []byte(m.jwtSecret), nil
		})

		if err == nil {
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if userID, ok := claims["user_id"].(float64); ok {
					c.Set("user_id", uint64(userID))
				}
				if platform, ok := claims["platform"].(string); ok {
					c.Set("platform", platform)
				}
			}
		}

		c.Next()
	}
}

// RequireRole 角色权限中间件
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 首先需要认证
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

		// TODO: 实现角色检查逻辑
		// 这里应该从数据库查询用户角色并验证
		_ = userID

		c.Next()
	}
}

// AdminAuth 管理员认证中间件
func (m *AuthMiddleware) AdminAuth() gin.HandlerFunc {
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

		// 验证管理员Token
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

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
				"error":   "Invalid token: " + err.Error(),
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 检查是否为管理员平台
			if platform, ok := claims["platform"].(string); ok && platform != "admin" {
				c.JSON(http.StatusForbidden, gin.H{
					"code":    403,
					"message": "权限不足",
					"error":   "Admin access required",
				})
				c.Abort()
				return
			}

			// 提取用户ID
			if userID, ok := claims["user_id"].(float64); ok {
				c.Set("admin_user_id", uint64(userID))
			}

			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
				"error":   "Invalid token claims",
			})
			c.Abort()
			return
		}
	}
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Device-ID, Device-Type, App-Version")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RequestLogger 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// RateLimiter 限流中间件
func RateLimiter() gin.HandlerFunc {
	// TODO: 实现基于Redis的限流逻辑
	return func(c *gin.Context) {
		c.Next()
	}
}

// Recovery 错误恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务器内部错误",
			"error":   "Internal server error",
		})
	})
}

// 上下文工具函数
func SetContextValue(c *gin.Context, key string, value interface{}) context.Context {
	ctx := c.Request.Context()
	return context.WithValue(ctx, key, value)
}

func GetContextValue(ctx context.Context, key string) interface{} {
	return ctx.Value(key)
}
