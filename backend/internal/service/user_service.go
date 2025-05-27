package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend/internal/data"
	"backend/pkg"

	"go.uber.org/zap"
)

// UserService 用户服务
type UserService struct {
	data       *data.Data
	jwtManager *pkg.JWTManager
	log        *zap.Logger
}

// NewUserService 创建用户服务
func NewUserService(data *data.Data, jwtManager *pkg.JWTManager, log *zap.Logger) *UserService {
	return &UserService{
		data:       data,
		jwtManager: jwtManager,
		log:        log,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required,len=11" validate:"required"`
	Password string `json:"password" binding:"required,min=6" validate:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	UserID    string `json:"user_id"`
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	UserID    string    `json:"user_id"`
	Phone     string    `json:"phone"`
	Nickname  string    `json:"nickname"`
	AvatarURL string    `json:"avatar_url"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	Nickname  string `json:"nickname" validate:"max=100"`
	AvatarURL string `json:"avatar_url" validate:"url"`
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*LoginResponse, error) {
	// 检查手机号是否已注册
	var existUser data.User
	if err := s.data.DB.Where("phone = ?", req.Phone).First(&existUser).Error; err == nil {
		return nil, errors.New("手机号已注册")
	}

	// 加密密码
	hashedPassword, err := pkg.HashPassword(req.Password)
	if err != nil {
		s.log.Error("密码加密失败", zap.Error(err))
		return nil, errors.New("注册失败")
	}

	// 创建用户
	user := data.User{
		UserID:       pkg.GenerateUserID(),
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		Nickname:     fmt.Sprintf("用户%s", req.Phone[7:11]),
		Status:       0,
		RegisteredAt: time.Now(),
	}

	if err := s.data.DB.Create(&user).Error; err != nil {
		s.log.Error("创建用户失败", zap.Error(err))
		return nil, errors.New("注册失败")
	}

	// 生成token
	token, err := s.jwtManager.GenerateToken(user.UserID, "user", "")
	if err != nil {
		s.log.Error("生成token失败", zap.Error(err))
		return nil, errors.New("登录失败")
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	s.data.DB.Save(&user)

	return &LoginResponse{
		UserID:    user.UserID,
		Token:     token,
		ExpiresIn: int64(24 * time.Hour / time.Second), // 24小时
	}, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	var user data.User
	if err := s.data.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	if !pkg.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("密码错误")
	}

	// 检查用户状态
	if user.Status != 0 {
		return nil, errors.New("用户已被禁用")
	}

	// 生成token
	token, err := s.jwtManager.GenerateToken(user.UserID, "user", "")
	if err != nil {
		s.log.Error("生成token失败", zap.Error(err))
		return nil, errors.New("登录失败")
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	s.data.DB.Save(&user)

	return &LoginResponse{
		UserID:    user.UserID,
		Token:     token,
		ExpiresIn: int64(24 * time.Hour / time.Second), // 24小时
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, userID string) (*UserInfoResponse, error) {
	var user data.User
	if err := s.data.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 脱敏处理手机号
	phone := user.Phone
	if len(phone) == 11 {
		phone = phone[:3] + "****" + phone[7:]
	}

	return &UserInfoResponse{
		UserID:    user.UserID,
		Phone:     phone,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}, nil
}

// UpdateUserInfo 更新用户信息
func (s *UserService) UpdateUserInfo(ctx context.Context, userID string, req *UpdateUserRequest) error {
	var user data.User
	if err := s.data.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}

	if len(updates) > 0 {
		if err := s.data.DB.Model(&user).Updates(updates).Error; err != nil {
			s.log.Error("更新用户信息失败", zap.Error(err))
			return errors.New("更新失败")
		}
	}

	return nil
}

// verifyCode 验证短信验证码（简化实现）
func (s *UserService) verifyCode(ctx context.Context, phone, code string) bool {
	// 从Redis获取验证码
	key := fmt.Sprintf("sms_code:%s", phone)
	storedCode, err := s.data.Redis.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	return storedCode == code
}

// deleteVerificationCode 删除验证码
func (s *UserService) deleteVerificationCode(ctx context.Context, phone string) {
	key := fmt.Sprintf("sms_code:%s", phone)
	s.data.Redis.Del(ctx, key)
}

// SendVerificationCode 发送验证码（模拟实现）
func (s *UserService) SendVerificationCode(ctx context.Context, phone string) error {
	// 生成6位验证码
	code := fmt.Sprintf("%06d", time.Now().Unix()%1000000)

	// 存储到Redis，5分钟过期
	key := fmt.Sprintf("sms_code:%s", phone)
	if err := s.data.Redis.Set(ctx, key, code, 5*time.Minute).Err(); err != nil {
		s.log.Error("存储验证码失败", zap.Error(err))
		return errors.New("发送验证码失败")
	}

	// 这里应该调用短信服务发送验证码，暂时只记录日志
	s.log.Info("发送验证码", zap.String("phone", phone), zap.String("code", code))

	return nil
}
