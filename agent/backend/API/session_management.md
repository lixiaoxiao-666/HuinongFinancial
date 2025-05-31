# 会话管理系统 - API 指南

## 🎯 系统概述

数字惠农系统采用基于Redis的统一分布式会话管理机制。该系统为所有平台（惠农APP、惠农Web、OA后台）提供安全、高效的认证与会话控制服务。

### ✨ 核心特性
-   **统一认证**: 所有平台共享一套核心认证逻辑，简化开发与维护。
-   **分布式会话**: 用户会话状态存储于Redis，支持多后端实例水平扩展。
-   **多平台支持**: 用户可以在不同设备和平台（`app`, `web`, `oa`）同时登录，会话独立管理。
-   **Token机制**: 使用 `access_token` (短效) 和 `refresh_token` (长效)保障安全性与用户体验。
-   **实时控制**: 支持管理员强制下线指定会话或用户所有会话。
-   **安全增强**: Token哈希存储、会话自动清理、设备信息记录等。

---

## 🔑 认证API (`/api/auth/*` 和 `/api/oa/auth/*`)

系统提供两组主要的认证入口，分别服务于惠农端和OA端。

### 1. 惠农APP/Web端认证接口

**接口路径前缀**: `/api/auth`
**适用平台**: `app`, `web`
**说明**: 这些接口用于惠农APP和Web端的普通用户注册、登录、Token管理等。

#### 1.1 用户注册

```http
POST /api/auth/register
Content-Type: application/json

{
    "phone": "13800138000",
    "password": "password123",
    "verification_code": "123456", // 可选，视系统配置
    "user_type": "farmer",         // 用户类型: farmer, farm_owner, etc.
    "real_name": "张三",
    "platform": "app",             // 注册平台: app, web
    "device_info": {              // 设备信息 (可选但推荐)
        "device_id": "unique_device_identifier",
        "device_type": "ios", // ios, android, web
        "device_name": "张三的iPhone",
        "app_version": "1.0.1"
    }
}
```

**响应 (成功):**
```json
{
    "code": 200,
    "message": "注册成功",
    "data": {
        "user": { "id": 101, "phone": "13800138000", ... },
        "session": {
            "access_token": "eyJhbGci...",
            "refresh_token": "eyJhbGci...",
            "expires_in": 86400
        }
    }
}
```

#### 1.2 用户登录 (手机号/邮箱 + 密码)

```http
POST /api/auth/login
Content-Type: application/json

{
    "phone": "13800138000", // 或 "email": "user@example.com"
    "password": "password123",
    "platform": "web",
    "device_info": { ... }
}
```

**响应 (成功):** (同上，包含 `user` 和 `session` 信息)

#### 1.3 刷新Access Token

```http
POST /api/auth/refresh
Content-Type: application/json

{
    "refresh_token": "eyJhbGci..."
}
```

**响应 (成功):**
```json
{
    "code": 200,
    "message": "Token刷新成功",
    "data": {
        "access_token": "new_eyJhbGci...",
        "refresh_token": "potentially_new_eyJhbGci...", // Refresh token也可能被轮换
        "expires_in": 86400
    }
}
```

#### 1.4 验证当前Token有效性

```http
GET /api/auth/validate
Authorization: Bearer {access_token}
```

**响应 (成功 - Token有效):**
```json
{
    "code": 200,
    "message": "Token有效",
    "data": {
        "valid": true,
        "user_id": 101,
        "session_id": "sess_xyz",
        "platform": "app",
        "expires_at": "2024-01-16T10:30:00Z"
    }
}
```

### 2. OA系统认证接口

**接口路径前缀**: `/api/oa/auth`
**适用平台**: `oa`
**说明**: 这些接口专用于OA系统用户的登录和Token管理。

#### 2.1 OA用户登录

```http
POST /api/oa/auth/login
Content-Type: application/json

{
    "username": "oa_admin_user", // OA用户名或邮箱
    "password": "oa_password123",
    "platform": "oa", // 固定为 "oa"
    "device_info": { ... } // OA前端的设备信息
}
```

**响应 (成功):**
```json
{
    "code": 200,
    "message": "登录成功",
    "data": {
        "user": { "id": 201, "username": "oa_admin_user", "role": "admin", ... }, // OA用户信息
        "session": {
            "access_token": "oa_eyJhbGci...",
            "refresh_token": "oa_eyJhbGci...",
            "expires_in": 86400
        }
    }
}
```

#### 2.2 刷新OA Access Token

```http
POST /api/oa/auth/refresh
Content-Type: application/json

{
    "refresh_token": "oa_eyJhbGci..."
}
```

**响应 (成功):** (结构同惠农端刷新，但Token为OA专用)

#### 2.3 验证当前OA Token有效性

```http
GET /api/oa/auth/validate
Authorization: Bearer {oa_access_token}
```

**响应 (成功 - Token有效且为OA平台):**
```json
{
    "code": 200,
    "message": "OA Token有效",
    "data": {
        "valid": true,
        "user_id": 201, // OA User ID
        "session_id": "oa_sess_abc",
        "platform": "oa",
        "role": "admin", // OA用户角色
        "expires_at": "2024-01-16T11:00:00Z"
    }
}
```

---

## 💼 会话操作API (`/api/user/session/*` 和 `/api/oa/admin/sessions/*`)

这些接口用于管理用户的活动会话。

### 3. 惠农APP/Web端 - 用户会话管理

**接口路径前缀**: `/api/user/session`
**适用平台**: `app`, `web`
**认证要求**: `RequireAuth`

#### 3.1 获取当前用户的所有活动会话

```http
GET /api/user/session/list // 或 /api/user/session/info (保持兼容性)
Authorization: Bearer {access_token}
```

**响应 (成功):**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": [
        {
            "session_id": "sess_abc",
            "platform": "app",
            "device_info": { "device_name": "张三的iPhone", ... },
            "ip_address": "1.2.3.4",
            "location": "北京市",
            "last_active_at": "2024-01-15T14:00:00Z",
            "is_current": true // 标识是否为当前请求的会话
        },
        {
            "session_id": "sess_def",
            "platform": "web",
            "device_info": { "device_name": "Chrome浏览器", ... },
            "ip_address": "2.3.4.5",
            "location": "上海市",
            "last_active_at": "2024-01-14T10:00:00Z",
            "is_current": false
        }
    ]
}
```

#### 3.2 用户主动登出 (注销当前会话)

```http
POST /api/user/logout
Authorization: Bearer {access_token}
```

**响应 (成功):**
```json
{
    "code": 200,
    "message": "登出成功"
}
```

#### 3.3 注销指定会话 (例如：在设备管理列表中操作)

```http
POST /api/user/session/revoke  // 建议用 POST 或 DELETE
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "session_id_to_revoke": "sess_def"
}
```

**响应 (成功):**
```json
{
    "code": 200,
    "message": "会话 sess_def 已注销"
}
```

#### 3.4 注销除当前会话外的其他所有会话

```http
POST /api/user/session/revoke-others
Authorization: Bearer {access_token}
```

**响应 (成功):**
```json
{
    "code": 200,
    "message": "其他会话已成功注销",
    "data": {
        "revoked_count": 1
    }
}
```

### 4. OA系统 - 管理员会话管理

**接口路径前缀**: `/api/oa/admin/sessions`
**适用平台**: `oa`
**认证要求**: `RequireAuth`, `CheckPlatform("oa")`, `RequireRole("admin")`

#### 4.1 获取系统所有活跃会话 (可筛选)

```http
GET /api/oa/admin/sessions/active?platform=app&user_id=101&page=1&limit=20
Authorization: Bearer {oa_access_token}
```

**Query Parameters**:
-   `user_id` (uint64, 可选): 按用户ID筛选 (可以是惠农用户ID或OA用户ID，取决于 `user_id_type`)
-   `user_id_type` (string, 可选, `app_user` 或 `oa_user`): 当提供 `user_id` 时，指明其类型。
-   `platform` (string, 可选): 按平台筛选 (`app`, `web`, `oa`)
-   `ip_address` (string, 可选): 按IP地址筛选
-   `page`, `limit` (int, 可选): 分页参数

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "获取活跃会话列表",
    "data": {
        "total": 150,
        "sessions": [
            {
                "session_id": "sess_abc",
                "user_id": 101,
                "user_real_name": "张三 (惠农用户)",
                "platform": "app",
                "device_name": "张三的iPhone",
                "ip_address": "1.2.3.4",
                "location": "北京市",
                "login_time": "2024-01-15T09:00:00Z",
                "last_active_at": "2024-01-15T14:00:00Z",
                "duration_minutes": 300,
                "user_agent": "HuinongApp/1.3.1 (iOS 17.0)"
            },
            {
                "session_id": "oa_sess_def",
                "user_id": 201,
                "user_real_name": "管理员李四 (OA用户)",
                "platform": "oa",
                "device_name": "Chrome浏览器",
                "ip_address": "192.168.1.100",
                "location": "内网",
                "login_time": "2024-01-15T08:30:00Z",
                "last_active_at": "2024-01-15T14:15:00Z",
                "duration_minutes": 345,
                "user_agent": "Mozilla/5.0..."
            }
        ]
    }
}
```

#### 4.2 获取会话统计信息

```http
GET /api/oa/admin/sessions/statistics
Authorization: Bearer {oa_access_token}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total_active_sessions": 150,
        "platform_distribution": {
            "app": 100,
            "web": 30,
            "oa": 20
        },
        "daily_peak_users": 120,
        "average_session_duration_minutes": 30,
        "today_stats": {
            "new_sessions": 45,
            "expired_sessions": 28,
            "active_users": 156
        },
        "hourly_distribution": {
            "00": 5, "01": 2, "02": 1, "03": 1,
            "08": 25, "09": 45, "10": 65, "11": 78,
            "14": 89, "15": 95, "16": 87, "17": 76,
            "20": 45, "21": 32, "22": 18, "23": 12
        }
    }
}
```

#### 4.3 手动清理过期会话

```http
POST /api/oa/admin/sessions/cleanup
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "cleanup_type": "expired", // expired, inactive, all
    "inactive_threshold_hours": 24 // 可选，指定不活跃时间阈值
}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "会话清理完成",
    "data": {
        "cleaned": 23,
        "cleanup_type": "expired",
        "cleanup_time": "2024-01-15T16:30:00Z",
        "details": {
            "expired_sessions": 15,
            "inactive_sessions": 8,
            "total_before": 173,
            "total_after": 150
        }
    }
}
```

#### 4.4 管理员强制注销指定会话

```http
DELETE /api/oa/admin/sessions/{session_id}
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "reason": "安全检查",
    "notify_user": true // 是否通知用户
}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "会话已强制注销",
    "data": {
        "session_id": "sess_abc",
        "user_info": {
            "user_id": 101,
            "real_name": "张三",
            "platform": "app"
        },
        "revoked_at": "2024-01-15T16:45:00Z",
        "revoked_by": {
            "admin_id": 201,
            "admin_name": "管理员李四"
        },
        "reason": "安全检查",
        "user_notified": true
    }
}
```

#### 4.5 批量强制注销会话

```http
POST /api/oa/admin/sessions/batch-revoke
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "session_ids": ["sess_abc", "sess_def", "sess_ghi"],
    "reason": "系统维护",
    "notify_users": true
}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "批量注销完成",
    "data": {
        "total_requested": 3,
        "successful_revokes": 2,
        "failed_revokes": 1,
        "results": [
            {
                "session_id": "sess_abc",
                "status": "success",
                "revoked_at": "2024-01-15T17:00:00Z"
            },
            {
                "session_id": "sess_def", 
                "status": "success",
                "revoked_at": "2024-01-15T17:00:00Z"
            },
            {
                "session_id": "sess_ghi",
                "status": "failed",
                "error": "Session not found"
            }
        ],
        "revoked_by": {
            "admin_id": 201,
            "admin_name": "管理员李四"
        }
    }
}
```

#### 4.6 管理员强制注销指定用户的所有会话

```http
POST /api/oa/admin/sessions/revoke-user
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "user_id_to_revoke": 101, // 惠农用户ID
    "user_id_type": "app_user", // 或 "oa_user" 及对应的OA用户ID
    "reason": "账户异常活动",
    "exclude_current": true, // 是否排除当前会话（如果是OA管理员注销自己）
    "notify_user": true
}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "用户 惠农用户ID:101 的所有会话已被强制注销",
    "data": {
        "user_id": 101,
        "user_id_type": "app_user",
        "user_real_name": "张三",
        "revoked_count": 2,
        "revoked_sessions": [
            {
                "session_id": "sess_abc",
                "platform": "app",
                "device_name": "张三的iPhone"
            },
            {
                "session_id": "sess_def",
                "platform": "web", 
                "device_name": "Chrome浏览器"
            }
        ],
        "revoked_at": "2024-01-15T17:15:00Z",
        "revoked_by": {
            "admin_id": 201,
            "admin_name": "管理员李四"
        },
        "reason": "账户异常活动"
    }
}
```

#### 4.7 获取会话详细信息 (管理员)

```http
GET /api/oa/admin/sessions/{session_id}/details
Authorization: Bearer {oa_access_token}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "session_id": "sess_abc",
        "user_info": {
            "user_id": 101,
            "user_id_type": "app_user",
            "real_name": "张三",
            "phone": "138****8000",
            "user_type": "farmer"
        },
        "session_details": {
            "platform": "app",
            "device_info": {
                "device_id": "iPhone_12_ABC123",
                "device_type": "ios",
                "device_name": "张三的iPhone",
                "app_version": "1.3.1",
                "os_version": "iOS 17.0"
            },
            "network_info": {
                "ip_address": "1.2.3.4",
                "location": "北京市海淀区",
                "isp": "中国联通",
                "ip_type": "mobile"
            },
            "session_timeline": {
                "created_at": "2024-01-15T09:00:00Z",
                "last_active_at": "2024-01-15T14:00:00Z",
                "expires_at": "2024-01-16T09:00:00Z",
                "total_duration_minutes": 300,
                "idle_duration_minutes": 15
            },
            "security_info": {
                "login_method": "password",
                "is_trusted_device": true,
                "risk_score": "low",
                "security_events": []
            }
        },
        "activity_summary": {
            "api_calls_count": 156,
            "last_endpoint": "/api/user/loan/applications",
            "most_used_features": ["profile", "loan_application", "file_upload"],
            "pages_visited": 25,
            "files_uploaded": 3
        }
    }
}
```

---

## ⚙️ 底层机制说明 (供后端参考)

-   **Redis键结构 (示例)**:
    -   会话详情: `session:{session_id}` (HASH)
        -   `user_id`, `platform`, `device_json`, `ip`, `login_at`, `last_active_at`, `access_token_hash`, `refresh_token_hash`
    -   用户所有会话ID列表: `user_sessions:{user_id_type}:{user_id}` (SET of session_ids)
        -   `user_id_type`可以是 `app` (对应 `User` 模型ID) 或 `oa` (对应 `OAUser` 模型ID)
    -   Access Token 到 Session ID 映射: `token_access:{access_token_hash}` (STRING, value: session_id)
    -   Refresh Token 到 Session ID 映射: `token_refresh:{refresh_token_hash}` (STRING, value: session_id)
    -   会话统计缓存: `session_stats:daily:{date}` (HASH)
        -   包含每日的会话统计数据
    -   活跃用户计数: `active_users:{platform}:{date}` (SET)
        -   存储每日活跃用户ID集合

-   **Token哈希**: 存储在Redis中的Token均为哈希值 (如SHA256)，不存储明文Token。

-   **会话清理**: 定期任务清理Redis中过期的会话数据，管理员也可手动触发清理。

-   **安全监控**: 记录异常登录行为，如异地登录、设备变更等。

-   **性能优化**: 
    -   使用Redis管道批量操作
    -   会话统计数据缓存
    -   分页查询优化

**此文档旨在提供清晰的API使用说明，帮助前端工程师快速接入。** 