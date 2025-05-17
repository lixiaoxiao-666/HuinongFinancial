<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import AppFooter from './components/footer.vue'
import ServiceIcon from './components/ServiceIcon.vue'

const router = useRouter()

// 轮播图数据
const swiperImages = ref([
  { 
    id: 1, 
    url: '/src/assets/images/banner1.png',
    title: '惠农金融服务', 
    subtitle: '助力乡村振兴'
  },
  { 
    id: 2, 
    url: '/src/assets/images/banner2.png',
    title: '农业补贴政策', 
    subtitle: '一键查询'
  }
])

// 服务入口数据
const services = ref([
  { id: 1, icon: 'loan', name: '惠农贷', path: '/loan' },
  { id: 2, icon: 'insurance', name: '惠农医保', path: '/insurance' },
  { id: 3, icon: 'subsidy', name: '农业补贴', path: '/subsidy' },
  { id: 4, icon: 'pension', name: '个人养老金', path: '/pension' },
  { id: 5, icon: 'expert', name: '农技专家', path: '/expert' },
  { id: 6, icon: 'machinery', name: '农机租赁', path: '/machinery' },
  { id: 7, icon: 'digital-currency', name: '数字人民币', path: '/digital-currency' },
  { id: 8, icon: 'welfare', name: '低保服务', path: '/welfare' }
])

// 农业资讯数据
const newsItems = ref([
  {
    id: 1, 
    title: '2023年农业补贴政策解读',
    intro: '农业农村部最新发布的农业补贴政策，种粮农民有望获得更多支持...',
    image: '/src/assets/images/news1.png',
    date: '2023-05-01',
    reads: '2456',
    url: '/news/1'
  },
  {
    id: 2, 
    title: '夏粮增产增收 农业"半年报"亮眼',
    intro: '今年夏粮再获丰收，全国夏粮总产量比上年增加60多亿斤...',
    image: '/src/assets/images/news2.png',
    date: '2023-05-15',
    reads: '3782',
    url: '/news/2'
  },
  {
    id: 3, 
    title: '数字技术助力乡村振兴',
    intro: '各地积极推动数字技术与农业生产深度融合，助力乡村产业振兴...',
    image: '/src/assets/images/news3.png',
    date: '2023-05-20',
    reads: '1845',
    url: '/news/3'
  }
])

// 惠农商城数据
const shopItems = ref([
  {
    id: 1,
    title: '农场直供新鲜鸡蛋',
    price: '28.8',
    originalPrice: '32.5',
    sales: '862',
    image: new URL('../assets/images/农产品.png', import.meta.url).href
  },
  {
    id: 2,
    title: '有机蔬菜套餐',
    price: '98.0',
    originalPrice: '128.0',
    sales: '325',
    image: new URL('../assets/images/农产品2.png', import.meta.url).href
  }
])

// 当前选中的导航栏
const activeTab = ref('home')

// 跳转函数
const navigateTo = (path: string) => {
  router.push(path)
}

// 搜索相关
const searchValue = ref('')
const onSearch = () => {
  console.log('搜索:', searchValue.value)
}
</script>

<template>
  <div class="index-page">
    <!-- 顶部导航区域 -->
    <div class="header">
      <div class="nav-bar">
        <div class="notification">
          <el-badge :value="2" class="notification-badge">
            <el-icon><Bell /></el-icon>
          </el-badge>
        </div>
      </div>
      
      <!-- 搜索框 -->
      <div class="search-container">
        <el-input 
          v-model="searchValue" 
          placeholder="搜索农业政策、金融产品" 
          prefix-icon="el-icon-search"
          class="search-input"
          @keyup.enter="onSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
    </div>
    
    <!-- 内容区域 -->
    <div class="content">
      <!-- 轮播广告区 -->
      <div class="swiper-container">
        <el-carousel height="200px" indicator-position="none" :interval="4000">
          <el-carousel-item v-for="item in swiperImages" :key="item.id">
            <div class="carousel-item">
              <img :src="item.url" alt="轮播图" class="carousel-img"/>
              <div class="carousel-text">
                <div class="carousel-subtitle">{{ item.subtitle }}</div>
                <div class="carousel-title">{{ item.title }}</div>
              </div>
            </div>
          </el-carousel-item>
        </el-carousel>
      </div>
      
      <!-- 惠农服务标题 -->
      <div class="section-title-container">
        <h2 class="section-title">惠农服务</h2>
      </div>
      
      <!-- 主要服务入口 -->
      <div class="services-container">
        <div class="services-grid">
          <div class="service-item" v-for="service in services" :key="service.id" @click="navigateTo(service.path)">
            <service-icon :type="service.icon" />
            <div class="service-name">{{ service.name }}</div>
          </div>
        </div>
      </div>
      
      <!-- 农业资讯板块 -->
      <div class="news-section">
        <div class="section-header">
          <div class="section-title">农业资讯</div>
          <div class="section-more" @click="navigateTo('/news')">更多</div>
        </div>
        
        <div class="news-list">
          <div class="news-item" v-for="news in newsItems" :key="news.id" @click="navigateTo(news.url)">
            <div class="news-content">
              <div class="news-title">{{ news.title }}</div>
              <div class="news-intro">{{ news.intro }}</div>
              <div class="news-meta">
                <span class="news-date">{{ news.date }}</span>
                <span class="news-reads">{{ news.reads }} 阅读</span>
              </div>
            </div>
            <div class="news-image">
              <img :src="news.image" alt="资讯图片" class="news-img"/>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 惠农商城板块 -->
    <div class="shop-section">
      <div class="section-header">
        <div class="section-title">惠农商城</div>
        <div class="section-more" @click="navigateTo('/shop')">更多</div>
      </div>
      
      <div class="shop-list">
        <div class="shop-item" v-for="item in shopItems" :key="item.id" @click="navigateTo(`/shop/product/${item.id}`)">
          <div class="shop-image">
            <img :src="item.image" alt="商品图片" class="shop-img"/>
          </div>
          <div class="shop-content">
            <div class="shop-title">{{ item.title }}</div>
            <div class="shop-price-info">
              <span class="shop-price">¥{{ item.price }}</span>
              <span class="shop-original-price">¥{{ item.originalPrice }}</span>
            </div>
            <div class="shop-sales">销量: {{ item.sales }}</div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 底部导航栏 -->
    <app-footer v-model:active-tab="activeTab" />
  </div>
</template>

<style scoped>
/* 整体布局 */
.index-page {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: #f5f5f5;
  font-family: Arial, sans-serif;
  position: relative;
  padding-bottom: 60px; /* 为底部导航栏留出空间 */
}

/* 头部导航区域 */
.header {
  background-color: #4CAF50; /* 主题色 */
  padding: 10px 15px;
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-bar {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  height: 30px;
}

.notification {
  color: white;
  font-size: 22px;
}

/* 搜索框 */
.search-container {
  padding: 8px 0 5px;
}

.search-input {
  border-radius: 20px;
  background-color: rgba(255, 255, 255, 0.9);
}

.search-input :deep(.el-input__inner) {
  height: 36px;
  border-radius: 20px;
  padding-left: 40px;
}

/* 内容区域 */
.content {
  flex: 1;
  padding: 0 0 10px 0;
}

/* 轮播图 */
.swiper-container {
  margin-bottom: 15px;
}

.carousel-item {
  width: 100%;
  height: 100%;
  position: relative;
}

.carousel-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.carousel-text {
  position: absolute;
  bottom: 20px;
  left: 20px;
  right: 20px;
}

.carousel-title {
  color: white;
  font-size: 18px;
  font-weight: bold;
  text-shadow: 0 1px 3px rgba(0, 0, 0, 0.5);
}

.carousel-subtitle {
  color: white;
  font-size: 24px;
  margin-bottom: 8px;
  font-weight: bold;
  text-shadow: 0 1px 3px rgba(0, 0, 0, 0.5);
}

/* 服务区域标题 */
.section-title-container {
  padding: 10px 15px 5px;
}

.section-title-container h2 {
  font-size: 16px;
  margin: 0;
  color: #333;
  font-weight: 500;
}

/* 服务入口 */
.services-container {
  background-color: white;
  padding: 15px;
  margin-bottom: 15px;
}

.services-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 15px 10px;
}

.service-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
}

.service-name {
  font-size: 13px;
  color: #333;
  margin-top: 6px;
}

/* 农业资讯 */
.news-section {
  background-color: white;
  padding: 15px;
  margin: 0 15px;
  border-radius: 10px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.section-title {
  font-size: 16px;
  font-weight: bold;
  color: #333;
}

.section-more {
  color: #999;
  font-size: 13px;
  cursor: pointer;
}

.news-list {
  display: flex;
  flex-direction: column;
}

.news-item {
  display: flex;
  cursor: pointer;
  padding-bottom: 15px;
  margin-bottom: 15px;
  border-bottom: 1px solid #eee;
}

.news-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.news-content {
  flex: 1;
  margin-right: 12px;
  display: flex;
  flex-direction: column;
}

.news-title {
  font-size: 15px;
  color: #333;
  margin-bottom: 5px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.news-intro {
  font-size: 12px;
  color: #999;
  margin-bottom: 5px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  flex: 1;
}

.news-meta {
  font-size: 12px;
  color: #999;
  display: flex;
  justify-content: space-between;
}

.news-image {
  width: 90px;
  height: 70px;
  flex-shrink: 0;
}

.news-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 3px;
}

/* 惠农商城样式 */
.shop-section {
  background-color: white;
  padding: 15px;
  margin: 15px 15px 0;
  border-radius: 10px;
}

.shop-list {
  display: flex;
  gap: 15px;
  overflow-x: auto;
  padding: 10px 0;
  -webkit-overflow-scrolling: touch;
}

.shop-item {
  flex: 0 0 auto;
  width: 160px;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  background: white;
  cursor: pointer;
}

.shop-image {
  height: 160px;
  width: 100%;
}

.shop-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.shop-content {
  padding: 10px;
}

.shop-title {
  font-size: 14px;
  color: #333;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  height: 40px;
}

.shop-price-info {
  display: flex;
  align-items: baseline;
  gap: 5px;
  margin-bottom: 5px;
}

.shop-price {
  color: #f44336;
  font-weight: bold;
  font-size: 16px;
}

.shop-original-price {
  color: #999;
  font-size: 12px;
  text-decoration: line-through;
}

.shop-sales {
  font-size: 12px;
  color: #999;
}
</style>
