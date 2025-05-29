# 通用功能数据模型

## 文件概述

`common.go` 是数字惠农系统通用功能的数据模型文件，定义了内容管理、专家服务、系统配置、文件管理、离线队列和API日志等通用功能的数据结构。

## 核心数据模型

### 1. Article 文章表
内容管理系统的核心模型，支持多种文章类型和完整的发布流程。

```go
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
    
    // 作者信息
    AuthorID   uint64 `json:"author_id"`
    AuthorName string `json:"author_name"`
    
    // 发布状态：draft(草稿)、published(已发布)、archived(已归档)
    Status string `json:"status"`
    
    // 显示控制
    IsTop      bool `json:"is_top"`      // 是否置顶
    IsFeatured bool `json:"is_featured"` // 是否推荐
    
    // 统计信息
    ViewCount    int64 `json:"view_count"`
    LikeCount    int64 `json:"like_count"`
    ShareCount   int64 `json:"share_count"`
    CommentCount int64 `json:"comment_count"`
}
```

**文章分类 (Category)**:
- `policy`: 政策资讯 - 农业政策解读、补贴信息等
- `technology`: 技术指导 - 种植技术、养殖技术等
- `market`: 市场信息 - 农产品价格、市场行情等
- `news`: 行业新闻 - 农业行业动态、企业新闻等

**发布状态 (Status)**:
- `draft`: 草稿 - 正在编辑中的文章
- `published`: 已发布 - 正式发布的文章
- `archived`: 已归档 - 过期或下线的文章

### 2. Category 文章分类表
支持层级分类的分类管理系统。

```go
type Category struct {
    ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
    DisplayName string `gorm:"type:varchar(100);not null" json:"display_name"`
    Description string `gorm:"type:text" json:"description"`
    
    // 父分类ID(支持层级分类)
    ParentID *uint64 `json:"parent_id"`
    
    // 排序
    SortOrder int `json:"sort_order"`
    
    // 图标
    Icon string `json:"icon"`
    
    // 状态：active(启用)、inactive(禁用)
    Status string `json:"status"`
}
```

**层级分类示例**:
```
农业技术 (parent_id: null)
├── 种植技术 (parent_id: 1)
│   ├── 粮食作物 (parent_id: 2)
│   └── 经济作物 (parent_id: 2)
└── 养殖技术 (parent_id: 1)
    ├── 畜禽养殖 (parent_id: 5)
    └── 水产养殖 (parent_id: 5)
```

### 3. Expert 专家信息表
专家服务系统的核心模型，支持专家认证和服务管理。

```go
type Expert struct {
    ID           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    Name         string `gorm:"type:varchar(50);not null" json:"name"`
    Title        string `gorm:"type:varchar(100)" json:"title"`
    Organization string `gorm:"type:varchar(100)" json:"organization"`
    
    // 专业领域(JSON数组)
    Specialties string `gorm:"type:json" json:"specialties"`
    
    // 联系方式
    Phone  string `json:"phone"`
    Email  string `json:"email"`
    WeChat string `json:"wechat"`
    
    // 个人简介
    Biography string `json:"biography"`
    
    // 从业经验
    ExperienceYears int `json:"experience_years"`
    
    // 服务地区(JSON数组)
    ServiceAreas string `gorm:"type:json" json:"service_areas"`
    
    // 评分
    Rating      float32 `json:"rating"`
    RatingCount int     `json:"rating_count"`
    
    // 认证状态
    IsVerified bool       `json:"is_verified"`
    VerifiedAt *time.Time `json:"verified_at"`
}
```

**专业领域示例**:
```json
["种植技术", "植物保护", "土壤改良", "病虫害防治"]
```

**服务地区示例**:
```json
["山东省", "河南省", "江苏省"]
```

### 4. SystemConfig 系统配置表
系统运行参数的集中管理，支持不同类型的配置项。

```go
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
    IsEditable bool `json:"is_editable"`
    
    // 是否加密存储
    IsEncrypted bool `json:"is_encrypted"`
}
```

**配置分组 (ConfigGroup)**:
- `system`: 系统配置 - 基础系统参数
- `business`: 业务配置 - 业务规则参数
- `ui`: 界面配置 - 前端显示参数
- `integration`: 集成配置 - 第三方服务参数

### 5. FileUpload 文件上传管理表
统一的文件上传和管理系统，支持多种存储方式。

```go
type FileUpload struct {
    ID           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    FileName     string `gorm:"type:varchar(255);not null" json:"file_name"`
    OriginalName string `gorm:"type:varchar(255);not null" json:"original_name"`
    FilePath     string `gorm:"type:varchar(500);not null" json:"file_path"`
    FileURL      string `gorm:"type:varchar(500)" json:"file_url"`
    
    // 文件信息
    FileSize int64  `json:"file_size"`
    FileType string `json:"file_type"`
    MimeType string `json:"mime_type"`
    FileHash string `json:"file_hash"`
    
    // 上传信息
    UploaderID   uint64 `json:"uploader_id"`
    UploaderType string `json:"uploader_type"` // user、oa_user
    
    // 业务关联
    BusinessType string `json:"business_type"` // loan_application、machine、article等
    BusinessID   uint64 `json:"business_id"`
    
    // 存储信息
    StorageType string `json:"storage_type"` // local、oss、qiniu
    BucketName  string `json:"bucket_name"`
    
    // 访问控制
    IsPublic    bool  `json:"is_public"`
    AccessCount int64 `json:"access_count"`
}
```

**存储类型 (StorageType)**:
- `local`: 本地存储 - 文件存储在服务器本地
- `oss`: 阿里云OSS - 使用阿里云对象存储
- `qiniu`: 七牛云存储 - 使用七牛云对象存储

### 6. OfflineQueue 离线队列表
支持离线环境的操作队列管理，确保网络恢复后数据同步。

```go
type OfflineQueue struct {
    ID         uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID     uint64 `gorm:"not null;index" json:"user_id"`
    ActionType string `gorm:"type:varchar(50);not null" json:"action_type"` // create_loan_application、update_user_profile等
    
    // 请求数据(JSON格式)
    RequestData string `gorm:"type:json;not null" json:"request_data"`
    
    // 状态：pending(待处理)、processing(处理中)、completed(已完成)、failed(失败)
    Status string `json:"status"`
    
    // 重试信息
    RetryCount int        `json:"retry_count"`
    MaxRetries int        `json:"max_retries"`
    NextRetry  *time.Time `json:"next_retry"`
    
    // 错误信息
    ErrorMessage string `json:"error_message"`
    
    // 处理结果
    ResultData string `gorm:"type:json" json:"result_data"`
}
```

**动作类型 (ActionType)**:
- `create_loan_application`: 创建贷款申请
- `update_user_profile`: 更新用户信息
- `upload_file`: 上传文件
- `submit_rental_order`: 提交租赁订单

### 7. APILog API调用日志表
记录所有API调用的详细日志，用于监控、调试和审计。

```go
type APILog struct {
    ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    RequestID string `gorm:"type:varchar(100);index" json:"request_id"`
    
    // 请求信息
    Method    string `json:"method"`
    URL       string `json:"url"`
    Headers   string `json:"headers"`
    Query     string `json:"query"`
    Body      string `json:"body"`
    IPAddress string `json:"ip_address"`
    UserAgent string `json:"user_agent"`
    
    // 响应信息
    StatusCode   int    `json:"status_code"`
    ResponseBody string `json:"response_body"`
    ResponseTime int    `json:"response_time"` // 响应时间(毫秒)
    
    // 用户信息
    UserID   *uint64 `json:"user_id"`
    UserType string  `json:"user_type"` // user、oa_user、anonymous
    
    // 错误信息
    ErrorCode    string `json:"error_code"`
    ErrorMessage string `json:"error_message"`
}
```

## 业务功能设计

### 内容管理系统
支持完整的内容创建、编辑、发布和管理流程:

```go
// 发布文章
func PublishArticle(articleID uint64, publisherID uint64) error {
    var article model.Article
    if err := db.First(&article, articleID).Error; err != nil {
        return err
    }
    
    // 更新发布状态
    now := time.Now()
    article.Status = "published"
    article.PublishedAt = &now
    
    // 重置统计数据
    article.ViewCount = 0
    article.LikeCount = 0
    article.ShareCount = 0
    article.CommentCount = 0
    
    return db.Save(&article).Error
}

// 文章浏览计数
func IncrementViewCount(articleID uint64) error {
    return db.Model(&model.Article{}).Where("id = ?", articleID).
        UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}
```

### 专家服务系统
专家认证和服务管理功能:

```go
// 搜索专家
func SearchExperts(specialties []string, serviceArea string, page, limit int) ([]model.Expert, error) {
    var experts []model.Expert
    
    query := db.Where("is_verified = ? AND status = ?", true, "active")
    
    // 专业领域筛选
    if len(specialties) > 0 {
        for _, specialty := range specialties {
            query = query.Where("JSON_CONTAINS(specialties, ?)", fmt.Sprintf(`"%s"`, specialty))
        }
    }
    
    // 服务地区筛选
    if serviceArea != "" {
        query = query.Where("JSON_CONTAINS(service_areas, ?)", fmt.Sprintf(`"%s"`, serviceArea))
    }
    
    // 按评分排序
    query = query.Order("rating DESC, rating_count DESC")
    
    // 分页
    offset := (page - 1) * limit
    return experts, query.Offset(offset).Limit(limit).Find(&experts).Error
}
```

### 文件管理系统
统一的文件上传和管理:

```go
// 上传文件
func UploadFile(file multipart.File, header *multipart.FileHeader, uploaderID uint64, businessType string, businessID uint64) (*model.FileUpload, error) {
    // 生成文件名
    fileName := generateFileName(header.Filename)
    filePath := fmt.Sprintf("uploads/%s/%s", businessType, fileName)
    
    // 计算文件哈希
    fileHash, err := calculateFileHash(file)
    if err != nil {
        return nil, err
    }
    
    // 保存文件
    if err := saveFile(file, filePath); err != nil {
        return nil, err
    }
    
    // 创建文件记录
    fileUpload := &model.FileUpload{
        FileName:     fileName,
        OriginalName: header.Filename,
        FilePath:     filePath,
        FileURL:      fmt.Sprintf("/api/files/%s", fileName),
        FileSize:     header.Size,
        FileType:     getFileType(header.Filename),
        MimeType:     header.Header.Get("Content-Type"),
        FileHash:     fileHash,
        UploaderID:   uploaderID,
        UploaderType: "user",
        BusinessType: businessType,
        BusinessID:   businessID,
        StorageType:  "local",
        IsPublic:     false,
        AccessCount:  0,
    }
    
    return fileUpload, db.Create(fileUpload).Error
}
```

### 离线队列处理
离线环境下的操作队列管理:

```go
// 添加离线操作
func AddOfflineAction(userID uint64, actionType string, requestData interface{}) error {
    dataJSON, err := json.Marshal(requestData)
    if err != nil {
        return err
    }
    
    action := &model.OfflineQueue{
        UserID:      userID,
        ActionType:  actionType,
        RequestData: string(dataJSON),
        Status:      "pending",
        RetryCount:  0,
        MaxRetries:  3,
    }
    
    return db.Create(action).Error
}

// 处理离线队列
func ProcessOfflineQueue() error {
    var actions []model.OfflineQueue
    if err := db.Where("status = ?", "pending").
        Order("created_at ASC").
        Limit(100).
        Find(&actions).Error; err != nil {
        return err
    }
    
    for _, action := range actions {
        if err := processOfflineAction(&action); err != nil {
            // 处理失败，增加重试次数
            action.RetryCount++
            if action.RetryCount >= action.MaxRetries {
                action.Status = "failed"
                action.ErrorMessage = err.Error()
            } else {
                // 设置下次重试时间
                nextRetry := time.Now().Add(time.Duration(action.RetryCount*5) * time.Minute)
                action.NextRetry = &nextRetry
            }
            db.Save(&action)
        } else {
            // 处理成功
            action.Status = "completed"
            db.Save(&action)
        }
    }
    
    return nil
}
```

## 数据库索引设计

### 关键索引
1. **文章表 (articles)**:
   - `category` 普通索引
   - `status` 普通索引
   - `author_id` 普通索引
   - `published_at` 普通索引
   - 联合索引: (`category`, `status`)

2. **分类表 (categories)**:
   - `name` 唯一索引
   - `parent_id` 普通索引
   - `status` 普通索引

3. **专家表 (experts)**:
   - `is_verified` 普通索引
   - `status` 普通索引
   - 联合索引: (`is_verified`, `status`)

4. **系统配置表 (system_configs)**:
   - `config_key` 唯一索引
   - `config_group` 普通索引

5. **文件上传表 (file_uploads)**:
   - `file_hash` 普通索引
   - `uploader_id` 普通索引
   - `business_type` 普通索引
   - 联合索引: (`business_type`, `business_id`)

6. **离线队列表 (offline_queue)**:
   - `user_id` 普通索引
   - `status` 普通索引
   - `action_type` 普通索引

7. **API日志表 (api_logs)**:
   - `request_id` 普通索引
   - `user_id` 普通索引
   - `created_at` 普通索引

## 使用示例

### 内容管理
```go
// 创建文章
article := &model.Article{
    Title:      "春耕备播技术要点",
    Subtitle:   "确保春季播种质量的关键技术",
    Content:    "详细的春耕技术指导内容...",
    Summary:    "本文介绍了春耕备播的关键技术要点",
    Category:   "technology",
    Tags:       `["春耕", "播种", "技术指导"]`,
    AuthorID:   expertID,
    AuthorName: "农业专家张三",
    Status:     "draft",
    IsTop:      false,
    IsFeatured: true,
}

// SEO优化
article.SEOTitle = "春耕备播技术要点 - 数字惠农"
article.SEODescription = "专业的春耕备播技术指导，帮助农民朋友提高播种质量"
article.SEOKeywords = "春耕,播种,农业技术,种植指导"

if err := db.Create(article).Error; err != nil {
    return fmt.Errorf("创建文章失败: %w", err)
}
```

### 专家注册
```go
// 注册农业专家
expert := &model.Expert{
    Name:         "李农技",
    Title:        "高级农艺师",
    Organization: "市农业技术推广中心",
    Phone:        "13800138000",
    Email:        "linongji@example.com",
    Biography:    "从事农业技术推广工作20年，专注于粮食作物种植技术",
    ExperienceYears: 20,
    Status:       "active",
    IsVerified:   false,
}

// 设置专业领域
specialties := []string{"小麦种植", "玉米种植", "植物保护", "土壤改良"}
specialtiesJSON, _ := json.Marshal(specialties)
expert.Specialties = string(specialtiesJSON)

// 设置服务地区
serviceAreas := []string{"山东省", "河南省"}
serviceAreasJSON, _ := json.Marshal(serviceAreas)
expert.ServiceAreas = string(serviceAreasJSON)

if err := db.Create(expert).Error; err != nil {
    return fmt.Errorf("注册专家失败: %w", err)
}
```

### 系统配置管理
```go
// 设置系统配置
configs := []model.SystemConfig{
    {
        ConfigKey:   "loan.auto_approval_enabled",
        ConfigValue: "true",
        ConfigType:  "bool",
        ConfigGroup: "business",
        Description: "是否启用贷款自动审批功能",
        IsEditable:  true,
        IsEncrypted: false,
    },
    {
        ConfigKey:   "file.max_upload_size",
        ConfigValue: "10485760",
        ConfigType:  "int",
        ConfigGroup: "system",
        Description: "文件上传最大大小限制(字节)",
        IsEditable:  true,
        IsEncrypted: false,
    },
}

for _, config := range configs {
    if err := db.Create(&config).Error; err != nil {
        log.Printf("创建配置失败: %v", err)
    }
}
```

## 性能优化

### 内容检索优化
1. **全文搜索**: 对文章标题和内容建立全文索引
2. **标签检索**: 优化JSON字段的标签查询
3. **分类缓存**: 热门分类和文章缓存到Redis

### 文件存储优化
1. **CDN加速**: 静态文件使用CDN分发
2. **压缩存储**: 图片自动压缩和格式转换
3. **重复检测**: 基于文件哈希的重复文件检测

### 离线队列优化
1. **批量处理**: 批量处理离线队列任务
2. **优先级队列**: 支持不同优先级的任务处理
3. **失败重试**: 指数退避的重试策略

## 安全考虑

1. **内容审核**: 文章发布前的内容审核机制
2. **文件安全**: 文件类型验证和病毒扫描
3. **访问控制**: 基于角色的文件访问权限
4. **数据脱敏**: 日志中的敏感信息脱敏
5. **配置加密**: 敏感配置项的加密存储 