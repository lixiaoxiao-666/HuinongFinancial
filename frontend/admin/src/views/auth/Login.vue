<script setup lang="ts">
import { reactive, ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Form, message } from 'ant-design-vue'
import { UserOutlined, LockOutlined, EyeInvisibleOutlined, EyeTwoTone } from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/modules/auth'
import type { LoginCredentials } from '@/api/auth'

/**
 * 组件状态
 */
const loading = ref(false)
const router = useRouter()
const authStore = useAuthStore()

/**
 * 应用版本
 */
const appVersion = computed(() => import.meta.env.VITE_APP_VERSION || '1.0.0')

/**
 * 开发环境检查
 */
const isDevelopment = computed(() => import.meta.env.MODE === 'development')

/**
 * Mock API 检查
 */
const useMockApi = computed(() => import.meta.env.VITE_USE_MOCK === 'true')

/**
 * 环境信息
 */
const envInfo = computed(() => ({
  mode: import.meta.env.MODE,
  useMock: import.meta.env.VITE_USE_MOCK
}))

/**
 * 生成设备信息
 */
const generateDeviceInfo = () => {
  return {
    device_id: `OA_Web_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
    device_type: 'web' as const,
    device_name: `${getBrowserName()} - ${getOperatingSystem()}`,
    user_agent: navigator.userAgent,
    app_version: appVersion.value
  }
}

/**
 * 获取浏览器名称
 */
const getBrowserName = (): string => {
  const userAgent = navigator.userAgent
  if (userAgent.includes('Chrome')) return 'Chrome'
  if (userAgent.includes('Firefox')) return 'Firefox'
  if (userAgent.includes('Safari')) return 'Safari'
  if (userAgent.includes('Edge')) return 'Edge'
  return 'Unknown Browser'
}

/**
 * 获取操作系统
 */
const getOperatingSystem = (): string => {
  const userAgent = navigator.userAgent
  if (userAgent.includes('Windows')) return 'Windows'
  if (userAgent.includes('Mac')) return 'macOS'
  if (userAgent.includes('Linux')) return 'Linux'
  return 'Unknown OS'
}

/**
 * 表单数据
 */
const formState = reactive<LoginCredentials & { remember: boolean }>({
  username: '',
  password: '',
  remember: false,
  platform: 'oa',
  device_info: generateDeviceInfo()
})

/**
 * 表单验证规则
 */
const rules = {
  username: [
    { required: true, message: '请输入用户名或邮箱', trigger: 'blur' },
    { min: 3, message: '用户名长度不能少于3位', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

/**
 * 表单实例
 */
const { validate, validateInfos } = Form.useForm(formState, rules)

/**
 * 登录处理
 */
const handleLogin = async () => {
  try {
    await validate()
    
    loading.value = true
    console.log('🚀 开始登录流程', { username: formState.username })
    
    // 调用auth store的登录方法
    await authStore.login({
      username: formState.username,
      password: formState.password
    })
    
    // 登录成功，跳转到仪表盘
    router.push('/dashboard')
    
  } catch (error: any) {
    console.error('❌ 登录失败:', error)
    
    // 显示错误信息
    if (error.message) {
      message.error(error.message)
    } else {
      message.error('登录失败，请重试')
    }
    
    // 如果是401错误，清空密码字段
    if (error.response?.status === 401) {
      formState.password = ''
    }
  } finally {
    loading.value = false
  }
}

/**
 * 处理忘记密码
 */
const handleForgotPassword = () => {
  message.info('忘记密码功能暂未开放，请联系系统管理员重置密码')
}

/**
 * 处理回车键登录
 */
const handleKeyPress = (event: KeyboardEvent) => {
  if (event.key === 'Enter') {
    handleLogin()
  }
}

/**
 * 组件挂载时初始化
 */
onMounted(() => {
  // 如果已经登录，直接跳转到仪表盘
  if (authStore.isLoggedIn) {
    router.push('/dashboard')
  }
  
  // 开发环境下预填充表单（方便测试）
  if (import.meta.env.MODE === 'development') {
    formState.username = 'admin'
    formState.password = 'admin123'
  }
})
</script>

<template>
  <div class="login-container">
    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="decoration-circle circle-1"></div>
      <div class="decoration-circle circle-2"></div>
      <div class="decoration-circle circle-3"></div>
    </div>

    <!-- 登录卡片 -->
    <div class="login-card">
      <!-- 品牌区域 -->
      <div class="brand-section">
        <div class="brand-logo">
          <div class="logo-icon">
            <UserOutlined />
          </div>
        </div>
        <h1 class="brand-title">惠农OA管理系统</h1>
        <p class="brand-subtitle">智能高效 · 数字化办公</p>
      </div>

      <!-- 表单区域 -->
      <div class="form-section">
        <a-form
          :model="formState"
          @keypress="handleKeyPress"
          autocomplete="off"
        >
          <!-- 用户名输入 -->
          <a-form-item
            v-bind="validateInfos.username"
            class="form-item"
          >
            <a-input
              v-model:value="formState.username"
              size="large"
              placeholder="请输入用户名或邮箱"
              autocomplete="username"
            >
              <template #prefix>
                <UserOutlined class="input-icon" />
              </template>
            </a-input>
          </a-form-item>

          <!-- 密码输入 -->
          <a-form-item
            v-bind="validateInfos.password"
            class="form-item"
          >
            <a-input-password
              v-model:value="formState.password"
              size="large"
              placeholder="请输入密码"
              autocomplete="current-password"
            >
              <template #prefix>
                <LockOutlined class="input-icon" />
              </template>
            </a-input-password>
          </a-form-item>

          <!-- 选项行 -->
          <div class="form-options">
            <a-checkbox v-model:checked="formState.remember">
              记住我
            </a-checkbox>
            <a @click="handleForgotPassword" class="forgot-password">
              忘记密码？
            </a>
          </div>

          <!-- 登录按钮 -->
          <a-form-item class="form-item">
            <a-button
              type="primary"
              size="large"
              :loading="loading"
              :disabled="loading"
              block
              @click="handleLogin"
              class="login-button"
            >
              {{ loading ? '登录中...' : '立即登录' }}
            </a-button>
          </a-form-item>
        </a-form>

        <!-- 底部信息 -->
        <div class="footer-info">
          <!-- 开发环境提示 -->
          <div v-if="isDevelopment" class="dev-hint">
            <p class="dev-title">🚀 开发环境测试</p>
            <p class="dev-info">
              测试账号: <strong>admin</strong> / <strong>admin123</strong>
            </p>
            <p class="dev-info">
              API模式: {{ useMockApi ? '✅ Mock API (推荐)' : '⚠️ 真实后端API' }}
            </p>
            <p class="dev-info">
              环境: {{ envInfo.mode }} | Mock变量: {{ envInfo.useMock }}
            </p>
          </div>
          
          <p class="copyright">
            © 2024 数字惠农金融系统 · 版权所有
          </p>
          <p class="version">
            Version {{ appVersion }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #001529 0%, #1890ff 100%);
  position: relative;
  overflow: hidden;
  padding: 24px;
}

// 背景装饰
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

// 登录卡片
.login-card {
  width: 100%;
  max-width: 420px;
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

// 品牌区域
.brand-section {
  text-align: center;
  padding: 48px 32px 32px;
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
  margin-bottom: 16px;
  
  .logo-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 64px;
    height: 64px;
    background: rgba(255, 255, 255, 0.15);
    border-radius: 50%;
    font-size: 28px;
    backdrop-filter: blur(10px);
    border: 2px solid rgba(255, 255, 255, 0.2);
  }
}

.brand-title {
  font-size: 28px;
  font-weight: 600;
  margin: 0 0 8px;
  background: linear-gradient(135deg, #ffffff 0%, #e6f7ff 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.brand-subtitle {
  font-size: 14px;
  opacity: 0.9;
  margin: 0;
  font-weight: 300;
}

// 表单区域
.form-section {
  padding: 32px;
}

.form-item {
  margin-bottom: 24px;
  
  &:last-child {
    margin-bottom: 0;
  }
  
  :deep(.ant-input) {
    height: 48px;
    border-radius: 8px;
    border-color: #e8e8e8;
    transition: all 0.3s ease;
    
    &:hover {
      border-color: #1890ff;
    }
    
    &:focus {
      border-color: #1890ff;
      box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
    }
  }
  
  :deep(.ant-input-password) {
    .ant-input {
      padding-right: 40px;
    }
  }
}

.input-icon {
  color: rgba(0, 0, 0, 0.4);
  font-size: 16px;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  
  .forgot-password {
    color: #1890ff;
    text-decoration: none;
    font-size: 14px;
    transition: color 0.3s ease;
    
    &:hover {
      color: #40a9ff;
      text-decoration: underline;
    }
  }
}

.login-button {
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

// 底部信息
.footer-info {
  text-align: center;
  padding-top: 24px;
  border-top: 1px solid #f0f0f0;
  margin-top: 24px;
  
  .dev-hint {
    margin-bottom: 16px;
    
    .dev-title {
      font-size: 14px;
      font-weight: 600;
      margin: 0 0 8px;
    }
    
    .dev-info {
      font-size: 12px;
      color: rgba(0, 0, 0, 0.6);
      margin: 0;
    }
  }
  
  .copyright {
    font-size: 12px;
    color: rgba(0, 0, 0, 0.6);
    margin: 0 0 4px;
  }
  
  .version {
    font-size: 11px;
    color: rgba(0, 0, 0, 0.4);
    margin: 0;
  }
}

// 响应式设计
@media (max-width: 480px) {
  .login-container {
    padding: 16px;
  }
  
  .login-card {
    max-width: 100%;
  }
  
  .brand-section {
    padding: 32px 24px 24px;
  }
  
  .brand-title {
    font-size: 24px;
  }
  
  .form-section {
    padding: 24px;
  }
}
</style> 