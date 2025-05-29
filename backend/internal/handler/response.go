package handler

// SuccessResponse 成功响应结构体
type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse 错误响应结构体
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ListResponse 列表响应结构体
type ListResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(message string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string, err string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
		Error:   err,
	}
}

// NewListResponse 创建列表响应
func NewListResponse(message string, data interface{}, total int64, page, limit int) *ListResponse {
	return &ListResponse{
		Code:    200,
		Message: message,
		Data:    data,
		Total:   total,
		Page:    page,
		Limit:   limit,
	}
}
