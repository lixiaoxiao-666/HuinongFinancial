# 任务管理模块 - API 接口文档

## 📋 模块概述

任务管理模块为OA系统管理员提供待处理任务的管理功能，包括任务的创建、分配、处理、进度跟踪等。该模块主要用于管理业务流程中的待处理事项，如审批任务、处理任务等。

### 核心功能
-   **任务管理**: 创建、查看、更新、删除任务
-   **任务分配**: 将任务分配给特定的处理人员
-   **任务处理**: 处理任务并更新状态
-   **进度跟踪**: 查看任务处理进度和历史记录

---

## 🛠️ OA系统 - 任务管理接口 (管理员)

**接口路径前缀**: `/api/oa/admin/tasks`
**认证要求**: `RequireAuth`, `CheckPlatform("oa")`, `RequireRole("admin")`
**适用平台**: `oa`

### 1.1 获取任务列表

```http
GET /api/oa/admin/tasks?status=pending&assignee_id=201&page=1&limit=20
Authorization: Bearer {oa_access_token}
```

**Query Parameters**:
-   `status` (string, 可选): 任务状态筛选 (`pending`, `processing`, `completed`, `cancelled`)
-   `assignee_id` (uint64, 可选): 按任务分配人筛选
-   `creator_id` (uint64, 可选): 按任务创建人筛选
-   `task_type` (string, 可选): 任务类型 (`loan_review`, `auth_review`, `system_maintenance`)
-   `priority` (string, 可选): 优先级 (`low`, `normal`, `high`, `urgent`)
-   `date_range_start`, `date_range_end` (string, 可选): 按创建时间筛选
-   `page`, `limit` (int, 可选): 分页参数

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total": 45,
        "tasks": [
            {
                "id": 1001,
                "task_type": "loan_review",
                "title": "贷款申请审核 - LA20240115001",
                "description": "需要审核张三的农业创业贷申请",
                "status": "pending",
                "priority": "high",
                "creator": {
                    "id": 201,
                    "real_name": "管理员张",
                    "username": "admin_zhang"
                },
                "assignee": {
                    "id": 202,
                    "real_name": "审核员李",
                    "username": "reviewer_li"
                },
                "related_entity": {
                    "type": "loan_application",
                    "id": "LA20240115001"
                },
                "due_date": "2024-01-17T18:00:00Z",
                "created_at": "2024-01-15T09:00:00Z",
                "updated_at": "2024-01-15T09:00:00Z"
            }
        ]
    }
}
```

### 1.2 创建任务

```http
POST /api/oa/admin/tasks
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "task_type": "auth_review",
    "title": "实名认证审核 - 用户王五",
    "description": "需要审核用户王五提交的实名认证材料",
    "priority": "normal",
    "assignee_id": 203,
    "related_entity": {
        "type": "user_auth",
        "id": "auth_realname_uuid789"
    },
    "due_date": "2024-01-18T18:00:00Z",
    "metadata": {
        "user_id": 1005,
        "auth_type": "real_name"
    }
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "任务创建成功",
    "data": {
        "id": 1002,
        "task_type": "auth_review",
        "status": "pending",
        "created_at": "2024-01-15T10:30:00Z"
    }
}
```

### 1.3 获取任务详情

```http
GET /api/oa/admin/tasks/{task_id}
Authorization: Bearer {oa_access_token}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "id": 1001,
        "task_type": "loan_review",
        "title": "贷款申请审核 - LA20240115001",
        "description": "需要审核张三的农业创业贷申请",
        "status": "processing",
        "priority": "high",
        "progress": 50,
        "creator": {
            "id": 201,
            "real_name": "管理员张",
            "username": "admin_zhang"
        },
        "assignee": {
            "id": 202,
            "real_name": "审核员李",
            "username": "reviewer_li"
        },
        "related_entity": {
            "type": "loan_application",
            "id": "LA20240115001",
            "details": {
                "user_name": "张三",
                "amount": 100000,
                "product_name": "农业创业贷"
            }
        },
        "due_date": "2024-01-17T18:00:00Z",
        "started_at": "2024-01-15T14:00:00Z",
        "completed_at": null,
        "created_at": "2024-01-15T09:00:00Z",
        "updated_at": "2024-01-15T14:30:00Z",
        "logs": [
            {
                "id": 501,
                "action": "task_started",
                "operator": {
                    "id": 202,
                    "real_name": "审核员李"
                },
                "comment": "开始审核贷款申请",
                "created_at": "2024-01-15T14:00:00Z"
            }
        ]
    }
}
```

### 1.4 更新任务

```http
PUT /api/oa/admin/tasks/{task_id}
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "title": "贷款申请审核 - LA20240115001 (优先)",
    "description": "需要优先审核张三的农业创业贷申请",
    "priority": "urgent",
    "due_date": "2024-01-16T18:00:00Z"
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "任务更新成功",
    "data": {
        "updated_fields": ["title", "description", "priority", "due_date"],
        "updated_at": "2024-01-15T15:00:00Z"
    }
}
```

### 1.5 删除任务

```http
DELETE /api/oa/admin/tasks/{task_id}
Authorization: Bearer {oa_access_token}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "任务删除成功"
}
```

### 1.6 获取待处理任务列表

```http
GET /api/oa/admin/tasks/pending?assignee_id=202&limit=10
Authorization: Bearer {oa_access_token}
```

**说明**: 此接口为获取所有状态为 `pending` 的任务的快捷接口，响应格式与1.1相同。

---

## 🔄 任务处理操作接口

### 2.1 处理任务

```http
POST /api/oa/admin/tasks/{task_id}/process
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "action": "approve", // approve, reject, return, complete
    "comment": "审核通过，贷款申请符合要求",
    "result_data": {
        "approved_amount": 100000,
        "approved_term": 12,
        "interest_rate": 0.065
    },
    "next_step": "contract_generation" // 可选，指定下一步流程
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "任务处理成功",
    "data": {
        "task_id": 1001,
        "status": "completed",
        "action": "approve",
        "processed_at": "2024-01-15T16:00:00Z",
        "next_task_id": 1003 // 如果生成了后续任务
    }
}
```

### 2.2 分配任务

```http
POST /api/oa/admin/tasks/{task_id}/assign
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "assignee_id": 204,
    "comment": "分配给专业审核员"
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "任务分配成功",
    "data": {
        "task_id": 1001,
        "old_assignee_id": 202,
        "new_assignee_id": 204,
        "assigned_at": "2024-01-15T16:30:00Z"
    }
}
```

### 2.3 取消任务分配

```http
POST /api/oa/admin/tasks/{task_id}/unassign
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "comment": "重新安排处理人员"
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "任务取消分配成功",
    "data": {
        "task_id": 1001,
        "previous_assignee_id": 204,
        "unassigned_at": "2024-01-15T17:00:00Z"
    }
}
```

### 2.4 重新分配任务

```http
POST /api/oa/admin/tasks/{task_id}/reassign
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "new_assignee_id": 205,
    "comment": "重新分配给资深审核员"
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "任务重新分配成功",
    "data": {
        "task_id": 1001,
        "old_assignee_id": 204,
        "new_assignee_id": 205,
        "reassigned_at": "2024-01-15T17:30:00Z"
    }
}
```

### 2.5 获取任务进度

```http
GET /api/oa/admin/tasks/{task_id}/progress
Authorization: Bearer {oa_access_token}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "task_id": 1001,
        "current_status": "processing",
        "progress_percentage": 75,
        "steps": [
            {
                "step_name": "task_created",
                "step_title": "任务创建",
                "status": "completed",
                "completed_at": "2024-01-15T09:00:00Z"
            },
            {
                "step_name": "task_assigned",
                "step_title": "任务分配",
                "status": "completed",
                "completed_at": "2024-01-15T09:30:00Z"
            },
            {
                "step_name": "task_started",
                "step_title": "开始处理",
                "status": "completed",
                "completed_at": "2024-01-15T14:00:00Z"
            },
            {
                "step_name": "review_in_progress",
                "step_title": "审核进行中",
                "status": "in_progress",
                "started_at": "2024-01-15T14:00:00Z"
            },
            {
                "step_name": "task_completed",
                "step_title": "任务完成",
                "status": "pending",
                "estimated_completion": "2024-01-17T16:00:00Z"
            }
        ],
        "time_spent_minutes": 150,
        "estimated_remaining_minutes": 60
    }
}
```

---

## 🔧 错误码说明

| 错误码 | 说明 | 处理建议 |
|-------|------|---------|
| 5001 | 任务不存在 | 检查任务ID是否正确 |
| 5002 | 任务状态不允许操作 | 检查任务当前状态 |
| 5003 | 分配的用户不存在 | 检查用户ID是否有效 |
| 5004 | 任务类型无效 | 使用有效的任务类型 |
| 5005 | 优先级参数无效 | 检查优先级参数 |
| 5006 | 任务已被其他人处理 | 刷新任务状态后重试 |
| 5007 | 缺少必要的处理权限 | 确认用户权限 |
| 5008 | 截止日期无效 | 检查日期格式和逻辑 |
| 5009 | 关联实体不存在 | 检查关联的业务对象 |
| 5010 | 任务处理数据无效 | 检查处理结果数据 |

---

## 📝 接口调用示例

### JavaScript示例
```javascript
// 获取待处理任务列表
const getPendingTasks = async (token) => {
    const response = await fetch('/api/oa/admin/tasks/pending?limit=10', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};

// 处理任务
const processTask = async (token, taskId, processData) => {
    const response = await fetch(`/api/oa/admin/tasks/${taskId}/process`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(processData)
    });
    return response.json();
};

// 分配任务
const assignTask = async (token, taskId, assigneeId, comment) => {
    const response = await fetch(`/api/oa/admin/tasks/${taskId}/assign`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            assignee_id: assigneeId,
            comment: comment
        })
    });
    return response.json();
};

// 获取任务进度
const getTaskProgress = async (token, taskId) => {
    const response = await fetch(`/api/oa/admin/tasks/${taskId}/progress`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};
```

### 业务流程说明
1. **任务创建**: 系统自动创建或管理员手动创建任务
2. **任务分配**: 将任务分配给合适的处理人员
3. **任务处理**: 处理人员执行任务并记录处理结果
4. **进度跟踪**: 实时跟踪任务处理进度和状态变化
5. **任务完成**: 处理完成后更新任务状态并可能触发后续流程

### 注意事项
1. **权限控制**: 只有管理员可以创建、分配和查看所有任务
2. **状态流转**: 任务状态变化需要遵循预定义的流程
3. **日志记录**: 所有任务操作都会记录详细的操作日志
4. **截止时间**: 超过截止时间的任务会自动标记为延期
5. **关联业务**: 任务与具体的业务对象关联，便于快速定位和处理 