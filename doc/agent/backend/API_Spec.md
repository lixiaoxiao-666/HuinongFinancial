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

## 9. AI智能体服务 (AIAgentService)

> **重要说明**: 以下接口专门为Dify AI工作流设计，用于实现智能审批功能。

### 9.1 获取审批申请信息 (供Dify工作流调用)

- **URL**: `/api/v1/ai-agent/applications/{application_id}/info`
- **Method**: `GET`
- **认证**: Required (AI Agent Token)
- **用途**: Dify工作流通过此接口获取待审批申请的完整信息，包括用户信息、申请详情、上传文件等
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "application_id": "loan_app_uuid_789",
    "product_info": {
      "product_id": "loan_prod_001",
      "name": "春耕助力贷",
      "category": "种植贷",
      "max_amount": 50000,
      "interest_rate_yearly": "4.5% - 6.0%"
    },
    "application_info": {
      "amount": 30000,
      "term_months": 12,
      "purpose": "购买化肥和种子",
      "submitted_at": "2024-03-10T10:00:00Z",
      "status": "SUBMITTED"
    },
    "applicant_info": {
      "user_id": "user_uuid_123",
      "real_name": "张三",
      "id_card_number": "310***********1234",
      "phone": "138****8000",
      "address": "XX省XX市XX村",
      "age": 35,
      "is_verified": true
    },
    "financial_info": {
      "annual_income": 80000,
      "existing_loans": 0,
      "credit_score": 750,
      "account_balance": 15000,
      "land_area": "10亩",
      "farming_experience": "10年"
    },
    "uploaded_documents": [
      {
        "doc_type": "id_card_front",
        "file_id": "file_uuid_001",
        "file_url": "https://example.com/files/id_front.jpg",
        "ocr_result": {
          "name": "张三",
          "id_number": "310***********1234",
          "address": "XX省XX市XX村"
        }
      },
      {
        "doc_type": "land_contract",
        "file_id": "file_uuid_002",
        "file_url": "https://example.com/files/land_contract.pdf",
        "extracted_info": {
          "land_area": "10亩",
          "contract_period": "30年",
          "location": "XX省XX市XX村"
        }
      }
    ],
    "external_data": {
      "credit_bureau_score": 750,
      "blacklist_check": false,
      "previous_loan_history": [],
      "land_ownership_verified": true
    }
  }
}
```

### 9.2 提交AI审批结果 (供Dify工作流调用)

- **URL**: `/api/v1/ai-agent/applications/{application_id}/ai-decision`
- **Method**: `POST`
- **认证**: Required (AI Agent Token)
- **用途**: Dify工作流完成分析后，通过此接口提交AI审批决策结果
- **Request Body**:
```json
{
  "ai_analysis": {
    "risk_level": "LOW",
    "risk_score": 0.25,
    "confidence_score": 0.92,
    "analysis_summary": "申请人信用记录良好，有稳定农业收入，土地承包合同有效，综合风险较低。",
    "detailed_analysis": {
      "identity_verification": {
        "status": "PASSED",
        "score": 1.0,
        "details": "身份证信息与申请信息一致，OCR识别准确"
      },
      "financial_assessment": {
        "status": "PASSED",
        "score": 0.85,
        "details": "年收入8万元，负债率低，还款能力较强"
      },
      "credit_evaluation": {
        "status": "PASSED",
        "score": 0.9,
        "details": "信用评分750分，无不良记录"
      },
      "document_verification": {
        "status": "PASSED",
        "score": 0.95,
        "details": "所有必要文件已提交且有效"
      },
      "fraud_detection": {
        "status": "PASSED",
        "score": 0.98,
        "details": "未发现欺诈风险"
      }
    },
    "risk_factors": [
      {
        "factor": "申请金额",
        "level": "LOW",
        "description": "申请金额30000元，在产品限额范围内"
      },
      {
        "factor": "还款能力",
        "level": "LOW", 
        "description": "年收入8万元，月还款压力较小"
      }
    ],
    "recommendations": [
      "建议批准申请",
      "可按申请金额全额批准",
      "建议采用等额本息还款方式"
    ]
  },
  "ai_decision": {
    "decision": "AUTO_APPROVED",
    "approved_amount": 30000,
    "approved_term_months": 12,
    "suggested_interest_rate": "5.2%",
    "conditions": [],
    "next_action": "DIRECT_APPROVAL"
  },
  "processing_info": {
    "ai_model_version": "v1.2.0",
    "processing_time_ms": 2500,
    "workflow_id": "dify_workflow_001",
    "processed_at": "2024-03-10T10:05:30Z"
  }
}
```
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "AI审批结果已成功处理",
  "data": {
    "application_id": "loan_app_uuid_789",
    "new_status": "AI_APPROVED",
    "next_step": "AWAIT_FINAL_CONFIRMATION"
  }
}
```

### 9.3 触发AI审批工作流

- **URL**: `/api/v1/ai-agent/applications/{application_id}/trigger-workflow`
- **Method**: `POST`
- **认证**: Required (System Token)
- **用途**: 当用户提交申请后，系统调用此接口触发Dify AI工作流
- **Request Body**:
```json
{
  "workflow_type": "LOAN_APPROVAL",
  "priority": "NORMAL",
  "callback_url": "https://backend.example.com/api/v1/ai-agent/applications/{application_id}/ai-decision"
}
```
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "AI审批工作流已启动",
  "data": {
    "workflow_execution_id": "dify_exec_uuid_456",
    "estimated_completion_time": "2024-03-10T10:08:00Z",
    "status": "RUNNING"
  }
}
```

### 9.4 获取AI模型配置信息 (供Dify调用)

- **URL**: `/api/v1/ai-agent/config/models`
- **Method**: `GET`
- **认证**: Required (AI Agent Token)
- **用途**: Dify工作流获取当前可用的AI模型配置和规则参数
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "active_models": [
      {
        "model_id": "risk_assessment_v2",
        "model_type": "RISK_EVALUATION",
        "version": "2.1.0",
        "status": "ACTIVE",
        "thresholds": {
          "low_risk": 0.3,
          "medium_risk": 0.7,
          "high_risk": 0.9
        }
      },
      {
        "model_id": "fraud_detection_v1",
        "model_type": "FRAUD_DETECTION",
        "version": "1.5.2",
        "status": "ACTIVE",
        "sensitivity": 0.85
      }
    ],
    "approval_rules": {
      "auto_approval_threshold": 0.3,
      "auto_rejection_threshold": 0.8,
      "max_auto_approval_amount": 50000,
      "required_human_review_conditions": [
        "申请金额超过5万元",
        "信用评分低于700分",
        "存在潜在欺诈风险"
      ]
    },
    "business_parameters": {
      "max_debt_to_income_ratio": 0.5,
      "min_credit_score": 600,
      "max_loan_amount_by_category": {
        "种植贷": 50000,
        "设备贷": 200000,
        "其他": 30000
      }
    }
  }
}
```

### 9.5 查询外部数据 (供Dify调用)

- **URL**: `/api/v1/ai-agent/external-data/{user_id}`
- **Method**: `GET`
- **认证**: Required (AI Agent Token)
- **用途**: Dify工作流调用此接口获取用户的外部征信、银行流水等信息
- **Query Parameters**:
  - `data_types` (string): 需要查询的数据类型，用逗号分隔，如 "credit,bank_flow,blacklist"
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "user_id": "user_uuid_123",
    "credit_report": {
      "score": 750,
      "grade": "优秀",
      "report_date": "2024-03-01",
      "loan_history": [],
      "overdue_records": 0
    },
    "bank_flow": {
      "average_monthly_income": 6500,
      "account_stability": "稳定",
      "last_6_months_flow": [
        {"month": "2024-02", "income": 7200, "expense": 4800},
        {"month": "2024-01", "income": 6800, "expense": 5100}
      ]
    },
    "blacklist_check": {
      "is_blacklisted": false,
      "check_time": "2024-03-10T10:05:00Z"
    },
    "government_subsidy": {
      "received_subsidies": [
        {"year": 2023, "type": "种粮补贴", "amount": 1200},
        {"year": 2022, "type": "农机购置补贴", "amount": 3000}
      ]
    }
  }
}
```

### 9.6 更新申请状态 (供系统内部调用)

- **URL**: `/api/v1/ai-agent/applications/{application_id}/status`
- **Method**: `PUT`
- **认证**: Required (System Token)
- **用途**: 系统内部更新申请状态，记录审批流程的各个阶段
- **Request Body**:
```json
{
  "status": "AI_REVIEWING",
  "operator": "AI_SYSTEM",
  "remarks": "AI智能体正在分析申请信息",
  "metadata": {
    "workflow_execution_id": "dify_exec_uuid_456",
    "stage": "RISK_ASSESSMENT",
    "progress": 60
  }
}
```
- **Response (200 OK)**:
```json
{
  "code": 0,
  "message": "申请状态更新成功",
  "data": null
}
```

## 10. AI智能体认证说明

### 10.1 AI Agent Token

AI智能体使用专用的Token进行认证，具有以下特点：

- **Token类型**: `AI-Agent-Token`
- **Header格式**: `Authorization: AI-Agent-Token <token>`
- **有效期**: 长期有效（可配置刷新策略）
- **权限范围**: 仅限智能体相关接口
- **安全措施**: IP白名单限制、请求频率限制

### 10.2 Token获取方式

```bash
# 示例：通过系统配置或环境变量获取
AI_AGENT_TOKEN="ai_agent_secure_token_123456"
```

### 10.3 工作流认证示例

```bash
# Dify工作流调用示例
curl -X GET "http://localhost:8080/api/v1/ai-agent/applications/loan_app_uuid_789/info" \
  -H "Authorization: AI-Agent-Token ai_agent_secure_token_123456" \
  -H "Content-Type: application/json"
```

## 11. 错误码扩展

| 错误码 | HTTP状态码 | 说明 |
|--------|------------|------|
| 9001   | 400        | AI分析参数错误 |
| 9002   | 401        | AI Agent Token无效 |
| 9003   | 403        | AI Agent权限不足 |
| 9004   | 404        | 申请信息不存在 |
| 9005   | 409        | 申请状态冲突，无法处理 |
| 9006   | 500        | AI模型调用失败 |
| 9007   | 503        | AI服务暂时不可用 |

# AI智能体接口API文档

## 概述

本文档描述了慧农金融系统中AI智能体相关的API接口，用于支持Dify AI工作流与后端系统的集成。这些接口实现了自动化的贷款审批流程，包括数据获取、AI决策处理和状态更新。

## 认证方式

### AI Agent Token认证
用于Dify等外部AI服务调用：
```
Authorization: AI-Agent-Token your_ai_agent_token
```

### System Token认证
用于系统内部调用：
```
Authorization: System-Token your_system_token
```

## 接口列表

### 1. 获取申请信息
**接口功能**：供Dify AI工作流调用，获取贷款申请的详细信息用于AI分析

**接口地址**：`GET /api/v1/ai-agent/applications/{application_id}/info`

**请求参数**：
- `application_id` (path, string, required) - 申请ID

**请求头**：
```
Authorization: AI-Agent-Token your_token
Content-Type: application/json
```

**响应示例**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "application_id": "app_20240301_001",
    "product_info": {
      "product_id": "SEED_LOAN_001",
      "name": "种植贷",
      "category": "种植贷",
      "max_amount": 50000,
      "interest_rate_yearly": "6.5%-8.5%"
    },
    "application_info": {
      "amount": 30000,
      "term_months": 12,
      "purpose": "购买种子和农药",
      "submitted_at": "2024-03-01T10:30:00Z",
      "status": "AI_TRIGGERED"
    },
    "applicant_info": {
      "user_id": "user_20240301_001",
      "real_name": "张三",
      "id_card_number": "310***********1234",
      "phone": "138****5678",
      "address": "XX省XX市XX村",
      "age": 35,
      "is_verified": true
    },
    "financial_info": {
      "annual_income": 80000,
      "existing_loans": 0,
      "credit_score": 750,
      "account_balance": 15000,
      "land_area": "10亩",
      "farming_experience": "10年"
    },
    "uploaded_documents": [
      {
        "doc_type": "id_card_front",
        "file_id": "file_001",
        "file_url": "https://example.com/files/id_front.jpg",
        "ocr_result": {
          "name": "张三",
          "id_number": "310***********1234",
          "address": "XX省XX市XX村"
        }
      }
    ],
    "external_data": {
      "credit_bureau_score": 750,
      "blacklist_check": false,
      "previous_loan_history": [],
      "land_ownership_verified": true
    }
  }
}
```

### 2. 提交AI决策结果
**接口功能**：接收Dify AI工作流分析后的决策结果

**接口地址**：`POST /api/v1/ai-agent/applications/{application_id}/ai-decision`

**请求参数**：
- `application_id` (path, string, required) - 申请ID

**请求体**：
```json
{
  "ai_analysis": {
    "risk_level": "LOW",
    "risk_score": 0.25,
    "confidence_score": 0.92,
    "analysis_summary": "申请人信用状况良好，还款能力强",
    "detailed_analysis": {
      "income_analysis": "年收入8万元，稳定",
      "credit_analysis": "征信良好，无不良记录",
      "asset_analysis": "拥有10亩土地，资产充足"
    },
    "risk_factors": [
      {
        "factor": "credit_score",
        "value": 750,
        "weight": 0.3,
        "risk_contribution": 0.05
      }
    ],
    "recommendations": [
      "建议批准贷款",
      "可给予优惠利率"
    ]
  },
  "ai_decision": {
    "decision": "AUTO_APPROVED",
    "approved_amount": 30000,
    "approved_term_months": 12,
    "suggested_interest_rate": "6.5%",
    "conditions": [],
    "next_action": "AWAIT_FINAL_CONFIRMATION"
  },
  "processing_info": {
    "ai_model_version": "v2.1.0",
    "processing_time_ms": 1500,
    "workflow_id": "workflow_001",
    "processed_at": "2024-03-01T10:35:00Z"
  }
}
```

**响应示例**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "application_id": "app_20240301_001",
    "new_status": "AI_APPROVED",
    "next_step": "AWAIT_FINAL_CONFIRMATION"
  }
}
```

### 3. 触发AI审批工作流
**接口功能**：系统内部调用，触发AI审批流程

**接口地址**：`POST /api/v1/ai-agent/applications/{application_id}/trigger-workflow`

**请求参数**：
- `application_id` (path, string, required) - 申请ID

**请求体**：
```json
{
  "workflow_type": "LOAN_APPROVAL",
  "priority": "NORMAL",
  "callback_url": "https://example.com/api/callback"
}
```

**响应示例**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "workflow_execution_id": "exec_001",
    "estimated_completion_time": "2024-03-01T10:40:00Z",
    "status": "RUNNING"
  }
}
```

### 4. 获取AI模型配置
**接口功能**：获取当前活跃的AI模型配置信息

**接口地址**：`GET /api/v1/ai-agent/config/models`

**响应示例**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "active_models": [
      {
        "model_id": "risk_assessment_v2",
        "model_type": "RISK_EVALUATION",
        "version": "2.1.0",
        "status": "ACTIVE",
        "thresholds": {
          "low_risk": 0.3,
          "medium_risk": 0.7,
          "high_risk": 0.9
        }
      }
    ],
    "approval_rules": {
      "auto_approval_threshold": 0.3,
      "auto_rejection_threshold": 0.8,
      "max_auto_approval_amount": 50000,
      "required_human_review_conditions": [
        "申请金额超过5万元",
        "信用评分低于700分",
        "存在潜在欺诈风险"
      ]
    },
    "business_parameters": {
      "max_debt_to_income_ratio": 0.5,
      "min_credit_score": 600,
      "max_loan_amount_by_category": {
        "种植贷": 50000,
        "设备贷": 200000,
        "其他": 30000
      }
    }
  }
}
```

### 5. 获取外部数据
**接口功能**：获取用户相关的外部数据（征信、银行流水等）

**接口地址**：`GET /api/v1/ai-agent/external-data/{user_id}`

**请求参数**：
- `user_id` (path, string, required) - 用户ID
- `data_types` (query, string, optional) - 数据类型，逗号分隔，默认："credit,bank_flow,blacklist"

**响应示例**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user_id": "user_20240301_001",
    "credit_report": {
      "score": 750,
      "grade": "优秀",
      "report_date": "2024-03-01",
      "loan_history": [],
      "overdue_records": 0
    },
    "bank_flow": {
      "average_monthly_income": 6500,
      "account_stability": "稳定",
      "last_6_months_flow": [
        {
          "month": "2024-02",
          "income": 7200,
          "expense": 4800
        }
      ]
    },
    "blacklist_check": {
      "is_blacklisted": false,
      "check_time": "2024-03-01T10:30:00Z"
    },
    "government_subsidy": {
      "received_subsidies": [
        {
          "year": 2023,
          "type": "种粮补贴",
          "amount": 1200
        }
      ]
    }
  }
}
```

### 6. 更新申请状态
**接口功能**：更新贷款申请状态

**接口地址**：`PUT /api/v1/ai-agent/applications/{application_id}/status`

**请求参数**：
- `application_id` (path, string, required) - 申请ID

**请求体**：
```json
{
  "status": "MANUAL_REVIEW_REQUIRED",
  "operator": "ai_system",
  "remarks": "需要人工审核",
  "metadata": {
    "reason": "高风险申请",
    "review_priority": "HIGH"
  }
}
```

**响应示例**：
```json
{
  "code": 200,
  "message": "success",
  "data": null
}
```

## 错误响应

所有接口在发生错误时返回统一的错误格式：

```json
{
  "code": 400,
  "message": "参数错误",
  "data": null
}
```

### 错误代码说明

- `400` - 请求参数错误
- `401` - 认证失败
- `403` - 权限不足
- `404` - 资源不存在
- `500` - 服务器内部错误

## 状态码说明

### 申请状态 (Application Status)
- `SUBMITTED` - 已提交
- `AI_TRIGGERED` - AI处理中
- `AI_APPROVED` - AI自动通过
- `AI_REJECTED` - AI自动拒绝
- `MANUAL_REVIEW_REQUIRED` - 需人工审核
- `APPROVED` - 最终通过
- `REJECTED` - 最终拒绝

### AI决策类型 (AI Decision)
- `AUTO_APPROVED` - 自动通过
- `AUTO_REJECTED` - 自动拒绝
- `REQUIRE_HUMAN_REVIEW` - 需要人工审核

### 风险等级 (Risk Level)
- `LOW` - 低风险
- `MEDIUM` - 中风险
- `HIGH` - 高风险

## 工作流程

1. **申请提交** → 用户提交贷款申请
2. **触发AI** → 系统调用"触发AI审批工作流"接口
3. **获取数据** → Dify工作流调用"获取申请信息"接口
4. **AI分析** → Dify执行AI分析和决策
5. **提交决策** → Dify调用"提交AI决策结果"接口
6. **状态更新** → 系统根据AI决策更新申请状态

## 接口调用示例

### Dify工作流调用示例

```bash
# 1. 获取申请信息
curl -X GET "https://api.example.com/api/v1/ai-agent/applications/app_001/info" \
  -H "Authorization: AI-Agent-Token your_token"

# 2. 提交AI决策
curl -X POST "https://api.example.com/api/v1/ai-agent/applications/app_001/ai-decision" \
  -H "Authorization: AI-Agent-Token your_token" \
  -H "Content-Type: application/json" \
  -d '{
    "ai_analysis": {...},
    "ai_decision": {...},
    "processing_info": {...}
  }'
```

### 系统内部调用示例

```bash
# 触发工作流
curl -X POST "https://api.example.com/api/v1/ai-agent/applications/app_001/trigger-workflow" \
  -H "Authorization: System-Token your_token" \
  -H "Content-Type: application/json" \
  -d '{
    "workflow_type": "LOAN_APPROVAL",
    "priority": "NORMAL"
  }'
``` 