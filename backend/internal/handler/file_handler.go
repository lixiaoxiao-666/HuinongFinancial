package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"huinong-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// FileHandler 文件处理器
type FileHandler struct {
	systemService service.SystemService
}

// NewFileHandler 创建文件处理器
func NewFileHandler(systemService service.SystemService) *FileHandler {
	return &FileHandler{
		systemService: systemService,
	}
}

// UploadFile 上传文件
// @Summary 上传文件
// @Description 上传文件到服务器，支持多种业务类型
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param business_type formData string true "业务类型" Enums(avatar,document,image,loan_material,auth_material)
// @Param business_id formData string false "业务ID"
// @Param description formData string false "文件描述"
// @Success 200 {object} StandardResponse{data=service.UploadFileResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 413 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/files/upload [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "文件获取失败", err.Error()))
		return
	}
	defer file.Close()

	// 获取业务参数
	businessType := c.PostForm("business_type")
	businessIDStr := c.PostForm("business_id")
	// description := c.PostForm("description")  // TODO: 后续可用于文件描述

	if businessType == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "业务类型不能为空", "business_type is required"))
		return
	}

	// 验证业务类型
	validBusinessTypes := []string{"avatar", "document", "image", "loan_material", "auth_material"}
	isValidType := false
	for _, validType := range validBusinessTypes {
		if businessType == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的业务类型", fmt.Sprintf("business_type must be one of: %s", strings.Join(validBusinessTypes, ", "))))
		return
	}

	// 解析业务ID
	var businessID uint64
	if businessIDStr != "" {
		businessID, err = strconv.ParseUint(businessIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的业务ID", err.Error()))
			return
		}
	}

	// 检查文件大小 (最大10MB)
	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if header.Size > maxFileSize {
		c.JSON(http.StatusRequestEntityTooLarge, NewErrorResponse(http.StatusRequestEntityTooLarge, "文件大小超出限制", "文件大小不能超过10MB"))
		return
	}

	// 检查文件类型
	allowedExtensions := map[string][]string{
		"avatar":        {".jpg", ".jpeg", ".png", ".gif"},
		"document":      {".pdf", ".doc", ".docx", ".xls", ".xlsx", ".txt"},
		"image":         {".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"},
		"loan_material": {".jpg", ".jpeg", ".png", ".pdf", ".doc", ".docx"},
		"auth_material": {".jpg", ".jpeg", ".png", ".pdf"},
	}

	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	if allowedExts, exists := allowedExtensions[businessType]; exists {
		isValidExt := false
		for _, ext := range allowedExts {
			if fileExt == ext {
				isValidExt = true
				break
			}
		}
		if !isValidExt {
			c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "不支持的文件类型", fmt.Sprintf("业务类型 %s 支持的文件格式: %s", businessType, strings.Join(allowedExts, ", "))))
			return
		}
	}

	// 创建上传请求
	req := &service.UploadFileRequest{
		File:         file,
		FileName:     header.Filename,
		BusinessType: businessType,
		BusinessID:   businessID,
		IsPublic:     businessType == "avatar" || businessType == "image", // 头像和图片设为公开
	}

	// 调用服务层上传文件
	response, err := h.systemService.UploadFile(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "文件上传失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("文件上传成功", response))
}

// GetFile 获取文件信息
// @Summary 获取文件信息
// @Description 根据文件ID获取文件详细信息
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} StandardResponse{data=model.FileUpload}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/files/{id} [get]
func (h *FileHandler) GetFile(c *gin.Context) {
	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的文件ID", err.Error()))
		return
	}

	file, err := h.systemService.GetFile(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "文件不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "文件不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取文件信息失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", file))
}

// DeleteFile 删除文件
// @Summary 删除文件
// @Description 删除指定的文件
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/files/{id} [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的文件ID", err.Error()))
		return
	}

	err = h.systemService.DeleteFile(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "文件不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "文件不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "删除文件失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("文件删除成功", nil))
}

// UploadMultipleFiles 批量上传文件
// @Summary 批量上传文件
// @Description 一次性上传多个文件
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "文件列表"
// @Param business_type formData string true "业务类型"
// @Param business_id formData string false "业务ID"
// @Success 200 {object} StandardResponse{data=[]service.UploadFileResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/files/upload/batch [post]
func (h *FileHandler) UploadMultipleFiles(c *gin.Context) {
	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "用户未登录", "用户认证信息缺失"))
		return
	}

	// 获取表单数据
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "表单解析失败", err.Error()))
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "没有选择文件", "至少需要选择一个文件"))
		return
	}

	// 限制批量上传数量
	if len(files) > 10 {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "文件数量超出限制", "一次最多上传10个文件"))
		return
	}

	businessType := c.PostForm("business_type")
	businessIDStr := c.PostForm("business_id")

	if businessType == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "业务类型不能为空", "business_type is required"))
		return
	}

	var businessID uint64
	if businessIDStr != "" {
		businessID, err = strconv.ParseUint(businessIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的业务ID", err.Error()))
			return
		}
	}

	var responses []service.UploadFileResponse
	var errors []string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("文件 %s 打开失败: %v", fileHeader.Filename, err))
			continue
		}

		req := &service.UploadFileRequest{
			File:         file,
			FileName:     fileHeader.Filename,
			BusinessType: businessType,
			BusinessID:   businessID,
			IsPublic:     businessType == "avatar" || businessType == "image",
		}

		response, err := h.systemService.UploadFile(c.Request.Context(), req)
		file.Close()

		if err != nil {
			errors = append(errors, fmt.Sprintf("文件 %s 上传失败: %v", fileHeader.Filename, err))
			continue
		}

		responses = append(responses, *response)
	}

	if len(errors) > 0 && len(responses) == 0 {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "所有文件上传失败", strings.Join(errors, "; ")))
		return
	}

	result := map[string]interface{}{
		"uploaded_files": responses,
		"success_count":  len(responses),
		"total_count":    len(files),
	}

	if len(errors) > 0 {
		result["errors"] = errors
		result["error_count"] = len(errors)
	}

	message := fmt.Sprintf("成功上传 %d/%d 个文件", len(responses), len(files))
	c.JSON(http.StatusOK, NewSuccessResponse(message, result))
}
