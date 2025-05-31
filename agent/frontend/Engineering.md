# 数字惠农系统 - 前端工程化文档

## 📋 项目概述

数字惠农系统前端采用现代化前端技术栈，基于Vue3 + TypeScript构建，包含三个主要应用：惠农APP（移动端）、惠农Web（PC端）和OA后台管理系统。本文档定义了统一的工程化标准、开发规范和最佳实践。

### 🎯 技术选型

#### 核心技术栈
- **框架**: Vue 3.3+ (Composition API)
- **语言**: TypeScript 5.0+
- **构建工具**: Vite 4.0+
- **状态管理**: Pinia 2.0+
- **路由**: Vue Router 4.0+
- **UI组件库**: 
  - 移动端：Vant 4.0+
  - PC端：Ant Design Vue 4.0+
- **CSS预处理器**: Sass/SCSS
- **HTTP客户端**: Axios
- **工具库**: Lodash-es, Day.js

#### 工程化工具
- **包管理**: pnpm 8.0+
- **代码规范**: ESLint + Prettier
- **Git提交**: Husky + Commitizen
- **类型检查**: TypeScript + Vue-tsc
- **测试**: Vitest + @vue/test-utils
- **文档**: VitePress
- **监控**: Sentry

---

## 🏗️ 项目架构

### 1. 总体架构

```
HuinongFinancial/
├── frontend/                           # 前端项目根目录
│   ├── admin/                         # OA后台管理系统 (PC端)
│   │   ├── public/                    # 静态资源
│   │   ├── src/                       # 源代码
│   │   ├── package.json               # 依赖配置
│   │   ├── vite.config.ts             # Vite配置
│   │   └── tsconfig.json              # TypeScript配置
│   ├── users/                         # 惠农APP/Web (移动端/PC端)
│   │   ├── public/                    # 静态资源
│   │   ├── src/                       # 源代码
│   │   ├── package.json               # 依赖配置
│   │   ├── vite.config.ts             # Vite配置
│   │   └── tsconfig.json              # TypeScript配置
│   ├── shared/                        # 共享代码库
│   │   ├── components/                # 通用组件
│   │   ├── utils/                     # 工具函数
│   │   ├── types/                     # 类型定义
│   │   ├── constants/                 # 常量定义
│   │   └── api/                       # API接口定义
│   ├── docs/                          # 文档目录
│   ├── tools/                         # 工具脚本
│   └── package.json                   # 根级别依赖
```

### 2. 单个应用架构

```
src/
├── api/                               # API接口层
│   ├── modules/                       # 按模块分组的API
│   │   ├── auth.ts                    # 认证相关API
│   │   ├── loan.ts                    # 贷款相关API
│   │   ├── machine.ts                 # 农机相关API
│   │   ├── content.ts                 # 内容相关API
│   │   └── user.ts                    # 用户相关API
│   ├── request.ts                     # Axios配置和拦截器
│   ├── types.ts                       # API类型定义
│   └── index.ts                       # API导出
├── assets/                            # 静态资源
│   ├── fonts/                         # 字体文件
│   ├── icons/                         # 图标文件 (SVG)
│   ├── images/                        # 图片文件
│   └── styles/                        # 样式文件
│       ├── variables.scss             # SCSS变量
│       ├── mixins.scss                # SCSS混入
│       ├── reset.scss                 # 样式重置
│       └── global.scss                # 全局样式
├── components/                        # 组件库
│   ├── basic/                         # 基础组件
│   │   ├── Button/                    # 按钮组件
│   │   │   ├── index.vue
│   │   │   ├── types.ts
│   │   │   └── styles.scss
│   │   ├── Input/                     # 输入框组件
│   │   ├── Card/                      # 卡片组件
│   │   └── index.ts                   # 组件导出
│   ├── business/                      # 业务组件
│   │   ├── LoanCard/                  # 贷款卡片组件
│   │   ├── MachineList/               # 农机列表组件
│   │   ├── StatusProgress/            # 状态进度组件
│   │   └── index.ts
│   └── layout/                        # 布局组件
│       ├── Header/                    # 头部组件
│       ├── Footer/                    # 底部组件
│       ├── Sidebar/                   # 侧边栏组件
│       └── index.ts
├── composables/                       # 组合式函数
│   ├── useAuth.ts                     # 认证逻辑
│   ├── useRequest.ts                  # 请求逻辑
│   ├── useForm.ts                     # 表单逻辑
│   ├── useDevice.ts                   # 设备检测
│   └── index.ts
├── router/                            # 路由配置
│   ├── modules/                       # 路由模块
│   │   ├── auth.ts                    # 认证路由
│   │   ├── loan.ts                    # 贷款路由
│   │   ├── machine.ts                 # 农机路由
│   │   └── user.ts                    # 用户路由
│   ├── guards.ts                      # 路由守卫
│   ├── index.ts                       # 路由配置
│   └── types.ts                       # 路由类型
├── stores/                            # 状态管理
│   ├── modules/                       # Store模块
│   │   ├── auth.ts                    # 认证状态
│   │   ├── user.ts                    # 用户状态
│   │   ├── loan.ts                    # 贷款状态
│   │   └── app.ts                     # 应用全局状态
│   ├── index.ts                       # Store配置
│   └── types.ts                       # Store类型
├── utils/                             # 工具函数
│   ├── auth.ts                        # 认证工具
│   ├── storage.ts                     # 存储工具
│   ├── validate.ts                    # 验证工具
│   ├── format.ts                      # 格式化工具
│   ├── device.ts                      # 设备工具
│   ├── request.ts                     # 请求工具
│   └── index.ts
├── views/                             # 页面组件
│   ├── auth/                          # 认证相关页面
│   │   ├── Login/                     # 登录页面
│   │   │   ├── index.vue
│   │   │   ├── components/            # 页面私有组件
│   │   │   └── composables/           # 页面私有逻辑
│   │   ├── Register/                  # 注册页面
│   │   └── ForgotPassword/            # 忘记密码页面
│   ├── loan/                          # 贷款相关页面
│   │   ├── ProductList/               # 产品列表页
│   │   ├── Application/               # 申请页面
│   │   ├── Status/                    # 状态查询页
│   │   └── History/                   # 历史记录页
│   ├── machine/                       # 农机相关页面
│   │   ├── Search/                    # 搜索页面
│   │   ├── Detail/                    # 详情页面
│   │   ├── Booking/                   # 预约页面
│   │   └── Orders/                    # 订单页面
│   ├── user/                          # 用户相关页面
│   │   ├── Profile/                   # 个人资料
│   │   ├── Settings/                  # 设置页面
│   │   └── Verification/              # 认证页面
│   └── home/                          # 首页
│       └── index.vue
├── types/                             # 类型定义
│   ├── api.ts                         # API类型
│   ├── components.ts                  # 组件类型
│   ├── router.ts                      # 路由类型
│   ├── store.ts                       # Store类型
│   └── global.d.ts                    # 全局类型声明
├── App.vue                            # 根组件
├── main.ts                            # 应用入口
└── env.d.ts                           # 环境变量类型
```

---

## ⚙️ 开发环境配置

### 1. Node.js 环境

```bash
# 推荐使用 Node.js 18.x LTS 版本
node --version  # >= 18.0.0
pnpm --version  # >= 8.0.0
```

### 2. 项目初始化

```bash
# 克隆项目
git clone https://github.com/company/HuinongFinancial.git
cd HuinongFinancial/frontend

# 安装依赖（根目录）
pnpm install

# 进入具体应用
cd users
pnpm install
pnpm dev

# 或者运行OA后台
cd admin
pnpm install
pnpm dev
```

### 3. 环境变量配置

#### 用户端（users）环境变量
```bash
# .env.development
VITE_APP_TITLE=数字惠农APP
VITE_API_BASE_URL=http://localhost:8080/api
VITE_UPLOAD_URL=http://localhost:8080/upload
VITE_APP_ENV=development
VITE_APP_VERSION=1.0.0

# .env.production
VITE_APP_TITLE=数字惠农APP
VITE_API_BASE_URL=https://api.huinong.com/api
VITE_UPLOAD_URL=https://cdn.huinong.com/upload
VITE_APP_ENV=production
VITE_APP_VERSION=1.0.0
```

#### OA后台（admin）环境变量
```bash
# .env.development
VITE_APP_TITLE=惠农OA管理系统
VITE_API_BASE_URL=http://localhost:8080/api/oa
VITE_APP_ENV=development

# .env.production
VITE_APP_TITLE=惠农OA管理系统
VITE_API_BASE_URL=https://api.huinong.com/api/oa
VITE_APP_ENV=production
```

---

## 🔧 构建配置

### 1. Vite 配置 (vite.config.ts)

```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'
import Components from 'unplugin-vue-components/vite'
import { VantResolver } from 'unplugin-vue-components/resolvers'

export default defineConfig({
  plugins: [
    vue(),
    // 自动导入组件
    Components({
      resolvers: [VantResolver()],
    }),
    // SVG图标插件
    createSvgIconsPlugin({
      iconDirs: [resolve(process.cwd(), 'src/assets/icons')],
      symbolId: 'icon-[dir]-[name]',
    }),
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '@shared': resolve(__dirname, '../shared'),
    },
  },
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@import "@/assets/styles/variables.scss";`,
      },
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['vue', 'vue-router', 'pinia'],
          ui: ['vant', 'ant-design-vue'],
        },
      },
    },
  },
})
```

### 2. TypeScript 配置 (tsconfig.json)

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "preserve",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"],
      "@shared/*": ["../shared/*"]
    }
  },
  "include": ["src/**/*.ts", "src/**/*.d.ts", "src/**/*.tsx", "src/**/*.vue"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

---

## 📝 编码规范

### 1. Vue组件规范

#### 1.1 组件文件命名
```bash
# 使用 PascalCase 命名
components/
├── UserProfile.vue        # ✅ 正确
├── user-profile.vue       # ❌ 错误
└── userProfile.vue        # ❌ 错误
```

#### 1.2 组件结构规范
```vue
<template>
  <div class="user-profile">
    <!-- 模板内容 -->
  </div>
</template>

<script setup lang="ts">
// 1. 导入外部依赖
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

// 2. 导入类型定义
import type { UserInfo } from '@/types/user'

// 3. 导入组合式函数
import { useAuth } from '@/composables/useAuth'

// 4. 导入组件
import UserAvatar from '@/components/basic/UserAvatar.vue'

// 5. 定义Props
interface Props {
  userId: string
  editable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  editable: false
})

// 6. 定义Emits
interface Emits {
  update: [userInfo: UserInfo]
  save: [void]
}

const emit = defineEmits<Emits>()

// 7. 响应式数据
const userInfo = ref<UserInfo>()
const loading = ref(false)

// 8. 计算属性
const displayName = computed(() => {
  return userInfo.value?.name || '未设置'
})

// 9. 方法定义
const saveUserInfo = async () => {
  try {
    loading.value = true
    // 保存逻辑
    emit('save')
  } catch (error) {
    console.error('保存失败:', error)
  } finally {
    loading.value = false
  }
}

// 10. 生命周期
onMounted(() => {
  // 初始化逻辑
})
</script>

<style lang="scss" scoped>
.user-profile {
  padding: 16px;
  
  &__avatar {
    margin-bottom: 16px;
  }
  
  &__form {
    // 表单样式
  }
}
</style>
```

### 2. TypeScript 规范

#### 2.1 类型定义
```typescript
// types/user.ts
export interface UserInfo {
  id: string
  name: string
  phone: string
  email?: string
  avatar?: string
  status: UserStatus
  createdAt: string
  updatedAt: string
}

export enum UserStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  SUSPENDED = 'suspended'
}

export type UserRole = 'farmer' | 'farm_owner' | 'cooperative' | 'enterprise'

// API响应类型
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  meta?: {
    total?: number
    page?: number
    limit?: number
  }
}
```

#### 2.2 API接口规范
```typescript
// api/modules/user.ts
import type { UserInfo, ApiResponse } from '@/types'
import { request } from '../request'

export const userApi = {
  // 获取用户信息
  getUserInfo(): Promise<ApiResponse<UserInfo>> {
    return request.get('/user/profile')
  },

  // 更新用户信息
  updateUserInfo(data: Partial<UserInfo>): Promise<ApiResponse<UserInfo>> {
    return request.put('/user/profile', data)
  },

  // 上传头像
  uploadAvatar(file: File): Promise<ApiResponse<{ url: string }>> {
    const formData = new FormData()
    formData.append('avatar', file)
    return request.post('/user/avatar', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  }
}
```

### 3. 样式规范

#### 3.1 BEM命名规范
```scss
// 块(Block)、元素(Element)、修饰符(Modifier)
.loan-card {                    // Block
  padding: 16px;
  border-radius: 8px;

  &__header {                   // Element
    display: flex;
    justify-content: space-between;
    margin-bottom: 12px;
  }

  &__title {                    // Element
    font-size: 16px;
    font-weight: 600;
  }

  &__amount {                   // Element
    font-size: 18px;
    color: var(--primary-color);
  }

  &--featured {                 // Modifier
    border: 2px solid var(--primary-color);
  }

  &--disabled {                 // Modifier
    opacity: 0.6;
    pointer-events: none;
  }
}
```

#### 3.2 SCSS变量使用
```scss
// assets/styles/variables.scss
:root {
  // 颜色变量
  --primary-color: #52C41A;
  --secondary-color: #FAAD14;
  --success-color: #52C41A;
  --warning-color: #FAAD14;
  --error-color: #FF4D4F;
  --info-color: #1890FF;

  // 字体变量
  --font-size-xs: 10px;
  --font-size-sm: 12px;
  --font-size-base: 14px;
  --font-size-lg: 16px;
  --font-size-xl: 18px;

  // 间距变量
  --spacing-xs: 4px;
  --spacing-sm: 8px;
  --spacing-base: 16px;
  --spacing-lg: 24px;
  --spacing-xl: 32px;

  // 圆角变量
  --border-radius-sm: 4px;
  --border-radius-base: 6px;
  --border-radius-lg: 8px;
  --border-radius-xl: 12px;
}
```

---

## 🔒 状态管理规范

### 1. Pinia Store 结构

```typescript
// stores/modules/auth.ts
import { defineStore } from 'pinia'
import type { UserInfo } from '@/types'
import { authApi } from '@/api'

interface AuthState {
  token: string | null
  userInfo: UserInfo | null
  isLoggedIn: boolean
  permissions: string[]
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: localStorage.getItem('token'),
    userInfo: null,
    isLoggedIn: false,
    permissions: []
  }),

  getters: {
    // 是否已认证
    isAuthenticated: (state) => !!state.token && state.isLoggedIn,
    
    // 用户角色
    userRole: (state) => state.userInfo?.role,
    
    // 是否有特定权限
    hasPermission: (state) => (permission: string) => {
      return state.permissions.includes(permission)
    }
  },

  actions: {
    // 登录
    async login(credentials: LoginCredentials) {
      try {
        const response = await authApi.login(credentials)
        const { token, user } = response.data
        
        this.token = token
        this.userInfo = user
        this.isLoggedIn = true
        
        localStorage.setItem('token', token)
        
        return response
      } catch (error) {
        this.logout()
        throw error
      }
    },

    // 登出
    logout() {
      this.token = null
      this.userInfo = null
      this.isLoggedIn = false
      this.permissions = []
      
      localStorage.removeItem('token')
    },

    // 获取用户信息
    async fetchUserInfo() {
      if (!this.token) return
      
      try {
        const response = await authApi.getUserInfo()
        this.userInfo = response.data
        this.isLoggedIn = true
      } catch (error) {
        this.logout()
        throw error
      }
    }
  }
})
```

### 2. 组合式函数规范

```typescript
// composables/useAuth.ts
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores'

export function useAuth() {
  const authStore = useAuthStore()
  const router = useRouter()

  // 计算属性
  const isLoggedIn = computed(() => authStore.isAuthenticated)
  const userInfo = computed(() => authStore.userInfo)

  // 登录方法
  const login = async (credentials: LoginCredentials) => {
    try {
      await authStore.login(credentials)
      await router.push('/home')
    } catch (error) {
      console.error('登录失败:', error)
      throw error
    }
  }

  // 登出方法
  const logout = async () => {
    authStore.logout()
    await router.push('/login')
  }

  // 检查权限
  const hasPermission = (permission: string) => {
    return authStore.hasPermission(permission)
  }

  return {
    isLoggedIn,
    userInfo,
    login,
    logout,
    hasPermission
  }
}
```

---

## 🛡️ 错误处理与监控

### 1. 全局错误处理

```typescript
// utils/error.ts
export class ApiError extends Error {
  constructor(
    public code: number,
    public message: string,
    public data?: any
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

export function handleApiError(error: any): ApiError {
  if (error.response) {
    const { status, data } = error.response
    return new ApiError(status, data.message || '请求失败', data)
  } else if (error.request) {
    return new ApiError(0, '网络连接失败')
  } else {
    return new ApiError(-1, error.message || '未知错误')
  }
}

// 全局错误处理器
export function setupErrorHandler(app: App) {
  app.config.errorHandler = (err, vm, info) => {
    console.error('Vue错误:', err, vm, info)
    
    // 发送错误到监控服务
    if (import.meta.env.PROD) {
      // Sentry.captureException(err)
    }
  }

  window.addEventListener('unhandledrejection', (event) => {
    console.error('未处理的Promise拒绝:', event.reason)
    
    // 发送错误到监控服务
    if (import.meta.env.PROD) {
      // Sentry.captureException(event.reason)
    }
  })
}
```

### 2. 请求拦截器

```typescript
// api/request.ts
import axios from 'axios'
import type { AxiosResponse, AxiosError } from 'axios'
import { useAuthStore } from '@/stores'
import { handleApiError } from '@/utils/error'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000,
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    
    // 添加认证头
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    
    // 添加设备信息
    config.headers['X-Device-Type'] = 'web'
    config.headers['X-App-Version'] = import.meta.env.VITE_APP_VERSION
    
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const { code, message, data } = response.data
    
    // 统一处理业务错误
    if (code !== 200) {
      const error = new ApiError(code, message, data)
      return Promise.reject(error)
    }
    
    return response.data
  },
  (error: AxiosError) => {
    const apiError = handleApiError(error)
    
    // 401 未授权，跳转登录
    if (apiError.code === 401) {
      const authStore = useAuthStore()
      authStore.logout()
      window.location.href = '/login'
    }
    
    return Promise.reject(apiError)
  }
)

export { request }
```

---

## 🧪 测试规范

### 1. 单元测试

```typescript
// tests/components/Button.test.ts
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Button from '@/components/basic/Button.vue'

describe('Button', () => {
  it('renders properly', () => {
    const wrapper = mount(Button, {
      props: { type: 'primary' },
      slots: { default: 'Click me' }
    })
    
    expect(wrapper.text()).toContain('Click me')
    expect(wrapper.classes()).toContain('btn-primary')
  })

  it('emits click event', async () => {
    const wrapper = mount(Button)
    await wrapper.trigger('click')
    
    expect(wrapper.emitted()).toHaveProperty('click')
  })

  it('is disabled when loading', () => {
    const wrapper = mount(Button, {
      props: { loading: true }
    })
    
    expect(wrapper.find('button').attributes('disabled')).toBeDefined()
  })
})
```

### 2. E2E测试

```typescript
// tests/e2e/login.spec.ts
import { test, expect } from '@playwright/test'

test.describe('登录流程', () => {
  test('用户可以成功登录', async ({ page }) => {
    await page.goto('/login')
    
    // 填写登录信息
    await page.fill('[data-testid="phone-input"]', '13800138000')
    await page.fill('[data-testid="password-input"]', 'password123')
    
    // 点击登录按钮
    await page.click('[data-testid="login-button"]')
    
    // 验证跳转到首页
    await expect(page).toHaveURL('/home')
    await expect(page.locator('[data-testid="user-name"]')).toBeVisible()
  })

  test('输入错误密码显示错误信息', async ({ page }) => {
    await page.goto('/login')
    
    await page.fill('[data-testid="phone-input"]', '13800138000')
    await page.fill('[data-testid="password-input"]', 'wrongpassword')
    await page.click('[data-testid="login-button"]')
    
    await expect(page.locator('[data-testid="error-message"]')).toContainText('密码错误')
  })
})
```

---

## 📦 构建与部署

### 1. 构建脚本

```json
{
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc --noEmit && vite build",
    "build:dev": "vite build --mode development",
    "build:test": "vite build --mode testing",
    "build:prod": "vite build --mode production",
    "preview": "vite preview",
    "test": "vitest",
    "test:ui": "vitest --ui",
    "test:e2e": "playwright test",
    "lint": "eslint . --ext .vue,.js,.jsx,.cjs,.mjs,.ts,.tsx,.cts,.mts --fix",
    "type-check": "vue-tsc --noEmit",
    "analyze": "vite-bundle-analyzer"
  }
}
```

### 2. Docker配置

```dockerfile
# Dockerfile
FROM node:18-alpine as builder

WORKDIR /app

# 复制package文件
COPY package*.json pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install

# 复制源代码
COPY . .

# 构建应用
RUN pnpm build

# 生产环境
FROM nginx:alpine

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制Nginx配置
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### 3. CI/CD 配置

```yaml
# .github/workflows/deploy.yml
name: Deploy Frontend

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          
      - name: Install pnpm
        run: npm install -g pnpm
        
      - name: Install dependencies
        run: pnpm install
        
      - name: Run tests
        run: pnpm test
        
      - name: Type check
        run: pnpm type-check
        
      - name: Lint
        run: pnpm lint

  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          
      - name: Install pnpm
        run: npm install -g pnpm
        
      - name: Install dependencies
        run: pnpm install
        
      - name: Build
        run: pnpm build:prod
        
      - name: Deploy to OSS
        run: |
          # 部署到阿里云OSS或其他CDN
          echo "部署到生产环境"
```

---

## 🔍 代码质量保证

### 1. ESLint 配置

```javascript
// .eslintrc.js
module.exports = {
  root: true,
  env: {
    node: true,
    'vue/setup-compiler-macros': true
  },
  extends: [
    'plugin:vue/vue3-essential',
    '@vue/eslint-config-typescript',
    '@vue/eslint-config-prettier'
  ],
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'vue/multi-word-component-names': 'off',
    'vue/component-tags-order': [
      'error',
      {
        order: ['template', 'script', 'style']
      }
    ],
    '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }]
  }
}
```

### 2. Prettier 配置

```json
{
  "semi": false,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5",
  "printWidth": 80,
  "endOfLine": "lf",
  "vueIndentScriptAndStyle": true
}
```

### 3. Git Hooks

```javascript
// .husky/pre-commit
#!/usr/bin/env sh
. "$(dirname -- "$0")/_/husky.sh"

pnpm lint-staged
```

```json
{
  "lint-staged": {
    "*.{js,jsx,ts,tsx,vue}": [
      "eslint --fix",
      "prettier --write"
    ],
    "*.{css,scss,less}": [
      "prettier --write"
    ]
  }
}
```

---

## 📊 性能优化

### 1. 代码分割

```typescript
// router/index.ts
const routes = [
  {
    path: '/loan',
    component: () => import('@/views/loan/ProductList.vue'),
    meta: { title: '贷款产品' }
  },
  {
    path: '/machine',
    component: () => import('@/views/machine/Search.vue'),
    meta: { title: '农机租赁' }
  }
]
```

### 2. 组件懒加载

```vue
<template>
  <div>
    <Suspense>
      <template #default>
        <AsyncComponent />
      </template>
      <template #fallback>
        <div>加载中...</div>
      </template>
    </Suspense>
  </div>
</template>

<script setup lang="ts">
import { defineAsyncComponent } from 'vue'

const AsyncComponent = defineAsyncComponent(
  () => import('@/components/business/LoanCard.vue')
)
</script>
```

### 3. 图片优化

```typescript
// utils/image.ts
export function generateImageUrl(
  url: string,
  options: {
    width?: number
    height?: number
    quality?: number
    format?: 'webp' | 'jpeg' | 'png'
  } = {}
) {
  const { width, height, quality = 80, format = 'webp' } = options
  
  // 如果是OSS链接，添加处理参数
  if (url.includes('aliyuncs.com')) {
    const params = []
    if (width) params.push(`w_${width}`)
    if (height) params.push(`h_${height}`)
    params.push(`q_${quality}`)
    params.push(`f_${format}`)
    
    return `${url}?x-oss-process=image/resize,${params.join(',')}`
  }
  
  return url
}
```

---

## 🛠️ 开发工具

### 1. VS Code 配置

```json
// .vscode/settings.json
{
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  },
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "typescript.preferences.importModuleSpecifier": "relative",
  "typescript.suggest.autoImports": true,
  "vue.format.enable": false
}
```

### 2. Chrome 开发者工具插件

- Vue DevTools
- Redux DevTools（用于Pinia）
- axe DevTools（无障碍检测）
- Lighthouse（性能检测）

### 3. 调试配置

```json
// .vscode/launch.json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Chrome",
      "type": "chrome",
      "request": "launch",
      "url": "http://localhost:3000",
      "webRoot": "${workspaceFolder}/src",
      "sourceMapPathOverrides": {
        "webpack:///src/*": "${webRoot}/*"
      }
    }
  ]
}
```

---

## 📖 文档规范

### 1. 组件文档

```vue
<!-- components/basic/Button.vue -->
<template>
  <!-- 组件模板 -->
</template>

<script setup lang="ts">
/**
 * 基础按钮组件
 * 
 * @component Button
 * @description 提供多种样式的按钮组件，支持加载状态、禁用状态等
 * 
 * @example
 * ```vue
 * <Button type="primary" @click="handleClick">
 *   点击我
 * </Button>
 * ```
 */

interface Props {
  /** 按钮类型 */
  type?: 'primary' | 'secondary' | 'danger'
  /** 按钮大小 */
  size?: 'small' | 'medium' | 'large'
  /** 是否加载中 */
  loading?: boolean
  /** 是否禁用 */
  disabled?: boolean
}

interface Emits {
  /** 点击事件 */
  click: [event: MouseEvent]
}
</script>
```

### 2. API文档

```typescript
/**
 * 用户相关API
 * @module UserAPI
 */

/**
 * 获取用户信息
 * @param userId - 用户ID
 * @returns 用户信息
 * 
 * @example
 * ```typescript
 * const userInfo = await getUserInfo('123')
 * console.log(userInfo.name)
 * ```
 */
export async function getUserInfo(userId: string): Promise<UserInfo> {
  // 实现
}
```

---

## 🔄 版本管理

### 1. 版本号规范

- 主版本号：不兼容的API修改
- 次版本号：向下兼容的功能性新增
- 修订号：向下兼容的问题修正

### 2. 发布流程

```bash
# 1. 确保代码已经提交
git add .
git commit -m "feat: 新增用户管理功能"

# 2. 更新版本号
npm version patch  # 修订版本
npm version minor  # 次版本
npm version major  # 主版本

# 3. 推送代码和标签
git push origin main
git push origin --tags

# 4. 创建发布说明
gh release create v1.0.0 --notes "发布说明"
```

---

## 🎯 最佳实践

### 1. 组件设计原则

- **单一职责**：每个组件只负责一个功能
- **可复用性**：组件应该易于在不同场景中复用
- **可配置性**：通过props提供灵活的配置选项
- **可扩展性**：支持插槽和事件扩展功能

### 2. 性能优化建议

- 使用`v-memo`优化重复渲染
- 合理使用`shallowRef`和`shallowReactive`
- 避免在模板中进行复杂计算
- 使用`keep-alive`缓存组件状态

### 3. 代码组织建议

- 按功能模块组织文件结构
- 使用绝对路径导入避免相对路径混乱
- 统一命名规范和代码风格
- 编写有意义的注释和文档

---

本工程化文档将随着项目发展持续更新，确保开发规范的时效性和实用性。所有团队成员都应遵循本文档的规范，以保证代码质量和项目的可维护性。 

## 📋 项目架构

本项目采用现代化前端开发架构，基于Vue 3 + TypeScript构建，支持用户端和管理端的双端应用。

## 🏗️ 目录结构

```
frontend/
├── users/                    # 用户端应用
│   ├── src/
│   │   ├── views/           # 页面组件
│   │   │   ├── FinancePage.vue          # 理财页面 ✨新增完善
│   │   │   ├── LoanApplicationPage.vue   # 贷款申请页面
│   │   │   ├── IndexPage.vue            # 首页
│   │   │   └── ...
│   │   ├── services/        # API服务
│   │   │   └── api.ts       # API接口定义 ✨更新完善
│   │   ├── stores/          # 状态管理
│   │   ├── router/          # 路由配置
│   │   └── components/      # 通用组件
│   └── package.json
├── admin/                   # 管理端应用
│   └── ...
```

## 🔧 技术栈

### 核心框架
- **Vue 3**: 响应式框架，支持Composition API
- **TypeScript**: 类型安全的JavaScript超集
- **Vite**: 现代化构建工具，支持热重载
- **Element Plus**: UI组件库

### 状态管理与路由
- **Pinia**: 轻量级状态管理
- **Vue Router 4**: 客户端路由

### 开发工具
- **ESLint**: 代码质量检查
- **Prettier**: 代码格式化
- **Husky**: Git钩子管理

## 📱 页面功能模块

### 1. 理财页面模块 ✨完善重点

#### 功能架构
```
FinancePage.vue
├── 用户状态管理
│   ├── 登录状态检测
│   ├── 用户信息展示
│   └── 权限验证
├── 产品数据管理
│   ├── 智能数据源 (开发/生产环境)
│   ├── 产品分类筛选
│   └── 实时数据加载
├── UI交互组件
│   ├── 产品卡片展示
│   ├── 分类选择器
│   └── 操作按钮组
└── 路由集成
    ├── 产品详情跳转
    ├── 申请页面跳转
    └── 登录页面引导
```

#### 数据流管理
```typescript
// 数据源智能切换
const getProducts = async () => {
  // 开发环境: 模拟数据
  if (isDevelopment) {
    return getMockProducts()
  }
  
  // 生产环境: 真实API
  return await loanApi.getProducts()
}
```

### 2. API服务架构

#### 服务分层
```typescript
// services/api.ts
export const loanApi = {
  // 贷款产品管理
  getProducts(category?: string): Promise<ApiResponse<LoanProduct[]>>
  getProductDetail(id: string): Promise<ApiResponse<LoanProduct>>
  
  // 贷款申请管理
  submitApplication(data: LoanApplicationRequest): Promise<ApiResponse>
  getMyApplications(): Promise<PaginatedResponse<LoanApplication>>
  
  // 模拟数据系统 ✨新增
  getMockProducts(): Promise<ApiResponse<LoanProduct[]>>
}
```

#### 模拟数据系统 ✨新增特性
```typescript
// 完善的模拟数据
const mockProducts = [
  {
    product_id: 'NYCD001',
    name: '农业创业贷',
    category: '创业贷',
    min_amount: 50000,
    max_amount: 500000,
    interest_rate_yearly: '6.5%',
    // ... 完整产品信息
  }
  // ... 6款产品
]

// 智能环境检测
const isDev = import.meta.env?.DEV || window.location.hostname === 'localhost'
```

## 🎨 UI/UX设计规范

### 设计原则
1. **用户优先**: 以用户体验为核心
2. **响应式**: 适配多端设备
3. **一致性**: 统一的视觉语言
4. **可访问性**: 满足无障碍标准

### 色彩系统
```scss
// 主色调
$primary-color: #27ae60;      // 惠农绿
$primary-light: #2ecc71;      // 浅绿色
$primary-dark: #219a52;       // 深绿色

// 辅助色
$success-color: #67c23a;      // 成功绿
$warning-color: #e6a23c;      // 警告橙
$danger-color: #f56c6c;       // 危险红
$info-color: #909399;         // 信息灰

// 中性色
$text-primary: #2c3e50;       // 主文字
$text-regular: #606266;       // 常规文字
$text-secondary: #909399;     // 次要文字
$border-color: #dcdfe6;       // 边框色
$bg-color: #f5f5f5;          // 背景色
```

### 组件规范
```vue
<!-- 产品卡片组件标准 -->
<template>
  <div class="product-card">
    <div class="card-header">
      <h4 class="product-name">{{ product.name }}</h4>
      <el-tag type="success">{{ product.category }}</el-tag>
    </div>
    
    <div class="card-body">
      <div class="product-details">
        <div class="detail-item">
          <span class="label">贷款金额</span>
          <span class="value">{{ formatAmount(product.min_amount) }} - {{ formatAmount(product.max_amount) }}元</span>
        </div>
        <!-- ... 其他信息 -->
      </div>
    </div>
    
    <div class="card-footer">
      <el-button type="info" size="small">查看详情</el-button>
      <el-button type="primary" size="small">立即申请</el-button>
    </div>
  </div>
</template>

<style scoped>
.product-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.product-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}
</style>
```

## 🔗 路由配置

### 页面路由
```typescript
// router/index.ts
const routes = [
  // 理财相关路由
  {
    path: '/finance',
    name: 'finance',
    component: () => import('../views/FinancePage.vue'),
    meta: { requiresAuth: false, title: '惠农金融' }
  },
  
  // 贷款申请流程
  {
    path: '/loan/apply/:productId',
    name: 'loanApplication',
    component: () => import('../views/LoanApplicationPage.vue'),
    meta: { requiresAuth: true, title: '贷款申请' }
  },
  
  {
    path: '/loan/products/:productId',
    name: 'loanProductDetail',
    component: () => import('../views/LoanProductDetailPage.vue'),
    meta: { requiresAuth: false, title: '产品详情' }
  },
  
  // 申请管理
  {
    path: '/loan/my-applications',
    name: 'myLoanApplications',
    component: () => import('../views/MyLoanApplicationsPage.vue'),
    meta: { requiresAuth: true, title: '我的申请' }
  }
]
```

### 路由守卫
```typescript
// 认证守卫
router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
  } else {
    next()
  }
})
```

## 📊 状态管理

### 用户状态管理
```typescript
// stores/user.ts
export const useUserStore = defineStore('user', () => {
  const userInfo = ref<UserInfo | null>(null)
  const isLoggedIn = computed(() => !!userInfo.value)
  
  // 用户信息展示
  const getUserTypeDisplay = computed(() => {
    const typeMap = {
      farmer: '农户',
      farm_owner: '农场主',
      cooperative: '合作社',
      enterprise: '企业'
    }
    return typeMap[userInfo.value?.user_type] || '普通用户'
  })
  
  // 认证状态
  const getAuthStatus = computed(() => {
    // 返回认证状态信息
  })
  
  return {
    userInfo,
    isLoggedIn,
    getUserTypeDisplay,
    getAuthStatus
  }
})
```

## 🛠️ 开发工具配置

### Vite配置
```typescript
// vite.config.ts
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src')
    }
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://172.18.120.10:8080',
        changeOrigin: true
      }
    }
  }
})
```

### TypeScript配置
```json
// tsconfig.json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "module": "ESNext",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "preserve",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"]
    }
  },
  "include": ["src/**/*.ts", "src/**/*.d.ts", "src/**/*.tsx", "src/**/*.vue"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

## 🚀 构建与部署

### 开发环境
```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 类型检查
npm run type-check

# 代码检查
npm run lint
```

### 生产环境
```bash
# 构建项目
npm run build

# 预览构建结果
npm run preview

# 部署到服务器
npm run deploy
```

### 环境变量
```env
# .env.development
VITE_API_BASE_URL=http://172.18.120.10:8080/api
VITE_APP_TITLE=惠农金融-开发环境

# .env.production
VITE_API_BASE_URL=https://api.huinong.com/api
VITE_APP_TITLE=惠农金融
```

## 📈 性能优化

### 代码分割
```typescript
// 路由懒加载
const routes = [
  {
    path: '/finance',
    component: () => import('../views/FinancePage.vue')
  }
]

// 组件懒加载
const AsyncComponent = defineAsyncComponent(() => import('./HeavyComponent.vue'))
```

### 数据缓存
```typescript
// API响应缓存
const cache = new Map()

const getCachedData = async (key: string, fetcher: () => Promise<any>) => {
  if (cache.has(key)) {
    return cache.get(key)
  }
  
  const data = await fetcher()
  cache.set(key, data)
  return data
}
```

## 🔍 测试策略

### 单元测试
```typescript
// 使用Vitest进行单元测试
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import FinancePage from '@/views/FinancePage.vue'

describe('FinancePage', () => {
  it('渲染产品列表', () => {
    const wrapper = mount(FinancePage)
    expect(wrapper.find('.product-list').exists()).toBe(true)
  })
})
```

### 端到端测试
```typescript
// 使用Playwright进行E2E测试
import { test, expect } from '@playwright/test'

test('贷款申请流程', async ({ page }) => {
  await page.goto('/finance')
  await page.click('.product-card .apply-button')
  await expect(page).toHaveURL(/\/loan\/apply/)
})
```

## 📱 移动端适配

### 响应式设计
```scss
// 移动端优先
.finance-page {
  padding: 16px;
  
  @media (min-width: 768px) {
    padding: 24px;
    max-width: 600px;
    margin: 0 auto;
  }
  
  @media (min-width: 1200px) {
    max-width: 800px;
  }
}
```

### 移动端优化
- Touch事件支持
- 滚动性能优化
- 图片懒加载
- 网络状态感知

## 🔐 安全考虑

### 前端安全
```typescript
// XSS防护
const sanitizeHtml = (html: string) => {
  // 使用DOMPurify清理HTML
  return DOMPurify.sanitize(html)
}

// CSRF防护
const apiClient = axios.create({
  headers: {
    'X-Requested-With': 'XMLHttpRequest'
  }
})
```

### 数据验证
```typescript
// 输入验证
const validateAmount = (amount: number, min: number, max: number) => {
  if (amount < min || amount > max) {
    throw new Error(`金额必须在${min}-${max}之间`)
  }
}
```

## 📚 开发规范

### 组件命名
- 页面组件: `PascalCase` + `Page` 后缀
- 通用组件: `PascalCase`
- 工具组件: `camelCase`

### 文件组织
```
src/
├── views/           # 页面组件
├── components/      # 通用组件
├── services/        # API服务
├── stores/          # 状态管理
├── utils/           # 工具函数
├── types/           # TypeScript类型定义
└── styles/          # 全局样式
```

### 代码风格
- 使用ESLint + Prettier
- 遵循Vue 3 Style Guide
- TypeScript严格模式
- 组件Props类型定义

## 🔄 持续集成

### CI/CD流程
```yaml
# .github/workflows/deploy.yml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - run: npm install
      - run: npm run build
      - run: npm run test
      - name: Deploy
        run: npm run deploy
```

## 📋 待办事项

### 近期优化
- [ ] 产品对比功能
- [ ] 收藏/关注功能
- [ ] 申请进度可视化
- [ ] 消息推送集成

### 长期规划
- [ ] PWA支持
- [ ] 国际化(i18n)
- [ ] 主题切换
- [ ] 无障碍优化

---

**文档版本**: v2.0  
**最后更新**: 2024年1月15日  
**维护者**: 前端开发团队 