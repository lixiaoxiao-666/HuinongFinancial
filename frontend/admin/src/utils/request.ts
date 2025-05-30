import axios, { 
  type AxiosInstance, 
  type AxiosRequestConfig, 
  type AxiosResponse,
  type InternalAxiosRequestConfig
} from 'axios'
import { ElMessage, ElLoading } from 'element-plus'
import router from '@/router'
import { useAuthStore } from '@/stores/auth'
import type { ApiResponse } from '@/types/auth'

// 创建axios实例
const request: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 加载状态管理
let loadingInstance: any = null
let requestCount = 0

// 显示加载状态
const showLoading = () => {
  if (requestCount === 0) {
    loadingInstance = ElLoading.service({
      text: '加载中...',
      background: 'rgba(0, 0, 0, 0.7)'
    })
  }
  requestCount++
}

// 隐藏加载状态
const hideLoading = () => {
  requestCount--
  if (requestCount <= 0) {
    requestCount = 0
    if (loadingInstance) {
      loadingInstance.close()
      loadingInstance = null
    }
  }
}

// 请求拦截器
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 显示加载状态（排除某些不需要loading的请求）
    const silentRequests = ['/api/auth/validate', '/api/auth/refresh']
    if (!silentRequests.some(url => config.url?.includes(url))) {
      showLoading()
    }

    // 公开接口列表（不需要token认证）
    const publicEndpoints = [
      '/api/oa/auth/login',
      '/api/auth/register',
      '/api/auth/login',
      '/api/auth/send-sms',
      '/api/auth/reset-password',
      '/api/public/'
    ]

    // 检查是否为公开接口
    const isPublicEndpoint = publicEndpoints.some(endpoint => 
      config.url?.includes(endpoint)
    )

    // 只对非公开接口添加token
    if (!isPublicEndpoint) {
      const authStore = useAuthStore()
      if (authStore.accessToken) {
        config.headers.Authorization = `Bearer ${authStore.accessToken}`
      }
    }

    // 添加设备信息
    const getDeviceName = () => {
      const ua = navigator.userAgent
      if (ua.includes('Edg/')) return 'Microsoft Edge'
      if (ua.includes('Chrome/')) return 'Google Chrome'
      if (ua.includes('Firefox/')) return 'Mozilla Firefox'
      if (ua.includes('Safari/') && !ua.includes('Chrome/')) return 'Apple Safari'
      if (ua.includes('Opera/')) return 'Opera'
      return 'Unknown Browser'
    }

    config.headers['X-Device-Info'] = JSON.stringify({
      device_type: 'web',
      device_name: getDeviceName(),
      user_agent: navigator.userAgent,
      platform: 'oa'
    })

    // 添加请求时间戳
    config.headers['X-Request-Time'] = Date.now().toString()

    return config
  },
  (error) => {
    hideLoading()
    ElMessage.error('请求配置错误')
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    hideLoading()
    
    const { data } = response
    
    // 处理业务状态码
    if (data.code !== 200) {
      // 特殊错误码处理
      switch (data.code) {
        case 401:
          handleAuthError('登录已过期，请重新登录')
          break
        case 403:
          ElMessage.error('权限不足')
          break
        case 429:
          ElMessage.error('请求过于频繁，请稍后再试')
          break
        default:
          ElMessage.error(data.message || '请求失败')
      }
      return Promise.reject(new Error(data.message || '请求失败'))
    }

    // 返回完整的response对象，保持拦截器类型一致性
    return response
  },
  async (error) => {
    hideLoading()
    
    if (error.response) {
      const { status, data } = error.response
      
      switch (status) {
        case 401:
          // Token过期，尝试刷新
          const refreshed = await handleTokenRefresh()
          if (refreshed) {
            // 重试原请求
            return request.request(error.config)
          } else {
            handleAuthError('登录已过期，请重新登录')
          }
          break
        case 403:
          ElMessage.error('权限不足')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 429:
          ElMessage.error('请求过于频繁，请稍后再试')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        case 502:
        case 503:
        case 504:
          ElMessage.error('服务器暂时不可用，请稍后再试')
          break
        default:
          ElMessage.error(data?.message || `请求失败 (${status})`)
      }
    } else if (error.code === 'ECONNABORTED') {
      ElMessage.error('请求超时，请检查网络连接')
    } else if (error.message === 'Network Error') {
      ElMessage.error('网络连接失败，请检查网络设置')
    } else {
      ElMessage.error('网络错误，请稍后重试')
    }
    
    return Promise.reject(error)
  }
)

// 处理Token刷新
const handleTokenRefresh = async (): Promise<boolean> => {
  try {
    const authStore = useAuthStore()
    if (!authStore.refreshToken) {
      return false
    }

    // 调用OA专用刷新Token接口 - 使用JSON格式
    const response = await axios.post('/api/oa/auth/refresh', {
      refresh_token: authStore.refreshToken
    }, {
      headers: {
        'Content-Type': 'application/json'
      }
    })

    if (response.data.code === 200) {
      // 更新Token
      authStore.updateTokens(
        response.data.data.access_token,
        response.data.data.refresh_token
      )
      return true
    }
    return false
  } catch (error) {
    console.error('Token刷新失败:', error)
    return false
  }
}

// 处理认证错误
const handleAuthError = (message: string) => {
  ElMessage.error(message)
  const authStore = useAuthStore()
  authStore.logout()
  
  // 跳转到登录页，避免无限重定向
  if (router.currentRoute.value.path !== '/login') {
    router.push({
      path: '/login',
      query: { redirect: router.currentRoute.value.fullPath }
    })
  }
}

// 导出带有类型提示的请求方法
interface RequestMethods {
  get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T>
  post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T>
  put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T>
  delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<T>
  patch<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T>
}

// 包装请求方法，自动提取data
const requestMethods: RequestMethods = {
  get: async (url, config) => {
    const response = await request.get(url, config)
    return response.data
  },
  post: async (url, data, config) => {
    const response = await request.post(url, data, config)
    return response.data
  },
  put: async (url, data, config) => {
    const response = await request.put(url, data, config)
    return response.data
  },
  delete: async (url, config) => {
    const response = await request.delete(url, config)
    return response.data
  },
  patch: async (url, data, config) => {
    const response = await request.patch(url, data, config)
    return response.data
  }
}

export default requestMethods 