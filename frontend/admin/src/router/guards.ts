import type { Router } from 'vue-router'
import { useAuthStore } from '@/stores/modules/auth'
import { message } from 'ant-design-vue'

// 白名单路由（不需要认证）
const whiteList = ['/login', '/404', '/403']

/**
 * 设置路由守卫
 */
export function setupRouterGuards(router: Router) {
  // 全局前置守卫
  router.beforeEach(async (to, from, next) => {
    console.log('🚀 路由守卫:', { from: from.path, to: to.path })
    
    const authStore = useAuthStore()
    
    // 初始化认证状态（仅在应用启动时执行一次）
    if (!authStore.token && !authStore.userInfo) {
      try {
        await authStore.initAuth()
      } catch (error) {
        console.error('认证状态初始化失败:', error)
      }
    }
    
    // 检查是否在白名单中
    if (whiteList.includes(to.path)) {
      // 如果已登录且访问登录页，重定向到仪表盘
      if (to.path === '/login' && authStore.isLoggedIn) {
        next('/dashboard')
        return
      }
      next()
      return
    }
    
    // 检查是否已登录
    if (!authStore.isLoggedIn) {
      console.log('❌ 未登录，重定向到登录页')
      message.warning('请先登录')
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
      return
    }
    
    // 验证token有效性
    try {
      const isValid = await authStore.validateCurrentToken()
      if (!isValid) {
        console.log('❌ Token无效，重定向到登录页')
        message.error('登录已过期，请重新登录')
        next({
          path: '/login',
          query: { redirect: to.fullPath }
        })
        return
      }
    } catch (error) {
      console.error('Token验证失败:', error)
      message.error('认证验证失败，请重新登录')
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
      return
    }
    
    // 角色权限检查（可选）
    if (to.meta?.roles) {
      const requiredRoles = to.meta.roles as string[]
      if (!authStore.hasRole(requiredRoles)) {
        console.log('❌ 权限不足:', { required: requiredRoles, current: authStore.userRole })
        message.error('权限不足，无法访问该页面')
        next('/403')
        return
      }
    }
    
    console.log('✅ 路由守卫通过')
    next()
  })
  
  // 全局后置守卫
  router.afterEach((to, from) => {
    // 设置页面标题
    const title = to.meta?.title as string
    if (title) {
      document.title = `${title} - 惠农OA管理系统`
    } else {
      document.title = '惠农OA管理系统'
    }
    
    // 移除加载状态
    const loadingElement = document.getElementById('app-loading')
    if (loadingElement) {
      loadingElement.style.display = 'none'
    }
    
    console.log('🎯 路由跳转完成:', { from: from.path, to: to.path, title })
  })
  
  // 路由错误处理
  router.onError((error) => {
    console.error('🚨 路由错误:', error)
    message.error('页面加载失败')
  })
} 