# 数字惠农系统 - API 接口文档

## 1. 文档概述

本文档描述了数字惠农系统的完整API接口，包括用户管理、贷款服务、农机租赁、内容管理、OA后台等所有模块的接口定义。

### 接口规范
- **协议**: HTTPS
- **数据格式**: JSON
- **字符编码**: UTF-8
- **认证方式**: JWT Bearer Token
- **版本控制**: URL路径版本控制(/api/v1/)

### 基础响应格式

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {}, 
  "timestamp": 1640995200000,
  "request_id": "req_123456789"
}
```

### 错误码定义

| 错误码 | 说明 | HTTP状态码 |
|--------|------|------------|
| 200 | 操作成功 | 200 |
| 400 | 请求参数错误 | 400 |
| 401 | 未认证或Token无效 | 401 |
| 403 | 权限不足 | 403 |
| 404 | 资源不存在 | 404 |
| 409 | 资源冲突 | 409 |
| 422 | 业务逻辑错误 | 422 |
| 500 | 服务器内部错误 | 500 |

## 2. 认证授权

### 2.1 用户注册

**接口地址**: `POST /api/v1/auth/register`

**请求参数**:
```json
{
  "phone": "13800138000",
  "password": "password123",
  "code": "123456",
  "user_type": "farmer"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "user_id": 12345,
    "uuid": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

### 2.2 用户登录

**接口地址**: `POST /api/v1/auth/login`

**请求参数**:
```json
{
  "phone": "13800138000",
  "password": "password123",
  "platform": "app",
  "device_id": "device_123",
  "device_type": "android"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400,
    "user_info": {
      "id": 12345,
      "uuid": "550e8400-e29b-41d4-a716-446655440000",
      "phone": "13800138000",
      "user_type": "farmer",
      "real_name": "张三",
      "avatar": "https://cdn.example.com/avatar.jpg"
    }
  }
}
```

### 2.3 刷新Token

**接口地址**: `POST /api/v1/auth/refresh`

**请求参数**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## 3. 用户管理

### 3.1 获取用户信息

**接口地址**: `GET /api/v1/user/profile`

**请求头**: `Authorization: Bearer {access_token}`

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "id": 12345,
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "username": "farmer123",
    "phone": "13800138000",
    "email": "farmer@example.com",
    "user_type": "farmer",
    "status": "active",
    "real_name": "张三",
    "id_card": "110101199001011234",
    "avatar": "https://cdn.example.com/avatar.jpg",
    "gender": "male",
    "birthday": "1990-01-01",
    "province": "北京市",
    "city": "北京市",
    "county": "朝阳区",
    "address": "三环路123号",
    "is_real_name_verified": true,
    "is_bank_card_verified": true,
    "is_credit_verified": false,
    "balance": 10000,
    "credit_score": 750,
    "credit_level": "A",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 3.2 更新用户信息

**接口地址**: `PUT /api/v1/user/profile`

**请求参数**:
```json
{
  "real_name": "张三",
  "email": "farmer@example.com",
  "gender": "male",
  "birthday": "1990-01-01",
  "province": "北京市",
  "city": "北京市",
  "county": "朝阳区",
  "address": "三环路123号"
}
```

### 3.3 实名认证

**接口地址**: `POST /api/v1/user/auth/realname`

**请求参数**:
```json
{
  "real_name": "张三",
  "id_card_number": "110101199001011234",
  "id_card_front_img": "https://cdn.example.com/id_front.jpg",
  "id_card_back_img": "https://cdn.example.com/id_back.jpg",
  "face_verify_img": "https://cdn.example.com/face.jpg"
}
```

### 3.4 获取认证状态

**接口地址**: `GET /api/v1/user/auth/status`

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "real_name_auth": {
      "status": "approved",
      "submitted_at": "2024-01-01T00:00:00Z",
      "reviewed_at": "2024-01-02T00:00:00Z"
    },
    "bank_card_auth": {
      "status": "pending",
      "submitted_at": "2024-01-01T00:00:00Z"
    },
    "credit_auth": {
      "status": "none"
    }
  }
}
```

## 4. 贷款服务

### 4.1 获取贷款产品列表

**接口地址**: `GET /api/v1/loans/products`

**查询参数**:
- `product_type`: 产品类型(可选)
- `user_type`: 用户类型(可选)

**响应示例**:
```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "product_code": "NZDK001",
      "product_name": "农资贷",
      "description": "专为农户采购农资提供的贷款产品",
      "product_type": "agricultural_material",
      "min_amount": 100000,
      "max_amount": 10000000,
      "min_term": 30,
      "max_term": 365,
      "interest_rate": 0.12,
      "interest_type": "fixed",
      "repayment_method": "equal_installment",
      "partner_name": "XX银行",
      "status": "active"
    }
  ]
}
```

### 4.2 获取产品详情

**接口地址**: `GET /api/v1/loans/products/{product_id}`

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "product_code": "NZDK001",
    "product_name": "农资贷",
    "description": "专为农户采购农资提供的贷款产品",
    "product_type": "agricultural_material",
    "min_amount": 100000,
    "max_amount": 10000000,
    "min_term": 30,
    "max_term": 365,
    "interest_rate": 0.12,
    "eligibility_criteria": {
      "min_age": 18,
      "max_age": 65,
      "required_credit_score": 600,
      "required_documents": ["身份证", "银行流水", "土地承包合同"]
    },
    "required_documents": ["身份证正反面", "银行卡", "收入证明"],
    "applicable_user_types": ["farmer", "farm_owner"]
  }
}
```

### 4.3 提交贷款申请

**接口地址**: `POST /api/v1/loans/applications`

**请求参数**:
```json
{
  "product_id": 1,
  "applied_amount": 500000,
  "applied_term": 180,
  "purpose": "购买农资",
  "applicant_info": {
    "annual_income": 100000,
    "land_area": 10,
    "crop_types": ["水稻", "小麦"]
  },
  "uploaded_documents": [
    {
      "type": "id_card_front",
      "url": "https://cdn.example.com/doc1.jpg"
    },
    {
      "type": "bank_statement",
      "url": "https://cdn.example.com/doc2.pdf"
    }
  ]
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "申请提交成功",
  "data": {
    "application_id": 12345,
    "application_no": "APP202401010001",
    "status": "pending_ai",
    "estimated_review_time": "24小时内"
  }
}
```

### 4.4 获取我的申请列表

**接口地址**: `GET /api/v1/loans/applications`

**查询参数**:
- `status`: 申请状态(可选)
- `page`: 页码(默认1)
- `limit`: 每页数量(默认20)

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "applications": [
      {
        "id": 12345,
        "application_no": "APP202401010001",
        "product_name": "农资贷",
        "applied_amount": 500000,
        "applied_term": 180,
        "status": "approved",
        "status_text": "已通过",
        "approved_amount": 450000,
        "approved_rate": 0.12,
        "submitted_at": "2024-01-01T00:00:00Z",
        "approved_at": "2024-01-02T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 20
  }
}
```

### 4.5 获取申请详情

**接口地址**: `GET /api/v1/loans/applications/{application_id}`

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "id": 12345,
    "application_no": "APP202401010001",
    "product": {
      "id": 1,
      "product_name": "农资贷",
      "interest_rate": 0.12
    },
    "applied_amount": 500000,
    "applied_term": 180,
    "purpose": "购买农资",
    "status": "approved",
    "status_text": "已通过",
    "ai_score": 85.5,
    "ai_risk_level": "medium",
    "ai_comments": "综合评分良好，建议通过",
    "approved_amount": 450000,
    "approved_term": 180,
    "approved_rate": 0.12,
    "submitted_at": "2024-01-01T00:00:00Z",
    "ai_processed_at": "2024-01-01T12:00:00Z",
    "approved_at": "2024-01-02T00:00:00Z"
  }
}
```

## 5. 农机租赁

### 5.1 搜索附近农机

**接口地址**: `GET /api/v1/machines/nearby`

**查询参数**:
- `longitude`: 经度
- `latitude`: 纬度
- `radius`: 搜索半径(公里，默认10)
- `machine_type`: 设备类型(可选)
- `page`: 页码(默认1)
- `limit`: 每页数量(默认20)

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "machines": [
      {
        "id": 12345,
        "machine_code": "MC20240001",
        "machine_name": "东方红拖拉机",
        "brand": "东方红",
        "model": "LX904",
        "machine_type": "tractor",
        "status": "available",
        "images": ["https://cdn.example.com/machine1.jpg"],
        "hourly_rate": 8000,
        "daily_rate": 60000,
        "deposit_amount": 500000,
        "province": "北京市",
        "city": "北京市",
        "county": "朝阳区",
        "distance": 3.5,
        "average_rating": 4.8,
        "review_count": 25,
        "owner": {
          "id": 67890,
          "real_name": "李四",
          "avatar": "https://cdn.example.com/avatar2.jpg"
        }
      }
    ],
    "total": 15,
    "page": 1,
    "limit": 20
  }
}
```

### 5.2 获取农机详情

**接口地址**: `GET /api/v1/machines/{machine_id}`

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "id": 12345,
    "machine_code": "MC20240001",
    "machine_name": "东方红拖拉机",
    "brand": "东方红",
    "model": "LX904",
    "machine_type": "tractor",
    "status": "available",
    "images": ["https://cdn.example.com/machine1.jpg"],
    "specifications": {
      "power": "90马力",
      "weight": "3500kg",
      "length": "4.2m",
      "width": "2.1m"
    },
    "description": "性能优良的中型拖拉机，适合中等规模农田作业",
    "hourly_rate": 8000,
    "daily_rate": 60000,
    "deposit_amount": 500000,
    "min_rental_hours": 2,
    "max_rental_days": 30,
    "service_radius": 15,
    "manufacture_year": 2020,
    "working_hours": 1500,
    "average_rating": 4.8,
    "review_count": 25,
    "total_orders": 150,
    "success_orders": 148,
    "is_verified": true,
    "owner": {
      "id": 67890,
      "real_name": "李四",
      "avatar": "https://cdn.example.com/avatar2.jpg",
      "phone": "138****8000"
    },
    "location": {
      "province": "北京市",
      "city": "北京市",
      "county": "朝阳区",
      "detail_address": "XX路XX号",
      "longitude": 116.4074,
      "latitude": 39.9042
    }
  }
}
```

### 5.3 创建租赁订单

**接口地址**: `POST /api/v1/machines/{machine_id}/orders`

**请求参数**:
```json
{
  "start_time": "2024-01-10T08:00:00Z",
  "end_time": "2024-01-10T18:00:00Z",
  "rental_type": "hourly",
  "use_address": "北京市朝阳区XX农田",
  "use_longitude": 116.4074,
  "use_latitude": 39.9042,
  "renter_notes": "需要提前30分钟到达"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "订单创建成功",
  "data": {
    "order_id": 54321,
    "order_no": "ORD202401100001",
    "total_amount": 580000,
    "rental_amount": 80000,
    "deposit_amount": 500000,
    "status": "pending",
    "payment_deadline": "2024-01-09T18:00:00Z"
  }
}
```

### 5.4 获取我的订单

**接口地址**: `GET /api/v1/orders`

**查询参数**:
- `role`: 角色类型(renter/owner)
- `status`: 订单状态(可选)
- `page`: 页码(默认1)
- `limit`: 每页数量(默认20)

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "orders": [
      {
        "id": 54321,
        "order_no": "ORD202401100001",
        "machine": {
          "id": 12345,
          "machine_name": "东方红拖拉机",
          "machine_type": "tractor",
          "images": ["https://cdn.example.com/machine1.jpg"]
        },
        "start_time": "2024-01-10T08:00:00Z",
        "end_time": "2024-01-10T18:00:00Z",
        "rental_type": "hourly",
        "rental_amount": 80000,
        "deposit_amount": 500000,
        "total_amount": 580000,
        "status": "completed",
        "status_text": "已完成",
        "created_at": "2024-01-09T10:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 20
  }
}
```

### 5.5 订单支付

**接口地址**: `POST /api/v1/orders/{order_id}/pay`

**请求参数**:
```json
{
  "payment_method": "alipay",
  "payment_type": "deposit"
}
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "payment_no": "PAY202401100001",
    "payment_url": "https://pay.example.com/pay?order=xxx",
    "qr_code": "data:image/png;base64,xxx"
  }
}
```

### 5.6 提交评价

**接口地址**: `POST /api/v1/orders/{order_id}/reviews`

**请求参数**:
```json
{
  "overall_rating": 5,
  "device_rating": 5,
  "service_rating": 4,
  "delivery_rating": 5,
  "content": "设备性能很好，老板服务态度也不错",
  "images": ["https://cdn.example.com/review1.jpg"]
}
```

## 6. 内容管理

### 6.1 获取文章列表

**接口地址**: `GET /api/v1/articles`

**查询参数**:
- `category`: 分类(可选)
- `keyword`: 搜索关键词(可选)
- `page`: 页码(默认1)
- `limit`: 每页数量(默认20)

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "articles": [
      {
        "id": 1,
        "title": "春季农作物种植技术要点",
        "summary": "详细介绍春季主要农作物的种植技术和注意事项",
        "cover_image": "https://cdn.example.com/article1.jpg",
        "category": "种植技术",
        "author": "专家张三",
        "view_count": 1250,
        "like_count": 45,
        "is_top": false,
        "is_featured": true,
        "published_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "limit": 20
  }
}
```

### 6.2 获取文章详情

**接口地址**: `GET /api/v1/articles/{article_id}`

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "title": "春季农作物种植技术要点",
    "content": "文章详细内容...",
    "summary": "详细介绍春季主要农作物的种植技术和注意事项",
    "cover_image": "https://cdn.example.com/article1.jpg",
    "category": "种植技术",
    "author": "专家张三",
    "view_count": 1251,
    "like_count": 45,
    "share_count": 12,
    "is_top": false,
    "is_featured": true,
    "published_at": "2024-01-01T00:00:00Z"
  }
}
```

### 6.3 获取专家列表

**接口地址**: `GET /api/v1/experts`

**查询参数**:
- `specialty`: 专业领域(可选)
- `service_area`: 服务区域(可选)

**响应示例**:
```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "expert_name": "张教授",
      "avatar": "https://cdn.example.com/expert1.jpg",
      "title": "农业技术专家",
      "specialty": "水稻种植",
      "service_area": "华北地区",
      "experience_years": 20,
      "description": "从事水稻种植研究20年，发表论文50余篇",
      "consultation_count": 500,
      "average_rating": 4.9,
      "is_verified": true,
      "is_online": true
    }
  ]
}
```

### 6.4 提交咨询

**接口地址**: `POST /api/v1/consultations`

**请求参数**:
```json
{
  "expert_id": 1,
  "question": "水稻种植过程中出现叶片发黄现象，请问是什么原因？",
  "images": ["https://cdn.example.com/question1.jpg"],
  "contact_info": "微信：farmer123"
}
```

## 7. 系统功能

### 7.1 文件上传

**接口地址**: `POST /api/v1/files/upload`

**请求参数**: `multipart/form-data`
- `file`: 文件
- `business_type`: 业务类型(avatar/document/image)
- `business_id`: 业务ID(可选)

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "file_id": 12345,
    "file_url": "https://cdn.example.com/uploads/xxx.jpg",
    "file_name": "photo.jpg",
    "file_size": 102400,
    "file_type": "image/jpeg"
  }
}
```

### 7.2 获取系统配置

**接口地址**: `GET /api/v1/system/configs`

**查询参数**:
- `group`: 配置组(可选)

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "app_version": "1.0.0",
    "customer_service_phone": "400-123-4567",
    "customer_service_hours": "9:00-18:00",
    "about_us": "关于我们的介绍信息",
    "privacy_policy_url": "https://example.com/privacy",
    "terms_of_service_url": "https://example.com/terms"
  }
}
```

## 8. OA后台管理

### 8.1 OA用户登录

**接口地址**: `POST /api/v1/admin/auth/login`

**请求参数**:
```json
{
  "username": "admin",
  "password": "password123",
  "captcha": "ABCD",
  "captcha_id": "cap_123456"
}
```

### 8.2 用户管理

**接口地址**: `GET /api/v1/admin/users`

**查询参数**:
- `user_type`: 用户类型(可选)
- `status`: 用户状态(可选)
- `keyword`: 搜索关键词(可选)
- `page`: 页码(默认1)
- `limit`: 每页数量(默认20)

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "users": [
      {
        "id": 12345,
        "uuid": "550e8400-e29b-41d4-a716-446655440000",
        "phone": "13800138000",
        "user_type": "farmer",
        "status": "active",
        "real_name": "张三",
        "is_real_name_verified": true,
        "is_bank_card_verified": true,
        "created_at": "2024-01-01T00:00:00Z",
        "last_login_time": "2024-01-10T08:00:00Z"
      }
    ],
    "total": 1000,
    "page": 1,
    "limit": 20
  }
}
```

### 8.3 贷款申请管理

**接口地址**: `GET /api/v1/admin/loans/applications`

**查询参数**:
- `status`: 申请状态(可选)
- `product_id`: 产品ID(可选)
- `start_date`: 开始日期(可选)
- `end_date`: 结束日期(可选)
- `page`: 页码(默认1)
- `limit`: 每页数量(默认20)

### 8.4 审批贷款申请

**接口地址**: `POST /api/v1/admin/loans/applications/{application_id}/approve`

**请求参数**:
```json
{
  "approved_amount": 450000,
  "approved_term": 180,
  "approved_rate": 0.12,
  "comments": "申请人资质良好，同意放款"
}
```

### 8.5 认证审核

**接口地址**: `POST /api/v1/admin/auths/{auth_id}/review`

**请求参数**:
```json
{
  "action": "approve",
  "review_note": "材料齐全，审核通过"
}
```

### 8.6 统计数据

**接口地址**: `GET /api/v1/admin/statistics`

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "overview": {
      "total_users": 10000,
      "new_users_today": 50,
      "total_applications": 5000,
      "pending_applications": 120,
      "total_orders": 8000,
      "active_orders": 200
    },
    "charts": {
      "user_growth": [
        {"date": "2024-01-01", "count": 100},
        {"date": "2024-01-02", "count": 150}
      ],
      "application_trends": [
        {"date": "2024-01-01", "applications": 50, "approvals": 40},
        {"date": "2024-01-02", "applications": 60, "approvals": 45}
      ]
    }
  }
}
```

## 9. 错误处理

### 9.1 常见错误响应

**参数错误**:
```json
{
  "code": 400,
  "message": "请求参数错误",
  "errors": {
    "phone": ["手机号格式不正确"],
    "amount": ["金额必须大于0"]
  }
}
```

**认证错误**:
```json
{
  "code": 401,
  "message": "Token已过期，请重新登录"
}
```

**权限错误**:
```json
{
  "code": 403,
  "message": "权限不足，无法访问该资源"
}
```

**业务逻辑错误**:
```json
{
  "code": 422,
  "message": "设备当前不可用，无法创建订单"
}
```

### 9.2 错误码对照表

| 业务错误码 | 说明 |
|-----------|------|
| 10001 | 用户不存在 |
| 10002 | 密码错误 |
| 10003 | 验证码错误 |
| 10004 | 手机号已注册 |
| 20001 | 产品不存在 |
| 20002 | 申请金额超出限制 |
| 20003 | 用户已有待审批申请 |
| 30001 | 设备不存在 |
| 30002 | 设备不可用 |
| 30003 | 时间段冲突 |
| 40001 | 支付失败 |
| 40002 | 余额不足 |

## 10. 接口限流

### 10.1 限流规则

| 接口类型 | 限制频率 | 时间窗口 |
|---------|---------|----------|
| 登录接口 | 5次/分钟 | 1分钟 |
| 发送验证码 | 1次/分钟 | 1分钟 |
| 文件上传 | 10次/分钟 | 1分钟 |
| 普通查询 | 100次/分钟 | 1分钟 |
| 提交申请 | 5次/小时 | 1小时 |

### 10.2 限流响应

```json
{
  "code": 429,
  "message": "请求过于频繁，请稍后再试",
  "retry_after": 60
}
```

## 11. 版本控制

### 11.1 版本策略

- **URL版本控制**: `/api/v1/`, `/api/v2/`
- **向后兼容**: 新版本保持向后兼容性
- **废弃通知**: 提前3个月通知接口废弃
- **版本生命周期**: 每个版本维护2年

### 11.2 版本信息

**当前版本**: v1.0
**支持版本**: v1.0
**计划版本**: v1.1 (2024年6月)

## 12. 接口安全

### 12.1 HTTPS要求

所有API接口必须使用HTTPS协议，确保数据传输安全。

### 12.2 签名验证

部分敏感接口需要额外的签名验证：

```
签名算法: HMAC-SHA256
签名字符串: HTTP_METHOD + URL_PATH + TIMESTAMP + REQUEST_BODY
```

### 12.3 IP白名单

OA后台接口支持IP白名单限制，只允许指定IP访问。

### 12.4 数据脱敏

API响应中的敏感信息将自动脱敏处理：
- 手机号: `138****8000`
- 身份证号: `110101********1234`
- 银行卡号: `6222********1234` 