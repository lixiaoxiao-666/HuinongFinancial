package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/api"
	"backend/internal/conf"
	"backend/internal/data"
	"backend/internal/service"
	"backend/pkg"

	"go.uber.org/zap"
)

func main() {
	// 加载配置
	config, err := conf.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	logger, err := initLogger(config)
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.Sync()

	logger.Info("数字惠农APP后端服务启动中...",
		zap.String("name", config.App.Name),
		zap.String("version", config.App.Version),
		zap.String("env", config.App.Env),
		zap.Int("port", config.Server.Port),
	)

	// 初始化数据层
	dataLayer, err := data.NewData(config, logger)
	if err != nil {
		logger.Fatal("初始化数据层失败", zap.Error(err))
	}
	defer dataLayer.Close()

	// 自动迁移数据库表结构
	if err := dataLayer.AutoMigrate(); err != nil {
		logger.Fatal("数据库迁移失败", zap.Error(err))
	}
	logger.Info("数据库迁移完成")

	// 初始化示例数据
	if err := initSampleData(dataLayer, logger); err != nil {
		logger.Error("初始化示例数据失败", zap.Error(err))
	}

	// 初始化JWT管理器
	jwtManager := pkg.NewJWTManager(config.JWT.Secret, config.JWT.Expire)

	// 创建AdminService并初始化默认OA用户
	adminService := service.NewAdminService(dataLayer, jwtManager, logger)
	if err := adminService.CreateDefaultOAUsers(); err != nil {
		logger.Error("创建默认OA用户失败", zap.Error(err))
	}

	// 初始化路由器
	router := api.NewRouter(config, dataLayer, jwtManager, logger)
	engine := router.SetupRoutes()

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", config.Server.Port)
	logger.Info("服务器启动成功", zap.String("addr", serverAddr))

	// 优雅关闭
	go func() {
		if err := engine.Run(serverAddr); err != nil {
			logger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")
	logger.Info("服务器已关闭")
}

// initLogger 初始化日志
func initLogger(config *conf.Config) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	if config.App.Env == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}

	return logger, nil
}

// initSampleData 初始化示例数据
func initSampleData(dataLayer *data.Data, logger *zap.Logger) error {
	// 检查是否已有贷款产品数据
	var count int64
	dataLayer.DB.Model(&data.LoanProduct{}).Count(&count)
	if count > 0 {
		logger.Info("示例数据已存在，跳过初始化")
		return nil
	}

	logger.Info("正在初始化示例数据...")

	// 创建示例贷款产品
	products := []data.LoanProduct{
		{
			ProductID:             pkg.GenerateProductID(),
			Name:                  "春耕助力贷",
			Description:           "专为春耕生产设计，利率优惠，快速审批",
			Category:              "种植贷",
			MinAmount:             5000,
			MaxAmount:             50000,
			MinTermMonths:         6,
			MaxTermMonths:         24,
			InterestRateYearly:    "4.5% - 6.0%",
			RepaymentMethods:      []byte(`["等额本息", "先息后本"]`),
			ApplicationConditions: "1. 年满18周岁的农户；2. 有稳定的农业收入；3. 信用记录良好",
			RequiredDocuments:     []byte(`[{"type": "ID_CARD", "desc": "申请人身份证"}, {"type": "LAND_CONTRACT", "desc": "土地承包合同"}]`),
			Status:                0,
		},
		{
			ProductID:             pkg.GenerateProductID(),
			Name:                  "农机购置贷",
			Description:           "支持农户购买农业机械，助力农业现代化",
			Category:              "设备贷",
			MinAmount:             10000,
			MaxAmount:             200000,
			MinTermMonths:         12,
			MaxTermMonths:         60,
			InterestRateYearly:    "5.0% - 7.0%",
			RepaymentMethods:      []byte(`["等额本息", "等额本金"]`),
			ApplicationConditions: "1. 年满18周岁的农户；2. 有购机需求证明；3. 有还款能力",
			RequiredDocuments:     []byte(`[{"type": "ID_CARD", "desc": "申请人身份证"}, {"type": "PURCHASE_CONTRACT", "desc": "农机购买合同"}]`),
			Status:                0,
		},
		{
			ProductID:             pkg.GenerateProductID(),
			Name:                  "丰收种植贷",
			Description:           "支持大棚种植、果蔬种植等现代农业项目",
			Category:              "种植贷",
			MinAmount:             8000,
			MaxAmount:             80000,
			MinTermMonths:         12,
			MaxTermMonths:         36,
			InterestRateYearly:    "4.8% - 6.5%",
			RepaymentMethods:      []byte(`["等额本息", "季末还息到期还本"]`),
			ApplicationConditions: "1. 有种植经验；2. 有土地使用权证明；3. 无不良信用记录",
			RequiredDocuments:     []byte(`[{"type": "ID_CARD", "desc": "身份证"}, {"type": "LAND_USE_CERT", "desc": "土地使用权证"}]`),
			Status:                0,
		},
		{
			ProductID:             pkg.GenerateProductID(),
			Name:                  "养殖创业贷",
			Description:           "支持家禽、家畜养殖业发展，助力规模化养殖",
			Category:              "养殖贷",
			MinAmount:             15000,
			MaxAmount:             150000,
			MinTermMonths:         6,
			MaxTermMonths:         48,
			InterestRateYearly:    "5.2% - 7.5%",
			RepaymentMethods:      []byte(`["等额本息", "按季付息到期还本"]`),
			ApplicationConditions: "1. 有养殖场地；2. 有相关养殖技术；3. 有销售渠道",
			RequiredDocuments:     []byte(`[{"type": "ID_CARD", "desc": "身份证"}, {"type": "FARM_CERT", "desc": "养殖场证明"}]`),
			Status:                0,
		},
		{
			ProductID:             pkg.GenerateProductID(),
			Name:                  "智慧农业贷",
			Description:           "支持智能灌溉、物联网设备等现代农业技术应用",
			Category:              "设备贷",
			MinAmount:             20000,
			MaxAmount:             300000,
			MinTermMonths:         18,
			MaxTermMonths:         72,
			InterestRateYearly:    "4.8% - 6.8%",
			RepaymentMethods:      []byte(`["等额本息", "等额本金"]`),
			ApplicationConditions: "1. 有现代农业发展规划；2. 有技术团队支持；3. 有稳定收入来源",
			RequiredDocuments:     []byte(`[{"type": "ID_CARD", "desc": "身份证"}, {"type": "BUSINESS_PLAN", "desc": "项目计划书"}]`),
			Status:                0,
		},
		{
			ProductID:             pkg.GenerateProductID(),
			Name:                  "农村电商贷",
			Description:           "支持农产品电商平台、直播带货等新业态发展",
			Category:              "经营贷",
			MinAmount:             3000,
			MaxAmount:             50000,
			MinTermMonths:         3,
			MaxTermMonths:         24,
			InterestRateYearly:    "6.0% - 8.5%",
			RepaymentMethods:      []byte(`["等额本息", "按月付息到期还本"]`),
			ApplicationConditions: "1. 有电商运营经验；2. 有农产品货源；3. 有良好信用记录",
			RequiredDocuments:     []byte(`[{"type": "ID_CARD", "desc": "身份证"}, {"type": "BUSINESS_LICENSE", "desc": "营业执照"}]`),
			Status:                0,
		},
		{
			ProductID:             pkg.GenerateProductID(),
			Name:                  "合作社发展贷",
			Description:           "支持农民专业合作社扩大经营规模，提升服务能力",
			Category:              "经营贷",
			MinAmount:             50000,
			MaxAmount:             500000,
			MinTermMonths:         12,
			MaxTermMonths:         60,
			InterestRateYearly:    "4.2% - 6.2%",
			RepaymentMethods:      []byte(`["等额本息", "按季付息到期还本"]`),
			ApplicationConditions: "1. 注册满1年的合作社；2. 有稳定的经营收入；3. 无重大违法记录",
			RequiredDocuments:     []byte(`[{"type": "COOP_LICENSE", "desc": "合作社执照"}, {"type": "FINANCIAL_REPORT", "desc": "财务报表"}]`),
			Status:                0,
		},
		{
			ProductID:             pkg.GenerateProductID(),
			Name:                  "绿色农业贷",
			Description:           "支持有机农业、生态农业等绿色环保项目",
			Category:              "种植贷",
			MinAmount:             12000,
			MaxAmount:             120000,
			MinTermMonths:         12,
			MaxTermMonths:         48,
			InterestRateYearly:    "4.0% - 5.8%",
			RepaymentMethods:      []byte(`["等额本息", "季末还息到期还本"]`),
			ApplicationConditions: "1. 有绿色认证或申请中；2. 有环保设施；3. 有市场销路",
			RequiredDocuments:     []byte(`[{"type": "ID_CARD", "desc": "身份证"}, {"type": "GREEN_CERT", "desc": "绿色认证书"}]`),
			Status:                0,
		},
		{
			ProductID:             pkg.GenerateProductID(),
			Name:                  "水产养殖贷",
			Description:           "支持鱼类、虾类、蟹类等水产养殖业发展",
			Category:              "养殖贷",
			MinAmount:             8000,
			MaxAmount:             100000,
			MinTermMonths:         6,
			MaxTermMonths:         36,
			InterestRateYearly:    "5.5% - 7.8%",
			RepaymentMethods:      []byte(`["等额本息", "按季付息到期还本"]`),
			ApplicationConditions: "1. 有水产养殖经验；2. 有养殖场地；3. 有销售渠道",
			RequiredDocuments:     []byte(`[{"type": "ID_CARD", "desc": "身份证"}, {"type": "POND_CERT", "desc": "养殖场地证明"}]`),
			Status:                0,
		},
		{
			ProductID:             pkg.GenerateProductID(),
			Name:                  "温室大棚贷",
			Description:           "支持现代化温室大棚建设和设备采购",
			Category:              "设备贷",
			MinAmount:             30000,
			MaxAmount:             500000,
			MinTermMonths:         24,
			MaxTermMonths:         84,
			InterestRateYearly:    "4.5% - 6.5%",
			RepaymentMethods:      []byte(`["等额本息", "等额本金"]`),
			ApplicationConditions: "1. 有大棚建设规划；2. 有土地使用权；3. 有技术支持",
			RequiredDocuments:     []byte(`[{"type": "ID_CARD", "desc": "身份证"}, {"type": "CONSTRUCTION_PLAN", "desc": "建设规划书"}]`),
			Status:                0,
		},
		{
			ProductID:             pkg.GenerateProductID(),
			Name:                  "乡村旅游贷",
			Description:           "支持农家乐、民宿等乡村旅游项目发展",
			Category:              "经营贷",
			MinAmount:             20000,
			MaxAmount:             300000,
			MinTermMonths:         12,
			MaxTermMonths:         60,
			InterestRateYearly:    "5.8% - 8.0%",
			RepaymentMethods:      []byte(`["等额本息", "按季付息到期还本"]`),
			ApplicationConditions: "1. 有旅游资源；2. 有经营许可；3. 有市场调研",
			RequiredDocuments:     []byte(`[{"type": "ID_CARD", "desc": "身份证"}, {"type": "TOURISM_LICENSE", "desc": "旅游经营许可证"}]`),
			Status:                0,
		},
		{
			ProductID:             pkg.GenerateProductID(),
			Name:                  "农产品加工贷",
			Description:           "支持农产品深加工设备采购和技术升级",
			Category:              "设备贷",
			MinAmount:             25000,
			MaxAmount:             400000,
			MinTermMonths:         18,
			MaxTermMonths:         72,
			InterestRateYearly:    "4.8% - 6.8%",
			RepaymentMethods:      []byte(`["等额本息", "等额本金"]`),
			ApplicationConditions: "1. 有加工经验；2. 有原料供应；3. 有销售渠道",
			RequiredDocuments:     []byte(`[{"type": "ID_CARD", "desc": "身份证"}, {"type": "PROCESSING_LICENSE", "desc": "食品加工许可证"}]`),
			Status:                0,
		},
		{
			ProductID:             pkg.GenerateProductID(),
			Name:                  "果树种植贷",
			Description:           "支持果园建设、果树种植和果品销售",
			Category:              "种植贷",
			MinAmount:             10000,
			MaxAmount:             150000,
			MinTermMonths:         12,
			MaxTermMonths:         60,
			InterestRateYearly:    "4.5% - 6.8%",
			RepaymentMethods:      []byte(`["等额本息", "季末还息到期还本"]`),
			ApplicationConditions: "1. 有种植经验；2. 有土地承包权；3. 有销售计划",
			RequiredDocuments:     []byte(`[{"type": "ID_CARD", "desc": "身份证"}, {"type": "ORCHARD_PLAN", "desc": "果园种植计划"}]`),
			Status:                0,
		},
		{
			ProductID:             pkg.GenerateProductID(),
			Name:                  "畜牧设备贷",
			Description:           "支持现代化畜牧设备采购和养殖场建设",
			Category:              "设备贷",
			MinAmount:             15000,
			MaxAmount:             250000,
			MinTermMonths:         12,
			MaxTermMonths:         60,
			InterestRateYearly:    "5.0% - 7.2%",
			RepaymentMethods:      []byte(`["等额本息", "等额本金"]`),
			ApplicationConditions: "1. 有畜牧经验；2. 有养殖场地；3. 有环保手续",
			RequiredDocuments:     []byte(`[{"type": "ID_CARD", "desc": "身份证"}, {"type": "LIVESTOCK_PERMIT", "desc": "畜牧养殖许可证"}]`),
			Status:                0,
		},
		{
			ProductID:             pkg.GenerateProductID(),
			Name:                  "农业科技贷",
			Description:           "支持农业科技创新项目和技术研发投入",
			Category:              "经营贷",
			MinAmount:             30000,
			MaxAmount:             600000,
			MinTermMonths:         18,
			MaxTermMonths:         72,
			InterestRateYearly:    "4.2% - 6.0%",
			RepaymentMethods:      []byte(`["等额本息", "按季付息到期还本"]`),
			ApplicationConditions: "1. 有技术团队；2. 有创新项目；3. 有市场前景",
			RequiredDocuments:     []byte(`[{"type": "ID_CARD", "desc": "身份证"}, {"type": "TECH_PLAN", "desc": "技术创新计划书"}]`),
			Status:                0,
		},
	}

	for _, product := range products {
		if err := dataLayer.DB.Create(&product).Error; err != nil {
			return fmt.Errorf("创建示例贷款产品失败: %w", err)
		}
	}

	// 创建示例OA用户
	hashedPassword, _ := pkg.HashPassword("admin123")
	oaUser := data.OAUser{
		OAUserID:     "oa_admin001",
		Username:     "admin",
		PasswordHash: hashedPassword,
		Role:         "ADMIN",
		DisplayName:  "系统管理员",
		Email:        "admin@example.com",
		Status:       0,
	}

	if err := dataLayer.DB.Create(&oaUser).Error; err != nil {
		return fmt.Errorf("创建示例OA用户失败: %w", err)
	}

	// 创建示例普通用户
	userPasswordHash, _ := pkg.HashPassword("user123")
	testUser := data.User{
		UserID:       pkg.GenerateUserID(),
		Phone:        "13800138000",
		PasswordHash: userPasswordHash,
		Nickname:     "测试农户",
		Status:       0,
	}

	if err := dataLayer.DB.Create(&testUser).Error; err != nil {
		logger.Warn("创建示例用户失败", zap.Error(err))
	}

	// 创建用户详情
	userProfile := data.UserProfile{
		UserID:           testUser.UserID,
		RealName:         "张三",
		IDCardNumber:     "31010119900101****",
		Address:          "上海市浦东新区XX镇XX村",
		CreditAuthAgreed: true,
	}

	if err := dataLayer.DB.Create(&userProfile).Error; err != nil {
		logger.Warn("创建示例用户详情失败", zap.Error(err))
	}

	// 添加AI Agent测试专用数据
	// 创建user_001用户（用于AI Agent测试）
	aiTestUser := data.User{
		UserID:       "user_001",
		Phone:        "13800138001",
		PasswordHash: userPasswordHash,
		Nickname:     "AI测试用户",
		Status:       0,
	}
	if err := dataLayer.DB.Create(&aiTestUser).Error; err != nil {
		logger.Warn("创建AI测试用户失败", zap.Error(err))
	}

	// 创建AI测试用户详情
	birthDate, _ := time.Parse("2006-01-02", "1990-01-01")
	aiTestUserProfile := data.UserProfile{
		UserID:           "user_001",
		RealName:         "张三",
		IDCardNumber:     "110101199001011234",
		Address:          "北京市朝阳区测试街道123号",
		Gender:           1, // 1为男性, 0为女性
		BirthDate:        &birthDate,
		Occupation:       "软件工程师",
		AnnualIncome:     200000.00,
		CreditAuthAgreed: true,
	}
	if err := dataLayer.DB.Create(&aiTestUserProfile).Error; err != nil {
		logger.Warn("创建AI测试用户详情失败", zap.Error(err))
	}

	// 创建product_001产品（如果不存在）
	aiTestProduct := data.LoanProduct{
		ProductID:             "product_001",
		Name:                  "个人信用贷",
		Description:           "专为个人信用贷款设计，无需抵押，快速审批",
		Category:              "personal_credit",
		MinAmount:             10000,
		MaxAmount:             500000,
		MinTermMonths:         6,
		MaxTermMonths:         36,
		InterestRateYearly:    "8.5%",
		RepaymentMethods:      []byte(`["等额本息", "先息后本"]`),
		ApplicationConditions: "年收入不低于10万，征信良好",
		RequiredDocuments:     []byte(`[{"type": "ID_CARD", "desc": "身份证"}, {"type": "INCOME_PROOF", "desc": "收入证明"}]`),
		Status:                0,
	}
	if err := dataLayer.DB.Create(&aiTestProduct).Error; err != nil {
		logger.Warn("创建AI测试产品失败", zap.Error(err))
	}

	// 创建test_app_001申请（用于AI Agent测试）
	aiTestApplicantSnapshot, _ := json.Marshal(map[string]interface{}{
		"real_name":      "张三",
		"id_card_number": "110101199001011234",
		"address":        "北京市朝阳区测试街道123号",
		"phone":          "13800138001",
		"gender":         "male",
		"birth_date":     "1990-01-01",
		"occupation":     "软件工程师",
		"annual_income":  200000.00,
		"work_years":     5,
		"education":      "bachelor",
		"marital_status": "single",
		"has_house":      true,
		"has_car":        false,
		"credit_level":   "good",
	})

	aiTestApplication := data.LoanApplication{
		ApplicationID:     "test_app_001",
		UserID:            "user_001",
		ProductID:         "product_001",
		AmountApplied:     100000,
		TermMonthsApplied: 24,
		Purpose:           "装修",
		Status:            "pending_review",
		ApplicantSnapshot: aiTestApplicantSnapshot,
	}
	if err := dataLayer.DB.Create(&aiTestApplication).Error; err != nil {
		logger.Warn("创建AI测试申请失败", zap.Error(err))
	}

	// 创建AI测试相关的上传文件记录
	aiTestFiles := []data.UploadedFile{
		{
			FileID:      "file_001",
			UserID:      "user_001",
			FileName:    "身份证正面.jpg",
			FileType:    "image/jpeg",
			FileSize:    1024000,
			StoragePath: "/uploads/user_001/id_card_front.jpg",
			Purpose:     "id_card",
		},
		{
			FileID:      "file_002",
			UserID:      "user_001",
			FileName:    "收入证明.pdf",
			FileType:    "application/pdf",
			FileSize:    2048000,
			StoragePath: "/uploads/user_001/income_proof.pdf",
			Purpose:     "income_proof",
		},
		{
			FileID:      "file_003",
			UserID:      "user_001",
			FileName:    "银行流水.pdf",
			FileType:    "application/pdf",
			FileSize:    3072000,
			StoragePath: "/uploads/user_001/bank_statement.pdf",
			Purpose:     "bank_statement",
		},
	}

	for _, file := range aiTestFiles {
		if err := dataLayer.DB.Create(&file).Error; err != nil {
			logger.Warn("创建AI测试文件记录失败", zap.Error(err), zap.String("file_id", file.FileID))
		}
	}

	// 初始化AI审批开关配置
	aiConfig := data.SystemConfiguration{
		ConfigKey:   "ai_approval_enabled",
		ConfigValue: "true",
		Description: "AI审批功能开关",
	}
	dataLayer.DB.Save(&aiConfig)

	logger.Info("示例数据初始化完成")
	return nil
}
