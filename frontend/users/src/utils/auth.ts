import { useUserStore } from '@/stores/user'
import { userApi } from '@/services/api'
import { ElMessage } from 'element-plus'
import type { Router } from 'vue-router'

// 需要登录的路由列表
const authRequiredRoutes = [
  '/loan/apply',
  '/loan/application',
  '/loan/my-applications',
  '/me',
  '/user'
]

// 检查路由是否需要登录
export const isAuthRequired = (path: string): boolean => {
  return authRequiredRoutes.some(route => path.startsWith(route))
}

/**
 * 认证工具类
 * 提供Token刷新、自动登录检查等功能
 */
export class AuthManager {
  private static instance: AuthManager
  private refreshPromise: Promise<boolean> | null = null
  private router: Router | null = null

  private constructor() {}

  public static getInstance(): AuthManager {
    if (!AuthManager.instance) {
      AuthManager.instance = new AuthManager()
    }
    return AuthManager.instance
  }

  /**
   * 设置路由实例
   */
  public setRouter(router: Router) {
    this.router = router
  }

  /**
   * 刷新访问Token
   */
  public async refreshToken(): Promise<boolean> {
    // 如果已经有刷新请求在进行中，返回同一个Promise
    if (this.refreshPromise) {
      return this.refreshPromise
    }

    const userStore = useUserStore()
    
    if (!userStore.refreshToken) {
      console.warn('没有可用的刷新Token')
      return false
    }

    this.refreshPromise = this._performRefresh()
    const result = await this.refreshPromise
    this.refreshPromise = null
    
    return result
  }

  /**
   * 执行Token刷新
   */
  private async _performRefresh(): Promise<boolean> {
    const userStore = useUserStore()
    
    try {
      const response = await userApi.refreshToken(userStore.refreshToken)
      
      if (response.code === 200) {
        // 更新Token
        userStore.setToken(
          response.data.access_token,
          response.data.refresh_token,
          response.data.expires_in
        )
        
        console.log('Token刷新成功')
        return true
      } else {
        throw new Error(response.message || 'Token刷新失败')
      }
    } catch (error) {
      console.error('Token刷新失败:', error)
      // 刷新失败，清除用户状态并跳转到登录页
      this.handleAuthFailure()
      return false
    }
  }

  /**
   * 处理认证失败
   */
  private handleAuthFailure() {
    const userStore = useUserStore()
    userStore.logout()
    
    // 显示提示消息
    ElMessage.warning('登录已过期，请重新登录')
    
    // 跳转到登录页
    if (this.router) {
      this.router.push('/login')
    }
  }

  /**
   * 检查并自动刷新Token
   */
  public async checkAndRefreshToken(): Promise<boolean> {
    const userStore = useUserStore()
    
    // 如果用户未登录，不需要刷新
    if (!userStore.isLoggedIn) {
      return false
    }

    // 如果Token有效且不需要刷新，直接返回true
    if (userStore.isTokenValid() && !userStore.shouldRefreshToken()) {
      return true
    }

    // 如果Token完全过期且没有刷新Token，处理认证失败
    if (!userStore.isTokenValid() && !userStore.refreshToken) {
      this.handleAuthFailure()
      return false
    }

    // 尝试刷新Token
    return await this.refreshToken()
  }

  /**
   * 验证当前Token
   */
  public async validateToken(): Promise<boolean> {
    const userStore = useUserStore()
    
    if (!userStore.token) {
      return false
    }

    try {
      const response = await userApi.validateToken()
      return response.code === 200 && response.data.valid
    } catch (error) {
      console.error('Token验证失败:', error)
      return false
    }
  }

  /**
   * 初始化认证状态
   * 从localStorage恢复状态并验证Token
   */
  public async initializeAuth(): Promise<boolean> {
    const userStore = useUserStore()
    
    // 从localStorage恢复状态
    userStore.restoreFromStorage()
    
    // 如果没有Token，返回false
    if (!userStore.token) {
      return false
    }

    // 检查并刷新Token
    return await this.checkAndRefreshToken()
  }

  /**
   * 安全登出
   */
  public async logout(): Promise<void> {
    const userStore = useUserStore()
    
    try {
      // 调用服务端登出接口
      await userApi.logout()
    } catch (error) {
      console.warn('服务端登出失败:', error)
    } finally {
      // 无论服务端登出是否成功，都清除本地状态
      userStore.logout()
      
      // 跳转到登录页
      if (this.router) {
        this.router.push('/login')
      }
    }
  }

  /**
   * 获取认证头
   */
  public getAuthHeader(): string | null {
    const userStore = useUserStore()
    return userStore.token ? `Bearer ${userStore.token}` : null
  }
}

/**
 * 获取认证管理器实例
 */
export const authManager = AuthManager.getInstance()

/**
 * 路由守卫：检查认证状态
 */
export async function authGuard(to: any, from: any, next: any) {
  const userStore = useUserStore()
  
  // 定义不需要认证的路由
  const publicRoutes = ['/login', '/register', '/forgot-password']
  
  // 如果是公开路由，直接通过
  if (publicRoutes.includes(to.path)) {
    next()
    return
  }

  // 检查认证状态
  const isAuthenticated = await authManager.checkAndRefreshToken()
  
  if (isAuthenticated) {
    next()
  } else {
    // 未认证，跳转到登录页
    next('/login')
  }
}

/**
 * 获取设备指纹信息
 */
export function getDeviceFingerprint(): string {
  const canvas = document.createElement('canvas')
  const ctx = canvas.getContext('2d')
  ctx!.textBaseline = 'top'
  ctx!.font = '14px Arial'
  ctx!.fillText('Device fingerprint', 2, 2)
  
  const fingerprint = [
    navigator.userAgent,
    navigator.language,
    screen.width + 'x' + screen.height,
    new Date().getTimezoneOffset(),
    canvas.toDataURL()
  ].join('|')
  
  // 简单hash
  let hash = 0
  for (let i = 0; i < fingerprint.length; i++) {
    const char = fingerprint.charCodeAt(i)
    hash = ((hash << 5) - hash) + char
    hash = hash & hash // 转换为32位整数
  }
  
  return Math.abs(hash).toString(36)
}

/**
 * 检查密码强度
 */
export function checkPasswordStrength(password: string): {
  score: number
  level: 'weak' | 'medium' | 'strong'
  feedback: string[]
} {
  const feedback: string[] = []
  let score = 0

  // 长度检查
  if (password.length >= 8) {
    score += 2
  } else if (password.length >= 6) {
    score += 1
  } else {
    feedback.push('密码长度至少6位')
  }

  // 包含数字
  if (/\d/.test(password)) {
    score += 1
  } else {
    feedback.push('包含数字')
  }

  // 包含小写字母
  if (/[a-z]/.test(password)) {
    score += 1
  } else {
    feedback.push('包含小写字母')
  }

  // 包含大写字母
  if (/[A-Z]/.test(password)) {
    score += 1
  } else {
    feedback.push('包含大写字母')
  }

  // 包含特殊字符
  if (/[^A-Za-z0-9]/.test(password)) {
    score += 1
  } else {
    feedback.push('包含特殊字符')
  }

  // 判断强度等级
  let level: 'weak' | 'medium' | 'strong'
  if (score >= 5) {
    level = 'strong'
  } else if (score >= 3) {
    level = 'medium'
  } else {
    level = 'weak'
  }

  return { score, level, feedback }
}

/**
 * 格式化手机号显示
 */
export function formatPhoneNumber(phone: string): string {
  if (!phone) return ''
  
  // 脱敏显示
  if (phone.length === 11) {
    return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
  }
  
  return phone
}

/**
 * 验证手机号格式
 */
export function validatePhoneNumber(phone: string): boolean {
  return /^1[3-9]\d{9}$/.test(phone)
}

/**
 * 验证身份证号格式
 */
export function validateIdCard(idCard: string): boolean {
  return /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/.test(idCard)
}

export default authManager 