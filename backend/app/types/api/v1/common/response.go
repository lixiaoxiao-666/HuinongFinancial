package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Response(c *gin.Context, err error, data interface{}) {
	code := 200
	message := "success"
	if err != nil{
		code = 500
		message = err.Error()
	}
	c.JSON(http.StatusOK, BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Success(c *gin.Context, data interface{}) {
	Response(c, nil, data)
}

func Fail(c *gin.Context, err error) {
	Response(c, err, nil)
}
