import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { pinia } from './stores'

// 导入Ant Design Vue
import Ant from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'

// 导入全局样式
import '@/assets/styles/main.scss'

// 导入Ant Design Vue图标
import * as Icons from '@ant-design/icons-vue'

// 创建应用实例
const app = createApp(App)

// 注册Ant Design Vue图标
Object.keys(Icons).forEach(key => {
  app.component(key, (Icons as any)[key])
})

// 使用插件
app.use(pinia)
app.use(router)
app.use(Ant)

// 挂载应用
app.mount('#app')
