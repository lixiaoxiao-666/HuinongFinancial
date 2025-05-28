# AI智能体工作流程文档

## 概述

本文档详细描述了慧农金融系统中AI智能体的完整工作流程，包括与Dify AI工作流的集成方案、接口调用关系和数据流转过程。

## 系统架构图

```
用户APP ──────→ 后端API ──────→ AI智能体服务
   │                │              │
   │                │              │
   ▼                ▼              ▼
提交申请      触发AI工作流     调用Dify API
   │                │              │
   │                │              │
   ▼                ▼              ▼
数据入库      工作流管理       AI分析处理
   │                │              │
   │                │              │
   ▼                ▼              ▼
状态更新      状态追踪       决策结果返回
```

## 详细工作流程

### 阶段一：申请提交与触发

#### 1. 用户提交申请
**触发条件**：用户在APP端完成贷款申请提交

**执行步骤**：
1. 用户填写申请信息（金额、期限、用途等）
2. 上传必要文档（身份证、土地证明等）
3. 提交申请到后端API

**接口调用**：
```http
POST /api/v1/loans/applications
Authorization: Bearer {user_jwt_token}
Content-Type: application/json

{
  "product_id": "SEED_LOAN_001",
  "amount": 30000,
  "term_months": 12,
  "purpose": "购买种子和农药",
  "applicant_info": {
    "annual_income": 80000,
    "land_area": "10亩"
  },
  "uploaded_documents": [
    {
      "doc_type": "id_card_front",
      "file_id": "file_001"
    }
  ]
}
```

**处理逻辑**：
- 验证申请数据完整性
- 创建申请记录（状态：SUBMITTED）
- 关联上传文档
- 记录申请历史

#### 2. 系统触发AI工作流
**触发条件**：申请成功创建后自动触发

**执行步骤**：
1. 后端服务异步调用AI工作流触发接口
2. 创建工作流执行记录
3. 更新申请状态为"AI_TRIGGERED"

**接口调用**：
```http
POST /api/v1/ai-agent/applications/{application_id}/trigger-workflow
Authorization: System-Token {system_token}
Content-Type: application/json

{
  "workflow_type": "LOAN_APPROVAL",
  "priority": "NORMAL",
  "callback_url": "https://api.example.com/ai-agent/callback"
}
```

**代码实现**：
```go
// 在loan_service.go中
func (s *LoanService) SubmitLoanApplication(ctx context.Context, userID string, req *SubmitLoanApplicationRequest) (*LoanApplicationResponse, error) {
    // ... 创建申请逻辑 ...
    
    // 异步触发AI审批流程
    go s.triggerAIProcessing(application.ApplicationID)
    
    return response, nil
}

func (s *LoanService) triggerAIProcessing(applicationID string) {
    // 调用AI智能体服务触发工作流
    // 实际实现中会调用HTTP接口
}
```

### 阶段二：AI数据获取与分析

#### 3. Dify工作流获取申请数据
**触发条件**：AI工作流收到触发请求

**执行步骤**：
1. Dify工作流开始执行
2. 调用后端接口获取申请详细信息
3. 获取外部数据（征信、银行流水等）

**接口调用序列**：

```http
# 获取申请详细信息
GET /api/v1/ai-agent/applications/{application_id}/info
Authorization: AI-Agent-Token {dify_token}

# 获取外部数据
GET /api/v1/ai-agent/external-data/{user_id}?data_types=credit,bank_flow,blacklist
Authorization: AI-Agent-Token {dify_token}

# 获取AI模型配置
GET /api/v1/ai-agent/config/models
Authorization: AI-Agent-Token {dify_token}
```

**数据结构**：
```json
{
  "application_id": "app_20240301_001",
  "product_info": {
    "product_id": "SEED_LOAN_001",
    "name": "种植贷",
    "max_amount": 50000,
    "interest_rate_yearly": "6.5%-8.5%"
  },
  "applicant_info": {
    "user_id": "user_001",
    "real_name": "张三",
    "credit_score": 750,
    "annual_income": 80000
  },
  "external_data": {
    "credit_bureau_score": 750,
    "blacklist_check": false,
    "land_ownership_verified": true
  }
}
```

#### 4. AI模型分析与决策
**处理过程**：
1. **数据预处理**：清洗和标准化申请数据
2. **特征工程**：提取关键风险特征
3. **风险评估**：计算风险分数和等级
4. **决策生成**：根据规则生成审批决策

**AI分析维度**：
- **信用评估**：征信报告、历史还款记录
- **收入分析**：收入稳定性、债务收入比
- **资产评估**：土地价值、固定资产
- **行为分析**：申请行为、历史申请记录
- **欺诈检测**：异常行为识别

**决策逻辑**：
```python
def ai_decision_logic(risk_score, amount, user_profile):
    if risk_score < 0.3:
        return "AUTO_APPROVED"  # 低风险自动通过
    elif risk_score < 0.7:
        return "REQUIRE_HUMAN_REVIEW"  # 中风险人工审核
    else:
        return "AUTO_REJECTED"  # 高风险自动拒绝
```

### 阶段三：决策提交与状态更新

#### 5. Dify提交AI决策结果
**触发条件**：AI分析完成

**执行步骤**：
1. 整理AI分析结果
2. 生成决策建议
3. 调用后端接口提交决策

**接口调用**：
```http
POST /api/v1/ai-agent/applications/{application_id}/ai-decision
Authorization: AI-Agent-Token {dify_token}
Content-Type: application/json

{
  "ai_analysis": {
    "risk_level": "LOW",
    "risk_score": 0.25,
    "confidence_score": 0.92,
    "analysis_summary": "申请人信用状况良好，还款能力强",
    "detailed_analysis": {
      "income_analysis": "年收入8万元，稳定",
      "credit_analysis": "征信良好，无不良记录",
      "asset_analysis": "拥有10亩土地，资产充足"
    },
    "risk_factors": [
      {
        "factor": "credit_score",
        "value": 750,
        "weight": 0.3,
        "risk_contribution": 0.05
      }
    ],
    "recommendations": [
      "建议批准贷款",
      "可给予优惠利率"
    ]
  },
  "ai_decision": {
    "decision": "AUTO_APPROVED",
    "approved_amount": 30000,
    "approved_term_months": 12,
    "suggested_interest_rate": "6.5%",
    "conditions": [],
    "next_action": "AWAIT_FINAL_CONFIRMATION"
  },
  "processing_info": {
    "ai_model_version": "v2.1.0",
    "processing_time_ms": 1500,
    "workflow_id": "workflow_001",
    "processed_at": "2024-03-01T10:35:00Z"
  }
}
```

#### 6. 后端处理AI决策
**处理逻辑**：
1. 验证决策数据完整性
2. 保存AI分析结果
3. 更新申请状态
4. 记录状态变更历史
5. 触发后续流程

**状态流转规则**：
```go
switch request.AIDecision.Decision {
case "AUTO_APPROVED":
    newStatus = "AI_APPROVED"
    // 设置批准金额和期限
    // 记录为最终决策
case "REQUIRE_HUMAN_REVIEW":
    newStatus = "MANUAL_REVIEW_REQUIRED"
    // 进入人工审核队列
case "AUTO_REJECTED":
    newStatus = "AI_REJECTED"
    // 记录拒绝原因
    // 通知用户
}
```

### 阶段四：后续处理流程

#### 7. 不同决策的后续处理

##### 7.1 自动通过 (AUTO_APPROVED)
**处理流程**：
1. 更新申请状态为"AI_APPROVED"
2. 记录批准金额和期限
3. 生成放款准备
4. 发送通知给用户
5. 进入合同签署流程

**通知内容**：
```json
{
  "type": "LOAN_APPROVED",
  "application_id": "app_001",
  "message": "恭喜！您的贷款申请已通过AI审批",
  "approved_amount": 30000,
  "interest_rate": "6.5%",
  "next_steps": ["签署电子合同", "等待放款"]
}
```

##### 7.2 需要人工审核 (REQUIRE_HUMAN_REVIEW)
**处理流程**：
1. 更新申请状态为"MANUAL_REVIEW_REQUIRED"
2. 加入人工审核队列
3. 分配给审批员
4. 发送待审核通知
5. 提供AI分析建议

**审批员界面数据**：
```json
{
  "application_id": "app_001",
  "applicant_name": "张三",
  "applied_amount": 30000,
  "ai_analysis": {
    "risk_score": 0.65,
    "main_concerns": ["申请金额较大", "需要更多收入证明"],
    "recommendations": ["建议人工审核", "可考虑降低金额"]
  },
  "supporting_documents": [...],
  "external_data": {...}
}
```

##### 7.3 自动拒绝 (AUTO_REJECTED)
**处理流程**：
1. 更新申请状态为"AI_REJECTED"
2. 记录拒绝原因
3. 发送拒绝通知
4. 提供改进建议
5. 结束申请流程

**拒绝通知**：
```json
{
  "type": "LOAN_REJECTED",
  "application_id": "app_001",
  "message": "很抱歉，您的贷款申请未通过审核",
  "reasons": [
    "信用评分较低",
    "收入证明不足",
    "债务收入比过高"
  ],
  "suggestions": [
    "提高信用评分后重新申请",
    "提供更详细的收入证明",
    "降低申请金额"
  ]
}
```

## 错误处理机制

### 1. 网络异常处理
```go
func (s *AIAgentService) handleNetworkError(err error) {
    // 记录错误日志
    s.log.Error("网络请求失败", zap.Error(err))
    
    // 根据错误类型进行重试或降级处理
    if isRetryableError(err) {
        // 执行重试逻辑
        s.retryRequest()
    } else {
        // 降级到人工审核
        s.fallbackToManualReview()
    }
}
```

### 2. AI服务异常处理
- **超时处理**：设置合理的超时时间，超时后转人工审核
- **模型异常**：模型返回异常结果时降级处理
- **数据异常**：关键数据缺失时要求补充材料

### 3. 数据一致性保证
- 使用数据库事务确保数据一致性
- 记录完整的操作日志
- 实现幂等性处理

## 监控与日志

### 1. 性能监控
- **处理时间**：记录每个阶段的耗时
- **成功率**：统计AI决策的准确率
- **吞吐量**：监控系统处理能力

### 2. 业务监控
- **审批通过率**：按产品、地区统计
- **风险分布**：分析风险等级分布
- **人工干预率**：需要人工审核的比例

### 3. 日志记录
```go
type AIAgentLog struct {
    LogID         string    `json:"log_id"`
    ApplicationID string    `json:"application_id"`
    ActionType    string    `json:"action_type"`    // GET_INFO, SUBMIT_DECISION
    AgentType     string    `json:"agent_type"`     // DIFY, SYSTEM
    RequestData   string    `json:"request_data"`   // 请求数据
    ResponseData  string    `json:"response_data"`  // 响应数据
    Status        string    `json:"status"`         // SUCCESS, ERROR
    Duration      int       `json:"duration"`       // 处理时间(ms)
    OccurredAt    time.Time `json:"occurred_at"`
}
```

## 配置管理

### 1. AI模型配置
```json
{
  "active_models": [
    {
      "model_id": "risk_assessment_v2",
      "version": "2.1.0",
      "thresholds": {
        "low_risk": 0.3,
        "medium_risk": 0.7,
        "high_risk": 0.9
      }
    }
  ],
  "approval_rules": {
    "auto_approval_threshold": 0.3,
    "auto_rejection_threshold": 0.8,
    "max_auto_approval_amount": 50000
  }
}
```

### 2. 业务参数配置
```json
{
  "business_parameters": {
    "max_debt_to_income_ratio": 0.5,
    "min_credit_score": 600,
    "max_loan_amount_by_category": {
      "种植贷": 50000,
      "设备贷": 200000,
      "其他": 30000
    }
  }
}
```

## 安全考虑

### 1. 认证与授权
- **AI Agent Token**：用于Dify等外部服务认证
- **System Token**：用于系统内部服务调用
- **Token轮换**：定期更换Token确保安全

### 2. 数据安全
- **数据脱敏**：敏感信息在日志中脱敏处理
- **传输加密**：使用HTTPS确保数据传输安全
- **访问控制**：严格的接口访问权限控制

### 3. 审计跟踪
- 记录所有AI决策过程
- 保留完整的数据变更历史
- 支持决策过程的追溯和审计

## 部署与运维

### 1. 服务部署
```yaml
# docker-compose.yml
services:
  ai-agent-service:
    image: huinong-backend:latest
    environment:
      - AI_AGENT_ENABLED=true
      - DIFY_API_URL=https://api.dify.ai
      - AI_AGENT_TOKEN=secure_token
    depends_on:
      - mysql
      - redis
```

### 2. 健康检查
```go
func (h *AIAgentHandler) HealthCheck(c *gin.Context) {
    health := map[string]interface{}{
        "status": "healthy",
        "ai_service": "connected",
        "database": "connected",
        "last_check": time.Now(),
    }
    c.JSON(200, health)
}
```

### 3. 扩展性考虑
- 支持多AI服务提供商
- 可配置的决策规则
- 水平扩展能力

## 总结

AI智能体工作流程实现了从用户申请到AI决策的完整自动化流程，具备以下特点：

1. **完整性**：覆盖申请提交到决策输出的全流程
2. **可靠性**：具备完善的错误处理和降级机制
3. **可扩展性**：支持多种AI服务和决策规则
4. **可监控性**：提供全面的日志和监控能力
5. **安全性**：严格的认证授权和数据保护机制

该方案为慧农金融系统提供了高效、智能的自动化审批能力，显著提升了业务处理效率和用户体验。 