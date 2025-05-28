 # AI Agent Token 获取和管理指南

## 概述

AI Agent Token 是用于Dify等外部AI平台调用慧农金融后端API接口的认证凭证。本文档详细说明如何获取、配置和管理这些Token。

## 1. Token类型说明

### 1.1 AI Agent Token
- **用途**：供Dify工作流调用后端AI智能体接口
- **格式**：`AI-Agent-Token your_token_here`
- **权限**：访问AI智能体相关接口（获取申请信息、提交决策等）

### 1.2 System Token  
- **用途**：供系统内部服务调用
- **格式**：`System-Token your_token_here`
- **权限**：系统级操作权限

## 2. Token配置方法

### 2.1 配置文件方式（推荐）

在 `backend/configs/config.yaml` 中配置：

```yaml
ai:
  difyApiUrl: "http://172.18.120.57/v1/workflows/run"
  difyApiKey: "app-your-dify-api-key-here"
  agentTokens:
    - "ai_agent_secure_token_2024_v1"
    - "dify_huinong_finance_token_001"
    - "ai_risk_assessment_token_456"
  systemTokens:
    - "system_internal_api_token_2024"
    - "huinong_backend_system_token"
```

### 2.2 环境变量方式

```bash
# 设置环境变量
export AI_AGENT_TOKENS="ai_agent_secure_token_2024_v1,dify_huinong_finance_token_001"
export SYSTEM_TOKENS="system_internal_api_token_2024,huinong_backend_system_token"
```

## 3. 当前可用Token

### 3.1 预配置的AI Agent Token

以下Token已在配置文件中预配置，可直接使用：

| Token名称 | Token值 | 用途 | 状态 |
|-----------|---------|------|------|
| Dify主Token | `ai_agent_secure_token_2024_v1` | Dify工作流调用 | ✅ 活跃 |
| 风险评估Token | `dify_huinong_finance_token_001` | 风险评估服务 | ✅ 活跃 |
| 备用Token | `ai_risk_assessment_token_456` | 备用/测试 | ✅ 活跃 |

### 3.2 系统Token

| Token名称 | Token值 | 用途 | 状态 |
|-----------|---------|------|------|
| 系统主Token | `system_internal_api_token_2024` | 内部系统调用 | ✅ 活跃 |
| 后端服务Token | `huinong_backend_system_token` | 后端服务间调用 | ✅ 活跃 |

## 4. Dify平台配置步骤

### 4.1 在Dify中配置AI Agent Token

1. **进入自定义工具配置**
   - 登录Dify平台：`http://172.18.120.57`
   - 进入工具→自定义工具→创建工具

2. **配置认证信息**
   - 认证方式：选择 `API Key`
   - Header名称：`Authorization`
   - API Key值：`AI-Agent-Token ai_agent_secure_token_2024_v1`

3. **完整的认证配置示例**
   ```
   Authentication Type: API Key
   API Key Header Name: Authorization
   API Key Value: AI-Agent-Token ai_agent_secure_token_2024_v1
   ```

### 4.2 验证Token配置

使用curl命令测试Token有效性：

```bash
# 测试获取申请信息接口
curl -X GET "http://localhost:8080/api/v1/ai-agent/applications/test_app_001/info" \
  -H "Authorization: AI-Agent-Token ai_agent_secure_token_2024_v1" \
  -H "Content-Type: application/json"

# 测试获取外部数据接口
curl -X GET "http://localhost:8080/api/v1/ai-agent/external-data?user_id=user_001&data_types=credit_report,bank_flow" \
  -H "Authorization: AI-Agent-Token ai_agent_secure_token_2024_v1" \
  -H "Content-Type: application/json"

# 测试获取AI模型配置接口
curl -X GET "http://localhost:8080/api/v1/ai-agent/config/models" \
  -H "Authorization: AI-Agent-Token ai_agent_secure_token_2024_v1" \
  -H "Content-Type: application/json"
```

## 5. Token安全管理

### 5.1 Token生成建议

```bash
# 生成安全的Token（示例脚本）
#!/bin/bash

# 生成AI Agent Token
AI_TOKEN="ai_agent_$(date +%Y%m%d)_$(openssl rand -hex 8)"
echo "新的AI Agent Token: $AI_TOKEN"

# 生成System Token  
SYSTEM_TOKEN="system_$(date +%Y%m%d)_$(openssl rand -hex 8)"
echo "新的System Token: $SYSTEM_TOKEN"
```

### 5.2 Token轮换策略

1. **定期轮换**：建议每3个月轮换一次Token
2. **废弃旧Token**：新Token生效后，保留旧Token 7天缓冲期
3. **紧急轮换**：发现Token泄露时立即轮换

### 5.3 Token存储安全

- ✅ 使用配置文件存储（已加入.gitignore）
- ✅ 支持环境变量注入
- ✅ 生产环境使用密钥管理服务
- ❌ 不要硬编码在代码中
- ❌ 不要明文存储在数据库中

## 6. 权限控制

### 6.1 AI Agent Token权限

AI Agent Token只能访问以下接口：

- `GET /api/v1/ai-agent/applications/{id}/info` - 获取申请信息
- `GET /api/v1/ai-agent/external-data` - 获取外部数据  
- `POST /api/v1/ai-agent/applications/{id}/ai-decision` - 提交AI决策
- `GET /api/v1/ai-agent/config/models` - 获取AI模型配置

### 6.2 System Token权限

System Token具有系统级权限，可访问：

- 所有AI Agent Token的权限
- `POST /api/v1/ai-agent/applications/{id}/trigger-workflow` - 触发工作流
- `PUT /api/v1/ai-agent/applications/{id}/status` - 更新申请状态

## 7. 故障排查

### 7.1 常见错误

**错误1：401 Unauthorized - Missing authorization header**
```json
{
  "code": 401,
  "message": "Missing authorization header"
}
```
**解决方案**：检查请求头中是否包含Authorization字段

**错误2：401 Unauthorized - Invalid AI Agent token**
```json
{
  "code": 401, 
  "message": "Invalid AI Agent token"
}
```
**解决方案**：检查Token是否在配置文件的agentTokens列表中

**错误3：401 Unauthorized - Invalid authorization header format**
```json
{
  "code": 401,
  "message": "Invalid authorization header format"
}
```
**解决方案**：确保使用正确格式：`AI-Agent-Token your_token`

### 7.2 调试步骤

1. **检查配置文件**
   ```bash
   # 查看当前配置
   cat backend/configs/config.yaml | grep -A 10 "agentTokens"
   ```

2. **查看应用日志**
   ```bash
   # 查看认证相关日志
   tail -f logs/app.log | grep -i "agent.*token"
   ```

3. **验证服务运行状态**
   ```bash
   # 检查服务是否正常运行
   curl http://localhost:8080/health
   ```

## 8. 最佳实践

### 8.1 开发环境

- 使用较短的Token便于调试
- 在配置文件中明确标注Token用途
- 定期清理无用的测试Token

### 8.2 生产环境

- 使用长度至少32位的随机Token
- 启用Token访问日志记录
- 设置Token过期时间
- 实施IP白名单限制

### 8.3 监控告警

建议监控以下指标：

- Token使用频率异常
- 来自非预期IP的Token使用
- Token认证失败率过高
- Token未授权访问尝试

## 9. 总结

通过本指南，您应该能够：

1. ✅ 了解AI Agent Token的类型和用途
2. ✅ 在配置文件中正确配置Token
3. ✅ 在Dify平台中正确使用Token
4. ✅ 验证Token配置的有效性
5. ✅ 实施Token安全管理策略
6. ✅ 排查Token相关问题

如有疑问，请参考项目文档或联系技术支持团队。