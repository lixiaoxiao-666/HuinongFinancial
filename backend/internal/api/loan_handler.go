package api

import (
	"strconv"

	"backend/internal/service"
	"backend/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoanHandler 贷款处理器
type LoanHandler struct {
	loanService *service.LoanService
	log         *zap.Logger
}

// NewLoanHandler 创建贷款处理器
func NewLoanHandler(loanService *service.LoanService, log *zap.Logger) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
		log:         log,
	}
}

// GetLoanProducts 获取贷款产品列表
func (h *LoanHandler) GetLoanProducts(c *gin.Context) {
	category := c.Query("category")

	products, err := h.loanService.GetLoanProducts(c.Request.Context(), category)
	if err != nil {
		h.log.Error("获取贷款产品失败", zap.Error(err))
		pkg.InternalError(c, err.Error())
		return
	}

	pkg.Success(c, products)
}

// GetLoanProduct 获取贷款产品详情
func (h *LoanHandler) GetLoanProduct(c *gin.Context) {
	productID := c.Param("product_id")
	if productID == "" {
		pkg.BadRequest(c, "产品ID不能为空")
		return
	}

	product, err := h.loanService.GetLoanProduct(c.Request.Context(), productID)
	if err != nil {
		h.log.Error("获取贷款产品详情失败", zap.Error(err))
		pkg.NotFound(c, err.Error())
		return
	}

	pkg.Success(c, product)
}

// SubmitLoanApplication 提交贷款申请
func (h *LoanHandler) SubmitLoanApplication(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		pkg.Unauthorized(c, "用户信息不存在")
		return
	}

	var req service.SubmitLoanApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ValidationError(c, err.Error())
		return
	}

	result, err := h.loanService.SubmitLoanApplication(c.Request.Context(), userID.(string), &req)
	if err != nil {
		h.log.Error("提交贷款申请失败", zap.Error(err))
		pkg.BadRequest(c, err.Error())
		return
	}

	pkg.Created(c, result)
}

// GetLoanApplication 获取贷款申请详情
func (h *LoanHandler) GetLoanApplication(c *gin.Context) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		pkg.BadRequest(c, "申请ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	var userIDStr string
	if exists {
		userIDStr = userID.(string)
	}

	result, err := h.loanService.GetLoanApplication(c.Request.Context(), applicationID, userIDStr)
	if err != nil {
		h.log.Error("获取贷款申请详情失败", zap.Error(err))
		pkg.NotFound(c, err.Error())
		return
	}

	pkg.Success(c, result)
}

// GetMyLoanApplications 获取我的贷款申请列表
func (h *LoanHandler) GetMyLoanApplications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		pkg.Unauthorized(c, "用户信息不存在")
		return
	}

	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	applications, total, err := h.loanService.GetMyLoanApplications(
		c.Request.Context(),
		userID.(string),
		status,
		page,
		limit,
	)
	if err != nil {
		h.log.Error("获取贷款申请列表失败", zap.Error(err))
		pkg.InternalError(c, err.Error())
		return
	}

	pkg.ListSuccess(c, applications, total)
}

// RegisterLoanRoutes 注册贷款路由
func RegisterLoanRoutes(r *gin.RouterGroup, handler *LoanHandler, authMiddleware gin.HandlerFunc) {
	// 公开路由（产品信息）
	r.GET("/products", handler.GetLoanProducts)
	r.GET("/products/:product_id", handler.GetLoanProduct)

	// 需要认证的路由
	auth := r.Group("", authMiddleware)
	{
		auth.POST("/applications", handler.SubmitLoanApplication)
		auth.GET("/applications/my", handler.GetMyLoanApplications)
		auth.GET("/applications/:application_id", handler.GetLoanApplication)
	}
}
