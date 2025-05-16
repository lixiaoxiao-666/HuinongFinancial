import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '@/components/HomePage.vue'
import ShopPage from '@/components/ShopPage.vue'
import FinancePage from '../components/FinancePage.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: HomePage
  },
  {
    path: '/shop',
    name: 'Shop',
    component: ShopPage
  },
  {
    path: '/finance',
    name: 'Finance',
    component: FinancePage
  },
  {
    path: '/finance/product/:id',
    name: 'FinanceProduct',
    component: () => import('@/components/finance/ProductDetail.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
