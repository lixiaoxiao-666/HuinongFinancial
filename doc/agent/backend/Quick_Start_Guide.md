# AI智能体接口快速启动指南

## 快速开始

本指南将帮助您在5分钟内配置好AI智能体接口，并与Dify平台成功集成。

## 步骤1：获取AI Agent Token ⏱️ 1分钟

### 当前可用Token

您可以直接使用以下预配置的Token（已在配置文件中设置）：

```
AI Agent Token: ai_agent_secure_token_2024_v1
完整格式: AI-Agent-Token ai_agent_secure_token_2024_v1
```

## 步骤2：启动后端服务 ⏱️ 1分钟

```bash
# 进入后端目录
cd backend

# 启动服务
go run main.go

# 验证服务运行
curl http://localhost:8080/health
```

预期响应：
```json
{
  "status": "ok",
  "service": "digital-agriculture-backend", 
  "version": "1.0.0"
}
```

## 步骤3：验证AI Agent接口 ⏱️ 2分钟

### 3.1 测试Token认证

```bash
# 测试获取AI模型配置接口
curl -X GET "http://localhost:8080/api/v1/ai-agent/config/models" \
  -H "Authorization: AI-Agent-Token ai_agent_secure_token_2024_v1" \
  -H "Content-Type: application/json"
```

预期响应：
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "available_models": [...],
    "risk_thresholds": {...},
    "decision_rules": {...}
  }
}
```

### 3.2 测试获取申请信息

```bash
# 测试获取申请信息接口（使用测试数据）
curl -X GET "http://localhost:8080/api/v1/ai-agent/applications/test_app_001/info" \
  -H "Authorization: AI-Agent-Token ai_agent_secure_token_2024_v1" \
  -H "Content-Type: application/json"
```

### 3.3 测试获取外部数据

```bash
# 测试获取外部数据接口
curl -X GET "http://localhost:8080/api/v1/ai-agent/external-data?user_id=user_001&data_types=credit_report,bank_flow,blacklist_check" \
  -H "Authorization: AI-Agent-Token ai_agent_secure_token_2024_v1" \
  -H "Content-Type: application/json"
```

## 步骤4：配置Dify工具 ⏱️ 1分钟

### 4.1 导入OpenAPI Schema

1. 登录Dify平台：`http://172.18.120.57`
2. 进入 **工具** → **自定义工具** → **创建工具**
3. 选择 **OpenAPI** 导入方式
4. 复制并粘贴以下Schema：

```json
{
  "openapi": "3.1.0",
  "info": {
    "title": "慧农金融AI智能体接口",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "http://localhost:8080",
      "description": "本地环境"
    }
  ],
  "paths": {
    "/api/v1/ai-agent/applications/{application_id}/info": {
      "get": {
        "summary": "获取申请信息",
        "operationId": "getApplicationInfo",
        "parameters": [
          {
            "name": "application_id",
            "in": "path",
            "required": true,
            "schema": {"type": "string"}
          }
        ],
        "responses": {
          "200": {
            "description": "成功",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    },
    "/api/v1/ai-agent/external-data": {
      "get": {
        "summary": "获取外部数据",
        "operationId": "getExternalData",
        "parameters": [
          {
            "name": "user_id",
            "in": "query",
            "required": true,
            "schema": {"type": "string"}
          },
          {
            "name": "data_types",
            "in": "query",
            "required": true,
            "schema": {"type": "string"}
          }
        ],
        "responses": {
          "200": {
            "description": "成功",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    },
    "/api/v1/ai-agent/applications/{application_id}/ai-decision": {
      "post": {
        "summary": "提交AI决策",
        "operationId": "submitAIDecision",
        "parameters": [
          {
            "name": "application_id",
            "in": "path",
            "required": true,
            "schema": {"type": "string"}
          }
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "decision": {"type": "string"},
                  "risk_score": {"type": "number"},
                  "risk_level": {"type": "string"}
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "成功",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "AIAgentAuth": {
        "type": "apiKey",
        "in": "header",
        "name": "Authorization"
      }
    }
  },
  "security": [{"AIAgentAuth": []}]
}
```

### 4.2 配置认证信息

在Dify工具配置中：
- **认证方式**：API Key
- **Header名称**：Authorization  
- **API Key值**：`AI-Agent-Token ai_agent_secure_token_2024_v1`

### 4.3 测试工具连接

点击 **测试连接** 按钮，确保所有接口都能正常调用。

## 验证成功标志

✅ **后端服务**：health接口返回正常状态  
✅ **Token认证**：AI Agent接口正常响应，无401错误  
✅ **Dify工具**：工具导入成功，连接测试通过  
✅ **接口调用**：所有测试接口都返回预期的JSON格式响应

## 常见问题排查

### 问题1：401 Unauthorized

**现象**：
```json
{"code": 401, "message": "Invalid AI Agent token"}
```

**解决**：
- 检查Token格式：确保使用 `AI-Agent-Token ai_agent_secure_token_2024_v1`
- 检查配置文件：确认Token在 `configs/config.yaml` 的 `agentTokens` 列表中

### 问题2：连接拒绝

**现象**：
```
curl: (7) Failed to connect to localhost port 8080
```

**解决**：
- 确认后端服务正在运行：`go run main.go`
- 检查端口占用：`lsof -i :8080`
- 查看服务日志：检查启动错误信息

### 问题3：Dify无法连接

**现象**：Dify平台工具测试连接失败

**解决**：
- 确认网络连通性：Dify平台能否访问后端服务地址
- 检查防火墙设置：确保8080端口对外开放
- 确认服务地址：在Dify中配置正确的后端服务URL

## 下一步

配置完成后，您可以：

1. 📖 查看 [Dify工作流配置指南](./Dify_Setup_Guide.md) 创建完整的AI审批工作流
2. 📋 参考 [API实现状态报告](./API_Implementation_Status.md) 了解所有可用接口  
3. 🔐 阅读 [Token管理指南](./AI_Agent_Token_Guide.md) 了解Token安全管理
4. 🔄 查看 [AI工作流文档](./AI_Agent_Workflow.md) 了解完整的业务流程

## 技术支持

如遇到问题，请：
1. 查看服务日志：`tail -f logs/app.log`
2. 检查配置文件：`backend/configs/config.yaml`
3. 参考完整文档：`doc/agent/backend/` 目录下的详细文档

---

🎉 **恭喜！** 您已成功配置AI智能体接口系统，现在可以开始构建智能化的贷款审批工作流了！ 