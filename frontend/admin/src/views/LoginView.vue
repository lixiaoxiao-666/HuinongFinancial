<template>
  <div class="login-container">
    <div class="login-content">
      <div class="login-header">
        <h1>数字惠农OA管理系统</h1>
        <p>Digital Agriculture OA Management System</p>
      </div>
      
      <el-card class="login-card" shadow="always">
        <template #header>
          <div class="card-header">
            <span>管理员登录</span>
          </div>
        </template>
        
        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          size="large"
          @submit.prevent="handleLogin"
        >
          <el-form-item prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="请输入用户名"
              :prefix-icon="User"
              clearable
              autocomplete="username"
            />
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              :prefix-icon="Lock"
              show-password
              clearable
              autocomplete="current-password"
              @keyup.enter="handleLogin"
            />
          </el-form-item>
          
          <el-form-item>
            <div class="login-options">
              <el-checkbox v-model="rememberMe">记住用户名</el-checkbox>
            </div>
          </el-form-item>
          
          <el-form-item>
            <el-button
              type="primary"
              size="large"
              style="width: 100%"
              :loading="loading"
              @click="handleLogin"
            >
              <template #loading>
                <el-icon><Loading /></el-icon>
              </template>
              {{ loading ? '登录中...' : '登录' }}
            </el-button>
          </el-form-item>
        </el-form>
        
        <div class="login-tips" v-if="isDev">
          <h4>🧪 开发环境测试账号</h4>
          <div class="test-accounts">
            <div class="account-item">
              <el-tag type="danger" size="small">超级管理员</el-tag>
              <span>admin / admin123</span>
              <el-button 
                type="primary" 
                link 
                size="small"
                @click="quickLogin('admin', 'admin123')"
              >
                快速登录
              </el-button>
            </div>
            <div class="account-item">
              <el-tag type="warning" size="small">审批员</el-tag>
              <span>reviewer / reviewer123</span>
              <el-button 
                type="primary" 
                link 
                size="small"
                @click="quickLogin('reviewer', 'reviewer123')"
              >
                快速登录
              </el-button>
            </div>
          </div>
        </div>

        <div class="system-info">
          <el-divider>
            <el-icon><InfoFilled /></el-icon>
          </el-divider>
          <el-text size="small" type="info">
            数字惠农OA管理系统采用基于Redis的分布式会话管理，支持多设备登录和实时会话控制。
          </el-text>
          <br>
          <el-text size="small" type="warning">
            请使用企业内网环境访问，确保数据安全。
          </el-text>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { 
  User, 
  Lock, 
  Loading, 
  InfoFilled 
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const loginFormRef = ref<FormInstance>()
const loading = ref(false)
const rememberMe = ref(false)
const isDev = import.meta.env.DEV

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '用户名长度在 2 到 20 个字符', trigger: 'blur' },
    { 
      pattern: /^[a-zA-Z0-9_]{2,20}$/, 
      message: '用户名只能包含字母、数字和下划线', 
      trigger: 'blur' 
    }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    loading.value = true
    
    await authStore.login(loginForm.username, loginForm.password)
    
    if (rememberMe.value) {
      authStore.rememberUsername(loginForm.username)
    } else {
      authStore.clearRememberedUsername()
    }
    
    const redirect = route.query.redirect as string || '/dashboard'
    router.push(redirect)
  } catch (error: any) {
    ElMessage.error(error.message || '登录失败，请检查用户名和密码')
  } finally {
    loading.value = false
  }
}

const quickLogin = (username: string, password: string) => {
  loginForm.username = username
  loginForm.password = password
  handleLogin()
}

onMounted(() => {
  const rememberedUsername = authStore.getRememberedUsername()
  if (rememberedUsername) {
    loginForm.username = rememberedUsername
    rememberMe.value = true
  }
  
  if (authStore.isAuthenticated) {
    const redirect = route.query.redirect as string || '/dashboard'
    router.push(redirect)
  }
})
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  position: relative;
  overflow: hidden;
}

.login-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: 
    radial-gradient(circle at 25% 25%, rgba(255,255,255,0.1) 0%, transparent 50%),
    radial-gradient(circle at 75% 75%, rgba(255,255,255,0.1) 0%, transparent 50%);
  z-index: 1;
}

.login-content {
  width: 100%;
  max-width: 420px;
  position: relative;
  z-index: 2;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
  color: white;
}

.login-header h1 {
  font-size: 32px;
  font-weight: 600;
  margin-bottom: 10px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.login-header p {
  font-size: 16px;
  opacity: 0.9;
  margin: 0;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
}

.login-card {
  border-radius: 20px;
  border: none;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(10px);
  background-color: rgba(255, 255, 255, 0.95);
  overflow: hidden;
  transition: all 0.3s ease;
}

.login-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 25px 70px rgba(0, 0, 0, 0.4);
}

.card-header {
  text-align: center;
  font-size: 22px;
  font-weight: 600;
  color: #333;
  padding: 15px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.login-options {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.login-tips {
  margin-top: 20px;
  padding: 15px;
  background: linear-gradient(45deg, #e8f5e8, #f0f9ff);
  border-radius: 12px;
  border: 1px solid #d4e6f1;
}

.login-tips h4 {
  margin: 0 0 10px 0;
  font-size: 14px;
  color: #2c3e50;
  font-weight: 600;
}

.test-accounts {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.account-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 8px;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 6px;
  font-size: 13px;
}

.account-item span {
  flex: 1;
  margin: 0 10px;
  font-family: 'Courier New', monospace;
  color: #34495e;
}

.system-info {
  margin-top: 20px;
  text-align: center;
  padding: 15px;
  background: rgba(52, 73, 94, 0.05);
  border-radius: 8px;
}

.system-info .el-text {
  display: block;
  margin-bottom: 5px;
  line-height: 1.5;
}

@media (max-width: 480px) {
  .login-content {
    max-width: 100%;
    padding: 0 10px;
  }
  
  .login-header h1 {
    font-size: 28px;
  }
  
  .account-item {
    flex-direction: column;
    gap: 5px;
    align-items: stretch;
  }
  
  .account-item span {
    margin: 0;
    text-align: center;
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-content {
  animation: fadeInUp 0.6s ease-out;
}

@media (prefers-color-scheme: dark) {
  .login-card {
    background-color: rgba(30, 30, 30, 0.95);
    color: #ffffff;
  }
  
  .card-header {
    color: #ffffff;
  }
  
  .system-info {
    background: rgba(255, 255, 255, 0.1);
  }
}
</style> 