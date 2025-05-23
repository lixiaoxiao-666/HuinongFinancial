package router

import (
	"github.com/gin-gonic/gin"
	"huinongfinancial/router/middleware"
)

var (
	router *gin.Engine
	v1 *gin.RouterGroup
)

func InitRouter() {
	router := gin.Default()
	// 使用中间件
	router.Use(middleware.Cors())
	// 创建 v1 路由组
	v1 = router.Group("/api/v1")

	// 创建 /livez 的路由
	router.GET("/livez", func(c *gin.Context) {
		c.String(200, "livez")
	})
	// 创建 /readyz 的路由
	router.GET("/readyz", func(c *gin.Context) {
		c.String(200, "readyz")
	})
}

// V1 返回 /api/v1 的路由组，用于控制平面资源的 API 端点。
func V1() *gin.RouterGroup {
	return v1
}

// Router returns the main Gin engine instance.
// Router 返回主 Gin 引擎实例。
func Router() *gin.Engine {
	return router
}