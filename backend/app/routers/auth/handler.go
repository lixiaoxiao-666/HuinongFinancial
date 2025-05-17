package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"huinongfinancial/app/routers"
	v1 "huinongfinancial/app/types/api/v1"
	"huinongfinancial/app/types/api/v1/common"
)

// 注册路由
func handleRegister(c *gin.Context) {
	registerRequest := new(v1.RegisterRequest)
	if err := c.ShouldBindJSON(registerRequest); err != nil {
		logrus.Error("register request error: ", err)
		common.Fail(c, err)
		return
	}
	
}


//  登录路由
func handleLogin(c *gin.Context) {
	// 获取请求参数
}

func init() {
	// 注册路由
	router.V1.POST("/register", handleRegister)
	// 登录路由
	router.V1.POST("/login", handleLogin)
}