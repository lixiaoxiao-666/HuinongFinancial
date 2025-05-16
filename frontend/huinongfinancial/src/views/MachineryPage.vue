<script setup lang="ts">
import { ref } from 'vue'
import AppFooter from './components/footer.vue'

// 当前选中的导航栏
const activeTab = ref('machinery')
</script>

<template>
  <div class="machinery-container">
    <h1>农机租赁页面</h1>
    <p>此页面正在开发中...</p>
    
    <!-- 底部导航栏 -->
    <app-footer v-model:active-tab="activeTab" />
  </div>
</template>

<style scoped>
.machinery-container {
  padding: 20px;
  text-align: center;
  min-height: 100vh;
  padding-bottom: 60px; /* 为底部导航栏留出空间 */
}
</style> 