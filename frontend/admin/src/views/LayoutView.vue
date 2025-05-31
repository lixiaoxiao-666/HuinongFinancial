<template>
  <div class="common-layout">
    <el-container>
      <!-- 顶部导航栏 -->
      <el-header height="60px" class="header">
        <div class="header-left">
          <div class="logo-container">
            <img src="@/assets/icons/blue-oa-logo.png" alt="Logo" class="logo" />
            <span class="logo-text">数字惠农OA</span>
          </div>
          <el-icon @click="toggleCollapse" class="collapse-btn">
            <Expand v-if="isCollapse" />
            <Fold v-else />
          </el-icon>
        </div>
        
        <div class="header-center">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="breadcrumb">{{ breadcrumb }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="header-right">
          <!-- AI审批状态指示器 -->
          <div class="ai-status" v-if="hasPermission('system:manage')">
            <el-tooltip content="AI审批功能状态" placement="bottom">
              <el-tag
                :type="aiApprovalEnabled ? 'success' : 'danger'"
                size="small"
                effect="dark"
              >
                <el-icon><Cpu /></el-icon>
                AI审批{{ aiApprovalEnabled ? '开启' : '关闭' }}
              </el-tag>
            </el-tooltip>
          </div>
          
          <!-- 用户头像和下拉菜单 -->
          <el-dropdown @command="handleCommand">
            <span class="user-dropdown">
              <el-avatar :size="32" :src="userAvatar">
                <el-icon><UserFilled /></el-icon>
              </el-avatar>
              <span class="username">{{ currentUser?.display_name || currentUser?.username }}</span>
              <el-icon class="dropdown-icon"><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>个人信息
                </el-dropdown-item>
                <el-dropdown-item command="password">
                  <el-icon><Key /></el-icon>修改密码
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-container>
        <!-- 侧边栏 -->
        <el-aside :width="isCollapse ? '64px' : '240px'" class="sidebar">
          <el-menu
            :default-active="activeMenu"
            :collapse="isCollapse"
            router
            background-color="#304156"
            text-color="#bfcbd9"
            active-text-color="#409eff"
            class="sidebar-menu"
          >
            <el-menu-item index="/dashboard">
              <el-icon><House /></el-icon>
              <template #title>工作台</template>
            </el-menu-item>
            
            <el-menu-item index="/approval" v-if="hasPermission('approval:view')">
              <el-icon><DocumentChecked /></el-icon>
              <template #title>审批看板</template>
            </el-menu-item>
            
            <el-sub-menu index="management" v-if="hasPermission('user:manage')">
              <template #title>
                <el-icon><Setting /></el-icon>
                <span>系统管理</span>
              </template>
              <el-menu-item index="/users">
                <el-icon><User /></el-icon>
                <template #title>用户管理</template>
              </el-menu-item>
              <el-menu-item index="/logs">
                <el-icon><Document /></el-icon>
                <template #title>操作日志</template>
              </el-menu-item>
              <el-menu-item index="/system">
                <el-icon><Tools /></el-icon>
                <template #title>系统设置</template>
              </el-menu-item>
            </el-sub-menu>
          </el-menu>
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
  Cpu
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { getDashboard } from '@/api/admin'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const isCollapse = ref(false)
const aiApprovalEnabled = ref(true)
const userAvatar = ref('')

// 计算属性
const currentUser = computed(() => authStore.user)
const hasPermission = computed(() => authStore.hasPermission)

const activeMenu = computed(() => {
  const path = route.path
  if (path.startsWith('/approval')) {
    return '/approval'
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

// 获取系统状态
const fetchAIStatus = async () => {
  try {
    const data = await getDashboard()
    aiApprovalEnabled.value = data.ai_enabled
  } catch (error) {
    console.error('获取AI状态失败:', error)
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
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 12px;
  transition: all 0.3s ease;
}

.logo-container:hover .logo {
  transform: rotate(10deg);
}

.logo {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  transition: all 0.3s ease;
  filter: drop-shadow(0 2px 6px rgba(0, 0, 0, 0.15));
  border: 2px solid #e6e6e6;
  padding: 2px;
  background-color: #ffffff;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  color: #1976D2;
  white-space: nowrap;
  letter-spacing: 0.5px;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.collapse-btn {
  font-size: 18px;
  cursor: pointer;
  transition: color 0.3s;
  color: #666;
}

.collapse-btn:hover {
  color: #409eff;
}

.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
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

.user-dropdown {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 6px;
  transition: background-color 0.3s;
}

.user-dropdown:hover {
  background-color: #f5f7fa;
}

.username {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.dropdown-icon {
  font-size: 12px;
  color: #666;
}

.sidebar {
  background-color: #304156;
  transition: width 0.3s;
  overflow: hidden;
}

.sidebar-menu {
  border-right: none;
  height: 100%;
}

.main-content {
  background-color: #f4f5f7;
  padding: 0;
  overflow: hidden;
}

.content-wrapper {
  height: 100%;
  padding: 20px;
  overflow-y: auto;
}

/* Element Plus 组件样式覆写 */
:deep(.el-container) {
  height: 100%;
}

:deep(.el-header) {
  height: 60px !important;
  line-height: 60px;
}

:deep(.el-aside) {
  height: 100%;
  overflow: hidden;
}

:deep(.el-main) {
  height: calc(100vh - 60px);
  padding: 0;
}

:deep(.el-menu) {
  border-right: none;
}

:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  height: 48px;
  line-height: 48px;
}

:deep(.el-breadcrumb__inner) {
  color: #666;
}

:deep(.el-breadcrumb__inner.is-link):hover {
  color: #409eff;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header {
    padding: 0 16px;
  }
  
  .header-left {
    gap: 12px;
  }
  
  .logo-text {
    display: none;
  }
  
  .header-center {
    display: none;
  }
  
  .content-wrapper {
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .header-right {
    gap: 12px;
  }
  
  .ai-status {
    display: none;
  }
  
  .username {
    display: none;
  }
}
</style> 