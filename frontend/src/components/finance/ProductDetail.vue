<template>
  <div class="product-detail">
    <div class="header">
      <Icon icon="tabler:arrow-left" @click="$router.back()" />
      <h3>产品详情</h3>
      <div class="tag-group">
        <span class="tag risk-tag">低风险</span>
        <span class="tag recommend-tag">推荐</span>
      </div>
    </div>
    
    <div class="rate-card">
      <h2 class="product-title">惠农稳健理财产品</h2>
      <p class="product-desc">为农户提供的稳健型理财方案</p>
      <div class="rate">
        <span class="number">4.8</span>
        <span class="unit">%</span>
      </div>
      <div class="label">预期年化收益率</div>
    </div>

    <div class="info-card">
      <div class="info-item">
        <span class="label">投资期限</span>
        <span class="value">90天</span>
      </div>
      <div class="info-item">
        <span class="label">起投金额</span>
        <span class="value">1000元</span>
      </div>
    </div>

    <div class="product-details">
      <h4>产品说明</h4>
      <div class="detail-item">
        <h5>收益说明</h5>
        <p>本产品为固定收益类产品，按月付息，到期还本</p>
      </div>
      <div class="detail-item">
        <h5>投资建议</h5>
        <p>适合稳健型投资者，追求稳定收益</p>
      </div>
    </div>

    <div class="action-bar">
      <div class="agreement">
        <input type="checkbox" v-model="agreed">
        <span>我已阅读并同意《投资协议》</span>
      </div>
      <button class="invest-btn" :disabled="!agreed">立即投资</button>
    </div>
  </div>
</template>

<script>
import { Icon } from '@iconify/vue'

export default {
  name: 'ProductDetail',
  components: { Icon },
  data() {
    return {
      agreed: false
    }
  }
}
</script>

<style scoped>
.product-detail {
  min-height: 100vh;
  background: #f9f9f9;
}

.header {
  position: relative;
  display: flex;
  align-items: center;
  padding: 15px;
  background: #fff;
  height: 56px;
}

.header h3 {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  margin: 0;
}

.tag-group {
  position: absolute;
  right: 15px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  gap: 8px;
  z-index: 1;
}

.tag {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.risk-tag {
  color: #2C5530;
  background: #e8f5e9;
}

.recommend-tag {
  color: #f57c00;
  background: #fff3e0;
}

.rate-card {
  text-align: center;
  padding: 30px;
  background: #2C5530;
  color: white;
}

.product-title {
  font-size: 24px;
  margin-bottom: 8px;
}

.product-desc {
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: 20px;
}

.rate .number {
  font-size: 48px;
  font-weight: bold;
}

.rate .unit {
  font-size: 24px;
}

.info-card {
  margin: 15px;
  padding: 20px;
  background: white;
  border-radius: 8px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
}

.low-risk {
  color: #2C5530;
  background: #e8f5e9;
  padding: 2px 8px;
  border-radius: 4px;
}

.product-details {
  margin: 15px;
  padding: 20px;
  background: white;
  border-radius: 8px;
}

.detail-item {
  margin: 15px 0;
}

.detail-item h5 {
  color: #666;
  margin-bottom: 8px;
}

.action-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 15px;
  background: white;
  box-shadow: 0 -2px 10px rgba(0,0,0,0.05);
}

.agreement {
  margin-bottom: 10px;
  font-size: 14px;
  color: #666;
}

.invest-btn {
  width: 100%;
  padding: 12px;
  background: #2C5530;
  color: white;
  border: none;
  border-radius: 24px;
  font-size: 16px;
}

.invest-btn:disabled {
  background: #ccc;
}
</style>
