<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../stores/user'
import { loanApi } from '../services/api'
import type { LoanApplication } from '../services/api'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const refreshing = ref(false)

// 申请列表
const applications = ref<LoanApplication[]>([])
const total = ref(0)
const page = ref(1)
const limit = ref(10)

// 状态筛选
const selectedStatus = ref('')
const statusOptions = [
  { label: '全部', value: '' },
  { label: '已提交', value: 'SUBMITTED' },
  { label: 'AI审核中', value: 'AI_REVIEWING' },
  { label: 'AI审核通过', value: 'AI_APPROVED' },
  { label: 'AI审核拒绝', value: 'AI_REJECTED' },
  { label: '人工复审中', value: 'MANUAL_REVIEWING' },
  { label: '审核通过', value: 'APPROVED' },
  { label: '审核拒绝', value: 'REJECTED' },
  { label: '需要补充资料', value: 'NEED_MORE_INFO' }
]

// 状态映射
const statusMap = {
  'SUBMITTED': { text: '已提交', color: '#409EFF' },
  'AI_REVIEWING': { text: 'AI审核中', color: '#E6A23C' },
  'AI_APPROVED': { text: 'AI审核通过', color: '#67C23A' },
  'AI_REJECTED': { text: 'AI审核拒绝', color: '#F56C6C' },
  'MANUAL_REVIEWING': { text: '人工复审中', color: '#E6A23C' },
  'APPROVED': { text: '审核通过', color: '#67C23A' },
  'REJECTED': { text: '审核拒绝', color: '#F56C6C' },
  'NEED_MORE_INFO': { text: '需要补充资料', color: '#E6A23C' }
}

// 是否有更多数据
const hasMore = computed(() => {
  return applications.value.length < total.value
})

// 格式化时间
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

// 获取状态信息
const getStatusInfo = (status: string) => {
  return statusMap[status as keyof typeof statusMap] || {
    text: status,
    color: '#909399'
  }
}

// 加载申请列表
const loadApplications = async (reset = false) => {
  try {
    if (reset) {
      loading.value = true
      page.value = 1
      applications.value = []
    } else {
      refreshing.value = true
    }

    const params = {
      page: page.value,
      limit: limit.value,
      ...(selectedStatus.value && { status: selectedStatus.value })
    }

    const response = await loanApi.getMyApplications(params)
    
    if (reset) {
      applications.value = response.data
    } else {
      applications.value.push(...response.data)
    }
    
    total.value = response.total

  } catch (error: any) {
    console.error('加载申请列表失败:', error)
    ElMessage.error('加载申请列表失败')
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

// 加载更多
const loadMore = () => {
  if (hasMore.value && !refreshing.value) {
    page.value++
    loadApplications()
  }
}

// 筛选状态变化
const handleStatusChange = () => {
  loadApplications(true)
}

// 查看申请详情
const viewDetail = (applicationId: string) => {
  router.push(`/loan/application/${applicationId}`)
}

// 返回上一页
const goBack = () => {
  router.go(-1)
}

// 申请新贷款
const applyNewLoan = () => {
  router.push('/finance')
}

// 组件挂载时加载数据
onMounted(() => {
  // 检查登录状态
  if (!userStore.isLoggedIn) {
    ElMessage.error('请先登录')
    router.push('/login')
    return
  }
  
  loadApplications(true)
})
</script>

<template>
  <div class="my-applications-page">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <div class="nav-left" @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </div>
      <div class="nav-title">我的申请</div>
      <div class="nav-right"></div>
    </div>

    <div class="page-content">
      <!-- 筛选栏 -->
      <div class="filter-bar">
        <el-select 
          v-model="selectedStatus" 
          placeholder="筛选状态"
          @change="handleStatusChange"
          style="width: 200px"
        >
          <el-option
            v-for="option in statusOptions"
            :key="option.value"
            :label="option.label"
            :value="option.value"
          />
        </el-select>
      </div>

      <!-- 申请列表 -->
      <div v-if="!loading" class="applications-list">
        <div 
          v-for="app in applications" 
          :key="app.application_id"
          class="application-card"
          @click="viewDetail(app.application_id)"
        >
          <div class="card-header">
            <div class="app-id">申请编号：{{ app.application_id.slice(-8) }}</div>
            <div 
              class="status-tag"
              :style="{ 
                backgroundColor: getStatusInfo(app.status).color,
                color: 'white'
              }"
            >
              {{ getStatusInfo(app.status).text }}
            </div>
          </div>
          
          <div class="card-body">
            <div class="amount-info">
              <span class="amount-label">申请金额</span>
              <span class="amount-value">¥{{ app.amount.toLocaleString() }}</span>
            </div>
            
            <div class="info-row">
              <div class="info-item">
                <span class="label">贷款期限</span>
                <span class="value">{{ app.term_months }}个月</span>
              </div>
              <div class="info-item">
                <span class="label">申请时间</span>
                <span class="value">{{ formatDate(app.submitted_at) }}</span>
              </div>
            </div>
            
            <div class="purpose">
              <span class="purpose-label">用途：</span>
              <span class="purpose-text">{{ app.purpose }}</span>
            </div>

            <div v-if="app.approved_amount" class="approved-amount">
              <span class="approved-label">批准金额：</span>
              <span class="approved-value">¥{{ app.approved_amount.toLocaleString() }}</span>
            </div>
          </div>

          <div class="card-footer">
            <div class="update-time">
              更新时间：{{ formatDate(app.updated_at) }}
            </div>
            <el-icon class="arrow-icon"><ArrowRight /></el-icon>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-if="applications.length === 0" class="empty-state">
          <el-empty description="暂无申请记录">
            <el-button type="primary" @click="applyNewLoan">
              立即申请
            </el-button>
          </el-empty>
        </div>

        <!-- 加载更多 -->
        <div v-if="hasMore" class="load-more">
          <el-button 
            @click="loadMore" 
            :loading="refreshing"
            type="text"
            class="load-more-btn"
          >
            {{ refreshing ? '加载中...' : '加载更多' }}
          </el-button>
        </div>

        <!-- 已加载全部 -->
        <div v-else-if="applications.length > 0" class="load-complete">
          <p>已加载全部 {{ total }} 条记录</p>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="loading-container">
        <el-icon class="is-loading"><Loading /></el-icon>
        <p>加载中...</p>
      </div>
    </div>

    <!-- 底部操作栏 -->
    <div class="bottom-actions">
      <el-button 
        type="primary" 
        size="large" 
        @click="applyNewLoan"
        class="apply-btn"
      >
        申请新贷款
      </el-button>
    </div>
  </div>
</template>

<style scoped>
.my-applications-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;
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
  flex: 1;
  padding: 16px;
  max-width: 600px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}

.filter-bar {
  margin-bottom: 16px;
  display: flex;
  justify-content: flex-end;
}

.applications-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.application-card {
  background: white;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.application-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.app-id {
  font-size: 14px;
  color: #7f8c8d;
}

.status-tag {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.card-body {
  margin-bottom: 12px;
}

.amount-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.amount-label {
  font-size: 14px;
  color: #7f8c8d;
}

.amount-value {
  font-size: 20px;
  font-weight: 600;
  color: #27ae60;
}

.info-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.info-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.label {
  font-size: 12px;
  color: #7f8c8d;
  margin-bottom: 2px;
}

.value {
  font-size: 14px;
  color: #2c3e50;
  font-weight: 500;
}

.purpose {
  margin-bottom: 8px;
}

.purpose-label {
  font-size: 12px;
  color: #7f8c8d;
}

.purpose-text {
  font-size: 14px;
  color: #2c3e50;
  margin-left: 4px;
}

.approved-amount {
  background: #f0f9ff;
  padding: 8px;
  border-radius: 6px;
  margin-top: 8px;
}

.approved-label {
  font-size: 12px;
  color: #7f8c8d;
}

.approved-value {
  font-size: 16px;
  font-weight: 600;
  color: #27ae60;
  margin-left: 4px;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 8px;
  border-top: 1px solid #f0f0f0;
}

.update-time {
  font-size: 12px;
  color: #7f8c8d;
}

.arrow-icon {
  color: #c0c4cc;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
}

.load-more {
  text-align: center;
  padding: 20px;
}

.load-more-btn {
  color: #27ae60;
}

.load-complete {
  text-align: center;
  padding: 20px;
  color: #7f8c8d;
  font-size: 14px;
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

.bottom-actions {
  padding: 16px;
  background: white;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
}

.apply-btn {
  width: 100%;
  height: 50px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 25px;
}

:deep(.el-empty__description) {
  color: #7f8c8d;
}
</style> 