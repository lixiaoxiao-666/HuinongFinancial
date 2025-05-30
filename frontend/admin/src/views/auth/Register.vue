<script setup lang="ts">
import { reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Form, message } from 'ant-design-vue'
import { 
  UserOutlined, 
  LockOutlined, 
  MailOutlined, 
  TeamOutlined,
  HomeOutlined,
  PhoneOutlined,
  EyeInvisibleOutlined, 
  EyeTwoTone 
} from '@ant-design/icons-vue'

/**
 * 组件状态
 */
const loading = ref(false)
const router = useRouter()

/**
 * 表单数据
 */
const formState = reactive({
  username: '',
  real_name: '',
  email: '',
  password: '',
  confirmPassword: '',
  department: '',
  position: '',
  phone: '',
  role: 'operator' as 'admin' | 'manager' | 'operator',
  agree: false
})

/**
 * 部门选项
 */
const departmentOptions = [
  { label: '信贷部', value: '信贷部' },
  { label: '风控部', value: '风控部' },
  { label: '运营部', value: '运营部' },
  { label: '技术部', value: '技术部' },
  { label: '客服部', value: '客服部' },
  { label: '财务部', value: '财务部' },
  { label: '人事部', value: '人事部' },
  { label: '市场部', value: '市场部' }
]

/**
 * 角色选项
 */
const roleOptions = [
  { label: '操作员', value: 'operator' },
  { label: '经理', value: 'manager' },
  { label: '管理员', value: 'admin' }
]

/**
 * 密码强度检测
 */
const passwordStrength = computed(() => {
  const password = formState.password
  if (!password) return { level: 0, text: '' }
  
  let score = 0
  const checks = {
    length: password.length >= 8,
    lowercase: /[a-z]/.test(password),
    uppercase: /[A-Z]/.test(password),
    number: /\d/.test(password),
    special: /[!@#$%^&*(),.?":{}|<>]/.test(password)
  }
  
  score = Object.values(checks).filter(Boolean).length
  
  if (score < 2) return { level: 1, text: '弱', color: '#ff4d4f' }
  if (score < 4) return { level: 2, text: '中', color: '#faad14' }
  return { level: 3, text: '强', color: '#52c41a' }
})

/**
 * 表单验证规则
 */
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度为3-20个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  real_name: [
    { required: true, message: '请输入真实姓名', trigger: 'blur' },
    { min: 2, max: 10, message: '真实姓名长度为2-10个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
    { 
      validator: (_: any, value: string) => {
        if (!value) return Promise.reject('请输入密码')
        if (value.length < 6) return Promise.reject('密码长度不能少于6位')
        return Promise.resolve()
      }, 
      trigger: 'blur' 
    }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { 
      validator: (_: any, value: string) => {
        if (!value) return Promise.reject('请确认密码')
        if (value !== formState.password) return Promise.reject('两次输入的密码不一致')
        return Promise.resolve()
      }, 
      trigger: 'blur' 
    }
  ],
  department: [
    { required: true, message: '请选择部门', trigger: 'change' }
  ],
  position: [
    { required: true, message: '请输入职位', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号码', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ],
  agree: [
    { 
      validator: (_: any, value: boolean) => {
        if (!value) return Promise.reject('请阅读并同意用户协议')
        return Promise.resolve()
      }, 
      trigger: 'change' 
    }
  ]
}

/**
 * 表单实例
 */
const { validate, validateInfos } = Form.useForm(formState, rules)

/**
 * 处理注册
 */
const handleRegister = async () => {
  try {
    // 表单验证
    await validate()
    
    loading.value = true
    
    // 模拟API调用（实际应该调用注册接口）
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    message.success('注册申请已提交，请等待管理员审核！')
    
    // 跳转到登录页
    router.push('/login')
    
  } catch (error: any) {
    console.error('注册失败:', error)
    
    if (error.fields) {
      // 表单验证错误
      const firstError = Object.values(error.fields)[0] as any[]
      message.error(firstError[0]?.message || '请检查输入信息')
    } else {
      // API错误
      message.error(error.message || '注册失败，请重试')
    }
  } finally {
    loading.value = false
  }
}

/**
 * 返回登录页
 */
const backToLogin = () => {
  router.push('/login')
}
</script>

<template>
  <div class="register-container">
    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="decoration-circle circle-1"></div>
      <div class="decoration-circle circle-2"></div>
      <div class="decoration-circle circle-3"></div>
    </div>

    <!-- 注册卡片 -->
    <div class="register-card">
      <!-- 头部区域 -->
      <div class="header-section">
        <div class="brand-logo">
          <TeamOutlined />
        </div>
        <h1 class="page-title">账户注册</h1>
        <p class="page-subtitle">创建新的管理员账户</p>
      </div>

      <!-- 表单区域 -->
      <div class="form-section">
        <a-form
          :model="formState"
          layout="vertical"
          autocomplete="off"
        >
          <!-- 第一行：用户名和真实姓名 -->
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item
                label="用户名"
                v-bind="validateInfos.username"
              >
                <a-input
                  v-model:value="formState.username"
                  placeholder="请输入用户名"
                  autocomplete="username"
                >
                  <template #prefix>
                    <UserOutlined class="input-icon" />
                  </template>
                </a-input>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                label="真实姓名"
                v-bind="validateInfos.real_name"
              >
                <a-input
                  v-model:value="formState.real_name"
                  placeholder="请输入真实姓名"
                >
                  <template #prefix>
                    <UserOutlined class="input-icon" />
                  </template>
                </a-input>
              </a-form-item>
            </a-col>
          </a-row>

          <!-- 第二行：邮箱和手机号 -->
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item
                label="邮箱地址"
                v-bind="validateInfos.email"
              >
                <a-input
                  v-model:value="formState.email"
                  placeholder="请输入邮箱地址"
                  autocomplete="email"
                >
                  <template #prefix>
                    <MailOutlined class="input-icon" />
                  </template>
                </a-input>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                label="手机号码"
                v-bind="validateInfos.phone"
              >
                <a-input
                  v-model:value="formState.phone"
                  placeholder="请输入手机号码"
                >
                  <template #prefix>
                    <PhoneOutlined class="input-icon" />
                  </template>
                </a-input>
              </a-form-item>
            </a-col>
          </a-row>

          <!-- 第三行：部门和职位 -->
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item
                label="所属部门"
                v-bind="validateInfos.department"
              >
                <a-select
                  v-model:value="formState.department"
                  placeholder="请选择部门"
                  :options="departmentOptions"
                >
                  <template #suffixIcon>
                    <HomeOutlined class="input-icon" />
                  </template>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                label="职位"
                v-bind="validateInfos.position"
              >
                <a-input
                  v-model:value="formState.position"
                  placeholder="请输入职位"
                >
                  <template #prefix>
                    <TeamOutlined class="input-icon" />
                  </template>
                </a-input>
              </a-form-item>
            </a-col>
          </a-row>

          <!-- 角色选择 -->
          <a-form-item
            label="角色权限"
            v-bind="validateInfos.role"
          >
            <a-radio-group v-model:value="formState.role" button-style="solid">
              <a-radio-button
                v-for="option in roleOptions"
                :key="option.value"
                :value="option.value"
              >
                {{ option.label }}
              </a-radio-button>
            </a-radio-group>
          </a-form-item>

          <!-- 密码输入 -->
          <a-form-item
            label="登录密码"
            v-bind="validateInfos.password"
          >
            <a-input-password
              v-model:value="formState.password"
              placeholder="请输入密码"
              autocomplete="new-password"
              :icon-render="(visible: boolean) => visible ? EyeTwoTone : EyeInvisibleOutlined"
            >
              <template #prefix>
                <LockOutlined class="input-icon" />
              </template>
            </a-input-password>
            <!-- 密码强度指示器 -->
            <div v-if="formState.password" class="password-strength">
              <div class="strength-bar">
                <div 
                  class="strength-fill"
                  :style="{
                    width: `${(passwordStrength.level / 3) * 100}%`,
                    backgroundColor: passwordStrength.color
                  }"
                ></div>
              </div>
              <span 
                class="strength-text"
                :style="{ color: passwordStrength.color }"
              >
                密码强度：{{ passwordStrength.text }}
              </span>
            </div>
          </a-form-item>

          <!-- 确认密码 -->
          <a-form-item
            label="确认密码"
            v-bind="validateInfos.confirmPassword"
          >
            <a-input-password
              v-model:value="formState.confirmPassword"
              placeholder="请再次输入密码"
              autocomplete="new-password"
              :icon-render="(visible: boolean) => visible ? EyeTwoTone : EyeInvisibleOutlined"
            >
              <template #prefix>
                <LockOutlined class="input-icon" />
              </template>
            </a-input-password>
          </a-form-item>

          <!-- 协议同意 -->
          <a-form-item v-bind="validateInfos.agree">
            <a-checkbox v-model:checked="formState.agree">
              我已阅读并同意
              <a href="#" class="agreement-link">《用户服务协议》</a>
              和
              <a href="#" class="agreement-link">《隐私政策》</a>
            </a-checkbox>
          </a-form-item>

          <!-- 提交按钮 -->
          <a-form-item>
            <a-button
              type="primary"
              size="large"
              :loading="loading"
              :disabled="loading"
              block
              @click="handleRegister"
              class="register-button"
            >
              {{ loading ? '注册中...' : '立即注册' }}
            </a-button>
          </a-form-item>

          <!-- 返回登录 -->
          <div class="back-to-login">
            <span>已有账户？</span>
            <a @click="backToLogin" class="login-link">立即登录</a>
          </div>
        </a-form>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #001529 0%, #1890ff 100%);
  position: relative;
  overflow: hidden;
  padding: 24px;
}

// 背景装饰（与登录页相同）
.background-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  
  .decoration-circle {
    position: absolute;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.05);
    animation: float 20s infinite ease-in-out;
    
    &.circle-1 {
      width: 300px;
      height: 300px;
      top: -150px;
      left: -150px;
      animation-delay: 0s;
    }
    
    &.circle-2 {
      width: 200px;
      height: 200px;
      top: 50%;
      right: -100px;
      animation-delay: -8s;
    }
    
    &.circle-3 {
      width: 400px;
      height: 400px;
      bottom: -200px;
      left: 50%;
      transform: translateX(-50%);
      animation-delay: -16s;
    }
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  33% {
    transform: translateY(-30px) rotate(120deg);
  }
  66% {
    transform: translateY(30px) rotate(240deg);
  }
}

// 注册卡片
.register-card {
  width: 100%;
  max-width: 700px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
  animation: slideUp 0.8s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

// 头部区域
.header-section {
  text-align: center;
  padding: 32px 32px 24px;
  background: linear-gradient(135deg, #001529 0%, #1890ff 100%);
  color: white;
  position: relative;
  
  &::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  }
}

.brand-logo {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 50%;
  font-size: 24px;
  margin-bottom: 16px;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.2);
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 8px;
  background: linear-gradient(135deg, #ffffff 0%, #e6f7ff 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-subtitle {
  font-size: 14px;
  opacity: 0.9;
  margin: 0;
  font-weight: 300;
}

// 表单区域
.form-section {
  padding: 32px;
  max-height: 70vh;
  overflow-y: auto;
}

:deep(.ant-form-item) {
  margin-bottom: 20px;
  
  .ant-form-item-label {
    font-weight: 500;
    color: $text-color;
    
    > label {
      font-size: 14px;
      color: $text-color;
    }
  }
}

:deep(.ant-input),
:deep(.ant-input-password),
:deep(.ant-select-selector) {
  height: 40px;
  border-radius: 6px;
  border-color: #e8e8e8;
  transition: all 0.3s ease;
  
  &:hover {
    border-color: $primary-color;
  }
  
  &:focus,
  &.ant-input-focused,
  &.ant-select-focused {
    border-color: $primary-color;
    box-shadow: 0 0 0 2px rgba($primary-color, 0.2);
  }
}

.input-icon {
  color: rgba(0, 0, 0, 0.4);
  font-size: 14px;
}

// 密码强度指示器
.password-strength {
  margin-top: 8px;
  
  .strength-bar {
    width: 100%;
    height: 4px;
    background-color: #f0f0f0;
    border-radius: 2px;
    overflow: hidden;
    margin-bottom: 4px;
    
    .strength-fill {
      height: 100%;
      transition: all 0.3s ease;
      border-radius: 2px;
    }
  }
  
  .strength-text {
    font-size: 12px;
    font-weight: 500;
  }
}

// 角色选择
:deep(.ant-radio-group) {
  width: 100%;
  
  .ant-radio-button-wrapper {
    height: 40px;
    line-height: 38px;
    border-radius: 6px;
    margin-right: 8px;
    border-color: #e8e8e8;
    
    &:first-child {
      border-radius: 6px;
    }
    
    &:last-child {
      border-radius: 6px;
      margin-right: 0;
    }
    
    &.ant-radio-button-wrapper-checked {
      background: $primary-color;
      border-color: $primary-color;
    }
  }
}

// 协议链接
.agreement-link {
  color: $primary-color;
  text-decoration: none;
  
  &:hover {
    color: $primary-color-hover;
    text-decoration: underline;
  }
}

// 注册按钮
.register-button {
  height: 48px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  background: linear-gradient(135deg, #1890ff 0%, #001529 100%);
  border: none;
  transition: all 0.3s ease;
  
  &:hover {
    background: linear-gradient(135deg, #40a9ff 0%, #1890ff 100%);
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(24, 144, 255, 0.3);
  }
  
  &:active {
    transform: translateY(0);
  }
  
  &.ant-btn-loading {
    background: linear-gradient(135deg, #1890ff 0%, #001529 100%);
  }
}

// 返回登录
.back-to-login {
  text-align: center;
  margin-top: 24px;
  font-size: 14px;
  color: rgba(0, 0, 0, 0.6);
  
  .login-link {
    color: $primary-color;
    text-decoration: none;
    font-weight: 500;
    margin-left: 8px;
    
    &:hover {
      color: $primary-color-hover;
      text-decoration: underline;
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .register-container {
    padding: 16px;
  }
  
  .register-card {
    max-width: 100%;
  }
  
  .header-section {
    padding: 24px 24px 20px;
  }
  
  .page-title {
    font-size: 20px;
  }
  
  .form-section {
    padding: 24px;
    max-height: 75vh;
  }
  
  :deep(.ant-col) {
    margin-bottom: 0;
  }
  
  :deep(.ant-radio-group) {
    .ant-radio-button-wrapper {
      margin-bottom: 8px;
      margin-right: 4px;
    }
  }
}
</style> 