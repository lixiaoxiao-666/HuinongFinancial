<script setup lang="ts">
// 导入组合式API
import { reactive, ref } from 'vue'
// 导入 element-plus 表单实例和表单规则类型
import type { FormInstance, FormRules } from 'element-plus'
import { useUserStore } from '../../stores/user'
import { useRouter } from 'vue-router'

// 定义ref对象表单实例
const FormRef = ref<FormInstance>()
const router = useRouter()
const userStore = useUserStore()

// 定义表单数据
const FormData = reactive({
  phone: '',
  password: '',
  remember: false
})

// 定义表单验证
const validatePhone = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请输入手机号码'))
  } else {
    callback()
  }
}

const validatePassword = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请输入密码'))
  } else {
    callback()
  }
}

// 定义表单规则
const rules = reactive<FormRules<typeof FormData>>({
    phone: [{ validator: validatePhone, trigger: 'blur' }],
    password: [{ validator: validatePassword, trigger: 'blur' }]
})

// 定义表单提交
const submitForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return
  formEl.validate((valid) => {
    if (valid) {
      console.log('submit!')
      // 这里应该添加登录逻辑
      // 模拟登录成功
      userStore.setToken('mock-token')
      router.push('/')
    } else {
      console.log('error submit!')
    }
  })
}

// 忘记密码
const forgotPassword = () => {
  console.log('忘记密码')
}

// 立即注册
const register = () => {
  console.log('立即注册')
}
</script>

<template>
  <div class="login-container">
    <div class="login-box">
      <div class="logo">
        <img src="../../assets/images/logo.png" alt="数字惠农" />
      </div>
      
      <h1 class="system-title">数字惠农</h1>
      <p class="system-subtitle">金融助力，共富乡村</p>

      <div class="form-container">
        <div class="form-item">
          <div class="form-label">手机号码</div>
          <div class="input-wrapper">
            <el-input 
              v-model="FormData.phone" 
              placeholder="请输入手机号码" 
              prefix-icon="el-icon-phone"
            />
          </div>
        </div>

        <div class="form-item">
          <div class="form-label">密码</div>
          <div class="input-wrapper">
            <el-input
              v-model="FormData.password"
              type="password"
              placeholder="请输入密码"
              prefix-icon="el-icon-lock"
            />
          </div>
        </div>

        <div class="remember-forgot">
          <div class="remember">
            <el-checkbox v-model="FormData.remember">记住我</el-checkbox>
          </div>
          <div class="forgot">
            <a @click="forgotPassword" class="forgot-link">忘记密码？</a>
          </div>
        </div>

        <div class="submit-btn">
          <el-button type="primary" @click="submitForm(FormRef)" class="login-button">
            登录
          </el-button>
        </div>
        
        <div class="register-btn">
          <el-button @click="register" class="register-button">
            立即注册
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>


/* 登录页面 */
.login-container {
  width: 100vw;
  height: 100vh;
  /* 居中 */
  display: flex;
  /* 水平居中 */
  justify-content: center;
  /* 垂直居中 */
  align-items: center;
  /* 渐变背景 */
  background: linear-gradient(to bottom, #f0d080, #D4AF37);
  /* 背景图片 - 麦浪效果 */
  background-image: url('../../assets/images/wheat-field.svg');
  /* 背景图片大小 */
  background-size: cover;
  /* 背景图片位置 */
  background-position: center;
  /* 保持固定 */
  background-attachment: fixed;
  /* 添加过渡效果 */
  animation: bgShimmer 10s infinite alternate;
  position: relative;
  overflow: hidden;
}

@keyframes bgShimmer {
  0% {
    background-color: #f0d080;
  }
  50% {
    background-color: #e6a23c;
  }
  100% {
    background-color: #D4AF37;
  }
}

.login-box {
  /* 宽度 */
  width: 500px;
  /* 背景颜色 - 半透明 */
  background-color: rgba(255, 255, 255, 0.92);
  /* 圆角 */
  border-radius: 10px;
  /* 内边距 */
  padding: 40px;
  /* 文本居中 */
  text-align: center;
  /* 高级阴影效果 */
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1), 0 2px 10px rgba(0, 0, 0, 0.05);
  /* 边框效果 */
  border: 1px solid rgba(255, 255, 255, 0.3);
  /* 玻璃态模糊效果 */
  backdrop-filter: blur(5px);
  /* 动画效果 */
  animation: float 5s infinite ease-in-out;
  /* 相对定位 */
  position: relative;
  z-index: 1;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.logo {
  /* 下边距 */
  margin-bottom: 20px;
  /* 放大 */
  transform: scale(1.5);
}

.logo img {
  /* 宽度 */
  width: 100px;
  /* 高度 */
  height: 100px;
  /* 移除背景色以显示原始图片 */
  background-color: transparent;
  /* 保留原始图片形状 */
  border-radius: 0;
  /* 移除内边距 */
  padding: 0;
}

.system-title {
  /* 金色渐变 */
  background-image: linear-gradient(to right, #D4AF37, #FBD341, #D4AF37);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  /* 字体大小 */
  font-size: 28px;
  /* 下边距 */
  margin-bottom: 10px;
  /* 加粗显示更明显 */
  font-weight: bold;
}

.system-subtitle {
  color: #666;
  font-size: 20px;
  margin-bottom: 30px;
}

.form-container {
  margin: 30px 0;
}

.form-item {
  margin-bottom: 25px;
  text-align: left;
}

.form-label {
  color: #333;
  font-size: 20px;
  margin-bottom: 5px;
}

.input-wrapper {
  width: 100%;
}

.remember-forgot {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 25px;
}

.remember :deep(.el-checkbox__label) {
  color: #333;
  font-size: 20px;
}

.forgot-link {
  background-image: linear-gradient(to right, #e6a23c, #f5c069, #e6a23c);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  cursor: pointer;
  text-decoration: none;
  font-weight: bold;
  font-size: 20px;
}

.login-button {
  width: 100%;
  height: 45px;
  background-image: linear-gradient(to right, #e6a23c, #f5c069, #e6a23c);
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 18px;
  font-weight: bold;
  transition: all 0.3s ease;
}

.login-button:hover {
  background-image: linear-gradient(to right, #d39531, #e6a23c, #d39531);
  border: none;
  box-shadow: 0 4px 8px rgba(214, 146, 32, 0.3);
}

.register-button {
  width: 100%;
  height: 45px;
  background: linear-gradient(to right, #ffffff, #fdf6ec, #ffffff);
  color: #e6a23c;
  border: 1px solid #e6a23c;
  border-radius: 4px;
  font-size: 18px;
  font-weight: bold;
  margin-top: 15px;
  transition: all 0.3s ease;
}

.register-button:hover {
  background: linear-gradient(to right, #fdf6ec, #ffeed8, #fdf6ec);
  border: 1px solid #d39531;  
}

:deep(.el-input__wrapper) {
  /* 增大 */
  height: 60px;
}
</style>