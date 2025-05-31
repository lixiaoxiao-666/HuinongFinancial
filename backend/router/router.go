package router

import (
	"huinong-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// 跨域中间件
	router.Use(middlewares.CORSMiddleware())

	// api := router.Group("/api/v1")
	// auth := api.Group("/auth")

	// 需要认证的接口
	// authorized := api.Group("/")
	// authorized.Use(middlewares.AuthMiddleware())
	// {

	// }

	return router
}
