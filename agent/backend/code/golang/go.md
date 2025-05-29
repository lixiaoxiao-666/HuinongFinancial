# Go模块配置文件

## 文件概述

`go.mod` 是数字惠农后端项目的Go模块配置文件，定义了项目的模块名称和所有依赖包。

## 模块信息

- **模块名称**: `huinong-backend`
- **Go版本**: 1.21
- **项目类型**: 数字惠农后端服务

## 核心依赖

### Web框架
- `github.com/gin-gonic/gin v1.9.1` - 高性能的HTTP Web框架
- `github.com/gin-contrib/cors v1.4.0` - CORS跨域支持中间件

### 数据库相关
- `gorm.io/gorm v1.25.5` - ORM框架
- `gorm.io/driver/mysql v1.5.2` - MySQL数据库驱动
- `github.com/redis/go-redis/v9 v9.3.0` - Redis客户端

### 认证授权
- `github.com/golang-jwt/jwt/v5 v5.2.0` - JWT Token处理
- `golang.org/x/crypto v0.17.0` - 加密算法支持

### 配置管理
- `github.com/spf13/viper v1.17.0` - 配置文件管理
- `gopkg.in/yaml.v3 v3.0.1` - YAML文件解析

### API文档
- `github.com/swaggo/gin-swagger v1.6.0` - Swagger集成
- `github.com/swaggo/swag v1.16.2` - API文档生成

### 工具包
- `github.com/google/uuid v1.5.0` - UUID生成
- `github.com/go-playground/validator/v10 v10.16.0` - 数据验证

## 项目特色

1. **技术栈现代化**: 使用最新稳定版本的Go依赖包
2. **高性能**: 选择Gin作为Web框架，Redis作为缓存
3. **开发友好**: 集成Swagger自动生成API文档
4. **安全性**: 内置JWT认证和数据验证
5. **可维护性**: 使用GORM简化数据库操作

## 依赖更新

定期执行以下命令更新依赖：

```bash
go mod tidy
go mod download
```

## 注意事项

- 项目需要Go 1.21或更高版本
- 确保所有依赖包版本兼容
- 生产环境部署前需要进行依赖安全扫描 