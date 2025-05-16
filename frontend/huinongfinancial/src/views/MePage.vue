<!-- vue3 element-plus -->
<script setup lang="ts">
import { ref } from 'vue'
import {
  User,
  Service,
  Help,
  Document,
  Medal,
  Setting,
  Star,
  Location,
  ChatDotRound,
  QuestionFilled,
  Ticket,
  Warning,
  InfoFilled,
  Share,
  Money,
  Wallet,
  Clock,
  Present,
  Discount,
  Bell
} from '@element-plus/icons-vue'

const userName = ref('张三')
const userPhone = ref('138****1234')
const userLevel = ref('普通会员')
</script>

<template>
  <div class="me-page">
    <!-- Header with User Info -->
    <div class="user-header">
      <div class="logo-container">
        <img src="/images/logo.png" alt="数字惠农" class="logo-image" />
      </div>
      
      <div class="user-details">
        <el-avatar :size="70" class="user-avatar">{{ userName.charAt(0) }}</el-avatar>
        <div class="user-info">
          <div class="user-name">{{ userName }}</div>
          <div class="user-phone">{{ userPhone }}</div>
          <div class="user-level">
            <el-tag size="small" type="warning">{{ userLevel }}</el-tag>
          </div>
        </div>
      </div>
      
      <el-button class="settings-btn" text>
        <el-icon><Setting /></el-icon>
      </el-button>
    </div>

    <!-- Module Sections -->
    <div class="modules-container">
      <!-- User Module (Combined with Assets) -->
      <div class="module-section">
        <div class="module-title">
          <el-icon><User /></el-icon>
          <span>个人中心</span>
        </div>
        
        <!-- User Assets -->
        <div class="asset-row">
          <el-row :gutter="20" class="mb-15">
            <el-col :span="8">
              <div class="asset-item">
                <div class="asset-value">0.00</div>
                <div class="asset-label">账户余额(元)</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="asset-item">
                <div class="asset-value">0</div>
                <div class="asset-label">优惠券(张)</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="asset-item">
                <div class="asset-value">0</div>
                <div class="asset-label">积分</div>
              </div>
            </el-col>
          </el-row>
        </div>
        
        <el-divider class="my-divider" />
        
        <!-- User Functions -->
        <el-row :gutter="0" class="module-grid">
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Document /></el-icon></div>
              <div class="grid-label">我的订单</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Star /></el-icon></div>
              <div class="grid-label">我的收藏</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Location /></el-icon></div>
              <div class="grid-label">地址管理</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Present /></el-icon></div>
              <div class="grid-label">会员特权</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Discount /></el-icon></div>
              <div class="grid-label">优惠活动</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Bell /></el-icon></div>
              <div class="grid-label">消息通知</div>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- Service & Protection Module -->
      <div class="module-section">
        <div class="module-title">
          <el-icon><Service /></el-icon>
          <span>服务与保障</span>
        </div>
        <el-row :gutter="0" class="module-grid">
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><ChatDotRound /></el-icon></div>
              <div class="grid-label">在线客服</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><QuestionFilled /></el-icon></div>
              <div class="grid-label">常见问题</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Ticket /></el-icon></div>
              <div class="grid-label">投诉反馈</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Warning /></el-icon></div>
              <div class="grid-label">风险提示</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><InfoFilled /></el-icon></div>
              <div class="grid-label">消费者权益</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Share /></el-icon></div>
              <div class="grid-label">举报中心</div>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- Business Module -->
      <div class="module-section">
        <div class="module-title">
          <el-icon><Money /></el-icon>
          <span>业务办理</span>
        </div>
        <el-row :gutter="0" class="module-grid">
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Money /></el-icon></div>
              <div class="grid-label">贷款申请</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Wallet /></el-icon></div>
              <div class="grid-label">账户管理</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Clock /></el-icon></div>
              <div class="grid-label">还款计划</div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>
  </div>
</template>

<style scoped>
.me-page {
  background-color: #f5f5f5;
  min-height: 100vh;
  padding-bottom: 20px;
}

.user-header {
  background: linear-gradient(135deg, #2e8b57, #4CAF50);
  color: white;
  padding: 0;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.logo-container {
  width: 100%;
  background-color: white;
  padding: 10px 0;
  display: flex;
  justify-content: center;
}

.logo-image {
  height: 60px;
  width: auto;
}

.user-details {
  width: 100%;
  padding: 20px;
  display: flex;
  align-items: center;
}

.user-avatar {
  background-color: #fff;
  color: #2e8b57;
  font-weight: bold;
}

.user-info {
  margin-left: 15px;
  flex: 1;
}

.user-name {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 5px;
}

.user-phone {
  font-size: 14px;
  opacity: 0.8;
  margin-bottom: 5px;
}

.settings-btn {
  color: white;
  font-size: 20px;
  position: absolute;
  top: 70px;
  right: 15px;
}

.modules-container {
  padding: 0 10px;
}

.module-section {
  background-color: white;
  border-radius: 8px;
  margin-bottom: 12px;
  padding: 15px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.08);
}

.module-title {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.module-title .el-icon {
  margin-right: 8px;
  font-size: 18px;
  color: #2e8b57;
}

.module-grid {
  text-align: center;
}

.asset-item {
  text-align: center;
}

.asset-value {
  font-size: 18px;
  font-weight: bold;
  color: #333;
}

.asset-label {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.grid-item {
  padding: 10px 0;
  cursor: pointer;
}

.grid-icon {
  display: flex;
  justify-content: center;
  margin-bottom: 8px;
}

.grid-icon .el-icon {
  font-size: 24px;
  color: #2e8b57;
  background-color: rgba(46, 139, 87, 0.08);
  padding: 10px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.grid-label {
  font-size: 12px;
  color: #666;
}

.mb-15 {
  margin-bottom: 15px;
}

.my-divider {
  margin: 15px 0;
}
</style>
