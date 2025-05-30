<template>
  <div class="sider-menu">
    <a-menu
      v-model:selectedKeys="selectedKeys"
      v-model:openKeys="openKeys"
      mode="inline"
      theme="dark"
      :inline-collapsed="collapsed"
      :style="{ border: 'none', background: 'transparent' }"
      @click="handleMenuClick"
    >
      <!-- 工作台 -->
      <a-menu-item key="/dashboard">
        <template #icon>
          <DashboardOutlined />
        </template>
        <span>工作台</span>
      </a-menu-item>

      <!-- 贷款管理 -->
      <a-sub-menu key="loan" :title="collapsed ? '' : '贷款管理'">
        <template #icon>
          <BankOutlined />
        </template>
        <a-menu-item key="/loan/applications">
          <FileTextOutlined />
          <span>申请审批</span>
        </a-menu-item>
        <a-menu-item key="/loan/products">
          <ShoppingOutlined />
          <span>产品管理</span>
        </a-menu-item>
        <a-menu-item key="/loan/contracts">
          <FileProtectOutlined />
          <span>合同管理</span>
        </a-menu-item>
        <a-menu-item key="/loan/statistics">
          <BarChartOutlined />
          <span>数据统计</span>
        </a-menu-item>
      </a-sub-menu>

      <!-- 农机管理 -->
      <a-sub-menu key="machine" :title="collapsed ? '' : '农机管理'">
        <template #icon>
          <CarOutlined />
        </template>
        <a-menu-item key="/machine/inventory">
          <DatabaseOutlined />
          <span>设备库存</span>
        </a-menu-item>
        <a-menu-item key="/machine/rentals">
          <CalendarOutlined />
          <span>租赁订单</span>
        </a-menu-item>
        <a-menu-item key="/machine/maintenance">
          <ToolOutlined />
          <span>维护保养</span>
        </a-menu-item>
      </a-sub-menu>

      <!-- 用户管理 -->
      <a-sub-menu key="user" :title="collapsed ? '' : '用户管理'">
        <template #icon>
          <UserOutlined />
        </template>
        <a-menu-item key="/user/list">
          <TeamOutlined />
          <span>用户列表</span>
        </a-menu-item>
        <a-menu-item key="/user/verification">
          <SafetyCertificateOutlined />
          <span>实名认证</span>
        </a-menu-item>
        <a-menu-item key="/user/tags">
          <TagsOutlined />
          <span>用户标签</span>
        </a-menu-item>
      </a-sub-menu>

      <!-- 内容管理 -->
      <a-sub-menu key="content" :title="collapsed ? '' : '内容管理'">
        <template #icon>
          <FileOutlined />
        </template>
        <a-menu-item key="/content/articles">
          <ReadOutlined />
          <span>资讯文章</span>
        </a-menu-item>
        <a-menu-item key="/content/policies">
          <BookOutlined />
          <span>政策信息</span>
        </a-menu-item>
        <a-menu-item key="/content/experts">
          <UserOutlined />
          <span>专家管理</span>
        </a-menu-item>
        <a-menu-item key="/content/notifications">
          <BellOutlined />
          <span>通知管理</span>
        </a-menu-item>
      </a-sub-menu>

      <!-- AI 工作流 -->
      <a-sub-menu key="ai" :title="collapsed ? '' : 'AI 工作流'">
        <template #icon>
          <RobotOutlined />
        </template>
        <a-menu-item key="/ai/workflows">
          <ForkOutlined />
          <span>工作流管理</span>
        </a-menu-item>
        <a-menu-item key="/ai/logs">
          <FileSearchOutlined />
          <span>调用日志</span>
        </a-menu-item>
        <a-menu-item key="/ai/settings">
          <SettingOutlined />
          <span>AI 配置</span>
        </a-menu-item>
      </a-sub-menu>

      <!-- 系统管理 -->
      <a-sub-menu key="system" :title="collapsed ? '' : '系统管理'" v-if="hasAdminRole">
        <template #icon>
          <SettingOutlined />
        </template>
        <a-menu-item key="/system/users">
          <UserOutlined />
          <span>用户管理</span>
        </a-menu-item>
        <a-menu-item key="/system/roles">
          <SafetyCertificateOutlined />
          <span>角色权限</span>
        </a-menu-item>
        <a-menu-item key="/system/config">
          <ControlOutlined />
          <span>系统配置</span>
        </a-menu-item>
        <a-menu-item key="/system/logs">
          <FileTextOutlined />
          <span>操作日志</span>
        </a-menu-item>
      </a-sub-menu>

      <!-- 分割线 -->
      <a-menu-divider />

      <!-- 帮助中心 -->
      <a-menu-item key="/help">
        <template #icon>
          <QuestionCircleOutlined />
        </template>
        <span>帮助中心</span>
      </a-menu-item>
    </a-menu>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/modules/auth'
import {
  DashboardOutlined,
  BankOutlined,
  CarOutlined,
  UserOutlined,
  FileOutlined,
  RobotOutlined,
  SettingOutlined,
  QuestionCircleOutlined,
  FileTextOutlined,
  ShoppingOutlined,
  FileProtectOutlined,
  BarChartOutlined,
  DatabaseOutlined,
  CalendarOutlined,
  ToolOutlined,
  TeamOutlined,
  SafetyCertificateOutlined,
  TagsOutlined,
  ReadOutlined,
  BookOutlined,
  BellOutlined,
  ForkOutlined,
  FileSearchOutlined,
  ControlOutlined
} from '@ant-design/icons-vue'

interface Props {
  collapsed: boolean
}

const props = defineProps<Props>()

/**
 * 组件状态
 */
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const selectedKeys = ref<string[]>([])
const openKeys = ref<string[]>([])

/**
 * 权限检查
 */
const hasAdminRole = computed(() => {
  return authStore.hasRole(['admin', 'super_admin'])
})

/**
 * 根据当前路由设置选中的菜单项
 */
const updateSelectedKeys = () => {
  const path = route.path
  selectedKeys.value = [path]
  
  // 设置展开的子菜单
  if (path.startsWith('/loan')) {
    openKeys.value = ['loan']
  } else if (path.startsWith('/machine')) {
    openKeys.value = ['machine']
  } else if (path.startsWith('/user')) {
    openKeys.value = ['user']
  } else if (path.startsWith('/content')) {
    openKeys.value = ['content']
  } else if (path.startsWith('/ai')) {
    openKeys.value = ['ai']
  } else if (path.startsWith('/system')) {
    openKeys.value = ['system']
  }
}

/**
 * 处理菜单点击
 */
const handleMenuClick = ({ key }: { key: string }) => {
  if (key !== route.path) {
    router.push(key)
  }
}

/**
 * 监听路由变化
 */
watch(route, updateSelectedKeys, { immediate: true })

/**
 * 监听折叠状态变化
 */
watch(() => props.collapsed, (collapsed) => {
  if (collapsed) {
    openKeys.value = []
  } else {
    updateSelectedKeys()
  }
})
</script>

<style lang="scss" scoped>
.sider-menu {
  height: calc(100vh - 64px);
  overflow-y: auto;
  overflow-x: hidden;
  
  // 自定义滚动条
  &::-webkit-scrollbar {
    width: 4px;
  }
  
  &::-webkit-scrollbar-track {
    background: transparent;
  }
  
  &::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.2);
    border-radius: 2px;
    
    &:hover {
      background: rgba(255, 255, 255, 0.3);
    }
  }
  
  :deep(.ant-menu) {
    // 菜单项样式
    .ant-menu-item {
      height: 48px;
      line-height: 48px;
      margin: 4px 8px;
      border-radius: 6px;
      overflow: hidden;
      
      &.ant-menu-item-selected {
        background: #1890ff !important;
        color: #fff;
        
        &::after {
          display: none;
        }
      }
      
      &:hover {
        background: rgba(255, 255, 255, 0.1) !important;
        color: #fff;
      }
      
      .anticon {
        font-size: 16px;
        margin-right: 12px;
      }
    }
    
    // 子菜单样式
    .ant-menu-submenu {
      margin: 4px 8px;
      border-radius: 6px;
      overflow: hidden;
      
      .ant-menu-submenu-title {
        height: 48px;
        line-height: 48px;
        margin: 0;
        border-radius: 6px;
        
        &:hover {
          background: rgba(255, 255, 255, 0.1) !important;
          color: #fff;
        }
        
        .anticon {
          font-size: 16px;
          margin-right: 12px;
        }
      }
      
      .ant-menu-sub {
        background: rgba(0, 0, 0, 0.2) !important;
        
        .ant-menu-item {
          height: 40px;
          line-height: 40px;
          margin: 2px 4px;
          padding-left: 48px !important;
          
          .anticon {
            font-size: 14px;
            margin-right: 8px;
          }
        }
      }
    }
    
    // 分割线样式
    .ant-menu-item-divider {
      background: rgba(255, 255, 255, 0.15);
      margin: 16px 24px;
    }
    
    // 折叠状态样式调整
    &.ant-menu-inline-collapsed {
      .ant-menu-item,
      .ant-menu-submenu-title {
        padding: 0 20px !important;
        text-align: center;
        
        .anticon {
          margin-right: 0;
        }
      }
      
      .ant-menu-submenu-arrow {
        display: none;
      }
    }
  }
}
</style> 