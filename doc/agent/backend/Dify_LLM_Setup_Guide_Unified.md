# Dify LLM智能审批工作流配置指南 - 统一处理架构（v5.0）

## 概述

本文档基于**真正统一的处理架构**设计，实现一套工作流处理所有申请类型，消除接口冗余，提高系统维护性。

### 🎯 核心设计理念

- **单一入口**：一个统一接口处理所有申请类型
- **智能路由**：自动识别申请类型并路由到对应逻辑
- **统一输出**：标准化的响应格式，降低前端处理复杂度
- **易于扩展**：新增申请类型只需扩展而不需要新接口

### 📋 支持的申请类型
- 🏦 **金融贷款申请**：传统的贷款审批流程
- 🚜 **农机租赁申请**：农业机械设备租赁审批
- 🔮 **未来扩展**：保险申请、担保申请等

### 🏗️ 统一处理架构

```
单一入口 → 类型识别 → 智能路由 → 业务处理 → 统一输出
    ↓        ↓         ↓         ↓         ↓
  统一接口 → 申请分析 → 专业逻辑 → AI决策 → 标准响应
```

## 第一步：优化后的OpenAPI Schema（v5.0统一版）

### 1.1 精简的统一接口配置

```json
{
  "openapi": "3.1.0",
  "info": {
    "title": "慧农金融统一AI智能体接口（v5.0）",
    "description": "基于统一处理架构的AI智能体审批工作流接口，支持多种申请类型的统一处理",
    "version": "5.0.0"
  },
  "servers": [
    {
      "url": "http://172.18.120.10:8080",
      "description": "开发环境"
    }
  ],
  "paths": {
    "/api/v1/ai-agent/applications/{application_id}/info": {
      "get": {
        "summary": "获取申请信息（统一处理）",
        "description": "智能识别申请类型并返回对应的完整申请信息，支持贷款申请和农机租赁申请",
        "operationId": "getApplicationInfoUnified",
        "tags": ["统一处理"],
        "parameters": [
          {
            "name": "application_id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "申请ID，系统自动识别类型：贷款申请(test_app_*, app_*, loan_*)，农机租赁(ml_*, leasing_*)",
            "example": "test_app_001"
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
                      "example": "success"
                    },
                    "data": {
                      "$ref": "#/components/schemas/UnifiedApplicationInfo"
                    }
                  }
                }
              }
            }
          }
        },
        "security": [{"AIAgentToken": []}]
      }
    },
    "/api/v1/ai-agent/external-data/{user_id}": {
      "get": {
        "summary": "获取外部数据（智能适配）",
        "description": "根据用户类型和申请上下文智能获取相关外部数据",
        "operationId": "getExternalDataUnified",
        "tags": ["统一处理"],
        "parameters": [
          {
            "name": "user_id",
            "in": "path",
            "required": true,
            "schema": {"type": "string"},
            "description": "用户ID"
          },
          {
            "name": "data_types",
            "in": "query",
            "required": true,
            "schema": {"type": "string"},
            "description": "数据类型，系统会根据申请类型智能过滤：credit_report,bank_flow,blacklist_check,government_subsidy,farming_qualification"
          },
          {
            "name": "application_id",
            "in": "query",
            "required": false,
            "schema": {"type": "string"},
            "description": "申请ID，用于上下文识别"
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
                    "code": {"type": "integer", "example": 0},
                    "message": {"type": "string", "example": "success"},
                    "data": {"$ref": "#/components/schemas/UnifiedExternalDataResponse"}
                  }
                }
              }
            }
          }
        },
        "security": [{"AIAgentToken": []}]
      }
    },
    "/api/v1/ai-agent/applications/{application_id}/decisions": {
      "post": {
        "summary": "提交AI决策（智能路由）",
        "description": "接收LLM分析后的决策结果，系统自动识别申请类型并路由到对应的业务处理逻辑",
        "operationId": "submitAIDecisionUnified",
        "tags": ["统一处理"],
        "parameters": [
          {
            "name": "application_id",
            "in": "path",
            "required": true,
            "schema": {"type": "string"},
            "description": "申请ID"
          },
          {
            "name": "decision",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "enum": ["AUTO_APPROVED", "AUTO_REJECTED", "REQUIRE_HUMAN_REVIEW", "AUTO_APPROVE", "AUTO_REJECT", "REQUIRE_DEPOSIT_ADJUSTMENT"]
            },
            "description": "AI决策结果，系统会根据申请类型验证决策有效性"
          },
          {
            "name": "risk_score",
            "in": "query",
            "required": true,
            "schema": {"type": "number", "minimum": 0, "maximum": 1},
            "description": "风险分数(0-1)"
          },
          {
            "name": "risk_level",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "enum": ["LOW", "MEDIUM", "HIGH"]
            }
          },
          {
            "name": "confidence_score",
            "in": "query",
            "required": true,
            "schema": {"type": "number", "minimum": 0, "maximum": 1}
          },
          {
            "name": "analysis_summary",
            "in": "query",
            "required": true,
            "schema": {"type": "string"}
          },
          {
            "name": "approved_amount",
            "in": "query",
            "required": false,
            "schema": {"type": "number", "minimum": 0},
            "description": "批准金额（贷款申请）或建议租金（农机租赁）"
          },
          {
            "name": "approved_term_months",
            "in": "query",
            "required": false,
            "schema": {"type": "integer", "minimum": 1},
            "description": "批准期限（月，仅贷款申请需要）"
          },
          {
            "name": "suggested_interest_rate",
            "in": "query",
            "required": false,
            "schema": {"type": "string"},
            "description": "建议利率（仅贷款申请需要）"
          },
          {
            "name": "suggested_deposit",
            "in": "query",
            "required": false,
            "schema": {"type": "number", "minimum": 0},
            "description": "建议押金（仅农机租赁需要）"
          },
          {
            "name": "detailed_analysis",
            "in": "query",
            "required": false,
            "schema": {"type": "string"},
            "description": "详细分析JSON字符串"
          },
          {
            "name": "recommendations",
            "in": "query",
            "required": false,
            "schema": {"type": "string"},
            "description": "建议列表，逗号分隔"
          },
          {
            "name": "conditions",
            "in": "query",
            "required": false,
            "schema": {"type": "string"},
            "description": "条件列表，逗号分隔"
          },
          {
            "name": "ai_model_version",
            "in": "query",
            "required": false,
            "schema": {"type": "string"},
            "description": "AI模型版本"
          },
          {
            "name": "workflow_id",
            "in": "query",
            "required": false,
            "schema": {"type": "string"},
            "description": "工作流ID"
          }
        ],
        "responses": {
          "200": {
            "description": "AI决策处理成功",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "code": {"type": "integer", "example": 0},
                    "message": {"type": "string", "example": "AI决策提交成功"},
                    "data": {"$ref": "#/components/schemas/UnifiedDecisionResponse"}
                  }
                }
              }
            }
          }
        },
        "security": [{"AIAgentToken": []}]
      }
    },
    "/api/v1/ai-agent/config/models": {
      "get": {
        "summary": "获取AI模型配置（动态适配）",
        "description": "获取当前可用的AI模型配置，根据申请类型动态调整阈值和规则",
        "operationId": "getAIModelConfigUnified",
        "tags": ["统一处理"],
        "parameters": [
          {
            "name": "application_type",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "enum": ["LOAN_APPLICATION", "MACHINERY_LEASING", "AUTO_DETECT"]
            },
            "description": "申请类型，不传则返回所有类型的配置"
          }
        ],
        "responses": {
          "200": {
            "description": "成功获取模型配置",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "code": {"type": "integer", "example": 0},
                    "message": {"type": "string", "example": "success"},
                    "data": {"$ref": "#/components/schemas/UnifiedAIModelConfigResponse"}
                  }
                }
              }
            }
          }
        },
        "security": [{"AIAgentToken": []}]
      }
    },
    "/api/v1/ai-agent/logs": {
      "get": {
        "summary": "获取AI操作日志（统一查询）",
        "description": "查询AI操作的详细日志，支持多种申请类型的统一查询和过滤",
        "operationId": "getAIOperationLogs",
        "tags": ["统一处理"],
        "parameters": [
          {
            "name": "application_id",
            "in": "query",
            "required": false,
            "schema": {"type": "string"}
          },
          {
            "name": "application_type",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "enum": ["LOAN_APPLICATION", "MACHINERY_LEASING", "ALL"]
            }
          },
          {
            "name": "operation_type",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "enum": ["GET_INFO", "SUBMIT_DECISION", "GET_EXTERNAL_DATA", "ALL"]
            }
          },
          {
            "name": "page",
            "in": "query",
            "required": false,
            "schema": {"type": "integer", "minimum": 1, "default": 1}
          },
          {
            "name": "limit",
            "in": "query",
            "required": false,
            "schema": {"type": "integer", "minimum": 1, "maximum": 100, "default": 20}
          }
        ],
        "responses": {
          "200": {
            "description": "成功获取AI操作日志",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "code": {"type": "integer", "example": 0},
                    "message": {"type": "string", "example": "success"},
                    "data": {"$ref": "#/components/schemas/UnifiedAIOperationLogsResponse"}
                  }
                }
              }
            }
          }
        },
        "security": [{"AIAgentToken": []}]
      }
    }
  },
  "components": {
    "schemas": {
      "UnifiedApplicationInfo": {
        "type": "object",
        "description": "统一申请信息响应，根据申请类型动态调整字段",
        "properties": {
          "application_type": {
            "type": "string",
            "enum": ["LOAN_APPLICATION", "MACHINERY_LEASING"],
            "description": "申请类型标识"
          },
          "application_id": {
            "type": "string",
            "description": "申请ID"
          },
          "user_id": {
            "type": "string",
            "description": "申请人用户ID"
          },
          "status": {
            "type": "string",
            "description": "申请状态"
          },
          "submitted_at": {
            "type": "string",
            "format": "date-time",
            "description": "提交时间"
          },
          "basic_info": {
            "type": "object",
            "description": "基础信息，根据申请类型包含不同字段"
          },
          "business_info": {
            "type": "object",
            "description": "业务信息，贷款申请包含产品信息，农机租赁包含设备信息"
          },
          "applicant_info": {
            "type": "object",
            "description": "申请人信息，贷款申请为单人，农机租赁为承租方和出租方"
          },
          "financial_info": {
            "type": "object",
            "description": "财务信息，根据申请类型包含不同的财务数据"
          },
          "risk_assessment": {
            "type": "object",
            "description": "风险评估信息"
          },
          "documents": {
            "type": "array",
            "items": {
              "type": "object"
            },
            "description": "相关文档"
          }
        }
      },
      "UnifiedExternalDataResponse": {
        "type": "object",
        "description": "统一外部数据响应，根据申请类型智能过滤数据",
        "properties": {
          "user_id": {
            "type": "string"
          },
          "application_type": {
            "type": "string",
            "enum": ["LOAN_APPLICATION", "MACHINERY_LEASING"]
          },
          "data_types": {
            "type": "array",
            "items": {"type": "string"}
          },
          "credit_data": {
            "type": "object",
            "description": "征信数据（两种申请类型都需要）"
          },
          "bank_data": {
            "type": "object",
            "description": "银行流水数据"
          },
          "blacklist_data": {
            "type": "object",
            "description": "黑名单检查数据"
          },
          "government_data": {
            "type": "object",
            "description": "政府补贴数据（主要用于农业相关申请）"
          },
          "farming_data": {
            "type": "object",
            "description": "农业资质数据（主要用于农机租赁）"
          },
          "retrieved_at": {
            "type": "string",
            "format": "date-time"
          }
        }
      },
      "UnifiedDecisionResponse": {
        "type": "object",
        "description": "统一决策响应",
        "properties": {
          "application_id": {
            "type": "string"
          },
          "application_type": {
            "type": "string",
            "enum": ["LOAN_APPLICATION", "MACHINERY_LEASING"]
          },
          "decision": {
            "type": "string"
          },
          "new_status": {
            "type": "string",
            "description": "新的申请状态"
          },
          "next_step": {
            "type": "string",
            "description": "下一步操作"
          },
          "decision_id": {
            "type": "string",
            "description": "决策记录ID"
          },
          "ai_operation_id": {
            "type": "string",
            "description": "AI操作日志ID"
          },
          "processing_summary": {
            "type": "object",
            "description": "处理摘要信息"
          }
        }
      },
      "UnifiedAIModelConfigResponse": {
        "type": "object",
        "description": "统一AI模型配置响应",
        "properties": {
          "models": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "model_id": {"type": "string"},
                "model_type": {"type": "string"},
                "version": {"type": "string"},
                "status": {"type": "string"},
                "supported_application_types": {
                  "type": "array",
                  "items": {"type": "string"}
                }
              }
            }
          },
          "business_rules": {
            "type": "object",
            "properties": {
              "loan_application": {
                "type": "object",
                "properties": {
                  "auto_approval_threshold": {"type": "number"},
                  "auto_rejection_threshold": {"type": "number"},
                  "max_auto_approval_amount": {"type": "number"},
                  "required_human_review_conditions": {
                    "type": "array",
                    "items": {"type": "string"}
                  }
                }
              },
              "machinery_leasing": {
                "type": "object",
                "properties": {
                  "auto_approval_threshold": {"type": "number"},
                  "auto_rejection_threshold": {"type": "number"},
                  "max_auto_approval_deposit": {"type": "number"},
                  "required_human_review_conditions": {
                    "type": "array",
                    "items": {"type": "string"}
                  }
                }
              }
            }
          },
          "risk_thresholds": {
            "type": "object",
            "description": "风险阈值配置，根据申请类型动态应用"
          },
          "updated_at": {
            "type": "string",
            "format": "date-time"
          }
        }
      },
      "UnifiedAIOperationLogsResponse": {
        "type": "object",
        "properties": {
          "logs": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "operation_id": {"type": "string"},
                "application_id": {"type": "string"},
                "application_type": {"type": "string"},
                "operation_type": {"type": "string"},
                "decision": {"type": "string"},
                "risk_score": {"type": "number"},
                "confidence_score": {"type": "number"},
                "processing_time_ms": {"type": "integer"},
                "workflow_id": {"type": "string"},
                "ai_model_version": {"type": "string"},
                "created_at": {"type": "string", "format": "date-time"}
              }
            }
          },
          "pagination": {
            "type": "object",
            "properties": {
              "current_page": {"type": "integer"},
              "total_pages": {"type": "integer"},
              "total_count": {"type": "integer"},
              "limit": {"type": "integer"}
            }
          },
          "summary": {
            "type": "object",
            "description": "操作日志统计摘要",
            "properties": {
              "total_operations": {"type": "integer"},
              "by_application_type": {"type": "object"},
              "by_operation_type": {"type": "object"}
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
        "description": "AI Agent Token格式：'AI-Agent-Token your_token_here'"
      }
    }
  }
}
```

## 第二步：统一工作流配置（v5.0版）

### 2.1 工作流开始节点

**输入变量配置**：
```json
{
  "application_id": {
    "type": "text",
    "required": true,
    "description": "申请ID，系统自动识别类型"
  },
  "callback_url": {
    "type": "text", 
    "required": false,
    "description": "处理完成后的回调地址"
  }
}
```

### 2.2 核心节点配置

#### 节点1：获取申请信息（统一处理）
- **节点类型**：工具
- **工具选择**：慧农金融统一AI智能体 → getApplicationInfoUnified
- **参数配置**：
  - application_id: `{{start.application_id}}`

#### 节点2：获取外部数据（智能适配）
- **节点类型**：工具
- **工具选择**：慧农金融统一AI智能体 → getExternalDataUnified
- **参数配置**：
  - user_id: `{{#获取申请信息.text | jq '.data.user_id' | trim}}`
  - data_types: `credit_report,bank_flow,blacklist_check,government_subsidy,farming_qualification`
  - application_id: `{{start.application_id}}`

#### 节点3：获取AI模型配置（动态适配）
- **节点类型**：工具
- **工具选择**：慧农金融统一AI智能体 → getAIModelConfigUnified
- **参数配置**：
  - application_type: `{{#获取申请信息.text | jq '.data.application_type' | trim}}`

#### 节点4：LLM统一智能分析（增强版）
- **节点类型**：LLM
- **模型选择**：Claude-3.5-sonnet（推荐）
- **结构化输出**：启用

**增强版JSON Schema**：
```json
{
  "type": "object",
  "properties": {
    "application_type": {
      "type": "string",
      "enum": ["LOAN_APPLICATION", "MACHINERY_LEASING"],
      "description": "申请类型识别结果"
    },
    "type_confidence": {
      "type": "number",
      "minimum": 0,
      "maximum": 1,
      "description": "类型识别置信度"
    },
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
      "enum": ["AUTO_APPROVED", "AUTO_REJECTED", "REQUIRE_HUMAN_REVIEW", "AUTO_APPROVE", "AUTO_REJECT", "REQUIRE_DEPOSIT_ADJUSTMENT"],
      "description": "AI决策结果"
    },
    "business_specific_fields": {
      "type": "object",
      "description": "业务特定字段，根据申请类型包含不同内容",
      "properties": {
        "approved_amount": {
          "type": "number",
          "minimum": 0,
          "description": "批准金额（贷款）或建议租金（租赁）"
        },
        "approved_term_months": {
          "type": "integer",
          "minimum": 1,
          "description": "批准期限（仅贷款申请）"
        },
        "suggested_interest_rate": {
          "type": "string",
          "description": "建议利率（仅贷款申请）"
        },
        "suggested_deposit": {
          "type": "number",
          "minimum": 0,
          "description": "建议押金（仅农机租赁）"
        }
      }
    },
    "detailed_analysis": {
      "type": "object",
      "properties": {
        "primary_analysis": {"type": "string"},
        "secondary_analysis": {"type": "string"},
        "risk_factors": {
          "type": "array",
          "items": {"type": "string"}
        },
        "strengths": {
          "type": "array",
          "items": {"type": "string"}
        },
        "application_specific": {
          "type": "object",
          "description": "申请类型特定的分析"
        }
      },
      "required": ["primary_analysis", "secondary_analysis", "risk_factors", "strengths"]
    },
    "recommendations": {
      "type": "array",
      "items": {"type": "string"},
      "description": "建议事项"
    },
    "conditions": {
      "type": "array",
      "items": {"type": "string"},
      "description": "批准条件"
    }
  },
  "required": [
    "application_type",
    "type_confidence",
    "analysis_summary",
    "risk_score", 
    "risk_level",
    "confidence_score",
    "decision",
    "business_specific_fields",
    "detailed_analysis",
    "recommendations",
    "conditions"
  ]
}
```

**优化版系统提示词**：
```
你是慧农金融的统一AI智能审批专家（v5.0版），负责对多种类型的申请进行全面的风险评估和决策建议。

## 核心任务
1. **准确识别申请类型**：基于申请ID和数据结构特征
2. **应用专业分析框架**：根据申请类型使用对应的评估逻辑
3. **生成统一决策输出**：确保所有申请类型的输出格式一致

## 申请类型识别规则

### 贷款申请标识：
- ID格式：test_app_*, app_*, loan_*
- 数据特征：包含product_info, applicant_info, amount, term_months
- 关键字段：interest_rate, credit_score, annual_income

### 农机租赁申请标识：
- ID格式：ml_*, leasing_*
- 数据特征：包含lessee_info, lessor_info, machinery_info
- 关键字段：rental_days, deposit_amount, machinery_type

## 统一分析框架

### 通用评估要素：
1. **申请人信用分析**：信用历史、还款能力、风险记录
2. **财务状况评估**：收入稳定性、负债情况、资产状况
3. **外部环境因素**：行业风险、政策影响、市场环境

### 贷款申请专业逻辑：
- **决策枚举**：AUTO_APPROVED, AUTO_REJECTED, REQUIRE_HUMAN_REVIEW
- **核心指标**：债务收入比、信用分数、抵押物价值
- **输出重点**：approved_amount, approved_term_months, suggested_interest_rate

### 农机租赁专业逻辑：
- **决策枚举**：AUTO_APPROVE, AUTO_REJECT, REQUIRE_HUMAN_REVIEW, REQUIRE_DEPOSIT_ADJUSTMENT
- **核心指标**：设备状况、双方信用、租赁历史
- **输出重点**：suggested_deposit, rental_conditions

## 决策一致性要求
1. type_confidence ≥ 0.9 才能进行后续分析
2. risk_score 与 risk_level 必须匹配（<0.3=LOW, 0.3-0.7=MEDIUM, >0.7=HIGH）
3. decision 必须符合对应申请类型的枚举值
4. business_specific_fields 只包含申请类型相关的字段

现在请分析以下申请：
```

#### 节点5：智能决策提交（统一路由）
- **节点类型**：工具
- **工具选择**：慧农金融统一AI智能体 → submitAIDecisionUnified
- **参数配置**：
  - application_id: `{{start.application_id}}`
  - decision: `{{#LLM统一智能分析.structured_output.decision}}`
  - risk_score: `{{#LLM统一智能分析.structured_output.risk_score}}`
  - risk_level: `{{#LLM统一智能分析.structured_output.risk_level}}`
  - confidence_score: `{{#LLM统一智能分析.structured_output.confidence_score}}`
  - analysis_summary: `{{#LLM统一智能分析.structured_output.analysis_summary}}`
  - approved_amount: `{{#LLM统一智能分析.structured_output.business_specific_fields.approved_amount}}`
  - approved_term_months: `{{#LLM统一智能分析.structured_output.business_specific_fields.approved_term_months}}`
  - suggested_interest_rate: `{{#LLM统一智能分析.structured_output.business_specific_fields.suggested_interest_rate}}`
  - suggested_deposit: `{{#LLM统一智能分析.structured_output.business_specific_fields.suggested_deposit}}`
  - detailed_analysis: `{{#LLM统一智能分析.structured_output.detailed_analysis | json_encode}}`
  - recommendations: `{{#LLM统一智能分析.structured_output.recommendations | join(','')}}`
  - conditions: `{{#LLM统一智能分析.structured_output.conditions | join(',')}}`
  - ai_model_version: `LLM-v5.0-unified`
  - workflow_id: `dify-unified-v5`

#### 节点6：结束节点
- **输出变量配置**：
```json
{
  "application_id": "{{start.application_id}}",
  "application_type": "{{#LLM统一智能分析.structured_output.application_type}}",
  "type_confidence": "{{#LLM统一智能分析.structured_output.type_confidence}}",
  "decision": "{{#LLM统一智能分析.structured_output.decision}}",
  "risk_score": "{{#LLM统一智能分析.structured_output.risk_score}}",
  "risk_level": "{{#LLM统一智能分析.structured_output.risk_level}}",
  "processing_status": "completed",
  "workflow_version": "v5.0_unified",
  "analysis_summary": "基于统一处理架构的智能审批完成"
}
```

## 第三步：后端优化建议

### 3.1 统一服务接口设计

建议在后端创建一个真正的统一处理服务：

```go
// UnifiedApplicationProcessor 统一申请处理器
type UnifiedApplicationProcessor struct {
    loanService      *LoanService
    leasingService   *MachineryLeasingService
    aiService        *AIAgentService
    log              *zap.Logger
}

// ProcessApplicationUnified 统一处理申请
func (p *UnifiedApplicationProcessor) ProcessApplicationUnified(applicationID string) (*UnifiedApplicationResponse, error) {
    // 1. 自动识别申请类型
    appType, confidence := p.detectApplicationType(applicationID)
    
    // 2. 根据类型路由到专门的处理器
    switch appType {
    case "LOAN_APPLICATION":
        return p.processLoanApplication(applicationID)
    case "MACHINERY_LEASING":
        return p.processMachineryLeasing(applicationID)
    default:
        return nil, errors.New("unsupported application type")
    }
}

// detectApplicationType 智能识别申请类型
func (p *UnifiedApplicationProcessor) detectApplicationType(applicationID string) (string, float64) {
    // 基于ID模式识别
    if strings.HasPrefix(applicationID, "ml_") || strings.HasPrefix(applicationID, "leasing_") {
        return "MACHINERY_LEASING", 0.95
    }
    if strings.HasPrefix(applicationID, "test_app_") || strings.HasPrefix(applicationID, "app_") || strings.HasPrefix(applicationID, "loan_") {
        return "LOAN_APPLICATION", 0.95
    }
    
    // 基于数据库查询进一步确认
    // ... 实现数据库查询逻辑
    
    return "UNKNOWN", 0.0
}
```

### 3.2 接口迁移策略

1. **保留专用接口作为兼容性**：现有的专用接口可以保留，但标记为 `deprecated`
2. **统一接口作为主要入口**：新的集成都使用统一接口
3. **逐步迁移**：给现有用户时间迁移到新接口

## 第四步：优势总结

### ✅ 统一架构优势

1. **简化维护**：
   - 单一工作流处理所有申请类型
   - 减少接口数量和维护成本
   - 统一的错误处理和日志记录

2. **提高一致性**：
   - 标准化的响应格式
   - 一致的业务逻辑处理
   - 统一的监控和告警

3. **易于扩展**：
   - 新增申请类型只需扩展现有逻辑
   - 不需要新增接口和工作流
   - 配置驱动的业务规则

4. **更好的用户体验**：
   - 前端只需对接一套接口
   - 自动类型识别，无需用户指定
   - 统一的错误处理和提示

### 🚀 建议实施步骤

1. **Phase 1**：实现统一后端处理器
2. **Phase 2**：更新Dify工作流使用统一架构
3. **Phase 3**：前端迁移到统一接口
4. **Phase 4**：逐步下线专用接口

这样的设计真正实现了"统一处理"的目标，消除了接口冗余，提高了系统的可维护性和扩展性。 