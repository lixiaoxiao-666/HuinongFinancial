import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import type { Permission } from '@/types/auth'

// 配置NProgress
NProgress.configure({
  showSpinner: false,
  speed: 200,
  minimum: 0.1
})

// 路由组件懒加载
const LayoutView = () => import('@/views/LayoutView.vue')
const LoginView = () => import('@/views/LoginView.vue')
const DashboardView = () => import('@/views/DashboardView.vue')
const ApprovalView = () => import('@/views/ApprovalView.vue')
const ApprovalDetailView = () => import('@/views/ApprovalDetailView.vue')
const UsersView = () => import('@/views/UsersView.vue')
const SystemView = () => import('@/views/SystemView.vue')
const LogsView = () => import('@/views/LogsView.vue')
const NotFoundView = () => import('@/views/NotFoundView.vue')

// 路由配置
const routes = [
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
    meta: {
      title: '管理员登录',
      requiresAuth: false,
      hideInMenu: true
    }
  },
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    component: LayoutView,
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: DashboardView,
        meta: {
          title: '工作台',
          icon: 'Dashboard',
          requiresAuth: true,
          permissions: ['dashboard:view']
        }
      }
    ]
  },
  {
    path: '/approval',
    component: LayoutView,
    meta: {
      title: '贷款审批',
      icon: 'DocumentChecked',
      requiresAuth: true,
      permissions: ['loan_approve']
    },
    children: [
      {
        path: '',
        name: 'ApprovalList',
        component: ApprovalView,
        meta: {
          title: '申请列表',
          requiresAuth: true,
          permissions: ['loan_approve']
        }
      },
      {
        path: 'detail/:id',
        name: 'ApprovalDetail',
        component: ApprovalDetailView,
        meta: {
          title: '审批详情',
          requiresAuth: true,
          permissions: ['loan_approve'],
          hideInMenu: true
        }
      }
    ]
  },
  {
    path: '/users',
    component: LayoutView,
    children: [
      {
        path: '',
        name: 'Users',
        component: UsersView,
        meta: {
          title: '用户管理',
          icon: 'User',
          requiresAuth: true,
          permissions: ['user_manage']
        }
      }
    ]
  },
  {
    path: '/system',
    component: LayoutView,
    meta: {
      title: '系统管理',
      icon: 'Setting',
      requiresAuth: true,
      permissions: ['system_config']
    },
    children: [
      {
        path: '',
        name: 'System',
        component: SystemView,
        meta: {
          title: '系统配置',
          requiresAuth: true,
          permissions: ['system_config']
        }
      },
      {
        path: 'logs',
        name: 'Logs',
        component: LogsView,
        meta: {
          title: '操作日志',
          requiresAuth: true,
          permissions: ['log_view']
        }
      }
    ]
  },
  {
    path: '/session',
    component: LayoutView,
    children: [
      {
        path: 'management',
        name: 'SessionManagement',
        component: () => import('@/views/SessionManagementView.vue'),
        meta: {
          title: '会话管理',
          icon: 'Connection',
          requiresAuth: true,
          permissions: ['system_config']
        }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFoundView,
    meta: {
      title: '404 - 页面不存在',
      hideInMenu: true
    }
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// 全局路由守卫
router.beforeEach(async (to, from, next) => {
  // 开始进度条
  NProgress.start()
  
  // 设置页面标题
  const title = to.meta?.title as string
  document.title = title ? `${title} - 数字惠农OA管理系统` : '数字惠农OA管理系统'
  
  const authStore = useAuthStore()
  
  // 如果路由不需要认证，直接通过
  if (!to.meta?.requiresAuth) {
    // 如果已登录且要访问登录页，重定向到首页
    if (to.name === 'Login' && authStore.isAuthenticated) {
      next('/dashboard')
      return
    }
    next()
    return
  }
  
  // 检查是否已登录
  if (!authStore.isAuthenticated) {
    ElMessage.warning('请先登录')
    next({
      path: '/login',
      query: { redirect: to.fullPath }
    })
    return
  }
  
  // 验证Token有效性
  try {
    const isValid = await authStore.verifyAuth()
    if (!isValid) {
      ElMessage.error('登录已过期，请重新登录')
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
      return
    }
  } catch (error) {
    console.error('验证认证状态失败:', error)
    next({
      path: '/login',
      query: { redirect: to.fullPath }
    })
    return
  }
  
  // 检查权限
  const requiredPermissions = to.meta?.permissions as Permission[]
  if (requiredPermissions && requiredPermissions.length > 0) {
    const hasPermission = authStore.hasAnyPermission(requiredPermissions)
    if (!hasPermission) {
      ElMessage.error('权限不足，无法访问该页面')
      // 如果是从其他页面跳转过来的，返回上一页
      if (from.name && from.name !== 'Login') {
        next(false)
      } else {
        // 否则跳转到首页
        next('/dashboard')
      }
      return
    }
  }
  
  next()
})

// 路由完成后的处理
router.afterEach((to, from) => {
  // 结束进度条
  NProgress.done()
  
  // 记录页面访问日志（开发环境）
  if (import.meta.env.DEV) {
    console.log(`路由跳转: ${from.path} → ${to.path}`)
  }
})

// 路由错误处理
router.onError((error) => {
  console.error('路由错误:', error)
  NProgress.done()
  ElMessage.error('页面加载失败，请刷新重试')
})

// 生成面包屑导航
export const generateBreadcrumb = (route: any) => {
  const matched = route.matched.filter((item: any) => item.meta && item.meta.title)
  const breadcrumbs = matched.map((item: any) => ({
    title: item.meta.title,
    path: item.path
  }))
  
  return breadcrumbs
}

// 获取菜单权限
export const getMenuPermissions = () => {
  const authStore = useAuthStore()
  
  return routes
    .filter(route => !route.meta?.hideInMenu && route.meta?.requiresAuth)
    .filter(route => {
      const permissions = route.meta?.permissions as Permission[]
      if (!permissions || permissions.length === 0) return true
      return authStore.hasAnyPermission(permissions)
    })
    .map(route => ({
      path: route.path,
      name: route.name,
      title: route.meta?.title,
      icon: route.meta?.icon,
      children: route.children?.filter(child => 
        !child.meta?.hideInMenu && 
        (!child.meta?.permissions || authStore.hasAnyPermission(child.meta.permissions as Permission[]))
      )
    }))
}

export default router
