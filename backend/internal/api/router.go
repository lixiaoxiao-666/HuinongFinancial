package api

import (
	"backend/internal/conf"
	"backend/internal/data"
	"backend/internal/service"
	"backend/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Router API路由器
type Router struct {
	config     *conf.Config
	data       *data.Data
	jwtManager *pkg.JWTManager
	log        *zap.Logger
}

// NewRouter 创建路由器
func NewRouter(config *conf.Config, data *data.Data, jwtManager *pkg.JWTManager, log *zap.Logger) *Router {
	return &Router{
		config:     config,
		data:       data,
		jwtManager: jwtManager,
		log:        log,
	}
}

// SetupRoutes 设置路由
func (r *Router) SetupRoutes() *gin.Engine {
	// 设置Gin模式
	gin.SetMode(r.config.Server.Mode)

	// 创建Gin引擎
	engine := gin.New()

	// 添加中间件
	engine.Use(gin.Recovery())
	engine.Use(CORSMiddleware())
	engine.Use(RequestLoggerMiddleware(r.log))
	engine.Use(ErrorHandlerMiddleware(r.log))
	engine.Use(RateLimitMiddleware())

	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		pkg.Success(c, map[string]string{
			"status":  "ok",
			"service": "digital-agriculture-backend",
			"version": r.config.App.Version,
		})
	})

	// API v1 路由组
	v1 := engine.Group("/api/v1")

	// 创建服务实例
	userService := service.NewUserService(r.data, r.jwtManager, r.log)
	loanService := service.NewLoanService(r.data, r.log)
	fileService := service.NewFileService(r.data, r.log)
	adminService := service.NewAdminService(r.data, r.jwtManager, r.log)
	aiAgentService := service.NewAIAgentService(r.data, r.log)
	leasingService := service.NewMachineryLeasingApprovalService(r.data, r.log)

	// 🔥 关键：创建统一处理器，整合所有服务
	unifiedProcessor := service.NewUnifiedApplicationProcessor(
		r.data,
		loanService,
		leasingService,
		aiAgentService,
		r.log,
	)

	// 实例化所有处理器
	userHandler := NewUserHandler(userService, r.log)
	loanHandler := NewLoanHandler(loanService, r.log)
	fileHandler := NewFileHandler(fileService, r.log)
	adminHandler := NewAdminHandler(adminService, loanService, r.log)
	aiAgentHandler := NewAIAgentHandler(aiAgentService, unifiedProcessor, r.log)
	machineryLeasingApprovalHandler := NewMachineryLeasingApprovalHandler(leasingService, r.log)

	// 创建认证中间件
	authMiddleware := AuthMiddleware(r.jwtManager)
	adminAuthMiddleware := AdminAuthMiddleware(r.jwtManager)
	aiAgentAuthMiddleware := AIAgentAuthMiddleware(&r.config.AI)

	// 注册路由
	r.registerUserRoutes(v1, userHandler, authMiddleware)
	r.registerLoanRoutes(v1, loanHandler, authMiddleware)
	r.registerFileRoutes(v1, fileHandler, authMiddleware)
	r.registerAdminRoutes(v1, adminHandler, adminAuthMiddleware)
	r.registerAIAgentRoutes(v1, aiAgentHandler, aiAgentAuthMiddleware)
	r.registerMachineryLeasingApprovalRoutes(v1, machineryLeasingApprovalHandler, authMiddleware)

	return engine
}

// registerUserRoutes 注册用户路由
func (r *Router) registerUserRoutes(v1 *gin.RouterGroup, handler *UserHandler, authMiddleware gin.HandlerFunc) {
	users := v1.Group("/users")
	RegisterUserRoutes(users, handler, authMiddleware)
}

// registerLoanRoutes 注册贷款路由
func (r *Router) registerLoanRoutes(v1 *gin.RouterGroup, handler *LoanHandler, authMiddleware gin.HandlerFunc) {
	loans := v1.Group("/loans")
	RegisterLoanRoutes(loans, handler, authMiddleware)
}

// registerFileRoutes 注册文件路由
func (r *Router) registerFileRoutes(v1 *gin.RouterGroup, handler *FileHandler, authMiddleware gin.HandlerFunc) {
	files := v1.Group("/files")
	files.Use(authMiddleware)
	{
		files.POST("/upload", r.handleFileUpload)
	}
}

// registerAdminRoutes 注册管理员路由
func (r *Router) registerAdminRoutes(v1 *gin.RouterGroup, adminHandler *AdminHandler, adminAuthMiddleware gin.HandlerFunc) {
	admin := v1.Group("/admin")
	RegisterAdminRoutes(admin, adminHandler, adminAuthMiddleware)
}

// registerAIAgentRoutes 注册AI智能体路由
func (r *Router) registerAIAgentRoutes(v1 *gin.RouterGroup, aiAgentHandler *AIAgentHandler, aiAgentAuthMiddleware gin.HandlerFunc) {
	RegisterAIAgentRoutes(v1, aiAgentHandler, aiAgentAuthMiddleware)
}

// registerMachineryLeasingApprovalRoutes 注册农机租赁审批路由
func (r *Router) registerMachineryLeasingApprovalRoutes(v1 *gin.RouterGroup, handler *MachineryLeasingApprovalHandler, authMiddleware gin.HandlerFunc) {
	leasingApprovals := v1.Group("/machinery-leasing-approvals")
	RegisterMachineryLeasingApprovalRoutes(leasingApprovals, handler, authMiddleware)
}

// handleFileUpload 文件上传处理器（简化实现）
func (r *Router) handleFileUpload(c *gin.Context) {
	userID, _ := c.Get("user_id")
	purpose := c.PostForm("purpose")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		pkg.BadRequest(c, "文件上传失败")
		return
	}
	defer file.Close()

	// 这里应该实现真正的文件上传逻辑
	fileID := pkg.GenerateFileID()
	fileURL := "https://example.com/files/" + fileID + ".jpg"

	// 保存文件记录到数据库
	uploadedFile := data.UploadedFile{
		FileID:      fileID,
		UserID:      userID.(string),
		FileName:    header.Filename,
		FileType:    header.Header.Get("Content-Type"),
		FileSize:    header.Size,
		StoragePath: fileURL,
		Purpose:     purpose,
	}

	if err := r.data.DB.Create(&uploadedFile).Error; err != nil {
		r.log.Error("保存文件记录失败", zap.Error(err))
		pkg.InternalError(c, "文件上传失败")
		return
	}

	pkg.Success(c, map[string]interface{}{
		"file_id":   fileID,
		"file_url":  fileURL,
		"file_name": header.Filename,
		"file_size": header.Size,
	})
}
