# OA管理系统前端开发规范

## 概述

本文档定义了数字惠农OA管理系统前端开发的技术规范、代码标准和最佳实践，确保代码的一致性、可维护性和高质量。

## 技术架构

### 核心技术栈
- **框架**: Vue 3 (Composition API)
- **构建工具**: Vite 4+
- **语言**: TypeScript 5+
- **UI库**: Element Plus 2.4+
- **路由**: Vue Router 4+
- **状态管理**: Pinia 2+
- **HTTP客户端**: Axios 1.6+
- **工具库**: Day.js、lodash-es

### 项目结构规范

```
src/
├── api/                    # API接口层
│   ├── index.ts           # Axios配置和拦截器
│   ├── admin.ts           # 管理系统API
│   └── types.ts           # API相关类型定义
├── assets/                # 静态资源
│   ├── fonts/            # 字体文件
│   ├── images/           # 图片资源
│   └── styles/           # 全局样式
├── components/            # 通用组件
│   ├── common/           # 基础组件
│   └── business/         # 业务组件
├── composables/           # 组合式函数
├── directives/            # 自定义指令
├── hooks/                 # 可复用逻辑钩子
├── router/                # 路由配置
├── stores/                # 状态管理
├── types/                 # TypeScript类型定义
├── utils/                 # 工具函数
├── views/                 # 页面组件
├── App.vue               # 根组件
└── main.ts               # 应用入口
```

## 代码规范

### 命名规范

#### 文件命名
- **组件文件**: PascalCase，如 `UserList.vue`
- **页面文件**: PascalCase + View后缀，如 `DashboardView.vue`
- **工具文件**: camelCase，如 `formatUtils.ts`
- **常量文件**: UPPER_SNAKE_CASE，如 `API_CONSTANTS.ts`

#### 变量命名
```typescript
// 常量 - UPPER_SNAKE_CASE
const API_BASE_URL = 'http://localhost:8080/api/v1'
const MAX_FILE_SIZE = 5 * 1024 * 1024

// 变量和函数 - camelCase
const userName = ref('')
const isLoading = ref(false)

// 组件名 - PascalCase
const UserProfile = defineComponent({})

// 类型定义 - PascalCase
interface UserInfo {
  id: string
  name: string
}
```

### TypeScript规范

#### 类型定义
```typescript
// 接口定义
interface AdminUser {
  admin_user_id: string
  username: string
  role: 'ADMIN' | '审批员'
  display_name: string
  email: string
  status: number
  created_at: string
  updated_at: string
}

// 类型别名
type UserRole = 'ADMIN' | '审批员'
type LoadingState = 'idle' | 'loading' | 'success' | 'error'

// 枚举
enum ApplicationStatus {
  PENDING = 'pending',
  APPROVED = 'approved',
  REJECTED = 'rejected'
}
```

#### 函数签名
```typescript
// API函数类型定义
const getApplicationDetail = (
  applicationId: string
): Promise<ApplicationDetail> => {
  return api.get(`/admin/loans/applications/${applicationId}`)
}

// 组件props类型
interface Props {
  title: string
  data: ApplicationDetail[]
  loading?: boolean
  onRefresh?: () => void
}
```

### Vue组件规范

#### Composition API使用
```typescript
<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'

// 响应式数据
const loading = ref(false)
const formData = reactive({
  username: '',
  password: ''
})

// 计算属性
const isFormValid = computed(() => {
  return formData.username.length > 0 && formData.password.length > 0
})

// 方法
const handleSubmit = async () => {
  loading.value = true
  try {
    await submitForm(formData)
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    loading.value = false
  }
}

// 生命周期
onMounted(() => {
  initializeData()
})
</script>
```

#### 模板规范
```vue
<template>
  <div class="user-management">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2 class="page-title">用户管理</h2>
      <div class="header-actions">
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          添加用户
        </el-button>
      </div>
    </div>

    <!-- 内容区域 -->
    <el-card class="content-card" shadow="never">
      <el-table :data="users" v-loading="loading">
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button 
              size="small" 
              type="danger" 
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>
```

#### 样式规范
```vue
<style scoped>
.user-management {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  margin: 0;
  color: #333;
  font-size: 24px;
  font-weight: 600;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.content-card {
  border-radius: 8px;
}

/* 深度选择器使用 :deep() */
:deep(.el-table__header) {
  background-color: #f8f9fa;
}
</style>
```

### API集成规范

#### 接口定义
```typescript
// api/admin.ts
export const getUsers = (params: {
  page?: number
  limit?: number
  role?: string
}): Promise<PaginationResponse<AdminUser>> => {
  return api.get('/admin/users', { params })
}

export const createUser = (data: {
  username: string
  password: string
  role: string
  display_name: string
  email: string
}): Promise<AdminUser> => {
  return api.post('/admin/users', data)
}
```

#### 错误处理
```typescript
// api/index.ts - 响应拦截器
api.interceptors.response.use(
  (response) => {
    const { code, message, data } = response.data
    if (code === 0) {
      return data
    } else {
      ElMessage.error(message || '请求失败')
      return Promise.reject(new Error(message))
    }
  },
  (error) => {
    if (error.response?.status === 401) {
      ElMessage.error('登录已过期，请重新登录')
      router.push('/login')
    } else if (error.response?.status === 403) {
      ElMessage.error('没有权限访问该资源')
    } else {
      ElMessage.error(error.message || '网络错误')
    }
    return Promise.reject(error)
  }
)
```

### 状态管理规范

#### Pinia Store定义
```typescript
// stores/auth.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { AdminUser, LoginResponse } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  // State
  const token = ref<string>(localStorage.getItem('admin_token') || '')
  const user = ref<AdminUser | null>(null)

  // Getters
  const isLoggedIn = computed(() => !!token.value)
  const userRole = computed(() => user.value?.role)

  // Actions
  const login = async (credentials: LoginRequest) => {
    const response = await adminLogin(credentials)
    token.value = response.token
    user.value = response.user
    localStorage.setItem('admin_token', response.token)
  }

  const logout = () => {
    token.value = ''
    user.value = null
    localStorage.removeItem('admin_token')
  }

  const hasPermission = (permission: string) => {
    if (!user.value) return false
    // 权限检查逻辑
    return true
  }

  return {
    token,
    user,
    isLoggedIn,
    userRole,
    login,
    logout,
    hasPermission
  }
})
```

### 路由配置规范

#### 路由定义
```typescript
// router/index.ts
const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/LoginView.vue'),
    meta: { 
      requiresAuth: false,
      title: '登录'
    }
  },
  {
    path: '/',
    name: 'layout',
    component: () => import('@/views/LayoutView.vue'),
    meta: { requiresAuth: true },
    redirect: '/dashboard',
    children: [
      {
        path: '/dashboard',
        name: 'dashboard',
        component: () => import('@/views/DashboardView.vue'),
        meta: { 
          requiresAuth: true,
          title: '工作台',
          permission: 'dashboard:view'
        }
      }
    ]
  }
]
```

#### 路由守卫
```typescript
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // 认证检查
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next('/login')
    return
  }
  
  // 权限检查
  if (to.meta.permission && !authStore.hasPermission(to.meta.permission)) {
    next('/403')
    return
  }
  
  next()
})
```

## 组件开发规范

### 组件设计原则
1. **单一职责**: 每个组件只负责一个功能
2. **可复用性**: 通过props和插槽实现复用
3. **可组合性**: 组件可以组合使用
4. **可测试性**: 组件逻辑清晰，易于测试

### 通用组件示例
```vue
<!-- components/common/DataTable.vue -->
<template>
  <div class="data-table">
    <el-table :data="data" v-loading="loading" v-bind="$attrs">
      <slot></slot>
    </el-table>
    
    <div v-if="showPagination" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="pageSizes"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  data: any[]
  loading?: boolean
  showPagination?: boolean
  total?: number
  pageSizes?: number[]
}

interface Emits {
  (e: 'page-change', page: number, size: number): void
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  showPagination: true,
  total: 0,
  pageSizes: () => [10, 20, 50, 100]
})

const emit = defineEmits<Emits>()

const currentPage = ref(1)
const pageSize = ref(20)

const handleSizeChange = () => {
  emit('page-change', currentPage.value, pageSize.value)
}

const handleCurrentChange = () => {
  emit('page-change', currentPage.value, pageSize.value)
}
</script>
```

## 性能优化规范

### 懒加载
```typescript
// 路由懒加载
const Dashboard = () => import('@/views/DashboardView.vue')

// 组件懒加载
const AsyncComponent = defineAsyncComponent(() => 
  import('@/components/HeavyComponent.vue')
)
```

### 防抖和节流
```typescript
import { debounce, throttle } from 'lodash-es'

// 搜索防抖
const handleSearch = debounce((keyword: string) => {
  searchUsers(keyword)
}, 300)

// 滚动节流
const handleScroll = throttle(() => {
  updateScrollPosition()
}, 100)
```

### 虚拟滚动
```vue
<template>
  <el-virtual-list
    :data="largeDataList"
    :height="400"
    :item-size="50"
  >
    <template #default="{ item }">
      <div class="list-item">{{ item.name }}</div>
    </template>
  </el-virtual-list>
</template>
```

## 错误处理规范

### 全局错误处理
```typescript
// main.ts
app.config.errorHandler = (err, vm, info) => {
  console.error('Vue错误:', err)
  console.error('错误信息:', info)
  // 发送错误报告
}

// 捕获Promise错误
window.addEventListener('unhandledrejection', event => {
  console.error('未处理的Promise拒绝:', event.reason)
  event.preventDefault()
})
```

### 组件错误边界
```vue
<template>
  <div>
    <div v-if="error" class="error-boundary">
      <h3>抱歉，出现了错误</h3>
      <p>{{ error.message }}</p>
      <el-button @click="retry">重试</el-button>
    </div>
    <slot v-else></slot>
  </div>
</template>

<script setup lang="ts">
const error = ref<Error | null>(null)

const retry = () => {
  error.value = null
  // 重新加载组件
}

onErrorCaptured((err) => {
  error.value = err
  return false
})
</script>
```

## 测试规范

### 单元测试
```typescript
// tests/components/UserList.test.ts
import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import UserList from '@/components/UserList.vue'

describe('UserList', () => {
  it('renders user list correctly', () => {
    const users = [
      { id: '1', name: 'John Doe', email: 'john@example.com' }
    ]
    
    const wrapper = mount(UserList, {
      props: { users }
    })
    
    expect(wrapper.text()).toContain('John Doe')
    expect(wrapper.text()).toContain('john@example.com')
  })
  
  it('emits edit event when edit button clicked', async () => {
    const wrapper = mount(UserList)
    await wrapper.find('.edit-button').trigger('click')
    
    expect(wrapper.emitted()).toHaveProperty('edit')
  })
})
```

### E2E测试
```typescript
// tests/e2e/login.spec.ts
import { test, expect } from '@playwright/test'

test('user can login successfully', async ({ page }) => {
  await page.goto('/login')
  
  await page.fill('[data-testid="username"]', 'admin')
  await page.fill('[data-testid="password"]', 'admin123')
  await page.click('[data-testid="login-button"]')
  
  await expect(page).toHaveURL('/dashboard')
  await expect(page.locator('h2')).toContainText('工作台')
})
```

## 构建和部署规范

### 环境配置
```typescript
// vite.config.ts
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['vue', 'vue-router', 'pinia'],
          element: ['element-plus'],
          utils: ['axios', 'dayjs', 'lodash-es']
        }
      }
    }
  },
  optimizeDeps: {
    include: ['element-plus/es/locale/lang/zh-cn']
  }
})
```

### 代码质量检查
```json
// package.json
{
  "scripts": {
    "lint": "eslint src --ext .vue,.js,.ts --fix",
    "type-check": "vue-tsc --noEmit",
    "test": "vitest",
    "test:e2e": "playwright test"
  }
}
```

## 安全规范

### XSS防护
```typescript
// 使用v-html时进行内容清理
import DOMPurify from 'dompurify'

const sanitizeHtml = (dirty: string) => {
  return DOMPurify.sanitize(dirty)
}
```

### CSRF防护
```typescript
// api/index.ts
api.defaults.xsrfCookieName = 'XSRF-TOKEN'
api.defaults.xsrfHeaderName = 'X-XSRF-TOKEN'
```

### 敏感信息处理
```typescript
// 敏感信息脱敏
const maskSensitiveInfo = (info: string, start = 3, end = 4) => {
  if (info.length <= start + end) return info
  return info.slice(0, start) + '*'.repeat(info.length - start - end) + info.slice(-end)
}

// 身份证号脱敏
const maskedIdCard = maskSensitiveInfo(idCard, 4, 4)
```

## 监控和日志规范

### 性能监控
```typescript
// utils/performance.ts
export const trackPageLoad = (pageName: string) => {
  const loadTime = performance.now()
  console.log(`页面 ${pageName} 加载时间: ${loadTime}ms`)
}

export const trackUserAction = (action: string, data?: any) => {
  console.log(`用户操作: ${action}`, data)
}
```

### 错误日志
```typescript
// utils/logger.ts
export const logger = {
  error: (message: string, error?: Error) => {
    console.error(`[ERROR] ${message}`, error)
    // 发送到日志服务
  },
  warn: (message: string) => {
    console.warn(`[WARN] ${message}`)
  },
  info: (message: string) => {
    console.info(`[INFO] ${message}`)
  }
}
```

## 版本管理和发布规范

### Git提交规范
```
feat: 新功能
fix: 修复bug
docs: 文档更新
style: 代码格式调整
refactor: 代码重构
test: 测试相关
chore: 构建工具或依赖库的更新
```

### 版本发布流程
1. 更新版本号 (package.json)
2. 更新 CHANGELOG.md
3. 创建Git标签
4. 执行构建和测试
5. 部署到生产环境

## 总结

本规范文档涵盖了OA管理系统前端开发的各个方面，包括代码规范、组件开发、性能优化、测试、安全等。遵循这些规范可以确保代码质量，提高开发效率，降低维护成本。

开发团队应该定期审查和更新这些规范，以适应技术发展和项目需求的变化。 