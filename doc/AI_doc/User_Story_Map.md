# 用户故事地图 (User Story Map): 数字惠农APP及OA后台管理系统

## 1. 用户故事地图概述

本用户故事地图旨在通过可视化的方式，梳理"数字惠农APP及OA后台管理系统"的用户活动、任务和具体的用户故事，并将其与产品路线图中的版本规划相对应。这有助于团队理解用户需求的全貌，并确保开发过程聚焦于为用户创造价值。考虑到本项目为演示项目，故事的粒度和深度会更侧重于支撑核心演示流程。

## 2. 用户故事地图结构

- **横向轴 (用户活动 Backbone)**: 代表用户为达成其目标所进行的一系列高层次活动。
- **纵向轴 (用户任务 Task)**: 在每个用户活动下，分解出的具体用户任务。
- **用户故事 (Story)**: 对每个任务的进一步细化，通常采用 "作为 [用户类型], 我想要 [完成某事], 以便 [获得价值]" 的格式。
- **版本映射 (Release Slices)**: 将用户故事分配到不同的产品版本（MVP, V1.1等），体现迭代交付的思路。

## 3. 用户故事地图详解

```mermaid
graph LR
    subgraph 用户App/网页端
        subgraph 用户活动: 账户管理 (Activity_AccMgt)
            direction LR
            Task_Reg[任务: 注册账号] --> Story_Reg_Email(作为新用户, 我想用手机号注册, 以便使用APP功能 - MVP V0)
            Task_Login[任务: 登录系统] --> Story_Login_User(作为注册用户, 我想用账号密码登录, 以便访问我的信息 - MVP V0)
            Task_Profile[任务: 管理个人资料] --> Story_ViewProfile(作为用户, 我想查看我的基本资料, 如昵称头像 - MVP V1)
            Task_Profile --> Story_EditProfile(作为用户, 我想修改我的昵称或头像 - V1.1 P1)
        end

        subgraph 用户活动: 惠农贷服务 (Activity_Loan)
            direction LR
            Task_Loan_Browse[任务: 浏览贷款产品] --> Story_Loan_List(作为农户, 我想浏览可选的惠农贷产品列表, 了解基本信息 - MVP V0)
            Task_Loan_Browse --> Story_Loan_Detail(作为农户, 我想查看特定贷款产品的详细说明和申请条件 - MVP V0)
            Task_Loan_Apply[任务: 申请贷款] --> Story_Loan_FillForm(作为农户, 我想在线填写贷款申请表单, 包含个人和财务信息 - MVP V0)
            Task_Loan_Apply --> Story_Loan_UploadDocs(作为农户, 我想方便地上传身份证、收入证明等所需文件 - MVP V0)
            Task_Loan_Apply --> Story_Loan_Submit(作为农户, 我想提交我的贷款申请 - MVP V0)
            Task_Loan_Track[任务: 跟踪申请进度] --> Story_Loan_ViewStatus(作为农户, 我想实时查询我的贷款申请状态 - MVP V0)
            Task_Loan_Track --> Story_Loan_ViewReason(作为农户, 若申请被拒, 我想查看拒绝原因 - MVP V0)
            Task_Loan_Track --> Story_Loan_ReceiveNotification(作为农户, 我想在申请状态变更时收到通知 - V1.1 P1)
        end

        subgraph 用户活动: 农机租赁服务 (Activity_Leasing)
            direction LR
            Task_Leasing_Browse[任务: 浏览/搜索农机 (承租方)] --> Story_Leasing_List(作为农户, 我想浏览可租赁的农机列表 - MVP V1)
            Task_Leasing_Browse --> Story_Leasing_Search(作为农户, 我想按类型或区域搜索农机 - V1.1 P0)
            Task_Leasing_Browse --> Story_Leasing_Detail(作为农户, 我想查看农机详情, 包括价格和图片 - MVP V1)
            Task_Leasing_Apply[任务: 申请租赁 (承租方)] --> Story_Leasing_SubmitRequest(作为农户, 我想向农机主提交租赁请求 - MVP V1)
            Task_Leasing_Publish[任务: 发布/管理农机 (出租方)] --> Story_Leasing_PublishNew(作为农机主, 我想发布我的农机信息以供出租 - V1.1 P0)
            Task_Leasing_Publish --> Story_Leasing_ManageMy(作为农机主, 我想管理我发布的农机信息 - V1.1 P0)
            Task_Leasing_Manage_Orders[任务: 管理租赁订单] --> Story_Leasing_ViewOrders_Lessee(作为承租方, 我想查看我的租赁订单状态 - V1.1 P0)
            Task_Leasing_Manage_Orders --> Story_Leasing_HandleRequest_Lessor(作为出租方, 我想处理收到的租赁请求 - V1.1 P0)
        end

        subgraph 用户活动: 政策与资讯 (Activity_Info)
            direction LR
            Task_Info_Browse[任务: 查看惠农政策] --> Story_Info_List(作为用户, 我想在首页看到最新的惠农政策公告列表 - MVP V0)
            Task_Info_Browse --> Story_Info_Detail(作为用户, 我想点击查看政策详情 - MVP V0)
        end
    end

    subgraph OA后台管理系统
        subgraph 用户活动: OA系统登录与导航 (Activity_OA_LoginNav)
            direction LR
            Task_OA_Login[任务: OA管理员登录] --> Story_OA_Login(作为审批员/管理员, 我想使用账号密码登录OA后台 - MVP V0)
            Task_OA_Navigate[任务: OA系统导航] --> Story_OA_Homepage(作为审批员/管理员, 我想看到OA系统首页/工作台, 包含主要功能入口和统计概览 - MVP V1)
        end

        subgraph 用户活动: 贷款智能审批 (Activity_OA_Approval)
            direction LR
            Task_OA_ViewDashboard[任务: 查看审批看板] --> Story_OA_PendingList(作为审批员, 我想在审批看板上看到待处理的贷款申请列表 - MVP V0)
            Task_OA_ViewDashboard --> Story_OA_FilterSort(作为审批员, 我想对待处理申请进行筛选和排序 - V1.1 P1)
            Task_OA_ProcessApp[任务: 处理单个申请] --> Story_OA_ViewAIDetails(作为审批员, 我想查看AI对单个申请的分析详情, 包括信息、验证、风险和建议 - MVP V0)
            Task_OA_ProcessApp --> Story_OA_ManualReview(作为审批员, 我想对需要人工复核的申请进行审批操作(通过/拒绝)并添加意见 - MVP V0)
            Task_OA_ControlAI[任务: 控制AI审批流程] --> Story_OA_ToggleAI(作为管理员, 我想通过控制面板开启或关闭整个智能审批流程 - MVP V0)
            Task_OA_ControlAI --> Story_OA_ToggleAISteps(作为管理员, 我想(演示目的)可以独立控制AI审批的某个环节 - V1.1 P1, 可选)
        end
        
        subgraph 用户活动: 系统与模型管理 (Activity_OA_SysMgt) - V1.2 技术演示
            direction LR
            Task_OA_ModelMgt[任务: AI模型管理演示] --> Story_OA_Higress_MultiModel(作为技术演示者, 我想展示Higress网关管理多个本地AI模型并实现路由切换 - V1.2 P0)
            Task_OA_AgentMgt[任务: AI智能体管理演示] --> Story_OA_Dify_AgentOrchestration(作为技术演示者, 我想展示Dify平台多个智能体协同工作的场景 - V1.2 P1)
        end
    end

    %% Styling for clarity
    classDef activity fill:#f9f,stroke:#333,stroke-width:2px;
    classDef task fill:#ccf,stroke:#333,stroke-width:2px;
    classDef story fill:#cfc,stroke:#333,stroke-width:1px;

    class Activity_AccMgt,Activity_Loan,Activity_Leasing,Activity_Info activity;
    class Task_Reg,Task_Login,Task_Profile task;
    class Story_Reg_Email,Story_Login_User,Story_ViewProfile,Story_EditProfile story;
    class Task_Loan_Browse,Task_Loan_Apply,Task_Loan_Track task;
    class Story_Loan_List,Story_Loan_Detail,Story_Loan_FillForm,Story_Loan_UploadDocs,Story_Loan_Submit,Story_Loan_ViewStatus,Story_Loan_ViewReason,Story_Loan_ReceiveNotification story;
    class Task_Leasing_Browse,Task_Leasing_Apply,Task_Leasing_Publish,Task_Leasing_Manage_Orders task;
    class Story_Leasing_List,Story_Leasing_Search,Story_Leasing_Detail,Story_Leasing_SubmitRequest,Story_Leasing_PublishNew,Story_Leasing_ManageMy,Story_Leasing_ViewOrders_Lessee,Story_Leasing_HandleRequest_Lessor story;
    class Task_Info_Browse task;
    class Story_Info_List,Story_Info_Detail story;

    class Activity_OA_LoginNav,Activity_OA_Approval,Activity_OA_SysMgt activity;
    class Task_OA_Login,Task_OA_Navigate task;
    class Story_OA_Login,Story_OA_Homepage story;
    class Task_OA_ViewDashboard,Task_OA_ProcessApp,Task_OA_ControlAI task;
    class Story_OA_PendingList,Story_OA_FilterSort,Story_OA_ViewAIDetails,Story_OA_ManualReview,Story_OA_ToggleAI,Story_OA_ToggleAISteps story;
    class Task_OA_ModelMgt,Task_OA_AgentMgt task;
    class Story_OA_Higress_MultiModel,Story_OA_Dify_AgentOrchestration story;
```

## 4. 故事优先级与版本映射说明

- **MVP V0 / MVP V1**: 代表这些故事是构成最小可行产品 (MVP) 的核心部分，其中V0表示最高优先级，V1次之，但都属于MVP范畴。这些故事将优先在MVP版本中实现，以确保核心演示流程的完整性。
- **V1.1 P0 / V1.1 P1**: 代表这些故事计划在V1.1版本（演示优化版）中实现，P0优先级高于P1。这些故事旨在优化用户体验、增强演示效果或初步引入扩展功能概念。
- **V1.2 P0 / V1.2 P1**: 代表这些故事计划在V1.2版本（技术特性深化版）中实现，用于更深入地演示特定技术点。
- **(概念)** / **(可选)**: 标记的故事可能更多是概念性的演示，或根据演示需求和时间灵活调整。

本用户故事地图将作为开发团队进行Sprint规划和任务分解的重要参考。随着项目的进展和对用户需求的更深入理解，本地图也会相应迭代更新。

--- 