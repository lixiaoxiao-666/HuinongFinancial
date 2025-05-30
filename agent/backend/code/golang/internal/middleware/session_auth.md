# 会话认证中间件 - 代码文档

## 1. 概述

`session_auth.go` 实现了基于Redis会话管理的认证中间件，替代传统的JWT认证中间件，提供更强大的会话控制能力和分布式支持。

## 2. 核心功能

### 2.1 主要特性
- 🔐 **基于会话的认证**: 集成Redis会话管理，支持实时会话状态检查
- 📱 **多平台支持**: 区分不同平台（app, web, oa），提供差异化认证
- ⚡ **高性能验证**: 毫秒级Token验证，支持高并发请求
- 🔄 **自动续期**: 自动更新会话活跃时间
- 🛡️ **安全控制**: 支持强制下线、角色权限、管理员认证等

### 2.2 中间件类型

#### SessionAuthMiddleware - 会话认证中间件
```go
type SessionAuthMiddleware struct {
    sessionService service.SessionService
}
```

## 3. 认证中间件方法

### 3.1 RequireAuth - 强制认证
```go
func (m *SessionAuthMiddleware) RequireAuth() gin.HandlerFunc
```

**功能说明:**
- 强制要求用户认证
- 验证Authorization头中的Bearer Token
- 检查会话有效性和状态
- 将用户信息存储到Gin上下文

**使用场景:**
- 需要登录才能访问的API
- 用户个人信息相关接口
- 核心业务功能接口

**示例代码:**
```go
// 创建认证中间件
authMiddleware := middleware.NewSessionAuthMiddleware(sessionService)

// 应用到路由组
api := router.Group("/api")
api.Use(authMiddleware.RequireAuth())
{
    api.GET("/user/profile", getUserProfile)
    api.POST("/user/update", updateUserProfile)
}
```

**上下文存储:**
- `user_id`: 用户ID (uint64)
- `session_id`: 会话ID (string)
- `platform`: 登录平台 (string)
- `device_info`: 设备信息 (*DeviceInfo)
- `network_info`: 网络信息 (*NetworkInfo)

### 3.2 OptionalAuth - 可选认证
```go
func (m *SessionAuthMiddleware) OptionalAuth() gin.HandlerFunc
```

**功能说明:**
- 不强制要求认证，有Token则验证
- 验证成功时存储用户信息到上下文
- 验证失败时继续执行，不阻断请求

**使用场景:**
- 公开内容，登录用户有额外权限
- 首页推荐内容（登录用户个性化）
- 商品列表（登录用户显示收藏状态）

**示例代码:**
```go
// 可选认证的路由
public := router.Group("/public")
public.Use(authMiddleware.OptionalAuth())
{
    public.GET("/products", getProductList)    // 登录用户显示个性化内容
    public.GET("/articles", getArticleList)    // 登录用户显示收藏状态
}
```

### 3.3 AdminAuth - 管理员认证
```go
func (m *SessionAuthMiddleware) AdminAuth() gin.HandlerFunc
```

**功能说明:**
- 验证管理员身份
- 检查登录平台是否为"oa"
- 提供最高级别的访问控制

**使用场景:**
- 后台管理系统接口
- 系统配置修改
- 用户管理操作

**示例代码:**
```go
// 管理员专用路由
admin := router.Group("/admin")
admin.Use(authMiddleware.AdminAuth())
{
    admin.GET("/users", getUserList)
    admin.POST("/users/freeze", freezeUser)
    admin.GET("/system/config", getSystemConfig)
}
```

**上下文存储:**
- `user_id`: 管理员用户ID
- `session_id`: 管理员会话ID
- `platform`: 固定为"oa"
- `is_admin`: 设置为true

### 3.4 RequireRole - 角色权限认证
```go
func (m *SessionAuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc
```

**功能说明:**
- 基于角色的访问控制（RBAC）
- 支持多角色权限检查
- 灵活的权限管理

**使用场景:**
- 分级管理系统
- 不同角色不同权限
- 精细化权限控制

**示例代码:**
```go
// 角色权限路由
api.Use(authMiddleware.RequireRole("admin", "manager")).POST("/loan/approve", approveLoan)
api.Use(authMiddleware.RequireRole("auditor")).GET("/audit/list", getAuditList)
```

## 4. 会话操作中间件

### 4.1 RefreshToken - Token刷新
```go
func (m *SessionAuthMiddleware) RefreshToken() gin.HandlerFunc
```

**功能说明:**
- 处理Token刷新请求
- 验证RefreshToken有效性
- 返回新的Token对

**请求参数:**
- `refresh_token`: 刷新令牌（POST表单）

**响应格式:**
```json
{
    "code": 200,
    "message": "刷新成功",
    "data": {
        "access_token": "new_access_token",
        "refresh_token": "new_refresh_token",
        "expires_in": 86400
    }
}
```

**使用示例:**
```go
// Token刷新路由
router.POST("/auth/refresh", authMiddleware.RefreshToken())
```

### 4.2 Logout - 用户登出
```go
func (m *SessionAuthMiddleware) Logout() gin.HandlerFunc
```

**功能说明:**
- 注销当前用户会话
- 清理Redis中的会话数据
- 使当前Token立即失效

**前置条件:**
- 需要先通过认证中间件

**响应格式:**
```json
{
    "code": 200,
    "message": "注销成功"
}
```

**使用示例:**
```go
// 登出路由
api.Use(authMiddleware.RequireAuth()).POST("/auth/logout", authMiddleware.Logout())
```

### 4.3 SessionInfo - 获取会话信息
```go
func (m *SessionAuthMiddleware) SessionInfo() gin.HandlerFunc
```

**功能说明:**
- 获取用户所有活跃会话
- 显示设备信息和登录时间
- 支持会话管理功能

**响应格式:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": [
        {
            "session_id": "sess_abc123",
            "platform": "app",
            "device_info": {
                "device_type": "ios",
                "device_name": "John's iPhone"
            },
            "created_at": "2024-01-15T10:30:00Z",
            "last_active_at": "2024-01-15T14:25:30Z"
        }
    ]
}
```

**使用示例:**
```go
// 会话信息路由
api.Use(authMiddleware.RequireAuth()).GET("/auth/sessions", authMiddleware.SessionInfo())
```

### 4.4 RevokeOtherSessions - 注销其他会话
```go
func (m *SessionAuthMiddleware) RevokeOtherSessions() gin.HandlerFunc
```

**功能说明:**
- 注销用户的其他会话
- 保留当前会话不变
- 实现"在其他设备强制下线"功能

**使用场景:**
- 账号安全管理
- 密码修改后强制重新登录
- 发现异常登录时的安全措施

**响应格式:**
```json
{
    "code": 200,
    "message": "注销其他会话成功"
}
```

**使用示例:**
```go
// 强制下线路由
api.Use(authMiddleware.RequireAuth()).POST("/auth/revoke-others", authMiddleware.RevokeOtherSessions())
```

## 5. 集成使用指南

### 5.1 创建中间件实例
```go
// 初始化会话服务
sessionService := service.NewSessionService(
    redisCache,
    sessionRepo,
    "jwt-secret-key",
    24*time.Hour,
    7*24*time.Hour,
    5,
)

// 创建认证中间件
authMiddleware := middleware.NewSessionAuthMiddleware(sessionService)
```

### 5.2 路由配置示例
```go
func SetupRoutes(router *gin.Engine, authMiddleware *middleware.SessionAuthMiddleware) {
    // 公开路由（无需认证）
    public := router.Group("/public")
    {
        public.GET("/health", healthCheck)
        public.POST("/auth/login", userLogin)
        public.POST("/auth/register", userRegister)
    }

    // 可选认证路由
    open := router.Group("/open")
    open.Use(authMiddleware.OptionalAuth())
    {
        open.GET("/products", getProductList)
        open.GET("/articles", getArticleList)
    }

    // 认证路由
    api := router.Group("/api")
    api.Use(authMiddleware.RequireAuth())
    {
        // 用户相关
        user := api.Group("/user")
        {
            user.GET("/profile", getUserProfile)
            user.PUT("/profile", updateUserProfile)
            user.POST("/avatar", uploadAvatar)
        }

        // 认证管理
        auth := api.Group("/auth")
        {
            auth.POST("/logout", authMiddleware.Logout())
            auth.GET("/sessions", authMiddleware.SessionInfo())
            auth.POST("/revoke-others", authMiddleware.RevokeOtherSessions())
        }

        // 业务功能
        business := api.Group("/business")
        {
            business.GET("/loans", getUserLoans)
            business.POST("/loans/apply", applyLoan)
        }
    }

    // 管理员路由
    admin := router.Group("/admin")
    admin.Use(authMiddleware.AdminAuth())
    {
        admin.GET("/users", getUserList)
        admin.POST("/users/:id/freeze", freezeUser)
        admin.GET("/system/stats", getSystemStats)
    }

    // Token刷新（独立路由）
    router.POST("/auth/refresh", authMiddleware.RefreshToken())
}
```

### 5.3 错误处理
```go
// 自定义错误处理中间件
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        // 处理认证错误
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            switch err.Type {
            case gin.ErrorTypePublic:
                c.JSON(http.StatusUnauthorized, gin.H{
                    "code":    401,
                    "message": "认证失败",
                    "error":   err.Error(),
                })
            case gin.ErrorTypePrivate:
                c.JSON(http.StatusInternalServerError, gin.H{
                    "code":    500,
                    "message": "服务器内部错误",
                })
            }
        }
    }
}
```

## 6. 最佳实践

### 6.1 性能优化
- 使用连接池减少Redis连接开销
- 设置合理的超时时间（5秒）
- 异步更新活跃时间避免阻塞
- 缓存用户权限信息

### 6.2 安全建议
- 定期轮换JWT密钥
- 监控异常登录行为
- 实施IP白名单机制
- 限制并发会话数量

### 6.3 监控告警
- 认证失败率监控
- 会话创建/销毁QPS
- Redis连接状态
- Token刷新成功率

### 6.4 日志记录
```go
// 认证日志中间件
func AuthLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        c.Next()
        
        // 记录认证相关日志
        if userID, exists := c.Get("user_id"); exists {
            log.Printf("Auth Success - User: %v, Path: %s, Duration: %v", 
                userID, path, time.Since(start))
        } else if c.Writer.Status() == 401 {
            log.Printf("Auth Failed - Path: %s, IP: %s, Duration: %v", 
                path, c.ClientIP(), time.Since(start))
        }
    }
}
```

## 7. 故障排查

### 7.1 常见问题
1. **Token验证失败**
   - 检查JWT密钥配置
   - 确认Token格式正确
   - 验证Token是否过期

2. **会话不存在**
   - 检查Redis连接状态
   - 确认会话是否被清理
   - 验证SessionID格式

3. **权限验证失败**
   - 检查用户角色配置
   - 确认权限策略正确
   - 验证角色权限映射

### 7.2 调试技巧
```go
// 启用调试模式
func DebugAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 打印请求头信息
        auth := c.GetHeader("Authorization")
        log.Printf("Debug - Auth Header: %s", auth)
        
        c.Next()
        
        // 打印上下文信息
        if userID, exists := c.Get("user_id"); exists {
            log.Printf("Debug - User ID: %v", userID)
        }
    }
}
```

通过以上会话认证中间件的设计和实现，系统可以提供完整的认证授权解决方案，支持多种认证模式和灵活的权限控制，确保系统的安全性和可扩展性。 