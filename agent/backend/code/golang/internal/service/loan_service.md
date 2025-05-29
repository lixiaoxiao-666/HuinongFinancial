# 贷款服务代码文档

## 概述

`loan_service.go` 实现了贷款申请的核心业务逻辑，包括产品管理、申请处理、审批流程和AI智能评估等功能。

## 主要功能

### 1. 依赖注入

```go
type loanService struct {
    loanRepo    repository.LoanRepository
    userRepo    repository.UserRepository
    difyService DifyService
}
```

贷款服务集成了以下依赖：
- `loanRepo`: 贷款数据访问层
- `userRepo`: 用户数据访问层  
- `difyService`: Dify AI工作流服务

### 2. 申请创建流程

#### CreateApplication 方法

```go
func (s *loanService) CreateApplication(ctx context.Context, req *CreateApplicationRequest) (*CreateApplicationResponse, error)
```

**核心流程**：
1. 验证产品是否存在
2. 验证申请金额是否在产品限制范围内
3. 生成唯一申请编号
4. 创建申请记录，初始状态为 `pending`
5. **异步触发AI评估** - 关键新增功能

**AI工作流触发**：
```go
// 异步触发AI审批流程
go func() {
    assessmentCtx := context.Background()
    if triggerErr := s.TriggerAIAssessment(assessmentCtx, application.ID, "loan_approval"); triggerErr != nil {
        fmt.Printf("触发AI评估失败 - 申请ID: %d, 错误: %v\n", application.ID, triggerErr)
    }
}()
```

**设计要点**：
- 使用 goroutine 异步执行，避免阻塞用户申请提交
- 使用新的 context 避免原请求 context 被取消影响AI评估
- 错误处理不影响申请创建的成功响应

### 3. AI智能评估

#### TriggerAIAssessment 方法

```go
func (s *loanService) TriggerAIAssessment(ctx context.Context, applicationID uint64, workflowType string) error
```

**完整评估流程**：

1. **获取申请数据**
   ```go
   application, err := s.loanRepo.GetApplicationByID(ctx, uint(applicationID))
   ```

2. **更新状态为处理中**
   ```go
   application.Status = "ai_processing"
   ```

3. **记录评估开始日志**
   ```go
   startLog := &model.ApprovalLog{
       ApplicationID: application.ID,
       ApproverID:    0, // AI系统
       Action:        "submit",
       Step:          "ai_assessment",
       // ...
   }
   ```

4. **调用Dify工作流**
   ```go
   response, err := s.difyService.CallLoanApprovalWorkflow(uint(applicationID), uint(application.UserID))
   ```

5. **处理评估结果**
   ```go
   if err := s.processAIAssessmentResult(ctx, application, response); err != nil {
       return fmt.Errorf("处理AI评估结果失败: %v", err)
   }
   ```

#### processAIAssessmentResult 方法

**结果处理逻辑**：

1. **解析AI响应**
   ```go
   result, ok := response.Data["result"]
   resultMap, ok := result.(map[string]interface{})
   ```

2. **提取关键信息**
   ```go
   // 获取AI建议
   recommendation := resultMap["recommendation"]
   // 获取AI决策
   decision := resultMap["decision"]
   // 获取信用评分
   application.CreditScore = int(resultMap["credit_score"])
   // 获取风险等级  
   application.RiskLevel = resultMap["risk_level"]
   ```

3. **根据决策更新状态**
   ```go
   switch decision {
   case "approve", "approved":
       newStatus = "ai_approved"
       application.AutoApprovalPassed = true
   case "reject", "rejected":
       newStatus = "ai_rejected"
       application.RejectionReason = resultMap["rejection_reason"]
   case "manual_review", "manual":
       newStatus = "manual_review"
       application.ApprovalLevel = 1
   }
   ```

4. **记录完成日志**
   ```go
   resultLog := &model.ApprovalLog{
       ApplicationID: application.ID,
       Action:        actionResult,
       Comment:       actionComment,
       Note:          fmt.Sprintf("AI评估建议: %s", recommendation),
       // ...
   }
   ```

### 4. 状态管理

#### 申请状态流转

```
pending → ai_processing → ai_approved/ai_rejected/manual_review
                      ↓
                   ai_failed (错误情况)
```

**状态说明**：
- `pending`: 初始提交状态
- `ai_processing`: AI评估进行中
- `ai_approved`: AI建议批准
- `ai_rejected`: AI建议拒绝
- `ai_failed`: AI评估失败
- `manual_review`: 需要人工审核

### 5. 错误处理

#### 多层错误处理机制

1. **工作流调用失败**
   ```go
   if err != nil {
       application.Status = "ai_failed"
       application.RejectionReason = fmt.Sprintf("AI评估失败: %v", err)
       // 记录失败日志
       failLog := &model.ApprovalLog{...}
   }
   ```

2. **结果处理失败**
   ```go
   if response.Status != "succeeded" {
       application.Status = "ai_failed"
       application.RejectionReason = "AI评估过程中出现错误"
   }
   ```

3. **日志记录失败**
   ```go
   if logErr := s.loanRepo.CreateApprovalLog(ctx, startLog); logErr != nil {
       fmt.Printf("创建审批日志失败: %v\n", logErr)
   }
   ```

**容错设计**：
- 日志记录失败不影响主流程
- AI评估失败有明确的状态标识
- 错误信息详细记录便于调试

### 6. 性能优化

#### 异步处理

- **申请提交响应快**: 用户提交后立即返回，不等待AI评估完成
- **资源利用高效**: 使用goroutine异步处理，不阻塞主线程
- **并发安全**: 每个评估使用独立的context和事务

#### 数据传递优化

- **最小化数据库查询**: 在AI评估中一次性获取所需数据
- **结构化数据传递**: 使用标准化的工作流输入格式
- **缓存友好**: 状态更新操作支持缓存机制

### 7. 监控和日志

#### 审批日志记录

每个关键步骤都有详细的审批日志：
```go
type ApprovalLog struct {
    ApplicationID  uint64    // 申请ID
    ApproverID     uint64    // 审批人ID (0表示AI系统)
    Action         string    // 审批动作
    Step           string    // 审批步骤
    Result         string    // 审批结果
    Comment        string    // 审批意见
    PreviousStatus string    // 前置状态
    NewStatus      string    // 新状态
    ActionTime     time.Time // 操作时间
}
```

#### Dify工作流日志

详细记录AI工作流的执行情况：
```go
type DifyWorkflowLog struct {
    ApplicationID   uint64  // 申请ID
    WorkflowID      string  // 工作流ID
    ConversationID  string  // 对话ID
    WorkflowType    string  // 工作流类型
    Status          string  // 执行状态
    Result          string  // 执行结果
    Recommendation  string  // AI建议
    RequestData     string  // 请求数据
    ResponseData    string  // 响应数据
}
```

### 8. 扩展性设计

#### 工作流类型扩展

```go
// 支持多种工作流类型
const (
    WorkflowLoanApproval   = "loan_approval"
    WorkflowRiskAssessment = "risk_assessment" 
    WorkflowDocumentReview = "document_review"
)
```

#### 决策规则扩展

```go
// 可扩展的决策处理
switch decision {
case "approve", "approved":
    // 批准逻辑
case "reject", "rejected": 
    // 拒绝逻辑
case "manual_review", "manual":
    // 人工审核逻辑
default:
    // 默认处理逻辑
}
```

## 使用示例

### 客户端调用

```go
// 创建申请
req := &CreateApplicationRequest{
    ProductID:     1,
    LoanAmount:    5000000, // 5万元，以分为单位
    TermMonths:    12,
    LoanPurpose:   "农作物种植资金",
    ContactPhone:  "13800138000",
    MaterialsJSON: `{"id_card":"path/to/id.jpg"}`,
}

response, err := loanService.CreateApplication(ctx, req)
if err != nil {
    log.Fatal(err)
}

// 申请创建成功，AI评估将在后台自动进行
fmt.Printf("申请编号: %s\n", response.ApplicationNo)
```

### 状态查询

```go
// 查询申请详情和AI评估结果
details, err := loanService.GetApplicationDetails(ctx, &GetApplicationDetailsRequest{
    ID: response.ID,
})

// 检查AI评估状态
switch details.Application.Status {
case "ai_approved":
    fmt.Printf("AI建议批准: %s\n", details.Application.AIRecommendation)
case "ai_rejected":
    fmt.Printf("AI建议拒绝: %s\n", details.Application.RejectionReason)
case "manual_review":
    fmt.Println("需要人工审核")
case "ai_processing":
    fmt.Println("AI评估进行中")
}
```

## 配置要求

### Dify配置

```yaml
dify:
  api_url: "http://172.18.120.57/v1"
  api_key: "app-d2ELQ0kfvzVv1m84LGQR4pKs"
  workflows:
    loan_approval: "loan-approval-workflow-id"
```

### 数据库表

- `loan_applications`: 贷款申请主表
- `approval_logs`: 审批日志表
- `dify_workflow_logs`: Dify工作流日志表

## 注意事项

1. **异步处理**: AI评估是异步的，客户端需要轮询或使用WebSocket获取结果
2. **错误恢复**: AI评估失败时有重试机制和人工介入流程
3. **数据一致性**: 使用事务确保状态更新的一致性
4. **性能监控**: 关注AI工作流的响应时间和成功率
5. **日志审计**: 完整记录所有操作便于问题排查和合规要求