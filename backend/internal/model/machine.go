package model

import (
	"time"

	"gorm.io/gorm"
)

// Machine 农机设备表
type Machine struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	MachineCode string `gorm:"type:varchar(30);uniqueIndex;not null" json:"machine_code"`
	MachineName string `gorm:"type:varchar(100);not null" json:"machine_name"`

	// 设备分类：tractor(拖拉机)、harvester(收割机)、planter(播种机)、sprayer(喷药机)
	MachineType string `gorm:"type:varchar(30);not null" json:"machine_type"`

	// 品牌和型号
	Brand string `gorm:"type:varchar(50);not null" json:"brand"`
	Model string `gorm:"type:varchar(50);not null" json:"model"`

	// 规格参数(JSON格式)
	Specifications string `gorm:"type:json" json:"specifications"`

	// 设备描述
	Description string `gorm:"type:text" json:"description"`

	// 设备图片(JSON数组)
	Images string `gorm:"type:json" json:"images"`

	// 所有者信息
	OwnerID   uint64 `gorm:"not null;index" json:"owner_id"`
	OwnerType string `gorm:"type:varchar(20);not null" json:"owner_type"` // individual(个人)、company(公司)

	// 位置信息
	Province  string  `gorm:"type:varchar(50);not null" json:"province"`
	City      string  `gorm:"type:varchar(50);not null" json:"city"`
	County    string  `gorm:"type:varchar(50);not null" json:"county"`
	Address   string  `gorm:"type:varchar(200)" json:"address"`
	Longitude float64 `gorm:"type:decimal(10,6)" json:"longitude"`
	Latitude  float64 `gorm:"type:decimal(10,6)" json:"latitude"`

	// 租赁定价
	HourlyRate  int64 `json:"hourly_rate"`   // 小时租金(分)
	DailyRate   int64 `json:"daily_rate"`    // 日租金(分)
	PerAcreRate int64 `json:"per_acre_rate"` // 按亩收费(分)

	// 押金
	DepositAmount int64 `gorm:"not null" json:"deposit_amount"` // 押金金额(分)

	// 设备状态：available(可租)、rented(已租出)、maintenance(维护中)、offline(下线)
	Status string `gorm:"type:varchar(20);not null;default:'available'" json:"status"`

	// 可用时间(JSON格式)
	AvailableSchedule string `gorm:"type:json" json:"available_schedule"`

	// 设备评分
	Rating      float32 `gorm:"type:decimal(3,1);default:0.0" json:"rating"`
	RatingCount int     `gorm:"default:0" json:"rating_count"`

	// 租赁条件
	MinRentalHours int `gorm:"default:1" json:"min_rental_hours"`  // 最少租赁小时数
	MaxAdvanceDays int `gorm:"default:30" json:"max_advance_days"` // 最多提前预订天数

	// 认证信息
	IsVerified bool       `gorm:"default:false" json:"is_verified"` // 是否已认证
	VerifiedAt *time.Time `json:"verified_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Owner        User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	RentalOrders []RentalOrder `gorm:"foreignKey:MachineID" json:"rental_orders,omitempty"`
}

// RentalOrder 租赁订单表
type RentalOrder struct {
	ID      uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo string `gorm:"type:varchar(30);uniqueIndex;not null" json:"order_no"`

	// 关联信息
	MachineID uint64 `gorm:"not null;index" json:"machine_id"`
	RenterID  uint64 `gorm:"not null;index" json:"renter_id"`
	OwnerID   uint64 `gorm:"not null;index" json:"owner_id"`

	// 租赁时间
	StartTime      time.Time `gorm:"not null" json:"start_time"`
	EndTime        time.Time `gorm:"not null" json:"end_time"`
	RentalDuration int       `json:"rental_duration"` // 租赁时长(小时)

	// 租赁地点
	RentalLocation string `gorm:"type:varchar(200)" json:"rental_location"`
	ContactPerson  string `gorm:"type:varchar(50)" json:"contact_person"`
	ContactPhone   string `gorm:"type:varchar(20)" json:"contact_phone"`

	// 计费方式：hourly(按小时)、daily(按天)、per_acre(按亩)
	BillingMethod string `gorm:"type:varchar(20);not null" json:"billing_method"`

	// 费用计算
	UnitPrice      int64   `gorm:"not null" json:"unit_price"`      // 单价(分)
	Quantity       float64 `gorm:"not null" json:"quantity"`        // 数量(小时/天/亩)
	SubtotalAmount int64   `gorm:"not null" json:"subtotal_amount"` // 小计(分)
	DepositAmount  int64   `gorm:"not null" json:"deposit_amount"`  // 押金(分)
	TotalAmount    int64   `gorm:"not null" json:"total_amount"`    // 总金额(分)

	// 订单状态：pending(待确认)、confirmed(已确认)、paid(已支付)、in_progress(进行中)、
	// completed(已完成)、cancelled(已取消)、disputed(有争议)
	Status string `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`

	// 支付信息
	PaymentMethod string     `gorm:"type:varchar(20)" json:"payment_method"`
	PaymentStatus string     `gorm:"type:varchar(20)" json:"payment_status"`
	PaidAmount    int64      `json:"paid_amount"`
	PaidAt        *time.Time `json:"paid_at"`

	// 确认信息
	OwnerConfirmedAt  *time.Time `json:"owner_confirmed_at"`
	RenterConfirmedAt *time.Time `json:"renter_confirmed_at"`

	// 服务评价
	RenterRating  float32 `gorm:"type:decimal(3,1)" json:"renter_rating"`
	RenterComment string  `gorm:"type:text" json:"renter_comment"`
	OwnerRating   float32 `gorm:"type:decimal(3,1)" json:"owner_rating"`
	OwnerComment  string  `gorm:"type:text" json:"owner_comment"`

	// 备注
	Remarks string `gorm:"type:text" json:"remarks"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	Machine Machine `gorm:"foreignKey:MachineID" json:"machine,omitempty"`
	Renter  User    `gorm:"foreignKey:RenterID" json:"renter,omitempty"`
	Owner   User    `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

// MachineSpecs 设备规格参数结构
type MachineSpecs struct {
	Power        string            `json:"power"`         // 功率
	WorkingWidth string            `json:"working_width"` // 工作幅宽
	Weight       string            `json:"weight"`        // 重量
	FuelType     string            `json:"fuel_type"`     // 燃料类型
	YearOfMake   int               `json:"year_of_make"`  // 制造年份
	WorkingSpeed string            `json:"working_speed"` // 工作速度
	Capacity     string            `json:"capacity"`      // 容量
	Other        map[string]string `json:"other"`         // 其他参数
}

// AvailableSchedule 可用时间表结构
type AvailableSchedule struct {
	Monday    []TimeSlot `json:"monday"`
	Tuesday   []TimeSlot `json:"tuesday"`
	Wednesday []TimeSlot `json:"wednesday"`
	Thursday  []TimeSlot `json:"thursday"`
	Friday    []TimeSlot `json:"friday"`
	Saturday  []TimeSlot `json:"saturday"`
	Sunday    []TimeSlot `json:"sunday"`
}

// TimeSlot 时间段结构
type TimeSlot struct {
	StartTime string `json:"start_time"` // "08:00"
	EndTime   string `json:"end_time"`   // "18:00"
}

// TableName 设置表名
func (Machine) TableName() string {
	return "machines"
}

func (RentalOrder) TableName() string {
	return "rental_orders"
}

// RentalRating 租赁评价表
type RentalRating struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID    uint64    `gorm:"not null;index" json:"order_id"`
	RaterID    uint64    `gorm:"not null;index" json:"rater_id"`
	RatingType string    `gorm:"type:varchar(20);not null" json:"rating_type"` // "renter" 或 "owner"
	Rating     float32   `gorm:"type:decimal(3,1);not null" json:"rating"`
	Comment    string    `gorm:"type:text" json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// RentalApplication 租赁申请表
type RentalApplication struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ApplicationNo  string     `gorm:"type:varchar(30);uniqueIndex;not null" json:"application_no"`
	UserID         uint64     `gorm:"not null;index" json:"user_id"`
	MachineID      uint64     `gorm:"not null;index" json:"machine_id"`
	StartTime      time.Time  `gorm:"not null" json:"start_time"`
	EndTime        time.Time  `gorm:"not null" json:"end_time"`
	RentalLocation string     `gorm:"type:varchar(200)" json:"rental_location"`
	ContactPerson  string     `gorm:"type:varchar(50)" json:"contact_person"`
	ContactPhone   string     `gorm:"type:varchar(20)" json:"contact_phone"`
	BillingMethod  string     `gorm:"type:varchar(20);not null" json:"billing_method"`
	Quantity       float64    `gorm:"not null" json:"quantity"`
	TotalAmount    int64      `gorm:"not null" json:"total_amount"`
	DepositAmount  int64      `gorm:"not null" json:"deposit_amount"`
	Status         string     `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	RiskLevel      string     `gorm:"type:varchar(20)" json:"risk_level"`
	RiskFactors    string     `gorm:"type:json" json:"risk_factors"`
	AIAssessment   string     `gorm:"type:text" json:"ai_assessment"`
	ReviewerID     *uint64    `gorm:"index" json:"reviewer_id"`
	ReviewedAt     *time.Time `json:"reviewed_at"`
	ReviewNote     string     `gorm:"type:text" json:"review_note"`
	Remarks        string     `gorm:"type:text" json:"remarks"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// RentalReviewLog 租赁审核日志表
type RentalReviewLog struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ApplicationID uint64    `gorm:"not null;index" json:"application_id"`
	ReviewerID    uint64    `gorm:"not null;index" json:"reviewer_id"`
	Action        string    `gorm:"type:varchar(20);not null" json:"action"` // "approve", "reject", "request_more_info"
	Note          string    `gorm:"type:text" json:"note"`
	CreatedAt     time.Time `json:"created_at"`
}

// TableName 设置表名
func (RentalRating) TableName() string {
	return "rental_ratings"
}

func (RentalApplication) TableName() string {
	return "rental_applications"
}

func (RentalReviewLog) TableName() string {
	return "rental_review_logs"
}
