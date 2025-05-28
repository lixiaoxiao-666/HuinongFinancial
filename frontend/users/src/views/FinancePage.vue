<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import AppFooter from './components/footer.vue'
import { useUserStore } from '../stores/user'
import { loanApi } from '../services/api'
import type { LoanProduct } from '../services/api'

const router = useRouter()
const userStore = useUserStore()
const activeTab = ref('finance')
const loading = ref(false)
const refreshing = ref(false)

// 贷款产品列表
const loanProducts = ref<LoanProduct[]>([])
const selectedCategory = ref('')

// 产品分类
const categories = [
  { value: '', label: '全部产品', icon: 'all' },
  { value: '种植贷', label: '种植贷', icon: 'plant' },
  { value: '设备贷', label: '设备贷', icon: 'machine' },
  { value: '养殖贷', label: '养殖贷', icon: 'livestock' },
  { value: '经营贷', label: '经营贷', icon: 'business' }
]

// 筛选后的产品列表
const filteredProducts = computed(() => {
  if (!selectedCategory.value) return loanProducts.value
  return loanProducts.value.filter(product => product.category === selectedCategory.value)
})

// 格式化金额
const formatAmount = (amount: number) => {
  if (amount >= 10000) {
    return `${(amount / 10000).toFixed(1)}万`
  }
  return amount.toLocaleString()
}

// 加载贷款产品
const loadLoanProducts = async () => {
  try {
    loading.value = true
    const response = await loanApi.getProducts(selectedCategory.value || undefined)
    loanProducts.value = response.data
  } catch (error: any) {
    console.error('加载贷款产品失败:', error)
    ElMessage.error('加载产品信息失败')
  } finally {
    loading.value = false
  }
}

// 刷新数据
const refreshData = async () => {
  try {
    refreshing.value = true
    await loadLoanProducts()
    ElMessage.success('刷新成功')
  } catch (error: any) {
    ElMessage.error('刷新失败')
  } finally {
    refreshing.value = false
  }
}

// 切换分类
const switchCategory = (category: string) => {
  selectedCategory.value = category
  loadLoanProducts()
}

// 申请贷款
const applyLoan = (product: LoanProduct) => {
  // 检查登录状态
  if (!userStore.isLoggedIn || !userStore.isTokenValid()) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  // 跳转到申请页面
  router.push(`/loan/apply/${product.product_id}`)
}

// 查看产品详情
const viewProductDetail = (productId: string) => {
  router.push(`/loan/products/${productId}`)
}

// 查看我的申请
const viewMyApplications = () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  router.push('/loan/my-applications')
}

// 组件挂载时加载数据
onMounted(() => {
  loadLoanProducts()
})
</script>

<template>
  <div class="finance-page">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <div class="nav-left">
        <el-icon @click="router.go(-1)"><ArrowLeft /></el-icon>
      </div>
      <div class="nav-title">惠农金融</div>
      <div class="nav-right">
        <el-icon @click="refreshData" :class="{ 'is-loading': refreshing }">
          <Refresh />
        </el-icon>
      </div>
    </div>

    <div class="page-content">
      <!-- 用户快捷操作 -->
      <div class="quick-actions" v-if="userStore.isLoggedIn">
        <div class="action-card primary" @click="viewMyApplications">
          <div class="card-icon">
            <svg viewBox="0 0 24 24" width="22" height="22" fill="white">
              <path d="M14,2H6C4.9,2,4,2.9,4,4v16c0,1.1,0.9,2,2,2h12c1.1,0,2-0.9,2-2V8L14,2z M16,18H8v-2h8V18z M16,14H8v-2h8V14z M13,9V3.5L18.5,9H13z"/>
            </svg>
          </div>
          <div class="card-content">
            <h3>我的申请</h3>
            <p>查看贷款申请进度</p>
          </div>
          <div class="card-arrow">
            <svg viewBox="0 0 24 24" width="18" height="18" fill="white" opacity="0.8">
              <path d="M8.59,16.59L13.17,12L8.59,7.41L10,6l6,6l-6,6L8.59,16.59z"/>
            </svg>
          </div>
        </div>
        
        <div class="stats-row">
          <div class="stat-item">
            <div class="stat-value">4</div>
            <div class="stat-label">申请笔数</div>
            <div class="stat-icon">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="#27ae60">
                <path d="M19,3H14.82C14.4,1.84,13.3,1,12,1S9.6,1.84,9.18,3H5C3.9,3,3,3.9,3,5v14c0,1.1,0.9,2,2,2h14c1.1,0,2-0.9,2-2V5 C21,3.9,20.1,3,19,3z M12,3c0.55,0,1,0.45,1,1s-0.45,1-1,1s-1-0.45-1-1S11.45,3,12,3z M10,17l-4-4l1.41-1.41L10,14.17l6.59-6.59 L18,9L10,17z"/>
              </svg>
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-value">3.2万</div>
            <div class="stat-label">申请总额</div>
            <div class="stat-icon">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="#27ae60">
                <path d="M11.8 10.9c-2.27-.59-3-1.2-3-2.15 0-1.09 1.01-1.85 2.7-1.85 1.78 0 2.44.85 2.5 2.1h2.21c-.07-1.72-1.12-3.3-3.21-3.81V3h-3v2.16c-1.94.42-3.5 1.68-3.5 3.61 0 2.31 1.91 3.46 4.7 4.13 2.5.6 3 1.48 3 2.41 0 .69-.49 1.79-2.7 1.79-2.06 0-2.87-.92-2.98-2.1h-2.2c.12 2.19 1.76 3.42 3.68 3.83V21h3v-2.15c1.95-.37 3.5-1.5 3.5-3.55 0-2.84-2.43-3.81-4.7-4.4z"/>
              </svg>
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-value">85</div>
            <div class="stat-label">信用评分</div>
            <div class="stat-icon">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="#27ae60">
                <path d="M12 17.27L18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z"/>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- 登录提示 -->
      <div class="login-prompt" v-else>
        <div class="prompt-content">
          <el-icon class="prompt-icon"><User /></el-icon>
          <p>登录后享受更多金融服务</p>
          <el-button type="primary" @click="router.push('/login')">
            立即登录
          </el-button>
        </div>
      </div>

      <!-- 产品分类筛选 -->
      <div class="category-filter">
        <div class="filter-header">
          <span class="filter-icon">
            <svg viewBox="0 0 24 24" width="20" height="20" fill="#27ae60">
              <path d="M10,18V16H8V14H10V12H12V14H14V16H12V18H10M3,4H21V8H19.5C19.5,8 19,8.5 19,9C19,9.5 19.5,10 19.5,10H21V14H19.5C19.5,14 19,14.5 19,15C19,15.5 19.5,16 19.5,16H21V20H3V16H4.5C4.5,16 5,15.5 5,15C5,14.5 4.5,14 4.5,14H3V10H4.5C4.5,10 5,9.5 5,9C5,8.5 4.5,8 4.5,8H3V4Z"/>
            </svg>
          </span>
          <span class="filter-title">产品分类</span>
        </div>
        <div class="category-tabs">
          <div 
            v-for="category in categories" 
            :key="category.value"
            class="category-tab"
            :class="{ 'active': selectedCategory === category.value }"
            @click="switchCategory(category.value)"
          >
            <span class="tab-icon" v-if="category.icon">
              <svg v-if="category.icon === 'all'" viewBox="0 0 24 24" width="16" height="16">
                <path d="M4,6H20V16H4V6M20,18V20H4V18H20M20,4V6H4V4H20M2,6V20H22V6H2Z" :fill="selectedCategory === category.value ? '#fff' : '#666'" />
              </svg>
              <svg v-else-if="category.icon === 'plant'" viewBox="0 0 24 24" width="16" height="16">
                <path d="M15,4V6H18V10H15.5V12H13.5V10H11V8H9V10H7V6H10V4H15M7,14H21V16H7V14M7,18H21V20H7V18Z" :fill="selectedCategory === category.value ? '#fff' : '#666'" />
              </svg>
              <svg v-else-if="category.icon === 'machine'" viewBox="0 0 24 24" width="16" height="16">
                <path d="M7,2V4H8V18H5V16H3V22H21V16H19V18H16V4H17V2H7M5,14H13V16H5V14Z" :fill="selectedCategory === category.value ? '#fff' : '#666'" />
              </svg>
              <svg v-else-if="category.icon === 'livestock'" viewBox="0 0 24 24" width="16" height="16">
                <path d="M19,5H17V3H19V5M16,5H8V3H16V5M7,5H5V3H7V5M22,3V14H19.82L20,15C20.67,15 21.77,15.29 21.77,16.56C21.71,17.73 20.86,18.06 20.42,18.11C20.31,18.12 20.19,18.16 20.07,18.2C19.94,18.24 19.8,18.31 19.66,18.36C19.59,18.39 19.54,18.43 19.43,18.47C19.18,18.58 19,18.59 18.97,18.77L18.06,22H18.03C17.62,22 17.39,21.55 17.37,21.06L17.28,18.16C17.28,17.23 16.62,17 16.05,17H13V21H12C11.23,21 10.62,20.91 9,19L7.32,15.91C7.16,15.68 6.69,15.34 6.3,15.31C5.93,15.28 5.74,15.44 5.7,15.75C5.67,16.04 5.96,16.44 6.25,16.64L7.95,18.22C8.55,18.76 9,21 9,21H8C6,21 3,18.88 3,16.92V15C3,12.9 4.35,12.07 5.74,12L8,8H17L19,13H22V3H22Z" :fill="selectedCategory === category.value ? '#fff' : '#666'" />
              </svg>
              <svg v-else-if="category.icon === 'business'" viewBox="0 0 24 24" width="16" height="16">
                <path d="M18,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V4A2,2 0 0,0 18,2M10,4H12V12L11,11.27L10,12V4M4,20H3V4H4V20M21,20H20V4H21V20Z" :fill="selectedCategory === category.value ? '#fff' : '#666'" />
              </svg>
            </span>
            <span class="tab-text">{{ category.label }}</span>
          </div>
        </div>
      </div>

      <!-- 产品列表 -->
      <div class="products-section">
        <div class="section-header">
          <h3>贷款产品</h3>
          <span class="product-count">{{ filteredProducts.length }}款产品</span>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="loading-container">
          <el-icon class="is-loading"><Loading /></el-icon>
          <p>加载中...</p>
        </div>

        <!-- 产品卡片列表 -->
        <div v-else class="product-list">
          <div 
            v-for="product in filteredProducts" 
            :key="product.product_id"
            class="product-card"
          >
            <div class="card-header">
              <div class="product-info">
                <h4 class="product-name">{{ product.name }}</h4>
                <p class="product-desc">{{ product.description }}</p>
                <div class="product-category">
                  <el-tag size="small" type="success">{{ product.category }}</el-tag>
                </div>
              </div>
            </div>

            <div class="card-body">
              <div class="product-details">
                <div class="detail-item">
                  <span class="label">贷款金额</span>
                  <span class="value">{{ formatAmount(product.min_amount) }} - {{ formatAmount(product.max_amount) }}元</span>
                </div>
                <div class="detail-item">
                  <span class="label">贷款期限</span>
                  <span class="value">{{ product.min_term_months }} - {{ product.max_term_months }}个月</span>
                </div>
                <div class="detail-item">
                  <span class="label">年利率</span>
                  <span class="value highlight">{{ product.interest_rate_yearly }}</span>
                </div>
              </div>
            </div>

            <div class="card-footer">
              <el-button 
                type="info" 
                size="small" 
                @click="viewProductDetail(product.product_id)"
              >
                查看详情
              </el-button>
              <el-button 
                type="primary" 
                size="small"
                @click="applyLoan(product)"
              >
                立即申请
              </el-button>
            </div>
          </div>

          <!-- 空状态 -->
          <div v-if="filteredProducts.length === 0" class="empty-state">
            <el-empty description="暂无相关产品">
              <el-button type="primary" @click="switchCategory('')">
                查看全部产品
              </el-button>
            </el-empty>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部导航 -->
    <app-footer v-model:active-tab="activeTab" />
  </div>
</template>

<style scoped>
.finance-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 80px;
}

.top-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: white;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-left, .nav-right {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #2c3e50;
}

.nav-title {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
}

.page-content {
  padding: 16px;
  max-width: 600px;
  margin: 0 auto;
}

.quick-actions {
  margin-bottom: 16px;
}

.action-card {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  border-radius: 16px;
  padding: 22px;
  color: white;
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 16px;
  box-shadow: 0 8px 16px rgba(39, 174, 96, 0.2);
  position: relative;
  overflow: hidden;
}

.action-card:before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(45deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0) 100%);
  z-index: 1;
}

.action-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 12px 28px rgba(39, 174, 96, 0.3);
}

.card-icon {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  position: relative;
  z-index: 2;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.card-content {
  flex: 1;
  position: relative;
  z-index: 2;
}

.card-content h3 {
  margin: 0 0 6px;
  font-size: 20px;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.card-content p {
  margin: 0;
  font-size: 14px;
  opacity: 0.9;
  font-weight: 400;
}

.card-arrow {
  opacity: 0.8;
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 8px 16px rgba(0,0,0,0.08);
  position: relative;
  overflow: hidden;
}

.stats-row:before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #27ae60, #2ecc71);
}

.stat-item {
  text-align: center;
  position: relative;
  padding: 8px 0;
}

.stat-item:not(:last-child):after {
  content: '';
  position: absolute;
  right: -8px;
  top: 20%;
  height: 60%;
  width: 1px;
  background: #f0f0f0;
}

.stat-value {
  font-size: 26px;
  font-weight: 700;
  color: #27ae60;
  margin-bottom: 8px;
  position: relative;
}

.stat-label {
  font-size: 13px;
  color: #666;
  font-weight: 500;
}

.stat-icon {
  position: absolute;
  top: 8px;
  right: 16px;
  opacity: 0.2;
}

.login-prompt {
  background: white;
  border-radius: 12px;
  padding: 24px;
  text-align: center;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.prompt-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.prompt-icon {
  font-size: 32px;
  color: #27ae60;
}

.category-filter {
  background: white;
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);
  position: relative;
  overflow: hidden;
}

.category-filter:before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
  background: linear-gradient(to bottom, #27ae60, #2ecc71);
}

.filter-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.filter-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 10px;
}

.filter-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
}

.category-tabs {
  display: flex;
  gap: 12px;
  overflow-x: auto;
  scrollbar-width: none;
  padding-bottom: 4px;
  -webkit-overflow-scrolling: touch;
}

.category-tabs::-webkit-scrollbar {
  display: none;
}

.category-tab {
  padding: 8px 16px;
  border-radius: 12px;
  border: 2px solid #e1e1e1;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.3s;
  background: white;
  color: #666;
  display: flex;
  align-items: center;
  gap: 6px;
}

.category-tab.active {
  background: #27ae60;
  color: white;
  border-color: #27ae60;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(39, 174, 96, 0.2);
}

.category-tab:hover:not(.active) {
  border-color: #27ae60;
  color: #27ae60;
  background-color: rgba(39, 174, 96, 0.05);
}

.tab-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.tab-text {
  display: inline-block;
  position: relative;
}

.category-tab.active .tab-text:after {
  content: '';
  position: absolute;
  bottom: -3px;
  left: 0;
  width: 100%;
  height: 2px;
  background-color: rgba(255, 255, 255, 0.5);
  transform: scaleX(0);
  transition: transform 0.3s ease;
}

.category-tab.active:hover .tab-text:after {
  transform: scaleX(1);
}

.products-section {
  background: white;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.section-header h3 {
  margin: 0;
  font-size: 16px;
  color: #2c3e50;
}

.product-count {
  font-size: 12px;
  color: #7f8c8d;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: #7f8c8d;
}

.loading-container .el-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.product-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.product-card {
  border: 1px solid #e1e1e1;
  border-radius: 12px;
  padding: 16px;
  transition: all 0.3s;
}

.product-card:hover {
  border-color: #27ae60;
  box-shadow: 0 4px 12px rgba(39, 174, 96, 0.1);
}

.card-header {
  margin-bottom: 12px;
}

.product-name {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.product-desc {
  font-size: 14px;
  color: #7f8c8d;
  margin: 0 0 8px 0;
  line-height: 1.4;
}

.card-body {
  margin-bottom: 16px;
}

.product-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.label {
  font-size: 14px;
  color: #7f8c8d;
}

.value {
  font-size: 14px;
  color: #2c3e50;
  font-weight: 500;
}

.value.highlight {
  color: #27ae60;
  font-weight: 600;
}

.card-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
}

@media (max-width: 768px) {
  .page-content {
    padding: 12px;
  }
  
  .action-card {
    padding: 18px;
  }
  
  .card-icon {
    width: 42px;
    height: 42px;
  }
  
  .card-content h3 {
    font-size: 18px;
  }
  
  .stats-row {
    padding: 16px 10px;
  }
  
  .stat-value {
    font-size: 22px;
  }
  
  .stat-label {
    font-size: 12px;
  }
  
  .stat-icon {
    display: none;
  }
  
  .card-footer {
    flex-direction: column;
    gap: 8px;
  }
  
  .card-footer .el-button {
    width: 100%;
  }
  
  .category-filter {
    padding: 16px;
    margin-bottom: 16px;
  }
  
  .category-tab {
    padding: 6px 14px;
    font-size: 13px;
  }
}
</style> 