package model

import (
	"time"

	"gorm.io/gorm"
)

// Article 文章表
type Article struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Title    string `gorm:"type:varchar(200);not null" json:"title"`
	Subtitle string `gorm:"type:varchar(300)" json:"subtitle"`
	Content  string `gorm:"type:longtext;not null" json:"content"`
	Summary  string `gorm:"type:text" json:"summary"`

	// 分类：policy(政策)、technology(技术)、market(市场)、news(新闻)
	Category string `gorm:"type:varchar(30);not null" json:"category"`

	// 标签(JSON数组)
	Tags string `gorm:"type:json" json:"tags"`

	// 封面图片
	CoverImage string `gorm:"type:varchar(255)" json:"cover_image"`

	// 作者信息
	AuthorID   uint64 `gorm:"not null;index" json:"author_id"`
	AuthorName string `gorm:"type:varchar(50);not null" json:"author_name"`

	// 发布状态：draft(草稿)、published(已发布)、archived(已归档)
	Status string `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`

	// 显示控制
	IsTop      bool `gorm:"default:false" json:"is_top"`      // 是否置顶
	IsFeatured bool `gorm:"default:false" json:"is_featured"` // 是否推荐

	// 统计信息
	ViewCount    int64 `gorm:"default:0" json:"view_count"`
	LikeCount    int64 `gorm:"default:0" json:"like_count"`
	ShareCount   int64 `gorm:"default:0" json:"share_count"`
	CommentCount int64 `gorm:"default:0" json:"comment_count"`

	// SEO优化
	SEOTitle       string `gorm:"type:varchar(200)" json:"seo_title"`
	SEODescription string `gorm:"type:text" json:"seo_description"`
	SEOKeywords    string `gorm:"type:varchar(500)" json:"seo_keywords"`

	// 发布时间
	PublishedAt *time.Time `json:"published_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Author OAUser `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}

// Category 文章分类表
type Category struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	DisplayName string `gorm:"type:varchar(100);not null" json:"display_name"`
	Description string `gorm:"type:text" json:"description"`

	// 父分类ID(支持层级分类)
	ParentID *uint64 `json:"parent_id"`

	// 排序
	SortOrder int `gorm:"default:0" json:"sort_order"`

	// 图标
	Icon string `gorm:"type:varchar(100)" json:"icon"`

	// 状态：active(启用)、inactive(禁用)
	Status string `gorm:"type:varchar(20);not null;default:'active'" json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// Expert 专家信息表
type Expert struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string `gorm:"type:varchar(50);not null" json:"name"`
	Title        string `gorm:"type:varchar(100)" json:"title"`
	Organization string `gorm:"type:varchar(100)" json:"organization"`

	// 专业领域(JSON数组)
	Specialties string `gorm:"type:json" json:"specialties"`

	// 联系方式
	Phone  string `gorm:"type:varchar(20)" json:"phone"`
	Email  string `gorm:"type:varchar(100)" json:"email"`
	WeChat string `gorm:"type:varchar(50)" json:"wechat"`

	// 头像
	Avatar string `gorm:"type:varchar(255)" json:"avatar"`

	// 个人简介
	Biography string `gorm:"type:text" json:"biography"`

	// 从业经验
	ExperienceYears int `json:"experience_years"`

	// 服务地区(JSON数组)
	ServiceAreas string `gorm:"type:json" json:"service_areas"`

	// 评分
	Rating      float32 `gorm:"type:decimal(3,1);default:0.0" json:"rating"`
	RatingCount int     `gorm:"default:0" json:"rating_count"`

	// 状态：active(活跃)、inactive(不活跃)
	Status string `gorm:"type:varchar(20);not null;default:'active'" json:"status"`

	// 认证状态
	IsVerified bool       `gorm:"default:false" json:"is_verified"`
	VerifiedAt *time.Time `json:"verified_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// SystemConfig 系统配置表
type SystemConfig struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigKey   string `gorm:"type:varchar(100);uniqueIndex;not null" json:"config_key"`
	ConfigValue string `gorm:"type:text;not null" json:"config_value"`
	ConfigType  string `gorm:"type:varchar(20);not null" json:"config_type"` // string、int、float、bool、json

	// 配置分组：system(系统)、business(业务)、ui(界面)、integration(集成)
	ConfigGroup string `gorm:"type:varchar(30);not null" json:"config_group"`

	// 配置描述
	Description string `gorm:"type:text" json:"description"`

	// 是否可通过界面修改
	IsEditable bool `gorm:"default:true" json:"is_editable"`

	// 是否加密存储
	IsEncrypted bool `gorm:"default:false" json:"is_encrypted"`

	// 最后修改人
	UpdatedBy uint64 `json:"updated_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FileUpload 文件上传管理表
type FileUpload struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	FileName     string `gorm:"type:varchar(255);not null" json:"file_name"`
	OriginalName string `gorm:"type:varchar(255);not null" json:"original_name"`
	FilePath     string `gorm:"type:varchar(500);not null" json:"file_path"`
	FileURL      string `gorm:"type:varchar(500)" json:"file_url"`

	// 文件信息
	FileSize int64  `gorm:"not null" json:"file_size"`
	FileType string `gorm:"type:varchar(50);not null" json:"file_type"`
	MimeType string `gorm:"type:varchar(100);not null" json:"mime_type"`
	FileHash string `gorm:"type:varchar(64);index" json:"file_hash"`

	// 上传信息
	UploaderID   uint64 `gorm:"not null;index" json:"uploader_id"`
	UploaderType string `gorm:"type:varchar(20);not null" json:"uploader_type"` // user、oa_user

	// 业务关联
	BusinessType string `gorm:"type:varchar(30)" json:"business_type"` // loan_application、machine、article等
	BusinessID   uint64 `json:"business_id"`

	// 存储信息
	StorageType string `gorm:"type:varchar(20);not null;default:'local'" json:"storage_type"` // local、oss、qiniu
	BucketName  string `gorm:"type:varchar(100)" json:"bucket_name"`

	// 状态：uploaded(已上传)、processing(处理中)、processed(已处理)、failed(失败)
	Status string `gorm:"type:varchar(20);not null;default:'uploaded'" json:"status"`

	// 访问控制
	IsPublic    bool  `gorm:"default:false" json:"is_public"`
	AccessCount int64 `gorm:"default:0" json:"access_count"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// OfflineQueue 离线队列表
type OfflineQueue struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint64 `gorm:"not null;index" json:"user_id"`
	ActionType string `gorm:"type:varchar(50);not null" json:"action_type"` // create_loan_application、update_user_profile等

	// 请求数据(JSON格式)
	RequestData string `gorm:"type:json;not null" json:"request_data"`

	// 状态：pending(待处理)、processing(处理中)、completed(已完成)、failed(失败)
	Status string `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`

	// 重试信息
	RetryCount int        `gorm:"default:0" json:"retry_count"`
	MaxRetries int        `gorm:"default:3" json:"max_retries"`
	NextRetry  *time.Time `json:"next_retry"`

	// 错误信息
	ErrorMessage string `gorm:"type:text" json:"error_message"`

	// 处理结果
	ResultData string `gorm:"type:json" json:"result_data"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// APILog API调用日志表
type APILog struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	RequestID string `gorm:"type:varchar(100);index" json:"request_id"`

	// 请求信息
	Method    string `gorm:"type:varchar(10);not null" json:"method"`
	URL       string `gorm:"type:varchar(500);not null" json:"url"`
	Headers   string `gorm:"type:json" json:"headers"`
	Query     string `gorm:"type:json" json:"query"`
	Body      string `gorm:"type:text" json:"body"`
	IPAddress string `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent string `gorm:"type:varchar(500)" json:"user_agent"`

	// 响应信息
	StatusCode   int    `json:"status_code"`
	ResponseBody string `gorm:"type:text" json:"response_body"`
	ResponseTime int    `json:"response_time"` // 响应时间(毫秒)

	// 用户信息
	UserID   *uint64 `json:"user_id"`
	UserType string  `gorm:"type:varchar(20)" json:"user_type"` // user、oa_user、anonymous

	// 错误信息
	ErrorCode    string `gorm:"type:varchar(50)" json:"error_code"`
	ErrorMessage string `gorm:"type:text" json:"error_message"`

	CreatedAt time.Time `json:"created_at"`
}

// TableName 设置表名
func (Article) TableName() string {
	return "articles"
}

func (Category) TableName() string {
	return "categories"
}

func (Expert) TableName() string {
	return "experts"
}

func (SystemConfig) TableName() string {
	return "system_configs"
}

func (FileUpload) TableName() string {
	return "file_uploads"
}

func (OfflineQueue) TableName() string {
	return "offline_queue"
}

func (APILog) TableName() string {
	return "api_logs"
}
