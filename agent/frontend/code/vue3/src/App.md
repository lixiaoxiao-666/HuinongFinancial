# App.vue - 应用根组件文档

## 📋 组件概述

`App.vue` 是数字惠农系统的根组件，负责整个应用的初始化、全局状态管理、路由配置和错误处理。该组件为所有子组件提供统一的上下文环境。

### 🎯 主要职责

- **应用初始化**: 设置全局配置和初始化状态
- **认证检查**: 验证用户登录状态并处理自动登录
- **错误边界**: 捕获和处理全局错误
- **主题设置**: 管理应用主题和样式变量
- **路由容器**: 提供路由出口和导航守卫

---

## 🔧 组件实现

### Vue组件代码

```vue
<template>
  <div 
    id="app" 
    :class="[
      'app-container',
      `app-theme--${currentTheme}`,
      { 'app-loading': isInitializing }
    ]"
  >
    <!-- 全局加载遮罩 -->
    <div v-if="isInitializing" class="app-loading-overlay">
      <div class="loading-spinner">
        <van-loading type="spinner" size="24px" color="#52C41A">
          初始化中...
        </van-loading>
      </div>
    </div>

    <!-- 网络状态提示 -->
    <div v-if="!isOnline" class="network-status">
      <van-notice-bar 
        color="#ed6a0c" 
        background="#fffbe8" 
        left-icon="info-o"
      >
        网络连接已断开，请检查网络设置
      </van-notice-bar>
    </div>

    <!-- 主要内容区域 -->
    <div class="app-content">
      <router-view v-slot="{ Component, route }">
        <transition :name="transitionName" mode="out-in">
          <keep-alive :include="keepAliveComponents">
            <component 
              :is="Component" 
              :key="route.fullPath"
              @error="handleComponentError"
            />
          </keep-alive>
        </transition>
      </router-view>
    </div>

    <!-- 全局弹窗容器 -->
    <van-popup
      v-model:show="globalDialog.visible"
      :round="true"
      :closeable="true"
      @close="closeGlobalDialog"
    >
      <div class="global-dialog-content">
        <div class="dialog-header">
          <h3>{{ globalDialog.title }}</h3>
        </div>
        <div class="dialog-body">
          <p>{{ globalDialog.message }}</p>
        </div>
        <div class="dialog-footer">
          <van-button 
            v-if="globalDialog.showCancel"
            plain
            @click="closeGlobalDialog"
          >
            取消
          </van-button>
          <van-button 
            type="primary" 
            @click="confirmGlobalDialog"
          >
            确定
          </van-button>
        </div>
      </div>
    </van-popup>

    <!-- 全局Toast容器 -->
    <div id="toast-container"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast } from 'vant'

// 导入stores
import { useAuthStore } from '@/stores/modules/auth'
import { useAppStore } from '@/stores/modules/app'
import { useUserStore } from '@/stores/modules/user'

// 导入组合式函数
import { useAuth } from '@/composables/useAuth'
import { useDevice } from '@/composables/useDevice'
import { useNetworkStatus } from '@/composables/useNetworkStatus'

// 导入工具函数
import { setupErrorHandler } from '@/utils/error'
import { initializeApp } from '@/utils/app'

/**
 * 应用状态
 */
const isInitializing = ref(true)
const currentTheme = ref<'light' | 'dark'>('light')

/**
 * Store实例
 */
const authStore = useAuthStore()
const appStore = useAppStore()
const userStore = useUserStore()

/**
 * 组合式函数
 */
const { isLoggedIn, fetchUserInfo } = useAuth()
const { deviceInfo, isMobile } = useDevice()
const { isOnline } = useNetworkStatus()

/**
 * 路由实例
 */
const router = useRouter()
const route = useRoute()

/**
 * 缓存的组件列表
 */
const keepAliveComponents = computed(() => [
  'HomePage',
  'LoanProductList',
  'MachineSearch',
  'UserProfile'
])

/**
 * 页面转场动画名称
 */
const transitionName = computed(() => {
  // 根据路由层级和设备类型确定动画
  if (!isMobile.value) return 'fade'
  
  const toDepth = route.path.split('/').length
  const fromDepth = router.options.history.state.current?.split('/').length || 0
  
  return toDepth > fromDepth ? 'slide-left' : 'slide-right'
})

/**
 * 全局对话框状态
 */
interface GlobalDialog {
  visible: boolean
  title: string
  message: string
  showCancel: boolean
  onConfirm?: () => void
  onCancel?: () => void
}

const globalDialog = ref<GlobalDialog>({
  visible: false,
  title: '',
  message: '',
  showCancel: false
})

/**
 * 应用初始化
 */
const initializeApplication = async () => {
  try {
    // 1. 初始化应用配置
    await initializeApp()
    
    // 2. 设置错误处理器
    setupErrorHandler()
    
    // 3. 设置设备信息
    appStore.setDeviceInfo(deviceInfo.value)
    
    // 4. 检查并恢复用户会话
    if (authStore.token) {
      try {
        await fetchUserInfo()
        console.log('用户会话恢复成功')
      } catch (error) {
        console.warn('用户会话恢复失败:', error)
        authStore.logout()
      }
    }
    
    // 5. 设置主题
    const savedTheme = localStorage.getItem('app-theme') as 'light' | 'dark'
    if (savedTheme) {
      currentTheme.value = savedTheme
      applyTheme(savedTheme)
    }
    
    // 6. 初始化完成
    isInitializing.value = false
    
  } catch (error) {
    console.error('应用初始化失败:', error)
    showToast('应用初始化失败，请刷新重试')
    isInitializing.value = false
  }
}

/**
 * 应用主题设置
 */
const applyTheme = (theme: 'light' | 'dark') => {
  currentTheme.value = theme
  document.documentElement.setAttribute('data-theme', theme)
  localStorage.setItem('app-theme', theme)
}

/**
 * 全局对话框操作
 */
const showGlobalDialog = (options: Partial<GlobalDialog>) => {
  Object.assign(globalDialog.value, {
    visible: true,
    showCancel: false,
    ...options
  })
}

const closeGlobalDialog = () => {
  globalDialog.value.visible = false
  globalDialog.value.onCancel?.()
}

const confirmGlobalDialog = () => {
  globalDialog.value.visible = false
  globalDialog.value.onConfirm?.()
}

/**
 * 组件错误处理
 */
const handleComponentError = (error: Error, componentName: string) => {
  console.error(`组件 ${componentName} 发生错误:`, error)
  
  showGlobalDialog({
    title: '组件错误',
    message: '页面加载遇到问题，请刷新重试',
    onConfirm: () => {
      window.location.reload()
    }
  })
}

/**
 * 路由错误处理
 */
router.onError((error) => {
  console.error('路由错误:', error)
  showToast('页面跳转失败')
})

/**
 * 监听网络状态变化
 */
watch(isOnline, (online) => {
  if (online) {
    showToast('网络连接已恢复')
  } else {
    showToast('网络连接已断开')
  }
})

/**
 * 监听认证状态变化
 */
watch(isLoggedIn, (loggedIn) => {
  if (!loggedIn && route.meta.requiresAuth) {
    router.push('/auth/login')
  }
})

/**
 * 应用可见性变化处理
 */
const handleVisibilityChange = () => {
  if (!document.hidden && isLoggedIn.value) {
    // 应用重新可见时，检查token是否过期
    authStore.validateToken()
  }
}

/**
 * 生命周期钩子
 */
onMounted(async () => {
  // 初始化应用
  await initializeApplication()
  
  // 监听页面可见性变化
  document.addEventListener('visibilitychange', handleVisibilityChange)
  
  // 监听内存警告
  if ('memory' in performance) {
    const checkMemory = () => {
      const memory = (performance as any).memory
      if (memory.usedJSHeapSize / memory.jsHeapSizeLimit > 0.9) {
        console.warn('内存使用过高，建议刷新页面')
      }
    }
    setInterval(checkMemory, 60000) // 每分钟检查一次
  }
})

onBeforeUnmount(() => {
  document.removeEventListener('visibilitychange', handleVisibilityChange)
})

// 暴露给全局使用的方法
defineExpose({
  showGlobalDialog,
  closeGlobalDialog,
  applyTheme
})
</script>

<style lang="scss" scoped>
.app-container {
  min-height: 100vh;
  background-color: var(--background-color);
  color: var(--text-color);
  position: relative;
  
  &.app-loading {
    overflow: hidden;
  }
}

.app-loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(255, 255, 255, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  
  .loading-spinner {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
  }
}

.network-status {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
}

.app-content {
  position: relative;
  min-height: 100vh;
}

// 页面转场动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-left-enter-active,
.slide-left-leave-active,
.slide-right-enter-active,
.slide-right-leave-active {
  transition: transform 0.3s ease;
}

.slide-left-enter-from {
  transform: translateX(100%);
}

.slide-left-leave-to {
  transform: translateX(-100%);
}

.slide-right-enter-from {
  transform: translateX(-100%);
}

.slide-right-leave-to {
  transform: translateX(100%);
}

// 全局对话框样式
.global-dialog-content {
  padding: 24px;
  min-width: 280px;
  
  .dialog-header {
    margin-bottom: 16px;
    
    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      color: var(--text-primary);
    }
  }
  
  .dialog-body {
    margin-bottom: 24px;
    
    p {
      margin: 0;
      line-height: 1.5;
      color: var(--text-secondary);
    }
  }
  
  .dialog-footer {
    display: flex;
    gap: 12px;
    justify-content: flex-end;
  }
}

// 主题相关样式
.app-theme--light {
  --background-color: #fafafa;
  --text-color: #262626;
  --text-primary: #262626;
  --text-secondary: #595959;
}

.app-theme--dark {
  --background-color: #1a1a1a;
  --text-color: #ffffff;
  --text-primary: #ffffff;
  --text-secondary: #cccccc;
}

// 响应式设计
@media (max-width: 768px) {
  .global-dialog-content {
    padding: 20px;
    min-width: 260px;
  }
}
</style>
```

---

## 📚 依赖关系

### Store依赖
- `useAuthStore`: 认证状态管理
- `useAppStore`: 应用全局状态管理  
- `useUserStore`: 用户信息状态管理

### Composables依赖
- `useAuth`: 认证相关逻辑
- `useDevice`: 设备信息检测
- `useNetworkStatus`: 网络状态监听

### 工具函数依赖
- `setupErrorHandler`: 全局错误处理设置
- `initializeApp`: 应用初始化逻辑

---

## 🔧 核心功能

### 1. 应用初始化流程

```typescript
// 初始化步骤
1. 初始化应用配置
2. 设置错误处理器
3. 获取设备信息
4. 恢复用户会话
5. 应用主题设置
6. 标记初始化完成
```

### 2. 认证状态管理

- **会话恢复**: 应用启动时自动恢复用户登录状态
- **Token验证**: 定期验证Token有效性
- **自动登出**: Token过期时自动清除状态并跳转登录

### 3. 错误处理机制

- **全局错误捕获**: 捕获未处理的Promise拒绝和Vue错误
- **组件错误边界**: 处理子组件抛出的错误
- **网络错误处理**: 监听网络状态变化并提示用户

### 4. 性能优化

- **组件缓存**: 使用keep-alive缓存指定组件
- **懒加载**: 路由组件按需加载
- **内存监控**: 监控内存使用情况并预警

---

## 🎨 样式特性

### 主题系统
- **双主题支持**: 支持明亮和暗黑两种主题
- **动态切换**: 运行时切换主题并持久化存储
- **CSS变量**: 使用CSS自定义属性实现主题

### 动画效果
- **路由转场**: 根据导航方向使用不同的滑动动画
- **加载动画**: 应用启动时的加载指示器
- **状态过渡**: 各种状态变化的平滑过渡

### 响应式设计
- **移动优先**: 优先适配移动端设备
- **断点适配**: 在不同屏幕尺寸下的布局调整
- **触摸友好**: 移动端触摸交互优化

---

## 🚀 最佳实践

### 性能优化建议
1. **合理使用keep-alive**: 只缓存必要的组件
2. **控制bundle大小**: 按需导入第三方库
3. **图片懒加载**: 使用Intersection Observer API
4. **内存泄漏预防**: 及时清理事件监听器和定时器

### 错误处理建议
1. **错误边界**: 为关键组件设置错误边界
2. **用户友好提示**: 将技术错误转换为用户能理解的信息
3. **错误上报**: 在生产环境中上报错误到监控系统
4. **降级方案**: 为关键功能提供降级处理

### 代码维护建议
1. **模块化设计**: 将功能拆分为独立的模块
2. **类型安全**: 使用TypeScript确保类型安全
3. **单元测试**: 为核心逻辑编写单元测试
4. **文档更新**: 及时更新组件文档和API说明

---

## 🔍 调试指南

### 开发环境调试
```javascript
// 启用Vue DevTools
if (process.env.NODE_ENV === 'development') {
  window.__VUE_DEVTOOLS_GLOBAL_HOOK__.Vue = Vue
}

// 调试应用状态
console.log('Auth Store:', authStore.$state)
console.log('App Store:', appStore.$state)
console.log('Device Info:', deviceInfo.value)
```

### 错误排查
1. **检查网络状态**: 确认API请求是否正常
2. **验证Token状态**: 检查Token是否过期或无效
3. **查看控制台错误**: 关注Vue和JavaScript错误信息
4. **检查路由配置**: 确认路由守卫和权限设置

---

## 📝 更新记录

| 版本 | 日期 | 更新内容 |
|------|------|----------|
| 1.0.0 | 2024-01-15 | 初始版本，实现基础应用框架 |
| 1.1.0 | 2024-02-01 | 添加主题切换功能 |
| 1.2.0 | 2024-03-01 | 增强错误处理机制 |
| 1.3.0 | 2024-04-01 | 优化性能和内存管理 |

本文档将随着组件功能的迭代持续更新，确保文档与代码实现的一致性。 