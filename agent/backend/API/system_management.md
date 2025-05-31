# 系统管理模块 - API 接口文档

## 📋 模块概述

系统管理模块提供系统配置、监控、统计等功能。包含公开API和管理员专用API两部分。

### 核心功能
-   **公开API**: 系统版本信息、公共配置等
-   **系统配置**: 系统参数设置和管理
-   **系统监控**: 健康检查、性能统计
-   **数据统计**: 业务数据统计和分析

---

## 🌐 公开API (无需认证)

**接口路径前缀**: `/api/public`
**认证要求**: 无
**适用平台**: `app`, `web`, `oa`

### 1.1 获取系统版本信息

```http
GET /api/public/version
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "api_version": "v1.3.1",
        "build_time": "2024-01-15T08:00:00Z",
        "build_commit": "abc1234567890",
        "go_version": "1.21.0",
        "system_name": "数字惠农API服务",
        "environment": "production",
        "features": {
            "session_enabled": true,
            "ai_risk_assessment": true,
            "file_upload": true,
            "multi_platform": true
        },
        "supported_platforms": ["app", "web", "oa"],
        "api_endpoints": {
            "health_check": "/health",
            "swagger_docs": "/swagger/index.html"
        }
    }
}
```

### 1.2 获取公共配置

```http
GET /api/public/configs
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "app_config": {
            "app_name": "数字惠农",
            "app_version": "1.3.1",
            "logo_url": "https://static.huinong.com/logo.png",
            "support_email": "support@huinong.com",
            "support_phone": "400-1234-567"
        },
        "file_upload": {
            "max_file_size_mb": 50,
            "supported_extensions": [".jpg", ".jpeg", ".png", ".pdf", ".doc", ".docx"],
            "image_max_size_mb": 10,
            "document_max_size_mb": 20
        },
        "loan_config": {
            "min_loan_amount": 5000,
            "max_loan_amount": 1000000,
            "default_interest_rate": 0.065,
            "max_loan_term_months": 60
        },
        "machine_rental": {
            "booking_advance_days": 7,
            "max_rental_days": 30,
            "cancellation_hours": 24
        },
        "contact_info": {
            "company_name": "数字惠农科技有限公司",
            "address": "山东省济南市历城区创新大厦",
            "business_hours": "周一至周五 9:00-18:00",
            "emergency_contact": "400-1234-567"
        },
        "privacy_policy_url": "https://www.huinong.com/privacy",
        "terms_of_service_url": "https://www.huinong.com/terms",
        "last_updated": "2024-01-15T10:00:00Z"
    }
}
```

---

## ⚡ 系统健康检查

### 2.1 健康检查接口

```http
GET /health
```

**说明**: 此接口无需认证，用于负载均衡器和监控系统检查服务状态。

**响应示例 (服务正常):**
```json
{
    "status": "ok",
    "message": "数字惠农API服务正在运行",
    "version": "1.3.1",
    "timestamp": "2024-01-15T14:30:00Z",
    "session_enabled": true,
    "uptime_seconds": 86400,
    "checks": {
        "database": "ok",
        "redis": "ok", 
        "file_storage": "ok",
        "external_services": "ok"
    }
}
```

**响应示例 (服务异常):**
```json
{
    "status": "error",
    "message": "服务部分功能异常",
    "version": "1.3.1",
    "timestamp": "2024-01-15T14:30:00Z",
    "session_enabled": true,
    "uptime_seconds": 86400,
    "checks": {
        "database": "ok",
        "redis": "error",
        "file_storage": "ok", 
        "external_services": "warning"
    },
    "errors": [
        {
            "component": "redis",
            "error": "Connection timeout",
            "last_check": "2024-01-15T14:29:30Z"
        }
    ]
}
```

---

## 🛠️ OA系统 - 系统管理接口 (管理员)

**接口路径前缀**: `/api/oa/admin/system`
**认证要求**: `RequireAuth`, `CheckPlatform("oa")`, `RequireRole("admin")`
**适用平台**: `oa`

### 3.1 获取系统配置

```http
GET /api/oa/admin/system/config
Authorization: Bearer {oa_access_token}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "general": {
            "site_name": "数字惠农金融平台",
            "site_description": "专注于农业金融服务的数字化平台",
            "admin_email": "admin@huinong.com",
            "maintenance_mode": false,
            "registration_enabled": true
        },
        "loan_settings": {
            "default_interest_rate": 0.065,
            "min_loan_amount": 5000,
            "max_loan_amount": 1000000,
            "max_loan_term_months": 60,
            "auto_approval_threshold": 50000,
            "ai_assessment_enabled": true
        },
        "file_settings": {
            "max_file_size_mb": 50,
            "allowed_extensions": [".jpg", ".jpeg", ".png", ".pdf", ".doc", ".docx"],
            "storage_type": "oss", // local, oss, s3
            "auto_cleanup_days": 365
        },
        "notification_settings": {
            "sms_enabled": true,
            "email_enabled": true,
            "push_enabled": true,
            "daily_report_enabled": true
        },
        "security_settings": {
            "session_timeout_hours": 24,
            "password_min_length": 8,
            "password_require_special_chars": true,
            "max_login_attempts": 5,
            "lockout_duration_minutes": 30
        },
        "integration_settings": {
            "dify_api_enabled": true,
            "external_credit_check": true,
            "bank_integration": true
        }
    }
}
```

### 3.2 更新系统配置

```http
PUT /api/oa/admin/system/config
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "general": {
        "site_name": "新数字惠农金融平台",
        "maintenance_mode": false
    },
    "loan_settings": {
        "default_interest_rate": 0.055,
        "auto_approval_threshold": 60000
    },
    "security_settings": {
        "session_timeout_hours": 12,
        "max_login_attempts": 3
    }
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "系统配置更新成功",
    "data": {
        "updated_sections": ["general", "loan_settings", "security_settings"],
        "updated_at": "2024-01-15T16:00:00Z",
        "updated_by": {
            "id": 201,
            "username": "admin_zhang",
            "real_name": "管理员张三"
        }
    }
}
```

### 3.3 获取配置列表 (分类查看)

```http
GET /api/oa/admin/system/configs?section=loan_settings
Authorization: Bearer {oa_access_token}
```

**Query Parameters**:
-   `section` (string, 可选): 配置分类 (`general`, `loan_settings`, `file_settings`, `notification_settings`, `security_settings`, `integration_settings`)

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "sections": [
            {
                "section_name": "loan_settings",
                "section_title": "贷款设置",
                "configs": [
                    {
                        "key": "default_interest_rate",
                        "name": "默认利率",
                        "value": 0.065,
                        "type": "float",
                        "description": "新贷款产品的默认年利率",
                        "min_value": 0.01,
                        "max_value": 0.20,
                        "required": true,
                        "last_updated": "2024-01-15T16:00:00Z"
                    },
                    {
                        "key": "min_loan_amount", 
                        "name": "最小贷款金额",
                        "value": 5000,
                        "type": "integer",
                        "description": "用户可申请的最小贷款金额",
                        "min_value": 1000,
                        "max_value": 50000,
                        "required": true,
                        "last_updated": "2024-01-10T10:00:00Z"
                    }
                ]
            }
        ]
    }
}
```

### 3.4 详细健康检查 (管理员)

```http
GET /api/oa/admin/system/health
Authorization: Bearer {oa_access_token}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "overall_status": "healthy",
        "last_check": "2024-01-15T16:30:00Z",
        "uptime": {
            "seconds": 864000,
            "readable": "10天"
        },
        "components": {
            "database": {
                "status": "healthy",
                "response_time_ms": 12,
                "connections": {
                    "active": 8,
                    "idle": 12,
                    "max": 100
                },
                "last_error": null,
                "last_check": "2024-01-15T16:30:00Z"
            },
            "redis": {
                "status": "healthy",
                "response_time_ms": 2,
                "memory_usage": {
                    "used_mb": 256,
                    "max_mb": 1024,
                    "usage_percentage": 25.0
                },
                "connected_clients": 15,
                "last_check": "2024-01-15T16:30:00Z"
            },
            "file_storage": {
                "status": "healthy",
                "storage_usage": {
                    "used_gb": 45.6,
                    "total_gb": 1000,
                    "usage_percentage": 4.56
                },
                "upload_speed_mbps": 10.2,
                "download_speed_mbps": 25.8,
                "last_check": "2024-01-15T16:29:00Z"
            },
            "external_services": {
                "dify_api": {
                    "status": "healthy",
                    "response_time_ms": 156,
                    "success_rate_24h": 99.2,
                    "last_call": "2024-01-15T16:25:00Z"
                },
                "bank_api": {
                    "status": "warning",
                    "response_time_ms": 2500,
                    "success_rate_24h": 95.8,
                    "last_error": "Timeout on credit check",
                    "last_call": "2024-01-15T16:28:00Z"
                }
            }
        }
    }
}
```

### 3.5 获取系统统计

```http
GET /api/oa/admin/system/statistics?period=7d
Authorization: Bearer {oa_access_token}
```

**Query Parameters**:
-   `period` (string, 可选): 统计周期 (`1d`, `7d`, `30d`, `90d`, `1y`)

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "period": "7d",
        "generated_at": "2024-01-15T16:30:00Z",
        "api_statistics": {
            "total_requests": 125680,
            "average_requests_per_day": 17954,
            "success_rate": 99.3,
            "average_response_time_ms": 125,
            "error_rate": 0.7,
            "top_endpoints": [
                {
                    "endpoint": "/api/user/profile",
                    "calls": 8520,
                    "success_rate": 99.8
                },
                {
                    "endpoint": "/api/user/loan/applications",
                    "calls": 6240,
                    "success_rate": 98.5
                }
            ]
        },
        "user_statistics": {
            "total_users": 15420,
            "new_users_this_period": 320,
            "active_users_this_period": 8640,
            "daily_active_users_avg": 1234,
            "platform_distribution": {
                "app": 12500,
                "web": 2650,
                "oa": 270
            }
        },
        "business_statistics": {
            "loan_applications": {
                "total": 156,
                "approved": 89,
                "rejected": 45,
                "pending": 22,
                "approval_rate": 66.9
            },
            "file_uploads": {
                "total": 2350,
                "total_size_gb": 12.8,
                "average_size_mb": 5.6
            },
            "machine_rentals": {
                "total_orders": 45,
                "completed": 38,
                "cancelled": 4,
                "in_progress": 3
            }
        },
        "system_performance": {
            "cpu_usage_avg": 25.6,
            "memory_usage_avg": 62.4,
            "disk_usage_avg": 45.2,
            "network_in_mb": 1250.8,
            "network_out_mb": 2150.6
        }
    }
}
```

---

## 📊 工作台和仪表盘接口

### 4.1 获取OA工作台概览

```http
GET /api/oa/admin/dashboard
Authorization: Bearer {oa_access_token}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "overview": {
            "total_users": 15420,
            "new_users_today": 45,
            "active_sessions": 158,
            "pending_tasks": 23
        },
        "loan_overview": {
            "applications_today": 12,
            "pending_review": 8,
            "approved_today": 15,
            "total_amount_approved_today": 850000
        },
        "system_status": {
            "api_health": "healthy",
            "database_health": "healthy",
            "redis_health": "healthy",
            "storage_health": "healthy"
        },
        "recent_activities": [
            {
                "id": 1,
                "type": "loan_approved",
                "description": "贷款申请 LA20240115008 已批准",
                "operator": "审核员李四",
                "time": "2024-01-15T16:25:00Z"
            },
            {
                "id": 2,
                "type": "user_registered",
                "description": "新用户张三完成注册",
                "time": "2024-01-15T16:20:00Z"
            }
        ],
        "quick_stats": {
            "response_time_ms": 125,
            "error_rate": 0.3,
            "uptime_percentage": 99.9
        }
    }
}
```

### 4.2 获取风险监控数据

```http
GET /api/oa/admin/dashboard/risk-monitoring
Authorization: Bearer {oa_access_token}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "risk_alerts": [
            {
                "id": 1,
                "level": "medium",
                "type": "high_rejection_rate",
                "title": "贷款拒绝率异常",
                "description": "过去24小时贷款拒绝率为78%，超过正常水平",
                "created_at": "2024-01-15T14:00:00Z",
                "status": "active"
            },
            {
                "id": 2,
                "level": "low",
                "type": "api_slow_response",
                "title": "API响应时间偏高",
                "description": "征信查询接口平均响应时间为2.5秒",
                "created_at": "2024-01-15T15:30:00Z",
                "status": "active"
            }
        ],
        "metrics": {
            "fraud_detection": {
                "suspicious_applications": 3,
                "blocked_ips": 5,
                "failed_authentications": 25
            },
            "credit_risk": {
                "high_risk_applications": 8,
                "average_credit_score": 720,
                "default_rate_30d": 2.1
            },
            "operational_risk": {
                "system_errors_24h": 12,
                "slow_queries": 3,
                "storage_usage": 45.6
            }
        },
        "trend_analysis": {
            "application_volume": {
                "current_week": 156,
                "previous_week": 142,
                "change_percentage": 9.9
            },
            "approval_rate": {
                "current_week": 66.9,
                "previous_week": 72.1,
                "change_percentage": -7.2
            }
        }
    }
}
```

---

## 🔧 错误码说明

| 错误码 | 说明 | 处理建议 |
|-------|------|---------|
| 7001 | 配置项不存在 | 检查配置键名是否正确 |
| 7002 | 配置值格式无效 | 检查配置值类型和格式 |
| 7003 | 配置值超出范围 | 使用有效范围内的值 |
| 7004 | 系统维护模式中 | 等待维护完成后重试 |
| 7005 | 系统服务不可用 | 联系系统管理员 |
| 7006 | 统计数据计算中 | 稍后重试获取统计数据 |
| 7007 | 健康检查失败 | 检查系统组件状态 |
| 7008 | 配置更新权限不足 | 确认管理员权限 |
| 7009 | 配置锁定中 | 等待其他管理员完成配置 |
| 7010 | 备份创建失败 | 检查存储空间和权限 |

---

## 📝 接口调用示例

### JavaScript示例
```javascript
// 获取系统版本信息
const getSystemVersion = async () => {
    const response = await fetch('/api/public/version');
    return response.json();
};

// 获取公共配置
const getPublicConfigs = async () => {
    const response = await fetch('/api/public/configs');
    return response.json();
};

// 健康检查
const healthCheck = async () => {
    const response = await fetch('/health');
    return response.json();
};

// 获取系统配置 (管理员)
const getSystemConfig = async (token) => {
    const response = await fetch('/api/oa/admin/system/config', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};

// 更新系统配置 (管理员)
const updateSystemConfig = async (token, configData) => {
    const response = await fetch('/api/oa/admin/system/config', {
        method: 'PUT',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(configData)
    });
    return response.json();
};

// 获取系统统计 (管理员)
const getSystemStatistics = async (token, period = '7d') => {
    const response = await fetch(`/api/oa/admin/system/statistics?period=${period}`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};

// 获取工作台数据 (管理员)
const getDashboard = async (token) => {
    const response = await fetch('/api/oa/admin/dashboard', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};
```

### 配置管理示例
```javascript
// 配置更新工具类
class SystemConfigManager {
    constructor(token) {
        this.token = token;
        this.baseUrl = '/api/oa/admin/system';
    }
    
    // 获取指定分类的配置
    async getConfigSection(section) {
        const response = await fetch(`${this.baseUrl}/configs?section=${section}`, {
            headers: { 'Authorization': `Bearer ${this.token}` }
        });
        return response.json();
    }
    
    // 批量更新配置
    async updateConfigs(updates) {
        const response = await fetch(`${this.baseUrl}/config`, {
            method: 'PUT',
            headers: {
                'Authorization': `Bearer ${this.token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(updates)
        });
        return response.json();
    }
    
    // 验证配置值
    validateConfig(key, value, configSchema) {
        const schema = configSchema[key];
        if (!schema) return { valid: false, error: 'Unknown config key' };
        
        if (schema.required && !value) {
            return { valid: false, error: 'Value is required' };
        }
        
        if (schema.type === 'integer' && !Number.isInteger(value)) {
            return { valid: false, error: 'Value must be an integer' };
        }
        
        if (schema.min_value && value < schema.min_value) {
            return { valid: false, error: `Value must be >= ${schema.min_value}` };
        }
        
        if (schema.max_value && value > schema.max_value) {
            return { valid: false, error: `Value must be <= ${schema.max_value}` };
        }
        
        return { valid: true };
    }
}
```

### 注意事项
1. **权限控制**: 系统配置修改需要最高权限
2. **配置备份**: 重要配置变更前自动创建备份
3. **影响评估**: 配置修改可能影响系统行为，需谨慎操作
4. **版本管理**: 配置变更需要记录版本和操作日志
5. **热更新**: 部分配置支持热更新，无需重启服务
6. **监控告警**: 关键配置异常时及时告警 