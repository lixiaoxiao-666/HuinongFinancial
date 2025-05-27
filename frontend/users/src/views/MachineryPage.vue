<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import AppFooter from './components/footer.vue'
import { useUserStore } from '../stores/user'
import '../assets/icons/agri-icons.css'
import '../assets/icons/machinery-icons.css'
import MachineryCarousel from '../components/MachineryCarousel.vue'
import { carouselImages } from '../assets/images/machinery-images'

// 导入图片资源
import jiqi1 from '../assets/images/jiqi1.png';
import jiqi2 from '../assets/images/jiqi2.png';
import jiqi3 from '../assets/images/jiqi3.png';

const router = useRouter()
const userStore = useUserStore()
const activeTab = ref('machinery')
const loading = ref(false)
const refreshing = ref(false)

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
  { id: 1, name: '收割机', count: 5, icon: 'harvester' },
  { id: 2, name: '拖拉机', count: 2, icon: 'tractor' },
  { id: 3, name: '播种机', count: 2, icon: 'seeder' },
  { id: 4, name: '插秧机', count: 2, icon: 'transplanter' },
  { id: 5, name: '旋耕机', count: 2, icon: 'rotary-tiller' }
]

// 农机设备接口
interface MachineryItem {
  id: string
  name: string
  description: string
  category: string
  brand: string
  model: string
  power: string
  workWidth: string
  fuelType: string
  dailyPrice: number
  hourlyPrice: number
  available: boolean
  location: string
  contactPhone: string
  images: string[]
  specifications: Record<string, string>
}

// 农机设备列表
const machineryList = ref<MachineryItem[]>([])
const selectedCategory = ref('')

// 农机分类
const categories = [
  { value: '', label: '全部设备', icon: 'all' },
  { value: '收割机', label: '收割机', icon: 'harvester' },
  { value: '拖拉机', label: '拖拉机', icon: 'tractor' },
  { value: '播种机', label: '播种机', icon: 'seeder' },
  { value: '插秧机', label: '插秧机', icon: 'transplanter' },
  { value: '旋耕机', label: '旋耕机', icon: 'rotary' },
]

// 筛选后的农机列表
const filteredMachinery = computed(() => {
  if (!selectedCategory.value) {
    return machineryList.value
  }
  return machineryList.value.filter(item => item.category === selectedCategory.value)
})

// 格式化价格
const formatPrice = (price: number) => {
  return price.toLocaleString()
}

// 模拟农机数据
const mockMachineryData: MachineryItem[] = [
  {
    id: 'machinery_001',
    name: '雷沃谷神GE50收割机',
    description: '大型联合收割机，适合水稻、小麦收割，效率高，操作简便',
    category: '收割机',
    brand: '雷沃重工',
    model: 'GE50',
    power: '150马力',
    workWidth: '4.2米',
    fuelType: '柴油',
    dailyPrice: 800,
    hourlyPrice: 35,
    available: true,
    location: '泰安市农机服务中心',
    contactPhone: '18888888888',
    images: ['./src/assets/images/jiqi1.png'],
    specifications: {
      '发动机功率': '150马力',
      '作业幅宽': '4.2米',
      '燃料类型': '柴油',
      '整机重量': '8.5吨',
      '作业效率': '20-30亩/小时'
    }
  },
  {
    id: 'machinery_002',
    name: '东方红1104拖拉机',
    description: '中大型轮式拖拉机，动力强劲，适合耕地、播种等多种作业',
    category: '拖拉机',
    brand: '一拖集团',
    model: '1104',
    power: '110马力',
    workWidth: '2.5米',
    fuelType: '柴油',
    dailyPrice: 500,
    hourlyPrice: 25,
    available: true,
    location: '洛阳市农机租赁点',
    contactPhone: '18777777777',
    images: ['./src/assets/images/jiqi2.png'],
    specifications: {
      '发动机功率': '110马力',
      '驱动方式': '四轮驱动',
      '燃料类型': '柴油',
      '整机重量': '4.8吨',
      '适用作业': '耕地、播种、施肥'
    }
  },
  {
    id: 'machinery_003',
    name: '常发CF505插秧机',
    description: '水稻插秧专用机械，插秧速度快，秧苗成活率高',
    category: '插秧机',
    brand: '常发农机',
    model: 'CF505',
    power: '35马力',
    workWidth: '1.8米',
    fuelType: '汽油',
    dailyPrice: 400,
    hourlyPrice: 20,
    available: false,
    location: '常州市农机服务站',
    contactPhone: '18666666666',
    images: ['./src/assets/images/jiqi3.png'],
    specifications: {
      '发动机功率': '35马力',
      '插秧行数': '6行',
      '燃料类型': '汽油',
      '整机重量': '850公斤',
      '作业效率': '8-12亩/小时'
    }
  },
  {
    id: 'machinery_004',
    name: '德邦大为2BFX-12播种机',
    description: '小麦玉米通用播种机，播种均匀，深度可调',
    category: '播种机',
    brand: '德邦大为',
    model: '2BFX-12',
    power: '配套60-80马力拖拉机',
    workWidth: '3.0米',
    fuelType: '拖拉机带动',
    dailyPrice: 300,
    hourlyPrice: 15,
    available: true,
    location: '德州市农机合作社',
    contactPhone: '18555555555',
    images: ['./src/assets/images/jiqi1.png'],
    specifications: {
      '播种行数': '12行',
      '行距': '250mm',
      '作业幅宽': '3.0米',
      '整机重量': '1.2吨',
      '适用作物': '小麦、玉米、大豆'
    }
  },
  {
    id: 'machinery_005',
    name: '华德1GQN-200旋耕机',
    description: '多功能旋耕机，耕作深度可调，适合不同土壤类型',
    category: '旋耕机',
    brand: '华德农机',
    model: '1GQN-200',
    power: '配套40-60马力拖拉机',
    workWidth: '2.0米',
    fuelType: '拖拉机带动',
    dailyPrice: 250,
    hourlyPrice: 12,
    available: true,
    location: '潍坊市农机服务中心',
    contactPhone: '18444444444',
    images: ['./src/assets/images/jiqi2.png'],
    specifications: {
      '耕作深度': '15-25cm',
      '作业幅宽': '2.0米',
      '刀片数量': '32片',
      '整机重量': '680公斤',
      '适用土壤': '粘土、沙土、壤土'
    }
  },
  {
    id: 'machinery_006',
    name: '久保田DC70收割机',
    description: '进口联合收割机，性能稳定，收割质量高',
    category: '收割机',
    brand: '久保田',
    model: 'DC70',
    power: '70马力',
    workWidth: '2.8米',
    fuelType: '柴油',
    dailyPrice: 700,
    hourlyPrice: 30,
    available: true,
    location: '济南市农机租赁站',
    contactPhone: '18333333333',
    images: ['./src/assets/images/jiqi3.png'],
    specifications: {
      '发动机功率': '70马力',
      '作业幅宽': '2.8米',
      '燃料类型': '柴油',
      '整机重量': '3.8吨',
      '作业效率': '15-20亩/小时'
    }
  },
  {
    id: 'machinery_007',
    name: '约翰迪尔6B系列拖拉机',
    description: '国际知名品牌，高效可靠，适合大面积作业',
    category: '拖拉机',
    brand: '约翰迪尔',
    model: '6B-1454',
    power: '145马力',
    workWidth: '3.2米',
    fuelType: '柴油',
    dailyPrice: 650,
    hourlyPrice: 30,
    available: true,
    location: '青岛市农机中心',
    contactPhone: '18222222222',
    images: ['./src/assets/images/jiqi1.png'],
    specifications: {
      '发动机功率': '145马力',
      '驱动方式': '四轮驱动',
      '燃料类型': '柴油',
      '整机重量': '6.2吨',
      '适用作业': '深耕、播种、收获'
    }
  },
  {
    id: 'machinery_008',
    name: '沃得锐龙4LZ-8A收割机',
    description: '国产优质收割机，适合多种作物收割，性价比高',
    category: '收割机',
    brand: '沃得农机',
    model: '4LZ-8A',
    power: '125马力',
    workWidth: '3.8米',
    fuelType: '柴油',
    dailyPrice: 720,
    hourlyPrice: 32,
    available: true,
    location: '南京市农机服务站',
    contactPhone: '18111111111',
    images: ['./src/assets/images/jiqi2.png'],
    specifications: {
      '发动机功率': '125马力',
      '作业幅宽': '3.8米',
      '燃料类型': '柴油',
      '整机重量': '7.2吨',
      '作业效率': '18-25亩/小时'
    }
  },
  {
    id: 'machinery_009',
    name: '洋马VP8G插秧机',
    description: '日本进口高端插秧机，插秧精度高，作业效率佳',
    category: '插秧机',
    brand: '洋马农机',
    model: 'VP8G',
    power: '40马力',
    workWidth: '2.4米',
    fuelType: '汽油',
    dailyPrice: 480,
    hourlyPrice: 22,
    available: true,
    location: '杭州市农机租赁中心',
    contactPhone: '18999999999',
    images: ['./src/assets/images/jiqi3.png'],
    specifications: {
      '发动机功率': '40马力',
      '插秧行数': '8行',
      '燃料类型': '汽油',
      '整机重量': '1.1吨',
      '作业效率': '12-16亩/小时'
    }
  },
  {
    id: 'machinery_010',
    name: '大华宝来2BQX-200播种机',
    description: '气力式精密播种机，适合玉米、大豆等作物播种',
    category: '播种机',
    brand: '大华宝来',
    model: '2BQX-200',
    power: '配套80-120马力拖拉机',
    workWidth: '4.0米',
    fuelType: '拖拉机带动',
    dailyPrice: 380,
    hourlyPrice: 18,
    available: true,
    location: '哈尔滨市农机合作社',
    contactPhone: '18188888888',
    images: ['./src/assets/images/jiqi1.png'],
    specifications: {
      '播种行数': '16行',
      '行距': '250mm',
      '作业幅宽': '4.0米',
      '整机重量': '1.8吨',
      '适用作物': '玉米、大豆、花生'
    }
  },
  {
    id: 'machinery_011',
    name: '金达威1GKN-300旋耕机',
    description: '重型旋耕机，适合坚硬土壤耕作，耐用性强',
    category: '旋耕机',
    brand: '金达威',
    model: '1GKN-300',
    power: '配套60-100马力拖拉机',
    workWidth: '3.0米',
    fuelType: '拖拉机带动',
    dailyPrice: 320,
    hourlyPrice: 16,
    available: true,
    location: '郑州市农机服务中心',
    contactPhone: '18177777777',
    images: ['./src/assets/images/jiqi2.png'],
    specifications: {
      '耕作深度': '18-30cm',
      '作业幅宽': '3.0米',
      '刀片数量': '48片',
      '整机重量': '950公斤',
      '适用土壤': '重粘土、硬质土壤'
    }
  },
  {
    id: 'machinery_012',
    name: '中联重科PL40水稻收割机',
    description: '专业水稻收割机，湿地作业能力强，收割损失率低',
    category: '收割机',
    brand: '中联重科',
    model: 'PL40',
    power: '90马力',
    workWidth: '2.2米',
    fuelType: '柴油',
    dailyPrice: 580,
    hourlyPrice: 26,
    available: false,
    location: '长沙市农机服务站',
    contactPhone: '18166666666',
    images: ['./src/assets/images/jiqi3.png'],
    specifications: {
      '发动机功率': '90马力',
      '作业幅宽': '2.2米',
      '燃料类型': '柴油',
      '整机重量': '4.5吨',
      '作业效率': '12-18亩/小时'
    }
  }
]

// 获取农机列表
const getMachineryList = async () => {
  loading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    machineryList.value = mockMachineryData
  } catch (error) {
    console.error('获取农机列表失败:', error)
    ElMessage.error('获取农机列表失败')
  } finally {
    loading.value = false
  }
}

// 刷新农机列表
const onRefresh = async () => {
  refreshing.value = true
  try {
    await getMachineryList()
    ElMessage.success('刷新成功')
  } finally {
    refreshing.value = false
  }
}

// 查看农机详情
const viewMachineryDetail = (machinery: MachineryItem) => {
  // 跳转到农机详情页
  router.push(`/machinery/detail/${machinery.id}`)
}

// 租赁农机
const rentMachinery = (machinery: MachineryItem) => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  if (!machinery.available) {
    ElMessage.warning('该设备暂不可租用')
    return
  }
  
  // 跳转到租赁申请页
  router.push(`/machinery/rent/${machinery.id}`)
}

// 联系租赁方
const contactRenter = (machinery: MachineryItem) => {
  const phone = machinery.contactPhone
  if (confirm(`确定要拨打 ${phone} 吗？`)) {
    window.location.href = `tel:${phone}`
  }
}

// 选择分类
const selectCategory = (category: string) => {
  selectedCategory.value = category
}

onMounted(() => {
  getMachineryList()
})
</script>

<template>
  <div class="machinery-container">
    <!-- 顶部导航 -->
    <div class="top-header">
      <div class="header-content">
        <div class="back-btn" @click="router.go(-1)">
          <i class="el-icon-arrow-left"></i>
        </div>
        <h1 class="header-title">农机租赁</h1>
        <div class="header-actions">
          <i class="el-icon-refresh" @click="onRefresh" :class="{ rotating: refreshing }"></i>
        </div>
      </div>
    </div>

    <!-- 我的租赁入口 -->
    <div class="my-rental-entry" v-if="userStore.isLoggedIn">
      <div class="entry-content" @click="router.push('/machinery/my-rentals')">
        <div class="entry-left">
          <div class="entry-title">我的租赁</div>
          <div class="entry-subtitle">查看租赁记录进度</div>
        </div>
        <div class="entry-right">
          <div class="entry-action">查看</div>
          <i class="el-icon-arrow-right"></i>
        </div>
      </div>
    </div>

    <!-- 分类筛选 -->
    <div class="category-section">
      <div class="category-title">设备分类</div>
      <div class="category-tabs">
        <div
          v-for="category in categories"
          :key="category.value"
          class="category-tab"
          :class="{ active: selectedCategory === category.value }"
          @click="selectCategory(category.value)"
        >
          <div class="category-icon" :class="`icon-${category.icon}`"></div>
          <div class="category-name">{{ category.label }}</div>
        </div>
      </div>
    </div>

    <!-- 农机列表 -->
    <div class="machinery-section">
      <div class="section-header">
        <div class="section-title">可租设备</div>
        <div class="section-count">{{ filteredMachinery.length }}台设备</div>
      </div>

      <div v-loading="loading" class="machinery-list">
        <div
          v-for="machinery in filteredMachinery"
          :key="machinery.id"
          class="machinery-card"
        >
          <div class="card-header">
            <div class="machinery-info">
              <h3 class="machinery-name">{{ machinery.name }}</h3>
              <div class="machinery-meta">
                <span class="category-tag">{{ machinery.category }}</span>
                <span class="brand-tag">{{ machinery.brand }}</span>
                <span class="availability-tag" :class="{ available: machinery.available, unavailable: !machinery.available }">
                  {{ machinery.available ? '可租用' : '已租出' }}
                </span>
              </div>
              <p class="machinery-desc">{{ machinery.description }}</p>
            </div>
          </div>

          <div class="card-content">
            <div class="spec-grid">
              <div class="spec-item">
                <span class="spec-label">动力</span>
                <span class="spec-value">{{ machinery.power }}</span>
              </div>
              <div class="spec-item">
                <span class="spec-label">幅宽</span>
                <span class="spec-value">{{ machinery.workWidth }}</span>
              </div>
              <div class="spec-item">
                <span class="spec-label">燃料</span>
                <span class="spec-value">{{ machinery.fuelType }}</span>
              </div>
              <div class="spec-item">
                <span class="spec-label">位置</span>
                <span class="spec-value">{{ machinery.location }}</span>
              </div>
            </div>

            <div class="price-section">
              <div class="price-info">
                <div class="daily-price">
                  <span class="price-label">日租</span>
                  <span class="price-value">¥{{ formatPrice(machinery.dailyPrice) }}</span>
                  <span class="price-unit">/天</span>
                </div>
                <div class="hourly-price">
                  <span class="price-label">时租</span>
                  <span class="price-value">¥{{ formatPrice(machinery.hourlyPrice) }}</span>
                  <span class="price-unit">/小时</span>
                </div>
              </div>
            </div>
          </div>

          <div class="card-actions">
            <button class="action-btn detail-btn" @click="viewMachineryDetail(machinery)">
              查看详情
            </button>
            <button
              class="action-btn contact-btn"
              @click="contactRenter(machinery)"
            >
              联系租赁
            </button>
            <button
              class="action-btn rent-btn"
              :class="{ disabled: !machinery.available }"
              @click="rentMachinery(machinery)"
              :disabled="!machinery.available"
            >
              {{ machinery.available ? '立即租赁' : '已租出' }}
            </button>
          </div>
        </div>
      </div>

      <div v-if="filteredMachinery.length === 0 && !loading" class="empty-state">
        <div class="empty-icon">🚜</div>
        <div class="empty-text">暂无可租用设备</div>
        <div class="empty-desc">请稍后再试或联系客服</div>
      </div>
    </div>

    <!-- 底部导航栏 -->
    <app-footer v-model:active-tab="activeTab" />
  </div>
</template>

<style scoped>
.machinery-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding-bottom: 80px;
}

/* 顶部导航 */
.top-header {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: white;
  padding: 20px 0 15px;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.back-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255,255,255,0.2);
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.3s ease;
}

.back-btn:hover {
  background: rgba(255,255,255,0.3);
  transform: scale(1.1);
}

.header-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
}

.header-actions {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.rotating {
  animation: rotate 1s linear infinite;
}

@keyframes rotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 我的租赁入口 */
.my-rental-entry {
  margin: 20px 20px 0;
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
}

.my-rental-entry:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(39, 174, 96, 0.3);
}

.entry-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px;
  color: white;
}

.entry-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 4px;
}

.entry-subtitle {
  font-size: 14px;
  opacity: 0.8;
}

.entry-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.entry-action {
  font-size: 16px;
  font-weight: 500;
}

/* 分类筛选 */
.category-section {
  margin: 20px 20px 0;
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.category-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 15px;
}

.category-tabs {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
  gap: 10px;
}

.category-tab {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 15px 10px;
  border-radius: 8px;
  background: #f8f9fa;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.category-tab:hover {
  background: #e9ecef;
  transform: translateY(-2px);
}

.category-tab.active {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: white;
  border-color: #27ae60;
  box-shadow: 0 4px 15px rgba(39, 174, 96, 0.3);
}

.category-icon {
  width: 24px;
  height: 24px;
  margin-bottom: 8px;
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
}

.category-name {
  font-size: 14px;
  font-weight: 500;
  text-align: center;
}

/* 农机列表 */
.machinery-section {
  margin: 20px 20px 0;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.section-count {
  font-size: 14px;
  color: #666;
}

.machinery-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.machinery-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  transition: all 0.3s ease;
}

.machinery-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0,0,0,0.15);
}

.card-header {
  padding: 20px 20px 0;
}

.machinery-name {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0 0 10px;
}

.machinery-meta {
  display: flex;
  gap: 8px;
  margin-bottom: 10px;
  flex-wrap: wrap;
}

.category-tag, .brand-tag {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.category-tag {
  background: #e3f2fd;
  color: #1976d2;
}

.brand-tag {
  background: #f3e5f5;
  color: #7b1fa2;
}

.availability-tag {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.availability-tag.available {
  background: #e8f5e8;
  color: #2e7d32;
}

.availability-tag.unavailable {
  background: #ffebee;
  color: #c62828;
}

.machinery-desc {
  font-size: 14px;
  color: #666;
  line-height: 1.5;
  margin: 0;
}

.card-content {
  padding: 15px 20px;
}

.spec-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
  gap: 10px;
  margin-bottom: 15px;
}

.spec-item {
  text-align: center;
}

.spec-label {
  display: block;
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
}

.spec-value {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.price-section {
  border-top: 1px solid #eee;
  padding-top: 15px;
}

.price-info {
  display: flex;
  justify-content: space-around;
  text-align: center;
}

.daily-price, .hourly-price {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.price-label {
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
}

.price-value {
  font-size: 18px;
  font-weight: 600;
  color: #27ae60;
}

.price-unit {
  font-size: 12px;
  color: #999;
}

.card-actions {
  display: flex;
  gap: 10px;
  padding: 15px 20px 20px;
}

.action-btn {
  flex: 1;
  padding: 10px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.detail-btn {
  background: #f8f9fa;
  color: #333;
}

.detail-btn:hover {
  background: #e9ecef;
}

.contact-btn {
  background: #fff3e0;
  color: #f57c00;
}

.contact-btn:hover {
  background: #ffe0b2;
}

.rent-btn {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: white;
}

.rent-btn:hover:not(.disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(39, 174, 96, 0.3);
}

.rent-btn.disabled {
  background: #e0e0e0;
  color: #999;
  cursor: not-allowed;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: #666;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 8px;
}

.empty-desc {
  font-size: 14px;
  opacity: 0.8;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .machinery-container {
    padding-bottom: 70px;
  }
  
  .header-content {
    padding: 0 15px;
  }
  
  .my-rental-entry,
  .category-section,
  .machinery-section {
    margin-left: 15px;
    margin-right: 15px;
  }
  
  .category-tabs {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .spec-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .card-actions {
    flex-direction: column;
  }
  
  .action-btn {
    flex: none;
  }
}
</style> 