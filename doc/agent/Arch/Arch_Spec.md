# 架构设计文档 (Arch_Spec.md)

## 1. 引言

本文档旨在详细描述"数字惠农APP及OA后台管理系统"的系统架构设计，包括系统组件、部署拓扑以及技术选型。本文档是根据项目需求文档 (`PRD.md`) 和用户故事地图 (`User_Story_Map.md`) 制定的，旨在为开发团队提供清晰的架构指导。

## 2. 系统架构图

### 2.1 组件视图 (Component View)

系统采用微服务架构思想，将不同的业务能力划分为独立的组件，通过明确定义的API进行交互。

```mermaid
graph TD
    subgraph 用户端 (User Facing Layer)
        direction LR
        MobileApp[移动App (HTML/CSS/JS)]
        WebApp[Web网页版 (HTML/CSS/JS)]
    end

    subgraph 应用层 (Application Layer - Golang Microservices)
        direction LR
        APIGateway[API网关 (Higress AI 网关)]
        UserService[用户服务]
        LoanService[惠农贷服务]
        MachineryService[农机租赁服务]
        NotificationService[通知服务 (可选)]
        AdminService[OA后台服务]
    end

    subgraph 智能处理层 (Intelligent Processing Layer)
        direction LR
        DifyPlatform[Dify AI平台 (智能体、模型编排)]
        LocalModels[本地AI大模型集群]
        MCPTools[MCP工具集]
    end

    subgraph 数据存储与支撑层 (Data & Support Layer)
        direction LR
        TiDB[TiDB数据库 (分布式)]
        Redis[Redis缓存 (会话、热点数据)]
        FileStorage[文件存储服务 (如MinIO)]
        LoggingService[日志服务 (ELK/EFK)]
        MonitoringService[监控服务 (Prometheus/Grafana)]
    end

    subgraph 基础设施层 (Infrastructure Layer)
        direction LR
        K8S[Kubernetes集群 (RKE2)]
        Karmada[Karmada (多集群管理)]
        Docker[Docker容器]
        AnolisOS[龙蜥Anolis OS]
    end

    MobileApp --> APIGateway
    WebApp --> APIGateway

    APIGateway --> UserService
    APIGateway --> LoanService
    APIGateway --> MachineryService
    APIGateway --> NotificationService
    APIGateway --> AdminService

    LoanService --> DifyPlatform
    AdminService --> DifyPlatform
    AdminService --> MCPTools
    DifyPlatform --> LocalModels
    DifyPlatform -- 调用 --> MCPTools

    UserService --> TiDB
    UserService --> Redis
    LoanService --> TiDB
    MachineryService --> TiDB
    AdminService --> TiDB
    NotificationService --> TiDB
    NotificationService --> Redis

    LoanService --> FileStorage
    MachineryService --> FileStorage

    UserService -- 写日志 --> LoggingService
    LoanService -- 写日志 --> LoggingService
    MachineryService -- 写日志 --> LoggingService
    AdminService -- 写日志 --> LoggingService
    APIGateway -- 写日志 --> LoggingService

    K8S -- 运行 --> APIGateway
    K8S -- 运行 --> UserService
    K8S -- 运行 --> LoanService
    K8S -- 运行 --> MachineryService
    K8S -- 运行 --> NotificationService
    K8S -- 运行 --> AdminService
    K8S -- 运行 --> DifyPlatform
    K8S -- 运行 --> LocalModels
    K8S -- 运行 --> TiDB
    K8S -- 运行 --> Redis
    K8S -- 运行 --> FileStorage
    K8S -- 运行 --> LoggingService
    K8S -- 运行 --> MonitoringService
    K8S -- 管理 --> Docker
    Karmada -- 管理 --> K8S
    AnolisOS -- 支撑 --> K8S
```

**组件说明:**

*   **用户端 (MobileApp, WebApp)**: 用户直接交互的界面，使用HTML、Tailwind CSS和FontAwesome构建。
*   **API网关 (Higress AI 网关)**: 所有客户端请求的统一入口，负责请求路由、认证鉴权、限流熔断、日志记录、多模型管理和AI安全防护。
*   **用户服务 (UserService)**: 处理用户注册、登录、个人信息管理等。
*   **惠农贷服务 (LoanService)**: 处理惠农贷产品展示、申请提交、进度查询、与AI审批流程对接等。
*   **农机租赁服务 (MachineryService)**: 处理农机信息发布、浏览、搜索、租赁订单管理等。
*   **通知服务 (NotificationService)**: (可选) 负责发送系统通知、短信、邮件等。
*   **OA后台服务 (AdminService)**: 为OA后台提供数据接口，支持审批看板、智能审批详情、人工复核、系统管理等功能。
*   **Dify AI平台**: 承载AI智能体，负责表单数据分析、风险评估模型的编排与调用。
*   **本地AI大模型集群**: 部署国产AI大模型，提供核心AI能力。
*   **MCP工具集**: 辅助AI进行数据库查询和特定决策逻辑。
*   **TiDB数据库**: 项目的主数据库，存储所有业务持久化数据。
*   **Redis缓存**: 存储用户会话信息、热点数据、分布式锁等，提升系统性能。
*   **文件存储服务**: 存储用户上传的图片、文档等非结构化数据。
*   **日志服务/监控服务**: 负责系统日志收集分析和性能监控告警。
*   **基础设施层**: 提供系统运行所需的容器化、编排、操作系统等基础环境。

### 2.2 部署视图 (Deployment View)

系统将部署在基于龙蜥Anolis OS的Kubernetes (RKE2) 集群中，并由Karmada进行多集群管理（如果需要）。

```mermaid
graph TD
    subgraph 用户访问 (User Access)
        UserDevice[用户设备 (手机/PC)]
    end

    subgraph 网络层 (Network Layer)
        Internet[互联网]
        LB[负载均衡器 (如Nginx/HAProxy或云LB)]
    end

    subgraph Kubernetes集群 (K8S Cluster - RKE2 on Anolis OS)
        KarmadaControlPlane[Karmada 控制平面 (可选, 用于多集群)]

        subgraph K8S集群1 (华北)
            MasterNodes1[Master节点 (Anolis OS)]
            WorkerNodes1[Worker节点 (Anolis OS)]

            subgraph WorkerNodes1
                direction LR
                Pod_Higress1[Pod: Higress AI 网关]
                Pod_UserService1[Pod: 用户服务]
                Pod_LoanService1[Pod: 惠农贷服务]
                Pod_AdminService1[Pod: OA后台服务]
                Pod_Redis1[Pod: Redis]
            end
        end

        subgraph K8S集群2 (华东 - 可选, 用于灾备或扩展)
            MasterNodes2[Master节点 (Anolis OS)]
            WorkerNodes2[Worker节点 (Anolis OS)]

            subgraph WorkerNodes2
                direction LR
                Pod_Higress2[Pod: Higress AI 网关]
                Pod_MachineryService2[Pod: 农机租赁服务]
                Pod_Dify2[Pod: Dify AI平台]
                Pod_LocalModels2[Pod: 本地AI大模型]
            end
        end

        subgraph 共享存储与数据库集群 (Shared Storage & DB)
            TiDBCluster[TiDB 分布式数据库集群 (跨可用区)]
            FileStorageCluster[分布式文件存储 (如MinIO集群)]
        end

        LoggingMonitoring[日志/监控集群 (ELK/Prometheus)]
    end

    UserDevice --> Internet
    Internet --> LB
    LB --> Pod_Higress1
    LB --> Pod_Higress2


    KarmadaControlPlane -.-> K8S集群1
    KarmadaControlPlane -.-> K8S集群2

    MasterNodes1 -- 管理 --> WorkerNodes1
    MasterNodes2 -- 管理 --> WorkerNodes2

    Pod_Higress1 --> Pod_UserService1
    Pod_Higress1 --> Pod_LoanService1
    Pod_Higress1 --> Pod_AdminService1

    Pod_Higress2 --> Pod_MachineryService2
    Pod_Higress2 --> Pod_Dify2

    Pod_LoanService1 --> Pod_Dify2
    Pod_AdminService1 --> Pod_Dify2
    Pod_Dify2 --> Pod_LocalModels2

    Pod_UserService1 --> TiDBCluster
    Pod_LoanService1 --> TiDBCluster
    Pod_AdminService1 --> TiDBCluster
    Pod_MachineryService2 --> TiDBCluster
    Pod_Dify2 --> TiDBCluster

    Pod_UserService1 --> Pod_Redis1
    Pod_LoanService1 --> Pod_Redis1
    Pod_AdminService1 --> Pod_Redis1

    Pod_LoanService1 --> FileStorageCluster
    Pod_MachineryService2 --> FileStorageCluster

    WorkerNodes1 --> LoggingMonitoring
    WorkerNodes2 --> LoggingMonitoring
```

**部署说明:**

*   **多集群 (可选)**: Karmada用于管理跨地域或跨可用区的多个K8S集群，实现高可用和灾备。Higress AI网关支持多集群部署和会话同步。
*   **服务部署**: 各微服务、AI平台、数据库、缓存等均以Docker容器形式运行在K8S Pod中。
*   **TiDB部署**: 采用TiDB Operator在K8S中部署和管理TiDB集群，保证其高可用性和可扩展性。
*   **Redis部署**: 可采用Redis Sentinel或Cluster模式部署在K8S中。
*   **负载均衡**: 对外暴露的服务通过负载均衡器分发流量到Higress AI网关。
*   **数据同步**: TiDB自身支持分布式事务和数据同步。对于文件存储，可采用MinIO等支持多副本或跨集群同步的方案。

## 3. 技术选型

### 3.1 后端技术栈

*   **编程语言**: **Golang** (高性能、并发友好、生态完善)
*   **Web框架**: (可选) `Gin` 或 `Echo` (轻量级、高性能的Go Web框架) 或标准库`net/http`
*   **ORM框架**: (可选) `GORM` (功能丰富的ORM库) 或直接使用 `database/sql` 进行原生SQL操作
*   **API规范**: RESTful API 或 gRPC (微服务间通信)
*   **依赖管理**: Go Modules

### 3.2 前端技术栈

*   **基础技术**: **HTML5, CSS3**
*   **CSS框架**: **Tailwind CSS** (实用优先的CSS框架，快速构建界面)
*   **图标库**: **FontAwesome** (丰富的矢量图标库)
*   **JavaScript**: 原生JavaScript或轻量级库 (如Alpine.js，可选，用于增强简单交互)

### 3.3 数据库与缓存

*   **主数据库**: **TiDB** (分布式HTAP数据库，兼容MySQL协议，支持水平扩展和数据同步)
*   **缓存数据库**: **Redis** (高性能键值存储，用于用户会话、热点数据缓存、消息队列等)

### 3.4 AI与智能处理

*   **AI应用开发平台**: **Dify** (用于AI智能体的开发、Prompt编排、数据集管理、模型微调与部署)
*   **AI大模型**: 本地部署的多个**国产AI大模型** (具体型号根据项目需求确定)
*   **辅助工具**: **MCP工具集** (用于辅助AI进行数据库查询、执行特定规则等)

### 3.5 API网关与服务治理

*   **API网关**: **Higress AI 网关** (基于Envoy，提供服务发现、路由、安全、可观测性，并特别强化了对AI场景的支持，如多模型管理、Token限流、内容安全等)

### 3.6 容器化与编排

*   **容器技术**: **Docker**
*   **容器编排**: **Kubernetes (K8S)**，具体发行版为 **RKE2**
*   **多集群管理**: **Karmada** (用于管理多个K8S集群)

### 3.7 操作系统

*   **服务器操作系统**: **龙蜥Anolis OS** (国产化Linux发行版)

### 3.8 其他关键组件

*   **文件存储**: MinIO (开源对象存储) 或其他兼容S3协议的存储服务。
*   **日志管理**: ELK Stack (Elasticsearch, Logstash, Kibana) 或 EFK Stack (Elasticsearch, Fluentd, Kibana)。
*   **监控告警**: Prometheus, Grafana, Alertmanager。
*   **CI/CD**: Jenkins, GitLab CI, ArgoCD等。
*   **消息队列 (可选)**: Kafka, RabbitMQ 或 NATS (用于服务解耦、异步处理)。

### 3.9 技术选型理由总结

*   **国产化与自主可控**: 优先选择国产化组件 (Anolis OS, TiDB, Dify, 本地大模型, Higress) 满足项目特定要求。
*   **高性能与可扩展性**: Golang、TiDB、K8S等技术能很好地支持高并发和大规模部署。
*   **云原生架构**: 基于Docker、K8S、微服务构建，易于部署、运维和弹性伸缩。
*   **AI能力整合**: Dify、本地大模型、Higress AI网关共同构成了强大的AI处理能力。
*   **成熟生态**: 所选技术大多拥有成熟的社区和丰富的生态资源。
*   **前端简洁高效**: HTML + Tailwind CSS + FontAwesome 能够快速构建简单且美观的前端页面，符合PRD中对前端的定位。

## 4. 架构决策与考虑

*   **微服务拆分粒度**: 初期可以按照核心业务领域（用户、贷款、农机、后台管理）进行拆分。随着业务发展，可以进一步细化。
*   **服务间通信**: API网关与后端服务间可采用HTTP/REST。服务内部间通信可考虑gRPC以提升性能。
*   **数据一致性**: 分布式事务是微服务架构中的难点。初期可采用最终一致性方案（如基于消息队列的事件驱动架构），对于强一致性场景审慎设计。TiDB自身支持分布式事务，可以简化部分场景。
*   **安全性**: 安全是金融系统的重中之重。从网络层、应用层到数据层都需要全面的安全防护措施。Higress AI网关提供了重要的安全能力。
*   **可观测性**: 完善的日志、监控和告警体系对于快速定位和解决问题至关重要。
*   **开发效率**: 虽然前端技术栈简化，但后端微服务架构需要良好的DevOps实践支持。

## 5. 未来演进方向

*   引入更完善的服务网格 (如Istio，尽管Higress已具备部分能力)。
*   深化AI在更多场景的应用，如智能客服、精准营销、贷后风险监控。
*   构建数据中台，更好地管理和利用业务数据。
*   扩展到更多农业服务领域。 