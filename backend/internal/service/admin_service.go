package service

import (
	"backend/internal/data"
	"backend/pkg"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminService OA管理服务
type AdminService struct {
	data       *data.Data
	jwtManager *pkg.JWTManager
	log        *zap.Logger
}

// ReviewDetails 审批详情
type ReviewDetails struct {
	ApprovedAmount      *float64
	ApprovedTermMonths  *int
	Comments            string
	RequiredInfoDetails *string
}

// SystemStats 系统统计信息
type SystemStats struct {
	TotalApplications    int64   `json:"total_applications"`
	PendingApplications  int64   `json:"pending_applications"`
	ApprovedApplications int64   `json:"approved_applications"`
	RejectedApplications int64   `json:"rejected_applications"`
	TodayApplications    int64   `json:"today_applications"`
	AIProcessedRate      float64 `json:"ai_processed_rate"`
	AvgProcessingTime    float64 `json:"avg_processing_time_hours"`
	AIApprovalEnabled    bool    `json:"ai_approval_enabled"`
}

// NewAdminService 创建OA管理服务
func NewAdminService(data *data.Data, jwtManager *pkg.JWTManager, log *zap.Logger) *AdminService {
	return &AdminService{
		data:       data,
		jwtManager: jwtManager,
		log:        log,
	}
}

// Login OA用户登录验证
func (s *AdminService) Login(username, password string) (*data.OAUser, string, error) {
	// 查找OA用户
	var oaUser data.OAUser
	err := s.data.DB.Where("username = ? AND status = 0", username).First(&oaUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("用户不存在")
		}
		s.log.Error("查询OA用户失败", zap.String("username", username), zap.Error(err))
		return nil, "", errors.New("查询用户失败")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(oaUser.PasswordHash), []byte(password))
	if err != nil {
		s.log.Warn("OA用户密码验证失败", zap.String("username", username))
		return nil, "", errors.New("密码错误")
	}

	// 生成JWT Token
	token, err := s.jwtManager.GenerateToken(oaUser.OAUserID, "oa_user", oaUser.Role)
	if err != nil {
		s.log.Error("生成JWT Token失败", zap.String("oa_user_id", oaUser.OAUserID), zap.Error(err))
		return nil, "", errors.New("生成Token失败")
	}

	return &oaUser, token, nil
}

// GetPendingApplications 获取待审批申请列表
func (s *AdminService) GetPendingApplications(filters map[string]interface{}, page, limit int) ([]data.LoanApplication, int64, error) {
	query := s.data.DB.Model(&data.LoanApplication{})

	// 应用过滤条件
	if status, exists := filters["status"]; exists {
		query = query.Where("status = ?", status)
	} else {
		// 默认查询需要人工审核的申请
		query = query.Where("status IN (?)", []string{"AI_REVIEWING", "MANUAL_REVIEW_REQUIRED", "MORE_INFO_REQUESTED"})
	}

	if applicationID, exists := filters["application_id"]; exists {
		query = query.Where("application_id LIKE ?", "%"+applicationID.(string)+"%")
	}

	// 申请人姓名过滤需要通过申请人快照数据进行
	if applicantName, exists := filters["applicant_name"]; exists {
		query = query.Where("JSON_EXTRACT(applicant_snapshot, '$.real_name') LIKE ?", "%"+applicantName.(string)+"%")
	}

	// 获取总数
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		s.log.Error("统计待审批申请数量失败", zap.Error(err))
		return nil, 0, err
	}

	// 分页查询
	var applications []data.LoanApplication
	offset := (page - 1) * limit
	err = query.Order("submitted_at DESC").Offset(offset).Limit(limit).Find(&applications).Error
	if err != nil {
		s.log.Error("查询待审批申请列表失败", zap.Error(err))
		return nil, 0, err
	}

	return applications, total, nil
}

// GetApplicationDetail 获取申请详情
func (s *AdminService) GetApplicationDetail(applicationID string) (*data.LoanApplication, error) {
	var application data.LoanApplication
	err := s.data.DB.Where("application_id = ?", applicationID).First(&application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("申请不存在")
		}
		s.log.Error("查询申请详情失败", zap.String("application_id", applicationID), zap.Error(err))
		return nil, err
	}

	return &application, nil
}

// GetApplicationHistory 获取申请历史记录
func (s *AdminService) GetApplicationHistory(applicationID string) ([]data.LoanApplicationHistory, error) {
	var history []data.LoanApplicationHistory
	err := s.data.DB.Where("application_id = ?", applicationID).
		Order("occurred_at ASC").
		Find(&history).Error
	if err != nil {
		s.log.Error("查询申请历史失败", zap.String("application_id", applicationID), zap.Error(err))
		return nil, err
	}

	return history, nil
}

// GetApplicationFiles 获取申请相关文件
func (s *AdminService) GetApplicationFiles(applicationID string) ([]data.UploadedFile, error) {
	var files []data.UploadedFile
	err := s.data.DB.Where("related_id = ? AND purpose = ?", applicationID, "loan_document").
		Order("uploaded_at ASC").
		Find(&files).Error
	if err != nil {
		s.log.Error("查询申请文件失败", zap.String("application_id", applicationID), zap.Error(err))
		return nil, err
	}

	return files, nil
}

// ReviewApplication 审批申请
func (s *AdminService) ReviewApplication(applicationID, adminUserID, decision string, details *ReviewDetails) error {
	// 开始事务
	tx := s.data.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查找申请
	var application data.LoanApplication
	err := tx.Where("application_id = ?", applicationID).First(&application).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("申请不存在")
		}
		return err
	}

	// 检查申请状态是否允许审批
	if !s.isStatusReviewable(application.Status) {
		tx.Rollback()
		return errors.New("当前申请状态不允许审批")
	}

	// 记录原状态
	originalStatus := application.Status

	// 更新申请状态和审批信息
	now := time.Now()
	updates := map[string]interface{}{
		"processed_by":    adminUserID,
		"processed_at":    &now,
		"decision_reason": details.Comments,
	}

	switch decision {
	case "approved":
		updates["status"] = "APPROVED"
		updates["final_decision"] = "approved"
		if details.ApprovedAmount != nil {
			updates["approved_amount"] = *details.ApprovedAmount
		} else {
			updates["approved_amount"] = application.AmountApplied
		}
		if details.ApprovedTermMonths != nil {
			updates["approved_term_months"] = *details.ApprovedTermMonths
		} else {
			updates["approved_term_months"] = application.TermMonthsApplied
		}

	case "rejected":
		updates["status"] = "REJECTED"
		updates["final_decision"] = "rejected"

	case "request_more_info":
		updates["status"] = "MORE_INFO_REQUESTED"
		if details.RequiredInfoDetails != nil {
			updates["decision_reason"] = *details.RequiredInfoDetails
		}
	}

	err = tx.Model(&application).Updates(updates).Error
	if err != nil {
		tx.Rollback()
		s.log.Error("更新申请审批信息失败", zap.String("application_id", applicationID), zap.Error(err))
		return err
	}

	// 记录审批历史
	history := data.LoanApplicationHistory{
		ApplicationID: applicationID,
		StatusFrom:    originalStatus,
		StatusTo:      updates["status"].(string),
		OperatorType:  "MANUAL",
		OperatorID:    adminUserID,
		Comments:      details.Comments,
		OccurredAt:    now,
	}

	err = tx.Create(&history).Error
	if err != nil {
		tx.Rollback()
		s.log.Error("记录审批历史失败", zap.String("application_id", applicationID), zap.Error(err))
		return err
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		s.log.Error("提交审批事务失败", zap.String("application_id", applicationID), zap.Error(err))
		return err
	}

	// 后续处理（如发送通知等）
	go s.postReviewProcess(applicationID, decision, updates)

	return nil
}

// ToggleAIApproval 切换AI审批开关
func (s *AdminService) ToggleAIApproval(enabled bool, adminUserID string) error {
	// 更新系统配置
	config := data.SystemConfiguration{
		ConfigKey:   "ai_approval_enabled",
		ConfigValue: fmt.Sprintf("%t", enabled),
		Description: "AI审批功能开关",
	}

	err := s.data.DB.Save(&config).Error
	if err != nil {
		s.log.Error("更新AI审批开关配置失败", zap.Bool("enabled", enabled), zap.Error(err))
		return err
	}

	s.log.Info("AI审批开关已更新",
		zap.Bool("enabled", enabled),
		zap.String("admin_user_id", adminUserID))

	return nil
}

// GetSystemStats 获取系统统计信息
func (s *AdminService) GetSystemStats() (*SystemStats, error) {
	stats := &SystemStats{}

	// 总申请数
	s.data.DB.Model(&data.LoanApplication{}).Count(&stats.TotalApplications)

	// 各状态申请数
	s.data.DB.Model(&data.LoanApplication{}).
		Where("status IN (?)", []string{"AI_REVIEWING", "MANUAL_REVIEW_REQUIRED", "MORE_INFO_REQUESTED"}).
		Count(&stats.PendingApplications)

	s.data.DB.Model(&data.LoanApplication{}).
		Where("status = ?", "APPROVED").
		Count(&stats.ApprovedApplications)

	s.data.DB.Model(&data.LoanApplication{}).
		Where("status = ?", "REJECTED").
		Count(&stats.RejectedApplications)

	// 今日申请数
	today := time.Now().Truncate(24 * time.Hour)
	s.data.DB.Model(&data.LoanApplication{}).
		Where("submitted_at >= ?", today).
		Count(&stats.TodayApplications)

	// AI处理率
	var aiProcessedCount int64
	s.data.DB.Model(&data.LoanApplication{}).
		Where("ai_risk_score IS NOT NULL").
		Count(&aiProcessedCount)

	if stats.TotalApplications > 0 {
		stats.AIProcessedRate = float64(aiProcessedCount) / float64(stats.TotalApplications) * 100
	}

	// 平均处理时间（已完成的申请）
	var avgProcessingTime sql.NullFloat64
	s.data.DB.Model(&data.LoanApplication{}).
		Select("AVG(TIMESTAMPDIFF(HOUR, submitted_at, processed_at))").
		Where("processed_at IS NOT NULL").
		Scan(&avgProcessingTime)

	if avgProcessingTime.Valid {
		stats.AvgProcessingTime = avgProcessingTime.Float64
	}

	// AI审批开关状态
	var config data.SystemConfiguration
	err := s.data.DB.Where("config_key = ?", "ai_approval_enabled").First(&config).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.log.Error("查询AI审批开关状态失败", zap.Error(err))
	} else if err == nil {
		stats.AIApprovalEnabled = config.ConfigValue == "true"
	}

	return stats, nil
}

// 创建默认OA用户
func (s *AdminService) CreateDefaultOAUsers() error {
	// 检查是否已存在默认用户
	var count int64
	s.data.DB.Model(&data.OAUser{}).Count(&count)
	if count > 0 {
		return nil // 已存在用户，跳过初始化
	}

	// 创建默认管理员账号
	adminPasswordHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	adminUser := data.OAUser{
		OAUserID:     "oa_admin_001",
		Username:     "admin",
		PasswordHash: string(adminPasswordHash),
		Role:         "ADMIN",
		DisplayName:  "系统管理员",
		Email:        "admin@example.com",
		Status:       0,
	}

	err := s.data.DB.Create(&adminUser).Error
	if err != nil {
		s.log.Error("创建默认管理员失败", zap.Error(err))
		return err
	}

	// 创建默认审批员账号
	reviewerPasswordHash, _ := bcrypt.GenerateFromPassword([]byte("reviewer123"), bcrypt.DefaultCost)
	reviewerUser := data.OAUser{
		OAUserID:     "oa_reviewer_001",
		Username:     "reviewer",
		PasswordHash: string(reviewerPasswordHash),
		Role:         "REVIEWER",
		DisplayName:  "审批员",
		Email:        "reviewer@example.com",
		Status:       0,
	}

	err = s.data.DB.Create(&reviewerUser).Error
	if err != nil {
		s.log.Error("创建默认审批员失败", zap.Error(err))
		return err
	}

	s.log.Info("默认OA用户创建成功")
	return nil
}

// 辅助方法

// isStatusReviewable 检查状态是否可审批
func (s *AdminService) isStatusReviewable(status string) bool {
	reviewableStatuses := []string{
		"AI_REVIEWING",
		"MANUAL_REVIEW_REQUIRED",
		"MORE_INFO_REQUESTED",
	}

	for _, reviewableStatus := range reviewableStatuses {
		if status == reviewableStatus {
			return true
		}
	}
	return false
}

// postReviewProcess 审批后续处理
func (s *AdminService) postReviewProcess(applicationID, decision string, updates map[string]interface{}) {
	// 这里可以添加审批后的后续处理逻辑
	// 例如：发送通知、更新其他系统状态等

	s.log.Info("审批后续处理完成",
		zap.String("application_id", applicationID),
		zap.String("decision", decision))

	// TODO: 发送通知给申请人
	// TODO: 集成其他系统
}

// SimulateAIReview 模拟AI审批（用于演示）
func (s *AdminService) SimulateAIReview(applicationID string) error {
	// 查找申请
	var application data.LoanApplication
	err := s.data.DB.Where("application_id = ?", applicationID).First(&application).Error
	if err != nil {
		return err
	}

	// 检查AI审批是否已启用
	var config data.SystemConfiguration
	err = s.data.DB.Where("config_key = ?", "ai_approval_enabled").First(&config).Error
	aiEnabled := (err == nil && config.ConfigValue == "true")

	if !aiEnabled {
		s.log.Info("AI审批已禁用，跳过AI审批", zap.String("application_id", applicationID))
		return nil
	}

	// 模拟AI分析过程
	time.Sleep(2 * time.Second) // 模拟AI处理时间

	// 生成模拟的AI评分和建议
	riskScore := s.generateMockRiskScore(application)
	suggestion := s.generateMockSuggestion(riskScore, application)

	// 更新申请的AI分析结果
	updates := map[string]interface{}{
		"ai_risk_score": riskScore,
		"ai_suggestion": suggestion,
	}

	// 根据风险评分决定是否需要人工复核
	if riskScore >= 80 {
		updates["status"] = "APPROVED"
		updates["final_decision"] = "approved"
		updates["approved_amount"] = application.AmountApplied
		updates["approved_term_months"] = application.TermMonthsApplied
	} else if riskScore <= 30 {
		updates["status"] = "REJECTED"
		updates["final_decision"] = "rejected"
		updates["decision_reason"] = "AI风险评估：申请风险过高"
	} else {
		updates["status"] = "MANUAL_REVIEW_REQUIRED"
	}

	err = s.data.DB.Model(&application).Updates(updates).Error
	if err != nil {
		s.log.Error("更新AI审批结果失败", zap.String("application_id", applicationID), zap.Error(err))
		return err
	}

	// 记录AI审批历史
	history := data.LoanApplicationHistory{
		ApplicationID: applicationID,
		StatusFrom:    "SUBMITTED",
		StatusTo:      updates["status"].(string),
		OperatorType:  "AI",
		OperatorID:    "ai_system",
		Comments:      suggestion,
		OccurredAt:    time.Now(),
	}

	err = s.data.DB.Create(&history).Error
	if err != nil {
		s.log.Error("记录AI审批历史失败", zap.String("application_id", applicationID), zap.Error(err))
	}

	s.log.Info("AI审批完成",
		zap.String("application_id", applicationID),
		zap.Int("risk_score", riskScore),
		zap.String("new_status", updates["status"].(string)))

	return nil
}

// generateMockRiskScore 生成模拟风险评分
func (s *AdminService) generateMockRiskScore(application data.LoanApplication) int {
	// 基于申请金额、用户信息等生成模拟评分
	score := 75 // 基础分

	// 申请金额影响
	if application.AmountApplied <= 30000 {
		score += 10
	} else if application.AmountApplied >= 100000 {
		score -= 15
	}

	// 随机因子
	score += rand.Intn(20) - 10

	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	return score
}

// generateMockSuggestion 生成模拟AI建议
func (s *AdminService) generateMockSuggestion(riskScore int, application data.LoanApplication) string {
	if riskScore >= 85 {
		return "风险评估良好，建议自动通过。申请人信用记录优秀，收入稳定。"
	} else if riskScore >= 70 {
		return "风险评估中等，建议人工复核。重点关注申请人还款能力证明。"
	} else if riskScore >= 50 {
		return "风险评估较高，建议谨慎审批。需要详细核实申请人财务状况和担保情况。"
	} else {
		return "风险评估高，建议拒绝申请。申请人信用记录存在问题或收入证明不足。"
	}
}

// GetPendingTasks 获取指定管理员的待办事项
func (s *AdminService) GetPendingTasks(adminUserID string) ([]map[string]interface{}, error) {
	// 查询该管理员相关的待办事项
	var applications []data.LoanApplication
	err := s.data.DB.Where("status IN (?)", []string{"MANUAL_REVIEW_REQUIRED", "MORE_INFO_REQUESTED"}).
		Limit(5).
		Order("submitted_at ASC").
		Find(&applications).Error
	if err != nil {
		return nil, err
	}

	var tasks []map[string]interface{}
	for _, app := range applications {
		var applicantInfo map[string]interface{}
		if app.ApplicantSnapshot != nil {
			json.Unmarshal(app.ApplicantSnapshot, &applicantInfo)
		}

		applicantName := ""
		if realName, exists := applicantInfo["real_name"]; exists {
			applicantName = realName.(string)
		}

		task := map[string]interface{}{
			"task_id":       app.ApplicationID,
			"task_type":     "loan_approval",
			"title":         fmt.Sprintf("审批 %s 的贷款申请", applicantName),
			"description":   fmt.Sprintf("申请金额：%.2f元", app.AmountApplied),
			"priority":      s.calculateTaskPriority(app),
			"submitted_at":  app.SubmittedAt,
			"waiting_hours": int(time.Since(app.SubmittedAt).Hours()),
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetOAUsers 获取OA用户列表
func (s *AdminService) GetOAUsers(page, limit int, role string) ([]map[string]interface{}, int64, error) {
	query := s.data.DB.Model(&data.OAUser{})

	if role != "" {
		query = query.Where("role = ?", role)
	}

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	var users []data.OAUser
	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	var result []map[string]interface{}
	for _, user := range users {
		item := map[string]interface{}{
			"oa_user_id":   user.OAUserID,
			"username":     user.Username,
			"role":         user.Role,
			"display_name": user.DisplayName,
			"email":        user.Email,
			"status":       user.Status,
			"created_at":   user.CreatedAt,
			"updated_at":   user.UpdatedAt,
		}
		result = append(result, item)
	}

	return result, total, nil
}

// CreateOAUser 创建OA用户
func (s *AdminService) CreateOAUser(req interface{}, creatorID string) error {
	// 通过反射获取结构体字段
	var username, password, role, displayName, email string

	// 尝试获取字段值
	switch r := req.(type) {
	case *struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		Role        string `json:"role"`
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
	}:
		username = r.Username
		password = r.Password
		role = r.Role
		displayName = r.DisplayName
		email = r.Email
	default:
		return errors.New("无效的请求参数")
	}

	// 检查用户名是否已存在
	var count int64
	err := s.data.DB.Model(&data.OAUser{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 生成密码哈希
	hashedPassword, err := pkg.HashPassword(password)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 创建OA用户
	oaUser := data.OAUser{
		OAUserID:     pkg.GenerateOAUserID(),
		Username:     username,
		PasswordHash: hashedPassword,
		Role:         role,
		DisplayName:  displayName,
		Email:        email,
		Status:       0,
	}

	return s.data.DB.Create(&oaUser).Error
}

// UpdateOAUserStatus 更新OA用户状态
func (s *AdminService) UpdateOAUserStatus(userID string, status int8, operatorID string) error {
	// 检查用户是否存在
	var user data.OAUser
	err := s.data.DB.Where("oa_user_id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 更新状态
	return s.data.DB.Model(&user).Update("status", status).Error
}

// GetOperationLogs 获取操作日志
func (s *AdminService) GetOperationLogs(req interface{}) ([]map[string]interface{}, int64, error) {
	// 模拟操作日志数据 - 在实际实现中应该从专门的日志表查询
	logs := []map[string]interface{}{
		{
			"id":            1,
			"operator_id":   "oa_admin001",
			"operator_name": "管理员李四",
			"action":        "审批申请",
			"target":        "贷款申请 la_app_001",
			"result":        "已批准",
			"ip_address":    "192.168.1.100",
			"user_agent":    "Mozilla/5.0...",
			"occurred_at":   time.Now().Add(-2 * time.Hour),
		},
		{
			"id":            2,
			"operator_id":   "oa_reviewer001",
			"operator_name": "审批员张三",
			"action":        "创建用户",
			"target":        "新建审批员账号",
			"result":        "创建成功",
			"ip_address":    "192.168.1.101",
			"user_agent":    "Mozilla/5.0...",
			"occurred_at":   time.Now().Add(-4 * time.Hour),
		},
		{
			"id":            3,
			"operator_id":   "oa_admin001",
			"operator_name": "管理员李四",
			"action":        "系统配置",
			"target":        "AI审批开关",
			"result":        "配置成功",
			"ip_address":    "192.168.1.100",
			"user_agent":    "Mozilla/5.0...",
			"occurred_at":   time.Now().Add(-6 * time.Hour),
		},
	}

	// 返回模拟数据
	return logs, int64(len(logs)), nil
}

// GetSystemConfigurations 获取系统配置
func (s *AdminService) GetSystemConfigurations() ([]map[string]interface{}, error) {
	var configs []data.SystemConfiguration
	err := s.data.DB.Find(&configs).Error
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, config := range configs {
		item := map[string]interface{}{
			"config_key":   config.ConfigKey,
			"config_value": config.ConfigValue,
			"description":  config.Description,
			"updated_at":   config.UpdatedAt,
		}
		result = append(result, item)
	}

	return result, nil
}

// UpdateSystemConfiguration 更新系统配置
func (s *AdminService) UpdateSystemConfiguration(configKey, configValue, operatorID string) error {
	// 检查配置是否存在，如果不存在则创建
	var config data.SystemConfiguration
	err := s.data.DB.Where("config_key = ?", configKey).First(&config).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建新配置
		newConfig := data.SystemConfiguration{
			ConfigKey:   configKey,
			ConfigValue: configValue,
			Description: "通过管理后台创建的配置",
		}
		return s.data.DB.Create(&newConfig).Error
	}

	if err != nil {
		return err
	}

	// 更新现有配置
	return s.data.DB.Model(&config).Update("config_value", configValue).Error
}

// calculateTaskPriority 计算任务优先级
func (s *AdminService) calculateTaskPriority(app data.LoanApplication) string {
	waitingHours := int(time.Since(app.SubmittedAt).Hours())

	if waitingHours > 48 {
		return "urgent"
	} else if waitingHours > 24 {
		return "high"
	} else if app.AmountApplied > 100000 {
		return "medium"
	}
	return "normal"
}
