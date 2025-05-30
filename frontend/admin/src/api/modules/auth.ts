import request from '../request'
import type { 
  ApiResponse, 
  LoginCredentials, 
  LoginResponse, 
  OAUser,
  SessionInfo 
} from '../types'

/**
 * 认证相关API
 */
export const authApi = {
  /**
   * OA用户登录
   */
  login(credentials: LoginCredentials): Promise<ApiResponse<LoginResponse>> {
    return request.post('/auth/login', {
      ...credentials,
      platform: 'oa'
    })
  },

  /**
   * 刷新Access Token
   */
  refresh(refresh_token: string): Promise<ApiResponse<{
    access_token: string
    refresh_token: string
    expires_in: number
  }>> {
    return request.post('/auth/refresh', { refresh_token })
  },

  /**
   * 验证当前Token有效性
   */
  validate(): Promise<ApiResponse<{
    valid: boolean
    user_id: number
    session_id: string
    platform: string
    role: string
    expires_at: string
  }>> {
    return request.get('/auth/validate')
  },

  /**
   * 获取当前用户信息
   */
  getCurrentUser(): Promise<ApiResponse<OAUser>> {
    return request.get('/auth/me')
  },

  /**
   * 用户登出
   */
  logout(): Promise<ApiResponse<null>> {
    return request.post('/auth/logout')
  },

  /**
   * 修改密码
   */
  changePassword(data: {
    old_password: string
    new_password: string
  }): Promise<ApiResponse<null>> {
    return request.post('/auth/change-password', data)
  }
}

/**
 * 会话管理相关API
 */
export const sessionApi = {
  /**
   * 获取当前用户的活跃会话列表
   */
  getCurrentUserSessions(): Promise<ApiResponse<SessionInfo[]>> {
    return request.get('/user/session/list')
  },

  /**
   * 注销指定会话
   */
  revokeSession(session_id: string): Promise<ApiResponse<null>> {
    return request.post('/user/session/revoke', {
      session_id_to_revoke: session_id
    })
  },

  /**
   * 注销除当前会话外的其他所有会话
   */
  revokeOtherSessions(): Promise<ApiResponse<{ revoked_count: number }>> {
    return request.post('/user/session/revoke-others')
  }
} 