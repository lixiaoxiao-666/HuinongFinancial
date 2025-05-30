# 数字惠农系统 - 前端 API 接口文档

## 1. 文档概述

### 1.1 文档目的
本文档旨在将《数字惠农系统 - 后端 API 接口文档》 (`agent/backend/API.md`) 中定义的后端接口，转化为前端开发友好、易于理解和使用的前端API接口规范。它将作为前端调用后端服务的主要参考，并指导前端API请求模块的封装。

### 1.2 基础约定
-   **请求基地址 (Base URL)**: 根据不同环境配置 (开发/测试/生产)，例如 `https://api.example.com/api/v1`。
-   **数据格式**: 前后端统一使用 JSON 数据格式进行交互。
-   **认证方式**: JWT Bearer Token。前端在调用需要认证的接口时，需在请求头中携带 `Authorization: Bearer {access_token}`。
-   **请求方法**: 遵循 RESTful 风格，使用 GET, POST, PUT, DELETE 等 HTTP 方法。
-   **响应结构**: 遵循后端定义的统一响应格式：
    ```json
    {
      "code": 200, // 业务状态码, 200 表示成功
      "message": "操作成功", // 提示信息
      "data": {}, // 实际响应数据
      "timestamp": 1640995200000, // 服务器时间戳
      "request_id": "req_123456789" // 请求ID
    }
    ```
-   **错误处理**: 前端需根据响应中的 `code` 字段判断业务是否成功。非200的 `code` 表示业务错误，`message` 字段通常包含错误信息。HTTP状态码如401, 403, 404, 500等也需要妥善处理。
-   **类型定义**: 前端将使用 TypeScript 定义请求参数和响应数据的类型，这些类型定义通常存放在 `src/types/api.ts` 或类似文件中。

## 2. API 模块划分

前端API将按照业务模块进行划分，与后端接口模块保持一致，方便管理和调用。例如：`authApi.ts`, `userApi.ts`, `loanApi.ts`, `machineApi.ts`, `articleApi.ts`, `adminApi.ts` 等。

## 3. 核心API接口详情

以下列出前端常用的核心API接口，并提供前端视角下的封装建议和类型定义参考。完整的后端接口定义请参考 `agent/backend/API.md`。

### 3.1 认证授权 (Auth)

#### 3.1.1 用户注册
-   **后端接口**: `POST /auth/register`
-   **前端封装示例 (e.g., `src/api/auth.ts`)**:
    ```typescript
    import apiClient from './client'; // 假设的axios实例封装
    import { RegisterPayload, RegisterResponse } from '@/types/api';

    export const register = (payload: RegisterPayload): Promise<RegisterResponse> => {
      return apiClient.post('/auth/register', payload);
    };
    ```
-   **请求参数类型 (`RegisterPayload`)**:
    ```typescript
    interface RegisterPayload {
      phone: string;       // 手机号
      password: string;    // 密码
      code: string;        // 短信验证码
      user_type: string;   // 用户类型 (e.g., 'farmer')
    }
    ```
-   **响应数据类型 (`RegisterResponseData`)**:
    ```typescript
    interface RegisterResponseData {
      user_id: number;
      uuid: string;
    }
    // RegisterResponse 为包含 code, message, data (类型为 RegisterResponseData) 的标准响应体
    ```

#### 3.1.2 用户登录
-   **后端接口**: `POST /auth/login`
-   **前端封装示例**:
    ```typescript
    import { LoginPayload, LoginResponse } from '@/types/api';

    export const login = (payload: LoginPayload): Promise<LoginResponse> => {
      return apiClient.post('/auth/login', payload);
    };
    ```
-   **请求参数类型 (`LoginPayload`)**:
    ```typescript
    interface LoginPayload {
      phone: string;
      password: string;
      platform?: string;    // e.g., 'app', 'web'
      device_id?: string;
      device_type?: string; // e.g., 'android', 'ios'
    }
    ```
-   **响应数据类型 (`LoginResponseData`)**:
    ```typescript
    interface UserInfo {
      id: number;
      uuid: string;
      phone: string;
      user_type: string;
      real_name?: string;
      avatar?: string;
      // ... 其他用户信息字段
    }

    interface LoginResponseData {
      access_token: string;
      refresh_token: string;
      expires_in: number;
      user_info: UserInfo;
    }
    ```

#### 3.1.3 刷新Token
-   **后端接口**: `POST /auth/refresh`
-   **前端封装示例**:
    ```typescript
    import { RefreshTokenPayload, RefreshTokenResponse } from '@/types/api';

    export const refreshToken = (payload: RefreshTokenPayload): Promise<RefreshTokenResponse> => {
      return apiClient.post('/auth/refresh', payload);
    };
    ```
-   **请求参数类型 (`RefreshTokenPayload`)**:
    ```typescript
    interface RefreshTokenPayload {
      refresh_token: string;
    }
    ```
-   **响应数据类型 (`RefreshTokenResponseData`)**: (通常包含新的 `access_token`, `refresh_token`, `expires_in`)
    ```typescript
    interface RefreshTokenResponseData {
      access_token: string;
      refresh_token: string;
      expires_in: number;
    }
    ```

### 3.2 用户管理 (User)

#### 3.2.1 获取当前用户信息
-   **后端接口**: `GET /user/profile`
-   **前端封装示例**:
    ```typescript
    import { UserProfileResponse } from '@/types/api';

    export const getUserProfile = (): Promise<UserProfileResponse> => {
      return apiClient.get('/user/profile'); // Authorization header 会由axios拦截器自动添加
    };
    ```
-   **响应数据类型 (`UserProfileResponseData`)**: (与 `LoginResponseData` 中的 `UserInfo` 类似，但更完整)
    ```typescript
    interface UserProfileResponseData {
      id: number;
      uuid: string;
      username?: string;
      phone: string;
      email?: string;
      user_type: string;
      status: string;
      real_name?: string;
      id_card?: string;
      avatar?: string;
      gender?: string;
      birthday?: string;
      province?: string;
      city?: string;
      county?: string;
      address?: string;
      is_real_name_verified: boolean;
      is_bank_card_verified: boolean;
      is_credit_verified: boolean;
      balance?: number; // 注意单位，如分
      credit_score?: number;
      credit_level?: string;
      created_at: string;
    }
    ```

#### 3.2.2 更新用户信息
-   **后端接口**: `PUT /user/profile`
-   **前端封装示例**:
    ```typescript
    import { UpdateUserProfilePayload, BaseApiResponse } from '@/types/api';

    export const updateUserProfile = (payload: UpdateUserProfilePayload): Promise<BaseApiResponse> => {
      return apiClient.put('/user/profile', payload);
    };
    ```
-   **请求参数类型 (`UpdateUserProfilePayload`)**: (只包含可更新的字段)
    ```typescript
    interface UpdateUserProfilePayload {
      real_name?: string;
      email?: string;
      gender?: string;
      birthday?: string;
      province?: string;
      city?: string;
      county?: string;
      address?: string;
      avatar?: string; // 通常头像更新是单独的文件上传接口
    }
    ```
-   **响应数据类型 (`BaseApiResponse`)**: 通用成功响应，data可能为空对象或包含更新后的部分信息。

#### 3.2.3 实名认证
-   **后端接口**: `POST /user/auth/realname`
-   **前端封装示例**:
    ```typescript
    import { RealNameAuthPayload, BaseApiResponse } from '@/types/api';

    export const submitRealNameAuth = (payload: RealNameAuthPayload): Promise<BaseApiResponse> => {
      return apiClient.post('/user/auth/realname', payload);
    };
    ```
-   **请求参数类型 (`RealNameAuthPayload`)**:
    ```typescript
    interface RealNameAuthPayload {
      real_name: string;
      id_card_number: string;
      id_card_front_img: string; // 图片URL，通过文件上传接口获取
      id_card_back_img: string;  // 图片URL
      face_verify_img?: string; // 图片URL (如果需要)
    }
    ```

### 3.3 贷款服务 (Loan)

#### 3.3.1 获取贷款产品列表
-   **后端接口**: `GET /loans/products`
-   **前端封装示例**:
    ```typescript
    import { LoanProductListParams, LoanProductListResponse } from '@/types/api';

    export const getLoanProducts = (params?: LoanProductListParams): Promise<LoanProductListResponse> => {
      return apiClient.get('/loans/products', { params });
    };
    ```
-   **请求查询参数类型 (`LoanProductListParams`)**: (可选)
    ```typescript
    interface LoanProductListParams {
      product_type?: string;
      user_type?: string;
    }
    ```
-   **响应数据类型 (`LoanProduct`, `LoanProductListResponseData`)**:
    ```typescript
    interface LoanProduct {
      id: number;
      product_code: string;
      product_name: string;
      description: string;
      product_type: string;
      min_amount: number; // 注意单位，如分
      max_amount: number;
      min_term: number; // 天
      max_term: number;
      interest_rate: number; // 年化利率，如 0.12 表示 12%
      interest_type: string; // 'fixed'
      repayment_method: string; // 'equal_installment'
      partner_name?: string;
      status: string; // 'active'
    }

    type LoanProductListResponseData = LoanProduct[];
    ```

#### 3.3.2 获取产品详情
-   **后端接口**: `GET /loans/products/{product_id}`
-   **前端封装示例**:
    ```typescript
    import { LoanProductDetailResponse } from '@/types/api';

    export const getLoanProductDetail = (productId: number): Promise<LoanProductDetailResponse> => {
      return apiClient.get(`/loans/products/${productId}`);
    };
    ```
-   **响应数据类型 (`LoanProductDetailResponseData`)**: (扩展自 `LoanProduct`)
    ```typescript
    interface EligibilityCriteria {
      min_age?: number;
      max_age?: number;
      required_credit_score?: number;
      required_documents?: string[];
    }

    interface LoanProductDetailResponseData extends LoanProduct {
      eligibility_criteria?: EligibilityCriteria;
      required_documents?: string[];
      applicable_user_types?: string[];
    }
    ```

#### 3.3.3 提交贷款申请
-   **后端接口**: `POST /loans/applications`
-   **前端封装示例**:
    ```typescript
    import { LoanApplicationPayload, LoanApplicationResponse } from '@/types/api';

    export const submitLoanApplication = (payload: LoanApplicationPayload): Promise<LoanApplicationResponse> => {
      return apiClient.post('/loans/applications', payload);
    };
    ```
-   **请求参数类型 (`LoanApplicationPayload`)**:
    ```typescript
    interface ApplicantInfo {
      annual_income: number;
      land_area?: number;
      crop_types?: string[];
      // ... 其他申请人补充信息
    }

    interface UploadedDocument {
      type: string; // e.g., 'id_card_front', 'bank_statement'
      url: string;  // 文件URL，通过文件上传接口获取
    }

    interface LoanApplicationPayload {
      product_id: number;
      applied_amount: number; // 申请金额，注意单位
      applied_term: number;   // 申请期限 (天)
      purpose: string;
      applicant_info?: ApplicantInfo;
      uploaded_documents?: UploadedDocument[];
    }
    ```
-   **响应数据类型 (`LoanApplicationResponseData`)**:
    ```typescript
    interface LoanApplicationResponseData {
      application_id: number;
      application_no: string;
      status: string; // e.g., 'pending_ai'
      estimated_review_time?: string;
    }
    ```

#### 3.3.4 获取我的申请列表
-   **后端接口**: `GET /loans/applications`
-   **前端封装示例**:
    ```typescript
    import { MyLoanApplicationsParams, MyLoanApplicationsResponse } from '@/types/api';

    export const getMyLoanApplications = (params?: MyLoanApplicationsParams): Promise<MyLoanApplicationsResponse> => {
      return apiClient.get('/loans/applications', { params });
    };
    ```
-   **请求查询参数类型 (`MyLoanApplicationsParams`)**: (可选)
    ```typescript
    interface MyLoanApplicationsParams {
      status?: string;
      page?: number;
      limit?: number;
    }
    ```
-   **响应数据类型 (`LoanApplicationItem`, `MyLoanApplicationsResponseData`)**:
    ```typescript
    interface LoanApplicationItem {
      id: number;
      application_no: string;
      product_name: string;
      applied_amount: number;
      applied_term: number;
      status: string;
      status_text: string;
      approved_amount?: number;
      approved_rate?: number;
      submitted_at: string;
      approved_at?: string;
    }

    interface MyLoanApplicationsResponseData {
      applications: LoanApplicationItem[];
      total: number;
      page: number;
      limit: number;
    }
    ```

### 3.4 农机租赁 (Machine)

(结构类似贷款服务，列举部分)

#### 3.4.1 搜索附近农机
-   **后端接口**: `GET /machines/nearby`
-   **请求查询参数类型 (`SearchNearbyMachinesParams`)**:
    ```typescript
    interface SearchNearbyMachinesParams {
      longitude: number;
      latitude: number;
      radius?: number; // 公里，默认10
      machine_type?: string;
      page?: number;
      limit?: number;
    }
    ```
-   **响应数据类型 (`MachineItem`, `NearbyMachinesResponseData`)**:
    ```typescript
    interface MachineOwner {
      id: number;
      real_name: string;
      avatar?: string;
    }
    interface MachineItem {
      id: number;
      machine_code: string;
      machine_name: string;
      brand?: string;
      model?: string;
      machine_type: string;
      status: string; // 'available'
      images?: string[];
      hourly_rate?: number; // 分/小时
      daily_rate?: number;  // 分/天
      deposit_amount?: number; // 分
      province?: string;
      city?: string;
      county?: string;
      distance?: number; // 公里
      average_rating?: number;
      review_count?: number;
      owner?: MachineOwner;
    }
    interface NearbyMachinesResponseData {
      machines: MachineItem[];
      total: number;
      page: number;
      limit: number;
    }
    ```

### 3.5 内容管理 (Content / Article)

#### 3.5.1 获取文章列表
-   **后端接口**: `GET /articles` (或 `/api/content/articles` 参考PRD)
-   **请求查询参数类型 (`ArticleListParams`)**: (可选)
    ```typescript
    interface ArticleListParams {
      category?: string;
      keyword?: string;
      page?: number;
      limit?: number;
    }
    ```
-   **响应数据类型 (`ArticleSummary`, `ArticleListResponseData`)**:
    ```typescript
    interface ArticleSummary {
      id: number;
      title: string;
      summary: string;
      cover_image?: string;
      category: string;
      author?: string;
      view_count: number;
      like_count?: number;
      is_top?: boolean;
      is_featured?: boolean;
      published_at: string;
    }
    interface ArticleListResponseData {
      articles: ArticleSummary[];
      total: number;
      page: number;
      limit: number;
    }
    ```

### 3.6 系统功能 (System)

#### 3.6.1 文件上传
-   **后端接口**: `POST /files/upload`
-   **前端封装注意**: 请求参数为 `multipart/form-data`。前端通常使用 `FormData` 对象来构建请求体。
    ```typescript
    // payload: FormData containing 'file', 'business_type', 'business_id' (optional)
    // export const uploadFile = (payload: FormData): Promise<FileUploadResponse> => {
    //   return apiClient.post('/files/upload', payload, {
    //     headers: {
    //       'Content-Type': 'multipart/form-data',
    //     },
    //   });
    // };
    ```
-   **响应数据类型 (`FileUploadResponseData`)**:
    ```typescript
    interface FileUploadResponseData {
      file_id: number;
      file_url: string;
      file_name: string;
      file_size: number;
      file_type: string;
    }
    ```

### 3.7 OA后台管理 (Admin)

OA后台的API调用与用户端类似，只是请求路径通常带有 `/admin/` 前缀，并且需要OA用户的认证Token。

#### 3.7.1 OA用户登录
-   **后端接口**: `POST /admin/auth/login`
-   **请求参数类型 (`AdminLoginPayload`)**:
    ```typescript
    interface AdminLoginPayload {
      username: string;
      password: string;
      captcha?: string;    // 如有验证码
      captcha_id?: string; // 验证码ID
    }
    ```
-   **响应数据类型 (`AdminLoginResponseData`)**: (类似于用户端登录，但可能包含OA用户特有信息和权限数据)

#### 3.7.2 获取用户列表 (OA)
-   **后端接口**: `GET /admin/users`
-   **请求查询参数类型 (`AdminUserListParams`)**:
    ```typescript
    interface AdminUserListParams {
      user_type?: string;
      status?: string;
      keyword?: string;
      page?: number;
      limit?: number;
    }
    ```
-   **响应数据类型 (`AdminUserItem`, `AdminUserListResponseData`)**: `AdminUserItem` 结构可能与 `UserProfileResponseData` 类似，但包含更多管理所需字段。

#### 3.7.3 审批贷款申请 (OA)
-   **后端接口**: `POST /admin/loans/applications/{application_id}/approve`
-   **请求参数类型 (`ApproveLoanPayload`)**:
    ```typescript
    interface ApproveLoanPayload {
      approved_amount: number;
      approved_term: number;
      approved_rate: number;
      comments?: string;
    }
    ```

## 4. API 请求封装与管理

-   **创建API Client实例**: 封装一个统一的API请求客户端 (如 `src/api/client.ts` 使用Axios)。
    -   设置基础URL (`baseURL`)。
    -   配置请求超时时间 (`timeout`)。
    -   **请求拦截器**: 自动添加 `Authorization` 请求头 (从状态管理或本地存储获取Token)。添加通用请求参数如 `request_id`。
    -   **响应拦截器**: 处理通用错误 (如401未授权时跳转登录页，Token过期时尝试刷新Token)，提取 `data` 字段，处理业务错误码等。
-   **API模块化**: 如第2节所述，按业务模块组织API调用函数。
-   **类型安全**: 所有请求参数和响应数据都应有明确的TypeScript类型定义。
-   **错误处理**: API调用函数应返回 `Promise`，并在调用处使用 `try...catch` 或 `.catch()` 处理可能发生的错误。
-   **状态管理集成**: 用户Token、用户信息等通常存储在状态管理库中，API请求时从中读取或更新。

## 5. 文档更新

当前端API封装或后端接口发生变更时，应及时更新本文档及相关的TypeScript类型定义，确保文档与实际代码一致。

本文档主要基于 `agent/backend/API.md` v1.0 版本进行前端化转换。后续后端接口版本迭代时，本文档也需要同步更新。 