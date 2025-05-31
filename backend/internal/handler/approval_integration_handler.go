package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"huinong-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// ApprovalIntegrationHandler 审批集成处理器
type ApprovalIntegrationHandler struct {
	loanService    service.LoanService
	machineService service.MachineService
	taskService    service.TaskService
	oaService      service.OAService
}

// NewApprovalIntegrationHandler 创建审批集成处理器
func NewApprovalIntegrationHandler(
	loanService service.LoanService,
	machineService service.MachineService,
	taskService service.TaskService,
	oaService service.OAService,
) *ApprovalIntegrationHandler {
	return &ApprovalIntegrationHandler{
		loanService:    loanService,
		machineService: machineService,
		taskService:    taskService,
		oaService:      oaService,
	}
}

// ============= 统一审批接口 =============

// GetPendingApprovals 获取待审批事项
// @Summary 获取所有待审批事项
// @Description 获取贷款申请、农机租赁等所有待审批事项的统一列表
// @Tags 审批管理
// @Accept json
// @Produce json
// @Param type query string false "审批类型 (loan/machine_rental/all)" default(all)
// @Param priority query string false "优先级筛选 (low/normal/high/urgent)"
// @Param assigned_to query uint64 false "分配给指定审核员"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Success 200 {object} StandardResponse{data=ApprovalListResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/approvals/pending [get]
func (h *ApprovalIntegrationHandler) GetPendingApprovals(c *gin.Context) {
	approvalType := c.DefaultQuery("type", "all")
	priority := c.Query("priority")
	assignedToStr := c.Query("assigned_to")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	var assignedTo *uint64
	if assignedToStr != "" {
		if id, err := strconv.ParseUint(assignedToStr, 10, 64); err == nil {
			assignedTo = &id
		}
	}

	// 构建统一的审批列表响应
	approvals := []map[string]interface{}{}

	// 获取贷款申请（如果类型匹配）
	if approvalType == "all" || approvalType == "loan" {
		loanReq := &service.GetAdminApplicationsRequest{
			Status: "pending",
			Page:   1, // 获取所有待审批项目进行合并
			Limit:  1000,
		}
		if loanApps, err := h.loanService.GetAdminApplications(c.Request.Context(), loanReq); err == nil {
			for _, app := range loanApps.Applications {
				approvals = append(approvals, map[string]interface{}{
					"id":           fmt.Sprintf("loan_%d", app.ID),
					"type":         "loan_application",
					"type_name":    "贷款申请",
					"title":        fmt.Sprintf("贷款申请 - %s", app.ApplicantName),
					"description":  fmt.Sprintf("申请金额: %.2f元，期限: %d个月", float64(app.ApplyAmount)/100, app.ApplyTermMonths),
					"business_id":  app.ID,
					"user_id":      app.UserID,
					"user_name":    app.ApplicantName,
					"amount":       float64(app.ApplyAmount) / 100,
					"status":       app.Status,
					"priority":     h.calculatePriority(app.RiskLevel, time.Since(app.CreatedAt)),
					"risk_level":   app.RiskLevel,
					"created_at":   app.CreatedAt,
					"due_date":     app.CreatedAt.Add(48 * time.Hour), // 48小时内处理
					"assigned_to":  app.AssignedTo,
					"is_overdue":   time.Since(app.CreatedAt) > 48*time.Hour,
				})
			}
		}
	}

	// 获取农机租赁申请（模拟数据，实际应调用service）
	if approvalType == "all" || approvalType == "machine_rental" {
		// 这里应该调用实际的农机租赁服务
		// 现在添加一些模拟数据
		rentalApprovals := []map[string]interface{}{
			{
				"id":           "rental_1001",
				"type":         "machine_rental",
				"type_name":    "农机租赁",
				"title":        "农机租赁申请 - 拖拉机",
				"description":  "租赁期限: 7天，租金: 1200元",
				"business_id":  1001,
				"user_id":      2001,
				"user_name":    "农户张三",
				"amount":       1200.0,
				"status":       "pending",
				"priority":     "normal",
				"risk_level":   "low",
				"created_at":   time.Now().Add(-3 * time.Hour),
				"due_date":     time.Now().Add(24 * time.Hour),
				"assigned_to":  nil,
				"is_overdue":   false,
			},
			{
				"id":           "rental_1002",
				"type":         "machine_rental",
				"type_name":    "农机租赁",
				"title":        "农机租赁申请 - 收割机",
				"description":  "租赁期限: 3天，租金: 2400元",
				"business_id":  1002,
				"user_id":      2002,
				"user_name":    "农户李四",
				"amount":       2400.0,
				"status":       "pending",
				"priority":     "high",
				"risk_level":   "medium",
				"created_at":   time.Now().Add(-26 * time.Hour),
				"due_date":     time.Now().Add(22 * time.Hour),
				"assigned_to":  nil,
				"is_overdue":   true,
			},
		}
		approvals = append(approvals, rentalApprovals...)
	}

	// 按优先级和时间排序
	// 这里简化处理，实际应该实现更复杂的排序逻辑

	// 应用筛选条件
	filteredApprovals := []map[string]interface{}{}
	for _, approval := range approvals {
		// 优先级筛选
		if priority != "" && approval["priority"] != priority {
			continue
		}

		// 分配人筛选
		if assignedTo != nil {
			if approval["assigned_to"] == nil {
				continue
			}
			if assignedID, ok := approval["assigned_to"].(uint64); !ok || assignedID != *assignedTo {
				continue
			}
		}

		filteredApprovals = append(filteredApprovals, approval)
	}

	// 分页处理
	total := len(filteredApprovals)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pagedApprovals := filteredApprovals[start:end]

	// 统计信息
	statistics := map[string]interface{}{
		"total_pending":    total,
		"overdue_count":    h.countOverdue(filteredApprovals),
		"high_priority":    h.countByPriority(filteredApprovals, "high"),
		"urgent_priority":  h.countByPriority(filteredApprovals, "urgent"),
		"unassigned_count": h.countUnassigned(filteredApprovals),
	}

	response := map[string]interface{}{
		"approvals": pagedApprovals,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + limit - 1) / limit,
		},
		"statistics": statistics,
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// BatchAssignApprovals 批量分配审批任务
// @Summary 批量分配审批任务
// @Description 将多个审批任务分配给指定审核员
// @Tags 审批管理
// @Accept json
// @Produce json
// @Param request body BatchAssignRequest true "批量分配信息"
// @Success 200 {object} StandardResponse{data=BatchAssignResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/approvals/batch-assign [post]
func (h *ApprovalIntegrationHandler) BatchAssignApprovals(c *gin.Context) {
	type BatchAssignRequest struct {
		ApprovalIDs []string `json:"approval_ids" validate:"required"`
		AssignedTo  uint64   `json:"assigned_to" validate:"required"`
		Priority    string   `json:"priority,omitempty"`
		DueDate     string   `json:"due_date,omitempty"`
		Notes       string   `json:"notes,omitempty"`
	}

	var req BatchAssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "请求参数错误", err.Error()))
		return
	}

	// 从上下文获取操作员ID
	operatorIDInterface, exists := c.Get("oaUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "操作员未登录", "OA用户认证信息缺失"))
		return
	}

	operatorID, ok := operatorIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "操作员ID格式错误", "operatorID type assertion failed"))
		return
	}

	successCount := 0
	failureCount := 0
	results := []map[string]interface{}{}

	for _, approvalID := range req.ApprovalIDs {
		result := map[string]interface{}{
			"approval_id": approvalID,
			"success":     false,
			"message":     "",
		}

		// 根据ID前缀判断审批类型并处理
		if len(approvalID) > 5 && approvalID[:5] == "loan_" {
			// 处理贷款申请分配
			if err := h.assignLoanApplication(approvalID[5:], req.AssignedTo, operatorID); err != nil {
				result["message"] = err.Error()
				failureCount++
			} else {
				result["success"] = true
				result["message"] = "分配成功"
				successCount++
			}
		} else if len(approvalID) > 7 && approvalID[:7] == "rental_" {
			// 处理农机租赁分配
			if err := h.assignMachineRental(approvalID[7:], req.AssignedTo, operatorID); err != nil {
				result["message"] = err.Error()
				failureCount++
			} else {
				result["success"] = true
				result["message"] = "分配成功"
				successCount++
			}
		} else {
			result["message"] = "未知的审批类型"
			failureCount++
		}

		results = append(results, result)
	}

	response := map[string]interface{}{
		"total_count":   len(req.ApprovalIDs),
		"success_count": successCount,
		"failure_count": failureCount,
		"results":       results,
		"assigned_to":   req.AssignedTo,
		"assigned_by":   operatorID,
		"assigned_at":   time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, NewSuccessResponse("批量分配完成", response))
}

// GetApprovalWorkload 获取审核员工作负载
// @Summary 获取审核员工作负载
// @Description 获取各审核员的工作负载统计
// @Tags 审批管理
// @Accept json
// @Produce json
// @Success 200 {object} StandardResponse{data=WorkloadResponse}
// @Failure 500 {object} ErrorResponse
// @Router /api/oa/admin/approvals/workload [get]
func (h *ApprovalIntegrationHandler) GetApprovalWorkload(c *gin.Context) {
	// 模拟审核员工作负载数据
	// 实际实现中应该从数据库查询
	workloadData := []map[string]interface{}{
		{
			"reviewer_id":      "OA001",
			"reviewer_name":    "张审核员",
			"department":       "贷款审批部",
			"speciality":       []string{"贷款审批", "风险评估"},
			"pending_count":    12,
			"in_progress_count": 3,
			"completed_today":  5,
			"completed_week":   28,
			"average_time_hours": 2.3,
			"efficiency_score": 92.5,
			"workload_level":   "中等",
			"last_activity":    time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
		},
		{
			"reviewer_id":      "OA002",
			"reviewer_name":    "李审核员",
			"department":       "贷款审批部",
			"speciality":       []string{"风险评估", "征信分析"},
			"pending_count":    8,
			"in_progress_count": 2,
			"completed_today":  7,
			"completed_week":   35,
			"average_time_hours": 1.8,
			"efficiency_score": 95.2,
			"workload_level":   "较低",
			"last_activity":    time.Now().Add(-15 * time.Minute).Format(time.RFC3339),
		},
		{
			"reviewer_id":      "OA003",
			"reviewer_name":    "王审核员",
			"department":       "农机管理部",
			"speciality":       []string{"农机管理", "设备审核"},
			"pending_count":    15,
			"in_progress_count": 4,
			"completed_today":  3,
			"completed_week":   18,
			"average_time_hours": 3.1,
			"efficiency_score": 88.7,
			"workload_level":   "较高",
			"last_activity":    time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
		},
	}

	// 计算总体统计
	totalPending := 0
	totalInProgress := 0
	totalCompletedToday := 0
	avgEfficiency := 0.0

	for _, reviewer := range workloadData {
		totalPending += reviewer["pending_count"].(int)
		totalInProgress += reviewer["in_progress_count"].(int)
		totalCompletedToday += reviewer["completed_today"].(int)
		avgEfficiency += reviewer["efficiency_score"].(float64)
	}

	if len(workloadData) > 0 {
		avgEfficiency /= float64(len(workloadData))
	}

	response := map[string]interface{}{
		"reviewers": workloadData,
		"summary": map[string]interface{}{
			"total_reviewers":      len(workloadData),
			"total_pending":        totalPending,
			"total_in_progress":    totalInProgress,
			"total_completed_today": totalCompletedToday,
			"average_efficiency":   avgEfficiency,
			"recommended_assignments": h.getRecommendedAssignments(workloadData),
		},
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", response))
}

// ============= 辅助函数 =============

// calculatePriority 计算审批优先级
func (h *ApprovalIntegrationHandler) calculatePriority(riskLevel string, duration time.Duration) string {
	hours := duration.Hours()

	if hours > 48 {
		return "urgent"
	}
	if hours > 24 {
		return "high"
	}
	if riskLevel == "high" {
		return "high"
	}
	if hours > 12 {
		return "normal"
	}
	return "low"
}

// countOverdue 统计超期数量
func (h *ApprovalIntegrationHandler) countOverdue(approvals []map[string]interface{}) int {
	count := 0
	for _, approval := range approvals {
		if isOverdue, ok := approval["is_overdue"].(bool); ok && isOverdue {
			count++
		}
	}
	return count
}

// countByPriority 按优先级统计
func (h *ApprovalIntegrationHandler) countByPriority(approvals []map[string]interface{}, priority string) int {
	count := 0
	for _, approval := range approvals {
		if p, ok := approval["priority"].(string); ok && p == priority {
			count++
		}
	}
	return count
}

// countUnassigned 统计未分配数量
func (h *ApprovalIntegrationHandler) countUnassigned(approvals []map[string]interface{}) int {
	count := 0
	for _, approval := range approvals {
		if approval["assigned_to"] == nil {
			count++
		}
	}
	return count
}

// assignLoanApplication 分配贷款申请
func (h *ApprovalIntegrationHandler) assignLoanApplication(loanIDStr string, assignedTo, operatorID uint64) error {
	// 这里应该调用贷款服务的分配方法
	// 现在返回模拟结果
	return nil
}

// assignMachineRental 分配农机租赁
func (h *ApprovalIntegrationHandler) assignMachineRental(rentalIDStr string, assignedTo, operatorID uint64) error {
	// 这里应该调用农机服务的分配方法
	// 现在返回模拟结果
	return nil
}

// getRecommendedAssignments 获取推荐分配
func (h *ApprovalIntegrationHandler) getRecommendedAssignments(workloadData []map[string]interface{}) []map[string]interface{} {
	recommendations := []map[string]interface{}{}

	for _, reviewer := range workloadData {
		if workloadLevel, ok := reviewer["workload_level"].(string); ok && workloadLevel == "较低" {
			recommendations = append(recommendations, map[string]interface{}{
				"reviewer_id":   reviewer["reviewer_id"],
				"reviewer_name": reviewer["reviewer_name"],
				"reason":        "工作负载较低，可承接更多任务",
				"capacity":      "可接受5-8个新任务",
			})
		}
	}

	return recommendations
}