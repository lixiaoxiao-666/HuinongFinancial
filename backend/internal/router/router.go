package router

import (
	"huinong-backend/internal/handler"
	"huinong-backend/internal/middleware"
	"huinong-backend/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RouterConfig 路由配置
type RouterConfig struct {
	UserService    service.UserService
	LoanService    service.LoanService
	MachineService service.MachineService
	ArticleService service.ArticleService
	ExpertService  service.ExpertService
	FileService    service.FileService
	SystemService  service.SystemService
	OAService      service.OAService
	DifyService    service.DifyService
	JWTSecret      string
	DifyAPIToken   string
}

// SetupRouter 设置路由
func SetupRouter(config *RouterConfig) *gin.Engine {
	// 创建Gin实例
	r := gin.New()

	// 基础中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.CORS())

	// 创建认证中间件
	authMiddleware := middleware.NewAuthMiddleware(config.JWTSecret)

	// 创建Handler实例
	userHandler := handler.NewUserHandler(config.UserService)
	difyHandler := handler.NewDifyHandler(
		config.UserService,
		config.LoanService,
		config.MachineService,
	)
	// TODO: 创建其他Handler实例
	// loanHandler := handler.NewLoanHandler(config.LoanService)
	// machineHandler := handler.NewMachineHandler(config.MachineService)
	// articleHandler := handler.NewArticleHandler(config.ArticleService)
	// expertHandler := handler.NewExpertHandler(config.ExpertService)
	// fileHandler := handler.NewFileHandler(config.FileService)
	// systemHandler := handler.NewSystemHandler(config.SystemService)
	// oaHandler := handler.NewOAHandler(config.OAService)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "数字惠农API服务正在运行",
		})
	})

	// API版本分组
	api := r.Group("/api")
	{
		// 公开API（无需认证）
		public := api.Group("/public")
		{
			public.GET("/version", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"version":     "1.0.0",
					"name":        "数字惠农API",
					"description": "数字惠农APP及OA后台管理系统API接口",
				})
			})
		}

		// 内部API（供Dify工作流调用）
		internal := api.Group("/internal")
		internal.Use(middleware.DifyAuthMiddleware(config.DifyAPIToken))
		{
			// Dify工作流回调接口
			dify := internal.Group("/dify")
			{
				// 贷款相关接口
				loan := dify.Group("/loan")
				{
					loan.POST("/get-application-details", difyHandler.GetLoanApplicationDetails)
					loan.POST("/submit-assessment", difyHandler.SubmitRiskAssessment)
				}

				// 农机相关接口
				machine := dify.Group("/machine")
				{
					machine.POST("/get-rental-details", difyHandler.GetMachineRentalDetails)
				}

				// 征信相关接口
				credit := dify.Group("/credit")
				{
					credit.POST("/query", difyHandler.QueryCreditReport)
				}
			}
		}

		// 认证相关API（无需认证）
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/refresh", userHandler.RefreshToken)
			// TODO: 添加其他认证接口
			// auth.POST("/forgot-password", userHandler.ForgotPassword)
			// auth.POST("/reset-password", userHandler.ResetPassword)
			// auth.POST("/verify-sms", userHandler.VerifySMS)
			// auth.POST("/send-sms", userHandler.SendSMS)
		}

		// 用户API（需要认证）
		user := api.Group("/user")
		user.Use(authMiddleware.RequireAuth())
		{
			// 用户资料管理
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.PUT("/password", userHandler.ChangePassword)
			user.POST("/logout", userHandler.Logout)

			// 用户认证
			userAuth := user.Group("/auth")
			{
				userAuth.POST("/real-name", userHandler.SubmitRealNameAuth)
				userAuth.POST("/bank-card", userHandler.SubmitBankCardAuth)
				// TODO: 添加其他认证接口
				// userAuth.GET("/status", userHandler.GetAuthStatus)
				// userAuth.POST("/credit", userHandler.SubmitCreditAuth)
			}

			// 用户标签
			user.GET("/tags", userHandler.GetUserTags)
			user.POST("/tags", userHandler.AddUserTag)
			user.DELETE("/tags/:tag_key", userHandler.RemoveUserTag)

			// TODO: 贷款相关API
			// loan := user.Group("/loan")
			// {
			//     loan.GET("/products", loanHandler.GetProducts)
			//     loan.POST("/applications", loanHandler.SubmitApplication)
			//     loan.GET("/applications", loanHandler.GetUserApplications)
			//     loan.GET("/applications/:id", loanHandler.GetApplication)
			//     loan.PUT("/applications/:id", loanHandler.UpdateApplication)
			//     loan.DELETE("/applications/:id", loanHandler.CancelApplication)
			// }

			// TODO: 农机相关API
			// machine := user.Group("/machine")
			// {
			//     machine.POST("/register", machineHandler.RegisterMachine)
			//     machine.GET("/list", machineHandler.GetUserMachines)
			//     machine.GET("/:id", machineHandler.GetMachine)
			//     machine.PUT("/:id", machineHandler.UpdateMachine)
			//     machine.DELETE("/:id", machineHandler.DeleteMachine)
			//     machine.GET("/search", machineHandler.SearchMachines)
			//     machine.POST("/orders", machineHandler.CreateOrder)
			//     machine.GET("/orders", machineHandler.GetUserOrders)
			//     machine.PUT("/orders/:id/confirm", machineHandler.ConfirmOrder)
			//     machine.PUT("/orders/:id/pay", machineHandler.PayOrder)
			//     machine.PUT("/orders/:id/complete", machineHandler.CompleteOrder)
			//     machine.PUT("/orders/:id/cancel", machineHandler.CancelOrder)
			//     machine.POST("/orders/:id/rate", machineHandler.RateOrder)
			// }

			// TODO: 文件上传API
			// file := user.Group("/file")
			// {
			//     file.POST("/upload", fileHandler.UploadFile)
			//     file.POST("/upload/multiple", fileHandler.UploadMultipleFiles)
			//     file.GET("/:id", fileHandler.GetFile)
			//     file.DELETE("/:id", fileHandler.DeleteFile)
			//     file.GET("/:id/download", fileHandler.DownloadFile)
			// }
		}

		// 公共内容API（可选认证）
		content := api.Group("/content")
		content.Use(authMiddleware.OptionalAuth())
		{
			// TODO: 文章相关API
			// content.GET("/articles", articleHandler.ListArticles)
			// content.GET("/articles/:id", articleHandler.GetArticle)
			// content.GET("/articles/:id/view", articleHandler.ViewArticle)
			// content.POST("/articles/:id/like", articleHandler.LikeArticle)
			// content.POST("/articles/:id/share", articleHandler.ShareArticle)
			// content.GET("/categories", articleHandler.ListCategories)

			// TODO: 专家相关API
			// content.GET("/experts", expertHandler.ListExperts)
			// content.GET("/experts/:id", expertHandler.GetExpert)
			// content.POST("/experts/:id/consult", expertHandler.RequestConsultation)
			// content.POST("/experts/:id/rate", expertHandler.RateExpert)
		}

		// 管理员API（需要管理员认证）
		admin := api.Group("/admin")
		admin.Use(authMiddleware.AdminAuth())
		{
			// 用户管理
			adminUser := admin.Group("/users")
			{
				adminUser.GET("", userHandler.ListUsers)
				adminUser.GET("/statistics", userHandler.GetUserStatistics)
				// TODO: 添加其他管理员用户接口
				// adminUser.PUT("/:id/freeze", userHandler.FreezeUser)
				// adminUser.PUT("/:id/unfreeze", userHandler.UnfreezeUser)
				// adminUser.GET("/:id/auth", userHandler.GetUserAuth)
				// adminUser.PUT("/auth/:id/review", userHandler.ReviewAuth)
			}

			// TODO: 贷款管理
			// adminLoan := admin.Group("/loan")
			// {
			//     adminLoan.POST("/products", loanHandler.CreateProduct)
			//     adminLoan.PUT("/products/:id", loanHandler.UpdateProduct)
			//     adminLoan.DELETE("/products/:id", loanHandler.DeleteProduct)
			//     adminLoan.GET("/applications", loanHandler.ListApplications)
			//     adminLoan.PUT("/applications/:id/approve", loanHandler.ApproveApplication)
			//     adminLoan.PUT("/applications/:id/reject", loanHandler.RejectApplication)
			//     adminLoan.PUT("/applications/:id/return", loanHandler.ReturnApplication)
			//     adminLoan.GET("/statistics", loanHandler.GetLoanStatistics)
			//     adminLoan.GET("/reports", loanHandler.GenerateReports)
			// }

			// TODO: 农机管理
			// adminMachine := admin.Group("/machine")
			// {
			//     adminMachine.GET("/list", machineHandler.ListAllMachines)
			//     adminMachine.GET("/statistics", machineHandler.GetMachineStatistics)
			//     adminMachine.GET("/orders", machineHandler.ListAllOrders)
			// }

			// TODO: 内容管理
			// adminContent := admin.Group("/content")
			// {
			//     adminContent.POST("/articles", articleHandler.CreateArticle)
			//     adminContent.PUT("/articles/:id", articleHandler.UpdateArticle)
			//     adminContent.DELETE("/articles/:id", articleHandler.DeleteArticle)
			//     adminContent.PUT("/articles/:id/publish", articleHandler.PublishArticle)
			//     adminContent.POST("/categories", articleHandler.CreateCategory)
			//     adminContent.PUT("/categories/:id", articleHandler.UpdateCategory)
			//     adminContent.DELETE("/categories/:id", articleHandler.DeleteCategory)
			// }

			// TODO: 专家管理
			// adminExpert := admin.Group("/expert")
			// {
			//     adminExpert.POST("/register", expertHandler.RegisterExpert)
			//     adminExpert.PUT("/:id", expertHandler.UpdateExpert)
			//     adminExpert.PUT("/:id/verify", expertHandler.VerifyExpert)
			//     adminExpert.GET("/list", expertHandler.ListAllExperts)
			// }

			// TODO: 系统管理
			// adminSystem := admin.Group("/system")
			// {
			//     adminSystem.GET("/config", systemHandler.GetConfigs)
			//     adminSystem.PUT("/config", systemHandler.SetConfig)
			//     adminSystem.GET("/logs", systemHandler.GetAPILogs)
			//     adminSystem.GET("/statistics", systemHandler.GetAPIStatistics)
			//     adminSystem.DELETE("/logs/cleanup", systemHandler.CleanupLogs)
			//     adminSystem.GET("/health", systemHandler.HealthCheck)
			//     adminSystem.GET("/status", systemHandler.GetSystemStatus)
			//     adminSystem.GET("/offline-queue", systemHandler.GetOfflineQueue)
			//     adminSystem.POST("/offline-queue/process", systemHandler.ProcessOfflineQueue)
			// }
		}

		// OA后台API（需要OA认证）
		oa := api.Group("/oa")
		oa.Use(authMiddleware.AdminAuth()) // TODO: 改为OA专用的认证中间件
		{
			// TODO: OA用户管理
			// oaUser := oa.Group("/users")
			// {
			//     oaUser.POST("", oaHandler.CreateOAUser)
			//     oaUser.GET("", oaHandler.ListOAUsers)
			//     oaUser.GET("/:id", oaHandler.GetOAUser)
			//     oaUser.PUT("/:id", oaHandler.UpdateOAUser)
			//     oaUser.DELETE("/:id", oaHandler.DeleteOAUser)
			//     oaUser.PUT("/:id/reset-password", oaHandler.ResetPassword)
			// }

			// TODO: OA角色管理
			// oaRole := oa.Group("/roles")
			// {
			//     oaRole.POST("", oaHandler.CreateRole)
			//     oaRole.GET("", oaHandler.ListRoles)
			//     oaRole.GET("/:id", oaHandler.GetRole)
			//     oaRole.PUT("/:id", oaHandler.UpdateRole)
			//     oaRole.DELETE("/:id", oaHandler.DeleteRole)
			// }

			// TODO: OA工作台
			// oa.GET("/dashboard", oaHandler.GetDashboard)
			// oa.GET("/tasks", oaHandler.GetWorkTasks)
		}
	}

	// Swagger文档（仅在开发环境）
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

// SetupAPIRoutes 设置API路由（用于测试）
func SetupAPIRoutes(r *gin.Engine, config *RouterConfig) {
	authMiddleware := middleware.NewAuthMiddleware(config.JWTSecret)
	userHandler := handler.NewUserHandler(config.UserService)

	api := r.Group("/api")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/refresh", userHandler.RefreshToken)
		}

		// 用户路由
		user := api.Group("/user")
		user.Use(authMiddleware.RequireAuth())
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.PUT("/password", userHandler.ChangePassword)
			user.POST("/logout", userHandler.Logout)
		}
	}
}
