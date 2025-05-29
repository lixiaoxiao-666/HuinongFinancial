# Dify 平台接入文档 - 数字惠农项目

## 1. Dify 平台介绍与准备

### 1.1 什么是 Dify
Dify 是一个开源的 LLM 应用开发平台，支持创建 AI 工作流、构建 AI 应用。我们将使用 Dify 来实现：
- 贷款申请的智能审批
- 农机租赁的智能处理
- 风险评估和决策建议

### 1.2 环境准备
1. 注册 Dify 账号：访问 [Dify 官网](https://dify.ai/) 
2. 创建工作区：选择适合的套餐（建议开发阶段使用免费版）
3. 获取 API 密钥：在设置页面获取 API Key

## 2. 工具创建 (Tools)

### 2.0 鉴权配置说明

#### API Token 配置
在数字惠农后端系统中配置专用的 Dify API Token：

**方法1: 配置文件方式**
```yaml
# backend/configs/config.yaml
dify:
  api_token: "dify-huinong-secure-token-2024"
```

**方法2: 环境变量方式**
```bash
export DIFY_API_TOKEN="dify-huinong-secure-token-2024"
```

#### Dify 工具鉴权配置
在 Dify 控制台创建工具时：
- **鉴权类型**: 选择 `Custom`
- **键名**: `Authorization`
- **键值**: `Bearer dify-huinong-secure-token-2024`

### 2.1 数字惠农 Dify 集成工具

创建一个包含所有必需API的综合工具，用于支持贷款审批和农机租赁的AI工作流。

```json
{
  "openapi": "3.1.0",
  "info": {
    "title": "数字惠农 Dify 集成工具",
    "description": "数字惠农后端系统 Dify AI 工作流集成工具，包含贷款审批和农机租赁相关API",
    "version": "v1.0.0"
  },
  "servers": [
    {
      "url": "{{BASE_URL}}"
    }
  ],
  "paths": {
    "/api/internal/dify/loan/get-application-details": {
      "post": {
        "description": "获取贷款申请的详细信息，包括用户基本信息、信用记录、财务状况等",
        "operationId": "GetLoanApplicationDetails",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "application_id": {
                    "type": "string",
                    "description": "贷款申请ID"
                  },
                  "user_id": {
                    "type": "string",
                    "description": "用户ID"
                  },
                  "include_credit": {
                    "type": "boolean",
                    "description": "是否包含征信信息",
                    "default": true
                  }
                },
                "required": ["application_id", "user_id"]
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "成功获取贷款申请详情",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "success": {"type": "boolean"},
                    "data": {
                      "type": "object",
                      "properties": {
                        "application": {
                          "type": "object",
                          "properties": {
                            "id": {"type": "string"},
                            "amount": {"type": "number"},
                            "term_months": {"type": "integer"},
                            "purpose": {"type": "string"},
                            "monthly_income": {"type": "number"},
                            "yearly_income": {"type": "number"},
                            "debt_amount": {"type": "number"}
                          }
                        },
                        "user": {
                          "type": "object",
                          "properties": {
                            "id": {"type": "string"},
                            "user_type": {"type": "string"},
                            "real_name_verified": {"type": "boolean"},
                            "bank_card_verified": {"type": "boolean"},
                            "credit_verified": {"type": "boolean"},
                            "years_of_experience": {"type": "integer"}
                          }
                        },
                        "credit_info": {
                          "type": "object",
                          "properties": {
                            "credit_score": {"type": "number"},
                            "debt_income_ratio": {"type": "number"},
                            "overdue_count": {"type": "integer"},
                            "max_overdue_days": {"type": "integer"}
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ]
      }
    },
    "/api/internal/dify/credit/query": {
      "post": {
        "description": "查询用户征信报告信息",
        "operationId": "QueryCreditReport",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "user_id": {
                    "type": "string",
                    "description": "用户ID"
                  },
                  "id_card": {
                    "type": "string",
                    "description": "身份证号"
                  },
                  "real_name": {
                    "type": "string",
                    "description": "真实姓名"
                  }
                },
                "required": ["user_id", "id_card", "real_name"]
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "成功获取征信报告",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "success": {"type": "boolean"},
                    "data": {
                      "type": "object",
                      "properties": {
                        "credit_score": {"type": "number"},
                        "debt_income_ratio": {"type": "number"},
                        "overdue_count": {"type": "integer"},
                        "max_overdue_days": {"type": "integer"},
                        "total_debt": {"type": "number"},
                        "monthly_payment": {"type": "number"},
                        "credit_history_months": {"type": "integer"},
                        "query_time": {"type": "string"}
                      }
                    }
                  }
                }
              }
            }
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ]
      }
    },
    "/api/internal/dify/loan/submit-assessment": {
      "post": {
        "description": "提交AI风险评估结果到后端系统",
        "operationId": "SubmitRiskAssessment",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "application_id": {
                    "type": "string",
                    "description": "贷款申请ID"
                  },
                  "risk_level": {
                    "type": "string",
                    "enum": ["low", "medium", "high"],
                    "description": "风险等级"
                  },
                  "decision": {
                    "type": "string",
                    "enum": ["approve", "reject", "manual"],
                    "description": "审批决策"
                  },
                  "recommended_amount": {
                    "type": "number",
                    "description": "建议批准金额"
                  },
                  "recommended_term": {
                    "type": "integer",
                    "description": "建议期限(月)"
                  },
                  "recommended_rate": {
                    "type": "number",
                    "description": "建议利率"
                  },
                  "risk_factors": {
                    "type": "array",
                    "items": {"type": "string"},
                    "description": "风险因素列表"
                  },
                  "comments": {
                    "type": "string",
                    "description": "评估意见"
                  },
                  "confidence_score": {
                    "type": "number",
                    "description": "置信度评分 0-1"
                  }
                },
                "required": ["application_id", "risk_level", "decision", "comments"]
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "成功提交风险评估结果",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "success": {"type": "boolean"},
                    "message": {"type": "string"},
                    "data": {
                      "type": "object",
                      "properties": {
                        "application_id": {"type": "string"},
                        "status": {"type": "string"},
                        "decision": {"type": "string"}
                      }
                    }
                  }
                }
              }
            }
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ]
      }
    },
    "/api/internal/dify/machine/get-rental-details": {
      "post": {
        "description": "获取农机租赁请求的详细信息",
        "operationId": "GetMachineRentalDetails",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "request_id": {
                    "type": "string",
                    "description": "租赁请求ID"
                  },
                  "user_id": {
                    "type": "string",
                    "description": "用户ID"
                  },
                  "machine_id": {
                    "type": "string",
                    "description": "农机ID"
                  }
                },
                "required": ["request_id", "user_id", "machine_id"]
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "成功获取农机租赁详情",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "success": {"type": "boolean"},
                    "data": {
                      "type": "object",
                      "properties": {
                        "request": {
                          "type": "object",
                          "properties": {
                            "id": {"type": "string"},
                            "start_time": {"type": "string"},
                            "end_time": {"type": "string"},
                            "location": {"type": "string"},
                            "has_conflict": {"type": "boolean"}
                          }
                        },
                        "user": {
                          "type": "object",
                          "properties": {
                            "id": {"type": "string"},
                            "user_type": {"type": "string"},
                            "real_name_verified": {"type": "boolean"},
                            "bank_card_verified": {"type": "boolean"},
                            "credit_verified": {"type": "boolean"}
                          }
                        },
                        "machine": {
                          "type": "object",
                          "properties": {
                            "id": {"type": "string"},
                            "name": {"type": "string"},
                            "type": {"type": "string"},
                            "status": {"type": "string"},
                            "hourly_rate": {"type": "number"},
                            "daily_rate": {"type": "number"}
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ]
      }
    }
  },
  "components": {
    "securitySchemes": {
      "BearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "{{API_TOKEN}}"
      }
    }
  }
}
```

### 2.2 环境变量配置

在 Dify 工作流中需要配置以下环境变量：

- **BASE_URL**: 您的后端服务地址 (如: `http://your-domain.com`)
- **API_TOKEN**: 在后端配置的 Dify API Token (如: `dify-huinong-secure-token-2024`)

### 2.3 可用操作说明

创建工具后，您可以在工作流中使用以下4个操作：

1. **GetLoanApplicationDetails** - 获取贷款申请详情
   - 用途：获取申请人基本信息、财务状况、征信记录
   - 输入：申请ID、用户ID
   - 输出：完整的申请数据

2. **QueryCreditReport** - 查询征信报告  
   - 用途：获取用户征信详细信息
   - 输入：用户ID、身份证、姓名
   - 输出：征信评分、负债情况等

3. **SubmitRiskAssessment** - 提交风险评估结果
   - 用途：将AI评估结果提交到后端系统
   - 输入：风险等级、决策、建议参数等
   - 输出：提交确认信息

4. **GetMachineRentalDetails** - 获取农机租赁详情
   - 用途：获取农机租赁请求信息
   - 输入：请求ID、用户ID、农机ID  
   - 输出：租赁详情、冲突检查结果

## 3. 工作流设计

### 3.1 贷款审批工作流设计

#### 工作流名称：`loan_risk_assessment`

#### 工作流步骤：

**步骤1: 开始节点**
- 节点类型：Start
- 输入变量：
  ```json
  {
    "application_id": {
      "type": "string",
      "required": true,
      "description": "贷款申请ID"
    },
    "user_id": {
      "type": "string", 
      "required": true,
      "description": "用户ID"
    }
  }
  ```

**步骤2: 获取申请详情**
- 节点类型：Tool
- 工具：get_loan_application_details
- 输入映射：
  ```json
  {
    "application_id": "{{#start.application_id#}}",
    "user_id": "{{#start.user_id#}}",
    "include_credit": true
  }
  ```
- 输出变量：`loan_details`

**步骤3: 解析申请数据**
- 节点类型：代码 (Python3)
- 作用：将工具返回的JSON数据解析为独立变量，便于后续节点使用
- 输入变量：`application_response`: `{{#get_loan_application_details.text#}}`

```python
import json

def main(application_response: str) -> dict:
    """解析贷款申请详情数据，提取关键字段"""
    try:
        data = json.loads(application_response)
        
        if not data.get('success', False):
            return {"error": data.get('error', '未知错误'), "success": False}
        
        app_data = data.get('data', {})
        application = app_data.get('application', {})
        user = app_data.get('user', {})
        credit = app_data.get('credit_info', {})
        
        return {
            "success": True,
            "application_id": application.get('id', ''),
            "loan_amount": float(application.get('amount', 0)),
            "term_months": int(application.get('term_months', 0)),
            "loan_purpose": application.get('purpose', ''),
            "monthly_income": float(application.get('monthly_income', 0)),
            "yearly_income": float(application.get('yearly_income', 0)),
            "debt_amount": float(application.get('debt_amount', 0)),
            "user_type": user.get('user_type', ''),
            "real_name_verified": bool(user.get('real_name_verified', False)),
            "bank_card_verified": bool(user.get('bank_card_verified', False)),
            "credit_verified": bool(user.get('credit_verified', False)),
            "years_of_experience": int(user.get('years_of_experience', 0)),
            "credit_score": float(credit.get('credit_score', 0)),
            "debt_income_ratio": float(credit.get('debt_income_ratio', 0)),
            "overdue_count": int(credit.get('overdue_count', 0)),
            "max_overdue_days": int(credit.get('max_overdue_days', 0)),
            "auth_completeness": calculate_auth_score(user),
            "risk_indicators": analyze_risk_factors(application, credit)
        }
    except Exception as e:
        return {"error": f"数据处理失败: {str(e)}", "success": False}

def calculate_auth_score(user_data: dict) -> str:
    """计算认证完整度"""
    score = sum([
        user_data.get('real_name_verified', False),
        user_data.get('bank_card_verified', False), 
        user_data.get('credit_verified', False)
    ])
    return ["incomplete", "partial", "partial", "complete"][score]

def analyze_risk_factors(application: dict, credit: dict) -> str:
    """分析风险因素"""
    risk_factors = []
    debt_ratio = credit.get('debt_income_ratio', 0)
    credit_score = credit.get('credit_score', 0)
    
    if debt_ratio > 0.5: risk_factors.append("高负债率")
    elif debt_ratio > 0.3: risk_factors.append("中等负债率")
    
    if credit_score < 600: risk_factors.append("信用分偏低")
    elif credit_score < 650: risk_factors.append("信用分中等")
    
    if credit.get('overdue_count', 0) > 0: risk_factors.append("存在逾期记录")
    
    return "; ".join(risk_factors) if risk_factors else "无明显风险"

**步骤4: 条件判断**
检查用户认证状态：
- 条件: `{{#parse_application_data.real_name_verified#}} == true AND {{#parse_application_data.bank_card_verified#}} == true`
- True分支: 继续AI评估
- False分支: 直接拒绝

**步骤5: LLM风险评估**
使用AI大模型进行智能风险评估：

系统提示词：
```
你是专业的农业金融风险评估专家。根据以下规则进行评估：

## 评估规则
**低风险 (approve)**：信用分>=700，债务收入比<=0.3，农业经验>=3年，认证完整
**中风险 (manual)**：信用分600-699，债务收入比0.3-0.5，农业经验>=1年  
**高风险 (reject)**：信用分<600，债务收入比>0.5，认证不完整

## 输出格式
返回标准JSON：
{
  "risk_level": "low|medium|high",
  "decision": "approve|manual|reject",
  "recommended_amount": 推荐金额,
  "recommended_term": 推荐期限,
  "recommended_rate": 推荐利率,
  "risk_factors": ["风险因素列表"],
  "comments": "详细评估意见",
  "confidence_score": 0.0-1.0置信度
}
```

用户提示词：
```
## 申请信息
- 申请金额：{{#parse_application_data.loan_amount#}}元
- 申请期限：{{#parse_application_data.term_months#}}个月
- 贷款用途：{{#parse_application_data.loan_purpose#}}

## 申请人状况  
- 用户类型：{{#parse_application_data.user_type#}}
- 农业经验：{{#parse_application_data.years_of_experience#}}年
- 认证状态：{{#parse_application_data.auth_completeness#}}
- 月收入：{{#parse_application_data.monthly_income#}}元
- 现有债务：{{#parse_application_data.debt_amount#}}元

## 征信评估
- 信用评分：{{#parse_application_data.credit_score#}}
- 债务收入比：{{#parse_application_data.debt_income_ratio#}}
- 逾期记录：{{#parse_application_data.overdue_count#}}次
- 风险指标：{{#parse_application_data.risk_indicators#}}

请进行全面风险评估并返回JSON结果。
```

**步骤6: 解析AI评估结果**
- 节点类型：代码 (Python3) 
- 作用：解析LLM返回的JSON结果，确保格式正确

```python
import json
import re

def main(ai_response: str) -> dict:
    """解析AI评估结果，确保格式正确"""
    try:
        # 提取JSON部分
        json_match = re.search(r'\{.*\}', ai_response, re.DOTALL)
        if json_match:
            result = json.loads(json_match.group())
        else:
            result = json.loads(ai_response)
        
        # 验证必需字段
        required_fields = ['risk_level', 'decision', 'comments']
        for field in required_fields:
            if field not in result:
                return {"error": f"缺少必需字段: {field}", "success": False}
        
        # 标准化数据类型
        return {
            "success": True,
            "risk_level": str(result.get('risk_level', 'medium')),
            "decision": str(result.get('decision', 'manual')),
            "recommended_amount": float(result.get('recommended_amount', 0)),
            "recommended_term": int(result.get('recommended_term', 12)),
            "recommended_rate": float(result.get('recommended_rate', 0.1)),
            "risk_factors": result.get('risk_factors', []),
            "comments": str(result.get('comments', '')),
            "confidence_score": float(result.get('confidence_score', 0.5))
        }
    except Exception as e:
        return {"error": f"解析失败: {str(e)}", "success": False}
```

**步骤7: 提交评估结果**
- 节点类型：Tool
- 工具：submit_risk_assessment
- 输入映射：
  ```json
  {
    "application_id": "{{#start.application_id#}}",
    "risk_level": "{{#parse_result.risk_level#}}",
    "decision": "{{#parse_result.decision#}}",
    "recommended_amount": "{{#parse_result.recommended_amount#}}",
    "recommended_term": "{{#parse_result.recommended_term#}}",
    "recommended_rate": "{{#parse_result.recommended_rate#}}",
    "risk_factors": "{{#parse_result.risk_factors#}}",
    "comments": "{{#parse_result.comments#}}",
    "confidence_score": "{{#parse_result.confidence_score#}}"
  }
  ```

**步骤8: 结束节点**
- 节点类型：End
- 输出：
  ```json
  {
    "success": true,
    "result": "{{#parse_result#}}",
    "message": "风险评估完成"
  }
  ```

### 3.2 农机租赁工作流设计

#### 工作流名称：`machine_rental_check`

#### 工作流步骤：

**步骤1: 开始节点**
- 输入变量：
  ```json
  {
    "request_id": {"type": "string", "required": true},
    "user_id": {"type": "string", "required": true},
    "machine_id": {"type": "string", "required": true}
  }
  ```

**步骤2: 获取租赁详情**
- 节点类型：Tool
- 工具：get_machine_rental_details

**步骤3: 简单规则检查**
- 节点类型：Code
- 代码类型：Python
- 代码：
  ```python
  def main(rental_details: dict) -> dict:
      data = rental_details.get('data', {})
      user = data.get('user', {})
      machine = data.get('machine', {})
      request = data.get('request', {})
      
      # 检查用户信用
      if not user.get('real_name_verified', False):
          return {
              "decision": "reject",
              "reason": "用户未完成实名认证"
          }
      
      # 检查农机可用性
      if machine.get('status') != 'available':
          return {
              "decision": "manual",
              "reason": "农机当前不可用，需人工协调"
          }
      
      # 检查时间冲突
      if request.get('has_conflict', False):
          return {
              "decision": "manual", 
              "reason": "租赁时间存在冲突"
          }
      
      # 通过基础检查
      return {
          "decision": "approve",
          "reason": "通过基础检查，可以租赁"
      }
  ```

## 4. 后端接口实现

### 4.1 Dify回调接口

您需要在后端实现以下接口供Dify调用：

```go
// /api/internal/dify/loan/get-application-details
func (h *DifyHandler) GetLoanApplicationDetails(c *gin.Context) {
    var req struct {
        ApplicationID string `json:"application_id"`
        UserID        string `json:"user_id"`
        IncludeCredit bool   `json:"include_credit"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"success": false, "error": err.Error()})
        return
    }
    
    // 实现获取贷款申请详情的逻辑
    // ...
    
    c.JSON(200, gin.H{
        "success": true,
        "data": map[string]interface{}{
            "application": application,
            "user": user,
            "credit_info": creditInfo,
        },
    })
}

// /api/internal/dify/loan/submit-assessment
func (h *DifyHandler) SubmitRiskAssessment(c *gin.Context) {
    var req struct {
        ApplicationID     string   `json:"application_id"`
        RiskLevel        string   `json:"risk_level"`
        Decision         string   `json:"decision"`
        RecommendedAmount float64 `json:"recommended_amount"`
        // ... 其他字段
    }
    
    // 实现提交风险评估结果的逻辑
    // ...
}
```

## 5. 环境变量配置

在Dify工作流中配置以下环境变量：

```
BASE_URL=https://your-backend-api.com
API_TOKEN=your-api-token
```

## 6. 测试工作流

### 6.1 单元测试
1. 在Dify控制台进入工作流编辑器
2. 点击"测试"按钮
3. 输入测试数据：
   ```json
   {
     "application_id": "LA20250530001",
     "user_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479"
   }
   ```
4. 查看执行结果和日志

### 6.2 集成测试
1. 在后端OA系统中调用Dify API
2. 监控工作流执行状态
3. 验证回调结果

## 7. 监控和维护

### 7.1 工作流监控
- 在Dify控制台查看工作流执行日志
- 监控成功率和执行时间
- 设置异常告警

### 7.2 性能优化
- 优化LLM提示词以提高准确性
- 调整工作流超时时间
- 缓存频繁查询的数据

## 8. 注意事项

1. **安全性**：确保API Token安全，设置IP白名单
2. **幂等性**：工作流可能会重试，确保接口支持幂等操作
3. **错误处理**：工作流中添加充分的错误处理逻辑
4. **日志记录**：记录详细的执行日志便于调试
5. **版本管理**：工作流变更时做好版本管理

这样就完成了完整的Dify接入配置。您可以根据实际需求调整工作流逻辑和参数。 