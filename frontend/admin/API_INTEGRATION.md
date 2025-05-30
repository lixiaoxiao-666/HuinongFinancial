# 惠农OA管理系统 - 后端API集成说明

## 📋 概述

本文档说明如何将前端系统从Mock数据模式切换到真实的后端API集成模式。

## 🔧 主要修改

### 1. API服务层重构

#### 请求工具 (`src/api/request.ts`)
- ✅ 重构axios配置，支持真实API调用
- ✅ 添加完整的错误处理和响应拦截
- ✅ 支持Token自动刷新机制
- ✅ 添加请求/响应日志记录

#### 认证API (`src/api/auth.ts`)
- ✅ 对接OA系统认证接口 `/api/oa/auth/*`
- ✅ 支持设备信息收集和会话管理
- ✅ 完整的TypeScript类型定义

#### 仪表盘API (`src/api/modules/dashboard.ts`)
- ✅ 集成后端统计数据接口
- ✅ 支持概览数据、风险监控、会话统计
- ✅ 待处理任务和最新申请数据

### 2. 状态管理重构

#### Auth Store (`src/stores/modules/auth.ts`)
- ✅ 重构认证状态管理
- ✅ 支持Token自动刷新和持久化
- ✅ 完整的用户信息管理
- ✅ 角色权限检查机制

### 3. 页面组件更新

#### Dashboard页面 (`src/views/dashboard/index.vue`)
- ✅ 使用真实API数据
- ✅ 添加加载状态和错误处理
- ✅ 响应式数据绑定

#### Login页面 (`src/views/auth/Login.vue`)
- ✅ 简化登录逻辑
- ✅ 移除Mock相关代码
- ✅ 更好的错误提示

### 4. 路由守卫增强

#### 认证守卫 (`src/router/guards.ts`)
- ✅ 完整的认证状态检查
- ✅ Token有效性验证
- ✅ 角色权限控制
- ✅ 自动重定向逻辑

## 🚀 配置指南

### 1. 环境变量配置

创建 `.env.development` 文件：

```bash
# 开发环境配置
NODE_ENV=development

# API配置 - 后端服务地址
VITE_API_BASE_URL=http://localhost:8080
VITE_API_TIMEOUT=30000

# 应用配置
VITE_APP_TITLE=惠农OA管理系统
VITE_APP_VERSION=1.0.0

# 功能开关
VITE_USE_MOCK=false  # 关闭Mock模式
VITE_ENABLE_DEVTOOLS=true

# 调试配置
VITE_DEBUG_API=true
VITE_SHOW_REQUEST_LOGS=true
```

### 2. 后端服务要求

确保后端服务提供以下API接口：

#### 认证接口
- `POST /api/oa/auth/login` - OA用户登录
- `POST /api/oa/auth/refresh` - Token刷新
- `GET /api/oa/auth/validate` - Token验证
- `POST /api/oa/auth/logout` - 用户登出

#### 用户管理接口
- `GET /api/oa/user/profile` - 获取当前用户信息
- `PUT /api/oa/user/profile` - 更新用户信息
- `PUT /api/oa/user/password` - 修改密码

#### 仪表盘接口
- `GET /api/oa/admin/dashboard/overview` - 概览数据
- `GET /api/oa/admin/dashboard/risk-monitoring` - 风险监控
- `GET /api/oa/admin/sessions/statistics` - 会话统计

### 3. 启动配置

1. **安装依赖**：
   ```bash
   cd frontend/admin
   pnpm install
   ```

2. **配置环境变量**：
   ```bash
   cp .env.example .env.development
   # 编辑 .env.development，设置正确的API地址
   ```

3. **启动开发服务器**：
   ```bash
   pnpm dev
   ```

4. **确保后端服务运行**：
   - 后端服务应运行在 `http://localhost:8080`
   - 或根据实际情况修改 `VITE_API_BASE_URL`

## 📝 API响应格式

所有API接口应返回统一的响应格式：

```typescript
interface ApiResponse<T = any> {
  code: number        // 状态码，200表示成功
  message: string     // 响应消息
  data: T            // 响应数据
  meta?: {           // 元数据（可选）
    timestamp: string
    request_id: string
    pagination?: {
      page: number
      limit: number
      total: number
    }
  }
}
```

## 🔐 认证流程

### 登录流程
1. 用户提交用户名/密码
2. 前端调用 `/api/oa/auth/login`
3. 后端验证并返回用户信息和Token
4. 前端保存Token和用户信息
5. 自动跳转到仪表盘

### Token管理
- `access_token`: 短期有效（24小时）
- `refresh_token`: 长期有效（7天）
- 自动刷新机制，无需用户重新登录

### 权限控制
- 基于用户角色的权限控制
- 路由级别的权限检查
- API级别的权限验证

## 🛠️ 调试工具

### 开发环境
- 开启 `VITE_DEBUG_API=true` 查看API请求日志
- 使用浏览器开发者工具检查网络请求
- 检查控制台输出的认证状态信息

### 错误处理
- 统一的错误提示机制
- 自动Token刷新重试
- 网络异常自动重连

## 📊 监控和日志

### 前端日志
- API请求/响应日志
- 认证状态变化日志
- 路由跳转日志
- 错误和警告日志

### 性能监控
- API响应时间监控
- 页面加载时间统计
- 用户操作行为追踪

## 🚨 常见问题

### 1. CORS跨域问题
```javascript
// vite.config.ts
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  }
})
```

### 2. Token过期处理
- 自动检测Token过期
- 自动使用refresh_token刷新
- 刷新失败时跳转登录页

### 3. 网络异常处理
- 请求超时重试机制
- 网络断开提示
- 服务器错误友好提示

## 🔄 Mock模式切换

如需切换回Mock模式进行开发：

1. 设置环境变量：`VITE_USE_MOCK=true`
2. 确保Mock API服务正常运行
3. 重启开发服务器

## 📈 后续优化

- [ ] 添加API缓存机制
- [ ] 实现离线模式支持
- [ ] 添加数据预加载
- [ ] 优化Bundle大小
- [ ] 添加Service Worker

---

## 💡 技术栈

- **前端框架**: Vue 3 + TypeScript
- **状态管理**: Pinia
- **UI组件**: Ant Design Vue
- **HTTP客户端**: Axios
- **路由管理**: Vue Router
- **构建工具**: Vite
- **样式处理**: SCSS

## 🤝 开发规范

- 遵循Vue 3 Composition API最佳实践
- 使用TypeScript严格模式
- 统一的代码格式和命名规范
- 完整的错误处理和用户反馈
- 响应式设计和无障碍支持 