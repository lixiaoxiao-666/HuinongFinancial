// 基于后端 session_management.md API 设计的认证类型定义

export interface LoginRequest {
  username: string
  password: string
  platform?: 'oa'
  device_id?: string
  device_type?: string
  device_name?: string
  app_version?: string
  ip_address?: string
  location?: string
}

export interface AdminUser {
  id: number
  username: string
  email: string
  real_name: string
  role: Role
  department: string
  position: string
  permissions: Permission[]
  status: string
  last_login_at: string
  created_at: string
  updated_at: string
}

export interface SessionInfo {
  session_id: string
  access_token: string
  refresh_token: string
  expires_in: number
  platform: string
}

export interface LoginResponse extends ApiResponse {
  data: {
    admin: AdminUser
    session: {
      access_token: string
      refresh_token: string
      session_id: string
      expires_in: number
    }
  }
}

export interface SessionDetail {
  session_id: string
  platform: string
  device_info: {
    device_id: string
    device_type: string
    device_name: string
    app_version: string
  }
  network_info: {
    ip_address: string
    location: string
  }
  status: string
  created_at: string
  last_active_at: string
}

export interface RefreshTokenResponse {
  code: number
  message: string
  data: {
    access_token: string
    refresh_token: string
    expires_in: number
  }
}

// 权限相关类型
export type Permission = 
  | 'user_manage' 
  | 'loan_approve' 
  | 'machine_manage'
  | 'content_manage'
  | 'system_config'
  | 'risk_monitor'
  | 'data_export'
  | 'log_view'
  | 'dashboard:view'

export type Role = 'super_admin' | 'admin' | 'reviewer' | 'operator'

// API响应基础类型
export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

// 错误类型
export interface AuthError {
  code: number
  message: string
  details?: string
} 