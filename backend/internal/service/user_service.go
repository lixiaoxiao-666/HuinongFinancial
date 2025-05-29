package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"huinong-backend/internal/model"
	"huinong-backend/internal/repository"
)

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
	// cache    cache.CacheInterface
	// sms      sms.SMSInterface
	jwtSecret string
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository, jwtSecret string) UserService {
	return &userService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register 用户注册
func (s *userService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// 检查手机号是否已存在
	existingUser, err := s.userRepo.GetByPhone(ctx, req.Phone)
	if err == nil && existingUser != nil {
		return nil, errors.New("手机号已被注册")
	}

	// 验证短信验证码
	// TODO: 实现短信验证码验证逻辑
	// if !s.sms.VerifyCode(req.Phone, req.SmsCode) {
	//     return nil, errors.New("验证码错误")
	// }

	// 生成密码哈希
	salt := generateSalt()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password+salt), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 创建用户对象
	user := &model.User{
		UUID:         uuid.New().String(),
		Username:     req.Username,
		Phone:        req.Phone,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		Salt:         salt,
		UserType:     req.UserType,
		RealName:     req.RealName,
		Province:     req.Province,
		City:         req.City,
		County:       req.County,
		Address:      req.Address,
		Status:       "active",
	}

	// 保存用户
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	// 生成JWT令牌
	accessToken, refreshToken, err := s.generateTokens(user.ID, "app")
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %v", err)
	}

	// TODO: 创建用户会话
	// session := &model.UserSession{
	// 	UserID:         user.ID,
	// 	SessionID:      generateSessionID(),
	// 	Platform:       "app",
	// 	DeviceID:       getDeviceIDFromContext(ctx),
	// 	DeviceType:     getDeviceTypeFromContext(ctx),
	// 	AppVersion:     getAppVersionFromContext(ctx),
	// 	IPAddress:      getIPFromContext(ctx),
	// 	AccessToken:    accessToken,
	// 	RefreshToken:   refreshToken,
	// 	TokenExpiresAt: timePtr(time.Now().Add(24 * time.Hour)),
	// 	Status:         "active",
	// }

	// TODO: 保存会话到数据库
	// err = s.sessionRepo.Create(ctx, session)

	return &RegisterResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    86400, // 24小时
	}, nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 根据手机号查找用户
	user, err := s.userRepo.GetByPhone(ctx, req.Phone)
	if err != nil || user == nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, errors.New("用户账户已被冻结")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password+user.Salt))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	// 生成JWT令牌
	accessToken, refreshToken, err := s.generateTokens(user.ID, req.Platform)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %v", err)
	}

	// 创建会话
	sessionID := generateSessionID()
	// TODO: 创建会话逻辑
	// session := &model.UserSession{
	// 	UserID:         user.ID,
	// 	SessionID:      sessionID,
	// 	Platform:       req.Platform,
	// 	DeviceID:       req.DeviceID,
	// 	DeviceType:     req.DeviceType,
	// 	AppVersion:     req.AppVersion,
	// 	IPAddress:      getIPFromContext(ctx),
	// 	AccessToken:    accessToken,
	// 	RefreshToken:   refreshToken,
	// 	TokenExpiresAt: timePtr(time.Now().Add(24 * time.Hour)),
	// 	Status:         "active",
	// }

	// TODO: 保存会话到数据库
	// err = s.sessionRepo.Create(ctx, session)

	// 更新用户登录信息
	err = s.userRepo.UpdateLoginInfo(ctx, user.ID, getIPFromContext(ctx))
	if err != nil {
		log.Printf("更新登录信息失败: %v", err)
	}

	return &LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    86400, // 24小时
		SessionID:    sessionID,
	}, nil
}

// Logout 用户登出
func (s *userService) Logout(ctx context.Context, sessionID string) error {
	// TODO: 实现会话注销逻辑
	// 1. 标记会话为已注销
	// 2. 清除相关缓存
	// 3. 记录日志

	// err := s.sessionRepo.UpdateStatus(ctx, sessionID, "revoked")
	// if err != nil {
	//     return fmt.Errorf("注销会话失败: %v", err)
	// }

	// s.cache.Del(fmt.Sprintf("session:%s", sessionID))

	return nil
}

// RefreshToken 刷新令牌
func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// 验证refresh token
	claims, err := s.validateToken(refreshToken)
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	// 检查令牌类型
	if tokenType, ok := claims["type"].(string); !ok || tokenType != "refresh" {
		return nil, errors.New("无效的令牌类型")
	}

	// 获取用户ID
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("无效的用户ID")
	}
	userID := uint64(userIDFloat)

	// 获取平台信息
	platform, _ := claims["platform"].(string)

	// 生成新的令牌对
	newAccessToken, newRefreshToken, err := s.generateTokens(userID, platform)
	if err != nil {
		return nil, fmt.Errorf("生成新令牌失败: %v", err)
	}

	return &TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    86400, // 24小时
	}, nil
}

// GetProfile 获取用户资料
func (s *userService) GetProfile(ctx context.Context, userID uint64) (*UserProfileResponse, error) {
	// 获取用户基本信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 获取认证信息
	authInfo := make(map[string]*model.UserAuth)
	// TODO: 实现获取认证信息逻辑
	// auths, err := s.authRepo.GetByUserID(ctx, userID)
	// if err == nil {
	//     for _, auth := range auths {
	//         authInfo[auth.AuthType] = auth
	//     }
	// }

	// 获取用户标签
	tags := []*model.UserTag{}
	// TODO: 实现获取用户标签逻辑
	// tags, err = s.tagRepo.GetByUserID(ctx, userID)

	return &UserProfileResponse{
		User:     user,
		AuthInfo: authInfo,
		Tags:     tags,
	}, nil
}

// UpdateProfile 更新用户资料
func (s *userService) UpdateProfile(ctx context.Context, userID uint64, req *UpdateProfileRequest) error {
	// 获取当前用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 更新字段
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Province != "" {
		user.Province = req.Province
	}
	if req.City != "" {
		user.City = req.City
	}
	if req.County != "" {
		user.County = req.County
	}
	if req.Address != "" {
		user.Address = req.Address
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.Birthday != "" {
		// 解析生日字符串
		if birthday, err := time.Parse("2006-01-02", req.Birthday); err == nil {
			user.Birthday = &birthday
		}
	}

	// 保存更新
	return s.userRepo.Update(ctx, user)
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, userID uint64, req *ChangePasswordRequest) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword+user.Salt))
	if err != nil {
		return errors.New("原密码错误")
	}

	// 生成新密码哈希
	newSalt := generateSalt()
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword+newSalt), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	// 更新密码
	user.PasswordHash = string(newPasswordHash)
	user.Salt = newSalt

	return s.userRepo.Update(ctx, user)
}

// SubmitRealNameAuth 提交实名认证
func (s *userService) SubmitRealNameAuth(ctx context.Context, userID uint64, req *RealNameAuthRequest) error {
	// TODO: 实现实名认证逻辑
	// 1. 验证身份证号格式
	// 2. 调用第三方实名认证API
	// 3. 保存认证信息
	// 4. 更新用户认证状态

	return nil
}

// FreezeUser 冻结用户
func (s *userService) FreezeUser(ctx context.Context, userID uint64, reason string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %v", err)
	}

	user.Status = "frozen"
	// TODO: 记录冻结原因到操作日志

	return s.userRepo.Update(ctx, user)
}

// UnfreezeUser 解冻用户
func (s *userService) UnfreezeUser(ctx context.Context, userID uint64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %v", err)
	}

	user.Status = "active"
	// TODO: 记录解冻操作到操作日志

	return s.userRepo.Update(ctx, user)
}

// SubmitBankCardAuth 提交银行卡认证
func (s *userService) SubmitBankCardAuth(ctx context.Context, userID uint64, req *BankCardAuthRequest) error {
	// TODO: 实现银行卡认证逻辑
	return nil
}

// ReviewAuth 认证审核
func (s *userService) ReviewAuth(ctx context.Context, authID uint64, req *ReviewAuthRequest) error {
	// TODO: 实现认证审核逻辑
	return nil
}

// AddUserTag 添加用户标签
func (s *userService) AddUserTag(ctx context.Context, userID uint64, req *AddTagRequest) error {
	// TODO: 实现添加用户标签逻辑
	return nil
}

// GetUserTags 获取用户标签
func (s *userService) GetUserTags(ctx context.Context, userID uint64, tagType string) ([]*model.UserTag, error) {
	// TODO: 实现获取用户标签逻辑
	return []*model.UserTag{}, nil
}

// RemoveUserTag 删除用户标签
func (s *userService) RemoveUserTag(ctx context.Context, userID uint64, tagKey string) error {
	// TODO: 实现删除用户标签逻辑
	return nil
}

// ListUsers 用户列表
func (s *userService) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	// TODO: 需要适配repository接口
	return &ListUsersResponse{}, nil
}

// GetUserStatistics 获取用户统计
func (s *userService) GetUserStatistics(ctx context.Context) (*UserStatistics, error) {
	// TODO: 实现用户统计逻辑
	return &UserStatistics{}, nil
}

// UploadAvatar 上传头像
func (s *userService) UploadAvatar(ctx context.Context, userID uint64, file io.Reader, filename string) (string, error) {
	// TODO: 实现头像上传逻辑
	return "", nil
}

// ==================== 工具函数 ====================

// generateSalt 生成随机盐值
func generateSalt() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 16)
	rand.Read(b)
	for i := range b {
		b[i] = letters[b[i]%byte(len(letters))]
	}
	return string(b)
}

// generateSessionID 生成会话ID
func generateSessionID() string {
	return uuid.New().String()
}

// generateTokens 生成JWT令牌
func (s *userService) generateTokens(userID uint64, platform string) (string, string, error) {
	now := time.Now()

	// Access Token (24小时有效)
	accessClaims := jwt.MapClaims{
		"user_id":  userID,
		"platform": platform,
		"type":     "access",
		"iat":      now.Unix(),
		"exp":      now.Add(24 * time.Hour).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	// Refresh Token (7天有效)
	refreshClaims := jwt.MapClaims{
		"user_id":  userID,
		"platform": platform,
		"type":     "refresh",
		"iat":      now.Unix(),
		"exp":      now.Add(7 * 24 * time.Hour).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// validateToken 验证JWT令牌
func (s *userService) validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token无效")
}

// ==================== 上下文辅助函数 ====================

// getIPFromContext 从上下文中获取IP地址
func getIPFromContext(ctx context.Context) string {
	if ip, ok := ctx.Value("client_ip").(string); ok {
		return ip
	}
	return "127.0.0.1"
}

// getDeviceIDFromContext 从上下文中获取设备ID
func getDeviceIDFromContext(ctx context.Context) string {
	if deviceID, ok := ctx.Value("device_id").(string); ok {
		return deviceID
	}
	return ""
}

// getDeviceTypeFromContext 从上下文中获取设备类型
func getDeviceTypeFromContext(ctx context.Context) string {
	if deviceType, ok := ctx.Value("device_type").(string); ok {
		return deviceType
	}
	return "unknown"
}

// getAppVersionFromContext 从上下文中获取应用版本
func getAppVersionFromContext(ctx context.Context) string {
	if version, ok := ctx.Value("app_version").(string); ok {
		return version
	}
	return ""
}

// timePtr 时间指针辅助函数
func timePtr(t time.Time) *time.Time {
	return &t
}
