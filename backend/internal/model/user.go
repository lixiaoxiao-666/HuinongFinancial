package model

import (
	"time"

	"gorm.io/gorm"
)

// User 核心用户表
type User struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID     string `gorm:"type:varchar(36);uniqueIndex;not null" json:"uuid"`
	Username string `gorm:"type:varchar(50);uniqueIndex" json:"username"`
	Phone    string `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`
	Email    string `gorm:"type:varchar(100);index" json:"email"`

	// 密码相关 (不在JSON中返回)
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Salt         string `gorm:"type:varchar(32);not null" json:"-"`

	// 用户类型：farmer(个体农户)、farm_owner(农场主)、cooperative(合作社)、enterprise(企业)
	UserType string `gorm:"type:varchar(20);not null;default:'farmer'" json:"user_type"`

	// 用户状态：active(正常)、frozen(冻结)、deleted(删除)
	Status string `gorm:"type:varchar(20);not null;default:'active'" json:"status"`

	// 基本信息
	RealName string     `gorm:"type:varchar(50)" json:"real_name"`
	IDCard   string     `gorm:"type:varchar(18);index" json:"id_card"`
	Avatar   string     `gorm:"type:varchar(255)" json:"avatar"`
	Gender   string     `gorm:"type:varchar(10)" json:"gender"` // male, female, unknown
	Birthday *time.Time `json:"birthday"`

	// 地址信息
	Province string `gorm:"type:varchar(50)" json:"province"`
	City     string `gorm:"type:varchar(50)" json:"city"`
	County   string `gorm:"type:varchar(50)" json:"county"`
	Address  string `gorm:"type:varchar(200)" json:"address"`

	// 认证状态
	IsRealNameVerified bool `gorm:"default:false" json:"is_real_name_verified"`
	IsBankCardVerified bool `gorm:"default:false" json:"is_bank_card_verified"`
	IsCreditVerified   bool `gorm:"default:false" json:"is_credit_verified"`

	// 登录信息
	LastLoginTime *time.Time `json:"last_login_time"`
	LastLoginIP   string     `gorm:"type:varchar(45)" json:"last_login_ip"`
	LoginCount    uint32     `gorm:"default:0" json:"login_count"`

	// 注册信息
	RegisterIP   string    `gorm:"type:varchar(45)" json:"register_ip"`
	RegisterTime time.Time `json:"register_time"`

	// 时间字段
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联数据
	UserAuths    []UserAuth    `gorm:"foreignKey:UserID" json:"user_auths,omitempty"`
	UserSessions []UserSession `gorm:"foreignKey:UserID" json:"user_sessions,omitempty"`
	UserTags     []UserTag     `gorm:"foreignKey:UserID" json:"user_tags,omitempty"`
}

// UserAuth 用户认证信息表
type UserAuth struct {
	ID     uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID uint64 `gorm:"not null;index" json:"user_id"`

	// 认证类型：real_name(实名认证)、bank_card(银行卡认证)、credit(征信认证)
	AuthType string `gorm:"type:varchar(20);not null" json:"auth_type"`

	// 认证状态：pending(待审核)、approved(通过)、rejected(拒绝)、expired(过期)
	AuthStatus string `gorm:"type:varchar(20);not null;default:'pending'" json:"auth_status"`

	// 认证数据(JSON格式存储)
	AuthData string `gorm:"type:json" json:"auth_data"`

	// 审核信息
	ReviewerID *uint64    `json:"reviewer_id"`
	ReviewNote string     `gorm:"type:text" json:"review_note"`
	ReviewedAt *time.Time `json:"reviewed_at"`

	// 有效期
	ExpiresAt *time.Time `json:"expires_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// UserSession 用户会话管理表
type UserSession struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64 `gorm:"not null;index" json:"user_id"`
	SessionID string `gorm:"type:varchar(64);uniqueIndex;not null" json:"session_id"`

	// 平台类型：app(移动应用)、web(网页)、oa(后台管理)
	Platform string `gorm:"type:varchar(10);not null" json:"platform"`

	// 设备信息
	DeviceID   string `gorm:"type:varchar(64)" json:"device_id"`
	DeviceType string `gorm:"type:varchar(20)" json:"device_type"` // ios, android, web
	DeviceName string `gorm:"type:varchar(500)" json:"device_name"`
	AppVersion string `gorm:"type:varchar(20)" json:"app_version"`
	UserAgent  string `gorm:"type:text" json:"user_agent"`

	// IP和地理位置
	IPAddress string `gorm:"type:varchar(45)" json:"ip_address"`
	Location  string `gorm:"type:varchar(100)" json:"location"`

	// JWT Token信息哈希 (不在JSON中返回)
	AccessToken      string     `gorm:"type:text" json:"-"`
	RefreshToken     string     `gorm:"type:text" json:"-"`
	AccessTokenHash  string     `gorm:"type:varchar(64);index" json:"-"`
	RefreshTokenHash string     `gorm:"type:varchar(64);index" json:"-"`
	TokenExpiresAt   *time.Time `json:"token_expires_at"`
	RefreshExpiresAt *time.Time `json:"refresh_expires_at"`

	// 会话状态：active(活跃)、expired(过期)、revoked(撤销)
	Status string `gorm:"type:varchar(20);not null;default:'active'" json:"status"`

	// 会话时间
	LoginTime    time.Time  `json:"login_time"`
	LastActiveAt time.Time  `json:"last_active_at"`
	LogoutTime   *time.Time `json:"logout_time"`

	// 登录方式
	LoginMethod string `gorm:"type:varchar(20)" json:"login_method"` // password, sms, oauth

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// UserTag 用户标签表
type UserTag struct {
	ID     uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID uint64 `gorm:"not null;index" json:"user_id"`

	// 标签类型：business(业务标签)、behavior(行为标签)、risk(风险标签)
	TagType  string `gorm:"type:varchar(20);not null" json:"tag_type"`
	TagKey   string `gorm:"type:varchar(50);not null" json:"tag_key"`
	TagValue string `gorm:"type:varchar(100)" json:"tag_value"`

	// 标签来源：manual(手动)、system(系统)、ai(AI标注)
	Source string `gorm:"type:varchar(20);not null;default:'system'" json:"source"`

	// 创建者(如果是手动标注)
	CreatorID *uint64 `json:"creator_id"`

	// 权重分数(0-100)
	Score int `gorm:"default:0" json:"score"`

	// 过期时间(可选)
	ExpiresAt *time.Time `json:"expires_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// OAUser OA后台管理用户表
type OAUser struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Phone    string `gorm:"type:varchar(20);index" json:"phone"`

	// 密码相关 (不在JSON中返回)
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Salt         string `gorm:"type:varchar(32);not null" json:"-"`

	// 基本信息
	RealName string `gorm:"type:varchar(50);not null" json:"real_name"`
	Avatar   string `gorm:"type:varchar(255)" json:"avatar"`

	// 角色权限
	RoleID     uint64 `gorm:"not null;index" json:"role_id"`
	Department string `gorm:"type:varchar(50)" json:"department"`
	Position   string `gorm:"type:varchar(50)" json:"position"`

	// 状态：active(正常)、frozen(冻结)、deleted(删除)
	Status string `gorm:"type:varchar(20);not null;default:'active'" json:"status"`

	// 登录信息
	LastLoginAt *time.Time `json:"last_login_at"`
	LastLoginIP string     `gorm:"type:varchar(45)" json:"last_login_ip"`
	LoginCount  uint32     `gorm:"default:0" json:"login_count"`

	// 两步验证
	TwoFactorEnabled bool   `gorm:"default:false" json:"two_factor_enabled"`
	TwoFactorSecret  string `gorm:"type:varchar(32)" json:"-"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Role OARole `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// OARole OA角色表
type OARole struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	DisplayName string `gorm:"type:varchar(100);not null" json:"display_name"`
	Description string `gorm:"type:text" json:"description"`

	// 权限配置(JSON格式)
	Permissions string `gorm:"type:json" json:"permissions"`

	// 是否为超级管理员角色
	IsSuper bool `gorm:"default:false" json:"is_super"`

	// 状态：active(正常)、inactive(禁用)
	Status string `gorm:"type:varchar(20);not null;default:'active'" json:"status"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// RealNameAuthData 实名认证数据结构
type RealNameAuthData struct {
	IDCardNumber   string `json:"id_card_number"`
	RealName       string `json:"real_name"`
	IDCardFrontImg string `json:"id_card_front_img"`
	IDCardBackImg  string `json:"id_card_back_img"`
	FaceVerifyImg  string `json:"face_verify_img"`
}

// BankCardAuthData 银行卡认证数据结构
type BankCardAuthData struct {
	BankCardNumber string `json:"bank_card_number"`
	BankName       string `json:"bank_name"`
	CardholderName string `json:"cardholder_name"`
}

// CreditAuthData 征信认证数据结构
type CreditAuthData struct {
	CreditScore     int    `json:"credit_score"`
	CreditLevel     string `json:"credit_level"`
	CreditReportUrl string `json:"credit_report_url"`
	ProviderName    string `json:"provider_name"`
}

// RolePermissions 角色权限配置结构
type RolePermissions struct {
	LoanManagement    []string `json:"loan_management"`    // ["view", "create", "update", "approve"]
	MachineManagement []string `json:"machine_management"` // ["view", "create", "update", "delete"]
	UserManagement    []string `json:"user_management"`    // ["view", "create", "update", "freeze"]
	ContentManagement []string `json:"content_management"` // ["view", "create", "update", "publish"]
	SystemSettings    []string `json:"system_settings"`    // ["view", "update"]
	DataAnalytics     []string `json:"data_analytics"`     // ["view", "export"]
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}

func (UserAuth) TableName() string {
	return "user_auths"
}

func (UserSession) TableName() string {
	return "user_sessions"
}

func (UserTag) TableName() string {
	return "user_tags"
}

func (OAUser) TableName() string {
	return "oa_users"
}

func (OARole) TableName() string {
	return "oa_roles"
}
