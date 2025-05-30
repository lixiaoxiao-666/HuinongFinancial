# 用户管理模块 - API 接口文档

## 📋 模块概述

用户管理模块是数字惠农系统的核心基础模块，提供用户注册、认证、权限管理、信息维护等功能。支持多种用户类型（个体农户、家庭农场主、合作社、企业等）以及完整的用户生命周期管理。

### 核心功能
- 👤 **用户注册登录**: 手机号注册、密码登录、短信验证码登录
- 🔐 **身份认证**: 实名认证、银行卡认证、征信认证
- 📱 **信息管理**: 个人信息、地址信息、头像上传
- 🎭 **权限控制**: 基于角色的权限管理
- 📊 **用户画像**: 标签管理、行为分析

---

## 🔐 用户认证相关

### 1.1 用户注册
```http
POST /api/auth/register
Content-Type: application/json

{
    "phone": "13800138000",
    "password": "password123",
    "verification_code": "123456",
    "user_type": "farmer",
    "real_name": "张三",
    "province": "山东省",
    "city": "济南市",
    "county": "历城区",
    "device_id": "iPhone_12_ABC123",
    "platform": "app"
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "注册成功",
    "data": {
        "user": {
            "id": 1001,
            "uuid": "uuid-abc-123",
            "phone": "13800138000",
            "user_type": "farmer",
            "real_name": "张三",
            "status": "active",
            "is_real_name_verified": false
        },
        "session": {
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "expires_in": 86400
        }
    }
}
```

### 1.2 发送验证码
```http
POST /api/auth/send-sms
Content-Type: application/json

{
    "phone": "13800138000",
    "type": "register"
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "验证码已发送",
    "data": {
        "expires_in": 300,
        "retry_after": 60
    }
}
```

### 1.3 验证码登录
```http
POST /api/auth/login-sms
Content-Type: application/json

{
    "phone": "13800138000",
    "verification_code": "123456",
    "device_id": "iPhone_12_ABC123",
    "platform": "app",
    "device_type": "ios",
    "app_version": "1.0.0"
}
```

### 1.4 密码重置
```http
POST /api/auth/reset-password
Content-Type: application/json

{
    "phone": "13800138000",
    "verification_code": "123456",
    "new_password": "newpassword123"
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "密码重置成功"
}
```

---

## 👤 用户信息管理

### 2.1 获取用户信息
```http
GET /api/user/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "id": 1001,
        "uuid": "uuid-abc-123",
        "username": "zhangsan",
        "phone": "13800138000",
        "email": "zhangsan@example.com",
        "user_type": "farmer",
        "status": "active",
        "real_name": "张三",
        "id_card": "370123199001011234",
        "avatar": "https://example.com/avatar.jpg",
        "gender": "male",
        "birthday": "1990-01-01",
        "province": "山东省",
        "city": "济南市",
        "county": "历城区",
        "address": "某某村123号",
        "longitude": 117.1234,
        "latitude": 36.5678,
        "is_real_name_verified": true,
        "is_bank_card_verified": true,
        "is_credit_verified": false,
        "balance": 50000,
        "credit_score": 750,
        "credit_level": "良好",
        "last_login_time": "2024-01-15T14:25:30Z",
        "created_at": "2024-01-01T10:00:00Z"
    }
}
```

### 2.2 更新用户信息
```http
PUT /api/user/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "real_name": "张三丰",
    "email": "zhangsan@newmail.com",
    "gender": "male",
    "birthday": "1990-01-01",
    "address": "新地址123号"
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "更新成功",
    "data": {
        "updated_fields": ["real_name", "email", "address"]
    }
}
```

### 2.3 头像上传
```http
POST /api/user/avatar
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: multipart/form-data

avatar: [图片文件]
```

**响应示例:**
```json
{
    "code": 200,
    "message": "头像上传成功",
    "data": {
        "avatar_url": "https://example.com/uploads/avatars/uuid-abc-123.jpg"
    }
}
```

### 2.4 修改密码
```http
POST /api/user/change-password
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "old_password": "oldpassword123",
    "new_password": "newpassword123"
}
```

---

## 🔐 身份认证管理

### 3.1 实名认证申请
```http
POST /api/user/auth/real-name
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "id_card_number": "370123199001011234",
    "real_name": "张三",
    "id_card_front_img": "https://example.com/front.jpg",
    "id_card_back_img": "https://example.com/back.jpg",
    "face_verify_img": "https://example.com/face.jpg"
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "实名认证申请已提交",
    "data": {
        "auth_id": 10001,
        "auth_status": "pending",
        "estimated_review_time": "24小时内"
    }
}
```

### 3.2 银行卡认证申请
```http
POST /api/user/auth/bank-card
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "bank_card_number": "6226090000000001",
    "bank_name": "中国工商银行",
    "cardholder_name": "张三",
    "card_type": "储蓄卡"
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "银行卡认证申请已提交",
    "data": {
        "auth_id": 10002,
        "auth_status": "pending"
    }
}
```

### 3.3 获取认证状态
```http
GET /api/user/auth/status
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "real_name_auth": {
            "status": "approved",
            "auth_time": "2024-01-15T10:30:00Z",
            "review_note": "认证通过"
        },
        "bank_card_auth": {
            "status": "pending",
            "submitted_time": "2024-01-15T14:20:00Z"
        },
        "credit_auth": {
            "status": "not_submitted"
        }
    }
}
```

### 3.4 征信报告查询
```http
POST /api/user/auth/credit-query
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "provider": "pboc",
    "query_reason": "loan_application"
}
```

---

## 📊 用户标签管理

### 4.1 获取用户标签
```http
GET /api/user/tags
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "behavior_tags": [
            {
                "tag_key": "active_level",
                "tag_value": "high",
                "weight": 0.8
            }
        ],
        "preference_tags": [
            {
                "tag_key": "crop_type",
                "tag_value": "蔬菜种植",
                "weight": 0.9
            }
        ],
        "attribute_tags": [
            {
                "tag_key": "farm_size",
                "tag_value": "小规模",
                "weight": 1.0
            }
        ]
    }
}
```

### 4.2 更新用户标签
```http
POST /api/user/tags
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "tag_type": "preference",
    "tag_key": "crop_type",
    "tag_value": "水稻种植",
    "weight": 0.9
}
```

---

## 🏠 地址管理

### 5.1 获取地址列表
```http
GET /api/user/addresses
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": [
        {
            "id": 1,
            "name": "家庭农场",
            "contact_name": "张三",
            "contact_phone": "13800138000",
            "province": "山东省",
            "city": "济南市",
            "county": "历城区",
            "address": "某某村123号",
            "longitude": 117.1234,
            "latitude": 36.5678,
            "is_default": true,
            "address_type": "farm"
        }
    ]
}
```

### 5.2 添加地址
```http
POST /api/user/addresses
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "name": "新农场",
    "contact_name": "张三",
    "contact_phone": "13800138000",
    "province": "山东省",
    "city": "济南市",
    "county": "历城区",
    "address": "新地址456号",
    "longitude": 117.5678,
    "latitude": 36.9012,
    "address_type": "farm",
    "is_default": false
}
```

### 5.3 更新地址
```http
PUT /api/user/addresses/{address_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "name": "更新后的农场名称",
    "address": "更新后的地址"
}
```

### 5.4 删除地址
```http
DELETE /api/user/addresses/{address_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 📱 设备管理

### 6.1 获取设备列表
```http
GET /api/user/devices
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": [
        {
            "device_id": "iPhone_12_ABC123",
            "device_type": "ios",
            "device_name": "张三的iPhone",
            "platform": "app",
            "app_version": "1.0.0",
            "last_active_time": "2024-01-15T14:25:30Z",
            "status": "active"
        }
    ]
}
```

### 6.2 注销设备
```http
DELETE /api/user/devices/{device_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 💰 账户管理

### 7.1 获取账户余额
```http
GET /api/user/balance
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "balance": 50000,
        "balance_yuan": "500.00",
        "frozen_amount": 0,
        "available_amount": 50000
    }
}
```

### 7.2 获取账户流水
```http
GET /api/user/transactions?page=1&limit=20&type=all
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total": 50,
        "page": 1,
        "limit": 20,
        "transactions": [
            {
                "id": 10001,
                "type": "income",
                "amount": 10000,
                "amount_yuan": "100.00",
                "description": "账户充值",
                "balance_after": 60000,
                "created_at": "2024-01-15T10:30:00Z"
            }
        ]
    }
}
```

---

## 🔧 错误码说明

| 错误码 | 说明 | 处理建议 |
|-------|------|---------|
| 1001 | 手机号已注册 | 提示用户直接登录 |
| 1002 | 验证码错误 | 重新发送验证码 |
| 1003 | 验证码已过期 | 重新发送验证码 |
| 1004 | 密码错误 | 提示用户检查密码 |
| 1005 | 用户不存在 | 提示用户注册 |
| 1006 | 账户已被冻结 | 联系客服处理 |
| 1007 | 实名认证信息不匹配 | 重新提交认证 |
| 1008 | 银行卡号无效 | 检查银行卡号 |
| 1009 | 文件上传失败 | 重新上传文件 |
| 1010 | 地址信息不完整 | 补充完整地址 |

---

## 📝 接口调用示例

### JavaScript示例
```javascript
// 用户登录
const login = async (phone, password) => {
    const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            phone,
            password,
            platform: 'app',
            device_id: getDeviceId(),
            device_type: 'ios'
        })
    });
    return response.json();
};

// 获取用户信息
const getUserProfile = async (token) => {
    const response = await fetch('/api/user/profile', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};

// 更新用户信息
const updateProfile = async (token, data) => {
    const response = await fetch('/api/user/profile', {
        method: 'PUT',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    });
    return response.json();
};
```

### 安全注意事项
1. **Token管理**: 安全存储访问令牌，定期刷新
2. **密码安全**: 密码需要包含大小写字母、数字和特殊字符
3. **数据验证**: 所有用户输入都需要进行验证
4. **HTTPS**: 生产环境必须使用HTTPS协议
5. **限流保护**: 对敏感接口实施频率限制 