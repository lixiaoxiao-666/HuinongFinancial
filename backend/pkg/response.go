package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse 统一API响应结构
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error_details,omitempty"`
}

// ListResponse 列表响应结构
type ListResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "Success",
		Data:    data,
	})
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// Created 创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Code:    0,
		Message: "Created successfully",
		Data:    data,
	})
}

// ListSuccess 列表成功响应
func ListSuccess(c *gin.Context, data interface{}, total int64) {
	c.JSON(http.StatusOK, ListResponse{
		Code:    0,
		Message: "Success",
		Data:    data,
		Total:   total,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string, details ...string) {
	response := APIResponse{
		Code:    code,
		Message: message,
	}

	if len(details) > 0 {
		response.Error = details[0]
	}

	var httpCode int
	switch {
	case code >= 1000 && code < 2000:
		httpCode = http.StatusBadRequest
	case code >= 2000 && code < 3000:
		httpCode = http.StatusUnauthorized
	case code >= 3000 && code < 4000:
		httpCode = http.StatusForbidden
	case code >= 4000 && code < 5000:
		httpCode = http.StatusNotFound
	case code >= 5000:
		httpCode = http.StatusInternalServerError
	default:
		httpCode = http.StatusBadRequest
	}

	c.JSON(httpCode, response)
}

// BadRequest 400错误响应
func BadRequest(c *gin.Context, message string, details ...string) {
	Error(c, 1001, message, details...)
}

// BadRequestWithMessage 400错误响应（自定义消息）
func BadRequestWithMessage(c *gin.Context, message string, details ...string) {
	Error(c, 1001, message, details...)
}

// Unauthorized 401错误响应
func Unauthorized(c *gin.Context, message string, details ...string) {
	Error(c, 2001, message, details...)
}

// UnauthorizedWithMessage 401错误响应（自定义消息）
func UnauthorizedWithMessage(c *gin.Context, message string, details ...string) {
	Error(c, 2001, message, details...)
}

// Forbidden 403错误响应
func Forbidden(c *gin.Context, message string, details ...string) {
	Error(c, 3001, message, details...)
}

// ForbiddenWithMessage 403错误响应（自定义消息）
func ForbiddenWithMessage(c *gin.Context, message string, details ...string) {
	Error(c, 3001, message, details...)
}

// NotFound 404错误响应
func NotFound(c *gin.Context, message string, details ...string) {
	Error(c, 4001, message, details...)
}

// NotFoundWithMessage 资源不存在响应（自定义消息）
func NotFoundWithMessage(c *gin.Context, message string) {
	Error(c, 4001, message)
}

// InternalError 500错误响应
func InternalError(c *gin.Context, message string, details ...string) {
	Error(c, 5001, message, details...)
}

// ValidationError 参数验证错误响应
func ValidationError(c *gin.Context, details string) {
	Error(c, 1002, "Invalid input parameters", details)
}

// InternalErrorWithMessage 内部错误响应（自定义消息）
func InternalErrorWithMessage(c *gin.Context, message string) {
	Error(c, 5001, message)
}
