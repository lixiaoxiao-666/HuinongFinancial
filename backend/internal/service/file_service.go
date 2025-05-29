package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"backend/internal/data"
	"backend/pkg"

	"go.uber.org/zap"
)

// FileService 文件服务
type FileService struct {
	data *data.Data
	log  *zap.Logger
}

// NewFileService 创建文件服务
func NewFileService(data *data.Data, log *zap.Logger) *FileService {
	return &FileService{
		data: data,
		log:  log,
	}
}

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	FileID     string    `json:"file_id"`
	FileName   string    `json:"file_name"`
	FileSize   int64     `json:"file_size"`
	FileType   string    `json:"file_type"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// UploadFile 上传文件
func (s *FileService) UploadFile(ctx context.Context, userID string, file *multipart.FileHeader, fileType, businessType string) (*FileUploadResponse, error) {
	// 生成文件ID
	fileID := pkg.GenerateFileID()

	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// 生成存储文件名
	storedFileName := fmt.Sprintf("%s_%s%s", userID, fileID, ext)

	// 这里应该实现真实的文件存储逻辑
	// 例如：保存到本地磁盘、OSS、AWS S3等
	// 目前返回模拟数据

	// 创建文件记录
	fileRecord := data.FileUpload{
		FileID:       fileID,
		UserID:       userID,
		FileName:     file.Filename,
		FileSize:     file.Size,
		FileType:     fileType,
		BusinessType: businessType,
		StoragePath:  fmt.Sprintf("/uploads/%s", storedFileName),
		UploadedAt:   time.Now(),
		Status:       1, // 1表示正常
	}

	// 保存到数据库
	if err := s.data.DB.Create(&fileRecord).Error; err != nil {
		s.log.Error("保存文件记录失败", zap.Error(err))
		return nil, fmt.Errorf("保存文件记录失败")
	}

	return &FileUploadResponse{
		FileID:     fileID,
		FileName:   file.Filename,
		FileSize:   file.Size,
		FileType:   fileType,
		UploadedAt: fileRecord.UploadedAt,
	}, nil
}
