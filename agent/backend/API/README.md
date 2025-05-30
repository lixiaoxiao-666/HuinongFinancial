# 数字惠农系统 - API 接口文档总览

## 📋 系统概述

数字惠农系统是一个综合性的智慧农业服务平台，为农户提供金融服务、农机租赁、政策资讯、技术指导等一站式服务。系统采用微服务架构，支持移动端APP和Web端OA后台管理。

### 🏗️ 系统架构

```
┌─────────────────────────────────────────────────────────┐
│                     前端应用层                            │
├──────────────────┬──────────────────┬──────────────────┤
│    移动端APP      │     Web端        │    OA管理后台     │
│   (iOS/Android)   │   (Vue3 + TS)    │   (Vue3 + TS)    │
└──────────────────┴──────────────────┴──────────────────┘
                              │
┌─────────────────────────────────────────────────────────┐
│                     API网关层                            │
│          统一路由、认证、限流、监控                        │
└─────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────┐
│                   业务服务层                              │
├──────────┬──────────┬──────────┬──────────┬──────────┤
│ 用户管理   │ 贷款服务   │ 农机租赁   │ 内容管理   │ OA管理   │
│ Service   │ Service   │ Service   │ Service   │ Service   │
└──────────┴──────────┴──────────┴──────────┴──────────┘
                              │
┌─────────────────────────────────────────────────────────┐
│                   数据存储层                              │
├──────────┬──────────┬──────────┬──────────┬──────────┤
│  MySQL    │  Redis    │ MongoDB   │   OSS     │  Kafka   │
│ 主数据库   │ 缓存/会话  │ 日志/统计  │ 文件存储   │ 消息队列 │
└──────────┴──────────┴──────────┴──────────┴──────────┘
```

### 🔧 技术栈

- **后端框架**: Golang + Gin + GORM
- **数据库**: MySQL 8.0 + Redis 6.0
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

**技术特性:**
- 🚀 **Redis集群**: 支持Redis集群部署，保证高可用
- 🔄 **实时同步**: 会话状态跨实例实时同步
- 🛡️ **安全机制**: Token哈希存储、设备指纹、防重放攻击
- 📈 **性能优化**: 连接池管理、批量操作、智能清理

**主要接口:**
```http
POST   /api/auth/login           # 用户登录 (智能设备检测)
POST   /api/auth/refresh         # 通用Token刷新
POST   /api/oa/auth/refresh      # OA专用Token刷新
GET    /api/auth/validate        # 会话验证
POST   /api/user/logout          # 用户登出
GET    /api/user/session/info    # 获取用户会话列表
DELETE /api/user/session/{id}    # 注销指定会话
POST   /api/user/session/revoke-others # 注销其他设备
GET    /api/admin/sessions/statistics # 会话统计
```

### 👤 2. 用户管理模块 [`user_management.md`](./user_management.md)
完整的用户生命周期管理，支持多种用户类型和身份认证。

**核心功能:**
- 📝 用户注册登录 (手机/密码/验证码)
- 🆔 身份认证 (实名/银行卡/征信)
- 📋 信息管理 (个人资料/地址/头像)
- 🏷️ 用户标签 (行为分析/精准推荐)
- 💰 账户管理 (余额/流水/统计)

**主要接口:**
```http
POST   /api/auth/register        # 用户注册
POST   /api/auth/send-sms        # 发送验证码
GET    /api/user/profile         # 获取用户信息
POST   /api/user/auth/real-name  # 实名认证申请
GET    /api/user/balance         # 获取账户余额
```

### 💰 3. 贷款管理模块 [`loan_management.md`](./loan_management.md)
智能化贷款服务，集成AI风险评估和自动化审批流程。

**核心功能:**
- 🏦 贷款产品管理 (产品目录/详情查询)
- 📋 申请流程 (在线申请/材料上传/状态跟踪)
- 🤖 AI智能审批 (风险评估/额度计算)
- 💳 合同管理 (电子签约/放款确认)
- 📊 还款管理 (计划查询/主动还款/提前还款)

**主要接口:**
```http
GET    /api/loans/products       # 获取贷款产品
POST   /api/loans/applications   # 提交贷款申请
GET    /api/loans/credit-limit   # 获取信用额度
POST   /api/loans/{id}/repayment # 主动还款
GET    /api/loans/overview       # 贷款概览
```

### 🚜 4. 农机租赁模块 [`machine_rental.md`](./machine_rental.md)
全流程农机租赁服务，支持设备搜索、预约、使用跟踪和归还管理。

**核心功能:**
- 🔍 设备搜索 (多维度筛选/地理位置匹配)
- 📅 预约管理 (时间规划/冲突检测/智能推荐)
- 💼 订单管理 (在线支付/合同签署/服务确认)
- 📍 设备跟踪 (GPS定位/状态监控/使用统计)
- 🔄 归还流程 (状态检查/费用结算/评价反馈)

**主要接口:**
```http
GET    /api/machines/search      # 搜索农机设备
POST   /api/machines/reservations # 创建预约
GET    /api/machines/{id}/location # 获取设备位置
POST   /api/machines/orders/{id}/return # 申请归还
GET    /api/machines/statistics  # 租赁统计
```

### 📰 5. 内容管理模块 [`content_management.md`](./content_management.md)
丰富的农业内容服务，提供资讯、政策、专家咨询等信息化支持。

**核心功能:**
- 📰 资讯管理 (农业新闻/技术资讯/市场行情)
- 📋 政策服务 (政策发布/解读/在线申请)
- 👨‍🎓 专家咨询 (在线问答/技术指导/知识库)
- 🔔 通知管理 (系统通知/消息推送/状态提醒)
- 🏷️ 内容分类 (智能标签/个性化推荐)

**主要接口:**
```http
GET    /api/content/articles     # 获取资讯列表
GET    /api/content/policies     # 获取政策列表
POST   /api/content/consultations # 提交咨询问题
GET    /api/content/notifications # 获取通知列表
POST   /api/content/articles/{id}/like # 点赞文章
```

### 🏢 6. OA后台管理模块 [`oa_management.md`](./oa_management.md)
全面的后台管理功能，支持业务审批、数据统计、系统配置等管理需求。

**核心功能:**
- 👥 用户管理 (用户列表/状态管理/权限控制)
- 📋 审批管理 (贷款审批/认证审核/AI辅助决策)
- 📊 数据统计 (业务报表/风险监控/用户分析)
- ⚙️ 系统配置 (参数设置/通知模板/权限管理)
- 🔍 日志审计 (操作记录/系统日志/安全监控)

**主要接口:**
```http
POST   /api/oa/auth/login        # 管理员登录 (已简化，无需验证码)
GET    /api/oa/users             # 获取用户列表
POST   /api/oa/loan-applications/{id}/review # 审批贷款申请
GET    /api/oa/dashboard/overview # 获取业务概览
GET    /api/oa/logs/operations   # 获取操作日志
```

**🆕 最新更新 (2024-01-15):**
- ✅ **简化登录流程**: 删除验证码功能，提升管理员登录体验
- ✅ **路由优化**: 调整认证接口路由配置，提升安全性
- ✅ **代码重构**: 修复函数重复声明，完善OA服务实现

---

## 🔗 API设计规范

### 📋 URL规范
```
{protocol}://{domain}/api/{version}/{module}/{resource}[/{id}][/{action}]

示例：
https://api.huinong.com/api/v1/loans/applications/LA20240115001/review
```

### 🏷️ HTTP方法规范
- `GET` - 查询数据
- `POST` - 创建资源
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
    *   `RequireAuth()`: 验证Token有效性，并将用户信息（如 `user_id`, `platform`）存入上下文。
    *   `CheckPlatform("{platform_name}")`: （可选）检查当前用户的 `platform` 是否符合接口要求。
    *   `RequireRole("{role_name}")`: （仅OA管理员接口）检查OA用户的角色是否符合要求。
4.  **Token刷新**:
    *   当 `access_token` 过期时，使用 `refresh_token` 通过 `/api/auth/refresh` (惠农端) 或 `/api/oa/auth/refresh` (OA端) 获取新的Token对。

### 🛡️ 权限模型

系统权限主要通过以下方式控制：

1.  **平台隔离 (`platform`)**: Token中会包含用户登录的平台信息 (`app`, `web`, `oa`)。后端通过 `CheckPlatform()` 中间件确保用户只能访问其对应平台的API。
    *   例如，OA用户不能访问惠农APP特有的接口。
2.  **路由分组**: 不同平台的API通过不同的路由前缀进行物理隔离。
    *   惠农APP/Web通用接口: `/api/user/*`
    *   OA系统普通用户接口: `/api/oa/user/*`
    *   OA系统管理员接口: `/api/oa/admin/*`
3.  **OA系统角色 (`role`)**: OA系统内部有明确的角色划分（如 `admin`, `staff` 等）。
    *   管理员接口 (`/api/oa/admin/*`) 会通过 `RequireRole("admin")` 中间件校验用户是否具有管理员角色。

### 🎫 Token与会话

-   **Access Token**: 短效期（例如：24小时），用于API接口调用时的身份验证。
-   **Refresh Token**: 长效期（例如：7天），用于在Access Token过期后，安全地获取新的Access Token，无需用户重新登录。
-   **Session**: 用户登录后，会在Redis中创建一条会话记录，包含用户ID、平台、设备信息、Token哈希等。Token的有效性最终通过查询Redis会话来确认。

### 👥 用户与角色

-   **惠农APP/Web用户 (`User` 模型)**:
    *   主要通过 `UserType` (如 `farmer`, `farm_owner`) 区分用户类型，业务逻辑上可能有所不同，但API权限上通常视为普通用户。
-   **OA系统用户 (`OAUser` 模型)**:
    *   具有 `RoleID` 字段，关联到 `OARole` 模型，实现细粒度的权限控制。
    *   **普通OA用户**: 登录OA系统，可以访问 `/api/oa/user/*` 下的接口，如查看个人信息、自己提交的申请等。
    *   **OA管理员**: 具有 `admin` 角色，可以访问 `/api/oa/admin/*` 下的所有管理接口，如用户管理、贷款审批、系统配置等。

### 🚦 API端点权限示例

| API端点                               | 认证要求                         | 平台要求     | 角色要求     | 说明                                     |
| :------------------------------------ | :------------------------------- | :----------- | :----------- | :--------------------------------------- |
| `POST /api/auth/login`                | 无                               | N/A          | N/A          | 惠农APP/Web用户登录                      |
| `POST /api/oa/auth/login`             | 无                               | N/A          | N/A          | OA系统用户登录                           |
| `GET /api/user/profile`               | `RequireAuth`                    | `app` 或 `web` | N/A          | 获取惠农APP/Web用户个人信息                |
| `GET /api/oa/user/profile`            | `RequireAuth`                    | `oa`         | N/A          | 获取OA系统普通用户个人信息                 |
| `GET /api/oa/admin/users`             | `RequireAuth`                    | `oa`         | `admin`      | OA管理员获取用户列表                     |
| `POST /api/oa/admin/loans/approve`    | `RequireAuth`                    | `oa`         | `admin`      | OA管理员审批贷款                         |
| `GET /api/content/articles`           | `OptionalAuth` (可选认证)        | N/A          | N/A          | 公开获取文章，登录用户可能看到个性化内容   |
| `POST /api/internal/dify/...`         | Dify专用Token                   | N/A          | N/A          | Dify工作流内部调用                       |

---

## 📊 数据模型关系

```
Users (用户表)
├── UserAuths (认证信息)
├── UserSessions (会话记录) 
├── UserTags (用户标签)
├── UserAddresses (地址信息)
└── UserTransactions (交易记录)

LoanApplications (贷款申请)
├── LoanProducts (贷款产品)
├── LoanContracts (贷款合同)
├── LoanRepayments (还款记录)
└── LoanReviews (审批记录)

MachineRentals (农机租赁)
├── Machines (农机设备)
├── MachineReservations (预约记录)
├── MachineOrders (租赁订单)
└── MachineReturns (归还记录)

Content (内容管理)
├── Articles (文章资讯)
├── Policies (政策信息)
├── Experts (专家信息)
├── Consultations (咨询记录)
└── Notifications (通知消息)
```

---

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