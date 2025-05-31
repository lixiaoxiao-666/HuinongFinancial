import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { useUserStore } from './stores/user'


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

app.mount('#app')
