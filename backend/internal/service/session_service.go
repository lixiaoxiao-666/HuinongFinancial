package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"huinong-backend/internal/cache"
	"huinong-backend/internal/model"
	"huinong-backend/internal/repository"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// SessionService 会话管理服务接口
type SessionService interface {
	// 会话创建
	CreateSession(ctx context.Context, userID uint64, loginInfo *LoginInfo) (*SessionInfo, error)

	// 会话验证
	ValidateSession(ctx context.Context, sessionID string) (*SessionInfo, error)
	ValidateToken(ctx context.Context, accessToken string) (*SessionInfo, error)

	// 会话更新
	UpdateLastActive(ctx context.Context, sessionID string) error
	RefreshSession(ctx context.Context, refreshToken string) (*TokenPair, error)

	// 会话注销
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeUserSessions(ctx context.Context, userID uint64, excludeSessionID string) error
	RevokeAllSessions(ctx context.Context, userID uint64) error

	// 会话查询
	GetUserSessions(ctx context.Context, userID uint64) ([]*SessionInfo, error)
	GetActiveSessionCount(ctx context.Context, userID uint64) (int, error)

	// 清理任务
	CleanupExpiredSessions(ctx context.Context) error
	CleanupUserSessions(ctx context.Context, userID uint64, keepCount int) error
}

// sessionService 会话管理服务实现
type sessionService struct {
	cache           cache.CacheInterface
	sessionRepo     repository.SessionRepository
	jwtSecret       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	maxSessions     int
	instanceID      string
}

// NewSessionService 创建会话管理服务实例
func NewSessionService(
	cache cache.CacheInterface,
	sessionRepo repository.SessionRepository,
	jwtSecret string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	maxSessions int,
) SessionService {
	return &sessionService{
		cache:           cache,
		sessionRepo:     sessionRepo,
		jwtSecret:       jwtSecret,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		maxSessions:     maxSessions,
		instanceID:      generateBackendInstanceID(),
	}
}

// LoginInfo 登录信息
type LoginInfo struct {
	Platform    string `json:"platform"`
	DeviceID    string `json:"device_id"`
	DeviceType  string `json:"device_type"`
	DeviceName  string `json:"device_name"`
	AppVersion  string `json:"app_version"`
	UserAgent   string `json:"user_agent"`
	IPAddress   string `json:"ip_address"`
	Location    string `json:"location"`
	LoginMethod string `json:"login_method"`
}

// SessionInfo 会话信息
type SessionInfo struct {
	SessionID    string       `json:"session_id"`
	UserID       uint64       `json:"user_id"`
	Platform     string       `json:"platform"`
	DeviceInfo   *DeviceInfo  `json:"device_info"`
	NetworkInfo  *NetworkInfo `json:"network_info"`
	TokenInfo    *TokenInfo   `json:"token_info"`
	Status       string       `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
	LastActiveAt time.Time    `json:"last_active_at"`
	ExpiresAt    time.Time    `json:"expires_at"`
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	DeviceID   string `json:"device_id"`
	DeviceType string `json:"device_type"`
	DeviceName string `json:"device_name"`
	AppVersion string `json:"app_version"`
	UserAgent  string `json:"user_agent"`
}

// NetworkInfo 网络信息
type NetworkInfo struct {
	IPAddress string `json:"ip_address"`
	Location  string `json:"location"`
	ISP       string `json:"isp"`
}

// TokenInfo Token信息
type TokenInfo struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

// TokenPair Token对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// CreateSession 创建会话
func (s *sessionService) CreateSession(ctx context.Context, userID uint64, loginInfo *LoginInfo) (*SessionInfo, error) {
	// 检查用户会话数量限制
	activeCount, err := s.GetActiveSessionCount(ctx, userID)
	if err != nil {
		log.Printf("获取活跃会话数量失败: %v", err)
	} else if activeCount >= s.maxSessions {
		// 清理最旧的会话
		err = s.CleanupUserSessions(ctx, userID, s.maxSessions-1)
		if err != nil {
			log.Printf("清理用户会话失败: %v", err)
		}
	}

	// 生成会话ID
	sessionID := generateSessionUniqueID()
	now := time.Now()

	// 生成JWT令牌
	accessToken, refreshToken, err := s.generateTokens(userID, sessionID, loginInfo.Platform)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %v", err)
	}

	// 创建会话信息
	sessionInfo := &SessionInfo{
		SessionID: sessionID,
		UserID:    userID,
		Platform:  loginInfo.Platform,
		DeviceInfo: &DeviceInfo{
			DeviceID:   loginInfo.DeviceID,
			DeviceType: loginInfo.DeviceType,
			DeviceName: loginInfo.DeviceName,
			AppVersion: loginInfo.AppVersion,
			UserAgent:  loginInfo.UserAgent,
		},
		NetworkInfo: &NetworkInfo{
			IPAddress: loginInfo.IPAddress,
			Location:  loginInfo.Location,
		},
		TokenInfo: &TokenInfo{
			AccessToken:      accessToken,
			RefreshToken:     refreshToken,
			AccessExpiresAt:  now.Add(s.accessTokenTTL),
			RefreshExpiresAt: now.Add(s.refreshTokenTTL),
		},
		Status:       "active",
		CreatedAt:    now,
		LastActiveAt: now,
		ExpiresAt:    now.Add(s.refreshTokenTTL),
	}

	// 存储到Redis
	err = s.storeSessionToRedis(ctx, sessionInfo)
	if err != nil {
		return nil, fmt.Errorf("存储会话到Redis失败: %v", err)
	}

	// 存储到数据库
	// 截断设备名称以确保符合数据库字段长度限制
	deviceName := loginInfo.DeviceName
	if len(deviceName) > 500 {
		deviceName = deviceName[:500]
	}

	dbSession := &model.UserSession{
		UserID:           userID,
		SessionID:        sessionID,
		Platform:         loginInfo.Platform,
		DeviceID:         loginInfo.DeviceID,
		DeviceType:       loginInfo.DeviceType,
		DeviceName:       deviceName,
		AppVersion:       loginInfo.AppVersion,
		UserAgent:        loginInfo.UserAgent,
		IPAddress:        loginInfo.IPAddress,
		Location:         loginInfo.Location,
		AccessTokenHash:  hashTokenString(accessToken),
		RefreshTokenHash: hashTokenString(refreshToken),
		TokenExpiresAt:   &sessionInfo.TokenInfo.AccessExpiresAt,
		RefreshExpiresAt: &sessionInfo.TokenInfo.RefreshExpiresAt,
		Status:           "active",
		LoginTime:        now,
		LastActiveAt:     now,
		LoginMethod:      loginInfo.LoginMethod,
	}

	err = s.sessionRepo.Create(ctx, dbSession)
	if err != nil {
		log.Printf("保存会话到数据库失败: %v", err)
		// 不返回错误，因为Redis已经存储成功
	}

	// 发布会话创建事件
	s.publishSessionEvent(ctx, EventSessionCreated, sessionID, userID, loginInfo)

	return sessionInfo, nil
}

// ValidateSession 验证会话
func (s *sessionService) ValidateSession(ctx context.Context, sessionID string) (*SessionInfo, error) {
	// 从Redis获取会话信息
	sessionInfo, err := s.getSessionFromRedis(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("获取会话信息失败: %v", err)
	}

	// 检查会话状态
	if sessionInfo.Status != "active" {
		return nil, fmt.Errorf("会话已失效")
	}

	// 检查会话是否过期
	if time.Now().After(sessionInfo.ExpiresAt) {
		// 标记会话为过期
		s.expireSession(ctx, sessionID)
		return nil, fmt.Errorf("会话已过期")
	}

	return sessionInfo, nil
}

// ValidateToken 验证Token
func (s *sessionService) ValidateToken(ctx context.Context, accessToken string) (*SessionInfo, error) {
	// 从Token映射获取SessionID
	tokenHash := hashTokenString(accessToken)
	sessionID, err := s.cache.Get(ctx, fmt.Sprintf("token:access:%s", tokenHash))
	if err != nil {
		return nil, fmt.Errorf("Token不存在或已过期")
	}

	// 验证JWT Token
	claims, err := s.validateJWTToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("Token验证失败: %v", err)
	}

	// 检查Token类型
	if tokenType, ok := claims["type"].(string); !ok || tokenType != "access" {
		return nil, fmt.Errorf("无效的Token类型")
	}

	// 验证会话
	sessionInfo, err := s.ValidateSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// 更新最后活跃时间
	go s.UpdateLastActive(context.Background(), sessionID)

	return sessionInfo, nil
}

// UpdateLastActive 更新最后活跃时间
func (s *sessionService) UpdateLastActive(ctx context.Context, sessionID string) error {
	now := time.Now()

	// 更新Redis中的最后活跃时间
	key := fmt.Sprintf("session:%s", sessionID)
	err := s.cache.HSet(ctx, key, "last_active_at", now.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("更新Redis最后活跃时间失败: %v", err)
	}

	// 更新活跃会话排行
	err = s.cache.ZAdd(ctx, "sessions:active", float64(now.Unix()), sessionID)
	if err != nil {
		log.Printf("更新活跃会话排行失败: %v", err)
	}

	// 异步更新数据库
	go func() {
		err := s.sessionRepo.UpdateLastActive(context.Background(), sessionID, now)
		if err != nil {
			log.Printf("更新数据库最后活跃时间失败: %v", err)
		}
	}()

	return nil
}

// RefreshSession 刷新会话
func (s *sessionService) RefreshSession(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// 验证refresh token
	claims, err := s.validateJWTToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("刷新令牌验证失败: %v", err)
	}

	// 检查令牌类型
	if tokenType, ok := claims["type"].(string); !ok || tokenType != "refresh" {
		return nil, fmt.Errorf("无效的令牌类型")
	}

	// 获取会话信息
	sessionID := claims["session_id"].(string)
	userID := uint64(claims["user_id"].(float64))
	platform := claims["platform"].(string)

	sessionInfo, err := s.ValidateSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// 生成新的Token对
	newAccessToken, newRefreshToken, err := s.generateTokens(userID, sessionID, platform)
	if err != nil {
		return nil, fmt.Errorf("生成新令牌失败: %v", err)
	}

	now := time.Now()

	// 更新Redis中的Token信息
	sessionKey := fmt.Sprintf("session:%s", sessionID)
	updates := map[string]interface{}{
		"access_token":       newAccessToken,
		"refresh_token":      newRefreshToken,
		"access_expires_at":  now.Add(s.accessTokenTTL).Format(time.RFC3339),
		"refresh_expires_at": now.Add(s.refreshTokenTTL).Format(time.RFC3339),
		"last_active_at":     now.Format(time.RFC3339),
	}

	for field, value := range updates {
		err = s.cache.HSet(ctx, sessionKey, field, value)
		if err != nil {
			log.Printf("更新会话字段 %s 失败: %v", field, err)
		}
	}

	// 清除旧Token映射
	oldAccessHash := hashTokenString(sessionInfo.TokenInfo.AccessToken)
	oldRefreshHash := hashTokenString(sessionInfo.TokenInfo.RefreshToken)
	s.cache.Delete(ctx, fmt.Sprintf("token:access:%s", oldAccessHash))
	s.cache.Delete(ctx, fmt.Sprintf("token:refresh:%s", oldRefreshHash))

	// 创建新Token映射
	newAccessHash := hashTokenString(newAccessToken)
	newRefreshHash := hashTokenString(newRefreshToken)
	s.cache.Set(ctx, fmt.Sprintf("token:access:%s", newAccessHash), sessionID, s.accessTokenTTL)
	s.cache.Set(ctx, fmt.Sprintf("token:refresh:%s", newRefreshHash), sessionID, s.refreshTokenTTL)

	// 异步更新数据库
	go func() {
		err := s.sessionRepo.UpdateTokens(context.Background(), sessionID, hashTokenString(newAccessToken), hashTokenString(newRefreshToken), now.Add(s.accessTokenTTL), now.Add(s.refreshTokenTTL))
		if err != nil {
			log.Printf("更新数据库Token信息失败: %v", err)
		}
	}()

	// 发布Token刷新事件
	s.publishSessionEvent(ctx, EventTokenRefreshed, sessionID, userID, map[string]interface{}{
		"old_access_token": sessionInfo.TokenInfo.AccessToken,
		"new_access_token": newAccessToken,
		"refresh_time":     now,
	})

	return &TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.accessTokenTTL.Seconds()),
	}, nil
}

// RevokeSession 注销会话
func (s *sessionService) RevokeSession(ctx context.Context, sessionID string) error {
	// 获取会话信息
	sessionInfo, err := s.getSessionFromRedis(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("获取会话信息失败: %v", err)
	}

	// 从Redis删除会话数据
	err = s.deleteSessionFromRedis(ctx, sessionInfo)
	if err != nil {
		log.Printf("从Redis删除会话失败: %v", err)
	}

	// 更新数据库会话状态
	err = s.sessionRepo.UpdateStatus(ctx, sessionID, "revoked", time.Now())
	if err != nil {
		log.Printf("更新数据库会话状态失败: %v", err)
	}

	// 发布会话注销事件
	s.publishSessionEvent(ctx, EventSessionRevoked, sessionID, sessionInfo.UserID, nil)

	return nil
}

// RevokeUserSessions 注销用户的其他会话
func (s *sessionService) RevokeUserSessions(ctx context.Context, userID uint64, excludeSessionID string) error {
	// 获取用户所有会话
	sessionIDs, err := s.cache.SMembers(ctx, fmt.Sprintf("user:sessions:%d", userID))
	if err != nil {
		return fmt.Errorf("获取用户会话列表失败: %v", err)
	}

	for _, sessionID := range sessionIDs {
		if sessionID != excludeSessionID {
			err = s.RevokeSession(ctx, sessionID)
			if err != nil {
				log.Printf("注销会话 %s 失败: %v", sessionID, err)
			}
		}
	}

	return nil
}

// RevokeAllSessions 注销用户所有会话
func (s *sessionService) RevokeAllSessions(ctx context.Context, userID uint64) error {
	return s.RevokeUserSessions(ctx, userID, "")
}

// GetUserSessions 获取用户会话列表
func (s *sessionService) GetUserSessions(ctx context.Context, userID uint64) ([]*SessionInfo, error) {
	// 获取用户所有会话ID
	sessionIDs, err := s.cache.SMembers(ctx, fmt.Sprintf("user:sessions:%d", userID))
	if err != nil {
		return nil, fmt.Errorf("获取用户会话列表失败: %v", err)
	}

	var sessions []*SessionInfo
	for _, sessionID := range sessionIDs {
		sessionInfo, err := s.getSessionFromRedis(ctx, sessionID)
		if err != nil {
			log.Printf("获取会话 %s 信息失败: %v", sessionID, err)
			continue
		}
		sessions = append(sessions, sessionInfo)
	}

	return sessions, nil
}

// GetActiveSessionCount 获取用户活跃会话数量
func (s *sessionService) GetActiveSessionCount(ctx context.Context, userID uint64) (int, error) {
	sessionIDs, err := s.cache.SMembers(ctx, fmt.Sprintf("user:sessions:%d", userID))
	if err != nil {
		return 0, fmt.Errorf("获取用户会话列表失败: %v", err)
	}

	activeCount := 0
	for _, sessionID := range sessionIDs {
		exists, err := s.cache.Exists(ctx, fmt.Sprintf("session:%s", sessionID))
		if err == nil && exists {
			activeCount++
		}
	}

	return activeCount, nil
}

// CleanupExpiredSessions 清理过期会话
func (s *sessionService) CleanupExpiredSessions(ctx context.Context) error {
	// 从活跃会话排行中获取所有会话
	sessions, err := s.cache.ZRange(ctx, "sessions:active", 0, -1)
	if err != nil {
		return fmt.Errorf("获取活跃会话列表失败: %v", err)
	}

	now := time.Now()
	expiredCount := 0

	for _, sessionID := range sessions {
		sessionInfo, err := s.getSessionFromRedis(ctx, sessionID)
		if err != nil {
			// 会话不存在，从排行中移除
			s.cache.ZRem(ctx, "sessions:active", sessionID)
			continue
		}

		// 检查是否过期
		if now.After(sessionInfo.ExpiresAt) {
			err = s.expireSession(ctx, sessionID)
			if err != nil {
				log.Printf("过期会话 %s 失败: %v", sessionID, err)
			} else {
				expiredCount++
			}
		}
	}

	log.Printf("清理了 %d 个过期会话", expiredCount)
	return nil
}

// CleanupUserSessions 清理用户会话（保留指定数量）
func (s *sessionService) CleanupUserSessions(ctx context.Context, userID uint64, keepCount int) error {
	// 获取用户所有会话，按最后活跃时间排序
	sessionIDs, err := s.cache.SMembers(ctx, fmt.Sprintf("user:sessions:%d", userID))
	if err != nil {
		return fmt.Errorf("获取用户会话列表失败: %v", err)
	}

	if len(sessionIDs) <= keepCount {
		return nil
	}

	// 按最后活跃时间排序
	type sessionWithTime struct {
		SessionID    string
		LastActiveAt time.Time
	}

	var sessions []sessionWithTime
	for _, sessionID := range sessionIDs {
		sessionInfo, err := s.getSessionFromRedis(ctx, sessionID)
		if err != nil {
			continue
		}
		sessions = append(sessions, sessionWithTime{
			SessionID:    sessionID,
			LastActiveAt: sessionInfo.LastActiveAt,
		})
	}

	// 按时间排序，保留最新的
	if len(sessions) > keepCount {
		// 简单排序，实际应用中可以使用更高效的排序算法
		toRemove := len(sessions) - keepCount
		for i := 0; i < toRemove; i++ {
			oldest := 0
			for j := 1; j < len(sessions); j++ {
				if sessions[j].LastActiveAt.Before(sessions[oldest].LastActiveAt) {
					oldest = j
				}
			}

			// 注销最旧的会话
			err = s.RevokeSession(ctx, sessions[oldest].SessionID)
			if err != nil {
				log.Printf("注销旧会话 %s 失败: %v", sessions[oldest].SessionID, err)
			}

			// 从列表中移除
			sessions = append(sessions[:oldest], sessions[oldest+1:]...)
		}
	}

	return nil
}

// 辅助方法

// storeSessionToRedis 存储会话到Redis
func (s *sessionService) storeSessionToRedis(ctx context.Context, sessionInfo *SessionInfo) error {
	sessionKey := fmt.Sprintf("session:%s", sessionInfo.SessionID)
	userSessionsKey := fmt.Sprintf("user:sessions:%d", sessionInfo.UserID)

	// 存储会话详情
	sessionData := map[string]interface{}{
		"user_id":            sessionInfo.UserID,
		"platform":           sessionInfo.Platform,
		"device_id":          sessionInfo.DeviceInfo.DeviceID,
		"device_type":        sessionInfo.DeviceInfo.DeviceType,
		"device_name":        sessionInfo.DeviceInfo.DeviceName,
		"app_version":        sessionInfo.DeviceInfo.AppVersion,
		"user_agent":         sessionInfo.DeviceInfo.UserAgent,
		"ip_address":         sessionInfo.NetworkInfo.IPAddress,
		"location":           sessionInfo.NetworkInfo.Location,
		"access_token":       sessionInfo.TokenInfo.AccessToken,
		"refresh_token":      sessionInfo.TokenInfo.RefreshToken,
		"access_expires_at":  sessionInfo.TokenInfo.AccessExpiresAt.Format(time.RFC3339),
		"refresh_expires_at": sessionInfo.TokenInfo.RefreshExpiresAt.Format(time.RFC3339),
		"status":             sessionInfo.Status,
		"created_at":         sessionInfo.CreatedAt.Format(time.RFC3339),
		"last_active_at":     sessionInfo.LastActiveAt.Format(time.RFC3339),
		"expires_at":         sessionInfo.ExpiresAt.Format(time.RFC3339),
	}

	for field, value := range sessionData {
		err := s.cache.HSet(ctx, sessionKey, field, value)
		if err != nil {
			return fmt.Errorf("设置会话字段 %s 失败: %v", field, err)
		}
	}

	// 设置会话过期时间
	err := s.cache.Expire(ctx, sessionKey, s.refreshTokenTTL)
	if err != nil {
		return fmt.Errorf("设置会话过期时间失败: %v", err)
	}

	// 添加到用户会话集合
	err = s.cache.SAdd(ctx, userSessionsKey, sessionInfo.SessionID)
	if err != nil {
		return fmt.Errorf("添加到用户会话集合失败: %v", err)
	}

	// 设置用户会话集合过期时间
	err = s.cache.Expire(ctx, userSessionsKey, s.refreshTokenTTL)
	if err != nil {
		log.Printf("设置用户会话集合过期时间失败: %v", err)
	}

	// 创建Token映射
	accessTokenHash := hashTokenString(sessionInfo.TokenInfo.AccessToken)
	refreshTokenHash := hashTokenString(sessionInfo.TokenInfo.RefreshToken)

	err = s.cache.Set(ctx, fmt.Sprintf("token:access:%s", accessTokenHash), sessionInfo.SessionID, s.accessTokenTTL)
	if err != nil {
		return fmt.Errorf("创建访问令牌映射失败: %v", err)
	}

	err = s.cache.Set(ctx, fmt.Sprintf("token:refresh:%s", refreshTokenHash), sessionInfo.SessionID, s.refreshTokenTTL)
	if err != nil {
		return fmt.Errorf("创建刷新令牌映射失败: %v", err)
	}

	// 添加到活跃会话排行
	err = s.cache.ZAdd(ctx, "sessions:active", float64(sessionInfo.LastActiveAt.Unix()), sessionInfo.SessionID)
	if err != nil {
		log.Printf("添加到活跃会话排行失败: %v", err)
	}

	return nil
}

// getSessionFromRedis 从Redis获取会话信息
func (s *sessionService) getSessionFromRedis(ctx context.Context, sessionID string) (*SessionInfo, error) {
	sessionKey := fmt.Sprintf("session:%s", sessionID)

	// 检查会话是否存在
	exists, err := s.cache.Exists(ctx, sessionKey)
	if err != nil {
		return nil, fmt.Errorf("检查会话存在性失败: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("会话不存在")
	}

	// 获取所有会话字段
	sessionData, err := s.cache.HGetAll(ctx, sessionKey)
	if err != nil {
		return nil, fmt.Errorf("获取会话数据失败: %v", err)
	}

	// 解析会话信息
	sessionInfo := &SessionInfo{
		SessionID:   sessionID,
		DeviceInfo:  &DeviceInfo{},
		NetworkInfo: &NetworkInfo{},
		TokenInfo:   &TokenInfo{},
	}

	// 解析基本字段
	if userIDStr, ok := sessionData["user_id"]; ok {
		var userID uint64
		fmt.Sscanf(userIDStr, "%d", &userID)
		sessionInfo.UserID = userID
	}

	sessionInfo.Platform = sessionData["platform"]
	sessionInfo.Status = sessionData["status"]

	// 解析设备信息
	sessionInfo.DeviceInfo.DeviceID = sessionData["device_id"]
	sessionInfo.DeviceInfo.DeviceType = sessionData["device_type"]
	sessionInfo.DeviceInfo.DeviceName = sessionData["device_name"]
	sessionInfo.DeviceInfo.AppVersion = sessionData["app_version"]
	sessionInfo.DeviceInfo.UserAgent = sessionData["user_agent"]

	// 解析网络信息
	sessionInfo.NetworkInfo.IPAddress = sessionData["ip_address"]
	sessionInfo.NetworkInfo.Location = sessionData["location"]

	// 解析Token信息
	sessionInfo.TokenInfo.AccessToken = sessionData["access_token"]
	sessionInfo.TokenInfo.RefreshToken = sessionData["refresh_token"]

	// 解析时间字段
	if createdAtStr, ok := sessionData["created_at"]; ok {
		sessionInfo.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	}
	if lastActiveAtStr, ok := sessionData["last_active_at"]; ok {
		sessionInfo.LastActiveAt, _ = time.Parse(time.RFC3339, lastActiveAtStr)
	}
	if expiresAtStr, ok := sessionData["expires_at"]; ok {
		sessionInfo.ExpiresAt, _ = time.Parse(time.RFC3339, expiresAtStr)
	}
	if accessExpiresAtStr, ok := sessionData["access_expires_at"]; ok {
		sessionInfo.TokenInfo.AccessExpiresAt, _ = time.Parse(time.RFC3339, accessExpiresAtStr)
	}
	if refreshExpiresAtStr, ok := sessionData["refresh_expires_at"]; ok {
		sessionInfo.TokenInfo.RefreshExpiresAt, _ = time.Parse(time.RFC3339, refreshExpiresAtStr)
	}

	return sessionInfo, nil
}

// deleteSessionFromRedis 从Redis删除会话数据
func (s *sessionService) deleteSessionFromRedis(ctx context.Context, sessionInfo *SessionInfo) error {
	sessionKey := fmt.Sprintf("session:%s", sessionInfo.SessionID)
	userSessionsKey := fmt.Sprintf("user:sessions:%d", sessionInfo.UserID)

	// 删除会话数据
	err := s.cache.Delete(ctx, sessionKey)
	if err != nil {
		log.Printf("删除会话数据失败: %v", err)
	}

	// 从用户会话集合中移除
	err = s.cache.SRem(ctx, userSessionsKey, sessionInfo.SessionID)
	if err != nil {
		log.Printf("从用户会话集合移除失败: %v", err)
	}

	// 删除Token映射
	accessTokenHash := hashTokenString(sessionInfo.TokenInfo.AccessToken)
	refreshTokenHash := hashTokenString(sessionInfo.TokenInfo.RefreshToken)

	s.cache.Delete(ctx, fmt.Sprintf("token:access:%s", accessTokenHash))
	s.cache.Delete(ctx, fmt.Sprintf("token:refresh:%s", refreshTokenHash))

	// 从活跃会话排行中移除
	err = s.cache.ZRem(ctx, "sessions:active", sessionInfo.SessionID)
	if err != nil {
		log.Printf("从活跃会话排行移除失败: %v", err)
	}

	return nil
}

// expireSession 过期会话
func (s *sessionService) expireSession(ctx context.Context, sessionID string) error {
	// 获取会话信息
	sessionInfo, err := s.getSessionFromRedis(ctx, sessionID)
	if err != nil {
		return err
	}

	// 删除Redis数据
	err = s.deleteSessionFromRedis(ctx, sessionInfo)
	if err != nil {
		log.Printf("删除过期会话Redis数据失败: %v", err)
	}

	// 更新数据库状态
	err = s.sessionRepo.UpdateStatus(ctx, sessionID, "expired", time.Now())
	if err != nil {
		log.Printf("更新过期会话数据库状态失败: %v", err)
	}

	// 发布会话过期事件
	s.publishSessionEvent(ctx, EventSessionExpired, sessionID, sessionInfo.UserID, nil)

	return nil
}

// generateTokens 生成JWT令牌对
func (s *sessionService) generateTokens(userID uint64, sessionID, platform string) (string, string, error) {
	now := time.Now()

	// 生成访问令牌
	accessClaims := jwt.MapClaims{
		"user_id":    userID,
		"session_id": sessionID,
		"platform":   platform,
		"type":       "access",
		"iat":        now.Unix(),
		"exp":        now.Add(s.accessTokenTTL).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("生成访问令牌失败: %v", err)
	}

	// 生成刷新令牌
	refreshClaims := jwt.MapClaims{
		"user_id":    userID,
		"session_id": sessionID,
		"platform":   platform,
		"type":       "refresh",
		"iat":        now.Unix(),
		"exp":        now.Add(s.refreshTokenTTL).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("生成刷新令牌失败: %v", err)
	}

	return accessTokenString, refreshTokenString, nil
}

// validateJWTToken 验证JWT令牌
func (s *sessionService) validateJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("无效的令牌")
}

// publishSessionEvent 发布会话事件
func (s *sessionService) publishSessionEvent(ctx context.Context, eventType, sessionID string, userID uint64, data interface{}) {
	event := SessionEvent{
		Type:      eventType,
		SessionID: sessionID,
		UserID:    userID,
		Data:      data,
		Timestamp: time.Now().Unix(),
		Source:    s.instanceID,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		log.Printf("序列化会话事件失败: %v", err)
		return
	}

	// 发布到Redis
	err = s.cache.Publish(ctx, "session:events", string(eventData))
	if err != nil {
		log.Printf("发布会话事件失败: %v", err)
	}
}

// 工具函数

// generateSessionUniqueID 生成会话ID
func generateSessionUniqueID() string {
	return fmt.Sprintf("sess_%s", uuid.New().String()[:16])
}

// generateBackendInstanceID 生成实例ID
func generateBackendInstanceID() string {
	return fmt.Sprintf("inst_%s", uuid.New().String()[:8])
}

// hashTokenString 对Token进行哈希
func hashTokenString(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// SessionEvent 会话事件
type SessionEvent struct {
	Type      string      `json:"type"`
	SessionID string      `json:"session_id,omitempty"`
	UserID    uint64      `json:"user_id"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
	Source    string      `json:"source"`
}

// 事件类型常量
const (
	EventSessionCreated = "session_created"
	EventSessionUpdated = "session_updated"
	EventSessionRevoked = "session_revoked"
	EventUserLogout     = "user_logout"
	EventTokenRefreshed = "token_refreshed"
	EventSessionExpired = "session_expired"
)
