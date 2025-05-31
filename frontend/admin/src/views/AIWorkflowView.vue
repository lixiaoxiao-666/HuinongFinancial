<template>
  <div class="ai-workflow-container">
    <el-card class="header-card">
      <template #header>
        <div class="card-header">
          <span>AI审批流程管理</span>
          <div class="header-actions">
            <el-button type="primary" @click="createWorkflow">
              <el-icon><Plus /></el-icon>
              管理Dify AI工作流
            </el-button>
            <el-button @click="refreshData">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>
      
      <div class="workflow-stats">
        <el-row :gutter="20">
          <el-col :span="6">
            <div class="stat-card">
              <div class="stat-icon active">
                <el-icon><Cpu /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">{{ stats.activeWorkflows }}</div>
                <div class="stat-label">活跃流程</div>
              </div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-card">
              <div class="stat-icon total">
                <el-icon><Document /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">{{ stats.totalWorkflows }}</div>
                <div class="stat-label">总流程数</div>
              </div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-card">
              <div class="stat-icon processing">
                <el-icon><Loading /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">{{ stats.processingTasks }}</div>
                <div class="stat-label">处理中任务</div>
              </div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-card">
              <div class="stat-icon completed">
                <el-icon><Select /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">{{ stats.completedToday }}</div>
                <div class="stat-label">今日完成</div>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </el-card>

    <el-card class="content-card">
      <template #header>
        <div class="card-header">
          <span>审批流程列表</span>
          <div class="header-filters">
            <el-select v-model="selectedStatus" placeholder="流程状态" style="width: 120px" @change="handleStatusChange">
              <el-option label="全部" value="" />
              <el-option label="活跃" value="active" />
              <el-option label="暂停" value="paused" />
              <el-option label="已停止" value="stopped" />
            </el-select>
            <el-select v-model="selectedCategory" placeholder="流程类型" style="width: 140px" @change="handleCategoryChange">
              <el-option label="全部类型" value="" />
              <el-option label="贷款审批" value="loan" />
              <el-option label="补贴申请" value="subsidy" />
              <el-option label="保险理赔" value="insurance" />
              <el-option label="其他" value="other" />
            </el-select>
          </div>
        </div>
      </template>

      <el-table
        :data="workflowData"
        v-loading="loading"
        stripe
        style="width: 100%"
        @sort-change="handleSortChange"
      >
        <el-table-column prop="id" label="流程ID" width="100" />
        <el-table-column prop="name" label="流程名称" min-width="200" />
        <el-table-column prop="category" label="类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getCategoryTagType(row.category)">
              {{ getCategoryText(row.category) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag
              :type="getStatusTagType(row.status)"
              effect="dark"
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="version" label="版本" width="80" />
        <el-table-column prop="ai_confidence_threshold" label="AI置信度阈值" width="120">
          <template #default="{ row }">
            <span class="confidence-threshold">{{ row.ai_confidence_threshold }}%</span>
          </template>
        </el-table-column>
        <el-table-column prop="processed_count" label="已处理" width="100" sortable="custom">
          <template #default="{ row }">
            <span class="processed-count">{{ row.processed_count }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="success_rate" label="成功率" width="100" sortable="custom">
          <template #default="{ row }">
            <el-progress
              :percentage="row.success_rate"
              :color="getSuccessRateColor(row.success_rate)"
              :stroke-width="6"
              :show-text="false"
            />
            <span class="success-rate">{{ row.success_rate }}%</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160" sortable="custom">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间" width="160" sortable="custom">
          <template #default="{ row }">
            {{ formatTime(row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click="viewWorkflow(row)"
            >
              查看
            </el-button>
            <el-button
              type="warning"
              size="small"
              @click="editWorkflow(row)"
            >
              编辑
            </el-button>
            <el-button
              v-if="row.status === 'active'"
              type="info"
              size="small"
              @click="pauseWorkflow(row)"
            >
              暂停
            </el-button>
            <el-button
              v-else-if="row.status === 'paused'"
              type="success"
              size="small"
              @click="resumeWorkflow(row)"
            >
              恢复
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="deleteWorkflow(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 流程详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="审批流程详情"
      width="1000px"
      destroy-on-close
    >
      <div v-if="selectedWorkflow" class="workflow-detail">
        <el-descriptions :column="3" border style="margin-bottom: 20px">
          <el-descriptions-item label="流程ID">{{ selectedWorkflow.id }}</el-descriptions-item>
          <el-descriptions-item label="流程名称">{{ selectedWorkflow.name }}</el-descriptions-item>
          <el-descriptions-item label="类型">{{ getCategoryText(selectedWorkflow.category) }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusTagType(selectedWorkflow.status)">
              {{ getStatusText(selectedWorkflow.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="版本">{{ selectedWorkflow.version }}</el-descriptions-item>
          <el-descriptions-item label="AI置信度阈值">{{ selectedWorkflow.ai_confidence_threshold }}%</el-descriptions-item>
          <el-descriptions-item label="已处理任务">{{ selectedWorkflow.processed_count }}</el-descriptions-item>
          <el-descriptions-item label="成功率">{{ selectedWorkflow.success_rate }}%</el-descriptions-item>
          <el-descriptions-item label="平均处理时间">{{ selectedWorkflow.avg_processing_time }}分钟</el-descriptions-item>
        </el-descriptions>

        <div class="workflow-steps">
          <h4>审批流程步骤</h4>
          <el-steps :active="selectedWorkflow.steps?.length || 0" finish-status="success" align-center>
            <el-step
              v-for="(step, index) in selectedWorkflow.steps"
              :key="index"
              :title="step.name"
              :description="step.description"
            />
          </el-steps>
        </div>

        <div class="workflow-config" v-if="selectedWorkflow.config">
          <h4>配置信息</h4>
          <el-card>
            <pre>{{ JSON.stringify(selectedWorkflow.config, null, 2) }}</pre>
          </el-card>
        </div>
      </div>
    </el-dialog>

    <!-- 创建/编辑流程对话框 -->
    <el-dialog
      v-model="formDialogVisible"
      :title="isEdit ? '编辑审批流程' : '创建审批流程'"
      width="800px"
      destroy-on-close
    >
      <el-form
        ref="workflowFormRef"
        :model="workflowForm"
        :rules="workflowRules"
        label-width="120px"
      >
        <el-form-item label="流程名称" prop="name">
          <el-input v-model="workflowForm.name" placeholder="请输入流程名称" />
        </el-form-item>
        <el-form-item label="流程类型" prop="category">
          <el-select v-model="workflowForm.category" placeholder="请选择流程类型">
            <el-option label="贷款审批" value="loan" />
            <el-option label="补贴申请" value="subsidy" />
            <el-option label="保险理赔" value="insurance" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="AI置信度阈值" prop="ai_confidence_threshold">
          <el-slider
            v-model="workflowForm.ai_confidence_threshold"
            :min="50"
            :max="95"
            :step="5"
            show-stops
            show-input
          />
        </el-form-item>
        <el-form-item label="流程描述" prop="description">
          <el-input
            v-model="workflowForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入流程描述"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="workflowForm.status">
            <el-radio value="active">活跃</el-radio>
            <el-radio value="paused">暂停</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="formDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitWorkflow">
            {{ isEdit ? '更新' : '创建' }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Refresh,
  Cpu,
  Document,
  Loading,
  Select
} from '@element-plus/icons-vue'
import dayjs from 'dayjs'

// 接口定义
interface WorkflowStep {
  name: string
  description: string
}

interface WorkflowRecord {
  id: string
  name: string
  category: string
  status: string
  version: string
  ai_confidence_threshold: number
  processed_count: number
  success_rate: number
  avg_processing_time: number
  created_at: string
  updated_at: string
  description?: string
  steps?: WorkflowStep[]
  config?: any
}

interface WorkflowForm {
  name: string
  category: string
  ai_confidence_threshold: number
  description: string
  status: string
}

// 响应式数据
const loading = ref(false)
const workflowData = ref<WorkflowRecord[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const selectedStatus = ref('')
const selectedCategory = ref('')
const detailDialogVisible = ref(false)
const formDialogVisible = ref(false)
const selectedWorkflow = ref<WorkflowRecord | null>(null)
const isEdit = ref(false)

// 统计数据
const stats = reactive({
  activeWorkflows: 0,
  totalWorkflows: 0,
  processingTasks: 0,
  completedToday: 0
})

// 表单数据
const workflowForm = reactive<WorkflowForm>({
  name: '',
  category: '',
  ai_confidence_threshold: 80,
  description: '',
  status: 'active'
})

// 表单验证规则
const workflowRules = {
  name: [
    { required: true, message: '请输入流程名称', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请选择流程类型', trigger: 'change' }
  ],
  ai_confidence_threshold: [
    { required: true, message: '请设置AI置信度阈值', trigger: 'blur' }
  ]
}

// 方法
const refreshData = () => {
  fetchWorkflows()
  fetchStats()
}

const fetchWorkflows = async () => {
  loading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // 模拟数据
    const mockData = Array.from({ length: pageSize.value }, (_, index) => ({
      id: `WF${String(currentPage.value * 100 + index + 1).padStart(4, '0')}`,
      name: `${['贷款审批', '补贴申请', '保险理赔'][index % 3]}流程${index + 1}`,
      category: ['loan', 'subsidy', 'insurance'][index % 3],
      status: ['active', 'paused', 'stopped'][Math.floor(Math.random() * 3)],
      version: `v${Math.floor(Math.random() * 3) + 1}.${Math.floor(Math.random() * 10)}`,
      ai_confidence_threshold: [75, 80, 85, 90][Math.floor(Math.random() * 4)],
      processed_count: Math.floor(Math.random() * 5000) + 100,
      success_rate: Math.floor(Math.random() * 30) + 70,
      avg_processing_time: Math.floor(Math.random() * 30) + 5,
      created_at: new Date(Date.now() - Math.random() * 365 * 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
      description: '这是一个智能审批流程的描述信息',
      steps: [
        { name: '材料收集', description: '收集申请材料' },
        { name: 'AI初审', description: 'AI自动初步审核' },
        { name: '人工复审', description: '人工复审确认' },
        { name: '结果通知', description: '通知审批结果' }
      ],
      config: {
        maxAmount: 100000,
        minCreditScore: 600,
        autoApprovalThreshold: 85
      }
    }))
    
    workflowData.value = mockData
    total.value = 150
  } catch (error) {
    ElMessage.error('获取流程数据失败')
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  try {
    // 模拟统计数据
    stats.activeWorkflows = 25
    stats.totalWorkflows = 38
    stats.processingTasks = 156
    stats.completedToday = 89
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const handleStatusChange = () => {
  currentPage.value = 1
  fetchWorkflows()
}

const handleCategoryChange = () => {
  currentPage.value = 1
  fetchWorkflows()
}

const handleSortChange = (sortInfo: any) => {
  console.log('排序变化:', sortInfo)
  fetchWorkflows()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchWorkflows()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchWorkflows()
}

const viewWorkflow = (row: WorkflowRecord) => {
  selectedWorkflow.value = row
  detailDialogVisible.value = true
}

const createWorkflow = () => {
  // 跳转到Dify AI工作流管理页面
  window.open('http://10.10.10.5/app/c827412a-c7dd-412a-9c3f-333ef5bb70d7/workflow', '_blank')
}

const editWorkflow = (row: WorkflowRecord) => {
  isEdit.value = true
  workflowForm.name = row.name
  workflowForm.category = row.category
  workflowForm.ai_confidence_threshold = row.ai_confidence_threshold
  workflowForm.description = row.description || ''
  workflowForm.status = row.status
  selectedWorkflow.value = row
  formDialogVisible.value = true
}

const pauseWorkflow = async (row: WorkflowRecord) => {
  try {
    await ElMessageBox.confirm('确定要暂停此审批流程吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    ElMessage.success('流程已暂停')
    refreshData()
  } catch {
    // 用户取消
  }
}

const resumeWorkflow = async (row: WorkflowRecord) => {
  ElMessage.success('流程已恢复')
  refreshData()
}

const deleteWorkflow = async (row: WorkflowRecord) => {
  try {
    await ElMessageBox.confirm('确定要删除此审批流程吗？删除后无法恢复！', '警告', {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'error'
    })
    ElMessage.success('流程已删除')
    refreshData()
  } catch {
    // 用户取消
  }
}

const submitWorkflow = () => {
  // 这里应该调用表单验证和提交逻辑
  ElMessage.success(isEdit.value ? '流程更新成功' : '流程创建成功')
  formDialogVisible.value = false
  refreshData()
}

const resetForm = () => {
  workflowForm.name = ''
  workflowForm.category = ''
  workflowForm.ai_confidence_threshold = 80
  workflowForm.description = ''
  workflowForm.status = 'active'
}

// 工具方法
const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

const getStatusTagType = (status: string) => {
  switch (status) {
    case 'active': return 'success'
    case 'paused': return 'warning'
    case 'stopped': return 'danger'
    default: return 'info'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'active': return '活跃'
    case 'paused': return '暂停'
    case 'stopped': return '已停止'
    default: return '未知'
  }
}

const getCategoryTagType = (category: string) => {
  switch (category) {
    case 'loan': return 'primary'
    case 'subsidy': return 'success'
    case 'insurance': return 'warning'
    default: return 'info'
  }
}

const getCategoryText = (category: string) => {
  switch (category) {
    case 'loan': return '贷款审批'
    case 'subsidy': return '补贴申请'
    case 'insurance': return '保险理赔'
    case 'other': return '其他'
    default: return '未知'
  }
}

const getSuccessRateColor = (rate: number) => {
  if (rate >= 90) return '#67c23a'
  if (rate >= 80) return '#e6a23c'
  return '#f56c6c'
}

onMounted(() => {
  refreshData()
})
</script>

<style scoped>
.ai-workflow-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 0 20px;
  overflow: hidden;
}

.header-card {
  flex-shrink: 0;
}

.content-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.header-filters {
  display: flex;
  gap: 12px;
  align-items: center;
}

.workflow-stats {
  margin-top: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  font-size: 24px;
  color: #fff;
}

.stat-icon.active {
  background: linear-gradient(135deg, #67c23a, #85ce61);
}

.stat-icon.total {
  background: linear-gradient(135deg, #409eff, #66b1ff);
}

.stat-icon.processing {
  background: linear-gradient(135deg, #e6a23c, #ebb563);
}

.stat-icon.completed {
  background: linear-gradient(135deg, #f56c6c, #f78989);
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #333;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.confidence-threshold {
  font-weight: 600;
  color: #e6a23c;
}

.processed-count {
  font-weight: 600;
  color: #409eff;
}

.success-rate {
  margin-left: 8px;
  font-size: 12px;
  color: #666;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  flex-shrink: 0;
}

.workflow-detail {
  max-height: 600px;
  overflow-y: auto;
  padding-right: 10px;
  /* 详情对话框滚动条样式 */
  scrollbar-width: thin;
  scrollbar-color: #e6a23c #f1f1f1;
}

.workflow-detail::-webkit-scrollbar {
  width: 12px;
}

.workflow-detail::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 6px;
}

.workflow-detail::-webkit-scrollbar-thumb {
  background: #e6a23c;
  border-radius: 6px;
  border: 2px solid #f1f1f1;
}

.workflow-detail::-webkit-scrollbar-thumb:hover {
  background: #cf9236;
}

.workflow-steps {
  margin: 20px 0;
}

.workflow-steps h4 {
  margin-bottom: 16px;
  color: #333;
}

.workflow-config {
  margin-top: 20px;
}

.workflow-config h4 {
  margin-bottom: 12px;
  color: #333;
}

.workflow-config pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.5;
  color: #555;
}

/* Element Plus卡片内容区域样式 */
:deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 20px;
}

/* 表格容器样式 */
:deep(.el-table) {
  flex: 1;
  overflow: hidden;
}

/* 表格主体滚动区域 */
:deep(.el-table__body-wrapper) {
  max-height: calc(100vh - 400px);
  overflow-y: auto !important;
  overflow-x: auto !important;
  /* 表格滚动条样式 */
  scrollbar-width: thin;
  scrollbar-color: #e6a23c #f1f1f1;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar) {
  width: 12px;
  height: 12px;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar-track) {
  background: #f1f1f1;
  border-radius: 6px;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar-thumb) {
  background: #e6a23c;
  border-radius: 6px;
  border: 2px solid #f1f1f1;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar-thumb:hover) {
  background: #cf9236;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar-corner) {
  background: #f1f1f1;
}

/* 表格头部固定 */
:deep(.el-table__header-wrapper) {
  overflow: visible;
}

:deep(.el-steps) {
  margin: 20px 0;
}

/* 分页组件位置调整 */
:deep(.el-pagination) {
  margin-top: 20px;
  justify-content: center;
}
</style> 