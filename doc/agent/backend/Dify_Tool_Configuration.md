# Dify工具配置与工作流集成文档

## 概述

本文档详细介绍如何在Dify平台中配置AI智能体接口工具，并创建完整的AI审批工作流。通过OpenAPI Schema导入自定义工具，实现与慧农金融后端系统的无缝集成。

## 1. OpenAPI Schema配置

### 1.1 AI智能体接口OpenAPI Schema

基于项目的API文档，为Dify创建完整的OpenAPI 3.1规范文件：

```json
{
  "openapi": "3.1.0",
  "info": {
    "title": "慧农金融AI智能体接口",
    "description": "AI智能体审批工作流相关接口，支持Dify平台调用",
    "version": "1.0.0",
    "contact": {
      "name": "慧农金融技术支持",
      "url": "http://localhost:8080"
    }
  },
  "servers": [
    {
      "url": "http://localhost:8080",
      "description": "本地开发环境"
    },
    {
      "url": "http://172.18.120.57:8080",
      "description": "内网测试环境"
    }
  ],
  "paths": {
    "/api/v1/ai-agent/applications/{application_id}/info": {
      "get": {
        "summary": "获取申请信息",
        "description": "获取贷款申请的详细信息用于AI分析",
        "operationId": "getApplicationInfo",
        "tags": ["AI智能体"],
        "parameters": [
          {
            "name": "application_id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "申请ID"
          }
        ],
        "responses": {
          "200": {
            "description": "成功获取申请信息",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "code": {
                      "type": "integer",
                      "example": 0
                    },
                    "message": {
                      "type": "string",
                      "example": "获取成功"
                    },
                    "data": {
                      "$ref": "#/components/schemas/ApplicationInfo"
                    }
                  }
                }
              }
            }
          }
        },
        "security": [
          {
            "AIAgentToken": []
          }
        ]
      }
    },
    "/api/v1/ai-agent/external-data": {
      "get": {
        "summary": "获取外部数据",
        "description": "获取征信、银行流水等外部数据用于AI分析",
        "operationId": "getExternalData",
        "tags": ["AI智能体"],
        "parameters": [
          {
            "name": "user_id",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "用户ID"
          },
          {
            "name": "data_types",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "数据类型，逗号分隔：credit_report,bank_flow,blacklist_check"
          }
        ],
        "responses": {
          "200": {
            "description": "成功获取外部数据",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "code": {
                      "type": "integer",
                      "example": 0
                    },
                    "message": {
                      "type": "string",
                      "example": "获取成功"
                    },
                    "data": {
                      "$ref": "#/components/schemas/ExternalData"
                    }
                  }
                }
              }
            }
          }
        },
        "security": [
          {
            "AIAgentToken": []
          }
        ]
      }
    },
    "/api/v1/ai-agent/applications/{application_id}/ai-decision": {
      "post": {
        "summary": "提交AI决策结果",
        "description": "Dify工作流完成分析后，提交AI决策结果",
        "operationId": "submitAIDecision",
        "tags": ["AI智能体"],
        "parameters": [
          {
            "name": "application_id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "申请ID"
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/AIDecisionRequest"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "AI决策结果处理成功",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "code": {
                      "type": "integer",
                      "example": 0
                    },
                    "message": {
                      "type": "string",
                      "example": "AI审批结果已成功处理"
                    },
                    "data": {
                      "type": "object",
                      "properties": {
                        "application_id": {
                          "type": "string"
                        },
                        "new_status": {
                          "type": "string",
                          "enum": ["AI_APPROVED", "AI_REJECTED", "REQUIRE_HUMAN_REVIEW"]
                        },
                        "next_step": {
                          "type": "string"
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
            "AIAgentToken": []
          }
        ]
      }
    },
    "/api/v1/ai-agent/config/models": {
      "get": {
        "summary": "获取AI模型配置",
        "description": "获取当前可用的AI模型配置和规则参数",
        "operationId": "getAIModelConfig",
        "tags": ["AI智能体"],
        "responses": {
          "200": {
            "description": "成功获取模型配置",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "code": {
                      "type": "integer",
                      "example": 0
                    },
                    "message": {
                      "type": "string",
                      "example": "获取成功"
                    },
                    "data": {
                      "$ref": "#/components/schemas/AIModelConfig"
                    }
                  }
                }
              }
            }
          }
        },
        "security": [
          {
            "AIAgentToken": []
          }
        ]
      }
    }
  },
  "components": {
    "schemas": {
      "ApplicationInfo": {
        "type": "object",
        "properties": {
          "application_id": {
            "type": "string",
            "description": "申请ID"
          },
          "product_info": {
            "type": "object",
            "properties": {
              "product_id": {
                "type": "string"
              },
              "name": {
                "type": "string"
              },
              "max_amount": {
                "type": "number"
              },
              "interest_rate_yearly": {
                "type": "string"
              }
            }
          },
          "applicant_info": {
            "type": "object",
            "properties": {
              "user_id": {
                "type": "string"
              },
              "real_name": {
                "type": "string"
              },
              "credit_score": {
                "type": "number"
              },
              "annual_income": {
                "type": "number"
              }
            }
          },
          "amount": {
            "type": "number",
            "description": "申请金额"
          },
          "term_months": {
            "type": "integer",
            "description": "申请期限（月）"
          },
          "purpose": {
            "type": "string",
            "description": "申请用途"
          }
        }
      },
      "ExternalData": {
        "type": "object",
        "properties": {
          "credit_report": {
            "type": "object",
            "properties": {
              "score": {
                "type": "number"
              },
              "level": {
                "type": "string"
              },
              "detail": {
                "type": "string"
              }
            }
          },
          "bank_flow": {
            "type": "object",
            "properties": {
              "average_monthly_income": {
                "type": "number"
              },
              "debt_to_income_ratio": {
                "type": "number"
              },
              "account_stability": {
                "type": "string"
              }
            }
          },
          "blacklist_check": {
            "type": "object",
            "properties": {
              "is_blacklisted": {
                "type": "boolean"
              },
              "risk_level": {
                "type": "string"
              }
            }
          }
        }
      },
      "AIDecisionRequest": {
        "type": "object",
        "required": ["decision", "risk_score", "risk_level"],
        "properties": {
          "decision": {
            "type": "string",
            "enum": ["AUTO_APPROVED", "AUTO_REJECTED", "REQUIRE_HUMAN_REVIEW"],
            "description": "AI决策结果"
          },
          "risk_score": {
            "type": "number",
            "minimum": 0,
            "maximum": 1,
            "description": "风险分数"
          },
          "risk_level": {
            "type": "string",
            "enum": ["LOW", "MEDIUM", "HIGH"],
            "description": "风险等级"
          },
          "confidence": {
            "type": "number",
            "minimum": 0,
            "maximum": 1,
            "description": "置信度"
          },
          "analysis_result": {
            "type": "object",
            "properties": {
              "credit_analysis": {
                "type": "string"
              },
              "income_analysis": {
                "type": "string"
              },
              "risk_factors": {
                "type": "array",
                "items": {
                  "type": "string"
                }
              },
              "approval_amount": {
                "type": "number"
              },
              "recommended_rate": {
                "type": "string"
              }
            }
          },
          "processing_info": {
            "type": "object",
            "properties": {
              "ai_model_version": {
                "type": "string"
              },
              "processing_time_ms": {
                "type": "number"
              },
              "workflow_id": {
                "type": "string"
              }
            }
          }
        }
      },
      "AIModelConfig": {
        "type": "object",
        "properties": {
          "available_models": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "model_id": {
                  "type": "string"
                },
                "model_name": {
                  "type": "string"
                },
                "version": {
                  "type": "string"
                },
                "status": {
                  "type": "string"
                }
              }
            }
          },
          "risk_thresholds": {
            "type": "object",
            "properties": {
              "auto_approve_threshold": {
                "type": "number"
              },
              "auto_reject_threshold": {
                "type": "number"
              }
            }
          },
          "decision_rules": {
            "type": "object",
            "properties": {
              "min_credit_score": {
                "type": "number"
              },
              "max_debt_ratio": {
                "type": "number"
              },
              "blacklist_auto_reject": {
                "type": "boolean"
              }
            }
          }
        }
      }
    },
    "securitySchemes": {
      "AIAgentToken": {
        "type": "apiKey",
        "in": "header",
        "name": "Authorization",
        "description": "AI Agent Token格式：'AI-Agent-Token your_token'"
      }
    }
  }
}
```

## 2. Dify工具配置步骤

### 2.1 导入自定义工具

1. **登录Dify平台**
   - 访问：`http://172.18.120.57`
   - 使用管理员账号登录

2. **创建自定义工具**
   - 进入 `工具` → `自定义工具`
   - 点击 `添加工具`
   - 选择 `OpenAPI`

3. **导入OpenAPI Schema**
   - 方式一：直接粘贴上述JSON内容
   - 方式二：上传JSON文件
   - 方式三：提供Schema URL（如果后端提供OpenAPI文档接口）

4. **配置工具认证**
   - 认证方式：选择 `API Key`
   - 在 `API Key` 字段中配置：
   ```
   Header Name: Authorization
   API Key: AI-Agent-Token your_actual_token_here
   ```

5. **测试工具连接**
   - 点击 `测试连接`
   - 验证每个接口是否正常响应

### 2.2 工具参数配置

配置完成后，Dify会自动解析出以下工具：

#### 工具1：获取申请信息
- **名称**：getApplicationInfo
- **参数**：application_id (必填)
- **用途**：获取申请详细信息供AI分析

#### 工具2：获取外部数据
- **名称**：getExternalData  
- **参数**：
  - user_id (必填)
  - data_types (必填，示例：`credit_report,bank_flow,blacklist_check`)
- **用途**：获取征信等外部数据

#### 工具3：提交AI决策
- **名称**：submitAIDecision
- **参数**：
  - application_id (必填)
  - Request Body (复杂对象，包含决策结果)
- **用途**：提交AI分析结果

#### 工具4：获取AI模型配置
- **名称**：getAIModelConfig
- **参数**：无
- **用途**：获取AI模型配置信息

## 3. Dify工作流配置

### 3.1 创建AI审批工作流

1. **新建工作流**
   - 进入 `工作室` → `创建应用` → `工作流`
   - 名称：`AI智能审批工作流`
   - 描述：`慧农金融贷款申请AI智能审批`

2. **配置开始节点**
   ```json
   输入参数：
   - application_id (文本，必填)：申请ID
   - callback_url (文本，可选)：回调地址
   ```

### 3.2 工作流节点配置

#### 节点1：获取申请信息
- **节点类型**：工具
- **选择工具**：慧农金融AI智能体接口 → getApplicationInfo
- **参数配置**：
  - application_id: `{{start.application_id}}`

#### 节点2：代码执行（解析申请信息）
- **节点类型**：代码执行
- **语言**：Python
- **代码**：
```python
import json

def main(application_info: str) -> dict:
    """解析申请信息，提取关键字段"""
    try:
        data = json.loads(application_info)['data']
        
        return {
            'user_id': data['applicant_info']['user_id'],
            'amount': data['amount'],
            'credit_score': data['applicant_info']['credit_score'],
            'annual_income': data['applicant_info']['annual_income'],
            'application_data': application_info
        }
    except Exception as e:
        return {
            'error': str(e),
            'user_id': '',
            'amount': 0,
            'credit_score': 0,
            'annual_income': 0,
            'application_data': application_info
        }
```

#### 节点3：获取外部数据
- **节点类型**：工具
- **选择工具**：慧农金融AI智能体接口 → getExternalData
- **参数配置**：
  - user_id: `{{code_parse.user_id}}`
  - data_types: `credit_report,bank_flow,blacklist_check`

#### 节点4：代码执行（AI风险分析）
- **节点类型**：代码执行
- **语言**：Python
- **代码**：
```python
import json
import math

def main(application_data: str, external_data: str) -> dict:
    """AI风险分析算法"""
    try:
        app_data = json.loads(application_data)['data']
        ext_data = json.loads(external_data)['data']
        
        # 基础信息
        amount = app_data['amount']
        credit_score = app_data['applicant_info']['credit_score']
        annual_income = app_data['applicant_info']['annual_income']
        
        # 外部数据
        credit_report = ext_data.get('credit_report', {})
        bank_flow = ext_data.get('bank_flow', {})
        blacklist = ext_data.get('blacklist_check', {})
        
        # 风险评分算法
        risk_score = 0.0
        risk_factors = []
        
        # 1. 信用分数评估 (权重: 0.4)
        if credit_score < 600:
            risk_score += 0.3
            risk_factors.append("信用分数过低")
        elif credit_score < 700:
            risk_score += 0.15
            risk_factors.append("信用分数中等")
        
        # 2. 收入债务比评估 (权重: 0.3)
        debt_ratio = bank_flow.get('debt_to_income_ratio', 0)
        if debt_ratio > 0.5:
            risk_score += 0.25
            risk_factors.append("债务收入比过高")
        elif debt_ratio > 0.3:
            risk_score += 0.1
            risk_factors.append("债务收入比偏高")
        
        # 3. 黑名单检查 (权重: 0.3)
        if blacklist.get('is_blacklisted', False):
            risk_score += 0.5
            risk_factors.append("存在黑名单记录")
        
        # 4. 申请金额与收入比 (权重: 0.2)
        amount_income_ratio = amount / max(annual_income, 1)
        if amount_income_ratio > 0.5:
            risk_score += 0.15
            risk_factors.append("申请金额与收入比过高")
        
        # 确保风险分数在0-1之间
        risk_score = min(risk_score, 1.0)
        
        # 决策逻辑
        if risk_score < 0.3:
            decision = "AUTO_APPROVED"
            risk_level = "LOW"
        elif risk_score < 0.7:
            decision = "REQUIRE_HUMAN_REVIEW"
            risk_level = "MEDIUM"
        else:
            decision = "AUTO_REJECTED"
            risk_level = "HIGH"
        
        # 计算建议额度
        if decision == "AUTO_APPROVED":
            approval_amount = amount
        elif decision == "REQUIRE_HUMAN_REVIEW":
            approval_amount = min(amount, annual_income * 0.3)
        else:
            approval_amount = 0
        
        return {
            'decision': decision,
            'risk_score': round(risk_score, 3),
            'risk_level': risk_level,
            'confidence': round(1 - risk_score * 0.5, 3),
            'risk_factors': risk_factors,
            'approval_amount': approval_amount,
            'analysis_summary': f"风险评分: {risk_score:.3f}, 决策: {decision}"
        }
        
    except Exception as e:
        return {
            'decision': 'REQUIRE_HUMAN_REVIEW',
            'risk_score': 0.5,
            'risk_level': 'MEDIUM',
            'confidence': 0.0,
            'risk_factors': [f"分析错误: {str(e)}"],
            'approval_amount': 0,
            'analysis_summary': f"系统分析错误: {str(e)}"
        }
```

#### 节点5：提交AI决策结果
- **节点类型**：工具
- **选择工具**：慧农金融AI智能体接口 → submitAIDecision
- **参数配置**：
  - application_id: `{{start.application_id}}`
  - Request Body:
```json
{
  "decision": "{{ai_analysis.decision}}",
  "risk_score": {{ai_analysis.risk_score}},
  "risk_level": "{{ai_analysis.risk_level}}",
  "confidence": {{ai_analysis.confidence}},
  "analysis_result": {
    "credit_analysis": "基于信用分数{{code_parse.credit_score}}的综合分析",
    "income_analysis": "年收入{{code_parse.annual_income}}，申请金额{{code_parse.amount}}",
    "risk_factors": {{ai_analysis.risk_factors}},
    "approval_amount": {{ai_analysis.approval_amount}},
    "recommended_rate": "根据风险等级确定"
  },
  "processing_info": {
    "ai_model_version": "v1.2.0",
    "processing_time_ms": 2500,
    "workflow_id": "{{sys.workflow_id}}"
  }
}
```

#### 节点6：结束节点
- **节点类型**：结束
- **输出变量**：
```json
{
  "decision": "{{ai_analysis.decision}}",
  "risk_score": "{{ai_analysis.risk_score}}",
  "analysis_summary": "{{ai_analysis.analysis_summary}}",
  "processing_status": "completed"
}
```

### 3.3 工作流发布与API调用

1. **调试工作流**
   - 使用测试数据验证每个节点
   - 确保数据流转正确

2. **发布工作流**
   - 点击右上角 `发布`
   - 选择发布为 `API`

3. **获取API信息**
   - 发布后获得工作流API地址
   - 配置到后端系统的`difyApiUrl`配置项

## 4. 错误处理与监控

### 4.1 异常处理配置

为关键节点启用错误处理：

1. **工具节点错误处理**
   - 启用 `失败时重试`
   - 设置重试次数：3次
   - 重试间隔：1000ms

2. **异常分支处理**
   - 为每个工具节点配置失败分支
   - 失败时返回默认值或错误信息

### 4.2 日志监控

1. **工作流日志**
   - 在Dify平台查看执行日志
   - 监控API调用成功率

2. **后端API日志**
   - 监控AI Agent接口调用情况
   - 记录决策结果和处理时间

## 5. 测试验证

### 5.1 工具连接测试

```bash
# 测试获取申请信息
curl -X GET "http://localhost:8080/api/v1/ai-agent/applications/test_app_001/info" \
  -H "Authorization: AI-Agent-Token your_token"

# 测试获取外部数据  
curl -X GET "http://localhost:8080/api/v1/ai-agent/external-data?user_id=user_001&data_types=credit_report,bank_flow" \
  -H "Authorization: AI-Agent-Token your_token"
```

### 5.2 工作流端到端测试

1. **准备测试数据**
   - 创建测试申请记录
   - 确保申请状态为 `SUBMITTED`

2. **触发工作流**
   - 通过Dify API调用工作流
   - 或在Dify平台手动触发

3. **验证结果**
   - 检查申请状态是否正确更新
   - 验证AI决策结果是否合理

## 6. 部署注意事项

### 6.1 环境配置

1. **网络连通性**
   - 确保Dify平台能访问后端API
   - 配置防火墙规则

2. **认证Token管理**
   - 安全存储AI Agent Token
   - 定期轮换Token

3. **性能优化**
   - 设置合理的超时时间
   - 配置连接池和重试机制

### 6.2 安全考虑

1. **API安全**
   - 使用HTTPS传输
   - 验证请求来源

2. **数据脱敏**
   - 敏感信息加密传输
   - 日志中避免记录敏感数据

这样配置完成后，您就有了一个完整的Dify AI智能审批工作流，能够自动处理贷款申请的风险评估和决策。 