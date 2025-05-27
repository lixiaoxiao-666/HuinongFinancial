# API 定义文档 (API_Spec.md)

## 1. 概述

本文档定义了"数字惠农APP及OA后台管理系统"后端服务的主要API接口。所有API都应通过API网关 (Higress AI 网关) 进行访问。

**通用约定:**

*   **Base URL**: `https://api.digital-agriculture.com/v1` (示例)
*   **认证**: 大部分接口需要Token认证 (JWT)，在HTTP Header中传递 `Authorization: Bearer <token>`。
*   **请求格式**: JSON (`Content-Type: application/json`)
*   **响应格式**: JSON
*   **成功响应**: HTTP状态码 `200 OK` 或 `201 Created`。响应体包含 `code`, `message`, `data` 字段。
    ```json
    {
      "code": 0, // 0表示成功，其他表示错误码
      "message": "Success",
      "data": { ... } // 业务数据
    }
    ```
*   **错误响应**: HTTP状态码 `4xx` (客户端错误) 或 `5xx` (服务器错误)。响应体包含 `code`, `message`, `error_details` (可选)。
    ```json
    {
      "code": 1001, // 具体错误码
      "message": "Invalid input parameters",
      "error_details": "Field 'username' is required."
    }
    ```
*   **分页**: 对于列表接口，使用 `page` 和 `limit` 参数进行分页，响应中包含 `total` 记录总数。

## 2. 用户服务 (UserService)

### 2.1 用户注册

*   **URL**: `/users/register`
*   **Method**: `POST`
*   **Request Body**:
    ```json
    {
      "phone": "13800138000", // 手机号
      "password": "your_password", // 密码 (前端应加密传输)
      "verification_code": "123456" // 短信验证码 (如需要)
    }
    ```
*   **Response (201 Created)**:
    ```json
    {
      "code": 0,
      "message": "User registered successfully",
      "data": {
        "user_id": "user_uuid_123"
      }
    }
    ```

### 2.2 用户登录

*   **URL**: `/users/login`
*   **Method**: `POST`
*   **Request Body**:
    ```json
    {
      "phone": "13800138000",
      "password": "your_password"
    }
    ```
*   **Response (200 OK)**:
    ```json
    {
      "code": 0,
      "message": "Login successful",
      "data": {
        "user_id": "user_uuid_123",
        "token": "jwt_auth_token_string",
        "expires_in": 7200
      }
    }
    ```

### 2.3 获取用户基本信息

*   **URL**: `/users/me`
*   **Method**: `GET`
*   **Authentication**: Required
*   **Response (200 OK)**:
    ```json
    {
      "code": 0,
      "message": "Success",
      "data": {
        "user_id": "user_uuid_123",
        "phone": "138****8000", // 脱敏显示
        "nickname": "农户小张",
        "avatar_url": "https://example.com/avatar.jpg"
        // ... 其他基本信息
      }
    }
    ```

### 2.4 更新用户基本信息

*   **URL**: `/users/me`
*   **Method**: `PUT`
*   **Authentication**: Required
*   **Request Body**:
    ```json
    {
      "nickname": "农户大张",
      "avatar_url": "https://example.com/new_avatar.jpg"
      // ... 其他可修改信息
    }
    ```
*   **Response (200 OK)**:
    ```json
    {
      "code": 0,
      "message": "User profile updated successfully"
    }
    ```

## 3. 惠农贷服务 (LoanService)

### 3.1 获取贷款产品列表

*   **URL**: `/loans/products`
*   **Method**: `GET`
*   **Query Parameters**:
    *   `category` (string, optional): 产品分类
*   **Response (200 OK)**:
    ```json
    {
      "code": 0,
      "message": "Success",
      "data": [
        {
          "product_id": "loan_prod_001",
          "name": "春耕助力贷",
          "description": "专为春耕生产设计，利率优惠",
          "min_amount": 5000,
          "max_amount": 50000,
          "min_term_months": 6,
          "max_term_months": 24,
          "interest_rate_yearly": "4.5% - 6.0%"
        }
        // ...更多产品
      ]
    }
    ```

### 3.2 获取贷款产品详情

*   **URL**: `/loans/products/{product_id}`
*   **Method**: `GET`
*   **Response (200 OK)**:
    ```json
    {
      "code": 0,
      "message": "Success",
      "data": {
        "product_id": "loan_prod_001",
        "name": "春耕助力贷",
        // ...更详细的产品信息，申请条件，所需材料列表等
      }
    }
    ```

### 3.3 提交贷款申请

*   **URL**: `/loans/applications`
*   **Method**: `POST`
*   **Authentication**: Required
*   **Request Body**:
    ```json
    {
      "product_id": "loan_prod_001",
      "amount": 30000, // 申请金额
      "term_months": 12, // 申请期限 (月)
      "purpose": "购买化肥和种子",
      "applicant_info": { // 申请人信息，部分可从用户服务获取，此处为补充
        "real_name": "张三",
        "id_card_number": "310...",
        "address": "XX省XX市XX村"
      },
      "uploaded_documents": [
        {"doc_type": "id_card_front", "file_id": "file_uuid_001"},
        {"doc_type": "land_contract", "file_id": "file_uuid_002"}
      ]
    }
    ```
*   **Response (201 Created)**:
    ```json
    {
      "code": 0,
      "message": "Loan application submitted successfully",
      "data": {
        "application_id": "loan_app_uuid_789"
      }
    }
    ```

### 3.4 查询贷款申请进度

*   **URL**: `/loans/applications/{application_id}`
*   **Method**: `GET`
*   **Authentication**: Required (申请人或有权限的管理员)
*   **Response (200 OK)**:
    ```json
    {
      "code": 0,
      "message": "Success",
      "data": {
        "application_id": "loan_app_uuid_789",
        "status": "AI_审批中", // 例如: 待提交, AI_审批中, 待人工复核, 已批准, 已拒绝
        "submitted_at": "2024-03-10T10:00:00Z",
        "updated_at": "2024-03-10T11:30:00Z",
        "approved_amount": null, // 批准后填充
        "remarks": "AI系统正在分析您的申请信息",
        "history": [
          {"status": "已提交", "timestamp": "2024-03-10T10:00:00Z", "operator": "用户"},
          {"status": "AI_审批中", "timestamp": "2024-03-10T10:05:00Z", "operator": "系统"}
        ]
      }
    }
    ```

### 3.5 获取我的贷款申请列表

*   **URL**: `/loans/applications/my`
*   **Method**: `GET`
*   **Authentication**: Required
*   **Query Parameters**:
    *   `status` (string, optional): 按状态筛选
    *   `page` (int, optional, default: 1)
    *   `limit` (int, optional, default: 10)
*   **Response (200 OK)**: (结构类似 3.1，包含申请列表和total)

### 3.6 文件上传接口 (通用)

*   **URL**: `/files/upload` (可能由专门的文件服务提供，或通过网关代理)
*   **Method**: `POST`
*   **Authentication**: Required
*   **Request Body**: `multipart/form-data` 包含文件本身
    *   `file`: (file)
    *   `purpose`: (string, optional, e.g., "loan_document", "machinery_image")
*   **Response (200 OK)**:
    ```json
    {
      "code": 0,
      "message": "File uploaded successfully",
      "data": {
        "file_id": "file_uuid_xyz",
        "file_url": "https://cdn.example.com/files/file_uuid_xyz.jpg",
        "file_name": "身份证正面.jpg",
        "file_size": 102400 // bytes
      }
    }
    ```

## 4. 农机租赁服务 (MachineryService)

### 4.1 发布农机信息

*   **URL**: `/machinery`
*   **Method**: `POST`
*   **Authentication**: Required (农机主)
*   **Request Body**:
    ```json
    {
      "type": "拖拉机",
      "brand_model": "东方红 LX904",
      "description": "90马力，四轮驱动，附带旋耕机",
      "images": ["file_uuid_img1", "file_uuid_img2"],
      "daily_rent": 300,
      "deposit": 1000,
      "location": {"latitude": 39.9, "longitude": 116.3, "address_text": "XX农场"},
      "availability_dates": [{"start_date": "2024-04-01", "end_date": "2024-04-15"}]
    }
    ```
*   **Response (201 Created)**:
    ```json
    {
      "code": 0,
      "message": "Machinery published successfully",
      "data": {
        "machinery_id": "mach_uuid_abc"
      }
    }
    ```

### 4.2 获取农机列表 (搜索与浏览)

*   **URL**: `/machinery`
*   **Method**: `GET`
*   **Query Parameters**:
    *   `type` (string, optional): 农机类型
    *   `location_area` (string, optional): 区域 (如城市、区县)
    *   `min_rent` (int, optional)
    *   `max_rent` (int, optional)
    *   `available_date` (string, optional, YYYY-MM-DD): 查询某天可用的农机
    *   `page` (int, optional, default: 1)
    *   `limit` (int, optional, default: 10)
*   **Response (200 OK)**: (结构类似 3.1，包含农机列表和total)
    ```json
    {
      "code": 0,
      "message": "Success",
      "data": [
        {
          "machinery_id": "mach_uuid_abc",
          "type": "拖拉机",
          "brand_model": "东方红 LX904",
          "main_image_url": "https://example.com/tractor.jpg",
          "daily_rent": 300,
          "location_text": "XX农场",
          "owner_user_id": "user_uuid_owner"
        }
      ],
      "total": 50
    }
    ```

### 4.3 获取农机详情

*   **URL**: `/machinery/{machinery_id}`
*   **Method**: `GET`
*   **Response (200 OK)**: (包含农机完整信息，用户评价等)

### 4.4 创建租赁订单

*   **URL**: `/machinery/orders`
*   **Method**: `POST`
*   **Authentication**: Required (承租方)
*   **Request Body**:
    ```json
    {
      "machinery_id": "mach_uuid_abc",
      "start_date": "2024-04-05",
      "end_date": "2024-04-07",
      "lessee_notes": "希望能提前半天取机"
    }
    ```
*   **Response (201 Created)**:
    ```json
    {
      "code": 0,
      "message": "Leasing order created successfully, waiting for owner confirmation",
      "data": {
        "order_id": "order_uuid_123"
      }
    }
    ```

### 4.5 获取我的租赁订单 (承租方/出租方)

*   **URL**: `/machinery/orders/my`
*   **Method**: `GET`
*   **Authentication**: Required
*   **Query Parameters**:
    *   `role` (string, required, "lessee" or "lessor"): 查询角色
    *   `status` (string, optional): 订单状态 (待确认, 已确认, 进行中, 已完成, 已取消)
    *   `page`, `limit`
*   **Response (200 OK)**: (包含订单列表和total)

### 4.6 更新租赁订单状态 (出租方)

*   **URL**: `/machinery/orders/{order_id}/status`
*   **Method**: `PUT`
*   **Authentication**: Required (出租方或管理员)
*   **Request Body**:
    ```json
    {
      "status": "confirmed", // confirmed, rejected, completed, cancelled
      "reason": "同意租赁" // (可选，拒绝时填写)
    }
    ```
*   **Response (200 OK)**:
    ```json
    {
      "code": 0,
      "message": "Order status updated successfully"
    }
    ```

## 5. OA后台服务 (AdminService)

### 5.1 OA用户登录

*   **URL**: `/admin/login`
*   **Method**: `POST`
*   **Request Body**:
    ```json
    {
      "username": "admin_user",
      "password": "admin_password"
    }
    ```
*   **Response (200 OK)**: (类似用户登录，返回OA用户的token和角色信息)
    ```json
    {
      "code": 0,
      "message": "Admin login successful",
      "data": {
        "admin_user_id": "admin_uuid_xyz",
        "username": "admin_user",
        "role": "审批员", // 或 "系统管理员"
        "token": "admin_jwt_auth_token",
        "expires_in": 3600
      }
    }
    ```

### 5.2 获取待审批贷款申请列表 (审批看板)

*   **URL**: `/admin/loans/applications/pending`
*   **Method**: `GET`
*   **Authentication**: Required (审批员/管理员)
*   **Query Parameters**:
    *   `status_filter` (string, optional, e.g., "AI_审批中,待人工复核")
    *   `applicant_name` (string, optional)
    *   `application_id` (string, optional)
    *   `sort_by` (string, optional, e.g., "submission_time_desc")
    *   `page`, `limit`
*   **Response (200 OK)**:
    ```json
    {
      "code": 0,
      "message": "Success",
      "data": [
        {
          "application_id": "loan_app_uuid_789",
          "applicant_name": "张三",
          "amount": 30000,
          "status": "待人工复核",
          "submission_time": "2024-03-10T10:00:00Z",
          "ai_risk_score": 75, // AI给出的风险评分 (示例)
          "ai_suggestion": "建议人工复核，关注还款能力证明。"
        }
        // ... 更多待审批申请
      ],
      "total": 25
    }
    ```

### 5.3 获取贷款申请详情 (OA)

*   **URL**: `/admin/loans/applications/{application_id}`
*   **Method**: `GET`
*   **Authentication**: Required (审批员/管理员)
*   **Response (200 OK)**: (返回比用户端更详细的信息，包括完整的AI分析报告、用户所有提交材料的链接等)
    ```json
    {
        "code": 0,
        "message": "Success",
        "data": {
            // ... 申请基本信息 (同用户端查询结果) ...
            "applicant_details": {
                "real_name": "张三",
                "id_card_number": "310...",
                // ... 更多申请时提交的详细信息
            },
            "uploaded_documents_details": [
                {"doc_type": "id_card_front", "file_id": "file_uuid_001", "file_url": "...", "ocr_result": "..."},
                {"doc_type": "land_contract", "file_id": "file_uuid_002", "file_url": "..."}
            ],
            "ai_analysis_report": {
                "overall_risk_score": 75,
                "risk_factors": ["收入证明材料不足", "近期多笔小额借贷记录"],
                "data_verification_results": [
                    {"item": "身份证号真实性", "result": "通过"},
                    {"item": "手机号实名", "result": "通过"}
                ],
                "credit_score_details": { ... }, // 模拟外部征信或内部评分详情
                "suggestion": "建议人工复核，重点核实收入来源和近期负债情况。"
            },
            "mcp_tool_outputs": [ // MCP工具执行结果 (如有)
                {"tool_name": "数据库查询-历史贷款", "result": {"count": 1, "status": "逾期"}}
            ]
            // ... 审批历史记录
        }
    }
    ```

### 5.4 提交人工审批决策

*   **URL**: `/admin/loans/applications/{application_id}/review`
*   **Method**: `POST`
*   **Authentication**: Required (审批员)
*   **Request Body**:
    ```json
    {
      "decision": "approved", // approved, rejected, require_more_info
      "approved_amount": 25000, // (可选, 如果是approved且金额调整)
      "comments": "申请人信用良好，但考虑到当前负债，略微调整批准金额。",
      "required_info_details": null // (可选, 如果是require_more_info，说明需要补充什么)
    }
    ```
*   **Response (200 OK)**:
    ```json
    {
      "code": 0,
      "message": "Review decision submitted successfully"
    }
    ```

### 5.5 控制AI审批流程开关 (示例)

*   **URL**: `/admin/system/ai-approval/toggle`
*   **Method**: `POST`
*   **Authentication**: Required (系统管理员)
*   **Request Body**:
    ```json
    {
      "enabled": false // true 或 false
    }
    ```
*   **Response (200 OK)**:
    ```json
    {
      "code": 0,
      "message": "AI approval process status updated successfully"
    }
    ```

## 6. Dify及本地AI模型调用 (内部，通过Higress AI网关)

这部分API通常是后端服务内部调用，由Higress AI网关统一管理和暴露。具体接口格式取决于Dify智能体和本地模型的定义。

**示例：调用Dify智能体进行贷款申请预分析**

*   **Higress 路由**: `/ai/dify/loan_pre_analysis` (内部服务地址映射到Dify智能体API)
*   **Method**: `POST`
*   **Request Body (由LoanService或AdminService构造)**:
    ```json
    {
      "application_data": { // 用户提交的申请数据
        "amount": 30000,
        "purpose": "购买化肥",
        "description": "今年计划扩大种植面积，需要购买一批优质化肥。有三年种植经验。"
        // ...其他相关字段
      }
    }
    ```
*   **Response (由Dify智能体返回，通过Higress透传)**:
    ```json
    {
      "analysis_id": "dify_analysis_xyz",
      "risk_level_suggestion": "medium", // low, medium, high
      "key_info_extracted": {
        "loan_purpose_category": "农资购买",
        "experience_years": 3
      },
      "potential_risks": ["未提供明确的还款来源说明"],
      "next_step_suggestion": "建议人工关注还款能力证明"
    }
    ```

**注意**: 上述API定义为初步设计，具体字段和结构可能在详细设计阶段进行调整。错误码和消息文本也需要进一步细化。文件上传接口可能需要更复杂的处理，如分片上传、断点续传等，此处为简化版本。 