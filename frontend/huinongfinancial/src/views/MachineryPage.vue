<script setup lang="ts">
import { ref } from 'vue'
import AppFooter from './components/footer.vue'
import '../assets/icons/agri-icons.css'
import '../assets/icons/machinery-icons.css'
import MachineryCarousel from '../components/MachineryCarousel.vue'
import { carouselImages } from '../assets/images/machinery-images'

// 导入图片资源
import jiqi1 from '../assets/images/jiqi1.png';
import jiqi2 from '../assets/images/jiqi2.png';
import jiqi3 from '../assets/images/jiqi3.png';

// 当前选中的导航栏
const activeTab = ref('machinery')

// 当前订单信息
const orderInfo = ref({
  machineId: 'HN20240517',
  machineName: '马达5100联合收割机',
  machineType: '大型',
  power: '100马力',
  workWidth: '3.6米',
  fuelType: '柴油',
  rentalPrice: '600元/天',
  company: 'XX农业机械有限公司',
  pickupLocation: 'XX市XX区科学园惠农农机服务中心',
  returnLocation: 'XX市XX区科学园惠农农机服务中心',
  returnTime: '2024-06-20',
  estimatedFee: '1800',
  contactName: '',
  contactPhone: '',
  agreeTerms: false
})

// 农机类型数据
const machineryTypes = [
  { id: 1, name: '收割机', count: 24, icon: 'harvester' },
  { id: 2, name: '拖拉机', count: 18, icon: 'tractor' },
  { id: 3, name: '播种机', count: 12, icon: 'seeder' },
  { id: 4, name: '插秧机', count: 9, icon: 'transplanter' },
  { id: 5, name: '旋耕机', count: 15, icon: 'rotary-tiller' }
]

// 热门农机
interface MachineryItem {
  id: string;
  name: string;
  type: string;
  power: string;
  workWidth: string;
  fuelType: string;
  price: string;
  available: boolean;
  image: any;  // 使用any类型以接受图片资源
}

const popularMachinery: MachineryItem[] = [
  {
    id: 'HN20240517',
    name: '惠农5100联合收割机',
    type: '大型',
    power: '100马力',
    workWidth: '3.6米',
    fuelType: '柴油',
    price: '600元/天',
    available: true,
    image: jiqi1
  },
  {
    id: 'HN20240518',
    name: '惠农3088手扶式拖拉机',
    type: '中型',
    power: '45马力',
    workWidth: '1.2米',
    fuelType: '柴油',
    price: '350元/天',
    available: true,
    image: jiqi2
  },
  {
    id: 'HN20240519',
    name: '惠农2560水稻插秧机',
    type: '小型',
    power: '30马力',
    workWidth: '2.4米',
    fuelType: '汽油',
    price: '280元/天',
    available: false,
    image: jiqi3
  }
]

// 确认订单
const confirmOrder = () => {
  // 表单验证
  if (!orderInfo.value.contactName || !orderInfo.value.contactPhone) {
    alert('请填写联系人信息')
    return
  }
  
  if (!orderInfo.value.agreeTerms) {
    alert('请阅读并同意《用机须知》')
    return
  }
  
  // 提交订单逻辑
  alert('订单提交成功，即将跳转到支付页面')
}

// 选择农机类型
const selectMachineType = (typeId: number) => {
  // 切换农机类型的实现逻辑
  console.log('选择农机类型:', typeId)
}

// 租赁农机
const rentMachine = (machine: MachineryItem) => {
  if (!machine.available) {
    alert('此农机当前不可租用')
    return
  }
  
  // 更新当前订单信息
  orderInfo.value = {
    ...orderInfo.value,
    machineId: machine.id,
    machineName: machine.name,
    machineType: machine.type,
    power: machine.power,
    workWidth: machine.workWidth,
    fuelType: machine.fuelType,
    rentalPrice: machine.price
  }
  
  // 显示订单确认页
  showOrderConfirm.value = true
}

// 是否显示订单确认页
const showOrderConfirm = ref(true)
const showMachineList = ref(false)

// 返回农机列表
const goBackToList = () => {
  showOrderConfirm.value = false
  showMachineList.value = true
}

// 初始化页面视图
const initPage = () => {
  // 这里可以放置页面初始化逻辑，如获取农机数据等
}
</script>

<template>
  <div class="machinery-container">
    <!-- 订单确认页面 -->
    <div v-if="showOrderConfirm" class="order-confirm-page">
      <!-- 顶部导航 -->
      <div class="top-nav">
        <div class="nav-left" @click="goBackToList">
          <i class="nav-icon-back"></i>
        </div>
        <div class="nav-title">农机租赁</div>
        <div class="nav-right">
          <div class="nav-actions">
            <span class="share-btn"></span>
            <span class="more-btn"></span>
          </div>
        </div>
      </div>
      
      <!-- 农机轮播图 -->
      <div class="machine-carousel">
        <machinery-carousel :images="carouselImages" height="180px" :autoplay-delay="3000" />
      </div>
      
      <!-- 农机基本信息 -->
      <div class="machine-info">
        <h3 class="machine-name">{{ orderInfo.machineName }}</h3>
        <div class="machine-specs">
          <span class="spec-tag">{{ orderInfo.machineType }}</span>
          <span class="spec-tag">{{ orderInfo.power }}</span>
          <span class="spec-tag">{{ orderInfo.workWidth }}</span>
          <span class="spec-tag">{{ orderInfo.fuelType }}</span>
        </div>
        <div class="machine-description" v-if="carouselImages[0].description">
          {{ carouselImages[0].description }}
        </div>
      </div>
      
      <!-- 租赁信息 -->
      <div class="rental-info">
        <div class="info-item">
          <span class="info-icon machinery-icon machinery-icon-pickup"></span>
          <span class="info-label">取机地点</span>
          <span class="info-value">{{ orderInfo.company }}</span>
        </div>
        
        <div class="info-item">
          <span class="info-icon machinery-icon machinery-icon-return"></span>
          <span class="info-label">归还地点</span>
          <span class="info-value">{{ orderInfo.returnLocation }}</span>
        </div>
        
        <div class="info-item">
          <span class="info-icon machinery-icon machinery-icon-time"></span>
          <span class="info-label">归还时间</span>
          <span class="info-value">{{ orderInfo.returnTime }}</span>
        </div>
      </div>
      
      <!-- 费用信息 -->
      <div class="fee-section">
        <div class="fee-header">
          <span class="fee-title">预估费用</span>
          <span class="fee-amount">{{ orderInfo.estimatedFee }}元</span>
        </div>
        
        <!-- 联系人信息 -->
        <div class="contact-info">
          <div class="contact-item">
            <span class="contact-label">联系人</span>
            <input type="text" placeholder="请输入" v-model="orderInfo.contactName" class="contact-input">
          </div>
          
          <div class="contact-item">
            <span class="contact-label">联系电话</span>
            <input type="tel" placeholder="请输入" v-model="orderInfo.contactPhone" class="contact-input">
          </div>
        </div>
      </div>
      
      <!-- 服务说明 -->
      <div class="service-info">
        <div class="service-header">此次农机服务由惠农农机合作社提供</div>
        <div class="service-rules">
          <p>1.此农机价格为{{ orderInfo.rentalPrice }}/100公里，超出100公里按5元/公里收费，超出24小时按30元/小时收费。</p>
          <p>2.预估价格根据当前订单使用时间估算，仅作参考，实际费用将根据实际时间和里程数收取。</p>
        </div>
      </div>
      
      <!-- 提交订单 -->
      <div class="submit-section">
        <button class="submit-btn" @click="confirmOrder">去选机</button>
      </div>
      
      <!-- 用户协议 -->
      <div class="agreement-section">
        <label class="agreement-label">
          <input type="checkbox" v-model="orderInfo.agreeTerms">
          <span>已阅读并同意《用机须知》</span>
        </label>
      </div>
    </div>
    
    <!-- 农机列表页面 -->
    <div v-if="showMachineList" class="machine-list-page">
      <!-- 农机类型导航 -->
      <div class="machine-types">
        <div v-for="type in machineryTypes" :key="type.id" 
             class="type-item" 
             @click="selectMachineType(type.id)">
          <span class="machinery-icon" :class="`machinery-icon-${type.icon}`"></span>
          <span class="type-name">{{ type.name }}</span>
          <span class="type-count">({{ type.count }})</span>
        </div>
      </div>
      
      <!-- 农机列表 -->
      <div class="machines-list">
        <div v-for="machine in popularMachinery" :key="machine.id" 
             class="machine-card" 
             @click="rentMachine(machine)">
          <img :src="machine.image" class="machine-card-image" :alt="machine.name">
          <div class="machine-card-info">
            <h4 class="machine-card-name">{{ machine.name }}</h4>
            <div class="machine-card-specs">
              <span>{{ machine.type }}</span>
              <span>{{ machine.power }}</span>
              <span>{{ machine.workWidth }}</span>
            </div>
            <div class="machine-card-price">
              <span class="price-value">{{ machine.price }}</span>
              <span class="availability" :class="{ 'available': machine.available, 'unavailable': !machine.available }">
                {{ machine.available ? '可租用' : '已租出' }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 底部导航栏 -->
    <app-footer v-model:active-tab="activeTab" />
  </div>
</template>

<style scoped>
.machinery-container {
  padding: 0;
  text-align: left;
  min-height: 100vh;
  padding-bottom: 60px; /* 为底部导航栏留出空间 */
  background-color: #f5f5f5;
  color: #333;
  font-size: 13px;
}

@media screen and (max-width: 375px) {
  .machinery-container {
    font-size: 12px;
  }
}

/* 顶部导航栏 */
.top-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 44px;
  background-color: #fff;
  padding: 0 15px;
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-left, .nav-right {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
}

.nav-title {
  font-size: 22px;
  font-weight: 500;
}

.nav-icon-back {
  width: 24px;
  height: 24px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23333' d='M20,11H7.83l5.59-5.59L12,4l-8,8l8,8l1.41-1.41L7.83,13H20V11z'/%3E%3C/svg%3E");
  background-size: cover;
}

.nav-actions {
  display: flex;
  gap: 15px;
}

.share-btn, .more-btn {
  width: 24px;
  height: 24px;
  background-size: cover;
}

.share-btn {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23333' d='M18,16.08c-0.76,0-1.44,0.3-1.96,0.77L8.91,12.7C8.96,12.47,9,12.24,9,12s-0.04-0.47-0.09-0.7l7.05-4.11C16.5,7.69,17.21,8,18,8c1.66,0,3-1.34,3-3s-1.34-3-3-3s-3,1.34-3,3c0,0.24,0.04,0.47,0.09,0.7L8.04,9.81C7.5,9.31,6.79,9,6,9c-1.66,0-3,1.34-3,3s1.34,3,3,3c0.79,0,1.5-0.31,2.04-0.81l7.12,4.16c-0.05,0.21-0.08,0.43-0.08,0.65c0,1.61,1.31,2.92,2.92,2.92s2.92-1.31,2.92-2.92C20.92,17.39,19.61,16.08,18,16.08z'/%3E%3C/svg%3E");
}

.more-btn {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23333' d='M12,8c1.1,0,2-0.9,2-2s-0.9-2-2-2s-2,0.9-2,2S10.9,8,12,8z M12,10c-1.1,0-2,0.9-2,2s0.9,2,2,2s2-0.9,2-2S13.1,10,12,10z M12,16c-1.1,0-2,0.9-2,2s0.9,2,2,2s2-0.9,2-2S13.1,16,12,16z'/%3E%3C/svg%3E");
}

/* 农机轮播图样式 */
.machine-carousel {
  margin-bottom: 10px;
  width: 100%;
  height: 180px;
}

/* 农机基本信息 */
.machine-info {
  background-color: #fff;
  padding: 12px 15px;
  width: 100%;
  box-sizing: border-box;
}

@media screen and (max-width: 375px) {
  .machine-info {
    padding: 10px;
  }
}

@media screen and (min-width: 480px) {
  .machine-info {
    padding: 15px 20px;
    border-radius: 8px;
    margin: 0 12px;
    width: calc(100% - 24px);
  }
}

.machine-name {
  font-size: 16px;
  font-weight: bold;
  margin: 0 0 8px 0;
}

@media screen and (max-width: 375px) {
  .machine-name {
    font-size: 15px;
    margin: 0 0 6px 0;
  }
}

@media screen and (min-width: 480px) {
  .machine-name {
    font-size: 18px;
    margin: 0 0 10px 0;
  }
}

.machine-specs {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

@media screen and (max-width: 375px) {
  .machine-specs {
    gap: 6px;
  }
}

.spec-tag {
  background-color: #e8f5e9;
  color: #4CAF50;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  transition: background-color 0.2s;
}

.spec-tag:hover {
  background-color: #d5ecd8;
}

/* 农机描述 */
.machine-description {
  margin-top: 10px;
  color: #666;
  font-size: 14px;
  line-height: 1.5;
}

@media screen and (max-width: 375px) {
  .machine-description {
    font-size: 13px;
  }
}

/* 租赁信息 */
.rental-info {
  background-color: #fff;
  margin-top: 6px;
  padding: 6px 12px;
  width: 100%;
  box-sizing: border-box;
}

@media screen and (max-width: 375px) {
  .rental-info {
    margin-top: 6px;
    padding: 6px 10px;
  }
}

@media screen and (min-width: 480px) {
  .rental-info {
    padding: 15px;
    border-radius: 8px;
    margin: 10px 12px;
    width: calc(100% - 24px);
    box-shadow: 0 1px 4px rgba(0,0,0,0.05);
  }
}

.info-item {
  display: flex;
  align-items: flex-start;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
}

.info-item:last-child {
  border-bottom: none;
}

.info-icon {
  width: 22px;
  height: 22px;
  margin-right: 10px;
}

.info-label {
  width: 70px;
  color: #666;
}

.info-value {
  flex: 1;
  word-break: break-all;
}

/* 费用信息 */
.fee-section {
  background-color: #fff;
  margin-top: 6px;
  padding: 10px 12px;
  width: 100%;
  box-sizing: border-box;
}

@media screen and (max-width: 375px) {
  .fee-section {
    margin-top: 6px;
    padding: 10px;
  }
}

@media screen and (min-width: 480px) {
  .fee-section {
    padding: 15px;
    border-radius: 8px;
    margin: 10px 12px;
    width: calc(100% - 24px);
    box-shadow: 0 1px 4px rgba(0,0,0,0.05);
  }
}

.fee-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 15px;
  border-bottom: 1px solid #f0f0f0;
}

.fee-title {
  font-weight: bold;
}

.fee-amount {
  font-weight: bold;
  font-size: 18px;
  color: #f44336;
}

.contact-info {
  padding-top: 15px;
}

.contact-item {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.contact-label {
  width: 70px;
  color: #666;
}

.contact-input {
  flex: 1;
  border: none;
  background-color: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
  font-size: 14px;
}

/* 服务说明 */
.service-info {
  background-color: #fff;
  margin-top: 6px;
  padding: 10px 12px;
  border-radius: 4px;
  border: 1px solid #e0f2f1;
  width: 100%;
  box-sizing: border-box;
}

@media screen and (max-width: 375px) {
  .service-info {
    margin-top: 6px;
    padding: 10px;
  }
}

@media screen and (min-width: 480px) {
  .service-info {
    padding: 15px;
    border-radius: 8px;
    margin: 10px 12px;
    width: calc(100% - 24px);
    box-shadow: 0 1px 4px rgba(0,0,0,0.05);
    border: none;
  }
}

.service-header {
  color: #666;
  margin-bottom: 10px;
}

.service-rules {
  color: #4CAF50;
  font-size: 13px;
  line-height: 1.5;
}

/* 提交订单 */
.submit-section {
  padding: 12px;
  margin-top: 16px;
}

@media screen and (max-width: 375px) {
  .submit-section {
    padding: 10px;
    margin-top: 12px;
  }
}

@media screen and (min-width: 480px) {
  .submit-section {
    padding: 0 12px;
    margin: 20px auto;
    max-width: 400px;
  }
}

.submit-btn {
  width: 100%;
  background-color: #4CAF50;
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 12px 0;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  transition: background-color 0.2s;
}

.submit-btn:hover {
  background-color: #3e8e41;
}

/* 用户协议 */
.agreement-section {
  padding: 0 15px;
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.agreement-label {
  display: flex;
  align-items: center;
  font-size: 12px;
  color: #666;
}

.agreement-label input {
  margin-right: 5px;
}

/* 农机列表页样式 */
.machine-list-page {
  padding: 12px;
}

@media screen and (max-width: 375px) {
  .machine-list-page {
    padding: 10px;
  }
}

@media screen and (min-width: 480px) {
  .machine-list-page {
    padding: 15px;
    max-width: 900px;
    margin: 0 auto;
  }
}

.machine-types {
  display: flex;
  overflow-x: auto;
  padding-bottom: 10px;
  margin-bottom: 15px;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none; /* Firefox */
}

.machine-types::-webkit-scrollbar {
  display: none; /* Chrome, Safari, Edge */
}

.type-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px 12px;
  margin-right: 15px;
  white-space: nowrap;
  transition: transform 0.2s;
}

.type-item:active {
  transform: scale(0.95);
}

.type-item .machinery-icon {
  width: 32px;
  height: 32px;
  margin-bottom: 5px;
}

.type-name {
  font-size: 14px;
  color: #333;
  margin-bottom: 3px;
}

.type-count {
  font-size: 12px;
  color: #999;
}

.machines-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

@media screen and (min-width: 768px) {
  .machines-list {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
  }
}

.machine-card {
  background-color: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 6px rgba(0,0,0,0.1);
  transition: transform 0.2s, box-shadow 0.2s;
}

.machine-card:active {
  transform: scale(0.98);
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.machine-card-image {
  width: 100%;
  height: 100px;
  object-fit: cover;
}

@media screen and (min-width: 480px) {
  .machine-card-image {
    height: 120px;
  }
}

.machine-card-info {
  padding: 12px;
}

.machine-card-name {
  margin: 0 0 8px 0;
  font-size: 16px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.machine-card-specs {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  color: #666;
  font-size: 13px;
  margin-bottom: 10px;
}

.machine-card-specs span {
  background-color: #f5f5f5;
  padding: 2px 8px;
  border-radius: 4px;
}

.machine-card-price {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.price-value {
  color: #f44336;
  font-weight: bold;
  font-size: 16px;
}

.availability {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.available {
  background-color: #e8f5e9;
  color: #4CAF50;
}

.unavailable {
  background-color: #f5f5f5;
  color: #999;
}
</style> 