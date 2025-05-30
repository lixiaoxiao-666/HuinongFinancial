<template>
  <div class="loan-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">贷款申请管理</h1>
        <div class="header-actions">
          <a-space>
            <a-button type="primary" @click="exportData">
              <template #icon>
                <DownloadOutlined />
              </template>
              导出数据
            </a-button>
            <a-button @click="refreshData">
              <template #icon>
                <ReloadOutlined />
              </template>
              刷新
            </a-button>
          </a-space>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="statistics-section">
      <a-row :gutter="[24, 24]">
        <a-col :xs="12" :sm="6" v-for="stat in statistics" :key="stat.key">
          <div class="stat-card" :class="stat.type">
            <div class="stat-header">
              <div class="stat-icon" :style="{ background: stat.color }">
                <component :is="getIconComponent(stat.icon)" />
              </div>
              <div class="stat-trend" :class="stat.trend?.type">
                <component :is="getIconComponent(stat.trend?.icon)" v-if="stat.trend?.icon" />
                <span>{{ stat.trend?.value }}</span>
              </div>
            </div>
            <div class="stat-content">
              <div class="stat-value">
                {{ formatStatValue(stat.value, stat.precision) }}<span class="stat-suffix">{{ stat.suffix }}</span>
              </div>
              <div class="stat-title">{{ stat.title }}</div>
            </div>
          </div>
        </a-col>
      </a-row>
    </div>

    <!-- 筛选区域 -->
    <div class="filter-section">
      <a-card :bodyStyle="{ padding: '16px' }">
        <a-form :model="filters" layout="inline" @finish="handleSearch">
          <a-form-item label="申请状态">
            <a-select 
              v-model:value="filters.status" 
              placeholder="请选择状态" 
              style="width: 150px"
              allowClear
            >
              <a-select-option value="pending">待审批</a-select-option>
              <a-select-option value="under_review">审核中</a-select-option>
              <a-select-option value="approved">已批准</a-select-option>
              <a-select-option value="rejected">已拒绝</a-select-option>
              <a-select-option value="returned">已退回</a-select-option>
            </a-select>
          </a-form-item>

          <a-form-item label="产品类型">
            <a-select 
              v-model:value="filters.product_id" 
              placeholder="请选择产品" 
              style="width: 180px"
              allowClear
            >
              <a-select-option v-for="product in loanProducts" :key="product.id" :value="product.id">
                {{ product.product_name }}
              </a-select-option>
            </a-select>
          </a-form-item>

          <a-form-item label="风险等级">
            <a-select 
              v-model:value="filters.risk_level" 
              placeholder="请选择风险等级" 
              style="width: 150px"
              allowClear
            >
              <a-select-option value="low">低风险</a-select-option>
              <a-select-option value="medium">中风险</a-select-option>
              <a-select-option value="high">高风险</a-select-option>
            </a-select>
          </a-form-item>

          <a-form-item label="申请日期">
            <a-range-picker 
              v-model:value="filters.dateRange" 
              format="YYYY-MM-DD"
              style="width: 240px"
            />
          </a-form-item>

          <a-form-item label="申请人">
            <a-input 
              v-model:value="filters.applicant" 
              placeholder="姓名/手机号" 
              style="width: 150px"
              allowClear
            />
          </a-form-item>

          <a-form-item>
            <a-space>
              <a-button type="primary" html-type="submit" :loading="loading">
                <template #icon>
                  <SearchOutlined />
                </template>
                查询
              </a-button>
              <a-button @click="resetFilters">
                重置
              </a-button>
            </a-space>
          </a-form-item>
        </a-form>
      </a-card>
    </div>

    <!-- 申请列表 -->
    <div class="table-section">
      <a-card>
        <template #title>
          <div class="table-header">
            <span>贷款申请列表</span>
            <div class="table-actions">
              <a-space>
                <a-button 
                  type="primary" 
                  :disabled="!hasSelectedItems"
                  @click="batchApprove"
                >
                  批量审批
                </a-button>
                <a-dropdown>
                  <template #overlay>
                    <a-menu @click="handleBatchAction">
                      <a-menu-item key="approve">批量通过</a-menu-item>
                      <a-menu-item key="reject">批量拒绝</a-menu-item>
                      <a-menu-item key="assign">批量分配</a-menu-item>
                    </a-menu>
                  </template>
                  <a-button>
                    批量操作
                    <DownOutlined />
                  </a-button>
                </a-dropdown>
              </a-space>
            </div>
          </div>
        </template>

        <a-table
          :columns="columns"
          :data-source="applicationList"
          :loading="loading"
          :pagination="pagination"
          :row-selection="rowSelection"
          :scroll="{ x: 1200 }"
          @change="handleTableChange"
          row-key="id"
          size="middle"
        >
          <!-- 申请人信息 -->
          <template #applicantInfo="{ record }">
            <div class="applicant-info">
              <a-avatar :src="record.user_info.avatar" :size="32">
                <template #icon>
                  <UserOutlined />
                </template>
              </a-avatar>
              <div class="applicant-details">
                <div class="applicant-name">{{ record.user_info.real_name }}</div>
                <div class="applicant-phone">{{ record.user_info.phone }}</div>
              </div>
            </div>
          </template>

          <!-- 申请信息 -->
          <template #applicationInfo="{ record }">
            <div class="application-info">
              <div class="application-id">{{ record.id }}</div>
              <div class="product-name">{{ record.product_name }}</div>
              <div class="application-amount">
                ¥{{ formatCurrency(record.amount) }}
              </div>
            </div>
          </template>

          <!-- 状态 -->
          <template #status="{ record }">
            <a-tag :color="getStatusColor(record.status)">
              {{ getStatusText(record.status) }}
            </a-tag>
          </template>

          <!-- AI建议 -->
          <template #aiSuggestion="{ record }">
            <div class="ai-suggestion" v-if="record.ai_assessment">
              <a-tag 
                :color="getAiSuggestionColor(record.ai_assessment.suggestion)"
                style="margin-bottom: 4px;"
              >
                {{ getAiSuggestionText(record.ai_assessment.suggestion) }}
              </a-tag>
              <div class="risk-score">
                风险评分: {{ record.ai_assessment.risk_score }}
              </div>
            </div>
            <span v-else class="text-gray">待评估</span>
          </template>

          <!-- 操作 -->
          <template #action="{ record }">
            <a-space>
              <a-button 
                type="link" 
                size="small" 
                @click="viewDetails(record)"
              >
                查看详情
              </a-button>
              
              <a-dropdown v-if="canOperate(record.status)">
                <template #overlay>
                  <a-menu @click="({ key }) => handleAction(key, record)">
                    <a-menu-item key="approve" v-if="record.status === 'pending'">
                      <CheckOutlined />
                      批准
                    </a-menu-item>
                    <a-menu-item key="reject" v-if="['pending', 'under_review'].includes(record.status)">
                      <CloseOutlined />
                      拒绝
                    </a-menu-item>
                    <a-menu-item key="return" v-if="['pending', 'under_review'].includes(record.status)">
                      <RollbackOutlined />
                      退回
                    </a-menu-item>
                    <a-menu-item key="assign" v-if="record.status === 'pending'">
                      <UserSwitchOutlined />
                      分配审核员
                    </a-menu-item>
                    <a-menu-item key="retry-ai" v-if="!record.ai_assessment">
                      <ReloadOutlined />
                      重新AI评估
                    </a-menu-item>
                  </a-menu>
                </template>
                <a-button type="link" size="small">
                  更多操作
                  <DownOutlined />
                </a-button>
              </a-dropdown>
            </a-space>
          </template>
        </a-table>
      </a-card>
    </div>

    <!-- 申请详情抽屉 -->
    <a-drawer
      v-model:open="detailDrawerVisible"
      title="贷款申请详情"
      width="800"
      :destroyOnClose="true"
    >
      <LoanApplicationDetail 
        v-if="detailDrawerVisible && selectedApplication"
        :application="selectedApplication"
        @refresh="refreshData"
        @close="detailDrawerVisible = false"
      />
    </a-drawer>

    <!-- 审批操作模态框 -->
    <ApprovalModal
      v-model:visible="approvalModalVisible"
      :application="selectedApplication"
      :action="currentAction"
      @success="handleApprovalSuccess"
    />

    <!-- 分配审核员模态框 -->
    <AssignReviewerModal
      v-model:visible="assignModalVisible"
      :applications="selectedApplications"
      @success="handleAssignSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  SearchOutlined,
  ReloadOutlined,
  DownloadOutlined,
  DownOutlined,
  UserOutlined,
  CheckOutlined,
  CloseOutlined,
  RollbackOutlined,
  UserSwitchOutlined,
  ArrowUpOutlined,
  ArrowDownOutlined,
  FileTextOutlined,
  BankOutlined,
  ExclamationCircleOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import type { TableColumnsType, TableProps } from 'ant-design-vue'
import { loanApi } from '@/api/modules/loan'
import { useAuthStore } from '@/stores/modules/auth'
import LoanApplicationDetail from './components/LoanApplicationDetail.vue'
import ApprovalModal from './components/ApprovalModal.vue'
import AssignReviewerModal from './components/AssignReviewerModal.vue'

/**
 * 类型定义
 */
interface LoanApplication {
  id: string
  user_info: {
    user_id: number
    real_name: string
    phone: string
    avatar?: string
    is_verified: boolean
    bank_verified: boolean
    id_card: string
    user_type: string
    created_at: string
    credit_score: number
  }
  product_name: string
  amount: number
  term: number
  status: string
  applied_at: string
  ai_assessment?: {
    suggestion: string
    risk_score: number
    confidence: number
    risk_factors: { name: string; score: number }[]
    analysis_text: string
  }
  materials: { id: string; name: string; file_url: string; verified: boolean }[]
  approval_records: { id: string; action: string; operator_name: string; comments: string; created_at: string }[]
  reviewer?: {
    id: number
    name: string
  }
}

interface FilterForm {
  status?: string
  product_id?: number
  risk_level?: string
  dateRange?: [string, string]
  applicant?: string
}

/**
 * 响应式数据
 */
const authStore = useAuthStore()
const loading = ref(false)
const applicationList = ref<LoanApplication[]>([])
const loanProducts = ref<any[]>([])
const selectedRowKeys = ref<string[]>([])
const selectedApplications = ref<LoanApplication[]>([])
const detailDrawerVisible = ref(false)
const approvalModalVisible = ref(false)
const assignModalVisible = ref(false)
const selectedApplication = ref<LoanApplication | null>(null)
const currentAction = ref<string>('')

// 筛选表单
const filters = reactive<FilterForm>({
  status: 'pending'
})

// 分页配置
const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number) => `共 ${total} 条记录`
})

// 统计数据
const statistics = ref([
  {
    key: 'pending',
    title: '待审批申请',
    value: 0,
    color: '#faad14',
    icon: 'FileTextOutlined',
    suffix: '',
    precision: 0,
    type: 'warning',
    trend: {
      type: 'increase',
      icon: 'ArrowUpOutlined',
      value: '+5',
      period: '较昨日'
    }
  },
  {
    key: 'total_amount',
    title: '申请总金额',
    value: 0,
    color: '#52c41a',
    icon: 'BankOutlined',
    suffix: '万',
    precision: 1,
    type: 'success',
    trend: {
      type: 'increase',
      icon: 'ArrowUpOutlined',
      value: '+12%',
      period: '较上月'
    }
  },
  {
    key: 'approval_rate',
    title: '审批通过率',
    value: 0,
    color: '#1890ff',
    icon: 'CheckOutlined',
    suffix: '%',
    precision: 1,
    type: 'info',
    trend: {
      type: 'decrease',
      icon: 'ArrowDownOutlined',
      value: '-2%',
      period: '较上月'
    }
  },
  {
    key: 'risk_alerts',
    title: '风险预警',
    value: 0,
    color: '#f5222d',
    icon: 'ExclamationCircleOutlined',
    suffix: '',
    precision: 0,
    type: 'danger',
    trend: {
      type: 'normal',
      icon: '',
      value: '正常',
      period: '当前状态'
    }
  }
])

/**
 * 表格列配置
 */
const columns: TableColumnsType = [
  {
    title: '申请人',
    key: 'applicant',
    width: 160,
    fixed: 'left',
    slots: { customRender: 'applicantInfo' }
  },
  {
    title: '申请信息',
    key: 'application',
    width: 200,
    slots: { customRender: 'applicationInfo' }
  },
  {
    title: '申请时间',
    dataIndex: 'applied_at',
    width: 120,
    customRender: ({ text }) => dayjs(text).format('MM-DD HH:mm')
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    slots: { customRender: 'status' }
  },
  {
    title: 'AI建议',
    key: 'ai_suggestion',
    width: 140,
    slots: { customRender: 'aiSuggestion' }
  },
  {
    title: '审核员',
    key: 'reviewer',
    width: 100,
    customRender: ({ record }) => record.reviewer?.name || '未分配'
  },
  {
    title: '操作',
    key: 'action',
    width: 160,
    fixed: 'right',
    slots: { customRender: 'action' }
  }
]

/**
 * 计算属性
 */
const hasSelectedItems = computed(() => selectedRowKeys.value.length > 0)

const rowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys: string[], rows: LoanApplication[]) => {
    selectedRowKeys.value = keys
    selectedApplications.value = rows
  },
  getCheckboxProps: (record: LoanApplication) => ({
    disabled: !canOperate(record.status)
  })
}))

/**
 * 方法定义
 */
const getIconComponent = (iconName: string) => {
  const iconMap: Record<string, any> = {
    FileTextOutlined,
    BankOutlined,
    CheckOutlined,
    ExclamationCircleOutlined,
    ArrowUpOutlined,
    ArrowDownOutlined
  }
  return iconMap[iconName] || FileTextOutlined
}

const formatStatValue = (value: number, precision: number = 0) => {
  return Number(value).toLocaleString('zh-CN', {
    minimumFractionDigits: precision,
    maximumFractionDigits: precision
  })
}

const formatCurrency = (amount: number) => {
  return (amount / 10000).toFixed(1)
}

const getStatusColor = (status: string) => {
  const colorMap: Record<string, string> = {
    'pending': 'orange',
    'under_review': 'blue',
    'approved': 'green',
    'rejected': 'red',
    'returned': 'purple'
  }
  return colorMap[status] || 'default'
}

const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    'pending': '待审批',
    'under_review': '审核中',
    'approved': '已批准',
    'rejected': '已拒绝',
    'returned': '已退回'
  }
  return textMap[status] || status
}

const getAiSuggestionColor = (suggestion: string) => {
  const colorMap: Record<string, string> = {
    'approve': 'green',
    'reject': 'red',
    'manual_review': 'orange'
  }
  return colorMap[suggestion] || 'default'
}

const getAiSuggestionText = (suggestion: string) => {
  const textMap: Record<string, string> = {
    'approve': '建议通过',
    'reject': '建议拒绝',
    'manual_review': '人工审核'
  }
  return textMap[suggestion] || suggestion
}

const canOperate = (status: string) => {
  return ['pending', 'under_review'].includes(status)
}

// 加载数据
const loadData = async () => {
  try {
    loading.value = true
    
    // 临时使用模拟数据，替换为真实API调用
    const mockApplications: LoanApplication[] = [
      {
        id: 'LA202501001',
        user_info: {
          user_id: 1001,
          real_name: '王小农',
          phone: '138****5678',
          avatar: '',
          is_verified: true,
          bank_verified: true,
          id_card: '420106199001015432',
          user_type: 'farmer',
          created_at: '2024-12-01T10:30:00Z',
          credit_score: 85
        },
        product_name: '农业生产贷',
        amount: 300000,
        term: 12,
        status: 'pending',
        applied_at: '2025-01-31T09:30:00Z',
        ai_assessment: {
          suggestion: 'approve',
          risk_score: 78,
          confidence: 0.85,
          risk_factors: [
            { name: '收入稳定性', score: 0.8 },
            { name: '信用记录', score: 0.9 },
            { name: '负债比率', score: 0.7 }
          ],
          analysis_text: '申请人信用记录良好，收入稳定，建议通过审批。'
        },
        materials: [
          { id: '1', name: '身份证', file_url: '/uploads/id_card.jpg', verified: true },
          { id: '2', name: '收入证明', file_url: '/uploads/income.pdf', verified: true }
        ],
        approval_records: [
          {
            id: '1',
            action: 'submit',
            operator_name: '王小农',
            comments: '提交贷款申请',
            created_at: '2025-01-31T09:30:00Z'
          }
        ]
      },
      {
        id: 'LA202501002',
        user_info: {
          user_id: 1002,
          real_name: '李种植',
          phone: '139****6789',
          avatar: '',
          is_verified: true,
          bank_verified: false,
          id_card: '430102198502158765',
          user_type: 'farmer',
          created_at: '2024-11-15T14:20:00Z',
          credit_score: 72
        },
        product_name: '设备采购贷',
        amount: 500000,
        term: 24,
        status: 'under_review',
        applied_at: '2025-01-30T16:45:00Z',
        ai_assessment: {
          suggestion: 'manual_review',
          risk_score: 65,
          confidence: 0.72,
          risk_factors: [
            { name: '收入稳定性', score: 0.6 },
            { name: '信用记录', score: 0.8 },
            { name: '负债比率', score: 0.5 }
          ],
          analysis_text: '申请人收入波动较大，建议人工审核。'
        },
        materials: [
          { id: '3', name: '身份证', file_url: '/uploads/id_card2.jpg', verified: true },
          { id: '4', name: '设备采购合同', file_url: '/uploads/contract.pdf', verified: false }
        ],
        approval_records: [
          {
            id: '2',
            action: 'submit',
            operator_name: '李种植',
            comments: '提交设备采购贷款申请',
            created_at: '2025-01-30T16:45:00Z'
          },
          {
            id: '3',
            action: 'assign',
            operator_name: '系统',
            comments: '已分配给张审核员',
            created_at: '2025-01-31T08:00:00Z'
          }
        ],
        reviewer: {
          id: 2001,
          name: '张审核员'
        }
      },
      {
        id: 'LA202501003',
        user_info: {
          user_id: 1003,
          real_name: '陈农场主',
          phone: '137****4567',
          avatar: '',
          is_verified: true,
          bank_verified: true,
          id_card: '510104197812234321',
          user_type: 'farm_owner',
          created_at: '2024-08-20T11:15:00Z',
          credit_score: 92
        },
        product_name: '扩建贷款',
        amount: 1200000,
        term: 36,
        status: 'approved',
        applied_at: '2025-01-29T13:20:00Z',
        ai_assessment: {
          suggestion: 'approve',
          risk_score: 88,
          confidence: 0.95,
          risk_factors: [
            { name: '收入稳定性', score: 0.9 },
            { name: '信用记录', score: 0.95 },
            { name: '负债比率', score: 0.85 }
          ],
          analysis_text: '申请人资质优秀，强烈建议通过。'
        },
        materials: [
          { id: '5', name: '身份证', file_url: '/uploads/id_card3.jpg', verified: true },
          { id: '6', name: '营业执照', file_url: '/uploads/license.pdf', verified: true },
          { id: '7', name: '财务报表', file_url: '/uploads/financial.pdf', verified: true }
        ],
        approval_records: [
          {
            id: '4',
            action: 'submit',
            operator_name: '陈农场主',
            comments: '提交扩建贷款申请',
            created_at: '2025-01-29T13:20:00Z'
          },
          {
            id: '5',
            action: 'approve',
            operator_name: '王主管',
            comments: '申请人资质优秀，同意放款',
            created_at: '2025-01-30T10:30:00Z'
          }
        ],
        reviewer: {
          id: 2002,
          name: '王主管'
        }
      }
    ]
    
    // 根据筛选条件过滤数据
    let filteredData = mockApplications
    
    if (filters.status) {
      filteredData = filteredData.filter(app => app.status === filters.status)
    }
    
    if (filters.applicant) {
      filteredData = filteredData.filter(app => 
        app.user_info.real_name.includes(filters.applicant!) ||
        app.user_info.phone.includes(filters.applicant!)
      )
    }
    
    // 分页处理
    const startIndex = (pagination.current - 1) * pagination.pageSize
    const endIndex = startIndex + pagination.pageSize
    
    applicationList.value = filteredData.slice(startIndex, endIndex)
    pagination.total = filteredData.length
    
    // 真实项目中的API调用代码（注释掉）
    /*
    const params = {
      ...filters,
      page: pagination.current,
      limit: pagination.pageSize,
      date_range_start: filters.dateRange?.[0],
      date_range_end: filters.dateRange?.[1]
    }
    
    const response = await loanApi.getApplicationList(params)
    applicationList.value = response.data.applications || []
    pagination.total = response.data.total || 0
    */
    
  } catch (error) {
    console.error('加载申请列表失败:', error)
    message.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

// 加载统计数据
const loadStatistics = async () => {
  try {
    // 使用模拟数据
    const mockStats = {
      pending_applications: 15,
      total_applications: 156,
      approved_applications: 124,
      rejected_applications: 17,
      total_loan_amount: 45600000,
      approval_rate: 79.5,
      risk_alerts: 3
    }
    
    statistics.value[0].value = mockStats.pending_applications
    statistics.value[1].value = mockStats.total_loan_amount / 10000
    statistics.value[2].value = mockStats.approval_rate
    statistics.value[3].value = mockStats.risk_alerts
    
    // 真实项目中的API调用代码（注释掉）
    /*
    const response = await loanApi.getStatistics()
    const data = response.data
    
    statistics.value[0].value = data.pending_applications || 0
    statistics.value[1].value = (data.total_loan_amount || 0) / 10000
    statistics.value[2].value = ((data.approved_applications || 0) / (data.total_applications || 1)) * 100
    statistics.value[3].value = data.risk_alerts || 0
    */
    
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

// 加载贷款产品
const loadLoanProducts = async () => {
  try {
    // 使用模拟数据
    loanProducts.value = [
      { id: 1, product_name: '农业生产贷', description: '支持农业生产的贷款产品', min_amount: 10000, max_amount: 500000, interest_rate: 6.5, max_term: 24 },
      { id: 2, product_name: '设备采购贷', description: '农机设备采购专用贷款', min_amount: 50000, max_amount: 1000000, interest_rate: 7.2, max_term: 36 },
      { id: 3, product_name: '扩建贷款', description: '农场扩建改造贷款', min_amount: 100000, max_amount: 2000000, interest_rate: 7.8, max_term: 60 },
      { id: 4, product_name: '季节性周转贷', description: '季节性资金周转贷款', min_amount: 5000, max_amount: 200000, interest_rate: 8.5, max_term: 12 }
    ]
    
    // 真实项目中的API调用代码（注释掉）
    /*
    const response = await loanApi.getLoanProducts()
    loanProducts.value = response.data || []
    */
  } catch (error) {
    console.error('加载贷款产品失败:', error)
  }
}

// 搜索处理
const handleSearch = () => {
  pagination.current = 1
  loadData()
}

// 重置筛选
const resetFilters = () => {
  Object.keys(filters).forEach(key => {
    delete filters[key as keyof FilterForm]
  })
  filters.status = 'pending'
  handleSearch()
}

// 表格变化处理
const handleTableChange: TableProps['onChange'] = (page) => {
  if (page) {
    pagination.current = page.current || 1
    pagination.pageSize = page.pageSize || 20
  }
  loadData()
}

// 查看详情
const viewDetails = (record: LoanApplication) => {
  selectedApplication.value = record
  detailDrawerVisible.value = true
}

// 操作处理
const handleAction = async (action: string, record: LoanApplication) => {
  selectedApplication.value = record
  currentAction.value = action

  if (action === 'retry-ai') {
    await handleRetryAI(record)
  } else if (action === 'assign') {
    selectedApplications.value = [record]
    assignModalVisible.value = true
  } else {
    approvalModalVisible.value = true
  }
}

// 重试AI评估
const handleRetryAI = async (record: LoanApplication) => {
  try {
    loading.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 更新record的AI评估数据
    record.ai_assessment = {
      suggestion: 'approve',
      risk_score: Math.floor(Math.random() * 30) + 70,
      confidence: Math.random() * 0.3 + 0.7,
      risk_factors: [
        { name: '收入稳定性', score: Math.random() },
        { name: '信用记录', score: Math.random() },
        { name: '负债比率', score: Math.random() }
      ],
      analysis_text: 'AI重新评估完成，建议审核通过。'
    }
    
    message.success('AI评估已重新开始')
    await loadData()
    
    // 真实项目中的API调用代码（注释掉）
    /*
    await loanApi.retryAIAssessment(record.id)
    message.success('AI评估已重新开始')
    await loadData()
    */
  } catch (error) {
    console.error('重试AI评估失败:', error)
    message.error('重试AI评估失败')
  } finally {
    loading.value = false
  }
}

// 批量操作
const handleBatchAction = ({ key }: { key: string }) => {
  if (selectedApplications.value.length === 0) {
    message.warning('请先选择要操作的申请')
    return
  }

  if (key === 'assign') {
    assignModalVisible.value = true
  } else {
    Modal.confirm({
      title: '确认批量操作',
      content: `确定要对 ${selectedApplications.value.length} 个申请执行${key === 'approve' ? '批准' : '拒绝'}操作吗？`,
      onOk: () => performBatchAction(key)
    })
  }
}

// 执行批量操作
const performBatchAction = async (action: string) => {
  try {
    loading.value = true
    
    // 模拟批量操作
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    // 更新选中项目的状态
    selectedApplications.value.forEach(app => {
      if (action === 'approve') {
        app.status = 'approved'
      } else if (action === 'reject') {
        app.status = 'rejected'
      }
    })
    
    if (action === 'approve') {
      message.success('批量审批成功')
    } else if (action === 'reject') {
      message.success('批量拒绝成功')
    }
    
    selectedRowKeys.value = []
    selectedApplications.value = []
    await loadData()
    
    // 真实项目中的API调用代码（注释掉）
    /*
    const applicationIds = selectedApplications.value.map(app => app.id)
    
    if (action === 'approve') {
      await loanApi.batchApprove(applicationIds)
      message.success('批量审批成功')
    } else if (action === 'reject') {
      await loanApi.batchReject(applicationIds)
      message.success('批量拒绝成功')
    }
    
    selectedRowKeys.value = []
    selectedApplications.value = []
    await loadData()
    */
    
  } catch (error) {
    console.error('批量操作失败:', error)
    message.error('批量操作失败')
  } finally {
    loading.value = false
  }
}

// 审批成功回调
const handleApprovalSuccess = () => {
  approvalModalVisible.value = false
  selectedApplication.value = null
  loadData()
  loadStatistics()
}

// 分配成功回调
const handleAssignSuccess = () => {
  assignModalVisible.value = false
  selectedApplications.value = []
  selectedRowKeys.value = []
  loadData()
}

// 批量审批
const batchApprove = () => {
  handleBatchAction({ key: 'approve' })
}

// 刷新数据
const refreshData = () => {
  loadData()
  loadStatistics()
}

// 导出数据
const exportData = async () => {
  try {
    loading.value = true
    
    // 模拟导出功能
    await new Promise(resolve => setTimeout(resolve, 2000))
    message.success('导出成功')
    
    // 真实项目中的API调用代码（注释掉）
    /*
    const params = { ...filters }
    await loanApi.exportApplications(params)
    message.success('导出成功')
    */
  } catch (error) {
    console.error('导出失败:', error)
    message.error('导出失败')
  } finally {
    loading.value = false
  }
}

/**
 * 生命周期
 */
onMounted(() => {
  Promise.all([
    loadData(),
    loadStatistics(),
    loadLoanProducts()
  ])
})
</script>

<style lang="scss" scoped>
.loan-management {
  padding: 0;
  background: transparent;

  .page-header {
    margin-bottom: 24px;
    
    .header-content {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .page-title {
        margin: 0;
        font-size: 24px;
        font-weight: 600;
        color: #262626;
      }
    }
  }

  .statistics-section {
    margin-bottom: 24px;
    
    .stat-card {
      background: white;
      border-radius: 12px;
      padding: 20px;
      border: 1px solid rgba(0, 0, 0, 0.06);
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      height: 100%;
      
      &:hover {
        transform: translateY(-4px);
        box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
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
    }
  }

  .filter-section {
    margin-bottom: 24px;
  }

  .table-section {
    .table-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .table-actions {
        margin-left: auto;
      }
    }
    
    .applicant-info {
      display: flex;
      align-items: center;
      gap: 8px;
      
      .applicant-details {
        .applicant-name {
          font-weight: 500;
          color: #262626;
        }
        
        .applicant-phone {
          font-size: 12px;
          color: #8c8c8c;
        }
      }
    }
    
    .application-info {
      .application-id {
        font-size: 12px;
        color: #8c8c8c;
        margin-bottom: 2px;
      }
      
      .product-name {
        font-weight: 500;
        color: #262626;
        margin-bottom: 2px;
      }
      
      .application-amount {
        font-size: 14px;
        color: #52c41a;
        font-weight: 600;
      }
    }
    
    .ai-suggestion {
      .risk-score {
        font-size: 12px;
        color: #8c8c8c;
        margin-top: 2px;
      }
    }
  }
}

.text-gray {
  color: #8c8c8c;
}

// 响应式设计
@media (max-width: 768px) {
  .loan-management {
    .page-header .header-content {
      flex-direction: column;
      align-items: flex-start;
      gap: 16px;
    }
    
    .statistics-section .stat-card {
      padding: 16px;
      
      .stat-content .stat-value {
        font-size: 24px;
      }
    }
  }
}

// 暗色主题适配
:deep([data-theme="dark"]) {
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
  
  .page-title {
    color: rgba(255, 255, 255, 0.85) !important;
  }
}
</style> 