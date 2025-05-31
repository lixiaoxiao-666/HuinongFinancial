<template>
  <div class="common-layout">
    <el-container>
      <!-- 顶部导航栏 -->
      <el-header height="70px" class="header">
        <div class="header-left">
          <div class="logo-container">
            <span class="logo-text">数字OA管理系统</span>
          </div>
          <el-tooltip content="折叠/展开菜单" placement="bottom">
            <div class="collapse-btn" @click="toggleCollapse">
              <el-icon>
                <Expand v-if="isCollapse" />
                <Fold v-else />
              </el-icon>
            </div>
          </el-tooltip>
        </div>
        
        <div class="header-center">
          <el-breadcrumb separator="/" class="custom-breadcrumb">
            <el-breadcrumb-item :to="{ path: '/dashboard' }">
              <el-icon><House /></el-icon>
              <span>首页</span>
            </el-breadcrumb-item>
            <el-breadcrumb-item v-if="breadcrumb">{{ breadcrumb }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="header-right">
          <!-- AI审批状态指示器 -->
          <div class="ai-status" v-if="hasPermission('system:manage')">
            <el-tooltip content="AI审批功能状态 - 始终运行" placement="bottom">
              <div class="status-indicator active">
                <div class="status-icon">
                  <el-icon><Cpu /></el-icon>
                </div>
                <div class="status-text">
                  <span class="status-label">AI审批</span>
                  <span class="status-value">运行中</span>
                </div>
                <div class="status-dot active"></div>
              </div>
            </el-tooltip>
          </div>
          
          <!-- 通知中心 -->
          <div class="notification-center">
            <el-tooltip content="通知中心" placement="bottom">
              <div class="notification-btn">
                <el-icon><Bell /></el-icon>
                <div class="notification-badge" v-if="notificationCount > 0">{{ notificationCount }}</div>
              </div>
            </el-tooltip>
          </div>
          
          <!-- 用户头像和下拉菜单 -->
          <el-dropdown @command="handleCommand" trigger="hover">
            <div class="user-dropdown">
              <el-avatar :size="36" :src="userAvatar" class="user-avatar">
                <el-icon><UserFilled /></el-icon>
              </el-avatar>
              <div class="user-info" v-if="!isCollapse">
                <span class="username">{{ currentUser?.display_name || currentUser?.username }}</span>
                <span class="user-role">{{ getRoleName(currentUser?.role) }}</span>
              </div>
              <el-icon class="dropdown-icon" v-if="!isCollapse"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu class="user-dropdown-menu">
                <div class="dropdown-header">
                  <el-avatar :size="50" :src="userAvatar" class="dropdown-avatar">
                    <el-icon><UserFilled /></el-icon>
                  </el-avatar>
                  <div class="dropdown-user-info">
                    <div class="dropdown-username">{{ currentUser?.display_name || currentUser?.username }}</div>
                    <div class="dropdown-user-role">{{ getRoleName(currentUser?.role) }}</div>
                  </div>
                </div>
                <el-divider style="margin: 12px 0;" />
                <el-dropdown-item command="profile" class="dropdown-item">
                  <el-icon><User /></el-icon>
                  <span>个人信息</span>
                </el-dropdown-item>
                <el-dropdown-item command="password" class="dropdown-item">
                  <el-icon><Key /></el-icon>
                  <span>修改密码</span>
                </el-dropdown-item>
                <el-dropdown-item command="settings" class="dropdown-item">
                  <el-icon><Setting /></el-icon>
                  <span>偏好设置</span>
                </el-dropdown-item>
                <el-dropdown-item divided command="logout" class="dropdown-item logout">
                  <el-icon><SwitchButton /></el-icon>
                  <span>退出登录</span>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-container>
        <!-- 侧边栏 -->
        <el-aside :width="isCollapse ? '89px' : '259px'" class="sidebar">
          <div class="sidebar-header" v-if="!isCollapse">
            <div class="sidebar-brand">
              <div class="brand-icon">
                <el-icon><Cpu /></el-icon>
              </div>
              <div class="brand-info">
                <div class="brand-title">惠农管理</div>
                <div class="brand-subtitle">智能审批系统</div>
              </div>
            </div>
          </div>
          
          <el-menu
            :default-active="activeMenu"
            :collapse="isCollapse"
            router
            background-color="transparent"
            text-color="rgba(255, 255, 255, 0.85)"
            active-text-color="#ffffff"
            class="sidebar-menu"
          >
            <el-menu-item index="/dashboard" class="menu-item">
              <div class="menu-icon">
                <el-icon><House /></el-icon>
              </div>
              <template #title>
                <span class="menu-title">工作台</span>
              </template>
            </el-menu-item>
            
            <el-menu-item index="/approval" v-if="hasPermission('approval:view')" class="menu-item">
              <div class="menu-icon">
                <el-icon><DocumentChecked /></el-icon>
              </div>
              <template #title>
                <span class="menu-title">审批看板</span>
              </template>
            </el-menu-item>
            
            <el-menu-item index="/smart-approval" v-if="hasPermission('approval:view')" class="menu-item">
              <div class="menu-icon">
                <el-icon><Cpu /></el-icon>
              </div>
              <template #title>
                <span class="menu-title">贷款审批</span>
                <el-tag size="small" class="menu-tag">AI</el-tag>
              </template>
            </el-menu-item>
            
            <el-menu-item index="/lease-approval" v-if="hasPermission('approval:view')" class="menu-item">
              <div class="menu-icon">
                <el-icon><Van /></el-icon>
              </div>
              <template #title>
                <span class="menu-title">租赁审批</span>
              </template>
            </el-menu-item>
            
            <el-menu-item index="/ai-workflow" v-if="hasPermission('approval:manage')" class="menu-item">
              <div class="menu-icon">
                <el-icon><Promotion /></el-icon>
              </div>
              <template #title>
                <span class="menu-title">AI审批流</span>
                <el-tag size="small" class="menu-tag new">NEW</el-tag>
              </template>
            </el-menu-item>
            
            <div class="menu-divider" v-if="!isCollapse"></div>
            
            <el-sub-menu index="management" v-if="hasPermission('user:manage')" class="sub-menu">
              <template #title>
                <div class="menu-icon">
                  <el-icon><Setting /></el-icon>
                </div>
                <span class="menu-title">系统管理</span>
              </template>
              <el-menu-item index="/users" class="sub-menu-item">
                <div class="sub-menu-icon">
                  <el-icon><User /></el-icon>
                </div>
                <template #title>
                  <span class="sub-menu-title">用户管理</span>
                </template>
              </el-menu-item>
              <el-menu-item index="/logs" class="sub-menu-item">
                <div class="sub-menu-icon">
                  <el-icon><Document /></el-icon>
                </div>
                <template #title>
                  <span class="sub-menu-title">操作日志</span>
                </template>
              </el-menu-item>
              <el-menu-item index="/system" class="sub-menu-item">
                <div class="sub-menu-icon">
                  <el-icon><Tools /></el-icon>
                </div>
                <template #title>
                  <span class="sub-menu-title">系统设置</span>
                </template>
              </el-menu-item>
            </el-sub-menu>
          </el-menu>
          
          <!-- 侧边栏底部信息 -->
          <div class="sidebar-footer" v-if="!isCollapse">
            <div class="footer-content">
              <div class="footer-version">
                <el-icon><InfoFilled /></el-icon>
                <span>版本 v1.0.0</span>
              </div>
            </div>
          </div>
        </el-aside>
        
        <!-- 主内容区 -->
        <el-main class="main-content">
          <div class="content-wrapper">
            <router-view />
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  House,
  DocumentChecked,
  Setting,
  User,
  Document,
  Tools,
  Expand,
  Fold,
  UserFilled,
  ArrowDown,
  Key,
  SwitchButton,
  Cpu,
  Promotion,
  Van,
  Bell,
  InfoFilled
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { getDashboard } from '@/api/admin'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const isCollapse = ref(false)
const aiApprovalEnabled = ref(true)
const userAvatar = ref('')
const notificationCount = ref(5)

// 计算属性
const currentUser = computed(() => authStore.user)
const hasPermission = computed(() => authStore.hasPermission)

const activeMenu = computed(() => {
  const path = route.path
  if (path.startsWith('/approval')) {
    return '/approval'
  }
  if (path === '/smart-approval') {
    return '/smart-approval'
  }
  if (path === '/ai-workflow') {
    return '/ai-workflow'
  }
  if (path === '/lease-approval') {
    return '/lease-approval'
  }
  return path
})

const breadcrumb = computed(() => {
  return route.meta.title as string
})

// 方法
const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const handleCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      ElMessage.info('个人信息功能开发中...')
      break
    case 'password':
      ElMessage.info('修改密码功能开发中...')
      break
    case 'settings':
      ElMessage.info('偏好设置功能开发中...')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        authStore.logout()
        router.push('/login')
        ElMessage.success('退出登录成功')
      } catch {
        // 用户取消
      }
      break
  }
}

const getRoleName = (role?: string) => {
  const roleMap: Record<string, string> = {
    'ADMIN': '管理员',
    '审批员': '审批员'
  }
  return roleMap[role || ''] || '用户'
}

const handleAIToggle = async (enabled: boolean) => {
  // AI功能始终保持运行状态，不允许切换
  aiApprovalEnabled.value = true
  ElMessage.info('AI审批功能始终保持运行状态')
}

// 获取系统状态
const fetchAIStatus = async () => {
  try {
    const data = await getDashboard()
    // AI状态始终设置为运行中
    aiApprovalEnabled.value = true
  } catch (error) {
    console.error('获取AI状态失败:', error)
    // 即使获取失败，也保持运行状态
    aiApprovalEnabled.value = true
  }
}

// 监听路由变化
watch(() => route.path, () => {
  // 可以在这里添加路由变化的逻辑
})

onMounted(() => {
  fetchAIStatus()
})
</script>

<style scoped>
.common-layout {
  height: 100vh;
  width: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  backdrop-filter: blur(10px);
  position: relative;
  z-index: 100;
}

.header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(90deg, transparent 0%, rgba(79, 172, 254, 0.05) 50%, transparent 100%);
  pointer-events: none;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 24px;
}

.logo-container {
  display: flex;
  align-items: center;
  padding: 0;
  border-radius: 0;
  background: none;
  box-shadow: none;
  transition: none;
  width: 259px;
  height: 48px;
  justify-content: flex-start;
  margin-left: -8px;
}

.logo-container:hover {
  transform: none;
  box-shadow: none;
}

.logo-text {
  font-size: 24px;
  font-weight: 700;
  color: #333;
  letter-spacing: 0.5px;
  text-shadow: none;
  white-space: nowrap;
}

.collapse-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  background: rgba(79, 172, 254, 0.1);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #4facfe;
  font-size: 18px;
}

.collapse-btn:hover {
  background: rgba(79, 172, 254, 0.2);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.2);
}

.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
  max-width: 500px;
  margin: 0 auto;
}

.custom-breadcrumb {
  background: rgba(79, 172, 254, 0.08);
  padding: 12px 20px;
  border-radius: 25px;
  border: 1px solid rgba(79, 172, 254, 0.15);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.ai-status {
  display: flex;
  align-items: center;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 14px;
  background: rgba(234, 84, 85, 0.1);
  border: 1px solid rgba(234, 84, 85, 0.2);
  border-radius: 12px;
  transition: all 0.3s ease;
  min-height: 44px;
}

.status-indicator.active {
  background: rgba(40, 199, 111, 0.1);
  border-color: rgba(40, 199, 111, 0.2);
}

.status-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: rgba(234, 84, 85, 0.2);
  border-radius: 8px;
  color: #ea5455;
  flex-shrink: 0;
}

.status-indicator.active .status-icon {
  background: rgba(40, 199, 111, 0.2);
  color: #28c76f;
}

.status-text {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.status-label {
  font-size: 16px;
  color: #666;
  font-weight: 500;
  line-height: 1.1;
}

.status-value {
  font-size: 18px;
  color: #ea5455;
  font-weight: 600;
  line-height: 1.1;
  margin-top: 1px;
}

.status-indicator.active .status-value {
  color: #28c76f;
}

.status-dot {
  width: 8px;
  height: 8px;
  background: #ea5455;
  border-radius: 50%;
  animation: pulse-red 2s infinite;
}

.status-dot.active {
  background: #28c76f;
  animation: pulse-green 2s infinite;
}

@keyframes pulse-red {
  0% { box-shadow: 0 0 0 0 rgba(234, 84, 85, 0.7); }
  70% { box-shadow: 0 0 0 10px rgba(234, 84, 85, 0); }
  100% { box-shadow: 0 0 0 0 rgba(234, 84, 85, 0); }
}

@keyframes pulse-green {
  0% { box-shadow: 0 0 0 0 rgba(40, 199, 111, 0.7); }
  70% { box-shadow: 0 0 0 10px rgba(40, 199, 111, 0); }
  100% { box-shadow: 0 0 0 0 rgba(40, 199, 111, 0); }
}

.notification-center {
  position: relative;
}

.notification-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  background: rgba(255, 159, 67, 0.1);
  border: 1px solid rgba(255, 159, 67, 0.15);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #ff9f43;
  position: relative;
}

.notification-btn:hover {
  background: rgba(255, 159, 67, 0.2);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255, 159, 67, 0.3);
}

.notification-badge {
  position: absolute;
  top: -6px;
  right: -6px;
  background: linear-gradient(135deg, #ff4757 0%, #ff3742 100%);
  color: white;
  border-radius: 12px;
  padding: 3px 7px;
  font-size: 10px;
  font-weight: 700;
  min-width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid white;
  box-shadow: 0 2px 6px rgba(255, 71, 87, 0.3);
}

.user-dropdown {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 12px;
  transition: all 0.3s ease;
  background: rgba(79, 172, 254, 0.08);
  border: 1px solid rgba(79, 172, 254, 0.12);
  min-height: 44px;
}

.user-dropdown:hover {
  background: rgba(79, 172, 254, 0.15);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.2);
}

.user-avatar {
  border: 2px solid rgba(79, 172, 254, 0.3);
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.user-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  min-width: 0;
  flex: 1;
}

.username {
  font-size: 19px;
  color: #333;
  font-weight: 600;
  line-height: 1.2;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100px;
}

.user-role {
  font-size: 16px;
  color: #666;
  line-height: 1.2;
  margin-top: 1px;
}

.dropdown-icon {
  font-size: 12px;
  color: #4facfe;
  transition: transform 0.3s ease;
  flex-shrink: 0;
}

.user-dropdown:hover .dropdown-icon {
  transform: rotate(180deg);
}

.sidebar {
  background: linear-gradient(180deg, #1e293b 0%, #0f172a 100%);
  transition: all 0.3s ease;
  overflow: hidden;
  box-shadow: 4px 0 20px rgba(0, 0, 0, 0.15);
  position: relative;
  margin-left: -8px;
}

.sidebar::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.08) 0%, rgba(147, 51, 234, 0.08) 100%);
  pointer-events: none;
}

.sidebar-header {
  padding: 24px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  border-radius: 10px;
  color: white;
  font-size: 18px;
}

.brand-info {
  display: flex;
  flex-direction: column;
}

.brand-title {
  font-size: 22px;
  font-weight: 700;
  color: white;
  line-height: 1.2;
}

.brand-subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.65);
  line-height: 1.2;
  margin-top: 2px;
}

.sidebar-menu {
  border-right: none;
  height: calc(100vh - 200px);
  padding: 10px;
  overflow-y: auto;
}

.menu-item {
  margin-bottom: 4px;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;
}

.menu-item:hover {
  background: rgba(79, 172, 254, 0.1) !important;
}

.menu-item.is-active {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%) !important;
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.3);
}

.menu-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  margin-right: 12px;
}

.menu-title {
  flex: 1;
  font-weight: 500;
}

.menu-tag {
  margin-left: auto;
  background: rgba(79, 172, 254, 0.2);
  color: #4facfe;
  border: none;
  font-size: 10px;
  font-weight: 600;
}

.menu-tag.new {
  background: rgba(255, 107, 107, 0.2);
  color: #ff6b6b;
}

.menu-divider {
  height: 1px;
  background: rgba(255, 255, 255, 0.1);
  margin: 16px 12px;
}

.sub-menu {
  border-radius: 12px;
  overflow: hidden;
}

.sub-menu-item {
  margin-left: 20px;
  border-radius: 8px;
  margin-bottom: 2px;
}

.sub-menu-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  margin-right: 12px;
}

.sub-menu-title {
  font-weight: 400;
  font-size: 13px;
}

.sidebar-footer {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.footer-content {
  display: flex;
  align-items: center;
  justify-content: center;
}

.footer-version {
  display: flex;
  align-items: center;
  gap: 8px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
}

.main-content {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 0;
  overflow: hidden;
  position: relative;
  margin-left: -8px;
}

.content-wrapper {
  height: 100%;
  padding: 24px;
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: rgba(79, 172, 254, 0.3) transparent;
  position: relative;
}

/* WebKit浏览器滚动条样式 */
.content-wrapper::-webkit-scrollbar {
  width: 6px;
}

.content-wrapper::-webkit-scrollbar-track {
  background: transparent;
}

.content-wrapper::-webkit-scrollbar-thumb {
  background: rgba(79, 172, 254, 0.3);
  border-radius: 3px;
  transition: background 0.3s ease;
}

.content-wrapper::-webkit-scrollbar-thumb:hover {
  background: rgba(79, 172, 254, 0.5);
}

/* Element Plus 组件样式覆写 */
:deep(.el-container) {
  height: 100%;
}

:deep(.el-header) {
  height: 70px !important;
  line-height: 70px;
}

:deep(.el-aside) {
  height: 100%;
  overflow: hidden;
}

:deep(.el-main) {
  height: calc(100vh - 70px);
  padding: 0;
}

:deep(.el-menu) {
  border-right: none;
  background: transparent !important;
}

:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  height: 60px;
  line-height: 60px;
  border-radius: 12px;
  margin-bottom: 4px;
  font-size: 21px;
}

:deep(.el-menu-item:hover),
:deep(.el-sub-menu__title:hover) {
  background-color: rgba(79, 172, 254, 0.1) !important;
}

:deep(.el-menu-item.is-active) {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%) !important;
  color: white !important;
}

:deep(.el-sub-menu .el-menu-item) {
  height: 50px;
  line-height: 50px;
  font-size: 19px;
}

:deep(.el-breadcrumb__inner) {
  color: #4facfe;
  font-weight: 500;
  font-size: 21px;
}

:deep(.el-breadcrumb__inner.is-link):hover {
  color: #00f2fe;
}

:deep(.user-dropdown-menu) {
  border-radius: 12px;
  border: none;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  padding: 8px;
  min-width: 200px;
}

:deep(.dropdown-item) {
  border-radius: 8px;
  margin: 2px 0;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  transition: all 0.3s ease;
  font-size: 21px;
}

:deep(.dropdown-item:hover) {
  background: rgba(79, 172, 254, 0.1);
}

:deep(.dropdown-item.logout:hover) {
  background: rgba(234, 84, 85, 0.1);
  color: #ea5455;
}

/* 新增下拉菜单头部样式 */
:deep(.dropdown-header) {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
  border-radius: 8px;
  margin-bottom: 8px;
}

:deep(.dropdown-avatar) {
  border: 2px solid rgba(79, 172, 254, 0.3);
}

:deep(.dropdown-user-info) {
  flex: 1;
  display: flex;
  flex-direction: column;
}

:deep(.dropdown-username) {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  line-height: 1.2;
}

:deep(.dropdown-user-role) {
  font-size: 18px;
  color: #666;
  line-height: 1.2;
  margin-top: 2px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header {
    padding: 0 16px;
  }
  
  .header-left {
    gap: 16px;
  }
  
  .logo-container {
    width: 219px;
    height: 40px;
    padding: 0;
    margin-left: -25px;
  }
  
  .logo-text {
    font-size: 18px;
  }
  
  .header-center {
    display: none;
  }
  
  .content-wrapper {
    padding: 16px;
  }
  
  .user-info {
    display: none !important;
  }
  
  .ai-status .status-text {
    display: none;
  }
  
  .notification-center {
    order: -1;
  }
  
  .sidebar {
    margin-left: -6px;
  }
  
  .main-content {
    margin-left: -6px;
  }
}

@media (max-width: 480px) {
  .header-right {
    gap: 12px;
  }
  
  .ai-status {
    display: none;
  }
  
  .logo-text {
    font-size: 16px;
  }
}
</style> 