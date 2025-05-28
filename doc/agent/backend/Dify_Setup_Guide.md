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

#### 详细操作步骤：

1. **点击开始节点**
   - 在工作流画布中，点击蓝色的 **"开始"** 节点
   - 右侧会出现节点配置面板

2. **添加输入变量**（重要：这里是添加**输入参数**，不是赋值变量）
   
   在右侧配置面板中：
   - 找到 **"输入变量"** 或 **"变量"** 部分
   - 点击 **"+ 添加变量"** 按钮
   
3. **配置第一个变量：application_id**
   ```
   变量名：application_id
   变量类型：文本 (Text)
   是否必填：是 ✓
   描述：贷款申请ID
   默认值：（留空）
   ```

4. **配置第二个变量：callback_url**
   ```
   变量名：callback_url  
   变量类型：文本 (Text)
   是否必填：否
   描述：处理完成后的回调地址
   默认值：（留空）
   ```

#### 💡 如果仍显示"没有可用的赋值变量"：

**原因分析**：您可能在错误的位置尝试添加变量

**解决方案**：
1. 确保点击的是画布最左边的 **"开始"** 节点（蓝色圆形图标）
2. 在右侧面板中找到 **"输入变量"** 区域（不是"输出变量"）
3. 确保选择的是 **"输入"** 标签页，而不是其他标签

#### 🎯 正确的界面应该显示：
```
┌─ 开始节点配置 ─────────────┐
│                           │
│  📝 输入变量               │
│  ┌─────────────────────┐  │
│  │ + 添加变量           │  │
│  └─────────────────────┘  │
│                           │
│  💡 这里应该显示空的变量    │
│     列表，可以添加新变量    │
└───────────────────────────┘
```

#### 🚨 常见错误和解决方法：

**错误1**：在工具节点中尝试添加变量
- **解决**：只有开始节点需要添加输入变量

**错误2**：选择了错误的变量类型  
- **解决**：确保选择"文本"类型，不是"数字"或"文件"

**错误3**：在错误的标签页操作
- **解决**：确保在"输入"或"变量"标签页，不是在"输出"标签页

3. 点击 **"保存"** 完成开始节点配置

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
   
4. **输入变量配置**：
   - **变量名**：`application_info`
   - **类型**：字符串 (String)
   - **值**：`{{获取申请信息.text}}`

5. **代码内容**：
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
        applicant_info = app_data.get('applicant_info', {})
        
        return {
            'user_id': applicant_info.get('user_id', ''),
            'amount': float(app_data.get('amount', 0)),
            'credit_score': float(applicant_info.get('credit_score', 0)),
            'annual_income': float(applicant_info.get('annual_income', 0)),
            'application_data': application_info,
            'success': True,
            'error': None
        }
    except Exception as e:
        return {
            'user_id': '',
            'amount': 0.0,
            'credit_score': 0.0,
            'annual_income': 0.0,
            'application_data': application_info,
            'success': False,
            'error': str(e)
        }
```

6. **输出变量配置**（重要：请按照以下配置添加输出变量）：

| 变量名 | 变量类型 | 描述 | 示例值 |
|--------|----------|------|--------|
| `user_id` | String (文本) | 用户ID | "user_001" |
| `amount` | Number (数字) | 申请金额 | 100000.0 |
| `credit_score` | Number (数字) | 信用分数 | 750.0 |
| `annual_income` | Number (数字) | 年收入 | 200000.0 |
| `application_data` | String (文本) | 完整申请数据JSON | "{\"code\":0,\"data\":{...}}" |
| `success` | String (文本) | 处理是否成功 | "True" |
| `error` | String (文本) | 错误信息（如有） | null 或 错误描述 |

**⚠️ 配置注意事项**：
- 在代码节点的右侧面板中，找到 **"输出变量"** 区域
- 点击 **"+ 添加变量"** 按钮，逐一添加上述7个变量
- 确保变量名称完全匹配（区分大小写）
- Number类型用于数值计算，String类型用于文本和布尔值

7. 点击 **"保存"** 完成配置

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
   - **编程语言**：Python3
   
3. **输入变量配置**：

| 变量名 | 变量类型 | 引用值 | 描述 |
|--------|----------|--------|------|
| `application_data` | String | `{{解析申请数据.application_data}}` | 完整申请数据 |
| `external_data` | String | `{{获取外部数据.text}}` | 外部数据响应 |
| `credit_score` | Number | `{{解析申请数据.credit_score}}` | 信用分数 |
| `amount` | Number | `{{解析申请数据.amount}}` | 申请金额 |
| `annual_income` | Number | `{{解析申请数据.annual_income}}` | 年收入 |

4. **代码内容**：
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
        
        # 1. 信用分数评估 (40%权重)
        if credit_score < 600:
            risk_score += 0.4
            risk_factors.append("信用分数过低")
        elif credit_score < 700:
            risk_score += 0.2
            risk_factors.append("信用分数中等")
        
        # 2. 收入债务比评估 (30%权重)
        bank_flow = ext_data.get('bank_flow', {})
        debt_ratio = bank_flow.get('debt_to_income_ratio', 0.3)
        if debt_ratio > 0.5:
            risk_score += 0.3
            risk_factors.append("债务收入比过高")
        elif debt_ratio > 0.3:
            risk_score += 0.1
            risk_factors.append("债务收入比偏高")
        
        # 3. 黑名单检查 (50%权重)
        blacklist = ext_data.get('blacklist_check', {})
        if blacklist.get('is_blacklisted', False):
            risk_score += 0.5
            risk_factors.append("存在黑名单记录")
        
        # 4. 申请金额与收入比 (20%权重)
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
        
        # 计算置信度
        confidence = round(max(0.5, 1 - risk_score * 0.6), 3)
        
        # 生成分析摘要
        analysis_summary = f"风险评分: {risk_score:.3f}, 决策: {decision}, 风险等级: {risk_level}"
        
        return {
            'decision': decision,
            'risk_score': round(risk_score, 3),
            'risk_level': risk_level,
            'confidence': confidence,
            'risk_factors': json.dumps(risk_factors, ensure_ascii=False),
            'approval_amount': approval_amount,
            'analysis_summary': analysis_summary,
            'success': True,
            'error': None
        }
        
    except Exception as e:
        return {
            'decision': 'REQUIRE_HUMAN_REVIEW',
            'risk_score': 0.5,
            'risk_level': 'MEDIUM',
            'confidence': 0.1,
            'risk_factors': json.dumps([f"分析错误: {str(e)}"], ensure_ascii=False),
            'approval_amount': 0,
            'analysis_summary': f"系统分析错误: {str(e)}",
            'success': False,
            'error': str(e)
        }
```

5. **输出变量配置**（重要：请按照以下配置添加输出变量）：

| 变量名 | 变量类型 | 描述 | 可能值 |
|--------|----------|------|--------|
| `decision` | String (文本) | 审批决策 | "AUTO_APPROVED", "REQUIRE_HUMAN_REVIEW", "AUTO_REJECTED" |
| `risk_score` | Number (数字) | 风险评分 (0-1) | 0.245 |
| `risk_level` | String (文本) | 风险等级 | "LOW", "MEDIUM", "HIGH" |
| `confidence` | Number (数字) | 置信度 (0-1) | 0.753 |
| `risk_factors` | String (文本) | 风险因素列表JSON | ["信用分数中等", "债务收入比偏高"] |
| `approval_amount` | Number (数字) | 建议批准金额 | 80000.0 |
| `analysis_summary` | String (文本) | 分析摘要 | "风险评分: 0.245, 决策: AUTO_APPROVED, 风险等级: LOW" |
| `success` | String (文本) | 处理是否成功 | "True" |
| `error` | String (文本) | 错误信息（如有） | null 或 错误描述 |

**⚠️ 配置注意事项**：
- 在代码节点的右侧面板中，找到 **"输出变量"** 区域
- 点击 **"+ 添加变量"** 按钮，逐一添加上述9个变量
- 确保变量名称完全匹配（区分大小写）
- `risk_factors`使用JSON字符串格式存储数组数据
- Number类型变量用于数值计算和条件判断

6. 点击 **"保存"** 完成配置

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
    "credit_analysis": "基于信用分数和外部数据的综合分析",
    "income_analysis": "收入与申请金额比例分析",
    "risk_factors": "{{AI风险分析.risk_factors}}",
    "approval_amount": {{AI风险分析.approval_amount}},
    "recommended_rate": "根据风险等级确定利率"
  },
  "processing_info": {
    "ai_model_version": "v1.2.0",
    "processing_time_ms": 2000,
    "workflow_id": "{{sys.workflow_id}}",
    "analysis_summary": "{{AI风险分析.analysis_summary}}"
  }
}
```

**⚠️ JSON配置说明**：
- `{{AI风险分析.decision}}`：字符串类型，需要双引号包围
- `{{AI风险分析.risk_score}}`：数字类型，不需要引号
- `{{AI风险分析.confidence}}`：数字类型，不需要引号  
- `{{AI风险分析.approval_amount}}`：数字类型，不需要引号
- `"{{AI风险分析.risk_factors}}"`：JSON字符串，已在代码中转换为字符串

4. 节点命名：`提交AI决策`

5. 点击 **"保存"** 完成配置

### 2.8 添加结束节点

1. 添加 **"结束"** 节点
2. 配置输出变量：

| 变量名 | 引用值 | 类型 | 描述 |
|--------|--------|------|------|
| `decision` | `{{AI风险分析.decision}}` | String | 最终审批决策 |
| `risk_score` | `{{AI风险分析.risk_score}}` | Number | 风险评分 |
| `risk_level` | `{{AI风险分析.risk_level}}` | String | 风险等级 |
| `analysis_summary` | `{{AI风险分析.analysis_summary}}` | String | 分析摘要 |
| `approval_amount` | `{{AI风险分析.approval_amount}}` | Number | 建议批准金额 |
| `confidence` | `{{AI风险分析.confidence}}` | Number | 决策置信度 |
| `processing_status` | `completed` | String | 处理状态（固定值） |
| `application_id` | `{{start.application_id}}` | String | 申请ID |

**⚠️ 结束节点配置要点**：
- 结束节点的输出变量将作为整个工作流的最终返回值
- 这些变量可以被外部系统通过API调用获取
- 确保包含所有关键的决策信息，便于后续处理

3. 点击 **"保存"** 完成结束节点配置

## 🔍 变量配置汇总检查表

在开始测试之前，请使用以下表格检查所有节点的变量配置是否正确：

### 开始节点变量
| 变量名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| `application_id` | String | ✓ | 申请ID |
| `callback_url` | String | ✗ | 回调地址 |

### 解析申请数据节点
**输入变量**：
| 变量名 | 类型 | 引用值 |
|--------|------|--------|
| `application_info` | String | `{{获取申请信息.text}}` |

**输出变量**：
| 变量名 | 类型 | 描述 |
|--------|------|------|
| `user_id` | String | 用户ID |
| `amount` | Number | 申请金额 |
| `credit_score` | Number | 信用分数 |
| `annual_income` | Number | 年收入 |
| `application_data` | String | 完整申请数据JSON |
| `success` | String | 处理是否成功 |
| `error` | String | 错误信息 |

### AI风险分析节点
**输入变量**：
| 变量名 | 类型 | 引用值 |
|--------|------|--------|
| `application_data` | String | `{{解析申请数据.application_data}}` |
| `external_data` | String | `{{获取外部数据.text}}` |
| `credit_score` | Number | `{{解析申请数据.credit_score}}` |
| `amount` | Number | `{{解析申请数据.amount}}` |
| `annual_income` | Number | `{{解析申请数据.annual_income}}` |

**输出变量**：
| 变量名 | 类型 | 描述 |
|--------|------|------|
| `decision` | String | 审批决策 |
| `risk_score` | Number | 风险评分(0-1) |
| `risk_level` | String | 风险等级 |
| `confidence` | Number | 置信度(0-1) |
| `risk_factors` | String | 风险因素JSON |
| `approval_amount` | Number | 建议批准金额 |
| `analysis_summary` | String | 分析摘要 |
| `success` | String | 处理是否成功 |
| `error` | String | 错误信息 |

### 结束节点输出
| 变量名 | 类型 | 引用值 |
|--------|------|--------|
| `decision` | String | `{{AI风险分析.decision}}` |
| `risk_score` | Number | `{{AI风险分析.risk_score}}` |
| `risk_level` | String | `{{AI风险分析.risk_level}}` |
| `analysis_summary` | String | `{{AI风险分析.analysis_summary}}` |
| `approval_amount` | Number | `{{AI风险分析.approval_amount}}` |
| `confidence` | Number | `{{AI风险分析.confidence}}` |
| `processing_status` | String | `completed` |
| `application_id` | String | `{{start.application_id}}` |

**✅ 配置检查清单**：
- [ ] 所有变量名称正确拼写（区分大小写）
- [ ] Number类型变量不包含引号
- [ ] String类型变量包含双引号
- [ ] 变量引用格式为 `{{节点名.变量名}}`
- [ ] 所有必需的输出变量都已添加

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

3. **观察执行过程和预期结果**：

#### 节点执行顺序和预期输出：

**1. 获取申请信息节点**
- **预期状态**：✅ 成功
- **预期输出**：包含申请数据的JSON响应
- **检查要点**：`code: 0` 表示成功

**2. 解析申请数据节点** 
- **预期状态**：✅ 成功
- **预期输出**：
  ```json
  {
    "user_id": "user_001",
    "amount": 100000.0,
    "credit_score": 750.0,
    "annual_income": 200000.0,
    "success": "True"
  }
  ```
- **检查要点**：所有数值类型正确转换

**3. 获取外部数据节点**
- **预期状态**：✅ 成功
- **预期输出**：外部数据API响应
- **检查要点**：返回征信、银行流水等数据

**4. AI风险分析节点**
- **预期状态**：✅ 成功  
- **预期输出**：
  ```json
  {
    "decision": "AUTO_APPROVED",
    "risk_score": 0.2,
    "risk_level": "LOW",
    "confidence": 0.88,
    "approval_amount": 100000.0,
    "analysis_summary": "风险评分: 0.2, 决策: AUTO_APPROVED, 风险等级: LOW"
  }
  ```
- **检查要点**：决策逻辑正确，风险评分合理

**5. 提交AI决策节点**
- **预期状态**：✅ 成功
- **预期输出**：`{"code": 0, "message": "success"}`
- **检查要点**：成功提交到后端系统

**6. 结束节点**
- **预期状态**：✅ 成功
- **最终输出**：包含所有关键决策信息的完整结果

#### 🔍 故障排查指南：

**如果某个节点失败**：
1. 查看节点的详细错误信息
2. 检查输入变量的值是否正确
3. 验证变量引用格式
4. 确认后端服务正常运行

**常见错误及解决方案**：
- **变量未定义**：检查变量名拼写和引用格式
- **类型错误**：确认Number类型不包含引号  
- **API调用失败**：验证Token和网络连接
- **JSON解析错误**：检查代码节点的错误处理逻辑

4. 检查每个节点的输出数据格式和内容
5. 验证最终结果是否符合预期

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