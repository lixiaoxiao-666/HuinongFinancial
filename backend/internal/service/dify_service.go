package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"huinong-backend/internal/config"
	"huinong-backend/internal/model"
	"huinong-backend/internal/repository"
	"huinong-backend/internal/utils"
)

// DifyService Dify AI服务接口
type DifyService interface {
	// 调用贷款审批工作流
	CallLoanApprovalWorkflow(applicationID uint, userID uint) (*DifyWorkflowResponse, error)

	// 调用农机租赁检查工作流
	CallMachineRentalWorkflow(requestID uint, userID uint, machineID uint) (*DifyWorkflowResponse, error)

	// 记录工作流日志
	LogWorkflowExecution(log *model.DifyWorkflowLog) error
}

// difyService Dify服务实现
type difyService struct {
	cfg            *config.Config
	httpClient     *http.Client
	loanRepository repository.LoanRepository
}

// NewDifyService 创建Dify服务
func NewDifyService(cfg *config.Config, loanRepo repository.LoanRepository) DifyService {
	return &difyService{
		cfg:            cfg,
		httpClient:     &http.Client{Timeout: time.Duration(cfg.Dify.Timeout) * time.Second},
		loanRepository: loanRepo,
	}
}

// Dify工作流请求结构
type DifyWorkflowRequest struct {
	Inputs         map[string]interface{} `json:"inputs"`
	ResponseMode   string                 `json:"response_mode"`
	ConversationID string                 `json:"conversation_id,omitempty"`
	User           string                 `json:"user"`
}

// Dify工作流响应结构
type DifyWorkflowResponse struct {
	WorkflowRunID  string                 `json:"workflow_run_id"`
	TaskID         string                 `json:"task_id"`
	Data           map[string]interface{} `json:"data"`
	Error          string                 `json:"error,omitempty"`
	Status         string                 `json:"status"`
	ConversationID string                 `json:"conversation_id"`
	MessageID      string                 `json:"message_id"`
	Metadata       map[string]interface{} `json:"metadata"`
	CreatedAt      int64                  `json:"created_at"`
}

// CallLoanApprovalWorkflow 调用贷款审批工作流
func (s *difyService) CallLoanApprovalWorkflow(applicationID uint, userID uint) (*DifyWorkflowResponse, error) {
	// 构建请求数据
	request := DifyWorkflowRequest{
		Inputs: map[string]interface{}{
			"application_id": fmt.Sprintf("%d", applicationID),
			"user_id":        fmt.Sprintf("%d", userID),
		},
		ResponseMode: "blocking",
		User:         fmt.Sprintf("user_%d", userID),
	}

	// 获取工作流ID
	workflowID := s.cfg.Dify.Workflows["loan_approval"]
	if workflowID == "" {
		return nil, fmt.Errorf("贷款审批工作流ID未配置")
	}

	// 调用工作流
	response, err := s.callWorkflow(workflowID, request)
	if err != nil {
		return nil, fmt.Errorf("调用贷款审批工作流失败: %w", err)
	}

	// 记录工作流日志
	appID := uint64(applicationID)
	now := time.Now()
	duration := 0
	workflowLog := &model.DifyWorkflowLog{
		ApplicationID:  appID,
		WorkflowID:     workflowID,
		ConversationID: response.ConversationID,
		MessageID:      response.MessageID,
		WorkflowType:   "loan_approval",
		Status:         response.Status,
		StartTime:      now,
		EndTime:        &now,
		Duration:       &duration,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// 序列化请求和响应数据
	requestData, _ := json.Marshal(request)
	responseData, _ := json.Marshal(response)
	workflowLog.RequestData = string(requestData)
	workflowLog.ResponseData = string(responseData)

	// 根据响应状态设置结果
	if response.Status == "succeeded" {
		workflowLog.Result = "success"
		if result, ok := response.Data["result"]; ok {
			if resultMap, ok := result.(map[string]interface{}); ok {
				if recommendation, ok := resultMap["comments"]; ok {
					workflowLog.Recommendation = fmt.Sprintf("%v", recommendation)
				}
				if confidence, ok := resultMap["confidence_score"]; ok {
					if confFloat, ok := confidence.(float64); ok {
						workflowLog.ConfidenceScore = confFloat
					}
				}
			}
		}
	} else {
		workflowLog.Result = "failure"
		workflowLog.Recommendation = response.Error
	}

	// 保存日志
	if err := s.LogWorkflowExecution(workflowLog); err != nil {
		utils.LogError("保存Dify工作流日志失败", err)
	}

	return response, nil
}

// CallMachineRentalWorkflow 调用农机租赁检查工作流
func (s *difyService) CallMachineRentalWorkflow(requestID uint, userID uint, machineID uint) (*DifyWorkflowResponse, error) {
	// 构建请求数据
	request := DifyWorkflowRequest{
		Inputs: map[string]interface{}{
			"request_id": fmt.Sprintf("%d", requestID),
			"user_id":    fmt.Sprintf("%d", userID),
			"machine_id": fmt.Sprintf("%d", machineID),
		},
		ResponseMode: "blocking",
		User:         fmt.Sprintf("user_%d", userID),
	}

	// 获取工作流ID
	workflowID := s.cfg.Dify.Workflows["machine_rental"]
	if workflowID == "" {
		return nil, fmt.Errorf("农机租赁工作流ID未配置")
	}

	// 调用工作流
	response, err := s.callWorkflow(workflowID, request)
	if err != nil {
		return nil, fmt.Errorf("调用农机租赁工作流失败: %w", err)
	}

	// 记录工作流日志
	now := time.Now()
	duration := 0
	workflowLog := &model.DifyWorkflowLog{
		WorkflowID:     workflowID,
		ConversationID: response.ConversationID,
		MessageID:      response.MessageID,
		WorkflowType:   "machine_rental",
		Status:         response.Status,
		StartTime:      now,
		EndTime:        &now,
		Duration:       &duration,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// 序列化请求和响应数据
	requestData, _ := json.Marshal(request)
	responseData, _ := json.Marshal(response)
	workflowLog.RequestData = string(requestData)
	workflowLog.ResponseData = string(responseData)

	// 设置结果
	if response.Status == "succeeded" {
		workflowLog.Result = "success"
	} else {
		workflowLog.Result = "failure"
	}

	// 保存日志
	if err := s.LogWorkflowExecution(workflowLog); err != nil {
		utils.LogError("保存Dify工作流日志失败", err)
	}

	return response, nil
}

// callWorkflow 调用Dify工作流的通用方法
func (s *difyService) callWorkflow(workflowID string, request DifyWorkflowRequest) (*DifyWorkflowResponse, error) {
	// 序列化请求数据
	requestData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %w", err)
	}

	// 构建API URL
	apiURL := fmt.Sprintf("%s/workflows/%s/run", s.cfg.Dify.APIURL, workflowID)

	// 创建HTTP请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.Dify.APIKey)

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应数据
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应数据失败: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Dify API返回错误状态码: %d, 响应: %s", resp.StatusCode, string(responseData))
	}

	// 解析响应数据
	var response DifyWorkflowResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, fmt.Errorf("解析响应数据失败: %w", err)
	}

	// 检查业务错误
	if response.Error != "" {
		return nil, fmt.Errorf("Dify工作流执行失败: %s", response.Error)
	}

	return &response, nil
}

// LogWorkflowExecution 记录工作流执行日志
func (s *difyService) LogWorkflowExecution(log *model.DifyWorkflowLog) error {
	return s.loanRepository.CreateDifyLog(context.Background(), log)
}

// 重试调用工作流（带指数退避）
func (s *difyService) callWorkflowWithRetry(workflowID string, request DifyWorkflowRequest) (*DifyWorkflowResponse, error) {
	var lastErr error

	for i := 0; i < s.cfg.Dify.RetryTimes; i++ {
		response, err := s.callWorkflow(workflowID, request)
		if err == nil {
			return response, nil
		}

		lastErr = err
		utils.LogError(fmt.Sprintf("Dify工作流调用失败，第%d次重试", i+1), err)

		// 指数退避
		if i < s.cfg.Dify.RetryTimes-1 {
			time.Sleep(time.Duration(1<<i) * time.Second)
		}
	}

	return nil, fmt.Errorf("工作流调用失败，已重试%d次: %w", s.cfg.Dify.RetryTimes, lastErr)
}

// GetWorkflowStatus 获取工作流执行状态
func (s *difyService) GetWorkflowStatus(workflowRunID string) (*DifyWorkflowResponse, error) {
	apiURL := fmt.Sprintf("%s/workflows/run/%s", s.cfg.Dify.APIURL, workflowRunID)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.cfg.Dify.APIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应数据失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Dify API返回错误状态码: %d, 响应: %s", resp.StatusCode, string(responseData))
	}

	var response DifyWorkflowResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, fmt.Errorf("解析响应数据失败: %w", err)
	}

	return &response, nil
}

// CancelWorkflow 取消工作流执行
func (s *difyService) CancelWorkflow(workflowRunID string) error {
	apiURL := fmt.Sprintf("%s/workflows/run/%s/stop", s.cfg.Dify.APIURL, workflowRunID)

	req, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.cfg.Dify.APIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseData, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("取消工作流失败: %d, 响应: %s", resp.StatusCode, string(responseData))
	}

	return nil
}
