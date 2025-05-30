<template>
  <div class="dashboard-container">
    <!-- 欢迎区域 -->
    <div class="welcome-section">
      <div class="welcome-card">
        <div class="welcome-content">
          <div class="welcome-text">
            <h2>欢迎回来，{{ userName }}！</h2>
            <p class="welcome-desc">
              今天是 {{ currentDate }}，{{ getTimeGreeting() }}，注意休息
            </p>
            <div class="welcome-stats">
              您有 
              <span class="stat-badge pending">{{ pendingTasks }}</span> 个待处理任务，
              <span class="stat-badge new">{{ todayNewApps }}</span> 个新申请
            </div>
          </div>
          <div class="welcome-actions">
            <a-space :size="16">
              <a-button type="primary" size="large" @click="goToApproval">
                <template #icon>
                  <FileTextOutlined />
                </template>
                开始审批
              </a-button>
              <a-button size="large" @click="goToStatistics">
                <template #icon>
                  <BarChartOutlined />
                </template>
                查看报表
              </a-button>
            </a-space>
          </div>
        </div>
      </div>
    </div>

    <!-- 统计卡片区域 -->
    <div class="stats-section">
      <a-row :gutter="[24, 24]">
        <a-col :xs="12" :sm="6" v-for="stat in statistics" :key="stat.key">
          <div class="stat-card" :class="stat.type">
            <div class="stat-header">
              <div class="stat-icon" :style="{ background: stat.color }">
                <component :is="getIconComponent(stat.icon)" />
              </div>
              <div class="stat-trend" :class="stat.trend.type">
                <component :is="getIconComponent(stat.trend.icon)" />
                <span>{{ stat.trend.value }}</span>
              </div>
            </div>
            <div class="stat-content">
              <div class="stat-value">
                {{ formatStatValue(stat.value, stat.precision) }}<span class="stat-suffix">{{ stat.suffix }}</span>
              </div>
              <div class="stat-title">{{ stat.title }}</div>
            </div>
            <div class="stat-footer">
              <span class="stat-period">{{ stat.trend.period }}</span>
            </div>
          </div>
        </a-col>
      </a-row>
    </div>

    <!-- 数据可视化图表区域 -->
    <div class="charts-section">
      <a-row :gutter="[24, 24]">
        <!-- 贷款申请趋势图 -->
        <a-col :xs="24" :lg="12">
          <a-card title="贷款申请趋势" class="chart-card">
            <template #extra>
              <a-select 
                v-model:value="loanTrendPeriod" 
                size="small" 
                @change="loadLoanTrendData"
                :options="periodOptions"
              />
            </template>
            <BusinessChart
              type="line"
              :data="loanTrendData"
              :loading="loading.loanTrend"
              :height="300"
              subtitle="近期贷款申请数量变化趋势"
              @retry="loadLoanTrendData"
            />
          </a-card>
        </a-col>

        <!-- 审批状态分布 -->
        <a-col :xs="24" :lg="12">
          <a-card title="审批状态分布" class="chart-card">
            <BusinessChart
              type="pie"
              :data="approvalStatusData"
              :loading="loading.approvalStatus"
              :height="300"
              subtitle="当前申请的审批状态分布"
              @retry="loadApprovalStatusData"
            />
          </a-card>
        </a-col>

        <!-- 月度业务量统计 -->
        <a-col :xs="24" :lg="12">
          <a-card title="月度业务量统计" class="chart-card">
            <BusinessChart
              type="bar"
              :data="monthlyBusinessData"
              :loading="loading.monthlyBusiness"
              :height="300"
              subtitle="各类业务的月度处理量对比"
              @retry="loadMonthlyBusinessData"
            />
          </a-card>
        </a-col>

        <!-- 风险评估分布 -->
        <a-col :xs="24" :lg="12">
          <a-card title="风险评估分布" class="chart-card">
            <BusinessChart
              type="gauge"
              :data="riskDistributionData"
              :loading="loading.riskDistribution"
              :height="300"
              subtitle="AI风险评估结果分布情况"
              @retry="loadRiskDistributionData"
            />
          </a-card>
        </a-col>
      </a-row>
    </div>

    <!-- 主要内容区域 -->
    <a-row :gutter="[16, 16]">
      <!-- 待处理事项 -->
      <a-col :xs="24" :lg="12">
        <a-card class="content-card" size="small">
          <template #title>
            <div class="card-header">
              <h3>待处理事项</h3>
              <a @click="viewAllTasks" class="view-all-link">查看全部</a>
            </div>
          </template>
          <div class="card-body">
            <div v-if="loading.tasks" class="loading-state">
              <a-spin size="large" />
            </div>
            <div v-else-if="pendingTaskList.length === 0" class="empty-state">
              <div class="empty-icon">📋</div>
              <p>暂无待处理事项</p>
            </div>
            <div v-else class="task-list">
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
          </div>
        </a-card>
      </a-col>

      <!-- 最新申请 -->
      <a-col :xs="24" :lg="12">
        <a-card class="content-card" size="small">
          <template #title>
            <div class="card-header">
              <h3>最新申请</h3>
              <a @click="viewAllApplications" class="view-all-link">查看全部</a>
            </div>
          </template>
          <div class="card-body">
            <div v-if="loading.applications" class="loading-state">
              <a-spin size="large" />
            </div>
            <div v-else-if="recentApplications.length === 0" class="empty-state">
              <div class="empty-icon">📄</div>
              <p>暂无最新申请</p>
            </div>
            <div v-else class="application-list">
              <div v-for="app in recentApplications" :key="app.id" class="application-item">
                <a-avatar :src="app.user_info.user_avatar" :size="40" class="user-avatar">
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
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 系统状态和快捷操作 -->
    <a-row :gutter="[16, 16]" style="margin-top: 16px;">
      <!-- AI工作流状态 -->
      <a-col :xs="24" :lg="8">
        <a-card class="content-card" size="small">
          <template #title>
            <h3>AI工作流状态</h3>
          </template>
          <div class="card-body">
            <div class="ai-status">
              <div class="ai-item">
                <div class="ai-info">
                  <div class="ai-icon processing">🤖</div>
                  <div class="ai-label">贷款审批AI</div>
                </div>
                <a-badge status="processing" text="运行中" />
              </div>
              <div class="ai-item">
                <div class="ai-info">
                  <div class="ai-icon success">🔍</div>
                  <div class="ai-label">风险评估AI</div>
                </div>
                <a-badge status="success" text="正常" />
              </div>
              <div class="ai-item">
                <div class="ai-info">
                  <div class="ai-icon warning">⚙️</div>
                  <div class="ai-label">设备推荐AI</div>
                </div>
                <a-badge status="warning" text="维护中" />
              </div>
            </div>
          </div>
        </a-card>
      </a-col>

      <!-- 快捷操作 -->
      <a-col :xs="24" :lg="8">
        <a-card class="content-card" size="small">
          <template #title>
            <h3>快捷操作</h3>
          </template>
          <div class="card-body">
            <div class="quick-actions-grid">
              <div v-for="action in quickActions" :key="action.key" class="quick-action" @click="action.handler">
                <div class="quick-action-icon" :style="{ background: action.color }">
                  <component :is="getIconComponent(action.icon)" />
                </div>
                <div class="quick-action-text">{{ action.title }}</div>
              </div>
            </div>
          </div>
        </a-card>
      </a-col>

      <!-- 系统公告 -->
      <a-col :xs="24" :lg="8">
        <a-card class="content-card" size="small">
          <template #title>
            <div class="card-header">
              <h3>系统公告</h3>
              <a @click="viewAllAnnouncements" class="view-all-link">更多</a>
            </div>
          </template>
          <div class="card-body">
            <div v-if="loading.announcements" class="loading-state">
              <a-spin size="small" />
            </div>
            <div v-else-if="announcements.length === 0" class="empty-state">
              <div class="empty-icon">📢</div>
              <p>暂无系统公告</p>
            </div>
            <div v-else class="announcement-list">
              <div v-for="announcement in announcements" :key="announcement.id" class="announcement-item">
                <div class="announcement-title">{{ announcement.title }}</div>
                <div class="announcement-time">{{ formatTime(announcement.created_at) }}</div>
              </div>
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
import { BusinessChart } from '@/components/charts'

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
  announcements: false,
  loanTrend: false,
  approvalStatus: false,
  monthlyBusiness: false,
  riskDistribution: false
})

// 图表相关状态
const loanTrendPeriod = ref('30days')
const loanTrendData = ref<any[]>([])
const approvalStatusData = ref<any[]>([])
const monthlyBusinessData = ref<any[]>([])
const riskDistributionData = ref<any[]>([])

// 期间选择选项
const periodOptions = [
  { label: '近7天', value: '7days' },
  { label: '近30天', value: '30days' },
  { label: '近3个月', value: '3months' },
  { label: '近半年', value: '6months' }
]

// 数据状态（使用模拟数据）
const overviewData = ref({
  total_users: 15432,
  pending_applications: 128,
  total_loan_amount: 25600000,
  today_new_users: 45
})

const pendingTaskList = ref([
  {
    id: 1,
    title: '农业贷款申请审批',
    description: '张三的50万农业生产贷款申请',
    type: '贷款审批',
    priority: '高',
    created_at: '2025-01-31T10:30:00Z'
  },
  {
    id: 2,
    title: '用户实名认证审核',
    description: '李四的身份信息认证材料审核',
    type: '用户认证',
    priority: '中',
    created_at: '2025-01-31T09:15:00Z'
  },
  {
    id: 3,
    title: '设备采购贷款申请',
    description: '王五的农机设备购买贷款申请',
    type: '贷款审批',
    priority: '高',
    created_at: '2025-01-31T08:45:00Z'
  },
  {
    id: 4,
    title: '农户资质审核',
    description: '赵六的农业合作社资质认证',
    type: '资质审核',
    priority: '中',
    created_at: '2025-01-31T08:20:00Z'
  },
  {
    id: 5,
    title: '贷款额度调整申请',
    description: '钱七的信用额度提升申请',
    type: '额度调整',
    priority: '低',
    created_at: '2025-01-31T07:55:00Z'
  },
  {
    id: 6,
    title: '风险评估报告审核',
    description: '周八的高风险用户复核',
    type: '风险审核',
    priority: '高',
    created_at: '2025-01-31T07:30:00Z'
  }
])

const recentApplications = ref([
  {
    id: 1,
    user_info: {
      real_name: '王农户',
      user_avatar: ''
    },
    product_name: '农业生产贷',
    amount: 300000,
    status: 'pending',
    status_text: '待审批',
    created_at: '2025-01-31T14:20:00Z'
  },
  {
    id: 2,
    user_info: {
      real_name: '李合作社',
      user_avatar: ''
    },
    product_name: '设备采购贷',
    amount: 800000,
    status: 'approved',
    status_text: '已通过',
    created_at: '2025-01-31T13:45:00Z'
  },
  {
    id: 3,
    user_info: {
      real_name: '张种植户',
      user_avatar: ''
    },
    product_name: '季节性周转贷',
    amount: 150000,
    status: 'reviewing',
    status_text: '审核中',
    created_at: '2025-01-31T12:30:00Z'
  },
  {
    id: 4,
    user_info: {
      real_name: '陈农场主',
      user_avatar: ''
    },
    product_name: '扩建贷款',
    amount: 1200000,
    status: 'pending',
    status_text: '待审批',
    created_at: '2025-01-31T11:15:00Z'
  },
  {
    id: 5,
    user_info: {
      real_name: '刘养殖户',
      user_avatar: ''
    },
    product_name: '养殖业贷款',
    amount: 450000,
    status: 'rejected',
    status_text: '已拒绝',
    created_at: '2025-01-31T10:20:00Z'
  }
])

const announcements = ref([
  {
    id: 1,
    title: '系统维护通知',
    content: '系统将于今晚进行例行维护',
    created_at: '2025-01-31T08:00:00Z'
  },
  {
    id: 2,
    title: '新功能上线',
    content: 'AI风险评估功能正式上线',
    created_at: '2025-01-30T16:30:00Z'
  }
])

// 计算的统计数据
const statistics = computed(() => [
  {
    key: 'users',
    title: '总用户数',
    value: overviewData.value?.total_users || 0,
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
    value: overviewData.value?.pending_applications || 0,
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
    value: (overviewData.value?.total_loan_amount || 0) / 1000000,
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
])

// 计算属性
const userName = computed(() => authStore.userName || '超级管理员')
const currentDate = computed(() => dayjs().format('YYYY年MM月DD日 dddd'))
const pendingTasks = computed(() => pendingTaskList.value?.length || 0)
const todayNewApps = computed(() => overviewData.value?.today_new_users || 0)

// 快捷操作配置
const quickActions = computed(() => [
  {
    key: 'add_user',
    title: '新增用户',
    icon: 'PlusOutlined',
    color: '#1890ff',
    handler: () => router.push('/user/create')
  },
  {
    key: 'ai_config',
    title: 'AI配置',
    icon: 'SettingOutlined',
    color: '#52c41a',
    handler: () => router.push('/ai/config')
  },
  {
    key: 'view_logs',
    title: '查看日志',
    icon: 'FileSearchOutlined',
    color: '#faad14',
    handler: () => router.push('/system/logs')
  },
  {
    key: 'send_notification',
    title: '发送通知',
    icon: 'NotificationOutlined',
    color: '#f5222d',
    handler: () => router.push('/system/notifications')
  }
])

/**
 * 图表数据加载方法
 */
const loadLoanTrendData = async () => {
  try {
    loading.value.loanTrend = true
    // 模拟数据，实际项目中替换为API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    const days = loanTrendPeriod.value === '7days' ? 7 : 
                 loanTrendPeriod.value === '30days' ? 30 : 
                 loanTrendPeriod.value === '3months' ? 90 : 180
    
    loanTrendData.value = Array.from({ length: days }, (_, i) => ({
      name: dayjs().subtract(days - 1 - i, 'day').format('MM-DD'),
      value: Math.floor(Math.random() * 50) + 20
    }))
  } catch (error) {
    console.error('加载贷款趋势数据失败:', error)
    message.error('加载图表数据失败')
  } finally {
    loading.value.loanTrend = false
  }
}

const loadApprovalStatusData = async () => {
  try {
    loading.value.approvalStatus = true
    await new Promise(resolve => setTimeout(resolve, 800))
    
    approvalStatusData.value = [
      { name: '待审批', value: 128 },
      { name: '已通过', value: 356 },
      { name: '已拒绝', value: 42 },
      { name: '补充材料', value: 23 }
    ]
  } catch (error) {
    console.error('加载审批状态数据失败:', error)
  } finally {
    loading.value.approvalStatus = false
  }
}

const loadMonthlyBusinessData = async () => {
  try {
    loading.value.monthlyBusiness = true
    await new Promise(resolve => setTimeout(resolve, 600))
    
    monthlyBusinessData.value = [
      { name: '贷款申请', value: 456 },
      { name: '农机租赁', value: 234 },
      { name: '保险购买', value: 123 },
      { name: '咨询服务', value: 345 }
    ]
  } catch (error) {
    console.error('加载月度业务数据失败:', error)
  } finally {
    loading.value.monthlyBusiness = false
  }
}

const loadRiskDistributionData = async () => {
  try {
    loading.value.riskDistribution = true
    await new Promise(resolve => setTimeout(resolve, 700))
    
    riskDistributionData.value = [{
      name: '平均风险评分',
      value: 78,
      max: 100,
      unit: '分'
    }]
  } catch (error) {
    console.error('加载风险分布数据失败:', error)
  } finally {
    loading.value.riskDistribution = false
  }
}

/**
 * 工具方法
 */
const getTimeGreeting = () => {
  const hour = dayjs().hour()
  if (hour < 6) return '夜深了'
  if (hour < 9) return '早上好'
  if (hour < 12) return '上午好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  if (hour < 22) return '晚上好'
  return '夜深了'
}

const formatStatValue = (value: number, precision: number = 0) => {
  return Number(value).toLocaleString('zh-CN', {
    minimumFractionDigits: precision,
    maximumFractionDigits: precision
  })
}

const formatCurrency = (amount: number) => {
  return (amount / 10000).toFixed(1) + '万'
}

const formatTime = (time: string) => {
  return dayjs(time).fromNow()
}

const getIconComponent = (iconName: string) => {
  const iconMap: Record<string, any> = {
    UserOutlined,
    FileTextOutlined,
    BankOutlined,
    RobotOutlined,
    TeamOutlined,
    SettingOutlined,
    ArrowUpOutlined,
    ArrowDownOutlined,
    FileSearchOutlined,
    NotificationOutlined,
    PlusOutlined
  }
  return iconMap[iconName] || UserOutlined
}

const getTaskTypeColor = (type: string) => {
  const colorMap: Record<string, string> = {
    '贷款审批': 'blue',
    '用户认证': 'green',
    '风险评估': 'orange',
    '系统维护': 'purple'
  }
  return colorMap[type] || 'default'
}

const getPriorityColor = (priority: string) => {
  const colorMap: Record<string, string> = {
    '高': 'red',
    '中': 'orange',
    '低': 'green'
  }
  return colorMap[priority] || 'default'
}

const getStatusColor = (status: string) => {
  const colorMap: Record<string, string> = {
    'pending': 'orange',
    'approved': 'green',
    'rejected': 'red',
    'processing': 'blue'
  }
  return colorMap[status] || 'default'
}

/**
 * 事件处理方法
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
  router.push('/system/announcements')
}

const handleTaskAction = (task: any) => {
  router.push(`/task/${task.id}`)
}

/**
 * 生命周期
 */
onMounted(async () => {
  // 并发加载所有数据
  await Promise.all([
    loadLoanTrendData(),
    loadApprovalStatusData(),
    loadMonthlyBusinessData(),
    loadRiskDistributionData()
  ])
})
</script>

<style lang="scss" scoped>
.dashboard-container {
  padding: 0;
  background: transparent;
  min-height: calc(100vh - 200px);

  // 欢迎区域
  .welcome-section {
    margin-bottom: 24px;
    
    .welcome-card {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      border-radius: 16px;
      padding: 32px;
      position: relative;
      overflow: hidden;
      box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
      
      &::before {
        content: '';
        position: absolute;
        top: 0;
        right: 0;
        width: 200px;
        height: 200px;
        background: rgba(255, 255, 255, 0.1);
        border-radius: 50%;
        transform: translate(50%, -50%);
      }
      
      .welcome-content {
        position: relative;
        z-index: 1;
        display: flex;
        justify-content: space-between;
        align-items: center;
        
        .welcome-text {
          color: white;
          
          h2 {
            margin: 0 0 8px 0;
            font-size: 28px;
            font-weight: 600;
            color: white;
          }
          
          .welcome-desc {
            margin-bottom: 16px;
            font-size: 16px;
            opacity: 0.9;
          }
          
          .welcome-stats {
            font-size: 14px;
            opacity: 0.85;
            
            .stat-badge {
              display: inline-block;
              padding: 4px 12px;
              border-radius: 12px;
              font-weight: 600;
              margin: 0 4px;
              
              &.pending {
                background: rgba(250, 173, 20, 0.2);
                color: #faad14;
              }
              
              &.new {
                background: rgba(82, 196, 26, 0.2);
                color: #52c41a;
              }
            }
          }
        }
        
        .welcome-actions {
          flex-shrink: 0;
        }
      }
    }
  }

  // 统计卡片区域
  .stats-section {
    margin-bottom: 24px;
    
    .stat-card {
      background: white;
      border-radius: 12px;
      padding: 20px;
      border: 1px solid rgba(0, 0, 0, 0.06);
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      position: relative;
      overflow: hidden;
      
      &::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 3px;
        background: linear-gradient(90deg, transparent 0%, var(--primary-color, #1890ff) 50%, transparent 100%);
        opacity: 0;
        transition: opacity 0.3s ease;
      }
      
      &:hover {
        transform: translateY(-4px);
        box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
        
        &::before {
          opacity: 1;
        }
      }
      
      .stat-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16px;
        
        .stat-icon {
          width: 48px;
          height: 48px;
          border-radius: 12px;
          display: flex;
          align-items: center;
          justify-content: center;
          color: white;
          font-size: 20px;
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        }
        
        .stat-trend {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 12px;
          font-weight: 500;
          
          &.increase {
            color: #52c41a;
          }
          
          &.decrease {
            color: #ff4d4f;
          }
          
          &.normal {
            color: #1890ff;
          }
        }
      }
      
      .stat-content {
        .stat-value {
          font-size: 32px;
          font-weight: 700;
          color: #262626;
          line-height: 1;
          margin-bottom: 8px;
          
          .stat-suffix {
            font-size: 16px;
            font-weight: 500;
            color: #8c8c8c;
            margin-left: 2px;
          }
        }
        
        .stat-title {
          font-size: 14px;
          color: #8c8c8c;
          font-weight: 500;
        }
      }
      
      .stat-footer {
        margin-top: 12px;
        
        .stat-period {
          font-size: 12px;
          color: #bfbfbf;
        }
      }
    }
  }

  // 图表区域
  .charts-section {
    margin-bottom: 24px;
    
    .chart-card {
      height: 100%;
      min-height: 400px;
      
      :deep(.ant-card-head) {
        border-bottom: 1px solid #f0f0f0;
        padding: 12px 16px;
        min-height: auto;
      }
      
      :deep(.ant-card-body) {
        padding: 16px;
        max-height: 400px;
        overflow-y: auto;
        overflow-x: hidden;
        
        // 自定义滚动条样式
        &::-webkit-scrollbar {
          width: 6px;
        }
        
        &::-webkit-scrollbar-track {
          background: #f5f5f5;
          border-radius: 3px;
        }
        
        &::-webkit-scrollbar-thumb {
          background: #d9d9d9;
          border-radius: 3px;
          transition: background 0.3s ease;
          
          &:hover {
            background: #bfbfbf;
          }
        }
        
        // Firefox滚动条样式
        scrollbar-width: thin;
        scrollbar-color: #d9d9d9 #f5f5f5;
      }
    }
  }

  // 内容卡片样式
  .content-card {
    height: 100%;
    min-height: 500px;
    display: flex;
    flex-direction: column;
    
    :deep(.ant-card-body) {
      flex: 1;
      padding: 0;
    }
    
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin: 0;
      
      h3 {
        margin: 0;
        font-size: 16px;
        font-weight: 600;
        color: #262626;
      }
      
      .view-all-link {
        color: #1890ff;
        font-size: 14px;
        text-decoration: none;
        
        &:hover {
          color: #40a9ff;
        }
      }
    }
    
    .card-body {
      flex: 1;
      max-height: none;
      min-height: 420px;
      overflow-y: auto;
      overflow-x: hidden;
      padding: 16px;
      
      // 自定义滚动条样式
      &::-webkit-scrollbar {
        width: 6px;
      }
      
      &::-webkit-scrollbar-track {
        background: #f5f5f5;
        border-radius: 3px;
      }
      
      &::-webkit-scrollbar-thumb {
        background: #d9d9d9;
        border-radius: 3px;
        transition: background 0.3s ease;
        
        &:hover {
          background: #bfbfbf;
        }
      }
      
      // Firefox滚动条样式
      scrollbar-width: thin;
      scrollbar-color: #d9d9d9 #f5f5f5;
      
      .loading-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 40px 0;
        color: #8c8c8c;
      }
      
      .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 40px 0;
        color: #8c8c8c;
        
        .empty-icon {
          font-size: 48px;
          margin-bottom: 16px;
          opacity: 0.5;
        }
        
        p {
          margin: 0;
          font-size: 14px;
        }
      }
      
      // 任务列表
      .task-list {
        height: 100%;
        
        .task-item {
          display: flex;
          justify-content: space-between;
          align-items: flex-start;
          padding: 12px 0;
          border-bottom: 1px solid #f5f5f5;
          
          &:last-child {
            border-bottom: none;
          }
          
          .task-content {
            flex: 1;
            
            .task-title {
              font-size: 14px;
              font-weight: 500;
              color: #262626;
              margin-bottom: 4px;
            }
            
            .task-desc {
              font-size: 13px;
              color: #8c8c8c;
              margin-bottom: 8px;
            }
            
            .task-meta {
              display: flex;
              align-items: center;
              gap: 8px;
              
              .task-time {
                font-size: 12px;
                color: #bfbfbf;
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
        height: 100%;
        
        .application-item {
          display: flex;
          align-items: center;
          padding: 12px 0;
          border-bottom: 1px solid #f5f5f5;
          
          &:last-child {
            border-bottom: none;
          }
          
          .user-avatar {
            margin-right: 12px;
            flex-shrink: 0;
          }
          
          .application-content {
            flex: 1;
            
            .application-title {
              font-size: 14px;
              font-weight: 500;
              color: #262626;
              margin-bottom: 4px;
            }
            
            .application-desc {
              font-size: 13px;
              color: #8c8c8c;
              margin-bottom: 4px;
            }
            
            .application-time {
              font-size: 12px;
              color: #bfbfbf;
            }
          }
          
          .application-status {
            margin-left: 12px;
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
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
            overflow: hidden;
          }
          
          .announcement-time {
            font-size: 12px;
            color: #bfbfbf;
          }
        }
      }
    }
  }

  // AI状态
  .ai-status {
    height: 100%;
    overflow-y: auto;
    overflow-x: hidden;
    
    // 自定义滚动条样式
    &::-webkit-scrollbar {
      width: 6px;
    }
    
    &::-webkit-scrollbar-track {
      background: #f5f5f5;
      border-radius: 3px;
    }
    
    &::-webkit-scrollbar-thumb {
      background: #d9d9d9;
      border-radius: 3px;
      transition: background 0.3s ease;
      
      &:hover {
        background: #bfbfbf;
      }
    }
    
    // Firefox滚动条样式
    scrollbar-width: thin;
    scrollbar-color: #d9d9d9 #f5f5f5;
    
    .ai-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 16px 0;
      border-bottom: 1px solid #f5f5f5;
      
      &:last-child {
        border-bottom: none;
      }
      
      .ai-info {
        display: flex;
        align-items: center;
        
        .ai-icon {
          font-size: 24px;
          margin-right: 12px;
          
          &.processing {
            animation: pulse 2s infinite;
          }
        }
        
        .ai-label {
          font-size: 14px;
          color: #262626;
        }
      }
    }
  }

  // 快捷操作
  .quick-actions-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
    height: 100%;
    align-content: start;
    
    // 自定义滚动条样式
    &::-webkit-scrollbar {
      width: 6px;
    }
    
    &::-webkit-scrollbar-track {
      background: #f5f5f5;
      border-radius: 3px;
    }
    
    &::-webkit-scrollbar-thumb {
      background: #d9d9d9;
      border-radius: 3px;
      transition: background 0.3s ease;
      
      &:hover {
        background: #bfbfbf;
      }
    }
    
    // Firefox滚动条样式
    scrollbar-width: thin;
    scrollbar-color: #d9d9d9 #f5f5f5;
    
    .quick-action {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 20px 12px;
      border-radius: 8px;
      background: #fafafa;
      cursor: pointer;
      transition: all 0.3s ease;
      
      &:hover {
        background: #f0f0f0;
        transform: translateY(-2px);
      }
      
      .quick-action-icon {
        width: 40px;
        height: 40px;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        font-size: 18px;
        margin-bottom: 8px;
      }
      
      .quick-action-text {
        font-size: 12px;
        color: #262626;
        text-align: center;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .dashboard-container {
    .welcome-section .welcome-card {
      padding: 20px;
      
      .welcome-content {
        flex-direction: column;
        align-items: flex-start;
        gap: 20px;
        
        .welcome-text h2 {
          font-size: 24px;
        }
      }
    }
    
    .stats-section .stat-card {
      padding: 16px;
      
      .stat-content .stat-value {
        font-size: 24px;
      }
    }
    
    .quick-actions-grid {
      grid-template-columns: repeat(4, 1fr);
      gap: 12px;
      
      .quick-action {
        padding: 16px 8px;
        
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
}

// 暗色主题适配
:deep([data-theme="dark"]) {
  .welcome-section .welcome-card {
    background: linear-gradient(135deg, #434343 0%, #000000 100%);
  }
  
  .stat-card {
    background: #1f1f1f !important;
    border-color: #303030 !important;
    
    .stat-content .stat-value {
      color: rgba(255, 255, 255, 0.85) !important;
    }
    
    .stat-title {
      color: rgba(255, 255, 255, 0.45) !important;
    }
  }
  
  .content-card {
    background: #1f1f1f !important;
    border-color: #303030 !important;
    
    .card-header h3 {
      color: rgba(255, 255, 255, 0.85) !important;
    }
    
    .task-title,
    .application-title,
    .announcement-title,
    .ai-label {
      color: rgba(255, 255, 255, 0.85) !important;
    }
    
    .task-desc,
    .application-desc {
      color: rgba(255, 255, 255, 0.45) !important;
    }
  }
  
  .quick-action {
    background: #262626 !important;
    
    &:hover {
      background: #434343 !important;
    }
    
    .quick-action-text {
      color: rgba(255, 255, 255, 0.85) !important;
    }
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