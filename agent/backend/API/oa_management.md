# OA系统管理模块 - API 接口文档

## 📋 模块概述

OA后台管理模块为内部运营和管理人员提供系统管理功能。根据用户角色，提供不同级别的操作权限。

### 平台与角色

-   **适用平台**: `oa` (所有OA接口)
-   **用户模型**: `OAUser` (包含 `RoleID`)
-   **权限划分**:
    -   **普通OA用户**: 拥有基础操作权限，如查看个人信息、提交的申请等。
        -   API路径: `/api/oa/user/*`
        -   认证要求: `RequireAuth`, `CheckPlatform("oa")`
    -   **OA管理员**: 拥有高级管理权限，如用户管理、系统配置、业务审批等。
        -   API路径: `/api/oa/admin/*`
        -   认证要求: `RequireAuth`, `CheckPlatform("oa")`, `RequireRole("admin")`

### 核心功能 (按角色划分)

#### 普通OA用户 (`/api/oa/user/*`)
-   个人信息查看与修改
-   查看自己的业务数据（如贷款申请、农机订单等）

#### OA管理员 (`/api/oa/admin/*`)
-   **用户管理**: 管理所有惠农用户 (`User`) 和OA系统用户 (`OAUser`)。
-   **业务审批**: 贷款申请审批、实名认证审核等。
-   **内容管理**: 发布和管理资讯、政策等。
-   **农机管理**: 管理农机设备信息、租赁订单等。
-   **系统配置**: 系统参数设置、角色权限管理。
-   **数据统计与监控**: 查看业务报表、系统监控数据。
-   **会话管理**: 查看和管理用户会话。

---

## 🔐 OA系统 - 认证接口

**接口路径前缀**: `/api/oa/auth`
**适用平台**: `oa`
**认证要求**: 无 (部分接口如 /validate, /logout 需要先登录)

### 1.1 OA用户登录

```http
POST /api/oa/auth/login
Content-Type: application/json

{
    "username": "oa_admin_user", // 或 email
    "password": "password123",
    "platform": "oa", // 固定为 "oa"
    "device_info": { // 可选，用于审计和设备管理
        "device_id": "OA_WebApp_Session_XYZ",
        "device_type": "web",
        "user_agent": "Mozilla/5.0 (...)"
    }
}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "登录成功",
    "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expires_in": 86400,
        "user_info": { // 登录成功后返回的OA用户信息
            "id": 201,
            "username": "oa_admin_user",
            "real_name": "管理员张三",
            "role": "admin" // 用户角色
        }
    }
}
```

### 1.2 OA Token刷新

```http
POST /api/oa/auth/refresh
Content-Type: application/json

{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 1.3 OA Token验证

```http
GET /api/oa/auth/validate
Authorization: Bearer {access_token}
```

### 1.4 OA用户登出

```http
POST /api/oa/auth/logout
Authorization: Bearer {access_token}
```

---

## 🧑‍💼 OA系统 - 普通用户接口

**接口路径前缀**: `/api/oa/user`
**适用平台**: `oa`
**认证要求**: `RequireAuth`, `CheckPlatform("oa")`

### 2.1 获取当前OA用户信息

```http
GET /api/oa/user/profile
Authorization: Bearer {access_token}
```

(响应示例见 `user_management.md` 中OA用户部分)

### 2.2 更新当前OA用户信息

```http
PUT /api/oa/user/profile
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "email": "new_oa_user_email@example.com",
    "phone": "13900139002",
    "avatar": "https://new.oa_avatar.url/image.png"
}
```

### 2.3 OA用户修改密码

```http
PUT /api/oa/user/password
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "old_password": "oldpassword123",
    "new_password": "newpassword123"
}
```

### 2.4 查看自己提交的贷款申请

```http
GET /api/oa/user/loan/applications?status=pending&page=1&limit=10
Authorization: Bearer {access_token}
```

**说明**: 此接口复用 `/api/user/loan/applications` 的Handler，但通过OA认证和平台检查确保是OA用户访问自己的数据。具体参数和响应请参考 `loan_management.md`。

---

## 🛠️ OA系统 - 管理员接口

**接口路径前缀**: `/api/oa/admin`
**适用平台**: `oa`
**认证要求**: `RequireAuth`, `CheckPlatform("oa")`, `RequireRole("admin")`

### 3.1 用户管理 (管理员)

#### 3.1.1 获取用户列表 (惠农用户和OA用户)

```http
GET /api/oa/admin/users?page=1&limit=20&status=active&user_type=farmer&keyword=张三&platform_user_type=app_user
Authorization: Bearer {access_token}
```

**Query Parameters**:
-   `platform_user_type`: `app_user` (惠农用户), `oa_user` (OA用户), `all` (默认，所有)

**响应示例 (部分):**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total": 1250,
        "users": [
            {
                "id": 1001, // 如果是惠农用户，这里是 User ID
                "user_id_type": "app_user", // 标识用户来源
                "phone": "13800138000",
                "real_name": "张三 (惠农用户)",
                // ... 惠农User模型字段
            },
            {
                "id": 205, // 如果是OA用户，这里是 OAUser ID
                "user_id_type": "oa_user",
                "username": "oa_staff_li",
                "real_name": "李四 (OA员工)",
                "role": "staff",
                // ... OAUser模型字段
            }
        ]
    }
}
```

#### 3.1.2 获取指定用户详情 (惠农/OA)

```http
GET /api/oa/admin/users/{user_platform_id}?user_id_type=app_user
Authorization: Bearer {access_token}
```

**Query Parameters**:
-   `user_id_type`: 必填, `app_user` 或 `oa_user`，用于区分ID类型。

**Path Parameters**:
-   `user_platform_id`: 用户在对应平台上的ID。

#### 3.1.3 更新用户状态 (惠农/OA)

```http
PUT /api/oa/admin/users/{user_platform_id}/status?user_id_type=app_user
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "status": "frozen", // active, frozen
    "reason": "风险操作"
}
```

#### 3.1.4 创建OA系统用户

```http
POST /api/oa/admin/users/oa-user
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "username": "new_oa_staff",
    "password": "staffpassword",
    "email": "staff@example.com",
    "real_name": "新员工王五",
    "phone": "13700137000",
    "role_id": 2, // 对应 OARole 的 ID
    "department": "市场部",
    "position": "市场专员"
}
```

#### 3.1.5 更新OA系统用户信息

```http
PUT /api/oa/admin/users/oa-user/{oa_user_id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "email": "updated_staff@example.com",
    "real_name": "王五更新",
    "role_id": 3,
    "status": "active"
}
```

#### 3.1.6 删除OA系统用户

```http
DELETE /api/oa/admin/users/oa-user/{oa_user_id}
Authorization: Bearer {access_token}
```

### 3.2 贷款审批管理 (管理员)

(详细接口请参考 `loan_management.md` 中标记为管理员操作的部分，路径前缀为 `/api/oa/admin/loans/*`)

**示例接口**:
-   `GET /api/oa/admin/loans/applications` - 获取所有贷款申请列表 (可筛选)
-   `POST /api/oa/admin/loans/applications/{application_id}/approve` - 批准贷款申请
-   `POST /api/oa/admin/loans/applications/{application_id}/reject` - 拒绝贷款申请

### 3.3 实名认证审核 (管理员)

(详细接口请参考 `user_management.md` 或 `identity_auth.md` (如果单独创建) 中标记为管理员操作的部分)

**示例接口**:
-   `GET /api/oa/admin/auth/real-name/pending` - 获取待审核实名认证列表
-   `POST /api/oa/admin/auth/real-name/{auth_id}/approve` - 通过实名认证
-   `POST /api/oa/admin/auth/real-name/{auth_id}/reject` - 셔부实名认证

### 3.4 内容管理 (管理员)

(详细接口请参考 `content_management.md` 中标记为管理员操作的部分，路径前缀为 `/api/oa/admin/content/*`)

### 3.5 农机管理 (管理员)

(详细接口请参考 `machine_rental.md` 中标记为管理员操作的部分，路径前缀为 `/api/oa/admin/machines/*`)

### 3.6 系统配置与管理 (管理员)

#### 3.6.1 获取系统配置

```http
GET /api/oa/admin/system/config
Authorization: Bearer {access_token}
```

#### 3.6.2 更新系统配置

```http
PUT /api/oa/admin/system/config
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "site_name": "新惠农金融平台",
    "default_loan_interest_rate": 0.05
    // ... 其他配置项
}
```

#### 3.6.3 OA角色管理
-   `GET /api/oa/admin/system/roles` - 获取OA角色列表
-   `POST /api/oa/admin/system/roles` - 创建OA角色
-   `PUT /api/oa/admin/system/roles/{role_id}` - 更新OA角色
-   `DELETE /api/oa/admin/system/roles/{role_id}` - 删除OA角色

### 3.7 数据统计与仪表盘 (管理员)

```http
GET /api/oa/admin/dashboard/overview
Authorization: Bearer {access_token}
```

```http
GET /api/oa/admin/dashboard/risk-monitoring
Authorization: Bearer {access_token}
```

### 3.8 会话管理 (管理员)

(详细接口请参考 `session_management.md` 中标记为管理员操作的部分，路径前缀为 `/api/oa/admin/sessions/*`)

**示例接口**:
-   `GET /api/oa/admin/sessions/active` - 获取当前所有活跃会话
-   `POST /api/oa/admin/sessions/{session_id}/revoke` - 强制指定会话下线

---

## 🔧 错误码说明

| 错误码 | 说明 | 处理建议 |
|-------|------|---------|
| 4001 | 管理员账户不存在 | 检查用户名是否正确 |
| 4002 | 验证码错误 | 重新获取验证码 |
| 4003 | 权限不足 | 联系系统管理员 |
| 4004 | 用户不存在 | 检查用户ID |
| 4005 | 状态变更不允许 | 检查用户当前状态 |
| 4006 | 审批申请不存在 | 检查申请ID |
| 4007 | 审批状态不允许操作 | 检查申请状态 |
| 4008 | 配置参数无效 | 检查参数格式 |
| 4009 | 角色代码重复 | 使用不同的角色代码 |
| 4010 | 权限代码无效 | 检查权限代码 |
| 4011 | 农机设备不存在 | 检查设备ID |
| 4012 | 设备状态不允许操作 | 检查设备当前状态 |
| 4013 | 设备所有者不存在 | 检查所有者ID |
| 4014 | 维护记录不存在 | 检查维护记录ID |
| 4015 | 设备正在租赁中 | 等待租赁结束后操作 |

---

## 📝 接口调用示例

### JavaScript示例
```javascript
// 管理员登录
const adminLogin = async (credentials) => {
    const response = await fetch('/api/oa/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(credentials)
    });
    return response.json();
};

// 获取用户列表
const getUserList = async (token, params) => {
    const queryString = new URLSearchParams(params).toString();
    const response = await fetch(`/api/oa/users?${queryString}`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};

// 审批贷款申请
const approveLoan = async (token, applicationId, reviewData) => {
    const response = await fetch(`/api/oa/loan-applications/${applicationId}/review`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(reviewData)
    });
    return response.json();
};
```

### 权限验证中间件示例
```javascript
// 权限检查
const checkPermission = (requiredPermission) => {
    return (req, res, next) => {
        const userPermissions = req.user.permissions;
        if (userPermissions.includes('*') || userPermissions.includes(requiredPermission)) {
            next();
        } else {
            res.status(403).json({
                code: 4003,
                message: '权限不足'
            });
        }
    };
};
```

### 注意事项
1. **权限控制**: 严格的权限验证和操作审计
2. **数据安全**: 敏感信息脱敏处理和安全传输
3. **操作记录**: 所有管理操作都要记录日志
4. **审批流程**: 重要操作需要多级审批
5. **系统监控**: 实时监控系统运行状态和异常 