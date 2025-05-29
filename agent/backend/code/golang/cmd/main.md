# 主程序入口文件

## 文件概述

`main.go` 是数字惠农后端服务的主程序入口文件，负责应用程序的初始化、启动和优雅关闭。

## 主要功能

### 1. 应用程序初始化
- **配置加载**: 根据环境变量加载对应的配置文件
- **依赖注入**: 初始化数据库连接、Redis连接、日志系统等
- **数据库迁移**: 自动执行数据库表结构迁移
- **路由初始化**: 设置API路由和中间件

### 2. 服务器启动
```go
srv := &http.Server{
    Addr:           fmt.Sprintf(":%d", cfg.App.Port),
    Handler:        r,
    ReadTimeout:    30 * time.Second,
    WriteTimeout:   30 * time.Second,
    MaxHeaderBytes: 1 << 20, // 1MB
}
```

### 3. 优雅关闭
- 监听系统信号(SIGINT, SIGTERM)
- 5秒超时时间进行资源清理
- 确保正在处理的请求完成后再关闭

## 配置文件路径规则

### 环境变量配置
- `CONFIG_PATH`: 直接指定配置文件路径
- `APP_ENV`: 指定运行环境(development/production/test)

### 默认路径规则
| 环境 | 配置文件路径 |
|------|-------------|
| development | configs/config.yaml |
| production | configs/config.prod.yaml |
| test | configs/config.test.yaml |

## Swagger文档配置

主程序包含完整的Swagger API文档注释：

```go
// @title 数字惠农API文档
// @version 1.0
// @description 数字惠农系统后端API接口文档
// @termsOfService https://www.huinong.com/terms
// @contact.name API Support
// @contact.url https://www.huinong.com/support
// @contact.email support@huinong.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
```

## 依赖注入容器

### Container结构
```go
type Container struct {
    config *config.Config
    db     interface{} // GORM数据库连接
    redis  interface{} // Redis客户端
    logger Logger       // 日志接口
}
```

### 初始化流程
1. 加载配置文件
2. 初始化日志系统
3. 连接数据库并执行迁移
4. 连接Redis缓存
5. 创建依赖注入容器
6. 初始化路由和中间件
7. 启动HTTP服务器

## 日志接口设计

### Logger接口定义
```go
type Logger interface {
    Info(msg string)
    Infof(format string, args ...interface{})
    Error(msg string)
    Errorf(format string, args ...interface{})
    Fatal(msg string)
    Fatalf(format string, args ...interface{})
}
```

### 日志级别
- **Info**: 一般信息记录
- **Error**: 错误信息记录
- **Fatal**: 致命错误，程序退出

## 服务器配置

### HTTP服务器参数
- **ReadTimeout**: 30秒读取超时
- **WriteTimeout**: 30秒写入超时
- **MaxHeaderBytes**: 1MB最大请求头大小

### 安全配置
- 支持HTTPS协议
- 请求头大小限制
- 优雅关闭机制

## 运行方式

### 开发环境启动
```bash
go run cmd/server/main.go
```

### 生产环境启动
```bash
APP_ENV=production go run cmd/server/main.go
```

### Docker容器启动
```bash
docker run -p 8080:8080 -e APP_ENV=production huinong-backend
```

## 环境变量支持

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| CONFIG_PATH | 配置文件路径 | 根据APP_ENV自动选择 |
| APP_ENV | 运行环境 | development |
| HUINONG_APP_PORT | 服务端口 | 8080 |
| HUINONG_DATABASE_HOST | 数据库主机 | localhost |
| HUINONG_REDIS_HOST | Redis主机 | localhost |

## 错误处理

### 启动阶段错误
- 配置文件加载失败: 程序退出
- 数据库连接失败: 程序退出
- Redis连接失败: 程序退出
- 数据库迁移失败: 程序退出

### 运行时错误
- HTTP服务器错误: 记录日志并继续运行
- 优雅关闭超时: 强制关闭并记录日志

## 性能优化

1. **连接池管理**: 数据库和Redis连接池优化
2. **内存管理**: 及时释放资源，避免内存泄露
3. **并发处理**: 支持高并发请求处理
4. **超时控制**: 合理设置各种超时时间

## 监控和健康检查

- 服务器启动日志记录
- 环境信息日志输出
- 优雅关闭状态监控
- 资源释放确认

## 扩展性设计

- **模块化架构**: 支持功能模块独立扩展
- **配置化管理**: 通过配置文件控制功能开关
- **插件化设计**: 支持中间件和插件扩展
- **微服务准备**: 架构设计支持后续微服务拆分 