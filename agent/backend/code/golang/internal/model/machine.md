# 农机租赁数据模型

## 文件概述

`machine.go` 是数字惠农系统农机租赁业务的核心数据模型文件，定义了农机设备管理和租赁订单处理相关的所有数据结构。

## 核心数据模型

### 1. Machine 农机设备表
农机设备信息管理的核心模型，支持多种农机类型和灵活的定价策略。

```go
type Machine struct {
    ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    MachineCode string `gorm:"type:varchar(30);uniqueIndex;not null" json:"machine_code"`
    MachineName string `gorm:"type:varchar(100);not null" json:"machine_name"`
    
    // 设备分类：tractor(拖拉机)、harvester(收割机)、planter(播种机)、sprayer(喷药机)
    MachineType string `gorm:"type:varchar(30);not null" json:"machine_type"`
    
    // 品牌和型号
    Brand string `gorm:"type:varchar(50);not null" json:"brand"`
    Model string `gorm:"type:varchar(50);not null" json:"model"`
    
    // 规格参数(JSON格式)
    Specifications string `gorm:"type:json" json:"specifications"`
    
    // 位置信息
    Province  string  `json:"province"`  // 省份
    City      string  `json:"city"`      // 城市
    County    string  `json:"county"`    // 县区
    Longitude float64 `json:"longitude"` // 经度
    Latitude  float64 `json:"latitude"`  // 纬度
    
    // 租赁定价
    HourlyRate  int64 `json:"hourly_rate"`   // 小时租金(分)
    DailyRate   int64 `json:"daily_rate"`    // 日租金(分)
    PerAcreRate int64 `json:"per_acre_rate"` // 按亩收费(分)
    
    // 设备状态：available(可租)、rented(已租出)、maintenance(维护中)、offline(下线)
    Status string `json:"status"`
}
```

**设备分类 (MachineType)**:
- `tractor`: 拖拉机 - 农田作业的主要动力机械
- `harvester`: 收割机 - 作物收获专用设备
- `planter`: 播种机 - 种子播种专用设备
- `sprayer`: 喷药机 - 农药喷洒设备
- `cultivator`: 耕作机 - 土地耕作设备
- `fertilizer_spreader`: 施肥机 - 肥料撒播设备

**设备状态 (Status)**:
- `available`: 可租 - 设备正常，可以出租
- `rented`: 已租出 - 设备正在租赁中
- `maintenance`: 维护中 - 设备正在维修保养
- `offline`: 下线 - 设备暂时不可用

### 2. RentalOrder 租赁订单表
农机租赁订单的完整信息管理，支持多种计费方式和订单状态跟踪。

```go
type RentalOrder struct {
    ID      uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    OrderNo string `gorm:"type:varchar(30);uniqueIndex;not null" json:"order_no"`
    
    // 关联信息
    MachineID uint64 `gorm:"not null;index" json:"machine_id"`
    RenterID  uint64 `gorm:"not null;index" json:"renter_id"`
    OwnerID   uint64 `gorm:"not null;index" json:"owner_id"`
    
    // 租赁时间
    StartTime      time.Time `json:"start_time"`      // 开始时间
    EndTime        time.Time `json:"end_time"`        // 结束时间
    RentalDuration int       `json:"rental_duration"` // 租赁时长(小时)
    
    // 租赁地点
    RentalLocation string `json:"rental_location"` // 租赁地点
    ContactPerson  string `json:"contact_person"`  // 联系人
    ContactPhone   string `json:"contact_phone"`   // 联系电话
    
    // 计费方式：hourly(按小时)、daily(按天)、per_acre(按亩)
    BillingMethod string `json:"billing_method"`
    
    // 费用计算
    UnitPrice      int64   `json:"unit_price"`      // 单价(分)
    Quantity       float64 `json:"quantity"`        // 数量(小时/天/亩)
    SubtotalAmount int64   `json:"subtotal_amount"` // 小计(分)
    DepositAmount  int64   `json:"deposit_amount"`  // 押金(分)
    TotalAmount    int64   `json:"total_amount"`    // 总金额(分)
    
    // 订单状态：pending(待确认)、confirmed(已确认)、paid(已支付)、in_progress(进行中)、
    // completed(已完成)、cancelled(已取消)、disputed(有争议)
    Status string `json:"status"`
    
    // 服务评价
    RenterRating  float32 `json:"renter_rating"`  // 租客评分
    RenterComment string  `json:"renter_comment"` // 租客评价
    OwnerRating   float32 `json:"owner_rating"`   // 机主评分
    OwnerComment  string  `json:"owner_comment"`  // 机主评价
}
```

**计费方式 (BillingMethod)**:
- `hourly`: 按小时计费 - 适用于短期精确计费
- `daily`: 按天计费 - 适用于整天作业的情况
- `per_acre`: 按亩计费 - 适用于按作业面积计费

**订单状态流转**:
```
pending(待确认) → 
confirmed(已确认) → 
paid(已支付) → 
in_progress(进行中) → 
completed(已完成) / cancelled(已取消) / disputed(有争议)
```

## 辅助数据结构

### MachineSpecs 设备规格参数
```go
type MachineSpecs struct {
    Power        string            `json:"power"`         // 功率 (如: "120马力")
    WorkingWidth string            `json:"working_width"` // 工作幅宽 (如: "2.5米")
    Weight       string            `json:"weight"`        // 重量 (如: "3500公斤")
    FuelType     string            `json:"fuel_type"`     // 燃料类型 (柴油/汽油/电动)
    YearOfMake   int               `json:"year_of_make"`  // 制造年份
    WorkingSpeed string            `json:"working_speed"` // 工作速度 (如: "5-12公里/小时")
    Capacity     string            `json:"capacity"`      // 容量 (如: "3立方米")
    Other        map[string]string `json:"other"`         // 其他参数
}
```

### AvailableSchedule 可用时间表
```go
type AvailableSchedule struct {
    Monday    []TimeSlot `json:"monday"`    // 周一可用时间
    Tuesday   []TimeSlot `json:"tuesday"`   // 周二可用时间
    Wednesday []TimeSlot `json:"wednesday"` // 周三可用时间
    Thursday  []TimeSlot `json:"thursday"`  // 周四可用时间
    Friday    []TimeSlot `json:"friday"`    // 周五可用时间
    Saturday  []TimeSlot `json:"saturday"`  // 周六可用时间
    Sunday    []TimeSlot `json:"sunday"`    // 周日可用时间
}

type TimeSlot struct {
    StartTime string `json:"start_time"` // 开始时间 "08:00"
    EndTime   string `json:"end_time"`   // 结束时间 "18:00"
}
```

## 业务逻辑设计

### 设备搜索和筛选
支持多维度的农机搜索功能:

1. **地理位置搜索**: 基于GPS定位查找附近的农机
2. **设备类型筛选**: 按农机类型和品牌筛选
3. **价格范围筛选**: 按租金价格范围筛选
4. **时间可用性**: 查询指定时间段可用的设备
5. **评分排序**: 按设备评分和机主信誉排序

### 订单处理流程

1. **订单创建**: 租客选择设备并提交租赁申请
2. **机主确认**: 设备所有者确认订单和租赁条件
3. **支付处理**: 租客支付押金和租金
4. **设备交付**: 按约定时间地点交付设备
5. **作业进行**: 设备正常作业期间
6. **设备归还**: 作业完成后归还设备
7. **费用结算**: 根据实际使用情况结算费用
8. **服务评价**: 双方互相评价

### 定价策略

支持多种灵活的定价模式:

```go
// 根据计费方式计算费用
func CalculateRentalCost(machine *Machine, billingMethod string, quantity float64) int64 {
    switch billingMethod {
    case "hourly":
        return machine.HourlyRate * int64(quantity)
    case "daily":
        return machine.DailyRate * int64(quantity)
    case "per_acre":
        return machine.PerAcreRate * int64(quantity)
    default:
        return 0
    }
}

// 动态定价 - 根据需求和季节调整价格
func DynamicPricing(baseRate int64, demand float64, season string) int64 {
    var multiplier float64 = 1.0
    
    // 需求调节因子
    if demand > 0.8 {
        multiplier += 0.2 // 高需求时提价20%
    } else if demand < 0.3 {
        multiplier -= 0.1 // 低需求时降价10%
    }
    
    // 季节调节因子
    switch season {
    case "spring", "autumn": // 春耕秋收旺季
        multiplier += 0.15
    case "summer": // 夏季农忙
        multiplier += 0.1
    case "winter": // 冬季农闲
        multiplier -= 0.2
    }
    
    return int64(float64(baseRate) * multiplier)
}
```

## 数据库索引设计

### 关键索引
1. **农机设备表 (machines)**:
   - `machine_code` 唯一索引
   - `machine_type` 普通索引
   - `status` 普通索引
   - `owner_id` 普通索引
   - 地理位置索引: (`longitude`, `latitude`)
   - 联合索引: (`machine_type`, `status`)

2. **租赁订单表 (rental_orders)**:
   - `order_no` 唯一索引
   - `machine_id` 普通索引
   - `renter_id` 普通索引
   - `owner_id` 普通索引
   - `status` 普通索引
   - 联合索引: (`renter_id`, `status`)
   - 联合索引: (`machine_id`, `status`)

## 使用示例

### 农机设备注册
```go
// 创建农机设备
machine := &model.Machine{
    MachineCode: "TR001",
    MachineName: "约翰迪尔6B-1404拖拉机",
    MachineType: "tractor",
    Brand:       "约翰迪尔",
    Model:       "6B-1404",
    Description: "140马力大型拖拉机，适用于大面积农田作业",
    OwnerID:     ownerUserID,
    OwnerType:   "individual",
    Province:    "山东省",
    City:        "济南市",
    County:      "章丘区",
    Longitude:   117.5348,
    Latitude:    36.7140,
    HourlyRate:  15000, // 150元/小时
    DailyRate:   100000, // 1000元/天
    PerAcreRate: 8000,  // 80元/亩
    DepositAmount: 500000, // 5000元押金
    Status:      "available",
}

// 设置设备规格
specs := model.MachineSpecs{
    Power:        "140马力",
    WorkingWidth: "2.5米",
    Weight:       "4200公斤",
    FuelType:     "柴油",
    YearOfMake:   2022,
    WorkingSpeed: "6-15公里/小时",
    Capacity:     "4立方米",
    Other: map[string]string{
        "变速箱": "16前进档+8倒档",
        "驱动方式": "四轮驱动",
        "轮胎规格": "480/70R34",
    },
}
specsJSON, _ := json.Marshal(specs)
machine.Specifications = string(specsJSON)

// 设置可用时间
schedule := model.AvailableSchedule{
    Monday:    []model.TimeSlot{{"08:00", "18:00"}},
    Tuesday:   []model.TimeSlot{{"08:00", "18:00"}},
    Wednesday: []model.TimeSlot{{"08:00", "18:00"}},
    Thursday:  []model.TimeSlot{{"08:00", "18:00"}},
    Friday:    []model.TimeSlot{{"08:00", "18:00"}},
    Saturday:  []model.TimeSlot{{"08:00", "16:00"}},
    Sunday:    []model.TimeSlot{}, // 周日不可用
}
scheduleJSON, _ := json.Marshal(schedule)
machine.AvailableSchedule = string(scheduleJSON)

// 保存设备
if err := db.Create(machine).Error; err != nil {
    return fmt.Errorf("创建农机设备失败: %w", err)
}
```

### 创建租赁订单
```go
// 生成订单号
orderNo := generateOrderNo()

// 计算租赁费用
startTime := time.Date(2024, 4, 15, 8, 0, 0, 0, time.Local)
endTime := time.Date(2024, 4, 15, 18, 0, 0, 0, time.Local)
duration := int(endTime.Sub(startTime).Hours())
unitPrice := machine.HourlyRate
subtotal := unitPrice * int64(duration)
deposit := machine.DepositAmount
total := subtotal + deposit

// 创建租赁订单
order := &model.RentalOrder{
    OrderNo:        orderNo,
    MachineID:      machine.ID,
    RenterID:       renterUserID,
    OwnerID:        machine.OwnerID,
    StartTime:      startTime,
    EndTime:        endTime,
    RentalDuration: duration,
    RentalLocation: "山东省济南市章丘区某农场",
    ContactPerson:  "张农民",
    ContactPhone:   "13800138000",
    BillingMethod:  "hourly",
    UnitPrice:      unitPrice,
    Quantity:       float64(duration),
    SubtotalAmount: subtotal,
    DepositAmount:  deposit,
    TotalAmount:    total,
    Status:         "pending",
    Remarks:        "春耕期间土地翻耕作业",
}

// 保存订单
if err := db.Create(order).Error; err != nil {
    return fmt.Errorf("创建租赁订单失败: %w", err)
}

// 更新设备状态为已预订
if err := db.Model(&machine).Update("status", "rented").Error; err != nil {
    log.Printf("更新设备状态失败: %v", err)
}
```

### 地理位置搜索
```go
// 基于GPS位置搜索附近的农机
func SearchNearbyMachines(longitude, latitude float64, radius float64, machineType string) ([]model.Machine, error) {
    var machines []model.Machine
    
    // 使用经纬度范围查询 (简化版本，实际应使用地理空间索引)
    latRange := radius / 111.0 // 大约每度111公里
    lonRange := radius / (111.0 * math.Cos(latitude*math.Pi/180))
    
    query := db.Where("status = ?", "available").
        Where("latitude BETWEEN ? AND ?", latitude-latRange, latitude+latRange).
        Where("longitude BETWEEN ? AND ?", longitude-lonRange, longitude+lonRange)
    
    if machineType != "" {
        query = query.Where("machine_type = ?", machineType)
    }
    
    if err := query.Find(&machines).Error; err != nil {
        return nil, err
    }
    
    // 计算实际距离并排序
    for i := range machines {
        distance := calculateDistance(latitude, longitude, machines[i].Latitude, machines[i].Longitude)
        // 可以将距离信息添加到返回结果中
    }
    
    return machines, nil
}

// 计算两点间距离 (Haversine公式)
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
    const R = 6371 // 地球半径，单位：千米
    
    dLat := (lat2 - lat1) * math.Pi / 180
    dLon := (lon2 - lon1) * math.Pi / 180
    
    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
        math.Sin(dLon/2)*math.Sin(dLon/2)
    
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    
    return R * c
}
```

### 订单状态更新
```go
// 订单状态流转管理
func UpdateOrderStatus(orderID uint64, newStatus string, operatorID uint64) error {
    var order model.RentalOrder
    if err := db.First(&order, orderID).Error; err != nil {
        return err
    }
    
    // 验证状态转换的合法性
    if !isValidStatusTransition(order.Status, newStatus) {
        return fmt.Errorf("无效的状态转换: %s -> %s", order.Status, newStatus)
    }
    
    // 更新订单状态
    oldStatus := order.Status
    order.Status = newStatus
    
    // 根据状态更新相关信息
    switch newStatus {
    case "confirmed":
        // 机主确认订单
        order.OwnerConfirmedAt = &time.Time{}
        *order.OwnerConfirmedAt = time.Now()
        
    case "paid":
        // 支付完成
        order.PaidAmount = order.TotalAmount
        order.PaidAt = &time.Time{}
        *order.PaidAt = time.Now()
        
    case "in_progress":
        // 开始作业
        // 可以记录实际开始时间
        
    case "completed":
        // 作业完成
        // 可以记录实际结束时间和作业数据
        
    case "cancelled":
        // 订单取消，释放设备
        if err := db.Model(&model.Machine{}).Where("id = ?", order.MachineID).
            Update("status", "available").Error; err != nil {
            return err
        }
    }
    
    // 保存订单更新
    return db.Save(&order).Error
}

func isValidStatusTransition(fromStatus, toStatus string) bool {
    validTransitions := map[string][]string{
        "pending":     {"confirmed", "cancelled"},
        "confirmed":   {"paid", "cancelled"},
        "paid":        {"in_progress", "cancelled"},
        "in_progress": {"completed", "disputed"},
        "completed":   {"disputed"},
        "cancelled":   {},
        "disputed":    {"completed", "cancelled"},
    }
    
    allowedStates, exists := validTransitions[fromStatus]
    if !exists {
        return false
    }
    
    for _, state := range allowedStates {
        if state == toStatus {
            return true
        }
    }
    
    return false
}
```

## 性能优化

### 地理位置查询优化
1. **空间索引**: 使用MySQL的空间索引优化地理位置查询
2. **缓存机制**: 热门区域的农机信息缓存
3. **分页查询**: 大量结果使用分页避免性能问题

### 实时性优化
1. **状态同步**: 设备状态实时同步更新
2. **WebSocket**: 订单状态变更实时推送
3. **消息队列**: 异步处理费用结算和评价通知

### 数据分析
1. **使用统计**: 统计设备使用率和收益
2. **热力图**: 生成农机需求热力图
3. **价格分析**: 分析不同地区和时间的价格趋势

## 安全考虑

1. **设备认证**: 农机设备需要认证后才能上架
2. **用户验证**: 租赁双方身份验证
3. **保险保障**: 租赁期间的保险和责任划分
4. **争议处理**: 完善的争议处理机制
5. **数据保护**: 位置信息等敏感数据保护 