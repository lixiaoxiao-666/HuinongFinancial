package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"backend/internal/api"
	"backend/internal/conf"
	"backend/internal/data"
	"backend/internal/service"
	"backend/pkg"
)

// timePtr 返回时间指针
func timePtr(t time.Time) *time.Time {
	return &t
}

// float64Ptr 返回float64指针
func float64Ptr(f float64) *float64 {
	return &f
}

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

	// === 农机租赁相关模拟数据 ===

	// 1. 创建农机主用户
	lessorUsers := []data.User{
		{
			UserID:       "lessor_001",
			Phone:        "13900001001",
			PasswordHash: userPasswordHash,
			Nickname:     "农机主-李四",
			Status:       0,
		},
		{
			UserID:       "lessor_002",
			Phone:        "13900001002",
			PasswordHash: userPasswordHash,
			Nickname:     "农机主-王五",
			Status:       0,
		},
		{
			UserID:       "lessor_003",
			Phone:        "13900001003",
			PasswordHash: userPasswordHash,
			Nickname:     "农机主-赵六",
			Status:       0,
		},
	}

	for _, user := range lessorUsers {
		if err := dataLayer.DB.Create(&user).Error; err != nil {
			logger.Warn("创建农机主用户失败", zap.Error(err), zap.String("user_id", user.UserID))
		}
	}

	// 2. 创建承租方用户
	lesseeUsers := []data.User{
		{
			UserID:       "lessee_001",
			Phone:        "13800002001",
			PasswordHash: userPasswordHash,
			Nickname:     "承租方-陈七",
			Status:       0,
		},
		{
			UserID:       "lessee_002",
			Phone:        "13800002002",
			PasswordHash: userPasswordHash,
			Nickname:     "承租方-孙八",
			Status:       0,
		},
	}

	for _, user := range lesseeUsers {
		if err := dataLayer.DB.Create(&user).Error; err != nil {
			logger.Warn("创建承租方用户失败", zap.Error(err), zap.String("user_id", user.UserID))
		}
	}

	// 3. 创建农机主资质数据
	lessorQualifications := []data.LessorQualification{
		{
			UserID:                 "lessor_001",
			RealName:               "李四",
			IDCardNumber:           "370102198501151234",
			Phone:                  "13900001001",
			Address:                "山东省济南市历下区农机大户路88号",
			BusinessLicenseNumber:  "91370102MA3K123456",
			FarmScale:              "500亩",
			VerificationStatus:     "VERIFIED",
			CreditRating:           "AAA",
			TotalMachineryCount:    8,
			SuccessfulLeasingCount: 45,
			AverageRating:          float64Ptr(4.8),
			VerifiedAt:             timePtr(time.Now().AddDate(0, -2, 0)),
		},
		{
			UserID:                 "lessor_002",
			RealName:               "王五",
			IDCardNumber:           "410102198703201567",
			Phone:                  "13900001002",
			Address:                "河南省郑州市中原区农业园区北路66号",
			BusinessLicenseNumber:  "91410102MA45678901",
			FarmScale:              "300亩",
			VerificationStatus:     "VERIFIED",
			CreditRating:           "AA",
			TotalMachineryCount:    5,
			SuccessfulLeasingCount: 32,
			AverageRating:          float64Ptr(4.6),
			VerifiedAt:             timePtr(time.Now().AddDate(0, -1, -15)),
		},
		{
			UserID:                 "lessor_003",
			RealName:               "赵六",
			IDCardNumber:           "130102198909101890",
			Phone:                  "13900001003",
			Address:                "河北省石家庄市长安区现代农业示范区18号",
			BusinessLicenseNumber:  "",
			FarmScale:              "150亩",
			VerificationStatus:     "PENDING",
			CreditRating:           "A",
			TotalMachineryCount:    3,
			SuccessfulLeasingCount: 12,
			AverageRating:          float64Ptr(4.3),
			VerifiedAt:             nil,
		},
	}

	for _, qualification := range lessorQualifications {
		if err := dataLayer.DB.Create(&qualification).Error; err != nil {
			logger.Warn("创建农机主资质失败", zap.Error(err), zap.String("user_id", qualification.UserID))
		}
	}

	// 4. 创建农机设备数据
	farmMachinery := []data.FarmMachinery{
		{
			MachineryID:  pkg.GenerateMachineryID(),
			OwnerUserID:  "lessor_001",
			Type:         "拖拉机",
			BrandModel:   "约翰迪尔 5080R",
			Description:  "80马力四驱拖拉机，配备GPS导航，适用于耕地、播种、收割等多种作业",
			Images:       []byte(`["/images/tractor_001_1.jpg", "/images/tractor_001_2.jpg"]`),
			DailyRent:    500.00,
			Deposit:      float64Ptr(8000.00),
			LocationText: "山东省济南市历下区",
			LocationGeo:  "36.6512,117.1201",
			Status:       "AVAILABLE",
			PublishedAt:  timePtr(time.Now().AddDate(0, 0, -7)),
		},
		{
			MachineryID:  pkg.GenerateMachineryID(),
			OwnerUserID:  "lessor_001",
			Type:         "收割机",
			BrandModel:   "久保田 4LZ-2.5",
			Description:  "履带式联合收割机，适用于水稻、小麦收割，效率高，损失率低",
			Images:       []byte(`["/images/harvester_001_1.jpg", "/images/harvester_001_2.jpg"]`),
			DailyRent:    800.00,
			Deposit:      float64Ptr(15000.00),
			LocationText: "山东省济南市历下区",
			LocationGeo:  "36.6512,117.1201",
			Status:       "AVAILABLE",
			PublishedAt:  timePtr(time.Now().AddDate(0, 0, -5)),
		},
		{
			MachineryID:  pkg.GenerateMachineryID(),
			OwnerUserID:  "lessor_002",
			Type:         "播种机",
			BrandModel:   "东方红 2BQX-12",
			Description:  "气力式精密播种机，12行，适用于玉米、大豆等作物播种",
			Images:       []byte(`["/images/seeder_001_1.jpg"]`),
			DailyRent:    350.00,
			Deposit:      float64Ptr(5000.00),
			LocationText: "河南省郑州市中原区",
			LocationGeo:  "34.7578,113.6486",
			Status:       "AVAILABLE",
			PublishedAt:  timePtr(time.Now().AddDate(0, 0, -3)),
		},
		{
			MachineryID:  pkg.GenerateMachineryID(),
			OwnerUserID:  "lessor_002",
			Type:         "旋耕机",
			BrandModel:   "常发 1GQN-200",
			Description:  "重型旋耕机，2米工作幅宽，适用于深耕整地作业",
			Images:       []byte(`["/images/rotavator_001_1.jpg", "/images/rotavator_001_2.jpg"]`),
			DailyRent:    280.00,
			Deposit:      float64Ptr(3500.00),
			LocationText: "河南省郑州市中原区",
			LocationGeo:  "34.7578,113.6486",
			Status:       "RENTED",
			PublishedAt:  timePtr(time.Now().AddDate(0, 0, -10)),
		},
		{
			MachineryID:  pkg.GenerateMachineryID(),
			OwnerUserID:  "lessor_003",
			Type:         "植保无人机",
			BrandModel:   "大疆 T30",
			Description:  "农用植保无人机，30升载药量，高效精准喷洒",
			Images:       []byte(`["/images/drone_001_1.jpg"]`),
			DailyRent:    450.00,
			Deposit:      float64Ptr(12000.00),
			LocationText: "河北省石家庄市长安区",
			LocationGeo:  "38.0467,114.5143",
			Status:       "AVAILABLE",
			PublishedAt:  timePtr(time.Now().AddDate(0, 0, -1)),
		},
	}

	for _, machinery := range farmMachinery {
		if err := dataLayer.DB.Create(&machinery).Error; err != nil {
			logger.Warn("创建农机设备失败", zap.Error(err), zap.String("machinery_id", machinery.MachineryID))
		}
	}

	// 5. 创建农机租赁申请数据
	startDate1 := time.Now().AddDate(0, 0, 7)  // 7天后开始
	endDate1 := time.Now().AddDate(0, 0, 12)   // 12天后结束
	startDate2 := time.Now().AddDate(0, 0, 15) // 15天后开始
	endDate2 := time.Now().AddDate(0, 0, 20)   // 20天后结束

	leasingApplications := []data.MachineryLeasingApplication{
		{
			ApplicationID:      pkg.GenerateApplicationID(),
			LesseeUserID:       "lessee_001",
			LessorUserID:       "lessor_001",
			MachineryID:        farmMachinery[0].MachineryID, // 拖拉机
			RequestedStartDate: startDate1,
			RequestedEndDate:   endDate1,
			RentalDays:         5,
			TotalAmount:        2500.00, // 500 * 5天
			DepositAmount:      float64Ptr(8000.00),
			UsagePurpose:       "春季耕地作业",
			LesseeNotes:        "希望能够提供操作培训，预计作业面积50亩",
			ApplicationStatus:  "PENDING_REVIEW",
			SubmittedAt:        time.Now().AddDate(0, 0, -1),
		},
		{
			ApplicationID:      pkg.GenerateApplicationID(),
			LesseeUserID:       "lessee_002",
			LessorUserID:       "lessor_001",
			MachineryID:        farmMachinery[1].MachineryID, // 收割机
			RequestedStartDate: startDate2,
			RequestedEndDate:   endDate2,
			RentalDays:         5,
			TotalAmount:        4000.00, // 800 * 5天
			DepositAmount:      float64Ptr(15000.00),
			UsagePurpose:       "小麦收割作业",
			LesseeNotes:        "需要租赁期间包含周末，作业时间较紧",
			ApplicationStatus:  "PENDING_REVIEW",
			SubmittedAt:        time.Now().AddDate(0, 0, -2),
		},
		{
			ApplicationID:      "ml_test_001", // 固定ID用于AI测试
			LesseeUserID:       "lessee_001",
			LessorUserID:       "lessor_002",
			MachineryID:        farmMachinery[2].MachineryID, // 播种机
			RequestedStartDate: time.Now().AddDate(0, 0, 3),
			RequestedEndDate:   time.Now().AddDate(0, 0, 6),
			RentalDays:         3,
			TotalAmount:        1050.00, // 350 * 3天
			DepositAmount:      float64Ptr(5000.00),
			UsagePurpose:       "玉米播种作业",
			LesseeNotes:        "需要AI智能审批测试",
			ApplicationStatus:  "PENDING_REVIEW",
			SubmittedAt:        time.Now(),
		},
	}

	for _, application := range leasingApplications {
		if err := dataLayer.DB.Create(&application).Error; err != nil {
			logger.Warn("创建农机租赁申请失败", zap.Error(err), zap.String("application_id", application.ApplicationID))
		}
	}

	// 6. 创建承租方用户详情
	lesseeProfiles := []data.UserProfile{
		{
			UserID:           "lessee_001",
			RealName:         "陈七",
			IDCardNumber:     "320102199206151234",
			Address:          "江苏省南京市玄武区农业合作社路28号",
			Gender:           1,
			BirthDate:        timePtr(time.Date(1992, 6, 15, 0, 0, 0, 0, time.UTC)),
			Occupation:       "种植大户",
			AnnualIncome:     180000.00,
			CreditAuthAgreed: true,
		},
		{
			UserID:           "lessee_002",
			RealName:         "孙八",
			IDCardNumber:     "430102199408201567",
			Address:          "湖南省长沙市芙蓉区现代农业园区99号",
			Gender:           1,
			BirthDate:        timePtr(time.Date(1994, 8, 20, 0, 0, 0, 0, time.UTC)),
			Occupation:       "农民专业合作社成员",
			AnnualIncome:     120000.00,
			CreditAuthAgreed: true,
		},
	}

	for _, profile := range lesseeProfiles {
		if err := dataLayer.DB.Create(&profile).Error; err != nil {
			logger.Warn("创建承租方用户详情失败", zap.Error(err), zap.String("user_id", profile.UserID))
		}
	}

	// 7. 创建农机主用户详情
	lessorProfiles := []data.UserProfile{
		{
			UserID:           "lessor_001",
			RealName:         "李四",
			IDCardNumber:     "370102198501151234",
			Address:          "山东省济南市历下区农机大户路88号",
			Gender:           1,
			BirthDate:        timePtr(time.Date(1985, 1, 15, 0, 0, 0, 0, time.UTC)),
			Occupation:       "农机专业合作社理事长",
			AnnualIncome:     350000.00,
			CreditAuthAgreed: true,
		},
		{
			UserID:           "lessor_002",
			RealName:         "王五",
			IDCardNumber:     "410102198703201567",
			Address:          "河南省郑州市中原区农业园区北路66号",
			Gender:           1,
			BirthDate:        timePtr(time.Date(1987, 3, 20, 0, 0, 0, 0, time.UTC)),
			Occupation:       "农机服务公司经理",
			AnnualIncome:     280000.00,
			CreditAuthAgreed: true,
		},
		{
			UserID:           "lessor_003",
			RealName:         "赵六",
			IDCardNumber:     "130102198909101890",
			Address:          "河北省石家庄市长安区现代农业示范区18号",
			Gender:           1,
			BirthDate:        timePtr(time.Date(1989, 9, 10, 0, 0, 0, 0, time.UTC)),
			Occupation:       "农机大户",
			AnnualIncome:     200000.00,
			CreditAuthAgreed: true,
		},
	}

	for _, profile := range lessorProfiles {
		if err := dataLayer.DB.Create(&profile).Error; err != nil {
			logger.Warn("创建农机主用户详情失败", zap.Error(err), zap.String("user_id", profile.UserID))
		}
	}

	logger.Info("示例数据初始化完成")
	return nil
}
