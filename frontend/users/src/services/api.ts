import { useUserStore } from '@/stores/user'
import type { UserInfo, LoginResponse } from '@/stores/user'

// API基础配置
const API_BASE_URL = 'http://172.18.120.10:8080/api/v1'

// 统一响应格式
interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 分页响应格式
interface PaginatedResponse<T = any> extends ApiResponse<T[]> {
  total: number
}

// 请求拦截器
class ApiClient {
  private baseURL: string

  constructor(baseURL: string) {
    this.baseURL = baseURL
  }

  // 通用请求方法
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const userStore = useUserStore()
    const url = `${this.baseURL}${endpoint}`
    
    // 设置默认headers
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...((options.headers as Record<string, string>) || {})
    }

    // 如果有token且需要认证，添加Authorization header
    if (userStore.token && userStore.isTokenValid()) {
      headers.Authorization = `Bearer ${userStore.token}`
    }

    const config: RequestInit = {
      ...options,
      headers
    }

    try {
      const response = await fetch(url, config)
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const result = await response.json()
      
      // 检查API响应码
      if (result.code !== 0) {
        throw new Error(result.message || '请求失败')
      }

      return result
    } catch (error) {
      console.error('API请求失败:', error)
      throw error
    }
  }

  // GET请求
  get<T>(endpoint: string, params?: Record<string, any>): Promise<T> {
    const url = params ? `${endpoint}?${new URLSearchParams(params)}` : endpoint
    return this.request<T>(url, { method: 'GET' })
  }

  // POST请求
  post<T>(endpoint: string, data?: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined
    })
  }

  // PUT请求
  put<T>(endpoint: string, data?: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined
    })
  }

  // DELETE请求
  delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' })
  }

  // 文件上传
  async uploadFile(file: File, purpose?: string): Promise<ApiResponse<FileUploadResult>> {
    const userStore = useUserStore()
    const url = `${this.baseURL}/files/upload`
    
    const formData = new FormData()
    formData.append('file', file)
    if (purpose) {
      formData.append('purpose', purpose)
    }

    const headers: Record<string, string> = {}
    if (userStore.token && userStore.isTokenValid()) {
      headers.Authorization = `Bearer ${userStore.token}`
    }

    const response = await fetch(url, {
      method: 'POST',
      headers,
      body: formData
    })

    if (!response.ok) {
      throw new Error(`上传失败: ${response.status}`)
    }

    return response.json()
  }
}

// 创建API客户端实例
const apiClient = new ApiClient(API_BASE_URL)

// 用户相关接口
export interface RegisterRequest {
  phone: string
  password: string
}

export interface LoginRequest {
  phone: string
  password: string
}

export interface UpdateUserRequest {
  nickname?: string
  avatar_url?: string
  real_name?: string
  address?: string
}

// 贷款相关接口
export interface LoanProduct {
  product_id: string
  name: string
  description: string
  category: string
  min_amount: number
  max_amount: number
  min_term_months: number
  max_term_months: number
  interest_rate_yearly: string
  status: number
  repayment_methods?: string[]
  application_conditions?: string
  required_documents?: Array<{
    type: string
    desc: string
  }>
}

export interface LoanApplication {
  application_id: string
  product_id: string
  user_id: string
  amount: number
  term_months: number
  purpose: string
  status: string
  submitted_at: string
  updated_at: string
  approved_amount?: number
  remarks?: string
  history?: Array<{
    status: string
    timestamp: string
    operator: string
  }>
}

export interface LoanApplicationRequest {
  product_id: string
  amount: number
  term_months: number
  purpose: string
  applicant_info: {
    real_name: string
    id_card_number: string
    address: string
  }
  uploaded_documents: Array<{
    doc_type: string
    file_id: string
  }>
}

export interface FileUploadResult {
  file_id: string
  file_url: string
  file_name: string
  file_size: number
}

// 用户服务API
export const userApi = {
  // 用户注册
  register(data: RegisterRequest): Promise<ApiResponse<{ user_id: string }>> {
    return apiClient.post('/users/register', data)
  },

  // 用户登录
  login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    return apiClient.post('/users/login', data)
  },

  // 获取用户信息
  getUserInfo(): Promise<ApiResponse<UserInfo>> {
    return apiClient.get('/users/me')
  },

  // 更新用户信息
  updateUserInfo(data: UpdateUserRequest): Promise<ApiResponse> {
    return apiClient.put('/users/me', data)
  }
}

// 贷款服务API
export const loanApi = {
  // 获取贷款产品列表
  getProducts(category?: string): Promise<ApiResponse<LoanProduct[]>> {
    const params = category ? { category } : undefined
    return apiClient.get('/loans/products', params)
  },

  // 获取贷款产品详情
  getProductDetail(productId: string): Promise<ApiResponse<LoanProduct>> {
    return apiClient.get(`/loans/products/${productId}`)
  },

  // 提交贷款申请
  submitApplication(data: LoanApplicationRequest): Promise<ApiResponse<{ application_id: string }>> {
    return apiClient.post('/loans/applications', data)
  },

  // 获取贷款申请详情
  getApplicationDetail(applicationId: string): Promise<ApiResponse<LoanApplication>> {
    return apiClient.get(`/loans/applications/${applicationId}`)
  },

  // 获取我的贷款申请列表
  getMyApplications(params?: {
    status?: string
    page?: number
    limit?: number
  }): Promise<PaginatedResponse<LoanApplication>> {
    return apiClient.get('/loans/applications/my', params)
  }
}

// 文件服务API
export const fileApi = {
  // 文件上传
  upload(file: File, purpose?: string): Promise<ApiResponse<FileUploadResult>> {
    return apiClient.uploadFile(file, purpose)
  }
}

// 健康检查API
export const healthApi = {
  check(): Promise<ApiResponse<{ status: string; service: string; version: string }>> {
    return apiClient.get('/health')
  }
}

export default apiClient 