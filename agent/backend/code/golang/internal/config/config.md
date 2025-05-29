# 配置管理模块

## 文件概述

`config.go` 是数字惠农后端服务的配置管理核心模块，提供了完整的配置结构定义、加载、验证和管理功能。

## 核心功能

### 1. 配置结构体定义
- **Config**: 根配置结构体，整合所有子配置模块
- **AppConfig**: 应用基础配置（名称、版本、环境、端口、Host等）
- **DatabaseConfig**: 数据库连接配置
- **RedisConfig**: Redis缓存配置
- **JWTConfig**: JWT认证配置
- **DifyConfig**: Dify AI集成配置
- **SMSConfig**: 短信服务配置
- **FileConfig**: 文件上传配置
- **LogConfig**: 日志配置
- **CORSConfig**: 跨域配置
- **RateLimitConfig**: 限流配置
- **OfflineConfig**: 离线功能配置

### 2. 配置加载功能
```go
func LoadConfig(configPath string) (*Config, error)
```
- 支持YAML格式配置文件
- 环境变量覆盖机制
- 自动类型转换和映射

### 3. 配置验证机制
```go
func (c *Config) validate() error
```
- 必填项校验
- 数据类型验证
- 取值范围检查
- 依赖项验证

## 配置结构详解

### AppConfig 应用配置
```go
type AppConfig struct {
    Name    string `mapstructure:"name"`     // 应用名称
    Version string `mapstructure:"version"`  // 版本号
    Env     string `mapstructure:"env"`      // 运行环境
    Host    string `mapstructure:"host"`     // 服务器监听地址
    Port    int    `mapstructure:"port"`     // 监听端口
    Mode    string `mapstructure:"mode"`     // Gin运行模式
    Debug   bool   `mapstructure:"debug"`    // 调试模式
}
```

**环境类型**:
- `development`: 开发环境 - 启用调试功能和详细日志
- `testing`: 测试环境 - 用于自动化测试
- `staging`: 预发布环境 - 接近生产环境的测试
- `production`: 生产环境 - 正式运行环境

### DatabaseConfig 数据库配置
```go
type DatabaseConfig struct {
    Driver          string `mapstructure:"driver"`           // 数据库驱动
    Host            string `mapstructure:"host"`             // 主机地址
    Port            int    `mapstructure:"port"`             // 端口
    Username        string `mapstructure:"username"`         // 用户名
    Password        string `mapstructure:"password"`         // 密码
    Database        string `mapstructure:"database"`         // 数据库名
    Charset         string `mapstructure:"charset"`          // 字符集
    MaxIdleConns    int    `mapstructure:"max_idle_conns"`   // 最大空闲连接数
    MaxOpenConns    int    `mapstructure:"max_open_conns"`   // 最大打开连接数
    ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // 连接最大生存时间(秒)
}
```

**连接池配置说明**:
- `MaxIdleConns`: 连接池中保持的最大空闲连接数，建议值：10-50
- `MaxOpenConns`: 数据库的最大打开连接数，建议值：50-200
- `ConnMaxLifetime`: 连接的最大生存时间，建议值：300-3600秒

### RedisConfig Redis配置
```go
type RedisConfig struct {
    Host         string `mapstructure:"host"`          // Redis主机
    Port         int    `mapstructure:"port"`          // Redis端口
    Password     string `mapstructure:"password"`      // Redis密码
    Database     int    `mapstructure:"database"`      // 数据库索引
    PoolSize     int    `mapstructure:"pool_size"`     // 连接池大小
    MinIdleConns int    `mapstructure:"min_idle_conns"` // 最小空闲连接数
}
```

### JWTConfig JWT认证配置
```go
type JWTConfig struct {
    SecretKey      string `mapstructure:"secret_key"`       // JWT密钥
    ExpirationTime int    `mapstructure:"expiration_time"`  // 过期时间(小时)
    RefreshTime    int    `mapstructure:"refresh_time"`     // 刷新时间(小时)
    Issuer         string `mapstructure:"issuer"`           // 签发者
}
```

**安全建议**:
- `SecretKey`: 使用强随机密钥，长度至少32字符
- `ExpirationTime`: Token过期时间，建议2-24小时
- `RefreshTime`: 刷新Token有效期，建议7-30天

### DifyConfig Dify AI配置
```go
type DifyConfig struct {
    APIURL     string            `mapstructure:"api_url"`
    APIKey     string            `mapstructure:"api_key"`
    APIToken   string            `mapstructure:"api_token"`
    Timeout    int               `mapstructure:"timeout"`
    RetryTimes int               `mapstructure:"retry_times"`
    Workflows  map[string]string `mapstructure:"workflows"`
}
```

### SMSConfig 短信服务配置
```go
type SMSConfig struct {
    Provider        string `mapstructure:"provider"`         // 服务提供商
    AccessKeyID     string `mapstructure:"access_key_id"`    // 访问密钥ID
    AccessKeySecret string `mapstructure:"access_key_secret"` // 访问密钥
    SignName        string `mapstructure:"sign_name"`        // 签名
    TemplateCode    string `mapstructure:"template_code"`    // 模板代码
}
```

**支持的短信提供商**:
- `aliyun`: 阿里云短信服务
- `tencent`: 腾讯云短信服务
- `huawei`: 华为云短信服务

### FileConfig 文件配置
```go
type FileConfig struct {
    StorageType    string `mapstructure:"storage_type"`     // 存储类型
    MaxSize        int64  `mapstructure:"max_size"`         // 最大文件大小
    AllowedTypes   string `mapstructure:"allowed_types"`    // 允许的文件类型
    UploadPath     string `mapstructure:"upload_path"`      // 上传路径
    BaseURL        string `mapstructure:"base_url"`         // 文件访问基础URL
    
    // OSS配置(当storage_type为oss时)
    OSSEndpoint        string `mapstructure:"oss_endpoint"`
    OSSAccessKeyID     string `mapstructure:"oss_access_key_id"`
    OSSAccessKeySecret string `mapstructure:"oss_access_key_secret"`
    OSSBucketName      string `mapstructure:"oss_bucket_name"`
}
```

### LogConfig 日志配置
```go
type LogConfig struct {
    Level      string `mapstructure:"level"`       // 日志级别
    Format     string `mapstructure:"format"`      // 日志格式
    Output     string `mapstructure:"output"`      // 输出目标
    Filename   string `mapstructure:"filename"`    // 日志文件名
    MaxSize    int    `mapstructure:"max_size"`    // 单个文件最大大小(MB)
    MaxBackups int    `mapstructure:"max_backups"` // 保留的旧文件个数
    MaxAge     int    `mapstructure:"max_age"`     // 保留的最大天数
    Compress   bool   `mapstructure:"compress"`    // 是否压缩
}
```

**日志级别**:
- `debug`: 调试信息 - 包含所有日志
- `info`: 一般信息 - 包含info及以上级别
- `warn`: 警告信息 - 包含warn及以上级别
- `error`: 错误信息 - 仅错误和致命错误
- `fatal`: 致命错误 - 仅致命错误

### CORSConfig 跨域配置
```go
type CORSConfig struct {
    AllowOrigins     []string `mapstructure:"allow_origins"`      // 允许的来源
    AllowMethods     []string `mapstructure:"allow_methods"`      // 允许的方法
    AllowHeaders     []string `mapstructure:"allow_headers"`      // 允许的头部
    ExposeHeaders    []string `mapstructure:"expose_headers"`     // 暴露的头部
    AllowCredentials bool     `mapstructure:"allow_credentials"`  // 是否允许凭证
    MaxAge           int      `mapstructure:"max_age"`            // 预检请求缓存时间
}
```

### RateLimitConfig 限流配置
```go
type RateLimitConfig struct {
    Enable    bool `mapstructure:"enable"`     // 是否启用限流
    RPS       int  `mapstructure:"rps"`        // 每秒请求数
    Burst     int  `mapstructure:"burst"`      // 突发请求数
    WindowSize int `mapstructure:"window_size"` // 时间窗口大小(秒)
}
```

### OfflineConfig 离线功能配置
```go
type OfflineConfig struct {
    Enable         bool `mapstructure:"enable"`          // 是否启用离线功能
    QueueSize      int  `mapstructure:"queue_size"`      // 队列大小
    BatchSize      int  `mapstructure:"batch_size"`      // 批处理大小
    ProcessInterval int `mapstructure:"process_interval"` // 处理间隔(秒)
    MaxRetries     int  `mapstructure:"max_retries"`     // 最大重试次数
}
```

## 配置文件示例

### 完整的config.yaml示例
```yaml
# 应用配置
app:
  name: "数字惠农系统"
  version: "1.0.0"
  env: "production"
  host: "0.0.0.0"
  port: 8080
  mode: "debug"
  debug: false

# 数据库配置
database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "huinong"
  password: "your_password"
  database: "huinong_db"
  charset: "utf8mb4"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600

# Redis配置
redis:
  host: "localhost"
  port: 6379
  password: ""
  database: 0
  pool_size: 10
  min_idle_conns: 5

# JWT配置
jwt:
  secret_key: "your_super_secret_key_at_least_32_chars"
  expiration_time: 24
  refresh_time: 168
  issuer: "huinong-system"

# Dify AI配置
dify:
  api_url: "https://api.dify.ai/v1"
  api_key: "your_dify_api_key"
  api_token: "dify-huinong-secure-token-2024"
  timeout: 30
  retry_times: 3
  workflows:
    loan_approval: "loan-approval-workflow-id"

# 短信配置
sms:
  provider: "aliyun"
  access_key_id: "your_access_key_id"
  access_key_secret: "your_access_key_secret"
  sign_name: "数字惠农"
  template_code: "SMS_123456789"

# 文件配置
file:
  storage_type: "local"
  max_size: 10485760  # 10MB
  allowed_types: "jpg,jpeg,png,pdf,doc,docx"
  upload_path: "./uploads"
  base_url: "http://localhost:8080/files"

# 日志配置
log:
  level: "info"
  format: "json"
  output: "file"
  filename: "./logs/app.log"
  max_size: 100    # MB
  max_backups: 7
  max_age: 30      # days
  compress: true

# 跨域配置
cors:
  allow_origins: 
    - "http://localhost:3000"
    - "https://your-frontend-domain.com"
  allow_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allow_headers:
    - "Origin"
    - "Content-Type"
    - "Authorization"
  expose_headers:
    - "X-Total-Count"
  allow_credentials: true
  max_age: 300

# 限流配置
rate_limit:
  enable: true
  rps: 1000
  burst: 2000

# 离线功能配置
offline:
  enable: true
  queue_size: 1000
  batch_size: 50
  process_interval: 30
  max_retries: 3
```

## 配置加载流程

### LoadConfig函数详解
```go
func LoadConfig(configPath string) (*Config, error) {
    // 1. 设置配置文件
    viper.SetConfigFile(configPath)
    viper.SetConfigType("yaml")
    
    // 2. 设置环境变量前缀
    viper.SetEnvPrefix("HUINONG")
    viper.AutomaticEnv()
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
    // 3. 读取配置文件
    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("读取配置文件失败: %w", err)
    }
    
    // 4. 解析配置到结构体
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("配置解析失败: %w", err)
    }
    
    // 5. 验证配置
    if err := config.validate(); err != nil {
        return nil, fmt.Errorf("配置验证失败: %w", err)
    }
    
    return &config, nil
}
```

### 环境变量覆盖机制
配置支持通过环境变量覆盖，变量名规则：
```bash
# 格式: HUINONG_<section>_<key>
export HUINONG_APP_PORT=9090
export HUINONG_DATABASE_HOST=production-db.example.com
export HUINONG_JWT_SECRET_KEY=production_secret_key
export HUINONG_REDIS_PASSWORD=redis_password
```

### 配置验证规则
```go
func (c *Config) validate() error {
    // 应用配置验证
    if c.App.Name == "" {
        return errors.New("应用名称不能为空")
    }
    if c.App.Port <= 0 || c.App.Port > 65535 {
        return errors.New("端口号必须在1-65535之间")
    }
    
    // 数据库配置验证
    if c.Database.Host == "" {
        return errors.New("数据库主机不能为空")
    }
    if c.Database.Username == "" {
        return errors.New("数据库用户名不能为空")
    }
    if c.Database.Database == "" {
        return errors.New("数据库名不能为空")
    }
    
    // JWT配置验证
    if len(c.JWT.SecretKey) < 32 {
        return errors.New("JWT密钥长度不能少于32字符")
    }
    
    // 文件大小验证
    if c.File.MaxSize <= 0 {
        return errors.New("文件最大大小必须大于0")
    }
    
    return nil
}
```

## 使用示例

### 基本使用
```go
package main

import (
    "log"
    "huinong-backend/internal/config"
)

func main() {
    // 加载配置
    cfg, err := config.LoadConfig("./configs/config.yaml")
    if err != nil {
        log.Fatal("加载配置失败:", err)
    }
    
    // 使用配置
    log.Printf("应用名称: %s", cfg.App.Name)
    log.Printf("运行端口: %d", cfg.App.Port)
    log.Printf("数据库: %s@%s:%d/%s", 
        cfg.Database.Username,
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.Database)
}
```

### 配置热重载
```go
func WatchConfig(configPath string, callback func(*Config)) {
    viper.SetConfigFile(configPath)
    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        log.Println("配置文件发生变化:", e.Name)
        
        var newConfig Config
        if err := viper.Unmarshal(&newConfig); err != nil {
            log.Printf("重新解析配置失败: %v", err)
            return
        }
        
        if err := newConfig.validate(); err != nil {
            log.Printf("新配置验证失败: %v", err)
            return
        }
        
        callback(&newConfig)
    })
}
```

### 配置优先级
配置加载的优先级从高到低：
1. **环境变量**: `HUINONG_*`格式的环境变量
2. **配置文件**: YAML格式的配置文件
3. **默认值**: 代码中定义的默认值

### 敏感配置处理
```go
// 敏感配置脱敏显示
func (c *Config) SafeString() string {
    safeCfg := *c
    
    // 脱敏处理
    if safeCfg.Database.Password != "" {
        safeCfg.Database.Password = "***"
    }
    if safeCfg.Redis.Password != "" {
        safeCfg.Redis.Password = "***"
    }
    if safeCfg.JWT.SecretKey != "" {
        safeCfg.JWT.SecretKey = "***"
    }
    
    data, _ := yaml.Marshal(safeCfg)
    return string(data)
}
```

## 部署配置

### 开发环境配置
```yaml
app:
  env: "development"
  debug: true
  port: 8080

log:
  level: "debug"
  format: "text"
  output: "console"

cors:
  allow_origins: ["*"]
```

### 生产环境配置
```yaml
app:
  env: "production"
  debug: false
  port: 8080

log:
  level: "info"
  format: "json"
  output: "file"

cors:
  allow_origins: 
    - "https://app.yourdomain.com"
    - "https://admin.yourdomain.com"

rate_limit:
  enable: true
  rps: 1000
  burst: 2000
```

### Docker环境配置
```dockerfile
# Dockerfile
ENV HUINONG_APP_PORT=8080
ENV HUINONG_DATABASE_HOST=mysql
ENV HUINONG_REDIS_HOST=redis
```

```yaml
# docker-compose.yml
version: '3.8'
services:
  app:
    build: .
    environment:
      - HUINONG_DATABASE_HOST=mysql
      - HUINONG_DATABASE_PASSWORD=${DB_PASSWORD}
      - HUINONG_REDIS_HOST=redis
      - HUINONG_JWT_SECRET_KEY=${JWT_SECRET}
```

## 最佳实践

### 1. 配置分层
```
configs/
├── config.yaml          # 基础配置
├── development.yaml      # 开发环境配置
├── testing.yaml         # 测试环境配置
├── staging.yaml         # 预发布环境配置
└── production.yaml      # 生产环境配置
```

### 2. 安全配置
- 敏感信息通过环境变量传递
- 生产环境禁用调试模式
- 使用强密钥和安全的默认值
- 定期轮换密钥和密码

### 3. 性能调优
- 根据负载调整连接池大小
- 合理设置超时时间
- 启用适当的缓存策略
- 监控资源使用情况

### 4. 监控和告警
- 记录配置加载失败事件
- 监控配置变更
- 设置关键配置异常告警
- 定期备份配置文件

### 5. 文档维护
- 及时更新配置文档
- 记录配置变更历史
- 提供配置示例
- 说明各配置项的影响

## 重要更新记录

### Host配置支持 (2025-05-30)

#### 1. AppConfig结构体增强
```go
type AppConfig struct {
    Name    string `mapstructure:"name"`
    Version string `mapstructure:"version"`
    Env     string `mapstructure:"env"`
    Host    string `mapstructure:"host"`     // 新增Host字段
    Port    int    `mapstructure:"port"`
    Mode    string `mapstructure:"mode"`
}
```

#### 2. 新增服务器地址获取方法
```go
// GetServerAddr 获取服务器监听地址
func (c *AppConfig) GetServerAddr() string {
    // 如果Host为空，默认使用 "0.0.0.0"
    host := c.Host
    if host == "" {
        host = "0.0.0.0"
    }
    return fmt.Sprintf("%s:%d", host, c.Port)
}
```

#### 3. DifyConfig增强
```go
type DifyConfig struct {
    APIURL     string            `mapstructure:"api_url"`
    APIKey     string            `mapstructure:"api_key"`
    APIToken   string            `mapstructure:"api_token"`  // 新增APIToken字段
    Timeout    int               `mapstructure:"timeout"`
    RetryTimes int               `mapstructure:"retry_times"`
    Workflows  map[string]string `mapstructure:"workflows"`
}
```

## 配置方法说明

### AppConfig方法
- `GetServerAddr()`: 获取完整的服务器监听地址
- `IsDevelopment()`: 判断是否为开发环境
- `IsProduction()`: 判断是否为生产环境
- `IsTest()`: 判断是否为测试环境

### DatabaseConfig方法
- `GetDSN()`: 获取数据库连接字符串

### RedisConfig方法
- `GetRedisAddr()`: 获取Redis连接地址

### JWTConfig方法
- `GetJWTExpirationDuration()`: 获取JWT过期时间
- `GetRefreshExpirationDuration()`: 获取刷新令牌过期时间

## 注意事项

1. **Host配置**: 新增的Host配置允许指定服务器监听的具体地址
   - `0.0.0.0`: 监听所有网络接口
   - `127.0.0.1`: 仅监听本地回环接口
   - 具体IP: 监听指定网络接口

2. **配置优先级**: 环境变量 > 配置文件
3. **配置验证**: 启动时会自动验证配置的有效性
4. **目录创建**: 会自动创建必要的目录（如日志目录、上传目录）
5. **安全性**: 敏感配置（如JWT密钥）应通过环境变量设置 