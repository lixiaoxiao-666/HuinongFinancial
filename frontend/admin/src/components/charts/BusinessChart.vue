<template>
  <div class="business-chart">
    <div class="chart-header" v-if="title || subtitle">
      <h4 class="chart-title" v-if="title">{{ title }}</h4>
      <div class="chart-subtitle" v-if="subtitle">{{ subtitle }}</div>
    </div>
    
    <div class="chart-container" ref="chartContainer">
      <div 
        v-if="loading" 
        class="chart-loading"
      >
        <a-spin size="large" />
        <p>加载图表数据中...</p>
      </div>
      
      <div 
        v-else-if="hasError" 
        class="chart-error"
      >
        <div class="error-icon">📊</div>
        <p>图表加载失败</p>
        <a-button type="link" @click="$emit('retry')">重试</a-button>
      </div>
      
      <div 
        v-else
        ref="chartRef"
        :style="{ height: height + 'px', width: '100%' }"
      ></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import * as echarts from 'echarts'

interface Props {
  title?: string
  subtitle?: string
  type: 'line' | 'bar' | 'pie' | 'scatter' | 'radar' | 'gauge'
  data: any[]
  loading?: boolean
  error?: string
  height?: number | string
  width?: number | string
  options?: any
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  height: 400,
  width: '100%'
})

interface Emits {
  click: [data: any]
  retry: []
}

const emit = defineEmits<Emits>()

const chartContainer = ref<HTMLElement>()
const chartRef = ref<HTMLElement>()
const chartInstance = ref<echarts.ECharts>()
const hasError = ref(false)

// 图表颜色配置
const chartColors = [
  '#52C41A',  // 主绿
  '#1890FF',  // 蓝色
  '#FAAD14',  // 黄色
  '#F5222D',  // 红色
  '#722ED1',  // 紫色
  '#13C2C2',  // 青色
  '#FA8C16',  // 橙色
  '#A0D911'   // 柠檬绿
]

/**
 * 图表配置
 */
const getChartOption = () => {
  if (!props.data || props.data.length === 0) {
    return {}
  }

  const baseOption = {
    color: chartColors,
    backgroundColor: 'transparent',
    textStyle: {
      fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif'
    }
  }

  switch (props.type) {
    case 'line':
      return {
        ...baseOption,
        grid: {
          top: 40,
          left: 50,
          right: 30,
          bottom: 50
        },
        tooltip: {
          trigger: 'axis',
          backgroundColor: 'rgba(50, 50, 50, 0.9)',
          borderColor: 'transparent',
          textStyle: {
            color: '#fff'
          }
        },
        xAxis: {
          type: 'category',
          data: props.data.map(item => item.name),
          axisLine: { lineStyle: { color: '#e8e8e8' } },
          axisLabel: { color: '#8c8c8c' }
        },
        yAxis: {
          type: 'value',
          axisLine: { lineStyle: { color: '#e8e8e8' } },
          axisLabel: { color: '#8c8c8c' },
          splitLine: { lineStyle: { color: '#f5f5f5' } }
        },
        series: [{
          data: props.data.map(item => item.value),
          type: 'line',
          smooth: true,
          lineStyle: { color: chartColors[1] },
          itemStyle: { color: chartColors[1] },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0, y: 0, x2: 0, y2: 1,
              colorStops: [
                { offset: 0, color: 'rgba(24, 144, 255, 0.3)' },
                { offset: 1, color: 'rgba(24, 144, 255, 0.1)' }
              ]
            }
          }
        }]
      }

    case 'bar':
      return {
        ...baseOption,
        grid: {
          top: 40,
          left: 50,
          right: 30,
          bottom: 50
        },
        tooltip: {
          trigger: 'axis',
          backgroundColor: 'rgba(50, 50, 50, 0.9)',
          borderColor: 'transparent',
          textStyle: {
            color: '#fff'
          }
        },
        xAxis: {
          type: 'category',
          data: props.data.map(item => item.name),
          axisLine: { lineStyle: { color: '#e8e8e8' } },
          axisLabel: { color: '#8c8c8c' }
        },
        yAxis: {
          type: 'value',
          axisLine: { lineStyle: { color: '#e8e8e8' } },
          axisLabel: { color: '#8c8c8c' },
          splitLine: { lineStyle: { color: '#f5f5f5' } }
        },
        series: [{
          data: props.data.map(item => item.value),
          type: 'bar',
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: chartColors[1] },
              { offset: 1, color: chartColors[1] + '80' }
            ])
          },
          barWidth: '60%'
        }]
      }

    case 'pie':
      return {
        ...baseOption,
        tooltip: {
          trigger: 'item',
          formatter: '{a} <br/>{b}: {c} ({d}%)'
        },
        legend: {
          orient: 'vertical',
          left: 'right',
          textStyle: { color: '#8c8c8c' }
        },
        series: [{
          name: '数据分布',
          type: 'pie',
          radius: ['40%', '70%'],
          center: ['40%', '50%'],
          data: props.data,
          itemStyle: {
            borderRadius: 4,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: false
          },
          labelLine: {
            show: false
          }
        }]
      }

    case 'gauge':
      const gaugeData = Array.isArray(props.data) ? props.data[0] : props.data
      return {
        ...baseOption,
        series: [{
          name: '评分',
          type: 'gauge',
          center: ['50%', '60%'],
          startAngle: 200,
          endAngle: -20,
          min: 0,
          max: gaugeData?.max || 100,
          splitNumber: 10,
          itemStyle: {
            color: chartColors[0]
          },
          progress: {
            show: true,
            width: 30
          },
          pointer: {
            show: false
          },
          axisLine: {
            lineStyle: {
              width: 30
            }
          },
          axisTick: {
            distance: -45,
            splitNumber: 5,
            lineStyle: {
              width: 2,
              color: '#999'
            }
          },
          splitLine: {
            distance: -52,
            length: 14,
            lineStyle: {
              width: 3,
              color: '#999'
            }
          },
          axisLabel: {
            distance: -20,
            color: '#999',
            fontSize: 12
          },
          anchor: {
            show: false
          },
          title: {
            show: false
          },
          detail: {
            valueAnimation: true,
            width: '60%',
            lineHeight: 40,
            borderRadius: 8,
            offsetCenter: [0, '-15%'],
            fontSize: 24,
            fontWeight: 'bolder',
            formatter: '{value}' + (gaugeData?.unit || ''),
            color: 'inherit'
          },
          data: [gaugeData]
        }]
      }

    default:
      return baseOption
  }
}

/**
 * 初始化图表
 */
const initChart = async () => {
  if (!chartRef.value) return

  try {
    hasError.value = false
    
    // 销毁已存在的图表实例
    if (chartInstance.value) {
      chartInstance.value.dispose()
    }

    // 创建新的图表实例
    chartInstance.value = echarts.init(chartRef.value)
    
    const option = getChartOption()
    if (option && Object.keys(option).length > 0) {
      chartInstance.value.setOption(option)
    }
    
    // 响应式调整
    const resizeChart = () => {
      chartInstance.value?.resize()
    }
    
    window.addEventListener('resize', resizeChart)
    
    // 清理事件监听器
    const cleanup = () => {
      window.removeEventListener('resize', resizeChart)
    }
    
    chartRef.value.addEventListener('beforeunload', cleanup)

  } catch (error) {
    console.error('图表初始化失败:', error)
    hasError.value = true
  }
}

/**
 * 更新图表
 */
const updateChart = () => {
  if (!chartInstance.value) return
  
  try {
    const option = getChartOption()
    if (option && Object.keys(option).length > 0) {
      chartInstance.value.setOption(option, true)
    }
  } catch (error) {
    console.error('图表更新失败:', error)
    hasError.value = true
  }
}

// 处理图表点击事件
const handleChartClick = (params: any) => {
  emit('click', params)
}

// 监听数据变化，重新渲染图表
watch(() => props.data, () => {
  nextTick(() => {
    if (props.data && props.data.length > 0) {
      if (chartInstance.value) {
        updateChart()
      } else {
        initChart()
      }
    }
  })
}, { deep: true, immediate: true })

onMounted(() => {
  nextTick(() => {
    if (props.data && props.data.length > 0) {
      initChart()
    }
  })
})

onBeforeUnmount(() => {
  if (chartInstance.value) {
    chartInstance.value.dispose?.()
  }
})
</script>

<style lang="scss" scoped>
.business-chart {
  .chart-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 1px solid #f0f0f0;
    
    .chart-title {
      margin: 0;
      font-size: 16px;
      font-weight: 600;
      color: #262626;
      line-height: 1.4;
    }
    
    .chart-subtitle {
      font-size: 12px;
      color: #8c8c8c;
      margin-top: 4px;
    }
  }
  
  .chart-container {
    position: relative;
    min-height: 300px;
    
    .chart-loading {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 300px;
      color: #8c8c8c;
      
      p {
        margin-top: 16px;
        margin-bottom: 0;
      }
    }
    
    .chart-error {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 300px;
      color: #8c8c8c;
      
      .error-icon {
        font-size: 48px;
        margin-bottom: 16px;
        opacity: 0.5;
      }
      
      p {
        margin-bottom: 16px;
      }
    }
  }
}

// 暗色主题适配
:deep([data-theme="dark"]) {
  .chart-header {
    border-bottom-color: #303030;
    
    .chart-title {
      color: rgba(255, 255, 255, 0.85);
    }
    
    .chart-subtitle {
      color: rgba(255, 255, 255, 0.45);
    }
  }
  
  .chart-loading,
  .chart-error {
    color: rgba(255, 255, 255, 0.45);
  }
}
</style> 