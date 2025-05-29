package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"huinong-backend/internal/model"
	"huinong-backend/internal/repository"
)

// oaService OA服务实现
type oaService struct {
	oaRepo      repository.OARepository
	userRepo    repository.UserRepository
	loanRepo    repository.LoanRepository
	machineRepo repository.MachineRepository
	jwtSecret   string
}

// NewOAService 创建OA服务实例
func NewOAService(
	oaRepo repository.OARepository,
	userRepo repository.UserRepository,
	loanRepo repository.LoanRepository,
	machineRepo repository.MachineRepository,
	jwtSecret string,
) OAService {
	return &oaService{
		oaRepo:      oaRepo,
		userRepo:    userRepo,
		loanRepo:    loanRepo,
		machineRepo: machineRepo,
		jwtSecret:   jwtSecret,
	}
}

// ==================== OA用户管理 ====================

// CreateOAUser 创建OA用户
func (s *oaService) CreateOAUser(ctx context.Context, req *CreateOAUserRequest) (*model.OAUser, error) {
	// 检查用户名是否已存在
	existingUser, err := s.oaRepo.GetOAUserByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	// TODO: 添加按邮箱查询的方法

	// 生成密码哈希
	salt := generateSalt()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password+salt), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	oaUser := &model.OAUser{
		Username:     req.Username,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: string(passwordHash),
		Salt:         salt,
		RealName:     req.RealName,
		RoleID:       req.RoleID,
		Department:   req.Department,
		Position:     req.Position,
		Status:       "active",
	}

	err = s.oaRepo.CreateOAUser(ctx, oaUser)
	if err != nil {
		return nil, fmt.Errorf("创建OA用户失败: %v", err)
	}

	return oaUser, nil
}

// GetOAUser 获取OA用户
func (s *oaService) GetOAUser(ctx context.Context, userID uint64) (*model.OAUser, error) {
	return s.oaRepo.GetOAUserByID(ctx, userID)
}

// UpdateOAUser 更新OA用户
func (s *oaService) UpdateOAUser(ctx context.Context, userID uint64, req *UpdateOAUserRequest) error {
	oaUser, err := s.oaRepo.GetOAUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取OA用户失败: %v", err)
	}

	// 更新字段
	if req.Email != "" {
		oaUser.Email = req.Email
	}
	if req.Phone != "" {
		oaUser.Phone = req.Phone
	}
	if req.RealName != "" {
		oaUser.RealName = req.RealName
	}
	if req.RoleID > 0 {
		oaUser.RoleID = req.RoleID
	}
	if req.Department != "" {
		oaUser.Department = req.Department
	}
	if req.Position != "" {
		oaUser.Position = req.Position
	}
	if req.Status != "" {
		oaUser.Status = req.Status
	}

	return s.oaRepo.UpdateOAUser(ctx, oaUser)
}

// DeleteOAUser 删除OA用户
func (s *oaService) DeleteOAUser(ctx context.Context, userID uint64) error {
	return s.oaRepo.DeleteOAUser(ctx, userID)
}

// ListOAUsers 获取OA用户列表
func (s *oaService) ListOAUsers(ctx context.Context, req *ListOAUsersRequest) (*ListOAUsersResponse, error) {
	repoReq := &repository.ListOAUsersRequest{
		RoleID:     req.RoleID,
		Department: req.Department,
		Status:     req.Status,
		Keyword:    req.Keyword,
		Page:       req.Page,
		Limit:      req.Limit,
	}

	repoResp, err := s.oaRepo.ListOAUsers(ctx, repoReq)
	if err != nil {
		return nil, fmt.Errorf("获取OA用户列表失败: %v", err)
	}

	return &ListOAUsersResponse{
		Users: repoResp.Users,
		Total: repoResp.Total,
		Page:  repoResp.Page,
		Limit: repoResp.Limit,
	}, nil
}

// ==================== 角色管理 ====================

// CreateRole 创建角色
func (s *oaService) CreateRole(ctx context.Context, req *CreateRoleRequest) (*model.OARole, error) {
	// 检查角色名是否已存在
	existingRole, err := s.oaRepo.GetRoleByName(ctx, req.Name)
	if err == nil && existingRole != nil {
		return nil, errors.New("角色名已存在")
	}

	role := &model.OARole{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		IsSuper:     req.IsSuper,
		Status:      "active",
		// TODO: 序列化权限
		// Permissions: marshalToJSON(req.Permissions),
	}

	err = s.oaRepo.CreateRole(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("创建角色失败: %v", err)
	}

	return role, nil
}

// GetRole 获取角色
func (s *oaService) GetRole(ctx context.Context, roleID uint64) (*model.OARole, error) {
	return s.oaRepo.GetRoleByID(ctx, roleID)
}

// UpdateRole 更新角色
func (s *oaService) UpdateRole(ctx context.Context, roleID uint64, req *UpdateRoleRequest) error {
	role, err := s.oaRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("获取角色失败: %v", err)
	}

	// 更新字段
	if req.DisplayName != "" {
		role.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Status != "" {
		role.Status = req.Status
	}

	// TODO: 更新权限
	// if req.Permissions != nil {
	//     role.Permissions = marshalToJSON(req.Permissions)
	// }

	return s.oaRepo.UpdateRole(ctx, role)
}

// DeleteRole 删除角色
func (s *oaService) DeleteRole(ctx context.Context, roleID uint64) error {
	return s.oaRepo.DeleteRole(ctx, roleID)
}

// ListRoles 获取角色列表
func (s *oaService) ListRoles(ctx context.Context) ([]*model.OARole, error) {
	return s.oaRepo.ListRoles(ctx)
}

// ==================== 认证管理 ====================

// OALogin OA用户登录
func (s *oaService) OALogin(ctx context.Context, req *OALoginRequest) (*OALoginResponse, error) {
	// 根据用户名查找用户
	oaUser, err := s.oaRepo.GetOAUserByUsername(ctx, req.Username)
	if err != nil || oaUser == nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if oaUser.Status != "active" {
		return nil, errors.New("用户账户已被禁用")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(oaUser.PasswordHash), []byte(req.Password+oaUser.Salt))
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成JWT令牌
	userService := &userService{jwtSecret: s.jwtSecret}
	accessToken, refreshToken, err := userService.generateTokens(oaUser.ID, req.Platform)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %v", err)
	}

	// 更新登录信息
	now := time.Now()
	oaUser.LastLoginAt = &now
	oaUser.LastLoginIP = getIPFromContext(ctx)
	oaUser.LoginCount++

	// TODO: 保存会话信息

	s.oaRepo.UpdateOAUser(ctx, oaUser)

	return &OALoginResponse{
		User:         oaUser,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    86400, // 24小时
		SessionID:    generateSessionID(),
	}, nil
}

// OALogout OA用户登出
func (s *oaService) OALogout(ctx context.Context, sessionID string) error {
	// TODO: 实现OA用户登出逻辑
	// 1. 标记会话为已注销
	// 2. 清除相关缓存

	return nil
}

// ==================== 工作台 ====================

// GetDashboard 获取工作台数据
func (s *oaService) GetDashboard(ctx context.Context, userID uint64) (*DashboardResponse, error) {
	dashboard := &DashboardResponse{
		UserStats:        make(map[string]int64),
		LoanStats:        make(map[string]int64),
		MachineStats:     make(map[string]int64),
		RecentActivities: []ActivityLog{},
	}

	// 用户统计
	userCount, _ := s.userRepo.GetUserCount(ctx)
	dashboard.UserStats["total"] = userCount

	// 按用户类型统计
	farmerCount, _ := s.userRepo.GetUserCountByType(ctx, "farmer")
	ownerCount, _ := s.userRepo.GetUserCountByType(ctx, "farm_owner")
	dashboard.UserStats["farmer"] = farmerCount
	dashboard.UserStats["farm_owner"] = ownerCount

	// 贷款统计
	// TODO: 实现贷款统计
	// applications, _ := s.loanRepo.GetApplicationStatistics(ctx)
	// dashboard.LoanStats = applications

	// 农机统计
	machineCount, _ := s.machineRepo.GetMachineCount(ctx)
	availableMachineCount, _ := s.machineRepo.GetAvailableMachineCount(ctx)
	orderCount, _ := s.machineRepo.GetOrderCount(ctx)

	dashboard.MachineStats["total"] = machineCount
	dashboard.MachineStats["available"] = availableMachineCount
	dashboard.MachineStats["orders"] = orderCount

	// TODO: 获取最近活动
	// activities, _ := s.getRecentActivities(ctx, 10)
	// dashboard.RecentActivities = activities

	return dashboard, nil
}

// GetPendingTasks 获取待处理任务
func (s *oaService) GetPendingTasks(ctx context.Context, userID uint64) (*PendingTasksResponse, error) {
	response := &PendingTasksResponse{
		LoanApprovals:       []PendingLoanApproval{},
		UserAuthentications: []PendingUserAuth{},
	}

	// 获取待审批的贷款申请
	// TODO: 实现获取待审批贷款
	// pendingApplications, _ := s.loanRepo.GetPendingApplications(ctx, 20, 0)
	// for _, app := range pendingApplications {
	//     approval := PendingLoanApproval{
	//         ID:            uint64(app.ID),
	//         ApplicationNo: app.ApplicationNo,
	//         UserName:      app.UserName,
	//         ProductName:   app.ProductName,
	//         Amount:        app.LoanAmount,
	//         SubmittedAt:   app.CreatedAt,
	//     }
	//     response.LoanApprovals = append(response.LoanApprovals, approval)
	// }

	// 获取待审核的用户认证
	// TODO: 实现获取待审核用户认证

	response.TotalCount = int64(len(response.LoanApprovals) + len(response.UserAuthentications))

	return response, nil
}
