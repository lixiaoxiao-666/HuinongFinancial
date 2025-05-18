import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'login',
    component: () => import('../views/LoginPage.vue'),
  },
  {
    path: '/',
    component: () => import('../components/AppLayout.vue'),
    children: [
      {
        path: '/home',
        name: 'home',
        component: () => import('../views/HomePage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: '/rent',
        name: 'rent',
        component: () => import('../views/RentManagementPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: '/rent/application',
        name: 'rent-application',
        component: () => import('../views/RentApplicationPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: '/machinery',
        name: 'machinery',
        component: () => import('../views/MachineryPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: '/loan',
        name: 'loan',
        component: () => import('../views/LoanManagementPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: '/loan/apply',
        name: 'loan-apply',
        component: () => import('../views/LoanManagementPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: '/loan/status',
        name: 'loan-status',
        component: () => import('../views/LoanManagementPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: '/orders',
        name: 'orders',
        component: () => import('../views/OrdersPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: '/users',
        name: 'users',
        component: () => import('../views/UsersPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: '/settings',
        name: 'settings',
        component: () => import('../views/SystemSettingsPage.vue'),
        meta: { requiresAuth: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// 简单的导航守卫，确保只有登录后才能访问需要认证的页面
router.beforeEach((to, from, next) => {
  // 这里应该检查用户是否已登录，例如通过检查localStorage中的token
  // 目前简单实现，实际项目中应该通过调用API验证token有效性
  const isAuthenticated = localStorage.getItem('token') !== null
  
  if (to.meta.requiresAuth && !isAuthenticated) {
    // 如果需要认证但用户未登录，重定向到登录页
    next({ name: 'login' })
  } else {
    // 否则继续
    next()
  }
})

export default router
