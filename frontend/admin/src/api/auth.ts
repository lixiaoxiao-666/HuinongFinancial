import { request } from './request'

// OA用户信息接口
export interface OAUser {
  id: number
  username: string
  email: string
  phone: string
  real_name: string
  avatar?: string
  role_id: number
  role_name: string
  department?: string
  position?: string
  status: 'active' | 'frozen' | 'inactive'
  last_login_at?: string
  created_at: string
  updated_at: string
}

// 登录凭据接口
export interface LoginCredentials {
  username: string
  password: string
  platform: 'oa'
  device_info?: DeviceInfo
}

// 设备信息接口
export interface DeviceInfo {
  device_id: string
  device_type: 'web' | 'ios' | 'android'
  device_name?: string
  user_agent?: string
  app_version?: string
}

// 登录响应接口
export interface LoginResponse {
  user: OAUser
  access_token: string
  refresh_token: string
  expires_in: number
}

// Token刷新响应接口
export interface RefreshResponse {
  access_token: string
  refresh_token: string
  expires_in: number
}

// Token验证响应接口
export interface ValidateResponse {
  valid: boolean
  user_id: number
  platform: string
  role: string
  expires_at: string
}

/**
 * OA用户登录
 */
export const login = (credentials: LoginCredentials): Promise<LoginResponse> => {
  // 添加设备信息
  const deviceInfo: DeviceInfo = {
    device_id: `OA_Web_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
    device_type: 'web',
    device_name: 'OA管理后台',
    user_agent: navigator.userAgent,
    app_version: import.meta.env.VITE_APP_VERSION || '1.0.0'
  }

  return request<LoginResponse>({
    url: '/api/oa/auth/login',
    method: 'POST',
    data: {
      ...credentials,
      device_info: deviceInfo
    }
  })
}

/**
 * 刷新访问令牌
 */
export const refreshToken = (refresh_token: string): Promise<RefreshResponse> => {
  return request<RefreshResponse>({
    url: '/api/oa/auth/refresh',
    method: 'POST',
    data: { refresh_token }
  })
}

/**
 * 验证当前Token
 */
export const validateToken = (): Promise<ValidateResponse> => {
  return request<ValidateResponse>({
    url: '/api/oa/auth/validate',
    method: 'GET'
  })
}

/**
 * 退出登录
 */
export const logout = (): Promise<void> => {
  return request<void>({
    url: '/api/oa/auth/logout',
    method: 'POST'
  })
}

/**
 * 获取当前用户信息
 */
export const getCurrentUser = (): Promise<OAUser> => {
  return request<OAUser>({
    url: '/api/oa/user/profile',
    method: 'GET'
  })
}

/**
 * 更新用户信息
 */
export const updateProfile = (data: Partial<Pick<OAUser, 'email' | 'phone' | 'avatar'>>): Promise<OAUser> => {
  return request<OAUser>({
    url: '/api/oa/user/profile',
    method: 'PUT',
    data
  })
}

/**
 * 修改密码
 */
export const changePassword = (data: {
  old_password: string
  new_password: string
}): Promise<void> => {
  return request<void>({
    url: '/api/oa/user/password',
    method: 'PUT',
    data
  })
} 