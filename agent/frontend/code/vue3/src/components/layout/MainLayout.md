# MainLayout 组件文档

## 概述

`MainLayout` 是整个管理系统的主要布局组件，采用**左侧边栏 + 右侧布局（头部-内容-底部）**结构。该组件负责整体页面布局管理、侧边栏控制以及内容区域的滚动优化。

## 组件结构

```
MainLayout
├── Sider (左侧边栏)
│   ├── Logo区域
│   └── SiderMenu (导航菜单)
└── Layout (右侧布局)
    ├── Header (头部)
    ├── Content (内容区)
    └── Footer (底部)
```

## 主要功能

### 1. 布局管理
- **左侧固定边栏**：使用 `position: fixed` 完全占满左侧区域
- **右侧自适应布局**：通过 CSS 变量 `--sider-width` 动态调整边距
- **响应式设计**：移动端下侧边栏变为覆盖层模式

### 2. 侧边栏控制
- **折叠/展开**：支持侧边栏折叠，宽度在 240px 和 80px 之间切换
- **底部空白修复**：侧边栏高度设为 `100vh`，完全占满视口高度
- **滚动优化**：菜单区域支持独立滚动，带自定义滚动条样式

### 3. 内容区域优化
- **取消面包屑导航**：简化界面，直接展示页面内容
- **滚动栏实现**：内容区域支持垂直滚动，带自定义滚动条样式
- **高度计算**：内容区最大高度为 `calc(100vh - 128px)`，减去头部和底部高度

## 样式特性

### 侧边栏样式
```scss
.layout-sider {
  position: fixed;
  left: 0;
  top: 0;
  height: 100vh;
  z-index: 100;
  
  .sider-logo {
    flex-shrink: 0; // 防止Logo区域被压缩
  }
  
  .sider-menu-wrapper {
    flex: 1;
    overflow-y: auto; // 菜单区域独立滚动
  }
}
```

### 内容区样式
```scss
.layout-content {
  overflow: auto;
  max-height: calc(100vh - 128px);
  
  .page-content {
    padding: 16px;
    overflow-y: auto;
    min-height: calc(100vh - 128px);
  }
}
```

### 自定义滚动条
- **侧边栏滚动条**：白色半透明样式，适配深色背景
- **内容区滚动条**：灰色半透明样式，适配浅色背景
- **响应式交互**：hover 状态下滚动条变粗，提升可用性

## 响应式适配

### 桌面端 (> 768px)
- 左侧边栏固定定位
- 右侧布局使用 `margin-left` 避让
- 内容区域支持独立滚动

### 移动端 (≤ 768px)
- 侧边栏变为覆盖层模式
- 默认状态下隐藏在屏幕左侧
- 展开时覆盖在内容上方

## 性能优化

### 1. 过渡动画
- 使用 `cubic-bezier(0.4, 0, 0.2, 1)` 缓动函数
- 侧边栏折叠动画流畅自然
- 页面切换支持淡入淡出效果

### 2. 滚动优化
- 使用 `scrollbar-width: thin` 优化 Firefox
- WebKit 内核浏览器使用自定义滚动条
- 滚动区域独立，避免整页滚动

### 3. 样式隔离
- 使用 CSS 变量管理动态样式
- 全局滚动条样式统一管理
- 暗色主题适配完善

## 状态管理

### 响应式状态
```typescript
const siderCollapsed = ref(false) // 侧边栏折叠状态
```

### 计算属性
```typescript
const siderStyle = computed(() => ({
  height: '100vh',
  position: 'fixed' as const,
  zIndex: 100
}))

const contentStyle = computed(() => ({
  overflow: 'auto',
  maxHeight: 'calc(100vh - 128px)'
}))
```

## 使用方式

```vue
<template>
  <MainLayout>
    <!-- 页面内容通过 router-view 渲染 -->
  </MainLayout>
</template>
```

## 注意事项

1. **侧边栏高度**：现在使用 `100vh` 完全占满视口，修复了底部空白问题
2. **内容区滚动**：移除了面包屑导航，内容区域直接支持滚动显示
3. **z-index 层级**：侧边栏 z-index 为 100，确保在其他元素之上
4. **CSS 变量**：通过 `--sider-width` 变量控制右侧布局的位置调整

## 相关组件

- `SiderMenu` - 侧边栏导航菜单
- `HeaderContent` - 头部内容组件
- `FooterContent` - 底部内容组件

## 更新日志

- **v1.0.0**: 初始版本，基础布局实现
- **v1.1.0**: 新增响应式设计和暗色主题支持
- **v1.2.0**: 优化滚动条样式和过渡动画
- **v1.3.0**: 删除面包屑导航，优化内容区滚动，修复侧边栏底部空白 