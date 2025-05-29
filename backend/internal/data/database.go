package data

import (
	"context"
	"fmt"
	"time"

	"backend/internal/conf"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Data 数据层结构
type Data struct {
	DB    *gorm.DB
	Redis *redis.Client
	Log   *zap.Logger
}

// NewData 创建数据层实例
func NewData(config *conf.Config, log *zap.Logger) (*Data, error) {
	// 初始化数据库连接
	db, err := initDB(config, log)
	if err != nil {
		return nil, fmt.Errorf("初始化数据库失败: %w", err)
	}

	// 初始化Redis连接
	rdb, err := initRedis(config, log)
	if err != nil {
		return nil, fmt.Errorf("初始化Redis失败: %w", err)
	}

	return &Data{
		DB:    db,
		Redis: rdb,
		Log:   log,
	}, nil
}

// initDB 初始化数据库连接
func initDB(config *conf.Config, log *zap.Logger) (*gorm.DB, error) {
	dsn := config.Database.GetDSN()

	// 配置GORM日志
	gormLogger := logger.Default
	if config.App.Env == "production" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.Database.ConnMaxLifetime)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	log.Info("数据库连接成功", zap.String("host", config.Database.Host), zap.Int("port", config.Database.Port))

	return db, nil
}

// initRedis 初始化Redis连接
func initRedis(config *conf.Config, log *zap.Logger) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         config.Redis.GetRedisAddr(),
		Password:     config.Redis.Password,
		DB:           config.Redis.DB,
		PoolSize:     config.Redis.PoolSize,
		MinIdleConns: config.Redis.MinIdleConns,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis连接测试失败: %w", err)
	}

	log.Info("Redis连接成功", zap.String("addr", config.Redis.GetRedisAddr()))

	return rdb, nil
}

// AutoMigrate 自动迁移数据库表结构
func (d *Data) AutoMigrate() error {
	return d.DB.AutoMigrate(
		&User{},
		&UserProfile{},
		&LoanProduct{},
		&LoanApplication{},
		&LoanApplicationHistory{},
		&UploadedFile{},
		&FarmMachinery{},
		&MachineryLeasingOrder{},
		&OAUser{},
		&SystemConfiguration{},
		// AI智能体相关表
		&AIAnalysisResult{},
		&WorkflowExecution{},
		&AIModelConfig{},
		&ExternalDataQuery{},
		&AIAgentLog{},
		// 农机租赁审批相关表
		&MachineryLeasingApplication{},
		&MachineryLeasingApprovalHistory{},
		&LessorQualification{},
		&MachineryLeasingAIResult{},
	)
}

// Close 关闭数据库和Redis连接
func (d *Data) Close() error {
	// 关闭Redis连接
	if err := d.Redis.Close(); err != nil {
		d.Log.Error("关闭Redis连接失败", zap.Error(err))
	}

	// 关闭数据库连接
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
