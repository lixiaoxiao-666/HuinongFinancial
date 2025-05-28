# Dify平台实操配置指南

## 前提条件

1. **后端服务运行状态确认**
   ```bash
   # 确认后端服务正常运行
   curl http://localhost:8080/livez
   curl http://localhost:8080/readyz
   ```

2. **获取AI Agent Token**
   ```bash
   # 从配置文件或环境变量获取
   echo $AI_AGENT_TOKEN
   # 或查看配置文件中的token设置
   ```

3. **Dify平台访问**
   - 本地Dify地址：`http://172.18.120.57/app/c063c707-cbaa-4a4d-9412-88661aaf7753/workflow`
   - 确保可以正常访问并登录

## 第一步：创建自定义工具

### 1.1 进入工具管理页面

1. 登录Dify平台
2. 在左侧导航栏点击 **"工具"**
3. 点击 **"自定义工具"** 标签页
4. 点击右上角 **"创建工具"** 按钮

### 1.2 配置工具基本信息

1. **工具名称**：`慧农金融AI智能体`
2. **描述**：`慧农金融贷款申请AI审批相关接口工具`
3. **导入方式**：选择 `OpenAPI Schema`

### 1.3 导入OpenAPI Schema

将以下完整的OpenAPI Schema复制到输入框中：

```json
{
  "openapi": "3.1.0",
  "info": {
    "title": "慧农金融AI智能体接口",
    "description": "AI智能体审批工作流相关接口",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "http://localhost:8080",
      "description": "本地开发环境"
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

### 1.4 配置认证信息

1. **认证方式**：选择 `API Key`
2. **API Key设置**：
   - **Header名称**：`Authorization`
   - **API Key值**：`AI-Agent-Token your_actual_token_here`
   
   ⚠️ **注意**：请将`your_actual_token_here`替换为实际的Token值

3. 点击 **"测试连接"** 验证配置是否正确

### 1.5 保存工具

1. 确认所有配置无误
2. 点击 **"保存"** 按钮
3. 等待工具创建完成

## 第二步：创建AI审批工作流

### 2.1 新建工作流应用

1. 点击左侧导航 **"工作室"**
2. 点击 **"创建应用"**
3. 选择 **"工作流"**
4. 填写应用信息：
   - **应用名称**：`AI智能审批工作流`
   - **应用描述**：`慧农金融贷款申请AI智能审批系统`
   - **应用图标**：选择合适的图标

### 2.2 配置开始节点

1. 点击 **开始** 节点进行配置
2. 添加输入变量：

| 变量名 | 类型 | 是否必填 | 描述 |
|--------|------|----------|------|
| application_id | 文本 | 是 | 申请ID |
| callback_url | 文本 | 否 | 回调地址 |

3. 点击 **"保存"**

### 2.3 添加工具节点：获取申请信息

1. 点击开始节点右侧的 **"+"** 按钮
2. 选择 **"工具"**
3. 选择刚创建的工具 **"慧农金融AI智能体"**
4. 选择操作 **"getApplicationInfo"**
5. 配置参数：
   - **application_id**：点击右侧变量选择器，选择 `{{start.application_id}}`
6. 节点命名：`获取申请信息`
7. 点击 **"保存"**

### 2.4 添加代码节点：解析申请数据

1. 点击上一个节点右侧的 **"+"** 按钮
2. 选择 **"代码执行"**
3. 配置代码节点：
   - **节点名称**：`解析申请数据`
   - **编程语言**：Python3
   - **输入变量**：
     - 变量名：`application_info`
     - 类型：字符串
     - 值：`{{获取申请信息.text}}`

4. **代码内容**：
```python
import json

def main(application_info: str) -> dict:
    """解析申请信息"""
    try:
        # 解析JSON响应
        data = json.loads(application_info)
        if data.get('code') != 0:
            raise ValueError(f"API返回错误: {data.get('message', '未知错误')}")
        
        app_data = data['data']
        
        return {
            'user_id': app_data.get('applicant_info', {}).get('user_id', ''),
            'amount': app_data.get('amount', 0),
            'credit_score': app_data.get('applicant_info', {}).get('credit_score', 0),
            'annual_income': app_data.get('applicant_info', {}).get('annual_income', 0),
            'application_data': application_info,
            'success': True
        }
    except Exception as e:
        return {
            'user_id': '',
            'amount': 0,
            'credit_score': 0,
            'annual_income': 0,
            'application_data': application_info,
            'success': False,
            'error': str(e)
        }
```

5. **输出变量**：
   - user_id
   - amount
   - credit_score
   - annual_income
   - application_data
   - success
   - error

6. 点击 **"保存"**

### 2.5 添加工具节点：获取外部数据

1. 继续添加工具节点
2. 选择 **"慧农金融AI智能体"** → **"getExternalData"**
3. 配置参数：
   - **user_id**：`{{解析申请数据.user_id}}`
   - **data_types**：`credit_report,bank_flow,blacklist_check`
4. 节点命名：`获取外部数据`

### 2.6 添加代码节点：AI风险分析

1. 添加代码执行节点
2. 配置：
   - **节点名称**：`AI风险分析`
   - **输入变量**：
     - `application_data`：`{{解析申请数据.application_data}}`
     - `external_data`：`{{获取外部数据.text}}`
     - `credit_score`：`{{解析申请数据.credit_score}}`
     - `amount`：`{{解析申请数据.amount}}`
     - `annual_income`：`{{解析申请数据.annual_income}}`

3. **代码内容**：
```python
import json

def main(application_data: str, external_data: str, credit_score: float, amount: float, annual_income: float) -> dict:
    """AI风险分析"""
    try:
        # 解析外部数据
        ext_data_response = json.loads(external_data)
        if ext_data_response.get('code') != 0:
            # 如果外部数据获取失败，使用默认值
            ext_data = {}
        else:
            ext_data = ext_data_response.get('data', {})
        
        # 风险评分算法
        risk_score = 0.0
        risk_factors = []
        
        # 1. 信用分数评估
        if credit_score < 600:
            risk_score += 0.4
            risk_factors.append("信用分数过低")
        elif credit_score < 700:
            risk_score += 0.2
            risk_factors.append("信用分数中等")
        
        # 2. 收入债务比评估
        bank_flow = ext_data.get('bank_flow', {})
        debt_ratio = bank_flow.get('debt_to_income_ratio', 0.3)
        if debt_ratio > 0.5:
            risk_score += 0.3
            risk_factors.append("债务收入比过高")
        elif debt_ratio > 0.3:
            risk_score += 0.1
            risk_factors.append("债务收入比偏高")
        
        # 3. 黑名单检查
        blacklist = ext_data.get('blacklist_check', {})
        if blacklist.get('is_blacklisted', False):
            risk_score += 0.5
            risk_factors.append("存在黑名单记录")
        
        # 4. 申请金额与收入比
        if annual_income > 0:
            amount_income_ratio = amount / annual_income
            if amount_income_ratio > 0.5:
                risk_score += 0.2
                risk_factors.append("申请金额与收入比过高")
        
        # 限制风险分数范围
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
            approval_amount = min(amount, annual_income * 0.3) if annual_income > 0 else amount * 0.5
        else:
            approval_amount = 0
        
        return {
            'decision': decision,
            'risk_score': round(risk_score, 3),
            'risk_level': risk_level,
            'confidence': round(max(0.5, 1 - risk_score * 0.6), 3),
            'risk_factors': risk_factors,
            'approval_amount': approval_amount,
            'analysis_summary': f"风险评分: {risk_score:.3f}, 决策: {decision}",
            'success': True
        }
        
    except Exception as e:
        return {
            'decision': 'REQUIRE_HUMAN_REVIEW',
            'risk_score': 0.5,
            'risk_level': 'MEDIUM',
            'confidence': 0.1,
            'risk_factors': [f"分析错误: {str(e)}"],
            'approval_amount': 0,
            'analysis_summary': f"系统分析错误: {str(e)}",
            'success': False,
            'error': str(e)
        }
```

### 2.7 添加工具节点：提交AI决策

1. 添加工具节点
2. 选择 **"慧农金融AI智能体"** → **"submitAIDecision"**
3. 配置参数：
   - **application_id**：`{{start.application_id}}`
   - **Request Body**：切换到JSON模式，输入：

```json
{
  "decision": "{{AI风险分析.decision}}",
  "risk_score": {{AI风险分析.risk_score}},
  "risk_level": "{{AI风险分析.risk_level}}",
  "confidence": {{AI风险分析.confidence}},
  "analysis_result": {
    "credit_analysis": "基于信用分数的综合分析",
    "income_analysis": "收入与申请金额比例分析",
    "risk_factors": {{AI风险分析.risk_factors}},
    "approval_amount": {{AI风险分析.approval_amount}},
    "recommended_rate": "根据风险等级确定利率"
  },
  "processing_info": {
    "ai_model_version": "v1.2.0",
    "processing_time_ms": 2000,
    "workflow_id": "{{sys.workflow_id}}"
  }
}
```

4. 节点命名：`提交AI决策`

### 2.8 添加结束节点

1. 添加 **"结束"** 节点
2. 配置输出：
   - **decision**：`{{AI风险分析.decision}}`
   - **risk_score**：`{{AI风险分析.risk_score}}`
   - **analysis_summary**：`{{AI风险分析.analysis_summary}}`
   - **approval_amount**：`{{AI风险分析.approval_amount}}`
   - **processing_status**：`completed`

## 第三步：测试工作流

### 3.1 调试单个节点

1. 点击每个节点的 **"运行"** 按钮进行单独测试
2. 检查输入输出是否正确
3. 修复任何配置错误

### 3.2 端到端测试

1. 点击右上角 **"运行"** 按钮
2. 输入测试数据：
   ```json
   {
     "application_id": "test_app_001",
     "callback_url": "http://localhost:8080/callback"
   }
   ```
3. 观察执行过程和结果
4. 检查每个节点的输出

### 3.3 错误处理配置

为关键节点启用错误处理：

1. 选择工具节点
2. 在右侧面板找到 **"错误处理"** 选项
3. 启用 **"失败时重试"**
4. 设置：
   - 最大重试次数：3
   - 重试间隔：1000ms

## 第四步：发布工作流

### 4.1 发布为API

1. 工作流测试无误后，点击右上角 **"发布"**
2. 选择发布类型：**"API"**
3. 配置发布设置：
   - **版本号**：v1.0.0
   - **更新说明**：初始版本发布

### 4.2 获取API信息

发布成功后，系统会生成：

1. **工作流API端点**：用于外部调用
2. **API密钥**：用于身份验证
3. **调用示例**：包含完整的调用代码

### 4.3 配置后端系统

将获得的API信息配置到后端系统：

```yaml
# 在配置文件中添加
dify:
  api_url: "https://your-dify-domain/v1/workflows/run"
  api_key: "your_dify_api_key"
  workflow_id: "your_workflow_id"
```

## 第五步：监控和维护

### 5.1 监控工作流执行

1. 在Dify平台的 **"日志"** 页面查看执行记录
2. 监控成功率和错误情况
3. 分析性能指标

### 5.2 日常维护

1. **定期检查Token有效性**
2. **监控API调用频率和限制**
3. **更新工作流配置以适应业务变化**
4. **备份重要的工作流配置**

## 常见问题排查

### 问题1：工具连接测试失败

**排查步骤**：
1. 检查后端服务是否正常运行
2. 验证Token是否正确配置
3. 确认网络连通性
4. 检查API接口路径是否正确

### 问题2：工作流执行失败

**排查步骤**：
1. 查看具体的错误日志
2. 检查输入参数是否正确
3. 验证每个节点的配置
4. 确认数据流转是否正常

### 问题3：返回数据格式错误

**排查步骤**：
1. 检查后端API返回格式
2. 验证代码节点的解析逻辑
3. 确认变量引用是否正确

## 总结

完成以上步骤后，您将拥有一个完整的Dify AI智能审批工作流，能够：

1. ✅ 自动获取申请信息
2. ✅ 调用外部数据进行风险评估
3. ✅ 基于AI算法做出审批决策
4. ✅ 将结果提交回后端系统
5. ✅ 提供完整的执行日志和监控

这个工作流可以大大提高贷款审批的效率和准确性，实现真正的AI智能化审批。 