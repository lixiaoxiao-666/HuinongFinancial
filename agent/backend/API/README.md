# 数字惠农系统 - API 接口文档总览

## 📋 系统概述

数字惠农系统是一个综合性的智慧农业服务平台，为农户提供金融服务、农机租赁、政策资讯、技术指导等一站式服务。系统采用微服务架构，支持移动端APP、Web端和OA后台管理三大平台。

### 🏗️ 系统架构

```
┌─────────────────────────────────────────────────────────────────────┐
│                          前端应用层                                    │
├──────────────────┬──────────────────┬──────────────────┬────────────┤
│    移动端APP      │     Web端        │    OA管理后台     │ Dify工作流  │
│   (iOS/Android)   │   (Vue3 + TS)    │   (Vue3 + TS)    │   AI系统   │
└──────────────────┴──────────────────┴──────────────────┴────────────┘
                              │
┌─────────────────────────────────────────────────────────────────────┐
│                      API网关层 (Gin + 中间件)                         │
│     统一路由、Redis会话认证、平台隔离、限流、监控、日志                    │
└─────────────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────────────┐
│                        业务服务层                                      │
├──────────┬──────────┬──────────┬──────────┬──────────┬─────────────┤
│ 用户管理   │ 贷款服务   │ 农机租赁   │ 内容管理   │ OA管理   │ 系统管理     │
│ Service   │ Service   │ Service   │ Service   │ Service   │ Service     │
├──────────┼──────────┼──────────┼──────────┼──────────┼─────────────┤
│ 会话管理   │ 文件管理   │ 任务管理   │ Dify集成   │ 监控统计  │ 日志审计     │
│ Service   │ Service   │ Service   │ Service   │ Service   │ Service     │
└──────────┴──────────┴──────────┴──────────┴──────────┴─────────────┘
                              │
┌─────────────────────────────────────────────────────────────────────┐
│                        数据存储层                                      │
├──────────┬──────────┬──────────┬──────────┬──────────┬─────────────┤
│  MySQL    │  Redis    │ MongoDB   │   OSS     │  Kafka   │ Elasticsearch│
│ 主数据库   │ 缓存/会话  │ 日志/统计  │ 文件存储   │ 消息队列 │ 搜索引擎     │
└──────────┴──────────┴──────────┴──────────┴──────────┴─────────────┘
```

### 🔧 技术栈

- **后端框架**: Golang + Gin + GORM
- **认证系统**: Redis分布式会话管理
- **数据库**: MySQL 8.0 + Redis 6.0
- **AI集成**: Dify工作流平台
- **文件存储**: 阿里云OSS / MinIO
- **消息队列**: Apache Kafka
- **搜索引擎**: Elasticsearch (可选)
- **监控**: Prometheus + Grafana
- **部署**: Docker + Kubernetes

---

## 📚 API模块列表

### 🔐 1. 会话管理模块 [`session_management.md`](./session_management.md)
基于Redis的高性能分布式会话管理系统，支持多平台、多后端实例的会话同步。

**核心功能:**
- 🔑 **智能认证**: 用户登录/登出/自动Token刷新
- 🌐 **分布式会话**: 多后端实例共享用户会话状态
- 📱 **多端登录**: 支持APP/Web/OA同时登录管理
- 🔒 **安全控制**: 强制下线/设备绑定/异常检测
- 📊 **智能设备识别**: 自动识别浏览器类型，优化设备显示
- 👥 **会话监控**: 管理员实时会话管理和统计
- ⚡ **高性能**: Redis毫秒级会话验证，支持高并发

**主要接口:**
```http
# 惠农用户认证
POST   /api/auth/login           # 用户登录
POST   /api/auth/register        # 用户注册
POST   /api/auth/refresh         # Token刷新
GET    /api/auth/validate        # 会话验证
POST   /api/user/logout          # 用户登出

# OA系统认证
POST   /api/oa/auth/login        # OA管理员登录
POST   /api/oa/auth/refresh      # OA Token刷新
GET    /api/oa/auth/validate     # OA会话验证
POST   /api/oa/auth/logout       # OA登出

# 会话管理
GET    /api/user/session/info    # 获取用户会话信息
GET    /api/user/session/list    # 获取用户所有会话
POST   /api/user/session/revoke-others # 注销其他设备

# 管理员会话监控
GET    /api/oa/admin/sessions/statistics # 会话统计
GET    /api/oa/admin/sessions/active     # 活跃会话列表
POST   /api/oa/admin/sessions/cleanup    # 清理过期会话
DELETE /api/oa/admin/sessions/{id}       # 强制注销会话
```

### 👤 2. 用户管理模块 [`user_management.md`](./user_management.md)
完整的用户生命周期管理，支持多种用户类型和身份认证。

**核心功能:**
- 📝 **用户注册登录** (手机/密码/验证码)
- 🆔 **身份认证** (实名/银行卡/征信)
- 📋 **信息管理** (个人资料/地址/头像)
- 🏷️ **用户标签** (行为分析/精准推荐)
- 👥 **后台管理** (用户列表/状态管理/批量操作)

**主要接口:**
```http
# 用户基础管理
GET    /api/user/profile         # 获取用户信息
PUT    /api/user/profile         # 更新用户信息
PUT    /api/user/password        # 修改密码

# 用户标签系统
GET    /api/user/tags            # 获取用户标签
POST   /api/user/tags            # 添加用户标签
DELETE /api/user/tags/{tag_key}  # 删除用户标签

# 身份认证
POST   /api/user/auth/real-name  # 实名认证申请
POST   /api/user/auth/bank-card  # 银行卡认证申请

# OA后台用户管理
GET    /api/oa/admin/users       # 获取用户列表
GET    /api/oa/admin/users/statistics # 用户统计
GET    /api/oa/admin/users/{id}  # 获取用户详情
PUT    /api/oa/admin/users/{id}/status # 更新用户状态
POST   /api/oa/admin/users/batch-operation # 批量操作用户
GET    /api/oa/admin/users/{id}/auth-status # 获取认证状态
```

### 💰 3. 贷款管理模块 [`loan_management.md`](./loan_management.md)
智能化贷款服务，集成AI风险评估和自动化审批流程。

**核心功能:**
- 🏦 **贷款产品管理** (产品目录/详情查询)
- 📋 **申请流程** (在线申请/材料上传/状态跟踪)
- 🤖 **AI智能审批** (风险评估/额度计算)
- 👨‍💼 **管理员审批** (申请审核/状态管理/统计分析)
- 📊 **数据统计** (申请统计/风险监控)

**主要接口:**
```http
# 用户端贷款功能
GET    /api/user/loan/products   # 获取贷款产品
GET    /api/user/loan/products/{id} # 获取产品详情
POST   /api/user/loan/applications # 提交贷款申请
GET    /api/user/loan/applications # 获取用户申请列表
GET    /api/user/loan/applications/{id} # 获取申请详情
DELETE /api/user/loan/applications/{id} # 取消申请

# OA后台贷款管理
GET    /api/oa/admin/loans/applications # 获取申请列表
GET    /api/oa/admin/loans/applications/{id} # 获取申请详情
POST   /api/oa/admin/loans/applications/{id}/approve # 批准申请
POST   /api/oa/admin/loans/applications/{id}/reject  # 拒绝申请
POST   /api/oa/admin/loans/applications/{id}/return  # 退回申请
POST   /api/oa/admin/loans/applications/{id}/start-review # 开始审核
POST   /api/oa/admin/loans/applications/{id}/retry-ai # 重试AI评估
GET    /api/oa/admin/loans/statistics # 获取贷款统计
```

### 🚜 4. 农机租赁模块 [`machine_rental.md`](./machine_rental.md)
全流程农机租赁服务，支持设备搜索、预约、使用跟踪和归还管理。

**核心功能:**
- 🔍 **设备管理** (农机注册/搜索/详情查询)
- 📅 **订单管理** (创建订单/支付/确认/完成)
- 💼 **流程控制** (订单状态跟踪/评价系统)
- 👨‍💼 **后台管理** (设备审核/状态监控)

**主要接口:**
```http
# 用户端农机功能
POST   /api/user/machines        # 注册农机设备
GET    /api/user/machines        # 获取用户农机列表
GET    /api/user/machines/search # 搜索农机设备
GET    /api/user/machines/{id}   # 获取农机详情
POST   /api/user/machines/{id}/orders # 创建租赁订单

# 订单管理
GET    /api/user/orders          # 获取用户订单列表
PUT    /api/user/orders/{id}/confirm  # 确认订单
POST   /api/user/orders/{id}/pay      # 支付订单
PUT    /api/user/orders/{id}/complete # 完成订单
PUT    /api/user/orders/{id}/cancel   # 取消订单
POST   /api/user/orders/{id}/rate     # 评价订单

# OA后台农机管理
GET    /api/oa/admin/machines    # 获取农机列表
GET    /api/oa/admin/machines/{id} # 获取农机详情
```

### 📰 5. 内容管理模块 [`content_management.md`](./content_management.md)
丰富的农业内容服务，提供资讯、政策、专家咨询等信息化支持。

**核心功能:**
- 📰 **文章管理** (农业资讯/技术文章/分类管理)
- 👨‍🎓 **专家咨询** (专家信息/在线咨询)
- 🔔 **公告管理** (系统公告/通知发布)
- 🏷️ **内容分类** (分类管理/标签系统)

**主要接口:**
```http
# 公共内容API (可选认证)
GET    /api/content/articles     # 获取文章列表
GET    /api/content/articles/featured # 获取推荐文章
GET    /api/content/articles/{id} # 获取文章详情
GET    /api/content/categories   # 获取文章分类
GET    /api/content/experts      # 获取专家列表
GET    /api/content/experts/{id} # 获取专家详情

# 用户端咨询功能
POST   /api/user/consultations   # 提交咨询问题
GET    /api/user/consultations   # 获取咨询记录

# OA后台内容管理
POST   /api/oa/admin/content/articles # 创建文章
PUT    /api/oa/admin/content/articles/{id} # 更新文章
DELETE /api/oa/admin/content/articles/{id} # 删除文章
POST   /api/oa/admin/content/articles/{id}/publish # 发布文章

# 公告管理
GET    /api/oa/admin/content/announcements # 获取公告列表
POST   /api/oa/admin/content/announcements # 创建公告
PUT    /api/oa/admin/content/announcements/{id} # 更新公告
DELETE /api/oa/admin/content/announcements/{id} # 删除公告

# 分类和专家管理
POST   /api/oa/admin/content/categories # 创建分类
PUT    /api/oa/admin/content/categories/{id} # 更新分类
DELETE /api/oa/admin/content/categories/{id} # 删除分类
POST   /api/oa/admin/content/experts   # 创建专家
PUT    /api/oa/admin/content/experts/{id} # 更新专家
DELETE /api/oa/admin/content/experts/{id} # 删除专家
```

### 🏢 6. OA后台管理模块 [`oa_management.md`](./oa_management.md)
全面的后台管理功能，支持业务审批、数据统计、系统配置等管理需求。

**核心功能:**
- 👥 **用户管理** (用户列表/状态管理/权限控制)
- 📋 **认证审核** (实名认证/银行卡审核/批量处理)
- 📊 **数据统计** (业务报表/风险监控/用户分析)
- ⚙️ **系统配置** (参数设置/健康检查)
- 📊 **工作台** (仪表盘/风险监控/概览统计)

**主要接口:**
```http
# OA用户个人功能
GET    /api/oa/user/profile      # 获取OA用户信息
PUT    /api/oa/user/profile      # 更新OA用户信息
PUT    /api/oa/user/password     # 修改密码
GET    /api/oa/user/loan/applications # 查看个人申请

# 认证审核管理
GET    /api/oa/admin/auth/list   # 获取认证申请列表
GET    /api/oa/admin/auth/{id}   # 获取认证详情
POST   /api/oa/admin/auth/{id}/review # 审核认证申请
POST   /api/oa/admin/auth/batch-review # 批量审核
GET    /api/oa/admin/auth/statistics # 认证统计
GET    /api/oa/admin/auth/export # 导出认证数据

# 工作台和统计
GET    /api/oa/admin/dashboard   # 获取工作台数据
GET    /api/oa/admin/dashboard/overview # 获取业务概览
GET    /api/oa/admin/dashboard/risk-monitoring # 风险监控数据
```

### 📁 7. 文件管理模块 [`file_management.md`](./file_management.md)
安全高效的文件上传下载服务，支持多种文件类型和批量操作。

**核心功能:**
- 📤 **文件上传** (单文件/批量上传/类型验证)
- 📥 **文件下载** (安全下载/权限控制)
- 🗂️ **文件管理** (文件删除/统计信息)
- 🔒 **安全控制** (文件类型验证/存储配额/权限检查)

**主要接口:**
```http
POST   /api/user/files/upload    # 单文件上传
POST   /api/user/files/upload/batch # 批量文件上传
GET    /api/user/files/{id}      # 获取文件信息/下载
DELETE /api/user/files/{id}      # 删除文件
```

### 📋 8. 任务管理模块 [`task_management.md`](./task_management.md)
全面的任务流程管理，支持任务分配、处理、进度跟踪。

**核心功能:**
- 📝 **任务管理** (创建/更新/删除/查询)
- 👥 **任务分配** (指派/重新分配/取消分配)
- 📊 **进度跟踪** (状态更新/进度查询)
- 📋 **待办管理** (待处理任务/优先级管理)

**主要接口:**
```http
GET    /api/oa/admin/tasks       # 获取任务列表
POST   /api/oa/admin/tasks       # 创建任务
GET    /api/oa/admin/tasks/{id}  # 获取任务详情
PUT    /api/oa/admin/tasks/{id}  # 更新任务
DELETE /api/oa/admin/tasks/{id}  # 删除任务
GET    /api/oa/admin/tasks/pending # 获取待处理任务

# 任务处理操作
POST   /api/oa/admin/tasks/{id}/process  # 处理任务
POST   /api/oa/admin/tasks/{id}/assign   # 分配任务
POST   /api/oa/admin/tasks/{id}/unassign # 取消分配
POST   /api/oa/admin/tasks/{id}/reassign # 重新分配
GET    /api/oa/admin/tasks/{id}/progress # 获取任务进度
```

### 🤖 9. Dify工作流集成模块 [`dify_integration.md`](./dify_integration.md)
与AI平台深度集成，提供智能风险评估和数据分析服务。

**核心功能:**
- 🔍 **数据查询** (贷款申请详情/农机租赁信息)
- 🤖 **AI评估** (风险评估提交/征信查询)
- 🔐 **安全认证** (专用Token认证/数据脱敏)
- 📊 **智能分析** (数据挖掘/预测分析)

**主要接口:**
```http
# 贷款相关AI接口
POST   /api/internal/dify/loan/get-application-details # 获取贷款申请详情
POST   /api/internal/dify/loan/submit-assessment # 提交风险评估结果

# 农机相关AI接口
POST   /api/internal/dify/machine/get-rental-details # 获取农机租赁详情

# 征信相关AI接口
POST   /api/internal/dify/credit/query # 征信查询
```

### ⚙️ 10. 系统管理模块 [`system_management.md`](./system_management.md)
系统级配置和监控功能，确保系统稳定运行。

**核心功能:**
- 🔧 **系统配置** (参数配置/公共配置)
- 💊 **健康检查** (服务状态/系统监控)
- 📊 **统计分析** (系统统计/性能监控)
- 📋 **版本管理** (版本信息/更新记录)

**主要接口:**
```http
# 公共系统接口
GET    /health                   # 健康检查
GET    /api/public/version       # 获取系统版本
GET    /api/public/configs       # 获取公共配置

# 管理员系统管理
GET    /api/oa/admin/system/config # 获取系统配置
PUT    /api/oa/admin/system/config # 设置系统配置
GET    /api/oa/admin/system/configs # 获取配置列表
GET    /api/oa/admin/system/health  # 系统健康检查
GET    /api/oa/admin/system/statistics # 获取系统统计
```

---

## 🔗 API设计规范

### 📋 URL规范
```
{protocol}://{domain}/api/{module}/{resource}[/{id}][/{action}]

示例：
https://api.huinong.com/api/user/profile
https://api.huinong.com/api/oa/admin/loans/applications/123/approve
https://api.huinong.com/api/internal/dify/loan/get-application-details
```

### 🏷️ 路由分组规范
- `/api/public/*` - 公开API，无需认证
- `/api/auth/*` - 惠农用户认证相关
- `/api/user/*` - 惠农用户功能 (需要认证)
- `/api/content/*` - 公共内容 (可选认证)
- `/api/oa/auth/*` - OA系统认证相关  
- `/api/oa/user/*` - OA普通用户功能
- `/api/oa/admin/*` - OA管理员功能
- `/api/internal/dify/*` - Dify工作流内部调用

### 🏷️ HTTP方法规范
- `GET` - 查询数据
- `POST` - 创建资源/执行操作
- `PUT` - 更新整个资源
- `PATCH` - 部分更新资源
- `DELETE` - 删除资源

### 📦 响应格式规范
```json
{
    "code": 200,
    "message": "操作成功",
    "data": {
        // 响应数据
    },
    "meta": {
        "timestamp": "2024-01-15T10:30:00Z",
        "request_id": "req_abc123",
        "pagination": {
            "page": 1,
            "limit": 20,
            "total": 100
        }
    }
}
```

### 🔒 认证方式
```http
Authorization: Bearer {access_token}
```

### ⚠️ 错误码规范
```
1xxx - 用户相关错误
2xxx - 贷款相关错误  
3xxx - 农机相关错误
4xxx - OA管理相关错误
5xxx - 内容相关错误
6xxx - 系统相关错误
```

---

## 🔐 认证与授权

系统采用基于Redis的统一会话管理机制，为不同平台（惠农APP、惠农Web、OA后台）提供安全、高效的认证服务。

### 🔑 核心认证流程

1.  **登录获取Token**:
    *   惠农APP/Web用户通过 `/api/auth/login` 登录。
    *   OA系统用户通过 `/api/oa/auth/login` 登录。
    *   成功登录后，返回 `access_token` 和 `refresh_token`。
2.  **访问受保护API**:
    *   在请求头中携带 `Authorization: Bearer {access_token}`。
3.  **Token验证与平台检查**:
    *   `RequireAuth()`: 验证Token有效性，并将用户信息存入上下文。
    *   `CheckPlatform("{platform}")`: 检查当前用户的平台是否符合接口要求。
    *   `RequireRole("{role}")`: (仅OA管理员接口) 检查OA用户角色。
4.  **Token刷新**:
    *   使用 `refresh_token` 通过 `/api/auth/refresh` 或 `/api/oa/auth/refresh` 获取新Token。

### 🛡️ 权限模型

1.  **平台隔离**: Token包含平台信息 (`app`, `web`, `oa`)，通过 `CheckPlatform()` 中间件控制访问。
2.  **路由分组**: 不同平台API通过路由前缀物理隔离。
3.  **角色权限**: OA系统通过 `RequireRole("admin")` 控制管理员权限。

### 🎫 中间件链

```go
// 基础中间件
Recovery() -> RequestLogger() -> CORS()

// 认证中间件
RequireAuth() -> CheckPlatform("oa") -> RequireRole("admin")

// 可选认证
OptionalAuth() // 为登录用户提供个性化内容
```

### 🚦 API端点权限示例

| API端点 | 认证要求 | 平台要求 | 角色要求 | 说明 |
|---------|----------|----------|----------|------|
| `POST /api/auth/login` | 无 | N/A | N/A | 惠农用户登录 |
| `GET /api/user/profile` | `RequireAuth` | `app`/`web` | N/A | 惠农用户信息 |
| `GET /api/oa/user/profile` | `RequireAuth` | `oa` | N/A | OA普通用户信息 |
| `GET /api/oa/admin/users` | `RequireAuth` | `oa` | `admin` | OA管理员功能 |
| `GET /api/content/articles` | `OptionalAuth` | N/A | N/A | 公开内容 |
| `POST /api/internal/dify/*` | Dify Token | N/A | N/A | AI工作流调用 |

---

## 📊 数据模型关系

```
Users (惠农用户表)
├── UserAuths (认证信息)
├── UserSessions (会话记录) 
├── UserTags (用户标签)
└── UserTransactions (交易记录)

OAUsers (OA用户表)
├── OARoles (角色权限)
├── OASessions (OA会话)
└── OAOperationLogs (操作日志)

LoanApplications (贷款申请)
├── LoanProducts (贷款产品)
├── LoanReviews (审批记录)
└── AIAssessments (AI评估记录)

MachineRentals (农机租赁)
├── Machines (农机设备)
├── MachineOrders (租赁订单)
└── OrderRatings (评价反馈)

Content (内容管理)
├── Articles (文章资讯)
├── Experts (专家信息)
├── Consultations (咨询记录)
├── Announcements (系统公告)
└── Categories (内容分类)

Tasks (任务管理)
├── TaskAssignments (任务分配)
├── TaskProgress (任务进度)
└── TaskLogs (操作记录)

Files (文件管理)
├── FileTypes (文件类型)
├── FilePermissions (文件权限)
## 🚀 部署与环境

### 🏗️ 环境配置
```yaml
# 开发环境
development:
  database_url: mysql://localhost:3306/huinong_dev
  redis_url: redis://localhost:6379/0
  file_storage: local
  log_level: debug

# 生产环境  
production:
  database_url: mysql://prod-db:3306/huinong_prod
  redis_url: redis://prod-redis:6379/1
  file_storage: oss
  log_level: info
```

### 🔧 API网关配置
```yaml
gateway:
  rate_limit: 1000/minute
  cors_enabled: true
  auth_required: true
  request_timeout: 30s
  response_cache: 300s
```

### 📈 监控指标
- **QPS**: 每秒请求数
- **响应时间**: P95/P99延迟
- **错误率**: 4XX/5XX错误比例
- **可用性**: 服务正常运行时间
- **资源使用**: CPU/内存/网络使用率

---

## 📝 开发指南

### 🛠️ 本地开发
```bash
# 1. 克隆代码
git clone https://github.com/company/huinong-backend.git

# 2. 安装依赖
go mod download

# 3. 配置环境变量
cp .env.example .env

# 4. 启动数据库
docker-compose up -d mysql redis

# 5. 运行迁移
go run cmd/migrate/main.go

# 6. 启动服务
go run cmd/server/main.go
```

### 🧪 测试指南
```bash
# 单元测试
go test ./...

# 集成测试
go test -tags=integration ./...

# API测试
newman run tests/api/postman_collection.json
```

### 📚 文档生成
```bash
# 生成API文档
swag init

# 生成数据模型文档
go run tools/model-doc/main.go
```

---

## 🔄 版本历史

| 版本 | 日期 | 变更说明 |
|------|------|----------|
| v1.3.1 | 2024-01-15 | **OA登录系统优化** - 删除验证码功能，简化登录流程，优化路由配置 |
| v1.3.0 | 2024-04-01 | 新增专家咨询功能 |
| v1.2.0 | 2024-03-01 | 优化贷款审批流程 |
| v1.1.0 | 2024-02-01 | 新增农机租赁模块 |
| v1.0.0 | 2024-01-15 | 初始版本发布 |

---

## 🤝 贡献指南

### 📋 代码规范
- 遵循 Go 官方代码规范
- 使用 gofmt 格式化代码
- 添加必要的单元测试
- 编写清晰的注释文档

### 🔀 提交流程
1. Fork 项目到个人仓库
2. 创建功能分支 `git checkout -b feature/new-feature`
3. 提交代码 `git commit -m "Add new feature"`
4. 推送分支 `git push origin feature/new-feature`
5. 创建 Pull Request

### 🐛 问题反馈
- 使用 GitHub Issues 报告问题
- 提供详细的错误信息和复现步骤
- 标明环境信息和版本号

---

## 📞 联系方式

- **项目负责人**: 技术团队
- **邮箱**: tech@huinong.com
- **文档维护**: API文档团队
- **更新频率**: 每周更新

---

## 📄 许可证

本项目采用 MIT 许可证，详情请参阅 [LICENSE](../../../LICENSE) 文件。 