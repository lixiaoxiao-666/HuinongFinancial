package repository

import (
	"context"
	"fmt"
	"time"

	"huinong-backend/internal/model"

	"gorm.io/gorm"
)

// machineRepository 农机Repository实现
type machineRepository struct {
	db *gorm.DB
}

// NewMachineRepository 创建农机Repository实例
func NewMachineRepository(db *gorm.DB) MachineRepository {
	return &machineRepository{
		db: db,
	}
}

// ==================== 设备管理 ====================

// Create 创建农机设备
func (r *machineRepository) Create(ctx context.Context, machine *model.Machine) error {
	return r.db.WithContext(ctx).Create(machine).Error
}

// GetByID 根据ID获取农机设备
func (r *machineRepository) GetByID(ctx context.Context, id uint64) (*model.Machine, error) {
	var machine model.Machine
	err := r.db.WithContext(ctx).
		Preload("Owner").
		First(&machine, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("农机设备不存在")
		}
		return nil, err
	}
	return &machine, nil
}

// GetByCode 根据设备编码获取农机设备
func (r *machineRepository) GetByCode(ctx context.Context, code string) (*model.Machine, error) {
	var machine model.Machine
	err := r.db.WithContext(ctx).
		Preload("Owner").
		Where("machine_code = ?", code).
		First(&machine).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("农机设备不存在")
		}
		return nil, err
	}
	return &machine, nil
}

// Update 更新农机设备
func (r *machineRepository) Update(ctx context.Context, machine *model.Machine) error {
	return r.db.WithContext(ctx).Save(machine).Error
}

// Delete 删除农机设备
func (r *machineRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Machine{}, id).Error
}

// List 农机设备列表查询
func (r *machineRepository) List(ctx context.Context, req *ListMachinesRequest) (*ListMachinesResponse, error) {
	var machines []*model.Machine
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Machine{}).
		Preload("Owner")

	// 条件筛选
	if req.MachineType != "" {
		query = query.Where("machine_type = ?", req.MachineType)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Province != "" {
		query = query.Where("province = ?", req.Province)
	}
	if req.City != "" {
		query = query.Where("city = ?", req.City)
	}
	if req.County != "" {
		query = query.Where("county = ?", req.County)
	}
	if req.Keyword != "" {
		query = query.Where("machine_name LIKE ? OR brand LIKE ? OR model LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 地理位置搜索
	if req.Longitude != 0 && req.Latitude != 0 && req.Radius > 0 {
		// 使用Haversine公式计算距离
		query = query.Where(`
			(6371 * acos(
				cos(radians(?)) * cos(radians(latitude)) * 
				cos(radians(longitude) - radians(?)) + 
				sin(radians(?)) * sin(radians(latitude))
			)) <= ?`,
			req.Latitude, req.Longitude, req.Latitude, req.Radius)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.Limit
	err = query.Order("created_at DESC").
		Offset(offset).Limit(req.Limit).
		Find(&machines).Error
	if err != nil {
		return nil, err
	}

	return &ListMachinesResponse{
		Machines: machines,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}

// ==================== 地理位置搜索 ====================

// SearchNearby 附近农机搜索
func (r *machineRepository) SearchNearby(ctx context.Context, longitude, latitude, radius float64, machineType string) ([]*model.Machine, error) {
	var machines []*model.Machine
	query := r.db.WithContext(ctx).
		Preload("Owner").
		Where("status = ?", "available")

	if machineType != "" {
		query = query.Where("machine_type = ?", machineType)
	}

	// 使用Haversine公式计算距离
	err := query.Where(`
		(6371 * acos(
			cos(radians(?)) * cos(radians(latitude)) * 
			cos(radians(longitude) - radians(?)) + 
			sin(radians(?)) * sin(radians(latitude))
		)) <= ?`,
		latitude, longitude, latitude, radius).
		Order("created_at DESC").
		Find(&machines).Error

	return machines, err
}

// GetByLocation 根据行政区域获取农机设备
func (r *machineRepository) GetByLocation(ctx context.Context, province, city, county string) ([]*model.Machine, error) {
	var machines []*model.Machine
	query := r.db.WithContext(ctx).
		Preload("Owner").
		Where("status = ?", "available")

	if province != "" {
		query = query.Where("province = ?", province)
	}
	if city != "" {
		query = query.Where("city = ?", city)
	}
	if county != "" {
		query = query.Where("county = ?", county)
	}

	err := query.Order("created_at DESC").Find(&machines).Error
	return machines, err
}

// ==================== 租赁订单 ====================

// CreateOrder 创建租赁订单
func (r *machineRepository) CreateOrder(ctx context.Context, order *model.RentalOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// GetOrderByID 根据ID获取租赁订单
func (r *machineRepository) GetOrderByID(ctx context.Context, id uint64) (*model.RentalOrder, error) {
	var order model.RentalOrder
	err := r.db.WithContext(ctx).
		Preload("Machine").
		Preload("Renter").
		Preload("Owner").
		First(&order, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("租赁订单不存在")
		}
		return nil, err
	}
	return &order, nil
}

// GetOrderByNo 根据订单编号获取租赁订单
func (r *machineRepository) GetOrderByNo(ctx context.Context, orderNo string) (*model.RentalOrder, error) {
	var order model.RentalOrder
	err := r.db.WithContext(ctx).
		Preload("Machine").
		Preload("Renter").
		Preload("Owner").
		Where("order_no = ?", orderNo).
		First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("租赁订单不存在")
		}
		return nil, err
	}
	return &order, nil
}

// UpdateOrder 更新租赁订单
func (r *machineRepository) UpdateOrder(ctx context.Context, order *model.RentalOrder) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// ListOrders 租赁订单列表查询
func (r *machineRepository) ListOrders(ctx context.Context, req *ListOrdersRequest) (*ListOrdersResponse, error) {
	var orders []*model.RentalOrder
	var total int64

	query := r.db.WithContext(ctx).Model(&model.RentalOrder{}).
		Preload("Machine").
		Preload("Renter").
		Preload("Owner")

	// 条件筛选
	if req.MachineID > 0 {
		query = query.Where("machine_id = ?", req.MachineID)
	}
	if req.RenterID > 0 {
		query = query.Where("renter_id = ?", req.RenterID)
	}
	if req.OwnerID > 0 {
		query = query.Where("owner_id = ?", req.OwnerID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.StartDate != "" && req.EndDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.Limit
	err = query.Order("created_at DESC").
		Offset(offset).Limit(req.Limit).
		Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return &ListOrdersResponse{
		Orders: orders,
		Total:  total,
		Page:   req.Page,
		Limit:  req.Limit,
	}, nil
}

// GetUserOrders 获取用户订单列表
func (r *machineRepository) GetUserOrders(ctx context.Context, userID uint64, userType string, limit, offset int) ([]*model.RentalOrder, error) {
	var orders []*model.RentalOrder
	query := r.db.WithContext(ctx).
		Preload("Machine").
		Preload("Renter").
		Preload("Owner")

	// 根据用户类型查询
	if userType == "renter" {
		query = query.Where("renter_id = ?", userID)
	} else if userType == "owner" {
		query = query.Where("owner_id = ?", userID)
	} else {
		// 默认查询所有相关订单
		query = query.Where("renter_id = ? OR owner_id = ?", userID, userID)
	}

	err := query.Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&orders).Error
	return orders, err
}

// ==================== 统计方法 ====================

// GetMachineCount 获取农机总数
func (r *machineRepository) GetMachineCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Machine{}).Count(&count).Error
	return count, err
}

// GetAvailableMachineCount 获取可用农机数量
func (r *machineRepository) GetAvailableMachineCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Machine{}).
		Where("status = ?", "available").
		Count(&count).Error
	return count, err
}

// GetOrderCount 获取订单总数
func (r *machineRepository) GetOrderCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.RentalOrder{}).Count(&count).Error
	return count, err
}

// ==================== 时间冲突检查 ====================

// CheckTimeConflict 检查租赁时间冲突
func (r *machineRepository) CheckTimeConflict(ctx context.Context, machineID uint64, startTime, endTime time.Time, excludeOrderID uint64) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.RentalOrder{}).
		Where("machine_id = ?", machineID).
		Where("status IN (?)", []string{"pending", "confirmed", "in_progress"}).
		Where("NOT (end_time <= ? OR start_time >= ?)", startTime, endTime)

	// 排除指定的订单ID（用于更新操作）
	if excludeOrderID > 0 {
		query = query.Where("id != ?", excludeOrderID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
