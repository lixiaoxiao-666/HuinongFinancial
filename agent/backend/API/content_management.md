# 内容管理模块 - API 接口文档

## 📋 模块概述

内容管理模块为系统提供丰富的农业资讯、政策信息、技术指导和专家咨询服务。支持多媒体内容发布、个性化推荐、互动交流等功能，为农户提供及时、准确、有价值的信息服务。

### 核心功能
- 📰 **资讯管理**: 农业新闻、技术资讯、市场行情
- 📋 **政策发布**: 政策解读、补贴申请、办事指南
- 👨‍🎓 **专家咨询**: 在线咨询、知识问答、技术指导
- 🔔 **通知管理**: 系统通知、个人消息、推送管理
- 🏷️ **内容分类**: 标签管理、分类体系、搜索优化

---

## 📰 资讯内容管理

### 1.1 获取资讯列表
```http
GET /api/content/articles?category=tech&tag=种植技术&page=1&limit=20&sort=latest
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**查询参数:**
- `category`: 资讯分类 (tech/market/policy/news)
- `tag`: 标签筛选
- `region`: 地区筛选
- `keyword`: 关键词搜索
- `sort`: 排序方式 (latest/popular/recommended)

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total": 156,
        "page": 1,
        "limit": 20,
        "categories": [
            {"code": "tech", "name": "技术资讯", "count": 45},
            {"code": "market", "name": "市场行情", "count": 32},
            {"code": "policy", "name": "政策资讯", "count": 28},
            {"code": "news", "name": "农业新闻", "count": 51}
        ],
        "hot_tags": [
            {"name": "病虫害防治", "count": 25},
            {"name": "智慧农业", "count": 18},
            {"name": "有机种植", "count": 15}
        ],
        "articles": [
            {
                "id": 30001,
                "title": "冬季小麦病虫害综合防治技术",
                "summary": "详细介绍冬季小麦常见病虫害的识别、预防和治疗方法",
                "content_type": "article",
                "category": "tech",
                "category_name": "技术资讯",
                "tags": ["病虫害防治", "小麦种植", "冬季管理"],
                "author": {
                    "id": 5001,
                    "name": "张教授",
                    "title": "农业技术专家",
                    "avatar": "https://example.com/avatars/expert1.jpg"
                },
                "cover_image": "https://example.com/articles/30001_cover.jpg",
                "images": [
                    "https://example.com/articles/30001_1.jpg",
                    "https://example.com/articles/30001_2.jpg"
                ],
                "view_count": 1250,
                "like_count": 89,
                "comment_count": 23,
                "share_count": 45,
                "is_featured": true,
                "is_recommended": true,
                "status": "published",
                "publish_time": "2024-01-15T09:00:00Z",
                "created_at": "2024-01-14T16:30:00Z",
                "region": ["山东省", "河南省", "河北省"]
            }
        ]
    }
}
```

### 1.2 获取资讯详情
```http
GET /api/content/articles/{article_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "id": 30001,
        "title": "冬季小麦病虫害综合防治技术",
        "summary": "详细介绍冬季小麦常见病虫害的识别、预防和治疗方法",
        "content": "冬季是小麦生长的关键时期，也是病虫害防治的重要阶段...",
        "content_html": "<p>冬季是小麦生长的关键时期...</p>",
        "content_type": "article",
        "category": "tech",
        "category_name": "技术资讯",
        "tags": ["病虫害防治", "小麦种植", "冬季管理"],
        "author": {
            "id": 5001,
            "name": "张教授",
            "title": "农业技术专家",
            "bio": "从事农业技术研究20年，专注作物病虫害防治",
            "avatar": "https://example.com/avatars/expert1.jpg",
            "expertise": ["病虫害防治", "作物栽培", "土壤改良"]
        },
        "cover_image": "https://example.com/articles/30001_cover.jpg",
        "images": [
            {
                "url": "https://example.com/articles/30001_1.jpg",
                "description": "病虫害症状图"
            },
            {
                "url": "https://example.com/articles/30001_2.jpg",
                "description": "防治效果对比"
            }
        ],
        "attachments": [
            {
                "name": "防治技术手册.pdf",
                "url": "https://example.com/files/handbook.pdf",
                "size": 2048000
            }
        ],
        "statistics": {
            "view_count": 1250,
            "like_count": 89,
            "comment_count": 23,
            "share_count": 45,
            "collect_count": 67
        },
        "user_interaction": {
            "is_liked": false,
            "is_collected": true,
            "view_time": "2024-01-15T14:30:00Z"
        },
        "related_articles": [
            {
                "id": 30002,
                "title": "春季小麦田间管理要点",
                "cover_image": "https://example.com/articles/30002_cover.jpg"
            }
        ],
        "region": ["山东省", "河南省", "河北省"],
        "publish_time": "2024-01-15T09:00:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
    }
}
```

### 1.3 搜索资讯
```http
GET /api/content/articles/search?q=病虫害防治&region=山东省&category=tech
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 1.4 获取推荐资讯
```http
GET /api/content/articles/recommended?user_location=济南市&user_interest=种植技术
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "recommendation_reason": "基于您的位置和兴趣推荐",
        "articles": [
            {
                "id": 30003,
                "title": "济南地区冬小麦种植技术指导",
                "reason": "地区相关",
                "match_score": 0.95
            }
        ]
    }
}
```

---

## 📋 政策信息管理

### 2.1 获取政策列表
```http
GET /api/content/policies?type=subsidy&region=山东省&status=active&page=1&limit=20
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total": 85,
        "page": 1,
        "limit": 20,
        "categories": [
            {"code": "subsidy", "name": "补贴政策", "count": 35},
            {"code": "finance", "name": "金融政策", "count": 20},
            {"code": "land", "name": "土地政策", "count": 18},
            {"code": "insurance", "name": "保险政策", "count": 12}
        ],
        "policies": [
            {
                "id": 40001,
                "title": "2024年农机购置补贴实施方案",
                "summary": "针对农机购置提供30-50%的财政补贴支持",
                "policy_type": "subsidy",
                "policy_type_name": "补贴政策",
                "issuer": "山东省农业农村厅",
                "policy_number": "鲁农机发〔2024〕1号",
                "issue_date": "2024-01-01",
                "effective_date": "2024-01-01",
                "expiry_date": "2024-12-31",
                "status": "active",
                "status_text": "有效",
                "region": ["山东省"],
                "target_group": ["个体农户", "农民合作社", "家庭农场"],
                "subsidy_rate": "30-50%",
                "max_amount": 50000,
                "application_deadline": "2024-11-30",
                "view_count": 2580,
                "application_count": 156,
                "is_featured": true,
                "created_at": "2024-01-01T08:00:00Z"
            }
        ]
    }
}
```

### 2.2 获取政策详情
```http
GET /api/content/policies/{policy_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "id": 40001,
        "title": "2024年农机购置补贴实施方案",
        "summary": "针对农机购置提供30-50%的财政补贴支持",
        "content": "为加快推进农业机械化发展，提高农业生产效率...",
        "policy_type": "subsidy",
        "issuer": "山东省农业农村厅",
        "policy_number": "鲁农机发〔2024〕1号",
        "issue_date": "2024-01-01",
        "effective_date": "2024-01-01",
        "expiry_date": "2024-12-31",
        "region": ["山东省"],
        "target_group": ["个体农户", "农民合作社", "家庭农场"],
        "eligibility_criteria": [
            "在山东省内从事农业生产的个人或组织",
            "购置符合补贴目录的农机产品",
            "遵守国家和省有关法律法规"
        ],
        "subsidy_details": {
            "subsidy_rate": "30-50%",
            "max_amount": 50000,
            "calculation_method": "按设备购置价格的一定比例计算",
            "payment_method": "财政直接拨付"
        },
        "application_process": [
            {
                "step": 1,
                "title": "提交申请",
                "description": "在线填写申请表并上传相关材料",
                "required_materials": ["身份证明", "购机发票", "银行账户信息"]
            },
            {
                "step": 2,
                "title": "资格审核",
                "description": "相关部门审核申请资格",
                "duration": "5-10个工作日"
            },
            {
                "step": 3,
                "title": "补贴发放",
                "description": "审核通过后发放补贴资金",
                "duration": "15-20个工作日"
            }
        ],
        "required_materials": [
            {
                "name": "身份证明",
                "description": "申请人身份证或营业执照",
                "format": "PDF/JPG",
                "required": true
            },
            {
                "name": "购机发票",
                "description": "农机购置发票原件",
                "format": "PDF/JPG",
                "required": true
            }
        ],
        "application_deadline": "2024-11-30",
        "contact_info": {
            "department": "山东省农机推广站",
            "phone": "0531-12345678",
            "email": "njbz@shandong.gov.cn",
            "address": "济南市历下区经十路12345号"
        },
        "faq": [
            {
                "question": "补贴资金何时到账？",
                "answer": "审核通过后15-20个工作日内到账"
            }
        ],
        "related_policies": [
            {
                "id": 40002,
                "title": "农业保险补贴政策"
            }
        ],
        "statistics": {
            "view_count": 2580,
            "application_count": 156,
            "approval_rate": 0.85
        }
    }
}
```

### 2.3 政策申请
```http
POST /api/content/policies/{policy_id}/apply
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "applicant_type": "individual",
    "purchase_info": {
        "equipment_name": "约翰迪尔拖拉机",
        "model": "6B-1204",
        "purchase_price": 120000,
        "purchase_date": "2024-01-10",
        "dealer_name": "农机销售公司"
    },
    "materials": [
        {
            "type": "id_card",
            "file_url": "https://example.com/files/id_card.pdf"
        },
        {
            "type": "invoice",
            "file_url": "https://example.com/files/invoice.pdf"
        }
    ],
    "bank_account": {
        "account_number": "6226090000000001",
        "bank_name": "中国工商银行",
        "account_holder": "张三"
    }
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "申请提交成功",
    "data": {
        "application_id": "PA20240115001",
        "policy_id": 40001,
        "status": "submitted",
        "estimated_subsidy": 36000,
        "review_deadline": "2024-01-25",
        "tracking_number": "PA20240115001"
    }
}
```

---

## 👨‍🎓 专家咨询管理

### 3.1 获取专家列表
```http
GET /api/content/experts?expertise=病虫害防治&region=山东省&available=true&page=1&limit=20
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total": 45,
        "page": 1,
        "limit": 20,
        "expertise_areas": [
            {"name": "病虫害防治", "count": 12},
            {"name": "土壤改良", "count": 8},
            {"name": "作物栽培", "count": 15},
            {"name": "农机使用", "count": 10}
        ],
        "experts": [
            {
                "id": 5001,
                "name": "张教授",
                "title": "农业技术专家",
                "organization": "山东农业大学",
                "expertise": ["病虫害防治", "作物栽培", "土壤改良"],
                "bio": "从事农业技术研究20年，专注作物病虫害防治和绿色种植技术",
                "avatar": "https://example.com/avatars/expert1.jpg",
                "rating": 4.8,
                "consultation_count": 156,
                "response_rate": 0.95,
                "avg_response_time": 2.5,
                "online_status": "available",
                "consultation_fee": 0,
                "service_region": ["山东省", "河南省"],
                "languages": ["中文"],
                "availability": {
                    "mon": ["09:00-12:00", "14:00-17:00"],
                    "tue": ["09:00-12:00", "14:00-17:00"],
                    "wed": ["09:00-12:00"],
                    "thu": ["09:00-12:00", "14:00-17:00"],
                    "fri": ["09:00-12:00", "14:00-17:00"]
                }
            }
        ]
    }
}
```

### 3.2 获取专家详情
```http
GET /api/content/experts/{expert_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 3.3 提交咨询问题
```http
POST /api/content/consultations
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "expert_id": 5001,
    "consultation_type": "text",
    "subject": "小麦叶片出现黄斑如何处理",
    "description": "我家的冬小麦最近叶片出现黄色斑点，请问这是什么病害，如何防治？",
    "category": "病虫害防治",
    "urgency": "normal",
    "images": [
        "https://example.com/questions/q1_img1.jpg",
        "https://example.com/questions/q1_img2.jpg"
    ],
    "location": {
        "province": "山东省",
        "city": "济南市",
        "county": "历城区"
    },
    "crop_info": {
        "crop_type": "小麦",
        "variety": "济麦22",
        "planting_date": "2023-10-15",
        "growth_stage": "分蘖期",
        "field_area": 10
    }
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "咨询提交成功",
    "data": {
        "consultation_id": "CON20240115001",
        "expert_id": 5001,
        "expert_name": "张教授",
        "status": "pending",
        "estimated_response_time": "2小时内",
        "created_at": "2024-01-15T14:30:00Z"
    }
}
```

### 3.4 获取咨询记录
```http
GET /api/content/consultations?status=all&page=1&limit=20
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total": 25,
        "page": 1,
        "limit": 20,
        "consultations": [
            {
                "id": "CON20240115001",
                "expert": {
                    "id": 5001,
                    "name": "张教授",
                    "title": "农业技术专家",
                    "avatar": "https://example.com/avatars/expert1.jpg"
                },
                "subject": "小麦叶片出现黄斑如何处理",
                "category": "病虫害防治",
                "status": "answered",
                "status_text": "已回复",
                "urgency": "normal",
                "created_at": "2024-01-15T14:30:00Z",
                "answered_at": "2024-01-15T16:45:00Z",
                "response_time": 2.25,
                "rating": 5,
                "is_public": false
            }
        ]
    }
}
```

### 3.5 获取咨询详情
```http
GET /api/content/consultations/{consultation_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 🔔 通知消息管理

### 4.1 获取通知列表
```http
GET /api/content/notifications?type=system&status=unread&page=1&limit=20
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total": 15,
        "unread_count": 8,
        "page": 1,
        "limit": 20,
        "notifications": [
            {
                "id": 60001,
                "type": "system",
                "category": "loan_approval",
                "title": "贷款审批结果通知",
                "content": "您的贷款申请LA20240115001已通过审批，金额80,000元",
                "priority": "high",
                "status": "unread",
                "action_required": true,
                "action_url": "/loans/applications/LA20240115001",
                "action_text": "查看详情",
                "sender": {
                    "type": "system",
                    "name": "系统通知"
                },
                "extra_data": {
                    "application_id": "LA20240115001",
                    "amount": 80000
                },
                "created_at": "2024-01-16T11:30:00Z",
                "read_at": null,
                "expires_at": "2024-01-23T23:59:59Z"
            },
            {
                "id": 60002,
                "type": "expert_reply",
                "category": "consultation",
                "title": "专家回复通知",
                "content": "张教授已回复您的咨询问题：小麦叶片出现黄斑如何处理",
                "priority": "normal",
                "status": "unread",
                "action_required": true,
                "action_url": "/consultations/CON20240115001",
                "action_text": "查看回复",
                "sender": {
                    "type": "expert",
                    "id": 5001,
                    "name": "张教授",
                    "avatar": "https://example.com/avatars/expert1.jpg"
                },
                "created_at": "2024-01-15T16:45:00Z"
            }
        ]
    }
}
```

### 4.2 标记通知已读
```http
PUT /api/content/notifications/{notification_id}/read
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 4.3 批量标记已读
```http
PUT /api/content/notifications/mark-all-read
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "notification_ids": [60001, 60002, 60003]
}
```

### 4.4 删除通知
```http
DELETE /api/content/notifications/{notification_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 🏷️ 内容分类管理

### 5.1 获取分类体系
```http
GET /api/content/categories?type=article
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "categories": [
            {
                "id": 1,
                "code": "tech",
                "name": "技术资讯",
                "description": "农业技术、种植指导、病虫害防治等",
                "icon": "https://example.com/icons/tech.png",
                "sort_order": 1,
                "children": [
                    {
                        "id": 11,
                        "code": "planting",
                        "name": "种植技术",
                        "parent_id": 1
                    },
                    {
                        "id": 12,
                        "code": "pest_control",
                        "name": "病虫害防治",
                        "parent_id": 1
                    }
                ]
            }
        ]
    }
}
```

### 5.2 获取热门标签
```http
GET /api/content/tags/popular?limit=20
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "tags": [
            {"name": "病虫害防治", "count": 156, "trend": "up"},
            {"name": "智慧农业", "count": 89, "trend": "up"},
            {"name": "有机种植", "count": 67, "trend": "stable"},
            {"name": "节水灌溉", "count": 45, "trend": "down"}
        ]
    }
}
```

---

## 💬 互动功能

### 6.1 点赞/取消点赞
```http
POST /api/content/articles/{article_id}/like
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 6.2 收藏/取消收藏
```http
POST /api/content/articles/{article_id}/collect
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 6.3 分享文章
```http
POST /api/content/articles/{article_id}/share
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "platform": "wechat",
    "message": "分享一篇很有用的文章"
}
```

### 6.4 提交评论
```http
POST /api/content/articles/{article_id}/comments
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "content": "文章写得很好，学到了很多实用的技术",
    "parent_id": null,
    "images": ["https://example.com/comments/img1.jpg"]
}
```

### 6.5 获取评论列表
```http
GET /api/content/articles/{article_id}/comments?page=1&limit=20&sort=latest
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 📊 内容统计分析

### 7.1 获取用户阅读统计
```http
GET /api/content/statistics/reading?period=month
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "reading_summary": {
            "total_articles_read": 45,
            "total_reading_time": 1250,
            "favorite_category": "技术资讯",
            "most_active_day": "周二"
        },
        "category_distribution": [
            {"category": "技术资讯", "count": 18, "percentage": 40},
            {"category": "市场行情", "count": 12, "percentage": 27},
            {"category": "政策资讯", "count": 8, "percentage": 18},
            {"category": "农业新闻", "count": 7, "percentage": 15}
        ],
        "reading_trend": [
            {"date": "2024-01-01", "articles": 3, "time": 45},
            {"date": "2024-01-02", "articles": 2, "time": 32}
        ]
    }
}
```

### 7.2 获取收藏列表
```http
GET /api/content/collections?type=article&page=1&limit=20
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 7.3 获取浏览历史
```http
GET /api/content/history?type=article&page=1&limit=20
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 🔧 错误码说明

| 错误码 | 说明 | 处理建议 |
|-------|------|---------|
| 5001 | 文章不存在 | 检查文章ID是否正确 |
| 5002 | 内容已下线 | 文章可能被删除或下架 |
| 5003 | 专家不存在 | 检查专家ID |
| 5004 | 专家暂不可用 | 选择其他在线专家 |
| 5005 | 咨询问题为空 | 输入咨询内容 |
| 5006 | 图片上传失败 | 重新上传图片 |
| 5007 | 通知不存在 | 检查通知ID |
| 5008 | 评论内容违规 | 修改评论内容 |
| 5009 | 政策申请已存在 | 不能重复申请 |
| 5010 | 申请材料不完整 | 补充必要材料 |

---

## 📝 接口调用示例

### JavaScript示例
```javascript
// 获取推荐资讯
const getRecommendedArticles = async (token) => {
    const response = await fetch('/api/content/articles/recommended', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};

// 提交咨询问题
const submitConsultation = async (token, consultationData) => {
    const response = await fetch('/api/content/consultations', {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(consultationData)
    });
    return response.json();
};

// 点赞文章
const likeArticle = async (token, articleId) => {
    const response = await fetch(`/api/content/articles/${articleId}/like`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};

// 获取通知列表
const getNotifications = async (token) => {
    const response = await fetch('/api/content/notifications?status=unread', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};
```

### 内容推荐算法示例
```javascript
// 基于用户行为的内容推荐
const getPersonalizedRecommendations = (userProfile) => {
    const factors = {
        location: userProfile.location,
        interests: userProfile.tags,
        readingHistory: userProfile.readingHistory,
        seasonality: getCurrentSeason(),
        trending: getTrendingTopics()
    };
    
    return calculateRecommendationScore(factors);
};
```

### 注意事项
1. **内容质量**: 确保发布的内容准确、及时、有价值
2. **个性化推荐**: 基于用户画像和行为数据推荐相关内容
3. **专家资质**: 严格审核专家资质和专业能力
4. **内容审核**: 建立完善的内容审核机制
5. **用户反馈**: 及时处理用户反馈和投诉
6. **数据统计**: 持续优化内容推荐算法 