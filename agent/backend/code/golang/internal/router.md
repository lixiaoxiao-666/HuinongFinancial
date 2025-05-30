# 路由系统文档

## 📋 概述

数字惠农后端路由系统基于Gin框架构建，实现了完整的RESTful API接口设计。**整个系统已统一使用Redis分布式会话认证**，只有Dify平台保持专用的Token认证。

## 🔗 路由架构

### 认证策略统一
- **Redis会话认证**: 所有用户、管理员、OA后台API统一使用
- **Dify专用认证**: 仅用于 `/api/internal/dify/*` 工作流接口
- **无需认证**: 公开API、健康检查、认证相关接口

## 📚 路由分组结构

### 1. 健康检查
```
GET /health
```
- 无需认证
- 返回服务状态和版本信息

### 2. 公开API (`/api/public`)
```
GET /api/public/version          # 获取系统版本
GET /api/public/configs          # 获取公开配置
```
- 无需认证
- 面向公众的系统信息

### 3. Dify工作流API (`/api/internal`)
```
POST /api/internal/dify/loan/get-application-details
POST /api/internal/dify/loan/submit-assessment
POST /api/internal/dify/machine/get-rental-details
POST /api/internal/dify/credit/query
```
- **使用Dify专用Token认证**
- 供Dify AI工作流调用
- 独立的认证机制

### 4. 用户认证API (`/api/auth`)
```
POST /api/auth/register          # 用户注册
POST /api/auth/login            # 用户登录
POST /api/auth/refresh          # Token刷新（Redis会话）
```
- 无需认证（登录前接口）
- Token刷新使用Redis会话管理

### 5. 用户API (`/api/user`) - **Redis会话认证**
```
# 用户资料
GET    /api/user/profile         # 获取用户资料
PUT    /api/user/profile         # 更新用户资料
PUT    /api/user/password        # 修改密码
POST   /api/user/logout          # 登出（Redis会话）

# 会话管理
GET    /api/user/session/info    # 获取会话信息
POST   /api/user/session/revoke-others  # 注销其他设备
GET    /api/user/session/list    # 获取会话列表

# 用户认证
POST   /api/user/auth/real-name  # 实名认证
POST   /api/user/auth/bank-card  # 银行卡认证

# 用户标签
GET    /api/user/tags           # 获取用户标签
POST   /api/user/tags           # 添加用户标签
DELETE /api/user/tags/:tag_key  # 删除用户标签

# 贷款功能
GET    /api/user/loan/products           # 获取贷款产品
GET    /api/user/loan/products/:id       # 获取产品详情
POST   /api/user/loan/applications       # 提交贷款申请
GET    /api/user/loan/applications       # 获取用户申请
GET    /api/user/loan/applications/:id   # 获取申请详情
DELETE /api/user/loan/applications/:id   # 取消申请

# 文件管理
POST   /api/user/files/upload           # 文件上传
POST   /api/user/files/upload/batch     # 批量上传
GET    /api/user/files/:id              # 获取文件
DELETE /api/user/files/:id              # 删除文件

# 农机管理
POST   /api/user/machines              # 注册农机
GET    /api/user/machines              # 获取用户农机
GET    /api/user/machines/search       # 搜索农机
GET    /api/user/machines/:id          # 获取农机详情
POST   /api/user/machines/:id/orders   # 创建订单

# 农机订单
GET    /api/user/orders                # 获取用户订单
PUT    /api/user/orders/:id/confirm    # 确认订单
POST   /api/user/orders/:id/pay        # 支付订单
PUT    /api/user/orders/:id/complete   # 完成订单
PUT    /api/user/orders/:id/cancel     # 取消订单
POST   /api/user/orders/:id/rate       # 评价订单

# 专家咨询
POST   /api/user/consultations         # 提交咨询
GET    /api/user/consultations         # 获取咨询列表
```

### 6. 公共内容API (`/api/content`) - **Redis可选认证**
```
# 文章相关
GET /api/content/articles              # 文章列表
GET /api/content/articles/featured     # 推荐文章
GET /api/content/articles/:id          # 文章详情
GET /api/content/categories            # 文章分类

# 专家相关
GET /api/content/experts               # 专家列表
GET /api/content/experts/:id           # 专家详情
```
- 可选认证：登录用户获得个性化内容

### 7. 管理员API (`/api/admin`) - **Redis管理员认证**
```
# 用户管理
GET /api/admin/users                          # 用户列表
GET /api/admin/users/statistics               # 用户统计
GET /api/admin/users/:user_id/auth-status     # 用户认证状态

# 会话管理
GET    /api/admin/sessions/active             # 活跃会话列表
POST   /api/admin/sessions/cleanup            # 清理过期会话
DELETE /api/admin/sessions/:session_id        # 强制注销会话

# 贷款审批
GET  /api/admin/loans/applications            # 申请列表
GET  /api/admin/loans/applications/:id        # 申请详情
POST /api/admin/loans/applications/:id/approve # 批准申请
POST /api/admin/loans/applications/:id/reject  # 拒绝申请
POST /api/admin/loans/applications/:id/return  # 退回申请
POST /api/admin/loans/applications/:id/start-review # 开始审核
POST /api/admin/loans/applications/:id/retry-ai    # 重试AI评估
GET  /api/admin/loans/statistics              # 贷款统计

# 认证审核
GET  /api/admin/auth/list                     # 认证列表
GET  /api/admin/auth/:id                      # 认证详情
POST /api/admin/auth/:id/review               # 审核认证
POST /api/admin/auth/batch-review             # 批量审核
GET  /api/admin/auth/statistics               # 认证统计
GET  /api/admin/auth/export                   # 导出认证数据

# 内容管理
POST   /api/admin/content/articles            # 创建文章
PUT    /api/admin/content/articles/:id        # 更新文章
DELETE /api/admin/content/articles/:id        # 删除文章
POST   /api/admin/content/articles/:id/publish # 发布文章
POST   /api/admin/content/categories          # 创建分类
PUT    /api/admin/content/categories/:id      # 更新分类
DELETE /api/admin/content/categories/:id      # 删除分类
POST   /api/admin/content/experts             # 创建专家
PUT    /api/admin/content/experts/:id         # 更新专家
DELETE /api/admin/content/experts/:id         # 删除专家

# 系统管理
GET /api/admin/system/config                  # 获取配置
PUT /api/admin/system/config                  # 设置配置
GET /api/admin/system/configs                 # 获取所有配置
GET /api/admin/system/health                  # 健康检查
GET /api/admin/system/statistics              # 系统统计
```

### 8. OA后台API (`/api/oa`) - **Redis管理员认证**
```
# OA专用接口（待实现）
# 用户管理、角色管理、工作台等
```
- 检查platform为"oa"的管理员会话

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