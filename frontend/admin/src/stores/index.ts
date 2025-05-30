import { createPinia } from 'pinia'

// 创建pinia实例
export const pinia = createPinia()

// 导出store模块
export { useAuthStore } from './modules/auth'
export { useAppStore } from './modules/app'

export default pinia 