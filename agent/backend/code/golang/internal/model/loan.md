# 贷款业务数据模型

## 文件概述

`loan.go` 是数字惠农系统贷款业务的核心数据模型文件，定义了贷款产品、贷款申请、审批流程以及Dify AI工作流集成相关的所有数据结构。

## 核心数据模型

### 1. LoanProduct 贷款产品表
贷款产品管理的核心模型，定义了各种农业贷款产品的基本信息和规则参数。

```go
type LoanProduct struct {
    ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    ProductCode string `gorm:"type:varchar(30);uniqueIndex;not null" json:"product_code"`
    ProductName string `gorm:"type:varchar(100);not null" json:"product_name"`
    
    // 产品类型：agricultural_material(农资贷)、equipment(设备贷)、operation(经营贷)、seasonal(季节性贷款)
    ProductType string `gorm:"type:varchar(30);not null" json:"product_type"`
    
    // 贷款参数
    MinAmount      int64   `json:"min_amount"`       // 最小贷款金额(分)
    MaxAmount      int64   `json:"max_amount"`       // 最大贷款金额(分) 
    InterestRate   float64 `json:"interest_rate"`    // 年利率
    LoanTermMonths int     `json:"loan_term_months"` // 贷款期限(月)
    
    // 风控参数
    AutoApprovalThreshold int64 `json:"auto_approval_threshold"` // 自动审批阈值(分)
    RiskLevel            string `json:"risk_level"`              // low, medium, high
    CreditScoreMin       int    `json:"credit_score_min"`        // 最低信用分数
    CollateralRequired   bool   `json:"collateral_required"`     // 是否需要抵押
}
```

**产品类型 (ProductType)**:
- `agricultural_material`: 农资贷 - 购买种子、化肥、农药等农业生产资料
- `equipment`: 设备贷 - 购买农机设备、生产设备
- `operation`: 经营贷 - 农业生产经营周转资金
- `seasonal`: 季节性贷款 - 按农业生产季节特点设计的贷款

**风险等级 (RiskLevel)**:
- `low`: 低风险 - 风控要求相对宽松
- `medium`: 中等风险 - 标准风控要求
- `high`: 高风险 - 严格风控要求

### 2. LoanApplication 贷款申请表
贷款申请的详细信息管理，包含申请人信息、农业经营情况、审批流程等完整数据。

```go
type LoanApplication struct {
    ID            uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    ApplicationNo string `gorm:"type:varchar(30);uniqueIndex;not null" json:"application_no"`
    
    // 关联信息
    UserID    uint64 `gorm:"not null;index" json:"user_id"`
    ProductID uint64 `gorm:"not null;index" json:"product_id"`
    
    // 申请基本信息
    ApplyAmount     int64      `json:"apply_amount"`      // 申请金额(分)
    ApplyTermMonths int        `json:"apply_term_months"` // 申请期限(月)
    LoanPurpose     string     `json:"loan_purpose"`      // 贷款用途
    ExpectedUseDate *time.Time `json:"expected_use_date"` // 预期用款日期
    
    // 经营信息(农业相关)
    FarmArea          float64 `json:"farm_area"`           // 经营面积(亩)
    CropTypes         string  `json:"crop_types"`          // 种植作物(JSON数组)
    YearsOfExperience int     `json:"years_of_experience"` // 从业年限
    
    // 审批流程
    Status              string     `json:"status"`                // 当前状态
    CurrentApprover     *uint64    `json:"current_approver"`      // 当前审批人
    ApprovalLevel       int        `json:"approval_level"`        // 当前审批层级
    AutoApprovalPassed  bool       `json:"auto_approval_passed"`  // 是否通过自动审批
    
    // Dify AI工作流
    DifyWorkflowID     string `json:"dify_workflow_id"`     // Dify工作流ID
    DifyConversationID string `json:"dify_conversation_id"` // 对话ID
    AIRecommendation   string `json:"ai_recommendation"`    // AI审批建议
}
```

**申请状态流转**:
```
pending(待审核) → 
reviewing(审核中) → 
auto_approved(自动审批通过) / manual_review(人工审核) →
approved(审批通过) / rejected(审批拒绝) →
disbursed(已放款) / cancelled(已取消)
```

### 3. ApprovalLog 审批日志表
记录贷款申请的完整审批流程，支持多级审批和审批参数调整。

```go
type ApprovalLog struct {
    ID            uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    ApplicationID uint64 `gorm:"not null;index" json:"application_id"`
    ApproverID    uint64 `gorm:"not null;index" json:"approver_id"`
    
    // 审批动作：submit(提交)、approve(通过)、reject(拒绝)、return(退回)、withdraw(撤回)
    Action string `gorm:"type:varchar(20);not null" json:"action"`
    
    // 审批层级
    ApprovalLevel int `gorm:"not null" json:"approval_level"`
    
    // 审批结果：pending(待处理)、approved(通过)、rejected(拒绝)、returned(退回)
    Result string `gorm:"type:varchar(20);not null" json:"result"`
    
    // 审批意见
    Comment string `gorm:"type:text" json:"comment"`
    
    // 审批参数调整
    AmountAdjustment *int64   `json:"amount_adjustment"` // 金额调整(分)
    TermAdjustment   *int     `json:"term_adjustment"`   // 期限调整(月)
    RateAdjustment   *float64 `json:"rate_adjustment"`   // 利率调整
}
```

**审批动作类型**:
- `submit`: 提交申请
- `approve`: 审批通过
- `reject`: 审批拒绝
- `return`: 退回修改
- `withdraw`: 申请人撤回

### 4. DifyWorkflowLog Dify工作流调用日志表
记录与Dify AI平台的所有交互日志，用于AI辅助审批和决策支持。

```go
type DifyWorkflowLog struct {
    ID            uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    ApplicationID uint64 `gorm:"not null;index" json:"application_id"`
    
    // Dify相关信息
    WorkflowID     string `json:"workflow_id"`     // Dify工作流ID
    ConversationID string `json:"conversation_id"` // 对话ID
    MessageID      string `json:"message_id"`      // 消息ID
    
    // 工作流类型：risk_assessment(风险评估)、document_review(文档审核)、decision_support(决策支持)
    WorkflowType string `json:"workflow_type"`
    
    // 调用信息
    RequestData  string `json:"request_data"`  // 请求数据(JSON)
    ResponseData string `json:"response_data"` // 响应数据(JSON)
    
    // 执行结果
    Result          string  `json:"result"`           // 工作流执行结果
    Recommendation  string  `json:"recommendation"`   // AI建议
    ConfidenceScore float64 `json:"confidence_score"` // 置信度评分
    
    // 性能指标
    Duration   *int  `json:"duration"`    // 执行时长(毫秒)
    TokenUsage int   `json:"token_usage"` // Token使用量
    CostAmount int64 `json:"cost_amount"` // 成本(分)
}
```

**工作流类型**:
- `risk_assessment`: 风险评估 - 基于申请信息进行风险分析
- `document_review`: 文档审核 - AI审核申请材料的完整性和真实性
- `decision_support`: 决策支持 - 为审批人员提供决策建议

## 辅助数据结构

### LoanApplyConditions 贷款申请条件
```go
type LoanApplyConditions struct {
    MinAge               int      `json:"min_age"`                // 最小年龄
    MaxAge               int      `json:"max_age"`                // 最大年龄
    MinYearsExperience   int      `json:"min_years_experience"`   // 最小从业年限
    MinFarmArea          float64  `json:"min_farm_area"`          // 最小经营面积
    RequiredRegions      []string `json:"required_regions"`       // 适用地区
    RequiredCertificates []string `json:"required_certificates"`  // 必需证件
    BlacklistIndustries  []string `json:"blacklist_industries"`   // 禁入行业
}
```

### LoanFeatures 贷款产品特色
```go
type LoanFeatures struct {
    FastApproval      bool   `json:"fast_approval"`      // 快速审批
    OnlineApplication bool   `json:"online_application"` // 在线申请
    FlexibleRepayment bool   `json:"flexible_repayment"` // 灵活还款
    NoCollateral      bool   `json:"no_collateral"`      // 免抵押
    LowInterestRate   bool   `json:"low_interest_rate"`  // 低利率
    SeasonalSupport   bool   `json:"seasonal_support"`   // 季节性支持
    Description       string `json:"description"`        // 特色描述
}
```

### ApplicationDocuments 申请材料
```go
type ApplicationDocuments struct {
    IDCardFront     string   `json:"id_card_front"`     // 身份证正面
    IDCardBack      string   `json:"id_card_back"`      // 身份证背面
    BankStatement   []string `json:"bank_statement"`    // 银行流水
    IncomeProof     []string `json:"income_proof"`      // 收入证明
    LandCertificate []string `json:"land_certificate"`  // 土地证明
    BusinessLicense string   `json:"business_license"`  // 营业执照
    CropRecords     []string `json:"crop_records"`      // 种植记录
    Other           []string `json:"other"`             // 其他材料
}
```

## Dify集成数据结构

### DifyRequestData Dify请求数据
```go
type DifyRequestData struct {
    Query     string                 `json:"query"`  // 查询内容
    Inputs    map[string]interface{} `json:"inputs"` // 输入变量
    UserID    string                 `json:"user"`   // 用户标识
    Files     []DifyFile            `json:"files"`  // 文件列表
}
```

### DifyResponseData Dify响应数据
```go
type DifyResponseData struct {
    ConversationID string                 `json:"conversation_id"` // 对话ID
    MessageID      string                 `json:"message_id"`      // 消息ID
    Answer         string                 `json:"answer"`          // 回答内容
    Metadata       map[string]interface{} `json:"metadata"`        // 元数据
    Usage          DifyUsage             `json:"usage"`           // 使用统计
}
```

## 业务流程设计

### 贷款申请流程
1. **申请提交**: 用户在APP中选择贷款产品，填写申请信息
2. **自动预审**: 系统根据产品规则进行自动预审
3. **AI风险评估**: 调用Dify工作流进行风险评估
4. **人工审核**: 根据评估结果决定是否需要人工审核
5. **审批决策**: 最终审批决定和参数调整
6. **放款处理**: 审批通过后的放款操作

### AI辅助审批流程
```go
// AI风险评估工作流
func CallDifyRiskAssessment(application *LoanApplication) error {
    // 构建请求数据
    requestData := DifyRequestData{
        Query: "请对以下贷款申请进行风险评估",
        Inputs: map[string]interface{}{
            "申请金额":    application.ApplyAmount,
            "申请期限":    application.ApplyTermMonths,
            "经营面积":    application.FarmArea,
            "从业年限":    application.YearsOfExperience,
            "月收入":     application.MonthlyIncome,
            "其他负债":    application.OtherDebts,
        },
        UserID: fmt.Sprintf("user_%d", application.UserID),
    }
    
    // 调用Dify API
    response, err := callDifyWorkflow("risk_assessment", requestData)
    if err != nil {
        return err
    }
    
    // 记录调用日志
    log := &DifyWorkflowLog{
        ApplicationID:   application.ID,
        WorkflowType:    "risk_assessment",
        RequestData:     string(requestDataJSON),
        ResponseData:    string(responseJSON),
        Result:          response.Answer,
        ConfidenceScore: extractConfidenceScore(response),
        Status:          "completed",
    }
    
    return db.Create(log).Error
}
```

## 数据库索引设计

### 关键索引
1. **贷款产品表 (loan_products)**:
   - `product_code` 唯一索引
   - `product_type` 普通索引
   - `status` 普通索引

2. **贷款申请表 (loan_applications)**:
   - `application_no` 唯一索引
   - `user_id` 普通索引
   - `product_id` 普通索引
   - `status` 普通索引
   - 联合索引: (`user_id`, `status`)

3. **审批日志表 (approval_logs)**:
   - `application_id` 普通索引
   - `approver_id` 普通索引
   - 联合索引: (`application_id`, `approval_level`)

4. **Dify工作流日志表 (dify_workflow_logs)**:
   - `application_id` 普通索引
   - `workflow_type` 普通索引

## 使用示例

### 创建贷款产品
```go
// 创建农资贷产品
product := &model.LoanProduct{
    ProductCode: "NZRD001",
    ProductName: "春耕农资贷",
    ProductType: "agricultural_material",
    Description: "专为春耕季节农资采购设计的贷款产品",
    MinAmount:   1000000,  // 10,000元
    MaxAmount:   10000000, // 100,000元
    InterestRate: 0.0680,  // 6.8%年利率
    LoanTermMonths: 12,
    AutoApprovalThreshold: 5000000, // 50,000元自动审批
    RiskLevel: "medium",
    CreditScoreMin: 600,
    CollateralRequired: false,
    Status: "active",
}

// 设置申请条件
conditions := model.LoanApplyConditions{
    MinAge: 18,
    MaxAge: 65,
    MinYearsExperience: 2,
    MinFarmArea: 5.0, // 最少5亩
    RequiredRegions: []string{"华北", "华东", "华中"},
    RequiredCertificates: []string{"身份证", "土地承包证"},
}
conditionsJSON, _ := json.Marshal(conditions)
product.ApplyConditions = string(conditionsJSON)

// 保存产品
if err := db.Create(product).Error; err != nil {
    return fmt.Errorf("创建贷款产品失败: %w", err)
}
```

### 提交贷款申请
```go
// 生成申请编号
applicationNo := generateApplicationNo()

// 创建贷款申请
application := &model.LoanApplication{
    ApplicationNo:     applicationNo,
    UserID:           userID,
    ProductID:        productID,
    ApplyAmount:      3000000, // 30,000元
    ApplyTermMonths:  12,
    LoanPurpose:      "购买春耕农资",
    ApplicantName:    "张三",
    ApplicantIDCard:  "110101199001011234",
    ApplicantPhone:   "13800138000",
    MonthlyIncome:    800000,  // 8,000元
    YearlyIncome:     9600000, // 96,000元
    FarmArea:         15.5,
    YearsOfExperience: 8,
    Status:           "pending",
}

// 设置种植作物
cropTypes := []string{"小麦", "玉米", "大豆"}
cropTypesJSON, _ := json.Marshal(cropTypes)
application.CropTypes = string(cropTypesJSON)

// 保存申请
if err := db.Create(application).Error; err != nil {
    return fmt.Errorf("创建贷款申请失败: %w", err)
}

// 创建提交日志
approvalLog := &model.ApprovalLog{
    ApplicationID: application.ID,
    ApproverID:    userID, // 申请人自己
    Action:        "submit",
    ApprovalLevel: 1,
    Result:        "pending",
    Comment:       "用户提交贷款申请",
    ActionTime:    time.Now(),
}

if err := db.Create(approvalLog).Error; err != nil {
    log.Printf("创建审批日志失败: %v", err)
}
```

### 审批流程处理
```go
// 审批通过
func ApproveApplication(applicationID uint64, approverID uint64, comment string) error {
    // 查询申请
    var application model.LoanApplication
    if err := db.First(&application, applicationID).Error; err != nil {
        return err
    }
    
    // 更新申请状态
    application.Status = "approved"
    application.FinalApprovedAt = &time.Time{}
    *application.FinalApprovedAt = time.Now()
    application.FinalApprover = &approverID
    
    // 设置审批结果
    application.ApprovedAmount = &application.ApplyAmount
    application.ApprovedTermMonths = &application.ApplyTermMonths
    approvedRate := 0.0680
    application.ApprovedRate = &approvedRate
    
    // 保存更新
    if err := db.Save(&application).Error; err != nil {
        return err
    }
    
    // 创建审批日志
    log := &model.ApprovalLog{
        ApplicationID:  applicationID,
        ApproverID:     approverID,
        Action:         "approve",
        ApprovalLevel:  application.ApprovalLevel,
        Result:         "approved",
        Comment:        comment,
        PreviousStatus: "reviewing",
        NewStatus:      "approved",
        ActionTime:     time.Now(),
    }
    
    return db.Create(log).Error
}
```

## 性能优化

### 查询优化
1. **分页查询**: 大量申请数据使用分页
2. **状态索引**: 基于申请状态的快速查询
3. **预加载**: 查询申请时预加载关联的产品和用户信息
4. **缓存策略**: 热门贷款产品信息缓存

### 审批流程优化
1. **异步处理**: AI风险评估使用异步任务
2. **批量操作**: 批量更新申请状态
3. **限流控制**: Dify API调用限流保护
4. **重试机制**: 失败请求自动重试

### 数据归档
1. **历史数据**: 定期归档已完成的贷款申请
2. **日志清理**: 清理过期的Dify调用日志
3. **统计预计算**: 预计算审批统计数据

## 安全考虑

1. **敏感信息加密**: 身份证号等敏感信息加密存储
2. **审批权限**: 基于角色的审批权限控制
3. **操作审计**: 完整的操作日志记录
4. **数据脱敏**: 日志中敏感信息脱敏处理
5. **API安全**: Dify API调用的安全认证 