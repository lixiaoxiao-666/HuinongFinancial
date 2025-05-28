# OA管理系统前端构建说明

## 项目概述

本项目是数字惠农OA管理系统的前端部分，基于Vue 3 + TypeScript + Vite构建，使用Element Plus作为UI组件库。

## 技术栈

- **框架**: Vue 3.3+
- **构建工具**: Vite 4+
- **语言**: TypeScript 5+
- **UI库**: Element Plus 2.4+
- **路由**: Vue Router 4+
- **状态管理**: Pinia 2+
- **HTTP客户端**: Axios 1.6+
- **日期处理**: Day.js 1.11+
- **包管理器**: pnpm 8+

## 环境要求

- Node.js >= 18.0.0
- pnpm >= 8.0.0

## 安装和启动

### 1. 安装依赖

```bash
cd frontend/admin
pnpm install
```

### 2. 开发环境启动

```bash
pnpm dev
```

开发服务器将在 `http://localhost:5173` 启动

### 3. 生产构建

```bash
pnpm build
```

构建产物将输出到 `dist` 目录

### 4. 本地预览

```bash
pnpm preview
```

预览生产构建，服务器将在 `http://localhost:4173` 启动

## 项目结构

```
frontend/admin/
├── public/                 # 静态资源
├── src/
│   ├── api/               # API接口
│   │   ├── index.ts       # Axios配置
│   │   └── admin.ts       # 管理系统API
│   │
│   ├── assets/            # 资源文件
│   ├── components/        # 通用组件
│   ├── router/            # 路由配置
│   ├── stores/            # 状态管理
│   ├── types/             # TypeScript类型定义
│   ├── views/             # 页面组件
│   ├── App.vue            # 根组件
│   └── main.ts            # 入口文件
├── index.html             # HTML入口
├── package.json           # 依赖配置
├── tsconfig.json          # TypeScript配置
└── vite.config.ts         # Vite配置
```

## 页面组件说明

### 核心页面
- `LoginView.vue` - 登录页面
- `LayoutView.vue` - 主布局组件
- `DashboardView.vue` - 工作台
- `ApprovalView.vue` - 审批看板
- `ApprovalDetailView.vue` - 审批详情
- `UsersView.vue` - 用户管理
- `LogsView.vue` - 操作日志
- `SystemView.vue` - 系统设置
- `NotFoundView.vue` - 404页面

## API配置

### 后端接口地址
默认后端API地址为 `http://localhost:8080/api/v1`，可在 `src/api/index.ts` 中修改。

### 认证机制
使用JWT Token进行认证，Token存储在localStorage中，请求拦截器自动添加到请求头。

## 部署说明

### 1. 构建优化

生产构建已启用以下优化：
- 代码分割和懒加载
- 资源压缩和混淆
- Tree Shaking移除无用代码
- Gzip压缩支持

### 2. 环境变量

可通过环境变量配置：
- `VITE_API_BASE_URL` - API基础地址
- `VITE_APP_TITLE` - 应用标题

### 3. 部署步骤

1. 执行生产构建：
   ```bash
   pnpm build
   ```

2. 将 `dist` 目录内容部署到Web服务器

3. 配置服务器支持SPA路由（History模式）

### 4. Nginx配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /path/to/dist;
    index index.html;

    # 支持Vue Router History模式
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API代理（可选）
    location /api/ {
        proxy_pass http://localhost:8080/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 开发指南

### 1. 代码规范

- 使用TypeScript进行类型检查
- 遵循Vue 3 Composition API规范
- 组件使用PascalCase命名
- 文件名使用kebab-case

### 2. 状态管理

使用Pinia进行状态管理，主要store：
- `auth.ts` - 用户认证状态

### 3. 路由配置

- 支持路由守卫和权限控制
- 懒加载优化性能
- 面包屑导航自动生成

### 4. API调用

- 统一使用`src/api/admin.ts`中的方法
- 自动处理请求/响应拦截
- 统一错误处理和用户提示

## 测试账号

- **管理员**: admin / admin123
- **审批员**: reviewer / reviewer123

## 功能特性

- ✅ 用户认证和权限控制
- ✅ 工作台数据概览
- ✅ 贷款申请审批管理
- ✅ AI智能风险评估
- ✅ 用户管理
- ✅ 操作日志记录
- ✅ 系统配置管理
- ✅ 响应式设计
- ✅ 现代化UI界面

## 浏览器支持

- Chrome >= 87
- Firefox >= 78
- Safari >= 14
- Edge >= 88

## 故障排除

### 常见问题

1. **依赖安装失败**
   ```bash
   rm -rf node_modules pnpm-lock.yaml
   pnpm install
   ```

2. **开发服务器启动失败**
   - 检查端口5173是否被占用
   - 检查Node.js版本是否符合要求

3. **API请求失败**
   - 确认后端服务是否启动
   - 检查API地址配置是否正确
   - 查看浏览器Network面板的错误信息

4. **登录失败**
   - 确认使用正确的测试账号
   - 检查后端认证服务是否正常

## 联系方式

如有问题请联系开发团队或查看项目文档。 