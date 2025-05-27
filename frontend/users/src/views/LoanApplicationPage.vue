<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules, UploadFile } from 'element-plus'
import { useUserStore } from '../stores/user'
import { loanApi, fileApi } from '../services/api'
import type { LoanProduct, LoanApplicationRequest, FileUploadResult } from '../services/api'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)
const uploading = ref(false)

// 贷款产品信息
const loanProduct = ref<LoanProduct | null>(null)

// 申请表单数据
const formData = reactive({
  amount: 0,
  term_months: 12,
  purpose: '',
  applicant_info: {
    real_name: '',
    id_card_number: '',
    address: ''
  }
})

// 文档上传
const uploadedFiles = ref<Record<string, FileUploadResult>>({})
const requiredDocs = computed(() => loanProduct.value?.required_documents || [])

// 表单验证规则
const rules = reactive<FormRules>({
  amount: [
    { required: true, message: '请输入贷款金额', trigger: 'blur' },
    { 
      validator: (rule: any, value: number, callback: any) => {
        if (!loanProduct.value) {
          callback()
          return
        }
        if (value < loanProduct.value.min_amount) {
          callback(new Error(`贷款金额不能少于${loanProduct.value.min_amount}元`))
        } else if (value > loanProduct.value.max_amount) {
          callback(new Error(`贷款金额不能超过${loanProduct.value.max_amount}元`))
        } else {
          callback()
        }
      }, 
      trigger: 'blur' 
    }
  ],
  term_months: [
    { required: true, message: '请选择贷款期限', trigger: 'change' }
  ],
  purpose: [
    { required: true, message: '请输入贷款用途', trigger: 'blur' },
    { min: 10, message: '贷款用途不能少于10个字符', trigger: 'blur' }
  ],
  'applicant_info.real_name': [
    { required: true, message: '请输入真实姓名', trigger: 'blur' }
  ],
  'applicant_info.id_card_number': [
    { required: true, message: '请输入身份证号码', trigger: 'blur' },
    { 
      pattern: /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/, 
      message: '请输入正确的身份证号码', 
      trigger: 'blur' 
    }
  ],
  'applicant_info.address': [
    { required: true, message: '请输入详细地址', trigger: 'blur' }
  ]
})

// 期限选项
const termOptions = computed(() => {
  if (!loanProduct.value) return []
  const options = []
  for (let i = loanProduct.value.min_term_months; i <= loanProduct.value.max_term_months; i += 6) {
    options.push({
      value: i,
      label: `${i}个月`
    })
  }
  return options
})

// 加载贷款产品信息
const loadProductInfo = async () => {
  const productId = route.params.productId as string
  if (!productId) {
    ElMessage.error('缺少产品信息')
    router.go(-1)
    return
  }

  try {
    const response = await loanApi.getProductDetail(productId)
    loanProduct.value = response.data
    
    // 设置默认值
    formData.amount = loanProduct.value.min_amount
    formData.term_months = loanProduct.value.min_term_months
  } catch (error: any) {
    console.error('加载产品信息失败:', error)
    ElMessage.error('加载产品信息失败')
    router.go(-1)
  }
}

// 文件上传处理
const handleFileUpload = async (file: File, docType: string) => {
  try {
    uploading.value = true
    const response = await fileApi.upload(file, 'loan_document')
    uploadedFiles.value[docType] = response.data
    ElMessage.success('文件上传成功')
  } catch (error: any) {
    console.error('文件上传失败:', error)
    ElMessage.error(error.message || '文件上传失败')
  } finally {
    uploading.value = false
  }
}

// 文件上传前的验证
const beforeUpload = (file: File) => {
  const isValidType = ['image/jpeg', 'image/png', 'image/jpg', 'application/pdf'].includes(file.type)
  const isLt10M = file.size / 1024 / 1024 < 10

  if (!isValidType) {
    ElMessage.error('只能上传 JPG/PNG/PDF 格式的文件!')
    return false
  }
  if (!isLt10M) {
    ElMessage.error('上传文件大小不能超过 10MB!')
    return false
  }
  return true
}

// 自定义上传
const customUpload = (options: any, docType: string) => {
  const file = options.file
  if (beforeUpload(file)) {
    handleFileUpload(file, docType)
  }
  return false // 阻止默认上传行为
}

// 删除文件
const removeFile = (docType: string) => {
  delete uploadedFiles.value[docType]
  ElMessage.success('文件已删除')
}

// 检查是否所有必需文档都已上传
const allDocsUploaded = computed(() => {
  return requiredDocs.value.every(doc => uploadedFiles.value[doc.type])
})

// 提交申请
const submitApplication = async () => {
  if (!formRef.value) return
  
  try {
    // 验证表单
    const valid = await formRef.value.validate()
    if (!valid) return
    
    // 检查文档上传
    if (!allDocsUploaded.value) {
      ElMessage.error('请上传所有必需的文档')
      return
    }
    
    // 确认提交
    await ElMessageBox.confirm(
      '请确认申请信息无误，提交后将进入审核流程',
      '确认提交',
      {
        confirmButtonText: '确认提交',
        cancelButtonText: '再检查一下',
        type: 'warning'
      }
    )
    
    loading.value = true
    
    // 构建申请数据
    const applicationData: LoanApplicationRequest = {
      product_id: loanProduct.value!.product_id,
      amount: formData.amount,
      term_months: formData.term_months,
      purpose: formData.purpose,
      applicant_info: formData.applicant_info,
      uploaded_documents: Object.entries(uploadedFiles.value).map(([docType, fileInfo]) => ({
        doc_type: docType,
        file_id: fileInfo.file_id
      }))
    }
    
    // 提交申请
    const response = await loanApi.submitApplication(applicationData)
    
    ElMessage.success('申请提交成功')
    
    // 跳转到申请详情页
    router.push(`/loan/application/${response.data.application_id}`)
    
  } catch (error: any) {
    if (error.message !== 'cancel') {
      console.error('提交申请失败:', error)
      ElMessage.error(error.message || '提交申请失败')
    }
  } finally {
    loading.value = false
  }
}

// 返回上一页
const goBack = () => {
  router.go(-1)
}

// 组件挂载时加载数据
onMounted(() => {
  // 检查登录状态
  if (!userStore.isLoggedIn) {
    ElMessage.error('请先登录')
    router.push('/login')
    return
  }
  
  loadProductInfo()
  
  // 从用户信息中预填充申请人信息
  if (userStore.userInfo) {
    formData.applicant_info.real_name = userStore.userInfo.real_name || ''
    formData.applicant_info.id_card_number = userStore.userInfo.id_card_number || ''
    formData.applicant_info.address = userStore.userInfo.address || ''
  }
})
</script>

<template>
  <div class="loan-application-page">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <div class="nav-left" @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </div>
      <div class="nav-title">贷款申请</div>
      <div class="nav-right"></div>
    </div>

    <div class="page-content" v-if="loanProduct">
      <!-- 产品信息卡片 -->
      <div class="product-card">
        <h3 class="product-name">{{ loanProduct.name }}</h3>
        <p class="product-desc">{{ loanProduct.description }}</p>
        <div class="product-info">
          <div class="info-item">
            <span class="label">贷款金额：</span>
            <span class="value">{{ loanProduct.min_amount.toLocaleString() }} - {{ loanProduct.max_amount.toLocaleString() }}元</span>
          </div>
          <div class="info-item">
            <span class="label">贷款期限：</span>
            <span class="value">{{ loanProduct.min_term_months }} - {{ loanProduct.max_term_months }}个月</span>
          </div>
          <div class="info-item">
            <span class="label">年利率：</span>
            <span class="value">{{ loanProduct.interest_rate_yearly }}</span>
          </div>
        </div>
      </div>

      <!-- 申请表单 -->
      <div class="form-container">
        <el-form
          ref="formRef"
          :model="formData"
          :rules="rules"
          label-width="100px"
          class="application-form"
        >
          <el-form-item label="贷款金额" prop="amount">
            <el-input-number
              v-model="formData.amount"
              :min="loanProduct.min_amount"
              :max="loanProduct.max_amount"
              :step="1000"
              controls-position="right"
              style="width: 100%"
            />
            <div class="form-tip">可申请金额：{{ loanProduct.min_amount.toLocaleString() }} - {{ loanProduct.max_amount.toLocaleString() }}元</div>
          </el-form-item>

          <el-form-item label="贷款期限" prop="term_months">
            <el-select v-model="formData.term_months" style="width: 100%">
              <el-option
                v-for="option in termOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="贷款用途" prop="purpose">
            <el-input
              v-model="formData.purpose"
              type="textarea"
              :rows="4"
              placeholder="请详细说明贷款用途，不少于10个字符"
              maxlength="200"
              show-word-limit
            />
          </el-form-item>

          <div class="section-title">申请人信息</div>
          
          <el-form-item label="真实姓名" prop="applicant_info.real_name">
            <el-input v-model="formData.applicant_info.real_name" placeholder="请输入真实姓名" />
          </el-form-item>

          <el-form-item label="身份证号" prop="applicant_info.id_card_number">
            <el-input v-model="formData.applicant_info.id_card_number" placeholder="请输入身份证号码" />
          </el-form-item>

          <el-form-item label="详细地址" prop="applicant_info.address">
            <el-input v-model="formData.applicant_info.address" placeholder="请输入详细地址" />
          </el-form-item>
        </el-form>
      </div>

      <!-- 文档上传 -->
      <div class="upload-container">
        <div class="section-title">上传申请材料</div>
        <div class="upload-list">
          <div 
            v-for="doc in requiredDocs" 
            :key="doc.type" 
            class="upload-item"
          >
            <div class="doc-info">
              <h4 class="doc-title">{{ doc.desc }}</h4>
              <p class="doc-tip">支持 JPG、PNG、PDF 格式，大小不超过10MB</p>
            </div>
            
            <div v-if="uploadedFiles[doc.type]" class="uploaded-file">
              <div class="file-info">
                <el-icon class="file-icon"><Document /></el-icon>
                <span class="file-name">{{ uploadedFiles[doc.type].file_name }}</span>
                <el-button 
                  type="danger" 
                  size="small" 
                  @click="removeFile(doc.type)"
                  :icon="Delete"
                >
                  删除
                </el-button>
              </div>
            </div>
            
            <div v-else class="upload-area">
              <el-upload
                :http-request="(options) => customUpload(options, doc.type)"
                :show-file-list="false"
                accept=".jpg,.jpeg,.png,.pdf"
                class="upload-component"
              >
                <el-button type="primary" :loading="uploading">
                  <el-icon><Upload /></el-icon>
                  {{ uploading ? '上传中...' : '选择文件' }}
                </el-button>
              </el-upload>
            </div>
          </div>
        </div>
      </div>

      <!-- 提交按钮 -->
      <div class="submit-container">
        <el-button 
          type="primary" 
          size="large" 
          @click="submitApplication"
          :loading="loading"
          :disabled="!allDocsUploaded"
          class="submit-btn"
        >
          {{ loading ? '提交中...' : '提交申请' }}
        </el-button>
        <p class="submit-tip">
          <el-icon><InfoFilled /></el-icon>
          提交后将进入AI智能审核，通常1-3个工作日内完成审核
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.loan-application-page {
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

.product-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.product-name {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 8px;
}

.product-desc {
  color: #7f8c8d;
  margin-bottom: 16px;
  line-height: 1.6;
}

.product-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-item {
  display: flex;
  align-items: center;
}

.label {
  color: #7f8c8d;
  width: 80px;
}

.value {
  color: #2c3e50;
  font-weight: 500;
}

.form-container {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 20px 0 16px 0;
  padding-bottom: 8px;
  border-bottom: 2px solid #27ae60;
}

.form-tip {
  font-size: 12px;
  color: #7f8c8d;
  margin-top: 4px;
}

.upload-container {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.upload-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.upload-item {
  border: 1px solid #e1e1e1;
  border-radius: 8px;
  padding: 16px;
}

.doc-info {
  margin-bottom: 12px;
}

.doc-title {
  font-size: 14px;
  font-weight: 500;
  color: #2c3e50;
  margin-bottom: 4px;
}

.doc-tip {
  font-size: 12px;
  color: #7f8c8d;
  margin: 0;
}

.uploaded-file {
  display: flex;
  align-items: center;
  gap: 12px;
}

.file-info {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.file-icon {
  color: #27ae60;
}

.file-name {
  color: #2c3e50;
  font-size: 14px;
}

.upload-area {
  display: flex;
  justify-content: center;
}

.submit-container {
  text-align: center;
  padding: 20px 0;
}

.submit-btn {
  width: 100%;
  height: 50px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 25px;
  margin-bottom: 12px;
}

.submit-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  color: #7f8c8d;
  font-size: 12px;
  margin: 0;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #2c3e50;
}

:deep(.el-input-number) {
  width: 100%;
}

:deep(.el-upload) {
  width: 100%;
}
</style> 