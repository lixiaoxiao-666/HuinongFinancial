<template>
  <div class="login-container">
    <div class="login-content">
      <div class="login-left">
        <div class="illustration-container">
          <img src="../assets/login-illustration.svg" alt="系统插图" class="login-illustration" />
        </div>
      </div>
      <div class="login-right">
        <div class="login-form-container">
          <div class="login-header">
            <div class="system-logo"></div>
            <h1 class="login-title">数字惠农后台管理系统</h1>
            <div class="login-subtitle">DIGITAL AGRICULTURE FINANCIAL ADMIN</div>
          </div>
          <div class="login-form">
            <div class="form-item">
              <div class="form-icon">
                <svg class="icon" aria-hidden="true">
                  <use xlink:href="#icon-user"></use>
                </svg>
              </div>
              <input 
                type="text" 
                class="form-input" 
                placeholder="请输入手机号或账号" 
                v-model="loginForm.username"
              />
            </div>
            <div class="form-item">
              <div class="form-icon">
                <svg class="icon" aria-hidden="true">
                  <use xlink:href="#icon-lock"></use>
                </svg>
              </div>
              <input 
                type="password" 
                class="form-input" 
                placeholder="请输入密码" 
                v-model="loginForm.password"
              />
              <div class="password-toggle" @click="togglePasswordVisibility">
                <svg class="icon" aria-hidden="true">
                  <use xlink:href="#icon-eye"></use>
                </svg>
              </div>
            </div>
            <button class="login-button" @click="handleLogin">
              <span>登 录</span>
            </button>
            <div class="login-options">
              <span class="vendor-login">商户登录</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="login-footer">
      <p>技术支持: 18887123</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loginForm = reactive({
  username: '',
  password: '',
})

const passwordVisible = ref(false)

const togglePasswordVisibility = () => {
  const passwordInput = document.querySelector('input[type="password"]') as HTMLInputElement
  if (passwordInput) {
    passwordVisible.value = !passwordVisible.value
    passwordInput.type = passwordVisible.value ? 'text' : 'password'
  }
}

const handleLogin = () => {
  // 简单的验证
  if (!loginForm.username || !loginForm.password) {
    alert('请填写完整登录信息')
    return
  }
  
  // 模拟存储token
  localStorage.setItem('token', 'demo-token')
  
  // 登录成功后跳转到首页
  router.push('/home')
}
</script>

<style scoped>
.login-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100%;
  background: linear-gradient(135deg, #f0f8ff 0%, #e6f3ff 100%);
  overflow: hidden;
  position: relative;
}

.login-content {
  display: flex;
  flex: 1;
  width: 100%;
  position: relative;
  z-index: 1;
}

.login-left {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  position: relative;
}

.illustration-container {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.login-illustration {
  max-width: 100%;
  max-height: 80%;
  object-fit: contain;
}

.login-right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-form-container {
  width: 400px;
  padding: 40px;
  background-color: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  position: relative;
  overflow: hidden;
}

.system-logo {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #4285f4 0%, #34a853 100%);
  border-radius: 12px;
  margin: 0 auto 16px;
  position: relative;
}

.system-logo::before {
  content: '';
  position: absolute;
  top: 15px;
  left: 15px;
  right: 15px;
  bottom: 15px;
  background-color: rgba(255, 255, 255, 0.9);
  border-radius: 6px;
}

.system-logo::after {
  content: '';
  position: absolute;
  top: 22px;
  left: 22px;
  width: 16px;
  height: 16px;
  border-radius: 3px;
  background-color: #4285f4;
}

.login-header {
  margin-bottom: 40px;
  text-align: center;
}

.login-title {
  font-size: 24px;
  color: #333;
  font-weight: bold;
  margin: 0;
}

.login-subtitle {
  margin-top: 8px;
  font-size: 12px;
  color: #999;
  letter-spacing: 1px;
}

.form-item {
  display: flex;
  align-items: center;
  background-color: #f5f7fa;
  border-radius: 8px;
  padding: 12px 16px;
  margin-bottom: 20px;
  position: relative;
}

.form-icon {
  margin-right: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
}

.form-input {
  flex: 1;
  border: none;
  background: transparent;
  outline: none;
  font-size: 14px;
  color: #333;
}

.form-input::placeholder {
  color: #999;
}

.password-toggle {
  cursor: pointer;
  color: #999;
}

.login-button {
  width: 100%;
  height: 44px;
  border-radius: 8px;
  border: none;
  background: linear-gradient(90deg, #4285f4 0%, #34a853 100%);
  color: white;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  margin-top: 10px;
}

.login-options {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.vendor-login {
  color: #4285f4;
  cursor: pointer;
  font-size: 14px;
}

.login-footer {
  padding: 20px 0;
  text-align: center;
  color: #999;
  font-size: 14px;
  position: relative;
  z-index: 1;
}

.icon {
  width: 16px;
  height: 16px;
}

@media (max-width: 768px) {
  .login-content {
    flex-direction: column;
  }
  
  .login-left {
    display: none;
  }
  
  .login-form-container {
    width: 90%;
    max-width: 400px;
  }
}
</style> 