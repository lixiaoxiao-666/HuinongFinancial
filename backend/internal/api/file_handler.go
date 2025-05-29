package api

import (
	"backend/internal/service"
	"backend/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// FileHandler 文件处理器
type FileHandler struct {
	fileService *service.FileService
	log         *zap.Logger
}

// NewFileHandler 创建文件处理器
func NewFileHandler(fileService *service.FileService, log *zap.Logger) *FileHandler {
	return &FileHandler{
		fileService: fileService,
		log:         log,
	}
}

// RegisterFileRoutes 注册文件路由
func RegisterFileRoutes(group *gin.RouterGroup, handler *FileHandler, authMiddleware gin.HandlerFunc) {
	// 需要认证的路由
	authenticated := group.Group("")
	authenticated.Use(authMiddleware)
	{
		authenticated.POST("/upload", handler.UploadFile)
	}
}

// UploadFile 文件上传
// @Summary 文件上传
// @Description 上传文件并返回文件ID，支持身份证、银行流水、工作证明等文档类型
// @Tags 文件服务
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param file_type formData string true "文件类型。可选值：id_card,bank_flow,work_certificate,income_proof,other"
// @Param business_type formData string false "业务类型。可选值：loan,machinery_leasing"
// @Success 200 {object} CommonResponse{data=FileUploadResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/files/upload [post]
// @Security BearerAuth
func (h *FileHandler) UploadFile(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		h.log.Warn("文件上传时未找到用户ID")
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	// 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		h.log.Warn("获取上传文件失败", zap.Error(err))
		pkg.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 获取文件类型
	fileType := c.PostForm("file_type")
	if fileType == "" {
		pkg.BadRequest(c, "请指定文件类型")
		return
	}

	// 获取业务类型（可选）
	businessType := c.PostForm("business_type")

	// 验证文件类型
	validFileTypes := []string{"id_card", "bank_flow", "work_certificate", "income_proof", "other"}
	isValidType := false
	for _, validType := range validFileTypes {
		if fileType == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		pkg.BadRequest(c, "不支持的文件类型")
		return
	}

	// 验证文件大小（例如：10MB限制）
	if file.Size > 10*1024*1024 {
		pkg.BadRequest(c, "文件大小不能超过10MB")
		return
	}

	// 验证文件格式（例如：只允许图片和PDF）
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".pdf"}
	filename := file.Filename
	isValidExtension := false
	for _, ext := range allowedExtensions {
		if len(filename) >= len(ext) && filename[len(filename)-len(ext):] == ext {
			isValidExtension = true
			break
		}
	}
	if !isValidExtension {
		pkg.BadRequest(c, "只支持JPG、PNG、PDF格式的文件")
		return
	}

	// 调用服务层上传文件
	fileResponse, err := h.fileService.UploadFile(c.Request.Context(), userID.(string), file, fileType, businessType)
	if err != nil {
		h.log.Error("文件上传失败", zap.Error(err),
			zap.String("user_id", userID.(string)),
			zap.String("file_type", fileType),
			zap.String("filename", file.Filename))

		switch err.Error() {
		case "文件格式不支持":
			pkg.BadRequest(c, "文件格式不支持")
		case "文件大小超出限制":
			pkg.BadRequest(c, "文件大小超出限制")
		case "存储服务不可用":
			pkg.InternalError(c, "存储服务暂时不可用，请稍后重试")
		default:
			pkg.InternalError(c, "文件上传失败，请稍后重试")
		}
		return
	}

	// 记录成功日志
	h.log.Info("文件上传成功",
		zap.String("file_id", fileResponse.FileID),
		zap.String("user_id", userID.(string)),
		zap.String("file_type", fileType),
		zap.String("filename", file.Filename),
		zap.Int64("file_size", file.Size))

	// 返回成功响应
	pkg.Success(c, FileUploadResponse{
		FileID:     fileResponse.FileID,
		FileName:   fileResponse.FileName,
		FileSize:   fileResponse.FileSize,
		UploadedAt: fileResponse.UploadedAt,
	})
}
