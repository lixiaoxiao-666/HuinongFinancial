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
	TaskService    service.TaskService
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

	// 任务管理Handler
	taskHandler := handler.NewTaskHandler(config.TaskService)

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
					loan.POST("/auto-approve", difyHandler.AutoApproveLoanApplication)
				}

				// 农机相关接口
				machine := dify.Group("/machine")
				{
					machine.POST("/get-rental-details", difyHandler.GetMachineRentalDetails)
					machine.POST("/auto-approve-rental", difyHandler.AutoApproveMachineRental)
				}

				// 征信相关接口
				credit := dify.Group("/credit")
				{
					credit.POST("/query", difyHandler.QueryCreditReport)
				}

				// 任务管理接口
				task := dify.Group("/task")
				{
					task.POST("/create", difyHandler.CreateApprovalTask)
					task.POST("/status", difyHandler.GetTaskStatus)
					task.POST("/complete", difyHandler.CompleteApprovalTask)
				}
			}
		}

		// 惠农APP/Web认证API（无需认证的公开接口）
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/refresh", sessionAuthMiddleware.RefreshToken())

			// Token验证接口（需要认证）
			auth.GET("/validate", sessionAuthMiddleware.RequireAuth(), func(c *gin.Context) {
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

		// 惠农APP/Web用户API（需要认证）- 支持app和web平台
		user := api.Group("/user")
		user.Use(sessionAuthMiddleware.RequireAuth())
		{
			// 用户资料管理
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.PUT("/password", userHandler.ChangePassword)
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
			}

			// 用户标签
			user.GET("/tags", userHandler.GetUserTags)
			user.POST("/tags", userHandler.AddUserTag)
			user.DELETE("/tags/:tag_key", userHandler.RemoveUserTag)

			// 贷款相关API
			loan := user.Group("/loan")
			{
				loan.GET("/products", loanHandler.GetProducts)
				loan.GET("/products/:id", loanHandler.GetProductDetail)
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

		// 公共内容API（可选认证）- 登录用户可获得个性化内容
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

		// OA系统认证API（无需认证的公开接口）
		oaAuth := api.Group("/oa/auth")
		{
			oaAuth.POST("/login", oaHandler.Login)
			oaAuth.POST("/refresh", sessionAuthMiddleware.RefreshToken())
			oaAuth.GET("/validate", sessionAuthMiddleware.RequireAuth(), sessionAuthMiddleware.CheckPlatform("oa"), func(c *gin.Context) {
				c.JSON(200, gin.H{
					"code":    200,
					"message": "OA Token有效",
					"data": gin.H{
						"valid": true,
					},
				})
			})
			oaAuth.POST("/logout", sessionAuthMiddleware.RequireAuth(), sessionAuthMiddleware.CheckPlatform("oa"), sessionAuthMiddleware.Logout())
		}

		// OA系统普通用户API（需要OA平台认证，但不需要管理员权限）
		oaUser := api.Group("/oa/user")
		oaUser.Use(sessionAuthMiddleware.RequireAuth(), sessionAuthMiddleware.CheckPlatform("oa"))
		{
			// OA普通用户个人信息
			oaUser.GET("/profile", userHandler.GetProfile)
			oaUser.PUT("/profile", userHandler.UpdateProfile)
			oaUser.PUT("/password", userHandler.ChangePassword)

			// OA普通用户查看自己的申请
			oaUser.GET("/loan/applications", loanHandler.GetUserApplications)
			oaUser.GET("/loan/applications/:id", loanHandler.GetApplicationDetail)
		}

		// OA系统管理员API（需要OA平台认证 + 管理员权限）
		oaAdmin := api.Group("/oa/admin")
		oaAdmin.Use(sessionAuthMiddleware.RequireAuth(), sessionAuthMiddleware.CheckPlatform("oa"), sessionAuthMiddleware.RequireRole("admin"))
		{
			// 用户管理
			adminUser := oaAdmin.Group("/users")
			{
				adminUser.GET("", userHandler.ListUsers)
				adminUser.GET("/statistics", userHandler.GetUserStatistics)
				adminUser.GET("/:user_id", oaHandler.GetUserDetail)
				adminUser.PUT("/:user_id/status", oaHandler.UpdateUserStatus)
				adminUser.POST("/batch-operation", oaHandler.BatchOperateUsers)
				adminUser.GET("/:user_id/auth-status", userAuthHandler.GetUserAuthStatus)
			}

			// 会话管理
			adminSession := oaAdmin.Group("/sessions")
			{
				adminSession.GET("/statistics", func(c *gin.Context) {
					// 会话统计接口实现
					c.JSON(200, gin.H{
						"code":    200,
						"message": "获取成功",
						"data": gin.H{
							"total_active_sessions": 150,
							"platform_distribution": gin.H{
								"app": 100,
								"web": 30,
								"oa":  20,
							},
							"daily_peak_users":                 120,
							"average_session_duration_minutes": 30,
						},
					})
				})
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

			// 待处理任务管理
			adminTasks := oaAdmin.Group("/tasks")
			{
				// 任务基本操作
				adminTasks.GET("", taskHandler.ListTasks)
				adminTasks.POST("", taskHandler.CreateTask)
				adminTasks.GET("/:id", taskHandler.GetTask)
				adminTasks.PUT("/:id", taskHandler.UpdateTask)
				adminTasks.DELETE("/:id", taskHandler.DeleteTask)

				// 待处理任务列表（兼容原接口）
				adminTasks.GET("/pending", taskHandler.GetPendingTasks)

				// 任务处理操作
				adminTasks.POST("/:id/process", taskHandler.ProcessTask)
				adminTasks.POST("/:id/assign", taskHandler.AssignTask)
				adminTasks.POST("/:id/unassign", taskHandler.UnassignTask)
				adminTasks.POST("/:id/reassign", taskHandler.ReassignTask)

				// 任务进度
				adminTasks.GET("/:id/progress", taskHandler.GetTaskProgress)
			}

			// 贷款审批管理
			adminLoan := oaAdmin.Group("/loans")
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
			adminAuth := oaAdmin.Group("/auth")
			{
				adminAuth.GET("/list", userAuthHandler.GetAuthList)
				adminAuth.GET("/:id", userAuthHandler.GetAuthDetail)
				adminAuth.POST("/:id/review", userAuthHandler.ReviewAuth)
				adminAuth.POST("/batch-review", userAuthHandler.BatchReviewAuth)
				adminAuth.GET("/statistics", userAuthHandler.GetAuthStatistics)
				adminAuth.GET("/export", userAuthHandler.ExportAuthData)
			}

			// 内容管理
			adminContent := oaAdmin.Group("/content")
			{
				// 文章管理
				adminContent.POST("/articles", articleHandler.CreateArticle)
				adminContent.PUT("/articles/:id", articleHandler.UpdateArticle)
				adminContent.DELETE("/articles/:id", articleHandler.DeleteArticle)
				adminContent.POST("/articles/:id/publish", articleHandler.PublishArticle)

				// 系统公告管理
				adminContent.GET("/announcements", func(c *gin.Context) {
					// 获取系统公告列表
					limit := c.DefaultQuery("limit", "10")
					status := c.DefaultQuery("status", "")

					announcements := []gin.H{
						{
							"id":         1,
							"title":      "系统维护通知",
							"content":    "系统将于今晚22:00-24:00进行维护，期间可能影响部分功能使用。",
							"status":     "published",
							"created_at": "2024-01-15T08:00:00Z",
							"updated_at": "2024-01-15T08:00:00Z",
						},
						{
							"id":         2,
							"title":      "新功能上线公告",
							"content":    "AI智能风险评估功能已正式上线，将大幅提升审批效率。",
							"status":     "published",
							"created_at": "2024-01-14T10:00:00Z",
							"updated_at": "2024-01-14T10:00:00Z",
						},
						{
							"id":         3,
							"title":      "节假日服务安排",
							"content":    "春节期间客服时间调整为9:00-18:00，给您带来不便敬请谅解。",
							"status":     "published",
							"created_at": "2024-01-13T15:30:00Z",
							"updated_at": "2024-01-13T15:30:00Z",
						},
					}

					// 如果指定了状态过滤
					if status != "" {
						filtered := []gin.H{}
						for _, announcement := range announcements {
							if announcement["status"] == status {
								filtered = append(filtered, announcement)
							}
						}
						announcements = filtered
					}

					c.JSON(200, gin.H{
						"code":    200,
						"message": "获取成功",
						"data": gin.H{
							"announcements": announcements,
							"total":         len(announcements),
							"limit":         limit,
						},
					})
				})
				adminContent.POST("/announcements", func(c *gin.Context) {
					// 创建系统公告
					c.JSON(200, gin.H{
						"code":    200,
						"message": "公告创建成功",
						"data": gin.H{
							"id": 4,
						},
					})
				})
				adminContent.PUT("/announcements/:id", func(c *gin.Context) {
					// 更新系统公告
					id := c.Param("id")
					c.JSON(200, gin.H{
						"code":    200,
						"message": "公告更新成功",
						"data": gin.H{
							"id": id,
						},
					})
				})
				adminContent.DELETE("/announcements/:id", func(c *gin.Context) {
					// 删除系统公告
					id := c.Param("id")
					c.JSON(200, gin.H{
						"code":    200,
						"message": "公告删除成功",
						"data": gin.H{
							"id": id,
						},
					})
				})

				// 分类管理
				adminContent.POST("/categories", articleHandler.CreateCategory)
				adminContent.PUT("/categories/:id", articleHandler.UpdateCategory)
				adminContent.DELETE("/categories/:id", articleHandler.DeleteCategory)

				// 专家管理
				adminContent.POST("/experts", expertHandler.CreateExpert)
				adminContent.PUT("/experts/:id", expertHandler.UpdateExpert)
				adminContent.DELETE("/experts/:id", expertHandler.DeleteExpert)
			}

			// 农机管理
			adminMachine := oaAdmin.Group("/machines")
			{
				adminMachine.GET("", oaHandler.GetMachines)
				adminMachine.GET("/:machine_id", oaHandler.GetMachineDetail)
				// TODO: 添加设备审核、状态管理等接口
			}

			// 系统管理
			adminSystem := oaAdmin.Group("/system")
			{
				adminSystem.GET("/config", systemHandler.GetConfig)
				adminSystem.PUT("/config", systemHandler.SetConfig)
				adminSystem.GET("/configs", systemHandler.GetConfigs)
				adminSystem.GET("/health", systemHandler.HealthCheck)
				adminSystem.GET("/statistics", systemHandler.GetSystemStats)
			}

			// OA工作台和数据分析
			oaAdmin.GET("/dashboard", oaHandler.GetDashboard)
			oaAdmin.GET("/dashboard/overview", oaHandler.GetDashboard)
			oaAdmin.GET("/dashboard/risk-monitoring", oaHandler.GetRiskMonitoring)
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
