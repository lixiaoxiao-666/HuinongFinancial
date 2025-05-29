# 数据库连接管理模块

## 文件概述

`connection.go` 是数字惠农后端服务的数据库连接管理核心模块，提供数据库连接创建、配置、迁移、初始化和健康检查等完整功能。

## 核心功能

### 1. 数据库连接管理
- **连接创建**: 基于配置创建GORM数据库连接
- **连接池配置**: 优化数据库连接池参数
- **连接测试**: 验证数据库连接可用性
- **DSN构建**: 自动构建MySQL数据源名称

### 2. 数据库迁移
- **自动迁移**: 基于模型自动创建/更新表结构
- **分模块迁移**: 按业务模块分组执行迁移
- **错误处理**: 详细的迁移失败信息

### 3. 默认数据初始化
- **系统配置**: 初始化基础系统配置项
- **用户角色**: 创建默认用户角色和权限
- **管理员账户**: 初始化超级管理员账户
- **内容分类**: 创建默认文章分类

### 4. 数据库监控
- **健康检查**: 检查数据库连接状态
- **连接统计**: 获取连接池使用情况
- **性能监控**: 连接池性能指标

## 主要函数详解

### NewConnection 创建数据库连接
```go
func NewConnection(cfg *config.DatabaseConfig) (*gorm.DB, error)
```

**功能说明**:
- 根据配置创建GORM数据库实例
- 配置连接池参数
- 设置GORM日志级别
- 执行连接测试

**关键配置**:
```go
gormConfig := &gorm.Config{
    NamingStrategy: schema.NamingStrategy{
        TablePrefix:   "",    // 表名前缀
        SingularTable: false, // 使用复数表名
    },
    Logger: getLogLevel(cfg),
    NowFunc: func() time.Time {
        return time.Now().Local()
    },
    DisableForeignKeyConstraintWhenMigrating: true,
}
```

**连接池配置**:
- `MaxIdleConns`: 最大空闲连接数
- `MaxOpenConns`: 最大打开连接数  
- `ConnMaxLifetime`: 连接最大生存时间

### AutoMigrate 自动迁移
```go
func AutoMigrate(db *gorm.DB) error
```

**迁移顺序**:
1. **用户相关表**: User, UserAuth, UserSession, UserTag, OAUser, OARole
2. **贷款相关表**: LoanProduct, LoanApplication, ApprovalLog, DifyWorkflowLog
3. **农机相关表**: Machine, RentalOrder
4. **通用表**: Article, Category, Expert, SystemConfig, FileUpload, OfflineQueue, APILog

**特点**:
- 按模块分组迁移，便于错误定位
- 支持表结构增量更新
- 自动处理索引和约束

### InitDefaultData 初始化默认数据
```go
func InitDefaultData(db *gorm.DB) error
```

**初始化内容**:
- 系统基础配置项
- 默认用户角色（超级管理员、信贷员、内容管理员）
- 管理员用户账户（username: admin, password: password）
- 文章分类（政策资讯、技术指导、市场信息、行业新闻）

## 数据库配置优化

### 连接池参数建议

| 环境 | MaxIdleConns | MaxOpenConns | ConnMaxLifetime |
|------|--------------|--------------|-----------------|
| 开发环境 | 5-10 | 20-50 | 1-2小时 |
| 测试环境 | 10-20 | 50-100 | 2-4小时 |
| 生产环境 | 20-50 | 100-200 | 4-8小时 |

### 日志级别设置
```go
func getLogLevel(cfg *config.DatabaseConfig) logger.Interface {
    logLevel := logger.Silent // 生产环境静默模式
    
    // 开发环境可以设置为Info级别以便调试
    // if cfg.Debug {
    //     logLevel = logger.Info
    // }
    
    return logger.Default.LogMode(logLevel)
}
```

## 默认数据详解

### 系统配置项
```go
configs := []model.SystemConfig{
    {
        ConfigKey:   "app.name",
        ConfigValue: "数字惠农系统",
        ConfigType:  "string",
        ConfigGroup: "system",
        Description: "应用程序名称",
    },
    {
        ConfigKey:   "loan.auto_approval_threshold",
        ConfigValue: "50000",
        ConfigType:  "int", 
        ConfigGroup: "business",
        Description: "贷款自动审批阈值(分)",
    },
    // ... 更多配置项
}
```

### 默认用户角色
```go
roles := []model.OARole{
    {
        Name:        "super_admin",
        DisplayName: "超级管理员", 
        Description: "拥有所有权限的超级管理员角色",
        Permissions: `{
            "loan_management": ["view", "create", "update", "approve", "delete"],
            "machine_management": ["view", "create", "update", "delete"],
            "user_management": ["view", "create", "update", "freeze", "delete"],
            "content_management": ["view", "create", "update", "publish", "delete"],
            "system_settings": ["view", "update"],
            "data_analytics": ["view", "export"]
        }`,
        IsSuper: true,
        Status:  "active",
    },
    // ... 更多角色
}
```

### 管理员账户
- **用户名**: admin
- **邮箱**: admin@huinong.com  
- **密码**: password (BCrypt加密)
- **角色**: super_admin
- **状态**: active

## 健康检查和监控

### HealthCheck 健康检查
```go
func HealthCheck(db *gorm.DB) error
```
- 检查数据库连接是否正常
- 执行简单的Ping测试
- 返回连接状态信息

### GetDatabaseStats 连接统计
```go
func GetDatabaseStats(db *gorm.DB) (map[string]interface{}, error)
```

**返回统计信息**:
```json
{
    "max_open_connections": 100,
    "open_connections": 5,
    "in_use": 2,
    "idle": 3,
    "wait_count": 0,
    "wait_duration": "0s",
    "max_idle_closed": 0,
    "max_lifetime_closed": 0
}
```

## 使用示例

### 基本使用
```go
// 创建数据库连接
db, err := database.NewConnection(&cfg.Database)
if err != nil {
    log.Fatal("数据库连接失败:", err)
}
defer func() {
    if sqlDB, err := db.DB(); err == nil {
        sqlDB.Close()
    }
}()

// 执行迁移
if err := database.AutoMigrate(db); err != nil {
    log.Fatal("数据库迁移失败:", err)
}

// 初始化默认数据
if err := database.InitDefaultData(db); err != nil {
    log.Fatal("初始化默认数据失败:", err)
}
```

### 健康检查
```go
// 检查数据库健康状态
if err := database.HealthCheck(db); err != nil {
    log.Printf("数据库连接异常: %v", err)
}

// 获取连接统计
stats, err := database.GetDatabaseStats(db)
if err != nil {
    log.Printf("获取数据库统计失败: %v", err)
} else {
    log.Printf("数据库连接统计: %+v", stats)
}
```

## 错误处理

### 常见错误类型
1. **连接失败**: 数据库服务不可用、网络问题
2. **认证失败**: 用户名密码错误、权限不足
3. **迁移失败**: 表结构冲突、字段类型不兼容
4. **初始化失败**: 默认数据插入失败、约束冲突

### 错误处理策略
```go
// 详细错误信息
if err := db.AutoMigrate(&model.User{}); err != nil {
    return fmt.Errorf("用户表迁移失败: %w", err)
}

// 事务处理
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&config).Error; err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```

## 性能优化

### 连接池优化
1. **连接数设置**: 根据并发量合理设置最大连接数
2. **连接生存时间**: 避免长时间占用连接
3. **空闲连接**: 保持适量空闲连接减少创建开销

### 查询优化
1. **索引使用**: 确保查询字段有适当索引
2. **批量操作**: 使用批量插入/更新减少数据库交互
3. **连接复用**: 避免频繁创建/销毁连接

### 监控指标
- 连接池使用率
- 查询响应时间
- 慢查询日志
- 连接等待时间

## 扩展功能

### 多数据库支持
```go
// 支持多个数据库连接
type DatabaseManager struct {
    Master *gorm.DB
    Slave  *gorm.DB
}

func (dm *DatabaseManager) GetReadDB() *gorm.DB {
    // 读操作使用从库
    if dm.Slave != nil {
        return dm.Slave
    }
    return dm.Master
}
```

### 数据库中间件
```go
// 自定义GORM插件
type LoggingPlugin struct{}

func (p *LoggingPlugin) Name() string {
    return "logging"
}

func (p *LoggingPlugin) Initialize(db *gorm.DB) error {
    return db.Callback().Query().Before("gorm:query").Register("logging:before_query", p.beforeQuery)
}
```

## 安全考虑

1. **连接加密**: 生产环境启用TLS连接
2. **权限控制**: 使用最小权限原则创建数据库用户
3. **SQL注入防护**: 使用参数化查询
4. **敏感信息**: 数据库密码使用环境变量
5. **审计日志**: 记录重要的数据库操作 