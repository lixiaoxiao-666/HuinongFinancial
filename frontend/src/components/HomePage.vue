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

<style scoped>
.home-page {
  /* 内边距 */
  padding: 100px;
  /* 最小高度 */
  min-height: 1000px;
  /* 背景颜色 */
  background: #f9f9f9;
}

.swiper {
  /* 宽度 */
  width: 100%;
  /* 高度 */
  height: 200px;
  /* 背景颜色 */
  background: #f0f0f0;
  /* 下边距 */
  margin-bottom: 20px;
}

.service-grid {
  /* 显示方式 */
  display: grid;
  /* 网格模板列 */
  grid-template-columns: repeat(4, 1fr);
  /* 间距 */
  gap: 8px;
  /* 内边距 */
  padding: 10px;
  /* 最大宽度 */
  max-width: 800px;
  /* 自动居中 */
  margin: 0 auto;
  /* 相对位置 */
  position: relative;
  /* 层级 */
  z-index: 1;
  /* 背景颜色 */
  background-color: #f9f9f9;
}

.service-item {
  /* 显示方式 */
  display: flex; 
  /* 弹性方向 */
  flex-direction: column;
  /* 对齐方式 */
  align-items: center;
  /* 内边距 */
  padding: 10px 5px;
  /* 背景颜色 */
  background: #fff;
  /* 圆角 */
  border-radius: 8px;
  /* 阴影 */
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  /* 过渡 */
  transition: transform 0.2s;
  /* 最小宽度 */
  min-width: 0; 
  /* 相对位置 */
  position: relative;
  /* 层级 */
  z-index: 2;
}

.service-item span {
  /* 字体大小 */
  font-size: 13px;
  /* 上边距 */
  margin-top: 5px;
  /* 宽度 */
  width: 100%;
  /* 文本对齐方式 */
  text-align: center;
  /* 文本换行 */
  white-space: nowrap;
  /* 溢出 */
  overflow: hidden;
  /* 文本溢出 */
  text-overflow: ellipsis;
  /* 内边距 */
  padding: 0 2px;
}

.service-item:hover {
  /* 过渡 */
  transform: translateY(-2px);
}

.service-icon {
  /* 颜色 */
  color: #2C5530;
  /* 下边距 */
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
  /* 文本对齐方式 */
  text-align: center;
  /* 最大宽度 */
  max-width: 800px;
  /* 自动居中 */
  margin: 0 auto;
  /* 内边距 */
  padding: 20px;
}

/* .news-lable {
  text-align: center;
  margin-bottom: 15px;
} */

.news-list {
  /* 列表样式 */
  list-style: none;
  /* 内边距 */
  padding: 0;
  /* 宽度 */
  width: 100%;
  /* 最大宽度 */
  max-width: 600px;
  /* 自动居中 */
  margin: 0 auto;
}

.news-list li {
  /* 内边距 */
  padding: 10px 15px;
  /* 下边框 */
  border-bottom: 1px solid #eee;
  /* 显示方式 */
  display: flex;
  /* 对齐方式 */
  align-items: center;
  /* 间距 */
  gap: 15px;
  /* 文本对齐方式 */
  text-align: left;
}

.news-image {
  /* 宽度 */
  width: 80px;
  /* 高度 */
  height: 60px;
  /* 对象适配 */
  object-fit: cover;
  /* 圆角 */
  border-radius: 4px;
}

.news-title {
  /* 弹性 */
  flex: 1;
  /* 对齐方式 */
  align-items: center;
}

.placeholder-banner {
  /* 宽度 */
  width: 100%;
  /* 最大宽度 */
  max-width: 800px;
  /* 自动居中 */
  margin: 0 auto;
  /* 背景颜色 */
  background: linear-gradient(45deg, #42b983, #2c3e50);
  /* 颜色 */
  color: white;
  /* 显示方式 */
  display: flex;
  /* 对齐方式 */
  align-items: center;
  /* 居中 */
  justify-content: center;
  /* 字体大小 */
  font-size: 1.2em;
}

h4 {
  /* 上边距 */
  margin-top: 0px;
}

</style>