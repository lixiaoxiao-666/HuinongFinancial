# 内容管理模块 - API 接口文档

## 📋 模块概述

内容管理模块负责农业资讯、政策解读、专家信息等内容的发布与管理。
部分接口为公开访问，部分管理接口需要OA管理员权限。

### 核心功能
-   **公共内容接口 (`/api/content/*`)**: 用户浏览文章、专家列表等 (可选认证)。
-   **OA管理员内容管理接口 (`/api/oa/admin/content/*`)**: 创建、编辑、删除、发布内容。

---

## 📰 公共内容接口 (惠农APP/Web)

**接口路径前缀**: `/api/content`
**认证要求**: `OptionalAuth` (可选认证 - 登录用户可能获取个性化推荐或额外功能)
**适用平台**: `app`, `web`, `oa` (OA用户也可以浏览公共内容)

### 1.1 获取文章列表

```http
GET /api/content/articles?category=tech&page=1&limit=10&sort=newest
Authorization: Bearer {access_token} // 可选
```

**Query Parameters**:
-   `category` (string, 可选): 文章分类代码 (如 `tech`, `policy`, `market`)
-   `tag` (string, 可选): 文章标签
-   `keyword` (string, 可选): 关键词搜索
-   `page`, `limit` (int, 可选): 分页
-   `sort` (string, 可选): `newest`, `popular`, `featured`

### 1.2 获取特色/推荐文章

```http
GET /api/content/articles/featured?count=5
Authorization: Bearer {access_token} // 可选
```

### 1.3 获取文章详情

```http
GET /api/content/articles/{article_id_or_slug}
Authorization: Bearer {access_token} // 可选
```

### 1.4 获取文章分类列表

```http
GET /api/content/categories
Authorization: Bearer {access_token} // 可选
```

### 1.5 获取专家列表

```http
GET /api/content/experts?expertise_area=crop_disease&page=1&limit=10
Authorization: Bearer {access_token} // 可选
```

### 1.6 获取专家详情

```http
GET /api/content/experts/{expert_id_or_username}
Authorization: Bearer {access_token} // 可选
```

### 1.7 用户提交专家咨询 (惠农APP/Web用户)

**认证要求**: `RequireAuth` (惠农APP/Web用户)
**适用平台**: `app`, `web`

```http
POST /api/user/consultations  // 注意：此接口放在 /api/user/ 下，更符合用户行为
Authorization: Bearer {access_token}
Content-Type: application/json

{
    "expert_id": 701,
    "title": "关于水稻稻瘟病的防治问题",
    "question_details": "我的水稻田出现了稻瘟病迹象，请问如何有效防治？附图...",
    "attachments": [
        {"file_name": "rice_disease_1.jpg", "file_url": "https://example.com/uploads/..."}
    ]
}
```

### 1.8 用户获取自己的咨询列表 (惠农APP/Web用户)

**认证要求**: `RequireAuth` (惠农APP/Web用户)
**适用平台**: `app`, `web`

```http
GET /api/user/consultations?status=pending_reply&page=1&limit=10
Authorization: Bearer {access_token}
```

---

## 🛠️ OA系统 - 内容管理接口 (管理员)

**接口路径前缀**: `/api/oa/admin/content`
**认证要求**: `RequireAuth`, `CheckPlatform("oa")`, `RequireRole("admin")`
**适用平台**: `oa`

### 2.1 文章管理 (管理员)

#### 2.1.1 创建文章

```http
POST /api/oa/admin/content/articles
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "title": "春季小麦田间管理要点",
    "content_markdown": "### 1. 肥水管理\n春季是小麦生长的关键时期...",
    "category_id": 10, // 分类ID
    "tags": ["小麦", "田间管理", "春季"],
    "author_id": 201, // OA系统用户ID (发布者)
    "status": "draft" // draft, published, archived
}
```

#### 2.1.2 更新文章

```http
PUT /api/oa/admin/content/articles/{article_id}
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "title": "春季小麦高效田间管理技术",
    "content_html": "<h1>春季小麦高效田间管理技术</h1><p>...</p>", // 可以支持HTML直接输入或由Markdown转换
    "status": "published"
}
```

#### 2.1.3 删除文章

```http
DELETE /api/oa/admin/content/articles/{article_id}
Authorization: Bearer {oa_access_token}
```

#### 2.1.4 发布/取消发布文章

```http
POST /api/oa/admin/content/articles/{article_id}/publish
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "publish": true // true为发布, false为取消发布 (变为草稿或归档)
}
```

### 2.2 分类管理 (管理员)

#### 2.2.1 创建分类

```http
POST /api/oa/admin/content/categories
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "name": "病虫害防治",
    "slug": "pest-control",
    "description": "各类农作物病虫害防治技术与资讯。",
    "parent_id": 5 // 可选，父分类ID
}
```

#### 2.2.2 更新分类

```http
PUT /api/oa/admin/content/categories/{category_id}
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "name": "常见病虫害防治技术"
}
```

#### 2.2.3 删除分类

```http
DELETE /api/oa/admin/content/categories/{category_id}
Authorization: Bearer {oa_access_token}
```

### 2.3 专家管理 (管理员)

#### 2.3.1 添加专家信息

```http
POST /api/oa/admin/content/experts
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "name": "王农艺师",
    "title": "高级农艺师",
    "expertise_areas": ["水稻种植", "病虫害防治"],
    "bio": "王农艺师拥有超过20年的水稻种植和病虫害防治经验...",
    "avatar_url": "https://example.com/experts/wang.jpg",
    "contact_phone": "13600136000", // 可选
    "oa_user_id": 205 // 如果专家也是OA系统用户，可以关联，方便回复咨询
}
```

#### 2.3.2 更新专家信息

```http
PUT /api/oa/admin/content/experts/{expert_id}
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "title": "资深高级农艺师",
    "status": "active" // active, inactive
}
```

#### 2.3.3 删除专家信息

```http
DELETE /api/oa/admin/content/experts/{expert_id}
Authorization: Bearer {oa_access_token}
```

### 2.4 专家咨询管理 (管理员)

#### 2.4.1 获取待回复/所有咨询列表

```http
GET /api/oa/admin/content/consultations?status=pending_reply&expert_id=701&page=1&limit=10
Authorization: Bearer {oa_access_token}
```

#### 2.4.2 回复用户咨询 (专家或管理员代回复)

```http
POST /api/oa/admin/content/consultations/{consultation_id}/reply
Authorization: Bearer {oa_access_token}
Content-Type: application/json

{
    "replier_id": 205, // 回复者OA用户ID (专家本人或管理员)
    "reply_content": "关于您提到的稻瘟病问题，建议采用以下防治措施：1... 2...",
    "is_public": false // 是否将此问答公开到知识库
}
```

---

**说明**: 请注意区分公共内容接口和OA管理员专用的管理接口。路径和参数请根据实际后端实现进行调整。 