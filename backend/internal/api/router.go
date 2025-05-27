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

	// 创建处理器实例
	userHandler := NewUserHandler(userService, r.log)
	loanHandler := NewLoanHandler(loanService, r.log)

	// 创建认证中间件
	authMiddleware := AuthMiddleware(r.jwtManager)
	adminAuthMiddleware := AdminAuthMiddleware(r.jwtManager)

	// 注册路由
	r.registerUserRoutes(v1, userHandler, authMiddleware)
	r.registerLoanRoutes(v1, loanHandler, authMiddleware)
	r.registerFileRoutes(v1, authMiddleware)
	r.registerAdminRoutes(v1, adminAuthMiddleware)

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
func (r *Router) registerFileRoutes(v1 *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	files := v1.Group("/files")
	files.Use(authMiddleware)
	{
		files.POST("/upload", r.handleFileUpload)
	}
}

// registerAdminRoutes 注册管理员路由
func (r *Router) registerAdminRoutes(v1 *gin.RouterGroup, adminAuthMiddleware gin.HandlerFunc) {
	admin := v1.Group("/admin")
	admin.Use(adminAuthMiddleware)
	{
		// OA用户登录（不需要管理员认证）
		v1.POST("/admin/login", r.handleAdminLogin)

		// 需要管理员认证的路由
		admin.GET("/loans/applications/pending", r.handleGetPendingApplications)
		admin.GET("/loans/applications/:application_id", r.handleGetApplicationDetail)
		admin.POST("/loans/applications/:application_id/review", r.handleReviewApplication)
		admin.POST("/system/ai-approval/toggle", r.handleToggleAIApproval)
	}
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

// handleAdminLogin OA用户登录（简化实现）
func (r *Router) handleAdminLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ValidationError(c, err.Error())
		return
	}

	// 这里应该验证OA用户密码
	// 简化实现，直接生成token
	token, err := r.jwtManager.GenerateToken("oa_admin001", "oa_user", "ADMIN")
	if err != nil {
		pkg.InternalError(c, "登录失败")
		return
	}

	pkg.Success(c, map[string]interface{}{
		"admin_user_id": "oa_admin001",
		"username":      req.Username,
		"role":          "ADMIN",
		"token":         token,
		"expires_in":    3600,
	})
}

// 以下是管理员相关的处理器（简化实现）
func (r *Router) handleGetPendingApplications(c *gin.Context) {
	pkg.Success(c, []interface{}{})
}

func (r *Router) handleGetApplicationDetail(c *gin.Context) {
	applicationID := c.Param("application_id")
	pkg.Success(c, map[string]string{"application_id": applicationID})
}

func (r *Router) handleReviewApplication(c *gin.Context) {
	pkg.SuccessWithMessage(c, "审批决策提交成功", nil)
}

func (r *Router) handleToggleAIApproval(c *gin.Context) {
	pkg.SuccessWithMessage(c, "AI审批状态更新成功", nil)
}
