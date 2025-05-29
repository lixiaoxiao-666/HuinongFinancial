package database

import (
	"fmt"
	"time"

	"huinong-backend/internal/config"
	"huinong-backend/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// NewConnection 创建数据库连接
func NewConnection(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	// 构建DSN
	dsn := buildDSN(cfg)

	// GORM配置
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "", // 表名前缀
			SingularTable: false, // 使用复数表名
		},
		Logger: getLogLevel(cfg),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层sql.DB对象进行连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	return db, nil
}

// buildDSN 构建数据源名称
func buildDSN(cfg *config.DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
	)
}

// getLogLevel 获取日志级别
func getLogLevel(cfg *config.DatabaseConfig) logger.Interface {
	// 根据环境设置不同的日志级别
	logLevel := logger.Silent // 默认静默模式
	
	// 这里可以根据配置来设置不同的日志级别
	// 在开发环境可以设置为Info级别以便调试
	// if cfg.Debug {
	//     logLevel = logger.Info
	// }
	
	return logger.Default.LogMode(logLevel)
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	// 用户相关表
	if err := db.AutoMigrate(
		&model.User{},
		&model.UserAuth{},
		&model.UserSession{},
		&model.UserTag{},
		&model.OAUser{},
		&model.OARole{},
	); err != nil {
		return fmt.Errorf("用户表迁移失败: %w", err)
	}

	// 贷款相关表
	if err := db.AutoMigrate(
		&model.LoanProduct{},
		&model.LoanApplication{},
		&model.ApprovalLog{},
		&model.DifyWorkflowLog{},
	); err != nil {
		return fmt.Errorf("贷款表迁移失败: %w", err)
	}

	// 农机相关表
	if err := db.AutoMigrate(
		&model.Machine{},
		&model.RentalOrder{},
	); err != nil {
		return fmt.Errorf("农机表迁移失败: %w", err)
	}

	// 通用表
	if err := db.AutoMigrate(
		&model.Article{},
		&model.Category{},
		&model.Expert{},
		&model.SystemConfig{},
		&model.FileUpload{},
		&model.OfflineQueue{},
		&model.APILog{},
	); err != nil {
		return fmt.Errorf("通用表迁移失败: %w", err)
	}

	return nil
}

// InitDefaultData 初始化默认数据
func InitDefaultData(db *gorm.DB) error {
	// 初始化系统配置
	if err := initSystemConfigs(db); err != nil {
		return fmt.Errorf("初始化系统配置失败: %w", err)
	}

	// 初始化默认角色
	if err := initDefaultRoles(db); err != nil {
		return fmt.Errorf("初始化默认角色失败: %w", err)
	}

	// 初始化管理员用户
	if err := initAdminUser(db); err != nil {
		return fmt.Errorf("初始化管理员用户失败: %w", err)
	}

	// 初始化文章分类
	if err := initCategories(db); err != nil {
		return fmt.Errorf("初始化文章分类失败: %w", err)
	}

	return nil
}

// initSystemConfigs 初始化系统配置
func initSystemConfigs(db *gorm.DB) error {
	configs := []model.SystemConfig{
		{
			ConfigKey:   "app.name",
			ConfigValue: "数字惠农系统",
			ConfigType:  "string",
			ConfigGroup: "system",
			Description: "应用程序名称",
		},
		{
			ConfigKey:   "app.version",
			ConfigValue: "1.0.0",
			ConfigType:  "string",
			ConfigGroup: "system",
			Description: "应用程序版本",
		},
		{
			ConfigKey:   "loan.auto_approval_threshold",
			ConfigValue: "50000",
			ConfigType:  "int",
			ConfigGroup: "business",
			Description: "贷款自动审批阈值(分)",
		},
		{
			ConfigKey:   "file.max_upload_size",
			ConfigValue: "10485760",
			ConfigType:  "int",
			ConfigGroup: "system",
			Description: "文件上传最大大小(字节)",
		},
	}

	for _, config := range configs {
		var count int64
		db.Model(&model.SystemConfig{}).Where("config_key = ?", config.ConfigKey).Count(&count)
		if count == 0 {
			if err := db.Create(&config).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// initDefaultRoles 初始化默认角色
func initDefaultRoles(db *gorm.DB) error {
	roles := []model.OARole{
		{
			Name:        "super_admin",
			DisplayName: "超级管理员",
			Description: "拥有所有权限的超级管理员角色",
			Permissions: `{
				"loan_management": ["view", "create", "update", "approve", "delete"],
				"machine_management": ["view", "create", "update", "delete"],
				"user_management": ["view", "create", "update", "freeze", "delete"],
				"content_management": ["view", "create", "update", "publish", "delete"],
				"system_settings": ["view", "update"],
				"data_analytics": ["view", "export"]
			}`,
			IsSuper: true,
			Status:  "active",
		},
		{
			Name:        "loan_officer",
			DisplayName: "信贷员",
			Description: "负责贷款业务审核的信贷员角色",
			Permissions: `{
				"loan_management": ["view", "update", "approve"],
				"user_management": ["view"],
				"data_analytics": ["view"]
			}`,
			IsSuper: false,
			Status:  "active",
		},
		{
			Name:        "content_manager",
			DisplayName: "内容管理员",
			Description: "负责内容管理的角色",
			Permissions: `{
				"content_management": ["view", "create", "update", "publish"],
				"user_management": ["view"],
				"data_analytics": ["view"]
			}`,
			IsSuper: false,
			Status:  "active",
		},
	}

	for _, role := range roles {
		var count int64
		db.Model(&model.OARole{}).Where("name = ?", role.Name).Count(&count)
		if count == 0 {
			if err := db.Create(&role).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// initAdminUser 初始化管理员用户
func initAdminUser(db *gorm.DB) error {
	var superAdminRole model.OARole
	if err := db.Where("name = ?", "super_admin").First(&superAdminRole).Error; err != nil {
		return fmt.Errorf("找不到超级管理员角色: %w", err)
	}

	var count int64
	db.Model(&model.OAUser{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		adminUser := model.OAUser{
			Username:     "admin",
			Email:        "admin@huinong.com",
			PasswordHash: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Salt:         "huinong_salt_2024",
			RealName:     "系统管理员",
			RoleID:       superAdminRole.ID,
			Department:   "技术部",
			Position:     "系统管理员",
			Status:       "active",
		}

		if err := db.Create(&adminUser).Error; err != nil {
			return err
		}
	}

	return nil
}

// initCategories 初始化文章分类
func initCategories(db *gorm.DB) error {
	categories := []model.Category{
		{
			Name:        "policy",
			DisplayName: "政策资讯",
			Description: "农业政策相关资讯",
			SortOrder:   1,
			Status:      "active",
		},
		{
			Name:        "technology",
			DisplayName: "技术指导",
			Description: "农业技术指导文章",
			SortOrder:   2,
			Status:      "active",
		},
		{
			Name:        "market",
			DisplayName: "市场信息",
			Description: "农产品市场信息",
			SortOrder:   3,
			Status:      "active",
		},
		{
			Name:        "news",
			DisplayName: "行业新闻",
			Description: "农业行业新闻资讯",
			SortOrder:   4,
			Status:      "active",
		},
	}

	for _, category := range categories {
		var count int64
		db.Model(&model.Category{}).Where("name = ?", category.Name).Count(&count)
		if count == 0 {
			if err := db.Create(&category).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// HealthCheck 数据库健康检查
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	return nil
}

// GetDatabaseStats 获取数据库连接统计信息
func GetDatabaseStats(db *gorm.DB) (map[string]interface{}, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}

	stats := sqlDB.Stats()
	
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":              stats.InUse,
		"idle":                stats.Idle,
		"wait_count":          stats.WaitCount,
		"wait_duration":       stats.WaitDuration.String(),
		"max_idle_closed":     stats.MaxIdleClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
	}, nil
} 