import { defineStore } from 'pinia'
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

export const useAIStatusStore = defineStore('aiStatus', () => {
  // AI审批全局状态 - 默认为暂停
  const globalAIStatus = ref(false)
  
  // 切换AI审批状态
  const toggleAIStatus = () => {
    globalAIStatus.value = !globalAIStatus.value
    
    if (globalAIStatus.value) {
      ElMessage.success('AI审批已启动，所有流程恢复运行')
    } else {
      ElMessage.warning('AI审批已暂停，所有流程停止处理')
    }
    
    return globalAIStatus.value
  }
  
  // 设置AI审批状态
  const setAIStatus = (status: boolean) => {
    globalAIStatus.value = status
  }
  
  // 获取AI审批状态
  const getAIStatus = () => {
    return globalAIStatus.value
  }
  
  return {
    globalAIStatus,
    toggleAIStatus,
    setAIStatus,
    getAIStatus
  }
}) 