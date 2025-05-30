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