package v1

import (
	"gorm.io/gorm"
)

// 机械结构体
type Machinery struct {
	gorm.Model
	// 机械ID
	ID       int    `json:"id" binding:"required" gorm:"unique"`
	// 机械名称
	Name     string `json:"name" binding:"required"`
	// 价格
	Price    float64 `json:"price" binding:"required"`
	// 数量
	Quantity int     `json:"quantity" binding:"required"`
	// 图片
	Image    string  `json:"image" binding:"required"`
	// 描述
	Description string  `json:"description" binding:"required"`
}