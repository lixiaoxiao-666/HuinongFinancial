import { useUserStore } from '@/stores/user'
import type { RouteLocationNormalized, NavigationGuardNext } from 'vue-router'

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

// 路由守卫
export const authGuard = (
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
) => {
  const userStore = useUserStore()
  
  // 检查是否需要登录
  if (isAuthRequired(to.path)) {
    // 检查是否已登录且token有效
    if (!userStore.isLoggedIn || !userStore.isTokenValid()) {
      // 未登录，跳转到登录页
      next({
        path: '/login',
        query: { redirect: to.fullPath } // 保存原始路径，登录后可以回到原页面
      })
      return
    }
  }
  
  // 如果已登录用户访问登录页，重定向到首页
  if (to.path === '/login' && userStore.isLoggedIn && userStore.isTokenValid()) {
    next('/home')
    return
  }
  
  next()
}

// 检查用户权限
export const hasPermission = (permission: string): boolean => {
  const userStore = useUserStore()
  
  if (!userStore.isLoggedIn || !userStore.userInfo) {
    return false
  }
  
  // 这里可以根据用户角色或权限进行判断
  // 目前简单返回true，因为普通用户都有基本权限
  return true
}

// 格式化错误信息
export const getErrorMessage = (error: any): string => {
  if (typeof error === 'string') {
    return error
  }
  
  if (error?.message) {
    return error.message
  }
  
  if (error?.response?.data?.message) {
    return error.response.data.message
  }
  
  return '操作失败，请稍后重试'
} 