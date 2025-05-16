<template>
  <div class="home-page">
    <!-- 轮播图 -->
    <div class="swiper">
      <Carousel />
    </div>
    
    <!-- 服务模块 -->
    <div class="service-grid">
      <div v-for="(service, index) in services" :key="index" class="service-item">
        <Icon :icon="service.icon" :width="36" class="service-icon" />
        <span>{{ service.name }}</span>
      </div>
    </div>

    <!-- 天气模块
    <div class="weather-section">
      <h3>今日天气</h3>
      <div class="weather-info">
        <span>晴 26°C</span>
        <span>适合农作物生长</span>
      </div>
    </div> -->

    <!-- 资讯模块 -->
    <div class="news-section">
      <h4 class="news-lable">农业资讯</h4>
      <ul class="news-list">
        <li v-for="(news, index) in newsItems" :key="index" class="news-item">
          <img :src="news.image" :alt="news.title" class="news-image">
          <div class="news-title">{{ news.title }}</div>
        </li>
      </ul>
    </div>

    <ShopPage />
  </div>
</template>

<script>
import { Icon } from '@iconify/vue'
import Carousel from './Carousel.vue'
import ShopPage from './ShopPage.vue'

export default {
  name: 'HomePage',
  components: {
    Carousel,
    Icon,
    ShopPage
  },
  data() {
    return {
      services: [
        { name: '惠农贷', icon: 'tabler:cash-banknote' },
        { name: '惠农医保', icon: 'tabler:heart-plus' },
        { name: '农业补贴', icon: 'tabler:plant-2' },
        { name: '个人养老金', icon: 'tabler:Money' },
        { name: '农技专家', icon: 'tabler:building-factory-2' },
        { name: '农机租赁', icon: 'tabler:tractor' },
        { name: '数字人民币', icon: 'tabler:coin-yuan' },
        { name: '低保', icon: 'tabler:shield-check' }
      ],
      newsItems: [
        { 
          title: '知农爱农 | 金融支持绿色农业发展的模式、难点及对策建议',
          image: require('@/assets/1.webp')
        },
        { 
          title: '新思想引领新征程丨加快建设数字乡村 推进农业农村现代化',
          image: require('@/assets/image.png')
        },
        { 
          title: '金融"贷"动 现代农业发展添活力',
          image: require('@/assets/3.webp')
        }
      ]
    }
  }
}
</script>

<style scoped>
.home-page {
  padding: 10px;
  min-height: 100vh;
  background: #f9f9f9;
}

.swiper {
  width: 100%;
  height: 200px;
  background: #f0f0f0;
  margin-bottom: 20px;
}

.service-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  padding: 10px;
  max-width: 800px;
  margin: 0 auto;
  position: relative;
  z-index: 1;
  background-color: #f9f9f9;
}

.service-item {
  display: flex; 
  flex-direction: column;
  align-items: center;
  padding: 10px 5px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  transition: transform 0.2s;
  min-width: 0; /* 防止内容溢出 */
  position: relative;
  z-index: 2;
}

.service-item span {
  font-size: 13px;
  margin-top: 5px;
  width: 100%;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  padding: 0 2px;
}

.service-item:hover {
  transform: translateY(-2px);
}

.service-icon {
  color: #2C5530;
  margin-bottom: 5px;
}

/* .weather-section, .news-section {
  margin: 0;
  padding: 5px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
} */

.news-section {
  text-align: center;
  max-width: 800px;
  margin: 0 auto;
  /* padding: 20px; */
}

/* .news-lable {
  text-align: center;
  margin-bottom: 15px;
} */

.news-list {
  list-style: none;
  padding: 0;
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
}

.news-list li {
  padding: 10px 15px;
  border-bottom: 1px solid #eee;
  display: flex;
  align-items: center;
  gap: 15px;
  text-align: left;
}

.news-image {
  width: 80px;
  height: 60px;
  object-fit: cover;
  border-radius: 4px;
}

.news-title {
  flex: 1;
  align-items: center;
}

.placeholder-banner {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  background: linear-gradient(45deg, #42b983, #2c3e50);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2em;
}

h4 {
margin: 0px;
}

</style>