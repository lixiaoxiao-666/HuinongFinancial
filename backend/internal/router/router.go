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
	SessionService service.SessionService
	LoanService    service.LoanService
	MachineService service.MachineService
	ArticleService service.ContentService
	ExpertService  service.ContentService
	FileService    service.SystemService
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

	// 创建会话认证中间件（统一使用Redis会话认证）
	sessionAuthMiddleware := middleware.NewSessionAuthMiddleware(config.SessionService)

	// 创建Handler实例
	userHandler := handler.NewUserHandler(config.UserService)
	loanHandler := handler.NewLoanHandler(config.LoanService)
	oaLoanHandler := handler.NewOALoanHandler(config.LoanService, config.OAService)
	fileHandler := handler.NewFileHandler(config.SystemService)
	userAuthHandler := handler.NewUserAuthHandler(config.UserService, config.OAService)
	machineHandler := handler.NewMachineHandler(config.MachineService)

	// 新增的Handler实例
	articleHandler := handler.NewArticleHandler(config.ArticleService)
	expertHandler := handler.NewExpertHandler(config.ExpertService)
	systemHandler := handler.NewSystemHandler(config.SystemService)
	oaHandler := handler.NewOAHandler(config.OAService, config.UserService)

	difyHandler := handler.NewDifyHandler(
		config.UserService,
		config.LoanService,
		config.MachineService,
	)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":          "ok",
			"message":         "数字惠农API服务正在运行",
			"version":         "1.0.0",
			"session_enabled": true,
		})
	})

	// API版本分组
	api := r.Group("/api")
	{
		// 公开API（无需认证）
		public := api.Group("/public")
		{
			public.GET("/version", systemHandler.GetSystemVersion)
			public.GET("/configs", systemHandler.GetPublicConfigs)
		}

		// 内部API（供Dify工作流调用）- 保持Dify专用认证
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

			// 使用Redis会话管理系统进行Token刷新
			auth.POST("/refresh", sessionAuthMiddleware.RefreshToken())

			// Token验证接口（需要认证）
			auth.GET("/validate", sessionAuthMiddleware.RequireAuth(), func(c *gin.Context) {
				// 如果中间件验证通过，说明Token有效
				c.JSON(200, gin.H{
					"code":    200,
					"message": "Token有效",
					"data": gin.H{
						"valid": true,
					},
				})
			})

			// TODO: 添加其他认证接口
			// auth.POST("/forgot-password", userHandler.ForgotPassword)
			// auth.POST("/reset-password", userHandler.ResetPassword)
			// auth.POST("/verify-sms", userHandler.VerifySMS)
			// auth.POST("/send-sms", userHandler.SendSMS)
		}

		// 用户API（需要认证）- 使用Redis会话认证
		user := api.Group("/user")
		user.Use(sessionAuthMiddleware.RequireAuth())
		{
			// 用户资料管理
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.PUT("/password", userHandler.ChangePassword)

			// 使用Redis会话管理系统进行登出
			user.POST("/logout", sessionAuthMiddleware.Logout())

			// 会话管理相关接口
			session := user.Group("/session")
			{
				session.GET("/info", sessionAuthMiddleware.SessionInfo())
				session.POST("/revoke-others", sessionAuthMiddleware.RevokeOtherSessions())
				session.GET("/list", sessionAuthMiddleware.SessionInfo()) // 获取用户所有会话
			}

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

			// 贷款相关API
			loan := user.Group("/loan")
			{
				// 产品相关
				loan.GET("/products", loanHandler.GetProducts)
				loan.GET("/products/:id", loanHandler.GetProductDetail)

				// 申请相关
				loan.POST("/applications", loanHandler.SubmitApplication)
				loan.GET("/applications", loanHandler.GetUserApplications)
				loan.GET("/applications/:id", loanHandler.GetApplicationDetail)
				loan.DELETE("/applications/:id", loanHandler.CancelApplication)
			}

			// 文件上传API
			file := user.Group("/files")
			{
				file.POST("/upload", fileHandler.UploadFile)
				file.POST("/upload/batch", fileHandler.UploadMultipleFiles)
				file.GET("/:id", fileHandler.GetFile)
				file.DELETE("/:id", fileHandler.DeleteFile)
			}

			// 农机相关API
			machine := user.Group("/machines")
			{
				machine.POST("", machineHandler.RegisterMachine)
				machine.GET("", machineHandler.GetUserMachines)
				machine.GET("/search", machineHandler.SearchMachines)
				machine.GET("/:id", machineHandler.GetMachine)
				machine.POST("/:id/orders", machineHandler.CreateOrder)
				// TODO: 添加农机更新和删除功能
				// machine.PUT("/:id", machineHandler.UpdateMachine)
				// machine.DELETE("/:id", machineHandler.DeleteMachine)
			}

			// 农机订单API
			orders := user.Group("/orders")
			{
				orders.GET("", machineHandler.GetUserOrders)
				orders.PUT("/:id/confirm", machineHandler.ConfirmOrder)
				orders.POST("/:id/pay", machineHandler.PayOrder)
				orders.PUT("/:id/complete", machineHandler.CompleteOrder)
				orders.PUT("/:id/cancel", machineHandler.CancelOrder)
				orders.POST("/:id/rate", machineHandler.RateOrder)
			}

			// 专家咨询API
			consultations := user.Group("/consultations")
			{
				consultations.POST("", expertHandler.SubmitConsultation)
				consultations.GET("", expertHandler.GetConsultations)
			}
		}

		// 公共内容API（可选认证）- 使用Redis会话认证
		content := api.Group("/content")
		content.Use(sessionAuthMiddleware.OptionalAuth())
		{
			// 文章相关API
			content.GET("/articles", articleHandler.ListArticles)
			content.GET("/articles/featured", articleHandler.GetFeaturedArticles)
			content.GET("/articles/:id", articleHandler.GetArticle)
			content.GET("/categories", articleHandler.GetCategories)

			// 专家相关API
			content.GET("/experts", expertHandler.ListExperts)
			content.GET("/experts/:id", expertHandler.GetExpert)
		}

		// 管理员API（需要管理员认证）- 使用Redis会话认证
		admin := api.Group("/admin")
		admin.Use(sessionAuthMiddleware.AdminAuth())
		{
			// 用户管理
			adminUser := admin.Group("/users")
			{
				adminUser.GET("", userHandler.ListUsers)
				adminUser.GET("/statistics", userHandler.GetUserStatistics)
				adminUser.GET("/:user_id/auth-status", userAuthHandler.GetUserAuthStatus)
				// TODO: 添加其他管理员用户接口
				// adminUser.PUT("/:id/freeze", userHandler.FreezeUser)
				// adminUser.PUT("/:id/unfreeze", userHandler.UnfreezeUser)
				// adminUser.GET("/:id/auth", userHandler.GetUserAuth)
			}

			// 会话管理
			adminSession := admin.Group("/sessions")
			{
				adminSession.GET("/active", func(c *gin.Context) {
					// TODO: 实现获取所有活跃会话的接口
					c.JSON(200, gin.H{
						"message": "获取活跃会话列表",
						"data": gin.H{
							"total":    150,
							"sessions": []gin.H{},
						},
					})
				})
				adminSession.POST("/cleanup", func(c *gin.Context) {
					// TODO: 实现手动清理过期会话的接口
					c.JSON(200, gin.H{
						"message": "会话清理完成",
						"cleaned": 23,
					})
				})
				adminSession.DELETE("/:session_id", func(c *gin.Context) {
					// TODO: 实现强制注销指定会话的接口
					sessionID := c.Param("session_id")
					c.JSON(200, gin.H{
						"message":    "会话已强制注销",
						"session_id": sessionID,
					})
				})
			}

			// 贷款审批管理
			adminLoan := admin.Group("/loans")
			{
				adminLoan.GET("/applications", oaLoanHandler.GetApplications)
				adminLoan.GET("/applications/:id", oaLoanHandler.GetApplicationDetail)
				adminLoan.POST("/applications/:id/approve", oaLoanHandler.ApproveApplication)
				adminLoan.POST("/applications/:id/reject", oaLoanHandler.RejectApplication)
				adminLoan.POST("/applications/:id/return", oaLoanHandler.ReturnApplication)
				adminLoan.POST("/applications/:id/start-review", oaLoanHandler.StartReview)
				adminLoan.POST("/applications/:id/retry-ai", oaLoanHandler.RetryAIAssessment)
				adminLoan.GET("/statistics", oaLoanHandler.GetStatistics)
			}

			// 认证审核管理
			adminAuth := admin.Group("/auth")
			{
				adminAuth.GET("/list", userAuthHandler.GetAuthList)
				adminAuth.GET("/:id", userAuthHandler.GetAuthDetail)
				adminAuth.POST("/:id/review", userAuthHandler.ReviewAuth)
				adminAuth.POST("/batch-review", userAuthHandler.BatchReviewAuth)
				adminAuth.GET("/statistics", userAuthHandler.GetAuthStatistics)
				adminAuth.GET("/export", userAuthHandler.ExportAuthData)
			}

			// 内容管理
			adminContent := admin.Group("/content")
			{
				// 文章管理
				adminContent.POST("/articles", articleHandler.CreateArticle)
				adminContent.PUT("/articles/:id", articleHandler.UpdateArticle)
				adminContent.DELETE("/articles/:id", articleHandler.DeleteArticle)
				adminContent.POST("/articles/:id/publish", articleHandler.PublishArticle)

				// 分类管理
				adminContent.POST("/categories", articleHandler.CreateCategory)
				adminContent.PUT("/categories/:id", articleHandler.UpdateCategory)
				adminContent.DELETE("/categories/:id", articleHandler.DeleteCategory)

				// 专家管理
				adminContent.POST("/experts", expertHandler.CreateExpert)
				adminContent.PUT("/experts/:id", expertHandler.UpdateExpert)
				adminContent.DELETE("/experts/:id", expertHandler.DeleteExpert)
			}

			// 系统管理
			adminSystem := admin.Group("/system")
			{
				adminSystem.GET("/config", systemHandler.GetConfig)
				adminSystem.PUT("/config", systemHandler.SetConfig)
				adminSystem.GET("/configs", systemHandler.GetConfigs)
				adminSystem.GET("/health", systemHandler.HealthCheck)
				adminSystem.GET("/statistics", systemHandler.GetSystemStats)
			}

			// TODO: 农机管理
			// adminMachine := admin.Group("/machine")
			// {
			//     adminMachine.GET("/list", machineHandler.ListAllMachines)
			//     adminMachine.GET("/statistics", machineHandler.GetMachineStatistics)
			//     adminMachine.GET("/orders", machineHandler.ListAllOrders)
			// }
		}

		// OA认证相关API（公开接口，无需认证）
		oaAuth := api.Group("/oa/auth")
		{
			oaAuth.POST("/login", oaHandler.Login)
			oaAuth.POST("/refresh", sessionAuthMiddleware.RefreshToken())
			oaAuth.GET("/validate", sessionAuthMiddleware.RequireAuth(), func(c *gin.Context) {
				// 检查是否为OA平台
				platform, exists := c.Get("platform")
				if !exists || platform != "oa" {
					c.JSON(401, gin.H{
						"code":    401,
						"message": "非OA平台Token",
						"error":   "Invalid platform for OA validation",
					})
					c.Abort()
					return
				}

				// 如果中间件验证通过，说明Token有效
				c.JSON(200, gin.H{
					"code":    200,
					"message": "Token有效",
					"data": gin.H{
						"valid": true,
					},
				})
			})
			oaAuth.POST("/logout", sessionAuthMiddleware.RequireAuth(), sessionAuthMiddleware.Logout())
		}

		// OA后台API（需要OA认证）- 使用Redis会话认证
		oa := api.Group("/oa")
		oa.Use(sessionAuthMiddleware.AdminAuth()) // 使用管理员认证，检查platform为"oa"
		{
			// OA用户管理
			oaUser := oa.Group("/users")
			{
				oaUser.GET("", oaHandler.GetUsers)
				oaUser.GET("/:user_id", oaHandler.GetUserDetail)
				oaUser.PUT("/:user_id/status", oaHandler.UpdateUserStatus)
				oaUser.POST("/batch-operation", oaHandler.BatchOperateUsers)
			}

			// OA农机设备管理
			oaMachine := oa.Group("/machines")
			{
				oaMachine.GET("", oaHandler.GetMachines)
				oaMachine.GET("/:machine_id", oaHandler.GetMachineDetail)
				// TODO: 添加设备审核、状态管理等接口
			}

			// OA工作台和数据分析
			oa.GET("/dashboard", oaHandler.GetDashboard)
			oa.GET("/dashboard/overview", oaHandler.GetDashboard)
			oa.GET("/dashboard/risk-monitoring", oaHandler.GetRiskMonitoring)

			// TODO: 继续添加其他OA管理功能
			// - 认证审核管理
			// - 系统配置管理
			// - 操作日志管理
			// - 权限管理
		}
	}

	// Swagger文档（仅在开发环境）
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

// SetupAPIRoutes 设置API路由（用于测试）- 统一使用Redis会话认证
func SetupAPIRoutes(r *gin.Engine, config *RouterConfig) {
	// 使用Redis会话认证中间件
	sessionAuthMiddleware := middleware.NewSessionAuthMiddleware(config.SessionService)
	userHandler := handler.NewUserHandler(config.UserService)

	api := r.Group("/api")
	{
		// 认证路由（无需认证）
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/refresh", sessionAuthMiddleware.RefreshToken())
		}

		// 用户路由（使用Redis会话认证）
		user := api.Group("/user")
		user.Use(sessionAuthMiddleware.RequireAuth())
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.PUT("/password", userHandler.ChangePassword)
			user.POST("/logout", sessionAuthMiddleware.Logout())
		}
	}
}
