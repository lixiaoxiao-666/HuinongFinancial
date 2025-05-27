# 数据设计文档 (Database_Spec.md)

## 1. 概述

本文档详细描述了"数字惠农APP及OA后台管理系统"的数据库设计，采用 **TiDB** 作为核心数据库。设计遵循关系型数据库范式，并针对分布式特性进行考量。

**命名约定:**

*   表名：小写下划线 (e.g., `users`, `loan_applications`)
*   字段名：小写驼峰 (e.g., `userId`, `applicationAmount`)
*   主键：通常为 `id` (BIGINT AUTO_RANDOM) 或业务相关的唯一ID (e.g., `applicationId` VARCHAR)
*   外键：`relatedTableId` (e.g., `userId`)
*   索引：明确定义普通索引、唯一索引和联合索引。
*   时间戳：使用 `createdAt`, `updatedAt` (DATETIME(3) 或 TIMESTAMP(3))。

## 2. 数据库表结构

### 2.1 用户表 (users)

存储APP端和网页端的用户信息。

| 字段名         | 类型             | 约束/索引                      | 描述                     | 示例值                  |
| -------------- | ---------------- | ------------------------------ | ------------------------ | ----------------------- |
| `id`           | BIGINT           | PRIMARY KEY, AUTO_RANDOM       | 用户唯一ID (系统生成)    | 1                       |
| `userId`       | VARCHAR(64)      | UNIQUE                         | 用户业务ID (UUID)        | `usr_abc123def456`      |
| `phone`        | VARCHAR(20)      | UNIQUE                         | 手机号码 (登录凭证)      | `13800138000`           |
| `passwordHash` | VARCHAR(255)     | NOT NULL                       | 哈希后的密码             | `bcrypt_hash_string`    |
| `nickname`     | VARCHAR(100)     |                                | 用户昵称                 | `农户小张`              |
| `avatarUrl`    | VARCHAR(512)     |                                | 用户头像URL              | `https://.../avatar.jpg` |
| `status`       | TINYINT          | NOT NULL, DEFAULT 0            | 用户状态 (0:正常, 1:禁用) | 0                       |
| `registeredAt` | DATETIME(3)      | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) | 注册时间                 | `2024-03-10 10:00:00.123` |
| `lastLoginAt`  | DATETIME(3)      |                                | 最后登录时间             | `2024-03-12 09:30:00.000` |
| `createdAt`    | DATETIME(3)      | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) | 记录创建时间             |                         |
| `updatedAt`    | DATETIME(3)      | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) | 记录更新时间             |                         |

**索引:**
*   `idx_phone` ON `phone`
*   `idx_userId` ON `userId` (UNIQUE)

### 2.2 用户画像/详情表 (user_profiles)

存储用户更详细的认证信息和画像数据，与 `users` 表一对一关联。

| 字段名             | 类型         | 约束/索引         | 描述                 | 示例值                      |
| ------------------ | ------------ | ----------------- | -------------------- | --------------------------- |
| `userId`           | VARCHAR(64)  | PRIMARY KEY       | 用户ID (外键关联 users.userId) | `usr_abc123def456`          |
| `realName`         | VARCHAR(100) |                   | 真实姓名             | `张三`                      |
| `idCardNumber`     | VARCHAR(30)  | INDEX             | 身份证号码           | `310...X`                   |
| `idCardFrontUrl`   | VARCHAR(512) |                   | 身份证正面图片URL    | `https://.../id_front.jpg`  |
| `idCardBackUrl`    | VARCHAR(512) |                   | 身份证背面图片URL    | `https://.../id_back.jpg`   |
| `address`          | VARCHAR(255) |                   | 联系地址             | `XX省XX市XX村XX组`         |
| `gender`           | TINYINT      |                   | 性别 (0:未知, 1:男, 2:女) | 1                           |
| `birthDate`        | DATE         |                   | 出生日期             | `1990-01-01`                |
| `occupation`       | VARCHAR(100) |                   | 职业                 | `种植户`                    |
| `annualIncome`     | DECIMAL(15,2)|                   | 年收入估算 (万元)    | `10.50`                     |
| `creditAuthAgreed` | BOOLEAN      | NOT NULL, DEFAULT 0 | 是否同意信用授权     | `true`                      |
| `createdAt`        | DATETIME(3)  | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) | 记录创建时间         |                             |
| `updatedAt`        | DATETIME(3)  | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) | 记录更新时间         |                             |

**索引:**
*   `idx_idCardNumber` ON `idCardNumber`

### 2.3 贷款产品表 (loan_products)

存储平台提供的贷款产品信息。

| 字段名               | 类型          | 约束/索引                      | 描述                     | 示例值                      |
| -------------------- | ------------- | ------------------------------ | ------------------------ | --------------------------- |
| `id`                 | BIGINT        | PRIMARY KEY, AUTO_RANDOM       | 产品唯一ID               | 1                           |
| `productId`          | VARCHAR(64)   | UNIQUE                         | 产品业务ID (UUID)        | `lp_spring2024`             |
| `name`               | VARCHAR(255)  | NOT NULL                       | 产品名称                 | `春耕助力贷`                |
| `description`        | TEXT          |                                | 产品详细描述             | `专为春季农业生产设计...`   |
| `category`           | VARCHAR(50)   | INDEX                          | 产品分类                 | `种植贷`                    |
| `minAmount`          | DECIMAL(15,2) | NOT NULL                       | 最低申请金额             | `5000.00`                   |
| `maxAmount`          | DECIMAL(15,2) | NOT NULL                       | 最高申请金额             | `50000.00`                  |
| `minTermMonths`      | INT           | NOT NULL                       | 最短贷款期限 (月)        | 6                           |
| `maxTermMonths`      | INT           | NOT NULL                       | 最长贷款期限 (月)        | 24                          |
| `interestRateYearly` | VARCHAR(50)   | NOT NULL                       | 年利率范围 (文本描述)    | `4.5% - 6.0%`               |
| `repaymentMethods`   | JSON          |                                | 支持的还款方式列表       | `["等额本息", "先息后本"]` |
| `applicationConditions`| TEXT          |                                | 申请条件说明             | `1. 年满18周岁...`          |
| `requiredDocuments`  | JSON          |                                | 所需材料清单 (类型和说明)| `[{"type": "ID_CARD", "desc": "申请人身份证"}]` |
| `status`             | TINYINT       | NOT NULL, DEFAULT 0            | 产品状态 (0:有效, 1:下架) | 0                           |
| `createdAt`          | DATETIME(3)   | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) | 记录创建时间             |                             |
| `updatedAt`          | DATETIME(3)   | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) | 记录更新时间             |                             |

**索引:**
*   `idx_productId` ON `productId` (UNIQUE)
*   `idx_category` ON `category`
*   `idx_status` ON `status`

### 2.4 贷款申请表 (loan_applications)

存储用户提交的贷款申请记录。

| 字段名              | 类型          | 约束/索引                      | 描述                     | 示例值                      |
| ------------------- | ------------- | ------------------------------ | ------------------------ | --------------------------- |
| `id`                | BIGINT        | PRIMARY KEY, AUTO_RANDOM       | 申请唯一ID (系统生成)    | 1                           |
| `applicationId`     | VARCHAR(64)   | UNIQUE                         | 申请业务ID (UUID)        | `la_xyz789abc123`           |
| `userId`            | VARCHAR(64)   | NOT NULL, INDEX                | 申请用户ID (外键 users.userId) | `usr_abc123def456`          |
| `productId`         | VARCHAR(64)   | NOT NULL, INDEX                | 申请产品ID (外键 loan_products.productId) | `lp_spring2024`             |
| `amountApplied`     | DECIMAL(15,2) | NOT NULL                       | 申请金额                 | `30000.00`                  |
| `termMonthsApplied` | INT           | NOT NULL                       | 申请期限 (月)            | 12                          |
| `purpose`           | VARCHAR(500)  |                                | 贷款用途说明             | `购买化肥和种子`            |
| `status`            | VARCHAR(50)   | NOT NULL, INDEX                | 申请状态                 | `AI_PROCESSING`             |
| `applicantSnapshot` | JSON          |                                | 提交时申请人信息快照     | `{ "realName": "张三", ... }` |
| `submittedAt`       | DATETIME(3)   | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) | 提交时间                 | `2024-03-10 10:00:00.123` |
| `aiRiskScore`       | INT           |                                | AI风险评分 (0-100)       | 75                          |
| `aiSuggestion`      | TEXT          |                                | AI审批建议               | `建议人工复核...`           |
| `approvedAmount`    | DECIMAL(15,2) |                                | 批准金额 (审批后)        | `28000.00`                  |
| `approvedTermMonths`| INT           |                                | 批准期限 (审批后)        | 12                          |
| `finalDecision`     | VARCHAR(50)   |                                | 最终审批决定             | `APPROVED`                  |
| `decisionReason`    | TEXT          |                                | 审批意见/拒绝原因        | `符合条件`                  |
| `processedBy`       | VARCHAR(64)   |                                | 处理人ID (OA用户)        | `oa_admin001`               |
| `processedAt`       | DATETIME(3)   |                                | 处理完成时间             | `2024-03-11 14:30:00.000` |
| `createdAt`         | DATETIME(3)   | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) | 记录创建时间             |                             |
| `updatedAt`         | DATETIME(3)   | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) | 记录更新时间             |                             |

**状态枚举 (status, finalDecision):**
*   `PENDING_SUBMISSION` (待用户提交完整)
*   `SUBMITTED` (已提交，待处理)
*   `AI_PROCESSING` (AI审批处理中)
*   `PENDING_MANUAL_REVIEW` (待人工复核)
*   `MANUAL_REVIEWING` (人工复核中)
*   `APPROVED` (已批准)
*   `REJECTED` (已拒绝)
*   `CANCELED_BY_USER` (用户取消)
*   `REQUIRE_MORE_INFO` (需补充材料)
*   `LOAN_DISBURSED` (已放款)
*   `REPAYING` (还款中)
*   `COMPLETED` (已结清)

**索引:**
*   `idx_applicationId` ON `applicationId` (UNIQUE)
*   `idx_userId_status` ON (`userId`, `status`)
*   `idx_productId_status` ON (`productId`, `status`)
*   `idx_status_submittedAt` ON (`status`, `submittedAt`)

### 2.5 贷款申请审批历史表 (loan_application_history)

记录贷款申请的每一个状态变更和操作历史。

| 字段名          | 类型        | 约束/索引                 | 描述                     | 示例值                      |
| --------------- | ----------- | ------------------------- | ------------------------ | --------------------------- |
| `id`            | BIGINT      | PRIMARY KEY, AUTO_RANDOM  | 历史记录ID               | 1                           |
| `applicationId` | VARCHAR(64) | NOT NULL, INDEX           | 贷款申请ID (外键 loan_applications.applicationId) | `la_xyz789abc123`           |
| `statusFrom`    | VARCHAR(50) |                           | 变更前状态               | `AI_PROCESSING`             |
| `statusTo`      | VARCHAR(50) | NOT NULL                  | 变更后状态               | `PENDING_MANUAL_REVIEW`     |
| `operatorType`  | VARCHAR(20) | NOT NULL                  | 操作者类型 (USER, SYSTEM, OA_USER) | `SYSTEM`                    |
| `operatorId`    | VARCHAR(64) |                           | 操作者ID (用户ID或OA用户ID) | `SYSTEM_AI`                 |
| `comments`      | TEXT        |                           | 操作备注/意见            | `AI评分75，转人工复核`      |
| `occurredAt`    | DATETIME(3) | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) | 发生时间                 | `2024-03-10 10:05:00.000` |

**索引:**
*   `idx_applicationId_occurredAt` ON (`applicationId`, `occurredAt`)

### 2.6 上传文件记录表 (uploaded_files)

存储用户上传的文件元信息。

| 字段名      | 类型        | 约束/索引                 | 描述                       | 示例值                        |
| ----------- | ----------- | ------------------------- | -------------------------- | ----------------------------- |
| `id`        | BIGINT      | PRIMARY KEY, AUTO_RANDOM  | 文件记录ID                 | 1                             |
| `fileId`    | VARCHAR(64) | UNIQUE                    | 文件业务ID (UUID)          | `file_doc001_uuid`            |
| `userId`    | VARCHAR(64) | NOT NULL, INDEX           | 上传用户ID (外键 users.userId) | `usr_abc123def456`            |
| `fileName`  | VARCHAR(255)| NOT NULL                  | 原始文件名                 | `身份证正面.jpg`              |
| `fileType`  | VARCHAR(50) |                           | 文件MIME类型               | `image/jpeg`                  |
| `fileSize`  | BIGINT      |                           | 文件大小 (bytes)           | 102400                        |
| `storagePath`| VARCHAR(512)| NOT NULL                  | 文件存储路径/URL           | `s3://bucket/path/to/file.jpg` |
| `purpose`   | VARCHAR(100)| INDEX                     | 文件用途 (如: LOAN_DOC, MACHINERY_IMG) | `LOAN_DOC_ID_CARD_FRONT`      |
| `relatedId` | VARCHAR(64) | INDEX                     | 关联业务ID (如: 贷款申请ID) | `la_xyz789abc123`             |
| `uploadedAt`| DATETIME(3) | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) | 上传时间                   | `2024-03-10 09:50:00.000`   |
| `createdAt` | DATETIME(3) | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) | 记录创建时间               |                               |
| `updatedAt` | DATETIME(3) | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) | 记录更新时间               |                               |

**索引:**
*   `idx_fileId` ON `fileId` (UNIQUE)
*   `idx_userId_purpose` ON (`userId`, `purpose`)
*   `idx_relatedId_purpose` ON (`relatedId`, `purpose`)

### 2.7 农机信息表 (farm_machinery)

存储农机主发布的农机信息。

| 字段名             | 类型          | 约束/索引                      | 描述                     | 示例值                      |
| ------------------ | ------------- | ------------------------------ | ------------------------ | --------------------------- |
| `id`               | BIGINT        | PRIMARY KEY, AUTO_RANDOM       | 农机唯一ID               | 1                           |
| `machineryId`      | VARCHAR(64)   | UNIQUE                         | 农机业务ID (UUID)        | `fm_tractor001`             |
| `ownerUserId`      | VARCHAR(64)   | NOT NULL, INDEX                | 机主用户ID (外键 users.userId) | `usr_owner789`              |
| `type`             | VARCHAR(100)  | NOT NULL, INDEX                | 农机类型 (如: 拖拉机, 收割机) | `拖拉机`                    |
| `brandModel`       | VARCHAR(255)  |                                | 品牌型号                 | `东方红 LX904`              |
| `description`      | TEXT          |                                | 详细描述                 | `90马力，四轮驱动...`       |
| `images`           | JSON          |                                | 图片文件ID列表 (关联 uploaded_files) | `["file_img001", "file_img002"]` |
| `dailyRent`        | DECIMAL(10,2) | NOT NULL                       | 日租金                   | `300.00`                    |
| `deposit`          | DECIMAL(10,2) |                                | 押金                     | `1000.00`                   |
| `locationText`     | VARCHAR(255)  |                                | 位置文字描述             | `XX省XX市XX农场`            |
| `locationGeo`      | POINT         | SPATIAL INDEX                  | 地理位置坐标 (可选)      | `POINT(116.3 39.9)`         |
| `status`           | VARCHAR(50)   | NOT NULL, DEFAULT 'AVAILABLE'  | 状态 (AVAILABLE, RENTED, MAINTENANCE) | `AVAILABLE`                 |
| `publishedAt`      | DATETIME(3)   |                                | 发布时间                 | `2024-03-12 10:00:00.000` |
| `createdAt`        | DATETIME(3)   | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) | 记录创建时间             |                             |
| `updatedAt`        | DATETIME(3)   | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) | 记录更新时间             |                             |

**索引:**
*   `idx_machineryId` ON `machineryId` (UNIQUE)
*   `idx_ownerUserId_status` ON (`ownerUserId`, `status`)
*   `idx_type_status` ON (`type`, `status`)
*   `idx_locationText` (FULLTEXT on `locationText` if needed for search)

### 2.8 农机租赁订单表 (machinery_leasing_orders)

存储农机租赁的订单信息。

| 字段名          | 类型          | 约束/索引                      | 描述                     | 示例值                      |
| --------------- | ------------- | ------------------------------ | ------------------------ | --------------------------- |
| `id`            | BIGINT        | PRIMARY KEY, AUTO_RANDOM       | 订单唯一ID               | 1                           |
| `orderId`       | VARCHAR(64)   | UNIQUE                         | 订单业务ID (UUID)        | `mlo_order456`              |
| `machineryId`   | VARCHAR(64)   | NOT NULL, INDEX                | 租赁农机ID (外键 farm_machinery.machineryId) | `fm_tractor001`             |
| `lesseeUserId`  | VARCHAR(64)   | NOT NULL, INDEX                | 承租用户ID (外键 users.userId) | `usr_lessee123`             |
| `lessorUserId`  | VARCHAR(64)   | NOT NULL, INDEX                | 出租用户ID (机主)        | `usr_owner789`              |
| `startDate`     | DATE          | NOT NULL                       | 租赁开始日期             | `2024-04-05`                |
| `endDate`       | DATE          | NOT NULL                       | 租赁结束日期             | `2024-04-07`                |
| `totalRent`     | DECIMAL(10,2) | NOT NULL                       | 总租金                   | `900.00`                    |
| `depositAmount` | DECIMAL(10,2) |                                | 支付的押金金额           | `1000.00`                   |
| `status`        | VARCHAR(50)   | NOT NULL, INDEX                | 订单状态                 | `PENDING_CONFIRMATION`      |
| `lesseeNotes`   | TEXT          |                                | 承租方备注               | `希望能提前半天取机`        |
| `lessorNotes`   | TEXT          |                                | 出租方备注               | `已消毒`                    |
| `confirmedAt`   | DATETIME(3)   |                                | 机主确认时间             |                             |
| `completedAt`   | DATETIME(3)   |                                | 订单完成时间             |                             |
| `createdAt`     | DATETIME(3)   | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) | 记录创建时间             |                             |
| `updatedAt`     | DATETIME(3)   | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) | 记录更新时间             |                             |

**状态枚举 (status):**
*   `PENDING_CONFIRMATION` (待机主确认)
*   `CONFIRMED` (已确认，待支付/待开始)
*   `REJECTED_BY_LESSOR` (机主拒绝)
*   `CANCELED_BY_LESSEE` (承租方取消)
*   `IN_PROGRESS` (租赁进行中)
*   `COMPLETED` (已完成)
*   `DISPUTED` (有争议)

**索引:**
*   `idx_orderId` ON `orderId` (UNIQUE)
*   `idx_machineryId_status` ON (`machineryId`, `status`)
*   `idx_lesseeUserId_status` ON (`lesseeUserId`, `status`)
*   `idx_lessorUserId_status` ON (`lessorUserId`, `status`)
*   `idx_startDate_endDate` ON (`startDate`, `endDate`)

### 2.9 OA后台用户表 (oa_users)

存储OA后台系统管理员和审批员的信息。

| 字段名         | 类型         | 约束/索引                      | 描述                       | 示例值                      |
| -------------- | ------------ | ------------------------------ | -------------------------- | --------------------------- |
| `id`           | BIGINT       | PRIMARY KEY, AUTO_RANDOM       | OA用户唯一ID               | 1                           |
| `oaUserId`     | VARCHAR(64)  | UNIQUE                         | OA用户业务ID (UUID)        | `oa_admin001`               |
| `username`     | VARCHAR(100) | UNIQUE                         | 登录用户名                 | `admin_approver`            |
| `passwordHash` | VARCHAR(255) | NOT NULL                       | 哈希后的密码               | `bcrypt_hash_string`        |
| `role`         | VARCHAR(50)  | NOT NULL, INDEX                | 角色 (如: APPROVER, ADMIN) | `APPROVER`                  |
| `displayName`  | VARCHAR(100) |                                | 显示名称                   | `审批员-王五`               |
| `email`        | VARCHAR(255) | UNIQUE                         | 电子邮箱                   | `wangwu@example.com`        |
| `status`       | TINYINT      | NOT NULL, DEFAULT 0            | 状态 (0:正常, 1:禁用)      | 0                           |
| `createdAt`    | DATETIME(3)  | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) | 记录创建时间               |                             |
| `updatedAt`    | DATETIME(3)  | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) | 记录更新时间               |                             |

**索引:**
*   `idx_oaUserId` ON `oaUserId` (UNIQUE)
*   `idx_username` ON `username` (UNIQUE)
*   `idx_role` ON `role`

### 2.10 系统配置表 (system_configurations) (可选)

存储系统级别的一些可配置参数。

| 字段名      | 类型        | 约束/索引     | 描述         | 示例值                                  |
| ----------- | ----------- | ------------- | ------------ | --------------------------------------- |
| `configKey` | VARCHAR(100)| PRIMARY KEY   | 配置项键     | `AI_APPROVAL_ENABLED`                   |
| `configValue`| TEXT        | NOT NULL      | 配置项值     | `true`                                  |
| `description`| VARCHAR(255)|               | 配置项描述   | `是否启用AI自动审批流程`                |
| `updatedAt` | DATETIME(3) | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) | 最后更新时间 |                                         |

## 3. 数据模型 E-R 图 (简化版，基于PRD)

```mermaid
erDiagram
    USERS ||--o{ LOAN_APPLICATIONS : "applies_for"
    USERS ||--o{ USER_PROFILES : "has_one"
    USERS ||--o{ FARM_MACHINERY : "owns"
    USERS ||--o{ MACHINERY_LEASING_ORDERS : "leases_as_lessee"
    USERS ||--o{ MACHINERY_LEASING_ORDERS : "rents_out_as_lessor"
    USERS ||--o{ UPLOADED_FILES : "uploads"

    LOAN_PRODUCTS ||--o{ LOAN_APPLICATIONS : "is_for_product"

    LOAN_APPLICATIONS ||--o{ LOAN_APPLICATION_HISTORY : "has_history"
    LOAN_APPLICATIONS ||--o{ UPLOADED_FILES : "has_related_document" (逻辑关联)

    FARM_MACHINERY ||--o{ MACHINERY_LEASING_ORDERS : "is_leased_in"
    FARM_MACHINERY ||--o{ UPLOADED_FILES : "has_image" (逻辑关联)

    OA_USERS ||--o{ LOAN_APPLICATIONS : "processes" (逻辑关联)
    OA_USERS ||--o{ LOAN_APPLICATION_HISTORY : "operates_on" (逻辑关联)

    USERS {
        varchar userId PK
        varchar phone UNIQUE
        varchar passwordHash
        varchar nickname
        datetime registeredAt
    }
    USER_PROFILES {
        varchar userId PK FK
        varchar realName
        varchar idCardNumber
    }
    LOAN_PRODUCTS {
        varchar productId PK
        varchar name
        decimal minAmount
        decimal maxAmount
    }
    LOAN_APPLICATIONS {
        varchar applicationId PK
        varchar userId FK
        varchar productId FK
        decimal amountApplied
        varchar status
        datetime submittedAt
    }
    LOAN_APPLICATION_HISTORY {
        bigint id PK
        varchar applicationId FK
        varchar statusTo
        datetime occurredAt
    }
    UPLOADED_FILES {
        varchar fileId PK
        varchar userId FK
        varchar fileName
        varchar storagePath
        varchar purpose
        varchar relatedId
    }
    FARM_MACHINERY {
        varchar machineryId PK
        varchar ownerUserId FK
        varchar type
        decimal dailyRent
        varchar status
    }
    MACHINERY_LEASING_ORDERS {
        varchar orderId PK
        varchar machineryId FK
        varchar lesseeUserId FK
        varchar lessorUserId FK
        date startDate
        date endDate
        varchar status
    }
    OA_USERS {
        varchar oaUserId PK
        varchar username UNIQUE
        varchar role
    }
```

## 4. 索引策略

*   **主键**: 优先使用 TiDB 的 `AUTO_RANDOM` 属性来打散主键写入热点，或者使用业务相关的UUID字符串作为主键。
*   **外键**: 为所有外键字段创建索引，以加速JOIN操作。
*   **高频查询字段**: 为经常出现在 `WHERE` 子句、`ORDER BY` 子句、`GROUP BY` 子句中的字段创建索引。
*   **组合索引**: 根据查询模式创建组合索引，注意索引列的顺序。
*   **唯一约束**: 使用唯一索引 (UNIQUE INDEX) 确保字段的唯一性 (如 `users.phone`, `loan_applications.applicationId`)。
*   **空间索引**: 对地理位置坐标字段 (如 `farm_machinery.locationGeo`) 使用空间索引 (SPATIAL INDEX)。
*   **TiDB特性**: 考虑使用 TiFlash (列式存储引擎) 进行HTAP场景的加速，对于某些分析查询可以创建异步复制到TiFlash的副本。

## 5. 数据完整性与约束

*   **非空约束 (NOT NULL)**: 关键字段设为NOT NULL。
*   **默认值 (DEFAULT)**: 为某些字段提供合理的默认值 (如 `status`, `createdAt`)。
*   **外键约束**: 虽然TiDB支持外键约束，但在大规模分布式系统中，有时会在应用层面保证数据一致性以获得更好的性能和灵活性。本项目初期可以考虑使用外键约束，后续根据性能评估调整。
*   **检查约束 (CHECK)**: TiDB支持检查约束，可用于更复杂的数据校验 (如金额必须大于0)。
*   **应用层校验**: 大部分业务逻辑校验应在应用层完成，减轻数据库压力。

## 6. 数据安全与备份

*   **数据加密**: TiDB支持透明数据加密 (TDE) 和传输层加密 (TLS)。
*   **备份与恢复**: 利用TiDB提供的备份恢复工具 (如BR - Backup & Restore) 制定定期备份策略。
*   **权限控制**: 严格控制数据库用户权限，遵循最小权限原则。

## 7. Redis 数据结构设计 (示例)

*   **用户会话**: `session:<session_id>` (STRING or HASH)
    *   Value: JSON字符串，包含 `userId`, `role`, `expire_time` 等。
*   **短信验证码**: `sms_code:<phone_number>` (STRING)
    *   Value: `<code>:<timestamp>` (设置TTL，如5分钟)
*   **热点贷款产品信息**: `loan_product_hot:<product_id>` (HASH or JSON String)
    *   缓存产品详情，减少DB查询。
*   **用户最近贷款申请状态缓存**: `user_loan_status:<user_id>:<application_id>` (STRING)
    *   Value: `status_code` (设置较短TTL)
*   **API访问频率限制**: `api_limit:<user_id>:<api_path>` (COUNTER with TTL)
*   **分布式锁**: `lock:<resource_name>` (STRING with SETNX and EXPIRE)
    *   用于防止并发操作关键资源。

本文档提供了数据库设计的核心内容，在开发过程中，可能需要根据具体实现和性能测试进行调整和优化。 