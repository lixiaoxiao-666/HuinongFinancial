# 农机租赁模块 - API 接口文档

## 📋 模块概述

农机租赁模块为惠农用户提供农机搜索、预订、订单管理服务，并为OA管理员提供农机信息管理、订单审核等功能。

### 核心功能
-   **用户端 (`/api/user/machines/*`, `/api/user/orders/*`)**: 农机浏览、在线预订、订单支付、状态跟踪。
-   **OA管理员端 (`/api/oa/admin/machines/*`)**: 农机信息管理、租赁订单审核、设备状态更新。

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