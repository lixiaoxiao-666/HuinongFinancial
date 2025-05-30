# request.ts - API请求配置文档

## 📋 模块概述

`request.ts` 是OA系统的核心HTTP请求配置文件，基于axios封装，提供统一的API请求处理、认证管理、错误处理和响应拦截功能。

### 🎯 主要功能

- **统一请求配置**: 设置baseURL、超时时间、请求头
- **认证头自动添加**: 自动携带Bearer Token
- **平台信息标识**: 添加OA平台标识头
- **统一错误处理**: HTTP状态码和业务错误的统一处理
- **自动登出**: Token失效时自动清除状态并跳转登录

---

## 🔧 技术实现

### 基础配置
```typescript
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/oa',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})
```

### 请求拦截器
- 自动添加 `Authorization: Bearer {token}` 认证头
- 添加平台标识 `X-Platform: oa`
- 添加设备类型 `X-Device-Type: web`
- 添加应用版本 `X-App-Version`

### 响应拦截器
- 统一处理业务成功响应（code: 200）
- 统一处理业务错误（code ≠ 200）
- HTTP状态码错误处理：
  - 401: 自动登出并跳转登录页
  - 403: 权限不足提示
  - 404: 资源不存在提示
  - 500: 服务器错误提示

---

## 📚 使用示例

### 基础使用
```typescript
import { request } from '@/api/request'

// GET请求
const response = await request.get('/user/profile')

// POST请求
const response = await request.post('/auth/login', {
  username: 'admin',
  password: 'password'
})
```

### 错误处理
```typescript
try {
  const response = await request.get('/api/data')
  console.log(response.data)
} catch (error) {
  // 错误已在拦截器中统一处理
  console.error('请求失败:', error)
}
```

---

## 🔗 依赖关系

### 外部依赖
- `axios` - HTTP请求库
- `ant-design-vue` - message提示组件

### 内部依赖
- `@/stores/modules/auth` - 认证状态管理
- `@/router` - 路由实例（用于登出跳转）

---

## ⚙️ 环境变量

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| `VITE_API_BASE_URL` | API基础URL | `/api/oa` |
| `VITE_APP_VERSION` | 应用版本号 | `1.0.0` |

---

## 🚨 错误码处理

| HTTP状态码 | 处理方式 | 用户提示 |
|------------|----------|----------|
| 401 | 自动登出 + 跳转登录 | "认证失败，请重新登录" |
| 403 | 显示错误信息 | "权限不足" |
| 404 | 显示错误信息 | "请求的资源不存在" |
| 500 | 显示错误信息 | "服务器内部错误" |
| 网络错误 | 显示错误信息 | "网络连接失败，请检查网络设置" |

---

## 🔄 扩展说明

### 添加新的请求拦截
```typescript
request.interceptors.request.use(
  (config) => {
    // 添加自定义请求头
    config.headers['X-Custom-Header'] = 'value'
    return config
  }
)
```

### 添加新的响应处理
```typescript
request.interceptors.response.use(
  (response) => {
    // 自定义响应处理逻辑
    return response
  },
  (error) => {
    // 自定义错误处理逻辑
    return Promise.reject(error)
  }
)
```

---

## 📝 最佳实践

1. **统一错误处理**: 不要在业务代码中重复处理通用错误
2. **Token管理**: 依赖store中的认证状态，不要手动管理Token
3. **请求超时**: 默认10秒超时，长耗时操作需要单独配置
4. **并发控制**: 对于频繁请求的接口考虑添加防抖处理

---

## 🔍 调试指南

### 开发环境调试
```typescript
// 查看请求配置
console.log('Request config:', config)

// 查看响应数据
console.log('Response data:', response.data)
```

### 网络问题排查
1. 检查 `VITE_API_BASE_URL` 配置
2. 确认后端服务运行状态
3. 检查浏览器网络面板
4. 验证Token是否有效

本模块为整个OA系统提供了稳定可靠的HTTP请求基础设施。 