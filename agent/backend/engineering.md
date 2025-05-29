# 数字惠农后端工程化文档

## 项目概述

数字惠农APP及OA后台管理系统是一个综合性的农业数字化平台，采用单体架构设计，支持离线运行环境。

### 核心特性
- **单体架构**: 简化部署和维护，避免微服务复杂性
- **离线优先**: 支持无网络环境运行，数据同步机制
- **AI集成**: 基于Dify平台的智能贷款审批工作流
- **权限管理**: 基于角色的访问控制(RBAC)系统
- **多平台支持**: APP端、Web端、OA后台统一API

## 技术栈

### 后端技术栈
- **语言**: Go 1.21
- **Web框架**: Gin v1.9.1
- **ORM**: GORM v1.25.5
- **数据库**: TiDB（MySQL兼容）
- **缓存**: Redis v9.2.1
- **认证**: JWT v4.5.0
- **配置管理**: Viper v1.17.0
- **文档**: Swagger/OpenAPI 3.0
- **加密**: bcrypt, AES, SHA256

### 外部集成
- **Dify AI平台**: 智能审批工作流
- **短信服务**: 验证码发送
- **文件存储**: 支持本地/云存储
- **地理位置**: GPS定位服务

## 项目结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # 应用程序入口
├── configs/
│   └── config.yaml              # 配置文件
├── internal/
│   ├── config/
│   │   └── config.go            # 配置管理
│   ├── database/
│   │   └── connection.go        # 数据库连接
│   ├── model/
│   │   ├── user.go              # 用户数据模型
│   │   ├── loan.go              # 贷款数据模型
│   │   ├── machine.go           # 农机数据模型
│   │   └── common.go            # 通用数据模型
│   ├── repository/
│   │   ├── interface.go         # Repository接口定义
│   │   ├── user_repository.go   # 用户数据访问层
│   │   └── loan_repository.go   # 贷款数据访问层
│   ├── service/
│   │   ├── interface.go         # Service接口定义
│   │   └── user_service.go      # 用户业务逻辑层
│   ├── handler/
│   │   └── user_handler.go      # HTTP请求处理层
│   ├── middleware/
│   │   └── auth.go              # 认证中间件
│   ├── router/
│   │   └── router.go            # 路由配置
│   ├── cache/
│   │   └── redis.go             # 缓存管理
│   └── utils/
│       └── crypto.go            # 加密工具函数
├── go.mod                       # Go模块依赖
└── go.sum                       # 依赖版本锁定
```

## 架构设计

### 分层架构
```
┌─────────────────────────────────────┐
│            Handler Layer            │  HTTP请求处理
├─────────────────────────────────────┤
│            Service Layer            │  业务逻辑处理
├─────────────────────────────────────┤
│          Repository Layer           │  数据访问抽象
├─────────────────────────────────────┤
│             Model Layer             │  数据模型定义
└─────────────────────────────────────┘
```

### 核心组件

#### 1. 数据模型层 (Model)
- **用户管理**: User、UserAuth、UserSession、UserTag
- **贷款业务**: LoanProduct、LoanApplication、ApprovalLog、DifyWorkflowLog
- **农机租赁**: Machine、RentalOrder
- **内容管理**: Article、Category、Expert
- **系统管理**: SystemConfig、FileUpload、OfflineQueue、APILog
- **OA管理**: OAUser、OARole

#### 2. 数据访问层 (Repository)
```go
type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    GetByID(ctx context.Context, id uint64) (*model.User, error)
    GetByPhone(ctx context.Context, phone string) (*model.User, error)
    List(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error)
    // ... 更多方法
}
```

#### 3. 业务逻辑层 (Service)
```go
type UserService interface {
    Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
    Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
    GetProfile(ctx context.Context, userID uint64) (*UserProfileResponse, error)
    // ... 更多方法
}
```

#### 4. HTTP处理层 (Handler)
```go
func (h *UserHandler) Register(c *gin.Context) {
    var req service.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // 参数验证错误处理
    }
    
    result, err := h.userService.Register(ctx, &req)
    if err != nil {
        // 业务错误处理
    }
    
    c.JSON(http.StatusOK, NewSuccessResponse("注册成功", result))
}
```

### 中间件设计

#### 1. 认证中间件
- **JWT验证**: 解析和验证访问令牌
- **会话管理**: 用户会话状态检查
- **权限控制**: 基于角色的访问控制

#### 2. 功能中间件
- **CORS处理**: 跨域请求支持
- **请求日志**: API调用记录
- **错误恢复**: 全局异常处理
- **限流控制**: 防止API滥用

### 缓存策略

#### Redis缓存设计
```
用户缓存:     user:{user_id}
会话缓存:     session:{session_id}
短信验证码:   sms_code:{phone}
文章缓存:     article:{article_id}
产品缓存:     loan_product:{product_id}
配置缓存:     config:{config_key}
限流计数:     limit:{rate_key}
分布式锁:     lock:{lock_key}
```

#### 缓存策略
- **用户信息**: 30分钟过期，写后失效
- **会话数据**: 24小时过期，活跃延期
- **系统配置**: 24小时过期，更新失效
- **文章内容**: 1小时过期，发布失效
- **产品信息**: 2小时过期，修改失效

## API接口设计

### 认证接口
```
POST /api/auth/register     # 用户注册
POST /api/auth/login        # 用户登录
POST /api/auth/refresh      # 刷新Token
POST /api/auth/logout       # 用户登出
```

### 用户管理接口
```
GET  /api/user/profile      # 获取用户资料
PUT  /api/user/profile      # 更新用户资料
PUT  /api/user/password     # 修改密码
POST /api/user/auth/real-name   # 实名认证
POST /api/user/auth/bank-card   # 银行卡认证
GET  /api/user/tags         # 获取用户标签
POST /api/user/tags         # 添加用户标签
```

### 管理员接口
```
GET  /api/admin/users       # 用户列表
GET  /api/admin/users/statistics # 用户统计
PUT  /api/admin/users/:id/freeze # 冻结用户
```

### 响应格式
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    // 具体数据
  },
  "error": null
}
```

## 数据库设计

### 核心表结构

#### 用户表 (users)
```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    uuid VARCHAR(36) UNIQUE NOT NULL,
    username VARCHAR(50),
    phone VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(100),
    password_hash VARCHAR(255) NOT NULL,
    salt VARCHAR(32) NOT NULL,
    user_type ENUM('farmer', 'owner', 'expert', 'agent') NOT NULL,
    status ENUM('active', 'frozen', 'deleted') DEFAULT 'active',
    -- ... 更多字段
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 贷款申请表 (loan_applications)
```sql
CREATE TABLE loan_applications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    application_no VARCHAR(32) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    status ENUM('draft', 'submitted', 'reviewing', 'approved', 'rejected', 'cancelled') DEFAULT 'draft',
    -- ... 更多字段
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES loan_products(id)
);
```

### 索引设计
```sql
-- 用户表索引
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_type_status ON users(user_type, status);

-- 贷款申请表索引
CREATE INDEX idx_loan_apps_user ON loan_applications(user_id);
CREATE INDEX idx_loan_apps_status ON loan_applications(status);
CREATE INDEX idx_loan_apps_created ON loan_applications(created_at);
```

## 安全设计

### 认证安全
- **密码加密**: bcrypt + 随机盐值
- **JWT令牌**: HS256签名算法
- **会话管理**: Redis存储，过期自动清理
- **双令牌机制**: Access Token + Refresh Token

### 数据安全
- **敏感数据加密**: AES-256-CFB模式
- **传输安全**: HTTPS强制
- **SQL注入防护**: GORM参数化查询
- **XSS防护**: 输入验证和输出编码

### 权限控制
```go
// 权限检查中间件
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := GetUserIDFromContext(c)
        // 检查用户角色
        hasPermission := checkUserRole(userID, roles)
        if !hasPermission {
            c.JSON(403, gin.H{"error": "权限不足"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

## 离线机制设计

### 离线队列
```go
type OfflineQueue struct {
    ID           uint64    `gorm:"primaryKey;autoIncrement"`
    UserID       uint64    `gorm:"not null;index"`
    ActionType   string    `gorm:"size:50;not null"`
    RequestData  string    `gorm:"type:longtext"`
    Status       string    `gorm:"size:20;default:'pending'"`
    RetryCount   int       `gorm:"default:0"`
    CreatedAt    time.Time
}
```

### 数据同步策略
1. **离线操作记录**: 所有用户操作记录到离线队列
2. **批量同步**: 网络恢复后批量处理队列
3. **冲突解决**: 时间戳优先原则
4. **失败重试**: 指数退避重试机制

## AI集成设计

### Dify工作流集成
```go
type DifyWorkflowLog struct {
    ID             uint64    `gorm:"primaryKey;autoIncrement"`
    ApplicationID  uint64    `gorm:"not null;index"`
    WorkflowType   string    `gorm:"size:50;not null"`
    RequestData    string    `gorm:"type:longtext"`
    ResponseData   string    `gorm:"type:longtext"`
    Status         string    `gorm:"size:20;not null"`
    ProcessingTime int       `gorm:"comment:处理时间(毫秒)"`
    ErrorMessage   string    `gorm:"size:500"`
    CreatedAt      time.Time
}
```

### 工作流类型
- **信用评估**: credit_assessment
- **风险分析**: risk_analysis
- **额度计算**: amount_calculation
- **决策建议**: decision_recommendation

## 性能优化

### 数据库优化
- **连接池**: 最大20个连接，空闲5个
- **查询优化**: 索引覆盖，避免全表扫描
- **分页查询**: Limit/Offset优化
- **批量操作**: 减少数据库往返

### 缓存优化
- **多级缓存**: 内存缓存 + Redis缓存
- **缓存预热**: 系统启动时预加载热数据
- **缓存穿透**: 空结果也缓存
- **缓存雪崩**: 随机过期时间

### API优化
- **响应压缩**: Gzip压缩
- **并发限制**: 协程池控制
- **超时控制**: 请求超时设置
- **熔断机制**: 故障快速失败

## 监控和日志

### API日志记录
```go
type APILog struct {
    ID            uint64    `gorm:"primaryKey;autoIncrement"`
    Method        string    `gorm:"size:10;not null"`
    URL           string    `gorm:"size:500;not null"`
    UserID        *uint64   `gorm:"index"`
    IPAddress     string    `gorm:"size:45;not null"`
    UserAgent     string    `gorm:"size:500"`
    RequestData   string    `gorm:"type:longtext"`
    ResponseCode  int       `gorm:"not null"`
    ResponseTime  int       `gorm:"comment:响应时间(毫秒)"`
    CreatedAt     time.Time
}
```

### 健康检查
```go
func HealthCheck(ctx context.Context) (*HealthCheckResponse, error) {
    return &HealthCheckResponse{
        Status:    "ok",
        Database:  checkDatabase(),
        Redis:     checkRedis(),
        Services:  checkServices(),
        Timestamp: time.Now().Unix(),
    }, nil
}
```

## 部署和运维

### 环境配置
```yaml
# config.yaml
app:
  name: "数字惠农API"
  version: "1.0.0"
  environment: "production"
  port: 8080

database:
  host: "localhost"
  port: 4000
  username: "huinong"
  password: "password"
  database: "huinong_db"

redis:
  host: "localhost"
  port: 6379
  password: ""
  database: 0

jwt:
  secret: "your-jwt-secret-key"
  expires_in: 86400  # 24小时
```

### Docker部署
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs
EXPOSE 8080
CMD ["./main"]
```

### 服务监控
- **系统指标**: CPU、内存、磁盘使用率
- **应用指标**: API响应时间、错误率、QPS
- **业务指标**: 用户注册量、贷款申请量、机器租赁量
- **告警机制**: 关键指标异常自动告警

## 开发规范

### 代码规范
- **命名规范**: 驼峰命名，见名知意
- **错误处理**: 统一错误格式，详细错误信息
- **注释规范**: 公开接口必须有注释
- **测试覆盖**: 核心逻辑测试覆盖率80%+

### Git工作流
```
master (主分支)
├── develop (开发分支)
├── feature/user-auth (功能分支)
├── hotfix/login-bug (热修复分支)
└── release/v1.0.0 (发布分支)
```

### API文档
- **Swagger注解**: 所有API接口都有Swagger文档
- **示例数据**: 提供完整的请求响应示例
- **错误码说明**: 详细的错误码对照表
- **版本管理**: API版本向后兼容

## 项目里程碑

### Phase 1: 基础架构 ✅
- [x] 项目初始化和依赖管理
- [x] 数据库连接和配置管理
- [x] 数据模型设计和实现
- [x] Repository层接口定义

### Phase 2: 核心功能 ✅
- [x] Service层业务逻辑实现
- [x] Handler层HTTP处理
- [x] 认证中间件和JWT
- [x] 路由配置和API设计

### Phase 3: 增强功能 ✅
- [x] Redis缓存管理
- [x] 工具函数和加密
- [x] 中间件完善
- [x] 错误处理和日志

### Phase 4: 待完成功能
- [ ] Repository层具体实现
- [ ] 贷款Service和Handler
- [ ] 农机Service和Handler
- [ ] 文件上传和管理
- [ ] Dify AI集成
- [ ] 短信验证码服务
- [ ] 单元测试和集成测试
- [ ] 性能测试和优化
- [ ] 部署脚本和文档

## 总结

数字惠农后端系统采用现代化的Go技术栈，遵循分层架构和领域驱动设计原则，具有以下特点：

1. **架构清晰**: 分层明确，职责单一，易于维护
2. **功能完整**: 覆盖用户管理、贷款业务、农机租赁等核心功能
3. **性能优化**: 多级缓存、数据库优化、并发控制
4. **安全可靠**: 多重安全机制，数据加密，权限控制
5. **易于扩展**: 接口驱动，松耦合设计，支持功能扩展
6. **运维友好**: 完善的监控、日志、健康检查机制

通过工程化的开发流程和规范，确保了代码质量和项目的可维护性，为数字农业平台的长期发展奠定了坚实的技术基础。