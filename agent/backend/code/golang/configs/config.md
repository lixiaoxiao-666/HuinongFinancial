# 配置文件说明

## 文件概述

`config.yaml` 是数字惠农后端服务的主配置文件，包含了应用运行所需的所有配置参数。

## 配置分类

### 应用基础配置 (app)
```yaml
app:
  name: "数字惠农后端服务"    # 应用名称
  version: "1.0.0"           # 版本号
  env: "development"         # 环境：development, production, test
  port: 8080                 # 服务端口
  mode: "debug"             # 运行模式：debug, release, test
```

### 数据库配置 (database)
```yaml
database:
  driver: "mysql"            # 数据库类型
  host: "localhost"          # 数据库主机
  port: 3306                # 数据库端口
  username: "huinong"        # 用户名
  password: "password"       # 密码
  database: "huinong_db"     # 数据库名
  charset: "utf8mb4"         # 字符集
  max_idle_conns: 10         # 最大空闲连接数
  max_open_conns: 100        # 最大打开连接数
  conn_max_lifetime: 3600    # 连接最大生存时间(秒)
```

### Redis缓存配置 (redis)
```yaml
redis:
  host: "localhost"          # Redis主机
  port: 6379                # Redis端口
  password: ""              # Redis密码
  database: 0               # 数据库索引
  pool_size: 100            # 连接池大小
  min_idle_conns: 10        # 最小空闲连接数
```

### JWT认证配置 (jwt)
```yaml
jwt:
  secret_key: "huinong-jwt-secret-key-2024"  # JWT密钥
  expires_in: 86400         # 访问令牌过期时间(秒)
  refresh_expires_in: 604800 # 刷新令牌过期时间(秒)
```

### Dify AI集成配置 (dify)
```yaml
dify:
  api_url: "https://api.dify.ai"    # Dify API地址
  api_key: "your-dify-api-key"      # API密钥
  timeout: 30                       # 请求超时时间(秒)
  retry_times: 3                    # 重试次数
  workflows:
    loan_approval: "loan-approval-workflow-id"  # 贷款审批工作流ID
```

### 短信服务配置 (sms)
```yaml
sms:
  provider: "aliyun"               # 服务提供商
  access_key: "your-access-key"    # 访问密钥
  access_secret: "your-access-secret"  # 访问秘钥
  sign_name: "数字惠农"            # 短信签名
  template_codes:                  # 短信模板编码
    register: "SMS_123456789"      # 注册验证码模板
    login: "SMS_123456790"         # 登录验证码模板
    reset_password: "SMS_123456791" # 重置密码模板
```

### 文件上传配置 (file)
```yaml
file:
  storage_type: "local"            # 存储类型：local, oss, qiniu
  upload_path: "./uploads"         # 上传路径
  max_file_size: 10485760         # 最大文件大小(字节)
  allowed_types: ["jpg", "jpeg", "png", "pdf", "doc", "docx"]  # 允许的文件类型
```

### 日志配置 (log)
```yaml
log:
  level: "debug"                   # 日志级别：debug, info, warn, error
  format: "json"                   # 日志格式：text, json
  output: "both"                   # 输出方式：console, file, both
  file_path: "./logs/app.log"      # 日志文件路径
  max_size: 100                    # 单个日志文件最大大小(MB)
  max_age: 7                       # 日志文件保留天数
  max_backups: 5                   # 最大备份文件数
```

### CORS跨域配置 (cors)
```yaml
cors:
  allow_origins: ["*"]             # 允许的来源
  allow_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]  # 允许的方法
  allow_headers: ["Origin", "Content-Type", "Accept", "Authorization"]  # 允许的请求头
  expose_headers: ["Content-Length"]  # 暴露的响应头
  allow_credentials: true          # 是否允许凭证
  max_age: 12                      # 预检请求缓存时间(小时)
```

### 限流配置 (rate_limit)
```yaml
rate_limit:
  enabled: true                    # 是否启用限流
  requests_per_minute: 1000        # 每分钟请求数限制
  burst: 100                       # 突发请求数
```

### 离线功能配置 (offline)
```yaml
offline:
  enabled: true                    # 是否启用离线功能
  sync_interval: 300               # 同步间隔(秒)
  queue_size: 10000               # 队列大小
  retry_times: 3                   # 重试次数
```

## 环境变量覆盖

配置文件中的值可以通过环境变量覆盖，环境变量命名规则：
- 将配置项层级用下划线连接
- 全部转为大写
- 添加 `HUINONG_` 前缀

示例：
- `app.port` → `HUINONG_APP_PORT`
- `database.host` → `HUINONG_DATABASE_HOST`
- `jwt.secret_key` → `HUINONG_JWT_SECRET_KEY`

## 安全注意事项

1. **生产环境**: 必须修改所有默认密钥和密码
2. **敏感信息**: API密钥、数据库密码等敏感信息建议使用环境变量
3. **JWT密钥**: 使用足够复杂的随机字符串
4. **CORS配置**: 生产环境不应使用 `["*"]` 作为允许来源

## 配置验证

系统启动时会自动验证配置的有效性，包括：
- 必填项检查
- 数据类型验证
- 取值范围校验
- 依赖项检查 