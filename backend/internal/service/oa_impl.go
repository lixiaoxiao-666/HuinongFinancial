package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"huinong-backend/internal/model"
	"huinong-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

// oaServiceImpl OA服务实现
type oaServiceImpl struct {
	oaRepo         repository.OARepository
	userRepo       repository.UserRepository
	loanRepo       repository.LoanRepository
	machineRepo    repository.MachineRepository
	jwtSecret      string
	sessionService SessionService
}

// NewOAService 创建OA服务实例
func NewOAService(
	oaRepo repository.OARepository,
	userRepo repository.UserRepository,
	loanRepo repository.LoanRepository,
	machineRepo repository.MachineRepository,
	jwtSecret string,
	sessionService SessionService,
) OAService {
	return &oaServiceImpl{
		oaRepo:         oaRepo,
		userRepo:       userRepo,
		loanRepo:       loanRepo,
		machineRepo:    machineRepo,
		jwtSecret:      jwtSecret,
		sessionService: sessionService,
	}
}

// OALogin OA用户登录 - 集成Redis会话管理
func (s *oaServiceImpl) OALogin(ctx context.Context, req *OALoginRequest) (*OALoginResponse, error) {
	// 简单的用户名密码验证
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("用户名和密码不能为空")
	}

	// 模拟用户验证 - 开发环境允许的测试账户
	var mockUser *model.OAUser
	switch req.Username {
	case "admin":
		if req.Password != "admin123" {
			return nil, errors.New("用户名或密码错误")
		}
		mockUser = &model.OAUser{
			ID:         1,
			Username:   "admin",
			Email:      "admin@huinong.com",
			Phone:      "13800138000",
			RealName:   "超级管理员",
			RoleID:     1,
			Department: "技术部",
			Position:   "系统管理员",
			Status:     "active",
			LoginCount: 1,
			CreatedAt:  time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:  time.Now(),
		}
	case "reviewer":
		if req.Password != "reviewer123" {
			return nil, errors.New("用户名或密码错误")
		}
		mockUser = &model.OAUser{
			ID:         2,
			Username:   "reviewer",
			Email:      "reviewer@huinong.com",
			Phone:      "13800138001",
			RealName:   "审批员",
			RoleID:     2,
			Department: "业务部",
			Position:   "贷款审批员",
			Status:     "active",
			LoginCount: 1,
			CreatedAt:  time.Now().Add(-15 * 24 * time.Hour),
			UpdatedAt:  time.Now(),
		}
	default:
		return nil, errors.New("用户名或密码错误")
	}

	// 获取客户端IP
	clientIP := req.IPAddress
	if clientIP == "" {
		if ginCtx, ok := ctx.(*gin.Context); ok {
			clientIP = ginCtx.ClientIP()
		} else {
			clientIP = "127.0.0.1"
		}
	}

	// 使用Session Service创建Redis会话
	loginInfo := &LoginInfo{
		Platform:    "oa", // 固定为oa平台
		DeviceID:    req.DeviceID,
		DeviceType:  req.DeviceType,
		DeviceName:  req.DeviceName,
		AppVersion:  req.AppVersion,
		UserAgent:   req.DeviceName,
		IPAddress:   clientIP,
		Location:    req.Location,
		LoginMethod: "password",
	}

	// 创建Redis会话
	sessionInfo, err := s.sessionService.CreateSession(ctx, uint64(mockUser.ID), loginInfo)
	if err != nil {
		return nil, fmt.Errorf("创建会话失败: %v", err)
	}

	// 更新最后登录信息
	now := time.Now()
	mockUser.LastLoginAt = &now
	mockUser.LastLoginIP = clientIP

	return &OALoginResponse{
		User:         mockUser,
		AccessToken:  sessionInfo.TokenInfo.AccessToken,
		RefreshToken: sessionInfo.TokenInfo.RefreshToken,
		ExpiresIn:    int(sessionInfo.TokenInfo.AccessExpiresAt.Sub(now).Seconds()),
		SessionID:    sessionInfo.SessionID,
	}, nil
}

// OALogout OA用户登出
func (s *oaServiceImpl) OALogout(ctx context.Context, sessionID string) error {
	return s.sessionService.RevokeSession(ctx, sessionID)
}

// ==================== 其他接口的默认实现 ====================

func (s *oaServiceImpl) CreateOAUser(ctx context.Context, req *CreateOAUserRequest) (*model.OAUser, error) {
	return nil, errors.New("功能开发中")
}

func (s *oaServiceImpl) GetOAUser(ctx context.Context, userID uint64) (*model.OAUser, error) {
	return nil, errors.New("功能开发中")
}

func (s *oaServiceImpl) UpdateOAUser(ctx context.Context, userID uint64, req *UpdateOAUserRequest) error {
	return errors.New("功能开发中")
}

func (s *oaServiceImpl) DeleteOAUser(ctx context.Context, userID uint64) error {
	return errors.New("功能开发中")
}

func (s *oaServiceImpl) ListOAUsers(ctx context.Context, req *ListOAUsersRequest) (*ListOAUsersResponse, error) {
	return nil, errors.New("功能开发中")
}

func (s *oaServiceImpl) CreateRole(ctx context.Context, req *CreateRoleRequest) (*model.OARole, error) {
	return nil, errors.New("功能开发中")
}

func (s *oaServiceImpl) GetRole(ctx context.Context, roleID uint64) (*model.OARole, error) {
	return nil, errors.New("功能开发中")
}

func (s *oaServiceImpl) UpdateRole(ctx context.Context, roleID uint64, req *UpdateRoleRequest) error {
	return errors.New("功能开发中")
}

func (s *oaServiceImpl) DeleteRole(ctx context.Context, roleID uint64) error {
	return errors.New("功能开发中")
}

func (s *oaServiceImpl) ListRoles(ctx context.Context) ([]*model.OARole, error) {
	return nil, errors.New("功能开发中")
}

func (s *oaServiceImpl) GetDashboard(ctx context.Context, userID uint64) (*DashboardResponse, error) {
	return nil, errors.New("功能开发中")
}

func (s *oaServiceImpl) GetPendingTasks(ctx context.Context, userID uint64) (*PendingTasksResponse, error) {
	return nil, errors.New("功能开发中")
}
