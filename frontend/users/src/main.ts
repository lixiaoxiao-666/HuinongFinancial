import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { useUserStore } from './stores/user'
import { authManager } from './utils/auth'


const app = createApp(App)

app.use(ElementPlus)
const pinia = createPinia()
app.use(pinia)

// 注册 Element Plus 图标组件
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}
app.use(router)

// 恢复用户状态
const userStore = useUserStore()
userStore.restoreFromStorage()

// 设置认证管理器的路由实例
authManager.setRouter(router)

// 应用启动时初始化认证状态
authManager.initializeAuth().then((isAuthenticated) => {
  console.log('认证状态初始化完成:', isAuthenticated)
}).catch((error) => {
  console.error('认证状态初始化失败:', error)
})

app.mount('#app')
