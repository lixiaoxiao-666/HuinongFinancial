package service

import (
	"context"
	"fmt"
	"time"

	"huinong-backend/internal/model"
	"huinong-backend/internal/repository"
)

// systemService 系统服务实现
type systemService struct {
	configRepo  repository.SystemConfigRepository
	fileRepo    repository.FileRepository
	userRepo    repository.UserRepository
	loanRepo    repository.LoanRepository
	machineRepo repository.MachineRepository
	articleRepo repository.ArticleRepository
}

// NewSystemService 创建系统服务实例
func NewSystemService(
	configRepo repository.SystemConfigRepository,
	fileRepo repository.FileRepository,
	userRepo repository.UserRepository,
	loanRepo repository.LoanRepository,
	machineRepo repository.MachineRepository,
	articleRepo repository.ArticleRepository,
) SystemService {
	return &systemService{
		configRepo:  configRepo,
		fileRepo:    fileRepo,
		userRepo:    userRepo,
		loanRepo:    loanRepo,
		machineRepo: machineRepo,
		articleRepo: articleRepo,
	}
}

// ==================== 配置管理 ====================

// GetConfig 获取配置
func (s *systemService) GetConfig(ctx context.Context, configKey string) (string, error) {
	config, err := s.configRepo.Get(ctx, configKey)
	if err != nil {
		return "", fmt.Errorf("获取配置失败: %v", err)
	}
	return config.ConfigValue, nil
}

// SetConfig 设置配置
func (s *systemService) SetConfig(ctx context.Context, configKey, configValue string) error {
	config := &model.SystemConfig{
		ConfigKey:   configKey,
		ConfigValue: configValue,
		ConfigType:  "string",
		Description: "",
		ConfigGroup: "custom",
		IsEditable:  true,
		IsEncrypted: false,
	}

	return s.configRepo.Set(ctx, config)
}

// GetConfigs 获取配置组
func (s *systemService) GetConfigs(ctx context.Context, configGroup string) (map[string]string, error) {
	configs, err := s.configRepo.GetByGroup(ctx, configGroup)
	if err != nil {
		return nil, fmt.Errorf("获取配置组失败: %v", err)
	}

	result := make(map[string]string)
	for _, config := range configs {
		result[config.ConfigKey] = config.ConfigValue
	}

	return result, nil
}

// ==================== 文件管理 ====================

// UploadFile 上传文件
func (s *systemService) UploadFile(ctx context.Context, req *UploadFileRequest) (*UploadFileResponse, error) {
	// TODO: 实现文件上传逻辑
	// 1. 验证文件类型和大小
	// 2. 生成唯一文件名
	// 3. 上传到存储服务（OSS/本地存储）
	// 4. 计算文件哈希
	// 5. 保存文件记录到数据库

	fileUpload := &model.FileUpload{
		FileName:     req.FileName,
		OriginalName: req.FileName,
		FilePath:     "/uploads/" + req.FileName, // 临时路径
		FileURL:      "/uploads/" + req.FileName, // 临时URL
		BusinessType: req.BusinessType,
		BusinessID:   req.BusinessID,
		IsPublic:     req.IsPublic,
		Status:       "uploaded",
		StorageType:  "local",
		// TODO: 设置其他字段
		// FileSize: fileSize,
		// FileHash: fileHash,
		// MimeType: mimeType,
		// UploaderID: getUserIDFromContext(ctx),
		// UploaderType: "user",
	}

	err := s.fileRepo.Create(ctx, fileUpload)
	if err != nil {
		return nil, fmt.Errorf("保存文件记录失败: %v", err)
	}

	return &UploadFileResponse{
		ID:       fileUpload.ID,
		FileName: fileUpload.OriginalName,
		FileURL:  fileUpload.FileURL,
		FileSize: fileUpload.FileSize,
	}, nil
}

// GetFile 获取文件信息
func (s *systemService) GetFile(ctx context.Context, fileID uint64) (*model.FileUpload, error) {
	file, err := s.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 增加访问次数
	s.fileRepo.IncrementAccessCount(ctx, fileID)

	return file, nil
}

// DeleteFile 删除文件
func (s *systemService) DeleteFile(ctx context.Context, fileID uint64) error {
	// TODO: 删除实际文件
	// 1. 获取文件信息
	// 2. 从存储服务中删除文件
	// 3. 删除数据库记录

	return s.fileRepo.Delete(ctx, fileID)
}

// ==================== 健康检查 ====================

// HealthCheck 健康检查
func (s *systemService) HealthCheck(ctx context.Context) (*HealthCheckResponse, error) {
	response := &HealthCheckResponse{
		Status:    "ok",
		Timestamp: time.Now().Unix(),
		Database:  make(map[string]interface{}),
		Redis:     make(map[string]interface{}),
		Services:  make(map[string]interface{}),
	}

	// 检查数据库连接
	_, err := s.userRepo.GetUserCount(ctx)
	if err != nil {
		response.Status = "error"
		response.Database["status"] = "error"
		response.Database["error"] = err.Error()
	} else {
		response.Database["status"] = "ok"
		response.Database["connection"] = "active"
	}

	// TODO: 检查Redis连接
	response.Redis["status"] = "ok"
	response.Redis["connection"] = "not_implemented"

	// TODO: 检查外部服务
	response.Services["sms"] = "not_implemented"
	response.Services["payment"] = "not_implemented"
	response.Services["storage"] = "not_implemented"

	return response, nil
}

// ==================== 系统统计 ====================

// GetSystemStats 获取系统统计
func (s *systemService) GetSystemStats(ctx context.Context) (*SystemStatsResponse, error) {
	stats := &SystemStatsResponse{}

	// 用户统计
	userCount, err := s.userRepo.GetUserCount(ctx)
	if err != nil {
		userCount = 0
	}
	stats.UserCount = userCount

	// 农机统计
	machineCount, err := s.machineRepo.GetMachineCount(ctx)
	if err != nil {
		machineCount = 0
	}
	stats.MachineCount = machineCount

	// 订单统计
	orderCount, err := s.machineRepo.GetOrderCount(ctx)
	if err != nil {
		orderCount = 0
	}
	stats.OrderCount = orderCount

	// TODO: 添加更多统计数据
	stats.ApplicationCount = 0 // 临时设置

	return stats, nil
}

// GetPublicConfigs 获取公开配置
func (s *systemService) GetPublicConfigs(ctx context.Context) (map[string]string, error) {
	// 获取可以公开的配置项
	publicConfigKeys := []string{
		"app.name",
		"app.version",
		"app.description",
		"customer_service.phone",
		"customer_service.hours",
		"privacy_policy.url",
		"terms_of_service.url",
		"about_us.content",
	}

	result := make(map[string]string)

	for _, key := range publicConfigKeys {
		config, err := s.configRepo.Get(ctx, key)
		if err == nil {
			result[key] = config.ConfigValue
		}
	}

	// 设置默认值
	if result["app.name"] == "" {
		result["app.name"] = "数字惠农"
	}
	if result["app.version"] == "" {
		result["app.version"] = "1.0.0"
	}
	if result["customer_service.phone"] == "" {
		result["customer_service.phone"] = "400-123-4567"
	}
	if result["customer_service.hours"] == "" {
		result["customer_service.hours"] = "9:00-18:00"
	}

	return result, nil
}

// GetSystemVersion 获取系统版本信息
func (s *systemService) GetSystemVersion(ctx context.Context) (*SystemVersionResponse, error) {
	version, err := s.configRepo.Get(ctx, "app.version")
	versionStr := "1.0.0"
	if err == nil {
		versionStr = version.ConfigValue
	}

	buildTime, err := s.configRepo.Get(ctx, "app.build_time")
	buildTimeStr := time.Now().Format("2006-01-02 15:04:05")
	if err == nil {
		buildTimeStr = buildTime.ConfigValue
	}

	gitCommit, err := s.configRepo.Get(ctx, "app.git_commit")
	gitCommitStr := "unknown"
	if err == nil {
		gitCommitStr = gitCommit.ConfigValue
	}

	return &SystemVersionResponse{
		Version:     versionStr,
		BuildTime:   buildTimeStr,
		GitCommit:   gitCommitStr,
		GoVersion:   "1.21",
		Environment: "production",
	}, nil
}
