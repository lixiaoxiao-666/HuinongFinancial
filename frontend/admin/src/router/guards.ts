import type { Router } from 'vue-router'
import { useAuthStore } from '@/stores'
import { message } from 'ant-design-vue'

/**
 * 设置路由守卫
 */
export function setupRouterGuards(router: Router) {
  // 全局前置守卫
  router.beforeEach(async (to, from, next) => {
    const authStore = useAuthStore()

    // 设置页面标题
    document.title = to.meta?.title 
      ? `${to.meta.title} - 惠农OA管理系统` 
      : '惠农OA管理系统'

    // 检查是否需要认证
    if (to.meta?.requiresAuth !== false && to.path !== '/login') {
      if (!authStore.isAuthenticated) {
        // 未认证，跳转到登录页
        next({
          path: '/login',
          query: { redirect: to.fullPath }
        })
        return
      }

      // 验证Token有效性
      try {
        const isValid = await authStore.validateToken()
        if (!isValid) {
          next({
            path: '/login',
            query: { redirect: to.fullPath }
          })
          return
        }
      } catch (error) {
        console.error('Token验证失败:', error)
        next({
          path: '/login',
          query: { redirect: to.fullPath }
        })
        return
      }

      // 检查权限
      if (to.meta?.requiredRoles) {
        const roles = Array.isArray(to.meta.requiredRoles) 
          ? to.meta.requiredRoles 
          : [to.meta.requiredRoles]
        
        if (!authStore.hasRole(roles)) {
          message.error('权限不足，无法访问此页面')
          next('/403')
          return
        }
      }

      if (to.meta?.requiredPermissions) {
        const permissions = Array.isArray(to.meta.requiredPermissions)
          ? to.meta.requiredPermissions
          : [to.meta.requiredPermissions]
        
        const hasAllPermissions = permissions.every(permission => 
          authStore.hasPermission(permission)
        )
        
        if (!hasAllPermissions) {
          message.error('权限不足，无法访问此页面')
          next('/403')
          return
        }
      }
    }

    // 已登录用户访问登录页，重定向到首页
    if (to.path === '/login' && authStore.isAuthenticated) {
      next('/dashboard')
      return
    }

    next()
  })

  // 全局后置钩子
  router.afterEach((to, from) => {
    // 这里可以处理页面切换后的逻辑
    // 比如埋点、分析等
  })

  // 全局错误处理
  router.onError((error) => {
    console.error('路由错误:', error)
    message.error('页面加载失败，请刷新重试')
  })
} 