# Dify LLM工具配置与智能工作流集成文档

## 概述

本文档详细介绍如何在Dify平台中配置基于LLM的AI智能体接口工具，并创建智能化的AI审批工作流。通过LLM的强大推理能力，简化传统的代码逻辑，实现更智能、更灵活的贷款审批系统。

## 1. LLM版OpenAPI Schema配置

### 1.1 完整的OpenAPI Schema（支持LLM）

相比传统版本，LLM版本的OpenAPI Schema更注重数据结构的清晰性和LLM的理解能力：

```json
{
  "openapi": "3.1.0",
  "info": {
    "title": "慧农金融AI智能体接口（LLM增强版）",
    "description": "为Dify LLM工作流优化的AI智能体审批接口，提供结构化数据支持大语言模型分析",
    "version": "2.1.0",
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
        "summary": "获取申请详细信息（LLM优化）",
        "description": "获取贷款申请的结构化详细信息，格式优化用于LLM理解和分析。包含申请人基本信息、财务状况、产品信息等全量数据",
        "operationId": "getApplicationInfo",
        "tags": ["LLM数据源"],
        "parameters": [
          {
            "name": "application_id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string",
              "pattern": "^[a-zA-Z0-9_-]+$"
            },
            "description": "申请ID，支持格式：app_20240301_001、test_app_001等",
            "example": "app_20240301_001"
          }
        ],
        "responses": {
          "200": {
            "description": "成功获取申请信息，数据结构化便于LLM分析",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ApplicationInfoLLMResponse"
                },
                "examples": {
                  "standard_application": {
                    "summary": "标准农业贷款申请",
                    "value": {
                      "code": 0,
                      "message": "获取成功",
                      "data": {
                        "application_id": "app_20240301_001",
                        "product_info": {
                          "product_id": "AGRI_LOAN_001",
                          "name": "农业种植贷",
                          "category": "种植贷",
                          "max_amount": 100000,
                          "interest_rate_yearly": "4.5%"
                        },
                        "application_info": {
                          "amount": 50000,
                          "term_months": 12,
                          "purpose": "购买种子和化肥",
                          "submitted_at": "2024-03-01T10:00:00Z",
                          "status": "SUBMITTED"
                        },
                        "applicant_info": {
                          "user_id": "user_12345",
                          "real_name": "张农民",
                          "age": 35,
                          "is_verified": true
                        },
                        "financial_info": {
                          "annual_income": 80000,
                          "credit_score": 750,
                          "existing_loans": 0,
                          "land_area": "10亩",
                          "farming_experience": "10年"
                        }
                      }
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
    "/api/v1/ai-agent/external-data": {
      "get": {
        "summary": "获取外部风险数据（LLM格式化）",
        "description": "获取征信、银行流水、黑名单等外部数据，数据格式针对LLM分析优化，包含详细的风险指标说明",
        "operationId": "getExternalData",
        "tags": ["LLM数据源"],
        "parameters": [
          {
            "name": "user_id",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "用户ID",
            "example": "user_12345"
          },
          {
            "name": "data_types",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "enum": [
                "credit_report",
                "bank_flow", 
                "blacklist_check",
                "government_subsidy",
                "credit_report,bank_flow,blacklist_check",
                "all"
              ]
            },
            "description": "数据类型，支持单个或组合查询",
            "example": "credit_report,bank_flow,blacklist_check"
          }
        ],
        "responses": {
          "200": {
            "description": "成功获取外部数据，格式化用于LLM风险分析",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ExternalDataLLMResponse"
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
        "summary": "获取AI模型配置（LLM决策参考）",
        "description": "获取业务规则、风险阈值等配置信息，为LLM提供决策依据和参考标准",
        "operationId": "getAIModelConfig",
        "tags": ["LLM配置"],
        "responses": {
          "200": {
            "description": "获取AI模型配置成功，包含LLM决策规则",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/AIModelConfigLLMResponse"
                }
              }
            }
          }
        },
        "security": [{"AIAgentToken": []}]
      }
    },
    "/api/v1/ai-agent/applications/{application_id}/ai-decision": {
      "post": {
        "summary": "提交LLM决策结果",
        "description": "接收经过LLM分析后的决策结果，包含详细的推理过程和建议",
        "operationId": "submitAIDecision",
        "tags": ["LLM决策提交"],
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
                "$ref": "#/components/schemas/LLMAIDecisionRequest"
              },
              "examples": {
                "auto_approved": {
                  "summary": "自动批准案例",
                  "value": {
                    "ai_analysis": {
                      "risk_level": "LOW",
                      "risk_score": 0.25,
                      "confidence_score": 0.88,
                      "analysis_summary": "申请人信用状况良好，财务状况稳定，农业经验丰富，风险较低。",
                      "detailed_analysis": {
                        "credit_analysis": "信用分数750分，属于优秀等级",
                        "financial_analysis": "年收入8万，申请金额5万，收入覆盖率良好",
                        "risk_factors": [],
                        "strengths": ["信用记录良好", "农业经验丰富", "收入稳定"]
                      },
                      "recommendations": ["建议按申请金额批准", "利率可适当优惠"]
                    },
                    "ai_decision": {
                      "decision": "AUTO_APPROVED",
                      "approved_amount": 50000,
                      "approved_term_months": 12,
                      "suggested_interest_rate": "4.5%",
                      "conditions": ["按时还款", "保持土地使用权"]
                    }
                  }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "LLM决策处理成功",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/DecisionResponse"
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
      "ApplicationInfoLLMResponse": {
        "type": "object",
        "description": "为LLM优化的申请信息响应格式",
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
            "type": "object",
            "properties": {
              "application_id": {
                "type": "string",
                "description": "申请唯一标识"
              },
              "product_info": {
                "type": "object",
                "description": "贷款产品信息",
                "properties": {
                  "product_id": {"type": "string"},
                  "name": {"type": "string", "description": "产品名称"},
                  "category": {"type": "string", "description": "产品类别"},
                  "max_amount": {"type": "number", "description": "最大贷款额度"},
                  "interest_rate_yearly": {"type": "string", "description": "年利率"}
                }
              },
              "application_info": {
                "type": "object",
                "description": "申请基本信息",
                "properties": {
                  "amount": {"type": "number", "description": "申请金额"},
                  "term_months": {"type": "integer", "description": "申请期限（月）"},
                  "purpose": {"type": "string", "description": "贷款用途"},
                  "submitted_at": {"type": "string", "format": "date-time"},
                  "status": {"type": "string", "description": "当前状态"}
                }
              },
              "applicant_info": {
                "type": "object",
                "description": "申请人信息",
                "properties": {
                  "user_id": {"type": "string"},
                  "real_name": {"type": "string", "description": "真实姓名"},
                  "age": {"type": "integer", "description": "年龄"},
                  "is_verified": {"type": "boolean", "description": "身份是否验证"}
                }
              },
              "financial_info": {
                "type": "object",
                "description": "财务信息",
                "properties": {
                  "annual_income": {"type": "number", "description": "年收入（元）"},
                  "credit_score": {"type": "integer", "description": "信用分数（300-850）"},
                  "existing_loans": {"type": "integer", "description": "现有贷款数量"},
                  "land_area": {"type": "string", "description": "土地面积"},
                  "farming_experience": {"type": "string", "description": "农业经验"}
                }
              }
            }
          }
        }
      },
      "ExternalDataLLMResponse": {
        "type": "object",
        "description": "为LLM优化的外部数据响应",
        "properties": {
          "code": {"type": "integer"},
          "message": {"type": "string"},
          "data": {
            "type": "object",
            "properties": {
              "user_id": {"type": "string"},
              "credit_report": {
                "type": "object",
                "description": "征信报告详情",
                "properties": {
                  "score": {"type": "integer", "description": "征信分数"},
                  "grade": {"type": "string", "description": "信用等级"},
                  "report_date": {"type": "string"},
                  "overdue_records": {"type": "integer", "description": "逾期记录数"},
                  "loan_history": {"type": "array", "description": "历史贷款记录"}
                }
              },
              "bank_flow": {
                "type": "object",
                "description": "银行流水分析",
                "properties": {
                  "average_monthly_income": {"type": "number", "description": "月均收入"},
                  "account_stability": {"type": "string", "description": "账户稳定性"},
                  "debt_to_income_ratio": {"type": "number", "description": "债务收入比"},
                  "last_6_months_flow": {
                    "type": "array",
                    "description": "近6个月流水",
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
                "description": "黑名单检查",
                "properties": {
                  "is_blacklisted": {"type": "boolean"},
                  "risk_level": {"type": "string", "description": "风险等级"},
                  "check_time": {"type": "string"}
                }
              },
              "government_subsidy": {
                "type": "object",
                "description": "政府补贴信息",
                "properties": {
                  "received_subsidies": {
                    "type": "array",
                    "description": "已获得补贴列表",
                    "items": {
                      "type": "object",
                      "properties": {
                        "year": {"type": "integer"},
                        "type": {"type": "string"},
                        "amount": {"type": "number"}
                      }
                    }
                  },
                  "total_amount": {"type": "number", "description": "补贴总额"}
                }
              }
            }
          }
        }
      },
      "AIModelConfigLLMResponse": {
        "type": "object",
        "description": "为LLM提供的决策配置",
        "properties": {
          "code": {"type": "integer"},
          "message": {"type": "string"},
          "data": {
            "type": "object",
            "properties": {
              "approval_rules": {
                "type": "object",
                "description": "审批规则",
                "properties": {
                  "auto_approval_threshold": {"type": "number", "description": "自动批准阈值"},
                  "auto_rejection_threshold": {"type": "number", "description": "自动拒绝阈值"},
                  "max_auto_approval_amount": {"type": "number", "description": "自动批准最大金额"},
                  "required_human_review_conditions": {
                    "type": "array",
                    "items": {"type": "string"},
                    "description": "需要人工审核的条件"
                  }
                }
              },
              "business_parameters": {
                "type": "object",
                "description": "业务参数",
                "properties": {
                  "max_debt_to_income_ratio": {"type": "number", "description": "最大债务收入比"},
                  "min_credit_score": {"type": "integer", "description": "最低信用分数要求"},
                  "max_loan_amount_by_category": {
                    "type": "object",
                    "description": "按类别的最大贷款额度"
                  }
                }
              },
              "risk_factors_weights": {
                "type": "object",
                "description": "风险因素权重",
                "properties": {
                  "credit_score_weight": {"type": "number"},
                  "income_stability_weight": {"type": "number"},
                  "debt_ratio_weight": {"type": "number"},
                  "blacklist_weight": {"type": "number"}
                }
              }
            }
          }
        }
      },
      "LLMAIDecisionRequest": {
        "type": "object",
        "description": "LLM决策请求格式",
        "required": ["ai_analysis", "ai_decision", "processing_info"],
        "properties": {
          "ai_analysis": {
            "type": "object",
            "description": "LLM分析结果",
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
                "description": "风险分数"
              },
              "confidence_score": {
                "type": "number",
                "minimum": 0,
                "maximum": 1,
                "description": "LLM决策置信度"
              },
              "analysis_summary": {
                "type": "string",
                "description": "LLM生成的分析摘要"
              },
              "detailed_analysis": {
                "type": "object",
                "description": "详细分析内容",
                "properties": {
                  "credit_analysis": {"type": "string"},
                  "financial_analysis": {"type": "string"},
                  "risk_factors": {
                    "type": "array",
                    "items": {"type": "string"}
                  },
                  "strengths": {
                    "type": "array", 
                    "items": {"type": "string"}
                  }
                }
              },
              "recommendations": {
                "type": "array",
                "items": {"type": "string"},
                "description": "LLM建议"
              }
            }
          },
          "ai_decision": {
            "type": "object",
            "description": "LLM决策结果",
            "properties": {
              "decision": {
                "type": "string",
                "enum": ["AUTO_APPROVED", "AUTO_REJECTED", "REQUIRE_HUMAN_REVIEW"]
              },
              "approved_amount": {"type": "number"},
              "approved_term_months": {"type": "integer"},
              "suggested_interest_rate": {"type": "string"},
              "conditions": {
                "type": "array",
                "items": {"type": "string"}
              },
              "next_action": {"type": "string"}
            }
          },
          "processing_info": {
            "type": "object",
            "properties": {
              "ai_model_version": {"type": "string", "example": "LLM-v2.0"},
              "processing_time_ms": {"type": "integer"},
              "workflow_id": {"type": "string"},
              "processed_at": {"type": "string", "format": "date-time"}
            }
          }
        }
      },
      "DecisionResponse": {
        "type": "object",
        "properties": {
          "code": {"type": "integer"},
          "message": {"type": "string"},
          "data": {
            "type": "object",
            "properties": {
              "application_id": {"type": "string"},
              "new_status": {"type": "string"},
              "next_step": {"type": "string"}
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

## 2. Dify LLM工作流配置步骤

### 2.1 导入优化的自定义工具

1. **创建LLM专用工具**
   - 工具名称：`慧农金融AI智能体（LLM增强版）`
   - 描述：`专为大语言模型优化的智能审批接口工具`
   - 导入上述OpenAPI Schema

2. **验证工具导入**
   - 检查4个主要接口是否正确导入
   - 测试API连接状态
   - 验证参数格式

### 2.2 LLM智能工作流架构设计

#### 工作流总体架构
```
开始节点
    ↓
[数据收集层]
    ├── 获取申请信息
    ├── 获取外部数据  
    └── 获取AI配置
    ↓
[LLM智能分析层]
    ├── LLM风险分析（核心）
    └── 结果解析与验证
    ↓
[决策执行层]
    ├── 提交AI决策
    └── 结果返回
    ↓
结束节点
```

#### 节点间数据流设计
```
申请信息 + 外部数据 + AI配置 
    ↓ (结构化输入)
LLM智能分析引擎
    ↓ (自然语言推理)
JSON格式化决策结果
    ↓ (API标准格式)
后端系统状态更新
```

### 2.3 详细节点配置

#### 🔧 节点配置模板

**节点1: 获取申请信息**
```yaml
节点类型: 工具调用
工具选择: 慧农金融AI智能体（LLM增强版）
操作: getApplicationInfo
参数配置:
  application_id: "{{start.application_id}}"
输出别名: application_data
错误处理: 启用重试（3次）
```

**节点2: 获取外部数据**
```yaml
节点类型: 工具调用
工具选择: 慧农金融AI智能体（LLM增强版）
操作: getExternalData
参数配置:
  user_id: "{{application_data.text | jq '.data.applicant_info.user_id' | trim}}"
  data_types: "credit_report,bank_flow,blacklist_check,government_subsidy"
输出别名: external_data
依赖节点: 节点1
```

**节点3: 获取AI配置**
```yaml
节点类型: 工具调用
工具选择: 慧农金融AI智能体（LLM增强版）
操作: getAIModelConfig
参数配置: 无
输出别名: ai_config
并行执行: 可与节点2并行
```

**节点4: LLM智能分析（核心）**
```yaml
节点类型: LLM
模型选择: GPT-4-turbo 或 Claude-3.5-Sonnet
温度设置: 0.1 (确保输出一致性)
最大tokens: 2000
系统提示词: [详见下文]
用户提示词: [详见下文]
输出别名: llm_analysis
```

**节点5: 解析LLM输出**
```yaml
节点类型: 代码执行
编程语言: Python3
输入变量:
  - llm_output: "{{llm_analysis.text}}"
代码逻辑: [JSON解析和格式化]
输出别名: formatted_decision
```

**节点6: 提交AI决策**
```yaml
节点类型: 工具调用
工具选择: 慧农金融AI智能体（LLM增强版）
操作: submitAIDecision
参数配置:
  application_id: "{{start.application_id}}"
  Request Body: "{{formatted_decision.api_request}}"
输出别名: submit_result
```

## 3. LLM提示词工程

### 3.1 系统提示词（System Prompt）

```markdown
# 慧农金融AI智能审批专家

## 身份定位
你是慧农金融的资深AI审批专家，专门负责农业贷款的风险评估和审批决策。你具备丰富的金融风险管理经验和对农业行业的深度理解。

## 核心职责
1. **全面风险评估**: 基于申请人信息、财务状况、外部数据进行综合风险分析
2. **智能决策建议**: 根据业务规则和风险阈值给出明确的审批建议
3. **详细分析报告**: 提供清晰的风险分析和决策依据
4. **合规性检查**: 确保所有决策符合监管要求和公司政策

## 分析框架

### 📊 财务健康度评估
- **收入稳定性**: 分析年收入水平、收入来源多样性
- **债务负担**: 计算债务收入比，评估还款能力  
- **资产状况**: 考虑土地资产、农业设备等担保物
- **现金流**: 分析银行流水，关注季节性收入特点

### 🛡️ 信用风险评估
- **信用历史**: 征信分数、历史还款记录
- **违约风险**: 逾期记录、黑名单状态
- **行业风险**: 农业市场波动、自然灾害风险
- **政策风险**: 农业政策变化影响

### 🎯 业务规则检查
- **准入门槛**: 最低信用分数、收入要求
- **额度限制**: 按产品类别的最大贷款额度
- **期限匹配**: 贷款期限与用途的合理性
- **利率定价**: 基于风险等级的利率建议

## 决策规则

### ✅ 自动批准条件（AUTO_APPROVED）
- 信用分数 ≥ 750分
- 债务收入比 ≤ 30%
- 无黑名单记录
- 申请金额 ≤ 年收入的40%
- 综合风险评分 < 0.3
- 有明确的还款来源

### ⚠️ 人工审核条件（REQUIRE_HUMAN_REVIEW）
- 信用分数 600-749分
- 债务收入比 30-50%
- 申请金额较大但在合理范围
- 存在轻微风险因素但可控
- 综合风险评分 0.3-0.7
- 需要额外担保或条件

### ❌ 自动拒绝条件（AUTO_REJECTED）
- 信用分数 < 600分
- 存在黑名单记录
- 债务收入比 > 50%
- 综合风险评分 > 0.7
- 收入与申请金额严重不匹配
- 存在欺诈风险

## 输出规范

### 📝 必须输出格式
你必须严格按照以下JSON格式输出，确保格式正确且包含所有必需字段：

```json
{
  "analysis_summary": "简明扼要的风险分析总结，控制在150字以内",
  "risk_score": 0.25,
  "risk_level": "LOW|MEDIUM|HIGH",
  "confidence_score": 0.88,
  "decision": "AUTO_APPROVED|REQUIRE_HUMAN_REVIEW|AUTO_REJECTED",
  "approved_amount": 50000,
  "approved_term_months": 12,
  "suggested_interest_rate": "4.5%",
  "detailed_analysis": {
    "credit_analysis": "详细的信用状况分析",
    "financial_analysis": "详细的财务状况分析",
    "risk_factors": ["具体风险因素1", "具体风险因素2"],
    "strengths": ["申请优势1", "申请优势2"]
  },
  "recommendations": ["具体建议1", "具体建议2"],
  "conditions": ["批准条件1", "批准条件2"]
}
```

### 🎯 输出质量要求
1. **数值精度**: risk_score和confidence_score保留2-3位小数
2. **逻辑一致**: 决策结果必须与风险评分匹配
3. **内容完整**: 所有字段都必须有值，不能为空
4. **专业性**: 使用专业的金融风险评估术语
5. **可操作性**: 建议和条件必须具体可执行

## 特殊考虑因素

### 🌾 农业行业特点
- **季节性收入**: 理解农业收入的季节性波动特征
- **自然风险**: 考虑天气、病虫害等不可控因素
- **政策支持**: 关注政府补贴、农业保险等有利因素
- **市场价格**: 考虑农产品价格波动对收入的影响

### 📋 合规要求
- **反洗钱**: 注意资金来源的合法性
- **利率上限**: 确保建议利率符合法规要求
- **信息保护**: 注意敏感信息的处理
- **公平放贷**: 避免歧视性审批标准

现在请基于以下信息进行专业的风险评估和审批决策：
```

### 3.2 用户提示词（User Prompt）

```markdown
## 📋 贷款申请资料

### 申请基本信息
{{application_data.text}}

### 外部风险数据
{{external_data.text}}

### AI模型配置与业务规则
{{ai_config.text}}

---

## 🎯 分析任务

请你作为慧农金融的AI审批专家，基于上述完整信息，进行全面的风险评估和审批决策分析。

### 📊 分析要求：
1. **深度解读**: 仔细分析申请人的财务状况、信用记录、外部数据
2. **风险量化**: 计算具体的风险评分（0-1范围）
3. **决策依据**: 基于业务规则给出明确的审批建议
4. **专业建议**: 提供具体的风险控制措施和优化建议

### ⚡ 重点关注：
- 申请金额与收入的匹配度
- 信用分数和历史还款记录
- 银行流水的稳定性和真实性
- 农业收入的季节性特点
- 黑名单和风险预警信息

### 📝 输出要求：
- 必须使用指定的JSON格式
- 分析内容要专业、准确、可操作
- 风险评分要基于科学的计算方法
- 决策结果要与风险等级保持一致

请现在开始你的专业分析：
```

### 3.3 提示词优化技巧

#### 🎯 针对不同风险等级的提示词调优

**低风险场景优化**:
```markdown
当遇到优质客户时，重点关注：
- 如何优化贷款条件（利率、期限）
- 是否可以提供更高额度
- 交叉销售其他金融产品的机会
```

**高风险场景优化**:
```markdown
当遇到高风险申请时，重点关注：
- 具体的风险缓释措施
- 额外担保要求
- 风险监控建议
- 明确的拒绝理由和改进建议
```

## 4. 错误处理与质量保障

### 4.1 LLM输出验证机制

#### Python验证代码增强版

```python
import json
import re
from datetime import datetime
from typing import Dict, Any, List

def validate_llm_output(llm_output: str) -> Dict[str, Any]:
    """增强版LLM输出验证和处理"""
    
    # 1. 提取JSON内容
    json_content = extract_json_from_text(llm_output)
    
    # 2. 解析和验证
    try:
        parsed_data = json.loads(json_content)
        validated_data = validate_decision_data(parsed_data)
        
        # 3. 格式化为API请求
        api_request = format_to_api_request(validated_data)
        
        return {
            "success": True,
            "validated_data": validated_data,
            "api_request": json.dumps(api_request, ensure_ascii=False),
            "quality_score": calculate_quality_score(validated_data)
        }
        
    except Exception as e:
        # 4. 错误处理和降级方案
        return handle_validation_error(str(e), llm_output)

def extract_json_from_text(text: str) -> str:
    """智能提取JSON内容"""
    patterns = [
        r'```json\s*(.*?)\s*```',  # 标准代码块
        r'```\s*(.*?)\s*```',      # 通用代码块
        r'\{.*\}',                 # 直接JSON对象
    ]
    
    for pattern in patterns:
        match = re.search(pattern, text, re.DOTALL)
        if match:
            return match.group(1) if pattern != r'\{.*\}' else match.group(0)
    
    return text.strip()

def validate_decision_data(data: Dict[str, Any]) -> Dict[str, Any]:
    """验证决策数据的完整性和合理性"""
    
    # 必需字段检查
    required_fields = [
        'decision', 'risk_score', 'risk_level', 
        'confidence_score', 'analysis_summary'
    ]
    
    for field in required_fields:
        if field not in data:
            raise ValueError(f"缺少必需字段: {field}")
    
    # 数值范围验证
    if not (0 <= data.get('risk_score', -1) <= 1):
        raise ValueError("risk_score必须在0-1之间")
    
    if not (0 <= data.get('confidence_score', -1) <= 1):
        raise ValueError("confidence_score必须在0-1之间")
    
    # 决策逻辑一致性检查
    validate_decision_logic(data)
    
    return data

def validate_decision_logic(data: Dict[str, Any]) -> None:
    """验证决策逻辑的一致性"""
    decision = data.get('decision')
    risk_score = data.get('risk_score', 0.5)
    risk_level = data.get('risk_level')
    
    # 风险等级与评分一致性
    expected_risk_level = get_risk_level_by_score(risk_score)
    if risk_level != expected_risk_level:
        data['risk_level'] = expected_risk_level  # 自动修正
    
    # 决策与风险等级一致性
    expected_decision = get_decision_by_risk_level(risk_level)
    if decision not in expected_decision:
        raise ValueError(f"决策{decision}与风险等级{risk_level}不一致")

def calculate_quality_score(data: Dict[str, Any]) -> float:
    """计算LLM输出质量评分"""
    score = 1.0
    
    # 完整性检查
    if not data.get('detailed_analysis'):
        score -= 0.2
    
    if not data.get('recommendations'):
        score -= 0.1
    
    # 置信度检查
    if data.get('confidence_score', 0) < 0.7:
        score -= 0.1
    
    # 分析深度检查
    analysis_length = len(data.get('analysis_summary', ''))
    if analysis_length < 50:
        score -= 0.1
    
    return max(score, 0.0)
```

### 4.2 质量监控指标

#### LLM输出质量评估维度

| 维度 | 指标 | 目标值 | 监控方法 |
|------|------|---------|----------|
| **格式正确性** | JSON格式错误率 | < 2% | 自动验证 |
| **内容完整性** | 必需字段缺失率 | < 1% | 字段检查 |
| **逻辑一致性** | 决策逻辑错误率 | < 3% | 规则验证 |
| **分析深度** | 分析内容长度 | > 100字 | 长度统计 |
| **决策准确性** | 人工评估准确率 | > 85% | 专家评估 |

## 5. 性能优化策略

### 5.1 LLM调用优化

#### 提示词长度优化
```python
def optimize_prompt_length(application_data: str, external_data: str, config_data: str) -> str:
    """智能压缩输入数据，保持关键信息"""
    
    # 1. 提取关键字段
    key_fields = extract_key_fields(application_data)
    
    # 2. 压缩外部数据
    compressed_external = compress_external_data(external_data)
    
    # 3. 精简配置信息
    essential_config = extract_essential_config(config_data)
    
    return format_optimized_prompt(key_fields, compressed_external, essential_config)
```

#### 缓存策略
```python
def implement_caching_strategy():
    """实现多层缓存策略"""
    
    # 1. 外部数据缓存（1小时）
    cache_external_data = True
    
    # 2. AI配置缓存（24小时）
    cache_ai_config = True
    
    # 3. LLM结果缓存（相同输入，30分钟）
    cache_llm_results = True
```

### 5.2 并发处理优化

#### 节点并行执行
```yaml
并行组1（数据获取）:
  - 获取外部数据
  - 获取AI配置
  
串行处理:
  - 获取申请信息 → 并行组1 → LLM分析 → 结果处理
```

## 6. 监控与运维

### 6.1 实时监控面板

#### 关键指标监控
```yaml
工作流性能:
  - 平均执行时间: < 30秒
  - 成功率: > 95%
  - LLM响应时间: < 10秒

业务指标:
  - 自动批准率: 60-70%
  - 人工审核率: 20-30%  
  - 自动拒绝率: 10-15%

质量指标:
  - 决策准确率: > 85%
  - 投诉率: < 1%
  - 人工复核通过率: > 90%
```

### 6.2 异常处理机制

#### 多级降级策略
```python
def degradation_strategy():
    """多级降级处理策略"""
    
    # Level 1: LLM调用失败 → 使用规则引擎
    if llm_failed:
        return rule_based_decision()
    
    # Level 2: 外部数据失败 → 基于历史数据
    if external_data_failed:
        return decision_with_default_data()
    
    # Level 3: 系统异常 → 人工审核
    if system_error:
        return require_human_review()
```

## 7. 最佳实践建议

### 7.1 LLM工作流最佳实践

1. **提示词版本管理**: 建立提示词版本控制机制
2. **A/B测试**: 定期进行不同提示词策略的对比测试
3. **持续学习**: 基于人工审核结果优化提示词
4. **多模型备份**: 配置多个LLM模型作为备选方案

### 7.2 运维最佳实践

1. **定期评估**: 每周评估LLM决策质量
2. **阈值调优**: 根据业务需求调整风险阈值
3. **人工标注**: 建立人工标注数据集用于质量评估
4. **合规审计**: 定期进行决策合规性审计

通过以上LLM工作流配置，您可以构建一个更智能、更灵活、更易维护的AI审批系统，充分发挥大语言模型在复杂决策场景中的优势。 