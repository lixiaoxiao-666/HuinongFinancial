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
            <i class="el-icon-document"></i>
          </div>
          <div class="card-content">
            <h3>我的申请</h3>
            <p>查看贷款申请进度</p>
          </div>
          <div class="card-arrow">
            <i class="el-icon-arrow-right"></i>
          </div>
        </div>
        
        <div class="stats-row">
          <div class="stat-item">
            <div class="stat-value">0</div>
            <div class="stat-label">申请笔数</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">0万</div>
            <div class="stat-label">申请总额</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">85</div>
            <div class="stat-label">信用评分</div>
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
        <div class="filter-title">产品分类</div>
        <div class="category-tabs">
          <div 
            v-for="category in categories" 
            :key="category.value"
            class="category-tab"
            :class="{ 'active': selectedCategory === category.value }"
            @click="switchCategory(category.value)"
          >
            {{ category.label }}
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
  border-radius: 12px;
  padding: 20px;
  color: white;
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 15px;
}

.action-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(39, 174, 96, 0.3);
}

.card-icon {
  width: 40px;
  height: 40px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 15px;
}

.card-icon i {
  font-size: 20px;
}

.card-content {
  flex: 1;
}

.card-content h3 {
  margin: 0 0 5px;
  font-size: 18px;
  font-weight: 600;
}

.card-content p {
  margin: 0;
  font-size: 14px;
  opacity: 0.9;
}

.card-arrow {
  opacity: 0.7;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 15px;
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #27ae60;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 12px;
  color: #666;
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
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.filter-title {
  font-size: 14px;
  color: #7f8c8d;
  margin-bottom: 12px;
}

.category-tabs {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding-bottom: 4px;
}

.category-tab {
  padding: 6px 16px;
  border-radius: 20px;
  border: 1px solid #e1e1e1;
  font-size: 14px;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.3s;
  background: white;
  color: #666;
}

.category-tab.active {
  background: #27ae60;
  color: white;
  border-color: #27ae60;
}

.category-tab:hover {
  border-color: #27ae60;
  color: #27ae60;
}

.category-tab.active:hover {
  color: white;
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
  
  .product-card {
    padding: 12px;
  }
  
  .card-footer {
    flex-direction: column;
    gap: 8px;
  }
  
  .card-footer .el-button {
    width: 100%;
  }
}
</style> 