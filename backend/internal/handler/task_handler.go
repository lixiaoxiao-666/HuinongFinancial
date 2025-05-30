package handler

import (
	"huinong-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TaskHandler 任务处理器
type TaskHandler struct {
	taskService service.TaskService
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// CreateTask 创建任务
// @Summary 创建任务
// @Description 创建新的待处理任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param task body service.CreateTaskRequest true "任务信息"
// @Success 200 {object} ResponseData{data=service.CreateTaskResponse}
// @Failure 400 {object} ErrorResponse
// @Router /api/oa/admin/tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req service.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "参数错误", err.Error()))
		return
	}

	ctx := c.Request.Context()
	resp, err := h.taskService.CreateTask(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "创建任务失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("任务创建成功", resp))
}

// GetTask 获取任务详情
// @Summary 获取任务详情
// @Description 根据任务ID获取任务详细信息
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path uint64 true "任务ID"
// @Success 200 {object} ResponseData{data=service.TaskDetailResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/oa/admin/tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "任务ID格式错误", err.Error()))
		return
	}

	ctx := c.Request.Context()
	task, err := h.taskService.GetTask(ctx, taskID)
	if err != nil {
		if err.Error() == "任务不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "任务不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取任务失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", task))
}

// UpdateTask 更新任务
// @Summary 更新任务
// @Description 更新任务信息
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path uint64 true "任务ID"
// @Param task body service.UpdateTaskRequest true "更新信息"
// @Success 200 {object} ResponseData
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/oa/admin/tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "任务ID格式错误", err.Error()))
		return
	}

	var req service.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "参数错误", err.Error()))
		return
	}

	ctx := c.Request.Context()
	if err := h.taskService.UpdateTask(ctx, taskID, &req); err != nil {
		if err.Error() == "任务不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "任务不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "更新任务失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("任务更新成功", nil))
}

// DeleteTask 删除任务
// @Summary 删除任务
// @Description 删除指定任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path uint64 true "任务ID"
// @Success 200 {object} ResponseData
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/oa/admin/tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "任务ID格式错误", err.Error()))
		return
	}

	ctx := c.Request.Context()
	if err := h.taskService.DeleteTask(ctx, taskID); err != nil {
		if err.Error() == "任务不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "任务不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "删除任务失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("任务删除成功", nil))
}

// ListTasks 获取任务列表
// @Summary 获取任务列表
// @Description 获取待处理任务列表，支持分页和筛选
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Param type query string false "任务类型"
// @Param status query string false "任务状态"
// @Param priority query string false "任务优先级"
// @Param assigned_to query uint64 false "分配人ID"
// @Param created_by query uint64 false "创建人ID"
// @Param business_type query string false "业务类型"
// @Param business_id query uint64 false "业务ID"
// @Param is_overdue query bool false "是否超时"
// @Param keyword query string false "关键词搜索"
// @Param sort_by query string false "排序字段" default(created_at)
// @Param sort_order query string false "排序方向" default(desc)
// @Success 200 {object} ResponseData{data=service.ListTasksResponse}
// @Failure 400 {object} ErrorResponse
// @Router /api/oa/admin/tasks [get]
func (h *TaskHandler) ListTasks(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	taskType := c.Query("type")
	status := c.Query("status")
	priority := c.Query("priority")
	businessType := c.Query("business_type")
	keyword := c.Query("keyword")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	var assignedTo *uint64
	if assignedToStr := c.Query("assigned_to"); assignedToStr != "" {
		if id, err := strconv.ParseUint(assignedToStr, 10, 64); err == nil {
			assignedTo = &id
		}
	}

	var createdBy *uint64
	if createdByStr := c.Query("created_by"); createdByStr != "" {
		if id, err := strconv.ParseUint(createdByStr, 10, 64); err == nil {
			createdBy = &id
		}
	}

	var businessID *uint64
	if businessIDStr := c.Query("business_id"); businessIDStr != "" {
		if id, err := strconv.ParseUint(businessIDStr, 10, 64); err == nil {
			businessID = &id
		}
	}

	var isOverdue *bool
	if isOverdueStr := c.Query("is_overdue"); isOverdueStr != "" {
		if overdue, err := strconv.ParseBool(isOverdueStr); err == nil {
			isOverdue = &overdue
		}
	}

	req := &service.ListTasksRequest{
		Page:         page,
		Limit:        limit,
		Type:         taskType,
		Status:       status,
		Priority:     priority,
		AssignedTo:   assignedTo,
		CreatedBy:    createdBy,
		BusinessType: businessType,
		BusinessID:   businessID,
		IsOverdue:    isOverdue,
		Keyword:      keyword,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
	}

	ctx := c.Request.Context()
	resp, err := h.taskService.ListTasks(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取任务列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", resp))
}

// GetPendingTasks 获取待处理任务列表（兼容现有接口）
// @Summary 获取待处理任务列表
// @Description 获取当前用户的待处理任务，用于工作台显示
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param limit query int false "数量限制" default(10)
// @Success 200 {object} ResponseData{data=service.ListTasksResponse}
// @Failure 400 {object} ErrorResponse
// @Router /api/oa/admin/tasks/pending [get]
func (h *TaskHandler) GetPendingTasks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := &service.ListTasksRequest{
		Page:      1,
		Limit:     limit,
		Status:    "pending", // 只获取待处理任务
		SortBy:    "priority,created_at",
		SortOrder: "desc",
	}

	ctx := c.Request.Context()
	resp, err := h.taskService.ListTasks(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取待处理任务失败", err.Error()))
		return
	}

	// 为了兼容原有接口格式，转换数据结构
	pendingTasks := gin.H{
		"tasks":       resp.Tasks,
		"total_count": resp.Total,
		"statistics":  resp.Statistics,
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", pendingTasks))
}

// ProcessTask 处理任务
// @Summary 处理任务
// @Description 处理指定任务（开始处理、完成、取消等操作）
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path uint64 true "任务ID"
// @Param request body service.ProcessTaskRequest true "处理请求"
// @Success 200 {object} ResponseData{data=service.ProcessTaskResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/oa/admin/tasks/{id}/process [post]
func (h *TaskHandler) ProcessTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "任务ID格式错误", err.Error()))
		return
	}

	var req service.ProcessTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "参数错误", err.Error()))
		return
	}

	ctx := c.Request.Context()

	var response *service.ProcessTaskResponse
	var processErr error

	// 根据操作类型执行不同的处理逻辑
	switch req.Action {
	case "start":
		// 开始处理任务
		updateReq := &service.UpdateTaskRequest{
			Status: "processing",
		}
		processErr = h.taskService.UpdateTask(ctx, taskID, updateReq)
		if processErr == nil {
			response = &service.ProcessTaskResponse{
				Success:   true,
				Message:   "任务已开始处理",
				TaskID:    taskID,
				NewStatus: "processing",
			}
		}

	case "complete":
		// 完成任务
		processErr = h.taskService.CompleteTask(ctx, taskID)
		if processErr == nil {
			response = &service.ProcessTaskResponse{
				Success:   true,
				Message:   "任务已完成",
				TaskID:    taskID,
				NewStatus: "completed",
			}
		}

	case "cancel":
		// 取消任务
		updateReq := &service.UpdateTaskRequest{
			Status: "cancelled",
		}
		processErr = h.taskService.UpdateTask(ctx, taskID, updateReq)
		if processErr == nil {
			response = &service.ProcessTaskResponse{
				Success:   true,
				Message:   "任务已取消",
				TaskID:    taskID,
				NewStatus: "cancelled",
			}
		}

	case "progress":
		// 更新进度
		if req.Progress < 0 || req.Progress > 100 {
			c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "进度值错误", "进度值必须在0-100之间"))
			return
		}
		processErr = h.taskService.UpdateTaskProgress(ctx, taskID, req.Progress)
		if processErr == nil {
			newStatus := "processing"
			if req.Progress == 100 {
				newStatus = "completed"
			}
			response = &service.ProcessTaskResponse{
				Success:   true,
				Message:   "任务进度已更新",
				TaskID:    taskID,
				NewStatus: newStatus,
			}
		}

	default:
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "无效的操作类型", "支持的操作：start, complete, cancel, progress"))
		return
	}

	if processErr != nil {
		if processErr.Error() == "任务不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "任务不存在", processErr.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "处理任务失败", processErr.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("处理成功", response))
}

// AssignTask 分配任务
// @Summary 分配任务
// @Description 将任务分配给指定用户
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path uint64 true "任务ID"
// @Param assignee_id query uint64 true "分配目标用户ID"
// @Success 200 {object} ResponseData
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/oa/admin/tasks/{id}/assign [post]
func (h *TaskHandler) AssignTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "任务ID格式错误", err.Error()))
		return
	}

	assigneeIDStr := c.Query("assignee_id")
	if assigneeIDStr == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "参数错误", "assignee_id不能为空"))
		return
	}

	assigneeID, err := strconv.ParseUint(assigneeIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "用户ID格式错误", err.Error()))
		return
	}

	ctx := c.Request.Context()
	if err := h.taskService.AssignTask(ctx, taskID, assigneeID); err != nil {
		if err.Error() == "任务不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "任务不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "分配任务失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("任务分配成功", nil))
}

// UnassignTask 取消分配任务
// @Summary 取消分配任务
// @Description 取消任务的当前分配
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path uint64 true "任务ID"
// @Success 200 {object} ResponseData
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/oa/admin/tasks/{id}/unassign [post]
func (h *TaskHandler) UnassignTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "任务ID格式错误", err.Error()))
		return
	}

	ctx := c.Request.Context()
	if err := h.taskService.UnassignTask(ctx, taskID); err != nil {
		if err.Error() == "任务不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "任务不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "取消分配失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("取消分配成功", nil))
}

// ReassignTask 重新分配任务
// @Summary 重新分配任务
// @Description 将任务重新分配给其他用户
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path uint64 true "任务ID"
// @Param new_assignee_id query uint64 true "新分配目标用户ID"
// @Success 200 {object} ResponseData
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/oa/admin/tasks/{id}/reassign [post]
func (h *TaskHandler) ReassignTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "任务ID格式错误", err.Error()))
		return
	}

	newAssigneeIDStr := c.Query("new_assignee_id")
	if newAssigneeIDStr == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "参数错误", "new_assignee_id不能为空"))
		return
	}

	newAssigneeID, err := strconv.ParseUint(newAssigneeIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "用户ID格式错误", err.Error()))
		return
	}

	ctx := c.Request.Context()
	if err := h.taskService.ReassignTask(ctx, taskID, newAssigneeID); err != nil {
		if err.Error() == "任务不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "任务不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "重新分配失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("重新分配成功", nil))
}

// GetTaskProgress 获取任务进度
// @Summary 获取任务进度
// @Description 获取指定任务的进度信息
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path uint64 true "任务ID"
// @Success 200 {object} ResponseData{data=service.TaskProgressResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/oa/admin/tasks/{id}/progress [get]
func (h *TaskHandler) GetTaskProgress(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, "任务ID格式错误", err.Error()))
		return
	}

	ctx := c.Request.Context()
	progress, err := h.taskService.GetTaskProgress(ctx, taskID)
	if err != nil {
		if err.Error() == "任务不存在" {
			c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "任务不存在", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "获取任务进度失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("获取成功", progress))
}
