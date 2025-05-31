import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { AdminUser, LoginResponse } from '@/types'
import { adminLogin } from '@/api/admin'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('admin_token') || '')
  const user = ref<AdminUser | null>(
    localStorage.getItem('admin_user') 
      ? JSON.parse(localStorage.getItem('admin_user') as string)
      : null
  )
  const isLoggedIn = ref<boolean>(!!token.value)

  // 登录
  const login = async (username: string, password: string) => {
    try {
      const response = await adminLogin({ username, password })
      
      token.value = response.token
      user.value = {
        admin_user_id: response.admin_user_id,
        username: response.username,
        role: response.role as 'ADMIN' | '审批员',
        display_name: response.username,
        email: '',
        status: 0,
        created_at: '',
        updated_at: ''
      }
      isLoggedIn.value = true

      // 保存到localStorage
      localStorage.setItem('admin_token', response.token)
      localStorage.setItem('admin_user', JSON.stringify(user.value))

      return response
    } catch (error) {
      throw error
    }
  }

  // 登出
  const logout = () => {
    token.value = ''
    user.value = null
    isLoggedIn.value = false
    
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_user')
  }

  // 检查权限
  const hasPermission = (permission: string) => {
    if (!user.value) return false
    
    // 管理员拥有所有权限
    if (user.value.role === 'ADMIN') return true
    
    // 审批员权限检查
    const reviewerPermissions = [
      'approval:view',
      'approval:review',
      'dashboard:view'
    ]
    
    return reviewerPermissions.includes(permission)
  }

  return {
    token,
    user,
    isLoggedIn,
    login,
    logout,
    hasPermission
  }
}) 