<template>
  <div class="dashboard-container">
    <!-- 欢迎区域 -->
    <div class="welcome-section">
      <a-card :bordered="false" class="welcome-card">
        <div class="welcome-content">
          <div class="welcome-text">
            <h2>欢迎回来，{{ userName }}！</h2>
            <p class="welcome-desc">
              今天是 {{ currentDate }}，{{ getTimeGreeting() }}
            </p>
            <p class="welcome-tips">
              您有 <a-tag color="orange">{{ pendingTasks }}</a-tag> 个待处理任务，
              <a-tag color="blue">{{ todayNewApps }}</a-tag> 个新申请
            </p>
          </div>
          <div class="welcome-actions">
            <a-space :size="16">
              <a-button type="primary" @click="goToApproval">
                <template #icon>
                  <FileTextOutlined />
                </template>
                开始审批
              </a-button>
              <a-button @click="goToStatistics">
                <template #icon>
                  <BarChartOutlined />
                </template>
                查看报表
              </a-button>
            </a-space>
          </div>
        </div>
      </a-card>
    </div>

    <!-- 统计概览 -->
    <div class="statistics-section">
      <a-row :gutter="[24, 24]">
        <a-col :xs="24" :sm="12" :lg="6" v-for="stat in statistics" :key="stat.key">
          <a-card :bordered="false" class="stat-card" :class="`stat-${stat.type}`">
            <a-statistic
              :title="stat.title"
              :value="stat.value"
              :precision="stat.precision"
              :value-style="{ color: stat.color }"
            >
              <template #prefix>
                <component :is="stat.icon" class="stat-icon" />
              </template>
              <template #suffix>
                <span class="stat-suffix">{{ stat.suffix }}</span>
              </template>
            </a-statistic>
            <div class="stat-trend">
              <span :class="stat.trend.type">
                <component :is="stat.trend.icon" />
                {{ stat.trend.value }}
              </span>
              <span class="stat-period">{{ stat.trend.period }}</span>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </div>

    <!-- 主要内容区域 -->
    <a-row :gutter="[24, 24]">
      <!-- 待处理事项 -->
      <a-col :xs="24" :lg="12">
        <a-card title="待处理事项" :bordered="false" class="content-card">
          <template #extra>
            <a @click="viewAllTasks">查看全部</a>
          </template>
          <div class="task-list">
            <div v-for="task in pendingTaskList" :key="task.id" class="task-item">
              <div class="task-content">
                <div class="task-title">{{ task.title }}</div>
                <div class="task-desc">{{ task.description }}</div>
                <div class="task-meta">
                  <a-tag :color="getTaskTypeColor(task.type)">{{ task.type }}</a-tag>
                  <a-tag v-if="task.priority" :color="getPriorityColor(task.priority)">{{ task.priority }}</a-tag>
                  <span class="task-time">{{ formatTime(task.created_at) }}</span>
                </div>
              </div>
              <div class="task-actions">
                <a-button size="small" type="primary" @click="handleTaskAction(task)">
                  处理
                </a-button>
              </div>
            </div>
          </div>
          <div v-if="pendingTaskList.length === 0 && !loading.tasks" class="empty-state">
            <a-empty description="暂无待处理事项" />
          </div>
          <div v-if="loading.tasks" class="loading-state">
            <a-spin size="large" />
          </div>
        </a-card>
      </a-col>

      <!-- 最新申请 -->
      <a-col :xs="24" :lg="12">
        <a-card title="最新申请" :bordered="false" class="content-card">
          <template #extra>
            <a @click="viewAllApplications">查看全部</a>
          </template>
          <div class="application-list">
            <div v-for="app in recentApplications" :key="app.id" class="application-item">
              <a-avatar :src="app.user_info.user_avatar" :size="40">
                <template #icon>
                  <UserOutlined />
                </template>
              </a-avatar>
              <div class="application-content">
                <div class="application-title">{{ app.user_info.real_name }} - {{ app.product_name }}</div>
                <div class="application-desc">申请金额: ¥{{ formatCurrency(app.amount) }}</div>
                <div class="application-time">{{ formatTime(app.created_at) }}</div>
              </div>
              <div class="application-status">
                <a-tag :color="getStatusColor(app.status)">{{ app.status_text }}</a-tag>
              </div>
            </div>
          </div>
          <div v-if="recentApplications.length === 0 && !loading.applications" class="empty-state">
            <a-empty description="暂无最新申请" />
          </div>
          <div v-if="loading.applications" class="loading-state">
            <a-spin size="large" />
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 系统状态和快捷操作 -->
    <a-row :gutter="[24, 24]" style="margin-top: 24px;">
      <!-- AI工作流状态 -->
      <a-col :xs="24" :lg="8">
        <a-card title="AI工作流状态" :bordered="false" class="content-card">
          <div class="ai-status">
            <div class="ai-item">
              <div class="ai-label">贷款审批AI</div>
              <a-badge status="processing" text="运行中" />
            </div>
            <div class="ai-item">
              <div class="ai-label">风险评估AI</div>
              <a-badge status="success" text="正常" />
            </div>
            <div class="ai-item">
              <div class="ai-label">设备推荐AI</div>
              <a-badge status="warning" text="维护中" />
            </div>
          </div>
        </a-card>
      </a-col>

      <!-- 快捷操作 -->
      <a-col :xs="24" :lg="8">
        <a-card title="快捷操作" :bordered="false" class="content-card">
          <div class="quick-actions-grid">
            <div v-for="action in quickActions" :key="action.key" class="quick-action" @click="action.handler">
              <div class="quick-action-icon" :style="{ background: action.color }">
                <component :is="action.icon" />
              </div>
              <div class="quick-action-text">{{ action.title }}</div>
            </div>
          </div>
        </a-card>
      </a-col>

      <!-- 系统公告 -->
      <a-col :xs="24" :lg="8">
        <a-card title="系统公告" :bordered="false" class="content-card">
          <template #extra>
            <a @click="viewAllAnnouncements">更多</a>
          </template>
          <div class="announcement-list">
            <div v-for="announcement in announcements" :key="announcement.id" class="announcement-item">
              <div class="announcement-title">{{ announcement.title }}</div>
              <div class="announcement-time">{{ formatTime(announcement.created_at) }}</div>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/modules/auth'
import { message } from 'ant-design-vue'
import {
  FileTextOutlined,
  BarChartOutlined,
  UserOutlined,
  BankOutlined,
  CarOutlined,
  TeamOutlined,
  RobotOutlined,
  ArrowUpOutlined,
  ArrowDownOutlined,
  PlusOutlined,
  SettingOutlined,
  FileSearchOutlined,
  NotificationOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'
import {
  getDashboardOverview,
  getRiskMonitoring,
  getSessionStatistics,
  getPendingTasks,
  getRecentApplications,
  getSystemAnnouncements,
  handleTask,
  type DashboardOverview,
  type RiskMonitoring,
  type SessionStatistics,
  type PendingTask,
  type RecentApplication
} from '@/api/modules/dashboard'

/**
 * 组件状态
 */
const router = useRouter()
const authStore = useAuthStore()

// 配置dayjs插件
dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

// 加载状态
const loading = ref({
  overview: false,
  tasks: false,
  applications: false,
  announcements: false
})

// 数据状态
const overviewData = ref<DashboardOverview | null>(null)
const riskData = ref<RiskMonitoring | null>(null)
const sessionData = ref<SessionStatistics | null>(null)
const pendingTaskList = ref<PendingTask[]>([])
const recentApplications = ref<RecentApplication[]>([])
const announcements = ref<Array<{
  id: number
  title: string
  content: string
  created_at: string
  updated_at: string
}>>([])

// 计算的统计数据
const statistics = computed(() => {
  if (!overviewData.value) return []
  
  const data = overviewData.value
  return [
    {
      key: 'users',
      title: '总用户数',
      value: data.total_users,
      color: '#1890ff',
      icon: 'UserOutlined',
      suffix: '',
      precision: 0,
      type: 'primary',
      trend: {
        type: 'increase',
        icon: 'ArrowUpOutlined',
        value: '12%',
        period: '较上月'
      }
    },
    {
      key: 'applications',
      title: '待审批申请',
      value: data.pending_applications,
      color: '#faad14',
      icon: 'FileTextOutlined',
      suffix: '',
      precision: 0,
      type: 'warning',
      trend: {
        type: 'decrease',
        icon: 'ArrowDownOutlined',
        value: '5%',
        period: '较昨日'
      }
    },
    {
      key: 'amount',
      title: '本月放款',
      value: data.total_loan_amount / 1000000,
      color: '#52c41a',
      icon: 'BankOutlined',
      suffix: 'M',
      precision: 1,
      type: 'success',
      trend: {
        type: 'increase',
        icon: 'ArrowUpOutlined',
        value: '8%',
        period: '较上月'
      }
    },
    {
      key: 'system',
      title: '系统状态',
      value: 99.9,
      color: '#722ed1',
      icon: 'RobotOutlined',
      suffix: '%',
      precision: 1,
      type: 'normal',
      trend: {
        type: 'normal',
        icon: 'ArrowUpOutlined',
        value: '正常运行',
        period: '7天'
      }
    }
  ]
})

// 计算属性
const userName = computed(() => authStore.userName || '管理员')
const currentDate = computed(() => dayjs().format('YYYY年MM月DD日 dddd'))
const pendingTasks = computed(() => pendingTaskList.value?.length || 0)
const todayNewApps = computed(() => overviewData.value?.today_new_users || 0)

/**
 * 数据加载方法
 */
const loadOverviewData = async () => {
  try {
    loading.value.overview = true
    const [overview, risk, session] = await Promise.allSettled([
      getDashboardOverview(),
      getRiskMonitoring().catch(() => null), // 风险监控可能不可用
      getSessionStatistics().catch(() => null) // 会话统计可能不可用
    ])
    
    // 处理概览数据
    if (overview.status === 'fulfilled') {
      overviewData.value = overview.value
    } else {
      console.error('加载概览数据失败:', overview.reason)
      message.error('加载概览数据失败')
    }
    
    // 处理风险监控数据
    if (risk.status === 'fulfilled' && risk.value) {
      riskData.value = risk.value
    }
    
    // 处理会话统计数据  
    if (session.status === 'fulfilled' && session.value) {
      sessionData.value = session.value
    }
    
  } catch (error) {
    console.error('加载概览数据失败:', error)
    message.error('加载概览数据失败')
  } finally {
    loading.value.overview = false
  }
}

const loadPendingTasks = async () => {
  try {
    loading.value.tasks = true
    const result = await getPendingTasks({ limit: 5 })
    pendingTaskList.value = result.tasks || []
  } catch (error: any) {
    console.error('加载待处理任务失败:', error)
    // 将tasks设为空数组，避免null错误
    pendingTaskList.value = []
    // 404错误不显示错误提示，因为后端可能还未实现
    if (error.response?.status !== 404) {
      message.error('加载待处理任务失败')
    }
  } finally {
    loading.value.tasks = false
  }
}

const loadRecentApplications = async () => {
  try {
    loading.value.applications = true
    const result = await getRecentApplications({ limit: 5 })
    recentApplications.value = result.applications || []
  } catch (error: any) {
    console.error('加载最新申请失败:', error)
    // 将applications设为空数组，避免null错误
    recentApplications.value = []
    if (error.response?.status !== 404) {
      message.error('加载最新申请失败')
    }
  } finally {
    loading.value.applications = false
  }
}

const loadAnnouncements = async () => {
  try {
    loading.value.announcements = true
    const result = await getSystemAnnouncements({ limit: 3, status: 'published' })
    announcements.value = result.announcements || []
  } catch (error) {
    console.error('加载系统公告失败:', error)
    // 将announcements设为空数组，避免null错误
    announcements.value = []
    // 公告失败不显示错误消息，因为这不是关键功能，且可能404
  } finally {
    loading.value.announcements = false
  }
}

/**
 * 事件处理方法
 */
const handleTaskAction = async (task: PendingTask) => {
  try {
    await handleTask(task.id, 'process')
    message.success(`开始处理任务: ${task.title}`)
    // 重新加载任务列表
    loadPendingTasks()
  } catch (error) {
    console.error('处理任务失败:', error)
    message.error('处理任务失败')
  }
}

const refreshData = async () => {
  await Promise.all([
    loadOverviewData(),
    loadPendingTasks(),
    loadRecentApplications(),
    loadAnnouncements()
  ])
  message.success('数据已刷新')
}

// 快捷操作
const quickActions = ref([
  {
    key: 'new-user',
    title: '新增用户',
    icon: 'PlusOutlined',
    color: '#1890ff',
    handler: () => router.push('/user/list')
  },
  {
    key: 'ai-config',
    title: 'AI配置',
    icon: 'SettingOutlined',
    color: '#52c41a',
    handler: () => router.push('/ai/settings')
  },
  {
    key: 'view-logs',
    title: '查看日志',
    icon: 'FileSearchOutlined',
    color: '#faad14',
    handler: () => router.push('/system/logs')
  },
  {
    key: 'send-notice',
    title: '发送通知',
    icon: 'NotificationOutlined',
    color: '#f5222d',
    handler: () => message.info('通知功能开发中...')
  }
])

/**
 * 工具方法
 */
const getTimeGreeting = () => {
  const hour = dayjs().hour()
  if (hour < 6) return '夜深了，注意休息'
  if (hour < 9) return '早上好'
  if (hour < 12) return '上午好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  if (hour < 22) return '晚上好'
  return '夜深了，注意休息'
}

const formatTime = (time: string) => {
  return dayjs(time).fromNow()
}

const formatCurrency = (amount: number) => {
  return amount.toLocaleString()
}

const getTaskTypeColor = (type: string) => {
  const colors: Record<string, string> = {
    '贷款审批': 'orange',
    '设备租赁': 'blue',
    '用户认证': 'green',
    'loan_approval': 'orange',
    'machine_rental': 'blue',
    'user_verification': 'green'
  }
  return colors[type] || 'default'
}

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    'pending': 'orange',
    'approved': 'green',
    'rejected': 'red',
    'under_review': 'blue'
  }
  return colors[status] || 'default'
}

const getPriorityColor = (priority: string) => {
  const colors: Record<string, string> = {
    'high': 'red',
    'medium': 'orange',
    'low': 'green'
  }
  return colors[priority] || 'default'
}

/**
 * 导航方法
 */
const goToApproval = () => {
  router.push('/loan/applications')
}

const goToStatistics = () => {
  router.push('/loan/statistics')
}

const viewAllTasks = () => {
  router.push('/tasks')
}

const viewAllApplications = () => {
  router.push('/loan/applications')
}

const viewAllAnnouncements = () => {
  router.push('/announcements')
}

/**
 * 生命周期
 */
onMounted(async () => {
  console.log('Dashboard mounted')
  
  // 并行加载所有数据
  await Promise.all([
    loadOverviewData(),
    loadPendingTasks(),
    loadRecentApplications(),
    loadAnnouncements()
  ])
})
</script>

<style lang="scss" scoped>
.dashboard-container {
  padding: 24px;
  background: #f0f2f5;
  min-height: calc(100vh - 152px);
}

// 欢迎区域
.welcome-section {
  margin-bottom: 24px;
  
  .welcome-card {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 12px;
    overflow: hidden;
    
    :deep(.ant-card-body) {
      padding: 32px;
    }
    
    .welcome-content {
      display: flex;
      justify-content: space-between;
      align-items: center;
      color: white;
      
      .welcome-text {
        h2 {
          color: white;
          margin: 0 0 8px;
          font-size: 28px;
          font-weight: 600;
        }
        
        .welcome-desc {
          font-size: 16px;
          opacity: 0.9;
          margin: 0 0 12px;
        }
        
        .welcome-tips {
          font-size: 14px;
          opacity: 0.8;
          margin: 0;
        }
      }
      
      .welcome-actions {
        .ant-btn {
          height: 40px;
          border-radius: 20px;
          font-weight: 500;
        }
      }
    }
  }
}

// 统计卡片
.stat-card {
  border-radius: 8px;
  transition: transform 0.2s, box-shadow 0.2s;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
  
  .stat-icon {
    font-size: 20px;
  }
  
  .stat-suffix {
    font-size: 14px;
    margin-left: 4px;
  }
  
  .stat-trend {
    margin-top: 8px;
    font-size: 12px;
    
    .increase {
      color: #52c41a;
    }
    
    .decrease {
      color: #ff4d4f;
    }
    
    .normal {
      color: #1890ff;
    }
    
    .stat-period {
      color: #999;
      margin-left: 8px;
    }
  }
}

// 内容卡片
.content-card {
  border-radius: 8px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
  
  :deep(.ant-card-head) {
    border-bottom: 1px solid #f5f5f5;
    
    .ant-card-head-title {
      font-weight: 600;
    }
  }
}

// 任务列表
.task-list {
  .task-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 0;
    border-bottom: 1px solid #f5f5f5;
    
    &:last-child {
      border-bottom: none;
    }
    
    .task-content {
      flex: 1;
      
      .task-title {
        font-size: 14px;
        font-weight: 500;
        margin-bottom: 4px;
        color: #262626;
      }
      
      .task-desc {
        font-size: 12px;
        color: #666;
        margin-bottom: 8px;
      }
      
      .task-meta {
        display: flex;
        align-items: center;
        gap: 8px;
        
        .task-time {
          font-size: 11px;
          color: #999;
        }
      }
    }
    
    .task-actions {
      margin-left: 16px;
    }
  }
}

// 申请列表
.application-list {
  .application-item {
    display: flex;
    align-items: center;
    padding: 16px 0;
    border-bottom: 1px solid #f5f5f5;
    
    &:last-child {
      border-bottom: none;
    }
    
    .application-content {
      flex: 1;
      margin-left: 12px;
      
      .application-title {
        font-size: 14px;
        font-weight: 500;
        margin-bottom: 4px;
        color: #262626;
      }
      
      .application-desc {
        font-size: 12px;
        color: #666;
        margin-bottom: 4px;
      }
      
      .application-time {
        font-size: 11px;
        color: #999;
      }
    }
    
    .application-status {
      margin-left: 16px;
    }
  }
}

// AI状态
.ai-status {
  .ai-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 0;
    border-bottom: 1px solid #f5f5f5;
    
    &:last-child {
      border-bottom: none;
    }
    
    .ai-label {
      font-size: 14px;
      color: #262626;
    }
  }
}

// 快捷操作
.quick-actions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  
  .quick-action {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20px 16px;
    border-radius: 8px;
    background: #fafafa;
    cursor: pointer;
    transition: all 0.2s;
    
    &:hover {
      background: #f0f0f0;
      transform: translateY(-1px);
    }
    
    .quick-action-icon {
      width: 40px;
      height: 40px;
      border-radius: 20px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 18px;
      margin-bottom: 8px;
    }
    
    .quick-action-text {
      font-size: 12px;
      color: #666;
      text-align: center;
    }
  }
}

// 公告列表
.announcement-list {
  .announcement-item {
    padding: 12px 0;
    border-bottom: 1px solid #f5f5f5;
    
    &:last-child {
      border-bottom: none;
    }
    
    .announcement-title {
      font-size: 14px;
      color: #262626;
      margin-bottom: 4px;
      line-height: 1.4;
    }
    
    .announcement-time {
      font-size: 11px;
      color: #999;
    }
  }
}

// 空状态
.empty-state {
  padding: 40px 0;
  text-align: center;
}

// 响应式设计
@media (max-width: 768px) {
  .dashboard-container {
    padding: 16px;
  }
  
  .welcome-content {
    flex-direction: column;
    align-items: flex-start !important;
    gap: 20px;
  }
  
  .quick-actions-grid {
    grid-template-columns: repeat(4, 1fr);
    gap: 8px;
    
    .quick-action {
      padding: 12px 8px;
      
      .quick-action-icon {
        width: 32px;
        height: 32px;
        font-size: 16px;
      }
      
      .quick-action-text {
        font-size: 11px;
      }
    }
  }
}
</style>