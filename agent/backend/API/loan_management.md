# 贷款管理模块 - API 接口文档

## 📋 模块概述

贷款管理模块是数字惠农系统的核心金融服务模块，为农户提供便捷的贷款申请、审批、放款、还款等全生命周期服务。集成AI审批系统，实现快速、智能的风险评估和决策。

### 核心功能
- 💰 **贷款申请**: 多种贷款产品申请，智能风险评估
- 🤖 **AI审批**: 基于Dify平台的智能审批工作流
- 📊 **额度评估**: 动态信用额度评估和管理
- 💳 **还款管理**: 灵活的还款方式和提醒服务
- 📈 **数据分析**: 贷款数据统计和风险监控

---

## 🏦 贷款产品管理

### 1.1 获取贷款产品列表
```http
GET /api/loans/products?user_type=farmer&amount_min=10000&amount_max=500000
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": [
        {
            "id": 1001,
            "product_code": "NYCD001",
            "product_name": "农业创业贷",
            "description": "专为新型农业经营主体提供的创业资金支持",
            "min_amount": 10000,
            "max_amount": 500000,
            "min_term": 6,
            "max_term": 36,
            "interest_rate": 0.065,
            "interest_type": "fixed",
            "user_types": ["farmer", "farm_owner"],
            "collateral_required": false,
            "guarantor_required": true,
            "features": ["免抵押", "快速审批", "随借随还"],
            "eligibility": {
                "min_age": 18,
                "max_age": 65,
                "min_credit_score": 600,
                "required_auth": ["real_name", "bank_card"]
            },
            "status": "active"
        }
    ]
}
```

### 1.2 获取产品详情
```http
GET /api/loans/products/{product_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "id": 1001,
        "product_code": "NYCD001",
        "product_name": "农业创业贷",
        "description": "专为新型农业经营主体提供的创业资金支持",
        "detailed_description": "该产品面向有稳定收入来源的农业从业者...",
        "min_amount": 10000,
        "max_amount": 500000,
        "min_term": 6,
        "max_term": 36,
        "interest_rate": 0.065,
        "repayment_methods": ["equal_installment", "interest_first"],
        "application_materials": [
            {
                "name": "身份证明",
                "description": "二代身份证正反面",
                "required": true
            },
            {
                "name": "收入证明",
                "description": "近3个月银行流水",
                "required": true
            }
        ],
        "approval_process": {
            "estimated_time": "24小时",
            "steps": ["资料审核", "征信查询", "AI风险评估", "人工复审", "放款"]
        }
    }
}
```

---

## 💰 贷款申请管理

### 2.1 提交贷款申请
```http
POST /api/loans/applications
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "product_id": 1001,
    "amount": 100000,
    "term": 12,
    "purpose": "购买农机设备",
    "repayment_method": "equal_installment",
    "business_info": {
        "business_type": "种植业",
        "annual_income": 200000,
        "planting_area": 50,
        "main_crops": ["水稻", "玉米"],
        "years_experience": 5
    },
    "guarantor_info": {
        "name": "李四",
        "id_card": "370123199001011235",
        "phone": "13800138001",
        "relationship": "亲属",
        "annual_income": 150000
    },
    "materials": [
        {
            "type": "id_card",
            "file_url": "https://example.com/files/id_card_front.jpg"
        },
        {
            "type": "income_proof",
            "file_url": "https://example.com/files/bank_statement.pdf"
        }
    ]
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "申请提交成功",
    "data": {
        "application_id": "LA20240115001",
        "status": "submitted",
        "estimated_approval_time": "24小时内",
        "next_steps": [
            "系统将自动进行初步审核",
            "请保持手机畅通，等待审核结果"
        ]
    }
}
```

### 2.2 获取贷款申请列表
```http
GET /api/loans/applications?status=all&page=1&limit=10
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total": 25,
        "page": 1,
        "limit": 10,
        "applications": [
            {
                "id": "LA20240115001",
                "product_name": "农业创业贷",
                "amount": 100000,
                "amount_yuan": "100,000.00",
                "term": 12,
                "status": "approved",
                "status_text": "已批准",
                "applied_at": "2024-01-15T10:30:00Z",
                "approved_at": "2024-01-16T14:20:00Z",
                "interest_rate": 0.065,
                "monthly_payment": 8849
            }
        ]
    }
}
```

### 2.3 获取申请详情
```http
GET /api/loans/applications/{application_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "id": "LA20240115001",
        "product": {
            "id": 1001,
            "name": "农业创业贷",
            "code": "NYCD001"
        },
        "amount": 100000,
        "term": 12,
        "purpose": "购买农机设备",
        "status": "approved",
        "status_text": "已批准",
        "interest_rate": 0.065,
        "repayment_method": "equal_installment",
        "monthly_payment": 8849,
        "applied_at": "2024-01-15T10:30:00Z",
        "approved_at": "2024-01-16T14:20:00Z",
        "approval_history": [
            {
                "stage": "初审",
                "status": "passed",
                "reviewer": "系统自动审核",
                "review_time": "2024-01-15T11:00:00Z",
                "comments": "基本资料齐全，信用良好"
            },
            {
                "stage": "终审",
                "status": "passed",
                "reviewer": "风控部门",
                "review_time": "2024-01-16T14:20:00Z",
                "comments": "综合评估通过，同意放款"
            }
        ]
    }
}
```

### 2.4 更新申请信息
```http
PUT /api/loans/applications/{application_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "business_info": {
        "annual_income": 250000,
        "planting_area": 60
    },
    "additional_materials": [
        {
            "type": "land_certificate",
            "file_url": "https://example.com/files/land_cert.pdf"
        }
    ]
}
```

### 2.5 取消申请
```http
DELETE /api/loans/applications/{application_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "reason": "暂不需要贷款"
}
```

---

## 💳 贷款合同管理

### 3.1 生成贷款合同
```http
POST /api/loans/contracts/{application_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "合同生成成功",
    "data": {
        "contract_id": "LC20240116001",
        "contract_url": "https://example.com/contracts/LC20240116001.pdf",
        "expires_at": "2024-01-23T23:59:59Z",
        "signing_required": true
    }
}
```

### 3.2 签署合同
```http
POST /api/loans/contracts/{contract_id}/sign
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "signature": "base64_encoded_signature",
    "signing_location": "山东省济南市",
    "ip_address": "192.168.1.100"
}
```

### 3.3 获取合同详情
```http
GET /api/loans/contracts/{contract_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 📊 额度评估管理

### 4.1 获取信用额度
```http
GET /api/loans/credit-limit
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total_limit": 500000,
        "available_limit": 400000,
        "used_limit": 100000,
        "credit_score": 750,
        "credit_level": "优秀",
        "limit_details": [
            {
                "product_code": "NYCD001",
                "product_name": "农业创业贷",
                "limit": 300000,
                "available": 200000,
                "interest_rate": 0.065
            }
        ],
        "factors": {
            "income_stability": 0.85,
            "credit_history": 0.90,
            "business_scale": 0.75,
            "collateral_value": 0.70
        },
        "suggestions": [
            "完善征信记录可提升额度",
            "增加抵押物可获得更优利率"
        ]
    }
}
```

### 4.2 申请额度提升
```http
POST /api/loans/credit-limit/increase
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "requested_limit": 800000,
    "reason": "业务扩展需要",
    "additional_materials": [
        {
            "type": "business_expansion_plan",
            "file_url": "https://example.com/files/expansion_plan.pdf"
        }
    ]
}
```

---

## 💸 还款管理

### 5.1 获取还款计划
```http
GET /api/loans/{loan_id}/repayment-schedule
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "loan_id": "LN20240116001",
        "total_amount": 100000,
        "remaining_amount": 75000,
        "total_installments": 12,
        "remaining_installments": 9,
        "next_payment": {
            "installment_no": 4,
            "due_date": "2024-02-15",
            "principal": 8333,
            "interest": 516,
            "total_amount": 8849,
            "status": "pending"
        },
        "schedule": [
            {
                "installment_no": 1,
                "due_date": "2023-11-15",
                "principal": 8333,
                "interest": 516,
                "total_amount": 8849,
                "paid_amount": 8849,
                "paid_date": "2023-11-14",
                "status": "paid"
            }
        ]
    }
}
```

### 5.2 主动还款
```http
POST /api/loans/{loan_id}/repayment
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "payment_method": "bank_card",
    "amount": 8849,
    "installment_nos": [4],
    "bank_card_id": "BC20240101001"
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "还款成功",
    "data": {
        "payment_id": "PAY20240115001",
        "paid_amount": 8849,
        "remaining_principal": 66667,
        "next_due_date": "2024-03-15",
        "transaction_time": "2024-01-15T15:30:00Z"
    }
}
```

### 5.3 提前还款计算
```http
POST /api/loans/{loan_id}/prepayment-calculate
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "prepayment_type": "partial",
    "amount": 30000
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "计算成功",
    "data": {
        "prepayment_amount": 30000,
        "interest_saved": 2500,
        "remaining_principal": 45000,
        "new_monthly_payment": 5283,
        "new_final_payment_date": "2024-08-15",
        "prepayment_fee": 300,
        "total_payment_required": 30300
    }
}
```

### 5.4 获取还款记录
```http
GET /api/loans/{loan_id}/payments?page=1&limit=20
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 🏦 放款管理

### 6.1 确认放款信息
```http
GET /api/loans/{loan_id}/disbursement-info
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "loan_id": "LN20240116001",
        "approved_amount": 100000,
        "disbursement_amount": 99500,
        "service_fee": 500,
        "bank_account": {
            "account_number": "6226090000000001",
            "bank_name": "中国工商银行",
            "account_holder": "张三"
        },
        "estimated_arrival_time": "2-24小时",
        "status": "ready_to_disburse"
    }
}
```

### 6.2 申请放款
```http
POST /api/loans/{loan_id}/disburse
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "confirm_amount": 99500,
    "bank_account_id": "BA20240101001"
}
```

---

## 📈 贷款统计查询

### 7.1 获取贷款概览
```http
GET /api/loans/overview
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "active_loans": 2,
        "total_borrowed": 200000,
        "total_repaid": 45000,
        "outstanding_balance": 155000,
        "credit_limit_used": 0.31,
        "next_payment": {
            "amount": 8849,
            "due_date": "2024-02-15",
            "days_until_due": 10
        },
        "payment_status": "current"
    }
}
```

### 7.2 获取历史贷款
```http
GET /api/loans/history?status=completed&page=1&limit=10
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 🔔 提醒服务

### 8.1 获取还款提醒设置
```http
GET /api/loans/reminder-settings
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "sms_reminder": true,
        "email_reminder": false,
        "push_notification": true,
        "reminder_days": [7, 3, 1],
        "reminder_time": "09:00"
    }
}
```

### 8.2 更新提醒设置
```http
PUT /api/loans/reminder-settings
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "sms_reminder": true,
    "push_notification": true,
    "reminder_days": [5, 2, 1],
    "reminder_time": "10:00"
}
```

---

## 🔧 错误码说明

| 错误码 | 说明 | 处理建议 |
|-------|------|---------|
| 2001 | 贷款产品不存在 | 选择有效的贷款产品 |
| 2002 | 申请金额超出限制 | 调整申请金额 |
| 2003 | 用户不符合产品要求 | 完善认证信息 |
| 2004 | 信用评分不足 | 提升信用状况 |
| 2005 | 申请材料不完整 | 补充必要材料 |
| 2006 | 重复申请 | 等待当前申请处理 |
| 2007 | 合同已过期 | 重新生成合同 |
| 2008 | 银行卡信息错误 | 验证银行卡信息 |
| 2009 | 还款金额不足 | 检查账户余额 |
| 2010 | 贷款已结清 | 无需再次还款 |

---

## 📝 接口调用示例

### JavaScript示例
```javascript
// 获取贷款产品
const getLoanProducts = async (token) => {
    const response = await fetch('/api/loans/products?user_type=farmer', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};

// 提交贷款申请
const submitLoanApplication = async (token, applicationData) => {
    const response = await fetch('/api/loans/applications', {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(applicationData)
    });
    return response.json();
};

// 获取还款计划
const getRepaymentSchedule = async (token, loanId) => {
    const response = await fetch(`/api/loans/${loanId}/repayment-schedule`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};
```

### 业务流程说明
1. **申请流程**: 选择产品 → 填写申请 → 上传材料 → 等待审批 → 签署合同 → 放款
2. **审批流程**: 资料初审 → AI风险评估 → 人工复审 → 额度审批 → 合同生成
3. **还款流程**: 查看计划 → 选择方式 → 确认还款 → 更新记录 → 生成凭证
4. **提前还款**: 计算金额 → 确认费用 → 执行还款 → 调整计划

### 注意事项
1. **数据安全**: 敏感信息加密传输和存储
2. **风险控制**: 实时监控异常操作和风险指标
3. **合规要求**: 遵守金融监管政策和法规
4. **用户体验**: 提供清晰的状态反馈和操作指引 