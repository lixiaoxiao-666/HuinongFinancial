# 数字惠农APP后端服务

基于Golang、Gin框架构建的微服务后端，支持用户管理、贷款申请、农机租赁等功能。

## 技术栈

- **编程语言**: Go 1.23
- **Web框架**: Gin
- **数据库**: TiDB (兼容MySQL)
- **缓存**: Redis
- **ORM**: GORM
- **认证**: JWT
- **日志**: Zap
- **配置管理**: Viper

## 项目结构

```
backend/
├── cmd/                    # 应用入口点
│   └── api/
│       └── main.go
├── internal/               # 内部应用代码
│   ├── api/               # API处理器和路由
│   │   ├── middleware.go  # 中间件
│   │   ├── router.go      # 路由配置
│   │   ├── user_handler.go # 用户API处理器
│   │   └── loan_handler.go # 贷款API处理器
│   ├── biz/               # 业务逻辑层
│   ├── data/              # 数据访问层
│   │   ├── database.go    # 数据库连接
│   │   └── models.go      # 数据模型
│   ├── conf/              # 配置管理
│   │   └── config.go
│   └── service/           # 服务层
│       ├── user_service.go # 用户服务
│       └── loan_service.go # 贷款服务
├── pkg/                   # 公共库代码
│   ├── response.go        # 统一响应格式
│   ├── jwt.go            # JWT工具
│   ├── password.go       # 密码加密
│   └── utils.go          # 通用工具
├── configs/               # 配置文件
│   └── config.yaml
├── go.mod
├── go.sum
├── main.go               # 主入口文件
└── README.md
```

## 环境要求

- Go 1.23+
- TiDB (或MySQL兼容数据库)
- Redis

## 快速开始

### 1. 克隆项目
```bash
cd backend
```

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 配置数据库和Redis

确保TiDB和Redis服务正在运行：
- TiDB: 10.10.20.10:4000
- Redis: 10.10.20.10:6379

### 4. 启动服务
```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动

### 5. 验证服务
```bash
curl http://localhost:8080/health
```

## API文档

### 用户相关接口

#### 发送验证码
```http
POST /api/v1/users/send-verification-code
Content-Type: application/json

{
  "phone": "13800138000"
}
```

#### 用户注册
```http
POST /api/v1/users/register
Content-Type: application/json

{
  "phone": "13800138000",
  "password": "password123",
  "verification_code": "123456"
}
```

#### 用户登录
```http
POST /api/v1/users/login
Content-Type: application/json

{
  "phone": "13800138000",
  "password": "password123"
}
```

#### 获取用户信息
```http
GET /api/v1/users/me
Authorization: Bearer <token>
```

### 贷款相关接口

#### 获取贷款产品列表
```http
GET /api/v1/loans/products
```

#### 获取贷款产品详情
```http
GET /api/v1/loans/products/{product_id}
```

#### 提交贷款申请
```http
POST /api/v1/loans/applications
Authorization: Bearer <token>
Content-Type: application/json

{
  "product_id": "lp_spring2024",
  "amount": 30000,
  "term_months": 12,
  "purpose": "购买化肥和种子",
  "applicant_info": {
    "real_name": "张三",
    "id_card_number": "310...",
    "address": "XX省XX市XX村"
  },
  "uploaded_documents": [
    {"doc_type": "id_card_front", "file_id": "file_uuid_001"}
  ]
}
```

#### 获取我的贷款申请列表
```http
GET /api/v1/loans/applications/my?status=SUBMITTED&page=1&limit=10
Authorization: Bearer <token>
```

### 文件上传接口

#### 上传文件
```http
POST /api/v1/files/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <binary>
purpose: "loan_document"
```

### OA后台接口

#### OA用户登录
```http
POST /api/v1/admin/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

## 配置说明

主要配置文件位于 `configs/config.yaml`：

```yaml
app:
  name: "digital-agriculture-backend"
  version: "1.0.0"
  env: "development"

server:
  port: 8080
  mode: "debug"

database:
  host: "10.10.20.10"
  port: 4000
  username: "root"
  password: ""
  database: "digital_agriculture"

redis:
  host: "10.10.20.10"
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "digital-agriculture-jwt-secret-key-2024"
  expire: "24h"
```

## 数据库迁移

项目启动时会自动执行数据库迁移，创建所需的表结构。

## 示例数据

项目首次启动时会自动创建示例数据：
- 2个贷款产品（春耕助力贷、农机购置贷）
- 1个OA管理员用户（用户名：admin，密码：admin123）

## 测试

运行单元测试：
```bash
go test ./...
```

## 构建

构建可执行文件：
```bash
go build -o digital-agriculture-backend main.go
```

## Docker化部署

创建Dockerfile：
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o main main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs
CMD ["./main"]
```

构建和运行：
```bash
docker build -t digital-agriculture-backend .
docker run -p 8080:8080 digital-agriculture-backend
```

## 日志

应用日志输出到控制台，生产环境可配置文件输出。

## 监控

- 健康检查端点：`GET /health`
- 请求日志中间件记录所有API调用
- 支持结构化日志输出

## 许可证

[MIT License](LICENSE) 