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
      <SiderMenu :collapsed="siderCollapsed" />
    </a-layout-sider>

    <!-- 主要布局区域 -->
    <a-layout class="main-content-layout">
      <!-- 头部 -->
      <a-layout-header :style="headerStyle" class="layout-header">
        <HeaderContent 
          :collapsed="siderCollapsed" 
          @toggle-sider="toggleSider"
        />
      </a-layout-header>

      <!-- 内容区域 -->
      <a-layout-content :style="contentStyle" class="layout-content">
        <!-- 面包屑导航 -->
        <div class="breadcrumb-container">
          <a-breadcrumb>
            <a-breadcrumb-item>
              <HomeOutlined />
              <span>首页</span>
            </a-breadcrumb-item>
            <a-breadcrumb-item v-for="item in breadcrumbItems" :key="item.path">
              <router-link v-if="item.path" :to="item.path">
                {{ item.title }}
              </router-link>
              <span v-else>{{ item.title }}</span>
            </a-breadcrumb-item>
          </a-breadcrumb>
        </div>

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
import { HomeOutlined } from '@ant-design/icons-vue'
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
  overflow: 'auto',
  height: '100vh',
  position: 'fixed' as const,
  left: 0,
  top: 0,
  bottom: 0,
  background: '#001529',
  zIndex: 100
}))

const headerStyle = computed(() => ({
  background: '#fff',
  padding: 0,
  height: '64px',
  lineHeight: '64px',
  marginLeft: siderCollapsed.value ? '80px' : '240px',
  borderBottom: '1px solid #f0f0f0',
  position: 'fixed' as const,
  top: 0,
  right: 0,
  left: siderCollapsed.value ? '80px' : '240px',
  zIndex: 99,
  transition: 'all 0.2s'
}))

const contentStyle = computed(() => ({
  margin: `88px 24px 24px ${(siderCollapsed.value ? 80 : 240) + 24}px`,
  background: '#f0f2f5',
  minHeight: 'calc(100vh - 152px)',
  transition: 'all 0.2s'
}))

const footerStyle = computed(() => ({
  textAlign: 'center' as const,
  background: '#fff',
  borderTop: '1px solid #f0f0f0',
  marginLeft: siderCollapsed.value ? '80px' : '240px',
  transition: 'all 0.2s'
}))

/**
 * 面包屑数据
 */
const breadcrumbItems = computed(() => {
  const matched = route.matched.filter(item => item.meta && item.meta.title)
  return matched.map(item => ({
    title: item.meta?.title as string,
    path: item.path === route.path ? undefined : item.path
  }))
})

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
</script>

<style lang="scss" scoped>
.main-layout {
  min-height: 100vh;
  background: #f0f2f5;
}

// 侧边栏样式
.layout-sider {
  box-shadow: 2px 0 6px rgba(0, 21, 41, 0.35);
  
  :deep(.ant-layout-sider-children) {
    display: flex;
    flex-direction: column;
    height: 100%;
  }
}

.sider-logo {
  display: flex;
  align-items: center;
  padding: 16px 24px;
  background: rgba(255, 255, 255, 0.1);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.3s;
  min-height: 64px;
  
  &.collapsed {
    padding: 16px 12px;
    justify-content: center;
  }
  
  .logo-icon {
    width: 32px;
    height: 32px;
    border-radius: 4px;
    background: #1890ff;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 12px;
    flex-shrink: 0;
    
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
    }
    
    span {
      font-size: 12px;
      opacity: 0.8;
      line-height: 1;
    }
  }
  
  &.collapsed .logo-text {
    display: none;
  }
}

// 主内容区域
.main-content-layout {
  margin-left: 240px;
  transition: margin-left 0.2s;
  
  .layout-sider.ant-layout-sider-collapsed + & {
    margin-left: 80px;
  }
}

// 头部样式
.layout-header {
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

// 内容区域样式
.layout-content {
  .breadcrumb-container {
    background: #fff;
    padding: 16px 24px;
    border-radius: 6px;
    margin-bottom: 16px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
    
    :deep(.ant-breadcrumb) {
      .ant-breadcrumb-link {
        color: #666;
        
        &:hover {
          color: #1890ff;
        }
      }
      
      .ant-breadcrumb-separator {
        color: #999;
      }
    }
  }
  
  .page-content {
    background: #fff;
    border-radius: 6px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
    overflow: hidden;
  }
}

// 底部样式
.layout-footer {
  box-shadow: 0 -1px 4px rgba(0, 21, 41, 0.08);
}

// 页面切换动画
.page-transition-enter-active,
.page-transition-leave-active {
  transition: all 0.3s ease;
}

.page-transition-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.page-transition-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

// 响应式设计
@media (max-width: 768px) {
  .main-content-layout {
    margin-left: 0;
  }
  
  .layout-sider {
    :deep(.ant-layout-sider-children) {
      position: fixed;
      top: 0;
      left: 0;
      bottom: 0;
      z-index: 1000;
    }
  }
  
  .layout-header {
    margin-left: 0 !important;
    left: 0 !important;
  }
  
  .layout-content {
    margin-left: 24px !important;
  }
  
  .layout-footer {
    margin-left: 0 !important;
  }
}
</style> 