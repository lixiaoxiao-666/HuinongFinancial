package v1

import (
	"gorm.io/gorm"
)

// 财务结构体
type Finance struct {
	gorm.Model
	// 用户ID
	ID       int    `json:"id" binding:"required" gorm:"unique"`
	// 用户名
	Username string `json:"username" binding:"required"`
	// 余额
	Money    float64 `json:"money" binding:"required" gorm:"default:0"`
}

