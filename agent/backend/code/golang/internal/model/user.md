# 用户管理数据模型

## 文件概述

`user.go` 是数字惠农系统用户管理的核心数据模型文件，定义了用户、认证、会话管理、用户标签以及OA后台管理相关的所有数据结构。

## 核心数据模型

### 1. User 核心用户表
数字惠农APP的核心用户模型，支持多种用户类型和完整的用户信息管理。

```go
type User struct {
    ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    UUID     string `gorm:"type:varchar(36);uniqueIndex;not null" json:"uuid"`
    Username string `gorm:"type:varchar(50);uniqueIndex" json:"username"`
    Phone    string `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`
    Email    string `gorm:"type:varchar(100);index" json:"email"`
    
    // 密码相关 (不在JSON中返回)
    PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
    Salt         string `gorm:"type:varchar(32);not null" json:"-"`
    
    // 用户类型和状态
    UserType string `gorm:"type:varchar(20);not null;default:'farmer'" json:"user_type"`
    Status   string `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
    
    // 基本信息和地址信息
    // 认证状态
    // 登录信息
    // 时间字段
    // 关联数据
}
```

**用户类型 (UserType)**:
- `farmer`: 个体农户
- `farm_owner`: 农场主
- `cooperative`: 合作社
- `enterprise`: 企业

**用户状态 (Status)**:
- `active`: 正常
- `frozen`: 冻结
- `deleted`: 删除

**认证状态字段**:
- `IsRealNameVerified`: 实名认证状态
- `IsBankCardVerified`: 银行卡认证状态
- `IsCreditVerified`: 征信认证状态

### 2. UserAuth 用户认证信息表
管理用户的各种认证信息，支持多种认证类型和完整的审核流程。

```go
type UserAuth struct {
    ID     uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID uint64 `gorm:"not null;index" json:"user_id"`
    
    // 认证类型和状态
    AuthType   string `gorm:"type:varchar(20);not null" json:"auth_type"`
    AuthStatus string `gorm:"type:varchar(20);not null;default:'pending'" json:"auth_status"`
    
    // 认证数据(JSON格式存储)
    AuthData string `gorm:"type:json" json:"auth_data"`
    
    // 审核信息
    ReviewerID *uint64    `json:"reviewer_id"`
    ReviewNote string     `gorm:"type:text" json:"review_note"`
    ReviewedAt *time.Time `json:"reviewed_at"`
    
    // 有效期
    ExpiresAt *time.Time `json:"expires_at"`
}
```

**认证类型 (AuthType)**:
- `real_name`: 实名认证
- `bank_card`: 银行卡认证
- `credit`: 征信认证

**认证状态 (AuthStatus)**:
- `pending`: 待审核
- `approved`: 通过
- `rejected`: 拒绝
- `expired`: 过期

### 3. UserSession 用户会话管理表
管理用户在不同平台的登录会话，支持JWT令牌管理和设备信息记录。

```go
type UserSession struct {
    ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID    uint64 `gorm:"not null;index" json:"user_id"`
    SessionID string `gorm:"type:varchar(64);uniqueIndex;not null" json:"session_id"`
    
    // 平台类型：app(移动应用)、web(网页)、oa(后台管理)
    Platform string `gorm:"type:varchar(10);not null" json:"platform"`
    
    // 设备信息
    DeviceID   string `gorm:"type:varchar(64)" json:"device_id"`
    DeviceType string `gorm:"type:varchar(20)" json:"device_type"` // ios, android, web
    AppVersion string `gorm:"type:varchar(20)" json:"app_version"`
    
    // JWT Token信息 (不在JSON中返回)
    AccessToken    string     `gorm:"type:text" json:"-"`
    RefreshToken   string     `gorm:"type:text" json:"-"`
    TokenExpiresAt *time.Time `json:"token_expires_at"`
    
    // 会话状态：active(活跃)、expired(过期)、revoked(撤销)
    Status string `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
}
```

**平台类型 (Platform)**:
- `app`: 移动应用
- `web`: 网页端
- `oa`: 后台管理

**设备类型 (DeviceType)**:
- `ios`: iOS设备
- `android`: Android设备
- `web`: 网页浏览器

### 4. UserTag 用户标签表
用于用户画像和业务分析的标签系统，支持多种标签来源和权重管理。

```go
type UserTag struct {
    ID     uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID uint64 `gorm:"not null;index" json:"user_id"`
    
    // 标签类型：business(业务标签)、behavior(行为标签)、risk(风险标签)
    TagType  string `gorm:"type:varchar(20);not null" json:"tag_type"`
    TagKey   string `gorm:"type:varchar(50);not null" json:"tag_key"`
    TagValue string `gorm:"type:varchar(100)" json:"tag_value"`
    
    // 标签来源：manual(手动)、system(系统)、ai(AI标注)
    Source string `gorm:"type:varchar(20);not null;default:'system'" json:"source"`
    
    // 权重分数(0-100)
    Score int `gorm:"default:0" json:"score"`
    
    // 过期时间(可选)
    ExpiresAt *time.Time `json:"expires_at"`
}
```

**标签类型 (TagType)**:
- `business`: 业务标签（贷款偏好、农机需求等）
- `behavior`: 行为标签（活跃度、使用习惯等）
- `risk`: 风险标签（信用评级、违约风险等）

**标签来源 (Source)**:
- `manual`: 手动标注
- `system`: 系统自动标注
- `ai`: AI算法标注

### 5. OAUser OA后台管理用户表
后台管理系统的用户模型，支持角色权限管理和二次验证。

```go
type OAUser struct {
    ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
    Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
    Phone    string `gorm:"type:varchar(20);index" json:"phone"`
    
    // 密码相关 (不在JSON中返回)
    PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
    Salt         string `gorm:"type:varchar(32);not null" json:"-"`
    
    // 角色权限
    RoleID     uint64 `gorm:"not null;index" json:"role_id"`
    Department string `gorm:"type:varchar(50)" json:"department"`
    Position   string `gorm:"type:varchar(50)" json:"position"`
    
    // 两步验证
    TwoFactorEnabled bool   `gorm:"default:false" json:"two_factor_enabled"`
    TwoFactorSecret  string `gorm:"type:varchar(32)" json:"-"`
}
```

### 6. OARole OA角色表
后台管理系统的角色权限模型，支持灵活的权限配置。

```go
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
}
```

## 认证数据结构

### RealNameAuthData 实名认证数据
```go
type RealNameAuthData struct {
    IDCardNumber   string `json:"id_card_number"`   // 身份证号
    RealName       string `json:"real_name"`        // 真实姓名
    IDCardFrontImg string `json:"id_card_front_img"` // 身份证正面照
    IDCardBackImg  string `json:"id_card_back_img"`  // 身份证背面照
    FaceVerifyImg  string `json:"face_verify_img"`   // 人脸识别照片
}
```

### BankCardAuthData 银行卡认证数据
```go
type BankCardAuthData struct {
    BankCardNumber string `json:"bank_card_number"` // 银行卡号
    BankName       string `json:"bank_name"`        // 银行名称
    CardholderName string `json:"cardholder_name"`  // 持卡人姓名
}
```

### CreditAuthData 征信认证数据
```go
type CreditAuthData struct {
    CreditScore     int    `json:"credit_score"`      // 信用分数
    CreditLevel     string `json:"credit_level"`      // 信用等级
    CreditReportUrl string `json:"credit_report_url"` // 征信报告链接
    ProviderName    string `json:"provider_name"`     // 征信提供商
}
```

## 权限配置结构

### RolePermissions 角色权限配置
```go
type RolePermissions struct {
    LoanManagement    []string `json:"loan_management"`    // 贷款管理权限
    MachineManagement []string `json:"machine_management"` // 农机管理权限  
    UserManagement    []string `json:"user_management"`    // 用户管理权限
    ContentManagement []string `json:"content_management"` // 内容管理权限
    SystemSettings    []string `json:"system_settings"`    // 系统设置权限
    DataAnalytics     []string `json:"data_analytics"`     // 数据分析权限
}
```

**权限操作类型**:
- `view`: 查看
- `create`: 创建
- `update`: 更新
- `delete`: 删除
- `approve`: 审批
- `freeze`: 冻结
- `publish`: 发布
- `export`: 导出

## 数据库索引设计

### 关键索引
1. **用户表 (users)**:
   - `uuid` 唯一索引
   - `phone` 唯一索引
   - `username` 唯一索引
   - `email` 普通索引
   - `id_card` 普通索引

2. **用户认证表 (user_auths)**:
   - `user_id` 普通索引
   - 联合索引: (`user_id`, `auth_type`)

3. **用户会话表 (user_sessions)**:
   - `user_id` 普通索引
   - `session_id` 唯一索引

4. **用户标签表 (user_tags)**:
   - `user_id` 普通索引
   - 联合索引: (`user_id`, `tag_type`)

## 使用示例

### 用户创建和认证
```go
// 创建新用户
user := &model.User{
    UUID:     uuid.New().String(),
    Username: "farmer001",
    Phone:    "13800138000",
    Email:    "farmer001@example.com",
    UserType: "farmer",
    Status:   "active",
    // ... 其他字段
}

// 保存用户
if err := db.Create(user).Error; err != nil {
    return fmt.Errorf("创建用户失败: %w", err)
}

// 创建实名认证记录
authData := model.RealNameAuthData{
    IDCardNumber:   "110101199001011234",
    RealName:       "张三",
    IDCardFrontImg: "/uploads/id_front.jpg",
    IDCardBackImg:  "/uploads/id_back.jpg",
    FaceVerifyImg:  "/uploads/face.jpg",
}

authDataJSON, _ := json.Marshal(authData)
userAuth := &model.UserAuth{
    UserID:     user.ID,
    AuthType:   "real_name",
    AuthStatus: "pending",
    AuthData:   string(authDataJSON),
}

if err := db.Create(userAuth).Error; err != nil {
    return fmt.Errorf("创建认证记录失败: %w", err)
}
```

### 用户会话管理
```go
// 创建用户会话
session := &model.UserSession{
    UserID:         user.ID,
    SessionID:      generateSessionID(),
    Platform:       "app",
    DeviceID:       "device_12345",
    DeviceType:     "android",
    AppVersion:     "1.0.0",
    IPAddress:      "192.168.1.100",
    AccessToken:    accessToken,
    RefreshToken:   refreshToken,
    TokenExpiresAt: &expiresAt,
    Status:         "active",
}

if err := db.Create(session).Error; err != nil {
    return fmt.Errorf("创建会话失败: %w", err)
}

// 查询用户活跃会话
var activeSessions []model.UserSession
db.Where("user_id = ? AND status = ?", userID, "active").Find(&activeSessions)
```

### 用户标签管理
```go
// 添加用户标签
tags := []model.UserTag{
    {
        UserID:   user.ID,
        TagType:  "business",
        TagKey:   "loan_preference",
        TagValue: "agricultural_material",
        Source:   "system",
        Score:    80,
    },
    {
        UserID:   user.ID,
        TagType:  "behavior",
        TagKey:   "activity_level",
        TagValue: "high",
        Source:   "ai",
        Score:    90,
    },
}

for _, tag := range tags {
    if err := db.Create(&tag).Error; err != nil {
        log.Printf("创建标签失败: %v", err)
    }
}

// 查询用户标签
var userTags []model.UserTag
db.Where("user_id = ? AND tag_type = ?", userID, "business").Find(&userTags)
```

### OA用户和角色管理
```go
// 创建角色权限
permissions := model.RolePermissions{
    LoanManagement:    []string{"view", "create", "update", "approve"},
    MachineManagement: []string{"view", "create", "update"},
    UserManagement:    []string{"view"},
    ContentManagement: []string{"view", "create", "update", "publish"},
    SystemSettings:    []string{"view"},
    DataAnalytics:     []string{"view", "export"},
}

permissionsJSON, _ := json.Marshal(permissions)
role := &model.OARole{
    Name:        "loan_officer",
    DisplayName: "信贷员",
    Description: "负责贷款业务审核的信贷员角色",
    Permissions: string(permissionsJSON),
    IsSuper:     false,
    Status:      "active",
}

// 创建OA用户
oaUser := &model.OAUser{
    Username:     "loan001",
    Email:        "loan001@huinong.com",
    PasswordHash: hashedPassword,
    Salt:         salt,
    RealName:     "李四",
    RoleID:       role.ID,
    Department:   "信贷部",
    Position:     "信贷员",
    Status:       "active",
}
```

## 数据安全

### 密码安全
1. **加密存储**: 使用BCrypt哈希算法
2. **加盐处理**: 每个用户独立的盐值
3. **密码强度**: 强制密码复杂度要求

### 敏感信息保护
1. **JSON序列化**: 密码字段添加 `json:"-"` 标签
2. **访问控制**: 认证数据仅授权人员可访问
3. **数据脱敏**: 日志中不记录敏感信息

### 会话安全
1. **JWT令牌**: 使用强随机密钥签名
2. **令牌过期**: 合理设置访问令牌和刷新令牌过期时间
3. **会话撤销**: 支持主动撤销用户会话

## 性能优化

### 查询优化
1. **索引使用**: 关键查询字段建立适当索引
2. **分页查询**: 大数据量查询使用分页
3. **预加载**: 使用GORM的Preload减少N+1查询

### 缓存策略
1. **用户信息缓存**: 高频访问的用户基本信息缓存到Redis
2. **会话缓存**: 活跃会话信息缓存
3. **权限缓存**: 用户权限信息缓存

### 数据归档
1. **历史数据**: 定期归档过期的认证记录和会话
2. **软删除**: 重要数据使用软删除保留审计痕迹
3. **数据清理**: 定期清理过期的标签和临时数据 