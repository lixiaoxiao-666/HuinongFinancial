# 数字惠农OA管理系统 - 后端API文档

## 系统概述

数字惠农OA管理系统是一个专为农业金融贷款审批而设计的后台管理系统，主要功能包括：

### 核心功能模块

1. **🔐 用户认证与权限管理**
   - OA用户登录认证
   - 基于JWT的会话管理
   - 角色权限控制（管理员/审批员）
   - 用户状态管理

2. **📊 工作台/首页**
   - 系统统计数据展示
   - 个人待办事项列表
   - 快捷操作菜单
   - 最近活动记录

3. **📋 贷款审批管理**
   - 待审批申请列表查询
   - 申请详情查看（含AI分析）
   - 人工审批决策
   - 申请历史记录追踪
   - 文件附件管理

4. **⚙️ 系统管理**
   - AI审批功能开关控制
   - 系统统计信息
   - 配置参数管理
   - 实时监控数据

5. **👥 用户管理**
   - OA用户创建
   - 用户状态控制
   - 角色权限分配
   - 用户列表查询

6. **📝 操作日志**
   - 全面的操作审计
   - 日志查询和筛选
   - 操作追溯
   - 合规性记录

### 技术特性

- **智能审批流程**：集成AI风险评估和人工审核的混合审批机制
- **权限安全控制**：基于角色的精细化权限管理
- **实时数据统计**：动态的系统指标和业务数据展示
- **操作审计追踪**：完整的操作日志记录
- **高性能设计**：支持高并发访问和大数据量处理

### 默认账号

**系统管理员**
- 用户名：`admin`
- 密码：`admin123`
- 权限：全部管理功能

**审批员**
- 用户名：`reviewer`
- 密码：`reviewer123`
- 权限：审批相关功能

## API接口文档

## 1. 概述

本文档定义了"数字惠农APP及OA后台管理系统"后端服务的所有API接口。

**服务信息:**
- **服务名称**: digital-agriculture-backend
- **版本**: v1.0.0
- **Base URL**: `http://localhost:8080/api/v1`
- **环境**: development

**通用约定:**
- **认证**: 需要Token认证的接口在HTTP Header中传递 `Authorization: Bearer <token>`
- **请求格式**: JSON (`Content-Type: application/json`)
- **响应格式**: JSON
- **成功响应**: HTTP状态码 `200 OK` 或 `201 Created`
- **错误响应**: HTTP状态码 `4xx` (客户端错误) 或 `5xx` (服务器错误)

**响应格式:**
```json
{
  "code": 0,           // 0表示成功，其他表示错误码
  "message": "Success", // 响应消息
  "data": { ... }      // 业务数据
}
```

**分页响应格式:**
```json
{
  "code": 0,
  "message": "Success",
  "data": [...],       // 数据列表
  "total": 100         // 总数量
}
```

## 2. 健康检查接口

### 2.1 健康检查

- **URL**: `/health`
- **Method**: `GET`
- **认证**: 无需认证
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "status": "ok",
    "service": "digital-agriculture-backend",
    "version": "1.0.0"
  }
}
```

## 3. 用户服务 (UserService)

### 3.1 用户注册

- **URL**: `/api/v1/users/register`
- **Method**: `POST`
- **认证**: 无需认证
- **Request Body**:
```json
{
  "phone": "13800138000",
  "password": "your_password"
}
```
- **Response (201 Created)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "user_id": "user_uuid_123"
  }
}
```

### 3.2 用户登录

- **URL**: `/api/v1/users/login`
- **Method**: `POST`
- **认证**: 无需认证
- **Request Body**:
```json
{
  "phone": "13800138000",
  "password": "your_password"
}
```
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "user_id": "user_uuid_123",
    "token": "jwt_auth_token_string",
    "expires_in": 7200,
    "user_type": "user"
  }
}
```

### 3.3 获取用户信息

- **URL**: `/api/v1/users/me`
- **Method**: `GET`
- **认证**: Required (Bearer Token)
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "user_id": "user_uuid_123",
    "phone": "138****8000",
    "nickname": "农户小张",
    "avatar_url": "https://example.com/avatar.jpg",
    "real_name": "张三",
    "id_card_number": "310***",
    "address": "XX省XX市XX村"
  }
}
```

### 3.4 更新用户信息

- **URL**: `/api/v1/users/me`
- **Method**: `PUT`
- **认证**: Required (Bearer Token)
- **Request Body**:
```json
{
  "nickname": "农户大张",
  "avatar_url": "https://example.com/new_avatar.jpg",
  "real_name": "张三",
  "address": "XX省XX市XX村"
}
```
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "用户信息更新成功",
  "data": null
}
```

## 4. 贷款服务 (LoanService)

### 4.1 获取贷款产品列表

- **URL**: `/api/v1/loans/products`
- **Method**: `GET`
- **认证**: 无需认证
- **Query Parameters**:
  - `category` (string, optional): 产品分类
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": [
    {
      "product_id": "loan_prod_001",
      "name": "春耕助力贷",
      "description": "专为春耕生产设计，利率优惠，快速审批",
      "category": "种植贷",
      "min_amount": 5000,
      "max_amount": 50000,
      "min_term_months": 6,
      "max_term_months": 24,
      "interest_rate_yearly": "4.5% - 6.0%",
      "status": 0
    },
    {
      "product_id": "loan_prod_002",
      "name": "农机购置贷",
      "description": "支持农户购买农业机械，助力农业现代化",
      "category": "设备贷",
      "min_amount": 10000,
      "max_amount": 200000,
      "min_term_months": 12,
      "max_term_months": 60,
      "interest_rate_yearly": "5.0% - 7.0%",
      "status": 0
    }
  ]
}
```

### 4.2 获取贷款产品详情

- **URL**: `/api/v1/loans/products/{product_id}`
- **Method**: `GET`
- **认证**: 无需认证
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "product_id": "loan_prod_001",
    "name": "春耕助力贷",
    "description": "专为春耕生产设计，利率优惠，快速审批",
    "category": "种植贷",
    "min_amount": 5000,
    "max_amount": 50000,
    "min_term_months": 6,
    "max_term_months": 24,
    "interest_rate_yearly": "4.5% - 6.0%",
    "repayment_methods": ["等额本息", "先息后本"],
    "application_conditions": "1. 年满18周岁的农户；2. 有稳定的农业收入；3. 信用记录良好",
    "required_documents": [
      {"type": "ID_CARD", "desc": "申请人身份证"},
      {"type": "LAND_CONTRACT", "desc": "土地承包合同"}
    ],
    "status": 0
  }
}
```

### 4.3 提交贷款申请

- **URL**: `/api/v1/loans/applications`
- **Method**: `POST`
- **认证**: Required (Bearer Token)
- **Request Body**:
```json
{
  "product_id": "loan_prod_001",
  "amount": 30000,
  "term_months": 12,
  "purpose": "购买化肥和种子",
  "applicant_info": {
    "real_name": "张三",
    "id_card_number": "310...",
    "address": "XX省XX市XX村"
  },
  "uploaded_documents": [
    {"doc_type": "id_card_front", "file_id": "file_uuid_001"},
    {"doc_type": "land_contract", "file_id": "file_uuid_002"}
  ]
}
```
- **Response (201 Created)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "application_id": "loan_app_uuid_789"
  }
}
```

### 4.4 获取贷款申请详情

- **URL**: `/api/v1/loans/applications/{application_id}`
- **Method**: `GET`
- **认证**: Required (Bearer Token)
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "application_id": "loan_app_uuid_789",
    "product_id": "loan_prod_001",
    "user_id": "user_uuid_123",
    "amount": 30000,
    "term_months": 12,
    "purpose": "购买化肥和种子",
    "status": "AI_REVIEWING",
    "submitted_at": "2024-03-10T10:00:00Z",
    "updated_at": "2024-03-10T11:30:00Z",
    "approved_amount": null,
    "remarks": "AI系统正在分析您的申请信息",
    "history": [
      {"status": "SUBMITTED", "timestamp": "2024-03-10T10:00:00Z", "operator": "用户"},
      {"status": "AI_REVIEWING", "timestamp": "2024-03-10T10:05:00Z", "operator": "系统"}
    ]
  }
}
```

### 4.5 获取我的贷款申请列表

- **URL**: `/api/v1/loans/applications/my`
- **Method**: `GET`
- **认证**: Required (Bearer Token)
- **Query Parameters**:
  - `status` (string, optional): 按状态筛选
  - `page` (int, optional, default: 1)
  - `limit` (int, optional, default: 10)
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": [
    {
      "application_id": "loan_app_uuid_789",
      "product_name": "春耕助力贷",
      "amount": 30000,
      "status": "AI_REVIEWING",
      "submitted_at": "2024-03-10T10:00:00Z"
    }
  ],
  "total": 1
}
```

## 5. 文件服务 (FileService)

### 5.1 文件上传

- **URL**: `/api/v1/files/upload`
- **Method**: `POST`
- **认证**: Required (Bearer Token)
- **Content-Type**: `multipart/form-data`
- **Form Data**:
  - `file`: (file) 上传的文件
  - `purpose`: (string, optional) 文件用途，如 "loan_document", "machinery_image"
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "file_id": "file_uuid_xyz",
    "file_url": "https://example.com/files/file_uuid_xyz.jpg",
    "file_name": "身份证正面.jpg",
    "file_size": 102400
  }
}
```

## 6. OA后台管理服务 (AdminService)

### 6.1 OA用户登录

- **URL**: `/api/v1/admin/login`
- **Method**: `POST`
- **认证**: 无需认证
- **Request Body**:
```json
{
  "username": "admin",
  "password": "admin123"
}
```
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "admin_user_id": "oa_admin001",
    "username": "admin",
    "role": "ADMIN",
    "token": "admin_jwt_auth_token",
    "expires_in": 3600
  }
}
```

### 6.2 获取待审批贷款申请列表

- **URL**: `/api/v1/admin/loans/applications/pending`
- **Method**: `GET`
- **认证**: Required (Admin Bearer Token)
- **Query Parameters**:
  - `status_filter` (string, optional)
  - `applicant_name` (string, optional)
  - `application_id` (string, optional)
  - `page` (int, optional, default: 1)
  - `limit` (int, optional, default: 10)
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": []
}
```

### 6.3 获取贷款申请详情 (管理员视角)

- **URL**: `/api/v1/admin/loans/applications/{application_id}`
- **Method**: `GET`
- **认证**: Required (Admin Bearer Token)
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "application_id": "application_id"
  }
}
```

### 6.4 提交审批决策

- **URL**: `/api/v1/admin/loans/applications/{application_id}/review`
- **Method**: `POST`
- **认证**: Required (Admin Bearer Token)
- **Request Body**:
```json
{
  "decision": "approved",
  "approved_amount": 25000,
  "comments": "申请人信用良好，但考虑到当前负债，略微调整批准金额。",
  "required_info_details": null
}
```
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "审批决策提交成功",
  "data": null
}
```

### 6.5 控制AI审批流程开关

- **URL**: `/api/v1/admin/system/ai-approval/toggle`
- **Method**: `POST`
- **认证**: Required (Admin Bearer Token)
- **Request Body**:
```json
{
  "enabled": true
}
```
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "AI审批流程状态更新成功",
  "data": null
}
```

## 7. 错误码说明

| 错误码 | HTTP状态码 | 说明 |
|--------|------------|------|
| 0      | 200/201    | 成功 |
| 1001   | 400        | 请求参数错误 |
| 1002   | 401        | 未授权 |
| 1003   | 403        | 禁止访问 |
| 1004   | 404        | 资源不存在 |
| 1005   | 429        | 请求频率超限 |
| 5000   | 500        | 服务器内部错误 |

## 8. 注意事项

1. 所有日期时间格式遵循 ISO 8601 标准 (YYYY-MM-DDTHH:MM:SSZ)
2. 所有金额字段以分为单位的整数表示
3. 文件上传大小限制为 10MB
4. 支持的文件类型：jpg, jpeg, png, pdf, doc, docx
5. JWT Token 有效期为 24 小时
6. API 请求频率限制：每分钟 100 次 