/* eslint-disable */
<template>
  <div class="carousel">
    <div class="carousel-container" :style="{ transform: `translateX(-${currentIndex * 100}%)` }">
      <img :src="require('@/assets/carousel1.png')" alt="轮播图1">
      <img :src="require('@/assets/carousel2.png')" alt="轮播图2">
      <img :src="require('@/assets/carousel3.png')" alt="轮播图3">
    </div>
    <div class="indicators">
      <span 
        v-for="(_, index) in 3" 
        :key="index"
        :class="{ active: currentIndex === index }"
        @click="setCurrentIndex(index)"
      ></span>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AppCarousel',
  data() {
    return {
      currentIndex: 0,
      timer: null
    }
  },
  mounted() {
    this.startAutoPlay()
  },
  beforeUnmount() {
    this.stopAutoPlay()
  },
  methods: {
    startAutoPlay() {
      this.timer = setInterval(() => {
        this.currentIndex = (this.currentIndex + 1) % 3
      }, 3000)
    },
    stopAutoPlay() {
      if (this.timer) {
        clearInterval(this.timer)
      }
    },
    setCurrentIndex(index) {
      this.currentIndex = index
    }
  }
}
</script>

<style scoped>
.carousel {
  width: 100%;
  overflow: hidden;
  position: relative;
}

.carousel-container {
  display: flex;
  transition: transform 0.5s ease;
}

.carousel-container img {
  width: 100%;
  flex-shrink: 0;
}

.indicators {
  position: absolute;
  bottom: 10px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 8px;
}

.indicators span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.5);
  cursor: pointer;
}

.indicators span.active {
  background-color: #fff;
}
</style>