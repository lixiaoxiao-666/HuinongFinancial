# 用户管理模块 - API 接口文档

## 📋 模块概述

用户管理模块负责处理惠农APP/Web端用户以及OA系统用户的相关操作。包括用户注册、登录、信息管理、身份认证等。

### 平台与用户类型

-   **惠农APP/Web端 (`platform: "app"` 或 `"web"`)**: 主要面向C端用户（农户、农场主等）。
    -   API路径: `/api/auth/*` (认证), `/api/user/*` (用户操作)
    -   用户模型: `User`
    -   权限: 普通用户权限，无特殊角色区分。
-   **OA系统 (`platform: "oa"`)**: 主要面向内部运营和管理人员。
    -   API路径: `/api/oa/auth/*` (认证), `/api/oa/user/*` (普通OA用户操作), `/api/oa/admin/users/*` (管理员用户管理操作)
    -   用户模型: `OAUser` (包含 `RoleID`)
    -   权限: 分为普通OA用户和OA管理员。

---

## 🔑 惠农APP/Web - 认证接口

**适用平台**: `app`, `web`

### 1.1 用户注册

```http
POST /api/auth/register
Content-Type: application/json

{
    "phone": "13800138000",
    "password": "password123",
    "verification_code": "123456", // 可选，取决于系统配置
    "user_type": "farmer", // 用户类型: farmer, farm_owner, cooperative, enterprise
    "real_name": "张三",
    "platform": "app", // 客户端平台: app, web
    "device_info": {
        "device_id": "iPhone_12_ABC123",
        "device_type": "ios", // ios, android, web
        "device_name": "张三的iPhone",
        "app_version": "1.0.0"
    }
}
```

**响应示例 (成功):**
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
            "status": "active"
        },
        "session": {
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "expires_in": 86400 // access_token有效期（秒）
        }
    }
}
```

### 1.2 用户登录 (密码)

```http
POST /api/auth/login
Content-Type: application/json

{
    "phone": "13800138000",
    "password": "password123",
    "platform": "web", // 客户端平台: app, web
    "device_info": {
        "device_id": "WebApp_XYZ789",
        "device_type": "web",
        "user_agent": "Mozilla/5.0 (...)"
    }
}
```

**响应示例 (成功):** (同注册成功响应中的 `session` 部分)

### 1.3 Token刷新

```http
POST /api/auth/refresh
Content-Type: application/json

{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "Token刷新成功",
    "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", // 新的access_token
        "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", // Refresh token也可能被轮换
        "expires_in": 86400
    }
}
```

### 1.4 Token验证

```http
GET /api/auth/validate
Authorization: Bearer {access_token}
```

**响应示例 (成功 - Token有效):**
```json
{
    "code": 200,
    "message": "Token有效",
    "data": {
        "valid": true,
        "user_id": 1001,
        "platform": "app"
    }
}
```

---

## 👤 惠农APP/Web - 用户信息接口

**认证要求**: `RequireAuth`
**适用平台**: `app`, `web`

### 2.1 获取当前用户信息

```http
GET /api/user/profile
Authorization: Bearer {access_token}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "id": 1001,
        "uuid": "uuid-abc-123",
        "username": "zhangsan", // 可能为空，根据系统设计
        "phone": "13800138000",
        "email": "zhangsan@example.com",
        "user_type": "farmer",
        "status": "active",
        "real_name": "张三",
        "id_card": "3701...1234", // 脱敏显示
        "avatar": "https://example.com/avatar.jpg",
        "gender": "male",
        "birthday": "1990-01-01",
        "province": "山东省",
        "city": "济南市",
        "county": "历城区",
        "address": "某某村123号",
        "is_real_name_verified": true,
        "is_bank_card_verified": true,
        "is_credit_verified": false,
        "last_login_time": "2024-01-15T14:25:30Z",
        "created_at": "2024-01-01T10:00:00Z"
    }
}
```

### 2.2 更新当前用户信息

```http
PUT /api/user/profile
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "real_name": "张三丰", // 用户可修改字段
    "email": "zhangsan@newmail.com",
    "gender": "male",
    "birthday": "1990-01-01",
    "avatar": "https://new.avatar.url/image.png",
    "province": "北京市",
    "city": "北京市",
    "county": "朝阳区",
    "address": "新地址123号"
}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "更新成功",
    "data": {
        "updated_fields": ["real_name", "email", "avatar", "address"]
    }
}
```

### 2.3 修改密码

```http
PUT /api/user/password  // 注意：原为POST /api/user/change-password，建议统一为PUT
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "old_password": "oldpassword123",
    "new_password": "newpassword123"
}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "密码修改成功"
}
```

### 2.4 用户登出

```http
POST /api/user/logout
Authorization: Bearer {access_token}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "登出成功"
}
```

---

## 🛡️ 惠农APP/Web - 用户认证流程接口

**认证要求**: `RequireAuth`
**适用平台**: `app`, `web`

### 3.1 实名认证申请

```http
POST /api/user/auth/real-name
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "id_card_number": "370123199001011234",
    "real_name": "张三",
    "id_card_front_img_url": "https://example.com/uploads/id_front.jpg", // 文件上传后得到的URL
    "id_card_back_img_url": "https://example.com/uploads/id_back.jpg",
    "face_verify_img_url": "https://example.com/uploads/face.jpg" // 人脸照片或活体检测凭证
}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "实名认证申请已提交，等待审核",
    "data": {
        "auth_id": "auth_realname_uuid123", // 认证记录ID
        "auth_status": "pending"
    }
}
```

### 3.2 银行卡认证申请

```http
POST /api/user/auth/bank-card
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "bank_card_number": "6222020000001234567",
    "bank_name": "中国工商银行",
    "cardholder_name": "张三", // 通常与实名认证姓名一致，后端校验
    "bank_reserved_phone": "13800138000" // 可选，用于四要素验证
}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "银行卡认证申请已提交",
    "data": {
        "auth_id": "auth_bankcard_uuid456",
        "auth_status": "pending"
    }
}
```

---

## 🏢 OA系统 - 认证接口

**适用平台**: `oa`

### 4.1 OA用户登录

```http
POST /api/oa/auth/login
Content-Type: application/json

{
    "username": "oa_admin", // 或 email
    "password": "password123",
    "platform": "oa", // 固定为 "oa"
    "device_info": {
        "device_id": "OA_WebApp_SessionID",
        "device_type": "web",
        "user_agent": "Mozilla/5.0 (...)"
    }
}
```

**响应示例 (成功):** (同惠农端登录成功响应中的 `session` 部分)

### 4.2 OA Token刷新

```http
POST /api/oa/auth/refresh
Content-Type: application/json

{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应示例 (成功):** (同惠农端Token刷新成功响应)

### 4.3 OA Token验证

```http
GET /api/oa/auth/validate
Authorization: Bearer {access_token}
```

**响应示例 (成功 - Token有效且为OA平台):**
```json
{
    "code": 200,
    "message": "OA Token有效",
    "data": {
        "valid": true,
        "user_id": 201, // OA User ID
        "platform": "oa",
        "role": "admin" // 用户角色，例如 admin, staff
    }
}
```

### 4.4 OA用户登出

```http
POST /api/oa/auth/logout
Authorization: Bearer {access_token}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "登出成功"
}
```

---

## 🧑‍💼 OA系统 - 普通用户信息接口

**认证要求**: `RequireAuth`, `CheckPlatform("oa")`
**适用平台**: `oa`
**适用角色**: 所有OA用户 (包括管理员和普通员工)

### 5.1 获取当前OA用户信息

```http
GET /api/oa/user/profile
Authorization: Bearer {access_token}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "id": 201, // OA User ID
        "username": "oa_admin",
        "email": "admin@example.com",
        "phone": "13900139000",
        "real_name": "管理张",
        "avatar": "https://example.com/oa_avatar.jpg",
        "role_id": 1,
        "role_name": "系统管理员", // 角色名称
        "department": "技术部",
        "position": "后端工程师",
        "status": "active",
        "last_login_at": "2024-01-15T10:00:00Z"
    }
}
```

### 5.2 更新当前OA用户信息

```http
PUT /api/oa/user/profile
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "email": "new_admin_email@example.com",
    "phone": "13900139001",
    "avatar": "https://new.oa_avatar.url/image.png"
    // real_name, department, position 通常由管理员修改
}
```

### 5.3 OA用户修改密码

```http
PUT /api/oa/user/password
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "old_password": "oldpassword123",
    "new_password": "newpassword123"
}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "密码修改成功"
}
```

**[更多OA管理员的用户管理接口，请参见 `oa_management.md` 中的用户管理部分]**

---

## 🏷️ 惠农APP/Web - 用户标签管理接口

**认证要求**: `RequireAuth`
**适用平台**: `app`, `web`

### 3.1 获取当前用户的标签列表

```http
GET /api/user/tags
Authorization: Bearer {access_token}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "user_id": 1001,
        "tags": [
            {
                "tag_key": "user_type",
                "tag_value": "farmer",
                "display_name": "农户",
                "category": "基础信息",
                "created_at": "2024-01-01T10:00:00Z",
                "is_system": true,
                "is_editable": false
            },
            {
                "tag_key": "crop_type",
                "tag_value": "rice,wheat",
                "display_name": "种植作物",
                "category": "业务信息",
                "created_at": "2024-01-15T14:30:00Z",
                "is_system": false,
                "is_editable": true
            },
            {
                "tag_key": "farm_scale",
                "tag_value": "medium",
                "display_name": "农场规模",
                "category": "业务信息",
                "created_at": "2024-01-10T09:20:00Z",
                "is_system": false,
                "is_editable": true
            },
            {
                "tag_key": "credit_level",
                "tag_value": "good",
                "display_name": "信用等级",
                "category": "风控标签",
                "created_at": "2024-01-12T16:45:00Z",
                "is_system": true,
                "is_editable": false
            }
        ],
        "categories": {
            "基础信息": ["user_type", "region"],
            "业务信息": ["crop_type", "farm_scale", "business_type"],
            "风控标签": ["credit_level", "risk_score", "loan_history"],
            "行为标签": ["login_frequency", "feature_usage", "active_level"]
        }
    }
}
```

### 3.2 添加用户标签

```http
POST /api/user/tags
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "tag_key": "preferred_contact",
    "tag_value": "sms",
    "category": "偏好设置",
    "description": "用户偏好的联系方式"
}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "标签添加成功",
    "data": {
        "tag_key": "preferred_contact",
        "tag_value": "sms",
        "display_name": "偏好联系方式",
        "category": "偏好设置",
        "created_at": "2024-01-15T17:00:00Z",
        "is_system": false,
        "is_editable": true
    }
}
```

### 3.3 更新用户标签

```http
PUT /api/user/tags/{tag_key}
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "tag_value": "email,sms",
    "description": "更新为多种联系方式"
}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "标签更新成功",
    "data": {
        "tag_key": "preferred_contact",
        "old_value": "sms",
        "new_value": "email,sms",
        "updated_at": "2024-01-15T17:30:00Z"
    }
}
```

### 3.4 删除用户标签

```http
DELETE /api/user/tags/{tag_key}
Authorization: Bearer {access_token}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "标签删除成功",
    "data": {
        "tag_key": "preferred_contact",
        "deleted_at": "2024-01-15T18:00:00Z"
    }
}
```

**错误响应 (系统标签不可删除):**
```json
{
    "code": 1011,
    "message": "系统标签不允许删除",
    "data": {
        "tag_key": "user_type",
        "is_system": true
    }
}
```

### 标签系统说明

#### 标签分类
- **基础信息**: 用户基本属性，如用户类型、地区等
- **业务信息**: 业务相关属性，如种植作物、农场规模等
- **风控标签**: 风险控制相关，如信用等级、风险评分等
- **行为标签**: 用户行为分析，如登录频率、功能使用等
- **偏好设置**: 用户个性化设置

#### 标签类型
- **系统标签** (`is_system: true`): 由系统自动生成和维护，用户不可修改
- **用户标签** (`is_system: false`): 用户可以自定义添加、修改和删除

#### 常用标签示例
```javascript
const COMMON_TAGS = {
    // 基础信息标签
    'user_type': {
        values: ['farmer', 'farm_owner', 'cooperative', 'enterprise'],
        display: '用户类型'
    },
    'region': {
        values: ['华北', '华东', '华南', '西北', '西南', '东北'],
        display: '所在地区'
    },
    
    // 业务信息标签  
    'crop_type': {
        values: ['rice', 'wheat', 'corn', 'soybean', 'vegetable', 'fruit'],
        display: '种植作物'
    },
    'farm_scale': {
        values: ['small', 'medium', 'large', 'extra_large'],
        display: '农场规模'
    },
    'business_type': {
        values: ['planting', 'breeding', 'mixed', 'processing'],
        display: '经营类型'
    },
    
    // 风控标签
    'credit_level': {
        values: ['excellent', 'good', 'fair', 'poor'],
        display: '信用等级'
    },
    'risk_score': {
        values: ['low', 'medium', 'high'],
        display: '风险评分'
    },
    
    // 行为标签
    'login_frequency': {
        values: ['daily', 'weekly', 'monthly', 'occasional'],
        display: '登录频率'
    },
    'active_level': {
        values: ['very_active', 'active', 'normal', 'inactive'],
    }
}
```

### 错误码说明

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