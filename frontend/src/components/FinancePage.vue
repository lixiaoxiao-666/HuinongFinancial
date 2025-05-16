<template>
  <div class="finance-page">
    <!-- 数据概览 -->
    <div class="overview-section">
      <div class="overview-item">
        <span class="label">总资产(元)</span>
        <span class="value">0.00</span>
      </div>
      <div class="overview-item">
        <span class="label">昨日收益(元)</span>
        <span class="value">0.00</span>
      </div>
    </div>

    <!-- 收益走势图 -->
    <div class="chart-container">
      <h4>收益走势</h4>
      <v-chart class="chart" :option="chartOption" autoresize />
    </div>
    
    <!-- 产品分类 -->
    <div class="category-tabs">
      <span v-for="(category, index) in categories" 
            :key="index"
            :class="['tab', { active: currentCategory === category.value }]"
            @click="selectCategory(category.value)">
        {{ category.label }}
      </span>
    </div>

    <h4>理财产品</h4>
    <div class="finance-grid">
      <div v-for="(product, index) in products" 
           :key="index" 
           class="finance-card">
        <div class="product-header">
          <div class="header-left">
            <span class="risk-level" :class="product.riskLevel">{{product.riskText}}</span>
            <span v-if="product.isRecommend" class="recommend-tag">推荐</span>
          </div>
        </div>
        <div class="product-content">
          <h5>{{product.name}}</h5>
          <div class="rate-info">
            <span class="rate">{{product.rate}}%</span>
            <span class="tag">预期年化收益率</span>
          </div>
          <div class="product-details">
            <p class="period">投资期限：{{product.period}}天</p>
            <p class="min-amount">起投金额：{{product.minAmount}}元</p>
          </div>
          <div class="progress-info">
            <div class="progress-bar">
              <div class="progress" :style="{width: product.progress + '%'}"></div>
            </div>
            <p class="progress-text">剩余额度：{{product.remainAmount}}万</p>
          </div>
          <button class="invest-btn">立即投资</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import VChart from 'vue-echarts'

use([
  CanvasRenderer,
  LineChart,
  GridComponent,
  TooltipComponent,
  LegendComponent
])

export default {
  name: 'FinancePage',
  components: {
    VChart
  },
  data() {
    return {
      chartOption: {
        tooltip: {
          trigger: 'axis'
        },
        legend: {
          data: ['七日年化收益率']
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
        },
        yAxis: {
          type: 'value',
          axisLabel: {
            formatter: '{value}%'
          }
        },
        series: [{
          name: '七日年化收益率',
          type: 'line',
          data: [4.2, 4.3, 4.5, 4.3, 4.6, 4.7, 4.8],
          itemStyle: {
            color: '#2C5530'
          }
        }]
      },
      products: [
        {
          name: '农户贷收益权',
          rate: '4.8',
          period: 90,
          minAmount: 1000,
          riskLevel: 'low',
          riskText: '低风险',
          isRecommend: true,
          progress: 85,
          remainAmount: 50
        },
        {
          name: '惠农定期理财',
          rate: '4.2',
          period: 180,
          minAmount: 500,
          riskLevel: 'low',
          riskText: '低风险',
          progress: 60,
          remainAmount: 100
        },
        {
          name: '乡村振兴债券',
          rate: '5.0',
          period: 365,
          minAmount: 2000,
          riskLevel: 'medium',
          riskText: '中风险',
          progress: 45,
          remainAmount: 200
        }
      ],
      currentCategory: 'all',
      categories: [
        { label: '全部', value: 'all' },
        { label: '稳健理财', value: 'stable' },
        { label: '定期理财', value: 'fixed' },
        { label: '基金', value: 'fund' }
      ]
    }
  },
  methods: {
    selectCategory(category) {
      this.currentCategory = category;
      // 这里可以添加筛选产品的逻辑
    }
  }
}
</script>

<style scoped>
.finance-page {
  padding: 15px;
  max-width: 800px;
  margin: 0 auto;
  min-height: 100vh;
  background: #f9f9f9;
}

.finance-page h4 {
  text-align: center;
  margin-bottom: 20px;
}

.finance-grid {
  display: flex;
  flex-direction: column;
  gap: 15px;
  max-width: 600px;
  margin: 0 auto;
}

.finance-card {
  background: #fff;
  border-radius: 8px;
  padding: 15px;
  display: flex;
  flex-direction: column;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.product-header {
  margin-bottom: 15px;
}

.header-left {
  display: flex;
  gap: 8px;
}

.product-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.product-content h5 {
  margin: 0;
  font-size: 18px;
}

.rate-info {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.product-details {
  display: flex;
  justify-content: space-between;
}

.period, .min-amount {
  margin: 0;
}

.progress-info {
  margin-top: 10px;
}

.rate-info {
  text-align: center;
  min-width: 120px;
}

.rate {
  font-size: 24px;
  color: #f56c6c;
  font-weight: bold;
}

.tag {
  font-size: 12px;
  color: #666;
  display: block;
}

.product-info {
  flex: 1;
  padding: 0 20px;
}

.product-info h5 {
  margin: 0 0 10px 0;
  font-size: 16px;
}

.period, .min-amount {
  margin: 5px 0;
  font-size: 14px;
  color: #666;
}

.invest-btn {
  background: #2C5530;
  color: white;
  border: none;
  padding: 8px 20px;
  border-radius: 20px;
  cursor: pointer;
}

.chart-container {
  background: #fff;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.chart {
  height: 300px;
  width: 100%;
}

.overview-section {
  display: flex;
  justify-content: space-around;
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  margin-bottom: 15px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.overview-item {
  text-align: center;
}

.overview-item .label {
  font-size: 14px;
  color: #666;
}

.overview-item .value {
  display: block;
  font-size: 24px;
  font-weight: bold;
  color: #2C5530;
}

.category-tabs {
  display: flex;
  justify-content: space-between;
  margin: 0 0 15px;
  padding: 0 15px;
  overflow: hidden;
}

.tab {
  flex: 1;
  text-align: center;
  padding: 12px 0;
  margin: 0 8px;
  border-radius: 20px;
  background: #f0f0f0;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s ease;
}

.tab.active {
  background: #2C5530;
  color: white;
  transform: scale(1.05);
}

.product-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.risk-level {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.risk-level.low {
  background: #e8f5e9;
  color: #2C5530;
}

.risk-level.medium {
  background: #fff3e0;
  color: #f57c00;
}

.recommend-tag {
  background: #ff6b6b;
  color: white;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.progress-bar {
  height: 4px;
  background: #f0f0f0;
  border-radius: 2px;
  margin: 10px 0;
}

.progress {
  height: 100%;
  background: #2C5530;
  border-radius: 2px;
}

.progress-text {
  font-size: 12px;
  color: #666;
  margin: 0;
}
</style>
