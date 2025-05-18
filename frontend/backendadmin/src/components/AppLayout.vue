<template>
  <div class="app-layout">
    <div class="sidebar">
      <div class="logo-container">
        <div class="logo">
          <h2>数字惠农后台管理系统</h2>
        </div>
        <div class="logo-collapse" @click="toggleSidebar">
          <svg class="collapse-icon" viewBox="0 0 24 24">
            <path fill="currentColor" d="M15.41 16.59L10.83 12l4.58-4.59L14 6l-6 6 6 6 1.41-1.41z" />
          </svg>
        </div>
      </div>
      
      <div class="nav-container">
        <nav class="main-nav">
          <div 
            v-for="item in navItems" 
            :key="item.path"
            class="nav-item"
            :class="{ active: activePath === item.path }"
            @click="navigateTo(item.path)"
          >
            <div class="nav-icon">
              <svg class="icon" viewBox="0 0 24 24">
                <path fill="currentColor" :d="item.icon" />
              </svg>
            </div>
            <span class="nav-text">{{ item.title }}</span>
            <div class="arrow" v-if="item.children && item.children.length">
              <svg viewBox="0 0 24 24" class="arrow-icon" :class="{ expanded: item.expanded }">
                <path fill="currentColor" d="M7.41 8.59L12 13.17l4.59-4.58L18 10l-6 6-6-6 1.41-1.41z" />
              </svg>
            </div>
          </div>
        </nav>
      </div>
      
      <div class="sidebar-footer">
        <div class="user-info" @click="handleLogout">
          <div class="avatar">管</div>
          <span class="username">管理员</span>
          <svg class="logout-icon" viewBox="0 0 24 24">
            <path fill="currentColor" d="M17 7l-1.41 1.41L18.17 11H8v2h10.17l-2.58 2.58L17 17l5-5zM4 5h8V3H4c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h8v-2H4V5z" />
          </svg>
        </div>
      </div>
    </div>
    
    <div class="main-content">
      <div class="header">
        <div class="breadcrumb">
          <span>{{ currentPageTitle }}</span>
        </div>
        <div class="header-actions">
          <div class="user-notifications">
            <svg class="notification-icon" viewBox="0 0 24 24">
              <path fill="currentColor" d="M12 22c1.1 0 2-.9 2-2h-4c0 1.1.9 2 2 2zm6-6v-5c0-3.07-1.63-5.64-4.5-6.32V4c0-.83-.67-1.5-1.5-1.5s-1.5.67-1.5 1.5v.68C7.64 5.36 6 7.92 6 11v5l-2 2v1h16v-1l-2-2zm-2 1H8v-6c0-2.48 1.51-4.5 4-4.5s4 2.02 4 4.5v6z" />
            </svg>
          </div>
        </div>
      </div>
      
      <div class="content">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const sidebarCollapsed = ref(false)
const activePath = ref('/home')

const navItems = ref([
  {
    title: '首页',
    path: '/home',
    icon: 'M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z',
  },
  {
    title: '惠农贷',
    path: '/loan',
    icon: 'M11.8 10.9c-2.27-.59-3-1.2-3-2.15 0-1.09 1.01-1.85 2.7-1.85 1.78 0 2.44.85 2.5 2.1h2.21c-.07-1.72-1.12-3.3-3.21-3.81V3h-3v2.16c-1.94.42-3.5 1.68-3.5 3.61 0 2.31 1.91 3.46 4.7 4.13 2.5.6 3 1.48 3 2.41 0 .69-.49 1.79-2.7 1.79-2.06 0-2.87-.92-2.98-2.1h-2.2c.12 2.19 1.76 3.42 3.68 3.83V21h3v-2.15c1.95-.37 3.5-1.5 3.5-3.55 0-2.84-2.43-3.81-4.7-4.4z',
    expanded: false,
    children: [
      { title: '贷款申请', path: '/loan/apply' },
      { title: '贷款状态', path: '/loan/status' }
    ]
  },
  {
    title: '农机设备',
    path: '/machinery',
    icon: 'M15 11H9V7H7v10h2v-4h6v4h2V7h-2v4zm6-7h-4.18C16.4 1.84 15.3 1 14 1c-1.3 0-2.4.84-2.82 2H3v18h18V4zm-2 16H5V6h14v14zm-7-7c-.55 0-1-.45-1-1s.45-1 1-1 1 .45 1 1-.45 1-1 1z',
  },
  {
    title: '设备租赁',
    path: '/rent',
    icon: 'M19 3h-4.18C14.4 1.84 13.3 1 12 1c-1.3 0-2.4.84-2.82 2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-7 0c.55 0 1 .45 1 1s-.45 1-1 1-1-.45-1-1 .45-1 1-1zm2 14H7v-2h7v2zm3-4H7v-2h10v2zm0-4H7V7h10v2z',
  },
  {
    title: '租赁申请',
    path: '/rent/application',
    icon: 'M19 3h-4.18C14.4 1.84 13.3 1 12 1c-1.3 0-2.4.84-2.82 2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-7 0c.55 0 1 .45 1 1s-.45 1-1 1-1-.45-1-1 .45-1 1-1zm-2 14h7v-2H10v2zm10-4H7v-2h13v2zm0-4H7V7h13v2z',
  },
  {
    title: '订单管理',
    path: '/orders',
    icon: 'M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-7 9h-2V7h2v5zm0 4h-2v-2h2v2zm4-4h-2V7h2v5zm0 4h-2v-2h2v2zm-8-4H6V7h2v5zm0 4H6v-2h2v2z',
  },
  {
    title: '用户管理',
    path: '/users',
    icon: 'M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z',
  },
  {
    title: '系统设置',
    path: '/settings',
    icon: 'M19.14 12.94c.04-.3.06-.61.06-.94 0-.32-.02-.64-.07-.94l2.03-1.58c.18-.14.23-.41.12-.61l-1.92-3.32c-.12-.22-.37-.29-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54c-.04-.24-.24-.41-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.05.3-.09.63-.09.94s.02.64.07.94l-2.03 1.58c-.18.14-.23.41-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z',
  },
])

const currentPageTitle = computed(() => {
  const currentItem = navItems.value.find(item => route.path.startsWith(item.path))
  return currentItem ? currentItem.title : '首页'
})

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

const navigateTo = (path) => {
  activePath.value = path
  router.push(path)
}

const handleLogout = () => {
  localStorage.removeItem('token')
  router.push('/')
}
</script>

<style scoped>
.app-layout {
  display: flex;
  height: 100vh;
  width: 100%;
  overflow: hidden;
}

.sidebar {
  width: 260px;
  background-color: #001529;
  color: #fff;
  display: flex;
  flex-direction: column;
  transition: width 0.3s;
  overflow: hidden;
}

.logo-container {
  height: 64px;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo {
  display: flex;
  align-items: center;
}

.logo h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  white-space: nowrap;
  padding-left: 4px;
}

.logo-collapse {
  cursor: pointer;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
}

.logo-collapse:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.collapse-icon {
  width: 18px;
  height: 18px;
}

.nav-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px 0;
}

.main-nav {
  display: flex;
  flex-direction: column;
}

.nav-item {
  height: 48px;
  padding: 0 20px;
  display: flex;
  align-items: center;
  color: rgba(255, 255, 255, 0.65);
  cursor: pointer;
  transition: background-color 0.3s;
  position: relative;
}

.nav-item:hover {
  color: #fff;
  background-color: rgba(255, 255, 255, 0.08);
}

.nav-item.active {
  color: #fff;
  background-color: rgba(255, 255, 255, 0.08);
}

.nav-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background-color: #1890ff;
}

.nav-icon {
  margin-right: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon {
  width: 20px;
  height: 20px;
}

.nav-text {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.arrow {
  display: flex;
  align-items: center;
}

.arrow-icon {
  width: 16px;
  height: 16px;
  transition: transform 0.3s;
}

.arrow-icon.expanded {
  transform: rotate(180deg);
}

.sidebar-footer {
  padding: 18px 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 10px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: rgba(255, 255, 255, 0.08);
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: #1890ff;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 14px;
  font-size: 15px;
}

.username {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.logout-icon {
  width: 16px;
  height: 16px;
  opacity: 0.65;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: #f0f2f5;
}

.header {
  height: 64px;
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.breadcrumb {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.header-actions {
  display: flex;
  align-items: center;
}

.user-notifications {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  cursor: pointer;
  transition: background-color 0.3s;
}

.user-notifications:hover {
  background-color: rgba(0, 0, 0, 0.04);
}

.notification-icon {
  width: 22px;
  height: 22px;
  color: #666;
}

.content {
  flex: 1;
  padding: 24px;
  overflow: auto;
}
</style> 