import request from '../request'
import { mockApi, shouldUseMock } from '../mock'
import type { 
  ApiResponse, 
  LoginCredentials, 
  LoginResponse, 
  OAUser,
  SessionInfo 
} from '../types'

/**
 * 生成设备ID
 */
const generateDeviceId = (): string => {
  const timestamp = Date.now()
  const random = Math.random().toString(36).substring(2, 8)
  return `OA_WebApp_${timestamp}_${random}`
}

/**
 * 认证相关API
 */
export const authApi = {
  /**
   * OA用户登录
   * API文档: POST /api/oa/auth/login
   */
  login(credentials: LoginCredentials): Promise<ApiResponse<LoginResponse>> {
    // 开发环境优先使用Mock API
    if (shouldUseMock()) {
      console.log('🚀 使用Mock API进行登录测试')
      return mockApi.login(credentials)
    }

    return request.post('/auth/login', {
      username: credentials.username,
      password: credentials.password,
      platform: 'oa',
      device_info: {
        device_id: generateDeviceId(),
        device_type: 'web',
        user_agent: navigator.userAgent,
        ...credentials.device_info
      }
    })
  },

  /**
   * 刷新Access Token  
   * API文档: POST /api/oa/auth/refresh
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
   * API文档: GET /api/oa/auth/validate
   */
  validate(): Promise<ApiResponse<{
    valid: boolean
    user_id: number
    platform: string
    role: string
  }>> {
    if (shouldUseMock()) {
      return mockApi.validate()
    }
    return request.get('/auth/validate')
  },

  /**
   * 获取当前用户信息
   * API文档: GET /api/oa/user/profile
   */
  getCurrentUser(): Promise<ApiResponse<OAUser>> {
    if (shouldUseMock()) {
      return mockApi.getCurrentUser()
    }
    return request.get('/user/profile')
  },

  /**
   * 用户登出
   * API文档: POST /api/oa/auth/logout
   */
  logout(): Promise<ApiResponse<null>> {
    if (shouldUseMock()) {
      return mockApi.logout()
    }
    return request.post('/auth/logout')
  },

  /**
   * 修改密码
   * API文档: PUT /api/oa/user/password
   */
  changePassword(data: {
    old_password: string
    new_password: string
  }): Promise<ApiResponse<null>> {
    return request.put('/user/password', data)
  }
}

/**
 * 会话管理相关API
 */
export const sessionApi = {
  /**
   * 获取当前用户的活跃会话列表
   * API文档: GET /api/user/session/info
   */
  getCurrentUserSessions(): Promise<ApiResponse<SessionInfo[]>> {
    return request.get('/user/session/info')
  },

  /**
   * 注销指定会话
   * API文档: DELETE /api/user/session/{id}
   */
  revokeSession(session_id: string): Promise<ApiResponse<null>> {
    return request.delete(`/user/session/${session_id}`)
  },

  /**
   * 注销除当前会话外的其他所有会话
   * API文档: POST /api/user/session/revoke-others
   */
  revokeOtherSessions(): Promise<ApiResponse<{ revoked_count: number }>> {
    return request.post('/user/session/revoke-others')
  }
} 