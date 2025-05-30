import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { 
  oaLogin, 
  logout as apiLogout, 
  refreshToken as apiRefreshToken,
  getSessionInfo,
  revokeOtherSessions,
  getCurrentAdmin,
  validateToken
} from '@/api/auth'
import type { 
  AdminUser, 
  LoginRequest, 
  SessionInfo, 
  SessionDetail,
  Permission,
  Role 
} from '@/types/auth'
import { ElMessage } from 'element-plus'

// 存储键名常量
const STORAGE_KEYS = {
  ACCESS_TOKEN: 'oa_access_token',
  REFRESH_TOKEN: 'oa_refresh_token',
  SESSION_ID: 'oa_session_id',
  USER_INFO: 'oa_user_info',
  REMEMBER_USERNAME: 'oa_remember_username'
} as const

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const accessToken = ref<string>(localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN) || '')
  const refreshToken = ref<string>(localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN) || '')
  const sessionId = ref<string>(localStorage.getItem(STORAGE_KEYS.SESSION_ID) || '')
  
  // 获取存储的用户信息
  const getUserFromStorage = (): AdminUser | null => {
    const stored = localStorage.getItem(STORAGE_KEYS.USER_INFO)
    return stored ? JSON.parse(stored) : null
  }
  
  const user = ref<AdminUser | null>(getUserFromStorage())
  const isLoggedIn = ref<boolean>(!!accessToken.value)
  const sessions = ref<SessionDetail[]>([])

  // 计算属性
  const isAuthenticated = computed(() => !!accessToken.value && !!user.value)
  const userRole = computed(() => user.value?.role as Role)
  const userPermissions = computed(() => user.value?.permissions || [])

  // 获取设备信息
  const getDeviceInfo = () => {
    return {
      device_type: 'web',
      device_name: navigator.userAgent,
      platform: 'oa' as const,
      device_id: `web_${Date.now()}`,
      app_version: '1.0.0',
      ip_address: '', // 后端会自动获取
      location: '后台管理系统'
    }
  }

  // 登录 - 简化版本，无需验证码
  const login = async (username: string, password: string) => {
    try {
      const loginData: LoginRequest = {
        username,
        password,
        ...getDeviceInfo()
      }

      const response = await oaLogin(loginData)
      
      if (response.code === 200 && response.data) {
        // 保存用户信息
        user.value = response.data.admin
        
        // 保存会话信息
        accessToken.value = response.data.session.access_token
        refreshToken.value = response.data.session.refresh_token
        sessionId.value = response.data.session.session_id
        isLoggedIn.value = true

        // 持久化存储
        saveToStorage()

        ElMessage.success('登录成功')
        return response.data
      } else {
        throw new Error(response.message || '登录失败')
      }
    } catch (error: any) {
      console.error('登录失败:', error)
      throw new Error(error.message || '登录失败，请检查用户名和密码')
    }
  }

  // 登出
  const logout = async (showMessage = true) => {
    try {
      // 调用后端登出接口
      if (accessToken.value) {
        await apiLogout()
      }
      
      if (showMessage) {
        ElMessage.success('退出登录成功')
      }
    } catch (error) {
      console.error('登出请求失败:', error)
      // 即使后端请求失败，也要清除本地状态
    } finally {
      // 清除所有状态
      clearAuthState()
    }
  }

  // 刷新Token
  const updateTokens = (newAccessToken: string, newRefreshToken: string) => {
    accessToken.value = newAccessToken
    refreshToken.value = newRefreshToken
    saveToStorage()
  }

  // 验证Token有效性
  const verifyAuth = async (): Promise<boolean> => {
    if (!accessToken.value) {
      return false
    }

    try {
      await validateToken()
      return true
    } catch (error) {
      console.error('Token验证失败:', error)
      clearAuthState()
      return false
    }
  }

  // 获取会话信息
  const fetchSessionInfo = async () => {
    try {
      const response = await getSessionInfo()
      if (response.code === 200 && response.data) {
        sessions.value = response.data
      }
      return response.data
    } catch (error) {
      console.error('获取会话信息失败:', error)
      throw error
    }
  }

  // 注销其他设备
  const logoutOtherDevices = async () => {
    try {
      await revokeOtherSessions()
      await fetchSessionInfo() // 刷新会话信息
      ElMessage.success('已注销其他设备')
    } catch (error) {
      console.error('注销其他设备失败:', error)
      throw error
    }
  }

  // 权限检查
  const hasPermission = (permission: Permission): boolean => {
    if (!user.value) return false
    
    // 超级管理员拥有所有权限
    if (user.value.role === 'super_admin') return true
    
    // 检查用户权限列表
    return user.value.permissions.includes(permission)
  }

  // 角色检查
  const hasRole = (role: Role): boolean => {
    if (!user.value) return false
    return user.value.role === role
  }

  // 检查是否有任一权限
  const hasAnyPermission = (permissions: Permission[]): boolean => {
    return permissions.some(permission => hasPermission(permission))
  }

  // 保存到本地存储
  const saveToStorage = () => {
    localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, accessToken.value)
    localStorage.setItem(STORAGE_KEYS.REFRESH_TOKEN, refreshToken.value)
    localStorage.setItem(STORAGE_KEYS.SESSION_ID, sessionId.value)
    if (user.value) {
      localStorage.setItem(STORAGE_KEYS.USER_INFO, JSON.stringify(user.value))
    }
  }

  // 清除认证状态
  const clearAuthState = () => {
    accessToken.value = ''
    refreshToken.value = ''
    sessionId.value = ''
    user.value = null
    isLoggedIn.value = false
    sessions.value = []

    // 清除本地存储
    Object.values(STORAGE_KEYS).forEach(key => {
      localStorage.removeItem(key)
    })
  }

  // 记住用户名
  const rememberUsername = (username: string) => {
    localStorage.setItem(STORAGE_KEYS.REMEMBER_USERNAME, username)
  }

  // 获取记住的用户名
  const getRememberedUsername = (): string => {
    return localStorage.getItem(STORAGE_KEYS.REMEMBER_USERNAME) || ''
  }

  // 清除记住的用户名
  const clearRememberedUsername = () => {
    localStorage.removeItem(STORAGE_KEYS.REMEMBER_USERNAME)
  }

  // 初始化认证状态
  const initAuth = async () => {
    if (accessToken.value && user.value) {
      const isValid = await verifyAuth()
      if (!isValid) {
        clearAuthState()
      }
    }
  }

  return {
    // 状态
    accessToken,
    refreshToken,
    sessionId,
    user,
    isLoggedIn,
    sessions,
    
    // 计算属性
    isAuthenticated,
    userRole,
    userPermissions,
    
    // 方法
    login,
    logout,
    updateTokens,
    verifyAuth,
    fetchSessionInfo,
    logoutOtherDevices,
    hasPermission,
    hasRole,
    hasAnyPermission,
    rememberUsername,
    getRememberedUsername,
    clearRememberedUsername,
    initAuth,
    clearAuthState
  }
}) 