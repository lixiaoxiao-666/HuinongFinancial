import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // 登录页
    {
      path: '/',
      name: 'login',
      component: () => import('../views/LoginPage.vue'),
    },
    // 访问 '/login' 重定向到登录页
    {
      path: '/login',
      redirect: '/',
    },
    // 首页
    {
      path: '/index',
      name: 'index',
      component: () => import('../views/IndexPage.vue'),
    },
    // 我的
    {
      path: '/me',
      name: 'me',
      component: () => import('../views/MePage.vue'),
    },
  ],
})

export default router
