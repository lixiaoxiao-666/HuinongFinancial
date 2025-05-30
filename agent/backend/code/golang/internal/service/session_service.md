# 会话管理服务 - 代码文档

## 1. 概述

`session_service.go` 实现了基于Redis的分布式会话管理服务，是数字惠农系统中负责用户认证状态维护和多后端会话保持的核心组件。

## 2. 主要功能

### 2.1 核心特性
- 🔐 **分布式会话管理**: 基于Redis实现跨多后端实例的会话共享
- 📱 **多端登录支持**: 支持用户在不同设备、平台同时登录
- ⚡ **高性能缓存**: 毫秒级会话验证，支持高并发访问
- 🔄 **自动清理机制**: 过期会话自动清理，防止内存泄漏
- 🔒 **安全机制**: Token哈希存储，会话状态实时同步

### 2.2 数据结构设计

#### LoginInfo - 登录信息
```go
type LoginInfo struct {
    Platform    string `json:"platform"`    // 平台类型: app, web, oa
    DeviceID    string `json:"device_id"`   // 设备唯一标识
    DeviceType  string `json:"device_type"` // 设备类型: ios, android, web
    DeviceName  string `json:"device_name"` // 设备名称
    AppVersion  string `json:"app_version"` // 应用版本
    UserAgent   string `json:"user_agent"`  // 用户代理
    IPAddress   string `json:"ip_address"`  // IP地址
    Location    string `json:"location"`    // 地理位置
    LoginMethod string `json:"login_method"` // 登录方式: password, sms, oauth
}
```

#### SessionInfo - 会话信息
```go
type SessionInfo struct {
    SessionID       string       `json:"session_id"`      // 会话唯一标识
    UserID          uint64       `json:"user_id"`         // 用户ID
    Platform        string       `json:"platform"`        // 登录平台
    DeviceInfo      *DeviceInfo  `json:"device_info"`     // 设备信息
    NetworkInfo     *NetworkInfo `json:"network_info"`    // 网络信息
    TokenInfo       *TokenInfo   `json:"token_info"`      // Token信息
    Status          string       `json:"status"`          // 会话状态
    CreatedAt       time.Time    `json:"created_at"`      // 创建时间
    LastActiveAt    time.Time    `json:"last_active_at"`  // 最后活跃时间
    ExpiresAt       time.Time    `json:"expires_at"`      // 过期时间
}
```

#### TokenPair - Token对
```go
type TokenPair struct {
    AccessToken  string `json:"access_token"`  // 访问令牌
    RefreshToken string `json:"refresh_token"` // 刷新令牌
    ExpiresIn    int64  `json:"expires_in"`    // 过期时间(秒)
}
```

## 3. 服务接口定义

### 3.1 SessionService 接口
```go
type SessionService interface {
    // 会话创建
    CreateSession(ctx context.Context, userID uint64, loginInfo *LoginInfo) (*SessionInfo, error)
    
    // 会话验证
    ValidateSession(ctx context.Context, sessionID string) (*SessionInfo, error)
    ValidateToken(ctx context.Context, accessToken string) (*SessionInfo, error)
    
    // 会话更新
    UpdateLastActive(ctx context.Context, sessionID string) error
    RefreshSession(ctx context.Context, refreshToken string) (*TokenPair, error)
    
    // 会话注销
    RevokeSession(ctx context.Context, sessionID string) error
    RevokeUserSessions(ctx context.Context, userID uint64, excludeSessionID string) error
    RevokeAllSessions(ctx context.Context, userID uint64) error
    
    // 会话查询
    GetUserSessions(ctx context.Context, userID uint64) ([]*SessionInfo, error)
    GetActiveSessionCount(ctx context.Context, userID uint64) (int, error)
    
    // 清理任务
    CleanupExpiredSessions(ctx context.Context) error
    CleanupUserSessions(ctx context.Context, userID uint64, keepCount int) error
}
```

## 4. Redis数据结构

### 4.1 存储模式

#### 用户会话集合
```redis
Key: user:sessions:{user_id}
Type: SET
Value: [session_id1, session_id2, ...]
TTL: 7天
```

#### 会话详情
```redis
Key: session:{session_id}
Type: HASH
Fields: user_id, platform, device_id, ip_address, access_token, etc.
TTL: 24小时
```

#### Token映射
```redis
Key: token:access:{token_hash}
Type: STRING
Value: session_id
TTL: 24小时

Key: token:refresh:{token_hash}
Type: STRING
Value: session_id
TTL: 7天
```

#### 活跃会话排行
```redis
Key: sessions:active
Type: ZSET
Score: last_active_timestamp
Member: session_id
```

### 4.2 事件发布订阅
```redis
Channel: session:events
Message: {"type":"session_created","session_id":"sess_123","user_id":1001}
```

## 5. 核心方法实现

### 5.1 CreateSession - 创建会话
```go
func (s *sessionService) CreateSession(ctx context.Context, userID uint64, loginInfo *LoginInfo) (*SessionInfo, error)
```

**功能说明:**
- 检查用户会话数量限制
- 生成唯一会话ID和JWT令牌对
- 存储会话信息到Redis和数据库
- 建立Token到会话的映射关系
- 发布会话创建事件

**关键步骤:**
1. 会话数量控制（超限时清理旧会话）
2. 生成会话ID: `sess_` + UUID前16位
3. JWT令牌生成（包含用户ID、会话ID、平台信息）
4. Redis存储（HASH + SET + STRING映射）
5. 数据库持久化
6. 事件发布

### 5.2 ValidateToken - 验证Token
```go
func (s *sessionService) ValidateToken(ctx context.Context, accessToken string) (*SessionInfo, error)
```

**功能说明:**
- 通过Token哈希查找对应会话
- 验证JWT Token有效性
- 检查会话状态和过期时间
- 更新最后活跃时间

**验证流程:**
1. Token哈希计算（SHA256）
2. Redis查询Token映射获取SessionID
3. JWT Token解析和验证
4. 会话状态检查
5. 异步更新活跃时间

### 5.3 RefreshSession - 刷新会话
```go
func (s *sessionService) RefreshSession(ctx context.Context, refreshToken string) (*TokenPair, error)
```

**功能说明:**
- 验证RefreshToken有效性
- 生成新的Token对
- 更新Redis和数据库中的Token信息
- 清理旧Token映射

**刷新流程:**
1. RefreshToken验证
2. 生成新AccessToken和RefreshToken
3. 更新Redis会话信息
4. 重建Token映射关系
5. 异步更新数据库
6. 发布Token刷新事件

### 5.4 RevokeSession - 注销会话
```go
func (s *sessionService) RevokeSession(ctx context.Context, sessionID string) error
```

**功能说明:**
- 从Redis删除会话相关数据
- 更新数据库会话状态
- 发布会话注销事件

**注销步骤:**
1. 获取会话完整信息
2. 删除Redis中的会话数据
3. 清理Token映射关系
4. 更新数据库状态为"revoked"
5. 发布注销事件

## 6. 安全机制

### 6.1 Token安全
- JWT Token采用HMAC-SHA256签名
- Token内容包含过期时间验证
- 敏感Token仅存储SHA256哈希值
- 支持Token黑名单机制

### 6.2 会话安全
- 会话ID使用UUID确保唯一性
- 支持IP地址验证
- 设备指纹检查
- 异常活动检测

### 6.3 并发控制
- 用户最大会话数限制
- 自动清理过期会话
- 防止会话泄漏
- 支持强制下线

## 7. 性能优化

### 7.1 缓存策略
- Redis作为主要存储，数据库作为持久化备份
- 采用异步更新减少响应延迟
- 合理设置TTL避免内存浪费
- 使用Pipeline批量操作

### 7.2 清理机制
- 定时清理过期会话
- 用户会话数量限制
- LRU策略淘汰旧会话
- 批量删除优化性能

## 8. 使用示例

### 8.1 创建会话服务实例
```go
sessionService := service.NewSessionService(
    redisCache,           // Redis缓存接口
    sessionRepo,          // 会话存储库
    "jwt-secret-key",     // JWT密钥
    24*time.Hour,         // AccessToken过期时间
    7*24*time.Hour,       // RefreshToken过期时间
    5,                    // 每用户最大会话数
)
```

### 8.2 用户登录创建会话
```go
loginInfo := &service.LoginInfo{
    Platform:    "app",
    DeviceID:    "iPhone_12_ABC123",
    DeviceType:  "ios",
    DeviceName:  "John's iPhone",
    AppVersion:  "1.0.0",
    UserAgent:   "HuinongApp/1.0.0",
    IPAddress:   "192.168.1.100",
    Location:    "北京市朝阳区",
    LoginMethod: "password",
}

sessionInfo, err := sessionService.CreateSession(ctx, userID, loginInfo)
if err != nil {
    return err
}

// 返回Token给客户端
accessToken := sessionInfo.TokenInfo.AccessToken
refreshToken := sessionInfo.TokenInfo.RefreshToken
```

### 8.3 中间件集成验证
```go
// 在Gin路由中使用会话认证中间件
authMiddleware := middleware.NewSessionAuthMiddleware(sessionService)
router.Use(authMiddleware.RequireAuth())

// 受保护的路由
router.GET("/api/user/profile", func(c *gin.Context) {
    userID := c.GetUint64("user_id")
    sessionID := c.GetString("session_id")
    // 处理业务逻辑
})
```

### 8.4 Token刷新
```go
tokenPair, err := sessionService.RefreshSession(ctx, refreshToken)
if err != nil {
    // 刷新失败，需要重新登录
    return err
}

// 返回新Token
newAccessToken := tokenPair.AccessToken
newRefreshToken := tokenPair.RefreshToken
```

## 9. 错误处理

### 9.1 常见错误
- `Token不存在或已过期`: Token映射不存在或已失效
- `会话已失效`: 会话状态非active或已过期
- `无效的Token类型`: JWT中type字段不匹配
- `会话不存在`: Redis中无对应会话数据

### 9.2 错误恢复
- Token过期自动清理
- 会话状态同步修复
- 重复登录覆盖机制
- 异常会话强制注销

## 10. 监控与日志

### 10.1 关键指标
- 会话创建/验证/注销QPS
- Redis连接池使用率
- 会话命中率和过期率
- Token刷新成功率

### 10.2 日志记录
- 会话生命周期事件
- 安全异常检测
- 性能指标统计
- 错误详情追踪

## 11. 部署配置

### 11.1 Redis配置
```yaml
session:
  redis:
    host: "127.0.0.1"
    port: 6379
    password: ""
    db: 1
    pool_size: 20
  settings:
    access_token_ttl: 24h
    refresh_token_ttl: 168h
    max_sessions_per_user: 5
    cleanup_interval: 1h
```

### 11.2 定时任务
```go
// 启动会话清理定时任务
go func() {
    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()
    
    for range ticker.C {
        err := sessionService.CleanupExpiredSessions(context.Background())
        if err != nil {
            log.Printf("清理过期会话失败: %v", err)
        }
    }
}()
```

通过以上设计和实现，会话管理服务提供了完整的分布式会话保持解决方案，确保多后端实例间的用户状态一致性和系统的高可用性。