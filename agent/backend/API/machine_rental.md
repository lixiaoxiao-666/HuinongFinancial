# 农机租赁模块 - API 接口文档

## 📋 模块概述

农机租赁模块为农户提供便捷的农机设备租赁服务，支持在线查找、预约、租赁、支付、归还等全流程服务。通过智能匹配算法，为用户推荐最适合的农机设备和服务商。

### 核心功能
- 🚜 **农机查询**: 多维度搜索和筛选农机设备
- 📅 **预约管理**: 在线预约、时间管理、冲突检测
- 💰 **租赁服务**: 订单管理、支付、保险、配送
- 📍 **位置服务**: GPS定位、距离计算、路径规划
- 📊 **数据统计**: 租赁记录、费用统计、设备利用率

---

## 🚜 农机设备管理

### 1.1 搜索农机设备
```http
GET /api/machines/search?keyword=拖拉机&category=耕地机械&location=济南市&available_date=2024-01-20&sort=distance
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**查询参数:**
- `keyword`: 关键词搜索
- `category`: 设备分类
- `location`: 地理位置
- `available_date`: 可用日期
- `min_power`: 最小功率
- `max_power`: 最大功率
- `min_price`: 最低租金
- `max_price`: 最高租金
- `sort`: 排序方式 (distance/price/rating)

**响应示例:**
```json
{
    "code": 200,
    "message": "搜索成功",
    "data": {
        "total": 45,
        "page": 1,
        "limit": 20,
        "filters": {
            "categories": [
                {"code": "tillage", "name": "耕地机械", "count": 25},
                {"code": "planting", "name": "播种机械", "count": 12}
            ],
            "power_range": {"min": 50, "max": 300},
            "price_range": {"min": 200, "max": 800}
        },
        "machines": [
            {
                "id": 10001,
                "name": "约翰迪尔 6B-1204拖拉机",
                "category": "耕地机械",
                "sub_category": "四轮拖拉机",
                "brand": "约翰迪尔",
                "model": "6B-1204",
                "power": 120,
                "year": 2022,
                "condition": "excellent",
                "images": [
                    "https://example.com/machines/10001_1.jpg"
                ],
                "specifications": {
                    "engine_power": "120马力",
                    "transmission": "动力换挡",
                    "fuel_capacity": "280L",
                    "weight": "4800kg"
                },
                "rental_info": {
                    "daily_rate": 500,
                    "weekly_rate": 3000,
                    "monthly_rate": 10000,
                    "deposit": 5000,
                    "min_rental_hours": 8
                },
                "availability": {
                    "status": "available",
                    "next_available_date": "2024-01-20",
                    "busy_dates": ["2024-01-25", "2024-01-26"]
                },
                "owner": {
                    "id": 2001,
                    "name": "农机合作社",
                    "rating": 4.8,
                    "total_orders": 156
                },
                "location": {
                    "province": "山东省",
                    "city": "济南市",
                    "district": "历城区",
                    "address": "农机服务站",
                    "longitude": 117.1234,
                    "latitude": 36.5678,
                    "distance": 12.5
                },
                "services": [
                    {
                        "type": "delivery",
                        "name": "设备配送",
                        "fee": 100,
                        "range": 50
                    },
                    {
                        "type": "operator",
                        "name": "操作员服务",
                        "fee": 300,
                        "description": "提供专业操作员"
                    }
                ]
            }
        ]
    }
}
```

### 1.2 获取农机详情
```http
GET /api/machines/{machine_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "id": 10001,
        "name": "约翰迪尔 6B-1204拖拉机",
        "description": "全新2022年四轮拖拉机，适用于各种农田作业",
        "category": "耕地机械",
        "sub_category": "四轮拖拉机",
        "brand": "约翰迪尔",
        "model": "6B-1204",
        "serial_number": "JD2022120401",
        "purchase_date": "2022-03-15",
        "condition": "excellent",
        "maintenance_status": "正常",
        "last_maintenance": "2024-01-10",
        "next_maintenance": "2024-04-10",
        "images": [
            {
                "url": "https://example.com/machines/10001_1.jpg",
                "type": "main",
                "description": "主视图"
            }
        ],
        "specifications": {
            "engine_power": "120马力",
            "engine_type": "柴油发动机",
            "transmission": "动力换挡",
            "fuel_capacity": "280L",
            "hydraulic_system": "双路液压",
            "pto_speed": "540/1000 rpm",
            "weight": "4800kg",
            "dimensions": {
                "length": "4.2m",
                "width": "2.1m",
                "height": "2.8m"
            }
        },
        "rental_info": {
            "daily_rate": 500,
            "weekly_rate": 3000,
            "monthly_rate": 10000,
            "deposit": 5000,
            "min_rental_hours": 8,
            "overtime_rate": 80,
            "fuel_policy": "租户自付",
            "insurance_included": true
        },
        "availability_calendar": [
            {
                "date": "2024-01-20",
                "status": "available",
                "time_slots": [
                    {"start": "08:00", "end": "12:00", "status": "available"},
                    {"start": "13:00", "end": "17:00", "status": "booked"}
                ]
            }
        ],
        "reviews": {
            "average_rating": 4.8,
            "total_reviews": 24,
            "rating_distribution": {
                "5": 18,
                "4": 4,
                "3": 2,
                "2": 0,
                "1": 0
            },
            "recent_reviews": [
                {
                    "user_name": "张***",
                    "rating": 5,
                    "comment": "设备状态很好，操作简单，效率高",
                    "date": "2024-01-10"
                }
            ]
        }
    }
}
```

### 1.3 获取设备分类
```http
GET /api/machines/categories
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": [
        {
            "code": "tillage",
            "name": "耕地机械",
            "icon": "https://example.com/icons/tillage.png",
            "subcategories": [
                {"code": "tractor", "name": "拖拉机"},
                {"code": "plow", "name": "犁具"},
                {"code": "harrow", "name": "耙具"}
            ]
        },
        {
            "code": "planting",
            "name": "播种机械",
            "icon": "https://example.com/icons/planting.png",
            "subcategories": [
                {"code": "seeder", "name": "播种机"},
                {"code": "transplanter", "name": "插秧机"}
            ]
        },
        {
            "code": "harvesting",
            "name": "收获机械",
            "icon": "https://example.com/icons/harvesting.png",
            "subcategories": [
                {"code": "combine", "name": "联合收割机"},
                {"code": "thresher", "name": "脱粒机"}
            ]
        }
    ]
}
```

---

## 📅 预约管理

### 2.1 创建预约
```http
POST /api/machines/reservations
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "machine_id": 10001,
    "start_date": "2024-01-20",
    "end_date": "2024-01-22",
    "start_time": "08:00",
    "end_time": "17:00",
    "work_location": {
        "address": "济南市历城区某某村",
        "longitude": 117.2345,
        "latitude": 36.6789,
        "work_area": 50,
        "crop_type": "水稻"
    },
    "services": [
        {
            "type": "delivery",
            "required": true
        },
        {
            "type": "operator",
            "required": false
        }
    ],
    "special_requirements": "需要早上7点开始作业",
    "emergency_contact": {
        "name": "张三",
        "phone": "13800138000"
    }
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "预约创建成功",
    "data": {
        "reservation_id": "RS20240115001",
        "status": "pending",
        "estimated_cost": {
            "rental_fee": 1500,
            "delivery_fee": 100,
            "deposit": 5000,
            "total": 6600
        },
        "confirmation_deadline": "2024-01-18T18:00:00Z",
        "next_steps": [
            "等待设备所有者确认",
            "确认后请在24小时内支付定金"
        ]
    }
}
```

### 2.2 获取预约列表
```http
GET /api/machines/reservations?status=all&page=1&limit=10
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total": 15,
        "page": 1,
        "limit": 10,
        "reservations": [
            {
                "id": "RS20240115001",
                "machine": {
                    "id": 10001,
                    "name": "约翰迪尔 6B-1204拖拉机",
                    "image": "https://example.com/machines/10001_1.jpg"
                },
                "start_date": "2024-01-20",
                "end_date": "2024-01-22",
                "duration_days": 3,
                "status": "confirmed",
                "status_text": "已确认",
                "total_cost": 6600,
                "deposit_paid": 5000,
                "created_at": "2024-01-15T10:30:00Z",
                "confirmed_at": "2024-01-16T09:15:00Z"
            }
        ]
    }
}
```

### 2.3 获取预约详情
```http
GET /api/machines/reservations/{reservation_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 2.4 修改预约
```http
PUT /api/machines/reservations/{reservation_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "start_date": "2024-01-21",
    "end_date": "2024-01-23",
    "special_requirements": "更新的特殊要求"
}
```

### 2.5 取消预约
```http
DELETE /api/machines/reservations/{reservation_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "reason": "时间冲突",
    "cancellation_type": "user_initiated"
}
```

---

## 💰 租赁订单管理

### 3.1 创建租赁订单
```http
POST /api/machines/orders
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "reservation_id": "RS20240115001",
    "payment_method": "alipay",
    "insurance_plan": "basic",
    "additional_services": [
        {
            "type": "fuel_service",
            "quantity": 200
        }
    ]
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "订单创建成功",
    "data": {
        "order_id": "MO20240116001",
        "status": "pending_payment",
        "payment_info": {
            "amount": 6600,
            "deposit": 5000,
            "payment_deadline": "2024-01-17T18:00:00Z",
            "payment_url": "https://pay.example.com/mopay/MO20240116001"
        },
        "contract_url": "https://example.com/contracts/MO20240116001.pdf"
    }
}
```

### 3.2 获取订单列表
```http
GET /api/machines/orders?status=all&page=1&limit=10
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total": 28,
        "page": 1,
        "limit": 10,
        "orders": [
            {
                "id": "MO20240116001",
                "machine": {
                    "id": 10001,
                    "name": "约翰迪尔 6B-1204拖拉机",
                    "image": "https://example.com/machines/10001_1.jpg"
                },
                "rental_period": {
                    "start_date": "2024-01-20",
                    "end_date": "2024-01-22",
                    "duration_days": 3
                },
                "status": "in_progress",
                "status_text": "使用中",
                "total_amount": 6600,
                "paid_amount": 6600,
                "created_at": "2024-01-16T14:30:00Z",
                "current_location": {
                    "address": "济南市历城区某某村",
                    "updated_at": "2024-01-20T10:30:00Z"
                }
            }
        ]
    }
}
```

### 3.3 获取订单详情
```http
GET /api/machines/orders/{order_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "id": "MO20240116001",
        "reservation_id": "RS20240115001",
        "status": "in_progress",
        "status_text": "使用中",
        "machine": {
            "id": 10001,
            "name": "约翰迪尔 6B-1204拖拉机",
            "model": "6B-1204",
            "serial_number": "JD2022120401"
        },
        "rental_period": {
            "start_date": "2024-01-20",
            "end_date": "2024-01-22",
            "start_time": "08:00",
            "end_time": "17:00",
            "duration_days": 3,
            "actual_start_time": "2024-01-20T08:00:00Z"
        },
        "cost_breakdown": {
            "rental_fee": 1500,
            "delivery_fee": 100,
            "operator_fee": 0,
            "fuel_fee": 400,
            "insurance_fee": 100,
            "deposit": 5000,
            "subtotal": 2100,
            "total": 7100
        },
        "payment_info": {
            "total_paid": 7100,
            "payment_method": "alipay",
            "payment_time": "2024-01-16T15:30:00Z",
            "transaction_id": "PAY20240116001"
        },
        "delivery_info": {
            "delivery_address": "济南市历城区某某村",
            "delivery_time": "2024-01-20T07:00:00Z",
            "pickup_time": "2024-01-22T18:00:00Z",
            "delivery_status": "delivered",
            "driver_info": {
                "name": "李师傅",
                "phone": "13900139000"
            }
        },
        "usage_tracking": {
            "current_location": {
                "longitude": 117.2345,
                "latitude": 36.6789,
                "address": "济南市历城区某某村田间",
                "updated_at": "2024-01-20T14:30:00Z"
            },
            "operating_hours": 12.5,
            "fuel_consumption": 45.5,
            "work_area_completed": 35.2
        }
    }
}
```

---

## 📍 设备跟踪管理

### 4.1 获取设备位置
```http
GET /api/machines/{machine_id}/location
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "machine_id": 10001,
        "current_location": {
            "longitude": 117.2345,
            "latitude": 36.6789,
            "address": "济南市历城区某某村田间",
            "accuracy": 5.2,
            "updated_at": "2024-01-20T14:30:00Z"
        },
        "status": "working",
        "engine_status": "running",
        "fuel_level": 65,
        "operating_data": {
            "engine_hours": 1250.5,
            "speed": 8.5,
            "working_depth": 25,
            "hydraulic_pressure": 180
        }
    }
}
```

### 4.2 上报设备状态
```http
POST /api/machines/{machine_id}/status
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "location": {
        "longitude": 117.2345,
        "latitude": 36.6789,
        "accuracy": 3.8
    },
    "engine_status": "running",
    "fuel_level": 60,
    "operating_data": {
        "speed": 8.5,
        "working_depth": 25,
        "hydraulic_pressure": 180
    },
    "work_progress": {
        "area_completed": 15.5,
        "total_area": 50
    }
}
```

### 4.3 获取轨迹记录
```http
GET /api/machines/{machine_id}/track?start_date=2024-01-20&end_date=2024-01-22
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 🔄 归还管理

### 5.1 申请归还
```http
POST /api/machines/orders/{order_id}/return
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "return_location": {
        "longitude": 117.1234,
        "latitude": 36.5678,
        "address": "农机服务站"
    },
    "return_time": "2024-01-22T17:00:00Z",
    "fuel_level": 45,
    "machine_condition": "normal",
    "work_summary": {
        "total_area": 50,
        "completed_area": 48,
        "operating_hours": 24
    },
    "issues_reported": [],
    "photos": [
        "https://example.com/return/photos/1.jpg"
    ]
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "归还申请提交成功",
    "data": {
        "return_id": "RT20240122001",
        "status": "pending_inspection",
        "inspection_appointment": "2024-01-22T18:00:00Z",
        "estimated_refund": {
            "deposit_refund": 5000,
            "unused_fuel_refund": 150,
            "total_refund": 5150
        }
    }
}
```

### 5.2 确认归还
```http
POST /api/machines/orders/{order_id}/return-confirm
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "return_id": "RT20240122001",
    "inspection_result": "passed",
    "actual_fuel_level": 43,
    "damage_assessment": [],
    "final_charges": {
        "overtime_fee": 0,
        "damage_fee": 0,
        "cleaning_fee": 0
    }
}
```

---

## 💳 支付管理

### 6.1 创建支付订单
```http
POST /api/machines/payments
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "order_id": "MO20240116001",
    "payment_type": "deposit",
    "amount": 5000,
    "payment_method": "alipay"
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "支付订单创建成功",
    "data": {
        "payment_id": "PAY20240116001",
        "payment_url": "https://pay.example.com/mopay/PAY20240116001",
        "qr_code": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
        "expires_at": "2024-01-16T18:00:00Z"
    }
}
```

### 6.2 查询支付状态
```http
GET /api/machines/payments/{payment_id}/status
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 6.3 申请退款
```http
POST /api/machines/payments/{payment_id}/refund
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "refund_amount": 5000,
    "refund_reason": "设备故障无法使用",
    "bank_account": {
        "account_number": "6226090000000001",
        "bank_name": "中国工商银行",
        "account_holder": "张三"
    }
}
```

---

## 📊 统计分析

### 7.1 获取租赁统计
```http
GET /api/machines/statistics?period=month&start_date=2024-01-01&end_date=2024-01-31
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "overview": {
            "total_rentals": 15,
            "total_spending": 45000,
            "total_hours": 120,
            "favorite_category": "耕地机械"
        },
        "monthly_trend": [
            {
                "month": "2024-01",
                "rentals": 15,
                "spending": 45000,
                "hours": 120
            }
        ],
        "category_breakdown": [
            {
                "category": "耕地机械",
                "rentals": 8,
                "spending": 28000,
                "percentage": 62.2
            },
            {
                "category": "播种机械",
                "rentals": 4,
                "spending": 12000,
                "percentage": 26.7
            }
        ],
        "cost_analysis": {
            "rental_fees": 35000,
            "delivery_fees": 3000,
            "service_fees": 7000,
            "average_daily_cost": 375
        }
    }
}
```

### 7.2 获取设备评价
```http
GET /api/machines/{machine_id}/reviews?page=1&limit=10
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 7.3 提交设备评价
```http
POST /api/machines/orders/{order_id}/review
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "rating": 5,
    "comment": "设备状态很好，操作简单，效率高",
    "tags": ["设备新", "效率高", "服务好"],
    "photos": [
        "https://example.com/reviews/photo1.jpg"
    ]
}
```

---

## 🔧 错误码说明

| 错误码 | 说明 | 处理建议 |
|-------|------|---------|
| 3001 | 农机设备不存在 | 检查设备ID是否正确 |
| 3002 | 设备不可用 | 选择其他可用时间或设备 |
| 3003 | 预约时间冲突 | 调整预约时间 |
| 3004 | 预约已过期 | 重新创建预约 |
| 3005 | 订单状态不允许操作 | 检查订单状态 |
| 3006 | 支付金额不匹配 | 检查支付金额 |
| 3007 | 设备正在使用中 | 等待设备归还 |
| 3008 | 超出服务范围 | 选择服务范围内的地址 |
| 3009 | 设备维护中 | 选择其他设备 |
| 3010 | 用户认证不足 | 完成必要的身份认证 |

---

## 📝 接口调用示例

### JavaScript示例
```javascript
// 搜索农机设备
const searchMachines = async (token, params) => {
    const queryString = new URLSearchParams(params).toString();
    const response = await fetch(`/api/machines/search?${queryString}`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};

// 创建预约
const createReservation = async (token, reservationData) => {
    const response = await fetch('/api/machines/reservations', {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(reservationData)
    });
    return response.json();
};

// 跟踪设备位置
const trackMachine = async (token, machineId) => {
    const response = await fetch(`/api/machines/${machineId}/location`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};
```

### 业务流程说明
1. **租赁流程**: 搜索设备 → 查看详情 → 创建预约 → 确认订单 → 支付定金 → 设备配送 → 使用设备 → 归还设备 → 结算费用
2. **预约管理**: 时间检查 → 冲突检测 → 自动确认 → 状态更新 → 取消处理
3. **设备跟踪**: GPS定位 → 状态监控 → 使用统计 → 异常报警
4. **归还流程**: 申请归还 → 设备检查 → 费用结算 → 押金退还

### 注意事项
1. **安全保护**: 设备定位信息仅限租赁期间可见
2. **时效性**: 预约和支付都有时效限制
3. **责任划分**: 明确设备使用期间的责任归属
4. **服务质量**: 建立完善的评价和投诉机制