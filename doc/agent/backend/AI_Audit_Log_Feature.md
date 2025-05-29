# AI智能体审批日志记录功能

## 📋 概述

本文档详细说明了慧农金融系统中AI智能体操作的完整日志记录功能，确保每一次AI相关的操作都有完整的审计追踪。

## 🎯 功能特性

### 1. 全链路日志记录
- ✅ **请求日志**: 记录每次AI Agent API调用的完整请求参数
- ✅ **响应日志**: 记录AI分析结果和决策输出
- ✅ **性能日志**: 记录每次操作的处理时间
- ✅ **错误日志**: 记录操作失败的详细错误信息
- ✅ **访问日志**: 记录调用方IP地址和User-Agent

### 2. 操作类型覆盖
| 操作类型 | 说明 | Agent类型 | 记录内容 |
|---------|------|-----------|----------|
| `GET_APPLICATION_INFO` | 获取申请信息 | DIFY_WORKFLOW | 申请ID、完整申请数据 |
| `SUBMIT_AI_DECISION` | 提交AI决策 | DIFY_WORKFLOW | 决策结果、风险评分、分析报告 |
| `SUBMIT_AI_DECISION_QUERY` | 查询方式提交决策 | DIFY_WORKFLOW | 查询参数、决策结果 |
| `GET_EXTERNAL_DATA` | 获取外部数据 | DIFY_WORKFLOW | 数据类型、查询结果 |
| `GET_AI_MODEL_CONFIG` | 获取模型配置 | DIFY_WORKFLOW | 模型参数、规则配置 |

### 3. 数据结构

#### AI Agent日志表 (ai_agent_logs)
```sql
CREATE TABLE `ai_agent_logs` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `log_id` varchar(64) NOT NULL,                    -- 唯一日志ID
  `application_id` varchar(64) NULL,                -- 关联申请ID（如果有）
  `action_type` varchar(50) NOT NULL,               -- 操作类型
  `agent_type` varchar(50) NOT NULL,                -- AI智能体类型
  `request_data` json NULL,                         -- 请求数据JSON
  `response_data` json NULL,                        -- 响应数据JSON
  `status` varchar(20) NOT NULL,                    -- 操作状态（SUCCESS/ERROR）
  `error_message` text NULL,                        -- 错误信息
  `duration` bigint NULL,                           -- 处理时间（毫秒）
  `ip_address` varchar(45) NULL,                    -- 客户端IP
  `user_agent` text NULL,                           -- User-Agent
  `occurred_at` datetime(3) NOT NULL,               -- 操作时间
  `created_at` datetime(3) NOT NULL,                -- 创建时间
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_ai_agent_logs_log_id` (`log_id`),
  INDEX `idx_ai_agent_logs_action_type` (`action_type`),
  INDEX `idx_ai_agent_logs_application_id` (`application_id`)
);
```

## 🔧 技术实现

### 1. Service层日志记录
```go
// LogAIAgentAction 记录AI智能体操作日志
func (s *AIAgentService) LogAIAgentAction(
    actionType, agentType, applicationID string, 
    requestData, responseData interface{}, 
    status string, 
    errorMessage string, 
    duration int, 
    req *http.Request,
) {
    logID := pkg.GenerateUUID()
    
    var reqBytes, respBytes []byte
    if requestData != nil {
        reqBytes, _ = json.Marshal(requestData)
    }
    if responseData != nil {
        respBytes, _ = json.Marshal(responseData)
    }
    
    ipAddress := ""
    userAgent := ""
    if req != nil {
        ipAddress = getClientIP(req)
        userAgent = req.UserAgent()
    }
    
    aiLog := &data.AIAgentLog{
        LogID:         logID,
        ApplicationID: applicationID,
        ActionType:    actionType,
        AgentType:     agentType,
        RequestData:   reqBytes,
        ResponseData:  respBytes,
        Status:        status,
        ErrorMessage:  errorMessage,
        Duration:      duration,
        IPAddress:     ipAddress,
        UserAgent:     userAgent,
        OccurredAt:    time.Now(),
        CreatedAt:     time.Now(),
    }
    
    if err := s.data.DB.Create(aiLog).Error; err != nil {
        s.log.Error("记录AI Agent日志失败", zap.Error(err))
    }
}
```

### 2. 方法级日志记录

#### 获取申请信息
```go
func (s *AIAgentService) GetApplicationInfoWithLog(applicationID string, req *http.Request) (*ApplicationInfo, error) {
    startTime := time.Now()
    
    // 记录请求日志
    requestData := map[string]interface{}{
        "application_id": applicationID,
    }
    
    info, err := s.getApplicationInfoInternal(applicationID)
    
    // 计算处理时间
    duration := int(time.Since(startTime).Milliseconds())
    
    // 记录操作日志
    status := "SUCCESS"
    errorMessage := ""
    if err != nil {
        status = "ERROR"
        errorMessage = err.Error()
    }
    
    s.LogAIAgentAction(
        "GET_APPLICATION_INFO",
        "DIFY_WORKFLOW", 
        applicationID,
        requestData,
        info,
        status,
        errorMessage,
        duration,
        req,
    )
    
    return info, err
}
```

#### 提交AI决策
```go
func (s *AIAgentService) SubmitAIDecisionWithLog(applicationID string, request *AIDecisionRequest, req *http.Request) error {
    startTime := time.Now()
    
    err := s.submitAIDecisionInternal(applicationID, request)
    
    duration := int(time.Since(startTime).Milliseconds())
    
    status := "SUCCESS"
    errorMessage := ""
    if err != nil {
        status = "ERROR"
        errorMessage = err.Error()
    }
    
    responseData := &AIDecisionResult{
        ApplicationID: applicationID,
        NewStatus:     getNewStatusFromDecision(request.AIDecision.Decision),
        NextStep:      request.AIDecision.NextAction,
    }
    
    s.LogAIAgentAction(
        "SUBMIT_AI_DECISION",
        "DIFY_WORKFLOW",
        applicationID,
        request,
        responseData,
        status,
        errorMessage,
        duration,
        req,
    )
    
    return err
}
```

### 3. Handler层集成
```go
// GetApplicationInfo handler
func (h *AIAgentHandler) GetApplicationInfo(c *gin.Context) {
    applicationID := c.Param("application_id")
    if applicationID == "" {
        pkg.BadRequest(c, "申请ID不能为空")
        return
    }

    // 使用带日志记录的方法
    info, err := h.aiAgentService.GetApplicationInfoWithLog(applicationID, c.Request)
    if err != nil {
        pkg.InternalError(c, err.Error())
        return
    }

    pkg.Success(c, info)
}
```

## 📊 日志数据示例

### 1. 获取申请信息日志
```json
{
  "log_id": "log_20241201_001",
  "application_id": "test_app_001",
  "action_type": "GET_APPLICATION_INFO",
  "agent_type": "DIFY_WORKFLOW",
  "request_data": {
    "application_id": "test_app_001"
  },
  "response_data": {
    "application_id": "test_app_001",
    "product_info": {
      "product_id": "lp_001",
      "name": "春耕助力贷",
      "category": "种植贷"
    },
    "applicant_info": {
      "user_id": "user_001",
      "real_name": "张*",
      "phone": "138****5678"
    }
  },
  "status": "SUCCESS",
  "duration": 150,
  "ip_address": "10.0.0.1",
  "user_agent": "Dify-Agent/1.0",
  "occurred_at": "2024-12-01T10:30:00.000Z"
}
```

### 2. AI决策提交日志
```json
{
  "log_id": "log_20241201_002",
  "application_id": "test_app_001",
  "action_type": "SUBMIT_AI_DECISION",
  "agent_type": "DIFY_WORKFLOW",
  "request_data": {
    "ai_analysis": {
      "risk_level": "MEDIUM",
      "risk_score": 0.5,
      "confidence_score": 0.85,
      "analysis_summary": "需要人工审核"
    },
    "ai_decision": {
      "decision": "REQUIRE_HUMAN_REVIEW",
      "approved_amount": 25000,
      "approved_term_months": 12
    }
  },
  "response_data": {
    "application_id": "test_app_001",
    "new_status": "MANUAL_REVIEW_REQUIRED",
    "next_step": "ASSIGN_TO_REVIEWER"
  },
  "status": "SUCCESS",
  "duration": 300,
  "ip_address": "10.0.0.1",
  "occurred_at": "2024-12-01T10:32:00.000Z"
}
```

## 🔍 日志查询与分析

### 1. 按申请ID查询
```sql
SELECT * FROM ai_agent_logs 
WHERE application_id = 'test_app_001' 
ORDER BY occurred_at ASC;
```

### 2. 按操作类型统计
```sql
SELECT 
    action_type,
    COUNT(*) as total_count,
    AVG(duration) as avg_duration,
    SUM(CASE WHEN status = 'SUCCESS' THEN 1 ELSE 0 END) as success_count,
    SUM(CASE WHEN status = 'ERROR' THEN 1 ELSE 0 END) as error_count
FROM ai_agent_logs 
WHERE occurred_at >= DATE_SUB(NOW(), INTERVAL 1 DAY)
GROUP BY action_type;
```

### 3. 性能分析
```sql
SELECT 
    action_type,
    MIN(duration) as min_duration,
    MAX(duration) as max_duration,
    AVG(duration) as avg_duration,
    PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY duration) as p95_duration
FROM ai_agent_logs 
WHERE status = 'SUCCESS' 
    AND occurred_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)
GROUP BY action_type;
```

## 🛠️ 测试验证

### 运行日志记录测试
```bash
# 执行AI日志记录功能测试
./doc/agent/backend/Test-API-AI-Logs.sh
```

### 预期结果
1. ✅ 所有AI Agent接口调用成功
2. ✅ 数据库中生成对应的日志记录
3. ✅ 日志包含完整的请求/响应数据
4. ✅ 记录准确的处理时间和状态

## 📈 监控告警

### 1. 性能监控
- 监控各操作类型的平均响应时间
- 设置响应时间阈值告警（如>5秒）
- 监控操作成功率

### 2. 错误监控
- 监控错误率变化
- 特定错误类型的告警
- 连续失败的告警

### 3. 业务监控
- AI决策分布统计
- 风险评分趋势分析
- 申请处理量统计

## 🔐 安全与合规

### 1. 数据脱敏
- 用户敏感信息自动脱敏
- 日志中不包含完整身份证号、银行卡号等
- 仅记录业务必要的数据

### 2. 访问控制
- 日志数据仅限授权人员访问
- 操作审计人员可查看完整日志
- 系统管理员可进行日志管理

### 3. 数据保留
- 日志数据保留期：3年
- 超期数据自动归档
- 重要日志永久保存

## 📋 总结

通过完整的AI智能体操作日志记录功能，系统现在能够：

1. **完整追踪** - 记录每一次AI操作的完整过程
2. **性能监控** - 实时监控AI服务的性能表现  
3. **错误诊断** - 快速定位和诊断问题
4. **合规审计** - 满足监管要求的操作审计
5. **业务分析** - 提供AI决策效果的数据支持

这确保了AI审批流程的完整可追溯性和系统的可靠性。 