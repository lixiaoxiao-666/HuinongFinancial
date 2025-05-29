package handler

import (
	"net/http"
	"strconv"

	"huinong-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// ArticleHandler 文章处理器
type ArticleHandler struct {
	contentService service.ContentService
}

// NewArticleHandler 创建文章处理器
func NewArticleHandler(contentService service.ContentService) *ArticleHandler {
	return &ArticleHandler{
		contentService: contentService,
	}
}

// CreateArticle 创建文章
// @Summary 创建文章
// @Description 管理员创建新文章
// @Tags 内容管理
// @Accept json
// @Produce json
// @Param request body service.CreateArticleRequest true "文章信息"
// @Success 200 {object} StandardResponse{data=service.CreateArticleResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/content/articles [post]
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req service.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	// 从上下文获取用户ID
	_, exists := c.Get("oaUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "未登录", "OA用户认证信息缺失"))
		return
	}

	response, err := h.contentService.CreateArticle(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "创建文章失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("文章创建成功", response))
}

// GetArticle 获取文章详情
// @Summary 获取文章详情
// @Description 获取指定文章的详细信息
// @Tags 内容管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} StandardResponse{data=service.ArticleDetailResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/content/articles/{id} [get]
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的文章ID", err.Error()))
		return
	}

	response, err := h.contentService.GetArticle(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "文章不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "文章不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取文章详情失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// UpdateArticle 更新文章
// @Summary 更新文章
// @Description 管理员更新文章信息
// @Tags 内容管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param request body service.UpdateArticleRequest true "更新信息"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/content/articles/{id} [put]
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的文章ID", err.Error()))
		return
	}

	var req service.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	err = h.contentService.UpdateArticle(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "文章不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "文章不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "更新文章失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("文章更新成功", nil))
}

// DeleteArticle 删除文章
// @Summary 删除文章
// @Description 管理员删除文章
// @Tags 内容管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/content/articles/{id} [delete]
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的文章ID", err.Error()))
		return
	}

	err = h.contentService.DeleteArticle(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "文章不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "文章不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "删除文章失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("文章删除成功", nil))
}

// ListArticles 获取文章列表
// @Summary 获取文章列表
// @Description 获取文章列表，支持分类、状态等筛选条件
// @Tags 内容管理
// @Accept json
// @Produce json
// @Param category query string false "分类"
// @Param status query string false "状态"
// @Param keyword query string false "关键词"
// @Param is_top query bool false "是否置顶"
// @Param is_featured query bool false "是否推荐"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Success 200 {object} StandardResponse{data=service.ListArticlesResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/content/articles [get]
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	category := c.Query("category")
	status := c.Query("status")
	keyword := c.Query("keyword")
	isTopStr := c.Query("is_top")
	isFeaturedStr := c.Query("is_featured")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	var isTop *bool
	if isTopStr != "" {
		if val, err := strconv.ParseBool(isTopStr); err == nil {
			isTop = &val
		}
	}

	var isFeatured *bool
	if isFeaturedStr != "" {
		if val, err := strconv.ParseBool(isFeaturedStr); err == nil {
			isFeatured = &val
		}
	}

	req := &service.ListArticlesRequest{
		Category:   category,
		Status:     status,
		Keyword:    keyword,
		IsTop:      isTop,
		IsFeatured: isFeatured,
		Page:       page,
		Limit:      limit,
	}

	response, err := h.contentService.ListArticles(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取文章列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// PublishArticle 发布文章
// @Summary 发布文章
// @Description 管理员发布文章
// @Tags 内容管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/content/articles/{id}/publish [post]
func (h *ArticleHandler) PublishArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的文章ID", err.Error()))
		return
	}

	err = h.contentService.PublishArticle(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "文章不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "文章不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "发布文章失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("文章已发布", nil))
}

// GetFeaturedArticles 获取推荐文章
// @Summary 获取推荐文章
// @Description 获取首页推荐文章
// @Tags 内容管理
// @Accept json
// @Produce json
// @Param limit query int false "数量限制" default(10)
// @Success 200 {object} StandardResponse{data=service.FeaturedArticlesResponse}
// @Failure 500 {object} ErrorResponse
// @Router /api/content/articles/featured [get]
func (h *ArticleHandler) GetFeaturedArticles(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 50 {
		limit = 10
	}

	response, err := h.contentService.GetFeaturedArticles(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取推荐文章失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// CreateCategory 创建分类
// @Summary 创建分类
// @Description 管理员创建文章分类
// @Tags 内容管理
// @Accept json
// @Produce json
// @Param request body service.CreateCategoryRequest true "分类信息"
// @Success 200 {object} StandardResponse{data=model.Category}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/content/categories [post]
func (h *ArticleHandler) CreateCategory(c *gin.Context) {
	var req service.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	category, err := h.contentService.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "创建分类失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("分类创建成功", category))
}

// GetCategories 获取分类列表
// @Summary 获取分类列表
// @Description 获取所有文章分类
// @Tags 内容管理
// @Accept json
// @Produce json
// @Success 200 {object} StandardResponse{data=[]model.Category}
// @Failure 500 {object} ErrorResponse
// @Router /api/content/categories [get]
func (h *ArticleHandler) GetCategories(c *gin.Context) {
	categories, err := h.contentService.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取分类列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", categories))
}

// UpdateCategory 更新分类
// @Summary 更新分类
// @Description 管理员更新分类信息
// @Tags 内容管理
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Param request body service.UpdateCategoryRequest true "更新信息"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/content/categories/{id} [put]
func (h *ArticleHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的分类ID", err.Error()))
		return
	}

	var req service.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	err = h.contentService.UpdateCategory(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "分类不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "分类不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "更新分类失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("分类更新成功", nil))
}

// DeleteCategory 删除分类
// @Summary 删除分类
// @Description 管理员删除分类
// @Tags 内容管理
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/content/categories/{id} [delete]
func (h *ArticleHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的分类ID", err.Error()))
		return
	}

	err = h.contentService.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "分类不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "分类不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "删除分类失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("分类删除成功", nil))
}
