# 农机租赁模块 - API 接口文档

## 📋 模块概述

农机租赁模块提供全流程的农机设备租赁服务，支持农机注册、设备搜索、订单管理、支付确认、使用跟踪和评价反馈。为农户提供便捷的农机设备共享服务，提高农业生产效率。

### 🎯 核心功能
- **农机管理**: 设备注册、信息维护、状态管理
- **智能搜索**: 多维度筛选、地理位置匹配、智能推荐
- **订单流程**: 创建订单、支付确认、使用跟踪、完成归还
- **评价系统**: 服务评价、信誉管理、质量反馈
- **后台管理**: 设备审核、订单监控、数据统计

### 🏗️ 业务流程
```
设备注册 → 信息审核 → 搜索发现 → 创建订单 → 支付确认 → 设备使用 → 完成归还 → 评价反馈
    ↓          ↓          ↓          ↓          ↓          ↓          ↓          ↓
  提交资料   管理审核   用户选择   确认租赁   在线支付   实际作业   检查归还   双向评价
```

### 📊 数据模型关系
```
Machines (农机设备)
├── MachineTypes (设备类型)
├── MachineSpecs (技术参数)
└── MachineImages (设备图片)

MachineOrders (租赁订单)
├── OrderItems (订单明细)
├── OrderPayments (支付记录)
├── OrderTracking (使用跟踪)
└── OrderRatings (评价反馈)

Users (用户)
├── MachineOwners (设备所有者)
├── MachineRenters (租赁用户)
└── UserRatings (用户信誉)
```

---

## 🚜 惠农APP/Web - 农机信息与搜索

**接口路径前缀**: `/api/user/machines`
**认证要求**: `RequireAuth` (惠农APP/Web用户，部分搜索功能可能允许匿名访问，由后端实现决定)
**适用平台**: `app`, `web`

### 1.1 搜索农机设备

```http
GET /api/user/machines/search?category=tillage&location_city=济南市&available_date=2024-03-01
Authorization: Bearer {access_token} // 或无认证，取决于实现
```

### 1.2 获取农机详情

```http
GET /api/user/machines/{machine_id}
Authorization: Bearer {access_token} // 或无认证
```

### 1.3 用户注册自己的农机 (如果支持C2C模式)

```http
POST /api/user/machines
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "name": "我的拖拉机",
    "category": "tillage",
    // ... 其他农机信息
}
```

### 1.4 获取用户自己注册的农机列表

```http
GET /api/user/machines
Authorization: Bearer {access_token}
```

---

## 🧾 惠农APP/Web - 农机订单管理

**接口路径前缀**: `/api/user/orders`
**认证要求**: `RequireAuth` (惠农APP/Web用户)
**适用平台**: `app`, `web`

### 2.1 创建农机租赁订单

```http
POST /api/user/machines/{machine_id}/orders // 注意：此路径也可能为 /api/user/orders，取决于设计
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "machine_id": 501,
    "start_date": "2024-03-01",
    "end_date": "2024-03-05",
    "rental_days": 5
    // ... 其他订单信息
}
```

### 2.2 获取用户的农机订单列表

```http
GET /api/user/orders?status=pending&page=1&limit=10
Authorization: Bearer {access_token}
```

### 2.3 获取用户指定农机订单详情

```http
GET /api/user/orders/{order_id}
Authorization: Bearer {access_token}
```

### 2.4 用户确认订单

```http
PUT /api/user/orders/{order_id}/confirm
Authorization: Bearer {access_token}
```

### 2.5 用户支付订单

```http
POST /api/user/orders/{order_id}/pay
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "payment_method": "wechat_pay",
    "amount": 2500
}
```

### 2.6 用户完成订单 (确认使用完毕)

```http
PUT /api/user/orders/{order_id}/complete
Authorization: Bearer {access_token}
```

### 2.7 用户取消订单

```http
PUT /api/user/orders/{order_id}/cancel
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "cancel_reason": "行程变更"
}
```

### 2.8 用户评价订单

```http
POST /api/user/orders/{order_id}/rate
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "rating": 5, // 1-5 星
    "comment": "非常好用的拖拉机，下次还租！"
}
```

---

## 🛠️ OA系统 - 农机管理接口 (管理员)

**接口路径前缀**: `/api/oa/admin/machines`
**认证要求**: `RequireAuth`, `CheckPlatform("oa")`, `RequireRole("admin")`
**适用平台**: `oa`

### 3.1 获取所有农机设备列表 (管理员视图)

```http
GET /api/oa/admin/machines?status=active&owner_id=201&page=1&limit=20
Authorization: Bearer {oa_access_token}
```

### 3.2 获取指定农机设备详情 (管理员视图)

```http
GET /api/oa/admin/machines/{machine_id}
Authorization: Bearer {oa_access_token}
```

### 3.3 添加新的农机设备 (由平台或合作方录入)

```http
POST /api/oa/admin/machines
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "name": "大型联合收割机",
    "serial_number": "XYZ12345",
    "owner_id": 305, // 设备所有者ID (可能是合作社或个人)
    // ... 其他详细信息
}
```

### 3.4 更新农机设备信息

```http
PUT /api/oa/admin/machines/{machine_id}
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "status": "maintenance", // active, inactive, maintenance, rented
    "rental_price_per_day": 600
}
```

### 3.5 删除农机设备

```http
DELETE /api/oa/admin/machines/{machine_id}
Authorization: Bearer {oa_access_token}
```

### 3.6 获取所有农机租赁订单列表 (管理员视图)

```http
GET /api/oa/admin/machines/orders?status=pending_approval&user_id=101&page=1&limit=20
Authorization: Bearer {oa_access_token}
```

### 3.7 审核农机租赁订单

```http
POST /api/oa/admin/machines/orders/{order_id}/approve // 或 /reject, /return
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "approval_comments": "同意租赁请求。"
}
```

---

**说明**: 路径和参数请根据实际后端实现进行调整。本文档提供核心接口的参考。