<template>
  <div class="dashboard">
    <div class="page-header">
      <h2 class="page-title">工作台</h2>
      <div class="header-actions">
        <el-button @click="fetchDashboardData" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新数据
        </el-button>
      </div>
    </div>
    
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon total">
            <el-icon><Files /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ dashboardData.pending_count + dashboardData.approved_count + dashboardData.rejected_count }}</div>
            <div class="stat-label">总申请数</div>
          </div>
        </div>
      </el-card>
      
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon pending">
            <el-icon><Clock /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ dashboardData.pending_count }}</div>
            <div class="stat-label">待处理</div>
          </div>
        </div>
      </el-card>
      
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon approved">
            <el-icon><CircleCheck /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ dashboardData.approved_count }}</div>
            <div class="stat-label">已批准</div>
          </div>
        </div>
      </el-card>
      
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon ai">
            <el-icon><Cpu /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ dashboardData.ai_processing_count }}</div>
            <div class="stat-label">AI处理中</div>
          </div>
        </div>
      </el-card>
    </div>
    
    <!-- 主要内容区域 -->
    <el-row :gutter="20" class="main-content">
      <!-- 左侧列 -->
      <el-col :xs="24" :sm="24" :md="16" :lg="16" :xl="16">
        <!-- 待办事项 -->
        <el-card class="section-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Bell /></el-icon>
                我的待办
              </span>
              <el-badge :value="pendingTasks.length" class="badge">
                <el-button type="primary" size="small" @click="goToApproval">
                  查看全部
                </el-button>
              </el-badge>
            </div>
          </template>
          
          <div v-if="pendingTasks.length === 0" class="empty-state">
            <el-empty description="暂无待办事项" />
          </div>
          
          <div v-else class="task-list">
            <div
              v-for="task in pendingTasks"
              :key="task.task_id"
              class="task-item"
              @click="handleTaskClick(task)"
            >
              <div class="task-priority" :class="task.priority"></div>
              <div class="task-content">
                <div class="task-title">{{ task.title }}</div>
                <div class="task-desc">{{ task.task_type }}</div>
                <div class="task-meta">
                  <span class="task-time">{{ formatRelativeTime(task.created_at) }}</span>
                  <el-tag size="small" :type="getPriorityType(task.priority)">
                    {{ task.priority === 'high' ? '高优先级' : '普通' }}
                  </el-tag>
                </div>
              </div>
              <div class="task-action">
                <el-icon><ArrowRight /></el-icon>
              </div>
            </div>
          </div>
        </el-card>
        
        <!-- AI审批统计 -->
        <el-card class="section-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Cpu /></el-icon>
                AI审批统计
              </span>
            </div>
          </template>
          
          <div class="ai-stats">
            <div class="ai-stat-item">
              <div class="ai-stat-content">
                <div class="ai-stat-label">AI处理率</div>
                <div class="ai-stat-value">{{ dashboardData.ai_processing_rate }}%</div>
              </div>
              <el-progress 
                :percentage="dashboardData.ai_processing_rate" 
                :stroke-width="8"
                :show-text="false"
                color="#67c23a"
              />
            </div>
            <div class="ai-stat-item">
              <div class="ai-stat-content">
                <div class="ai-stat-label">AI审批状态</div>
                <div class="ai-stat-value">
                  <el-switch
                    v-model="dashboardData.ai_enabled"
                    @change="handleAIToggle"
                    :disabled="!hasPermission('system:manage')"
                  />
                  <span class="ai-status-text">{{ dashboardData.ai_enabled ? '运行中' : '已关闭' }}</span>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <!-- 右侧列 -->
      <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
        <!-- 快捷操作 -->
        <el-card class="section-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Operation /></el-icon>
                快捷操作
              </span>
            </div>
          </template>
          
          <div class="quick-actions">
            <div class="quick-action-item" @click="goToPage('/approval')">
              <el-icon class="action-icon">
                <DocumentChecked />
              </el-icon>
              <span>审批管理</span>
            </div>
            <div class="quick-action-item" @click="goToPage('/users')">
              <el-icon class="action-icon">
                <User />
              </el-icon>
              <span>用户管理</span>
            </div>
            <div class="quick-action-item" @click="goToPage('/logs')">
              <el-icon class="action-icon">
                <Document />
              </el-icon>
              <span>操作日志</span>
            </div>
            <div class="quick-action-item" @click="goToPage('/system')">
              <el-icon class="action-icon">
                <Setting />
              </el-icon>
              <span>系统设置</span>
            </div>
          </div>
        </el-card>
        
        <!-- 最近活动 -->
        <el-card class="section-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Document /></el-icon>
                最近活动
              </span>
            </div>
          </template>
          
          <div class="activity-list">
            <div
              v-for="(activity, index) in recentActivities"
              :key="index"
              class="activity-item"
            >
              <div class="activity-time">{{ formatTime(activity.timestamp) }}</div>
              <div class="activity-content">
                <div class="activity-details">
                  <div class="activity-action">{{ activity.activity_type }}</div>
                  <div class="activity-target">{{ activity.description }}</div>
                </div>
                <div class="activity-result">
                  <el-tag
                    size="small"
                    type="success"
                  >
                    {{ activity.operator }}
                  </el-tag>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Files,
  Clock,
  CircleCheck,
  Cpu,
  Bell,
  Operation,
  Document,
  ArrowRight,
  DocumentChecked,
  Setting,
  User,
  Refresh
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { getDashboard, toggleAIApproval } from '@/api/admin'
import type { DashboardData } from '@/types'
import dayjs from 'dayjs'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const dashboardData = ref<DashboardData>({
  pending_count: 0,
  approved_count: 0,
  rejected_count: 0,
  ai_processing_count: 0,
  ai_enabled: true,
  ai_processing_rate: 85,
  pending_tasks: [],
  recent_activities: []
})

// 计算属性
const pendingTasks = computed(() => dashboardData.value.pending_tasks)
const recentActivities = computed(() => dashboardData.value.recent_activities)
const hasPermission = computed(() => authStore.hasPermission)

// 方法
const fetchDashboardData = async () => {
  try {
    loading.value = true
    const data = await getDashboard()
    dashboardData.value = data
  } catch (error) {
    ElMessage.error('获取工作台数据失败')
  } finally {
    loading.value = false
  }
}

const handleAIToggle = async (enabled: boolean) => {
  try {
    await toggleAIApproval(enabled)
    ElMessage.success(`AI审批功能已${enabled ? '开启' : '关闭'}`)
  } catch (error) {
    ElMessage.error('切换AI审批状态失败')
    // 回滚状态
    dashboardData.value.ai_enabled = !enabled
  }
}

const handleTaskClick = (task: any) => {
  if (task.application_id) {
    router.push(`/approval/${task.application_id}`)
  } else {
    router.push('/approval')
  }
}

const goToApproval = () => {
  router.push('/approval')
}

const goToPage = (path: string) => {
  router.push(path)
}

const getPriorityType = (priority: string) => {
  return priority === 'high' ? 'danger' : 'info'
}

const formatRelativeTime = (datetime: string) => {
  const now = dayjs()
  const time = dayjs(datetime)
  const hours = now.diff(time, 'hour')
  
  if (hours < 1) {
    const minutes = now.diff(time, 'minute')
    return `${minutes}分钟前`
  } else if (hours < 24) {
    return `${hours}小时前`
  } else {
    const days = now.diff(time, 'day')
    return `${days}天前`
  }
}

const formatTime = (time: string) => {
  return dayjs(time).format('HH:mm')
}

onMounted(() => {
  fetchDashboardData()
})
</script>

<style scoped>
.dashboard {
  padding: 0;
  min-height: 100%;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 0 20px 0;
  margin-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
}

.page-title {
  margin: 0;
  color: #333;
  font-size: 24px;
  font-weight: 600;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 16px;
  }
}

.stat-card {
  border-radius: 12px;
  border: none;
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.08);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
  flex-shrink: 0;
}

.stat-icon.total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.pending {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.approved {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.ai {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 28px;
  font-weight: 700;
  color: #333;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  color: #666;
  font-size: 14px;
}

.main-content {
  margin-top: 0;
}

@media (max-width: 768px) {
  .main-content {
    gap: 16px;
  }
}

.section-card {
  margin-bottom: 20px;
  border-radius: 12px;
  border: 1px solid #ebeef5;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
  color: #333;
  gap: 8px;
}

.badge {
  margin-right: 0;
}

.empty-state {
  padding: 40px 0;
}

.task-list {
  max-height: 400px;
  overflow-y: auto;
}

.task-item {
  display: flex;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid #f0f2f5;
  cursor: pointer;
  transition: all 0.3s;
}

.task-item:hover {
  background-color: #f8f9fa;
  margin: 0 -20px;
  padding-left: 20px;
  padding-right: 20px;
  border-radius: 8px;
}

.task-item:last-child {
  border-bottom: none;
}

.task-priority {
  width: 4px;
  height: 40px;
  border-radius: 2px;
  margin-right: 12px;
  background-color: #e4e7ed;
  flex-shrink: 0;
}

.task-priority.high {
  background-color: #f56c6c;
}

.task-priority.medium,
.task-priority.normal {
  background-color: #409eff;
}

.task-priority.low {
  background-color: #67c23a;
}

.task-content {
  flex: 1;
  min-width: 0;
}

.task-title {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.task-desc {
  color: #666;
  font-size: 13px;
  margin-bottom: 8px;
}

.task-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.task-time {
  color: #999;
  font-size: 12px;
}

.task-action {
  color: #c0c4cc;
  font-size: 16px;
  flex-shrink: 0;
}

.ai-stats {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.ai-stat-item {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ai-stat-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.ai-stat-label {
  color: #666;
  font-size: 14px;
}

.ai-stat-value {
  font-weight: 600;
  color: #333;
  display: flex;
  align-items: center;
}

.ai-status-text {
  margin-left: 12px;
  color: #666;
  font-size: 12px;
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

@media (max-width: 768px) {
  .quick-actions {
    grid-template-columns: 1fr;
  }
}

.quick-action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 12px;
  background: #f8f9fa;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  gap: 8px;
  text-align: center;
}

.quick-action-item:hover {
  background: #409eff;
  color: white;
  transform: translateY(-2px);
}

.action-icon {
  font-size: 24px;
}

.activity-list {
  max-height: 300px;
  overflow-y: auto;
}

.activity-item {
  display: flex;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid #f0f2f5;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-time {
  color: #999;
  font-size: 12px;
  white-space: nowrap;
  min-width: 40px;
  flex-shrink: 0;
}

.activity-content {
  flex: 1;
  min-width: 0;
}

.activity-details {
  display: flex;
  flex-direction: column;
  margin-bottom: 8px;
}

.activity-action {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.activity-target {
  color: #666;
  font-size: 13px;
  word-break: break-all;
}

.activity-result {
  display: flex;
  justify-content: flex-end;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }
  
  .header-actions {
    justify-content: center;
  }
  
  .task-content {
    margin-right: 8px;
  }
  
  .task-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
  
  .ai-stat-content {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .activity-content {
    margin-right: 8px;
  }
}

/* 进度条样式 */
:deep(.el-progress-bar__outer) {
  background-color: #e4e7ed;
  border-radius: 4px;
}

:deep(.el-progress-bar__inner) {
  border-radius: 4px;
}

/* 卡片标题样式 */
:deep(.el-card__header) {
  padding: 18px 20px;
  border-bottom: 1px solid #f0f2f5;
}

:deep(.el-card__body) {
  padding: 20px;
}

/* 标签样式 */
:deep(.el-tag) {
  border: none;
}

/* 徽章样式 */
:deep(.el-badge__content) {
  border: 1px solid #fff;
}
</style> 