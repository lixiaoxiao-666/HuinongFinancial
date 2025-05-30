<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Spin as ASpin, message } from 'ant-design-vue'

// 导入stores
import { useAuthStore, useAppStore } from '@/stores'

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

/**
 * 路由实例
 */
const router = useRouter()
const route = useRoute()

/**
 * 页面转场动画名称
 */
const transitionName = computed(() => {
  // OA系统使用简单的淡入淡出动画
  return 'fade'
})

/**
 * 应用初始化
 */
const initializeApplication = async () => {
  try {
    console.log('🚀 惠农OA系统初始化开始...')

    // 1. 初始化应用配置
    appStore.initializeApp()
    
    // 2. 恢复主题设置
    const savedTheme = localStorage.getItem('oa_theme') as 'light' | 'dark'
    if (savedTheme) {
      currentTheme.value = savedTheme
      applyTheme(savedTheme)
    }

    // 3. 初始化认证状态
    await authStore.initializeAuth()
    
    // 4. 初始化完成
    isInitializing.value = false
    
    console.log('✅ 惠农OA系统初始化完成')
    
  } catch (error) {
    console.error('❌ 惠农OA系统初始化失败:', error)
    message.error('系统初始化失败，请刷新重试')
    isInitializing.value = false
  }
}

/**
 * 应用主题设置
 */
const applyTheme = (theme: 'light' | 'dark') => {
  currentTheme.value = theme
  document.documentElement.setAttribute('data-theme', theme)
  appStore.applyTheme(theme)
}

/**
 * 组件错误处理
 */
const handleComponentError = (error: Error, componentName: string) => {
  console.error(`组件 ${componentName} 发生错误:`, error)
  
  message.error({
    content: '页面加载遇到问题，请刷新重试',
    duration: 5,
    key: 'component-error'
  })
}

/**
 * 路由错误处理
 */
router.onError((error) => {
  console.error('路由错误:', error)
  message.error('页面跳转失败，请重试')
})

/**
 * 应用可见性变化处理
 */
const handleVisibilityChange = () => {
  if (!document.hidden && authStore.isAuthenticated) {
    // 应用重新可见时，验证token状态
    authStore.validateToken().catch(() => {
      // Token验证失败的处理在store中已经完成
    })
  }
}

/**
 * 窗口关闭前处理
 */
const handleBeforeUnload = (event: BeforeUnloadEvent) => {
  // 如果有未保存的更改，提示用户
  const hasUnsavedChanges = false // 暂时硬编码，后续可以从store获取
  if (hasUnsavedChanges) {
    event.preventDefault()
    event.returnValue = '您有未保存的更改，确定要离开吗？'
    return event.returnValue
  }
}

/**
 * 键盘快捷键处理
 */
const handleKeydown = (event: KeyboardEvent) => {
  // Ctrl/Cmd + K 打开全局搜索
  if ((event.ctrlKey || event.metaKey) && event.key === 'k') {
    event.preventDefault()
    appStore.toggleGlobalSearch()
  }
  
  // ESC 关闭模态框或搜索
  if (event.key === 'Escape') {
    if (appStore.globalSearchVisible) {
      appStore.toggleGlobalSearch()
    }
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
  
  // 监听窗口关闭前事件
  window.addEventListener('beforeunload', handleBeforeUnload)
  
  // 监听键盘事件
  document.addEventListener('keydown', handleKeydown)
  
  // 设置全局错误处理
  window.addEventListener('unhandledrejection', (event) => {
    console.error('未处理的Promise拒绝:', event.reason)
    if (import.meta.env.PROD) {
      // 生产环境可以上报到监控服务
      // Sentry.captureException(event.reason)
    }
  })
})

onBeforeUnmount(() => {
  // 清理事件监听器
  document.removeEventListener('visibilitychange', handleVisibilityChange)
  window.removeEventListener('beforeunload', handleBeforeUnload)
  document.removeEventListener('keydown', handleKeydown)
})

// 暴露给全局使用的方法
defineExpose({
  applyTheme,
  handleComponentError
})
</script>

<template>
  <div 
    id="app" 
    :class="[
      'oa-app',
      `oa-theme--${currentTheme}`,
      { 'oa-initializing': isInitializing }
    ]"
  >
    <!-- 全局加载遮罩 -->
    <div v-if="isInitializing" class="oa-loading-overlay">
      <div class="loading-content">
        <a-spin size="large" />
        <p class="loading-text">惠农OA系统初始化中...</p>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div v-else class="oa-content">
      <router-view v-slot="{ Component, route }">
        <transition :name="transitionName" mode="out-in">
          <component 
            :is="Component" 
            :key="route.fullPath"
            @error="handleComponentError"
          />
        </transition>
      </router-view>
    </div>

    <!-- 全局Modal容器 -->
    <div id="modal-container"></div>
  </div>
</template>

<style lang="scss" scoped>
.oa-app {
  height: 100vh;
  overflow: hidden;
  background-color: $background-color;
  color: $text-color;
  
  &.oa-initializing {
    overflow: hidden;
  }
}

.oa-loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, 
    rgba(0, 21, 41, 0.95) 0%, 
    rgba(24, 144, 255, 0.95) 100%
  );
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  
  .loading-content {
    text-align: center;
    color: white;
    
    .loading-text {
      margin-top: 16px;
      font-size: 16px;
      font-weight: 500;
    }
  }
}

.oa-content {
  height: 100vh;
  overflow: hidden;
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

// 主题相关样式
.oa-theme--light {
  --app-bg-color: #{$background-color};
  --app-text-color: #{$text-color};
}

.oa-theme--dark {
  --app-bg-color: #{$dark-bg-color};
  --app-text-color: #{$dark-text-color};
  
  background-color: var(--app-bg-color);
  color: var(--app-text-color);
}

// 响应式设计
@include responsive(xs) {
  .oa-loading-overlay {
    .loading-content .loading-text {
      font-size: 14px;
    }
  }
}
</style>

<style lang="scss">
// 全局样式（不使用scoped）

// 调整 Ant Design 组件在暗色主题下的样式
[data-theme="dark"] {
  // 这里可以添加暗色主题下的全局样式调整
  
  .ant-layout {
    background: $dark-bg-color;
  }
  
  .ant-layout-header {
    background: $dark-component-bg;
    border-bottom-color: $dark-border-color;
  }
  
  .ant-layout-sider {
    background: $dark-component-bg;
  }
  
  .ant-menu {
    background: transparent;
    color: $dark-text-color;
    
    .ant-menu-item {
      color: $dark-text-color-secondary;
      
      &:hover {
        color: $dark-text-color;
      }
      
      &.ant-menu-item-selected {
        color: $primary-color;
        background-color: rgba($primary-color, 0.1);
      }
    }
  }
}

// 全局滚动条优化
* {
  scrollbar-width: thin;
  scrollbar-color: rgba(0, 0, 0, 0.2) transparent;
}
</style>
