package repository

import (
	"context"
	"fmt"
	"huinong-backend/internal/model"

	"gorm.io/gorm"
)

// loanRepository 贷款Repository实现
type loanRepository struct {
	db *gorm.DB
}

// NewLoanRepository 创建贷款Repository实例
func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &loanRepository{
		db: db,
	}
}

// ==================== 贷款产品相关 ====================

// GetProductByID 根据ID获取产品
func (r *loanRepository) GetProductByID(ctx context.Context, id uint) (*model.LoanProduct, error) {
	var product model.LoanProduct
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("产品不存在")
		}
		return nil, err
	}
	return &product, nil
}

// GetProductByCode 根据产品代码获取产品
func (r *loanRepository) GetProductByCode(ctx context.Context, code string) (*model.LoanProduct, error) {
	var product model.LoanProduct
	err := r.db.WithContext(ctx).Where("product_code = ?", code).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("产品不存在")
		}
		return nil, err
	}
	return &product, nil
}

// GetProductsByUserType 根据用户类型获取产品列表
func (r *loanRepository) GetProductsByUserType(ctx context.Context, userType string, page, limit int) ([]*model.LoanProduct, int64, error) {
	var products []*model.LoanProduct
	var total int64

	query := r.db.WithContext(ctx).Model(&model.LoanProduct{}).
		Where("is_active = ? AND (eligible_user_type = ? OR eligible_user_type = '' OR eligible_user_type IS NULL)",
			true, userType)

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * limit
	err = query.Order("sort_order DESC, created_at DESC").
		Offset(offset).Limit(limit).
		Find(&products).Error

	return products, total, err
}

// GetAllProducts 获取所有产品
func (r *loanRepository) GetAllProducts(ctx context.Context) ([]*model.LoanProduct, error) {
	var products []*model.LoanProduct
	err := r.db.WithContext(ctx).Order("sort_order DESC, created_at DESC").Find(&products).Error
	return products, err
}

// GetActiveProducts 获取活跃产品
func (r *loanRepository) GetActiveProducts(ctx context.Context, userType string) ([]*model.LoanProduct, error) {
	var products []*model.LoanProduct
	query := r.db.WithContext(ctx).Where("is_active = ?", true)

	if userType != "" {
		query = query.Where("eligible_user_type = ? OR eligible_user_type = '' OR eligible_user_type IS NULL", userType)
	}

	err := query.Order("sort_order DESC, created_at DESC").Find(&products).Error
	return products, err
}

// CreateProduct 创建产品
func (r *loanRepository) CreateProduct(ctx context.Context, product *model.LoanProduct) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// UpdateProduct 更新产品
func (r *loanRepository) UpdateProduct(ctx context.Context, product *model.LoanProduct) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// DeleteProduct 删除产品
func (r *loanRepository) DeleteProduct(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.LoanProduct{}, id).Error
}

// ListProducts 产品列表查询
func (r *loanRepository) ListProducts(ctx context.Context, req *ListProductsRequest) (*ListProductsResponse, error) {
	var products []*model.LoanProduct
	var total int64

	query := r.db.WithContext(ctx).Model(&model.LoanProduct{})

	// 条件筛选
	if req.ProductType != "" {
		query = query.Where("product_type = ?", req.ProductType)
	}
	if req.Status != "" {
		isActive := req.Status == "active"
		query = query.Where("is_active = ?", isActive)
	}
	if req.Keyword != "" {
		query = query.Where("product_name LIKE ? OR product_code LIKE ? OR description LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.Limit
	err = query.Order("sort_order DESC, created_at DESC").
		Offset(offset).Limit(req.Limit).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	return &ListProductsResponse{
		Products: products,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}

// ==================== 贷款申请相关 ====================

// CreateApplication 创建申请
func (r *loanRepository) CreateApplication(ctx context.Context, application *model.LoanApplication) error {
	return r.db.WithContext(ctx).Create(application).Error
}

// GetApplicationByID 根据ID获取申请
func (r *loanRepository) GetApplicationByID(ctx context.Context, id uint) (*model.LoanApplication, error) {
	var application model.LoanApplication
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		First(&application, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("申请不存在")
		}
		return nil, err
	}
	return &application, nil
}

// GetApplicationByNo 根据申请编号获取申请
func (r *loanRepository) GetApplicationByNo(ctx context.Context, applicationNo string) (*model.LoanApplication, error) {
	var application model.LoanApplication
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		Where("application_no = ?", applicationNo).
		First(&application).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("申请不存在")
		}
		return nil, err
	}
	return &application, nil
}

// GetUserApplications 获取用户申请列表
func (r *loanRepository) GetUserApplications(ctx context.Context, userID uint, page, limit int, status string) ([]*model.LoanApplication, int64, error) {
	var applications []*model.LoanApplication
	var total int64

	query := r.db.WithContext(ctx).Model(&model.LoanApplication{}).
		Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	err = query.Preload("Product").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&applications).Error

	return applications, total, err
}

// GetPendingApplications 获取待处理申请
func (r *loanRepository) GetPendingApplications(ctx context.Context, limit, offset int) ([]*model.LoanApplication, error) {
	var applications []*model.LoanApplication
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		Where("status IN ?", []string{"pending", "reviewing"}).
		Order("created_at ASC").
		Offset(offset).Limit(limit).
		Find(&applications).Error
	return applications, err
}

// UpdateApplication 更新申请
func (r *loanRepository) UpdateApplication(ctx context.Context, application *model.LoanApplication) error {
	return r.db.WithContext(ctx).Save(application).Error
}

// UpdateApplicationStatus 更新申请状态
func (r *loanRepository) UpdateApplicationStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&model.LoanApplication{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// GetApplicationsForAdmin 获取管理员申请列表
func (r *loanRepository) GetApplicationsForAdmin(ctx context.Context, page, limit int, status string) ([]*model.LoanApplication, int64, error) {
	var applications []*model.LoanApplication
	var total int64

	query := r.db.WithContext(ctx).Model(&model.LoanApplication{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	err = query.Preload("User").
		Preload("Product").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&applications).Error

	return applications, total, err
}

// GetApplicationStatistics 获取申请统计
func (r *loanRepository) GetApplicationStatistics(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总申请数
	var totalApplications int64
	err := r.db.WithContext(ctx).Model(&model.LoanApplication{}).Count(&totalApplications).Error
	if err != nil {
		return nil, err
	}
	stats["total_applications"] = totalApplications

	// 本月申请数
	var monthlyApplications int64
	err = r.db.WithContext(ctx).Model(&model.LoanApplication{}).
		Where("DATE_FORMAT(created_at, '%Y-%m') = DATE_FORMAT(NOW(), '%Y-%m')").
		Count(&monthlyApplications).Error
	if err != nil {
		return nil, err
	}
	stats["monthly_applications"] = monthlyApplications

	// 状态统计
	var statusStats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	err = r.db.WithContext(ctx).Model(&model.LoanApplication{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&statusStats).Error
	if err != nil {
		return nil, err
	}

	statusMap := make(map[string]int64)
	for _, stat := range statusStats {
		statusMap[stat.Status] = stat.Count
	}
	stats["status_statistics"] = statusMap

	return stats, nil
}

// ==================== 审批日志相关 ====================

// CreateApprovalLog 创建审批日志
func (r *loanRepository) CreateApprovalLog(ctx context.Context, log *model.ApprovalLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetApprovalLogs 获取审批日志
func (r *loanRepository) GetApprovalLogs(ctx context.Context, applicationID uint) ([]*model.ApprovalLog, error) {
	var logs []*model.ApprovalLog
	err := r.db.WithContext(ctx).
		Where("application_id = ?", applicationID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// ==================== Dify工作流日志相关 ====================

// CreateDifyLog 创建Dify日志
func (r *loanRepository) CreateDifyLog(ctx context.Context, log *model.DifyWorkflowLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetDifyLogs 获取Dify日志
func (r *loanRepository) GetDifyLogs(ctx context.Context, applicationID uint) ([]*model.DifyWorkflowLog, error) {
	var logs []*model.DifyWorkflowLog
	err := r.db.WithContext(ctx).
		Where("application_id = ?", applicationID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}
