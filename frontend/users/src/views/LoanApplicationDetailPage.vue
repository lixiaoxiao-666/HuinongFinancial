<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../stores/user'
import { loanApi } from '../services/api'
import type { LoanApplication } from '../services/api'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const loading = ref(false)

// 申请详情
const application = ref<LoanApplication | null>(null)

// 状态映射
const statusMap = {
  'SUBMITTED': { text: '已提交', color: '#409EFF', step: 1 },
  'AI_REVIEWING': { text: 'AI审核中', color: '#E6A23C', step: 2 },
  'AI_APPROVED': { text: 'AI审核通过', color: '#67C23A', step: 3 },
  'AI_REJECTED': { text: 'AI审核拒绝', color: '#F56C6C', step: 3 },
  'MANUAL_REVIEWING': { text: '人工复审中', color: '#E6A23C', step: 4 },
  'APPROVED': { text: '审核通过', color: '#67C23A', step: 5 },
  'REJECTED': { text: '审核拒绝', color: '#F56C6C', step: 5 },
  'NEED_MORE_INFO': { text: '需要补充资料', color: '#E6A23C', step: 3 }
}

// 当前状态信息
const currentStatus = computed(() => {
  if (!application.value) return null
  return statusMap[application.value.status as keyof typeof statusMap] || {
    text: application.value.status,
    color: '#909399',
    step: 1
  }
})

// 进度步骤
const steps = [
  { title: '提交申请', status: 'finish' },
  { title: 'AI智能审核', status: 'process' },
  { title: '审核完成', status: 'wait' }
]

// 根据申请状态计算步骤状态
const processedSteps = computed(() => {
  if (!application.value) return steps

  const status = application.value.status
  const newSteps = [...steps]

  switch (status) {
    case 'SUBMITTED':
      newSteps[0].status = 'finish'
      newSteps[1].status = 'process'
      newSteps[2].status = 'wait'
      break
    case 'AI_REVIEWING':
      newSteps[0].status = 'finish'
      newSteps[1].status = 'process'
      newSteps[2].status = 'wait'
      break
    case 'AI_APPROVED':
    case 'MANUAL_REVIEWING':
      newSteps[0].status = 'finish'
      newSteps[1].status = 'finish'
      newSteps[2].status = 'process'
      break
    case 'APPROVED':
      newSteps[0].status = 'finish'
      newSteps[1].status = 'finish'
      newSteps[2].status = 'finish'
      break
    case 'AI_REJECTED':
    case 'REJECTED':
      newSteps[0].status = 'finish'
      newSteps[1].status = 'error'
      newSteps[2].status = 'wait'
      break
    case 'NEED_MORE_INFO':
      newSteps[0].status = 'finish'
      newSteps[1].status = 'error'
      newSteps[2].status = 'wait'
      break
  }

  return newSteps
})

// 格式化时间
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 加载申请详情
const loadApplicationDetail = async () => {
  const applicationId = route.params.applicationId as string
  if (!applicationId) {
    ElMessage.error('缺少申请信息')
    router.go(-1)
    return
  }

  try {
    loading.value = true
    const response = await loanApi.getApplicationDetail(applicationId)
    application.value = response.data
  } catch (error: any) {
    console.error('加载申请详情失败:', error)
    ElMessage.error('加载申请详情失败')
    router.go(-1)
  } finally {
    loading.value = false
  }
}

// 返回上一页
const goBack = () => {
  router.go(-1)
}

// 查看我的申请
const goToMyApplications = () => {
  router.push('/loan/my-applications')
}

// 组件挂载时加载数据
onMounted(() => {
  // 检查登录状态
  if (!userStore.isLoggedIn) {
    ElMessage.error('请先登录')
    router.push('/login')
    return
  }
  
  loadApplicationDetail()
})
</script>

<template>
  <div class="application-detail-page">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <div class="nav-left" @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </div>
      <div class="nav-title">申请详情</div>
      <div class="nav-right"></div>
    </div>

    <div class="page-content" v-if="application && !loading">
      <!-- 状态卡片 -->
      <div class="status-card">
        <div class="status-header">
          <div class="status-icon" :style="{ backgroundColor: currentStatus?.color }">
            <el-icon v-if="application.status === 'APPROVED'"><Check /></el-icon>
            <el-icon v-else-if="application.status === 'REJECTED' || application.status === 'AI_REJECTED'"><Close /></el-icon>
            <el-icon v-else><Clock /></el-icon>
          </div>
          <h2 class="status-text" :style="{ color: currentStatus?.color }">
            {{ currentStatus?.text }}
          </h2>
        </div>
        
        <div class="status-desc">
          <p v-if="application.remarks">{{ application.remarks }}</p>
          <p v-else-if="application.status === 'SUBMITTED'">您的申请已提交，请耐心等待审核</p>
          <p v-else-if="application.status === 'AI_REVIEWING'">AI系统正在智能分析您的申请信息</p>
          <p v-else-if="application.status === 'APPROVED'">恭喜！您的贷款申请已通过审核</p>
          <p v-else-if="application.status === 'REJECTED'">很遗憾，您的申请未通过审核</p>
        </div>

        <!-- 进度条 -->
        <div class="progress-section">
          <el-steps :active="currentStatus?.step || 1" align-center>
            <el-step
              v-for="(step, index) in processedSteps"
              :key="index"
              :title="step.title"
              :status="step.status"
            />
          </el-steps>
        </div>
      </div>

      <!-- 申请信息 -->
      <div class="info-card">
        <div class="card-title">申请信息</div>
        <div class="info-list">
          <div class="info-item">
            <span class="label">申请编号:</span>
            <span class="value">{{ application.application_id }}</span>
          </div>
          <div class="info-item">
            <span class="label">申请金额:</span>
            <span class="value amount">¥{{ application.amount.toLocaleString() }}</span>
          </div>
          <div class="info-item">
            <span class="label">贷款期限:</span>
            <span class="value">{{ application.term_months }}个月</span>
          </div>
          <div class="info-item">
            <span class="label">申请用途:</span>
            <span class="value">{{ application.purpose }}</span>
          </div>
          <div class="info-item">
            <span class="label">提交时间:</span>
            <span class="value">{{ formatDate(application.submitted_at) }}</span>
          </div>
          <div class="info-item">
            <span class="label">更新时间:</span>
            <span class="value">{{ formatDate(application.updated_at) }}</span>
          </div>
          <div v-if="application.approved_amount" class="info-item">
            <span class="label">批准金额:</span>
            <span class="value amount">¥{{ application.approved_amount.toLocaleString() }}</span>
          </div>
        </div>
      </div>

      <!-- 审核历史 -->
      <div v-if="application.history && application.history.length > 0" class="history-card">
        <div class="card-title">审核历史</div>
        <div class="timeline">
          <div 
            v-for="(item, index) in application.history" 
            :key="index" 
            class="timeline-item"
          >
            <div class="timeline-dot"></div>
            <div class="timeline-content">
              <div class="timeline-status">{{ statusMap[item.status as keyof typeof statusMap]?.text || item.status }}</div>
              <div class="timeline-time">{{ formatDate(item.timestamp) }}</div>
              <div class="timeline-operator">操作人：{{ item.operator }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <el-button 
          type="primary" 
          @click="goToMyApplications"
          size="large"
          class="action-btn"
        >
          查看我的申请
        </el-button>
        
        <el-button 
          v-if="application.status === 'NEED_MORE_INFO'"
          type="warning"
          size="large"
          class="action-btn"
        >
          补充资料
        </el-button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-icon class="is-loading"><Loading /></el-icon>
      <p>加载中...</p>
    </div>
  </div>
</template>

<style scoped>
.application-detail-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.top-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: white;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-left {
  cursor: pointer;
  padding: 8px;
}

.nav-title {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
}

.nav-right {
  width: 32px;
}

.page-content {
  padding: 16px;
  max-width: 600px;
  margin: 0 auto;
}

.status-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.status-header {
  margin-bottom: 20px;
}

.status-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 12px;
  color: white;
  font-size: 24px;
}

.status-text {
  font-size: 24px;
  font-weight: 600;
  margin: 0;
}

.status-desc {
  color: #7f8c8d;
  margin-bottom: 24px;
}

.progress-section {
  margin-top: 20px;
}

.info-card, .history-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 2px solid #27ae60;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.info-item:last-child {
  border-bottom: none;
}

.label {
  color: #7f8c8d;
  font-size: 14px;
}

.value {
  color: #2c3e50;
  font-weight: 500;
  text-align: right;
}

.value.amount {
  color: #27ae60;
  font-weight: 600;
  font-size: 16px;
}

.timeline {
  position: relative;
}

.timeline::before {
  content: '';
  position: absolute;
  left: 20px;
  top: 0;
  bottom: 0;
  width: 2px;
  background: #e1e1e1;
}

.timeline-item {
  position: relative;
  padding-left: 50px;
  margin-bottom: 20px;
}

.timeline-item:last-child {
  margin-bottom: 0;
}

.timeline-dot {
  position: absolute;
  left: 14px;
  top: 5px;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #27ae60;
  border: 3px solid white;
  box-shadow: 0 0 0 2px #27ae60;
}

.timeline-content {
  background: #f8f9fa;
  padding: 12px;
  border-radius: 8px;
}

.timeline-status {
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 4px;
}

.timeline-time {
  font-size: 12px;
  color: #7f8c8d;
  margin-bottom: 2px;
}

.timeline-operator {
  font-size: 12px;
  color: #7f8c8d;
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 20px 0;
}

.action-btn {
  height: 50px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 25px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: #7f8c8d;
}

.loading-container .el-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

:deep(.el-steps) {
  padding: 0;
}

:deep(.el-step__title) {
  font-size: 12px;
}
</style> 