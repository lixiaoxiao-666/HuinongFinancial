<template>
  <div class="header-content">
    <!-- 左侧区域 -->
    <div class="header-left">
      <!-- 侧边栏折叠按钮 -->
      <a-button
        type="text"
        class="collapse-btn"
        @click="$emit('toggle-sider')"
      >
        <template #icon>
          <MenuUnfoldOutlined v-if="collapsed" />
          <MenuFoldOutlined v-else />
        </template>
      </a-button>

      <!-- 全局搜索 -->
      <div class="global-search">
        <a-input-search
          v-model:value="searchText"
          placeholder="搜索用户、申请、订单..."
          style="width: 300px"
          @search="handleSearch"
        >
          <template #prefix>
            <SearchOutlined />
          </template>
        </a-input-search>
      </div>
    </div>

    <!-- 右侧区域 -->
    <div class="header-right">
      <!-- 快捷操作 -->
      <div class="quick-actions">
        <!-- 刷新页面 -->
        <a-tooltip title="刷新页面">
          <a-button type="text" class="action-btn" @click="handleRefresh">
            <template #icon>
              <ReloadOutlined :class="{ spinning: isRefreshing }" />
            </template>
          </a-button>
        </a-tooltip>

        <!-- 全屏切换 -->
        <a-tooltip :title="isFullscreen ? '退出全屏' : '全屏'">
          <a-button type="text" class="action-btn" @click="toggleFullscreen">
            <template #icon>
              <FullscreenExitOutlined v-if="isFullscreen" />
              <FullscreenOutlined v-else />
            </template>
          </a-button>
        </a-tooltip>
      </div>

      <!-- 通知消息 -->
      <div class="notifications">
        <a-badge :count="unreadCount" :offset="[10, 0]">
          <a-dropdown :trigger="['click']" placement="bottomRight">
            <a-button type="text" class="action-btn">
              <template #icon>
                <BellOutlined />
              </template>
            </a-button>
            <template #overlay>
              <div class="notification-dropdown">
                <div class="notification-header">
                  <span>通知消息</span>
                  <a @click="markAllAsRead">全部已读</a>
                </div>
                <div class="notification-list">
                  <div
                    v-for="notification in notifications"
                    :key="notification.id"
                    class="notification-item"
                    :class="{ unread: !notification.read }"
                    @click="handleNotificationClick(notification)"
                  >
                    <div class="notification-content">
                      <div class="notification-title">{{ notification.title }}</div>
                      <div class="notification-desc">{{ notification.content }}</div>
                      <div class="notification-time">{{ formatTime(notification.created_at) }}</div>
                    </div>
                  </div>
                  <div v-if="notifications.length === 0" class="notification-empty">
                    暂无通知消息
                  </div>
                </div>
                <div class="notification-footer">
                  <a @click="viewAllNotifications">查看全部</a>
                </div>
              </div>
            </template>
          </a-dropdown>
        </a-badge>
      </div>

      <!-- 用户菜单 -->
      <div class="user-menu">
        <a-dropdown :trigger="['click']" placement="bottomRight">
          <div class="user-info">
            <a-avatar :src="userInfo?.avatar" :size="32">
              <template #icon>
                <UserOutlined />
              </template>
            </a-avatar>
            <div class="user-details">
              <div class="user-name">{{ userInfo?.real_name || userInfo?.username }}</div>
              <div class="user-role">{{ userInfo?.role_name }}</div>
            </div>
            <DownOutlined class="dropdown-icon" />
          </div>
          <template #overlay>
            <a-menu class="user-dropdown-menu">
              <a-menu-item key="profile" @click="viewProfile">
                <UserOutlined />
                <span>个人资料</span>
              </a-menu-item>
              <a-menu-item key="settings" @click="openSettings">
                <SettingOutlined />
                <span>个人设置</span>
              </a-menu-item>
              <a-menu-divider />
              <a-menu-item key="help" @click="openHelp">
                <QuestionCircleOutlined />
                <span>帮助中心</span>
              </a-menu-item>
              <a-menu-divider />
              <a-menu-item key="logout" @click="handleLogout">
                <LogoutOutlined />
                <span>退出登录</span>
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/modules/auth'
import { message, Modal } from 'ant-design-vue'
import {
  MenuUnfoldOutlined,
  MenuFoldOutlined,
  SearchOutlined,
  ReloadOutlined,
  FullscreenOutlined,
  FullscreenExitOutlined,
  BellOutlined,
  UserOutlined,
  DownOutlined,
  SettingOutlined,
  QuestionCircleOutlined,
  LogoutOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'

interface Props {
  collapsed: boolean
}

interface Notification {
  id: number
  title: string
  content: string
  type: 'info' | 'warning' | 'error' | 'success'
  read: boolean
  created_at: string
}

const props = defineProps<Props>()

/**
 * 组件状态
 */
const router = useRouter()
const authStore = useAuthStore()

const searchText = ref('')
const isRefreshing = ref(false)
const isFullscreen = ref(false)

// 模拟通知数据
const notifications = ref<Notification[]>([
  {
    id: 1,
    title: '新的贷款申请',
    content: '用户李明提交了农业创业贷申请，请及时审核',
    type: 'info',
    read: false,
    created_at: '2024-01-15T10:30:00Z'
  },
  {
    id: 2,
    title: '系统维护通知',
    content: '系统将于今晚23:00进行维护升级，预计持续2小时',
    type: 'warning',
    read: false,
    created_at: '2024-01-15T09:00:00Z'
  },
  {
    id: 3,
    title: 'AI审批异常',
    content: 'Dify工作流调用失败，请检查配置',
    type: 'error',
    read: true,
    created_at: '2024-01-15T08:15:00Z'
  }
])

/**
 * 计算属性
 */
const userInfo = computed(() => authStore.userInfo)
const unreadCount = computed(() => notifications.value.filter(n => !n.read).length)

/**
 * 搜索处理
 */
const handleSearch = (value: string) => {
  if (!value.trim()) {
    message.warning('请输入搜索关键词')
    return
  }
  
  // 这里可以实现全局搜索逻辑
  console.log('Global search:', value)
  message.info(`搜索: ${value}`)
}

/**
 * 刷新页面
 */
const handleRefresh = () => {
  isRefreshing.value = true
  setTimeout(() => {
    location.reload()
  }, 500)
}

/**
 * 全屏切换
 */
const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
    isFullscreen.value = true
  } else {
    document.exitFullscreen()
    isFullscreen.value = false
  }
}

/**
 * 监听全屏状态变化
 */
const handleFullscreenChange = () => {
  isFullscreen.value = !!document.fullscreenElement
}

/**
 * 格式化时间
 */
const formatTime = (time: string) => {
  return dayjs(time).fromNow()
}

/**
 * 处理通知点击
 */
const handleNotificationClick = (notification: Notification) => {
  if (!notification.read) {
    notification.read = true
  }
  // 可以根据通知类型跳转到相应页面
  console.log('Notification clicked:', notification)
}

/**
 * 标记所有通知为已读
 */
const markAllAsRead = () => {
  notifications.value.forEach(n => n.read = true)
  message.success('已标记所有通知为已读')
}

/**
 * 查看所有通知
 */
const viewAllNotifications = () => {
  router.push('/notifications')
}

/**
 * 用户操作
 */
const viewProfile = () => {
  router.push('/profile')
}

const openSettings = () => {
  router.push('/settings')
}

const openHelp = () => {
  router.push('/help')
}

const handleLogout = () => {
  Modal.confirm({
    title: '确认退出',
    content: '确定要退出登录吗？',
    okText: '确定',
    cancelText: '取消',
    onOk: async () => {
      try {
        await authStore.logout()
        message.success('已安全退出')
        router.push('/login')
      } catch (error) {
        console.error('退出登录失败:', error)
        message.error('退出登录失败')
      }
    }
  })
}

/**
 * 生命周期
 */
onMounted(() => {
  document.addEventListener('fullscreenchange', handleFullscreenChange)
})

onUnmounted(() => {
  document.removeEventListener('fullscreenchange', handleFullscreenChange)
})
</script>

<style lang="scss" scoped>
.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
  padding: 0 24px;
  background: #fff;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
  
  .collapse-btn {
    font-size: 18px;
    
    &:hover {
      background: #f5f5f5;
    }
  }
  
  .global-search {
    :deep(.ant-input-search) {
      .ant-input {
        border-radius: 20px;
        border-color: #e8e8e8;
        
        &:hover, &:focus {
          border-color: #1890ff;
        }
      }
      
      .ant-input-search-button {
        border-radius: 0 20px 20px 0;
        border-left: none;
      }
    }
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.quick-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  
  .action-btn {
    font-size: 16px;
    
    &:hover {
      background: #f5f5f5;
    }
    
    .spinning {
      animation: spin 1s linear infinite;
    }
  }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.notifications {
  .action-btn {
    font-size: 16px;
    
    &:hover {
      background: #f5f5f5;
    }
  }
}

.notification-dropdown {
  width: 320px;
  max-height: 400px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  
  .notification-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    border-bottom: 1px solid #f0f0f0;
    font-weight: 600;
    
    a {
      color: #1890ff;
      font-size: 12px;
      
      &:hover {
        text-decoration: underline;
      }
    }
  }
  
  .notification-list {
    max-height: 280px;
    overflow-y: auto;
    
    .notification-item {
      padding: 12px 20px;
      border-bottom: 1px solid #f5f5f5;
      cursor: pointer;
      transition: background 0.2s;
      
      &:hover {
        background: #f9f9f9;
      }
      
      &.unread {
        background: #f6ffed;
        
        &::before {
          content: '';
          position: absolute;
          left: 8px;
          top: 50%;
          transform: translateY(-50%);
          width: 6px;
          height: 6px;
          background: #52c41a;
          border-radius: 50%;
        }
      }
      
      .notification-content {
        .notification-title {
          font-size: 14px;
          font-weight: 500;
          margin-bottom: 4px;
          color: #262626;
        }
        
        .notification-desc {
          font-size: 12px;
          color: #666;
          margin-bottom: 4px;
          line-height: 1.4;
        }
        
        .notification-time {
          font-size: 11px;
          color: #999;
        }
      }
    }
    
    .notification-empty {
      padding: 40px 20px;
      text-align: center;
      color: #999;
      font-size: 14px;
    }
  }
  
  .notification-footer {
    padding: 12px 20px;
    text-align: center;
    border-top: 1px solid #f0f0f0;
    
    a {
      color: #1890ff;
      font-size: 14px;
      
      &:hover {
        text-decoration: underline;
      }
    }
  }
}

.user-menu {
  .user-info {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 4px 8px;
    border-radius: 6px;
    cursor: pointer;
    transition: background 0.2s;
    
    &:hover {
      background: #f5f5f5;
    }
    
    .user-details {
      display: flex;
      flex-direction: column;
      
      .user-name {
        font-size: 14px;
        font-weight: 500;
        color: #262626;
        line-height: 1.2;
      }
      
      .user-role {
        font-size: 12px;
        color: #666;
        line-height: 1.2;
      }
    }
    
    .dropdown-icon {
      font-size: 12px;
      color: #999;
      margin-left: 4px;
    }
  }
}

.user-dropdown-menu {
  min-width: 160px;
  
  :deep(.ant-menu-item) {
    display: flex;
    align-items: center;
    gap: 8px;
    
    .anticon {
      font-size: 14px;
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .header-content {
    padding: 0 16px;
  }
  
  .global-search {
    display: none;
  }
  
  .quick-actions {
    gap: 4px;
  }
  
  .user-details {
    display: none;
  }
}
</style> 