# Dify 快速部署指南 - 数字惠农项目（优化版）

## 🚀 快速开始（5分钟上手）

### 1. 注册并配置 Dify

1. **注册账号**
   - 访问 [Dify 官网](https://dify.ai/)
   - 注册账号并创建工作区

2. **获取 API Key**
   - 进入 Dify 控制台
   - 在 "设置" -> "API Keys" 中生成新的 API Key
   - 保存这个 Key，后面会用到

### 2. 创建工具（Tools）

#### 🔐 鉴权配置说明

**API Token 获取**：
1. 在后端配置文件 `backend/configs/config.yaml` 中设置：
   ```yaml
   dify:
     api_token: "dify-huinong-secure-token-2024"
   ```

2. 或通过环境变量设置：
   ```bash
   export DIFY_API_TOKEN="dify-huinong-secure-token-2024"
   ```

**在 Dify 工具中配置**：
- 鉴权类型：选择 `Custom`
- 键：`Authorization` 
- 值：`Bearer dify-huinong-secure-token-2024`

#### 🌐 网络配置重要说明

**服务器地址配置**：
- **后端服务器**：监听在 `0.0.0.0:8080`，可以接受来自任何IP的连接
- **Dify工具配置**：需要使用能够从Dify服务器访问到后端的IP地址

**常见网络配置场景**：

1. **本地开发环境**：
   ```json
   "servers": [
     {
       "url": "http://localhost:8080"
     }
   ]
   ```

2. **Docker环境**：
   ```json
   "servers": [
     {
       "url": "http://host.docker.internal:8080"
     }
   ]
   ```

3. **局域网环境**：
   ```json
   "servers": [
     {
       "url": "http://192.168.1.100:8080"  // 替换为实际的内网IP
     }
   ]
   ```

4. **云服务器环境**：
   ```json
   "servers": [
     {
       "url": "http://your-server-ip:8080"  // 替换为服务器的公网或内网IP
     }
   ]
   ```

**如何确定正确的IP地址**：

1. **查看本机IP地址**：
   ```bash
   # Windows
   ipconfig
   
   # Linux/Mac
   ifconfig
   ip addr show
   ```

2. **测试网络连通性**：
   ```bash
   # 测试端口是否开放
   telnet your-ip 8080
   
   # 或使用curl测试API
   curl -X GET http://your-ip:8080/health
   ```

3. **防火墙配置**：
   确保8080端口在防火墙中已开放：
   ```bash
   # Windows防火墙
   netsh advfirewall firewall add rule name="HuinongAPI" dir=in action=allow protocol=TCP localport=8080
   
   # Linux iptables
   iptables -A INPUT -p tcp --dport 8080 -j ACCEPT
   
   # Ubuntu ufw
   ufw allow 8080
   ```

#### 🛠️ 整合工具定义

在 Dify 控制台的 "工具" 页面，创建一个包含所有4个API的整合工具：

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
      "url": "http://172.18.120.10:8080"
    }
  ],
  "paths": {
    "/api/internal/dify/loan/get-application-details": {
      "post": {
        "description": "获取贷款申请的详细信息，包括用户基本信息、信用记录、财务状况等",
        "operationId": "get_loan_application_details",
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
    "/api/internal/dify/loan/submit-assessment": {
      "post": {
        "description": "提交AI风险评估结果到后端系统",
        "operationId": "submit_risk_assessment",
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
        "operationId": "get_machine_rental_details",
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
    },
    "/api/internal/dify/credit/query": {
      "post": {
        "description": "查询用户征信报告信息",
        "operationId": "query_credit_report",
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
    }
  },
  "components": {
    "securitySchemes": {
      "BearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    }
  }
}
```

#### 🔑 环境变量配置

在创建工具时，需要配置以下环境变量：

- **BASE_URL**: `http://your-backend-domain.com`（您的后端服务地址）
- **API_TOKEN**: `dify-huinong-secure-token-2024`（您设置的API Token）

#### ✅ 工具创建完成后

创建成功后，您就可以在 Dify 工作流中使用以下4个操作：
1. `get_loan_application_details` - 获取贷款申请详情
2. `submit_risk_assessment` - 提交风险评估结果  
3. `get_machine_rental_details` - 获取农机租赁详情
4. `query_credit_report` - 查询征信报告

### 3. 创建贷款审批工作流（优化版）

1. **新建工作流**
   - 在 Dify 控制台点击 "工作流"
   - 点击 "创建工作流"
   - 选择 "从空白开始"
   - 命名为 `loan_risk_assessment_v2`

2. **配置工作流节点（优化版 - 5个节点）**

   **节点1：开始节点 (start)**
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

   **节点2：获取申请详情 (get_application_details)**
   - 节点类型：Tool
   - 工具：`get_loan_application_details`
   - 输入映射：
     ```json
     {
       "application_id": "{{#start.application_id#}}",
       "user_id": "{{#start.user_id#}}",
       "include_credit": true
     }
     ```
   - 输出变量：`application_details`

   **节点3：AI智能评估 (ai_smart_assessment)**
   - 节点类型：LLM
   - 模型：GPT-4o
   - 输入变量：
     - `application_data`: `{{#get_application_details.text#}}`
   
   **系统提示词**：
   ```
   你是专业的农业金融风险评估专家。你将收到完整的贷款申请JSON数据，需要进行全面评估并做出审批决策。

   ## 评估流程
   1. **数据解析**: 从JSON中提取关键信息
   2. **认证检查**: 检查用户认证完整度
   3. **风险评估**: 分析信用、财务、经验等风险因素
   4. **决策判断**: 基于评估结果做出approve/manual/reject决策
   5. **输出建议**: 提供具体的审批参数和意见

   ## 评估规则

   ### 认证要求
   - 实名认证、银行卡认证、征信认证必须全部完成才能批准
   - 认证不完整直接拒绝

   ### 风险等级标准

   **低风险 (approve) - 直接批准**
   必须同时满足：
   - 认证完整度：完全认证
   - 信用评分 >= 700
   - 负债收入比 <= 0.3
   - 农业经验 >= 3年
   - 无逾期记录或逾期次数 <= 1且天数 <= 30
   - 月供压力 <= 30%

   **中风险 (manual) - 人工审核**
   认证完整且满足：
   - 信用评分 600-699
   - 负债收入比 0.3-0.5
   - 农业经验 1-3年
   - 逾期次数 <= 2且最大逾期天数 <= 90
   - 月供压力 30-50%

   **高风险 (reject) - 直接拒绝**
   满足以下任一条件：
   - 认证信息不完整
   - 信用评分 < 600
   - 负债收入比 > 0.5
   - 农业经验 < 1年
   - 逾期次数 > 2或最大逾期天数 > 90
   - 月供压力 > 50%

   ## 推荐策略
   - **低风险**: 批准全额或适当金额，给予利率优惠(8-10%)
   - **中风险**: 批准部分金额(50-80%)，标准利率(10-12%)，建议加强监控
   - **高风险**: 拒绝申请，详细说明原因和改进建议

   ## 输出要求
   必须严格按照以下JSON格式输出，不要包含任何其他文字：
   {
     "risk_level": "low|medium|high",
     "decision": "approve|manual|reject",
     "recommended_amount": 推荐金额数字(approve/manual时必填),
     "recommended_term": 推荐期限数字(approve/manual时必填),
     "recommended_rate": 推荐年利率数字(approve/manual时必填),
     "risk_factors": ["风险因素1", "风险因素2"],
     "comments": "详细评估意见，包括决策理由和建议",
     "confidence_score": 置信度数字(0-1之间),
     "auth_status": "complete|partial|incomplete",
     "key_metrics": {
       "credit_score": 信用评分,
       "debt_ratio": 负债收入比,
       "experience_years": 农业经验年数,
       "monthly_payment_ratio": 月供收入比
     }
   }
   ```
   
   **用户提示词**：
   ```
   请对以下贷款申请进行全面评估：

   ## 申请详情数据
   {{#get_application_details.text#}}

   ## 评估要求
   1. 仔细解析JSON数据，提取所有关键信息
   2. 首先检查认证状态(real_name_verified, bank_card_verified, credit_verified)
   3. 如果认证不完整，直接拒绝并说明原因
   4. 如果认证完整，进行详细的风险评估
   5. 计算关键指标：负债收入比、月供压力、信用风险等
   6. 根据评估规则做出决策
   7. 提供具体的审批建议和参数

   请严格按照JSON格式输出完整的评估结果。
   ```
   
   **输出变量**：`ai_assessment` (Object)

   **节点4：提交评估结果 (submit_assessment)**
   - 节点类型：Tool
   - 工具：`submit_risk_assessment`
   - 输入映射：
     ```json
     {
       "application_id": "{{#start.application_id#}}",
       "risk_level": "{{#ai_smart_assessment.risk_level#}}",
       "decision": "{{#ai_smart_assessment.decision#}}",
       "recommended_amount": "{{#ai_smart_assessment.recommended_amount#}}",
       "recommended_term": "{{#ai_smart_assessment.recommended_term#}}",
       "recommended_rate": "{{#ai_smart_assessment.recommended_rate#}}",
       "risk_factors": "{{#ai_smart_assessment.risk_factors#}}",
       "comments": "{{#ai_smart_assessment.comments#}}",
       "confidence_score": "{{#ai_smart_assessment.confidence_score#}}"
     }
     ```

   **节点5：结束节点 (end)**
   - 输出变量：
     ```json
     {
       "workflow_result": "{{#submit_assessment.text#}}",
       "ai_assessment": "{{#ai_smart_assessment#}}",
       "final_decision": "{{#ai_smart_assessment.decision#}}",
       "final_risk_level": "{{#ai_smart_assessment.risk_level#}}",
       "application_id": "{{#start.application_id#}}",
       "processing_status": "completed",
       "processing_timestamp": "{{sys.current_time}}"
     }
     ```

3. **连接节点流程（优化版）**
   ```
   Start (节点1)
     ↓
   Get Application Details (节点2)
     ↓
   AI Smart Assessment (节点3)
     ↓
   Submit Assessment (节点4)
     ↓
   End (节点5)
   ```

4. **优化亮点**
   - **流程简化**: 从原来的9个节点简化为5个节点
   - **AI统一决策**: 将认证检查、风险评估、决策判断全部交给AI处理
   - **减少错误点**: 减少数据传递环节，降低出错概率
   - **提高效率**: 减少节点间的数据解析和转换
   - **智能化程度高**: AI直接处理原始JSON数据，做出完整决策

### 4. 配置后端环境

1. **修改配置文件**
   
   编辑 `backend/configs/config.yaml`：
   ```yaml
   dify:
     api_url: "https://api.dify.ai/v1"
     api_key: "你的Dify API Key"
     timeout: 30
     retry_times: 3
     workflows:
       loan_approval_v2: "你的优化工作流ID"
       risk_assessment: "你的农机工作流ID"
   ```

2. **设置环境变量**
   ```bash
   export BASE_URL=http://localhost:8080
   export API_TOKEN=your-secure-api-token
   ```

3. **在Dify中配置环境变量**
   
   在工作流设置中添加：
   - `BASE_URL`: `http://your-backend-domain.com`
   - `API_TOKEN`: `your-secure-api-token`

### 5. 启动服务

1. **启动后端服务**
   ```bash
   cd backend
   go run cmd/server/main.go
   ```

2. **测试 Dify 集成**
   ```bash
   # 测试获取申请详情接口
   curl -X POST http://localhost:8080/api/internal/dify/loan/get-application-details \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer your-secure-api-token" \
     -d '{"application_id": "1", "user_id": "1", "include_credit": true}'
   ```

### 6. 测试工作流

1. **在 Dify 控制台测试**
   - 进入工作流编辑器
   - 点击 "测试" 按钮
   - 输入测试数据：
     ```json
     {
       "application_id": "1",
       "user_id": "1"
     }
     ```

2. **在后端代码中调用**
   ```go
   // 在贷款申请提交后调用
   response, err := difyService.CallLoanApprovalWorkflowV2(applicationID, userID)
   if err != nil {
       log.Printf("调用Dify工作流失败: %v", err)
   }
   ```

## ⚠️ 注意事项

1. **安全性**
   - API Token 必须保密
   - 建议设置 IP 白名单
   - 使用 HTTPS 协议

2. **错误处理**
   - 工作流可能超时，需要设置合理的超时时间
   - 添加重试机制
   - 记录详细的执行日志

3. **性能优化**
   - 缓存频繁查询的数据
   - 优化 LLM 提示词
   - 监控工作流执行时间

## 🔧 故障排除

### 常见问题

1. **工具调用失败**
   - 检查 API URL 是否正确
   - 验证 API Token 是否有效
   - 确认后端服务是否运行

2. **工作流执行超时**
   - 增加超时时间配置
   - 检查网络连接
   - 优化 LLM 响应速度

3. **权限错误**
   - 确认 Authorization 头格式正确
   - 检查 API Token 权限
   - 验证 IP 白名单设置

### 调试建议

1. **查看日志**
   ```bash
   tail -f backend/logs/app.log
   ```

2. **监控工作流**
   - 在 Dify 控制台查看执行日志
   - 检查每个节点的输入输出
   - 分析错误信息

3. **测试接口**
   ```bash
   # 健康检查
   curl http://localhost:8080/health
   
   # API 版本信息
   curl http://localhost:8080/api/public/version
   ```

## 🎉 优化成果

通过这次优化，我们实现了：

1. **节点数量减少**: 从9个节点减少到5个节点
2. **流程简化**: 去除了复杂的数据解析和条件分支
3. **AI智能化**: AI直接处理所有评估逻辑
4. **维护性提升**: 更少的节点意味着更容易维护和调试
5. **执行效率**: 减少数据传递环节，提高执行速度

完成以上步骤后，您的优化版 Dify AI 工作流就可以与数字惠农后端系统高效集成工作了！🎉