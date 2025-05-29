# 数字惠农后端服务

这是数字惠农项目的后端API服务，采用Go语言开发，基于Gin框架构建。

## 🚀 快速开始

### 环境要求

- Go 1.19+
- MySQL 5.7+ 或 8.0+
- Redis 5.0+ (可选)
- Git

### 安装依赖

1. 确保已安装Go语言环境
2. 克隆项目到本地
3. 进入后端目录：`cd backend`

### 配置文件

1. 复制配置文件模板：
   ```bash
   # Windows
   copy configs\config.example.yaml configs\config.yaml
   
   # Linux/macOS
   cp configs/config.example.yaml configs/config.yaml
   ```

2. 修改 `configs/config.yaml` 中的配置：
   - **数据库配置**: 修改 `database` 部分的主机、端口、用户名、密码、数据库名
   - **Dify配置**: 修改 `dify` 部分的API地址和密钥
   - **JWT密钥**: 修改 `jwt.secret_key` 为安全的密钥
   - **其他服务**: 根据需要配置Redis、短信服务等

### 🔨 编译和运行

#### Windows用户

1. **一键运行**（推荐）：
   ```cmd
   run.bat
   ```
   此脚本会自动检查是否需要编译，然后启动服务。

2. **手动编译**：
   ```cmd
   build.bat
   ```

3. **运行编译后的程序**：
   ```cmd
   bin\huinong-backend.exe
   ```

#### Linux/macOS用户

1. **一键运行**（推荐）：
   ```bash
   chmod +x run.sh build.sh
   ./run.sh
   ```

2. **手动编译**：
   ```bash
   chmod +x build.sh
   ./build.sh
   ```

3. **运行编译后的程序**：
   ```bash
   ./bin/huinong-backend
   ```

### 🧪 连接测试

在首次运行前，可以使用测试脚本验证数据库和Dify平台连接：

```bash
# Windows
test.bat

# Linux/macOS (TODO: 创建test.sh)
```

## 📁 项目结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # 应用程序入口
├── configs/
│   ├── config.yaml              # 配置文件
│   └── config.example.yaml      # 配置文件模板
├── internal/
│   ├── config/                  # 配置管理
│   ├── database/                # 数据库连接和迁移
│   ├── model/                   # 数据模型
│   ├── repository/              # 数据访问层
│   ├── service/                 # 业务逻辑层
│   ├── handler/                 # HTTP处理器
│   ├── middleware/              # 中间件
│   ├── router/                  # 路由配置
│   └── utils/                   # 工具函数
├── bin/                         # 编译输出目录
├── logs/                        # 日志文件目录
├── uploads/                     # 文件上传目录
├── build.bat / build.sh         # 编译脚本
├── run.bat / run.sh             # 运行脚本
├── test.bat                     # 测试脚本
├── go.mod                       # Go模块文件
└── README.md                    # 项目说明
```

## 🔧 配置说明

### 数据库配置

```yaml
database:
  driver: "mysql"
  host: "localhost"           # 数据库主机
  port: 3306                  # 数据库端口
  username: "root"            # 用户名
  password: "your_password"   # 密码
  database: "huinong_db"      # 数据库名
  charset: "utf8mb4"
```

### Dify AI平台配置

```yaml
dify:
  api_url: "https://api.dify.ai/v1"           # Dify API地址
  api_key: "app-your-dify-api-key"            # API密钥
  timeout: 30                                 # 超时时间(秒)
  retry_times: 3                              # 重试次数
  workflows:
    loan_approval: "workflow-id"              # 贷款审批工作流ID
```

## 🔌 API接口

服务启动后，可以访问以下地址：

- **健康检查**: `GET http://localhost:8080/health`
- **API文档**: `http://localhost:8080/swagger/index.html`
- **用户接口**: `http://localhost:8080/api/v1/users/*`
- **贷款接口**: `http://localhost:8080/api/v1/loans/*`
- **农机接口**: `http://localhost:8080/api/v1/machines/*`

## 🐛 常见问题

### 1. 编译失败

**问题**: `go: command not found` 或类似错误
**解决**: 确保已正确安装Go语言环境并配置PATH

**问题**: 依赖下载失败
**解决**: 检查网络连接，可以配置Go代理：
```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

### 2. 数据库连接失败

**问题**: `connection refused` 或 `access denied`
**解决**: 
- 检查数据库服务是否启动
- 验证配置文件中的连接信息是否正确
- 确保数据库用户有足够的权限
- 检查防火墙设置

### 3. Dify平台连接失败

**问题**: API调用失败
**解决**:
- 检查Dify API地址是否正确
- 验证API密钥是否有效
- 检查网络连接和防火墙设置
- 确认工作流ID是否正确

### 4. 端口占用

**问题**: `bind: address already in use`
**解决**: 
- 修改配置文件中的端口号
- 或者停止占用端口的其他程序

## 📈 性能优化

1. **数据库连接池**：根据实际负载调整连接池参数
2. **Redis缓存**：启用Redis缓存提高查询性能
3. **日志级别**：生产环境设置为`info`或`warn`级别
4. **文件上传**：配置文件上传大小限制

## 🔒 安全配置

1. **JWT密钥**：使用强密钥并定期更换
2. **数据库密码**：使用复杂密码
3. **HTTPS**：生产环境启用HTTPS
4. **限流**：启用API限流保护

## 📚 更多文档

- [API接口文档](../agent/backend/API.md)
- [数据模型文档](../agent/backend/models/)
- [系统架构文档](../agent/arch/arch.md)

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 📄 许可证

此项目采用 MIT 许可证。 