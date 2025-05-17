<script setup lang="ts">
import { ref } from 'vue'
import AppFooter from './components/footer.vue'
import FinancialProduct from './components/FinancialProduct.vue'
import '../assets/icons/iconfont.css'
import '../assets/icons/agri-icons.css'

// 理财页面待开发

// 当前选中的导航栏
const activeTab = ref('finance')

// 产品分类数据
const productCategories = [
  { icon: 'pension', name: '养老保障', color: '#4CAF50' },
  { icon: 'insurance', name: '保险产品', color: '#4CAF50' },
  { icon: 'bank', name: '银行产品', color: '#4CAF50' },
  { icon: 'securities', name: '券商产品', color: '#4CAF50' },
  { icon: 'xiaojin', name: '惠农小金', color: '#4CAF50' },
  { icon: 'salary', name: '存工资', color: '#4CAF50' },
  { icon: 'fixed', name: '定期精选', color: '#4CAF50' },
  { icon: 'zone', name: '惠农专区', color: '#4CAF50' },
]

// 市场指数数据
const marketIndices = [
  { name: '上证指数', value: '3104.40', change: '+1.63%', isUp: true },
  { name: '深证成指', value: '10004.60', change: '-0.82%', isUp: false },
]

// 表现最佳基金
const topFunds = [
  {
    rank: 1,
    returnRate: '188.11%',
    period: '近一周涨幅',
    name: '惠农致富证券公司指数分级',
    type: '股票型',
    tag: '',
  },
  {
    rank: 2,
    returnRate: '170.51%',
    period: '近一周涨幅',
    name: '惠农服务证券分级',
    type: '股票型',
    tag: '',
  },
  {
    rank: 3,
    returnRate: '50.56%',
    period: '近一年涨幅额',
    name: '惠农富中证全指证券公司',
    type: '指数型',
    tag: '1年新高',
  },
]

// 农村金融产品
const agriProducts = [
  {
    type: '理财产品',
    name: '惠农丰收理财180天',
    annualReturn: '4.80%',
    minPeriod: '180天',
    riskLevel: '低',
    agriTag: '农业产业',
    isRecommended: true
  },
  {
    type: '定期存款',
    name: '农资种植贷专享存款',
    annualReturn: '3.85%',
    minPeriod: '1年',
    riskLevel: '低',
    agriTag: '种植专享',
    isRecommended: false
  },
  {
    type: '基金组合',
    name: '惠农优选平衡型基金组合',
    annualReturn: '8.50%',
    minPeriod: '无固定期限',
    riskLevel: '中',
    agriTag: '乡村振兴',
    isRecommended: true
  },
  {
    type: '保险',
    name: '农业生产保障保险',
    annualReturn: '3.50%',
    minPeriod: '5年',
    riskLevel: '低',
    agriTag: '灾害保障',
    isRecommended: false
  }
]

// 当前选择的产品分类
const selectedCategory = ref('全部')

// 切换产品分类
const switchCategory = (category: string) => {
  selectedCategory.value = category
}
</script>

<template>
  <div class="finance-container">
    <!-- 顶部导航栏 -->
    <div class="top-nav">
      <div class="nav-left">
        <!-- 删除返回图标 -->
      </div>
      <div class="nav-title">理财</div>
      <div class="nav-right">
        <!-- 删除筛选图标 -->
      </div>
    </div>

    <!-- 产品分类 -->
    <div class="product-categories">
      <div class="category-item" v-for="(item, index) in productCategories" :key="index">
        <div class="category-icon">
          <span class="agri-icon" :class="`agri-icon-${item.icon}`"></span>
        </div>
        <div class="category-name">{{ item.name }}</div>
      </div>
    </div>

    <!-- 市场指数 -->
    <div class="market-indices">
      <div v-for="(item, index) in marketIndices" :key="index" class="market-index-item">
        <span class="index-name">{{ item.name }}</span>
        <span class="index-value">{{ item.value }}</span>
        <span class="index-change" :class="{ 'up': item.isUp, 'down': !item.isUp }">
          {{ item.change }}
        </span>
      </div>
    </div>

    <!-- 推荐卡片 -->
    <div class="recommendation-cards">
      <div class="card red-card">
        <div class="card-tag">
          <span class="agri-icon agri-icon-chart" style="width: 18px; height: 18px;"></span>
          <span>指数低谷</span>
        </div>
        <div class="card-title">惠农银行指数进入低位区</div>
        <div class="card-subtitle">上涨潜力大 值得持有</div>
        <div class="card-button">去看看</div>
      </div>

      <div class="card blue-card">
        <div class="card-tag">
          <span class="agri-icon agri-icon-chart" style="width: 18px; height: 18px;"></span>
          <span>为农服务</span>
        </div>
        <div class="card-title">惠农市场产品</div>
        <div class="card-subtitle">别错过</div>
        <div class="card-button">查看详情</div>
      </div>
    </div>

    <!-- 表现最佳 -->
    <div class="top-performers">
      <div class="section-header">
        <span class="section-title">表现最佳</span>
        <span class="more-link">更多</span>
      </div>

      <div class="fund-list">
        <div v-for="(fund, index) in topFunds" :key="index" class="fund-item">
          <div class="fund-rank">
            <span class="crown" :class="`rank-${fund.rank}`">
              <i class="iconfont icon-crown"></i>
            </span>
          </div>
          
          <div class="fund-return">
            <div class="return-rate">{{ fund.returnRate }}</div>
            <div class="return-period">{{ fund.period }}</div>
          </div>
          
          <div class="fund-info">
            <div class="fund-name">{{ fund.name }}</div>
            <div class="fund-details">
              <span class="fund-type">{{ fund.type }}</span>
              <span v-if="fund.tag" class="fund-tag">{{ fund.tag }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 农村金融产品 -->
    <div class="agri-products">
      <div class="section-header">
        <span class="section-title">惠农金融产品</span>
        <span class="category-tabs">
          <span 
            class="category-tab" 
            :class="{ 'active': selectedCategory === '全部' }" 
            @click="switchCategory('全部')"
          >
            全部
          </span>
          <span 
            class="category-tab" 
            :class="{ 'active': selectedCategory === '理财' }" 
            @click="switchCategory('理财')"
          >
            理财
          </span>
          <span 
            class="category-tab" 
            :class="{ 'active': selectedCategory === '存款' }" 
            @click="switchCategory('存款')"
          >
            存款
          </span>
        </span>
      </div>
      
      <div class="product-list">
        <financial-product 
          v-for="(product, index) in agriProducts" 
          :key="index"
          :type="product.type"
          :name="product.name"
          :annual-return="product.annualReturn"
          :min-period="product.minPeriod"
          :risk-level="product.riskLevel"
          :agri-tag="product.agriTag"
          :is-recommended="product.isRecommended"
        />
      </div>
    </div>
    
    <!-- 底部导航栏 -->
    <app-footer v-model:active-tab="activeTab" />
  </div>
</template>

<style scoped>
.finance-container {
  padding: 0;
  text-align: left;
  min-height: 100vh;
  padding-bottom: 60px; /* 为底部导航栏留出空间 */
  background-color: #f5f5f5;
  color: #333;
  font-size: 13px;
}

@media screen and (max-width: 375px) {
  .finance-container {
    font-size: 12px;
  }
}

/* 顶部导航栏 */
.top-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 44px;
  background-color: #fff;
  padding: 0 12px;
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-title {
  font-size: 22px;
  font-weight: 500;
}

.nav-back, .nav-filter {
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 产品分类 */
.product-categories {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  padding: 10px;
  background-color: #fff;
  margin: 0 8px 8px;
  gap: 6px;
  border-radius: 8px;
  width: calc(100% - 16px);
}

@media screen and (max-width: 375px) {
  .product-categories {
    padding: 10px;
    margin: 0 10px 6px;
    gap: 6px;
  }
}

.category-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
}

.category-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rgba(76, 175, 80, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 5px;
}

.agri-icon {
  width: 24px;
  height: 24px;
}

.category-name {
  font-size: 13px;
}

/* 市场指数 */
.market-indices {
  display: flex;
  background-color: #fff;
  padding: 10px 12px;
  gap: 15px;
  margin: 0 8px 8px;
  border-radius: 8px;
  width: calc(100% - 16px);
}

@media screen and (max-width: 375px) {
  .market-indices {
    padding: 10px;
    gap: 12px;
    margin: 0 10px 6px;
  }
}

.market-index-item {
  display: flex;
  gap: 8px;
  align-items: baseline;
}

.index-name {
  font-size: 14px;
  color: #666;
}

.index-value {
  font-size: 16px;
  font-weight: 500;
}

.index-change {
  font-size: 14px;
}

.index-change.up {
  color: #f44336;
}

.index-change.down {
  color: #4CAF50;
}

/* 推荐卡片 */
.recommendation-cards {
  display: flex;
  gap: 8px;
  padding: 0;
  margin: 0 8px 8px;
  width: calc(100% - 16px);
}

@media screen and (max-width: 375px) {
  .recommendation-cards {
    gap: 6px;
    padding: 0 10px;
    margin: 0 10px 6px;
  }
}

.card {
  flex: 1;
  border-radius: 8px;
  padding: 12px;
  color: #fff;
  min-height: 100px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  margin: 0;
}

.red-card {
  background-color: #f44336;
}

.blue-card {
  background-color: #3f6eea;
}

.card-tag {
  display: flex;
  align-items: center;
  gap: 5px;
  background-color: rgba(255, 255, 255, 0.2);
  padding: 2px 8px;
  border-radius: 4px;
  width: fit-content;
  font-size: 12px;
}

.card-title {
  margin-top: 8px;
  font-size: 14px;
  font-weight: 500;
}

.card-subtitle {
  font-size: 13px;
  opacity: 0.8;
}

.card-button {
  background-color: rgba(255, 255, 255, 0.3);
  width: fit-content;
  padding: 4px 12px;
  border-radius: 12px;
  margin-top: 8px;
  font-size: 13px;
}

/* 表现最佳 */
.top-performers {
  background-color: #fff;
  padding: 10px 12px;
  margin: 0 8px 8px;
  border-radius: 8px;
  width: calc(100% - 16px);
}

@media screen and (max-width: 375px) {
  .top-performers {
    padding: 10px;
    margin: 0 10px 6px;
  }
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.section-title {
  font-size: 18px;
  font-weight: 500;
}

.more-link {
  color: #999;
  font-size: 14px;
}

.fund-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.fund-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.fund-rank {
  width: 35px;
  height: 35px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.crown {
  font-size: 20px;
}

.rank-1 {
  color: #ffd700;
}

.rank-2 {
  color: #c0c0c0;
}

.rank-3 {
  color: #cd7f32;
}

.fund-return {
  width: 100px;
}

.return-rate {
  font-size: 16px;
  font-weight: bold;
  color: #f44336;
}

.return-period {
  font-size: 12px;
  color: #999;
}

.fund-info {
  flex: 1;
}

.fund-name {
  font-size: 15px;
  margin-bottom: 4px;
}

.fund-details {
  display: flex;
  gap: 10px;
}

.fund-type {
  background-color: #f0f0f0;
  color: #666;
  border-radius: 4px;
  padding: 2px 6px;
  font-size: 12px;
}

.fund-tag {
  background-color: #4CAF50;
  color: white;
  border-radius: 4px;
  padding: 2px 6px;
  font-size: 12px;
}

/* 农村金融产品 */
.agri-products {
  background-color: #fff;
  padding: 10px 12px;
  margin: 0 8px 60px;
  border-radius: 8px;
  width: calc(100% - 16px);
}

@media screen and (max-width: 375px) {
  .agri-products {
    padding: 10px;
    margin: 0 10px 60px;
  }
}

.category-tabs {
  display: flex;
  gap: 15px;
}

.category-tab {
  font-size: 14px;
  color: #666;
  position: relative;
}

.category-tab.active {
  color: #4CAF50;
  font-weight: 500;
}

.category-tab.active::after {
  content: '';
  position: absolute;
  width: 100%;
  height: 2px;
  background-color: #4CAF50;
  bottom: -5px;
  left: 0;
}

.product-list {
  margin-top: 15px;
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
}
</style> 