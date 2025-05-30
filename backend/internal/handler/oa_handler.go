package handler

import (
	"huinong-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OAHandler OA后台管理处理器
type OAHandler struct {
	oaService   service.OAService
	userService service.UserService
}

// NewOAHandler 创建OA处理器
func NewOAHandler(oaService service.OAService, userService service.UserService) *OAHandler {
	return &OAHandler{
		oaService:   oaService,
		userService: userService,
	}
}

// Login OA管理员登录
func (h *OAHandler) Login(c *gin.Context) {
	var req service.OALoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "参数错误", err.Error()))
		return
	}

	ctx := c.Request.Context()
	resp, err := h.oaService.OALogin(ctx, &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "登录失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("登录成功", resp))
}

// GetCaptcha 获取验证码
func (h *OAHandler) GetCaptcha(c *gin.Context) {
	// TODO: 实现验证码生成逻辑
	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", gin.H{
		"captcha_id":    "cap_123456",
		"captcha_image": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
	}))
}

// GetUsers 获取用户列表
func (h *OAHandler) GetUsers(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	userType := c.Query("user_type")
	keyword := c.Query("keyword")

	// TODO: 实现获取用户列表逻辑，需要在Service层添加相应方法
	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", gin.H{
		"total": 100,
		"page":  page,
		"limit": limit,
		"filters": gin.H{
			"user_types": []gin.H{
				{"value": "farmer", "label": "个体农户", "count": 800},
				{"value": "farm_owner", "label": "农场主", "count": 300},
				{"value": "cooperative", "label": "合作社", "count": 150},
			},
			"statuses": []gin.H{
				{"value": "active", "label": "正常", "count": 1100},
				{"value": "frozen", "label": "冻结", "count": 50},
				{"value": "deleted", "label": "删除", "count": 100},
			},
		},
		"users": []gin.H{
			{
				"id":                    10001,
				"phone":                 "13800138000",
				"real_name":             "张三",
				"user_type":             userType,
				"status":                status,
				"keyword":               keyword,
				"is_real_name_verified": true,
				"credit_score":          750,
				"total_loans":           3,
				"current_debt":          50000,
				"created_at":            "2023-08-15T10:00:00Z",
			},
		},
	}))
}

// GetUserDetail 获取用户详情
func (h *OAHandler) GetUserDetail(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "用户ID格式错误", err.Error()))
		return
	}

	// TODO: 实现获取用户详情逻辑，需要在Service层添加相应方法
	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", gin.H{
		"basic_info": gin.H{
			"id":        userID,
			"real_name": "张三",
			"phone":     "13800138000",
			"user_type": "farmer",
			"status":    "active",
		},
		"auth_status": gin.H{
			"is_real_name_verified": true,
			"is_bank_card_verified": true,
			"is_credit_verified":    false,
		},
		"credit_info": gin.H{
			"credit_score":    750,
			"credit_level":    "优秀",
			"balance":         5000,
			"available_limit": 450000,
		},
		"business_summary": gin.H{
			"total_loans":    3,
			"total_borrowed": 150000,
			"current_debt":   50000,
			"overdue_count":  0,
			"total_rentals":  15,
		},
	}))
}

// GetMachines 获取农机设备列表
func (h *OAHandler) GetMachines(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	category := c.Query("category")
	ownerType := c.Query("owner_type")

	// TODO: 实现获取农机设备列表逻辑
	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", gin.H{
		"total": 156,
		"page":  page,
		"limit": limit,
		"summary": gin.H{
			"total_machines":       156,
			"active_machines":      135,
			"rented_machines":      45,
			"maintenance_machines": 12,
		},
		"filters": gin.H{
			"categories": []gin.H{
				{"code": "tillage", "name": "耕地机械", "count": 68},
				{"code": "planting", "name": "播种机械", "count": 35},
				{"code": "harvesting", "name": "收获机械", "count": 53},
			},
			"status_options": []gin.H{
				{"value": "active", "label": "可用", "count": 135},
				{"value": "rented", "label": "租赁中", "count": 45},
				{"value": "maintenance", "label": "维护中", "count": 12},
			},
		},
		"machines": []gin.H{
			{
				"id":          10001,
				"name":        "约翰迪尔 6B-1204拖拉机",
				"category":    category,
				"status":      status,
				"owner_type":  ownerType,
				"daily_rate":  500,
				"utilization": 0.68,
				"owner": gin.H{
					"id":   2001,
					"name": "济南农机合作社",
					"type": "cooperative",
				},
				"created_at": "2023-03-15T08:00:00Z",
			},
		},
	}))
}

// GetMachineDetail 获取农机设备详情
func (h *OAHandler) GetMachineDetail(c *gin.Context) {
	machineID, err := strconv.ParseInt(c.Param("machine_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "设备ID格式错误", err.Error()))
		return
	}

	// TODO: 实现获取农机设备详情逻辑
	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", gin.H{
		"basic_info": gin.H{
			"id":            machineID,
			"name":          "约翰迪尔 6B-1204拖拉机",
			"brand":         "约翰迪尔",
			"model":         "6B-1204",
			"serial_number": "JD2022120401",
			"status":        "active",
			"condition":     "excellent",
		},
		"owner_info": gin.H{
			"id":   2001,
			"name": "济南农机合作社",
			"type": "cooperative",
		},
		"rental_statistics": gin.H{
			"total_orders":     25,
			"total_revenue":    90000,
			"utilization_rate": 0.68,
			"customer_rating":  4.8,
		},
	}))
}

// UpdateUserStatus 更新用户状态
func (h *OAHandler) UpdateUserStatus(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "用户ID格式错误", err.Error()))
		return
	}

	var req struct {
		Status     string `json:"status"`
		Reason     string `json:"reason"`
		Note       string `json:"note"`
		NotifyUser bool   `json:"notify_user"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "参数错误", err.Error()))
		return
	}

	// TODO: 实现更新用户状态逻辑
	c.JSON(http.StatusOK, NewSuccessResponse("状态更新成功", gin.H{
		"user_id":        userID,
		"old_status":     "active",
		"new_status":     req.Status,
		"operation_time": "2024-01-15T16:30:00Z",
		"operator":       "admin",
	}))
}

// BatchOperateUsers 批量操作用户
func (h *OAHandler) BatchOperateUsers(c *gin.Context) {
	var req struct {
		Operation   string  `json:"operation"`
		UserIDs     []int64 `json:"user_ids"`
		Reason      string  `json:"reason"`
		NotifyUsers bool    `json:"notify_users"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "参数错误", err.Error()))
		return
	}

	// TODO: 实现批量操作逻辑
	c.JSON(http.StatusOK, NewSuccessResponse("批量操作成功", gin.H{
		"operation":      req.Operation,
		"affected_count": len(req.UserIDs),
		"success_count":  len(req.UserIDs),
		"failed_count":   0,
	}))
}

// GetDashboard 获取工作台数据
func (h *OAHandler) GetDashboard(c *gin.Context) {
	// TODO: 实现获取工作台数据逻辑
	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", gin.H{
		"user_statistics": gin.H{
			"total_users":      12500,
			"new_users_today":  25,
			"active_users":     8500,
			"user_growth_rate": 0.12,
		},
		"loan_statistics": gin.H{
			"total_applications":   156,
			"pending_applications": 35,
			"approved_today":       12,
			"total_amount":         15600000,
			"approval_rate":        0.78,
		},
		"machine_statistics": gin.H{
			"total_machines":   850,
			"active_rentals":   125,
			"rental_revenue":   456000,
			"utilization_rate": 0.68,
		},
		"financial_overview": gin.H{
			"total_revenue":   2500000,
			"monthly_revenue": 450000,
			"overdue_amount":  125000,
			"bad_debt_rate":   0.025,
		},
	}))
}

// GetRiskMonitoring 获取风险监控数据
func (h *OAHandler) GetRiskMonitoring(c *gin.Context) {
	// TODO: 实现获取风险监控数据逻辑
	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", gin.H{
		"risk_alerts": []gin.H{
			{
				"type":        "overdue_alert",
				"level":       "high",
				"title":       "逾期预警",
				"description": "有3笔贷款即将逾期",
				"count":       3,
				"created_at":  "2024-01-15T09:00:00Z",
			},
		},
		"overdue_analysis": gin.H{
			"overdue_count":  15,
			"overdue_amount": 125000,
			"overdue_rate":   0.025,
		},
		"risk_distribution": gin.H{
			"low_risk":    85,
			"medium_risk": 12,
			"high_risk":   3,
		},
	}))
}
