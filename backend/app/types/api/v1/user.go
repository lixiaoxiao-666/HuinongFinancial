package v1

import (
	"gorm.io/gorm"
)

// 用户结构体
type User struct {
	gorm.Model
	// binding:"required" 是gin框架的验证规则，表示该字段必须存在
	// 用户ID
	ID       int    `json:"id" binding:"required" gorm:"unique"`
	// 用户名称
	Username string `json:"username" binding:"required"`
	// 用户密码
	Password string `json:"passw" binding:"required"`
	// 用户手机号
	Phone    string `json:"phone" binding:"required" gorm:"unique"`
	// 用户是否验证
	Verifed	 bool 	`json:"verifed" binding:"required" gorm:"default:false"`
}