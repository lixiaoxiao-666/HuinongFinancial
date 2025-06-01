<template>
  <div class="dashboard">
    <div class="page-header">
      <div class="page-title-container">
        <h2 class="page-title">数字惠农工作台</h2>
        <div class="page-subtitle">智能审批管理系统</div>
      </div>
      <div class="header-actions">
        <div class="time-info">
          <div class="current-time">{{ currentTime }}</div>
          <div class="current-date">{{ currentDate }}</div>
        </div>
        <el-button @click="fetchDashboardData" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新数据
        </el-button>
      </div>
    </div>
    
    <!-- 核心统计指标 -->
    <div class="core-stats">
      <el-row :gutter="16">
        <el-col :xs="12" :sm="6" :md="6" :lg="6">
          <div class="core-stat-card">
            <div class="stat-icon total">
              <el-icon><Files /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-number">{{ dashboardData.stats?.total_applications || 1234 }}</div>
              <div class="stat-label">总申请数</div>
              <div class="stat-trend positive">+12.5%</div>
            </div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6" :md="6" :lg="6">
          <div class="core-stat-card">
            <div class="stat-icon pending">
              <el-icon><Clock /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-number">{{ dashboardData.stats?.pending_review || 89 }}</div>
              <div class="stat-label">待处理</div>
              <div class="stat-trend negative">-3.2%</div>
            </div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6" :md="6" :lg="6">
          <div class="core-stat-card">
            <div class="stat-icon approved">
              <el-icon><CircleCheck /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-number">{{ dashboardData.stats?.approved_today || 23 }}</div>
              <div class="stat-label">今日批准</div>
              <div class="stat-trend positive">+8.7%</div>
            </div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6" :md="6" :lg="6">
          <div class="core-stat-card">
            <div class="stat-icon ai">
              <el-icon><Cpu /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-number">{{ dashboardData.stats?.ai_efficiency || 85 }}%</div>
              <div class="stat-label">AI效率</div>
              <div class="stat-trend positive">+2.1%</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
    
    <!-- 统计图表区域 -->
    <div class="charts-section">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="12" :lg="8">
          <el-card class="chart-card" shadow="hover">
            <div class="chart-header">
              <div class="chart-title">
                <el-icon><PieChart /></el-icon>
                <span>申请状态分布</span>
              </div>
              <el-dropdown @command="handleChartAction">
                <el-icon class="chart-menu"><MoreFilled /></el-icon>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="export">导出数据</el-dropdown-item>
                    <el-dropdown-item command="refresh">刷新图表</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
            <div class="chart-container">
              <div ref="pieChartRef" class="chart"></div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :xs="24" :sm="12" :md="12" :lg="8">
          <el-card class="chart-card" shadow="hover">
            <div class="chart-header">
              <div class="chart-title">
                <el-icon><Histogram /></el-icon>
                <span>部门申请统计</span>
              </div>
              <el-dropdown @command="handleChartAction">
                <el-icon class="chart-menu"><MoreFilled /></el-icon>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="export">导出数据</el-dropdown-item>
                    <el-dropdown-item command="refresh">刷新图表</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
            <div class="chart-container">
              <div ref="barChartRef" class="chart"></div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :xs="24" :sm="24" :md="24" :lg="8">
          <el-card class="chart-card" shadow="hover">
            <div class="chart-header">
              <div class="chart-title">
                <el-icon><TrendCharts /></el-icon>
                <span>AI处理效率</span>
              </div>
              <el-dropdown @command="handleChartAction">
                <el-icon class="chart-menu"><MoreFilled /></el-icon>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="export">导出数据</el-dropdown-item>
                    <el-dropdown-item command="refresh">刷新图表</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
            <div class="chart-container">
              <div ref="doughnutChartRef" class="chart"></div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 第二行图表 -->
      <el-row :gutter="20" style="margin-top: 20px;">
        <el-col :xs="24" :sm="24" :md="16" :lg="16">
          <el-card class="chart-card" shadow="hover">
            <div class="chart-header">
              <div class="chart-title">
                <el-icon><TrendCharts /></el-icon>
                <span>7天审批趋势</span>
              </div>
              <div class="chart-filters">
                <el-radio-group v-model="trendPeriod" @change="updateTrendChart">
                  <el-radio-button value="7d">7天</el-radio-button>
                  <el-radio-button value="30d">30天</el-radio-button>
                  <el-radio-button value="90d">90天</el-radio-button>
                </el-radio-group>
              </div>
            </div>
            <div class="chart-container line-chart">
              <div ref="lineChartRef" class="chart"></div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :xs="24" :sm="24" :md="8" :lg="8">
          <el-card class="chart-card system-monitor-card" shadow="hover">
            <div class="chart-header">
              <div class="chart-title">
                <el-icon><Monitor /></el-icon>
                <span>系统状态监控</span>
              </div>
            </div>
            <div class="system-status">
              <div class="status-item">
                <div class="status-indicator online"></div>
                <div class="status-info">
                  <div class="status-label">API服务</div>
                  <div class="status-value">正常运行</div>
                </div>
              </div>
              <div class="status-item">
                <div class="status-indicator online"></div>
                <div class="status-info">
                  <div class="status-label">AI引擎</div>
                  <div class="status-value">运行中</div>
                </div>
              </div>
              <div class="status-item">
                <div class="status-indicator online"></div>
                <div class="status-info">
                  <div class="status-label">数据库</div>
                  <div class="status-value">连接正常</div>
                </div>
              </div>
              <div class="status-item">
                <div class="status-indicator warning"></div>
                <div class="status-info">
                  <div class="status-label">存储空间</div>
                  <div class="status-value">78% 已使用</div>
                </div>
              </div>
              <div class="system-metrics">
                <div class="metric-item">
                  <div class="metric-label">CPU使用率</div>
                  <el-progress :percentage="cpuUsage" :stroke-width="6" />
                </div>
                <div class="metric-item">
                  <div class="metric-label">内存使用率</div>
                  <el-progress :percentage="memoryUsage" :stroke-width="6" color="#e6a23c" />
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
    
    <!-- 主要内容区域 -->
    <el-row :gutter="20" class="main-content">
      <!-- 左侧列 -->
      <el-col :xs="24" :sm="24" :md="16" :lg="16" :xl="16">
        <!-- 待办事项 -->
        <el-card class="section-card todo-section" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Bell /></el-icon>
                我的待办任务
              </span>
              <el-badge :value="pendingTasks.length" class="badge">
                <el-button type="primary" size="small" @click="goToApproval">
                  查看全部
                </el-button>
              </el-badge>
            </div>
          </template>
          
          <div v-if="pendingTasks.length === 0" class="empty-state">
            <el-empty description="暂无待办事项">
              <el-button type="primary" @click="goToApproval">去审批中心</el-button>
            </el-empty>
          </div>
          
          <div v-else class="task-list">
            <div
              v-for="task in pendingTasks"
              :key="task.task_id"
              class="task-item"
              @click="handleTaskClick(task)"
            >
              <div class="task-priority" :class="task.priority"></div>
              <div class="task-content">
                <div class="task-title">{{ task.title }}</div>
                <div class="task-desc">{{ task.task_type }}</div>
                <div class="task-meta">
                  <span class="task-time">{{ formatRelativeTime(task.created_at) }}</span>
                  <el-tag size="small" :type="getPriorityType(task.priority)">
                    {{ getPriorityText(task.priority) }}
                  </el-tag>
                </div>
              </div>
              <div class="task-action">
                <el-icon><ArrowRight /></el-icon>
              </div>
            </div>
          </div>
        </el-card>
        
        <!-- 审批效率分析 -->
        <el-card class="section-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><DataAnalysis /></el-icon>
                审批效率分析
              </span>
              <div class="efficiency-period">
                <el-radio-group v-model="efficiencyPeriod" @change="updateEfficiencyCharts" size="small">
                  <el-radio-button value="week">本周</el-radio-button>
                  <el-radio-button value="month">本月</el-radio-button>
                  <el-radio-button value="quarter">本季</el-radio-button>
                </el-radio-group>
              </div>
            </div>
          </template>
          
          <div class="efficiency-analysis">
            <!-- 核心指标概览 -->
            <el-row :gutter="16" class="efficiency-overview">
              <el-col :span="8">
                <div class="efficiency-item">
                  <div class="efficiency-number">2.3</div>
                  <div class="efficiency-label">平均处理时间</div>
                  <div class="efficiency-unit">小时</div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="efficiency-item">
                  <div class="efficiency-number">94.2%</div>
                  <div class="efficiency-label">综合准确率</div>
                  <div class="efficiency-unit">本月</div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="efficiency-item">
                  <div class="efficiency-number">156</div>
                  <div class="efficiency-label">今日处理量</div>
                  <div class="efficiency-unit">件</div>
                </div>
              </el-col>
            </el-row>
            
            <!-- AI和人工效率趋势图表 -->
            <div class="efficiency-charts">
              <el-row :gutter="20">
                <el-col :span="12">
                  <div class="efficiency-chart-card">
                    <div class="chart-header-mini">
                      <div class="chart-title-mini">
                        <div class="chart-icon ai-icon">
                          <el-icon><Cpu /></el-icon>
                        </div>
                        <div class="chart-info">
                          <div class="chart-name">AI处理效率</div>
                          <div class="chart-value">{{ dashboardData.stats?.ai_efficiency || 85 }}%</div>
                        </div>
                      </div>
                      <div class="chart-trend positive">
                        <el-icon><TrendCharts /></el-icon>
                        <span>+2.3%</span>
                      </div>
                    </div>
                    <div class="mini-chart-container">
                      <div ref="aiEfficiencyChartRef" class="mini-chart"></div>
                    </div>
                  </div>
                </el-col>
                <el-col :span="12">
                  <div class="efficiency-chart-card">
                    <div class="chart-header-mini">
                      <div class="chart-title-mini">
                        <div class="chart-icon human-icon">
                          <el-icon><User /></el-icon>
                        </div>
                        <div class="chart-info">
                          <div class="chart-name">人工审核效率</div>
                          <div class="chart-value">92.1%</div>
                        </div>
                      </div>
                      <div class="chart-trend positive">
                        <el-icon><TrendCharts /></el-icon>
                        <span>+1.8%</span>
                      </div>
                    </div>
                    <div class="mini-chart-container">
                      <div ref="humanEfficiencyChartRef" class="mini-chart"></div>
                    </div>
                  </div>
                </el-col>
              </el-row>
            </div>
            
            <!-- 效率对比分析 -->
            <div class="efficiency-comparison">
              <div class="comparison-header">
                <h4>效率对比分析</h4>
                <div class="comparison-legend">
                  <div class="legend-item">
                    <div class="legend-color ai-color"></div>
                    <span>AI处理</span>
                  </div>
                  <div class="legend-item">
                    <div class="legend-color human-color"></div>
                    <span>人工审核</span>
                  </div>
                </div>
              </div>
              <div class="comparison-chart-container">
                <div ref="comparisonChartRef" class="comparison-chart"></div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <!-- 右侧列 -->
      <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
        <!-- 快捷操作 -->
        <el-card class="section-card quick-actions-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Operation /></el-icon>
                快捷操作中心
              </span>
            </div>
          </template>
          
          <div class="quick-actions">
            <div class="quick-action-item primary" @click="goToPage('/approval')">
              <el-icon class="action-icon">
                <DocumentChecked />
              </el-icon>
              <span>审批看板</span>
              <div class="action-badge">{{ dashboardData.stats?.pending_review || 0 }}</div>
            </div>
            <div class="quick-action-item ai" @click="goToPage('/smart-approval')">
              <el-icon class="action-icon">
                <Cpu />
              </el-icon>
              <span>贷款审批</span>
              <div class="action-badge new">AI</div>
            </div>
            <div class="quick-action-item success" @click="goToPage('/lease-approval')">
              <el-icon class="action-icon">
                <Van />
              </el-icon>
              <span>租赁审批</span>
            </div>
            <div class="quick-action-item warning" @click="goToPage('/ai-workflow')">
              <el-icon class="action-icon">
                <Promotion />
              </el-icon>
              <span>AI审批流</span>
            </div>
            <div class="quick-action-item info" @click="goToPage('/users')">
              <el-icon class="action-icon">
                <User />
              </el-icon>
              <span>用户管理</span>
            </div>
            <div class="quick-action-item default" @click="goToPage('/logs')">
              <el-icon class="action-icon">
                <Document />
              </el-icon>
              <span>操作日志</span>
            </div>
          </div>
        </el-card>

        <!-- AI控制面板 -->
        <el-card class="section-card ai-control-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Cpu /></el-icon>
                AI控制面板
              </span>
            </div>
          </template>
          
          <div class="ai-control-panel">
            <div class="ai-status-card">
              <div class="ai-status-header">
                <div class="ai-status-indicator active"></div>
                <div class="ai-status-info">
                  <div class="ai-status-title">AI审批引擎</div>
                  <div class="ai-status-desc">正在工作</div>
                </div>
                <el-switch
                  v-model="aiEnabled"
                  @change="handleAIToggle"
                  :disabled="true"
                  size="large"
                />
              </div>
            </div>
            
            <div class="ai-metrics">
              <div class="ai-metric-item">
                <div class="metric-icon">
                  <el-icon><Lightning /></el-icon>
                </div>
                <div class="metric-content">
                  <div class="metric-value">{{ dashboardData.stats?.ai_efficiency || 85 }}%</div>
                  <div class="metric-label">处理效率</div>
                </div>
              </div>
              <div class="ai-metric-item">
                <div class="metric-icon">
                  <el-icon><Aim /></el-icon>
                </div>
                <div class="metric-content">
                  <div class="metric-value">94.2%</div>
                  <div class="metric-label">准确率</div>
                </div>
              </div>
              <div class="ai-metric-item">
                <div class="metric-icon">
                  <el-icon><Timer /></el-icon>
                </div>
                <div class="metric-content">
                  <div class="metric-value">1.2s</div>
                  <div class="metric-label">平均响应</div>
                </div>
              </div>
            </div>
          </div>
        </el-card>
        
        <!-- 数据统计展示 -->
        <el-card class="section-card data-display-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><DataAnalysis /></el-icon>
                数据统计概览
              </span>
            </div>
          </template>
          
          <div class="data-display-panel">
            <div class="data-overview">
              <div class="data-item">
                <div class="data-icon">
                  <el-icon><Files /></el-icon>
                </div>
                <div class="data-content">
                  <div class="data-value">{{ dashboardData.stats?.total_applications || 1234 }}</div>
                  <div class="data-label">总申请数</div>
                </div>
              </div>
              <div class="data-item">
                <div class="data-icon">
                  <el-icon><Clock /></el-icon>
                </div>
                <div class="data-content">
                  <div class="data-value">{{ dashboardData.stats?.pending_review || 89 }}</div>
                  <div class="data-label">待审批</div>
                </div>
              </div>
              <div class="data-item">
                <div class="data-icon">
                  <el-icon><CircleCheck /></el-icon>
                </div>
                <div class="data-content">
                  <div class="data-value">{{ dashboardData.stats?.approved_today || 23 }}</div>
                  <div class="data-label">今日通过</div>
                </div>
              </div>
            </div>
            
            <div class="progress-section">
              <div class="progress-item">
                <div class="progress-label">
                  <span>审批进度</span>
                  <span class="progress-value">75%</span>
                </div>
                <el-progress :percentage="75" :stroke-width="8" color="#4CAF50" />
              </div>
              <div class="progress-item">
                <div class="progress-label">
                  <span>AI处理率</span>
                  <span class="progress-value">{{ dashboardData.stats?.ai_efficiency || 85 }}%</span>
                </div>
                <el-progress :percentage="dashboardData.stats?.ai_efficiency || 85" :stroke-width="8" color="#2196F3" />
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Files,
  Clock,
  CircleCheck,
  Cpu,
  Bell,
  Operation,
  Document,
  ArrowRight,
  DocumentChecked,
  Setting,
  User,
  Refresh,
  Van,
  Promotion,
  PieChart,
  Histogram,
  TrendCharts,
  Monitor,
  MoreFilled,
  DataAnalysis,
  Lightning,
  Aim,
  Timer
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { getDashboard, toggleAIApproval } from '@/api/admin'
import type { DashboardData } from '@/types'
import dayjs from 'dayjs'
import * as echarts from 'echarts'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const aiEnabled = ref(true)
const trendPeriod = ref('7d')
const efficiencyPeriod = ref('week')
const currentTime = ref('')
const currentDate = ref('')
const cpuUsage = ref(45)
const memoryUsage = ref(68)

const dashboardData = ref<DashboardData>({
  stats: {
    total_applications: 0,
    pending_review: 0,
    approved_today: 0,
    ai_efficiency: 0
  },
  pending_tasks: [],
  recent_activities: []
})

// 图表DOM引用
const pieChartRef = ref<HTMLElement | null>(null)
const barChartRef = ref<HTMLElement | null>(null)
const lineChartRef = ref<HTMLElement | null>(null)
const doughnutChartRef = ref<HTMLElement | null>(null)
const aiEfficiencyChartRef = ref<HTMLElement | null>(null)
const humanEfficiencyChartRef = ref<HTMLElement | null>(null)
const comparisonChartRef = ref<HTMLElement | null>(null)

// 图表实例
let pieChart: echarts.ECharts | null = null
let barChart: echarts.ECharts | null = null
let lineChart: echarts.ECharts | null = null
let doughnutChart: echarts.ECharts | null = null
let aiEfficiencyChart: echarts.ECharts | null = null
let humanEfficiencyChart: echarts.ECharts | null = null
let comparisonChart: echarts.ECharts | null = null

// 更新时间显示
const updateTime = () => {
  const now = dayjs()
  currentTime.value = now.format('HH:mm:ss')
  currentDate.value = now.format('YYYY年MM月DD日 dddd')
}

// 模拟数据
const getWeeklyData = () => {
  return ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
}

const getPendingByDepartment = () => {
  return [
    { value: 32, name: '农业补贴' },
    { value: 24, name: '金融贷款' },
    { value: 18, name: '设备申请' },
    { value: 12, name: '技术支持' }
  ]
}

const getApprovedTrend = () => {
  return [15, 22, 28, 25, 35, 40, 32]
}

// 计算属性
const pendingTasks = computed(() => dashboardData.value.pending_tasks || [])
const hasPermission = computed(() => authStore.hasPermission)

// 初始化图表
const initCharts = () => {
  const stats = dashboardData.value.stats
  
  // 饼图 - 总申请数据分析
  if (pieChartRef.value) {
    pieChart = echarts.init(pieChartRef.value)
    const totalApps = stats?.total_applications || 1234
    const pendingApps = stats?.pending_review || 89
    const approvedApps = Math.floor(totalApps * 0.7) // 估算已批准
    const rejectedApps = totalApps - pendingApps - approvedApps // 估算已拒绝
    
    const option = {
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'horizontal',
        bottom: 0,
        itemWidth: 10,
        itemHeight: 10,
        textStyle: {
          fontSize: 12
        }
      },
      series: [
        {
          name: '申请状态',
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 5,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: false
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 14,
              fontWeight: 'bold'
            }
          },
          labelLine: {
            show: false
          },
          data: [
            { 
              value: pendingApps, 
              name: '待处理',
              itemStyle: {
                color: '#FF9F43'
              }
            },
            { 
              value: approvedApps, 
              name: '已批准',
              itemStyle: {
                color: '#28C76F'
              }
            },
            { 
              value: rejectedApps, 
              name: '已拒绝',
              itemStyle: {
                color: '#EA5455'
              }
            }
          ]
        }
      ]
    }
    pieChart.setOption(option)
  }
  
  // 柱状图 - 待处理申请统计
  if (barChartRef.value) {
    barChart = echarts.init(barChartRef.value)
    const departments = getPendingByDepartment()
    const option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '12%',
        top: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: departments.map(item => item.name),
        axisLine: {
          lineStyle: {
            color: '#E0E0E0'
          }
        },
        axisLabel: {
          interval: 0,
          rotate: 30,
          fontSize: 10,
          color: '#666'
        }
      },
      yAxis: {
        type: 'value',
        axisLine: {
          show: false
        },
        axisTick: {
          show: false
        },
        splitLine: {
          lineStyle: {
            color: '#EFEFEF'
          }
        }
      },
      series: [
        {
          name: '待处理数量',
          type: 'bar',
          data: departments.map(item => item.value),
          barWidth: '40%',
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#5E72E4' },
              { offset: 1, color: '#768BF9' }
            ]),
            borderRadius: [5, 5, 0, 0]
          }
        }
      ]
    }
    barChart.setOption(option)
  }
  
  // 趋势图 - 已批准趋势
  if (lineChartRef.value) {
    lineChart = echarts.init(lineChartRef.value)
    const weekdays = getWeeklyData()
    const approvedData = getApprovedTrend()
    const option = {
      tooltip: {
        trigger: 'axis'
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '12%',
        top: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: weekdays,
        axisLine: {
          lineStyle: {
            color: '#E0E0E0'
          }
        },
        axisLabel: {
          color: '#666'
        }
      },
      yAxis: {
        type: 'value',
        axisLine: {
          show: false
        },
        axisTick: {
          show: false
        },
        splitLine: {
          lineStyle: {
            color: '#EFEFEF'
          }
        }
      },
      series: [
        {
          name: '已批准',
          type: 'line',
          smooth: true,
          symbol: 'circle',
          symbolSize: 6,
          itemStyle: {
            color: '#4CAF50'
          },
          lineStyle: {
            width: 3,
            color: '#4CAF50'
          },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(76, 175, 80, 0.5)' },
              { offset: 1, color: 'rgba(76, 175, 80, 0.05)' }
            ])
          },
          data: approvedData
        }
      ]
    }
    lineChart.setOption(option)
  }
  
  // 圆环图 - AI处理占比
  if (doughnutChartRef.value) {
    doughnutChart = echarts.init(doughnutChartRef.value)
    const aiEfficiency = stats?.ai_efficiency || 85
    const option = {
      tooltip: {
        formatter: '{b}: {c}%'
      },
      series: [
        {
          name: 'AI处理',
          type: 'pie',
          radius: ['60%', '80%'],
          avoidLabelOverlap: false,
          label: {
            show: false
          },
          emphasis: {
            label: {
              show: false
            }
          },
          labelLine: {
            show: false
          },
          data: [
            { 
              value: aiEfficiency, 
              name: 'AI处理率',
              itemStyle: {
                color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                  { offset: 0, color: '#7367F0' },
                  { offset: 1, color: '#A9A2F6' }
                ])
              }
            },
            { 
              value: 100 - aiEfficiency, 
              name: '人工处理率',
              itemStyle: {
                color: '#F4F4F4'
              }
            }
          ]
        },
        {
          name: 'AI处理',
          type: 'pie',
          radius: ['0', '45%'],
          avoidLabelOverlap: false,
          label: {
            show: true,
            position: 'center',
            formatter: '{b}\n{c}%',
            fontSize: 14,
            fontWeight: 'bold',
            color: '#7367F0'
          },
          emphasis: {
            label: {
              show: true
            }
          },
          labelLine: {
            show: false
          },
          data: [
            { 
              value: aiEfficiency, 
              name: 'AI处理率',
              itemStyle: {
                color: 'transparent'
              }
            }
          ]
        }
      ]
    }
    doughnutChart.setOption(option)
  }
}

const resizeCharts = () => {
  pieChart?.resize()
  barChart?.resize()
  lineChart?.resize()
  doughnutChart?.resize()
  aiEfficiencyChart?.resize()
  humanEfficiencyChart?.resize()
  comparisonChart?.resize()
}

// 方法
const fetchDashboardData = async () => {
  try {
    loading.value = true
    const data = await getDashboard()
    dashboardData.value = data
    aiEnabled.value = true // 始终设置为运行中
    
    // 重新初始化图表
    setTimeout(() => {
      initCharts()
    }, 100)
  } catch (error) {
    ElMessage.error('获取工作台数据失败')
  } finally {
    loading.value = false
  }
}

const handleAIToggle = async (enabled: boolean) => {
  // AI功能始终保持运行状态，不允许切换
  aiEnabled.value = true
  ElMessage.info('AI审批功能始终保持运行状态')
}

const handleChartAction = (command: string) => {
  switch (command) {
    case 'export':
      ElMessage.success('导出功能开发中')
      break
    case 'refresh':
      initCharts()
      ElMessage.success('图表已刷新')
      break
    default:
      break
  }
}

const updateTrendChart = () => {
  // 重新绘制趋势图表
  if (lineChart) {
    const weekdays = getWeeklyData()
    const approvedData = getApprovedTrend()
    
    const option = {
      tooltip: {
        trigger: 'axis'
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '12%',
        top: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: weekdays,
        axisLine: {
          lineStyle: {
            color: '#E0E0E0'
          }
        },
        axisLabel: {
          color: '#666'
        }
      },
      yAxis: {
        type: 'value',
        axisLine: {
          show: false
        },
        axisTick: {
          show: false
        },
        splitLine: {
          lineStyle: {
            color: '#EFEFEF'
          }
        }
      },
      series: [
        {
          name: '已批准',
          type: 'line',
          smooth: true,
          symbol: 'circle',
          symbolSize: 6,
          itemStyle: {
            color: '#4CAF50'
          },
          lineStyle: {
            width: 3,
            color: '#4CAF50'
          },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(76, 175, 80, 0.5)' },
              { offset: 1, color: 'rgba(76, 175, 80, 0.05)' }
            ])
          },
          data: approvedData
        }
      ]
    }
    lineChart.setOption(option)
  }
}

const handleTaskClick = (task: any) => {
  if (task.application_id) {
    router.push(`/approval/${task.application_id}`)
  } else {
    router.push('/approval')
  }
}

const goToApproval = () => {
  router.push('/approval')
}

const goToPage = (path: string) => {
  router.push(path)
}

const getPriorityType = (priority: string) => {
  switch (priority) {
    case 'high': return 'danger'
    case 'medium': return 'warning'
    case 'low': return 'success'
    default: return 'info'
  }
}

const getPriorityText = (priority: string) => {
  switch (priority) {
    case 'high': return '高优先级'
    case 'medium': return '中优先级'
    case 'low': return '低优先级'
    default: return '普通优先级'
  }
}

const formatRelativeTime = (datetime: string) => {
  const now = dayjs()
  const time = dayjs(datetime)
  const hours = now.diff(time, 'hour')
  
  if (hours < 1) {
    const minutes = now.diff(time, 'minute')
    return `${minutes}分钟前`
  } else if (hours < 24) {
    return `${hours}小时前`
  } else {
    const days = now.diff(time, 'day')
    return `${days}天前`
  }
}

// 获取效率数据
const getEfficiencyData = (period: string) => {
  const data: Record<string, { ai: number[], human: number[], labels: string[] }> = {
    week: {
      ai: [82, 85, 87, 84, 88, 90, 85],
      human: [89, 91, 88, 92, 90, 89, 92],
      labels: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
    },
    month: {
      ai: [80, 82, 85, 87, 84, 88, 90, 85, 89, 91, 88, 92, 90, 89, 92, 85, 87, 89, 91, 88, 92, 90, 89, 92, 85, 87, 89, 91, 88, 92],
      human: [88, 89, 91, 88, 92, 90, 89, 92, 85, 87, 89, 91, 88, 92, 90, 89, 92, 85, 87, 89, 91, 88, 92, 90, 89, 92, 85, 87, 89, 91],
      labels: Array.from({length: 30}, (_, i) => `${i + 1}日`)
    },
    quarter: {
      ai: [78, 80, 82, 85, 87, 84, 88, 90, 85, 89, 91, 88],
      human: [86, 88, 89, 91, 88, 92, 90, 89, 92, 85, 87, 89],
      labels: ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月']
    }
  }
  return data[period] || data.week
}

// 初始化效率图表
const initEfficiencyCharts = () => {
  const efficiencyData = getEfficiencyData(efficiencyPeriod.value)
  
  // AI效率趋势图
  if (aiEfficiencyChartRef.value) {
    aiEfficiencyChart = echarts.init(aiEfficiencyChartRef.value)
    const option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'line'
        }
      },
      grid: {
        left: '5%',
        right: '5%',
        bottom: '10%',
        top: '10%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: efficiencyData.labels,
        axisLine: {
          show: false
        },
        axisTick: {
          show: false
        },
        axisLabel: {
          show: false
        }
      },
      yAxis: {
        type: 'value',
        show: false,
        min: 70,
        max: 100
      },
      series: [
        {
          type: 'line',
          data: efficiencyData.ai,
          smooth: true,
          symbol: 'none',
          lineStyle: {
            width: 3,
            color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
              { offset: 0, color: '#4facfe' },
              { offset: 1, color: '#00f2fe' }
            ])
          },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(79, 172, 254, 0.3)' },
              { offset: 1, color: 'rgba(79, 172, 254, 0.05)' }
            ])
          }
        }
      ]
    }
    aiEfficiencyChart.setOption(option)
  }
  
  // 人工效率趋势图
  if (humanEfficiencyChartRef.value) {
    humanEfficiencyChart = echarts.init(humanEfficiencyChartRef.value)
    const option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'line'
        }
      },
      grid: {
        left: '5%',
        right: '5%',
        bottom: '10%',
        top: '10%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: efficiencyData.labels,
        axisLine: {
          show: false
        },
        axisTick: {
          show: false
        },
        axisLabel: {
          show: false
        }
      },
      yAxis: {
        type: 'value',
        show: false,
        min: 70,
        max: 100
      },
      series: [
        {
          type: 'line',
          data: efficiencyData.human,
          smooth: true,
          symbol: 'none',
          lineStyle: {
            width: 3,
            color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
              { offset: 0, color: '#a8edea' },
              { offset: 1, color: '#fed6e3' }
            ])
          },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(168, 237, 234, 0.3)' },
              { offset: 1, color: 'rgba(168, 237, 234, 0.05)' }
            ])
          }
        }
      ]
    }
    humanEfficiencyChart.setOption(option)
  }
  
  // 效率对比图表
  if (comparisonChartRef.value) {
    comparisonChart = echarts.init(comparisonChartRef.value)
    const option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross'
        }
      },
      legend: {
        data: ['AI处理', '人工审核'],
        bottom: 0,
        itemWidth: 12,
        itemHeight: 12
      },
      grid: {
        left: '5%',
        right: '5%',
        bottom: '15%',
        top: '5%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: efficiencyData.labels,
        axisLine: {
          lineStyle: {
            color: '#E0E0E0'
          }
        },
        axisLabel: {
          color: '#666',
          fontSize: 10
        }
      },
      yAxis: {
        type: 'value',
        min: 70,
        max: 100,
        axisLine: {
          show: false
        },
        axisTick: {
          show: false
        },
        splitLine: {
          lineStyle: {
            color: '#EFEFEF'
          }
        },
        axisLabel: {
          color: '#666'
        }
      },
      series: [
        {
          name: 'AI处理',
          type: 'line',
          data: efficiencyData.ai,
          smooth: true,
          symbol: 'circle',
          symbolSize: 4,
          lineStyle: {
            width: 2,
            color: '#4facfe'
          },
          itemStyle: {
            color: '#4facfe'
          }
        },
        {
          name: '人工审核',
          type: 'line',
          data: efficiencyData.human,
          smooth: true,
          symbol: 'circle',
          symbolSize: 4,
          lineStyle: {
            width: 2,
            color: '#a8edea'
          },
          itemStyle: {
            color: '#a8edea'
          }
        }
      ]
    }
    comparisonChart.setOption(option)
  }
}

const updateEfficiencyCharts = () => {
  initEfficiencyCharts()
}

onMounted(() => {
  // 开始更新时间
  updateTime()
  setInterval(updateTime, 1000)
  
  fetchDashboardData()
  // 需要等待DOM更新后初始化图表
  setTimeout(() => {
    initCharts()
    initEfficiencyCharts()
    window.addEventListener('resize', resizeCharts)
  }, 200)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', resizeCharts)
  pieChart?.dispose()
  barChart?.dispose()
  lineChart?.dispose()
  doughnutChart?.dispose()
  aiEfficiencyChart?.dispose()
  humanEfficiencyChart?.dispose()
  comparisonChart?.dispose()
})
</script>

<style scoped>
.dashboard {
  padding: 0;
  min-height: 100%;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 0 20px 0;
  margin-bottom: 20px;
  margin-left: -8px;
  border-bottom: 1px solid #ebeef5;
}

.page-title-container {
  padding: 0;
}

.page-title-container:hover {
  transform: none;
  box-shadow: none;
}

.page-title {
  margin: 0;
  color: #1976D2;
  font-size: 24px;
  font-weight: 700;
  letter-spacing: 0.5px;
  text-shadow: none;
}

.page-subtitle {
  color: #666;
  font-size: 14px;
  margin-top: 2px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.time-info {
  display: flex;
  flex-direction: column;
  margin-right: 12px;
}

.current-time {
  font-size: 18px;
  font-weight: 600;
  color: #1976D2;
}

.current-date {
  font-size: 12px;
  color: #999;
}

.core-stats {
  margin-bottom: 24px;
}

.core-stat-card {
  display: flex;
  align-items: center;
  padding: 20px;
  background: white;
  border-radius: 12px;
  border: 1px solid #ebeef5;
  transition: transform 0.2s, box-shadow 0.2s;
  height: 100px;
}

.core-stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.08);
}

.stat-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 10px;
  margin-right: 12px;
}

.stat-icon.total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.stat-icon.pending {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.stat-icon.approved {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
}

.stat-icon.ai {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  color: white;
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-number {
  font-size: 24px;
  font-weight: 700;
  color: #333;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  color: #666;
  font-size: 14px;
  margin-bottom: 4px;
}

.stat-trend {
  font-size: 12px;
  font-weight: 500;
}

.stat-trend.positive {
  color: #28C76F;
}

.stat-trend.negative {
  color: #EA5455;
}

.charts-section {
  margin-bottom: 24px;
}

.chart-card {
  border-radius: 12px;
  border: none;
  transition: transform 0.2s, box-shadow 0.2s;
  overflow: hidden;
}

.chart-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.08);
}

.chart-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f2f5;
}

.chart-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #333;
}

.chart-menu {
  cursor: pointer;
  color: #c0c4cc;
  transition: color 0.2s;
}

.chart-menu:hover {
  color: #409eff;
}

.chart-filters {
  display: flex;
  align-items: center;
}

.chart-container {
  height: 240px;
  width: 100%;
  padding: 8px;
}

.chart-container.line-chart {
  height: 280px;
}

.chart {
  height: 100%;
  width: 100%;
}

.system-status {
  padding: 16px;
}

.system-monitor-card {
  max-height: calc(100% - 38px);
}

.system-monitor-card .system-status {
  padding: 8px 16px;
}

.system-monitor-card .status-item {
  padding: 8px 0;
}

.system-monitor-card .system-metrics {
  margin-top: 8px;
}

.system-monitor-card .metric-item {
  margin-bottom: 8px;
}

.status-item {
  display: flex;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f2f5;
}

.status-item:last-child {
  border-bottom: none;
}

.status-indicator {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  margin-right: 12px;
}

.status-indicator.online {
  background-color: #28C76F;
}

.status-indicator.warning {
  background-color: #FF9F43;
}

.status-indicator.offline {
  background-color: #EA5455;
}

.status-info {
  flex: 1;
  min-width: 0;
}

.status-label {
  color: #333;
  font-weight: 500;
  font-size: 14px;
}

.status-value {
  color: #666;
  font-size: 12px;
}

.system-metrics {
  margin-top: 16px;
}

.metric-item {
  margin-bottom: 12px;
}

.metric-label {
  color: #666;
  font-size: 12px;
  margin-bottom: 4px;
}

.efficiency-analysis {
  padding: 16px 0;
}

.efficiency-overview {
  margin-bottom: 24px;
}

.efficiency-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
}

.efficiency-number {
  font-size: 24px;
  font-weight: 700;
  color: #1976D2;
  line-height: 1;
  margin-bottom: 4px;
}

.efficiency-label {
  color: #666;
  font-size: 14px;
  margin-bottom: 2px;
}

.efficiency-unit {
  color: #999;
  font-size: 12px;
}

.efficiency-period {
  display: flex;
  align-items: center;
}

.efficiency-charts {
  margin-bottom: 24px;
}

.efficiency-chart-card {
  background: white;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  padding: 16px;
  transition: box-shadow 0.2s;
}

.efficiency-chart-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.chart-header-mini {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.chart-title-mini {
  display: flex;
  align-items: center;
  gap: 12px;
}

.chart-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 8px;
}

.chart-icon.ai-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
}

.chart-icon.human-icon {
  background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
  color: #333;
}

.chart-info {
  display: flex;
  flex-direction: column;
}

.chart-name {
  font-size: 14px;
  color: #666;
  margin-bottom: 4px;
}

.chart-value {
  font-size: 20px;
  font-weight: 700;
  color: #333;
  line-height: 1;
}

.chart-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 600;
}

.chart-trend.positive {
  color: #28C76F;
}

.chart-trend.negative {
  color: #EA5455;
}

.mini-chart-container {
  height: 60px;
  width: 100%;
}

.mini-chart {
  height: 100%;
  width: 100%;
}

.efficiency-comparison {
  margin-top: 24px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
}

.comparison-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.comparison-header h4 {
  margin: 0;
  font-size: 16px;
  color: #333;
}

.comparison-legend {
  display: flex;
  gap: 16px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #666;
}

.legend-color {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.legend-color.ai-color {
  background-color: #4facfe;
}

.legend-color.human-color {
  background-color: #a8edea;
}

.comparison-chart-container {
  height: 200px;
  width: 100%;
}

.comparison-chart {
  height: 100%;
  width: 100%;
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

@media (max-width: 768px) {
  .quick-actions {
    grid-template-columns: 1fr;
  }
}

.quick-action-item {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 12px;
  background: #f8f9fa;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  gap: 8px;
  text-align: center;
  border: 2px solid transparent;
}

.quick-action-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.quick-action-item.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.quick-action-item.ai {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  color: white;
}

.quick-action-item.success {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
}

.quick-action-item.warning {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.quick-action-item.info {
  background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
  color: #333;
}

.quick-action-item.default {
  background: #f8f9fa;
  color: #333;
  border-color: #ebeef5;
}

.action-icon {
  font-size: 24px;
}

.action-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  background: #ff4757;
  color: white;
  border-radius: 12px;
  padding: 2px 6px;
  font-size: 10px;
  font-weight: 600;
  min-width: 16px;
  text-align: center;
}

.action-badge.new {
  background: #ff6b6b;
}

.ai-control-panel {
  padding: 16px 0;
}

.quick-actions-card,
.ai-control-card {
  margin-top: 0;
}

.quick-actions-card {
  margin-top: -19px;
  width: calc(100% + 11px);
  margin-right: -11px;
}

.quick-actions-card .quick-actions {
  min-height: 269px;
}

.quick-actions-card .quick-action-item {
  padding: 20px 12px;
  min-height: 69px;
}

.ai-control-card {
  margin-top: 0;
}

.ai-control-card .ai-control-panel {
  padding: 20px 0;
  min-height: 200px;
}

.ai-control-card .ai-status-card {
  padding: 20px;
  margin-bottom: 20px;
}

.ai-control-card .ai-metric-item {
  padding: 16px;
  min-height: 90px;
}

.data-display-card {
  margin-top: 20px;
}

.data-display-panel {
  padding: 16px 0;
}

.data-overview {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 20px;
}

.data-item {
  flex: 1;
  display: flex;
  align-items: center;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
  gap: 12px;
}

.data-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-size: 16px;
}

.data-content {
  flex: 1;
  text-align: center;
}

.data-value {
  font-size: 18px;
  font-weight: 700;
  color: #333;
  line-height: 1;
  margin-bottom: 4px;
}

.data-label {
  color: #666;
  font-size: 12px;
}

.progress-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.progress-item {
  background: white;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #ebeef5;
}

.progress-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
  color: #333;
}

.progress-value {
  font-weight: 600;
  color: #1976D2;
}

.ai-status-card {
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
  margin-bottom: 16px;
}

.ai-status-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.ai-status-indicator {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background-color: #e4e7ed;
  margin-right: 12px;
  transition: background-color 0.3s;
}

.ai-status-indicator.active {
  background-color: #28C76F;
  box-shadow: 0 0 0 3px rgba(40, 199, 111, 0.2);
}

.ai-status-info {
  flex: 1;
  min-width: 0;
}

.ai-status-title {
  font-weight: 600;
  color: #333;
  font-size: 16px;
}

.ai-status-desc {
  color: #666;
  font-size: 12px;
  margin-top: 2px;
}

.ai-metrics {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.ai-metric-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px;
  background: white;
  border-radius: 8px;
  border: 1px solid #ebeef5;
}

.metric-icon {
  margin-bottom: 8px;
  color: #409eff;
}

.metric-content {
  text-align: center;
}

.metric-value {
  font-size: 16px;
  font-weight: 700;
  color: #333;
  line-height: 1;
}

.metric-label {
  color: #666;
  font-size: 12px;
  margin-top: 4px;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
    margin-left: -6px;
  }
  
  .page-title-container {
    padding: 6px 12px 6px 4px;
  }
  
  .page-title {
    font-size: 20px;
  }
  
  .page-subtitle {
    font-size: 12px;
  }
  
  .header-actions {
    justify-content: space-between;
  }
  
  .time-info {
    order: -1;
  }
  
  .core-stat-card {
    height: 80px;
    padding: 16px;
  }
  
  .stat-icon {
    width: 36px;
    height: 36px;
    margin-right: 12px;
  }
  
  .stat-number {
    font-size: 20px;
  }
  
  .charts-section .el-col {
    margin-bottom: 16px;
  }
  
  .ai-metrics {
    flex-direction: column;
    gap: 8px;
  }
  
  .ai-metric-item {
    flex-direction: row;
    text-align: left;
  }
  
  .metric-icon {
    margin-right: 12px;
    margin-bottom: 0;
  }
}

.main-content {
  margin-top: 0;
}

@media (max-width: 768px) {
  .main-content {
    gap: 16px;
  }
}

.section-card {
  margin-bottom: 20px;
  border-radius: 12px;
  border: 1px solid #ebeef5;
}

.todo-section {
  margin-top: -19px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
  color: #333;
  gap: 8px;
}

.badge {
  margin-right: 0;
}

.empty-state {
  padding: 40px 0;
}

.task-list {
  max-height: 400px;
  overflow-y: auto;
}

.task-item {
  display: flex;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid #f0f2f5;
  cursor: pointer;
  transition: all 0.3s;
}

.task-item:hover {
  background-color: #f8f9fa;
  margin: 0 -20px;
  padding-left: 20px;
  padding-right: 20px;
  border-radius: 8px;
}

.task-item:last-child {
  border-bottom: none;
}

.task-priority {
  width: 4px;
  height: 40px;
  border-radius: 2px;
  margin-right: 12px;
  background-color: #e4e7ed;
  flex-shrink: 0;
}

.task-priority.high {
  background-color: #f56c6c;
}

.task-priority.medium,
.task-priority.normal {
  background-color: #409eff;
}

.task-priority.low {
  background-color: #67c23a;
}

.task-content {
  flex: 1;
  min-width: 0;
}

.task-title {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.task-desc {
  color: #666;
  font-size: 13px;
  margin-bottom: 8px;
}

.task-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.task-time {
  color: #999;
  font-size: 12px;
}

.task-action {
  color: #c0c4cc;
  font-size: 16px;
  flex-shrink: 0;
}

/* 移动端适配补充 */
@media (max-width: 768px) {
  .task-content {
    margin-right: 8px;
  }
  
  .task-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
}

/* 卡片标题样式 */
:deep(.el-card__header) {
  padding: 0;
  border-bottom: none;
}

:deep(.el-card__body) {
  padding: 0 0 16px 0;
}

/* 标签样式 */
:deep(.el-tag) {
  border: none;
}

/* 徽章样式 */
:deep(.el-badge__content) {
  border: 1px solid #fff;
}
</style> 