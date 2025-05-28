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
	adminService := service.NewAdminService(r.data, r.jwtManager, r.log)

	// 创建处理器实例
	userHandler := NewUserHandler(userService, r.log)
	loanHandler := NewLoanHandler(loanService, r.log)
	adminHandler := NewAdminHandler(adminService, loanService, r.log)

	// 创建认证中间件
	authMiddleware := AuthMiddleware(r.jwtManager)
	adminAuthMiddleware := AdminAuthMiddleware(r.jwtManager)

	// 注册路由
	r.registerUserRoutes(v1, userHandler, authMiddleware)
	r.registerLoanRoutes(v1, loanHandler, authMiddleware)
	r.registerFileRoutes(v1, authMiddleware)
	r.registerAdminRoutes(v1, adminHandler, adminAuthMiddleware)

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
func (r *Router) registerAdminRoutes(v1 *gin.RouterGroup, adminHandler *AdminHandler, adminAuthMiddleware gin.HandlerFunc) {
	admin := v1.Group("/admin")
	RegisterAdminRoutes(admin, adminHandler, adminAuthMiddleware)
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
