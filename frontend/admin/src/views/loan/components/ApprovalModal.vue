<template>
  <a-modal
    v-model:open="modalVisible"
    :title="getModalTitle()"
    width="600"
    :confirm-loading="loading"
    @ok="handleSubmit"
    @cancel="handleCancel"
  >
    <a-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      layout="vertical"
    >
      <!-- 批准操作表单 -->
      <template v-if="action === 'approve'">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="批准金额" name="approved_amount">
              <a-input-number
                v-model:value="formData.approved_amount"
                :min="1000"
                :max="application?.amount"
                :step="1000"
                style="width: 100%"
                :formatter="(value) => `¥ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')"
                :parser="(value) => value.replace(/¥\s?|(,*)/g, '')"
                placeholder="请输入批准金额"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="批准期限" name="approved_term">
              <a-select
                v-model:value="formData.approved_term"
                placeholder="请选择期限"
              >
                <a-select-option :value="6">6个月</a-select-option>
                <a-select-option :value="12">12个月</a-select-option>
                <a-select-option :value="18">18个月</a-select-option>
                <a-select-option :value="24">24个月</a-select-option>
                <a-select-option :value="36">36个月</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="年利率 (%)" name="interest_rate">
          <a-input-number
            v-model:value="formData.interest_rate"
            :min="3"
            :max="15"
            :step="0.1"
            :precision="2"
            style="width: 100%"
            placeholder="请输入年利率"
          />
        </a-form-item>

        <a-form-item label="还款开始日期" name="repayment_start_date">
          <a-date-picker
            v-model:value="formData.repayment_start_date"
            style="width: 100%"
            :disabled-date="disabledDate"
            placeholder="请选择还款开始日期"
          />
        </a-form-item>
      </template>

      <!-- 拒绝操作表单 -->
      <template v-if="action === 'reject'">
        <a-form-item label="拒绝原因" name="rejection_reason_code">
          <a-select
            v-model:value="formData.rejection_reason_code"
            placeholder="请选择拒绝原因"
          >
            <a-select-option value="INSUFFICIENT_CREDIT">信用评分不足</a-select-option>
            <a-select-option value="INCOME_INSUFFICIENT">收入不足</a-select-option>
            <a-select-option value="DEBT_RATIO_HIGH">负债率过高</a-select-option>
            <a-select-option value="INVALID_MATERIALS">材料不合规</a-select-option>
            <a-select-option value="RISK_TOO_HIGH">风险等级过高</a-select-option>
            <a-select-option value="POLICY_VIOLATION">违反政策规定</a-select-option>
            <a-select-option value="OTHER">其他原因</a-select-option>
          </a-select>
        </a-form-item>
      </template>

      <!-- 退回操作表单 -->
      <template v-if="action === 'return'">
        <a-form-item label="退回原因" name="return_reason">
          <a-input
            v-model:value="formData.return_reason"
            placeholder="请简要说明退回原因"
          />
        </a-form-item>

        <a-form-item label="需要补充的材料" name="required_materials">
          <a-select
            v-model:value="formData.required_materials"
            mode="multiple"
            placeholder="请选择需要补充的材料"
          >
            <a-select-option value="income_proof">收入证明</a-select-option>
            <a-select-option value="bank_statement">银行流水</a-select-option>
            <a-select-option value="asset_proof">资产证明</a-select-option>
            <a-select-option value="business_license">营业执照</a-select-option>
            <a-select-option value="land_certificate">土地证明</a-select-option>
            <a-select-option value="tax_certificate">税收证明</a-select-option>
            <a-select-option value="guarantor_info">担保人信息</a-select-option>
          </a-select>
        </a-form-item>
      </template>

      <!-- 通用备注字段 -->
      <a-form-item 
        :label="getCommentsLabel()" 
        name="comments"
      >
        <a-textarea
          v-model:value="formData.comments"
          :rows="4"
          :placeholder="getCommentsPlaceholder()"
          :maxlength="500"
          show-count
        />
      </a-form-item>

      <!-- 是否通知用户 -->
      <a-form-item name="notify_user">
        <a-checkbox v-model:checked="formData.notify_user">
          立即通知用户审批结果
        </a-checkbox>
      </a-form-item>
    </a-form>

    <!-- AI建议展示 -->
    <div v-if="application?.ai_assessment" class="ai-suggestion-box">
      <div class="ai-suggestion-header">
        <RobotOutlined />
        <span>AI审批建议</span>
      </div>
      <div class="ai-suggestion-content">
        <a-tag :color="getAiSuggestionColor(application.ai_assessment.suggestion)">
          {{ getAiSuggestionText(application.ai_assessment.suggestion) }}
        </a-tag>
        <span class="ai-confidence">
          置信度: {{ (application.ai_assessment.confidence * 100).toFixed(0) }}%
        </span>
        <div class="ai-reason">
          {{ application.ai_assessment.analysis_text }}
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, computed, watch, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { RobotOutlined } from '@ant-design/icons-vue'
import dayjs, { type Dayjs } from 'dayjs'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'
import { loanApi } from '@/api/modules/loan'

interface Props {
  visible: boolean
  application: any
  action: string
}

interface Emits {
  'update:visible': [value: boolean]
  success: []
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

/**
 * 响应式数据
 */
const formRef = ref<FormInstance>()
const loading = ref(false)

const modalVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 表单数据
const formData = reactive({
  // 批准字段
  approved_amount: 0,
  approved_term: 12,
  interest_rate: 6.5,
  repayment_start_date: null as Dayjs | null,
  
  // 拒绝字段
  rejection_reason_code: '',
  
  // 退回字段
  return_reason: '',
  required_materials: [] as string[],
  
  // 通用字段
  comments: '',
  notify_user: true
})

// 表单验证规则
const formRules = computed<Record<string, Rule[]>>(() => {
  const baseRules: Record<string, Rule[]> = {
    comments: [
      { max: 500, message: '备注内容不能超过500字符' }
    ]
  }

  if (props.action === 'approve') {
    return {
      ...baseRules,
      approved_amount: [
        { required: true, message: '请输入批准金额' },
        { type: 'number', min: 1000, message: '批准金额不能少于1000元' },
        { 
          validator: (_, value) => {
            if (value > props.application?.amount) {
              return Promise.reject('批准金额不能超过申请金额')
            }
            return Promise.resolve()
          }
        }
      ],
      approved_term: [
        { required: true, message: '请选择批准期限' }
      ],
      interest_rate: [
        { required: true, message: '请输入年利率' },
        { type: 'number', min: 3, max: 15, message: '年利率应在3%-15%之间' }
      ],
      repayment_start_date: [
        { required: true, message: '请选择还款开始日期' }
      ]
    }
  }

  if (props.action === 'reject') {
    return {
      ...baseRules,
      rejection_reason_code: [
        { required: true, message: '请选择拒绝原因' }
      ],
      comments: [
        { required: true, message: '请填写拒绝说明' },
        { min: 10, message: '拒绝说明至少10个字符' },
        { max: 500, message: '拒绝说明不能超过500字符' }
      ]
    }
  }

  if (props.action === 'return') {
    return {
      ...baseRules,
      return_reason: [
        { required: true, message: '请填写退回原因' },
        { min: 5, message: '退回原因至少5个字符' }
      ],
      required_materials: [
        { required: true, message: '请选择需要补充的材料' }
      ]
    }
  }

  return baseRules
})

/**
 * 计算属性和方法
 */
const getModalTitle = () => {
  const titleMap: Record<string, string> = {
    'approve': '批准贷款申请',
    'reject': '拒绝贷款申请',
    'return': '退回贷款申请'
  }
  return titleMap[props.action] || '操作申请'
}

const getCommentsLabel = () => {
  const labelMap: Record<string, string> = {
    'approve': '批准说明',
    'reject': '拒绝说明',
    'return': '退回说明'
  }
  return labelMap[props.action] || '备注'
}

const getCommentsPlaceholder = () => {
  const placeholderMap: Record<string, string> = {
    'approve': '请填写批准说明（选填）',
    'reject': '请详细说明拒绝原因',
    'return': '请说明退回原因和要求'
  }
  return placeholderMap[props.action] || '请填写备注'
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

const disabledDate = (current: Dayjs) => {
  // 不能选择今天之前的日期
  return current && current < dayjs().startOf('day')
}

/**
 * 事件处理
 */
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true

    const apiData = {
      ...formData,
      repayment_start_date: formData.repayment_start_date?.format('YYYY-MM-DD')
    }

    let response
    if (props.action === 'approve') {
      response = await loanApi.approveApplication(props.application.id, apiData)
    } else if (props.action === 'reject') {
      response = await loanApi.rejectApplication(props.application.id, apiData)
    } else if (props.action === 'return') {
      response = await loanApi.returnApplication(props.application.id, apiData)
    }

    message.success(`${getModalTitle()}成功`)
    emit('success')
    handleCancel()

  } catch (error: any) {
    console.error('审批操作失败:', error)
    message.error(error.message || '操作失败')
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  formRef.value?.resetFields()
  modalVisible.value = false
}

// 监听props变化，初始化表单数据
watch(
  () => [props.visible, props.application, props.action],
  ([visible, application, action]) => {
    if (visible && application) {
      // 重置表单
      formRef.value?.resetFields()
      
      // 根据操作类型初始化表单数据
      if (action === 'approve') {
        formData.approved_amount = application.amount
        formData.approved_term = application.term || 12
        formData.interest_rate = 6.5
        formData.repayment_start_date = dayjs().add(3, 'day')
      }
      
      formData.notify_user = true
      formData.comments = ''
    }
  },
  { immediate: true }
)
</script>

<style lang="scss" scoped>
.ai-suggestion-box {
  margin-top: 16px;
  padding: 16px;
  background: #f6ffed;
  border: 1px solid #b7eb8f;
  border-radius: 6px;
  
  .ai-suggestion-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
    font-weight: 500;
    color: #389e0d;
    
    .anticon {
      font-size: 16px;
    }
  }
  
  .ai-suggestion-content {
    .ai-confidence {
      margin-left: 12px;
      font-size: 12px;
      color: #8c8c8c;
    }
    
    .ai-reason {
      margin-top: 8px;
      font-size: 13px;
      color: #595959;
      line-height: 1.4;
    }
  }
}

// 暗色主题适配
:deep([data-theme="dark"]) {
  .ai-suggestion-box {
    background: #162312 !important;
    border-color: #274916 !important;
    
    .ai-suggestion-header {
      color: #52c41a !important;
    }
    
    .ai-reason {
      color: rgba(255, 255, 255, 0.65) !important;
    }
  }
}
</style> 