import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { setupRouterGuards } from './router/guards'
import { useAuthStore } from './stores/modules/auth'

// 样式
import './assets/styles/main.scss'

// Ant Design Vue
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'

// dayjs
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

// 配置dayjs
dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

// 创建应用实例
const app = createApp(App)

// 安装插件
app.use(createPinia())
app.use(router)
app.use(Antd)

// 设置路由守卫
setupRouterGuards(router)

// 全局错误处理
app.config.errorHandler = (error, instance, info) => {
  console.error('Global error:', error, info)
  // 这里可以添加错误上报逻辑
}

// 全局警告处理
app.config.warnHandler = (msg, instance, trace) => {
  console.warn('Global warning:', msg, trace)
}

// 应用启动
const initApp = async () => {
  try {
    // 初始化认证状态
    const authStore = useAuthStore()
    await authStore.initAuth()
    
    console.log('✅ 应用初始化完成')
    
    // 挂载应用
    app.mount('#app')
  } catch (error) {
    console.error('❌ 应用初始化失败:', error)
    
    // 即使初始化失败也要挂载应用，让用户可以登录
    app.mount('#app')
  }
}

// 启动应用
initApp()
