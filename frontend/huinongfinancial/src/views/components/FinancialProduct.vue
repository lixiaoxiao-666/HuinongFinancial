<script setup lang="ts">
import { defineProps } from 'vue';

const props = defineProps({
  // 产品类型: 股票型、债券型等
  type: {
    type: String,
    default: '',
  },
  // 产品名称
  name: {
    type: String,
    required: true,
  },
  // 年化收益率
  annualReturn: {
    type: String,
    default: '0.00%',
  },
  // 最短投资期限
  minPeriod: {
    type: String,
    default: '无固定期限',
  },
  // 风险等级: 低、中、高
  riskLevel: {
    type: String,
    default: '中',
  },
  // 农业相关标签
  agriTag: {
    type: String,
    default: '',
  },
  // 是否为推荐产品
  isRecommended: {
    type: Boolean,
    default: false,
  }
});
</script>

<template>
  <div class="financial-product-card" :class="{ 'is-recommended': isRecommended }">
    <div class="product-header">
      <span class="product-type">{{ type }}</span>
      <span v-if="agriTag" class="product-tag">{{ agriTag }}</span>
    </div>
    
    <div class="product-name">{{ name }}</div>
    
    <div class="product-info">
      <div class="return-info">
        <div class="return-rate">{{ annualReturn }}</div>
        <div class="return-label">预期年化</div>
      </div>
      
      <div class="period-info">
        <div class="period">{{ minPeriod }}</div>
        <div class="period-label">投资期限</div>
      </div>
    </div>
    
    <div class="risk-level" :class="`risk-${riskLevel}`">
      风险等级：{{ riskLevel }}
    </div>
    
    <div class="product-button">
      立即购买
    </div>
  </div>
</template>

<style scoped>
.financial-product-card {
  width: 100%;
  background-color: #fff;
  border-radius: 10px;
  padding: 15px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  margin-bottom: 15px;
}

.is-recommended {
  border: 1px solid #4CAF50;
}

.product-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.product-type {
  background-color: #f0f0f0;
  color: #666;
  border-radius: 4px;
  padding: 2px 8px;
  font-size: 12px;
}

.product-tag {
  background-color: #4CAF50;
  color: #fff;
  border-radius: 4px;
  padding: 2px 8px;
  font-size: 12px;
}

.product-name {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 15px;
}

.product-info {
  display: flex;
  margin-bottom: 15px;
}

.return-info, .period-info {
  flex: 1;
}

.return-rate {
  color: #f44336;
  font-size: 22px;
  font-weight: bold;
}

.period {
  font-size: 16px;
  font-weight: bold;
}

.return-label, .period-label {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.risk-level {
  font-size: 12px;
  margin-bottom: 15px;
}

.risk-低 {
  color: #4CAF50;
}

.risk-中 {
  color: #FF9800;
}

.risk-高 {
  color: #f44336;
}

.product-button {
  background-color: #4CAF50;
  color: white;
  text-align: center;
  padding: 10px;
  border-radius: 20px;
  font-size: 14px;
}
</style>
