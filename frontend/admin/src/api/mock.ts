import type { ApiResponse, LoginResponse, OAUser } from './types'

/**
 * Mock 数据 - 仅用于开发环境
 */

// Mock OA用户数据
const mockOAUser: OAUser = {
  id: 1,
  username: 'admin',
  real_name: '系统管理员',
  email: 'admin@huinong.com',
  phone: '13800138000',
  avatar: 'https://avatars.githubusercontent.com/u/1?v=4',
  role_id: 1,
  role_name: '系统管理员',
  department: '技术部',
  position: '系统管理员',
  status: 'active',
  permissions: ['*'], // 所有权限
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-15T10:00:00Z',
  last_login_at: new Date().toISOString()
}

/**
 * Mock API 服务
 */
export const mockApi = {
  /**
   * 模拟登录
   */
  login(credentials: { username: string; password: string }): Promise<ApiResponse<LoginResponse>> {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        // 模拟登录验证
        if (credentials.username === 'admin' && credentials.password === 'admin123') {
          resolve({
            code: 200,
            message: '登录成功',
            data: {
              user: mockOAUser,
              session: {
                access_token: 'mock_access_token_' + Date.now(),
                refresh_token: 'mock_refresh_token_' + Date.now(),
                expires_in: 86400 // 24小时
              }
            }
          })
        } else {
          reject({
            response: {
              data: {
                code: 1004,
                message: '用户名或密码错误'
              }
            }
          })
        }
      }, 1000) // 模拟网络延迟
    })
  },

  /**
   * 模拟Token验证
   */
  validate(): Promise<ApiResponse<{
    valid: boolean
    user_id: number
    platform: string
    role: string
  }>> {
    return Promise.resolve({
      code: 200,
      message: '验证成功',
      data: {
        valid: true,
        user_id: mockOAUser.id,
        platform: 'oa',
        role: 'admin'
      }
    })
  },

  /**
   * 模拟获取用户信息
   */
  getCurrentUser(): Promise<ApiResponse<OAUser>> {
    return Promise.resolve({
      code: 200,
      message: '获取成功',
      data: mockOAUser
    })
  },

  /**
   * 模拟登出
   */
  logout(): Promise<ApiResponse<null>> {
    return Promise.resolve({
      code: 200,
      message: '登出成功',
      data: null
    })
  }
}

/**
 * 检查是否应该使用Mock API
 */
export const shouldUseMock = (): boolean => {
  // 强制在开发环境使用Mock API
  const forceMock = import.meta.env.MODE === 'development'
  const envMock = import.meta.env.VITE_USE_MOCK === 'true'
  
  console.log('🔍 Mock API 检查:', {
    mode: import.meta.env.MODE,
    envMock: import.meta.env.VITE_USE_MOCK,
    forceMock,
    envMock,
    shouldUse: forceMock || envMock
  })
  
  return forceMock || envMock
} 