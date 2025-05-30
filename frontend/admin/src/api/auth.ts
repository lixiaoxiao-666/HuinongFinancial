import request from '@/utils/request'
import type { 
  LoginRequest, 
  LoginResponse, 
  RefreshTokenResponse,
  SessionDetail,
  ApiResponse 
} from '@/types/auth'

// OA系统登录 - 直接调用后端OA登录接口
export const oaLogin = async (data: LoginRequest): Promise<LoginResponse> => {
  const requestData = {
    username: data.username,
    password: data.password,
    platform: 'oa',
    device_type: data.device_type || 'web',
    device_name: data.device_name || 'OA管理系统',
    device_id: data.device_id || `oa_web_${Date.now()}`,
    app_version: data.app_version || '1.0.0'
  }
  
  // 直接调用OA登录接口
  const response = await request.post<any>('/api/oa/auth/login', requestData)
  
  // 转换响应格式以匹配前端期望
  return {
    code: response.code || 200,
    message: response.message || '登录成功',
    data: {
      admin: {
        id: response.data.user.id,
        username: response.data.user.username,
        email: response.data.user.email,
        real_name: response.data.user.real_name,
        role: response.data.user.username === 'admin' ? 'super_admin' : 'reviewer',
        permissions: response.data.user.username === 'admin' 
          ? ['user_manage', 'loan_approve', 'machine_manage', 'content_manage', 'system_config', 'dashboard:view']
          : ['loan_approve', 'dashboard:view'],
        status: response.data.user.status,
        department: response.data.user.department,
        position: response.data.user.position,
        last_login_at: response.data.user.last_login_at || new Date().toISOString(),
        created_at: response.data.user.created_at || new Date().toISOString(),
        updated_at: response.data.user.updated_at || new Date().toISOString()
      },
      session: {
        access_token: response.data.access_token,
        refresh_token: response.data.refresh_token,
        session_id: response.data.session_id,
        expires_in: response.data.expires_in
      }
    }
  }
}

// 刷新Token - 使用OA专用接口
export const refreshToken = (refreshToken: string): Promise<RefreshTokenResponse> => {
  return request.post('/api/oa/auth/refresh', {
    refresh_token: refreshToken
  }, {
    headers: {
      'Content-Type': 'application/json'
    }
  })
}

// 用户登出 - 使用OA专用接口
export const logout = (): Promise<ApiResponse> => {
  return request.post('/api/oa/auth/logout')
}

// 获取会话信息
export const getSessionInfo = (): Promise<ApiResponse<SessionDetail[]>> => {
  return request.get('/api/user/session/info')
}

// 注销其他设备会话
export const revokeOtherSessions = (): Promise<ApiResponse> => {
  return request.post('/api/user/session/revoke-others')
}

// 获取当前管理员信息（通过session验证）
export const getCurrentAdmin = (): Promise<ApiResponse> => {
  return request.get('/api/oa/auth/me')
}

// 验证Token有效性 - 使用OA专用接口
export const validateToken = (): Promise<ApiResponse> => {
  return request.get('/api/oa/auth/validate')
} 