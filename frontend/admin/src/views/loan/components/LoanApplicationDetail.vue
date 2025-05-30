<template>
  <div class="loan-application-detail">
    <a-spin :spinning="loading" size="large">
      <!-- 基本信息区域 -->
      <div class="detail-section">
        <h3 class="section-title">申请基本信息</h3>
        <a-row :gutter="[24, 16]">
          <a-col :span="12">
            <div class="info-item">
              <label>申请编号：</label>
              <span>{{ application.id }}</span>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="info-item">
              <label>申请时间：</label>
              <span>{{ formatDateTime(application.applied_at) }}</span>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="info-item">
              <label>产品名称：</label>
              <span>{{ application.product_name }}</span>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="info-item">
              <label>申请金额：</label>
              <span class="amount">¥{{ formatCurrency(application.amount) }}</span>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="info-item">
              <label>申请期限：</label>
              <span>{{ application.term || 12 }}个月</span>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="info-item">
              <label>当前状态：</label>
              <a-tag :color="getStatusColor(application.status)">
                {{ getStatusText(application.status) }}
              </a-tag>
            </div>
          </a-col>
        </a-row>
      </div>

      <!-- 申请人信息 -->
      <div class="detail-section">
        <h3 class="section-title">申请人信息</h3>
        <a-row :gutter="[24, 16]">
          <a-col :span="6">
            <div class="user-avatar-section">
              <a-avatar :src="application.user_info.avatar" :size="80">
                <template #icon>
                  <UserOutlined />
                </template>
              </a-avatar>
              <div class="verification-badges">
                <a-tag color="green" v-if="application.user_info.is_verified">实名认证</a-tag>
                <a-tag color="blue" v-if="application.user_info.bank_verified">银行认证</a-tag>
              </div>
            </div>
          </a-col>
          <a-col :span="18">
            <a-row :gutter="[24, 16]">
              <a-col :span="12">
                <div class="info-item">
                  <label>姓名：</label>
                  <span>{{ application.user_info.real_name }}</span>
                </div>
              </a-col>
              <a-col :span="12">
                <div class="info-item">
                  <label>手机号：</label>
                  <span>{{ application.user_info.phone }}</span>
                </div>
              </a-col>
              <a-col :span="12">
                <div class="info-item">
                  <label>身份证号：</label>
                  <span>{{ maskIdCard(application.user_info.id_card) }}</span>
                </div>
              </a-col>
              <a-col :span="12">
                <div class="info-item">
                  <label>用户类型：</label>
                  <span>{{ getUserTypeText(application.user_info.user_type) }}</span>
                </div>
              </a-col>
              <a-col :span="12">
                <div class="info-item">
                  <label>注册时间：</label>
                  <span>{{ formatDateTime(application.user_info.created_at) }}</span>
                </div>
              </a-col>
              <a-col :span="12">
                <div class="info-item">
                  <label>信用等级：</label>
                  <a-rate 
                    :value="getCreditLevel(application.user_info.credit_score)" 
                    disabled 
                    :count="5"
                  />
                </div>
              </a-col>
            </a-row>
          </a-col>
        </a-row>
      </div>

      <!-- AI风险评估 -->
      <div class="detail-section" v-if="application.ai_assessment">
        <h3 class="section-title">
          AI风险评估
          <a-tag :color="getAiSuggestionColor(application.ai_assessment.suggestion)" style="margin-left: 12px;">
            {{ getAiSuggestionText(application.ai_assessment.suggestion) }}
          </a-tag>
        </h3>
        
        <a-row :gutter="[24, 16]">
          <a-col :span="8">
            <div class="score-display">
              <div class="score-circle">
                <a-progress 
                  type="circle" 
                  :percent="application.ai_assessment.risk_score" 
                  :size="120"
                  :stroke-color="getRiskScoreColor(application.ai_assessment.risk_score)"
                  :format="() => `${application.ai_assessment.risk_score}分`"
                />
              </div>
              <div class="score-label">风险评分</div>
            </div>
          </a-col>
          <a-col :span="8">
            <div class="score-display">
              <div class="score-circle">
                <a-progress 
                  type="circle" 
                  :percent="application.ai_assessment.confidence * 100" 
                  :size="120"
                  stroke-color="#52c41a"
                  :format="() => `${(application.ai_assessment.confidence * 100).toFixed(0)}%`"
                />
              </div>
              <div class="score-label">置信度</div>
            </div>
          </a-col>
          <a-col :span="8">
            <div class="risk-factors">
              <h4>风险因子分析</h4>
              <div class="factor-list">
                <div 
                  v-for="factor in application.ai_assessment.risk_factors" 
                  :key="factor.name"
                  class="factor-item"
                >
                  <div class="factor-name">{{ factor.name }}</div>
                  <div class="factor-score">
                    <a-progress 
                      :percent="factor.score * 100" 
                      :stroke-color="getFactorColor(factor.score)"
                      size="small"
                      :show-info="false"
                    />
                    <span class="score-text">{{ (factor.score * 100).toFixed(0) }}%</span>
                  </div>
                </div>
              </div>
            </div>
          </a-col>
        </a-row>

        <div class="ai-analysis-text">
          <h4>AI分析报告</h4>
          <div class="analysis-content">
            {{ application.ai_assessment.analysis_text }}
          </div>
        </div>
      </div>

      <!-- 申请材料 -->
      <div class="detail-section">
        <h3 class="section-title">申请材料</h3>
        <div class="materials-grid">
          <div 
            v-for="material in application.materials" 
            :key="material.id"
            class="material-item"
          >
            <div class="material-header">
              <FileOutlined />
              <span>{{ material.name }}</span>
              <a-tag 
                :color="material.verified ? 'green' : 'orange'"
                size="small"
              >
                {{ material.verified ? '已验证' : '待验证' }}
              </a-tag>
            </div>
            <div class="material-preview" @click="previewFile(material)">
              <img 
                v-if="isImage(material.file_url)" 
                :src="material.file_url" 
                :alt="material.name"
                loading="lazy"
              />
              <div v-else class="file-placeholder">
                <FileOutlined style="font-size: 24px;" />
                <span>{{ getFileExtension(material.file_url) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 审批历史 -->
      <div class="detail-section">
        <h3 class="section-title">审批历史</h3>
        <a-timeline>
          <a-timeline-item 
            v-for="record in application.approval_records" 
            :key="record.id"
            :color="getTimelineColor(record.action)"
          >
            <template #dot>
              <component :is="getTimelineIcon(record.action)" />
            </template>
            <div class="timeline-content">
              <div class="timeline-header">
                <span class="action">{{ getActionText(record.action) }}</span>
                <span class="operator">{{ record.operator_name }}</span>
                <span class="time">{{ formatDateTime(record.created_at) }}</span>
              </div>
              <div class="timeline-comment" v-if="record.comments">
                {{ record.comments }}
              </div>
            </div>
          </a-timeline-item>
        </a-timeline>
      </div>

      <!-- 操作区域 -->
      <div class="action-section" v-if="canOperate">
        <a-space size="large">
          <a-button 
            type="primary" 
            size="large"
            @click="handleApprove"
            :loading="actionLoading"
          >
            <CheckOutlined />
            批准申请
          </a-button>
          <a-button 
            danger 
            size="large"
            @click="handleReject"
            :loading="actionLoading"
          >
            <CloseOutlined />
            拒绝申请
          </a-button>
          <a-button 
            size="large"
            @click="handleReturn"
            :loading="actionLoading"
          >
            <RollbackOutlined />
            退回修改
          </a-button>
          <a-button 
            size="large"
            @click="handleAssign"
          >
            <UserSwitchOutlined />
            分配审核员
          </a-button>
        </a-space>
      </div>
    </a-spin>

    <!-- 文件预览模态框 -->
    <a-modal
      v-model:open="previewModalVisible"
      title="文件预览"
      width="800"
      footer={null}
      centered
    >
      <div class="file-preview-content">
        <img 
          v-if="currentPreviewFile && isImage(currentPreviewFile.file_url)"
          :src="currentPreviewFile.file_url"
          style="width: 100%; max-height: 600px; object-fit: contain;"
        />
        <div v-else-if="currentPreviewFile" class="pdf-preview">
          <iframe 
            :src="currentPreviewFile.file_url"
            style="width: 100%; height: 600px; border: none;"
          ></iframe>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { message } from 'ant-design-vue'
import {
  UserOutlined,
  FileOutlined,
  CheckOutlined,
  CloseOutlined,
  RollbackOutlined,
  UserSwitchOutlined,
  ClockCircleOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'

interface Props {
  application: any
}

interface Emits {
  refresh: []
  close: []
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

/**
 * 响应式数据
 */
const loading = ref(false)
const actionLoading = ref(false)
const previewModalVisible = ref(false)
const currentPreviewFile = ref<any>(null)

/**
 * 计算属性
 */
const canOperate = computed(() => {
  return ['pending', 'under_review'].includes(props.application.status)
})

/**
 * 工具方法
 */
const formatDateTime = (dateTime: string) => {
  return dayjs(dateTime).format('YYYY-MM-DD HH:mm:ss')
}

const formatCurrency = (amount: number) => {
  return amount.toLocaleString('zh-CN')
}

const maskIdCard = (idCard: string) => {
  if (!idCard) return ''
  return idCard.replace(/(\d{6})\d{8}(\d{4})/, '$1********$2')
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

const getUserTypeText = (userType: string) => {
  const typeMap: Record<string, string> = {
    'farmer': '农户',
    'farm_owner': '农场主',
    'cooperative': '合作社',
    'enterprise': '企业'
  }
  return typeMap[userType] || userType
}

const getCreditLevel = (score: number) => {
  if (score >= 90) return 5
  if (score >= 80) return 4
  if (score >= 70) return 3
  if (score >= 60) return 2
  return 1
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

const getRiskScoreColor = (score: number) => {
  if (score >= 80) return '#52c41a'
  if (score >= 60) return '#faad14'
  return '#ff4d4f'
}

const getFactorColor = (score: number) => {
  if (score >= 0.8) return '#52c41a'
  if (score >= 0.6) return '#faad14'
  return '#ff4d4f'
}

const getTimelineColor = (action: string) => {
  const colorMap: Record<string, string> = {
    'submit': 'blue',
    'approve': 'green',
    'reject': 'red',
    'return': 'orange',
    'assign': 'purple'
  }
  return colorMap[action] || 'blue'
}

const getTimelineIcon = (action: string) => {
  const iconMap: Record<string, any> = {
    'submit': ClockCircleOutlined,
    'approve': CheckCircleOutlined,
    'reject': CloseCircleOutlined,
    'return': RollbackOutlined,
    'assign': UserSwitchOutlined
  }
  return iconMap[action] || ClockCircleOutlined
}

const getActionText = (action: string) => {
  const textMap: Record<string, string> = {
    'submit': '提交申请',
    'approve': '审批通过',
    'reject': '审批拒绝',
    'return': '退回修改',
    'assign': '分配审核员'
  }
  return textMap[action] || action
}

const isImage = (url: string) => {
  const imageExts = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp']
  return imageExts.some(ext => url.toLowerCase().includes(ext))
}

const getFileExtension = (url: string) => {
  const ext = url.split('.').pop()
  return ext?.toUpperCase() || 'FILE'
}

const previewFile = (material: any) => {
  currentPreviewFile.value = material
  previewModalVisible.value = true
}

/**
 * 操作方法
 */
const handleApprove = () => {
  emit('close')
  // 这里应该打开审批模态框
}

const handleReject = () => {
  emit('close')
  // 这里应该打开拒绝模态框
}

const handleReturn = () => {
  emit('close')
  // 这里应该打开退回模态框
}

const handleAssign = () => {
  emit('close')
  // 这里应该打开分配审核员模态框
}
</script>

<style lang="scss" scoped>
.loan-application-detail {
  padding: 0;
  
  .detail-section {
    margin-bottom: 32px;
    
    .section-title {
      font-size: 16px;
      font-weight: 600;
      color: #262626;
      margin-bottom: 16px;
      border-bottom: 1px solid #f0f0f0;
      padding-bottom: 8px;
    }
    
    .info-item {
      display: flex;
      align-items: center;
      margin-bottom: 12px;
      
      label {
        font-weight: 500;
        color: #595959;
        min-width: 80px;
        margin-right: 8px;
      }
      
      .amount {
        color: #52c41a;
        font-weight: 600;
        font-size: 16px;
      }
    }
  }

  .user-avatar-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    
    .verification-badges {
      margin-top: 12px;
      display: flex;
      flex-direction: column;
      gap: 4px;
    }
  }

  .score-display {
    display: flex;
    flex-direction: column;
    align-items: center;
    
    .score-circle {
      margin-bottom: 12px;
    }
    
    .score-label {
      font-size: 14px;
      color: #8c8c8c;
      font-weight: 500;
    }
  }

  .risk-factors {
    .factor-list {
      .factor-item {
        margin-bottom: 12px;
        
        .factor-name {
          font-size: 13px;
          color: #595959;
          margin-bottom: 4px;
        }
        
        .factor-score {
          display: flex;
          align-items: center;
          gap: 8px;
          
          .score-text {
            font-size: 12px;
            color: #8c8c8c;
            min-width: 35px;
          }
        }
      }
    }
  }

  .ai-analysis-text {
    margin-top: 24px;
    
    h4 {
      font-size: 14px;
      font-weight: 600;
      color: #262626;
      margin-bottom: 8px;
    }
    
    .analysis-content {
      background: #fafafa;
      border-radius: 6px;
      padding: 16px;
      font-size: 14px;
      line-height: 1.6;
      color: #595959;
    }
  }

  .materials-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 16px;
    
    .material-item {
      border: 1px solid #e8e8e8;
      border-radius: 8px;
      overflow: hidden;
      transition: all 0.3s ease;
      
      &:hover {
        border-color: #1890ff;
        box-shadow: 0 2px 8px rgba(24, 144, 255, 0.2);
      }
      
      .material-header {
        padding: 12px;
        background: #fafafa;
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 13px;
        border-bottom: 1px solid #e8e8e8;
      }
      
      .material-preview {
        height: 120px;
        cursor: pointer;
        overflow: hidden;
        
        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }
        
        .file-placeholder {
          height: 100%;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          color: #8c8c8c;
          background: #f5f5f5;
          
          span {
            margin-top: 8px;
            font-size: 12px;
          }
        }
      }
    }
  }

  .timeline-content {
    .timeline-header {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 4px;
      
      .action {
        font-weight: 500;
        color: #262626;
      }
      
      .operator {
        color: #1890ff;
      }
      
      .time {
        color: #8c8c8c;
        font-size: 12px;
        margin-left: auto;
      }
    }
    
    .timeline-comment {
      color: #595959;
      font-size: 13px;
      background: #fafafa;
      padding: 8px 12px;
      border-radius: 4px;
      margin-top: 8px;
    }
  }

  .action-section {
    padding: 24px;
    background: #fafafa;
    border-radius: 8px;
    text-align: center;
    margin-top: 32px;
    border: 1px solid #e8e8e8;
  }

  .file-preview-content {
    text-align: center;
    
    .pdf-preview {
      border: 1px solid #e8e8e8;
      border-radius: 4px;
      overflow: hidden;
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .loan-application-detail {
    .materials-grid {
      grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
      gap: 12px;
    }
    
    .action-section {
      padding: 16px;
      
      :deep(.ant-space) {
        flex-wrap: wrap;
      }
    }
  }
}

// 暗色主题适配
:deep([data-theme="dark"]) {
  .detail-section .section-title {
    color: rgba(255, 255, 255, 0.85) !important;
    border-bottom-color: #303030 !important;
  }
  
  .info-item label {
    color: rgba(255, 255, 255, 0.65) !important;
  }
  
  .ai-analysis-text .analysis-content {
    background: #262626 !important;
    color: rgba(255, 255, 255, 0.65) !important;
  }
  
  .action-section {
    background: #262626 !important;
    border-color: #303030 !important;
  }
  
  .material-item {
    border-color: #303030 !important;
    
    .material-header {
      background: #262626 !important;
      border-bottom-color: #303030 !important;
    }
    
    .file-placeholder {
      background: #1f1f1f !important;
    }
  }
}
</style> 