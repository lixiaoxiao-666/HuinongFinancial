import axios from 'axios'
import type { AxiosResponse, AxiosError } from 'axios'
import { message } from 'ant-design-vue'
import { useAuthStore } from '@/stores/modules/auth'
import router from '@/router'

// 创建axios实例
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/oa',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    
    // 添加认证头
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    
    // 添加平台信息
    config.headers['X-Platform'] = 'oa'
    config.headers['X-Device-Type'] = 'web'
    config.headers['X-App-Version'] = import.meta.env.VITE_APP_VERSION || '1.0.0'
    
    return config
  },
  (error) => {
    console.error('请求配置错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const { code, message: msg, data } = response.data
    
    // 统一处理业务错误
    if (code !== 200) {
      message.error(msg || '操作失败')
      return Promise.reject(new Error(msg || '操作失败'))
    }
    
    return response.data
  },
  (error: AxiosError) => {
    console.error('响应错误:', error)
    
    // 处理HTTP状态码错误
    if (error.response) {
      const { status, data } = error.response
      
      switch (status) {
        case 401:
          message.error('认证失败，请重新登录')
          const authStore = useAuthStore()
          authStore.logout()
          router.push('/login')
          break
          
        case 403:
          message.error('权限不足')
          break
          
        case 404:
          message.error('请求的资源不存在')
          break
          
        case 500:
          message.error('服务器内部错误')
          break
          
        default:
          message.error((data as any)?.message || `请求失败 (${status})`)
      }
    } else if (error.request) {
      message.error('网络连接失败，请检查网络设置')
    } else {
      message.error('请求配置错误')
    }
    
    return Promise.reject(error)
  }
)

export { request }
export default request 