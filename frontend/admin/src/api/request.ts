import axios from 'axios'
import type { AxiosRequestConfig, AxiosResponse } from 'axios'
import { message } from 'ant-design-vue'
import { useAuthStore } from '@/stores/modules/auth'
import router from '@/router'

// 响应数据接口
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  meta?: {
    timestamp: string
    request_id: string
    pagination?: {
      page: number
      limit: number
      total: number
    }
  }
}

// 创建axios实例
const service = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    
    // 添加认证token
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    
    // 添加请求ID用于调试
    config.headers['X-Request-ID'] = generateRequestId()
    
    // 开发环境日志
    if (import.meta.env.DEV) {
      console.log('API Request:', config.method?.toUpperCase(), config.url, config.data || config.params)
    }
    
    return config
  },
  (error) => {
    console.error('Request Error:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const { data } = response
    
    // 开发环境日志
    if (import.meta.env.DEV) {
      console.log('API Response:', response.config.url, data)
    }
    
    // 检查业务状态码
    if (data.code === 200) {
      return response // 返回完整的response对象
    } else {
      // 处理业务错误
      handleBusinessError(data.code, data.message)
      return Promise.reject(new Error(data.message || '请求失败'))
    }
  },
  (error) => {
    console.error('Response Error:', error)
    
    // 处理HTTP错误
    if (error.response) {
      const { status, data } = error.response
      handleHttpError(status, data?.message || error.message)
    } else if (error.request) {
      message.error('网络连接失败，请检查网络')
    } else {
      message.error('请求配置错误')
    }
    
    return Promise.reject(error)
  }
)

// 处理业务错误
const handleBusinessError = (code: number, msg: string) => {
  switch (code) {
    case 401:
      message.error('登录已过期，请重新登录')
      handleTokenExpired()
      break
    case 403:
      message.error('权限不足')
      break
    case 404:
      message.error('请求的资源不存在')
      break
    case 422:
      message.error('请求参数验证失败')
      break
    case 500:
      message.error('服务器内部错误')
      break
    default:
      if (msg) {
        message.error(msg)
      }
  }
}

// 处理HTTP错误
const handleHttpError = (status: number, msg: string) => {
  switch (status) {
    case 401:
      handleTokenExpired()
      break
    case 403:
      message.error('访问被拒绝')
      break
    case 404:
      message.error('接口不存在')
      break
    case 500:
      message.error('服务器错误')
      break
    case 502:
      message.error('网关错误')
      break
    case 503:
      message.error('服务不可用')
      break
    default:
      message.error(msg || `请求失败 (${status})`)
  }
}

// 处理token过期
const handleTokenExpired = () => {
  const authStore = useAuthStore()
  authStore.logout()
  
  // 如果不在登录页，跳转到登录页
  if (router.currentRoute.value.path !== '/login') {
    router.push('/login')
  }
}

// 生成请求ID
const generateRequestId = () => {
  return `req_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
}

// 通用请求方法
export const request = <T = any>(config: AxiosRequestConfig): Promise<T> => {
  return service(config).then(response => response.data.data)
}

// GET请求
export const get = <T = any>(url: string, params?: any): Promise<T> => {
  return request({ method: 'GET', url, params })
}

// POST请求
export const post = <T = any>(url: string, data?: any): Promise<T> => {
  return request({ method: 'POST', url, data })
}

// PUT请求
export const put = <T = any>(url: string, data?: any): Promise<T> => {
  return request({ method: 'PUT', url, data })
}

// DELETE请求
export const del = <T = any>(url: string, params?: any): Promise<T> => {
  return request({ method: 'DELETE', url, params })
}

// 文件上传
export const upload = <T = any>(url: string, formData: FormData, onProgress?: (progress: number) => void): Promise<T> => {
  return service({
    method: 'POST',
    url,
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    onUploadProgress: (progressEvent) => {
      if (onProgress && progressEvent.total) {
        const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
        onProgress(progress)
      }
    }
  }).then(response => response.data.data)
}

// 下载文件
export const download = (url: string, filename?: string): Promise<void> => {
  return service({
    method: 'GET',
    url,
    responseType: 'blob'
  }).then((response) => {
    const blob = new Blob([response.data])
    const downloadUrl = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = filename || 'download'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(downloadUrl)
  })
}

export default service 