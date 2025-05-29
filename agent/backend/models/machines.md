# 农机租赁模块 - 数据模型设计文档

## 1. 模块概述

农机租赁模块是数字惠农系统的重要组成部分，提供农机设备注册、地理位置搜索、租赁预订、订单管理、支付结算等全流程服务。该模块支持"共享经济"模式，让农机所有者可以将闲置设备出租获得收益，农户可以就近租用所需设备。

### 主要功能特性
- 🚜 **设备管理**: 支持多种农机类型注册，包含详细技术参数和状态管理
- 📍 **地理搜索**: 基于GPS定位的就近农机搜索，支持距离筛选
- 📅 **智能预订**: 实时设备状态查询，支持预约和即时租赁
- 💳 **支付集成**: 集成第三方支付，支持押金、租金分离结算
- ⭐ **评价体系**: 双向评价机制，建立信用体系
- 📱 **实时追踪**: 设备使用状态实时监控和轨迹记录

## 2. 核心数据模型

### 2.1 machines - 农机设备表

```go
type Machine struct {
    ID              uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    MachineCode     string    `gorm:"type:varchar(30);uniqueIndex;not null" json:"machine_code"`
    MachineName     string    `gorm:"type:varchar(100);not null" json:"machine_name"`
    
    // 设备分类：tractor(拖拉机)、harvester(收割机)、planter(播种机)、sprayer(喷药机)
    MachineType     string    `gorm:"type:varchar(30);not null" json:"machine_type"`
    
    // 品牌和型号
    Brand           string    `gorm:"type:varchar(50);not null" json:"brand"`
    Model           string    `gorm:"type:varchar(50);not null" json:"model"`
    
    // 规格参数(JSON格式)
    Specifications  string    `gorm:"type:json" json:"specifications"`
    
    // 设备描述
    Description     string    `gorm:"type:text" json:"description"`
    
    // 设备图片(JSON数组)
    Images          string    `gorm:"type:json" json:"images"`
    
    // 所有者信息
    OwnerID         uint64    `gorm:"not null;index" json:"owner_id"`
    OwnerType       string    `gorm:"type:varchar(20);not null" json:"owner_type"` // individual(个人)、company(公司)
    
    // 位置信息
    Province        string    `gorm:"type:varchar(50);not null" json:"province"`
    City            string    `gorm:"type:varchar(50);not null" json:"city"`
    County          string    `gorm:"type:varchar(50);not null" json:"county"`
    Address         string    `gorm:"type:varchar(200)" json:"address"`
    Longitude       float64   `gorm:"type:decimal(10,6)" json:"longitude"`
    Latitude        float64   `gorm:"type:decimal(10,6)" json:"latitude"`
    
    // 租赁定价
    HourlyRate      int64     `json:"hourly_rate"`      // 小时租金(分)
    DailyRate       int64     `json:"daily_rate"`       // 日租金(分)
    PerAcreRate     int64     `json:"per_acre_rate"`    // 按亩收费(分)
    
    // 押金
    DepositAmount   int64     `gorm:"not null" json:"deposit_amount"`    // 押金金额(分)
    
    // 设备状态：available(可租)、rented(已租出)、maintenance(维护中)、offline(下线)
    Status          string    `gorm:"type:varchar(20);not null;default:'available'" json:"status"`
    
    // 可用时间(JSON格式)
    AvailableSchedule string  `gorm:"type:json" json:"available_schedule"`
    
    // 设备评分
    Rating          float32   `gorm:"type:decimal(3,1);default:0.0" json:"rating"`
    RatingCount     int       `gorm:"default:0" json:"rating_count"`
    
    // 租赁条件
    MinRentalHours  int       `gorm:"default:1" json:"min_rental_hours"`     // 最少租赁小时数
    MaxAdvanceDays  int       `gorm:"default:30" json:"max_advance_days"`    // 最多提前预订天数
    
    // 认证信息
    IsVerified      bool      `gorm:"default:false" json:"is_verified"`      // 是否已认证
    VerifiedAt      *time.Time `json:"verified_at"`
    
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
    
    // 关联
    Owner           User      `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

// 设备规格参数结构
type MachineSpecs struct {
    Power           string  `json:"power"`           // 功率
    WorkingWidth    string  `json:"working_width"`   // 工作幅宽
    Weight          string  `json:"weight"`          // 重量
    FuelType        string  `json:"fuel_type"`       // 燃料类型
    YearOfMake      int     `json:"year_of_make"`    // 制造年份
    WorkingSpeed    string  `json:"working_speed"`   // 工作速度
    Capacity        string  `json:"capacity"`        // 容量
    Other           map[string]string `json:"other"` // 其他参数
}

// 可用时间表结构
type AvailableSchedule struct {
    Monday    []TimeSlot `json:"monday"`
    Tuesday   []TimeSlot `json:"tuesday"`
    Wednesday []TimeSlot `json:"wednesday"`
    Thursday  []TimeSlot `json:"thursday"`
    Friday    []TimeSlot `json:"friday"`
    Saturday  []TimeSlot `json:"saturday"`
    Sunday    []TimeSlot `json:"sunday"`
}

type TimeSlot struct {
    StartTime string `json:"start_time"` // "08:00"
    EndTime   string `json:"end_time"`   // "18:00"
}
```

### 2.2 rental_orders - 租赁订单表

```go
type RentalOrder struct {
    ID              uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    OrderNo         string    `gorm:"type:varchar(30);uniqueIndex;not null" json:"order_no"`
    
    // 关联信息
    MachineID       uint64    `gorm:"not null;index" json:"machine_id"`
    RenterID        uint64    `gorm:"not null;index" json:"renter_id"`
    OwnerID         uint64    `gorm:"not null;index" json:"owner_id"`
    
    // 租赁时间
    StartTime       time.Time `gorm:"not null" json:"start_time"`
    EndTime         time.Time `gorm:"not null" json:"end_time"`
    RentalDuration  int       `json:"rental_duration"`      // 租赁时长(小时)
    
    // 租赁地点
    RentalLocation  string    `gorm:"type:varchar(200)" json:"rental_location"`
    ContactPerson   string    `gorm:"type:varchar(50)" json:"contact_person"`
    ContactPhone    string    `gorm:"type:varchar(20)" json:"contact_phone"`
    
    // 计费方式：hourly(按小时)、daily(按天)、per_acre(按亩)
    BillingMethod   string    `gorm:"type:varchar(20);not null" json:"billing_method"`
    
    // 费用计算
    UnitPrice       int64     `gorm:"not null" json:"unit_price"`       // 单价(分)
    Quantity        float64   `gorm:"not null" json:"quantity"`         // 数量(小时/天/亩)
    SubtotalAmount  int64     `gorm:"not null" json:"subtotal_amount"`  // 小计(分)
    DepositAmount   int64     `gorm:"not null" json:"deposit_amount"`   // 押金(分)
    TotalAmount     int64     `gorm:"not null" json:"total_amount"`     // 总金额(分)
    
    // 订单状态：pending(待确认)、confirmed(已确认)、paid(已支付)、in_progress(进行中)、
    // completed(已完成)、cancelled(已取消)、disputed(有争议)
    Status          string    `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
    
    // 支付信息
    PaymentMethod   string    `gorm:"type:varchar(20)" json:"payment_method"`
    PaymentStatus   string    `gorm:"type:varchar(20)" json:"payment_status"`
    PaidAmount      int64     `json:"paid_amount"`
    PaidAt          *time.Time `json:"paid_at"`
    
    // 确认信息
    OwnerConfirmedAt  *time.Time `json:"owner_confirmed_at"`
    RenterConfirmedAt *time.Time `json:"renter_confirmed_at"`
    
    // 服务评价
    RenterRating    float32   `gorm:"type:decimal(3,1)" json:"renter_rating"`
    RenterComment   string    `gorm:"type:text" json:"renter_comment"`
    OwnerRating     float32   `gorm:"type:decimal(3,1)" json:"owner_rating"`
    OwnerComment    string    `gorm:"type:text" json:"owner_comment"`
    
    // 备注
    Remarks         string    `gorm:"type:text" json:"remarks"`
    
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    
    // 关联
    Machine         Machine   `gorm:"foreignKey:MachineID" json:"machine,omitempty"`
    Renter          User      `gorm:"foreignKey:RenterID" json:"renter,omitempty"`
    Owner           User      `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}
```

## 3. 内容管理模块数据模型

### 3.1 articles - 文章表

```go
type Article struct {
    ID              uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    Title           string    `gorm:"type:varchar(200);not null" json:"title"`
    Subtitle        string    `gorm:"type:varchar(300)" json:"subtitle"`
    Content         string    `gorm:"type:longtext;not null" json:"content"`
    Summary         string    `gorm:"type:text" json:"summary"`
    
    // 分类：policy(政策)、technology(技术)、market(市场)、news(新闻)
    Category        string    `gorm:"type:varchar(30);not null" json:"category"`
    
    // 标签(JSON数组)
    Tags            string    `gorm:"type:json" json:"tags"`
    
    // 封面图片
    CoverImage      string    `gorm:"type:varchar(255)" json:"cover_image"`
    
    // 作者信息
    AuthorID        uint64    `gorm:"not null;index" json:"author_id"`
    AuthorName      string    `gorm:"type:varchar(50);not null" json:"author_name"`
    
    // 发布状态：draft(草稿)、published(已发布)、archived(已归档)
    Status          string    `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`
    
    // 显示控制
    IsTop           bool      `gorm:"default:false" json:"is_top"`        // 是否置顶
    IsFeatured      bool      `gorm:"default:false" json:"is_featured"`   // 是否推荐
    
    // 统计信息
    ViewCount       int64     `gorm:"default:0" json:"view_count"`
    LikeCount       int64     `gorm:"default:0" json:"like_count"`
    ShareCount      int64     `gorm:"default:0" json:"share_count"`
    CommentCount    int64     `gorm:"default:0" json:"comment_count"`
    
    // SEO优化
    SEOTitle        string    `gorm:"type:varchar(200)" json:"seo_title"`
    SEODescription  string    `gorm:"type:text" json:"seo_description"`
    SEOKeywords     string    `gorm:"type:varchar(500)" json:"seo_keywords"`
    
    // 发布时间
    PublishedAt     *time.Time `json:"published_at"`
    
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
    
    // 关联
    Author          OAUser    `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}
```

### 3.2 categories - 文章分类表

```go
type Category struct {
    ID              uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    Name            string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
    DisplayName     string    `gorm:"type:varchar(100);not null" json:"display_name"`
    Description     string    `gorm:"type:text" json:"description"`
    
    // 父分类ID(支持层级分类)
    ParentID        *uint64   `json:"parent_id"`
    
    // 排序
    SortOrder       int       `gorm:"default:0" json:"sort_order"`
    
    // 图标
    Icon            string    `gorm:"type:varchar(100)" json:"icon"`
    
    // 状态：active(启用)、inactive(禁用)
    Status          string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
    
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    
    // 关联
    Parent          *Category `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
    Children        []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}
```

### 3.3 experts - 专家信息表

```go
type Expert struct {
    ID              uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    Name            string    `gorm:"type:varchar(50);not null" json:"name"`
    Title           string    `gorm:"type:varchar(100)" json:"title"`
    Organization    string    `gorm:"type:varchar(100)" json:"organization"`
    
    // 专业领域(JSON数组)
    Specialties     string    `gorm:"type:json" json:"specialties"`
    
    // 联系方式
    Phone           string    `gorm:"type:varchar(20)" json:"phone"`
    Email           string    `gorm:"type:varchar(100)" json:"email"`
    WeChat          string    `gorm:"type:varchar(50)" json:"wechat"`
    
    // 头像
    Avatar          string    `gorm:"type:varchar(255)" json:"avatar"`
    
    // 个人简介
    Biography       string    `gorm:"type:text" json:"biography"`
    
    // 从业经验
    ExperienceYears int       `json:"experience_years"`
    
    // 服务地区(JSON数组)
    ServiceAreas    string    `gorm:"type:json" json:"service_areas"`
    
    // 评分
    Rating          float32   `gorm:"type:decimal(3,1);default:0.0" json:"rating"`
    RatingCount     int       `gorm:"default:0" json:"rating_count"`
    
    // 状态：active(活跃)、inactive(不活跃)
    Status          string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
    
    // 认证状态
    IsVerified      bool      `gorm:"default:false" json:"is_verified"`
    VerifiedAt      *time.Time `json:"verified_at"`
    
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
```

## 4. 系统配置模块数据模型

### 4.1 system_configs - 系统配置表

```go
type SystemConfig struct {
    ID              uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    ConfigKey       string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"config_key"`
    ConfigValue     string    `gorm:"type:text;not null" json:"config_value"`
    ConfigType      string    `gorm:"type:varchar(20);not null" json:"config_type"` // string、int、float、bool、json
    
    // 配置分组：system(系统)、business(业务)、ui(界面)、integration(集成)
    ConfigGroup     string    `gorm:"type:varchar(30);not null" json:"config_group"`
    
    // 配置描述
    Description     string    `gorm:"type:text" json:"description"`
    
    // 是否可通过界面修改
    IsEditable      bool      `gorm:"default:true" json:"is_editable"`
    
    // 是否加密存储
    IsEncrypted     bool      `gorm:"default:false" json:"is_encrypted"`
    
    // 最后修改人
    UpdatedBy       uint64    `json:"updated_by"`
    
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

// 常用配置示例
/*
系统配置示例：
- app.name: "数字惠农"
- app.version: "1.0.0"
- sms.provider: "aliyun"
- file.max_size: "10485760"
- loan.max_amount: "1000000"
- dify.api_url: "https://api.dify.ai"
- dify.api_key: "encrypted_key"
*/
```

### 4.2 file_uploads - 文件上传管理表

```go
type FileUpload struct {
    ID              uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    FileName        string    `gorm:"type:varchar(255);not null" json:"file_name"`
    OriginalName    string    `gorm:"type:varchar(255);not null" json:"original_name"`
    FilePath        string    `gorm:"type:varchar(500);not null" json:"file_path"`
    FileURL         string    `gorm:"type:varchar(500)" json:"file_url"`
    
    // 文件信息
    FileSize        int64     `gorm:"not null" json:"file_size"`
    FileType        string    `gorm:"type:varchar(50);not null" json:"file_type"`
    MimeType        string    `gorm:"type:varchar(100);not null" json:"mime_type"`
    FileHash        string    `gorm:"type:varchar(64);index" json:"file_hash"`
    
    // 上传信息
    UploaderID      uint64    `gorm:"not null;index" json:"uploader_id"`
    UploaderType    string    `gorm:"type:varchar(20);not null" json:"uploader_type"` // user、oa_user
    
    // 业务关联
    BusinessType    string    `gorm:"type:varchar(30)" json:"business_type"` // loan_application、machine、article等
    BusinessID      uint64    `json:"business_id"`
    
    // 存储信息
    StorageType     string    `gorm:"type:varchar(20);not null;default:'local'" json:"storage_type"` // local、oss、qiniu
    BucketName      string    `gorm:"type:varchar(100)" json:"bucket_name"`
    
    // 状态：uploaded(已上传)、processing(处理中)、processed(已处理)、failed(失败)
    Status          string    `gorm:"type:varchar(20);not null;default:'uploaded'" json:"status"`
    
    // 访问控制
    IsPublic        bool      `gorm:"default:false" json:"is_public"`
    AccessCount     int64     `gorm:"default:0" json:"access_count"`
    
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
```

## 5. Repository接口设计

### 5.1 MachineRepository接口

```go
type MachineRepository interface {
    // 基本CRUD
    Create(ctx context.Context, machine *Machine) error
    GetByID(ctx context.Context, id uint64) (*Machine, error)
    Update(ctx context.Context, machine *Machine) error
    Delete(ctx context.Context, id uint64) error
    
    // 查询方法
    List(ctx context.Context, req *ListMachinesRequest) (*ListMachinesResponse, error)
    GetByOwnerID(ctx context.Context, ownerID uint64) ([]*Machine, error)
    SearchByLocation(ctx context.Context, province, city, county string) ([]*Machine, error)
    GetAvailableMachines(ctx context.Context, machineType string, startTime, endTime time.Time) ([]*Machine, error)
    
    // 统计方法
    GetMachineCount(ctx context.Context) (int64, error)
    GetMachineCountByType(ctx context.Context, machineType string) (int64, error)
}
```

### 5.2 RentalOrderRepository接口

```go
type RentalOrderRepository interface {
    // 基本CRUD
    Create(ctx context.Context, order *RentalOrder) error
    GetByID(ctx context.Context, id uint64) (*RentalOrder, error)
    GetByOrderNo(ctx context.Context, orderNo string) (*RentalOrder, error)
    Update(ctx context.Context, order *RentalOrder) error
    UpdateStatus(ctx context.Context, id uint64, status string) error
    
    // 查询方法
    GetByRenterID(ctx context.Context, renterID uint64, limit, offset int) ([]*RentalOrder, error)
    GetByOwnerID(ctx context.Context, ownerID uint64, limit, offset int) ([]*RentalOrder, error)
    GetByMachineID(ctx context.Context, machineID uint64) ([]*RentalOrder, error)
    GetByStatus(ctx context.Context, status string, limit, offset int) ([]*RentalOrder, error)
    
    // 统计方法
    GetOrderCount(ctx context.Context) (int64, error)
    GetRevenueStats(ctx context.Context, days int) (*RevenueStats, error)
}
```

### 5.3 ArticleRepository接口

```go
type ArticleRepository interface {
    // 基本CRUD
    Create(ctx context.Context, article *Article) error
    GetByID(ctx context.Context, id uint64) (*Article, error)
    Update(ctx context.Context, article *Article) error
    Delete(ctx context.Context, id uint64) error
    
    // 查询方法
    List(ctx context.Context, req *ListArticlesRequest) (*ListArticlesResponse, error)
    GetByCategory(ctx context.Context, category string, limit, offset int) ([]*Article, error)
    GetFeaturedArticles(ctx context.Context, limit int) ([]*Article, error)
    Search(ctx context.Context, keyword string, limit, offset int) ([]*Article, error)
    
    // 统计方法
    IncrementViewCount(ctx context.Context, id uint64) error
    GetPopularArticles(ctx context.Context, days, limit int) ([]*Article, error)
}
```

## 6. 缓存策略设计

```go
// 缓存键名设计
const (
    MachineCachePrefix          = "machine:"            // machine:12345
    MachineListCachePrefix      = "machines:"           // machines:tractor:beijing
    RentalOrderCachePrefix      = "rental_order:"       // rental_order:12345
    ArticleCachePrefix          = "article:"            // article:12345
    ArticleListCachePrefix      = "articles:"           // articles:policy:1
    CategoryCachePrefix         = "category:"           // category:12345
    ExpertCachePrefix           = "expert:"             // expert:12345
    SystemConfigCachePrefix     = "config:"             // config:app.name
)

// 缓存时间配置
const (
    MachineCacheTTL         = 4 * time.Hour       // 农机信息缓存4小时
    RentalOrderCacheTTL     = 2 * time.Hour       // 租赁订单缓存2小时
    ArticleCacheTTL         = 12 * time.Hour      // 文章缓存12小时
    CategoryCacheTTL        = 24 * time.Hour      // 分类缓存24小时
    ExpertCacheTTL          = 6 * time.Hour       // 专家信息缓存6小时
    SystemConfigCacheTTL    = 1 * time.Hour       // 系统配置缓存1小时
)
```

## 7. 总结

农机租赁及其他模块的数据模型设计具有以下特点：

1. **农机租赁全流程**：从设备管理到订单处理的完整业务流程
2. **灵活的定价模式**：支持按小时、按天、按亩等多种计费方式
3. **内容管理体系**：完整的文章发布和分类管理系统
4. **专家服务支持**：专家信息管理和评价体系
5. **系统配置管理**：灵活的配置项管理，支持动态配置
6. **文件管理系统**：统一的文件上传和存储管理
7. **地理位置支持**：基于地理位置的服务匹配
8. **评价体系**：完整的用户评价和信用体系

这些模块的设计充分考虑了农村地区的实际需求，提供了完整的数字化农业服务体系。 