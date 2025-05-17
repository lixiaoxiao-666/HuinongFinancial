import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/IndexPage.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/login/LoginPage.vue'),
    },
    {
      path: '/finance',
      name: 'finance',
      component: () => import('../views/FinancePage.vue'),
    },
    {
      path: '/machinery',
      name: 'machinery',
      component: () => import('../views/MachineryPage.vue'),
    },
    {
      path: '/me',
      name: 'me',
      component: () => import('../views/MePage.vue'),
    },
    {
      path: '/news',
      name: 'news',
      component: () => import('../views/NewsPage.vue'),
    }
  ],
})

export default router
