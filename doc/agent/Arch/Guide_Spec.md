# 开发指南 (Guide_Spec.md)

## 1. 概述

本文档为"数字惠农APP及OA后台管理系统"的开发团队提供指导，旨在规范开发流程、统一编码风格、确保代码质量和项目顺利进行。本文档内容涵盖后端（Golang）和前端（HTML/CSS/JS）的开发规范及建议。

## 2. 通用开发规范

### 2.1 版本控制

*   **Git**: 所有代码均使用Git进行版本控制。
*   **分支策略**: 推荐使用 `Git Flow` 或类似的成熟分支模型。
    *   `main` (或 `master`): 稳定的生产环境代码，只接受来自 `release` 分支的合并。
    *   `develop`: 开发主分支，集成所有已完成的功能。
    *   `feature/<feature-name>`: 开发新功能的分支，从 `develop` 创建，完成后合并回 `develop`。
    *   `release/<version>`: 准备发布新版本的分支，从 `develop` 创建，用于Bug修复和版本元数据准备。
    *   `hotfix/<fix-name>`: 修复生产环境紧急Bug的分支，从 `main` 创建，完成后合并回 `main` 和 `develop`。
*   **Commit Message**: 遵循规范的Commit Message格式，例如 [Conventional Commits](https://www.conventionalcommits.org/)。
    *   示例: `feat: add user login API`
    *   示例: `fix: resolve issue with loan application submission`
*   **Code Review**: 所有代码变更（特别是合并到 `develop` 和 `main` 的）都必须经过至少一位其他团队成员的Code Review。

### 2.2 编码风格

*   **一致性**: 团队内部应遵循统一的编码风格。具体语言的风格指南见后续章节。
*   **可读性**: 代码应易于阅读和理解。使用有意义的变量名、函数名，并适当添加注释解释复杂逻辑。
*   **简洁性**: 避免不必要的复杂代码和冗余逻辑。
*   **DRY (Don't Repeat Yourself)**: 尽量复用代码，避免重复编写相同逻辑。

### 2.3 测试

*   **单元测试**: 针对核心模块和函数编写单元测试，确保其按预期工作。
*   **集成测试**: 测试多个组件协同工作的场景。
*   **端到端测试 (E2E)**: (如果条件允许) 测试完整的用户流程。
*   **测试覆盖率**: 鼓励提高测试覆盖率，特别是核心业务逻辑。

### 2.4 文档

*   **代码注释**: 对公共API、复杂算法、重要业务逻辑等进行清晰注释。
*   **API文档**: 后端API需有详细的文档说明 (已创建 `API_Spec.md`)。
*   **架构文档**: (已创建 `Arch_Spec.md`, `Database_Spec.md`)
*   **README**: 每个服务或主要模块应有README文件，说明其功能、如何构建、运行和测试。

### 2.5 环境管理

*   **开发环境 (Development)**: 本地开发环境。
*   **测试环境 (Testing)**: 用于功能测试和集成测试的共享环境。
*   **预发布环境 (Staging/UAT)**: 与生产环境配置一致，用于上线前最后验证。
*   **生产环境 (Production)**: 实际用户使用的环境。
*   使用配置文件或环境变量管理不同环境的配置参数，避免硬编码。

### 2.6 依赖管理

*   **后端 (Golang)**: 使用 Go Modules (`go.mod`, `go.sum`) 管理依赖。
*   **前端**: 由于是简单的HTML + Tailwind CSS + FontAwesome，主要通过CDN引入，版本需在HTML中明确指定。

## 3. 后端开发指南 (Golang)

### 3.1 项目结构 (建议)

推荐采用标准化的Go项目布局，例如 [Standard Go Project Layout](https://github.com/golang-standards/project-layout)。

对于单个微服务，可以简化为：

```
/<service-name>/
  ├── cmd/                # 应用主入口 (main.go)
  │   └── <service-name>/
  │       └── main.go
  ├── internal/           # 内部应用和库代码 (不对外暴露)
  │   ├── api/            # API处理程序 (HTTP handlers/gRPC services)
  │   ├── biz/            # 业务逻辑层 (business logic)
  │   ├── data/           # 数据访问层 (database, cache access)
  │   ├── conf/           # 配置文件定义和加载
  │   └── service/        # 服务层 (组织业务逻辑，供API层调用)
  ├── pkg/                # 对外暴露的库代码 (可被外部应用引用)
  ├── configs/            # 配置文件 (e.g., config.yaml)
  ├── api/                # API定义文件 (e.g., .proto, OpenAPI spec)
  ├── go.mod
  ├── go.sum
  └── README.md
```

### 3.2 编码风格与规范

*   遵循官方的 [Effective Go](https://go.dev/doc/effective_go) 和 [CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments)。
*   使用 `gofmt` 或 `goimports` 自动格式化代码。
*   **命名**: 使用驼峰式命名，首字母大写表示导出 (public)，小写表示私有 (internal)。
*   **错误处理**: 明确处理所有可能的错误。优先使用 `error` 类型返回值，而不是 `panic`。错误信息应清晰，便于排查。
    ```go
    result, err := someFunction()
    if err != nil {
        log.Printf("Error calling someFunction: %v", err)
        return err // 或者返回包装后的错误
    }
    ```
*   **并发**: 谨慎使用goroutine和channel，避免数据竞争和死锁。使用 `sync` 包提供的同步原语。
*   **日志**: 使用结构化日志 (structured logging)，例如 `log/slog` (Go 1.21+) 或第三方库 (如 `logrus`, `zap`)。
    *   日志级别: DEBUG, INFO, WARN, ERROR, FATAL.
    *   日志内容: 时间戳、级别、调用者信息、消息、以及相关的上下文数据 (如 `request_id`, `user_id`)。
*   **配置管理**: 使用 `viper` 或类似库加载和管理配置文件，支持多种格式 (YAML, JSON, TOML) 和环境变量覆盖。
*   **数据库操作**:
    *   使用 `database/sql` 包进行数据库交互，或者选择合适的ORM如 `GORM`。
    *   注意SQL注入风险，使用参数化查询或预编译语句。
    *   管理好数据库连接池。
*   **API设计**:
    *   遵循RESTful原则设计HTTP API。
    *   请求体验证：对API的输入参数进行严格校验。
    *   响应格式：统一API响应结构 (如 `API_Spec.md` 中定义的)。

### 3.3 微服务开发要点

*   **服务注册与发现**: 与Higress AI网关集成，实现服务的自动注册和发现。
*   **配置中心**: (可选) 如果配置复杂或需要动态更新，可以考虑引入配置中心 (如Nacos, Consul, etcd)。
*   **链路追踪**: (可选) 引入OpenTelemetry或类似方案，实现分布式链路追踪，便于排查微服务间的调用问题。
*   **容错处理**: 实现合理的超时、重试、熔断机制 (Higress网关已提供部分能力)。
*   **无状态服务**: 尽量设计无状态服务，便于水平扩展和故障恢复。会话信息可存储在Redis中。

### 3.4 推荐库与工具

*   **Web框架**: `Gin`, `Echo`
*   **ORM**: `GORM`
*   **配置**: `Viper`
*   **日志**: `log/slog` (Go 1.21+), `logrus`, `zap`
*   **校验**: `go-playground/validator`
*   **测试**: `testing` (标准库), `testify/assert`, `testify/mock`
*   **UUID**: `google/uuid`

## 4. 前端开发指南 (HTML/CSS/JS)

由于本项目前端技术栈相对简单 (HTML + Tailwind CSS + FontAwesome，可能辅以少量原生JS)，指南侧重于规范和最佳实践。

### 4.1 项目结构 (建议)

```
/frontend/  (或直接在各服务/模块的 static 目录下)
  ├── assets/
  │   ├── css/              # 自定义CSS (如果Tailwind不完全满足)
  │   │   └── style.css
  │   ├── js/               # 自定义JavaScript文件
  │   │   └── main.js
  │   └── images/           # 项目图片资源 (如有本地存储需求)
  ├── pages/                # HTML页面文件 (或按模块组织)
  │   ├── user/
  │   │   ├── login.html
  │   │   └── profile.html
  │   ├── loan/
  │   │   └── apply.html
  │   └── index.html
  └── README.md
```

### 4.2 HTML规范

*   **语义化**: 使用HTML5语义化标签 (e.g., `<header>`, `<footer>`, `<nav>`, `<article>`, `<section>`, `<aside>`)。
*   **可访问性 (A11y)**: 
    *   为图片提供有意义的 `alt` 属性。
    *   表单元素使用 `<label>` 关联。
    *   确保键盘可导航。
    *   必要时使用ARIA属性增强可访问性。
*   **结构清晰**: 保持HTML结构清晰、简洁，正确嵌套标签。
*   **字符编码**: 使用 `UTF-8` (`<meta charset="UTF-8">`)。
*   **Viewport**: 设置正确的视口 (`<meta name="viewport" content="width=device-width, initial-scale=1.0">`)。
*   **CDN资源**: 使用CDN加载Tailwind CSS和FontAwesome时，确保版本号固定，并考虑使用 `integrity` 和 `crossorigin` 属性增强安全性。
    ```html
    <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css" integrity="sha512-iecdLmaskl7CVkqkXNQ/ZH/XLlvWZOJyj7Yy7tcenmpD1ypASozpmT/E0iPtmFIB46ZmdtAc9eNBvH0H/ZpiBw==" crossorigin="anonymous" referrerpolicy="no-referrer" />
    ```

### 4.3 CSS规范 (Tailwind CSS为主)

*   **优先使用Tailwind类**: 尽可能利用Tailwind CSS的实用类 (utility classes) 来构建样式，减少自定义CSS。
*   **配置文件**: 如果需要自定义Tailwind (如颜色、字体、断点)，修改 `tailwind.config.js` (如果项目结构支持JS构建流程)。对于纯HTML项目，可能需要依赖Tailwind Play CDN的配置或默认配置。
*   **自定义CSS**: 如果必须编写自定义CSS，应保持其模块化和简洁性，避免与Tailwind类冲突。
*   **命名**: 自定义CSS类名使用小写短横线连接 (kebab-case)，如 `custom-card-header`。
*   **响应式设计**: 利用Tailwind的响应式前缀 (e.g., `sm:`, `md:`, `lg:`) 实现响应式布局。

### 4.4 JavaScript规范 (少量使用)

*   **原生优先**: 对于简单交互，优先考虑原生JavaScript，避免引入不必要的库。
*   **代码位置**: JS代码可以放在 `<script>` 标签中，或外部 `.js` 文件中。外部文件应放在 `</body>` 前引入，以避免阻塞页面渲染，除非有特定需求 (如使用了 `defer` 或 `async`)。
*   **ES6+语法**: 使用现代JavaScript语法 (ES6+)，但需注意浏览器兼容性 (本项目PRD中未明确要求兼容旧版浏览器，可适当放宽)。
*   **DOM操作**: 高效操作DOM，避免频繁重绘和回流。
*   **事件处理**: 合理使用事件委托，优化性能。
*   **错误处理**: 对可能出错的操作 (如API请求) 进行错误处理。
*   **API调用**: 使用 `fetch` API 或 `XMLHttpRequest` (较旧) 进行后端API调用。处理异步操作 (Promises, async/await)。
    ```javascript
    async function fetchData(url) {
      try {
        const response = await fetch(url, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('authToken')}` // 示例Token
          }
        });
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        console.log(data);
        // 更新UI
      } catch (error) {
        console.error('Error fetching data:', error);
        // 显示错误信息给用户
      }
    }
    ```
*   **安全性**: 
    *   防范XSS攻击：对用户输入和动态插入DOM的内容进行适当转义。
    *   不将敏感信息硬编码在前端代码中。

### 4.5 性能优化 (前端)

*   **图片优化**: 压缩图片，使用合适的图片格式 (JPEG, PNG, WebP)。
*   **减少HTTP请求**: 合并CSS/JS文件 (如果未使用CDN或有构建流程)。
*   **利用浏览器缓存**: 合理设置HTTP缓存头 (主要由后端API和CDN配置)。
*   **代码压缩**: HTML, CSS, JS代码在部署前进行压缩 (如果未使用已压缩的CDN资源)。
*   **懒加载**: 对于非首屏图片或内容，可以考虑使用懒加载技术。

## 5. 与AI组件的集成

*   **Dify平台**: 后端服务 (主要是惠农贷服务和OA后台服务) 将通过API与Dify平台部署的AI智能体进行交互。API调用细节需参照Dify智能体的具体定义。
*   **本地大模型**: 调用本地大模型通常也会通过Higress AI网关进行代理和管理。后端服务向网关发起请求，网关根据配置路由到相应的本地模型服务。
*   **MCP工具集**: AI智能体 (Dify) 或后端服务在需要时，可能调用MCP工具集提供的接口来辅助决策 (如查询数据库获取特定数据)。
*   **数据流转**: 用户在前端提交数据 -> 后端服务接收并进行初步处理 -> 调用Dify/本地模型进行AI分析 -> AI结果返回给后端服务 -> 后端服务根据AI结果进行业务决策 -> 结果返回前端展示。

## 6. 安全注意事项

*   **输入验证**: 前后端都必须对用户输入进行严格验证。
*   **SQL注入防护**: 后端使用参数化查询或ORM来防止SQL注入。
*   **XSS防护**: 前端对动态内容进行转义，后端对输出到HTML的内容也应注意。
*   **CSRF防护**: (如果前端有复杂表单提交且未使用JS框架的内置防护) 后端应实现CSRF Token机制。
*   **认证与授权**: API网关 (Higress) 和各服务层面都需要严格的认证和授权机制。
*   **敏感数据处理**: 
    *   密码哈希存储 (如bcrypt)。
    *   敏感信息传输加密 (HTTPS)。
    *   敏感信息在日志中脱敏。
*   **依赖安全**: 定期扫描项目依赖，及时更新有漏洞的库。
*   **AI安全**: 参考Higress AI网关提供的能力，如内容安全过滤、Prompt注入防范等。

## 7. CI/CD (持续集成/持续部署)

*   **自动化流程**: 建立自动化的构建、测试、部署流水线。
*   **工具**: Jenkins, GitLab CI, ArgoCD (用于GitOps) 等。
*   **K8S部署**: 通过CI/CD流水线将Docker镜像部署到Kubernetes集群。

本文档为初始版本，随着项目的进展会不断更新和完善。团队成员应定期回顾并遵循此指南。 