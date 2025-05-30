import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { setupRouterGuards } from './guards'

// 基础布局组件
const Layout = () => import('@/components/layout/MainLayout.vue')

// 路由配置
const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: {
      title: '登录',
      requiresAuth: false
    }
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: {
          title: '工作台',
          icon: 'dashboard',
          requiresAuth: true
        }
      }
    ]
  },
  {
    path: '/loans',
    component: Layout,
    meta: {
      title: '贷款管理',
      icon: 'bank',
      requiresAuth: true,
      requiredPermissions: ['loan:view']
    },
    children: [
      {
        path: 'applications',
        name: 'LoanApplications',
        component: () => import('@/views/loans/Applications.vue'),
        meta: {
          title: '申请列表',
          requiresAuth: true,
          requiredPermissions: ['loan:view']
        }
      },
      {
        path: 'applications/:id',
        name: 'LoanApplicationDetail',
        component: () => import('@/views/loans/ApplicationDetail.vue'),
        meta: {
          title: '申请详情',
          requiresAuth: true,
          requiredPermissions: ['loan:view'],
          hidden: true
        }
      },
      {
        path: 'approval',
        name: 'LoanApproval',
        component: () => import('@/views/loans/Approval.vue'),
        meta: {
          title: '审批工作台',
          requiresAuth: true,
          requiredPermissions: ['loan:approve']
        }
      },
      {
        path: 'statistics',
        name: 'LoanStatistics',
        component: () => import('@/views/loans/Statistics.vue'),
        meta: {
          title: '贷款统计',
          requiresAuth: true,
          requiredPermissions: ['loan:statistics']
        }
      }
    ]
  },
  {
    path: '/machines',
    component: Layout,
    meta: {
      title: '农机管理',
      icon: 'tool',
      requiresAuth: true,
      requiredPermissions: ['machine:view']
    },
    children: [
      {
        path: 'orders',
        name: 'MachineOrders',
        component: () => import('@/views/machines/Orders.vue'),
        meta: {
          title: '订单管理',
          requiresAuth: true,
          requiredPermissions: ['machine:view']
        }
      },
      {
        path: 'orders/:id',
        name: 'MachineOrderDetail',
        component: () => import('@/views/machines/OrderDetail.vue'),
        meta: {
          title: '订单详情',
          requiresAuth: true,
          requiredPermissions: ['machine:view'],
          hidden: true
        }
      }
    ]
  },
  {
    path: '/users',
    component: Layout,
    meta: {
      title: '用户管理',
      icon: 'user',
      requiresAuth: true,
      requiredPermissions: ['user:view']
    },
    children: [
      {
        path: 'list',
        name: 'UserList',
        component: () => import('@/views/users/List.vue'),
        meta: {
          title: '用户列表',
          requiresAuth: true,
          requiredPermissions: ['user:view']
        }
      },
      {
        path: 'detail/:id',
        name: 'UserDetail',
        component: () => import('@/views/users/Detail.vue'),
        meta: {
          title: '用户详情',
          requiresAuth: true,
          requiredPermissions: ['user:view'],
          hidden: true
        }
      }
    ]
  },
  {
    path: '/sessions',
    component: Layout,
    meta: {
      title: '会话管理',
      icon: 'global',
      requiresAuth: true,
      requiredRoles: ['admin']
    },
    children: [
      {
        path: 'active',
        name: 'ActiveSessions',
        component: () => import('@/views/sessions/Active.vue'),
        meta: {
          title: '活跃会话',
          requiresAuth: true,
          requiredRoles: ['admin']
        }
      },
      {
        path: 'statistics',
        name: 'SessionStatistics',
        component: () => import('@/views/sessions/Statistics.vue'),
        meta: {
          title: '会话统计',
          requiresAuth: true,
          requiredRoles: ['admin']
        }
      }
    ]
  },
  {
    path: '/system',
    component: Layout,
    meta: {
      title: '系统管理',
      icon: 'setting',
      requiresAuth: true,
      requiredRoles: ['admin']
    },
    children: [
      {
        path: 'monitor',
        name: 'SystemMonitor',
        component: () => import('@/views/system/Monitor.vue'),
        meta: {
          title: '系统监控',
          requiresAuth: true,
          requiredRoles: ['admin']
        }
      },
      {
        path: 'logs',
        name: 'SystemLogs',
        component: () => import('@/views/system/Logs.vue'),
        meta: {
          title: '操作日志',
          requiresAuth: true,
          requiredRoles: ['admin']
        }
      }
    ]
  },
  {
    path: '/profile',
    component: Layout,
    meta: {
      hidden: true
    },
    children: [
      {
        path: '',
        name: 'Profile',
        component: () => import('@/views/profile/index.vue'),
        meta: {
          title: '个人中心',
          requiresAuth: true
        }
      }
    ]
  },
  // 错误页面
  {
    path: '/403',
    name: 'Forbidden',
    component: () => import('@/views/error/403.vue'),
    meta: {
      title: '403 - 权限不足',
      requiresAuth: false,
      hidden: true
    }
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: {
      title: '404 - 页面不存在',
      requiresAuth: false,
      hidden: true
    }
  },
  {
    path: '/500',
    name: 'ServerError',
    component: () => import('@/views/error/500.vue'),
    meta: {
      title: '500 - 服务器错误',
      requiresAuth: false,
      hidden: true
    }
  },
  // 捕获所有未匹配的路由
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404'
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    return { top: 0 }
  }
})

// 设置路由守卫
setupRouterGuards(router)

export default router

// 导出路由类型声明
declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    icon?: string
    requiresAuth?: boolean
    requiredRoles?: string | string[]
    requiredPermissions?: string | string[]
    hidden?: boolean
  }
}
