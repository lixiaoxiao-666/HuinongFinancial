package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// DifyAuthMiddleware Dify API认证中间件
func DifyAuthMiddleware(expectedToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "缺少Authorization头",
			})
			c.Abort()
			return
		}

		// 检查Bearer token格式
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "无效的Authorization格式",
			})
			c.Abort()
			return
		}

		// 提取token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "token不能为空",
			})
			c.Abort()
			return
		}

		// 验证token
		if token != expectedToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "无效的API token",
			})
			c.Abort()
			return
		}

		// 检查Dify来源头（可选）
		difySource := c.GetHeader("X-Dify-Source")
		if difySource != "" && difySource != "workflow" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "无效的请求来源",
			})
			c.Abort()
			return
		}

		// 设置上下文标识，表明这是来自Dify的请求
		c.Set("is_dify_request", true)
		c.Set("dify_source", difySource)

		// 继续处理请求
		c.Next()
	}
}

// DifyIPWhitelistMiddleware IP白名单中间件（可选）
func DifyIPWhitelistMiddleware(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(allowedIPs) == 0 {
			// 如果没有配置白名单，则跳过检查
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		
		// 检查IP是否在白名单中
		allowed := false
		for _, ip := range allowedIPs {
			if ip == clientIP || ip == "*" {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "IP地址不在白名单中",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// DifyRequestValidationMiddleware 请求验证中间件
func DifyRequestValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查Content-Type
		contentType := c.GetHeader("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Content-Type必须为application/json",
			})
			c.Abort()
			return
		}

		// 检查请求方法
		if c.Request.Method != "POST" {
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"success": false,
				"error":   "仅支持POST方法",
			})
			c.Abort()
			return
		}

		c.Next()
	}
} 