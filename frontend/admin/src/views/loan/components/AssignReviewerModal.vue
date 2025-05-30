<template>
  <a-modal
    v-model:open="modalVisible"
    title="分配审核员"
    width="500"
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
      <!-- 待分配的申请列表 -->
      <div class="assign-info">
        <h4>待分配申请 ({{ applications.length }}个)</h4>
        <div class="application-list">
          <div 
            v-for="app in applications.slice(0, 5)" 
            :key="app.id"
            class="application-item"
          >
            <span class="app-id">{{ app.id }}</span>
            <span class="app-user">{{ app.user_info.real_name }}</span>
            <span class="app-amount">¥{{ formatCurrency(app.amount) }}</span>
          </div>
          <div v-if="applications.length > 5" class="more-applications">
            还有 {{ applications.length - 5 }} 个申请...
          </div>
        </div>
      </div>

      <!-- 审核员选择 -->
      <a-form-item label="选择审核员" name="reviewer_id">
        <a-select
          v-model:value="formData.reviewer_id"
          placeholder="请选择审核员"
          show-search
          :filter-option="filterReviewer"
          @change="handleReviewerChange"
        >
          <a-select-option 
            v-for="reviewer in reviewerList" 
            :key="reviewer.id" 
            :value="reviewer.id"
          >
            <div class="reviewer-option">
              <a-avatar :src="reviewer.avatar" :size="24">
                <template #icon>
                  <UserOutlined />
                </template>
              </a-avatar>
              <div class="reviewer-info">
                <div class="reviewer-name">{{ reviewer.name }}</div>
                <div class="reviewer-meta">
                  {{ reviewer.department }} · 
                  当前处理 {{ reviewer.current_cases || 0 }} 件
                </div>
              </div>
              <div class="reviewer-status">
                <a-tag 
                  :color="reviewer.status === 'online' ? 'green' : 'default'" 
                  size="small"
                >
                  {{ reviewer.status === 'online' ? '在线' : '离线' }}
                </a-tag>
              </div>
            </div>
          </a-select-option>
        </a-select>
      </a-form-item>

      <!-- 选中审核员的详细信息 -->
      <div v-if="selectedReviewer" class="reviewer-detail">
        <h4>审核员信息</h4>
        <a-row :gutter="16">
          <a-col :span="8">
            <div class="detail-item">
              <label>姓名：</label>
              <span>{{ selectedReviewer.name }}</span>
            </div>
          </a-col>
          <a-col :span="8">
            <div class="detail-item">
              <label>部门：</label>
              <span>{{ selectedReviewer.department }}</span>
            </div>
          </a-col>
          <a-col :span="8">
            <div class="detail-item">
              <label>职位：</label>
              <span>{{ selectedReviewer.position }}</span>
            </div>
          </a-col>
          <a-col :span="8">
            <div class="detail-item">
              <label>当前案件：</label>
              <span>{{ selectedReviewer.current_cases || 0 }} 件</span>
            </div>
          </a-col>
          <a-col :span="8">
            <div class="detail-item">
              <label>专业领域：</label>
              <span>{{ selectedReviewer.expertise?.join(', ') || '通用' }}</span>
            </div>
          </a-col>
          <a-col :span="8">
            <div class="detail-item">
              <label>审批效率：</label>
              <a-rate 
                :value="selectedReviewer.efficiency_rating || 0" 
                disabled 
                :count="5" 
                style="font-size: 12px;"
              />
            </div>
          </a-col>
        </a-row>
      </div>

      <!-- 分配原因 -->
      <a-form-item label="分配原因" name="assign_reason">
        <a-select
          v-model:value="formData.assign_reason"
          placeholder="请选择分配原因"
        >
          <a-select-option value="workload_balance">工作负载均衡</a-select-option>
          <a-select-option value="expertise_match">专业领域匹配</a-select-option>
          <a-select-option value="urgent_case">紧急案件处理</a-select-option>
          <a-select-option value="experience_requirement">经验要求</a-select-option>
          <a-select-option value="manual_assignment">手动指定</a-select-option>
        </a-select>
      </a-form-item>

      <!-- 优先级设置 -->
      <a-form-item label="优先级" name="priority">
        <a-radio-group v-model:value="formData.priority">
          <a-radio value="low">低</a-radio>
          <a-radio value="normal">普通</a-radio>
          <a-radio value="high">高</a-radio>
          <a-radio value="urgent">紧急</a-radio>
        </a-radio-group>
      </a-form-item>

      <!-- 截止时间 -->
      <a-form-item label="期望完成时间" name="deadline">
        <a-date-picker
          v-model:value="formData.deadline"
          show-time
          style="width: 100%"
          :disabled-date="disabledDate"
          placeholder="请选择期望完成时间"
          format="YYYY-MM-DD HH:mm"
        />
      </a-form-item>

      <!-- 备注 -->
      <a-form-item label="分配说明" name="notes">
        <a-textarea
          v-model:value="formData.notes"
          :rows="3"
          placeholder="请填写分配说明或特殊要求（选填）"
          :maxlength="200"
          show-count
        />
      </a-form-item>

      <!-- 通知设置 -->
      <a-form-item name="notification_settings">
        <a-space direction="vertical" style="width: 100%;">
          <a-checkbox v-model:checked="formData.send_email">
            发送邮件通知审核员
          </a-checkbox>
          <a-checkbox v-model:checked="formData.send_sms">
            发送短信通知审核员
          </a-checkbox>
          <a-checkbox v-model:checked="formData.create_reminder">
            创建系统提醒
          </a-checkbox>
        </a-space>
      </a-form-item>
    </a-form>

    <!-- 智能推荐 -->
    <div class="smart-recommendation" v-if="recommendedReviewers.length > 0">
      <h4>
        <BulbOutlined />
        智能推荐
      </h4>
      <div class="recommendation-list">
        <div 
          v-for="reviewer in recommendedReviewers" 
          :key="reviewer.id"
          class="recommendation-item"
          @click="selectRecommendedReviewer(reviewer)"
        >
          <a-avatar :src="reviewer.avatar" :size="32">
            <template #icon>
              <UserOutlined />
            </template>
          </a-avatar>
          <div class="recommendation-info">
            <div class="reviewer-name">{{ reviewer.name }}</div>
            <div class="recommendation-reason">{{ reviewer.recommendation_reason }}</div>
          </div>
          <div class="recommendation-score">
            匹配度: {{ (reviewer.match_score * 100).toFixed(0) }}%
          </div>
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, computed, watch, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { UserOutlined, BulbOutlined } from '@ant-design/icons-vue'
import dayjs, { type Dayjs } from 'dayjs'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'
import { loanApi } from '@/api/modules/loan'

interface Props {
  visible: boolean
  applications: any[]
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
const reviewerList = ref<any[]>([])
const recommendedReviewers = ref<any[]>([])

const modalVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 表单数据
const formData = reactive({
  reviewer_id: null as number | null,
  assign_reason: 'workload_balance',
  priority: 'normal',
  deadline: null as Dayjs | null,
  notes: '',
  send_email: true,
  send_sms: false,
  create_reminder: true
})

// 表单验证规则
const formRules: Record<string, Rule[]> = {
  reviewer_id: [
    { required: true, message: '请选择审核员' }
  ],
  assign_reason: [
    { required: true, message: '请选择分配原因' }
  ],
  priority: [
    { required: true, message: '请选择优先级' }
  ],
  notes: [
    { max: 200, message: '分配说明不能超过200字符' }
  ]
}

/**
 * 计算属性
 */
const selectedReviewer = computed(() => {
  return reviewerList.value.find(r => r.id === formData.reviewer_id)
})

/**
 * 工具方法
 */
const formatCurrency = (amount: number) => {
  return (amount / 10000).toFixed(1) + '万'
}

const filterReviewer = (input: string, option: any) => {
  const reviewer = reviewerList.value.find(r => r.id === option.value)
  if (!reviewer) return false
  
  return reviewer.name.toLowerCase().includes(input.toLowerCase()) ||
         reviewer.department.toLowerCase().includes(input.toLowerCase())
}

const disabledDate = (current: Dayjs) => {
  // 不能选择今天之前的日期
  return current && current < dayjs().startOf('day')
}

/**
 * 事件处理
 */
const handleReviewerChange = (reviewerId: number) => {
  // 可以在这里添加选择审核员后的逻辑
}

const selectRecommendedReviewer = (reviewer: any) => {
  formData.reviewer_id = reviewer.id
  formData.assign_reason = 'expertise_match'
  formData.notes = `AI推荐：${reviewer.recommendation_reason}`
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true

    const apiData = {
      ...formData,
      application_ids: props.applications.map(app => app.id),
      deadline: formData.deadline?.format('YYYY-MM-DD HH:mm:ss')
    }

    await loanApi.assignReviewer(apiData)

    message.success(`成功分配 ${props.applications.length} 个申请给审核员`)
    emit('success')
    handleCancel()

  } catch (error: any) {
    console.error('分配审核员失败:', error)
    message.error(error.message || '分配失败')
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  formRef.value?.resetFields()
  modalVisible.value = false
}

// 加载审核员列表
const loadReviewers = async () => {
  try {
    const response = await loanApi.getReviewers()
    reviewerList.value = response.data || []
  } catch (error) {
    console.error('加载审核员列表失败:', error)
  }
}

// 获取智能推荐
const getSmartRecommendations = async () => {
  if (props.applications.length === 0) return

  try {
    const response = await loanApi.getReviewerRecommendations({
      application_ids: props.applications.map(app => app.id)
    })
    recommendedReviewers.value = response.data || []
  } catch (error) {
    console.error('获取智能推荐失败:', error)
  }
}

// 监听props变化
watch(
  () => [props.visible, props.applications],
  ([visible, applications]) => {
    if (visible && applications.length > 0) {
      // 重置表单
      formRef.value?.resetFields()
      
      // 设置默认截止时间（3个工作日后）
      formData.deadline = dayjs().add(3, 'day').hour(18).minute(0).second(0)
      
      // 加载数据
      getSmartRecommendations()
    }
  },
  { immediate: true }
)

/**
 * 生命周期
 */
onMounted(() => {
  loadReviewers()
})
</script>

<style lang="scss" scoped>
.assign-info {
  margin-bottom: 24px;
  
  h4 {
    font-size: 14px;
    font-weight: 600;
    color: #262626;
    margin-bottom: 12px;
  }
  
  .application-list {
    background: #fafafa;
    border-radius: 6px;
    padding: 12px;
    max-height: 120px;
    overflow-y: auto;
    
    .application-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 4px 0;
      font-size: 13px;
      
      .app-id {
        color: #8c8c8c;
        min-width: 80px;
      }
      
      .app-user {
        color: #262626;
        flex: 1;
        margin: 0 8px;
      }
      
      .app-amount {
        color: #52c41a;
        font-weight: 500;
        min-width: 60px;
        text-align: right;
      }
    }
    
    .more-applications {
      text-align: center;
      color: #8c8c8c;
      font-size: 12px;
      margin-top: 8px;
      padding-top: 8px;
      border-top: 1px solid #e8e8e8;
    }
  }
}

.reviewer-option {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0;
  
  .reviewer-info {
    flex: 1;
    
    .reviewer-name {
      font-size: 14px;
      font-weight: 500;
      color: #262626;
    }
    
    .reviewer-meta {
      font-size: 12px;
      color: #8c8c8c;
    }
  }
  
  .reviewer-status {
    margin-left: auto;
  }
}

.reviewer-detail {
  margin: 16px 0;
  padding: 16px;
  background: #fafafa;
  border-radius: 6px;
  
  h4 {
    font-size: 14px;
    font-weight: 600;
    color: #262626;
    margin-bottom: 12px;
  }
  
  .detail-item {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
    font-size: 13px;
    
    label {
      font-weight: 500;
      color: #595959;
      min-width: 70px;
      margin-right: 8px;
    }
  }
}

.smart-recommendation {
  margin-top: 16px;
  padding: 16px;
  background: #f6ffed;
  border: 1px solid #b7eb8f;
  border-radius: 6px;
  
  h4 {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    font-weight: 600;
    color: #389e0d;
    margin-bottom: 12px;
  }
  
  .recommendation-list {
    .recommendation-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 8px 12px;
      border-radius: 4px;
      cursor: pointer;
      transition: background-color 0.3s;
      
      &:hover {
        background: rgba(82, 196, 26, 0.1);
      }
      
      .recommendation-info {
        flex: 1;
        
        .reviewer-name {
          font-size: 14px;
          font-weight: 500;
          color: #262626;
        }
        
        .recommendation-reason {
          font-size: 12px;
          color: #8c8c8c;
          margin-top: 2px;
        }
      }
      
      .recommendation-score {
        font-size: 12px;
        color: #389e0d;
        font-weight: 500;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .assign-info .application-list {
    .application-item {
      flex-direction: column;
      align-items: flex-start;
      gap: 4px;
      
      .app-amount {
        align-self: flex-end;
      }
    }
  }
}

// 暗色主题适配
:deep([data-theme="dark"]) {
  .assign-info .application-list {
    background: #262626 !important;
    
    .application-item {
      .app-user {
        color: rgba(255, 255, 255, 0.85) !important;
      }
    }
    
    .more-applications {
      border-top-color: #303030 !important;
      color: rgba(255, 255, 255, 0.45) !important;
    }
  }
  
  .reviewer-detail {
    background: #262626 !important;
    
    h4 {
      color: rgba(255, 255, 255, 0.85) !important;
    }
    
    .detail-item label {
      color: rgba(255, 255, 255, 0.65) !important;
    }
  }
  
  .smart-recommendation {
    background: #162312 !important;
    border-color: #274916 !important;
    
    h4 {
      color: #52c41a !important;
    }
    
    .recommendation-item {
      .reviewer-name {
        color: rgba(255, 255, 255, 0.85) !important;
      }
      
      .recommendation-reason {
        color: rgba(255, 255, 255, 0.45) !important;
      }
    }
  }
}
</style> 