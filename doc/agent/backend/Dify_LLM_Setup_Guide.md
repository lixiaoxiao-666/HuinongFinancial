# Dify LLM智能审批工作流配置指南 - 统一多类型审批版

## 概述

本文档基于LLM（大语言模型）设计的Dify统一智能审批工作流，支持**多种申请类型的自动识别和处理**：

### 支持的申请类型
- 🏦 **金融贷款申请**：传统的贷款审批流程
- 🚜 **农机租赁申请**：农业机械设备租赁审批
- 🔮 **未来扩展**：可轻松扩展支持其他审批类型

### 核心优势
- **统一工作流**：一套Dify工作流处理多种申请类型
- **智能识别**：根据申请ID自动判断申请类型
- **业务解耦**：不同申请类型使用独立的分析逻辑
- **决策统一**：标准化的AI决策输出格式
- **日志完整**：全链路操作审计和追踪

### 技术架构

```
Dify工作流 → 统一接口 → 类型识别 → 分支处理 → 统一决策
    ↓            ↓          ↓         ↓         ↓
  LLM分析   →  申请信息   →  贷款/租赁  →  业务逻辑  →  结果输出
```

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

## 第一步：更新自定义工具OpenAPI Schema（统一版）

### 1.1 完整的统一多类型OpenAPI Schema配置

基于后端统一接口实现，为Dify创建支持多种申请类型的OpenAPI 3.1规范：

```json
{
  "openapi": "3.1.0",
  "info": {
    "title": "慧农金融统一AI智能体接口",
    "description": "支持多种申请类型的统一AI智能体审批工作流接口，包括贷款申请和农机租赁申请",
    "version": "4.0.0",
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
        "summary": "获取申请信息（统一接口）",
        "description": "统一获取申请信息，自动识别贷款申请或农机租赁申请，为LLM提供完整的分析数据",
        "operationId": "getApplicationInfoUnified",
        "tags": ["统一AI智能体"],
        "parameters": [
          {
            "name": "application_id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "申请ID，支持多种格式：贷款申请(test_app_001, app_xxx, loan_xxx)，农机租赁(ml_xxx, leasing_xxx)",
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
                      "example": 0,
                      "description": "响应代码，0表示成功"
                    },
                    "message": {
                      "type": "string",
                      "example": "success"
                    },
                    "data": {
                      "oneOf": [
                        {"$ref": "#/components/schemas/LoanApplicationInfo"},
                        {"$ref": "#/components/schemas/MachineryLeasingApplicationInfo"}
                      ]
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
    "/api/v1/ai-agent/external-data/{user_id}": {
      "get": {
        "summary": "获取外部数据（多类型支持）",
        "description": "获取征信报告、银行流水、黑名单检查等外部数据，支持贷款申请和农机租赁用户",
        "operationId": "getExternalDataUnified",
        "tags": ["统一AI智能体"],
        "parameters": [
          {
            "name": "user_id",
            "in": "path",
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
            "description": "数据类型，逗号分隔。可选值：credit_report,bank_flow,blacklist_check,government_subsidy,farming_qualification",
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
                      "example": "success"
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
    "/api/v1/ai-agent/applications/{application_id}/decisions": {
      "post": {
        "summary": "提交AI决策结果（统一接口）",
        "description": "接收LLM分析后的AI决策结果，自动识别申请类型并处理相应的业务逻辑",
        "operationId": "submitAIDecisionUnified",
        "tags": ["统一AI智能体"],
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
            "name": "application_type",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "enum": ["LOAN_APPLICATION", "MACHINERY_LEASING"]
            },
            "description": "申请类型，系统会自动识别"
          },
          {
            "name": "decision",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "enum": ["AUTO_APPROVED", "AUTO_REJECTED", "REQUIRE_HUMAN_REVIEW", "AUTO_APPROVE", "AUTO_REJECT", "REQUIRE_DEPOSIT_ADJUSTMENT"]
            },
            "description": "AI决策结果，支持贷款和农机租赁的不同决策类型"
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
            "name": "analysis_summary",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "分析摘要"
          },
          {
            "name": "approved_amount",
            "in": "query",
            "required": false,
            "schema": {
              "type": "number",
              "minimum": 0
            },
            "description": "批准金额（贷款申请专用）"
          },
          {
            "name": "approved_term_months",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "minimum": 1
            },
            "description": "批准期限（月，贷款申请专用）"
          },
          {
            "name": "suggested_interest_rate",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string"
            },
            "description": "建议利率（贷款申请专用）"
          },
          {
            "name": "suggested_deposit",
            "in": "query",
            "required": false,
            "schema": {
              "type": "number",
              "minimum": 0
            },
            "description": "建议押金（农机租赁专用）"
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
                      "example": "AI决策提交成功"
                    },
                    "data": {
                      "type": "object",
                      "properties": {
                        "application_id": {
                          "type": "string"
                        },
                        "application_type": {
                          "type": "string",
                          "enum": ["LOAN_APPLICATION", "MACHINERY_LEASING"]
                        },
                        "new_status": {
                          "type": "string",
                          "enum": ["AI_APPROVED", "AI_REJECTED", "MANUAL_REVIEW_REQUIRED", "DEPOSIT_ADJUSTMENT_REQUIRED"]
                        },
                        "next_step": {
                          "type": "string"
                        },
                        "decision_id": {
                          "type": "string"
                        },
                        "ai_operation_id": {
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
        "summary": "获取AI模型配置（多类型支持）",
        "description": "获取当前可用的AI模型配置、风险阈值和业务规则，支持多种申请类型",
        "operationId": "getAIModelConfigUnified",
        "tags": ["统一AI智能体"],
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
                      "example": "success"
                    },
                    "data": {
                      "$ref": "#/components/schemas/UnifiedAIModelConfigResponse"
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
    "/api/v1/ai-agent/machinery-leasing/applications/{application_id}": {
      "get": {
        "summary": "获取农机租赁申请信息（专用接口）",
        "description": "专门用于农机租赁申请的信息获取，提供更详细的农机租赁相关数据",
        "operationId": "getMachineryLeasingApplicationInfo",
        "tags": ["农机租赁专用"],
        "parameters": [
          {
            "name": "application_id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "农机租赁申请ID",
            "example": "ml_test_001"
          }
        ],
        "responses": {
          "200": {
            "description": "成功获取农机租赁申请信息",
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
                      "$ref": "#/components/schemas/MachineryLeasingApplicationInfo"
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
    "/api/v1/ai-agent/logs": {
      "get": {
        "summary": "获取AI操作日志",
        "description": "查询AI操作的详细日志，支持多种申请类型",
        "operationId": "getAIOperationLogs",
        "tags": ["AI操作日志"],
        "parameters": [
          {
            "name": "application_id",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string"
            },
            "description": "申请ID"
          },
          {
            "name": "application_type",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "enum": ["LOAN_APPLICATION", "MACHINERY_LEASING"]
            },
            "description": "申请类型"
          },
          {
            "name": "page",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "minimum": 1,
              "default": 1
            },
            "description": "页码"
          },
          {
            "name": "limit",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "minimum": 1,
              "maximum": 100,
              "default": 20
            },
            "description": "每页数量"
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
                    "code": {
                      "type": "integer",
                      "example": 0
                    },
                    "message": {
                      "type": "string",
                      "example": "success"
                    },
                    "data": {
                      "$ref": "#/components/schemas/AIOperationLogsResponse"
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
      "LoanApplicationInfo": {
        "type": "object",
        "description": "贷款申请信息响应",
        "properties": {
          "application_type": {
            "type": "string",
            "enum": ["LOAN_APPLICATION"],
            "description": "申请类型标识"
          },
          "application_id": {
            "type": "string",
            "description": "申请ID"
          },
          "product_info": {
            "type": "object",
            "properties": {
              "product_id": {"type": "string"},
              "name": {"type": "string"},
              "category": {"type": "string"},
              "max_amount": {"type": "number"},
              "interest_rate_yearly": {"type": "string"}
            }
          },
          "application_info": {
            "type": "object",
            "properties": {
              "amount": {"type": "number", "description": "申请金额"},
              "term_months": {"type": "integer", "description": "申请期限（月）"},
              "purpose": {"type": "string", "description": "申请用途"},
              "submitted_at": {"type": "string", "format": "date-time"},
              "status": {"type": "string"}
            }
          },
          "applicant_info": {
            "type": "object",
            "properties": {
              "user_id": {"type": "string"},
              "real_name": {"type": "string"},
              "id_card_number": {"type": "string"},
              "phone": {"type": "string"},
              "address": {"type": "string"},
              "age": {"type": "integer"},
              "is_verified": {"type": "boolean"}
            }
          },
          "financial_info": {
            "type": "object",
            "properties": {
              "annual_income": {"type": "number", "description": "年收入"},
              "existing_loans": {"type": "integer", "description": "现有贷款数量"},
              "credit_score": {"type": "integer", "description": "信用分数"},
              "account_balance": {"type": "number", "description": "账户余额"},
              "land_area": {"type": "string", "description": "土地面积"},
              "farming_experience": {"type": "string", "description": "农业经验"}
            }
          },
          "uploaded_documents": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "doc_type": {"type": "string"},
                "file_id": {"type": "string"},
                "file_url": {"type": "string"},
                "ocr_result": {"type": "object"},
                "extracted_info": {"type": "object"}
              }
            }
          },
          "external_data": {
            "type": "object",
            "properties": {
              "credit_bureau_score": {"type": "integer"},
              "blacklist_check": {"type": "boolean"},
              "previous_loan_history": {"type": "array", "items": {}},
              "land_ownership_verified": {"type": "boolean"}
            }
          }
        }
      },
      "MachineryLeasingApplicationInfo": {
        "type": "object",
        "description": "农机租赁申请信息响应",
        "properties": {
          "application_type": {
            "type": "string",
            "enum": ["MACHINERY_LEASING"],
            "description": "申请类型标识"
          },
          "application_id": {
            "type": "string",
            "description": "申请ID"
          },
          "lessee_info": {
            "type": "object",
            "properties": {
              "user_id": {"type": "string"},
              "real_name": {"type": "string"},
              "id_card_number": {"type": "string"},
              "phone": {"type": "string"},
              "address": {"type": "string"},
              "occupation": {"type": "string"},
              "annual_income": {"type": "number"},
              "farming_experience": {"type": "string"},
              "credit_rating": {"type": "string"},
              "is_verified": {"type": "boolean"},
              "previous_leasing_count": {"type": "integer"}
            }
          },
          "lessor_info": {
            "type": "object",
            "properties": {
              "user_id": {"type": "string"},
              "real_name": {"type": "string"},
              "phone": {"type": "string"},
              "business_name": {"type": "string"},
              "business_license": {"type": "string"},
              "verification_status": {"type": "string"},
              "credit_rating": {"type": "string"},
              "established_date": {"type": "string", "format": "date"},
              "average_rating": {"type": "number"},
              "successful_leasing_count": {"type": "integer"},
              "total_machinery_count": {"type": "integer"}
            }
          },
          "machinery_info": {
            "type": "object",
            "properties": {
              "machinery_id": {"type": "string"},
              "type": {"type": "string"},
              "brand_model": {"type": "string"},
              "engine_power": {"type": "string"},
              "manufacturing_year": {"type": "integer"},
              "condition": {"type": "string"},
              "daily_rent": {"type": "number"},
              "deposit": {"type": "number"},
              "location": {"type": "string"},
              "availability": {"type": "boolean"},
              "last_maintenance": {"type": "string", "format": "date"},
              "insurance_status": {"type": "string"},
              "insurance_expiry": {"type": "string", "format": "date"}
            }
          },
          "leasing_details": {
            "type": "object",
            "properties": {
              "requested_start_date": {"type": "string", "format": "date"},
              "requested_end_date": {"type": "string", "format": "date"},
              "rental_days": {"type": "integer"},
              "total_amount": {"type": "number"},
              "deposit_amount": {"type": "number"},
              "usage_purpose": {"type": "string"},
              "work_location": {"type": "string"},
              "estimated_work_area": {"type": "string"},
              "special_requirements": {"type": "string"}
            }
          },
          "risk_assessment": {
            "type": "object",
            "properties": {
              "lessee_credit_score": {"type": "integer"},
              "lessor_reliability": {"type": "string"},
              "machinery_condition": {"type": "string"},
              "insurance_status": {"type": "string"},
              "seasonal_risk": {"type": "string"},
              "weather_forecast": {"type": "string"},
              "regional_activity": {"type": "string"}
            }
          }
        }
      },
      "ExternalDataResponse": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "credit_report": {
            "type": "object",
            "properties": {
              "score": {"type": "integer", "description": "征信分数"},
              "grade": {"type": "string", "description": "信用等级"},
              "report_date": {"type": "string", "description": "报告日期"},
              "loan_history": {"type": "array", "items": {}},
              "overdue_records": {"type": "integer", "description": "逾期记录数"}
            }
          },
          "bank_flow": {
            "type": "object",
            "properties": {
              "average_monthly_income": {"type": "number", "description": "月均收入"},
              "account_stability": {"type": "string", "description": "账户稳定性"},
              "last_6_months_flow": {
                "type": "array",
                "items": {
                  "type": "object",
                  "properties": {
                    "month": {"type": "string"},
                    "income": {"type": "number"},
                    "expense": {"type": "number"}
                  }
                }
              }
            }
          },
          "blacklist_check": {
            "type": "object",
            "properties": {
              "is_blacklisted": {"type": "boolean"},
              "check_time": {"type": "string"}
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
                    "year": {"type": "integer"},
                    "type": {"type": "string"},
                    "amount": {"type": "number"}
                  }
                }
              }
            }
          },
          "farming_qualification": {
            "type": "object",
            "description": "农业资质信息（农机租赁专用）",
            "properties": {
              "certification_level": {"type": "string"},
              "experience_years": {"type": "integer"},
              "machinery_operation_skills": {"type": "array", "items": {"type": "string"}}
            }
          }
        }
      },
      "UnifiedAIModelConfigResponse": {
        "type": "object",
        "properties": {
          "loan_approval": {
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
          },
          "risk_assessment_models": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "model_id": {"type": "string"},
                "model_type": {"type": "string"},
                "version": {"type": "string"},
                "status": {"type": "string"}
              }
            }
          }
        }
      },
      "AIOperationLogsResponse": {
        "type": "object",
        "properties": {
          "logs": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "id": {"type": "string"},
                "application_id": {"type": "string"},
                "application_type": {"type": "string", "enum": ["LOAN_APPLICATION", "MACHINERY_LEASING"]},
                "operation_type": {"type": "string"},
                "ai_model_version": {"type": "string"},
                "decision": {"type": "string"},
                "risk_score": {"type": "number"},
                "confidence_score": {"type": "number"},
                "processing_time_ms": {"type": "integer"},
                "workflow_id": {"type": "string"},
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
          }
        }
      },
      "ErrorResponse": {
        "type": "object",
        "properties": {
          "code": {"type": "integer", "description": "错误代码"},
          "message": {"type": "string", "description": "错误信息"}
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
   - 工具名称：`慧农金融统一AI智能体（v4.0多类型支持）`
   - 描述：`支持贷款申请和农机租赁申请的统一AI审批接口工具，含日志查询功能`

3. **导入OpenAPI Schema**
   - 选择 `OpenAPI Schema` 导入方式
   - 复制上述完整JSON内容

4. **配置认证**
   - 认证方式：`API Key`
   - Header名称：`Authorization`
   - API Key值：`AI-Agent-Token your_actual_token_here`

## 第二步：创建统一LLM智能审批工作流

### 2.1 新建工作流应用

1. **创建工作流**
   - 应用名称：`统一AI智能审批工作流（多类型）`
   - 应用描述：`基于大语言模型的统一智能审批系统，支持贷款申请和农机租赁申请`
   - 应用类型：`工作流`

### 2.2 配置开始节点

**输入变量配置**：
```json
{
  "application_id": {
    "type": "text",
    "required": true,
    "description": "申请ID，支持贷款申请(test_app_001, app_xxx)和农机租赁申请(ml_xxx, leasing_xxx)"
  },
  "callback_url": {
    "type": "text", 
    "required": false,
    "description": "处理完成后的回调地址"
  }
}
```

### 2.3 统一工作流节点配置

#### 节点1：获取申请信息（统一接口）
- **节点类型**：工具
- **工具选择**：慧农金融统一AI智能体（多类型支持） → getApplicationInfoUnified
- **参数配置**：
  - application_id: `{{start.application_id}}`

#### 节点2：获取外部数据（多类型支持）
- **节点类型**：工具
- **工具选择**：慧农金融统一AI智能体（多类型支持） → getExternalDataUnified
- **参数配置**：
  - user_id: `{{#获取申请信息.text | jq '.data.applicant_info.user_id // .data.lessee_info.user_id' | trim}}`
  - data_types: `credit_report,bank_flow,blacklist_check,government_subsidy,farming_qualification`

#### 节点3：获取AI模型配置（多类型支持）
- **节点类型**：工具
- **工具选择**：慧农金融统一AI智能体（多类型支持） → getAIModelConfigUnified
- **参数配置**：无需参数

#### 节点4：LLM统一智能分析（结构化输出版本）
- **节点类型**：LLM
- **模型选择**：GPT-4o 或 Claude-3.5-sonnet（推荐）
- **结构化输出**：启用
- **输出模式**：JSON Schema

- **JSON Schema配置（统一版）**：

```json
{
  "type": "object",
  "properties": {
    "application_type": {
      "type": "string",
      "enum": ["LOAN_APPLICATION", "MACHINERY_LEASING"],
      "description": "申请类型识别结果"
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
    "approved_amount": {
      "type": "number",
      "minimum": 0,
      "description": "批准金额（贷款申请）或建议租金（农机租赁）"
    },
    "approved_term_months": {
      "type": "integer",
      "minimum": 1,
      "maximum": 360,
      "description": "批准期限（月，贷款申请专用）"
    },
    "suggested_interest_rate": {
      "type": "string",
      "description": "建议利率，如'4.5%'（贷款申请专用）"
    },
    "suggested_deposit": {
      "type": "number",
      "minimum": 0,
      "description": "建议押金（农机租赁专用）"
    },
    "detailed_analysis": {
      "type": "object",
      "properties": {
        "primary_analysis": {
          "type": "string",
          "description": "主要分析（信用分析或承租方分析）"
        },
        "secondary_analysis": {
          "type": "string",
          "description": "次要分析（财务分析或出租方分析）"
        },
        "asset_analysis": {
          "type": "string",
          "description": "资产分析（抵押物或农机设备）"
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
      "required": ["primary_analysis", "secondary_analysis", "risk_factors", "strengths"]
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
    "application_type",
    "analysis_summary",
    "risk_score", 
    "risk_level",
    "confidence_score",
    "decision",
    "approved_amount",
    "detailed_analysis",
    "recommendations",
    "conditions"
  ]
}
```

- **系统提示词（统一多类型版）**：

```
你是慧农金融的统一AI智能审批专家，负责对多种类型的申请进行全面的风险评估和决策建议。

## 申请类型识别
首先识别申请类型：
- **贷款申请**：ID格式如 test_app_001, app_xxx, loan_xxx，包含product_info和applicant_info
- **农机租赁申请**：ID格式如 ml_xxx, leasing_xxx，包含lessee_info和lessor_info

## 贷款申请分析框架

### 分析要素：
1. **申请人基础信息**：身份信息完整性、年龄、职业稳定性
2. **财务状况分析**：年收入水平、债务收入比、资产负债状况
3. **信用风险分析**：信用分数、历史记录、黑名单检查
4. **农业特色分析**：农业经验、土地资源、政府补贴、季节性收入

### 决策规则：
- **自动批准(AUTO_APPROVED)**：信用分数≥750，债务收入比≤30%，无黑名单，风险评分<0.3
- **人工审核(REQUIRE_HUMAN_REVIEW)**：信用分数600-749，债务收入比30-50%，风险评分0.3-0.7
- **自动拒绝(AUTO_REJECTED)**：信用分数<600，存在黑名单，债务收入比>50%，风险评分>0.7

### 输出字段：
- approved_amount：不超过申请金额和产品最大额度
- approved_term_months：贷款期限（月）
- suggested_interest_rate：如"4.5%"

## 农机租赁申请分析框架

### 分析要素：
1. **承租方分析**：农业经验、信用记录、租赁历史、支付能力
2. **出租方分析**：资质认证、信用评级、设备维护记录、服务质量
3. **设备分析**：农机类型、状况、保险、市场价值
4. **租赁条件**：租期合理性、使用目的、季节性需求

### 决策规则：
- **自动通过(AUTO_APPROVE)**：双方信用良好，设备状况优秀，风险评分<0.4
- **调整押金(REQUIRE_DEPOSIT_ADJUSTMENT)**：有轻微风险，建议调整押金或条件
- **人工审核(REQUIRE_HUMAN_REVIEW)**：风险评分0.4-0.7，需要人工判断
- **自动拒绝(AUTO_REJECT)**：高风险情况，风险评分>0.7

### 输出字段：
- approved_amount：建议租金（可能调整原租金）
- suggested_deposit：建议押金金额
- 不需要：approved_term_months, suggested_interest_rate

## 通用要求：
1. application_type字段必须准确识别：LOAN_APPLICATION 或 MACHINERY_LEASING
2. risk_score为0-1之间的小数，confidence_score为0-1之间的小数
3. detailed_analysis中的字段根据申请类型灵活调整含义
4. 所有数组字段至少包含一个元素
5. 决策逻辑必须符合上述规则
6. 根据申请类型选择合适的decision枚举值

现在请分析以下申请：
```

- **用户提示词**：

```
## 申请信息
{{#获取申请信息.text}}

## 外部数据
{{#获取外部数据.text}}

## AI模型配置  
{{#获取AI模型配置.text}}

请根据上述信息进行全面的风险评估和决策分析。首先识别申请类型，然后使用对应的分析框架进行评估。
```

#### 节点5：格式化输出（统一多类型版）
- **节点类型**：代码执行
- **编程语言**：Python3
- **输入变量**：
  - structured_output (Object): `{{#LLM统一智能分析.structured_output}}`

- **代码内容（统一多类型版）**：

```python
import json
from datetime import datetime

def main(structured_output: dict) -> dict:
    """
    统一多类型申请处理器 - 支持贷款申请和农机租赁申请
    根据申请类型自动调整输出格式和参数
    """
    
    print(f"[DEBUG] 接收到结构化输出: {type(structured_output)}")
    print(f"[DEBUG] 包含字段: {list(structured_output.keys()) if isinstance(structured_output, dict) else 'Not a dict'}")
    
    try:
        # 验证输入数据
        if not isinstance(structured_output, dict):
            raise ValueError(f"输入不是字典类型，而是: {type(structured_output)}")
        
        # 识别申请类型
        application_type = structured_output.get("application_type", "LOAN_APPLICATION")
        print(f"[DEBUG] 识别申请类型: {application_type}")
        
        # 填充默认值（防御性编程）
        data = fill_default_values(structured_output, application_type)
        
        # 验证和清理数据
        cleaned_data = validate_and_clean_data(data, application_type)
        
        # 创建API响应格式
        result = create_api_response(cleaned_data, application_type)
        
        print(f"[SUCCESS] 处理完成，申请类型: {application_type}，决策: {cleaned_data.get('decision')}")
        return result
        
    except Exception as e:
        print(f"[ERROR] 处理异常: {str(e)}")
        return create_fallback_response(str(e))

def fill_default_values(data: dict, application_type: str) -> dict:
    """根据申请类型填充缺失的默认值"""
    
    # 通用默认值
    defaults = {
        "application_type": application_type,
        "analysis_summary": "AI风险分析",
        "risk_score": 0.5,
        "risk_level": "MEDIUM",
        "confidence_score": 0.5,
        "approved_amount": 0,
        "detailed_analysis": {
            "primary_analysis": "主要分析",
            "secondary_analysis": "次要分析",
            "asset_analysis": "资产分析",
            "risk_factors": ["待评估"],
            "strengths": ["待评估"]
        },
        "recommendations": ["建议审核"],
        "conditions": ["需要审核"]
    }
    
    # 根据申请类型设置特定默认值
    if application_type == "LOAN_APPLICATION":
        defaults.update({
            "decision": "REQUIRE_HUMAN_REVIEW",
            "approved_term_months": 12,
            "suggested_interest_rate": "5.0%",
            "suggested_deposit": 0  # 贷款申请不需要押金
        })
    else:  # MACHINERY_LEASING
        defaults.update({
            "decision": "REQUIRE_HUMAN_REVIEW",
            "approved_term_months": 0,  # 农机租赁不需要期限
            "suggested_interest_rate": "0%",  # 农机租赁不需要利率
            "suggested_deposit": 1000
        })
    
    # 创建新的数据字典，保留原有数据，补充缺失项
    result = defaults.copy()
    result.update(data)
    
    # 特殊处理嵌套的detailed_analysis
    if "detailed_analysis" in data and isinstance(data["detailed_analysis"], dict):
        result["detailed_analysis"].update(data["detailed_analysis"])
    
    return result

def validate_and_clean_data(data: dict, application_type: str) -> dict:
    """根据申请类型验证和清理数据"""
    
    # 数值验证和修正
    try:
        data["risk_score"] = max(0.0, min(1.0, float(data["risk_score"])))
        data["confidence_score"] = max(0.0, min(1.0, float(data["confidence_score"])))
        data["approved_amount"] = max(0.0, float(data["approved_amount"]))
        data["approved_term_months"] = max(0, int(data.get("approved_term_months", 0)))
        data["suggested_deposit"] = max(0.0, float(data.get("suggested_deposit", 0)))
    except (ValueError, TypeError) as e:
        print(f"[WARNING] 数值修正: {e}")
        data["risk_score"] = 0.5
        data["confidence_score"] = 0.5
        data["approved_amount"] = 0.0
    
    # 枚举值验证
    if data.get("risk_level") not in ["LOW", "MEDIUM", "HIGH"]:
        data["risk_level"] = "MEDIUM"
        print("[WARNING] risk_level修正为MEDIUM")
    
    # 根据申请类型验证决策枚举值
    if application_type == "LOAN_APPLICATION":
        valid_decisions = ["AUTO_APPROVED", "AUTO_REJECTED", "REQUIRE_HUMAN_REVIEW"]
        if data.get("decision") not in valid_decisions:
            data["decision"] = "REQUIRE_HUMAN_REVIEW"
            print(f"[WARNING] 贷款申请decision修正为REQUIRE_HUMAN_REVIEW")
    else:  # MACHINERY_LEASING
        valid_decisions = ["AUTO_APPROVE", "AUTO_REJECT", "REQUIRE_HUMAN_REVIEW", "REQUIRE_DEPOSIT_ADJUSTMENT"]
        if data.get("decision") not in valid_decisions:
            data["decision"] = "REQUIRE_HUMAN_REVIEW"
            print(f"[WARNING] 农机租赁decision修正为REQUIRE_HUMAN_REVIEW")
    
    # 数组验证
    for field in ["recommendations", "conditions"]:
        if not isinstance(data.get(field), list) or len(data.get(field, [])) == 0:
            data[field] = ["需要进一步评估"]
    
    # detailed_analysis验证
    if not isinstance(data.get("detailed_analysis"), dict):
        data["detailed_analysis"] = {
            "primary_analysis": "需要重新评估",
            "secondary_analysis": "需要重新评估",
            "asset_analysis": "需要重新评估",
            "risk_factors": ["数据不完整"],
            "strengths": ["待评估"]
        }
    else:
        # 验证嵌套数组
        for field in ["risk_factors", "strengths"]:
            if not isinstance(data["detailed_analysis"].get(field), list):
                data["detailed_analysis"][field] = ["需要评估"]
    
    return data

def create_api_response(data: dict, application_type: str) -> dict:
    """根据申请类型创建API响应格式"""
    
    # 构建通用参数
    common_params = {
        "risk_score": float(data["risk_score"]),
        "confidence_score": float(data["confidence_score"]),
        "risk_level": data["risk_level"],
        "analysis_summary": data["analysis_summary"],
        "detailed_analysis": json.dumps(data["detailed_analysis"], ensure_ascii=False),
        "recommendations": ",".join(data["recommendations"]),
        "conditions": ",".join(data["conditions"]),
        "ai_model_version": f"LLM-v4.0-unified-{application_type.lower()}",
        "workflow_id": "dify-unified-structured-output"
    }
    
    # 根据申请类型添加特定参数
    if application_type == "LOAN_APPLICATION":
        specific_params = {
            "decision": data["decision"],
            "approved_amount": float(data["approved_amount"]),
            "approved_term_months": int(data["approved_term_months"]),
            "suggested_interest_rate": data.get("suggested_interest_rate", "5.0%")
        }
    else:  # MACHINERY_LEASING
        specific_params = {
            "decision": data["decision"],
            "suggested_deposit": float(data.get("suggested_deposit", 0))
        }
    
    # 合并参数
    all_params = {**common_params, **specific_params}
    
    return {
        "success": 1,
        "application_type": application_type,
        "decision": str(data["decision"]),
        "risk_score": float(data["risk_score"]),
        "risk_level": str(data["risk_level"]),
        "confidence_score": float(data["confidence_score"]),
        "approved_amount": float(data.get("approved_amount", 0)),
        "approved_term_months": int(data.get("approved_term_months", 0)),
        "suggested_interest_rate": str(data.get("suggested_interest_rate", "5.0%")),
        "suggested_deposit": float(data.get("suggested_deposit", 0)),
        "analysis_summary": str(data["analysis_summary"]),
        "detailed_analysis": json.dumps(data["detailed_analysis"], ensure_ascii=False),
        "recommendations": ",".join(data["recommendations"]),
        "conditions": ",".join(data["conditions"]),
        "api_params": json.dumps(all_params, ensure_ascii=False),
        "error": ""
    }

def create_fallback_response(error_msg: str) -> dict:
    """创建降级响应"""
    
    fallback_params = {
        "decision": "REQUIRE_HUMAN_REVIEW",
        "risk_score": 0.6,
        "confidence_score": 0.1,
        "risk_level": "MEDIUM",
        "analysis_summary": f"系统处理异常: {error_msg}，建议人工审核",
        "detailed_analysis": json.dumps({
            "primary_analysis": "系统异常，无法完成分析",
            "secondary_analysis": "系统异常，无法完成分析",
            "asset_analysis": "系统异常，无法完成分析",
            "risk_factors": ["系统处理异常"],
            "strengths": ["需要人工评估"]
        }, ensure_ascii=False),
        "recommendations": "转人工审核,检查系统配置",
        "conditions": "系统异常，需要人工处理",
        "ai_model_version": "LLM-v4.0-fallback",
        "workflow_id": "dify-unified-error-handler"
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

## 第三步：测试与验证（v4.0版）

### 3.1 单节点测试

1. **测试工具连接**
   ```bash
   # 测试统一接口连通性
   curl -H "Authorization: AI-Agent-Token your_token" \
        http://172.18.120.10:8080/api/v1/ai-agent/applications/test_app_001/info
   
   curl -H "Authorization: AI-Agent-Token your_token" \
        http://172.18.120.10:8080/api/v1/ai-agent/external-data/user_001?data_types=credit_report
   
   curl -H "Authorization: AI-Agent-Token your_token" \
        http://172.18.120.10:8080/api/v1/ai-agent/config/models
   ```

2. **测试LLM节点**
   - 使用不同申请类型的模拟数据
   - 验证申请类型识别准确性
   - 检查输出格式完整性

### 3.2 端到端测试

**测试数据集（v4.0多类型）**：

**贷款申请测试**：
```json
{
  "application_id": "test_app_001",
  "callback_url": "http://172.18.120.10:8080/callback"
}
```

**农机租赁申请测试**：
```json
{
  "application_id": "ml_test_001", 
  "callback_url": "http://172.18.120.10:8080/callback"
}
```

**预期执行流程（v4.0版）**：

1. ✅ **获取申请信息** → 返回带application_type标识的完整申请数据
2. ✅ **获取外部数据** → 根据用户类型和申请类型返回相关征信数据  
3. ✅ **获取AI配置** → 返回多种业务规则和不同申请类型的阈值配置
4. ✅ **LLM统一智能分析** → 
   - 准确识别申请类型（LOAN_APPLICATION/MACHINERY_LEASING）
   - 应用对应的分析框架和决策规则
   - 输出符合JSON Schema的结构化数据
5. ✅ **格式化输出与验证** → 
   - 根据申请类型动态调整输出格式
   - 验证业务逻辑一致性（风险分数与决策的匹配）
   - 数据完整性检查和错误处理
6. ✅ **提交AI决策** → 
   - 自动路由到对应的业务处理逻辑
   - 记录AI操作日志
   - 更新申请状态
7. ✅ **操作日志记录** → 查询并记录本次AI操作的详细信息
8. ✅ **工作流完成** → 返回包含申请类型和决策结果的完整响应

### 3.3 性能验证（v4.0标准）

**关键指标验证**：
- ⏱️ **处理时间**：单个申请端到端处理应在8秒内完成
- 🎯 **准确率**：
  - 申请类型识别准确率 ≥ 99%
  - 风险评分与决策一致性 ≥ 95%
  - 数据格式验证通过率 ≥ 99%
- 🔄 **并发性**：支持同时处理10个不同类型的申请
- 💾 **数据完整性**：
  - 所有AI操作日志完整记录
  - 审计追踪链路完整
  - 敏感数据自动脱敏

**质量验证清单**：

| 验证项目 | 贷款申请 | 农机租赁 | 验证标准 |
|---------|---------|---------|---------|
| 申请类型识别 | ✅ | ✅ | 100%准确 |
| 风险分数合理性 | ✅ | ✅ | 0-1范围，精确到3位小数 |
| 决策逻辑一致性 | ✅ | ✅ | 风险分数与决策匹配 |
| 输出格式完整性 | ✅ | ✅ | 所有必需字段不为空 |
| 业务规则遵循 | ✅ | ✅ | 符合各自业务阈值 |
| 错误处理能力 | ✅ | ✅ | 异常情况下有降级响应 |

## 第四步：高级配置与优化（v4.0版）

### 4.1 多类型差异化配置策略

**按申请类型优化模型配置**：

1. **贷款申请场景**
   - 推荐模型：`Claude-3.5-sonnet`（金融风险分析专业性强）
   - 备选模型：`GPT-4o`（复杂推理能力强）
   - Temperature：`0.05`（确保决策一致性和准确性）
   - Max Tokens：`2000`（详细分析需要更多输出）

2. **农机租赁场景**
   - 推荐模型：`GPT-4o`（农业场景理解和设备评估能力好）
   - 备选模型：`Claude-3.5-sonnet`（结构化输出稳定）
   - Temperature：`0.1`（允许适当的灵活性判断）
   - Max Tokens：`1800`（租赁分析相对简洁）

### 4.2 提示词优化策略（v4.0版）

**分层次提示词架构**：

```
# 层次1：统一角色定位和版本标识
你是慧农金融的统一AI智能审批专家（v4.0版）...

# 层次2：动态类型识别和路由
if application_id.startswith("ml_") or "lessee_info" in data:
    -> 农机租赁申请处理流程
elif application_id.startswith("test_app_") or "applicant_info" in data:
    -> 贷款申请处理流程

# 层次3：专业分析框架应用
贷款申请 -> 信用+财务+农业资产分析
农机租赁 -> 双方信用+设备状况+租赁合理性分析

# 层次4：一致性决策输出
统一的JSON Schema结构化输出，确保格式一致性
```

**提示词质量控制**：
- **A/B测试**：对比不同提示词版本的决策准确性
- **持续优化**：基于人工审核反馈调整提示词逻辑
- **版本管理**：维护提示词的版本历史和变更记录

### 4.3 错误处理与回退机制（v4.0增强版）

**多级智能回退策略**：

1. **Level 1 - 数据修正**：LLM输出格式问题 → 自动补全缺失字段，修正数据类型
2. **Level 2 - 逻辑修正**：业务逻辑不一致 → 根据规则自动调整决策和风险等级
3. **Level 3 - 类型识别失败**：无法识别申请类型 → 默认为贷款申请，转人工审核
4. **Level 4 - 系统异常**：API调用失败 → 启用降级模式，记录详细错误日志
5. **Level 5 - 完全失败**：所有策略失败 → 强制转人工审核，触发告警通知

### 4.4 监控与日志配置（v4.0版）

**实时监控指标体系**：

| 类别 | 指标名称 | 预警阈值 | 处理动作 |
|------|----------|----------|----------|
| **性能指标** | 平均处理时间 | >12秒 | 优化模型参数/检查网络 |
| | LLM响应时间 | >8秒 | 切换备用模型 |
| | API调用成功率 | <98% | 检查服务状态 |
| **准确性指标** | 类型识别错误率 | >2% | 优化识别逻辑 |
| | 风险评分偏差 | >±0.2 | 校准评分模型 |
| | 决策逻辑一致性 | <95% | 调整业务规则 |
| **业务指标** | 自动通过率（贷款） | <25%或>75% | 调整决策阈值 |
| | 自动通过率（租赁） | <30%或>80% | 调整风险参数 |
| | 人工审核比例 | >60% | 优化自动化程度 |
| **系统指标** | 错误率 | >0.5% | 检查系统稳定性 |
| | 并发处理能力 | <10 QPS | 扩容或优化 |

**结构化日志配置（v4.0版）**：

```json
{
  "unified_ai_workflow_logs": {
    "version": "v4.0",
    "level": "INFO",
    "format": "[{timestamp}] [{application_type}] [{workflow_id}] [{node_id}] {message}",
    "required_fields": [
      "application_id",
      "application_type", 
      "llm_model",
      "processing_time_ms",
      "decision",
      "risk_score",
      "confidence_score",
      "workflow_version",
      "node_execution_status"
    ],
    "business_fields": [
      "approved_amount",
      "suggested_deposit",
      "risk_level",
      "human_review_required"
    ],
    "audit_fields": [
      "ai_operation_id",
      "decision_id", 
      "data_sources",
      "model_version"
    ]
  }
}
```

## 第五步：集成与部署（v4.0版）

### 5.1 生产环境配置清单（v4.0版）

**必要配置检查**：

- [ ] **后端服务**：
  - [ ] 健康检查接口正常 (`/livez`, `/readyz`)
  - [ ] AI Agent Token配置正确且有足够权限
  - [ ] 数据库连接池配置合理（支持并发处理）
  - [ ] Redis缓存配置（缓存AI模型配置和外部数据）

- [ ] **Dify平台**：
  - [ ] 模型配额充足（GPT-4o和Claude-3.5-sonnet）
  - [ ] 工作流版本管理和备份
  - [ ] API调用频率限制配置
  - [ ] 错误重试策略设置

- [ ] **监控告警**：
  - [ ] 关键指标监控规则设置
  - [ ] 告警通知渠道配置
  - [ ] 业务异常自动处理规则

- [ ] **数据安全**：
  - [ ] 敏感数据脱敏配置
  - [ ] 日志归档和加密策略
  - [ ] 访问权限控制和审计

### 5.2 业务集成要点（v4.0版）

**前端集成（多类型支持）**：
```javascript
// 统一调用接口，支持多种申请类型
const callUnifiedAIWorkflow = async (applicationId, applicationType) => {
  try {
    const response = await fetch('/api/ai-workflow/execute', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + userToken
      },
      body: JSON.stringify({
        application_id: applicationId,
        workflow_type: 'UNIFIED_LLM_BASED_V4',
        expected_application_type: applicationType // 可选，用于验证
      })
    });

    const result = await response.json();
    
    // 根据申请类型显示不同的结果
    if (result.application_type === 'LOAN_APPLICATION') {
      displayLoanResult(result);
    } else if (result.application_type === 'MACHINERY_LEASING') {
      displayLeasingResult(result);
    }
    
    console.log('AI决策:', result.decision);
    console.log('风险评分:', result.risk_score);
    console.log('处理时间:', result.processing_time_ms + 'ms');
    
  } catch (error) {
    console.error('AI工作流调用失败:', error);
    // 降级到人工审核
    fallbackToManualReview(applicationId);
  }
};

// 使用示例
callUnifiedAIWorkflow('test_app_001', 'LOAN_APPLICATION');
callUnifiedAIWorkflow('ml_test_001', 'MACHINERY_LEASING');
```

**后端集成（统一处理器）**：
```go
// 统一处理不同类型申请的AI决策回调
func HandleUnifiedAIDecisionCallback(c *gin.Context) {
    var callback UnifiedAIDecisionCallback
    if err := c.ShouldBindJSON(&callback); err != nil {
        c.JSON(400, gin.H{"error": "参数解析失败", "details": err.Error()})
        return
    }
    
    // 记录请求日志
    log.Info("收到AI决策回调", 
        "application_id", callback.ApplicationID,
        "application_type", callback.ApplicationType,
        "decision", callback.Decision,
        "risk_score", callback.RiskScore)
    
    // 根据申请类型路由到不同处理器
    switch callback.ApplicationType {
    case "LOAN_APPLICATION":
        err = processingLoanDecisionV4(callback)
    case "MACHINERY_LEASING":
        err = processingMachineryLeasingDecisionV4(callback)
    default:
        log.Warn("未知申请类型", "type", callback.ApplicationType)
        err = processUnknownTypeDecisionV4(callback)
    }
    
    if err != nil {
        log.Error("AI决策处理失败", "error", err)
        c.JSON(500, gin.H{"error": "处理失败", "details": err.Error()})
        return
    }
    
    // 记录成功日志
    recordAIDecisionLog(callback)
    
    c.JSON(200, gin.H{
        "code": 0,
        "message": "AI决策处理成功",
        "data": gin.H{
            "application_id": callback.ApplicationID,
            "processed_at": time.Now().Format(time.RFC3339),
        },
    })
}

// 贷款申请决策处理
func processingLoanDecisionV4(callback UnifiedAIDecisionCallback) error {
    // 更新申请状态
    err := updateLoanApplicationStatus(
        callback.ApplicationID, 
        callback.Decision,
        callback.ApprovedAmount,
        callback.ApprovedTermMonths,
        callback.SuggestedInterestRate,
    )
    if err != nil {
        return fmt.Errorf("更新贷款申请状态失败: %w", err)
    }
    
    // 如果自动批准，启动合同生成流程
    if callback.Decision == "AUTO_APPROVED" {
        go generateLoanContract(callback.ApplicationID)
    }
    
    // 发送通知
    go sendLoanDecisionNotification(callback)
    
    return nil
}

// 农机租赁申请决策处理
func processingMachineryLeasingDecisionV4(callback UnifiedAIDecisionCallback) error {
    // 更新租赁申请状态
    err := updateMachineryLeasingStatus(
        callback.ApplicationID,
        callback.Decision,
        callback.SuggestedDeposit,
    )
    if err != nil {
        return fmt.Errorf("更新农机租赁状态失败: %w", err)
    }
    
    // 如果自动批准，启动租赁合同流程
    if callback.Decision == "AUTO_APPROVE" {
        go generateLeasingContract(callback.ApplicationID)
    }
    
    // 发送通知
    go sendLeasingDecisionNotification(callback)
    
    return nil
}
```

### 5.3 部署架构建议（v4.0版）

**推荐部署架构**：

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   前端应用      │    │   Dify平台      │    │   后端服务      │
│  (React/Vue)    │    │  (AI工作流)     │    │  (Go/Java)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │ HTTPS请求             │ AI工作流调用          │ 数据处理
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   负载均衡器    │    │   模型服务集群  │    │   数据库集群    │
│  (Nginx/ALB)    │    │  (GPT-4/Claude) │    │ (MySQL/Redis)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   监控告警      │    │   日志系统      │    │   文件存储      │
│ (Prometheus)    │    │ (ELK/Fluentd)   │    │  (OSS/S3)       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

**高可用配置**：
- **服务冗余**：每个组件至少2个实例
- **数据备份**：数据库主从复制，定期备份
- **故障切换**：自动故障检测和切换机制
- **降级策略**：AI服务异常时自动转人工审核

### 5.4 运维最佳实践（v4.0版）

**日常运维检查清单**：

1. **每日检查**：
   - [ ] AI工作流成功率 ≥ 95%
   - [ ] 平均处理时间 ≤ 8秒
   - [ ] 错误日志数量 ≤ 10条/小时
   - [ ] 模型配额使用情况 ≤ 80%

2. **每周检查**：
   - [ ] 申请类型识别准确率统计
   - [ ] 不同申请类型的自动通过率分析
   - [ ] AI决策与人工审核结果对比分析
   - [ ] 系统性能趋势分析

3. **每月检查**：
   - [ ] 提示词效果评估和优化
   - [ ] 业务规则阈值调整
   - [ ] 成本效益分析
   - [ ] 安全审计和合规检查

**故障处理手册**：

| 故障类型 | 症状 | 处理步骤 | 预计恢复时间 |
|---------|------|---------|-------------|
| LLM响应超时 | 处理时间>15秒 | 1.切换备用模型 2.检查网络 3.联系模型服务商 | 5分钟 |
| 申请类型识别错误 | 类型识别率<90% | 1.检查提示词 2.验证测试数据 3.回滚到上个版本 | 10分钟 |
| API调用失败 | 500错误>5% | 1.检查服务状态 2.查看错误日志 3.重启相关服务 | 15分钟 |
| 数据格式异常 | 输出格式错误>10% | 1.验证JSON Schema 2.检查LLM模型状态 3.启用降级模式 | 10分钟 |

## 总结与最佳实践（v4.0版）

### ✅ 成功配置标准

1. **功能完整性**
   - ✅ 支持贷款申请和农机租赁申请的统一处理（类型自动识别准确率>99%）
   - ✅ 智能业务规则应用（不同类型使用专门的分析框架和决策阈值）
   - ✅ 完整的错误处理和降级机制（异常情况下自动转人工审核）
   - ✅ 全链路日志记录和审计追踪（每个决策都可追溯）

2. **技术可靠性**
   - ✅ 多级错误处理（数据修正→逻辑修正→类型降级→系统降级→强制人工）
   - ✅ 业务逻辑一致性验证（风险分数与决策自动匹配）
   - ✅ 高并发支持（>10 QPS）和性能优化（<8秒处理时间）
   - ✅ 结构化输出保证数据格式一致性

3. **业务适用性**
   - ✅ **贷款申请**：信用分析、财务评估、利率建议、期限确定
   - ✅ **农机租赁**：双方评估、设备分析、押金建议、租期合理性
   - ✅ 决策可解释（详细分析报告）和人工审核支持
   - ✅ 符合金融行业合规要求（数据脱敏、审计追踪）

### 🚀 推荐最佳实践

1. **模型选择与优化**
   - **贷款场景**：优先Claude-3.5-sonnet（金融分析专业性）+ 极低Temperature（0.05）
   - **农机租赁场景**：优先GPT-4o（农业场景理解）+ 低Temperature（0.1）
   - **备选策略**：配置备用模型，自动故障切换

2. **提示词工程最佳实践**
   - **版本化管理**：维护提示词版本历史，支持快速回滚
   - **分层架构**：角色定位→类型识别→专业分析→统一输出
   - **持续优化**：基于业务反馈和A/B测试结果不断改进

3. **监控运维策略**
   - **实时监控**：关键业务指标的实时监控和告警
   - **定期评估**：周度业务指标分析，月度效果评估
   - **预防性维护**：主动识别潜在问题，预防性优化

4. **扩展性设计**
   - **新类型扩展**：预留新申请类型的扩展接口和配置
   - **多租户支持**：支持不同业务线的独立配置
   - **国际化准备**：支持多语言和不同地区的业务规则

### 🔧 常见问题解决（v4.0版）

**Q1: LLM识别申请类型不准确怎么办？**
A1: 
- 检查申请ID格式规范是否正确
- 在提示词中加强类型识别的判断逻辑
- 增加数据结构特征的识别权重
- 考虑在开始节点添加类型预判断

**Q2: 不同申请类型的决策阈值如何调整？**
A2: 
- 通过AI模型配置接口动态调整阈值
- 基于历史数据分析设置合理的通过率目标
- 定期评估阈值效果，进行数据驱动的优化
- 考虑季节性因素和市场变化

**Q3: 如何处理新的申请类型扩展？**
A3: 
- 在后端添加新的申请类型识别逻辑
- 更新LLM提示词，增加新类型的分析框架
- 扩展JSON Schema定义，支持新类型的字段
- 增加对应的业务规则配置

**Q4: AI决策质量如何持续提升？**
A4: 
- 建立AI决策与人工审核结果的对比分析机制
- 收集业务专家的反馈，优化分析逻辑
- 定期更新外部数据源，提高分析准确性
- 使用强化学习方法持续优化决策模型

**Q5: 如何确保系统的安全性和合规性？**
A5: 
- 实施完善的数据脱敏策略
- 建立完整的操作审计日志
- 定期进行安全评估和渗透测试
- 遵循金融行业的合规要求和标准

---

**🎉 恭喜！通过以上配置，您将拥有一个功能完整、技术先进、业务适用的v4.0统一多类型AI智能审批系统！**

**系统特色**：
- 🎯 **高精度**：申请类型识别准确率>99%，决策逻辑一致性>95%
- ⚡ **高性能**：单申请处理时间<8秒，支持>10 QPS并发处理
- 🛡️ **高可靠**：5级错误处理机制，异常情况下自动降级
- 📊 **可观测**：全链路监控，完整的业务指标和技术指标
- 🔧 **易维护**：版本化管理，支持热更新和快速回滚
- 🚀 **可扩展**：预留扩展接口，支持新业务类型快速接入