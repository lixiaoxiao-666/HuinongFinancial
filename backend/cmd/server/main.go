package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"huinong-backend/internal/cache"
	"huinong-backend/internal/config"
	"huinong-backend/internal/database"
	"huinong-backend/internal/repository"
	"huinong-backend/internal/router"
	"huinong-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// @title 数字惠农API
// @version 1.0
// @description 数字惠农APP及OA后台管理系统API接口文档
// @termsOfService https://example.com/terms/

// @contact.name API 技术支持
// @contact.url https://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer Token Authentication

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 设置Gin模式
	if cfg.App.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库连接
	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 执行数据库迁移
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化Redis缓存
	var cacheClient cache.CacheInterface
	if cfg.Redis.Host != "" {
		cacheClient = cache.NewRedisClient(
			fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
			cfg.Redis.Password,
			0, // 使用默认数据库
		)
		log.Println("Redis缓存已连接")
	} else {
		log.Println("Redis缓存未启用")
	}

	// 初始化Repository层
	userRepo := repository.NewUserRepository(db)
	loanRepo := repository.NewLoanRepository(db)
	// machineRepo := repository.NewMachineRepository(db)
	// articleRepo := repository.NewArticleRepository(db)
	// expertRepo := repository.NewExpertRepository(db)
	// fileRepo := repository.NewFileRepository(db)
	// systemRepo := repository.NewSystemRepository(db)
	// oaRepo := repository.NewOARepository(db)

	// 初始化Service层
	userService := service.NewUserService(userRepo, cfg.JWT.SecretKey)
	// loanService := service.NewLoanService(loanRepo, userRepo, cacheClient)
	// machineService := service.NewMachineService(machineRepo, userRepo)
	// articleService := service.NewArticleService(articleRepo, userRepo, cacheClient)
	// expertService := service.NewExpertService(expertRepo, userRepo)
	// fileService := service.NewFileService(fileRepo, cfg.File.UploadPath)
	// systemService := service.NewSystemService(systemRepo, cacheClient)
	// oaService := service.NewOAService(oaRepo, cfg.JWT.SecretKey)

	// 路由配置
	routerConfig := &router.RouterConfig{
		UserService:    userService,
		LoanService:    nil, // loanService,
		MachineService: nil, // machineService,
		ArticleService: nil, // articleService,
		ExpertService:  nil, // expertService,
		FileService:    nil, // fileService,
		SystemService:  nil, // systemService,
		OAService:      nil, // oaService,
		JWTSecret:      cfg.JWT.SecretKey,
	}

	// 设置路由
	r := router.SetupRouter(routerConfig)

	// 启动HTTP服务器
	serverAddr := fmt.Sprintf(":%d", cfg.App.Port)
	log.Printf("服务器启动在 %s", serverAddr)
	log.Printf("环境: %s", cfg.App.Mode)
	log.Printf("Swagger文档: http://localhost%s/swagger/index.html", serverAddr)

	// 优雅关闭
	go func() {
		if err := r.Run(serverAddr); err != nil {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 清理资源
	if cacheClient != nil {
		log.Println("关闭Redis连接...")
	}

	// 避免未使用变量错误
	_ = loanRepo

	log.Println("服务器已关闭")
}
