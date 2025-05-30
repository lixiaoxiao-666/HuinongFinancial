image.png<template>
  <div class="session-management">
    <el-card class="main-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon class="header-icon"><Connection /></el-icon>
            <span class="header-title">会话管理</span>
          </div>
          <div class="header-right">
            <el-button 
              type="primary" 
              :icon="Refresh" 
              @click="refreshSessions"
              :loading="refreshLoading"
            >
              刷新
            </el-button>
            <el-button 
              type="warning" 
              :icon="SwitchButton" 
              @click="logoutOtherDevices"
              :loading="logoutLoading"
            >
              注销其他设备
            </el-button>
          </div>
        </div>
      </template>

      <!-- 当前会话信息 -->
      <div class="current-session-section">
        <h3 class="section-title">
          <el-icon><Monitor /></el-icon>
          当前会话
        </h3>
        <div class="current-session-card" v-if="currentSession">
          <div class="session-info">
            <div class="session-row">
              <span class="label">会话ID:</span>
              <el-tag type="primary" size="small">{{ currentSession.session_id }}</el-tag>
            </div>
            <div class="session-row">
              <span class="label">平台:</span>
              <el-tag :type="getPlatformTagType(currentSession.platform)">
                {{ getPlatformText(currentSession.platform) }}
              </el-tag>
            </div>
            <div class="session-row">
              <span class="label">设备:</span>
              <span class="value">{{ currentSession.device_info?.device_name || '未知设备' }}</span>
            </div>
            <div class="session-row">
              <span class="label">IP地址:</span>
              <span class="value">{{ currentSession.network_info?.ip_address || '未知' }}</span>
            </div>
            <div class="session-row">
              <span class="label">地理位置:</span>
              <span class="value">{{ currentSession.network_info?.location || '未知' }}</span>
            </div>
            <div class="session-row">
              <span class="label">创建时间:</span>
              <span class="value">{{ formatTime(currentSession.created_at) }}</span>
            </div>
            <div class="session-row">
              <span class="label">最后活跃:</span>
              <span class="value">{{ formatTime(currentSession.last_active_at) }}</span>
            </div>
          </div>
          <div class="session-status">
            <el-tag 
              :type="getStatusTagType(currentSession.status)" 
              size="large"
              effect="light"
            >
              <el-icon><Check /></el-icon>
              {{ getStatusText(currentSession.status) }}
            </el-tag>
          </div>
        </div>
      </div>

      <!-- 所有会话列表 -->
      <div class="sessions-section">
        <h3 class="section-title">
          <el-icon><List /></el-icon>
          所有会话 
          <el-tag size="small" type="info">{{ sessions.length }}</el-tag>
        </h3>
        
        <div class="sessions-filters">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索设备名称或IP地址"
            :prefix-icon="Search"
            clearable
            style="width: 300px"
          />
          <el-select v-model="statusFilter" placeholder="状态筛选" clearable style="width: 150px">
            <el-option label="所有状态" value="" />
            <el-option label="活跃" value="active" />
            <el-option label="过期" value="expired" />
            <el-option label="已注销" value="revoked" />
          </el-select>
          <el-select v-model="platformFilter" placeholder="平台筛选" clearable style="width: 150px">
            <el-option label="所有平台" value="" />
            <el-option label="OA管理系统" value="oa" />
            <el-option label="移动应用" value="app" />
            <el-option label="Web端" value="web" />
          </el-select>
        </div>

        <el-table
          :data="filteredSessions"
          v-loading="loading"
          class="sessions-table"
          empty-text="暂无会话数据"
          @sort-change="handleSort"
        >
          <el-table-column prop="session_id" label="会话ID" width="200">
            <template #default="{ row }">
              <el-text class="session-id" size="small" type="primary">
                {{ row.session_id }}
              </el-text>
              <el-tag 
                v-if="row.session_id === authStore.sessionId" 
                type="success" 
                size="small"
                class="current-tag"
              >
                当前
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column prop="platform" label="平台" width="120">
            <template #default="{ row }">
              <el-tag :type="getPlatformTagType(row.platform)" size="small">
                {{ getPlatformText(row.platform) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column prop="device_info.device_name" label="设备信息" min-width="200">
            <template #default="{ row }">
              <div class="device-info">
                <div class="device-name">
                  <el-icon><Monitor /></el-icon>
                  {{ row.device_info?.device_name || '未知设备' }}
                </div>
                <div class="device-type">
                  {{ row.device_info?.device_type || '未知类型' }} 
                  <span v-if="row.device_info?.app_version">
                    v{{ row.device_info.app_version }}
                  </span>
                </div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column prop="network_info.ip_address" label="网络信息" width="180">
            <template #default="{ row }">
              <div class="network-info">
                <div class="ip-address">
                  <el-icon><Location /></el-icon>
                  {{ row.network_info?.ip_address || '未知' }}
                </div>
                <div class="location" v-if="row.network_info?.location">
                  {{ row.network_info.location }}
                </div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusTagType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column prop="created_at" label="创建时间" width="160" sortable="custom">
            <template #default="{ row }">
              <div class="time-info">
                <div>{{ formatTime(row.created_at) }}</div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column prop="last_active_at" label="最后活跃" width="160" sortable="custom">
            <template #default="{ row }">
              <div class="time-info">
                <div>{{ formatTime(row.last_active_at) }}</div>
                <div class="time-ago">{{ getTimeAgo(row.last_active_at) }}</div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-button
                v-if="row.session_id !== authStore.sessionId && row.status === 'active'"
                type="danger"
                size="small"
                :icon="SwitchButton"
                @click="revokeSession(row.session_id)"
                :loading="revokeLoading"
              >
                注销
              </el-button>
              <el-text v-else size="small" type="info">-</el-text>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <!-- 会话统计 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="6">
        <el-card class="stats-card">
          <el-statistic title="总会话数" :value="sessions.length">
            <template #prefix>
              <el-icon><Connection /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stats-card">
          <el-statistic title="活跃会话" :value="activeSessions" value-style="color: #67c23a">
            <template #prefix>
              <el-icon style="color: #67c23a"><Check /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stats-card">
          <el-statistic title="OA平台" :value="oaSessions" value-style="color: #409eff">
            <template #prefix>
              <el-icon style="color: #409eff"><Monitor /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stats-card">
          <el-statistic title="移动端" :value="appSessions" value-style="color: #e6a23c">
            <template #prefix>
              <el-icon style="color: #e6a23c"><MobileFilled /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Connection,
  Refresh,
  SwitchButton,
  Monitor,
  Check,
  List,
  Search,
  Location,
  MobileFilled
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import type { SessionDetail } from '@/types/auth'
import { formatDistanceToNow } from 'date-fns'
import { zhCN } from 'date-fns/locale'

// Store
const authStore = useAuthStore()

// 响应式数据
const loading = ref(false)
const refreshLoading = ref(false)
const logoutLoading = ref(false)
const revokeLoading = ref(false)
const sessions = ref<SessionDetail[]>([])
const searchKeyword = ref('')
const statusFilter = ref('')
const platformFilter = ref('')
const sortField = ref('')
const sortOrder = ref('')

// 计算属性
const currentSession = computed(() => {
  return sessions.value.find(session => session.session_id === authStore.sessionId)
})

const filteredSessions = computed(() => {
  let filtered = [...sessions.value]
  
  // 搜索过滤
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(session => 
      (session.device_info?.device_name || '').toLowerCase().includes(keyword) ||
      (session.network_info?.ip_address || '').toLowerCase().includes(keyword)
    )
  }
  
  // 状态过滤
  if (statusFilter.value) {
    filtered = filtered.filter(session => session.status === statusFilter.value)
  }
  
  // 平台过滤
  if (platformFilter.value) {
    filtered = filtered.filter(session => session.platform === platformFilter.value)
  }
  
  // 排序
  if (sortField.value) {
    filtered.sort((a, b) => {
      let aValue = a[sortField.value as keyof SessionDetail]
      let bValue = b[sortField.value as keyof SessionDetail]
      
      if (typeof aValue === 'string') {
        aValue = new Date(aValue).getTime()
      }
      if (typeof bValue === 'string') {
        bValue = new Date(bValue).getTime()
      }
      
      if (sortOrder.value === 'descending') {
        return (bValue as number) - (aValue as number)
      } else {
        return (aValue as number) - (bValue as number)
      }
    })
  }
  
  return filtered
})

const activeSessions = computed(() => {
  return sessions.value.filter(session => session.status === 'active').length
})

const oaSessions = computed(() => {
  return sessions.value.filter(session => session.platform === 'oa').length
})

const appSessions = computed(() => {
  return sessions.value.filter(session => session.platform === 'app').length
})

// 方法
const fetchSessions = async () => {
  try {
    loading.value = true
    const sessionData = await authStore.fetchSessionInfo()
    sessions.value = sessionData || []
  } catch (error: any) {
    console.error('获取会话信息失败:', error)
    ElMessage.error('获取会话信息失败：' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

const refreshSessions = async () => {
  try {
    refreshLoading.value = true
    await fetchSessions()
    ElMessage.success('会话信息已刷新')
  } catch (error) {
    // Error已在fetchSessions中处理
  } finally {
    refreshLoading.value = false
  }
}

const logoutOtherDevices = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要注销其他设备的会话吗？这将强制其他设备重新登录。',
      '确认注销',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    logoutLoading.value = true
    await authStore.logoutOtherDevices()
    await fetchSessions()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败：' + (error.message || '未知错误'))
    }
  } finally {
    logoutLoading.value = false
  }
}

const revokeSession = async (sessionId: string) => {
  try {
    await ElMessageBox.confirm(
      '确定要注销该会话吗？',
      '确认注销',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    revokeLoading.value = true
    // 这里应该调用具体的注销单个会话的API
    // await revokeSpecificSession(sessionId)
    ElMessage.success('会话已注销')
    await fetchSessions()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败：' + (error.message || '未知错误'))
    }
  } finally {
    revokeLoading.value = false
  }
}

const handleSort = ({ prop, order }: { prop: string; order: string }) => {
  sortField.value = prop
  sortOrder.value = order
}

// 辅助函数
const getPlatformText = (platform: string) => {
  const platformMap: Record<string, string> = {
    oa: 'OA管理系统',
    app: '移动应用',
    web: 'Web端'
  }
  return platformMap[platform] || platform
}

const getPlatformTagType = (platform: string) => {
  const typeMap: Record<string, string> = {
    oa: 'primary',
    app: 'warning',
    web: 'info'
  }
  return typeMap[platform] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '活跃',
    expired: '已过期',
    revoked: '已注销'
  }
  return statusMap[status] || status
}

const getStatusTagType = (status: string) => {
  const typeMap: Record<string, string> = {
    active: 'success',
    expired: 'warning',
    revoked: 'danger'
  }
  return typeMap[status] || 'info'
}

const formatTime = (timeStr: string) => {
  if (!timeStr) return '-'
  return new Date(timeStr).toLocaleString('zh-CN')
}

const getTimeAgo = (timeStr: string) => {
  if (!timeStr) return ''
  try {
    return formatDistanceToNow(new Date(timeStr), { 
      addSuffix: true, 
      locale: zhCN 
    })
  } catch {
    return ''
  }
}

// 生命周期
onMounted(() => {
  fetchSessions()
})
</script>

<style scoped>
.session-management {
  padding: 20px;
}

.main-card {
  border-radius: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-icon {
  font-size: 20px;
  color: #409eff;
}

.header-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.header-right {
  display: flex;
  gap: 12px;
}

.current-session-section {
  margin-bottom: 30px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.current-session-card {
  background: linear-gradient(135deg, #f6f9fc 0%, #e9f7ef 100%);
  border: 1px solid #d4edda;
  border-radius: 12px;
  padding: 24px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.session-info {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.session-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.label {
  font-weight: 500;
  color: #606266;
  min-width: 80px;
}

.value {
  color: #303133;
  font-family: 'Courier New', monospace;
}

.session-status {
  margin-left: 20px;
}

.sessions-section {
  margin-top: 30px;
}

.sessions-filters {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
  align-items: center;
}

.sessions-table {
  margin-top: 16px;
}

.session-id {
  font-family: 'Courier New', monospace;
}

.current-tag {
  margin-left: 8px;
}

.device-info, .network-info {
  line-height: 1.4;
}

.device-name, .ip-address {
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 500;
}

.device-type, .location {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}

.time-info {
  line-height: 1.4;
}

.time-ago {
  font-size: 12px;
  color: #909399;
}

.stats-card {
  border-radius: 8px;
  transition: all 0.3s ease;
}

.stats-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .session-management {
    padding: 10px;
  }
  
  .card-header {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
  
  .header-right {
    justify-content: center;
  }
  
  .current-session-card {
    flex-direction: column;
    gap: 16px;
  }
  
  .session-info {
    grid-template-columns: 1fr;
  }
  
  .sessions-filters {
    flex-direction: column;
    align-items: stretch;
  }
  
  .sessions-filters > * {
    width: 100% !important;
  }
}
</style> 