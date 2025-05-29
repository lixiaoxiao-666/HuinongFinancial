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

# AI智能体接口

## 概述

本系统的AI智能体接口支持**统一的多类型审批工作流**，能够自动识别并处理不同类型的申请：
- **贷款申请**：传统的金融贷款审批
- **农机租赁申请**：农业机械设备租赁审批

通过统一的接口设计，Dify工作流可以使用相同的接口处理不同类型的申请，系统会根据申请ID自动识别类型并调用相应的业务逻辑。

## 核心特性

- 🔄 **统一接口设计**：一套接口处理多种申请类型
- 🤖 **智能类型识别**：根据申请ID前缀自动判断申请类型
- 📊 **完整日志记录**：所有AI操作都有详细的审计日志
- 🔒 **数据脱敏保护**：自动脱敏敏感信息（姓名、身份证、手机号）
- ⚡ **高性能处理**：支持并发处理多个审批请求
- 🎯 **决策可追溯**：每个AI决策都有完整的分析链路

## 申请ID格式规范

| 申请类型 | ID前缀 | 示例 | 说明 |
|---------|--------|------|------|
| 贷款申请 | `app_`, `test_app_`, `loan_` | `test_app_001` | 传统格式或以app_开头 |
| 农机租赁 | `ml_`, `leasing_` | `ml_test_001` | 以ml_或leasing_开头 |

### 获取申请信息（统一接口）

**接口说明**：Dify工作流调用此接口获取申请的详细信息，自动识别贷款申请或农机租赁申请类型

- **URL**: `/ai-agent/applications/{application_id}/info`
- **Method**: `GET`
- **Headers**: 需要AI Agent Token认证
- **Parameters**:
  - `application_id` (path): 申请ID，支持贷款申请ID和农机租赁申请ID

**响应示例（贷款申请）**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "application_type": "LOAN_APPLICATION",
    "application_id": "test_app_001",
    "product_info": {
      "product_id": "product_001",
      "name": "农业种植贷",
      "category": "农业贷款",
      "max_amount": 500000,
      "interest_rate_yearly": "4.5%"
    },
    "application_info": {
      "amount": 100000,
      "term_months": 24,
      "purpose": "购买农业设备",
      "submitted_at": "2024-05-29T10:00:00Z",
      "status": "PENDING_AI_REVIEW"
    },
    "applicant_info": {
      "user_id": "user_001",
      "real_name": "张**",
      "id_card_number": "320***********1234",
      "phone": "138****0001",
      "address": "江苏省南京市玄武区农业示范园",
      "age": 35,
      "is_verified": true
    },
    "financial_info": {
      "annual_income": 200000,
      "existing_loans": 0,
      "credit_score": 750,
      "account_balance": 50000,
      "land_area": "10亩",
      "farming_experience": "8年"
    },
    "uploaded_documents": [
      {
        "doc_type": "身份证",
        "file_id": "file_001",
        "file_url": "/api/files/file_001",
        "ocr_result": {},
        "extracted_info": {}
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

**响应示例（农机租赁申请）**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "application_type": "MACHINERY_LEASING",
    "application_id": "ml_test_001",
    "lessee_info": {
      "user_id": "lessee_001",
      "real_name": "陈**",
      "id_card_number": "320***********1234",
      "phone": "138****2001",
      "address": "江苏省南京市玄武区农业合作社路28号",
      "occupation": "种植大户",
      "annual_income": 180000,
      "farming_experience": "8年",
      "credit_rating": "良好",
      "is_verified": true,
      "previous_leasing_count": 3
    },
    "lessor_info": {
      "user_id": "lessor_001",
      "real_name": "农**合作社",
      "phone": "025****1001",
      "business_name": "江苏农机服务合作社",
      "business_license": "91320***********",
      "verification_status": "已认证",
      "credit_rating": "AAA",
      "established_date": "2018-03-15",
      "average_rating": 4.8,
      "successful_leasing_count": 156,
      "total_machinery_count": 12
    },
    "machinery_info": {
      "machinery_id": "machinery_001",
      "type": "拖拉机",
      "brand_model": "约翰迪尔6B-1404",
      "engine_power": "140马力",
      "manufacturing_year": 2022,
      "condition": "优秀",
      "daily_rent": 500,
      "deposit": 2000,
      "location": "江苏省南京市六合区",
      "availability": true,
      "last_maintenance": "2024-05-01",
      "insurance_status": "已投保",
      "insurance_expiry": "2025-03-15"
    },
    "leasing_details": {
      "requested_start_date": "2024-06-01",
      "requested_end_date": "2024-06-10",
      "rental_days": 10,
      "total_amount": 5000,
      "deposit_amount": 2000,
      "usage_purpose": "玉米种植",
      "work_location": "江苏省南京市六合区马鞍街道",
      "estimated_work_area": "50亩",
      "special_requirements": "需要配套播种机"
    },
    "risk_assessment": {
      "lessee_credit_score": 720,
      "lessor_reliability": "高",
      "machinery_condition": "优秀",
      "insurance_status": "已投保",
      "seasonal_risk": "低",
      "weather_forecast": "适宜作业",
      "regional_activity": "正常"
    }
  }
}
```

### 提交AI决策（统一接口）

**接口说明**：Dify工作流完成AI分析后，调用此接口提交决策结果，支持多种申请类型的统一处理

- **URL**: `/ai-agent/applications/{application_id}/decisions`
- **Method**: `POST`
- **Headers**: 需要AI Agent Token认证
- **Parameters**:
  - `application_id` (path): 申请ID

**请求参数（贷款申请）**：
```json
{
  "application_type": "LOAN_APPLICATION",
  "ai_analysis": {
    "risk_level": "LOW",
    "risk_score": 0.25,
    "confidence_score": 0.85,
    "analysis_summary": "申请人信用状况良好，收入稳定，风险较低",
    "detailed_analysis": {
      "credit_analysis": "信用分数750，无不良记录",
      "financial_analysis": "年收入20万，债务收入比0%",
      "agricultural_analysis": "具备8年农业经验，土地资源充足",
      "risk_factors": ["季节性收入波动"],
      "strengths": ["信用评分高", "无现有贷款", "农业经验丰富"]
    },
    "recommendations": ["建议监控季节性风险", "关注农业市场变化"]
  },
  "ai_decision": {
    "decision": "AUTO_APPROVED",
    "approved_amount": 100000,
    "approved_term_months": 24,
    "suggested_interest_rate": "4.5%",
    "conditions": ["定期还款", "保持良好信用"],
    "next_action": "生成贷款合同"
  },
  "processing_info": {
    "ai_model_version": "LLM-v4.0-unified",
    "processing_time_ms": 2500,
    "workflow_id": "dify-unified-structured-output",
    "processed_at": "2024-05-29T10:05:00Z"
  }
}
```

**请求参数（农机租赁申请）**：
```json
{
  "application_type": "MACHINERY_LEASING",
  "ai_analysis": {
    "risk_level": "LOW",
    "risk_score": 0.3,
    "confidence_score": 0.8,
    "analysis_summary": "承租方信用良好，出租方资质优秀，设备状况良好",
    "detailed_analysis": {
      "lessee_analysis": "农业经验丰富，历史租赁记录良好",
      "lessor_analysis": "资质认证完整，服务评价优秀",
      "equipment_analysis": "设备状况优秀，保险齐全",
      "risk_factors": ["天气变化风险"],
      "strengths": ["双方信用良好", "设备状况优秀", "保险覆盖完整"]
    },
    "recommendations": ["关注天气预报", "确保设备操作培训"]
  },
  "ai_decision": {
    "decision": "AUTO_APPROVE",
    "suggested_deposit": 2000,
    "approved_rental_terms": {
      "daily_rate": 500,
      "rental_period": "2024-06-01至2024-06-10"
    },
    "conditions": ["提供操作证明", "购买意外保险"],
    "next_action": "生成租赁合同"
  },
  "processing_info": {
    "ai_model_version": "LLM-v4.0-unified",
    "processing_time_ms": 2200,
    "workflow_id": "dify-unified-structured-output",
    "processed_at": "2024-05-29T10:05:00Z"
  }
}
```

**响应示例**：
```json
{
  "code": 0,
  "message": "AI决策提交成功",
  "data": {
    "application_id": "test_app_001",
    "application_type": "LOAN_APPLICATION",
    "old_status": "PENDING_AI_REVIEW",
    "new_status": "AUTO_APPROVED",
    "decision_id": "decision_001",
    "processed_at": "2024-05-29T10:05:00Z",
    "ai_operation_id": "ai_op_001"
  }
}
```

### 获取外部数据（多类型支持）

**接口说明**：获取用户的外部数据，支持不同申请类型所需的数据

- **URL**: `/ai-agent/external-data/{user_id}`
- **Method**: `GET`
- **Headers**: 需要AI Agent Token认证
- **Parameters**:
  - `user_id` (path): 用户ID
  - `data_types` (query): 数据类型列表，逗号分隔

**支持的数据类型**：
- `credit_report`: 征信报告
- `bank_flow`: 银行流水
- `blacklist_check`: 黑名单检查
- `government_subsidy`: 政府补贴
- `farming_qualification`: 农业资质（农机租赁专用）

### 获取AI模型配置（多类型支持）

**接口说明**：获取AI模型的配置信息，支持多种申请类型的业务规则

- **URL**: `/ai-agent/config/models`
- **Method**: `GET`
- **Headers**: 需要AI Agent Token认证

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "loan_approval": {
      "auto_approval_threshold": 0.3,
      "auto_rejection_threshold": 0.7,
      "max_auto_approval_amount": 500000,
      "required_human_review_conditions": [
        "申请金额超过50万",
        "信用分数低于600",
        "存在黑名单记录"
      ]
    },
    "machinery_leasing": {
      "auto_approval_threshold": 0.4,
      "auto_rejection_threshold": 0.7,
      "max_auto_approval_deposit": 10000,
      "required_human_review_conditions": [
        "押金超过1万",
        "设备价值超过50万",
        "承租方信用评级低于良好"
      ]
    },
    "risk_assessment_models": [
      {
        "model_id": "credit_risk_v2.1",
        "model_type": "信用风险评估",
        "version": "2.1",
        "status": "active"
      },
      {
        "model_id": "machinery_risk_v1.3",
        "model_type": "农机租赁风险评估",
        "version": "1.3", 
        "status": "active"
      }
    ]
  }
}
```

## 农机租赁专用接口

### 获取农机租赁申请信息（专用接口）

**接口说明**：专门用于农机租赁申请的信息获取，提供更详细的农机租赁相关数据

- **URL**: `/ai-agent/machinery-leasing/applications/{application_id}`
- **Method**: `GET`
- **Headers**: 需要AI Agent Token认证
- **Parameters**:
  - `application_id` (path): 农机租赁申请ID

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "application_id": "ml_test_001",
    "lessee_info": {
      "user_id": "lessee_001",
      "real_name": "陈**",
      "id_card_number": "320***********1234",
      "phone": "138****2001",
      "address": "江苏省南京市玄武区农业合作社路28号",
      "occupation": "种植大户",
      "annual_income": 180000,
      "farming_experience": "8年",
      "credit_rating": "良好",
      "is_verified": true,
      "previous_leasing_count": 3
    },
    "lessor_info": {
      "user_id": "lessor_001",
      "real_name": "农**合作社",
      "phone": "025****1001",
      "business_name": "江苏农机服务合作社",
      "verification_status": "已认证",
      "credit_rating": "AAA",
      "average_rating": 4.8,
      "successful_leasing_count": 156
    },
    "machinery_info": {
      "machinery_id": "machinery_001",
      "type": "拖拉机",
      "brand_model": "约翰迪尔6B-1404",
      "daily_rent": 500,
      "deposit": 2000,
      "location": "江苏省南京市六合区"
    },
    "leasing_details": {
      "requested_start_date": "2024-06-01",
      "requested_end_date": "2024-06-10",
      "rental_days": 10,
      "total_amount": 5000,
      "deposit_amount": 2000,
      "usage_purpose": "玉米种植"
    },
    "risk_assessment": {
      "lessee_credit_score": 720,
      "lessor_reliability": "高",
      "machinery_condition": "优秀",
      "insurance_status": "已投保"
    }
  }
}
```

### 提交农机租赁AI决策（专用接口）

**接口说明**：专门用于农机租赁的AI决策提交

- **URL**: `/ai-agent/machinery-leasing/applications/{application_id}/decisions`
- **Method**: `POST`
- **Headers**: 需要AI Agent Token认证

## AI操作日志接口

### 获取AI操作日志

**接口说明**：查询AI操作的详细日志，支持多种申请类型

- **URL**: `/ai-agent/logs`
- **Method**: `GET`
- **Headers**: 需要AI Agent Token认证
- **Parameters**:
  - `application_id` (query, optional): 申请ID
  - `application_type` (query, optional): 申请类型 (LOAN_APPLICATION/MACHINERY_LEASING)
  - `page` (query, optional): 页码，默认1
  - `limit` (query, optional): 每页数量，默认20

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "logs": [
      {
        "id": "ai_op_001",
        "application_id": "test_app_001",
        "application_type": "LOAN_APPLICATION",
        "operation_type": "DECISION_SUBMISSION",
        "ai_model_version": "LLM-v4.0-unified",
        "decision": "AUTO_APPROVED",
        "risk_score": 0.25,
        "confidence_score": 0.85,
        "processing_time_ms": 2500,
        "workflow_id": "dify-unified-structured-output",
        "created_at": "2024-05-29T10:05:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 1,
      "total_count": 1,
      "limit": 20
    }
  }
}
```

## AI模型配置接口

### 更新AI模型配置

**接口说明**：更新AI模型的业务配置，支持多种申请类型

- **URL**: `/ai-agent/config/models`
- **Method**: `PUT`
- **Headers**: 需要AI Agent Token认证

## 错误码说明

| 错误码 | 说明 | 处理建议 |
|-------|------|---------|
| 0 | 成功 | - |
| 400 | 请求参数错误 | 检查请求参数格式 |
| 401 | AI Agent Token无效 | 检查Token配置 |
| 404 | 申请不存在 | 确认申请ID正确 |
| 422 | 申请状态不允许AI处理 | 检查申请当前状态 |
| 500 | 服务器内部错误 | 联系技术支持 |

## 认证说明

所有AI智能体接口都需要在请求头中包含AI Agent Token：

```http
Authorization: AI-Agent-Token your_token_here
```

Token配置位置：
- 环境变量：`AI_AGENT_TOKEN`
- 配置文件：`configs/config.yaml`

## 数据脱敏说明

为保护用户隐私，API响应中的敏感信息会自动脱敏：

| 字段类型 | 脱敏规则 | 示例 |
|---------|---------|------|
| 姓名 | 保留姓，其他用*代替 | 张三 → 张* |
| 身份证号 | 保留前3位和后4位 | 320123199001011234 → 320***********1234 |
| 手机号 | 保留前3位和后4位 | 13812345678 → 138****5678 |

## 性能指标

| 指标 | 目标值 | 说明 |
|------|-------|------|
| 响应时间 | < 500ms | 获取申请信息接口 |
| 响应时间 | < 1000ms | 提交AI决策接口 |
| 并发支持 | 100+ QPS | 高峰期处理能力 |
| 可用性 | 99.9% | 服务可用性目标 |

---

**文档版本**：v4.0-统一多类型版  
**最后更新**：2024年5月29日  
**适用范围**：慧农金融后端API v1.0+ 