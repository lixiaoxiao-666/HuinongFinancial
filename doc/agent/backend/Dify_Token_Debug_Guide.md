# Dify Token配置调试指南

## 问题现象

从日志分析，发现两个主要问题：

1. **Postman请求**（172.18.120.10）：✅ Token认证成功，❌ 数据库无测试数据
2. **Dify请求**（172.20.0.9）：❌ Token格式错误 - "Invalid authorization header format"

## 原因分析

### Dify Token配置问题

在您的Dify配置截图中，我看到：
- **认证类型**：API Key ✅
- **认证头前缀**：Basic ✅ 
- **键**：Authorization ✅
- **值**：`AI-Agent-Token ai_agent_secure_token_2024_v1`

**问题**：Dify可能将这个值错误地处理为Basic认证格式。

## 解决方案

### 方案1：修改Dify配置（推荐）

1. **认证类型**：选择 `API Key`
2. **认证头前缀**：选择 `Custom` 而不是 `Basic`
3. **键**：`Authorization`
4. **值**：只填写 Token 部分：`ai_agent_secure_token_2024_v1`

这样Dify会发送：
```
Authorization: ai_agent_secure_token_2024_v1
```

然后我们需要修改后端中间件来适配这种格式。

### 方案2：修改后端中间件兼容多种格式

修改`AIAgentAuthMiddleware`以支持多种Token格式：

```go
// AIAgentAuthMiddleware AI智能体Token认证中间件
func AIAgentAuthMiddleware(config *conf.AIConfig) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		
		// 添加调试日志
		log := zap.L()
		log.Info("收到AI Agent认证请求",
			zap.String("client_ip", c.ClientIP()),
			zap.String("auth_header", authHeader),
			zap.String("user_agent", c.GetHeader("User-Agent")),
		)
		
		if authHeader == "" {
			log.Warn("缺少Authorization头")
			pkg.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		var token string
		
		// 支持多种Token格式
		if strings.HasPrefix(authHeader, "AI-Agent-Token ") {
			// 格式1: AI-Agent-Token your_token
			token = strings.TrimPrefix(authHeader, "AI-Agent-Token ")
			log.Info("检测到AI-Agent-Token格式", zap.String("token", token))
		} else if strings.HasPrefix(authHeader, "Bearer ") {
			// 格式2: Bearer your_token（Dify可能使用此格式）
			token = strings.TrimPrefix(authHeader, "Bearer ")
			log.Info("检测到Bearer格式", zap.String("token", token))
		} else if strings.HasPrefix(authHeader, "System-Token ") {
			// 格式3: System-Token your_token
			token = strings.TrimPrefix(authHeader, "System-Token ")
			log.Info("检测到System-Token格式", zap.String("token", token))
		} else {
			// 格式4: 直接是Token值（Custom API Key）
			token = authHeader
			log.Info("检测到直接Token格式", zap.String("token", token))
		}

		// 验证Token有效性
		isValid := false
		var agentType string
		
		// 检查AI Agent Token
		for _, validToken := range config.AgentTokens {
			if token == validToken {
				isValid = true
				agentType = "ai_agent"
				break
			}
		}
		
		// 如果不是AI Agent Token，检查System Token
		if !isValid {
			for _, validToken := range config.SystemTokens {
				if token == validToken {
					isValid = true
					agentType = "system"
					break
				}
			}
		}

		if !isValid {
			log.Warn("无效的Token", zap.String("token", token))
			pkg.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		// 设置上下文信息
		c.Set("agent_type", agentType)
		c.Set("agent_token", token)
		log.Info("Token验证成功", zap.String("agent_type", agentType))
		c.Next()
	})
}
```

### 方案3：网络配置调试

对于IP地址问题（172.20.0.9 vs 172.18.120.57），这可能是：

1. **Docker网络**：Dify运行在容器中，使用了不同的网络
2. **负载均衡**：有反向代理或负载均衡器
3. **Kubernetes网络**：如果使用K8s部署

**解决方案**：
- 在后端服务配置中添加信任的代理IP
- 在`configs/config.yaml`中添加：

```yaml
server:
  port: 8080
  mode: "debug"
  trusted_proxies:
    - "172.20.0.0/16"
    - "172.18.0.0/16"
  cors:
    allowed_origins:
      - "http://172.18.120.57"
      - "http://172.20.0.9"
```

## 立即解决步骤

### 步骤1：先解决数据问题

执行SQL插入测试数据（参考Test_Data_Setup.md）

### 步骤2：重启后端服务查看调试日志

```bash
cd backend
go run main.go
```

### 步骤3：测试Dify连接

在Dify中点击"测试"按钮，然后查看后端日志输出，确认：
1. 收到的Authorization头内容
2. Token验证过程
3. 失败的具体原因

### 步骤4：根据日志调整配置

根据日志中显示的`auth_header`内容，调整Dify配置或后端中间件。

## 常见的Dify Token配置模式

### 模式1：Custom API Key（推荐）
```
认证类型: API Key
前缀: Custom  
键: Authorization
值: ai_agent_secure_token_2024_v1
```

### 模式2：Bearer Token
```
认证类型: API Key
前缀: Bearer
键: Authorization  
值: ai_agent_secure_token_2024_v1
```

### 模式3：完整格式
```
认证类型: API Key
前缀: Custom
键: Authorization
值: AI-Agent-Token ai_agent_secure_token_2024_v1
```

## 验证步骤

配置完成后，验证Token是否正常工作：

1. 查看后端日志确认收到正确的Authorization头
2. 测试所有AI Agent接口
3. 确认Dify工具连接测试通过

如果仍有问题，请提供后端日志中的具体`auth_header`内容，我将进一步协助解决。 