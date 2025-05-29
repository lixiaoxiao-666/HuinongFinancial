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

// RegisterLoanRoutes 注册贷款路由
func RegisterLoanRoutes(group *gin.RouterGroup, handler *LoanHandler, authMiddleware gin.HandlerFunc) {
	// 无需认证的路由
	group.GET("/products", handler.GetLoanProducts)
	group.GET("/products/:product_id", handler.GetLoanProductDetail)

	// 需要认证的路由
	authenticated := group.Group("")
	authenticated.Use(authMiddleware)
	{
		authenticated.POST("/applications", handler.SubmitLoanApplication)
		authenticated.GET("/applications/:application_id", handler.GetLoanApplicationDetail)
		authenticated.GET("/applications/my", handler.GetMyLoanApplications)
	}
}

// GetLoanProducts 获取贷款产品列表
// @Summary 获取贷款产品列表
// @Description 获取可申请的贷款产品列表，支持按分类筛选
// @Tags 贷款服务
// @Produce json
// @Param category query string false "产品分类"
// @Success 200 {object} CommonResponse{data=[]LoanProduct}
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/loans/products [get]
func (h *LoanHandler) GetLoanProducts(c *gin.Context) {
	category := c.Query("category")

	// 调用服务层获取贷款产品列表
	products, err := h.loanService.GetLoanProducts(c.Request.Context(), category)
	if err != nil {
		h.log.Error("获取贷款产品列表失败", zap.Error(err))
		pkg.InternalError(c, "获取产品列表失败，请稍后重试")
		return
	}

	// 转换为API响应格式
	var responseProducts []LoanProduct
	for _, product := range products {
		responseProducts = append(responseProducts, LoanProduct{
			ProductID:          product.ProductID,
			Name:               product.Name,
			Description:        product.Description,
			Category:           product.Category,
			MinAmount:          int64(product.MinAmount),
			MaxAmount:          int64(product.MaxAmount),
			MinTermMonths:      product.MinTermMonths,
			MaxTermMonths:      product.MaxTermMonths,
			InterestRateYearly: product.InterestRateYearly,
			Status:             0, // 默认状态为有效
		})
	}

	pkg.Success(c, responseProducts)
}

// GetLoanProductDetail 获取贷款产品详情
// @Summary 获取贷款产品详情
// @Description 获取指定贷款产品的详细信息
// @Tags 贷款服务
// @Produce json
// @Param product_id path string true "产品ID"
// @Success 200 {object} CommonResponse{data=LoanProduct}
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/loans/products/{product_id} [get]
func (h *LoanHandler) GetLoanProductDetail(c *gin.Context) {
	productID := c.Param("product_id")
	if productID == "" {
		pkg.BadRequest(c, "产品ID不能为空")
		return
	}

	// 调用服务层获取产品详情
	product, err := h.loanService.GetLoanProduct(c.Request.Context(), productID)
	if err != nil {
		h.log.Error("获取贷款产品详情失败", zap.Error(err), zap.String("product_id", productID))

		switch err.Error() {
		case "产品不存在":
			pkg.NotFound(c, "产品不存在")
		default:
			pkg.InternalError(c, "获取产品详情失败，请稍后重试")
		}
		return
	}

	// 转换必需文档格式
	var requiredDocs []RequiredDocument
	if product.RequiredDocuments != nil {
		if docs, ok := product.RequiredDocuments.([]interface{}); ok {
			for _, doc := range docs {
				if docMap, ok := doc.(map[string]interface{}); ok {
					requiredDocs = append(requiredDocs, RequiredDocument{
						Type: getString(docMap, "type"),
						Desc: getString(docMap, "desc"),
					})
				}
			}
		}
	}

	// 转换还款方式
	var repaymentMethods []string
	if product.RepaymentMethods != nil {
		if methods, ok := product.RepaymentMethods.([]interface{}); ok {
			for _, method := range methods {
				if methodStr, ok := method.(string); ok {
					repaymentMethods = append(repaymentMethods, methodStr)
				}
			}
		}
	}

	// 返回详情响应
	pkg.Success(c, LoanProduct{
		ProductID:             product.ProductID,
		Name:                  product.Name,
		Description:           product.Description,
		Category:              product.Category,
		MinAmount:             int64(product.MinAmount),
		MaxAmount:             int64(product.MaxAmount),
		MinTermMonths:         product.MinTermMonths,
		MaxTermMonths:         product.MaxTermMonths,
		InterestRateYearly:    product.InterestRateYearly,
		RepaymentMethods:      repaymentMethods,
		ApplicationConditions: product.ApplicationConditions,
		RequiredDocuments:     requiredDocs,
		Status:                0, // 默认状态为有效
	})
}

// SubmitLoanApplication 提交贷款申请
// @Summary 提交贷款申请
// @Description 用户提交贷款申请
// @Tags 贷款服务
// @Accept json
// @Produce json
// @Param request body LoanApplicationRequest true "申请请求"
// @Success 201 {object} CommonResponse{data=LoanApplicationResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/loans/applications [post]
// @Security BearerAuth
func (h *LoanHandler) SubmitLoanApplication(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		h.log.Warn("提交贷款申请时未找到用户ID")
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	var req LoanApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("提交贷款申请参数绑定失败", zap.Error(err))
		pkg.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 转换上传文档格式
	var uploadedDocs []service.DocumentInfo
	for _, doc := range req.UploadedDocuments {
		uploadedDocs = append(uploadedDocs, service.DocumentInfo{
			DocType: doc.DocType,
			FileID:  doc.FileID,
		})
	}

	// 构建申请人信息
	applicantInfo := map[string]interface{}{
		"real_name":      req.ApplicantInfo.RealName,
		"id_card_number": req.ApplicantInfo.IDCardNumber,
		"address":        req.ApplicantInfo.Address,
	}

	// 构建服务层请求
	serviceReq := &service.SubmitLoanApplicationRequest{
		ProductID:         req.ProductID,
		Amount:            float64(req.Amount),
		TermMonths:        req.TermMonths,
		Purpose:           req.Purpose,
		ApplicantInfo:     applicantInfo,
		UploadedDocuments: uploadedDocs,
	}

	// 调用服务层提交申请
	result, err := h.loanService.SubmitLoanApplication(c.Request.Context(), userID.(string), serviceReq)
	if err != nil {
		h.log.Error("提交贷款申请失败", zap.Error(err), zap.String("user_id", userID.(string)))

		switch err.Error() {
		case "贷款产品不存在":
			pkg.BadRequest(c, "贷款产品不存在")
		case "申请金额超出产品范围":
			pkg.BadRequest(c, "申请金额超出产品限制")
		case "申请期限超出产品范围":
			pkg.BadRequest(c, "申请期限超出产品限制")
		case "用户已有进行中的申请":
			pkg.BadRequest(c, "您已有进行中的贷款申请，请等待审核完成后再次申请")
		default:
			pkg.InternalError(c, "提交申请失败，请稍后重试")
		}
		return
	}

	// 记录成功日志
	h.log.Info("贷款申请提交成功",
		zap.String("application_id", result.ApplicationID),
		zap.String("user_id", userID.(string)),
		zap.String("product_id", req.ProductID),
		zap.Int64("amount", req.Amount))

	// 返回成功响应
	c.JSON(201, CommonResponse{
		Code:    CodeSuccess,
		Message: "Success",
		Data: LoanApplicationResponse{
			ApplicationID: result.ApplicationID,
		},
	})
}

// GetLoanApplicationDetail 获取贷款申请详情
// @Summary 获取贷款申请详情
// @Description 获取指定贷款申请的详细信息
// @Tags 贷款服务
// @Produce json
// @Param application_id path string true "申请ID"
// @Success 200 {object} CommonResponse{data=LoanApplicationDetail}
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/loans/applications/{application_id} [get]
// @Security BearerAuth
func (h *LoanHandler) GetLoanApplicationDetail(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		h.log.Warn("获取贷款申请详情时未找到用户ID")
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	applicationID := c.Param("application_id")
	if applicationID == "" {
		pkg.BadRequest(c, "申请ID不能为空")
		return
	}

	// 调用服务层获取申请详情
	application, err := h.loanService.GetLoanApplication(c.Request.Context(), applicationID, userID.(string))
	if err != nil {
		h.log.Error("获取贷款申请详情失败", zap.Error(err),
			zap.String("application_id", applicationID),
			zap.String("user_id", userID.(string)))

		switch err.Error() {
		case "申请不存在":
			pkg.NotFound(c, "申请不存在")
		case "无权限查看此申请":
			pkg.Forbidden(c, "无权限查看此申请")
		default:
			pkg.InternalError(c, "获取申请详情失败，请稍后重试")
		}
		return
	}

	// 转换申请历史格式
	var history []ApplicationHistory
	for _, h := range application.History {
		history = append(history, ApplicationHistory{
			Status:    h.Status,
			Timestamp: h.Timestamp,
			Operator:  h.Operator,
		})
	}

	// 转换响应格式
	var approvedAmount *int64
	if application.ApprovedAmount != nil {
		amount := int64(*application.ApprovedAmount)
		approvedAmount = &amount
	}

	// 返回详情响应
	pkg.Success(c, LoanApplicationDetail{
		ApplicationID:  application.ApplicationID,
		ProductID:      application.ProductID,
		UserID:         application.UserID,
		Amount:         int64(application.AmountApplied),
		TermMonths:     application.TermMonthsApplied,
		Purpose:        application.Purpose,
		Status:         application.Status,
		SubmittedAt:    application.SubmittedAt,
		UpdatedAt:      application.SubmittedAt, // 如果没有UpdatedAt，使用SubmittedAt
		ApprovedAmount: approvedAmount,
		Remarks:        application.AISuggestion,
		History:        history,
	})
}

// GetMyLoanApplications 获取我的贷款申请列表
// @Summary 获取我的贷款申请列表
// @Description 获取当前用户的贷款申请列表，支持状态筛选和分页
// @Tags 贷款服务
// @Produce json
// @Param status query string false "申请状态"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} PaginationResponse{data=[]MyLoanApplication}
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/loans/applications/my [get]
// @Security BearerAuth
func (h *LoanHandler) GetMyLoanApplications(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		h.log.Warn("获取我的贷款申请列表时未找到用户ID")
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	// 解析查询参数
	status := c.Query("status")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// 调用服务层获取申请列表
	applications, total, err := h.loanService.GetMyLoanApplications(c.Request.Context(), userID.(string), status, page, limit)
	if err != nil {
		h.log.Error("获取我的贷款申请列表失败", zap.Error(err), zap.String("user_id", userID.(string)))
		pkg.InternalError(c, "获取申请列表失败，请稍后重试")
		return
	}

	// 转换响应格式
	var responseApplications []MyLoanApplication
	for _, app := range applications {
		responseApplications = append(responseApplications, MyLoanApplication{
			ApplicationID: app.ApplicationID,
			ProductName:   "产品名称", // TODO: 从产品服务获取真实产品名称
			Amount:        int64(app.AmountApplied),
			Status:        app.Status,
			SubmittedAt:   app.SubmittedAt,
		})
	}

	// 返回分页响应
	c.JSON(200, PaginationResponse{
		Code:    CodeSuccess,
		Message: "Success",
		Data:    responseApplications,
		Total:   total,
	})
}

// 辅助函数：从map中安全获取字符串值
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
