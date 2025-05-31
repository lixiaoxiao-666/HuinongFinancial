import { useUserStore } from '@/stores/user'
import type { UserInfo, LoginResponse } from '@/stores/user'

// API基础配置 - 修正为正确的路径
const API_BASE_URL = 'http://172.18.120.10:8080/api'

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

// 设备信息类型
interface DeviceInfo {
  device_id: string
  device_type: 'ios' | 'android' | 'web'
  device_name?: string
  app_version?: string
  user_agent?: string
}

// 获取设备信息的辅助函数
function getDeviceInfo(): DeviceInfo {
  const userAgent = navigator.userAgent || ''
  let deviceType: 'ios' | 'android' | 'web' = 'web'
  let deviceName = 'Web浏览器'

  // 检测设备类型
  if (/iPhone|iPad|iPod/i.test(userAgent)) {
    deviceType = 'ios'
    deviceName = /iPad/i.test(userAgent) ? 'iPad' : 'iPhone'
  } else if (/Android/i.test(userAgent)) {
    deviceType = 'android'
    deviceName = 'Android设备'
  } else {
    // Web浏览器检测
    if (userAgent.includes('Chrome')) {
      deviceName = 'Chrome浏览器'
    } else if (userAgent.includes('Firefox')) {
      deviceName = 'Firefox浏览器'
    } else if (userAgent.includes('Safari')) {
      deviceName = 'Safari浏览器'
    } else if (userAgent.includes('Edge')) {
      deviceName = 'Edge浏览器'
    }
  }

  const deviceInfo: DeviceInfo = {
    device_id: generateDeviceId(),
    device_type: deviceType,
    device_name: deviceName,
    app_version: '1.0.0', // 可以从环境变量或配置文件获取
    user_agent: userAgent
  }

  console.log('生成的设备信息:', deviceInfo) // 添加调试日志
  
  return deviceInfo
}

// 生成设备ID的辅助函数
function generateDeviceId(): string {
  // 尝试从localStorage获取已存在的设备ID
  let deviceId = localStorage.getItem('device_id')
  if (!deviceId) {
    // 生成新的设备ID
    deviceId = 'web_' + Math.random().toString(36).substr(2, 9) + '_' + Date.now()
    localStorage.setItem('device_id', deviceId)
  }
  return deviceId
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

    console.log('发送请求:', { url, method: options.method || 'GET', headers, body: options.body }) // 添加请求日志

    try {
      const response = await fetch(url, config)
      
      console.log('响应状态:', response.status, response.statusText) // 添加响应状态日志
      
      // 尝试解析响应体
      let result: any
      try {
        const responseText = await response.text()
        console.log('响应内容:', responseText) // 添加响应内容日志
        
        if (responseText) {
          result = JSON.parse(responseText)
        } else {
          result = {}
        }
      } catch (parseError) {
        console.error('解析响应JSON失败:', parseError)
        throw new Error(`服务器响应格式错误: ${response.status}`)
      }
      
      if (!response.ok) {
        // 提供更详细的错误信息
        const errorMessage = result.message || result.error || `HTTP错误: ${response.status} ${response.statusText}`
        console.error('请求失败详情:', { status: response.status, result })
        throw new Error(errorMessage)
      }
      
      // 检查API响应码 - 修改为文档中的200状态码
      if (result.code !== undefined && result.code !== 200) {
        console.error('API业务错误:', result)
        throw new Error(result.message || '请求失败')
      }

      console.log('请求成功:', result) // 添加成功日志
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

// 用户相关接口类型定义 - 根据用户管理文档更新
export interface RegisterRequest {
  phone: string
  password: string
  user_type: string          // 用户类型: farmer, farm_owner, cooperative, enterprise
  real_name: string          // 真实姓名
  platform: string          // 平台: app, web
  device_info?: DeviceInfo   // 设备信息
}

export interface LoginRequest {
  phone: string
  password: string
  platform: string          // 平台: app, web
  device_info?: DeviceInfo   // 设备信息
}

export interface UpdateUserRequest {
  real_name?: string
  email?: string
  gender?: string
  birthday?: string
  avatar?: string
  province?: string
  city?: string
  county?: string
  address?: string
}

// 会话信息类型
export interface SessionInfo {
  session_id: string
  platform: string
  device_info: DeviceInfo
  ip_address: string
  location: string
  last_active_at: string
  is_current: boolean
}

// Token刷新响应类型
export interface TokenRefreshResponse {
  access_token: string
  refresh_token: string
  expires_in: number
}

// 登录响应类型 - 根据文档修正
export interface LoginApiResponse {
  user: {
    id: number
    uuid: string
    phone: string
    user_type: string
    real_name: string
    status: string
  }
  session: {
    access_token: string
    refresh_token: string
    expires_in: number
  }
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

// 用户服务API - 根据用户管理文档更新接口路径
export const userApi = {
  // 用户注册 - 使用正确的接口路径 /api/auth/register
  register(data: Partial<RegisterRequest>): Promise<ApiResponse<LoginApiResponse>> {
    // 合并默认参数，确保所有必填字段都有值
    const registerData: RegisterRequest = {
      phone: data.phone || '',
      password: data.password || '',
      user_type: data.user_type || 'farmer', // 默认为农户
      real_name: data.real_name || '',
      platform: 'web', // 当前为web端
      device_info: getDeviceInfo()
    }
    
    // 验证必填字段
    if (!registerData.phone) {
      throw new Error('手机号不能为空')
    }
    if (!registerData.password) {
      throw new Error('密码不能为空')
    }
    if (!registerData.real_name) {
      throw new Error('真实姓名不能为空')
    }
    
    console.log('注册请求数据:', registerData) // 添加调试日志
    
    return apiClient.post('/auth/register', registerData)
  },

  // 用户登录 - 使用正确的接口路径 /api/auth/login
  login(data: Partial<LoginRequest>): Promise<ApiResponse<LoginApiResponse>> {
    // 合并默认参数，确保所有必填字段都有值
    const loginData: LoginRequest = {
      phone: data.phone || '',
      password: data.password || '',
      platform: 'web', // 当前为web端
      device_info: getDeviceInfo()
    }
    
    // 验证必填字段
    if (!loginData.phone) {
      throw new Error('手机号不能为空')
    }
    if (!loginData.password) {
      throw new Error('密码不能为空')
    }
    
    console.log('登录请求数据:', loginData) // 添加调试日志
    
    return apiClient.post('/auth/login', loginData)
  },

  // Token刷新 - /api/auth/refresh
  refreshToken(refreshToken: string): Promise<ApiResponse<TokenRefreshResponse>> {
    return apiClient.post('/auth/refresh', { refresh_token: refreshToken })
  },

  // Token验证 - /api/auth/validate
  validateToken(): Promise<ApiResponse<{ valid: boolean; user_id: number; platform: string; session_id: string; expires_at: string }>> {
    return apiClient.get('/auth/validate')
  },

  // 获取用户信息 - /api/user/profile
  getUserInfo(): Promise<ApiResponse<UserInfo>> {
    return apiClient.get('/user/profile')
  },

  // 更新用户信息 - /api/user/profile
  updateUserInfo(data: UpdateUserRequest): Promise<ApiResponse> {
    return apiClient.put('/user/profile', data)
  },

  // 修改密码 - /api/user/password
  changePassword(oldPassword: string, newPassword: string): Promise<ApiResponse> {
    return apiClient.put('/user/password', {
      old_password: oldPassword,
      new_password: newPassword
    })
  },

  // 用户登出 - /api/user/logout
  logout(): Promise<ApiResponse> {
    return apiClient.post('/user/logout')
  },

  // 获取用户会话列表 - /api/user/session/list
  getUserSessions(): Promise<ApiResponse<SessionInfo[]>> {
    return apiClient.get('/user/session/list')
  },

  // 注销指定会话 - /api/user/session/revoke
  revokeSession(sessionId: string): Promise<ApiResponse> {
    return apiClient.post('/user/session/revoke', { session_id_to_revoke: sessionId })
  },

  // 注销其他所有会话 - /api/user/session/revoke-others
  revokeOtherSessions(): Promise<ApiResponse<{ revoked_count: number }>> {
    return apiClient.post('/user/session/revoke-others')
  },

  // 实名认证申请 - /api/user/auth/real-name
  submitRealNameAuth(data: {
    id_card_number: string
    real_name: string
    id_card_front_img_url: string
    id_card_back_img_url: string
    face_verify_img_url: string
  }): Promise<ApiResponse<{ auth_id: string; auth_status: string }>> {
    return apiClient.post('/user/auth/real-name', data)
  },

  // 银行卡认证申请 - /api/user/auth/bank-card
  submitBankCardAuth(data: {
    bank_card_number: string
    bank_name: string
    cardholder_name: string
    bank_reserved_phone?: string
  }): Promise<ApiResponse<{ auth_id: string; auth_status: string }>> {
    return apiClient.post('/user/auth/bank-card', data)
  },

  // 获取用户标签 - /api/user/tags
  getUserTags(): Promise<ApiResponse<{
    user_id: number
    tags: Array<{
      tag_key: string
      tag_value: string
      display_name: string
      category: string
      created_at: string
      is_system: boolean
      is_editable: boolean
    }>
    categories: Record<string, string[]>
  }>> {
    return apiClient.get('/user/tags')
  },

  // 添加用户标签 - /api/user/tags
  addUserTag(data: {
    tag_key: string
    tag_value: string
    category: string
    description?: string
  }): Promise<ApiResponse> {
    return apiClient.post('/user/tags', data)
  },

  // 更新用户标签 - /api/user/tags/{tag_key}
  updateUserTag(tagKey: string, data: {
    tag_value: string
    description?: string
  }): Promise<ApiResponse> {
    return apiClient.put(`/user/tags/${tagKey}`, data)
  },

  // 删除用户标签 - /api/user/tags/{tag_key}
  deleteUserTag(tagKey: string): Promise<ApiResponse> {
    return apiClient.delete(`/user/tags/${tagKey}`)
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

// 工具函数导出
export { getDeviceInfo, generateDeviceId }

export default apiClient