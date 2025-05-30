# 数字惠农系统 - 前端工程化文档

## 1. 文档概述

### 1.1 文档目的
本文档旨在定义"数字惠农APP"及"惠农OA后台管理系统"前端项目的工程化规范，包括技术选型、目录结构、编码规范、代码管理、构建与部署等，以提高开发效率、保证代码质量、方便团队协作和项目维护。

### 1.2 参考文档
- 《需求说明文档 (PRD) - 数字惠农APP及OA后台管理系统》 (`agent/demand/PRD-2025-5-29.md`)
- 《数字惠农系统 - 前端 UI/UX 设计文档》 (`agent/frontend/UI-UX.md`)
- 《数字惠农系统 - 后端 API 接口文档》 (`agent/backend/API.md`)

## 2. 技术选型

根据项目需求和特点，拟采用以下技术栈：

### 2.1 数字惠农APP
-   **框架**: Vue 3 (使用 TypeScript) + Uni-app。
    -   *理由*: Vue 3 提供了更优的性能和开发体验。Uni-app 支持一套代码编译到多个移动端平台（iOS, Android）以及小程序等，有成熟的生态和UI库，符合国内技术趋势。TypeScript提供类型安全。
-   **状态管理**: Pinia。
    -   *理由*: Pinia 是 Vue 3 官方推荐的状态管理库，轻量、简单易用，且对 TypeScript 支持良好。
-   **导航**: Uni-app 内置路由 (基于 `pages.json` 配置)。
    -   *理由*: Uni-app 框架自带，配置简单，满足大部分APP导航需求。
-   **UI组件库**: uView UI / Thor UI (或其他 Uni-app 生态中成熟的组件库)。
    -   *理由*: 专为 Uni-app 设计，提供丰富的跨端组件，加速开发。结合《UI/UX设计文档》进行定制。
-   **数据请求**: Axios。
    -   *理由*: 功能丰富，支持拦截器、取消请求等，社区成熟，与Vue配合良好。
-   **图表库 (如需)**: uCharts / ECharts (可通过ucharts封装或特定uni-app插件使用)。
-   **地图库 (如需)**: Uni-app 内置地图组件 / 腾讯地图、高德地图等 SDK (按需引入)。

### 2.2 惠农OA后台管理系统 (Web)
-   **框架**: Vue 3 (使用 TypeScript) + Vite。
    -   *理由*: Vue 3 带来了显著的性能提升和更好的开发体验（如Composition API）。TypeScript提供类型安全。Vite提供极速的冷启动、热更新和构建体验。
-   **状态管理**: Pinia。
    -   *理由*: Vue 3 官方推荐，轻量、模块化、对TypeScript友好。
-   **路由**: Vue Router 4。
    -   *理由*: Vue 官方路由管理器，与 Vue 3 深度集成。
-   **UI组件库**: Element Plus.
    -   *理由*: Element Plus 作为成熟的企业级Vue 3 UI组件库，提供丰富、设计精美的组件。它将作为惠农OA后台管理系统主要的UI构建和界面美化方案，以确保统一、专业的视觉风格和用户体验。其强大的主题定制能力将用于实现《UI/UX设计文档》中定义的视觉要求。
-   **数据请求**: Axios。
-   **图表库**: ECharts for Vue / AntV G2Plot (Vue版)。

### 2.3 通用工具
-   **包管理器**: pnpm / yarn / npm (推荐 pnpm，速度快且节省磁盘空间)。
-   **代码规范**: ESLint, Prettier, Stylelint (针对CSS/SCSS)。
-   **提交规范**: Commitlint (配合 Husky)。
-   **版本控制**: Git。

## 3. 项目目录结构

为保证项目结构清晰、模块化，建议采用以下目录结构。这些结构是推荐的最佳实践，现有项目可逐步向此看齐或在新模块中遵循此结构。实际项目中可根据具体情况进行调整。

### 3.1 数字惠农APP (Vue 3 + Uni-app)

```
/app (uni-app 项目根目录)
├── common/                 # 公共模块、工具类
│   ├── api/                # API 请求模块 (按业务模块划分)
│   │   ├── auth.ts
│   │   └── loan.ts
│   ├── config/             # 应用配置 (环境变量, API地址等)
│   │   └── index.ts
│   ├── constants/          # 常量
│   ├── store/              # 状态管理 (Pinia)
│   │   ├── index.ts
│   │   └── modules/
│   ├── styles/             # 全局样式, 主题配置 (如 uni.scss)
│   ├── utils/              # 工具函数
│   └── types/              # TypeScript 类型定义
│       ├── api.ts
│       └── index.ts
├── components/             # 可复用的UI组件 (符合uni-app规范)
│   ├── Common/             # 通用基础组件 (如HnButton, HnInput, HnCard)
│   └── Business/           # 业务相关组件 (如LoanItem, MachineCard)
├── pages/                  # 业务页面 (符合pages.json规范)
│   ├── Auth/
│   │   ├── Login.vue
│   │   └── Register.vue
│   ├── Home/
│   │   └── Index.vue
│   └── Loan/
│       ├── LoanList.vue
│       └── LoanApply.vue
├── static/                 # 静态资源 (图片, 字体, 图标等)
│   ├── fonts/
│   ├── images/
│   └── icons/
├── hybrid/                 # App端存放本地html、js等静态资源 (如需)
├── platforms/              # 特定平台代码 (如App端原生插件)
├── services/               # 业务服务逻辑 (复杂的业务计算，数据处理)
├── App.vue                 # 应用配置，用来配置App全局样式以及监听 应用生命周期
├── main.ts                 # Vue初始化入口文件
├── manifest.json           # 应用配置清单，打包发布所需
├── pages.json              # 页面路由及窗口表现配置
├── uni.scss                # Uni-app内置的常用样式变量
├── .env.development
├── .env.production
├── .env.example
├── .eslintrc.js
├── .prettierrc.js
├── tsconfig.json
└── package.json
```

### 3.2 惠农OA后台管理系统 (Vue 3 + Vite)

```
/oa-admin
├── public/
│   └── index.html
├── src/
│   ├── api/                # API 请求模块 (按业务模块划分)
│   │   ├── auth.ts
│   │   └── loanApproval.ts
│   ├── assets/             # 静态资源
│   │   ├── images/
│   │   ├── icons/
│   │   └── styles/           # 全局静态样式文件 (如 reset.css)
│   ├── components/         # 可复用的UI组件
│   │   ├── Common/           # 通用基础组件 (如 SvgIcon, HnTable)
│   │   └── Business/         # 业务相关组件
│   ├── config/             # 应用配置
│   │   └── index.ts
│   ├── constants/          # 常量
│   ├── directives/         # 自定义指令
│   ├── hooks/              # 自定义Vue Hooks (Composition API)
│   ├── layouts/            # 布局组件 (如后台主布局: DefaultLayout.vue)
│   │   └── components/       # 布局内部组件 (如 Sidebar, Navbar, AppMain)
│   │   └── DefaultLayout.vue
│   ├── locales/            # 国际化语言包 (如需)
│   ├── router/             # 路由配置
│   │   └── index.ts
│   ├── services/           # 业务服务逻辑
│   ├── store/              # 状态管理 (Pinia)
│   │   ├── index.ts
│   │   └── modules/
│   ├── styles/             # 全局样式, 主题配置 (如 global.scss, variables.scss)
│   │   ├── global.scss
│   │   └── variables.scss
│   ├── types/              # TypeScript 类型定义
│   │   ├── api.ts
│   │   └── index.ts
│   ├── utils/              # 工具函数
│   ├── views/              # 页面级组件 (按业务模块划分)
│   │   ├── Auth/
│   │   │   └── Login.vue
│   │   ├── Dashboard/
│   │   │   └── Index.vue
│   │   └── LoanApproval/
│   │       ├── List.vue
│   │       └── Detail.vue
│   ├── App.vue               # 应用主组件
│   └── main.ts               # 应用主入口
├── .env.development
├── .env.production
├── .env.example
├── .eslintrc.cjs
├── .prettierrc.json
├── tsconfig.json
├── tsconfig.node.json
├── vite.config.ts
└── package.json
```

## 4. 编码规范

### 4.1 JavaScript / TypeScript
-   遵循 ESLint 配置规则 (推荐使用 `@vue/eslint-config-typescript` 结合 Prettier)。
-   使用 Prettier 进行代码格式化，确保风格统一。
-   **命名规范**:
    -   变量和函数名：小驼峰命名 (camelCase)。
    -   接口名和类型别名：大驼峰命名 (PascalCase)，可考虑以 `I` 或 `T` 开头 (如 `IUser`, `TConfig`)。
    -   类和组件名 (Vue 组件)：大驼峰命名 (PascalCase)。
    -   常量名：全大写，下划线分隔 (UPPER_CASE_SNAKE_CASE)。
    -   文件名：
        -   Vue 组件文件 (`.vue`): 大驼峰命名 (e.g., `MyComponent.vue`)。
        -   TypeScript 文件 (`.ts`): 业务逻辑、工具函数、API模块、Store模块等使用小驼峰命名 (e.g., `useCounter.ts`, `authApi.ts`, `userStore.ts`) 或大驼峰命名 (如果导出的是类或构造函数)。
-   **注释**: 对复杂逻辑、重要函数、公共组件等添加必要的注释。推荐使用 JSDoc 风格。
-   **模块化**: 使用 ES6 模块导入导出。
-   **类型**: 充分利用 TypeScript 的类型系统，为变量、函数参数、返回值等添加明确的类型定义。避免使用 `any` 类型，除非确实必要且有注释说明。
-   **Vue 3 组件编写**:
    -   强烈推荐使用 `<script setup>` 语法糖结合 Composition API 进行组件开发，以获得更好的类型推断和代码组织。
    -   组件拆分应遵循单一职责原则，保持组件精炼、易于维护。
    -   Props 定义使用 `defineProps<T>()`，并提供明确的 TypeScript 类型或接口。
    -   Emits 定义使用 `defineEmits<T>()`，并提供明确的事件签名。
    -   组件的 `name` 属性应在 `<script setup>` 外的 `<script>` 块中明确指定 (利于调试和组件递归)，并与文件名保持一致 (PascalCase)。
        ```vue
        <script lang="ts">
        export default {
          name: 'MyComponent'
        }
        </script>
        <script setup lang="ts">
        // ... setup logic ...
        </script>
        ```
    -   在模板中，指令 (`v-model`, `v-for`) 和特性名 (props) 推荐使用 kebab-case (如 `v-model`, `my-prop`)。 组件标签使用 PascalCase (如 `<MyComponent />`)。
    -   避免在模板中编写复杂的 JavaScript 表达式，应将其封装在 `computed` 属性或组件的 `methods` (在 Composition API 中通常是 `<script setup>` 内的函数) 中。
    -   合理使用 `v-if`, `v-else-if`, `v-else` 和 `v-show`，理解它们的渲染机制和适用场景。
    -   对于列表渲染 (`v-for`)，务必为每个元素提供唯一的 `key`，避免使用索引作为 `key`，除非列表是静态且不会改变顺序的。
    -   优先使用 `ref` 和 `reactive` 创建响应式数据。理解它们的区别。
    -   副作用操作 (如API请求, DOM操作) 应放在 `onMounted`, `onUpdated` 等生命周期钩子或 `watch`, `watchEffect` 中。

### 4.2 CSS / SCSS / Less
-   遵循 Stylelint 配置规则 (可集成 `stylelint-config-standard-vue` 或类似配置)。
-   使用 Prettier 进行代码格式化。
-   **强烈推荐**在 Vue 组件中使用 `<style scoped>` 来实现组件样式的局部作用域，以避免全局样式污染。
-   对于全局样式、主题变量或可复用的原子化 CSS 类，可以定义在单独的样式文件中 (如 `src/styles/global.scss`, `src/styles/variables.scss`) 并按需引入到主入口文件 (`main.ts`) 或根组件。
-   命名规范 (如 BEM, SMACSS) 可以辅助组织非 scoped 样式，但在 scoped 环境下，可以直接使用简洁的类名。
-   避免使用过深的 CSS 选择器嵌套，特别是在预处理器 (SCSS, Less) 中。
-   合理使用 CSS 预处理器 (如 SCSS, Less) 的变量、mixins、functions 等特性来提高样式的可维护性和复用性。
-   对于 Uni-app 项目，遵循其特定的 rpx 单位和样式规范。

### 4.3 文件与目录文档
-   对于每个主要的业务模块目录 (如 `src/api`, `src/screens/Loan`, `src/pages/LoanApproval`)，建议在该目录下创建一个 `README.md` 文件，简要说明该模块的功能、包含的主要组件/页面以及与其他模块的依赖关系。
-   对于复杂或核心的组件，可以在组件文件顶部添加注释说明其用途、Props 和使用示例。
-   代码实现应与 `agent/frontend/code/` 目录下的对应代码文件的 MD 文档同步。当代码发生更改时，也需要同步更改对应代码文件的 MD 文档，保证一代码一文档。

## 5. 代码管理 (Git)

-   **分支策略**: 推荐使用 Git Flow 或类似的简化版分支模型。
    -   `main` (或 `master`): 生产环境分支，稳定版本。
    -   `develop`: 开发主分支，集成新功能。
    -   `feature/xxx`: 功能开发分支，从 `develop` 切出，完成后合并回 `develop`。
    -   `release/vx.x.x`: 发布分支，用于准备发布版本，从 `develop` 切出，完成后合并到 `main` 和 `develop`。
    -   `hotfix/xxx`: 紧急修复分支，从 `main` 切出，修复后合并到 `main` 和 `develop`。
-   **提交信息规范**: 遵循 Conventional Commits 规范。
    -   格式: `<type>(<scope>): <subject>`
    -   示例: `feat(loan): add loan application form`
    -   `type` 可选值: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore` 等。
    -   使用 Commitlint 工具强制校验提交信息。
-   **Code Review**: 所有合并到 `develop` 和 `main` 分支的代码必须经过至少一位其他团队成员的 Code Review。
-   **合并请求 (Pull Request / Merge Request)**: 使用 PR/MR 进行代码合并，PR/MR 描述清晰，关联相关任务/Issue。

## 6. 构建与部署

### 6.1 数字惠农APP
-   **构建**: 使用 HBuilderX (Uni-app官方IDE) 或 Uni-app CLI (`@dcloudio/uni-cli-shared`) 进行构建。
    -   多端输出: 可编译为 Android (APK/AAB), iOS (需Xcode配合), 以及各类小程序包。
-   **环境配置**: 通过 `.env.[mode]` 文件 (如 `.env.development`, `.env.production`) 结合 `process.env` 或 `import.meta.env` (Vite构建时) 来区分不同环境的配置。
-   **持续集成/持续部署 (CI/CD)**: (可选，推荐) 使用 GitLab CI/CD, GitHub Actions, Jenkins 等，结合 Uni-app CLI 进行自动化构建和发布。
-   **发布**: 
    -   App: 发布到各大应用商店 (如 Apple App Store, Google Play Store) 和内部测试平台 (如蒲公英、Fir.im)。
    -   小程序: 发布到对应的小程序平台。

### 6.2 惠农OA后台管理系统 (Web)
-   **构建**: 使用 Vite (`vite build`) 进行打包，生成静态资源文件。
-   **环境配置**: 通过 Vite 的环境变量机制 (`.env.[mode]`) 配置不同环境。
-   **CI/CD**: (可选，推荐) 使用 GitLab CI/CD, GitHub Actions, Jenkins 等工具自动化构建和部署。
-   **部署**: 将构建后的静态文件部署到Web服务器 (如 Nginx, Apache) 或静态网站托管服务 (如 Vercel, Netlify, GitHub Pages, Firebase Hosting)。

## 7. 测试

-   **单元测试**: 
    -   Vue 组件: 使用 `@vue/test-utils` 结合 Jest 或 Vitest。
    -   TypeScript/JavaScript 逻辑: 使用 Jest 或 Vitest。
    -   Uni-app: 官方测试方案尚不完善，可针对核心逻辑进行单元测试，UI层面依赖E2E或人工测试较多。
-   **集成测试**: 测试模块间的交互，如 Pinia stores 和 API 服务的集成。
-   **端到端测试 (E2E)**: 
    -   Web (Vue 3): 使用 Cypress / Playwright。
    -   Uni-app: 可尝试 Appium, 或特定平台的自动化测试工具，或依赖人工测试。
-   **测试覆盖率**: 设定目标测试覆盖率 (如使用 `istanbul` 或 Vitest 内置覆盖率工具)，并持续监控。

## 8. 性能优化

-   **代码层面 (Vue 3 & Uni-app)**:
    -   避免不必要的组件重渲染 (合理使用 `computed`, `watch`, `shouldComponentUpdate` (选项式API中), 或 Composition API 的精细化响应控制)。
    -   优化列表性能 (Uni-app 使用 `<scroll-view>` 或虚拟列表组件，Web端使用虚拟滚动库)。
    -   合理使用异步组件 (`defineAsyncComponent`) 和路由懒加载。
    -   `v-show` vs `v-if`: 根据条件切换频率和初始渲染开销选择。
    -   `keep-alive`: 缓存失活的动态组件。
    -   `onActivated` / `onDeactivated`: 配合 `keep-alive` 管理组件状态。
    -   Uni-app: 注意分包加载优化小程序/App启动速度和包体积。
-   **资源优化 (Web & Uni-app)**: 
    -   代码分割 (Vite/Webpack自动处理，Uni-app有分包机制)。
    -   图片优化: 压缩图片 (如使用 TinyPNG, ImageOptim)，使用 WebP 等现代格式 (需考虑兼容性)，Uni-app 中合理使用 `image` 组件的 `mode` 属性。
    -   Gzip/Brotli 压缩 (Web服务器配置)。
    -   使用 CDN 加速静态资源访问。
-   **性能监控**: 
    -   浏览器开发者工具 (Performance, Network tabs)。
    -   Vue Devtools。
    -   Uni-app: HBuilderX 内置性能分析工具，或原生平台提供的分析工具。
    -   第三方服务 (如 Sentry, Fundebug) 用于错误监控和性能追踪。

## 9. 文档同步

-   项目开发过程中，应确保代码实现与 `agent/frontend/code/` 目录下各模块对应的 `.md` 文档保持同步更新。MD文档应包含模块功能说明、主要组件介绍、关键逻辑等。
-   API接口的使用应参照 `agent/frontend/API.md`，如有变更，需及时更新。

## 10. 团队协作

-   定期进行代码同步 (git pull/push)。
-   通过项目管理工具 (如 Jira, Trello, Asana) 跟踪任务进度。
-   定期举行简短的站会，同步进展和遇到的问题。
-   保持良好沟通，及时解决开发中的疑问和冲突。

本文档将作为前端团队开发工作的基本遵循，后续可根据项目实际进展和团队反馈进行修订和完善。 