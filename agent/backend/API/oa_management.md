# OA后台管理模块 - API 接口文档

## 📋 模块概述

OA后台管理模块为管理员提供完整的系统管理功能，包括用户管理、业务审批、数据统计、系统配置等。支持多级权限管理，实现精细化的权限控制和业务流程管理。

### 🚀 当前实现状态

#### ✅ **已实现的接口：**
- 🔐 **管理员登录**: `POST /api/oa/auth/login` 
- 🖼️ **获取验证码**: `GET /api/oa/auth/captcha`
- 👥 **用户管理**: 
  - `GET /api/oa/users` - 获取用户列表
  - `GET /api/oa/users/{user_id}` - 获取用户详情
  - `PUT /api/oa/users/{user_id}/status` - 更新用户状态
  - `POST /api/oa/users/batch-operation` - 批量操作用户
- 🚜 **农机设备管理**:
  - `GET /api/oa/machines` - 获取农机设备列表
  - `GET /api/oa/machines/{machine_id}` - 获取农机设备详情
- 📊 **数据统计**:
  - `GET /api/oa/dashboard` - 获取工作台数据
  - `GET /api/oa/dashboard/overview` - 获取业务概览
  - `GET /api/oa/dashboard/risk-monitoring` - 获取风险监控数据

#### ⚠️ **待实现的接口（当前返回模拟数据）：**
- 📋 认证审核管理
- ⚙️ 系统配置管理
- 🔍 操作日志管理
- 👑 权限管理
- 📈 报表导出

### 核心功能
- 👥 **用户管理**: 用户列表、权限管理、认证审核
- 📋 **审批管理**: 贷款审批、农机审批、工作流管理
- 📊 **数据统计**: 业务报表、用户分析、风险监控
- ⚙️ **系统配置**: 参数设置、通知管理、日志查看
- 🔐 **权限管理**: 角色管理、权限分配、操作审计

---

## 🔐 管理员认证

### 1.1 管理员登录
```http
POST /api/oa/auth/login
Content-Type: application/json

{
    "username": "admin",
    "password": "password123",
    "captcha": "ABCD",
    "captcha_id": "cap_123456"
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "登录成功",
    "data": {
        "admin": {
            "id": 1001,
            "username": "admin",
            "real_name": "系统管理员",
            "role": "super_admin",
            "department": "信息技术部",
            "permissions": ["user_manage", "loan_approve", "system_config"],
            "last_login_time": "2024-01-15T09:30:00Z"
        },
        "session": {
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "expires_in": 28800
        }
    }
}
```

### 1.2 获取验证码
```http
GET /api/oa/auth/captcha
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "captcha_id": "cap_123456",
        "captcha_image": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
    }
}
```

---

## 👥 用户管理

### 2.1 获取用户列表
```http
GET /api/oa/users?page=1&limit=20&status=active&user_type=farmer&keyword=张三
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total": 1250,
        "page": 1,
        "limit": 20,
        "filters": {
            "user_types": [
                {"value": "farmer", "label": "个体农户", "count": 800},
                {"value": "farm_owner", "label": "农场主", "count": 300},
                {"value": "cooperative", "label": "合作社", "count": 150}
            ],
            "statuses": [
                {"value": "active", "label": "正常", "count": 1100},
                {"value": "frozen", "label": "冻结", "count": 50},
                {"value": "deleted", "label": "删除", "count": 100}
            ]
        },
        "users": [
            {
                "id": 10001,
                "uuid": "uuid-abc-123",
                "phone": "13800138000",
                "real_name": "张三",
                "user_type": "farmer",
                "user_type_text": "个体农户",
                "status": "active",
                "status_text": "正常",
                "province": "山东省",
                "city": "济南市",
                "county": "历城区",
                "is_real_name_verified": true,
                "is_bank_card_verified": true,
                "credit_score": 750,
                "credit_level": "优秀",
                "total_loans": 3,
                "total_borrowed": 150000,
                "current_debt": 50000,
                "overdue_count": 0,
                "last_login_time": "2024-01-15T14:30:00Z",
                "created_at": "2023-08-15T10:00:00Z"
            }
        ]
    }
}
```

### 2.2 获取用户详情
```http
GET /api/oa/users/{user_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "basic_info": {
            "id": 10001,
            "uuid": "uuid-abc-123",
            "phone": "13800138000",
            "email": "zhangsan@example.com",
            "real_name": "张三",
            "id_card": "370123199001011234",
            "gender": "male",
            "birthday": "1990-01-01",
            "user_type": "farmer",
            "status": "active",
            "created_at": "2023-08-15T10:00:00Z"
        },
        "address_info": {
            "province": "山东省",
            "city": "济南市",
            "county": "历城区",
            "address": "某某村123号",
            "longitude": 117.1234,
            "latitude": 36.5678
        },
        "auth_status": {
            "is_real_name_verified": true,
            "is_bank_card_verified": true,
            "is_credit_verified": false,
            "real_name_auth_time": "2023-08-20T10:00:00Z",
            "bank_card_auth_time": "2023-08-25T15:30:00Z"
        },
        "credit_info": {
            "credit_score": 750,
            "credit_level": "优秀",
            "balance": 5000,
            "total_limit": 500000,
            "used_limit": 50000,
            "available_limit": 450000
        },
        "business_summary": {
            "total_loans": 3,
            "total_borrowed": 150000,
            "total_repaid": 100000,
            "current_debt": 50000,
            "overdue_count": 0,
            "total_rentals": 15,
            "total_rental_cost": 45000
        },
        "login_info": {
            "login_count": 156,
            "last_login_time": "2024-01-15T14:30:00Z",
            "last_login_ip": "192.168.1.100",
            "device_count": 2
        },
        "risk_assessment": {
            "risk_level": "low",
            "risk_score": 85,
            "risk_factors": [],
            "blacklist_status": false
        }
    }
}
```

### 2.3 用户状态管理
```http
PUT /api/oa/users/{user_id}/status
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "status": "frozen",
    "reason": "风险账户",
    "note": "多次逾期还款，暂时冻结账户",
    "notify_user": true
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "状态更新成功",
    "data": {
        "user_id": 10001,
        "old_status": "active",
        "new_status": "frozen",
        "operation_time": "2024-01-15T16:30:00Z",
        "operator": "admin"
    }
}
```

### 2.4 批量操作
```http
POST /api/oa/users/batch-operation
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "operation": "freeze",
    "user_ids": [10001, 10002, 10003],
    "reason": "批量风控处理",
    "notify_users": true
}
```

---

## 🚜 农机设备管理

### 5.1 获取农机设备列表
```http
GET /api/oa/machines?status=active&category=tillage&owner_type=cooperative&page=1&limit=20
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total": 156,
        "page": 1,
        "limit": 20,
        "summary": {
            "total_machines": 156,
            "active_machines": 135,
            "rented_machines": 45,
            "maintenance_machines": 12
        },
        "filters": {
            "categories": [
                {"code": "tillage", "name": "耕地机械", "count": 68},
                {"code": "planting", "name": "播种机械", "count": 35},
                {"code": "harvesting", "name": "收获机械", "count": 53}
            ],
            "status_options": [
                {"value": "active", "label": "可用", "count": 135},
                {"value": "rented", "label": "租赁中", "count": 45},
                {"value": "maintenance", "label": "维护中", "count": 12}
            ]
        },
        "machines": [
            {
                "id": 10001,
                "name": "约翰迪尔 6B-1204拖拉机",
                "category": "tillage",
                "category_name": "耕地机械",
                "brand": "约翰迪尔",
                "model": "6B-1204",
                "serial_number": "JD2022120401",
                "year": 2022,
                "status": "active",
                "status_text": "可用",
                "condition": "excellent",
                "owner": {
                    "id": 2001,
                    "name": "济南农机合作社",
                    "type": "cooperative",
                    "contact_phone": "0531-12345678"
                },
                "location": {
                    "province": "山东省",
                    "city": "济南市",
                    "district": "历城区",
                    "address": "农机服务站"
                },
                "rental_info": {
                    "daily_rate": 500,
                    "total_orders": 25,
                    "total_revenue": 37500,
                    "utilization_rate": 0.68
                },
                "maintenance": {
                    "last_maintenance": "2024-01-10",
                    "next_maintenance": "2024-04-10",
                    "maintenance_status": "正常"
                },
                "created_at": "2023-03-15T08:00:00Z",
                "updated_at": "2024-01-15T10:30:00Z"
            }
        ]
    }
}
```

### 5.2 获取农机设备详情
```http
GET /api/oa/machines/{machine_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "basic_info": {
            "id": 10001,
            "name": "约翰迪尔 6B-1204拖拉机",
            "category": "tillage",
            "brand": "约翰迪尔",
            "model": "6B-1204",
            "serial_number": "JD2022120401",
            "purchase_date": "2022-03-15",
            "purchase_price": 320000,
            "current_value": 280000,
            "depreciation_rate": 0.125,
            "status": "active",
            "condition": "excellent"
        },
        "owner_info": {
            "id": 2001,
            "name": "济南农机合作社",
            "type": "cooperative",
            "legal_person": "李四",
            "contact_phone": "0531-12345678",
            "business_license": "91370100123456789X",
            "registration_date": "2020-05-20"
        },
        "specifications": {
            "engine_power": "120马力",
            "engine_type": "柴油发动机",
            "transmission": "动力换挡",
            "fuel_capacity": "280L",
            "weight": "4800kg",
            "dimensions": {
                "length": "4.2m",
                "width": "2.1m", 
                "height": "2.8m"
            }
        },
        "rental_statistics": {
            "total_orders": 25,
            "total_days": 180,
            "total_revenue": 90000,
            "average_daily_rate": 500,
            "utilization_rate": 0.68,
            "customer_rating": 4.8,
            "return_rate": 0.96
        },
        "maintenance_records": [
            {
                "id": 1001,
                "type": "routine",
                "description": "例行保养检查",
                "date": "2024-01-10",
                "cost": 800,
                "technician": "张师傅",
                "status": "completed"
            }
        ],
        "current_status": {
            "location": {
                "longitude": 117.1234,
                "latitude": 36.5678,
                "address": "济南市历城区农机服务站",
                "updated_at": "2024-01-15T10:30:00Z"
            },
            "current_order": null,
            "next_booking": {
                "start_date": "2024-01-20",
                "customer_name": "张三"
            }
        }
    }
}
```

### 5.3 审批农机申请
```http
POST /api/oa/machines/{machine_id}/review
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "action": "approve",
    "review_note": "设备状态良好，服务商资质齐全",
    "approved_categories": ["耕地作业", "播种作业"],
    "restricted_areas": [],
    "rental_rate_approved": true
}
```

### 5.4 设备状态管理
```http
PUT /api/oa/machines/{machine_id}/status
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "status": "maintenance",
    "reason": "例行维护保养",
    "estimated_duration": 5,
    "maintenance_type": "routine",
    "notify_pending_customers": true
}
```

### 5.5 获取设备所有者列表
```http
GET /api/oa/machine-owners?type=cooperative&status=active&page=1&limit=20
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total": 45,
        "page": 1,
        "limit": 20,
        "owners": [
            {
                "id": 2001,
                "name": "济南农机合作社",
                "type": "cooperative",
                "legal_person": "李四",
                "contact_phone": "0531-12345678",
                "business_license": "91370100123456789X",
                "status": "active",
                "machine_count": 12,
                "total_revenue": 450000,
                "rating": 4.8,
                "registration_date": "2020-05-20",
                "last_active": "2024-01-15T14:30:00Z"
            }
        ]
    }
}
```

### 5.6 设备所有者审核
```http
POST /api/oa/machine-owners/{owner_id}/review
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "action": "approve",
    "review_note": "资质审核通过，允许上架设备",
    "approved_categories": ["耕地机械", "播种机械"],
    "max_machine_count": 50,
    "commission_rate": 0.05
}
```

### 5.7 获取设备维护记录
```http
GET /api/oa/machines/{machine_id}/maintenance?page=1&limit=20
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 5.8 创建维护记录
```http
POST /api/oa/machines/{machine_id}/maintenance
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "type": "repair",
    "description": "液压系统故障维修",
    "estimated_cost": 1500,
    "estimated_duration": 3,
    "technician": "王师傅",
    "parts_needed": ["液压油缸", "密封件"],
    "priority": "high"
}
```

### 5.9 设备利用率分析
```http
GET /api/oa/machines/analytics/utilization?period=month&start_date=2024-01-01&end_date=2024-01-31
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "overall_utilization": {
            "average_rate": 0.65,
            "total_machines": 156,
            "active_machines": 135,
            "high_utilization_count": 45,
            "low_utilization_count": 23
        },
        "category_analysis": [
            {
                "category": "tillage",
                "name": "耕地机械",
                "machine_count": 68,
                "utilization_rate": 0.72,
                "revenue": 280000
            }
        ],
        "top_performers": [
            {
                "machine_id": 10001,
                "name": "约翰迪尔 6B-1204拖拉机",
                "utilization_rate": 0.89,
                "revenue": 45000
            }
        ],
        "maintenance_alerts": [
            {
                "machine_id": 10015,
                "alert_type": "overdue_maintenance",
                "message": "设备超期未保养"
            }
        ]
    }
}
```

### 5.10 批量操作
```http
POST /api/oa/machines/batch-operation
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "operation": "status_change",
    "machine_ids": [10001, 10002, 10003],
    "new_status": "maintenance",
    "reason": "冬季集中保养",
    "estimated_duration": 7
}
```

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