# 慧农金融API规范文档 v4.1

## 概述

慧农金融数字化农业金融服务系统后端API，提供用户管理、贷款服务、文件管理、审批管理和AI智能体服务。

**版本**: v4.1
**基础URL**: `http://localhost:8080`
**更新时间**: 2024-12-19

## 新增功能 (v4.1)

### 主要改进
- ✅ 完善了所有服务层实现
- ✅ 统一了错误处理和响应格式
- ✅ 增强了AI智能体服务功能
- ✅ 完善了文件服务实现
- ✅ 优化了路由和中间件配置
- ✅ 增加了完整的接口测试脚本

### 技术栈
- **框架**: Gin (Go Web Framework)
- **数据库**: PostgreSQL/MySQL (支持GORM)
- **认证**: JWT Token
- **日志**: Zap Logger
- **文档**: Swagger/OpenAPI

## 1. 健康检查

### 1.1 应用健康检查
```http
GET /health
```

**响应示例**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "status": "ok",
    "service": "digital-agriculture-backend",
    "version": "v1.0.0"
  }
}
```

### 1.2 Dify工作流健康检查
```http
GET /livez
GET /readyz
```

**响应示例**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "status": "healthy"
  }
}
```

## 2. 用户服务 (UserService)

### 2.1 用户注册
```http
POST /api/v1/users/register
Content-Type: application/json
```

**请求体**:
```json
{
  "username": "testuser",
  "password": "password123",
  "confirm_password": "password123",
  "real_name": "张三",
  "phone": "13800138000",
  "email": "test@example.com"
}
```

**响应 (201)**:
```json
{
  "code": 0,
  "message": "Created successfully",
  "data": {
    "user_id": "user_1703001234567"
  }
}
```

### 2.2 用户登录
```http
POST /api/v1/users/login
Content-Type: application/json
```

**请求体**:
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**响应 (200)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_info": {
      "user_id": "user_1703001234567",
      "username": "testuser",
      "real_name": "张三",
      "phone": "13800138000",
      "email": "test@example.com",
      "is_verified": false,
      "created_at": "2024-12-19T10:00:00Z",
      "updated_at": "2024-12-19T10:00:00Z"
    }
  }
}
```

### 2.3 获取用户信息
```http
GET /api/v1/users/me
Authorization: Bearer <token>
```

**响应 (200)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "user_id": "user_1703001234567",
    "username": "testuser",
    "real_name": "张三",
    "phone": "13800138000",
    "email": "test@example.com",
    "is_verified": false,
    "created_at": "2024-12-19T10:00:00Z",
    "updated_at": "2024-12-19T10:00:00Z"
  }
}
```

### 2.4 更新用户信息
```http
PUT /api/v1/users/me
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**:
```json
{
  "real_name": "张三三",
  "phone": "13800138001",
  "email": "test2@example.com",
  "id_card": "123456789012345678",
  "address": "某省某市某县某村"
}
```

## 3. 贷款服务 (LoanService)

### 3.1 获取贷款产品列表
```http
GET /api/v1/loans/products
GET /api/v1/loans/products?category=agriculture
```

**响应 (200)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": [
    {
      "product_id": "product_001",
      "name": "农业生产贷款",
      "description": "专为农业生产提供的低息贷款",
      "category": "agriculture",
      "min_amount": 10000,
      "max_amount": 1000000,
      "min_term_months": 6,
      "max_term_months": 36,
      "interest_rate_yearly": "4.8%",
      "repayment_methods": ["等额本息", "等额本金"],
      "status": 0
    }
  ]
}
```

### 3.2 获取贷款产品详情
```http
GET /api/v1/loans/products/{product_id}
```

**响应 (200)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "product_id": "product_001",
    "name": "农业生产贷款",
    "description": "专为农业生产提供的低息贷款",
    "category": "agriculture",
    "min_amount": 10000,
    "max_amount": 1000000,
    "min_term_months": 6,
    "max_term_months": 36,
    "interest_rate_yearly": "4.8%",
    "repayment_methods": ["等额本息", "等额本金"],
    "application_conditions": "年满18周岁，有稳定收入来源",
    "required_documents": [
      {
        "type": "id_card",
        "desc": "身份证正反面"
      },
      {
        "type": "income_proof",
        "desc": "收入证明"
      }
    ],
    "status": 0
  }
}
```

### 3.3 提交贷款申请
```http
POST /api/v1/loans/applications
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**:
```json
{
  "product_id": "product_001",
  "amount": 100000,
  "term_months": 12,
  "purpose": "农业生产资金",
  "applicant_info": {
    "real_name": "张三",
    "id_card_number": "123456789012345678",
    "address": "某省某市某县某村"
  },
  "uploaded_documents": [
    {
      "doc_type": "id_card",
      "file_id": "file_1703001234567"
    }
  ]
}
```

**响应 (201)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "application_id": "app_1703001234567"
  }
}
```

### 3.4 获取贷款申请详情
```http
GET /api/v1/loans/applications/{application_id}
Authorization: Bearer <token>
```

**响应 (200)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "application_id": "app_1703001234567",
    "product_id": "product_001",
    "user_id": "user_1703001234567",
    "amount": 100000,
    "term_months": 12,
    "purpose": "农业生产资金",
    "status": "SUBMITTED",
    "submitted_at": "2024-12-19T10:00:00Z",
    "updated_at": "2024-12-19T10:00:00Z",
    "approved_amount": null,
    "remarks": "",
    "history": [
      {
        "status": "SUBMITTED",
        "timestamp": "2024-12-19T10:00:00Z",
        "operator": "系统"
      }
    ]
  }
}
```

### 3.5 获取我的贷款申请列表
```http
GET /api/v1/loans/applications/my
GET /api/v1/loans/applications/my?status=SUBMITTED&page=1&limit=10
Authorization: Bearer <token>
```

**响应 (200)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": [
    {
      "application_id": "app_1703001234567",
      "product_name": "农业生产贷款",
      "amount": 100000,
      "status": "SUBMITTED",
      "submitted_at": "2024-12-19T10:00:00Z"
    }
  ],
  "total": 1
}
```

## 4. 文件服务 (FileService)

### 4.1 文件上传
```http
POST /api/v1/files/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**表单参数**:
- `file`: 文件 (必需)
- `file_type`: 文件类型 (必需) - `id_card`, `bank_flow`, `work_certificate`, `income_proof`, `other`
- `business_type`: 业务类型 (可选) - `loan`, `machinery_leasing`

**响应 (200)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "file_id": "file_1703001234567",
    "file_name": "identity_card.jpg",
    "file_size": 1024000,
    "uploaded_at": "2024-12-19T10:00:00Z"
  }
}
```

## 5. OA后台管理服务 (AdminService)

### 5.1 管理员登录
```http
POST /api/v1/admin/login
Content-Type: application/json
```

**请求体**:
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**响应 (200)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "admin_info": {
      "admin_id": "admin_001",
      "username": "admin",
      "real_name": "管理员",
      "role": "admin",
      "created_at": "2024-12-19T10:00:00Z"
    }
  }
}
```

### 5.2 获取待审批申请列表
```http
GET /api/v1/admin/loan-applications
GET /api/v1/admin/loan-applications?status=SUBMITTED&page=1&limit=10
Authorization: Bearer <admin_token>
```

**响应 (200)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": [
    {
      "application_id": "app_1703001234567",
      "user_id": "user_1703001234567",
      "user_name": "张三",
      "product_name": "农业生产贷款",
      "amount": 100000,
      "term_months": 12,
      "status": "SUBMITTED",
      "submitted_at": "2024-12-19T10:00:00Z",
      "ai_risk_score": 75,
      "ai_suggestion": "建议批准",
      "requires_review": true
    }
  ],
  "total": 1
}
```

### 5.3 获取申请详情（管理员视角）
```http
GET /api/v1/admin/loan-applications/{application_id}
Authorization: Bearer <admin_token>
```

### 5.4 审批申请
```http
POST /api/v1/admin/loan-applications/{application_id}/approve
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**请求体**:
```json
{
  "action": "approve",
  "approved_amount": 80000,
  "comments": "批准申请，金额调整为8万元"
}
```

**响应 (200)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "application_id": "app_1703001234567",
    "action": "approve",
    "processed_at": "2024-12-19T10:00:00Z"
  }
}
```

## 6. AI智能体服务 (AIAgentService) v4.0

### 6.1 获取申请信息（统一接口）
```http
GET /api/v1/ai-agent/applications/{application_id}/info
X-AI-Agent-Token: <ai_token>
```

**描述**: 统一获取申请信息，自动识别贷款申请或农机租赁申请

**响应 (200)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "application_type": "LOAN_APPLICATION",
    "application_id": "app_1703001234567",
    "user_id": "user_1703001234567",
    "status": "SUBMITTED",
    "basic_info": {
      "amount": 100000,
      "term_months": 12,
      "purpose": "农业生产资金"
    },
    "financial_info": {
      "real_name": "张三",
      "id_card_number": "123456789012345678",
      "address": "某省某市某县某村"
    },
    "product_info": {
      "product_id": "product_001",
      "product_name": "农业生产贷款",
      "category": "agriculture"
    },
    "documents": [
      {
        "doc_type": "id_card",
        "file_id": "file_1703001234567"
      }
    ],
    "submitted_at": "2024-12-19T10:00:00Z"
  }
}
```

### 6.2 提交AI决策结果（统一接口）
```http
POST /api/v1/ai-agent/applications/{application_id}/decisions?application_type=LOAN_APPLICATION&decision=AUTO_APPROVED&risk_score=0.2&risk_level=LOW&confidence_score=0.95&analysis_summary=AI分析建议批准&approved_amount=80000&approved_term_months=12
X-AI-Agent-Token: <ai_token>
```

**查询参数**:
- `application_type`: 申请类型 (必需) - `LOAN_APPLICATION`, `MACHINERY_LEASING`
- `decision`: AI决策结果 (必需)
- `risk_score`: 风险分数 (必需) - 0-1之间
- `risk_level`: 风险等级 (必需) - `LOW`, `MEDIUM`, `HIGH`
- `confidence_score`: 置信度 (必需) - 0-1之间
- `analysis_summary`: 分析摘要 (必需)
- `approved_amount`: 批准金额 (可选，贷款申请专用)
- `approved_term_months`: 批准期限 (可选，贷款申请专用)
- `suggested_interest_rate`: 建议利率 (可选，贷款申请专用)
- `suggested_deposit`: 建议押金 (可选，农机租赁专用)
- `detailed_analysis`: 详细分析JSON字符串 (可选)
- `recommendations`: 建议列表，逗号分隔 (可选)
- `conditions`: 条件列表，逗号分隔 (可选)
- `ai_model_version`: AI模型版本 (可选)
- `workflow_id`: 工作流ID (可选)

**响应 (200)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "success": true,
    "ai_operation_id": "ai_op_1703001234567",
    "processed_at": "2024-12-19T10:00:00Z",
    "next_step_message": "AI决策已成功处理，系统将继续后续流程"
  }
}
```

### 6.3 获取外部数据（多类型支持）
```http
GET /api/v1/ai-agent/external-data/{user_id}?data_types=credit_report,bank_flow,blacklist_check
X-AI-Agent-Token: <ai_token>
```

**查询参数**:
- `data_types`: 数据类型，逗号分隔 (必需)
  - `credit_report`: 征信报告
  - `bank_flow`: 银行流水
  - `blacklist_check`: 黑名单检查
  - `government_subsidy`: 政府补贴
  - `farming_qualification`: 农业资质

**响应 (200)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "user_id": "user_1703001234567",
    "data_types": ["credit_report", "bank_flow", "blacklist_check"],
    "credit_data": {
      "credit_score": 750,
      "credit_grade": "A",
      "overdue_count": 0
    },
    "bank_data": {
      "monthly_income": 15000,
      "monthly_expense": 8000,
      "balance": 50000
    },
    "blacklist_data": {
      "is_blacklisted": false,
      "risk_level": "LOW"
    },
    "retrieved_at": "2024-12-19T10:00:00Z"
  }
}
```

### 6.4 获取AI模型配置（多类型支持）
```http
GET /api/v1/ai-agent/config/models
X-AI-Agent-Token: <ai_token>
```

**响应 (200)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "models": [
      {
        "model_name": "claude-3.5-sonnet",
        "version": "v1.0",
        "type": "LOAN_APPLICATION",
        "enabled": true
      },
      {
        "model_name": "gpt-4o",
        "version": "v1.0",
        "type": "MACHINERY_LEASING",
        "enabled": true
      }
    ],
    "risk_thresholds": {
      "low_risk_threshold": 0.3,
      "medium_risk_threshold": 0.7,
      "auto_approve_threshold": 0.2,
      "auto_reject_threshold": 0.8
    },
    "business_rules": {
      "max_loan_amount": 1000000,
      "min_credit_score": 600,
      "max_leasing_term_days": 365,
      "required_documents": ["id_card", "income_proof"]
    },
    "updated_at": "2024-12-19T10:00:00Z"
  }
}
```

### 6.5 获取农机租赁申请信息（专用接口）
```http
GET /api/v1/ai-agent/machinery-leasing/applications/{application_id}
X-AI-Agent-Token: <ai_token>
```

**响应 (200)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "application_id": "ml_1703001234567",
    "lessee_info": {
      "name": "张三",
      "id_card": "123456789012345678",
      "phone": "13800138000"
    },
    "lessor_info": {
      "company_name": "农机租赁公司",
      "contact": "李四"
    },
    "machinery_info": {
      "type": "拖拉机",
      "model": "东方红LX1000",
      "value": 150000
    },
    "leasing_details": {
      "rental_amount": 5000,
      "term_days": 30,
      "deposit": 10000
    },
    "status": "SUBMITTED",
    "submitted_at": "2024-12-19T10:00:00Z"
  }
}
```

### 6.6 获取AI操作日志
```http
GET /api/v1/ai-agent/logs
GET /api/v1/ai-agent/logs?application_id=app_123&application_type=LOAN_APPLICATION&page=1&limit=20
X-AI-Agent-Token: <ai_token>
```

**查询参数**:
- `application_id`: 申请ID (可选)
- `application_type`: 申请类型 (可选) - `LOAN_APPLICATION`, `MACHINERY_LEASING`
- `page`: 页码 (可选，默认1)
- `limit`: 每页数量 (可选，默认20)

**响应 (200)**:
```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "logs": [
      {
        "operation_id": "ai_op_1703001234567",
        "application_id": "app_1703001234567",
        "application_type": "LOAN_APPLICATION",
        "operation": "AI_DECISION",
        "result": "AUTO_APPROVED",
        "details": "AI分析建议批准",
        "timestamp": "2024-12-19T10:00:00Z"
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

## 7. 错误响应格式

### 7.1 标准错误响应
```json
{
  "code": 1001,
  "message": "请求参数错误",
  "error_details": "username字段不能为空"
}
```

### 7.2 错误码定义
- `0`: 成功
- `1001`: 请求参数错误
- `2001`: 未授权
- `3001`: 禁止访问
- `4001`: 资源不存在
- `4002`: AI申请信息不存在
- `4003`: AI申请状态冲突
- `4004`: AI分析参数错误
- `5001`: 服务器内部错误
- `5002`: AI服务不可用

### 7.3 HTTP状态码映射
- `200`: 成功
- `201`: 创建成功
- `400`: 请求参数错误 (1001-1999)
- `401`: 未授权 (2001-2999)
- `403`: 禁止访问 (3001-3999)
- `404`: 资源不存在 (4001-4999)
- `409`: 冲突 (特殊状态冲突)
- `500`: 服务器内部错误 (5001-5999)

## 8. 认证和安全

### 8.1 用户认证
```http
Authorization: Bearer <jwt_token>
```

### 8.2 管理员认证
```http
Authorization: Bearer <admin_jwt_token>
```

### 8.3 AI智能体认证
```http
X-AI-Agent-Token: <ai_agent_token>
```

## 9. 接口测试

使用提供的测试脚本进行全面的接口测试：

```bash
chmod +x doc/agent/backend/Test-API.sh
./doc/agent/backend/Test-API.sh
```

测试脚本包含：
- ✅ 健康检查测试
- ✅ 用户服务完整流程测试
- ✅ 贷款服务完整流程测试
- ✅ 文件服务接口验证
- ✅ 管理员服务测试
- ✅ AI智能体服务测试
- ✅ 错误场景测试
- ✅ 性能和并发测试

## 10. 更新日志

### v4.1 (2024-12-19)
- ✅ 完善所有处理器实现
- ✅ 统一错误处理和响应格式
- ✅ 增强AI智能体服务功能
- ✅ 完善文件服务实现
- ✅ 优化路由配置
- ✅ 增加完整接口测试脚本
- ✅ 改进API文档结构

### v4.0
- 初始化统一多类型AI智能体接口系统
- 支持贷款申请和农机租赁申请
- 自动申请类型识别
- AI风险评估和决策支持

---

**注意**: 此文档描述的是完整的API规范。实际部署时请确保：
1. 数据库正确初始化
2. JWT密钥和AI Token正确配置
3. 文件上传路径可写
4. 日志目录可写
5. 相关中间件正确配置 