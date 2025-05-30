<template>
  <a-layout class="main-layout">
    <!-- 侧边栏 -->
    <a-layout-sider
      v-model:collapsed="siderCollapsed"
      :trigger="null"
      collapsible
      :width="240"
      :collapsed-width="80"
      :style="siderStyle"
      class="layout-sider"
    >
      <!-- Logo区域 -->
      <div class="sider-logo" :class="{ collapsed: siderCollapsed }">
        <div class="logo-icon">
          <img src="/logo.svg" alt="Logo" class="logo-img" />
        </div>
        <div v-if="!siderCollapsed" class="logo-text">
          <h3>惠农OA</h3>
          <span>管理系统</span>
        </div>
      </div>

      <!-- 导航菜单 -->
      <div class="sider-menu-wrapper" ref="siderMenuRef">
        <SiderMenu :collapsed="siderCollapsed" />
      </div>
    </a-layout-sider>
    
    <!-- 右侧布局 -->
    <a-layout>
      <!-- 头部 -->
      <a-layout-header :style="headerStyle" class="layout-header">
        <HeaderContent 
          :collapsed="siderCollapsed" 
          @toggle-sider="toggleSider"
        />
      </a-layout-header>
      
      <!-- 内容区域 -->
      <a-layout-content :style="contentStyle" class="layout-content">
        <!-- 页面内容 -->
        <div class="page-content">
          <router-view v-slot="{ Component }">
            <transition name="page-transition" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </a-layout-content>
      
      <!-- 底部 -->
      <a-layout-footer :style="footerStyle" class="layout-footer">
        <FooterContent />
      </a-layout-footer>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import SiderMenu from './SiderMenu.vue'
import HeaderContent from './HeaderContent.vue'
import FooterContent from './FooterContent.vue'

/**
 * 组件状态
 */
const siderCollapsed = ref(false)
const route = useRoute()

/**
 * 样式定义
 */
const siderStyle = computed(() => ({
  overflow: 'hidden',
  height: '100vh', // 修改为占满整个视口高度
  background: '#001529',
  boxShadow: '2px 0 8px 0 rgba(29, 35, 41, 0.05)',
  position: 'fixed' as const,
  left: 0,
  top: 0,
  zIndex: 100
}))

const headerStyle = computed(() => ({
  background: '#fff',
  padding: 0,
  height: '64px',
  lineHeight: '64px',
  borderBottom: '1px solid #f0f0f0',
  boxShadow: '0 1px 4px 0 rgba(0, 21, 41, 0.08)',
  position: 'relative',
  zIndex: 99
}))

const contentStyle = computed(() => ({
  background: '#f0f2f5',
  minHeight: 'calc(100vh - 128px)', // 减去header和footer的高度
  borderRadius: '0',
  overflow: 'auto', // 添加滚动
  maxHeight: 'calc(100vh - 128px)' // 限制最大高度
}))

const footerStyle = computed(() => ({
  textAlign: 'center' as const,
  background: '#fff',
  borderTop: '1px solid #f0f0f0',
  padding: '12px 24px',
  height: '64px'
}))

/**
 * 切换侧边栏
 */
const toggleSider = () => {
  siderCollapsed.value = !siderCollapsed.value
}

/**
 * 监听路由变化，可以在这里添加页面切换逻辑
 */
watch(route, (newRoute) => {
  // 路由变化时的处理逻辑
  console.log('Route changed to:', newRoute.path)
})

/**
 * 监听侧边栏折叠状态变化，调整右侧布局
 */
watch(siderCollapsed, (collapsed) => {
  // 通过CSS变量来控制右侧布局的margin-left
  document.documentElement.style.setProperty(
    '--sider-width',
    collapsed ? '80px' : '240px'
  )
})
</script>

<style lang="scss" scoped>
.main-layout {
  min-height: 100vh;
  background: #f0f2f5;
  
  // 右侧布局需要留出左侧边栏的位置
  > .ant-layout {
    margin-left: var(--sider-width, 240px);
    transition: margin-left 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }
}

// 头部样式
.layout-header {
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.95);
}

// 侧边栏样式
.layout-sider {
  border-right: 1px solid rgba(255, 255, 255, 0.05);
  box-shadow: 2px 0 8px 0 rgba(29, 35, 41, 0.05);
  
  :deep(.ant-layout-sider-children) {
    display: flex;
    flex-direction: column;
    height: 100%;
    position: relative;
  }
}

.sider-logo {
  display: flex;
  align-items: center;
  padding: 16px 24px;
  background: rgba(255, 255, 255, 0.05);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  min-height: 64px;
  backdrop-filter: blur(10px);
  flex-shrink: 0; // 防止被压缩
  
  &.collapsed {
    padding: 16px 12px;
    justify-content: center;
  }
  
  .logo-icon {
    width: 32px;
    height: 32px;
    border-radius: 6px;
    background: linear-gradient(135deg, #1890ff 0%, #40a9ff 100%);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 12px;
    flex-shrink: 0;
    box-shadow: 0 2px 8px rgba(24, 144, 255, 0.3);
    
    .logo-img {
      width: 20px;
      height: 20px;
      filter: brightness(0) invert(1);
    }
  }
  
  .logo-text {
    color: white;
    
    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      line-height: 1.2;
      background: linear-gradient(135deg, #fff 0%, #f0f0f0 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }
    
    span {
      font-size: 12px;
      opacity: 0.8;
      line-height: 1;
      color: rgba(255, 255, 255, 0.7);
    }
  }
  
  &.collapsed .logo-text {
    display: none;
  }
}

.sider-menu-wrapper {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding-right: 2px;
  
  // 自定义滚动条样式 - 增强可见性
  &::-webkit-scrollbar {
    width: 8px;
  }
  
  &::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 4px;
    margin: 8px 0;
  }
  
  &::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.4);
    border-radius: 4px;
    transition: all 0.3s ease;
    border: 1px solid rgba(255, 255, 255, 0.1);
    
    &:hover {
      background: rgba(255, 255, 255, 0.6);
      transform: scaleX(1.2);
    }
    
    &:active {
      background: rgba(255, 255, 255, 0.8);
    }
  }
  
  // Firefox滚动条样式
  scrollbar-width: thin;
  scrollbar-color: rgba(255, 255, 255, 0.4) rgba(255, 255, 255, 0.05);
  
  // 滚动条始终显示
  &:hover::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.6);
  }
}

// 内容区域样式
.layout-content {
  .page-content {
    background: transparent;
    min-height: calc(100vh - 128px);
    overflow-y: auto;
    overflow-x: hidden;
    padding: 16px;
    
    // 自定义滚动条样式
    &::-webkit-scrollbar {
      width: 8px;
    }
    
    &::-webkit-scrollbar-track {
      background: rgba(0, 0, 0, 0.04);
      border-radius: 4px;
    }
    
    &::-webkit-scrollbar-thumb {
      background: rgba(0, 0, 0, 0.15);
      border-radius: 4px;
      transition: background 0.3s ease;
      
      &:hover {
        background: rgba(0, 0, 0, 0.25);
      }
    }
    
    // Firefox滚动条样式
    scrollbar-width: thin;
    scrollbar-color: rgba(0, 0, 0, 0.15) rgba(0, 0, 0, 0.04);
  }
}

// 页面转场动画
.page-transition-enter-active,
.page-transition-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.page-transition-enter-from {
  opacity: 0;
  transform: translateY(16px);
}

.page-transition-leave-to {
  opacity: 0;
  transform: translateY(-16px);
}

// 底部样式
.layout-footer {
  font-size: 14px;
  color: rgba(0, 0, 0, 0.65);
  
  :deep(.footer-content) {
    display: flex;
    justify-content: center;
    align-items: center;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .main-layout {
    > .ant-layout {
      margin-left: 0;
    }
  }
  
  .layout-content {
    margin: 8px;
  }
  
  .layout-sider {
    position: absolute;
    top: 0;
    left: 0;
    height: 100%;
    z-index: 1000;
    transform: translateX(-100%);
    transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    
    &:not(.ant-layout-sider-collapsed) {
      transform: translateX(0);
    }
  }
  
  .sider-logo {
    min-height: 56px;
    padding: 12px 16px;
  }
}

// 暗色主题适配
:deep([data-theme="dark"]) {
  .layout-content {
    background: #141414;
  }
  
  .layout-header {
    background: rgba(31, 31, 31, 0.95);
    border-bottom-color: #303030;
  }
  
  .layout-footer {
    background: #1f1f1f;
    border-top-color: #303030;
    color: rgba(255, 255, 255, 0.65);
  }
}
</style>

<style lang="scss">
// 全局滚动条优化
* {
  scrollbar-width: thin;
  scrollbar-color: rgba(0, 0, 0, 0.2) transparent;
}

*::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

*::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.04);
  border-radius: 4px;
}

*::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.15);
  border-radius: 4px;
  transition: background 0.3s ease;
}

*::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.25);
}

*::-webkit-scrollbar-corner {
  background: transparent;
}

// 优化卡片阴影效果
.ant-card {
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.03), 
              0 1px 6px -1px rgba(0, 0, 0, 0.02), 
              0 2px 4px 0 rgba(0, 0, 0, 0.02) !important;
  border: 1px solid rgba(0, 0, 0, 0.06) !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  
  &:hover {
    box-shadow: 0 4px 12px 0 rgba(0, 0, 0, 0.05), 
                0 2px 16px -1px rgba(0, 0, 0, 0.04), 
                0 4px 8px 0 rgba(0, 0, 0, 0.04) !important;
    transform: translateY(-2px);
  }
}

@keyframes pulse {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
  100% {
    opacity: 1;
  }
}
</style> 