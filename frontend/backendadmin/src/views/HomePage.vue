<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 统计数据
const stats = ref({
  totalLoans: 425,
  pendingLoans: 135,
  totalRentals: 252,
  pendingRentals: 175,
  totalMachinery: 489,
  availableMachinery: 185
})

// 最近贷款申请
const recentLoans = ref([
  {
    id: 1,
    applicantName: '张三',
    propertyType: '种植补贴',
    amount: 50000,
    status: '待审核',
    date: '2024-04-10',
    contact: '134 568 1234'
  },
  {
    id: 2,
    applicantName: '李四',
    propertyType: '养殖贷款',
    amount: 100000,
    status: '已批准',
    date: '2024-04-15',
    contact: '145 568 1234'
  },
  {
    id: 3,
    applicantName: '王五',
    propertyType: '农机购置',
    amount: 30000,
    status: '待审核',
    date: '2024-04-05',
    contact: '156 568 1234'
  },
  {
    id: 4,
    applicantName: '赵六',
    propertyType: '设施农业',
    amount: 200000,
    status: '已拒绝',
    date: '2024-04-12',
    contact: '165 854 1234'
  },
  {
    id: 5,
    applicantName: '钱七',
    propertyType: '农资贷款',
    amount: 80000,
    status: '待审核',
    date: '2024-04-20',
    contact: '132 567 8234'
  }
])

// 用户反馈数据
const userFeedback = ref([
  { category: "服务好评", value: 62 },
  { category: "服务投诉", value: 8 },
  { category: "功能建议", value: 15 },
  { category: "系统问题", value: 10 },
  { category: "其他", value: 5 }
])

// 惠农项目数据
const projectSales = ref([
  {
    id: 1,
    name: '水稻种植补贴',
    propertyType: '种植类',
    soldUnits: 85,
    unsoldUnits: 45,
    totalUnits: 130
  },
  {
    id: 2,
    name: '生猪养殖贷款',
    propertyType: '养殖类',
    soldUnits: 75,
    unsoldUnits: 65,
    totalUnits: 140
  },
  {
    id: 3,
    name: '智慧农业设施',
    propertyType: '设施类',
    soldUnits: 95,
    unsoldUnits: 45,
    totalUnits: 140
  },
  {
    id: 4,
    name: '农机具购置',
    propertyType: '机械类',
    soldUnits: 55,
    unsoldUnits: 45,
    totalUnits: 100
  },
  {
    id: 5,
    name: '农资采购',
    propertyType: '农资类',
    soldUnits: 154,
    unsoldUnits: 85,
    totalUnits: 239
  }
])

// 最近租赁申请
const recentRentals = ref([
  {
    id: 1,
    subject: '拖拉机X100',
    applicant: 'xiaofeifei',
    status: '待审核',
    date: '2024-04-10'
  },
  {
    id: 2,
    subject: '收割机Y200',
    applicant: 'wanglei',
    status: '申请通过',
    date: '2024-04-15'
  }
])

// 顶级代理人
const topAgent = ref({
  name: '张经理',
  title: '惠农金融专家',
  sales: 359,
  projects: 45,
  followers: 92,
  target: 500
})

const handleLogout = () => {
  localStorage.removeItem('token')
  router.push('/')
}
</script>

<template>
  <div class="home-container">
    <!-- 统计卡片区域 -->
    <div class="stats-cards">
      <div class="stats-card blue">
        <div class="stats-value">425</div>
        <div class="stats-title">农户贷款申请</div>
        <div class="stats-desc">已审批: 135</div>
        <div class="progress-bar">
          <div class="progress" style="width: 31.8%"></div>
        </div>
      </div>
      
      <div class="stats-card red">
        <div class="stats-value">185</div>
        <div class="stats-title">已完成放款</div>
        <div class="stats-desc">总申请: 230</div>
        <div class="progress-bar">
          <div class="progress" style="width: 80.4%"></div>
        </div>
      </div>
      
      <div class="stats-card purple">
        <div class="stats-value">252</div>
        <div class="stats-title">惠农补贴项目</div>
        <div class="stats-desc">已发放: 175</div>
        <div class="progress-bar">
          <div class="progress" style="width: 69.4%"></div>
        </div>
      </div>
      
      <div class="stats-card green">
        <div class="stats-value">489</div>
        <div class="stats-title">农机设备数量</div>
        <div class="stats-desc">已租赁: 185</div>
        <div class="progress-bar">
          <div class="progress" style="width: 37.8%"></div>
        </div>
      </div>
    </div>
    
    <!-- 近期贷款申请 - 独占一行 -->
    <div class="data-cards-row full-width-card">
      <!-- 最近贷款申请 -->
      <div class="data-card">
        <div class="card-header">
          <h3>近期惠农贷款申请</h3>
          <div class="card-actions">
            <button class="card-action"><i class="action-icon refresh"></i></button>
            <button class="card-action"><i class="action-icon minimize"></i></button>
            <button class="card-action"><i class="action-icon close"></i></button>
          </div>
        </div>
        <div class="card-body">
          <table class="data-table">
            <thead>
              <tr>
                <th>#</th>
                <th>申请人</th>
                <th>贷款类型</th>
                <th>日期</th>
                <th>状态</th>
                <th>联系方式</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, index) in recentLoans.slice(0, 5)" :key="item.id">
                <td>{{ index + 1 }}</td>
                <td>{{ item.applicantName }}</td>
                <td>
                  <span 
                    class="property-tag"
                    :class="{
                      'commercial': item.propertyType === '种植补贴' || item.propertyType === '农机购置',
                      'residential': item.propertyType === '养殖贷款' || item.propertyType === '设施农业' || item.propertyType === '农资贷款'
                    }"
                  >
                    {{ item.propertyType }}
                  </span>
                </td>
                <td>{{ item.date }}</td>
                <td>
                  <span 
                    class="status-tag"
                    :class="{
                      'status-pending': item.status === '待审核',
                      'status-approved': item.status === '已批准',
                      'status-rejected': item.status === '已拒绝'
                    }"
                  >
                    {{ item.status }}
                  </span>
                </td>
                <td>{{ item.contact }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
    
    <!-- 中间数据卡片区域 -->
    <div class="data-cards-row">
      <!-- 预订状态 -->
      <div class="data-card">
        <div class="card-header">
          <h3>惠农贷款申请趋势</h3>
          <div class="card-actions">
            <button class="card-action"><i class="action-icon refresh"></i></button>
            <button class="card-action"><i class="action-icon minimize"></i></button>
            <button class="card-action"><i class="action-icon close"></i></button>
          </div>
        </div>
        <div class="card-body chart-container">
          <div class="chart bar-chart">
            <div class="chart-bar" style="height: 30%; left: 5%;">
              <div class="bar-column primary" style="height: 100%;"></div>
              <div class="bar-column secondary" style="height: 80%;"></div>
              <div class="bar-label">1月</div>
            </div>
            <div class="chart-bar" style="height: 40%; left: 20%;">
              <div class="bar-column primary" style="height: 100%;"></div>
              <div class="bar-column secondary" style="height: 70%;"></div>
              <div class="bar-label">2月</div>
            </div>
            <div class="chart-bar" style="height: 70%; left: 35%;">
              <div class="bar-column primary" style="height: 100%;"></div>
              <div class="bar-column secondary" style="height: 90%;"></div>
              <div class="bar-label">3月</div>
            </div>
            <div class="chart-bar" style="height: 90%; left: 50%;">
              <div class="bar-column primary" style="height: 100%;"></div>
              <div class="bar-column secondary" style="height: 60%;"></div>
              <div class="bar-label">4月</div>
            </div>
            <div class="chart-bar" style="height: 60%; left: 65%;">
              <div class="bar-column primary" style="height: 100%;"></div>
              <div class="bar-column secondary" style="height: 80%;"></div>
              <div class="bar-label">5月</div>
            </div>
            <div class="chart-bar" style="height: 20%; left: 80%;">
              <div class="bar-column primary" style="height: 100%;"></div>
              <div class="bar-column secondary" style="height: 50%;"></div>
              <div class="bar-label">6月</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 用户反馈模块 -->
      <div class="data-card">
        <div class="card-header">
          <h3>用户反馈统计</h3>
          <div class="card-actions">
            <button class="card-action"><i class="action-icon refresh"></i></button>
            <button class="card-action"><i class="action-icon minimize"></i></button>
            <button class="card-action"><i class="action-icon close"></i></button>
          </div>
        </div>
        <div class="card-body chart-container">
          <div class="pie-chart-wrapper">
            <svg viewBox="0 0 200 200" class="pie-chart">
              <!-- 饼图扇区 - 动态计算位置 -->
              <g transform="translate(100, 100)">
                <!-- 服务好评 - 62% -->
                <path 
                  d="M0,0 L0,-80 A80,80 0 0,1 61.95,50.32 z" 
                  fill="#A5D6A7"
                />
                <!-- 服务投诉 - 8% -->
                <path 
                  d="M0,0 L61.95,50.32 A80,80 0 0,1 29.28,74.08 z" 
                  fill="#FFCDD2"
                />
                <!-- 功能建议 - 15% -->
                <path 
                  d="M0,0 L29.28,74.08 A80,80 0 0,1 -20.80,77.12 z" 
                  fill="#90CAF9"
                />
                <!-- 系统问题 - 10% -->
                <path 
                  d="M0,0 L-20.80,77.12 A80,80 0 0,1 -61.95,50.32 z" 
                  fill="#FFCC80"
                />
                <!-- 其他 - 5% -->
                <path 
                  d="M0,0 L-61.95,50.32 A80,80 0 0,1 -77.76,19.84 z" 
                  fill="#CE93D8"
                />
                <path 
                  d="M0,0 L-77.76,19.84 A80,80 0 0,1 0,-80 z" 
                  fill="#B0BEC5"
                />
              </g>
            </svg>

            <!-- 饼图图例 -->
            <div class="pie-chart-legend">
              <div class="legend-item">
                <div class="legend-color" style="background-color: #A5D6A7;"></div>
                <div class="legend-text">
                  <div class="legend-title">服务好评</div>
                  <div class="legend-value">62%</div>
                </div>
              </div>
              <div class="legend-item">
                <div class="legend-color" style="background-color: #FFCDD2;"></div>
                <div class="legend-text">
                  <div class="legend-title">服务投诉</div>
                  <div class="legend-value">8%</div>
                </div>
              </div>
              <div class="legend-item">
                <div class="legend-color" style="background-color: #90CAF9;"></div>
                <div class="legend-text">
                  <div class="legend-title">功能建议</div>
                  <div class="legend-value">15%</div>
                </div>
              </div>
              <div class="legend-item">
                <div class="legend-color" style="background-color: #FFCC80;"></div>
                <div class="legend-text">
                  <div class="legend-title">系统问题</div>
                  <div class="legend-value">10%</div>
                </div>
              </div>
              <div class="legend-item">
                <div class="legend-color" style="background-color: #CE93D8;"></div>
                <div class="legend-text">
                  <div class="legend-title">其他</div>
                  <div class="legend-value">5%</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div class="data-cards-row">
      <!-- 收入统计 -->
      <div class="data-card">
        <div class="card-header">
          <h3>惠农贷款放款统计</h3>
          <div class="card-actions">
            <button class="card-action"><i class="action-icon refresh"></i></button>
            <button class="card-action"><i class="action-icon minimize"></i></button>
            <button class="card-action"><i class="action-icon close"></i></button>
          </div>
        </div>
        <div class="card-body chart-container">
          <div class="chart-title">惠农贷款放款统计</div>
          <div class="chart line-chart">
            <svg viewBox="0 0 400 200" class="line-chart-svg">
              <!-- Y轴标签 -->
              <text x="35" y="180" font-size="12" text-anchor="end">0</text>
              <text x="35" y="140" font-size="12" text-anchor="end">5百万</text>
              <text x="35" y="100" font-size="12" text-anchor="end">10百万</text>
              <text x="35" y="60" font-size="12" text-anchor="end">15百万</text>
              <text x="35" y="20" font-size="12" text-anchor="end">20百万</text>
              
              <!-- X轴标签 -->
              <text x="70" y="195" font-size="12" text-anchor="middle">1季度</text>
              <text x="170" y="195" font-size="12" text-anchor="middle">2季度</text>
              <text x="270" y="195" font-size="12" text-anchor="middle">3季度</text>
              <text x="370" y="195" font-size="12" text-anchor="middle">4季度</text>
              
              <!-- 网格线 -->
              <line x1="40" y1="180" x2="380" y2="180" stroke="#eee" stroke-width="1"/>
              <line x1="40" y1="140" x2="380" y2="140" stroke="#eee" stroke-width="1"/>
              <line x1="40" y1="100" x2="380" y2="100" stroke="#eee" stroke-width="1"/>
              <line x1="40" y1="60" x2="380" y2="60" stroke="#eee" stroke-width="1"/>
              <line x1="40" y1="20" x2="380" y2="20" stroke="#eee" stroke-width="1"/>
              
              <!-- 红色线和区域 -->
              <path 
                d="M70,160 L170,80 L270,60 L370,100" 
                fill="none" 
                stroke="#ff6b81" 
                stroke-width="3"
                class="chart-line"
              />
              <path 
                d="M70,160 L170,80 L270,60 L370,100 L370,180 L70,180 Z" 
                fill="url(#redGradient)" 
                opacity="0.5"
                class="chart-area"
              />
              
              <!-- 蓝色线和区域 -->
              <path 
                d="M70,150 L170,100 L270,120 L370,90" 
                fill="none" 
                stroke="#54a0ff" 
                stroke-width="3"
                class="chart-line"
              />
              <path 
                d="M70,150 L170,100 L270,120 L370,90 L370,180 L70,180 Z" 
                fill="url(#blueGradient)" 
                opacity="0.5"
                class="chart-area"
              />
              
              <!-- 渐变定义 -->
              <defs>
                <linearGradient id="redGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                  <stop offset="0%" stop-color="#ff6b81" />
                  <stop offset="100%" stop-color="#ff6b81" stop-opacity="0" />
                </linearGradient>
                <linearGradient id="blueGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                  <stop offset="0%" stop-color="#54a0ff" />
                  <stop offset="100%" stop-color="#54a0ff" stop-opacity="0" />
                </linearGradient>
              </defs>
              
              <!-- 数据点 - 红色 -->
              <circle cx="70" cy="160" r="4" fill="#ff6b81" class="data-point" />
              <circle cx="170" cy="80" r="4" fill="#ff6b81" class="data-point" />
              <circle cx="270" cy="60" r="4" fill="#ff6b81" class="data-point" />
              <circle cx="370" cy="100" r="4" fill="#ff6b81" class="data-point" />
              
              <!-- 数据点 - 蓝色 -->
              <circle cx="70" cy="150" r="4" fill="#54a0ff" class="data-point" />
              <circle cx="170" cy="100" r="4" fill="#54a0ff" class="data-point" />
              <circle cx="270" cy="120" r="4" fill="#54a0ff" class="data-point" />
              <circle cx="370" cy="90" r="4" fill="#54a0ff" class="data-point" />
              
              <!-- 轴线 -->
              <line x1="40" y1="20" x2="40" y2="180" stroke="#ddd" stroke-width="1"/>
              <line x1="40" y1="180" x2="380" y2="180" stroke="#ddd" stroke-width="1"/>
            </svg>
            
            <!-- 图表说明 -->
            <div class="chart-legend">
              <div class="legend-item">
                <div class="legend-color" style="background-color: #ff6b81;"></div>
                <span>放款金额</span>
              </div>
              <div class="legend-item">
                <div class="legend-color" style="background-color: #54a0ff;"></div>
                <span>申请金额</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 项目销售 -->
      <div class="data-card">
        <div class="card-header">
          <h3>惠农贷款产品</h3>
          <div class="card-actions">
            <button class="card-action"><i class="action-icon refresh"></i></button>
            <button class="card-action"><i class="action-icon minimize"></i></button>
            <button class="card-action"><i class="action-icon close"></i></button>
          </div>
        </div>
        <div class="card-body">
          <table class="data-table">
            <thead>
              <tr>
                <th>#</th>
                <th>贷款产品</th>
                <th>类型</th>
                <th>已放款</th>
                <th>待审批</th>
                <th>总申请</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, index) in projectSales" :key="item.id">
                <td>{{ index + 1 }}</td>
                <td>{{ item.name }}</td>
                <td>
                  <span 
                    class="property-tag"
                    :class="{
                      'commercial': item.propertyType === '种植类' || item.propertyType === '机械类',
                      'residential': item.propertyType === '养殖类' || item.propertyType === '设施类' || item.propertyType === '农资类'
                    }"
                  >
                    {{ item.propertyType }}
                  </span>
                </td>
                <td>{{ item.soldUnits }}</td>
                <td>{{ item.unsoldUnits }}</td>
                <td>{{ item.totalUnits }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
    
    <!-- 底部卡片区域 -->
    <div class="data-cards-row">
      <!-- 顶级代理 -->
      <div class="data-card agent-card">
        <div class="card-header">
          <h3>优秀农金员</h3>
        </div>
        <div class="card-body">
          <div class="agent-profile">
            <div class="agent-avatar"></div>
            <div class="agent-info">
              <h4>{{ topAgent.name }}</h4>
              <p>{{ topAgent.title }}</p>
            </div>
          </div>
          <div class="agent-stats">
            <div class="stat-item">
              <div class="stat-label">贷款放款(万元)</div>
              <div class="stat-value blue">{{ topAgent.sales }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">服务农户</div>
              <div class="stat-value blue">{{ topAgent.projects }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">新增贷户</div>
              <div class="stat-value yellow">{{ topAgent.followers }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">目标(万元)</div>
              <div class="stat-value red">{{ topAgent.target }}</div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 每日查询 -->
      <div class="data-card inquiry-card">
        <div class="card-header">
          <h3>惠农贷</h3>
          <div class="card-actions">
            <button class="card-action"><i class="action-icon refresh"></i></button>
            <button class="card-action"><i class="action-icon minimize"></i></button>
            <button class="card-action"><i class="action-icon close"></i></button>
          </div>
        </div>
        <div class="card-body">
          <div class="chart-container">
            <div class="donut-chart">
              <div class="loan-reception-wrapper">
                <div class="loan-reception-label">今日受理</div>
                <div class="loan-reception-value">38</div>
                <svg viewBox="0 0 200 200" class="loan-reception-circle">
                  <!-- 背景圆 -->
                  <circle cx="100" cy="100" r="90" fill="none" stroke="#f1f1f1" stroke-width="15" />
                  
                  <!-- 进度圆 - 75%完成 -->
                  <circle 
                    cx="100" 
                    cy="100" 
                    r="90" 
                    fill="none" 
                    stroke="#6c5ce7" 
                    stroke-width="15" 
                    stroke-dasharray="565.5"
                    stroke-dashoffset="141.4"
                    transform="rotate(-90 100 100)"
                  />
                </svg>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home-container {
  padding: 0;
}

/* 统计卡片样式 */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.stats-card {
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  padding: 24px;
  display: flex;
  flex-direction: column;
  color: white;
  position: relative;
  overflow: hidden;
}

.stats-card.blue {
  background-color: #2196f3;
}

.stats-card.red {
  background-color: #f44336;
}

.stats-card.purple {
  background-color: #673ab7;
}

.stats-card.green {
  background-color: #2cc2a5;
}

.stats-value {
  font-size: 48px;
  font-weight: 700;
  margin-bottom: 5px;
  line-height: 1;
}

.stats-title {
  font-size: 16px;
  margin-bottom: 12px;
}

.stats-desc {
  font-size: 14px;
  margin-bottom: 20px;
  opacity: 0.8;
}

.progress-bar {
  height: 4px;
  background-color: rgba(255, 255, 255, 0.3);
  border-radius: 2px;
  overflow: hidden;
  position: relative;
}

.progress {
  height: 100%;
  background-color: white;
  position: absolute;
  top: 0;
  left: 0;
}

/* 保留其他样式，但移除 stats-chart-container 相关样式 */
.stats-chart-container {
  display: none;
}

/* 数据卡片样式 */
.data-cards-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(480px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.data-card {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.card-header {
  padding: 15px 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #f1f1f1;
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.card-actions {
  display: flex;
  gap: 5px;
}

.card-action {
  width: 24px;
  height: 24px;
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-action:hover {
  background-color: #f5f5f5;
}

.action-icon {
  width: 14px;
  height: 14px;
  display: block;
}

.card-body {
  padding: 20px;
}

/* 表格样式 */
.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th, 
.data-table td {
  padding: 10px;
  text-align: left;
  border-bottom: 1px solid #f1f1f1;
}

.data-table th {
  color: #666;
  font-weight: 500;
}

.property-tag {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.property-tag.commercial {
  background-color: #ede7f6;
  color: #673ab7;
}

.property-tag.residential {
  background-color: #e0f2f1;
  color: #009688;
}

.status-tag {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.status-pending {
  background-color: #e6f7ff;
  color: #1890ff;
}

.status-approved {
  background-color: #f6ffed;
  color: #52c41a;
}

.status-rejected {
  background-color: #fff2f0;
  color: #ff4d4f;
}

/* 图表容器 */
.chart-container {
  height: 250px;
  position: relative;
}

/* 柱状图样式 */
.bar-chart {
  height: 100%;
  position: relative;
  padding: 20px;
}

.chart-bar {
  position: absolute;
  width: 10%;
  bottom: 30px;
  display: flex;
  align-items: flex-end;
}

.bar-column {
  width: 45%;
  margin: 0 2.5%;
  border-radius: 4px 4px 0 0;
}

.bar-column.primary {
  background-color: #2196f3;
}

.bar-column.secondary {
  background-color: #4caf50;
}

.bar-label {
  position: absolute;
  bottom: -25px;
  width: 100%;
  text-align: center;
  font-size: 12px;
  color: #666;
}

/* 折线图样式 */
.line-chart {
  height: 100%;
  width: 100%;
  position: relative;
}

.line-chart-svg {
  width: 100%;
  height: 85%;
}

.chart-title {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin-bottom: 15px;
  padding-left: 10px;
}

.chart-line {
  transition: stroke-width 0.2s;
}

.chart-area {
  transition: opacity 0.2s;
}

.data-point {
  transition: r 0.2s;
}

.line-chart:hover .chart-line {
  stroke-width: 4;
}

.line-chart:hover .chart-area {
  opacity: 0.7;
}

.line-chart:hover .data-point {
  r: 5;
}

.chart-legend {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 10px;
}

.legend-item {
  display: flex;
  align-items: center;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  margin-right: 6px;
}

/* 代理卡片特殊样式 */
.agent-card .card-body {
  padding: 0;
}

.agent-profile {
  padding: 20px;
  display: flex;
  align-items: center;
  background: linear-gradient(to right, #f5f5f5, #e0e0e0);
}

.agent-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background-color: #fff;
  border: 3px solid #fff;
  margin-right: 15px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23ccc' d='M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z'/%3E%3C/svg%3E");
  background-position: center;
  background-size: 80%;
  background-repeat: no-repeat;
}

.agent-info h4 {
  margin: 0 0 5px;
  font-size: 18px;
}

.agent-info p {
  margin: 0;
  color: #666;
}

.agent-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  padding: 10px;
}

.stat-item {
  padding: 10px;
  text-align: center;
}

.stat-label {
  font-size: 12px;
  color: #666;
  margin-bottom: 5px;
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
}

.stat-value.blue {
  color: #2196f3;
}

.stat-value.yellow {
  color: #ffc107;
}

.stat-value.red {
  color: #f44336;
}

/* 环形图样式 */
.donut-chart {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

/* 惠农贷款受理模块样式更新 */
.inquiry-card .card-body {
  padding: 0;
}

.loan-reception-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.loan-reception-circle {
  width: 80%;
  height: auto;
  max-width: 280px;
}

.loan-reception-label {
  font-size: 24px;
  font-weight: bold;
  color: #333;
  position: absolute;
  top: 40%;
  left: 50%;
  transform: translate(-50%, -100%);
}

.loan-reception-value {
  font-size: 72px;
  font-weight: bold;
  color: #6c5ce7;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -20%);
}

/* 调整卡片头部样式 */
.inquiry-card .card-header {
  background-color: #f5f5f5;
}

.inquiry-card .card-header h3 {
  font-size: 18px;
  color: #333;
}

/* 响应式布局调整 */
@media (min-width: 992px) {
  .stats-cards {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-cards {
    grid-template-columns: 1fr;
  }
  
  .stats-value {
    font-size: 36px;
  }
}

/* 新增：全宽卡片样式 */
.full-width-card {
  grid-template-columns: 1fr;
}

/* 饼图样式 */
.pie-chart-wrapper {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.pie-chart {
  width: 50%;
  height: auto;
}

.pie-chart-legend {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-left: 20px;
}

.legend-item {
  display: flex;
  align-items: center;
}

.legend-color {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  margin-right: 8px;
}

.legend-text {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  width: 100%;
  min-width: 120px;
}

.legend-title {
  font-size: 13px;
  color: #333;
}

.legend-value {
  font-size: 13px;
  font-weight: 600;
  color: #333;
}
</style>
