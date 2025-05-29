package pkg

import (
	"crypto/rand"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	timestamp := time.Now().Unix()
	return fmt.Sprintf("app_%d", timestamp)
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
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	return fmt.Sprintf("file_%d_%x", timestamp, randomBytes)
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

// GetIntParam 从查询参数中获取整数值，如果无法解析则返回默认值
func GetIntParam(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}

	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultValue
}

// Contains 检查切片中是否包含指定的字符串
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetClientIP 获取客户端IP地址
func GetClientIP(r *http.Request) string {
	// 优先检查 X-Forwarded-For 头
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For 可能包含多个IP，取第一个
		if commaIdx := strings.Index(xff, ","); commaIdx != -1 {
			return strings.TrimSpace(xff[:commaIdx])
		}
		return strings.TrimSpace(xff)
	}

	// 检查 X-Real-IP 头
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// 检查 X-Forwarded-Proto 头
	if xfp := r.Header.Get("X-Forwarded-Proto"); xfp != "" {
		return strings.TrimSpace(xfp)
	}

	// 最后使用 RemoteAddr
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}

	return r.RemoteAddr
}

// MaskName 脱敏姓名
func MaskName(name string) string {
	if len(name) == 0 {
		return name
	}
	if len(name) == 1 {
		return "*"
	}
	if len(name) == 2 {
		return string([]rune(name)[0]) + "*"
	}
	runes := []rune(name)
	if len(runes) <= 3 {
		return string(runes[0]) + "*" + string(runes[len(runes)-1])
	}
	// 对于超过3个字符的姓名，保留首尾，中间用*代替
	return string(runes[0]) + strings.Repeat("*", len(runes)-2) + string(runes[len(runes)-1])
}

// MaskIDCard 脱敏身份证号
func MaskIDCard(idCard string) string {
	if len(idCard) < 8 {
		return idCard
	}
	if len(idCard) == 18 {
		return idCard[:3] + "***********" + idCard[14:]
	}
	// 15位身份证号
	if len(idCard) == 15 {
		return idCard[:3] + "********" + idCard[11:]
	}
	// 其他长度，保留前3位和后4位
	if len(idCard) > 7 {
		return idCard[:3] + strings.Repeat("*", len(idCard)-7) + idCard[len(idCard)-4:]
	}
	return idCard
}

// MaskPhone 脱敏手机号
func MaskPhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}
