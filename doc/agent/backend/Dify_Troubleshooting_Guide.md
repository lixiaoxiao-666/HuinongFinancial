# Dify工作流配置问题排查指南

## 🚨 常见问题1：开始节点显示"没有可用的赋值变量"

### 问题现象
在配置开始节点时，界面显示"没有可用的赋值变量"，无法添加输入参数。

### 🔍 原因分析
1. **应用类型错误**：创建的是聊天应用而不是工作流应用
2. **节点选择错误**：选择了错误的节点或位置
3. **界面状态错误**：浏览器缓存或界面异常

### ✅ 解决步骤

#### 方案1：重新创建工作流应用

1. **删除现有应用**（如果创建错误）
   - 在应用列表中找到错误的应用
   - 点击设置 → 删除应用

2. **正确创建工作流**
   ```
   步骤1：点击"创建应用"
   步骤2：选择"工作流"（重要：不是聊天助手）
   步骤3：填写应用名称：AI智能审批工作流
   步骤4：点击"创建"
   ```

#### 方案2：检查节点选择

1. **确认开始节点**
   - 画布最左边的蓝色圆形节点
   - 节点标签显示"开始"或"Start"
   - 点击后右侧显示节点配置面板

2. **确认配置区域**
   - 右侧面板显示"开始节点配置"
   - 找到"输入变量"或"输入参数"区域
   - 应该有"+ 添加变量"按钮

#### 方案3：浏览器问题解决

1. **清除浏览器缓存**
   ```
   Chrome: Ctrl+Shift+Delete
   Edge: Ctrl+Shift+Delete
   Safari: Cmd+Option+E
   ```

2. **刷新页面**
   - 按F5或Ctrl+R刷新页面
   - 重新点击开始节点配置

3. **尝试不同浏览器**
   - 推荐使用Chrome或Edge最新版本
   - 避免使用IE浏览器

### 📱 正确配置界面示例

```
Dify工作流编辑器
├── 画布区域
│   └── [开始] → [+] （蓝色开始节点）
└── 右侧配置面板
    ├── 📝 节点配置
    ├── 📝 输入变量
    │   ├── [+ 添加变量] 按钮
    │   └── 变量列表（初始为空）
    └── 💾 保存按钮
```

## 🚨 常见问题2：工具连接失败

### 问题现象
- 提示"连接失败"或"网络错误"
- Token认证不通过

### 解决方案

1. **检查后端服务**
   ```bash
   # 确认服务运行
   curl http://localhost:8080/livez
   
   # 检查AI Agent接口
   curl -H "Authorization: AI-Agent-Token ai_agent_secure_token_2024_v1" \
        http://localhost:8080/api/v1/ai-agent/applications/test_app_001/info
   ```

2. **验证Token配置**
   - 确认Token值正确
   - 检查认证头格式：`AI-Agent-Token your_token_here`

3. **网络连通性**
   - 确认Dify可以访问localhost:8080
   - 检查防火墙设置

## 🚨 常见问题3：节点参数引用错误

### 问题现象
- 节点执行失败
- 参数值显示为空或undefined

### 解决方案

1. **正确的参数引用格式**
   ```
   正确：{{start.application_id}}
   错误：{start.application_id}
   错误：{{application_id}}
   ```

2. **检查变量名一致性**
   - 开始节点变量名：application_id
   - 引用时：{{start.application_id}}

## 🚨 常见问题4：JSON解析错误

### 问题现象
- 代码节点执行失败
- JSON解析异常

### 解决方案

1. **添加错误处理**
   ```python
   try:
       data = json.loads(response)
   except json.JSONDecodeError as e:
       return {"error": f"JSON解析失败: {str(e)}"}
   ```

2. **检查API返回格式**
   - 确认后端返回标准JSON格式
   - 验证Content-Type为application/json

## 🆘 紧急解决方案

如果上述方案都无法解决，请尝试：

### 1. 使用简化版工作流

创建最小可工作版本：
```
[开始] → [获取申请信息] → [结束]
```

### 2. 逐步添加节点

1. 先配置并测试开始节点
2. 添加第一个工具节点并测试
3. 确认可工作后再添加其他节点

### 3. 导入预配置模板

使用以下简化的OpenAPI Schema：

```json
{
  "openapi": "3.1.0",
  "info": {
    "title": "AI Agent API",
    "version": "1.0.0"
  },
  "servers": [{"url": "http://localhost:8080"}],
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

## 📞 技术支持

如果问题仍然存在，请提供以下信息：

1. **环境信息**
   - Dify版本号
   - 浏览器类型和版本
   - 操作系统

2. **错误截图**
   - 完整的错误界面
   - 浏览器开发者工具中的错误信息

3. **操作步骤**
   - 详细的操作步骤
   - 何时出现错误

4. **后端服务状态**
   - 服务是否正常运行
   - API测试结果 