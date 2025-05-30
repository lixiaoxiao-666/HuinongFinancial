<template>
  <div class="footer-content">
    <!-- 主要信息 -->
    <div class="footer-main">
      <div class="copyright">
        <span>© 2024 数字惠农金融系统</span>
        <a-divider type="vertical" />
        <span>版权所有</span>
        <a-divider type="vertical" />
        <a href="#" @click.prevent="showTerms">服务条款</a>
        <a-divider type="vertical" />
        <a href="#" @click.prevent="showPrivacy">隐私政策</a>
      </div>
      
      <div class="system-info">
        <a-space :size="16">
          <span>版本 {{ version }}</span>
          <span class="system-status" :class="systemStatusClass">
            <a-badge :status="systemStatusType" />
            {{ systemStatusText }}
          </span>
          <span class="online-users">
            <UserOutlined />
            在线用户: {{ onlineUsers }}
          </span>
        </a-space>
      </div>
    </div>

    <!-- 扩展信息（可选显示） -->
    <div v-if="showExtendedInfo" class="footer-extended">
      <a-row :gutter="24">
        <a-col :span="6">
          <div class="footer-section">
            <h4>快速链接</h4>
            <ul class="footer-links">
              <li><a href="#" @click.prevent="goToHelp">帮助中心</a></li>
              <li><a href="#" @click.prevent="goToFeedback">意见反馈</a></li>
              <li><a href="#" @click.prevent="goToSupport">技术支持</a></li>
            </ul>
          </div>
        </a-col>
        
        <a-col :span="6">
          <div class="footer-section">
            <h4>系统监控</h4>
            <ul class="footer-links">
              <li>服务器状态: <span class="status-good">正常</span></li>
              <li>数据库状态: <span class="status-good">正常</span></li>
              <li>缓存状态: <span class="status-good">正常</span></li>
            </ul>
          </div>
        </a-col>
        
        <a-col :span="6">
          <div class="footer-section">
            <h4>联系方式</h4>
            <ul class="footer-links">
              <li>客服热线: 400-1234-5678</li>
              <li>邮箱: support@huinong.com</li>
              <li>工作时间: 9:00-18:00</li>
            </ul>
          </div>
        </a-col>
        
        <a-col :span="6">
          <div class="footer-section">
            <h4>数据统计</h4>
            <ul class="footer-links">
              <li>今日访问: {{ todayVisits }}</li>
              <li>在线时长: {{ onlineTime }}</li>
              <li>最后更新: {{ lastUpdate }}</li>
            </ul>
          </div>
        </a-col>
      </a-row>
    </div>

    <!-- 折叠/展开按钮 -->
    <div class="footer-toggle">
      <a-button
        type="text"
        size="small"
        @click="toggleExtendedInfo"
      >
        <template #icon>
          <UpOutlined v-if="showExtendedInfo" />
          <DownOutlined v-else />
        </template>
        {{ showExtendedInfo ? '收起' : '更多信息' }}
      </a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { UserOutlined, UpOutlined, DownOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'

/**
 * 组件状态
 */
const router = useRouter()
const showExtendedInfo = ref(false)
const systemStatus = ref<'normal' | 'warning' | 'error'>('normal')
const onlineUsers = ref(156)
const todayVisits = ref(2847)
const onlineTime = ref('')
const lastUpdate = ref('')

/**
 * 计算属性
 */
const version = computed(() => import.meta.env.VITE_APP_VERSION || '1.0.0')

const systemStatusClass = computed(() => ({
  'status-normal': systemStatus.value === 'normal',
  'status-warning': systemStatus.value === 'warning',
  'status-error': systemStatus.value === 'error'
}))

const systemStatusType = computed(() => {
  switch (systemStatus.value) {
    case 'normal': return 'success'
    case 'warning': return 'warning'
    case 'error': return 'error'
    default: return 'default'
  }
})

const systemStatusText = computed(() => {
  switch (systemStatus.value) {
    case 'normal': return '系统正常'
    case 'warning': return '系统警告'
    case 'error': return '系统异常'
    default: return '未知状态'
  }
})

/**
 * 方法定义
 */
const toggleExtendedInfo = () => {
  showExtendedInfo.value = !showExtendedInfo.value
}

const updateOnlineTime = () => {
  const startTime = localStorage.getItem('loginTime')
  if (startTime) {
    const duration = dayjs().diff(dayjs(startTime), 'minute')
    const hours = Math.floor(duration / 60)
    const minutes = duration % 60
    onlineTime.value = `${hours}小时${minutes}分钟`
  } else {
    onlineTime.value = '未知'
  }
}

const updateLastUpdate = () => {
  lastUpdate.value = dayjs().format('MM-DD HH:mm')
}

// 模拟系统状态检查
const checkSystemStatus = () => {
  // 这里可以调用实际的系统状态API
  // 模拟随机状态变化
  const statuses = ['normal', 'normal', 'normal', 'warning', 'error']
  const randomIndex = Math.floor(Math.random() * statuses.length)
  systemStatus.value = statuses[randomIndex] as 'normal' | 'warning' | 'error'
}

// 模拟在线用户数变化
const updateOnlineUsers = () => {
  const baseCount = 150
  const variation = Math.floor(Math.random() * 20) - 10
  onlineUsers.value = baseCount + variation
}

// 模拟今日访问量变化
const updateTodayVisits = () => {
  todayVisits.value += Math.floor(Math.random() * 5)
}

/**
 * 导航方法
 */
const goToHelp = () => {
  router.push('/help')
}

const goToFeedback = () => {
  message.info('意见反馈功能开发中...')
}

const goToSupport = () => {
  message.info('技术支持联系方式: support@huinong.com')
}

const showTerms = () => {
  message.info('服务条款页面开发中...')
}

const showPrivacy = () => {
  message.info('隐私政策页面开发中...')
}

/**
 * 生命周期
 */
let statusTimer: NodeJS.Timeout
let timeTimer: NodeJS.Timeout
let dataTimer: NodeJS.Timeout

onMounted(() => {
  // 初始化数据
  updateOnlineTime()
  updateLastUpdate()
  
  // 设置定时器
  statusTimer = setInterval(checkSystemStatus, 30000) // 30秒检查一次系统状态
  timeTimer = setInterval(updateOnlineTime, 60000) // 1分钟更新一次在线时间
  dataTimer = setInterval(() => {
    updateOnlineUsers()
    updateTodayVisits()
    updateLastUpdate()
  }, 10000) // 10秒更新一次数据
  
  // 记录登录时间（如果还没有记录）
  if (!localStorage.getItem('loginTime')) {
    localStorage.setItem('loginTime', dayjs().toISOString())
  }
})

onUnmounted(() => {
  if (statusTimer) clearInterval(statusTimer)
  if (timeTimer) clearInterval(timeTimer)
  if (dataTimer) clearInterval(dataTimer)
})
</script>

<style lang="scss" scoped>
.footer-content {
  padding: 16px 24px;
  background: #fff;
  border-top: 1px solid #f0f0f0;
  
  .footer-main {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 12px;
    color: #666;
    
    .copyright {
      display: flex;
      align-items: center;
      
      a {
        color: #1890ff;
        text-decoration: none;
        
        &:hover {
          text-decoration: underline;
        }
      }
    }
    
    .system-info {
      display: flex;
      align-items: center;
      
      .system-status {
        display: flex;
        align-items: center;
        
        &.status-normal {
          color: #52c41a;
        }
        
        &.status-warning {
          color: #faad14;
        }
        
        &.status-error {
          color: #ff4d4f;
        }
      }
      
      .online-users {
        display: flex;
        align-items: center;
        gap: 4px;
        color: #1890ff;
      }
    }
  }
  
  .footer-extended {
    margin-top: 24px;
    padding-top: 24px;
    border-top: 1px solid #f5f5f5;
    
    .footer-section {
      h4 {
        font-size: 14px;
        font-weight: 600;
        color: #262626;
        margin-bottom: 12px;
      }
      
      .footer-links {
        list-style: none;
        padding: 0;
        margin: 0;
        
        li {
          margin-bottom: 8px;
          font-size: 12px;
          color: #666;
          
          a {
            color: #1890ff;
            text-decoration: none;
            
            &:hover {
              text-decoration: underline;
            }
          }
          
          .status-good {
            color: #52c41a;
            font-weight: 500;
          }
          
          .status-warning {
            color: #faad14;
            font-weight: 500;
          }
          
          .status-error {
            color: #ff4d4f;
            font-weight: 500;
          }
        }
      }
    }
  }
  
  .footer-toggle {
    display: flex;
    justify-content: center;
    margin-top: 16px;
    
    .ant-btn {
      font-size: 12px;
      color: #666;
      
      &:hover {
        color: #1890ff;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .footer-content {
    padding: 12px 16px;
    
    .footer-main {
      flex-direction: column;
      gap: 8px;
      text-align: center;
      
      .copyright {
        justify-content: center;
        flex-wrap: wrap;
      }
      
      .system-info {
        justify-content: center;
        
        :deep(.ant-space) {
          flex-wrap: wrap;
          justify-content: center;
        }
      }
    }
    
    .footer-extended {
      .ant-row {
        .ant-col {
          margin-bottom: 16px;
        }
      }
    }
  }
}

@media (max-width: 576px) {
  .footer-content {
    .footer-main {
      .system-info {
        :deep(.ant-space-item) {
          margin-bottom: 4px;
        }
      }
    }
    
    .footer-extended {
      .ant-row {
        .ant-col {
          width: 100% !important;
          flex: none !important;
        }
      }
    }
  }
}
</style> 