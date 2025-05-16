<template>
  <div class="finance-overview">
    <div class="card total-assets">
      <div class="amount">
        <span class="label">总资产(元)</span>
        <span class="value">****.**</span>
        <span class="eye-icon" @click="toggleHidden">
          <Icon :icon="isHidden ? 'tabler:eye-off' : 'tabler:eye'" />
        </span>
      </div>
      <div class="detail-row">
        <div class="detail-item">
          <span class="label">昨日收益</span>
          <span class="value">+2.45</span>
        </div>
        <div class="detail-item">
          <span class="label">累计收益</span>
          <span class="value">+1,234.56</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { Icon } from '@iconify/vue'

export default {
  name: 'FinanceOverview',
  components: { Icon },
  data() {
    return {
      isHidden: false
    }
  },
  methods: {
    toggleHidden() {
      this.isHidden = !this.isHidden
    }
  }
}
</script>

<style scoped>
.finance-overview {
  padding: 15px;
}

.card {
  background: linear-gradient(135deg, #2C5530, #3a7144);
  border-radius: 12px;
  padding: 20px;
  color: white;
}

.amount {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.amount .value {
  font-size: 28px;
  font-weight: bold;
  margin: 0 10px;
}

.eye-icon {
  cursor: pointer;
}

.detail-row {
  display: flex;
  justify-content: space-between;
}

.detail-item {
  display: flex;
  flex-direction: column;
}

.detail-item .label {
  font-size: 12px;
  opacity: 0.8;
}

.detail-item .value {
  font-size: 16px;
  margin-top: 5px;
}
</style>
