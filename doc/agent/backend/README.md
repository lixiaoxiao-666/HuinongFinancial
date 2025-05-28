# 数字惠农OA管理系统 - 后端服务

## 项目概述

数字惠农OA管理系统是一个专为农业金融贷款审批而设计的智能化后台管理平台。系统集成了AI风险评估和人工审核的混合审批机制，为农户贷款提供高效、安全的审批流程管理。

## 核心功能特性

### 🔐 认证与权限管理
- **多角色权限控制**：支持系统管理员、审批员等不同角色
- **JWT Token认证**：安全的会话管理和API访问控制
- **用户状态管理**：支持账号启用/禁用、权限变更等操作

### 📊 智能工作台
- **实时数据监控**：展示贷款申请总数、通过率、处理时长等关键指标
- **个性化待办**：根据用户角色显示相关的待审批任务
- **快捷操作入口**：提供常用功能的快速访问通道
- **活动记录追踪**：实时显示系统操作动态

### 📋 智能审批管理
- **AI风险评估**：集成机器学习模型进行申请人信用风险评估
- **混合审批流程**：AI预筛选 + 人工复核的双重保障机制
- **申请全生命周期管理**：从提交到最终决策的完整流程追踪
- **多维度筛选查询**：支持按状态、申请人、金额等多个维度筛选
- **详细申请档案**：包含申请人信息、AI分析结果、历史记录等

### ⚙️ 系统管理与配置
- **AI审批开关控制**：灵活启用/禁用AI自动审批功能
- **配置参数管理**：支持动态调整系统配置参数
- **系统统计分析**：提供丰富的业务数据统计和分析
- **实时性能监控**：系统运行状态和性能指标监控

### 👥 用户管理
- **OA用户创建**：支持创建不同角色的系统用户
- **权限分配管理**：精细化的功能权限控制
- **用户状态控制**：账号启用、禁用、角色变更等操作
- **用户信息维护**：用户基本信息的查看和管理

### 📝 操作审计
- **全面操作记录**：记录所有系统操作的详细日志
- **多维度日志查询**：支持按操作者、操作类型、时间范围等筛选
- **操作追溯**：完整的操作链路追踪和审计
- **合规性保障**：满足金融行业的合规审计要求

## 技术架构

### 后端技术栈
- **编程语言**：Go 1.19+
- **Web框架**：Gin
- **数据库**：MySQL 8.0+
- **ORM框架**：GORM
- **认证授权**：JWT
- **密码加密**：bcrypt
- **日志组件**：Zap
- **配置管理**：Viper

### 项目结构
```
backend/
├── cmd/                    # 应用入口
│   └── api/               # API服务入口
├── internal/              # 内部代码
│   ├── api/              # API处理器
│   │   ├── admin_handler.go   # OA管理功能
│   │   ├── user_handler.go    # 用户功能
│   │   ├── loan_handler.go    # 贷款功能
│   │   ├── middleware.go      # 中间件
│   │   └── router.go          # 路由配置
│   ├── service/          # 业务逻辑层
│   │   ├── admin_service.go   # OA管理服务
│   │   ├── user_service.go    # 用户服务
│   │   └── loan_service.go    # 贷款服务
│   ├── data/             # 数据层
│   │   ├── models.go          # 数据模型
│   │   └── database.go        # 数据库配置
│   └── conf/             # 配置管理
├── pkg/                   # 公共组件
│   ├── jwt.go            # JWT工具
│   ├── password.go       # 密码工具
│   ├── utils.go          # 通用工具
│   └── response.go       # 响应工具
└── configs/              # 配置文件
```

### 数据库设计
核心数据表包括：
- `oa_users`: OA系统用户表
- `loan_applications`: 贷款申请表
- `loan_application_history`: 申请审批历史表
- `system_configurations`: 系统配置表
- `uploaded_files`: 文件上传记录表

## 快速开始

### 环境要求
- Go 1.19+
- MySQL 8.0+
- Redis 6.0+ (可选)

### 安装部署

1. **克隆项目**
```bash
git clone [项目地址]
cd HuinongFinancial/backend
```

2. **安装依赖**
```bash
go mod tidy
```

3. **配置数据库**
```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE huinong_finance;
```

4. **配置文件**
修改 `configs/config.yaml` 中的数据库连接信息：
```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: your_password
  database: huinong_finance
```

5. **启动服务**
```bash
go run cmd/api/main.go
```

6. **验证服务**
访问 `http://localhost:8080/health` 检查服务状态

## API接口文档

详细的API接口文档请参考：
- [API_Spec.md](./API_Spec.md) - 完整的接口文档
- Swagger UI: `http://localhost:8080/swagger/index.html` (如已配置)

### 核心接口概览

#### 认证接口
- `POST /admin/login` - OA用户登录

#### 工作台接口
- `GET /admin/dashboard` - 获取工作台信息

#### 审批管理
- `GET /admin/loans/applications/pending` - 待审批列表
- `GET /admin/loans/applications/{id}` - 申请详情
- `POST /admin/loans/applications/{id}/review` - 提交审批

#### 系统管理
- `GET /admin/system/stats` - 系统统计
- `POST /admin/system/ai-approval/toggle` - AI审批开关

#### 用户管理
- `GET /admin/users` - 用户列表
- `POST /admin/users` - 创建用户
- `PUT /admin/users/{id}/status` - 更新用户状态

#### 日志和配置
- `GET /admin/logs` - 操作日志
- `GET /admin/configs` - 系统配置
- `PUT /admin/configs/{key}` - 更新配置

## 测试

### 接口测试
我们提供了完整的API接口测试脚本：

**Linux/Mac:**
```bash
chmod +x doc/agent/backend/Test-API.sh
./doc/agent/backend/Test-API.sh
```

**Windows:**
```powershell
.\doc\agent\backend\Test-API.ps1
```

### 单元测试
```bash
go test ./...
```

## 默认账号

系统预设了以下测试账号：

| 角色 | 用户名 | 密码 | 权限 |
|------|--------|------|------|
| 系统管理员 | admin | admin123 | 全部功能 |
| 审批员 | reviewer | reviewer123 | 审批相关功能 |

⚠️ **生产环境部署时请务必修改默认密码**

## 安全配置

### JWT配置
- Token有效期：默认1小时
- 签名算法：HMAC-SHA256
- 刷新机制：支持token刷新

### 密码安全
- 使用bcrypt加密存储
- 最小密码长度：6位
- 支持密码强度验证

### API安全
- 所有管理接口均需要JWT认证
- 基于角色的权限控制
- 请求参数验证和SQL注入防护

## 性能优化

### 数据库优化
- 合理的索引设计
- 分页查询优化
- 连接池配置

### 应用优化
- 高效的JSON序列化
- 内存管理优化
- 并发处理能力

## 监控与运维

### 日志管理
- 结构化日志输出
- 不同级别的日志分类
- 日志轮转和清理策略

### 健康检查
- 应用健康状态检查：`/health`
- 数据库连接检查
- 依赖服务状态检查

### 指标监控
- 系统运行指标
- 业务指标统计
- 性能监控数据

## 开发规范

### 代码规范
- 遵循Go官方编码规范
- 使用gofmt格式化代码
- 完善的错误处理机制

### API设计规范
- RESTful API设计
- 统一的响应格式
- 清晰的错误码定义

### 提交规范
- 清晰的commit message
- 功能分支开发
- 代码review流程

## 常见问题

### Q: 服务启动失败？
A: 检查数据库连接配置，确保MySQL服务正常运行

### Q: JWT Token失效？
A: 检查系统时间同步，token默认1小时有效期

### Q: 权限验证失败？
A: 确认用户角色配置正确，检查middleware配置

### Q: 数据库连接超时？
A: 检查数据库连接池配置，调整超时参数

## 版本历史

- v1.0.0: 基础OA管理系统功能
- v1.1.0: 增加AI审批功能
- v1.2.0: 完善权限管理和操作审计

## 联系方式

如有问题或建议，请联系开发团队或提交Issue。

## 许可证

本项目采用 [MIT License](LICENSE) 许可证。 