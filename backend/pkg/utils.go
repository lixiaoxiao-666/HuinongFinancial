package pkg

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

// GenerateUUID 生成UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateUserID 生成用户ID
func GenerateUserID() string {
	return fmt.Sprintf("usr_%s", GenerateUUID()[:12])
}

// GenerateApplicationID 生成申请ID
func GenerateApplicationID() string {
	return fmt.Sprintf("la_%s", GenerateUUID()[:12])
}

// GenerateLoanApplicationID 生成贷款申请ID (别名)
func GenerateLoanApplicationID() string {
	return GenerateApplicationID()
}

// GenerateMachineryID 生成农机ID
func GenerateMachineryID() string {
	return fmt.Sprintf("fm_%s", GenerateUUID()[:12])
}

// GenerateOrderID 生成订单ID
func GenerateOrderID() string {
	return fmt.Sprintf("mlo_%s", GenerateUUID()[:12])
}

// GenerateFileID 生成文件ID
func GenerateFileID() string {
	return fmt.Sprintf("file_%s", GenerateUUID()[:12])
}

// GenerateOAUserID 生成OA用户ID
func GenerateOAUserID() string {
	return fmt.Sprintf("oa_%s", GenerateUUID()[:12])
}

// GenerateProductID 生成产品ID
func GenerateProductID() string {
	return fmt.Sprintf("lp_%s", GenerateUUID()[:12])
}

// StringToInt 字符串转整数
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// StringToInt64 字符串转int64
func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// GetPagination 计算分页参数
func GetPagination(page, limit int) (offset int, validLimit int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	offset = (page - 1) * limit
	return offset, limit
}

// GetCurrentTimestamp 获取当前时间戳
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// GetIntParam 从gin.Context中获取整数参数
func GetIntParam(c *gin.Context, paramName string, defaultValue int) int {
	paramStr := c.Query(paramName)
	if paramStr == "" {
		return defaultValue
	}
	
	value, err := strconv.Atoi(paramStr)
	if err != nil {
		return defaultValue
	}
	
	return value
}
