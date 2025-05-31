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
-   `POST /api/oa/admin/auth/real-name/{auth_id}/reject` - 拒绝实名认证

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

# OA后台管理模块 API 文档

## 📋 模块概述

OA后台管理模块为管理员提供全面的业务管理功能，包括用户管理、认证审核、数据统计、系统配置等。支持普通OA用户和管理员两种角色，实现分层权限管理和精细化业务控制。

### 🎯 核心功能
- **用户管理**: 用户列表、状态管理、权限控制、批量操作
- **认证审核**: 实名认证、银行卡审核、批量处理、数据导出
- **数据统计**: 业务报表、风险监控、用户分析、趋势预测
- **系统配置**: 参数设置、健康检查、性能监控
- **工作台**: 仪表盘、数据概览、风险监控、待办任务

### 🏗️ 权限架构
```
OA系统权限体系
├── 平台认证 (CheckPlatform: "oa")
├── 角色权限 (RequireRole: "admin")
└── 功能权限
    ├── 用户管理权限
    ├── 业务审核权限
    ├── 数据查看权限
    └── 系统配置权限
```

### 📊 数据模型关系
```
OAUsers (OA用户)
├── OARoles (角色权限)
├── OASessions (OA会话)
└── OAOperationLogs (操作日志)

BusinessData (业务数据)
├── UserApplications (用户申请)
├── AuthenticationRecords (认证记录)
├── StatisticsReports (统计报表)
└── SystemConfigurations (系统配置)
```

---

## 👤 OA用户个人功能

### 1. 获取OA用户信息
**接口路径**: `GET /api/oa/user/profile`  
**认证要求**: 需要认证 (OA用户)  
**功能描述**: 获取当前OA用户的个人信息

#### 响应示例
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "user_info": {
            "user_id": "OA20240115001",
            "username": "admin001",
            "name": "李管理员",
            "email": "admin001@huinong.com",
            "phone": "13800138000",
            "role": "admin",
            "role_name": "系统管理员",
            "department": "业务管理部",
            "position": "高级审核员",
            "avatar": "https://oss.example.com/avatars/admin_li.jpg",
            "status": "active",
            "created_at": "2023-06-15T00:00:00Z",
            "last_login": "2024-01-15T08:00:00Z",
            "login_count": 245
        },
        "permissions": [
            "user_management",
            "loan_approval", 
            "auth_review",
            "system_config",
            "data_export"
        ],
        "work_statistics": {
            "processed_applications": 156,
            "approved_applications": 134,
            "rejected_applications": 22,
            "approval_rate": 85.9,
            "average_process_time": "2.5小时",
            "work_efficiency_score": 94.2
        },
        "recent_activities": [
            {
                "activity": "审批贷款申请",
                "target": "申请编号 LA20240115001",
                "result": "批准",
                "timestamp": "2024-01-15T14:30:00Z"
            },
            {
                "activity": "审核实名认证",
                "target": "用户 张三",
                "result": "通过",
                "timestamp": "2024-01-15T13:45:00Z"
            }
        ]
    }
}
```

#### JavaScript调用示例
```javascript
// 获取OA用户信息
async function getOAUserProfile() {
    try {
        const response = await fetch('/api/oa/user/profile', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${oaToken}`,
                'Content-Type': 'application/json'
            }
        });
        
        const result = await response.json();
        if (result.code === 200) {
            console.log('OA用户信息:', result.data);
            return result.data;
        } else {
            throw new Error(result.message);
        }
    } catch (error) {
        console.error('获取OA用户信息失败:', error);
        throw error;
    }
}
```

### 2. 更新OA用户信息
**接口路径**: `PUT /api/oa/user/profile`  
**认证要求**: 需要认证 (OA用户)  
**功能描述**: 更新OA用户个人信息

#### 请求参数
```json
{
    "name": "李管理员",
    "email": "admin001@huinong.com",
    "phone": "13800138000",
    "avatar": "https://oss.example.com/avatars/admin_li_new.jpg",
    "signature": "认真负责，高效审核"
}
```

#### 响应示例
```json
{
    "code": 200,
    "message": "信息更新成功",
    "data": {
        "user_id": "OA20240115001",
        "updated_fields": ["email", "avatar"],
        "updated_at": "2024-01-15T15:00:00Z"
    }
}
```

### 3. 修改OA用户密码
**接口路径**: `PUT /api/oa/user/password`  
**认证要求**: 需要认证 (OA用户)  
**功能描述**: 修改OA用户登录密码

#### 请求参数
```json
{
    "old_password": "oldPassword123",
    "new_password": "newPassword456",
    "confirm_password": "newPassword456"
}
```

#### 响应示例
```json
{
    "code": 200,
    "message": "密码修改成功",
    "data": {
        "user_id": "OA20240115001",
        "updated_at": "2024-01-15T15:30:00Z",
        "security_tips": [
            "建议定期更换密码",
            "使用强密码组合",
            "避免在多个系统使用相同密码"
        ]
    }
}
```

### 4. 查看个人申请
**接口路径**: `GET /api/oa/user/loan/applications`  
**认证要求**: 需要认证 (OA用户)  
**功能描述**: OA用户查看自己提交的贷款申请

#### 请求参数
```
?page={page}              # 页码，默认1
&limit={limit}            # 每页数量，默认10
&status={status}          # 状态筛选
&date_from={date}         # 日期起始
&date_to={date}           # 日期结束
```

#### 响应示例
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "applications": [
            {
                "application_id": "LA20240115010",
                "application_number": "OA202401150001",
                "product_name": "员工专项贷",
                "application_amount": 50000,
                "status": "approved",
                "status_text": "已批准",
                "submitted_at": "2024-01-10T09:00:00Z",
                "approved_at": "2024-01-12T16:30:00Z",
                "approved_amount": 50000
            }
        ],
        "pagination": {
            "page": 1,
            "limit": 10,
            "total": 1,
            "pages": 1
        }
    }
}
```

---

## 🛡️ 认证审核管理

### 5. 获取认证申请列表
**接口路径**: `GET /api/oa/admin/auth/list`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 获取用户认证申请列表

#### 请求参数
```
?page={page}              # 页码，默认1
&limit={limit}            # 每页数量，默认20
&auth_type={type}         # 认证类型筛选 (real_name/bank_card)
&status={status}          # 状态筛选 (pending/approved/rejected)
&submitted_from={date}    # 提交日期起始
&submitted_to={date}      # 提交日期结束
&user_search={keyword}    # 用户搜索
&reviewer={reviewer}      # 审核员筛选
&sort_by={field}          # 排序字段
&sort_order={desc|asc}    # 排序方向
```

#### 响应示例
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "auth_applications": [
            {
                "auth_id": "AUTH20240115001",
                "user_info": {
                    "user_id": "HN20240115001",
                    "user_name": "张三",
                    "user_phone": "13800138000",
                    "registration_date": "2023-08-15T00:00:00Z"
                },
                "auth_type": "real_name",
                "auth_type_name": "实名认证",
                "status": "pending",
                "status_text": "待审核",
                "submitted_at": "2024-01-15T10:30:00Z",
                "priority": "normal",
                "documents_count": 2,
                "estimated_review_time": "2小时",
                "assigned_reviewer": null,
                "days_pending": 0
            },
            {
                "auth_id": "AUTH20240115002",
                "user_info": {
                    "user_id": "HN20240115002",
                    "user_name": "李四",
                    "user_phone": "13900139000",
                    "registration_date": "2023-09-20T00:00:00Z"
                },
                "auth_type": "bank_card",
                "auth_type_name": "银行卡认证",
                "status": "approved",
                "status_text": "已通过",
                "submitted_at": "2024-01-14T14:20:00Z",
                "reviewed_at": "2024-01-14T16:45:00Z",
                "reviewer": "李审核员",
                "review_time_hours": 2.4
            }
        ],
        "pagination": {
            "page": 1,
            "limit": 20,
            "total": 89,
            "pages": 5
        },
        "statistics": {
            "total_applications": 89,
            "pending_review": 23,
            "approved_today": 12,
            "rejected_today": 3,
            "average_review_time": "1.8小时",
            "approval_rate": 87.6
        }
    }
}
```

### 6. 获取认证详情
**接口路径**: `GET /api/oa/admin/auth/{auth_id}`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 获取指定认证申请的详细信息

#### 路径参数
- `auth_id`: 认证申请ID

#### 响应示例
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "auth_info": {
            "auth_id": "AUTH20240115001",
            "auth_type": "real_name",
            "auth_type_name": "实名认证",
            "status": "pending",
            "submitted_at": "2024-01-15T10:30:00Z",
            "auto_review_result": {
                "passed": true,
                "confidence": 0.95,
                "flags": []
            }
        },
        "user_profile": {
            "user_id": "HN20240115001",
            "user_name": "张三",
            "user_phone": "13800138000",
            "email": "zhangsan@example.com",
            "registration_date": "2023-08-15T00:00:00Z",
            "kyc_history": {
                "previous_applications": 0,
                "total_rejections": 0,
                "last_application": null
            },
            "user_activity": {
                "login_frequency": "活跃",
                "business_transactions": 5,
                "risk_level": "低"
            }
        },
        "submitted_data": {
            "real_name": "张三",
            "id_card_number": "370102199001151234",
            "documents": [
                {
                    "document_type": "id_card_front",
                    "document_name": "身份证正面",
                    "file_url": "https://oss.example.com/auth/id_front_123.jpg",
                    "upload_time": "2024-01-15T10:25:00Z",
                    "file_size": "2.3MB",
                    "ocr_result": {
                        "name": "张三",
                        "id_number": "370102199001151234",
                        "confidence": 0.98
                    }
                },
                {
                    "document_type": "id_card_back",
                    "document_name": "身份证背面", 
                    "file_url": "https://oss.example.com/auth/id_back_123.jpg",
                    "upload_time": "2024-01-15T10:26:00Z",
                    "file_size": "2.1MB",
                    "ocr_result": {
                        "issue_authority": "济南市公安局",
                        "valid_period": "2015.01.15-2025.01.15",
                        "confidence": 0.96
                    }
                }
            ]
        },
        "verification_results": {
            "document_quality": {
                "clarity": "优秀",
                "completeness": "完整",
                "authenticity": "真实"
            },
            "identity_verification": {
                "name_match": true,
                "id_number_valid": true,
                "face_match": true,
                "public_security_check": "通过"
            },
            "risk_assessment": {
                "risk_level": "低",
                "risk_score": 85,
                "risk_factors": []
            }
        },
        "review_suggestions": {
            "system_recommendation": "approve",
            "confidence_level": "高",
            "attention_points": [],
            "suggested_actions": ["批准认证申请"]
        }
    }
}
```

### 7. 审核认证申请
**接口路径**: `POST /api/oa/admin/auth/{auth_id}/review`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 审核用户认证申请

#### 路径参数
- `auth_id`: 认证申请ID

#### 请求参数
```json
{
    "review_result": "approved",
    "review_comments": "材料齐全，信息核实无误，同意通过认证",
    "verified_name": "张三",
    "verified_id_number": "370102199001151234",
    "notes": "用户身份信息真实有效",
    "quality_score": 95,
    "special_notes": []
}
```

#### 响应示例
```json
{
    "code": 200,
    "message": "审核完成",
    "data": {
        "auth_id": "AUTH20240115001",
        "review_result": "approved",
        "reviewed_by": "李审核员",
        "reviewed_at": "2024-01-15T16:30:00Z",
        "user_notification": {
            "notification_sent": true,
            "notification_methods": ["sms", "app_push"],
            "sent_at": "2024-01-15T16:31:00Z"
        },
        "next_steps": [
            "用户认证状态已更新",
            "用户可使用认证相关功能",
            "系统已自动发送通知"
        ]
    }
}
```

### 8. 批量审核
**接口路径**: `POST /api/oa/admin/auth/batch-review`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 批量审核多个认证申请

#### 请求参数
```json
{
    "auth_ids": ["AUTH20240115001", "AUTH20240115003", "AUTH20240115004"],
    "review_result": "approved",
    "review_comments": "材料完整，批量通过认证",
    "auto_notify": true
}
```

#### 响应示例
```json
{
    "code": 200,
    "message": "批量审核完成",
    "data": {
        "total_count": 3,
        "success_count": 3,
        "failed_count": 0,
        "results": [
            {
                "auth_id": "AUTH20240115001",
                "result": "approved",
                "status": "success"
            },
            {
                "auth_id": "AUTH20240115003",
                "result": "approved", 
                "status": "success"
            },
            {
                "auth_id": "AUTH20240115004",
                "result": "approved",
                "status": "success"
            }
        ],
        "batch_id": "BATCH_REV_001",
        "processed_by": "李审核员",
        "processed_at": "2024-01-15T17:00:00Z"
    }
}
```

### 9. 获取认证统计
**接口路径**: `GET /api/oa/admin/auth/statistics`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 获取认证审核统计数据

#### 请求参数
```
?period={period}          # 统计周期 (day/week/month/year)
&start_date={date}        # 统计开始日期
&end_date={date}          # 统计结束日期
&auth_type={type}         # 认证类型筛选
```

#### 响应示例
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "overview": {
            "total_applications": 456,
            "pending_applications": 23,
            "approved_applications": 378,
            "rejected_applications": 55,
            "approval_rate": 87.3,
            "average_review_time": "1.8小时"
        },
        "auth_type_breakdown": {
            "real_name": {
                "total": 234,
                "approved": 203,
                "rejected": 31,
                "approval_rate": 86.8
            },
            "bank_card": {
                "total": 222,
                "approved": 175,
                "rejected": 24, 
                "approval_rate": 87.9
            }
        },
        "daily_trends": [
            {
                "date": "2024-01-15",
                "submitted": 12,
                "reviewed": 15,
                "approved": 13,
                "rejected": 2,
                "pending": 23
            }
        ],
        "reviewer_performance": [
            {
                "reviewer": "李审核员",
                "reviewed_count": 89,
                "approval_count": 78,
                "rejection_count": 11,
                "average_review_time": "1.5小时",
                "accuracy_rate": 97.8
            }
        ],
        "quality_metrics": {
            "auto_review_accuracy": 94.5,
            "manual_review_consistency": 92.1,
            "user_appeal_rate": 2.3,
            "fraud_detection_rate": 0.8
        }
    }
}
```

### 10. 导出认证数据
**接口路径**: `GET /api/oa/admin/auth/export`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 导出认证数据报表

#### 请求参数
```
?format={format}          # 导出格式 (excel/csv)
&auth_type={type}         # 认证类型筛选
&status={status}          # 状态筛选
&date_from={date}         # 日期起始
&date_to={date}           # 日期结束
&fields={fields}          # 导出字段 (逗号分隔)
```

#### 响应示例
```json
{
    "code": 200,
    "message": "导出任务已创建",
    "data": {
        "export_id": "EXPORT_001",
        "export_format": "excel",
        "total_records": 456,
        "estimated_time": "2-3分钟",
        "status": "processing",
        "download_url": null,
        "created_at": "2024-01-15T17:30:00Z",
        "expires_at": "2024-01-16T17:30:00Z"
    }
}
```

---

## 📊 工作台和统计

### 11. 获取工作台数据
**接口路径**: `GET /api/oa/admin/dashboard`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 获取工作台仪表盘数据

#### 响应示例
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "overview_stats": {
            "total_users": 1567,
            "new_users_today": 23,
            "active_users": 1234,
            "pending_reviews": 45,
            "total_loans": 2340000,
            "loan_applications_today": 12
        },
        "pending_tasks": {
            "auth_reviews": 23,
            "loan_reviews": 15,
            "machine_approvals": 8,
            "user_appeals": 3,
            "system_alerts": 2
        },
        "recent_activities": [
            {
                "activity_type": "loan_approval",
                "description": "批准了贷款申请 LA20240115001",
                "operator": "李审核员",
                "timestamp": "2024-01-15T16:30:00Z"
            },
            {
                "activity_type": "auth_review", 
                "description": "审核通过实名认证申请",
                "operator": "王审核员",
                "timestamp": "2024-01-15T16:15:00Z"
            }
        ],
        "system_health": {
            "api_status": "normal",
            "database_status": "normal",
            "redis_status": "normal",
            "file_storage_status": "normal",
            "response_time": "150ms",
            "error_rate": "0.02%"
        },
        "performance_metrics": {
            "daily_transactions": 1234,
            "success_rate": 99.8,
            "average_response_time": "150ms",
            "concurrent_users": 156
        }
    }
}
```

### 12. 获取业务概览
**接口路径**: `GET /api/oa/admin/dashboard/overview`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 获取详细的业务概览数据

#### 请求参数
```
?period={period}          # 统计周期 (day/week/month)
&compare_previous={bool}  # 是否对比上期数据
```

#### 响应示例
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "user_metrics": {
            "total_users": 1567,
            "growth_rate": 12.5,
            "new_registrations": 89,
            "active_users": 1234,
            "verified_users": 987,
            "user_retention_rate": 78.6
        },
        "business_metrics": {
            "total_loan_amount": 2340000,
            "loan_applications": 156,
            "approval_rate": 85.6,
            "average_loan_amount": 125000,
            "machine_rentals": 89,
            "rental_revenue": 156000
        },
        "operational_metrics": {
            "pending_reviews": 45,
            "average_review_time": "1.8小时",
            "reviewer_efficiency": 94.2,
            "customer_satisfaction": 4.7,
            "system_uptime": 99.9
        },
        "comparison_data": {
            "previous_period": {
                "user_growth": 8.3,
                "loan_amount": 2100000,
                "approval_rate": 82.1
            },
            "trends": {
                "user_growth": "上升",
                "business_volume": "上升",
                "efficiency": "稳定"
            }
        },
        "charts_data": {
            "user_growth_chart": [
                {"date": "2024-01-08", "new_users": 15, "total_users": 1478},
                {"date": "2024-01-09", "new_users": 18, "total_users": 1496},
                {"date": "2024-01-10", "new_users": 22, "total_users": 1518}
            ],
            "loan_trend_chart": [
                {"date": "2024-01-08", "applications": 8, "amount": 980000},
                {"date": "2024-01-09", "applications": 12, "amount": 1450000},
                {"date": "2024-01-10", "applications": 15, "amount": 1820000}
            ]
        }
    }
}
```

### 13. 风险监控数据
**接口路径**: `GET /api/oa/admin/dashboard/risk-monitoring`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 获取风险监控和预警数据

#### 响应示例
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "risk_overview": {
            "total_risk_alerts": 5,
            "high_risk_users": 3,
            "suspicious_transactions": 2,
            "fraud_attempts": 1,
            "system_anomalies": 1
        },
        "risk_alerts": [
            {
                "alert_id": "ALERT_001",
                "alert_type": "suspicious_login",
                "severity": "medium",
                "description": "用户 HN20240115001 在异常地点登录",
                "affected_user": "张三",
                "occurred_at": "2024-01-15T14:30:00Z",
                "status": "investigating",
                "assigned_to": "安全团队"
            },
            {
                "alert_id": "ALERT_002",
                "alert_type": "high_risk_application",
                "severity": "high",
                "description": "贷款申请 LA20240115005 风险评分异常",
                "affected_application": "LA20240115005",
                "occurred_at": "2024-01-15T15:20:00Z",
                "status": "pending_review",
                "assigned_to": "风控审核员"
            }
        ],
        "risk_trends": {
            "fraud_trend": [
                {"date": "2024-01-08", "attempts": 2, "blocked": 2},
                {"date": "2024-01-09", "attempts": 1, "blocked": 1},
                {"date": "2024-01-10", "attempts": 3, "blocked": 3}
            ],
            "risk_distribution": {
                "low_risk": 89.2,
                "medium_risk": 8.5,
                "high_risk": 2.3
            }
        },
        "security_metrics": {
            "login_security_score": 94.5,
            "transaction_security_score": 96.8,
            "data_protection_score": 98.2,
            "overall_security_score": 96.5
        },
        "recommended_actions": [
            {
                "priority": "high",
                "action": "加强异常登录监控",
                "description": "建议对异地登录增加验证步骤"
            },
            {
                "priority": "medium",
                "action": "优化风险评估模型",
                "description": "根据最新数据调整AI评估参数"
            }
        ]
    }
}
```

---

## ⚠️ 错误码说明

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 4001 | OA用户不存在 | 检查用户ID是否正确 |
| 4002 | 权限不足 | 检查用户角色和权限配置 |
| 4003 | 认证申请不存在 | 检查认证申请ID是否正确 |
| 4004 | 审核状态不允许操作 | 检查认证申请当前状态 |
| 4005 | 批量操作部分失败 | 查看详细错误信息 |
| 4006 | 导出任务创建失败 | 检查参数和系统状态 |
| 4007 | 系统维护中 | 等待系统维护完成 |
| 4008 | 操作日志记录失败 | 检查日志系统状态 |
| 4009 | 数据访问限制 | 检查数据权限配置 |
| 4010 | 并发操作冲突 | 刷新页面后重试 |

---

## 🔄 最佳实践

### 审核管理
1. **及时处理**: 优先处理紧急和高优先级的审核申请
2. **详细记录**: 审核过程中记录详细的审核意见和依据
3. **一致性标准**: 保持审核标准的一致性和公正性
4. **质量控制**: 定期抽查审核质量，持续改进流程

### 数据安全
1. **权限控制**: 严格控制数据访问权限，实行最小权限原则
2. **操作日志**: 记录所有关键操作，便于审计和追溯
3. **数据脱敏**: 敏感数据在展示时进行适当脱敏
4. **安全监控**: 实时监控异常操作和安全威胁

### 系统运维
1. **性能监控**: 定期监控系统性能指标，及时发现问题
2. **数据备份**: 定期备份重要数据，制定灾难恢复方案
3. **版本管理**: 系统更新时做好版本控制和回滚准备
4. **用户培训**: 定期组织OA用户培训，提升操作技能

### 业务优化
1. **流程改进**: 根据业务数据分析，持续优化业务流程
2. **用户体验**: 关注用户反馈，改善后台操作体验
3. **决策支持**: 利用数据分析为业务决策提供支持
4. **风险防控**: 建立完善的风险预警和防控机制 