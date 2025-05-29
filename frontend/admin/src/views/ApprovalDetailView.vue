<template>
  <div class="approval-detail">
    <div class="page-header">
      <div class="header-left">
        <el-button @click="goBack" size="default">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <h2 class="page-title">审批详情</h2>
      </div>
      <div class="header-actions">
        <el-button @click="refreshData" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <div v-loading="loading" class="detail-content">
      <el-row :gutter="20">
        <!-- 左侧主要信息 -->
        <el-col :span="16">
          <!-- 基本信息 -->
          <el-card class="info-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Document /></el-icon>
                  申请基本信息
                </span>
                <el-tag :type="getStatusType(detail.status)" size="large">
                  {{ detail.status }}
                </el-tag>
              </div>
            </template>
            
            <el-descriptions :column="2" border>
              <el-descriptions-item label="申请编号">
                {{ detail.application_id }}
              </el-descriptions-item>
              <el-descriptions-item label="申请人">
                {{ detail.applicant_details?.real_name }}
              </el-descriptions-item>
              <el-descriptions-item label="身份证号">
                {{ maskIdCard(detail.applicant_details?.id_card_number) }}
              </el-descriptions-item>
              <el-descriptions-item label="联系地址">
                {{ detail.applicant_details?.address }}
              </el-descriptions-item>
              <el-descriptions-item label="申请金额">
                <span class="amount">¥{{ formatAmount(detail.amount) }}</span>
              </el-descriptions-item>
              <el-descriptions-item label="申请期限">
                {{ detail.term_months }}个月
              </el-descriptions-item>
              <el-descriptions-item label="贷款用途" span="2">
                {{ detail.purpose }}
              </el-descriptions-item>
              <el-descriptions-item label="提交时间">
                {{ formatDateTime(detail.submitted_at) }}
              </el-descriptions-item>
              <el-descriptions-item label="更新时间">
                {{ formatDateTime(detail.updated_at) }}
              </el-descriptions-item>
              <el-descriptions-item v-if="detail.approved_amount" label="批准金额">
                <span class="amount approved">¥{{ formatAmount(detail.approved_amount) }}</span>
              </el-descriptions-item>
            </el-descriptions>
          </el-card>

          <!-- AI分析报告 -->
          <el-card v-if="detail.ai_analysis_report" class="ai-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Cpu /></el-icon>
                  AI智能分析报告
                </span>
                <el-tag type="info" size="small">
                  风险评分: {{ detail.ai_analysis_report.overall_risk_score }}分
                </el-tag>
              </div>
            </template>
            
            <div class="ai-analysis">
              <!-- 风险评分 -->
              <div class="analysis-section">
                <h4>风险评估</h4>
                <div class="risk-score-display">
                  <el-progress
                    type="circle"
                    :percentage="detail.ai_analysis_report.overall_risk_score"
                    :color="getRiskColor(detail.ai_analysis_report.overall_risk_score)"
                    :width="120"
                  >
                    <template #default="{ percentage }">
                      <span class="risk-score-text">{{ percentage }}分</span>
                    </template>
                  </el-progress>
                  <div class="risk-level">
                    {{ getRiskLevel(detail.ai_analysis_report.overall_risk_score) }}
                  </div>
                </div>
              </div>

              <!-- 风险因素 -->
              <div class="analysis-section">
                <h4>识别的风险因素</h4>
                <div class="risk-factors">
                  <el-tag
                    v-for="factor in detail.ai_analysis_report.risk_factors"
                    :key="factor"
                    type="warning"
                    size="small"
                    class="risk-factor-tag"
                  >
                    {{ factor }}
                  </el-tag>
                </div>
              </div>

              <!-- 数据验证结果 -->
              <div class="analysis-section">
                <h4>数据验证结果</h4>
                <div class="verification-results">
                  <div
                    v-for="result in detail.ai_analysis_report.data_verification_results"
                    :key="result.item"
                    class="verification-item"
                  >
                    <span class="verification-label">{{ result.item }}</span>
                    <el-tag
                      :type="result.result === '通过' ? 'success' : 'danger'"
                      size="small"
                    >
                      {{ result.result }}
                    </el-tag>
                  </div>
                </div>
              </div>

              <!-- AI建议 -->
              <div class="analysis-section">
                <h4>AI处理建议</h4>
                <div class="ai-suggestion">
                  <el-alert
                    :title="detail.ai_analysis_report.suggestion"
                    type="info"
                    :closable="false"
                    show-icon
                  />
                </div>
              </div>
            </div>
          </el-card>

          <!-- 上传文件 -->
          <el-card class="documents-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Folder /></el-icon>
                  上传文件
                </span>
              </div>
            </template>
            
            <div class="documents-list">
              <div
                v-for="doc in detail.uploaded_documents_details"
                :key="doc.file_id"
                class="document-item"
              >
                <div class="doc-info">
                  <el-icon class="doc-icon"><Document /></el-icon>
                  <div class="doc-details">
                    <div class="doc-type">{{ getDocTypeName(doc.doc_type) }}</div>
                    <div v-if="doc.ocr_result" class="doc-ocr">
                      OCR识别: {{ doc.ocr_result }}
                    </div>
                  </div>
                </div>
                <div class="doc-actions">
                  <el-button size="small" @click="previewDocument(doc)">
                    预览
                  </el-button>
                  <el-button size="small" type="primary" @click="downloadDocument(doc)">
                    下载
                  </el-button>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 右侧操作区域 -->
        <el-col :span="8">
          <!-- 审批操作 -->
          <el-card v-if="canReview(detail.status)" class="review-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Check /></el-icon>
                  审批操作
                </span>
              </div>
            </template>
            
            <el-form
              ref="reviewFormRef"
              :model="reviewForm"
              :rules="reviewRules"
              label-width="80px"
            >
              <el-form-item label="审批决策" prop="decision" required>
                <el-radio-group v-model="reviewForm.decision">
                  <el-radio value="approved" class="decision-radio">
                    <el-icon><CircleCheck /></el-icon>
                    批准
                  </el-radio>
                  <el-radio value="rejected" class="decision-radio">
                    <el-icon><CircleClose /></el-icon>
                    拒绝
                  </el-radio>
                  <el-radio value="require_more_info" class="decision-radio">
                    <el-icon><Warning /></el-icon>
                    补充材料
                  </el-radio>
                </el-radio-group>
              </el-form-item>
              
              <el-form-item
                v-if="reviewForm.decision === 'approved'"
                label="批准金额"
                prop="approved_amount"
              >
                <el-input-number
                  v-model="reviewForm.approved_amount"
                  :min="1"
                  :max="detail.amount"
                  :step="1000"
                  style="width: 100%"
                />
                <div class="amount-hint">
                  申请金额: ¥{{ formatAmount(detail.amount) }}
                </div>
              </el-form-item>
              
              <el-form-item label="审批意见" prop="comments" required>
                <el-input
                  v-model="reviewForm.comments"
                  type="textarea"
                  :rows="4"
                  placeholder="请输入审批意见"
                />
              </el-form-item>
              
              <el-form-item
                v-if="reviewForm.decision === 'require_more_info'"
                label="补充说明"
                prop="required_info_details"
              >
                <el-input
                  v-model="reviewForm.required_info_details"
                  type="textarea"
                  :rows="3"
                  placeholder="请说明需要补充的材料或信息"
                />
              </el-form-item>
              
              <el-form-item>
                <el-button
                  type="primary"
                  @click="submitApprovalReview"
                  :loading="submitting"
                  style="width: 100%"
                >
                  提交审批
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>

          <!-- 审批历史 -->
          <el-card class="history-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Clock /></el-icon>
                  审批历史
                </span>
              </div>
            </template>
            
            <el-timeline>
              <el-timeline-item
                v-for="(item, index) in detail.history"
                :key="index"
                :timestamp="formatDateTime(item.timestamp)"
                placement="top"
                :type="getTimelineType(item.status)"
              >
                <div class="timeline-content">
                  <div class="timeline-status">{{ item.status }}</div>
                  <div class="timeline-operator">操作人: {{ item.operator }}</div>
                </div>
              </el-timeline-item>
            </el-timeline>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import {
  ArrowLeft,
  Refresh,
  Document,
  Cpu,
  Folder,
  Check,
  CircleCheck,
  CircleClose,
  Warning,
  Clock
} from '@element-plus/icons-vue'
import { getApplicationDetail, submitReview } from '@/api/admin'
import type { ApplicationDetail } from '@/types'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const submitting = ref(false)
const reviewFormRef = ref<FormInstance>()

const detail = ref<ApplicationDetail>({} as ApplicationDetail)

// 审批表单
const reviewForm = reactive({
  decision: '',
  approved_amount: 0,
  comments: '',
  required_info_details: ''
})

const reviewRules: FormRules = {
  decision: [
    { required: true, message: '请选择审批决策', trigger: 'change' }
  ],
  comments: [
    { required: true, message: '请输入审批意见', trigger: 'blur' }
  ],
  approved_amount: [
    { required: true, message: '请输入批准金额', trigger: 'blur' }
  ]
}

// 方法
const fetchDetail = async () => {
  try {
    loading.value = true
    const applicationId = route.params.id as string
    const data = await getApplicationDetail(applicationId)
    detail.value = data
    
    // 初始化审批表单
    if (canReview(data.status)) {
      reviewForm.approved_amount = data.amount
    }
  } catch (error) {
    ElMessage.error('获取申请详情失败')
  } finally {
    loading.value = false
  }
}

const refreshData = () => {
  fetchDetail()
}

const goBack = () => {
  router.back()
}

const canReview = (status: string) => {
  return status === '待人工复核'
}

const submitApprovalReview = async () => {
  if (!reviewFormRef.value) return
  
  try {
    await reviewFormRef.value.validate()
    submitting.value = true
    
    const submitData = {
      decision: reviewForm.decision as 'approved' | 'rejected' | 'require_more_info',
      approved_amount: reviewForm.decision === 'approved' ? reviewForm.approved_amount : undefined,
      comments: reviewForm.comments,
      required_info_details: reviewForm.decision === 'require_more_info' ? reviewForm.required_info_details : undefined
    }
    
    await submitReview(detail.value.application_id, submitData)
    
    ElMessage.success('审批提交成功')
    fetchDetail() // 重新获取详情
  } catch (error) {
    ElMessage.error('审批提交失败')
  } finally {
    submitting.value = false
  }
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    'AI_审批中': 'warning',
    '待人工复核': 'info',
    '已批准': 'success',
    '已拒绝': 'danger'
  }
  return statusMap[status] || 'info'
}

const getTimelineType = (status: string) => {
  const typeMap: Record<string, string> = {
    '已提交': 'primary',
    'AI_审批中': 'warning',
    '待人工复核': 'info',
    '已批准': 'success',
    '已拒绝': 'danger'
  }
  return typeMap[status] || 'primary'
}

const getRiskColor = (score: number) => {
  if (score <= 30) return '#67c23a'
  if (score <= 70) return '#e6a23c'
  return '#f56c6c'
}

const getRiskLevel = (score: number) => {
  if (score <= 30) return '低风险'
  if (score <= 70) return '中风险'
  return '高风险'
}

const getDocTypeName = (docType: string) => {
  const typeMap: Record<string, string> = {
    'id_card_front': '身份证正面',
    'id_card_back': '身份证背面',
    'land_contract': '土地承包合同',
    'income_proof': '收入证明',
    'asset_proof': '资产证明'
  }
  return typeMap[docType] || docType
}

const maskIdCard = (idCard?: string) => {
  if (!idCard) return '-'
  return idCard.replace(/(.{4}).*(.{4})/, '$1****$2')
}

const formatAmount = (amount: number | undefined | null) => {
  if (typeof amount !== 'number' || isNaN(amount)) return '-'
  return amount.toLocaleString()
}

const formatDateTime = (datetime: string) => {
  return dayjs(datetime).format('YYYY-MM-DD HH:mm:ss')
}

const previewDocument = (doc: any) => {
  // 打开文档预览
  window.open(doc.file_url, '_blank')
}

const downloadDocument = (doc: any) => {
  // 下载文档
  const link = document.createElement('a')
  link.href = doc.file_url
  link.download = getDocTypeName(doc.doc_type)
  link.click()
}

onMounted(() => {
  fetchDetail()
})
</script>

<style scoped>
.approval-detail {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.page-title {
  margin: 0;
  color: #333;
  font-size: 24px;
  font-weight: 600;
}

.info-card,
.ai-card,
.documents-card,
.review-card,
.history-card {
  margin-bottom: 20px;
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  gap: 8px;
}

.amount {
  font-weight: 600;
  color: #f56c6c;
}

.amount.approved {
  color: #67c23a;
}

.ai-analysis {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.analysis-section h4 {
  margin: 0 0 12px 0;
  color: #333;
  font-size: 16px;
  font-weight: 600;
}

.risk-score-display {
  display: flex;
  align-items: center;
  gap: 20px;
}

.risk-score-text {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.risk-level {
  font-size: 16px;
  font-weight: 500;
  color: #666;
}

.risk-factors {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.risk-factor-tag {
  margin: 0;
}

.verification-results {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.verification-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #f8f9fa;
  border-radius: 6px;
}

.verification-label {
  color: #333;
  font-size: 14px;
}

.documents-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.document-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  background: #fafafa;
}

.doc-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.doc-icon {
  font-size: 20px;
  color: #409eff;
}

.doc-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.doc-type {
  font-weight: 500;
  color: #333;
}

.doc-ocr {
  font-size: 12px;
  color: #666;
}

.doc-actions {
  display: flex;
  gap: 8px;
}

.decision-radio {
  display: flex;
  align-items: center;
  gap: 4px;
  width: 100%;
  margin-bottom: 12px;
}

.amount-hint {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.timeline-content {
  background: #fff;
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid #ebeef5;
}

.timeline-status {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.timeline-operator {
  font-size: 12px;
  color: #666;
}

:deep(.el-descriptions__label) {
  font-weight: 500;
}

:deep(.el-radio) {
  margin-right: 0;
  margin-bottom: 12px;
}

:deep(.el-timeline-item__content) {
  padding-bottom: 12px;
}
</style> 