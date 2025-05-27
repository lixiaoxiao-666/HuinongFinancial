# UI/UX设计规范文档

## 1. 设计概述

数字惠农APP及OA后台管理系统UI/UX设计规范，基于现代化、简洁、易用的设计理念，体现农业特色，为农户和管理员提供优质的用户体验。

### 1.1 设计目标
- **易用性**: 降低农户使用门槛，简化操作流程
- **专业性**: 管理端体现专业性和高效性
- **可信度**: 通过设计增强用户对金融服务的信任
- **一致性**: 保持平台设计的统一性和连贯性

### 1.2 目标用户
- **农户用户**: 30-55岁，中低学历，智能手机操作经验有限
- **审批员**: 25-45岁，大专以上学历，熟练使用办公软件
- **管理员**: 30-50岁，具备技术背景，需要高效的操作界面

## 2. 色彩系统

### 2.1 主色调
```css
/* 绿色系 - 体现农业特色 */
--primary-50: #f0fdf4
--primary-100: #dcfce7
--primary-200: #bbf7d0
--primary-300: #86efac
--primary-400: #4ade80
--primary-500: #22c55e  /* 主要绿色 */
--primary-600: #16a34a
--primary-700: #15803d
--primary-800: #166534
--primary-900: #14532d
```

### 2.2 辅助色彩
```css
/* 蓝色系 - 信任感 */
--secondary-500: #3b82f6
--secondary-600: #2563eb

/* 警告色 */
--warning-500: #f59e0b
--warning-600: #d97706

/* 错误色 */
--error-500: #ef4444
--error-600: #dc2626

/* 成功色 */
--success-500: #10b981
--success-600: #059669
```

### 2.3 中性色调
```css
/* 灰色系 */
--gray-50: #f9fafb
--gray-100: #f3f4f6
--gray-200: #e5e7eb
--gray-300: #d1d5db
--gray-400: #9ca3af
--gray-500: #6b7280
--gray-600: #4b5563
--gray-700: #374151
--gray-800: #1f2937
--gray-900: #111827
```

### 2.4 语义化颜色
```css
/* 背景色 */
--bg-primary: #ffffff
--bg-secondary: #f9fafb
--bg-tertiary: #f3f4f6

/* 文本色 */
--text-primary: #111827
--text-secondary: #6b7280
--text-tertiary: #9ca3af

/* 边框色 */
--border-light: #e5e7eb
--border-medium: #d1d5db
--border-dark: #9ca3af
```

## 3. 字体系统

### 3.1 字体族
```css
/* 主要字体 */
font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, 
             "Helvetica Neue", Arial, "Noto Sans", sans-serif, 
             "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", 
             "Noto Color Emoji";

/* 中文字体备选 */
font-family: "PingFang SC", "Microsoft YaHei", "Source Han Sans CN", 
             "Noto Sans CJK SC", sans-serif;
```

### 3.2 字体大小
```css
/* 标题 */
--text-xs: 0.75rem;     /* 12px */
--text-sm: 0.875rem;    /* 14px */
--text-base: 1rem;      /* 16px */
--text-lg: 1.125rem;    /* 18px */
--text-xl: 1.25rem;     /* 20px */
--text-2xl: 1.5rem;     /* 24px */
--text-3xl: 1.875rem;   /* 30px */
--text-4xl: 2.25rem;    /* 36px */
```

### 3.3 字重
```css
--font-light: 300;
--font-normal: 400;
--font-medium: 500;
--font-semibold: 600;
--font-bold: 700;
```

### 3.4 行高
```css
--leading-tight: 1.25;
--leading-snug: 1.375;
--leading-normal: 1.5;
--leading-relaxed: 1.625;
--leading-loose: 2;
```

## 4. 间距系统

### 4.1 基础间距
```css
--spacing-0: 0px;
--spacing-1: 0.25rem;   /* 4px */
--spacing-2: 0.5rem;    /* 8px */
--spacing-3: 0.75rem;   /* 12px */
--spacing-4: 1rem;      /* 16px */
--spacing-5: 1.25rem;   /* 20px */
--spacing-6: 1.5rem;    /* 24px */
--spacing-8: 2rem;      /* 32px */
--spacing-10: 2.5rem;   /* 40px */
--spacing-12: 3rem;     /* 48px */
--spacing-16: 4rem;     /* 64px */
--spacing-20: 5rem;     /* 80px */
```

### 4.2 组件间距
- **卡片内边距**: 24px (spacing-6)
- **表单字段间距**: 16px (spacing-4)
- **按钮内边距**: 8px 16px (spacing-2 spacing-4)
- **页面边距**: 16px (spacing-4) 手机端, 24px (spacing-6) 桌面端

## 5. 组件设计规范

### 5.1 按钮规范
```html
<!-- 主要按钮 -->
<button class="bg-green-600 hover:bg-green-700 active:bg-green-800 
               text-white font-medium py-2 px-4 rounded-lg 
               transition-all duration-200 
               focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-opacity-50
               disabled:opacity-50 disabled:cursor-not-allowed">
    主要操作
</button>

<!-- 次要按钮 -->
<button class="bg-white border border-gray-300 hover:bg-gray-50 active:bg-gray-100
               text-gray-700 font-medium py-2 px-4 rounded-lg
               transition-all duration-200
               focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-opacity-50">
    次要操作
</button>

<!-- 文本按钮 -->
<button class="text-green-600 hover:text-green-700 active:text-green-800
               font-medium py-2 px-4 rounded-lg
               transition-colors duration-200
               focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-opacity-50">
    文本按钮
</button>
```

### 5.2 表单控件规范
```html
<!-- 输入框 -->
<div class="relative">
    <input class="w-full px-3 py-2 border border-gray-300 rounded-lg
                  placeholder-gray-400 text-gray-900
                  focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent
                  disabled:bg-gray-50 disabled:text-gray-500"
           type="text" placeholder="请输入内容">
    <!-- 错误状态 -->
    <input class="w-full px-3 py-2 border border-red-300 rounded-lg
                  placeholder-gray-400 text-gray-900
                  focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent
                  bg-red-50"
           type="text" placeholder="请输入内容">
</div>

<!-- 选择框 -->
<select class="w-full px-3 py-2 border border-gray-300 rounded-lg
               text-gray-900 bg-white
               focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent">
    <option value="">请选择</option>
    <option value="1">选项1</option>
</select>
```

### 5.3 卡片组件规范
```html
<!-- 基础卡片 -->
<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">卡片标题</h3>
    <p class="text-gray-600">卡片内容</p>
</div>

<!-- 交互卡片 -->
<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6
            hover:shadow-md hover:border-gray-300
            transition-all duration-200 cursor-pointer">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">可点击卡片</h3>
    <p class="text-gray-600">卡片内容</p>
</div>
```

### 5.4 状态指示器
```html
<!-- 成功状态 -->
<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
             bg-green-100 text-green-800">
    <div class="w-1.5 h-1.5 bg-green-400 rounded-full mr-1.5"></div>
    审批通过
</span>

<!-- 处理中状态 -->
<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
             bg-yellow-100 text-yellow-800">
    <div class="w-1.5 h-1.5 bg-yellow-400 rounded-full mr-1.5"></div>
    审批中
</span>

<!-- 失败状态 -->
<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
             bg-red-100 text-red-800">
    <div class="w-1.5 h-1.5 bg-red-400 rounded-full mr-1.5"></div>
    审批拒绝
</span>
```

## 6. 图标规范

### 6.1 图标来源
- **FontAwesome 6.4.0**: 主要图标库
- **自定义图标**: 农业相关特色图标

### 6.2 图标尺寸
```css
/* 图标大小 */
--icon-xs: 12px;
--icon-sm: 16px;
--icon-base: 20px;
--icon-lg: 24px;
--icon-xl: 32px;
--icon-2xl: 48px;
```

### 6.3 常用图标
```html
<!-- 功能图标 -->
<i class="fas fa-home"></i>          <!-- 首页 -->
<i class="fas fa-user"></i>          <!-- 用户 -->
<i class="fas fa-credit-card"></i>   <!-- 贷款 -->
<i class="fas fa-tractor"></i>       <!-- 农机 -->
<i class="fas fa-chart-bar"></i>     <!-- 统计 -->
<i class="fas fa-cog"></i>           <!-- 设置 -->

<!-- 状态图标 -->
<i class="fas fa-check-circle text-green-500"></i>  <!-- 成功 -->
<i class="fas fa-clock text-yellow-500"></i>        <!-- 等待 -->
<i class="fas fa-times-circle text-red-500"></i>    <!-- 失败 -->
<i class="fas fa-info-circle text-blue-500"></i>    <!-- 信息 -->
```

## 7. 布局规范

### 7.1 用户端布局
```html
<!-- 移动端主布局 -->
<div class="min-h-screen bg-gray-50 pb-20">
    <!-- 顶部导航 -->
    <header class="bg-green-600 text-white px-4 py-3 sticky top-0 z-50">
        <div class="flex items-center justify-between">
            <h1 class="text-lg font-semibold">页面标题</h1>
            <button class="p-2">
                <i class="fas fa-bell"></i>
            </button>
        </div>
    </header>
    
    <!-- 主要内容 -->
    <main class="px-4 py-6">
        <!-- 页面内容 -->
    </main>
    
    <!-- 底部导航 -->
    <nav class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 z-50">
        <div class="grid grid-cols-4 gap-1">
            <!-- 导航项 -->
        </div>
    </nav>
</div>
```

### 7.2 管理端布局
```html
<!-- 桌面端主布局 -->
<div class="flex h-screen bg-gray-100">
    <!-- 侧边栏 -->
    <aside class="w-64 bg-gray-800 text-white">
        <!-- 侧边栏内容 -->
    </aside>
    
    <!-- 主内容区域 -->
    <div class="flex-1 flex flex-col">
        <!-- 顶部栏 -->
        <header class="bg-white border-b border-gray-200 px-6 py-4">
            <!-- 顶部栏内容 -->
        </header>
        
        <!-- 页面内容 -->
        <main class="flex-1 p-6 overflow-auto">
            <!-- 页面内容 -->
        </main>
    </div>
</div>
```

## 8. 交互规范

### 8.1 动画效果
```css
/* 基础过渡 */
.transition-base {
    transition: all 0.2s ease-in-out;
}

/* 缓慢过渡 */
.transition-slow {
    transition: all 0.3s ease-in-out;
}

/* 快速过渡 */
.transition-fast {
    transition: all 0.1s ease-in-out;
}
```

### 8.2 悬停效果
```css
/* 按钮悬停 */
.btn:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

/* 卡片悬停 */
.card:hover {
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
    border-color: #d1d5db;
}
```

### 8.3 加载状态
```html
<!-- 按钮加载状态 -->
<button class="bg-green-600 text-white px-4 py-2 rounded-lg disabled:opacity-50" disabled>
    <i class="fas fa-spinner fa-spin mr-2"></i>
    加载中...
</button>

<!-- 页面加载状态 -->
<div class="flex items-center justify-center py-12">
    <div class="text-center">
        <i class="fas fa-spinner fa-spin text-green-600 text-2xl mb-4"></i>
        <p class="text-gray-600">正在加载...</p>
    </div>
</div>
```

## 9. 响应式设计

### 9.1 断点策略
```css
/* 移动端优先 */
.mobile-first {
    /* 默认移动端样式 */
}

@media (min-width: 640px) {
    .mobile-first {
        /* 大屏手机样式 */
    }
}

@media (min-width: 768px) {
    .mobile-first {
        /* 平板样式 */
    }
}

@media (min-width: 1024px) {
    .mobile-first {
        /* 桌面样式 */
    }
}
```

### 9.2 网格系统
```html
<!-- 响应式网格 -->
<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
    <!-- 网格项 -->
</div>

<!-- 响应式弹性布局 -->
<div class="flex flex-col md:flex-row gap-6">
    <!-- 弹性项 -->
</div>
```

## 10. 可访问性规范

### 10.1 语义化HTML
```html
<!-- 正确的语义化结构 -->
<article class="bg-white rounded-lg p-6">
    <header>
        <h2 class="text-xl font-semibold mb-2">文章标题</h2>
        <time datetime="2024-03-10" class="text-sm text-gray-500">2024年3月10日</time>
    </header>
    <main>
        <p class="text-gray-700">文章内容...</p>
    </main>
</article>
```

### 10.2 焦点管理
```css
/* 自定义焦点样式 */
.focus-ring:focus {
    outline: none;
    box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.5);
}

/* 跳过链接 */
.skip-link {
    position: absolute;
    top: -40px;
    left: 6px;
    background: #000;
    color: #fff;
    padding: 8px;
    text-decoration: none;
    z-index: 1000;
}

.skip-link:focus {
    top: 6px;
}
```

### 10.3 ARIA属性
```html
<!-- 模态框 -->
<div role="dialog" aria-labelledby="modal-title" aria-modal="true">
    <h2 id="modal-title">模态框标题</h2>
    <!-- 模态框内容 -->
</div>

<!-- 表单标签 -->
<label for="username">用户名 <span aria-label="必填">*</span></label>
<input id="username" type="text" aria-required="true" aria-describedby="username-help">
<div id="username-help" class="text-sm text-gray-500">请输入您的用户名</div>
```

## 11. 暗色主题支持

### 11.1 CSS变量定义
```css
:root {
    --bg-primary: #ffffff;
    --bg-secondary: #f9fafb;
    --text-primary: #111827;
    --text-secondary: #6b7280;
}

[data-theme="dark"] {
    --bg-primary: #1f2937;
    --bg-secondary: #111827;
    --text-primary: #f9fafb;
    --text-secondary: #d1d5db;
}
```

### 11.2 主题切换
```html
<!-- 主题切换按钮 -->
<button id="theme-toggle" class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700">
    <i class="fas fa-sun hidden dark:inline" title="切换到亮色主题"></i>
    <i class="fas fa-moon dark:hidden" title="切换到暗色主题"></i>
</button>
```

这个设计规范文档确保了整个项目的视觉一致性和用户体验质量。 