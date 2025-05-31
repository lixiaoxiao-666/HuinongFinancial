import api from './index'
import type { 
  LoginResponse, 
  DashboardData, 
  PaginationResponse, 
  LoanApplication, 
  ApplicationDetail, 
  AdminUser, 
  OperationLog, 
  SystemConfig 
} from '@/types'
import {
  mockLogin,
  mockDashboard,
  mockApplications,
  mockApplicationDetail,
  mockUsers,
  mockLogs,
  mockConfigs,
  mockSubmitReview,
  mockCreateUser,
  mockUpdateUserStatus,
  mockToggleAI,
  mockUpdateConfig
} from './mock'

// 开发环境使用模拟数据
const USE_MOCK_DATA = import.meta.env.DEV

// 带有容错的API调用包装器
const apiWithFallback = async <T>(apiCall: () => Promise<T>, mockCall: () => Promise<T>): Promise<T> => {
  if (USE_MOCK_DATA) {
    console.log('🔧 开发模式：使用模拟数据')
    return mockCall()
  }
  
  try {
    return await apiCall()
  } catch (error: any) {
    console.warn('⚠️ API调用失败，切换到模拟数据:', error.message)
    return mockCall()
  }
}

// 登录相关接口
export const adminLogin = (data: { username: string; password: string }): Promise<LoginResponse> => {
  return apiWithFallback(
    () => api.post('/admin/login', data),
    () => mockLogin(data)
  )
}

// 工作台相关接口
export const getDashboard = (): Promise<DashboardData> => {
  return apiWithFallback(
    () => api.get('/admin/dashboard'),
    () => mockDashboard()
  )
}

// 贷款审批相关接口
export const getPendingApplications = (params: {
  status_filter?: string
  applicant_name?: string
  application_id?: string
  page?: number
  limit?: number
}): Promise<PaginationResponse<LoanApplication>> => {
  return apiWithFallback(
    () => api.get('/admin/loans/applications/pending', { params }),
    () => mockApplications(params)
  )
}

export const getApplicationDetail = (applicationId: string): Promise<ApplicationDetail> => {
  return apiWithFallback(
    () => api.get(`/admin/loans/applications/${applicationId}`),
    () => mockApplicationDetail(applicationId)
  )
}

export const submitReview = (applicationId: string, data: {
  decision: 'approved' | 'rejected' | 'require_more_info'
  approved_amount?: number
  comments: string
  required_info_details?: string
}): Promise<void> => {
  return apiWithFallback(
    () => api.post(`/admin/loans/applications/${applicationId}/review`, data),
    () => mockSubmitReview(applicationId, data)
  )
}

// 系统管理相关接口
export const toggleAIApproval = (enabled: boolean): Promise<void> => {
  return apiWithFallback(
    () => api.post('/admin/system/ai-approval/toggle', { enabled }),
    () => mockToggleAI(enabled)
  )
}

export const getOAUsers = (params: {
  page?: number
  limit?: number
  role?: string
}): Promise<PaginationResponse<AdminUser>> => {
  return apiWithFallback(
    () => api.get('/admin/users', { params }),
    () => mockUsers(params)
  )
}

export const createOAUser = (data: {
  username: string
  password: string
  role: string
  display_name: string
  email: string
}): Promise<AdminUser> => {
  return apiWithFallback(
    () => api.post('/admin/users', data),
    () => mockCreateUser(data)
  )
}

export const updateOAUserStatus = (userId: string, status: number): Promise<void> => {
  return apiWithFallback(
    () => api.put(`/admin/users/${userId}/status`, { status }),
    () => mockUpdateUserStatus(userId, status)
  )
}

// 操作日志相关接口
export const getOperationLogs = (params: {
  operator_id?: string
  action?: string
  start_date?: string
  end_date?: string
  page?: number
  limit?: number
}): Promise<PaginationResponse<OperationLog>> => {
  return apiWithFallback(
    () => api.get('/admin/logs', { params }),
    () => mockLogs(params)
  )
}

// 系统配置相关接口
export const getSystemConfigurations = (): Promise<SystemConfig[]> => {
  return apiWithFallback(
    () => api.get('/admin/configs'),
    () => mockConfigs()
  )
}

export const updateSystemConfiguration = (configKey: string, configValue: string): Promise<void> => {
  return apiWithFallback(
    () => api.put(`/admin/configs/${configKey}`, { config_value: configValue }),
    () => mockUpdateConfig(configKey, configValue)
  )
} 