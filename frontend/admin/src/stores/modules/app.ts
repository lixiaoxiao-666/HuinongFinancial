import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAppStore = defineStore('app', () => {
  // 侧边栏状态
  const sidebarCollapsed = ref(false)
  const sidebarMobile = ref(false)

  // 主题设置
  const theme = ref<'light' | 'dark'>('light')

  // 页面加载状态
  const pageLoading = ref(false)

  // 全局搜索状态
  const globalSearchVisible = ref(false)

  // 面包屑导航
  const breadcrumbs = ref<Array<{
    title: string
    path?: string
  }>>([])

  // 标签页状态
  const tabs = ref<Array<{
    key: string
    title: string
    path: string
    closable: boolean
  }>>([])
  const activeTabKey = ref<string>('')

  // 计算属性
  const isMobile = computed(() => {
    return window.innerWidth < 768
  })

  const sidebarWidth = computed(() => {
    if (sidebarMobile.value) return 0
    return sidebarCollapsed.value ? 80 : 256
  })

  // 切换侧边栏
  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value
    localStorage.setItem('oa_sidebar_collapsed', String(sidebarCollapsed.value))
  }

  // 设置移动端侧边栏
  const setSidebarMobile = (mobile: boolean) => {
    sidebarMobile.value = mobile
  }

  // 切换主题
  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
    localStorage.setItem('oa_theme', theme.value)
    applyTheme(theme.value)
  }

  // 应用主题
  const applyTheme = (themeMode: 'light' | 'dark') => {
    theme.value = themeMode
    document.documentElement.setAttribute('data-theme', themeMode)
  }

  // 设置页面加载状态
  const setPageLoading = (loading: boolean) => {
    pageLoading.value = loading
  }

  // 切换全局搜索
  const toggleGlobalSearch = () => {
    globalSearchVisible.value = !globalSearchVisible.value
  }

  // 设置面包屑
  const setBreadcrumbs = (crumbs: Array<{ title: string; path?: string }>) => {
    breadcrumbs.value = crumbs
  }

  // 添加标签页
  const addTab = (tab: {
    key: string
    title: string
    path: string
    closable?: boolean
  }) => {
    const existingTab = tabs.value.find(t => t.key === tab.key)
    if (!existingTab) {
      tabs.value.push({
        ...tab,
        closable: tab.closable !== false
      })
    }
    activeTabKey.value = tab.key
  }

  // 移除标签页
  const removeTab = (key: string) => {
    const index = tabs.value.findIndex(t => t.key === key)
    if (index > -1) {
      tabs.value.splice(index, 1)
      
      // 如果移除的是当前激活标签，需要切换到其他标签
      if (activeTabKey.value === key) {
        if (tabs.value.length > 0) {
          const newIndex = Math.min(index, tabs.value.length - 1)
          activeTabKey.value = tabs.value[newIndex].key
        } else {
          activeTabKey.value = ''
        }
      }
    }
  }

  // 设置激活标签
  const setActiveTab = (key: string) => {
    activeTabKey.value = key
  }

  // 关闭其他标签页
  const closeOtherTabs = (currentKey: string) => {
    tabs.value = tabs.value.filter(tab => 
      tab.key === currentKey || !tab.closable
    )
    activeTabKey.value = currentKey
  }

  // 关闭所有标签页
  const closeAllTabs = () => {
    tabs.value = tabs.value.filter(tab => !tab.closable)
    activeTabKey.value = tabs.value.length > 0 ? tabs.value[0].key : ''
  }

  // 初始化应用设置
  const initializeApp = () => {
    // 恢复侧边栏状态
    const savedCollapsed = localStorage.getItem('oa_sidebar_collapsed')
    if (savedCollapsed !== null) {
      sidebarCollapsed.value = savedCollapsed === 'true'
    }

    // 恢复主题设置
    const savedTheme = localStorage.getItem('oa_theme') as 'light' | 'dark'
    if (savedTheme) {
      applyTheme(savedTheme)
    }

    // 监听窗口大小变化
    const handleResize = () => {
      setSidebarMobile(window.innerWidth < 768)
      if (window.innerWidth < 768) {
        sidebarCollapsed.value = true
      }
    }
    
    window.addEventListener('resize', handleResize)
    handleResize() // 初始检查

    // 添加默认首页标签
    addTab({
      key: 'dashboard',
      title: '工作台',
      path: '/dashboard',
      closable: false
    })
  }

  return {
    // 状态
    sidebarCollapsed,
    sidebarMobile,
    theme,
    pageLoading,
    globalSearchVisible,
    breadcrumbs,
    tabs,
    activeTabKey,

    // 计算属性
    isMobile,
    sidebarWidth,

    // 方法
    toggleSidebar,
    setSidebarMobile,
    toggleTheme,
    applyTheme,
    setPageLoading,
    toggleGlobalSearch,
    setBreadcrumbs,
    addTab,
    removeTab,
    setActiveTab,
    closeOtherTabs,
    closeAllTabs,
    initializeApp
  }
}) 