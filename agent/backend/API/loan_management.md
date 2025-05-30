# 贷款管理模块 - API 接口文档

## 📋 模块概述

贷款管理模块为惠农APP/Web用户提供贷款产品查询、申请服务，并为OA系统管理员提供贷款审批和管理功能。

### 核心功能
-   **用户端 (`/api/user/loan/*`)**: 产品浏览、在线申请、申请状态跟踪。
-   **OA管理员端 (`/api/oa/admin/loans/*`)**: 申请列表、审批操作、风险评估、放款管理。

---

## 🏦 惠农APP/Web - 贷款产品接口

**接口路径前缀**: `/api/user/loan/products`
**认证要求**: `RequireAuth` (惠农APP/Web用户)
**适用平台**: `app`, `web`

### 1.1 获取贷款产品列表

```http
GET /api/user/loan/products?user_type=farmer&amount_min=10000&amount_max=500000
Authorization: Bearer {access_token}
```

**响应示例 (部分):**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": [
        {
            "id": 1001,
            "product_name": "农业创业贷",
            // ... 其他产品字段
        }
    ]
}
```

### 1.2 获取贷款产品详情

```http
GET /api/user/loan/products/{product_id}
Authorization: Bearer {access_token}
```

---

## 💰 惠农APP/Web - 贷款申请接口

**接口路径前缀**: `/api/user/loan/applications`
**认证要求**: `RequireAuth` (惠农APP/Web用户)
**适用平台**: `app`, `web`

### 2.1 提交贷款申请

```http
POST /api/user/loan/applications
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "product_id": 1001,
    "amount": 100000,
    // ... 其他申请信息
}
```

### 2.2 获取当前用户的贷款申请列表

```http
GET /api/user/loan/applications?status=all&page=1&limit=10
Authorization: Bearer {access_token}
```

### 2.3 获取当前用户指定贷款申请详情

```http
GET /api/user/loan/applications/{application_id}
Authorization: Bearer {access_token}
```

### 2.4 用户取消贷款申请 (在特定状态下)

```http
DELETE /api/user/loan/applications/{application_id}
Authorization: Bearer {access_token}
```

---

## 🛠️ OA系统 - 贷款审批与管理接口 (管理员)

**接口路径前缀**: `/api/oa/admin/loans`
**认证要求**: `RequireAuth`, `CheckPlatform("oa")`, `RequireRole("admin")`
**适用平台**: `oa`

### 3.1 获取所有贷款申请列表 (管理员视图)

```http
GET /api/oa/admin/loans/applications?status=pending&user_id=101&page=1&limit=20
Authorization: Bearer {oa_access_token}
```

**Query Parameters (示例)**:
-   `status`: `pending`, `approved`, `rejected`, `submitted`, `under_review`, `all`
-   `user_id`: 筛选特定惠农用户的申请 (User ID)
-   `product_id`: 筛选特定贷款产品的申请
-   `date_range_start`, `date_range_end`: 按申请日期筛选

**响应示例 (部分):**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total": 50,
        "applications": [
            {
                "id": "LA20240115001",
                "user_info": { "user_id": 101, "real_name": "张三", "phone": "138..." },
                "product_name": "农业创业贷",
                "amount": 100000,
                "status": "pending",
                "applied_at": "2024-01-15T10:30:00Z"
            }
        ]
    }
}
```

### 3.2 获取指定贷款申请详情 (管理员视图)

```http
GET /api/oa/admin/loans/applications/{application_id}
Authorization: Bearer {oa_access_token}
```

(响应会包含更详细的审核信息、用户信息等)

### 3.3 审批操作：批准贷款申请

```http
POST /api/oa/admin/loans/applications/{application_id}/approve
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "approved_amount": 100000, // 批准金额，可能与申请金额不同
    "approved_term": 12,       // 批准期限
    "interest_rate": 0.06,     // 最终利率
    "repayment_start_date": "2024-02-01",
    "approval_comments": "综合评估通过，同意放款。"
}
```

### 3.4 审批操作：拒绝贷款申请

```http
POST /api/oa/admin/loans/applications/{application_id}/reject
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "rejection_reason_code": "INSUFFICIENT_CREDIT", // 预定义的拒绝原因代码
    "rejection_comments": "申请人信用评分不足，且缺乏有效抵押物。",
    "notify_user": true
}
```

### 3.5 审批操作：退回贷款申请 (要求补充材料)

```http
POST /api/oa/admin/loans/applications/{application_id}/return
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "return_reason": "需要补充最新的银行流水证明。",
    "required_materials": ["近三个月银行流水"],
    "notify_user": true
}
```

### 3.6 开始人工审核 (标记申请进入审核流程)

```http
POST /api/oa/admin/loans/applications/{application_id}/start-review
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "reviewer_id": 301, // 当前审核员的OA User ID
    "review_department": "信贷审批部"
}
```

### 3.7 重试AI风险评估 (当Dify调用失败或需要重新评估时)

```http
POST /api/oa/admin/loans/applications/{application_id}/retry-ai
Authorization: Bearer {oa_access_token}
```

### 3.8 获取贷款统计数据 (管理员)

```http
GET /api/oa/admin/loans/statistics
Authorization: Bearer {oa_access_token}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total_applications": 1200,
        "pending_applications": 50,
        "approved_applications": 800,
        "rejected_applications": 350,
        "total_loan_amount": 50000000,
        "average_loan_amount": 41666,
        "default_rate": 0.025
    }
}
```

---

**说明**: 上述API仅为核心示例，实际项目中会包含更多参数、状态流转和错误处理。前端开发时请结合后端具体实现进行调整。

## 💳 贷款合同管理

### 4.1 生成贷款合同
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

### 4.2 签署合同
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

### 4.3 获取合同详情
```http
GET /api/loans/contracts/{contract_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 📊 额度评估管理

### 5.1 获取信用额度
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

### 5.2 申请额度提升
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

### 6.1 获取还款计划
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

### 6.2 主动还款
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

### 6.3 提前还款计算
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

### 6.4 获取还款记录
```http
GET /api/loans/{loan_id}/payments?page=1&limit=20
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 🏦 放款管理

### 7.1 确认放款信息
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

### 7.2 申请放款
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

### 8.1 获取贷款概览
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

### 8.2 获取历史贷款
```http
GET /api/loans/history?status=completed&page=1&limit=10
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 🔔 提醒服务

### 9.1 获取还款提醒设置
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

### 9.2 更新提醒设置
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
    const response = await fetch('/api/user/loan/products?user_type=farmer', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};

// 提交贷款申请
const submitLoanApplication = async (token, applicationData) => {
    const response = await fetch('/api/user/loan/applications', {
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