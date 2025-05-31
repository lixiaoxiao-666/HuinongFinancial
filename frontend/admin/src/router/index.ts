import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      name: 'layout',
      component: () => import('@/views/LayoutView.vue'),
      meta: { requiresAuth: true },
      redirect: '/dashboard',
      children: [
        {
          path: '/dashboard',
          name: 'dashboard',
          component: () => import('@/views/DashboardView.vue'),
          meta: { 
            requiresAuth: true,
            title: '工作台',
            permission: 'dashboard:view'
          }
        },
        {
          path: '/approval',
          name: 'approval',
          component: () => import('@/views/ApprovalView.vue'),
          meta: { 
            requiresAuth: true,
            title: '审批看板',
            permission: 'approval:view'
          }
        },
        {
          path: '/approval/:id',
          name: 'approval-detail',
          component: () => import('@/views/ApprovalDetailView.vue'),
          meta: { 
            requiresAuth: true,
            title: '审批详情',
            permission: 'approval:view'
          }
        },
        {
          path: '/smart-approval',
          name: 'smart-approval',
          component: () => import('@/views/SmartApprovalView.vue'),
          meta: { 
            requiresAuth: true,
            title: '智能审批',
            permission: 'approval:view'
          }
        },
        {
          path: '/ai-workflow',
          name: 'ai-workflow',
          component: () => import('@/views/AIWorkflowView.vue'),
          meta: { 
            requiresAuth: true,
            title: 'AI审批流管理',
            permission: 'approval:manage'
          }
        },
        {
          path: '/lease-approval',
          name: 'lease-approval',
          component: () => import('@/views/LeaseApprovalView.vue'),
          meta: { 
            requiresAuth: true,
            title: '租赁审批',
            permission: 'approval:view'
          }
        },
        {
          path: '/users',
          name: 'users',
          component: () => import('@/views/UsersView.vue'),
          meta: { 
            requiresAuth: true,
            title: '用户管理',
            permission: 'user:manage'
          }
        },
        {
          path: '/logs',
          name: 'logs',
          component: () => import('@/views/LogsView.vue'),
          meta: { 
            requiresAuth: true,
            title: '操作日志',
            permission: 'log:view'
          }
        },
        {
          path: '/system',
          name: 'system',
          component: () => import('@/views/SystemView.vue'),
          meta: { 
            requiresAuth: true,
            title: '系统设置',
            permission: 'system:manage'
          }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFoundView.vue')
    }
  ],
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // 检查是否需要认证
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next('/login')
    return
  }
  
  // 如果已登录且访问登录页，重定向到首页
  if (to.name === 'login' && authStore.isLoggedIn) {
    next('/dashboard')
    return
  }
  
  // 权限检查
  if (to.meta.permission && !authStore.hasPermission(to.meta.permission as string)) {
    // 如果没有权限，可以重定向到无权限页面或首页
    next('/dashboard')
    return
  }
  
  next()
})

export default router
