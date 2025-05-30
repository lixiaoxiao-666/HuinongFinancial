import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api'
import type { OAUser, LoginCredentials } from '@/api/types'
import { message } from 'ant-design-vue'

export const useAuthStore = defineStore('auth', () => {
  // 状态定义
  const token = ref<string | null>(localStorage.getItem('oa_token'))
  const refreshToken = ref<string | null>(localStorage.getItem('oa_refresh_token'))
  const userInfo = ref<OAUser | null>(null)
  const isLoggedIn = ref(false)
  const permissions = ref<string[]>([])
  const sessionId = ref<string | null>(null)

  // 计算属性
  const isAuthenticated = computed(() => !!token.value && isLoggedIn.value)
  const userRole = computed(() => userInfo.value?.role)
  const userName = computed(() => userInfo.value?.real_name || userInfo.value?.username)
  const userDepartment = computed(() => userInfo.value?.department)

  // 权限检查
  const hasPermission = (permission: string): boolean => {
    return permissions.value.includes(permission) || userRole.value === 'admin'
  }

  const hasRole = (role: string | string[]): boolean => {
    if (!userRole.value) return false
    if (Array.isArray(role)) {
      return role.includes(userRole.value)
    }
    return userRole.value === role
  }

  // 登录方法
  const login = async (credentials: LoginCredentials) => {
    try {
      const response = await authApi.login({
        ...credentials,
        device_info: {
          device_type: 'web',
          device_name: navigator.userAgent,
          app_version: import.meta.env.VITE_APP_VERSION || '1.0.0'
        }
      })

      const { user, session } = response.data

      // 保存认证信息
      token.value = session.access_token
      refreshToken.value = session.refresh_token
      userInfo.value = user
      isLoggedIn.value = true
      permissions.value = user.permissions || []

      // 持久化存储
      localStorage.setItem('oa_token', session.access_token)
      localStorage.setItem('oa_refresh_token', session.refresh_token)
      localStorage.setItem('oa_user_info', JSON.stringify(user))

      message.success(`欢迎回来，${user.real_name}!`)
      return response
    } catch (error) {
      console.error('登录失败:', error)
      throw error
    }
  }

  // 登出方法
  const logout = async (showMessage = true) => {
    try {
      // 调用后端登出接口
      if (token.value) {
        await authApi.logout()
      }
    } catch (error) {
      console.warn('后端登出失败:', error)
    } finally {
      // 清除本地状态
      token.value = null
      refreshToken.value = null
      userInfo.value = null
      isLoggedIn.value = false
      permissions.value = []
      sessionId.value = null

      // 清除本地存储
      localStorage.removeItem('oa_token')
      localStorage.removeItem('oa_refresh_token')
      localStorage.removeItem('oa_user_info')

      if (showMessage) {
        message.success('已安全登出')
      }
    }
  }

  // 刷新Token
  const refreshAccessToken = async () => {
    if (!refreshToken.value) {
      throw new Error('没有refresh token')
    }

    try {
      const response = await authApi.refresh(refreshToken.value)
      const { access_token, refresh_token, expires_in } = response.data

      token.value = access_token
      refreshToken.value = refresh_token

      localStorage.setItem('oa_token', access_token)
      localStorage.setItem('oa_refresh_token', refresh_token)

      return response
    } catch (error) {
      console.error('刷新Token失败:', error)
      await logout(false)
      throw error
    }
  }

  // 验证Token有效性
  const validateToken = async () => {
    if (!token.value) {
      return false
    }

    try {
      const response = await authApi.validate()
      const data = response.data

      if (data.valid) {
        sessionId.value = data.session_id
        return true
      } else {
        await logout(false)
        return false
      }
    } catch (error) {
      console.error('Token验证失败:', error)
      await logout(false)
      return false
    }
  }

  // 获取当前用户信息
  const fetchUserInfo = async () => {
    if (!token.value) {
      throw new Error('未登录')
    }

    try {
      const response = await authApi.getCurrentUser()
      userInfo.value = response.data
      permissions.value = response.data.permissions || []
      
      // 更新本地存储
      localStorage.setItem('oa_user_info', JSON.stringify(response.data))
      
      return response.data
    } catch (error) {
      console.error('获取用户信息失败:', error)
      throw error
    }
  }

  // 初始化认证状态（从本地存储恢复）
  const initializeAuth = async () => {
    const savedToken = localStorage.getItem('oa_token')
    const savedUserInfo = localStorage.getItem('oa_user_info')

    if (savedToken && savedUserInfo) {
      try {
        token.value = savedToken
        userInfo.value = JSON.parse(savedUserInfo)
        permissions.value = userInfo.value?.permissions || []

        // 验证Token是否仍然有效
        const isValid = await validateToken()
        if (isValid) {
          isLoggedIn.value = true
          console.log('认证状态恢复成功')
        }
      } catch (error) {
        console.error('认证状态恢复失败:', error)
        await logout(false)
      }
    }
  }

  // 修改密码
  const changePassword = async (oldPassword: string, newPassword: string) => {
    try {
      await authApi.changePassword({
        old_password: oldPassword,
        new_password: newPassword
      })
      message.success('密码修改成功')
    } catch (error) {
      console.error('密码修改失败:', error)
      throw error
    }
  }

  return {
    // 状态
    token,
    refreshToken,
    userInfo,
    isLoggedIn,
    permissions,
    sessionId,

    // 计算属性
    isAuthenticated,
    userRole,
    userName,
    userDepartment,

    // 方法
    login,
    logout,
    refreshAccessToken,
    validateToken,
    fetchUserInfo,
    initializeAuth,
    changePassword,
    hasPermission,
    hasRole
  }
}) 