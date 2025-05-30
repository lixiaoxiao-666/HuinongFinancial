# Router 路由管理模块

## 📋 模块概述

Router模块负责整个API的路由配置和中间件管理，实现了基于Redis会话的认证系统、权限控制、请求日志记录等功能。

## 🔗 路由结构

### 公开API（无需认证）
```go
// 健康检查
GET    /health

// 系统信息
GET    /api/public/version
GET    /api/public/configs

// 用户认证
POST   /api/auth/register
POST   /api/auth/login
POST   /api/auth/refresh

// OA管理员认证（无需验证码）
POST   /api/oa/auth/login
```

### 需要认证的API

#### 内部API（Dify工作流专用）
```go
POST   /api/internal/dify/loan/get-application-details
POST   /api/internal/dify/loan/submit-assessment
POST   /api/internal/dify/machine/get-rental-details
POST   /api/internal/dify/credit/query
```

#### 用户API（需要用户认证）
```go
// 用户资料管理
GET    /api/user/profile
PUT    /api/user/profile
PUT    /api/user/password
POST   /api/user/logout

// 会话管理
GET    /api/user/session/info
POST   /api/user/session/revoke-others

// 用户认证
POST   /api/user/auth/real-name
POST   /api/user/auth/bank-card

// 用户标签
GET    /api/user/tags
POST   /api/user/tags
DELETE /api/user/tags/:tag_key

// 贷款相关
GET    /api/user/loan/products
GET    /api/user/loan/products/:id
POST   /api/user/loan/applications
GET    /api/user/loan/applications
GET    /api/user/loan/applications/:id
DELETE /api/user/loan/applications/:id

// 文件上传
POST   /api/user/files/upload
POST   /api/user/files/upload/batch
GET    /api/user/files/:id
DELETE /api/user/files/:id

// 农机相关
POST   /api/user/machines
GET    /api/user/machines
GET    /api/user/machines/search
GET    /api/user/machines/:id
POST   /api/user/machines/:id/orders

// 农机订单
GET    /api/user/orders
PUT    /api/user/orders/:id/confirm
POST   /api/user/orders/:id/pay
PUT    /api/user/orders/:id/complete
PUT    /api/user/orders/:id/cancel
POST   /api/user/orders/:id/rate

// 专家咨询
POST   /api/user/consultations
GET    /api/user/consultations
```

#### 公共内容API（可选认证）
```go
// 文章相关
GET    /api/content/articles
GET    /api/content/articles/featured
GET    /api/content/articles/:id
GET    /api/content/categories

// 专家相关
GET    /api/content/experts
GET    /api/content/experts/:id
```

#### 管理员API（需要管理员认证）
```go
// 用户管理
GET    /api/admin/users
GET    /api/admin/users/statistics
GET    /api/admin/users/:user_id/auth-status

// 会话管理
GET    /api/admin/sessions/active
POST   /api/admin/sessions/cleanup
DELETE /api/admin/sessions/:session_id

// 贷款审批管理
GET    /api/admin/loans/applications
GET    /api/admin/loans/applications/:id
POST   /api/admin/loans/applications/:id/approve
POST   /api/admin/loans/applications/:id/reject
POST   /api/admin/loans/applications/:id/return
POST   /api/admin/loans/applications/:id/start-review
POST   /api/admin/loans/applications/:id/retry-ai
GET    /api/admin/loans/statistics

// 认证审核管理
GET    /api/admin/auth/list
GET    /api/admin/auth/:id
POST   /api/admin/auth/:id/review
POST   /api/admin/auth/batch-review
GET    /api/admin/auth/statistics
GET    /api/admin/auth/export

// 内容管理
POST   /api/admin/content/articles
PUT    /api/admin/content/articles/:id
DELETE /api/admin/content/articles/:id
POST   /api/admin/content/articles/:id/publish
POST   /api/admin/content/categories
PUT    /api/admin/content/categories/:id
DELETE /api/admin/content/categories/:id
POST   /api/admin/content/experts
PUT    /api/admin/content/experts/:id
DELETE /api/admin/content/experts/:id

// 系统管理
GET    /api/admin/system/config
PUT    /api/admin/system/config
GET    /api/admin/system/configs
GET    /api/admin/system/health
GET    /api/admin/system/statistics
```

#### OA后台API（需要OA管理员认证）
```go
// OA用户管理
GET    /api/oa/users
GET    /api/oa/users/:user_id
PUT    /api/oa/users/:user_id/status
POST   /api/oa/users/batch-operation

// OA农机设备管理
GET    /api/oa/machines
GET    /api/oa/machines/:machine_id

// OA工作台和数据分析
GET    /api/oa/dashboard
GET    /api/oa/dashboard/overview
GET    /api/oa/dashboard/risk-monitoring
```

## 🔐 认证中间件详解

### Redis会话认证中间件
```go
// 必须认证
sessionAuthMiddleware.RequireAuth()

// 可选认证  
sessionAuthMiddleware.OptionalAuth()

// 管理员认证
sessionAuthMiddleware.AdminAuth()

// 登出处理
sessionAuthMiddleware.Logout()

// Token刷新
sessionAuthMiddleware.RefreshToken()

// 会话信息
sessionAuthMiddleware.SessionInfo()

// 注销其他会话
sessionAuthMiddleware.RevokeOtherSessions()
```

### 特殊认证
```go
// Dify专用认证
middleware.DifyAuthMiddleware(config.DifyAPIToken)
```

## 📊 路由统计

| 路由组 | 认证方式 | 接口数量 | 状态 |
|--------|----------|----------|------|
| /health | 无需认证 | 1 | ✅ |
| /api/public | 无需认证 | 2 | ✅ |
| /api/internal | Dify专用 | 4 | ✅ |
| /api/auth | 无需认证 | 3 | ✅ |
| /api/user | Redis会话 | 25+ | ✅ |
| /api/content | Redis可选 | 6 | ✅ |
| /api/admin | Redis管理员 | 30+ | ✅ |
| /api/oa | Redis管理员 | 待实现 | 🚧 |
| Swagger | 无需认证 | 1 | ✅ |

## 🔧 配置和使用

### RouterConfig结构
```go
type RouterConfig struct {
    UserService    service.UserService
    SessionService service.SessionService      // Redis会话服务
    LoanService    service.LoanService
    MachineService service.MachineService
    ArticleService service.ContentService
    ExpertService  service.ContentService
    FileService    service.SystemService
    SystemService  service.SystemService
    OAService      service.OAService
    DifyService    service.DifyService
    JWTSecret      string                      // 用于JWT签名
    DifyAPIToken   string                      // Dify专用Token
}
```

### 中间件初始化
```go
// 统一使用Redis会话认证
sessionAuthMiddleware := middleware.NewSessionAuthMiddleware(config.SessionService)
```

## 🚀 特性优势

### 1. 统一认证架构
- **单一认证策略**: 除Dify外全部使用Redis会话
- **分布式支持**: 多后端实例共享会话状态
- **实时控制**: 支持强制下线、会话监控

### 2. 高性能设计
- **Redis缓存**: 毫秒级会话验证
- **连接池**: 高并发支持
- **自动清理**: 过期会话自动清理

### 3. 安全保障
- **Token哈希**: 安全存储访问令牌
- **设备绑定**: 防止令牌滥用
- **IP验证**: 可选的IP地址验证
- **会话限制**: 单用户最大会话数控制

### 4. 开发友好
- **清晰分组**: 功能模块化路由设计
- **文档完整**: 每个接口都有详细说明
- **测试支持**: 提供测试专用路由函数

## 📝 维护说明

### 添加新路由
1. 在相应的路由组中添加新路由
2. 选择合适的认证中间件
3. 更新本文档和API文档

### 认证策略变更
- **用户相关**: 使用 `sessionAuthMiddleware.RequireAuth()`
- **管理员功能**: 使用 `sessionAuthMiddleware.AdminAuth()`
- **可选认证**: 使用 `sessionAuthMiddleware.OptionalAuth()`
- **Dify专用**: 保持 `middleware.DifyAuthMiddleware()`

通过这种统一的认证架构，系统实现了高效、安全、可扩展的分布式会话管理！ 🎉 