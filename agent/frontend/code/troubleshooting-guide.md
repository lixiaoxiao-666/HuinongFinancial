# OA系统登录问题故障排除指南

## 🚨 当前状态

### 问题现象
1. ✅ 前端服务正常运行 (localhost:5173)
2. ✅ 后端服务正常运行 (localhost:8080)  
3. ❌ 验证码接口返回404错误
4. ⚠️ Go编译环境暂时不可用

### 最新修复
1. **后端路由修复**：已将OA登录接口移出认证中间件
2. **前端API恢复**：已将接口调用改回 `/api/oa/auth/*`
3. **错误处理优化**：验证码失败时不阻塞登录流程

## 🔧 临时解决方案

### 验证码处理
如果验证码接口不可用，系统会：
- 显示"验证码暂不可用"
- 隐藏验证码输入框
- 允许直接进行登录测试

### 测试账户
开发环境提供了快速登录按钮：
- **超级管理员**: admin / admin123
- **审批员**: reviewer / reviewer123

## 📋 验证步骤

### 1. 检查服务状态
```bash
# 检查后端服务
curl http://localhost:8080/health

# 检查前端服务  
curl http://localhost:5173
```

### 2. 测试API接口
```bash
# 测试验证码接口
curl http://localhost:8080/api/oa/auth/captcha

# 测试登录接口
curl -X POST http://localhost:8080/api/oa/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 3. 前端功能测试
1. 访问 http://localhost:5173
2. 尝试点击"快速登录"按钮
3. 或手动输入测试账户信息
4. 检查浏览器控制台错误信息

## 🔄 下一步行动计划

### 优先级1: 修复编译环境
1. 重新安装Go环境
2. 配置PATH环境变量
3. 测试 `go version` 命令

### 优先级2: 重新编译后端
1. 运行 `go mod tidy`
2. 执行 `go build -o server.exe ./cmd/server`
3. 重启后端服务应用路由修复

### 优先级3: 验证完整流程
1. 测试验证码接口
2. 测试登录认证
3. 测试会话管理
4. 测试权限控制

## 🛠️ 调试命令

### 后端服务管理
```bash
# 查看进程
tasklist | findstr server

# 停止服务
taskkill /f /im server.exe

# 启动服务
.\server.exe

# 查看日志
tail -f logs/app.log
```

### 前端开发服务
```bash
# 启动开发服务器
pnpm dev

# 清理缓存重启
pnpm dev --force

# 构建生产版本
pnpm build
```

### 网络诊断
```bash
# 检查端口占用
netstat -ano | findstr :8080
netstat -ano | findstr :5173

# 测试连通性
ping localhost
telnet localhost 8080
```

## 📊 错误码对照表

| 错误码 | 含义 | 解决方案 |
|--------|------|----------|
| 404 | 接口不存在 | 检查路由配置和URL路径 |
| 401 | 未授权 | 检查token或登录状态 |
| 403 | 权限不足 | 检查用户权限配置 |
| 500 | 服务器错误 | 检查后端日志和配置 |

## 🔐 安全检查清单

- [ ] 公开接口白名单配置正确
- [ ] 认证中间件应用到正确的路由组
- [ ] Session配置和Redis连接正常
- [ ] CORS设置允许前端域名
- [ ] 设备信息和平台标识正确传递

## 📝 开发提示

### 前端调试
1. 使用浏览器开发者工具查看网络请求
2. 检查Vue Devtools中的Pinia状态
3. 查看控制台错误和警告信息

### 后端调试  
1. 检查应用日志文件
2. 使用Postman或curl测试API接口
3. 验证数据库连接和Redis状态

### 环境配置
1. 确保配置文件正确 (`configs/config.yaml`)
2. 检查环境变量设置
3. 验证数据库和Redis连接字符串 