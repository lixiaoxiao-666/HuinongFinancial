# Dify LLM智能审批工作流配置指南

## 概述

本文档基于LLM（大语言模型）设计的Dify智能审批工作流，通过LLM的强大理解和推理能力，简化数据解析、风险分析等环节，提供更智能、更灵活的审批流程。

## 前提条件

1. **后端服务运行状态确认**
   ```bash
   # 确认后端服务正常运行
   curl http://172.18.120.10:8080/livez
   curl http://172.18.120.10:8080/readyz
   ```

2. **获取AI Agent Token**
   ```bash
   # 从配置文件或环境变量获取
   echo $AI_AGENT_TOKEN
   # 或查看配置文件中的token设置
   ```

3. **Dify平台配置**
   - 访问地址：`http://172.18.120.57`
   - 确保已配置合适的LLM模型（如GPT-4、Claude等）
   - 确保可以正常访问并登录

## 第一步：更新自定义工具OpenAPI Schema

### 1.1 完整的OpenAPI Schema配置

基于后端接口实现，为Dify创建完整的OpenAPI 3.1规范：

```json
{
  "openapi": "3.1.0",
  "info": {
    "title": "慧农金融AI智能体接口",
    "description": "AI智能体审批工作流相关接口，支持Dify平台LLM调用",
    "version": "2.0.0",
    "contact": {
      "name": "慧农金融技术支持",
      "url": "http://172.18.120.10:8080"
    }
  },
  "servers": [
    {
      "url": "http://172.18.120.10:8080",
      "description": "开发环境"
    },
    {
      "url": "http://localhost:8080",
      "description": "本地开发环境"
    }
  ],
  "paths": {
    "/api/v1/ai-agent/applications/{application_id}/info": {
      "get": {
        "summary": "获取申请信息",
        "description": "获取贷款申请的详细信息，包含申请人信息、财务状况、产品信息等，供LLM进行综合分析",
        "operationId": "getApplicationInfo",
        "tags": ["AI智能体-数据获取"],
        "parameters": [
          {
            "name": "application_id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "申请ID，格式如：app_20240301_001"
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
                      "example": 0,
                      "description": "响应代码，0表示成功"
                    },
                    "message": {
                      "type": "string",
                      "example": "获取成功"
                    },
                    "data": {
                      "$ref": "#/components/schemas/ApplicationInfoResponse"
                    }
                  }
                }
              }
            }
          },
          "400": {
            "description": "请求参数错误",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "404": {
            "description": "申请不存在",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
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
        "description": "获取征信报告、银行流水、黑名单检查等外部数据，为LLM提供全面的风险评估数据源",
        "operationId": "getExternalData",
        "tags": ["AI智能体-数据获取"],
        "parameters": [
          {
            "name": "user_id",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "用户ID",
            "example": "user_001"
          },
          {
            "name": "data_types",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "数据类型，逗号分隔。可选值：credit_report（征信报告）,bank_flow（银行流水）,blacklist_check（黑名单检查）,government_subsidy（政府补贴）",
            "example": "credit_report,bank_flow,blacklist_check"
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
                      "$ref": "#/components/schemas/ExternalDataResponse"
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
        "description": "接收LLM分析后的AI决策结果，包含详细的分析报告和审批建议",
        "operationId": "submitAIDecision",
        "tags": ["AI智能体-决策提交"],
        "parameters": [
          {
            "name": "application_id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "申请ID"
          },
          {
            "name": "decision",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "enum": ["AUTO_APPROVED", "AUTO_REJECTED", "REQUIRE_HUMAN_REVIEW"]
            },
            "description": "AI决策结果"
          },
          {
            "name": "risk_score",
            "in": "query",
            "required": true,
            "schema": {
              "type": "number",
              "minimum": 0,
              "maximum": 1
            },
            "description": "风险分数(0-1)"
          },
          {
            "name": "risk_level",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "enum": ["LOW", "MEDIUM", "HIGH"]
            },
            "description": "风险等级"
          },
          {
            "name": "confidence_score",
            "in": "query",
            "required": true,
            "schema": {
              "type": "number",
              "minimum": 0,
              "maximum": 1
            },
            "description": "置信度(0-1)"
          },
          {
            "name": "approved_amount",
            "in": "query",
            "required": true,
            "schema": {
              "type": "number",
              "minimum": 0
            },
            "description": "批准金额"
          },
          {
            "name": "approved_term_months",
            "in": "query",
            "required": true,
            "schema": {
              "type": "integer",
              "minimum": 1
            },
            "description": "批准期限（月）"
          },
          {
            "name": "suggested_interest_rate",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "建议利率"
          },
          {
            "name": "analysis_summary",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "分析摘要"
          },
          {
            "name": "detailed_analysis",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string"
            },
            "description": "详细分析JSON字符串"
          },
          {
            "name": "recommendations",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string"
            },
            "description": "建议列表，逗号分隔"
          },
          {
            "name": "conditions",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string"
            },
            "description": "条件列表，逗号分隔"
          },
          {
            "name": "ai_model_version",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string"
            },
            "description": "AI模型版本"
          },
          {
            "name": "workflow_id",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string"
            },
            "description": "工作流ID"
          }
        ],
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
                          "enum": ["AI_APPROVED", "AI_REJECTED", "MANUAL_REVIEW_REQUIRED"]
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
        "description": "获取当前可用的AI模型配置、风险阈值和业务规则，为LLM提供决策参考依据",
        "operationId": "getAIModelConfig",
        "tags": ["AI智能体-配置"],
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
                      "$ref": "#/components/schemas/AIModelConfigResponse"
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
      "ApplicationInfoResponse": {
        "type": "object",
        "description": "申请信息详细响应",
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
              "category": {
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
          "application_info": {
            "type": "object",
            "properties": {
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
              },
              "submitted_at": {
                "type": "string",
                "format": "date-time"
              },
              "status": {
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
              "id_card_number": {
                "type": "string"
              },
              "phone": {
                "type": "string"
              },
              "address": {
                "type": "string"
              },
              "age": {
                "type": "integer"
              },
              "is_verified": {
                "type": "boolean"
              }
            }
          },
          "financial_info": {
            "type": "object",
            "properties": {
              "annual_income": {
                "type": "number",
                "description": "年收入"
              },
              "existing_loans": {
                "type": "integer",
                "description": "现有贷款数量"
              },
              "credit_score": {
                "type": "integer",
                "description": "信用分数"
              },
              "account_balance": {
                "type": "number",
                "description": "账户余额"
              },
              "land_area": {
                "type": "string",
                "description": "土地面积"
              },
              "farming_experience": {
                "type": "string",
                "description": "农业经验"
              }
            }
          },
          "uploaded_documents": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "doc_type": {
                  "type": "string"
                },
                "file_id": {
                  "type": "string"
                },
                "file_url": {
                  "type": "string"
                },
                "ocr_result": {
                  "type": "object"
                },
                "extracted_info": {
                  "type": "object"
                }
              }
            }
          },
          "external_data": {
            "type": "object",
            "properties": {
              "credit_bureau_score": {
                "type": "integer"
              },
              "blacklist_check": {
                "type": "boolean"
              },
              "previous_loan_history": {
                "type": "array",
                "items": {}
              },
              "land_ownership_verified": {
                "type": "boolean"
              }
            }
          }
        }
      },
      "ExternalDataResponse": {
        "type": "object",
        "properties": {
          "user_id": {
            "type": "string"
          },
          "credit_report": {
            "type": "object",
            "properties": {
              "score": {
                "type": "integer",
                "description": "征信分数"
              },
              "grade": {
                "type": "string",
                "description": "信用等级"
              },
              "report_date": {
                "type": "string",
                "description": "报告日期"
              },
              "loan_history": {
                "type": "array",
                "items": {}
              },
              "overdue_records": {
                "type": "integer",
                "description": "逾期记录数"
              }
            }
          },
          "bank_flow": {
            "type": "object",
            "properties": {
              "average_monthly_income": {
                "type": "number",
                "description": "月均收入"
              },
              "account_stability": {
                "type": "string",
                "description": "账户稳定性"
              },
              "last_6_months_flow": {
                "type": "array",
                "items": {
                  "type": "object",
                  "properties": {
                    "month": {
                      "type": "string"
                    },
                    "income": {
                      "type": "number"
                    },
                    "expense": {
                      "type": "number"
                    }
                  }
                }
              }
            }
          },
          "blacklist_check": {
            "type": "object",
            "properties": {
              "is_blacklisted": {
                "type": "boolean"
              },
              "check_time": {
                "type": "string"
              }
            }
          },
          "government_subsidy": {
            "type": "object",
            "properties": {
              "received_subsidies": {
                "type": "array",
                "items": {
                  "type": "object",
                  "properties": {
                    "year": {
                      "type": "integer"
                    },
                    "type": {
                      "type": "string"
                    },
                    "amount": {
                      "type": "number"
                    }
                  }
                }
              }
            }
          }
        }
      },
      "AIDecisionRequest": {
        "type": "object",
        "required": ["ai_analysis", "ai_decision", "processing_info"],
        "properties": {
          "ai_analysis": {
            "type": "object",
            "properties": {
              "risk_level": {
                "type": "string",
                "enum": ["LOW", "MEDIUM", "HIGH"],
                "description": "风险等级"
              },
              "risk_score": {
                "type": "number",
                "minimum": 0,
                "maximum": 1,
                "description": "风险分数(0-1)"
              },
              "confidence_score": {
                "type": "number",
                "minimum": 0,
                "maximum": 1,
                "description": "置信度(0-1)"
              },
              "analysis_summary": {
                "type": "string",
                "description": "分析摘要"
              },
              "detailed_analysis": {
                "type": "object",
                "description": "详细分析结果"
              },
              "risk_factors": {
                "type": "array",
                "items": {
                  "type": "object"
                },
                "description": "风险因素列表"
              },
              "recommendations": {
                "type": "array",
                "items": {
                  "type": "string"
                },
                "description": "建议事项"
              }
            }
          },
          "ai_decision": {
            "type": "object",
            "properties": {
              "decision": {
                "type": "string",
                "enum": ["AUTO_APPROVED", "AUTO_REJECTED", "REQUIRE_HUMAN_REVIEW"],
                "description": "AI决策结果"
              },
              "approved_amount": {
                "type": "number",
                "description": "批准金额"
              },
              "approved_term_months": {
                "type": "integer",
                "description": "批准期限（月）"
              },
              "suggested_interest_rate": {
                "type": "string",
                "description": "建议利率"
              },
              "conditions": {
                "type": "array",
                "items": {
                  "type": "string"
                },
                "description": "附加条件"
              },
              "next_action": {
                "type": "string",
                "description": "下一步行动"
              }
            }
          },
          "processing_info": {
            "type": "object",
            "properties": {
              "ai_model_version": {
                "type": "string",
                "description": "AI模型版本"
              },
              "processing_time_ms": {
                "type": "integer",
                "description": "处理时间（毫秒）"
              },
              "workflow_id": {
                "type": "string",
                "description": "工作流ID"
              },
              "processed_at": {
                "type": "string",
                "format": "date-time",
                "description": "处理时间"
              }
            }
          }
        }
      },
      "AIModelConfigResponse": {
        "type": "object",
        "properties": {
          "active_models": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "model_id": {
                  "type": "string"
                },
                "model_type": {
                  "type": "string"
                },
                "version": {
                  "type": "string"
                },
                "status": {
                  "type": "string"
                },
                "thresholds": {
                  "type": "object"
                }
              }
            }
          },
          "approval_rules": {
            "type": "object",
            "properties": {
              "auto_approval_threshold": {
                "type": "number"
              },
              "auto_rejection_threshold": {
                "type": "number"
              },
              "max_auto_approval_amount": {
                "type": "number"
              },
              "required_human_review_conditions": {
                "type": "array",
                "items": {
                  "type": "string"
                }
              }
            }
          },
          "business_parameters": {
            "type": "object",
            "properties": {
              "max_debt_to_income_ratio": {
                "type": "number"
              },
              "min_credit_score": {
                "type": "integer"
              },
              "max_loan_amount_by_category": {
                "type": "object"
              }
            }
          }
        }
      },
      "ErrorResponse": {
        "type": "object",
        "properties": {
          "code": {
            "type": "integer",
            "description": "错误代码"
          },
          "message": {
            "type": "string",
            "description": "错误信息"
          }
        }
      }
    },
    "securitySchemes": {
      "AIAgentToken": {
        "type": "apiKey",
        "in": "header",
        "name": "Authorization",
        "description": "AI Agent Token格式：'AI-Agent-Token your_token_here'"
      }
    }
  }
}
```

### 1.2 导入工具到Dify

1. **登录Dify平台**
   - 访问：`http://172.18.120.57`

2. **创建自定义工具**
   - 进入 `工具` → `自定义工具`
   - 点击 `创建工具`
   - 工具名称：`慧农金融AI智能体（LLM版）`
   - 描述：`基于LLM的智能审批接口工具`

3. **导入OpenAPI Schema**
   - 选择 `OpenAPI Schema` 导入方式
   - 复制上述完整JSON内容

4. **配置认证**
   - 认证方式：`API Key`
   - Header名称：`Authorization`
   - API Key值：`AI-Agent-Token your_actual_token_here`

## 第二步：创建LLM智能审批工作流

### 2.1 新建工作流应用

1. **创建工作流**
   - 应用名称：`LLM智能审批工作流`
   - 应用描述：`基于大语言模型的智能贷款审批系统`
   - 应用类型：`工作流`

### 2.2 配置开始节点

**输入变量配置**：
```json
{
  "application_id": {
    "type": "text",
    "required": true,
    "description": "贷款申请ID"
  },
  "callback_url": {
    "type": "text", 
    "required": false,
    "description": "处理完成后的回调地址"
  }
}
```

### 2.3 工作流节点配置

#### 节点1：获取申请信息
- **节点类型**：工具
- **工具选择**：慧农金融AI智能体（LLM版） → getApplicationInfo
- **参数配置**：
  - application_id: `{{start.application_id}}`

#### 节点2：获取外部数据
- **节点类型**：工具
- **工具选择**：慧农金融AI智能体（LLM版） → getExternalData
- **参数配置**：
  - user_id: `{{#1731652140017.text | jq '.data.applicant_info.user_id' | trim}}`
  - data_types: `credit_report,bank_flow,blacklist_check,government_subsidy`

#### 节点3：获取AI模型配置
- **节点类型**：工具
- **工具选择**：慧农金融AI智能体（LLM版） → getAIModelConfig
- **参数配置**：无需参数

#### 节点4：LLM智能分析（结构化输出版本）
- **节点类型**：LLM
- **模型选择**：GPT-4 或 Claude-3.5-sonnet（推荐）
- **结构化输出**：启用
- **输出模式**：JSON Schema

- **JSON Schema配置**：

```json
{
  "type": "object",
  "properties": {
    "analysis_summary": {
      "type": "string",
      "description": "风险分析摘要，150字内"
    },
    "risk_score": {
      "type": "number",
      "minimum": 0,
      "maximum": 1,
      "description": "风险分数(0-1)"
    },
    "risk_level": {
      "type": "string",
      "enum": ["LOW", "MEDIUM", "HIGH"],
      "description": "风险等级"
    },
    "confidence_score": {
      "type": "number",
      "minimum": 0,
      "maximum": 1,
      "description": "决策置信度(0-1)"
    },
    "decision": {
      "type": "string",
      "enum": ["AUTO_APPROVED", "AUTO_REJECTED", "REQUIRE_HUMAN_REVIEW"],
      "description": "AI决策结果"
    },
    "approved_amount": {
      "type": "number",
      "minimum": 0,
      "description": "批准金额"
    },
    "approved_term_months": {
      "type": "integer",
      "minimum": 1,
      "maximum": 360,
      "description": "批准期限（月）"
    },
    "suggested_interest_rate": {
      "type": "string",
      "description": "建议利率，如'4.5%'"
    },
    "detailed_analysis": {
      "type": "object",
      "properties": {
        "credit_analysis": {
          "type": "string",
          "description": "信用状况分析"
        },
        "financial_analysis": {
          "type": "string",
          "description": "财务状况分析"
        },
        "risk_factors": {
          "type": "array",
          "items": {
            "type": "string"
          },
          "description": "风险因素列表"
        },
        "strengths": {
          "type": "array",
          "items": {
            "type": "string"
          },
          "description": "申请优势列表"
        }
      },
      "required": ["credit_analysis", "financial_analysis", "risk_factors", "strengths"]
    },
    "recommendations": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "description": "建议事项"
    },
    "conditions": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "description": "批准条件"
    }
  },
  "required": [
    "analysis_summary",
    "risk_score", 
    "risk_level",
    "confidence_score",
    "decision",
    "approved_amount",
    "approved_term_months",
    "suggested_interest_rate",
    "detailed_analysis",
    "recommendations",
    "conditions"
  ]
}
```

- **系统提示词（简化版）**：

```
你是慧农金融的AI智能审批专家，负责对农业贷款申请进行全面的风险评估和决策建议。

## 你的职责：
1. 综合分析申请人的基本信息、财务状况、信用记录等
2. 评估贷款风险并给出量化的风险评分
3. 基于业务规则和风险阈值做出审批建议
4. 提供详细的分析报告和建议

## 分析要素：
### 基础信息分析
- 申请人身份信息完整性和真实性
- 申请金额与个人收入的匹配度
- 贷款用途的合理性

### 财务状况分析  
- 年收入水平及稳定性
- 债务收入比计算
- 资产负债状况
- 银行流水分析

### 信用风险分析
- 信用分数评估
- 历史贷款记录
- 逾期还款情况
- 黑名单检查结果

### 农业特色分析
- 农业经验和土地资源
- 政府补贴收入情况
- 季节性收入特点
- 农业市场风险

## 决策规则：
### 自动批准条件（AUTO_APPROVED）：
- 信用分数 ≥ 750
- 债务收入比 ≤ 30%
- 无黑名单记录
- 申请金额 ≤ 年收入的40%
- 风险评分 < 0.3

### 人工审核条件（REQUIRE_HUMAN_REVIEW）：
- 信用分数 600-749
- 债务收入比 30-50%
- 申请金额较大但在合理范围内
- 风险评分 0.3-0.7

### 自动拒绝条件（AUTO_REJECTED）：
- 信用分数 < 600
- 存在黑名单记录
- 债务收入比 > 50%
- 风险评分 > 0.7

## 输出要求：
你的输出将自动符合预定义的JSON结构。请确保：
1. risk_score为0-1之间的小数
2. confidence_score为0-1之间的小数
3. approved_amount不超过申请金额和产品最大额度
4. 所有数组字段至少包含一个元素
5. 决策逻辑必须符合上述规则

现在请分析以下申请：
```

- **用户提示词**：

```
## 申请信息
{{#1731652140017.text}}

## 外部数据
{{#1731652175020.text}}

## AI模型配置  
{{#1731652193039.text}}

请根据上述信息进行全面的风险评估和决策分析。
```

#### 节点5：格式化输出（简化版）
- **节点类型**：代码执行
- **编程语言**：Python3
- **输入变量**：
  - structured_output (Object): `{{#1731652265082.structured_output}}`

- **代码内容（结构化输出版）**：

```python
import json
from datetime import datetime

def main(structured_output: dict) -> dict:
    """
    结构化输出处理器 - 直接处理LLM的JSON结构化输出
    不再需要复杂的字符串解析和JSON提取
    """
    
    print(f"[DEBUG] 接收到结构化输出: {type(structured_output)}")
    print(f"[DEBUG] 包含字段: {list(structured_output.keys()) if isinstance(structured_output, dict) else 'Not a dict'}")
    
    try:
        # 验证输入数据
        if not isinstance(structured_output, dict):
            raise ValueError(f"输入不是字典类型，而是: {type(structured_output)}")
        
        # 填充默认值（防御性编程）
        data = fill_default_values(structured_output)
        
        # 验证和清理数据
        cleaned_data = validate_and_clean_data(data)
        
        # 创建API响应格式
        result = create_api_response(cleaned_data)
        
        print(f"[SUCCESS] 处理完成，决策: {cleaned_data.get('decision')}")
        return result
        
    except Exception as e:
        print(f"[ERROR] 处理异常: {str(e)}")
        return create_fallback_response(str(e))

def fill_default_values(data: dict) -> dict:
    """填充缺失的默认值"""
    
    defaults = {
        "analysis_summary": "AI风险分析",
        "risk_score": 0.5,
        "risk_level": "MEDIUM",
        "confidence_score": 0.5,
        "decision": "REQUIRE_HUMAN_REVIEW",
        "approved_amount": 0,
        "approved_term_months": 12,
        "suggested_interest_rate": "5.0%",
        "detailed_analysis": {
            "credit_analysis": "信用分析",
            "financial_analysis": "财务分析",
            "risk_factors": ["待评估"],
            "strengths": ["待评估"]
        },
        "recommendations": ["建议审核"],
        "conditions": ["需要审核"]
    }
    
    # 创建新的数据字典，保留原有数据，补充缺失项
    result = defaults.copy()
    result.update(data)
    
    # 特殊处理嵌套的detailed_analysis
    if "detailed_analysis" in data and isinstance(data["detailed_analysis"], dict):
        result["detailed_analysis"].update(data["detailed_analysis"])
    
    return result

def validate_and_clean_data(data: dict) -> dict:
    """验证和清理数据"""
    
    # 数值验证和修正
    try:
        data["risk_score"] = max(0.0, min(1.0, float(data["risk_score"])))
        data["confidence_score"] = max(0.0, min(1.0, float(data["confidence_score"])))
        data["approved_amount"] = max(0.0, float(data["approved_amount"]))
        data["approved_term_months"] = max(1, int(data["approved_term_months"]))
    except (ValueError, TypeError) as e:
        print(f"[WARNING] 数值修正: {e}")
        data["risk_score"] = 0.5
        data["confidence_score"] = 0.5
        data["approved_amount"] = 0.0
        data["approved_term_months"] = 12
    
    # 枚举值验证
    if data.get("risk_level") not in ["LOW", "MEDIUM", "HIGH"]:
        data["risk_level"] = "MEDIUM"
        print("[WARNING] risk_level修正为MEDIUM")
    
    if data.get("decision") not in ["AUTO_APPROVED", "AUTO_REJECTED", "REQUIRE_HUMAN_REVIEW"]:
        data["decision"] = "REQUIRE_HUMAN_REVIEW"
        print("[WARNING] decision修正为REQUIRE_HUMAN_REVIEW")
    
    # 数组验证
    for field in ["recommendations", "conditions"]:
        if not isinstance(data.get(field), list) or len(data.get(field, [])) == 0:
            data[field] = ["需要进一步评估"]
    
    # detailed_analysis验证
    if not isinstance(data.get("detailed_analysis"), dict):
        data["detailed_analysis"] = {
            "credit_analysis": "需要重新评估",
            "financial_analysis": "需要重新评估",
            "risk_factors": ["数据不完整"],
            "strengths": ["待评估"]
        }
    else:
        # 验证嵌套数组
        for field in ["risk_factors", "strengths"]:
            if not isinstance(data["detailed_analysis"].get(field), list):
                data["detailed_analysis"][field] = ["需要评估"]
    
    return data

def create_api_response(data: dict) -> dict:
    """创建API响应格式"""
    
    # 构建API请求对象
    api_request = {
        "ai_analysis": {
            "risk_level": data["risk_level"],
            "risk_score": float(data["risk_score"]),
            "confidence_score": float(data["confidence_score"]),
            "analysis_summary": data["analysis_summary"],
            "detailed_analysis": data["detailed_analysis"],
            "risk_factors": [
                {"factor": factor, "impact": "medium"} 
                for factor in data["detailed_analysis"].get("risk_factors", [])
            ],
            "recommendations": data["recommendations"]
        },
        "ai_decision": {
            "decision": data["decision"],
            "approved_amount": float(data["approved_amount"]),
            "approved_term_months": int(data["approved_term_months"]),
            "suggested_interest_rate": data["suggested_interest_rate"],
            "conditions": data["conditions"],
            "next_action": get_next_action(data["decision"])
        },
        "processing_info": {
            "ai_model_version": "LLM-v4.0-structured",
            "processing_time_ms": 1500,
            "workflow_id": "dify-llm-structured-output",
            "processed_at": datetime.now().isoformat()
        }
    }
    
    return {
        "success": 1,
        "api_request": json.dumps(api_request, ensure_ascii=False),
        "analysis_result": json.dumps(data, ensure_ascii=False),
        "decision": str(data["decision"]),
        "risk_score": float(data["risk_score"]),
        "risk_level": str(data["risk_level"]),
        "confidence_score": float(data["confidence_score"]),
        "approved_amount": float(data["approved_amount"]),
        "approved_term_months": int(data["approved_term_months"]),
        "suggested_interest_rate": str(data["suggested_interest_rate"]),
        "analysis_summary": str(data["analysis_summary"]),
        "error": ""
    }

def create_fallback_response(error_msg: str) -> dict:
    """创建降级响应"""
    
    fallback_data = {
        "analysis_summary": f"系统处理异常: {error_msg}，建议人工审核",
        "risk_score": 0.6,
        "risk_level": "MEDIUM",
        "confidence_score": 0.1,
        "decision": "REQUIRE_HUMAN_REVIEW",
        "approved_amount": 0.0,
        "approved_term_months": 12,
        "suggested_interest_rate": "5.0%",
        "detailed_analysis": {
            "credit_analysis": "系统异常，无法完成分析",
            "financial_analysis": "系统异常，无法完成分析",
            "risk_factors": ["系统处理异常"],
            "strengths": ["需要人工评估"]
        },
        "recommendations": ["转人工审核", "检查系统配置"],
        "conditions": ["系统异常，需要人工处理"]
    }
    
    api_request = {
        "ai_analysis": {
            "risk_level": "MEDIUM",
            "risk_score": 0.6,
            "confidence_score": 0.1,
            "analysis_summary": fallback_data["analysis_summary"],
            "detailed_analysis": fallback_data["detailed_analysis"],
            "risk_factors": [{"factor": "系统处理异常", "impact": "high"}],
            "recommendations": fallback_data["recommendations"]
        },
        "ai_decision": {
            "decision": "REQUIRE_HUMAN_REVIEW",
            "approved_amount": 0.0,
            "approved_term_months": 12,
            "suggested_interest_rate": "5.0%",
            "conditions": fallback_data["conditions"],
            "next_action": "ASSIGN_TO_REVIEWER"
        },
        "processing_info": {
            "ai_model_version": "LLM-v4.0-fallback",
            "processing_time_ms": 500,
            "workflow_id": "dify-llm-error-handler",
            "processed_at": datetime.now().isoformat()
        }
    }
    
    return {
        "success": 0,
        "api_request": json.dumps(api_request, ensure_ascii=False),
        "analysis_result": json.dumps(fallback_data, ensure_ascii=False),
        "decision": "REQUIRE_HUMAN_REVIEW",
        "risk_score": 0.6,
        "risk_level": "MEDIUM",
        "confidence_score": 0.1,
        "approved_amount": 0.0,
        "approved_term_months": 12,
        "suggested_interest_rate": "5.0%",
        "analysis_summary": fallback_data["analysis_summary"],
        "error": str(error_msg)
    }

def get_next_action(decision: str) -> str:
    """根据决策确定下一步行动"""
    action_map = {
        "AUTO_APPROVED": "GENERATE_CONTRACT",
        "AUTO_REJECTED": "SEND_REJECTION_NOTICE",
        "REQUIRE_HUMAN_REVIEW": "ASSIGN_TO_REVIEWER"
    }
    return action_map.get(decision, "MANUAL_REVIEW")
```

- **输出变量配置（修复版）**：

| 变量名 | 类型 | 描述 |
|--------|------|------|
| `success` | Number | 处理是否成功 (1=成功, 0=失败) |
| `api_request` | String | 格式化的API请求JSON |
| `analysis_result` | String | LLM分析结果JSON |
| `decision` | String | 审批决策 |
| `risk_score` | Number | 风险分数 |
| `risk_level` | String | 风险等级 |
| `confidence_score` | Number | 置信度分数 |
| `approved_amount` | Number | 批准金额 |
| `approved_term_months` | Number | 批准期限（月） |
| `suggested_interest_rate` | String | 建议利率 |
| `analysis_summary` | String | 分析摘要 |
| `error` | String | 错误信息（可选） |

#### 节点6：提交AI决策
- **节点类型**：工具
- **工具选择**：慧农金融AI智能体（LLM版） → submitAIDecision
- **参数配置**：
  - application_id: `{{start.application_id}}`
  - decision: `{{#结果验证.decision}}`
  - risk_score: `{{#结果验证.risk_score}}`
  - risk_level: `{{#结果验证.risk_level}}`
  - confidence_score: `{{#结果验证.confidence_score}}`
  - approved_amount: `{{#结果验证.approved_amount}}`
  - approved_term_months: `{{#结果验证.approved_term_months}}`
  - suggested_interest_rate: `{{#结果验证.suggested_interest_rate}}`
  - analysis_summary: `{{#结果验证.analysis_summary}}`
  - detailed_analysis: `{{#结果验证.analysis_result}}`
  - ai_model_version: `LLM-v4.0-structured`
  - workflow_id: `dify-llm-structured-output`

#### 节点7：结束节点
- **输出变量配置**：

```json
{
  "application_id": "{{start.application_id}}",
  "decision": "{{#1731652324556.decision}}",
  "risk_score": "{{#1731652324556.risk_score}}",
  "risk_level": "{{#1731652324556.risk_level}}",
  "approved_amount": "{{#1731652324556.approved_amount}}",
  "processing_status": "completed",
  "workflow_type": "LLM_BASED",
  "analysis_summary": "基于大语言模型的智能审批完成"
}
```

## 第三步：测试与验证

### 3.1 单节点测试

1. **测试工具连接**
   - 逐个测试API工具是否正常响应
   - 检查数据格式是否正确

2. **测试LLM节点**
   - 使用模拟数据测试LLM分析能力
   - 验证输出格式是否符合要求

### 3.2 端到端测试

**测试数据**：
```json
{
  "application_id": "test_app_001",
  "callback_url": "http://172.18.120.10:8080/callback"
}
```

**预期执行流程**：

1. ✅ 获取申请信息 → 返回完整的申请数据
2. ✅ 获取外部数据 → 返回征信、银行流水等数据  
3. ✅ 获取AI配置 → 返回业务规则和阈值
4. ✅ LLM智能分析 → 基于提示词生成决策分析
5. ✅ 解析LLM输出 → 格式化为API标准格式
6. ✅ 提交AI决策 → 更新申请状态
7. ✅ 工作流完成 → 返回最终结果

### 3.3 修复版本验证测试

**问题场景测试**：

使用真实的嵌套JSON输入进行测试：

```json
{
  "llm_output": "```json\n{\n  \"analysis_summary\": \"申请人张三信用分数高，无债务记录，收入稳定且有政府补贴收入支持，申请金额合理，风险较低。\",\n  \"risk_score\": 0.25,\n  \"risk_level\": \"LOW\",\n  \"confidence_score\": 0.85,\n  \"decision\": \"AUTO_APPROVED\",\n  \"approved_amount\": 100000,\n  \"approved_term_months\": 24,\n  \"suggested_interest_rate\": \"4.5%\",\n  \"detailed_analysis\": {\n    \"credit_analysis\": \"申请人张三的信用评分750，历史贷款记录无逾期还款情况，且不在黑名单中。\",\n    \"financial_analysis\": \"年收入稳定为20万元，银行流水显示每月平均收入6500元，没有现有债务。此外，该申请人每年还获得了政府补贴1200至3000元不等。\",\n    \"risk_factors\": [\"季节性收入波动可能对还款能力有影响\"],\n    \"strengths\": [\"信用评分高\", \"无现有贷款\"]\n  },\n  \"recommendations\": [\"建议定期监控借款人财务状况以应对农业市场风险\", \"鼓励申请人继续申请政府补贴以增加现金流稳定性\"],\n  \"conditions\": [\"信用分数750以上\", \"债务收入比为0%\", \"不在黑名单中\", \"申请金额10万元在年收入的40%以内\"]\n}\n```"
}
```

**期望的成功输出**：

```json
{
  "success": true,
  "decision": "AUTO_APPROVED",
  "risk_score": 0.25,
  "risk_level": "LOW", 
  "approved_amount": 100000.0,
  "api_request": "{...完整的API请求JSON...}",
  "analysis_result": "{...完整的分析结果JSON...}"
}
```

**测试步骤**：

1. **在Dify工作流中创建测试节点**
2. **输入上述JSON数据到Python脚本节点**
3. **观察日志输出**：
   ```
   [DEBUG] 收到原始输入，长度: 1234
   [DEBUG] 检测到嵌套格式
   [DEBUG] 成功提取嵌套内容
   [DEBUG] 移除开头```json标记
   [DEBUG] 移除结尾```标记
   [DEBUG] 提取JSON: {"analysis_summary":...
   [DEBUG] 解析成功，包含字段: ['analysis_summary', 'risk_score', ...]
   [SUCCESS] 处理完成
   ```

4. **验证输出结果包含正确的决策信息**

### 3.5 结构化输出版本验证测试

**新方案优势**：

✅ **无需字符串解析** - LLM直接输出JSON对象
✅ **避免转义字符问题** - 不再处理复杂的字符串转义
✅ **类型安全** - Dify自动验证JSON Schema
✅ **格式一致性** - LLM必须严格按照Schema输出
✅ **减少错误** - 大幅降低解析失败的可能性

**测试输入（模拟LLM结构化输出）**：

当LLM节点启用结构化输出后，它会直接输出如下格式的对象：

```json
{
  "analysis_summary": "申请人张三信用分数高，无债务记录，收入稳定且有政府补贴收入支持，申请金额合理，风险较低。",
  "risk_score": 0.25,
  "risk_level": "LOW",
  "confidence_score": 0.85,
  "decision": "AUTO_APPROVED",
  "approved_amount": 100000,
  "approved_term_months": 24,
  "suggested_interest_rate": "4.5%",
  "detailed_analysis": {
    "credit_analysis": "申请人张三的信用评分750，历史贷款记录无逾期还款情况，且不在黑名单中。",
    "financial_analysis": "年收入稳定为20万元，银行流水显示每月平均收入6500元，没有现有债务。此外，该申请人每年还获得了政府补贴1200至3000元不等。",
    "risk_factors": ["季节性收入波动可能对还款能力有影响"],
    "strengths": ["信用评分高", "无现有贷款"]
  },
  "recommendations": ["建议定期监控借款人财务状况以应对农业市场风险", "鼓励申请人继续申请政府补贴以增加现金流稳定性"],
  "conditions": ["信用分数750以上", "债务收入比为0%", "不在黑名单中", "申请金额10万元在年收入的40%以内"]
}
```

**期望的Python处理输出**：

```json
{
  "success": 1,
  "decision": "AUTO_APPROVED",
  "risk_score": 0.25,
  "risk_level": "LOW",
  "approved_amount": 100000.0,
  "api_request": "{...完整的API请求JSON...}",
  "analysis_result": "{...完整的分析结果JSON...}",
  "error": ""
}
```

**测试步骤**：

1. **配置LLM节点结构化输出**：
   - 在Dify中进入LLM节点设置
   - 启用"结构化输出"
   - 选择"JSON Schema"模式
   - 粘贴提供的Schema配置

2. **更新Python节点**：
   - 修改输入变量为：`structured_output (Object)`
   - 更新变量引用为：`{{#LLM节点ID.structured_output}}`
   - 使用简化版Python代码

3. **执行测试**：
   ```
   [DEBUG] 接收到结构化输出: <class 'dict'>
   [DEBUG] 包含字段: ['analysis_summary', 'risk_score', 'risk_level', ...]
   [SUCCESS] 处理完成，决策: AUTO_APPROVED
   ```

4. **验证输出格式**：
   - 确认所有输出变量类型正确
   - 验证数值在合理范围内
   - 检查枚举值符合预期

**故障排除**：

| 问题 | 原因 | 解决方案 |
|------|------|---------|
| "structured_output is not defined" | 输入变量配置错误 | 检查变量名和引用路径 |
| "输入不是字典类型" | LLM输出格式异常 | 检查JSON Schema配置 |
| 枚举值被修正 | LLM输出不符合Schema | 调整提示词或Schema |
| 数组字段为空 | LLM未正确填充数组 | 在Schema中设置minItems |

**配置验证清单**：

- [ ] LLM节点已启用结构化输出
- [ ] JSON Schema已正确配置
- [ ] Python输入变量类型为Object
- [ ] 变量引用路径正确：`{{#节点ID.structured_output}}`
- [ ] 所有required字段在Schema中定义
- [ ] 枚举值与业务逻辑匹配
- [ ] 数值字段设置了合理的minimum/maximum
- [ ] 数组字段至少有一个元素

### 3.4 LLM输出质量验证

**检查要点**：
- 风险分析的准确性和合理性
- 决策逻辑是否符合业务规则
- JSON格式是否标准且完整
- 置信度评估是否可信

## 第四步：优化与监控

### 4.1 提示词优化

基于测试结果持续优化：

1. **补充业务规则细节**
2. **增加边界案例处理**
3. **优化分析维度和权重**
4. **提升输出格式的稳定性**

### 4.2 Python脚本故障排除指南

**常见问题与解决方案**：

#### 问题1：解析失败 - "Failed to parse result"

**现象**：输出显示success: 0，返回错误响应
**原因**：LLM输出格式包含转义字符或嵌套结构
**解决方案**：
1. 使用提供的修复版Python代码
2. 检查日志中的DEBUG信息
3. 确认LLM输出格式是否为预期的JSON

#### 问题2：字符串长度限制

**现象**：长文本处理失败
**解决方案**：
```yaml
# docker-compose.yaml中调整
environment:
  - CODE_MAX_STRING_LENGTH=2000000  # 增加到200万字符
```

#### 问题3：LLM输出格式不一致

**现象**：有时成功有时失败
**解决方案**：
1. 在LLM系统提示词末尾强调输出格式要求
2. 使用温度参数为0或很低的值
3. 添加输出格式示例

#### 问题4：处理超时

**现象**：Python脚本执行超时
**解决方案**：
1. 优化代码逻辑，减少复杂度
2. 调整Dify的超时设置
3. 考虑分批处理大型数据

### 4.3 监控与日志分析

**关键日志位置**：
- Dify工作流执行日志
- Python代码执行节点的控制台输出
- Sandbox容器日志（如果启用）

**监控指标**：
- 解析成功率：目标 > 95%
- 平均处理时间：< 3秒
- 错误类型分布
- LLM输出格式一致性

**日志分析命令**：
```bash
# 查看Dify容器日志
docker logs dify_api_1 | grep "code_node"

# 查看sandbox日志  
docker logs dify_sandbox_1

# 检查错误模式
docker logs dify_api_1 | grep "Failed to parse result" | tail -10
```

### 4.4 性能监控

## 方案对比：字符串解析 vs 结构化输出

| 方面 | 旧方案（字符串解析） | 新方案（结构化输出） |
|------|-------------------|---------------------|
| **复杂度** | 高（复杂字符串处理） | 低（直接对象处理） |
| **稳定性** | 较低（易受格式影响） | 高（Schema验证） |
| **错误率** | 较高（解析失败常见） | 低（自动类型检查） |
| **调试难度** | 困难（字符转义问题） | 简单（清晰的对象结构） |
| **代码维护** | 复杂（多层解析逻辑） | 简单（主要是数据验证） |
| **性能** | 较慢（字符串处理） | 快（直接对象操作） |
| **可靠性** | 依赖LLM输出格式 | Dify平台保证格式 |
| **开发效率** | 低（需处理各种边界情况） | 高（专注业务逻辑） |

**推荐方案**：✅ **结构化输出方案**

### 迁移指南：从字符串解析到结构化输出

如果你当前使用的是字符串解析方案，按以下步骤迁移：

1. **备份现有工作流**
2. **修改LLM节点**：
   - 启用结构化输出
   - 配置JSON Schema
   - 移除提示词中的JSON格式要求
3. **更新Python节点**：
   - 输入变量：`llm_output (String)` → `structured_output (Object)`
   - 变量引用：`{{#节点ID.text}}` → `{{#节点ID.structured_output}}`
   - 代码：使用简化版处理逻辑
4. **测试验证**：
   - 确认输出格式正确
   - 验证所有字段类型匹配
5. **部署上线**

## 优势总结

### 传统代码工作流 vs LLM智能工作流

| 特性 | 传统代码工作流 | LLM智能工作流 |
|------|---------------|--------------|
| **节点数量** | 7-10个节点 | 5-6个节点 |
| **代码复杂度** | 高（复杂算法逻辑） | 低（数据处理+LLM分析） |
| **维护成本** | 高（需要调整算法） | 低（调整提示词+Schema） |
| **分析深度** | 固定规则 | 灵活理解 |
| **适应性** | 较差（硬编码规则） | 强（自然语言理解） |
| **决策解释性** | 一般（规则输出） | 优秀（自然语言解释） |
| **处理时间** | 快（<1秒） | 中等（2-5秒） |
| **成本** | 低（计算资源） | 中等（LLM调用费用） |
| **扩展性** | 困难（需修改代码） | 容易（调整提示词） |

### LLM工作流的优势

1. **简化复杂逻辑**：无需编写复杂的风险评估算法
2. **自然语言理解**：能够理解复杂的业务描述和规则
3. **灵活适应**：容易调整业务规则和评估标准
4. **决策解释**：提供自然语言的决策解释
5. **快速迭代**：通过调整提示词快速优化
6. **异常处理**：LLM具备较强的异常情况理解能力

### 适用场景

- 需要灵活决策逻辑的业务场景
- 规则经常变化的业务环境
- 需要详细决策解释的合规要求
- 处理复杂非结构化信息的场景

通过LLM工作流，您可以构建更智能、更灵活的审批系统，同时降低开发和维护成本。 

## 快速修复指南 - 解决API调用问题

### 如果遇到"请求参数格式错误"或"Request failed with status code 400"错误：

**问题原因**：
- Dify工具调用不支持复杂的Request Body格式
- 需要将JSON对象拆分为单独的查询参数

**立即解决步骤**：

1. **更新OpenAPI Schema**：
   - 将submitAIDecision的requestBody改为query参数
   - 使用提供的修复版OpenAPI配置

2. **更新工具配置**：
   - 删除Request Body参数
   - 添加单独的查询参数：
     ```
     decision: {{#结果验证.decision}}
     risk_score: {{#结果验证.risk_score}}
     risk_level: {{#结果验证.risk_level}}
     confidence_score: {{#结果验证.confidence_score}}
     approved_amount: {{#结果验证.approved_amount}}
     approved_term_months: {{#结果验证.approved_term_months}}
     suggested_interest_rate: {{#结果验证.suggested_interest_rate}}
     analysis_summary: {{#结果验证.analysis_summary}}
     detailed_analysis: {{#结果验证.analysis_result}}
     ```

3. **更新Python输出变量**：
   - 在Dify中确保Python节点输出包含所有必需字段
   - 变量类型必须匹配（Number/String）

4. **重新导入工具**：
   - 在Dify中删除旧的自定义工具
   - 重新导入修复版的OpenAPI Schema
   - 重新配置工具参数

**验证方法**：
- 检查每个参数是否正确映射
- 确认变量引用路径正确
- 测试API调用是否成功返回200状态码

**常见参数映射错误**：

| 错误写法 | 正确写法 |
|---------|---------|
| `{{#节点ID.api_request}}` | `{{#节点ID.decision}}` |
| `Request Body: {...}` | 分别配置每个参数 |
| 使用JSON对象 | 使用单独的字符串/数值 |

### 技术支持

如果按照以上步骤仍无法解决问题，请提供：
1. 完整的错误日志
2. LLM的原始输出内容
3. Python脚本的调试日志
4. Dify版本和部署方式

---

**文档版本**：v3.0-修复版  
**最后更新**：2025年5月29日  
**适用于**：Dify LLM工作流解析问题修复 