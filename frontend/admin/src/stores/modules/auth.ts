import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { message } from 'ant-design-vue'
import router from '@/router'
import * as authApi from '@/api/auth'
import type { OAUser, LoginCredentials } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref<string>(localStorage.getItem('auth_token') || '')
  const refreshToken = ref<string>(localStorage.getItem('refresh_token') || '')
  const userInfo = ref<OAUser | null>(null)
  const isLoading = ref(false)

  // 计算属性
  const isLoggedIn = computed(() => !!token.value && !!userInfo.value)
  const userName = computed(() => userInfo.value?.real_name || userInfo.value?.username || '')
  const userRole = computed(() => userInfo.value?.role_name || '')
  const userAvatar = computed(() => userInfo.value?.avatar || '')

  // 检查是否有指定角色
  const hasRole = (roles: string | string[]) => {
    if (!userInfo.value) return false
    const userRoleName = userInfo.value.role_name
    const roleArray = Array.isArray(roles) ? roles : [roles]
    return roleArray.includes(userRoleName)
  }

  // 保存token到localStorage
  const saveTokens = (accessToken: string, refresh: string) => {
    token.value = accessToken
    refreshToken.value = refresh
    localStorage.setItem('auth_token', accessToken)
    localStorage.setItem('refresh_token', refresh)
  }

  // 清除token
  const clearTokens = () => {
    token.value = ''
    refreshToken.value = ''
    userInfo.value = null
    localStorage.removeItem('auth_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user_info')
  }

  // 保存用户信息
  const saveUserInfo = (user: OAUser) => {
    userInfo.value = user
    localStorage.setItem('user_info', JSON.stringify(user))
  }

  // 从localStorage恢复用户信息
  const restoreUserInfo = () => {
    const savedUserInfo = localStorage.getItem('user_info')
    if (savedUserInfo) {
      try {
        userInfo.value = JSON.parse(savedUserInfo)
      } catch (error) {
        console.error('解析用户信息失败:', error)
        localStorage.removeItem('user_info')
      }
    }
  }

  // 登录
  const login = async (credentials: Omit<LoginCredentials, 'platform'>) => {
    try {
      isLoading.value = true
      
      const loginData: LoginCredentials = {
        ...credentials,
        platform: 'oa'
      }

      const response = await authApi.login(loginData)
      
      // 保存tokens和用户信息
      saveTokens(response.access_token, response.refresh_token)
      saveUserInfo(response.user)
      
      message.success('登录成功')
      return response
      
    } catch (error: any) {
      console.error('登录失败:', error)
      const errorMessage = error?.message || '登录失败，请检查用户名和密码'
      message.error(errorMessage)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 刷新token
  const refreshAccessToken = async () => {
    if (!refreshToken.value) {
      throw new Error('没有refresh token')
    }

    try {
      const response = await authApi.refreshToken(refreshToken.value)
      saveTokens(response.access_token, response.refresh_token)
      return response.access_token
    } catch (error) {
      console.error('刷新token失败:', error)
      // 刷新失败，清除所有认证信息
      await logout()
      throw error
    }
  }

  // 验证token
  const validateCurrentToken = async () => {
    if (!token.value) return false

    try {
      const response = await authApi.validateToken()
      return response.valid
    } catch (error) {
      console.error('验证token失败:', error)
      return false
    }
  }

  // 获取用户信息
  const fetchUserInfo = async () => {
    try {
      const user = await authApi.getCurrentUser()
      saveUserInfo(user)
      return user
    } catch (error) {
      console.error('获取用户信息失败:', error)
      throw error
    }
  }

  // 更新用户信息
  const updateUserProfile = async (data: Partial<Pick<OAUser, 'email' | 'phone' | 'avatar'>>) => {
    try {
      const updatedUser = await authApi.updateProfile(data)
      saveUserInfo(updatedUser)
      message.success('用户信息更新成功')
      return updatedUser
    } catch (error) {
      console.error('更新用户信息失败:', error)
      message.error('更新用户信息失败')
      throw error
    }
  }

  // 修改密码
  const changePassword = async (oldPassword: string, newPassword: string) => {
    try {
      await authApi.changePassword({
        old_password: oldPassword,
        new_password: newPassword
      })
      message.success('密码修改成功，请重新登录')
      await logout()
    } catch (error) {
      console.error('修改密码失败:', error)
      message.error('修改密码失败')
      throw error
    }
  }

  // 登出
  const logout = async () => {
    try {
      if (token.value) {
        await authApi.logout()
      }
    } catch (error) {
      console.error('登出请求失败:', error)
      // 即使请求失败也要清除本地数据
    } finally {
      clearTokens()
      message.success('已安全退出')
      
      // 跳转到登录页面
      if (router.currentRoute.value.path !== '/login') {
        router.push('/login')
      }
    }
  }

  // 初始化认证状态
  const initAuth = async () => {
    // 恢复用户信息
    restoreUserInfo()
    
    // 如果有token但没有用户信息，尝试获取
    if (token.value && !userInfo.value) {
      try {
        await fetchUserInfo()
      } catch (error) {
        console.error('初始化时获取用户信息失败:', error)
        clearTokens()
      }
    }
    
    // 验证token有效性
    if (token.value) {
      const isValid = await validateCurrentToken()
      if (!isValid) {
        clearTokens()
      }
    }
  }

  // 返回store API
  return {
    // 状态
    token,
    refreshToken,
    userInfo,
    isLoading,
    
    // 计算属性
    isLoggedIn,
    userName,
    userRole,
    userAvatar,
    
    // 方法
    hasRole,
    login,
    logout,
    refreshAccessToken,
    validateCurrentToken,
    fetchUserInfo,
    updateUserProfile,
    changePassword,
    initAuth
  }
}) 