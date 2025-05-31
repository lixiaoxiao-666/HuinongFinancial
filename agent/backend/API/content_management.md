# 内容管理模块 - API 接口文档

## 📋 模块概述

内容管理模块提供丰富的农业信息服务，包括农业资讯、专家咨询、系统公告等功能。支持公开访问和登录用户个性化内容推荐，同时提供完善的后台内容管理功能。

### 🎯 核心功能
- **文章管理**: 农业资讯、技术文章、行业动态
- **专家服务**: 专家信息、在线咨询、知识问答
- **公告管理**: 系统公告、重要通知、政策发布
- **分类管理**: 内容分类、标签系统、个性化推荐
- **用户互动**: 咨询提交、专家回答、内容反馈

### 🏗️ 内容架构
```
内容体系
├── 文章资讯 (Articles)
│   ├── 农业新闻
│   ├── 技术指导
│   └── 市场行情
├── 专家服务 (Experts)
│   ├── 专家信息
│   ├── 在线咨询
│   └── 专业解答
├── 系统公告 (Announcements)
│   ├── 系统通知
│   ├── 政策公告
│   └── 功能更新
└── 内容分类 (Categories)
    ├── 按行业分类
    ├── 按内容类型
    └── 按用户标签
```

### 📊 数据模型关系
```
Articles (文章)
├── Categories (分类)
├── Tags (标签)
└── UserInteractions (用户互动)

Experts (专家)
├── Specialties (专业领域)
├── Certifications (资质认证)
└── ConsultationHistory (咨询历史)

Consultations (咨询)
├── Questions (问题)
├── Answers (回答)
└── FollowUps (追问)

Announcements (公告)
├── AnnouncementTypes (公告类型)
└── ReadStatistics (阅读统计)
```

---

## 📰 公共内容接口 (可选认证)

### 1. 获取文章列表
**接口路径**: `GET /api/content/articles`  
**认证要求**: 可选认证 (登录用户可获得个性化内容)  
**功能描述**: 获取文章列表，支持分类筛选和搜索

#### 请求参数
```
?page={page}              # 页码，默认1
&limit={limit}            # 每页数量，默认20
&category_id={id}         # 分类筛选
&tag={tag}               # 标签筛选
&keyword={keyword}        # 关键词搜索
&featured={boolean}       # 是否推荐文章
&sort_by={field}          # 排序字段 (created_at/views/likes)
&sort_order={desc|asc}    # 排序方向
```

#### 响应示例
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "articles": [
            {
                "article_id": "ART20240115001",
                "title": "春季小麦种植技术要点",
                "summary": "详细介绍春季小麦种植的关键技术和注意事项...",
                "cover_image": "https://oss.example.com/articles/wheat_planting.jpg",
                "category": {
                    "category_id": "CAT_001",
                    "category_name": "种植技术"
                },
                "tags": ["小麦", "春季种植", "农业技术"],
                "author": {
                    "author_id": "AUTH_001",
                    "author_name": "农业专家李老师",
                    "avatar": "https://oss.example.com/avatars/expert_li.jpg"
                },
                "publish_time": "2024-01-15T08:00:00Z",
                "read_count": 1234,
                "like_count": 89,
                "comment_count": 23,
                "is_featured": true,
                "content_preview": "春季是小麦种植的关键时期，正确的种植技术...",
                "estimated_read_time": "3分钟"
            },
            {
                "article_id": "ART20240115002",
                "title": "农机设备维护保养指南",
                "summary": "全面介绍农机设备的日常维护和保养方法...",
                "cover_image": "https://oss.example.com/articles/machine_maintenance.jpg",
                "category": {
                    "category_id": "CAT_002",
                    "category_name": "设备维护"
                },
                "tags": ["农机", "维护保养", "设备管理"],
                "author": {
                    "author_id": "AUTH_002",
                    "author_name": "机械专家王师傅",
                    "avatar": "https://oss.example.com/avatars/expert_wang.jpg"
                },
                "publish_time": "2024-01-14T14:30:00Z",
                "read_count": 856,
                "like_count": 67,
                "comment_count": 15,
                "is_featured": false,
                "content_preview": "农机设备的正确维护保养是确保设备长期稳定运行...",
                "estimated_read_time": "5分钟"
            }
        ],
        "pagination": {
            "page": 1,
            "limit": 20,
            "total": 156,
            "pages": 8
        },
        "personalized_recommendations": [
            {
                "article_id": "ART20240115003",
                "title": "智能灌溉系统应用案例",
                "reason": "基于您的关注领域推荐",
                "match_score": 95
            }
        ]
    }
}
```

#### JavaScript调用示例
```javascript
// 获取文章列表
async function getArticles(params = {}) {
    try {
        const queryParams = new URLSearchParams({
            page: 1,
            limit: 20,
            ...params
        });
        
        const headers = {
            'Content-Type': 'application/json'
        };
        
        // 如果用户已登录，添加认证头
        const token = localStorage.getItem('access_token');
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
        
        const response = await fetch(`/api/content/articles?${queryParams}`, {
            method: 'GET',
            headers: headers
        });
        
        const result = await response.json();
        if (result.code === 200) {
            console.log('文章列表:', result.data);
            return result.data;
        } else {
            throw new Error(result.message);
        }
    } catch (error) {
        console.error('获取文章列表失败:', error);
        throw error;
    }
}
```

### 2. 获取推荐文章
**接口路径**: `GET /api/content/articles/featured`  
**认证要求**: 可选认证  
**功能描述**: 获取编辑推荐的优质文章

#### 请求参数
```
?limit={limit}            # 数量限制，默认10
&category_id={id}         # 分类筛选
```

#### 响应示例
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "featured_articles": [
            {
                "article_id": "ART20240115001",
                "title": "春季小麦种植技术要点",
                "summary": "详细介绍春季小麦种植的关键技术和注意事项",
                "cover_image": "https://oss.example.com/articles/wheat_planting.jpg",
                "featured_reason": "编辑推荐：应季农业技术",
                "featured_at": "2024-01-15T08:00:00Z",
                "priority": 1
            }
        ],
        "total_featured": 8
    }
}
```

### 3. 获取文章详情
**接口路径**: `GET /api/content/articles/{article_id}`  
**认证要求**: 可选认证  
**功能描述**: 获取文章完整内容和详细信息

#### 路径参数
- `article_id`: 文章ID

#### 响应示例
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "article_info": {
            "article_id": "ART20240115001",
            "title": "春季小麦种植技术要点",
            "content": "春季是小麦种植的关键时期...[完整文章内容]",
            "summary": "详细介绍春季小麦种植的关键技术和注意事项",
            "cover_image": "https://oss.example.com/articles/wheat_planting.jpg",
            "images": [
                "https://oss.example.com/articles/wheat_1.jpg",
                "https://oss.example.com/articles/wheat_2.jpg"
            ],
            "category": {
                "category_id": "CAT_001",
                "category_name": "种植技术",
                "category_path": "农业技术 > 种植技术"
            },
            "tags": ["小麦", "春季种植", "农业技术"],
            "author": {
                "author_id": "AUTH_001",
                "author_name": "农业专家李老师",
                "avatar": "https://oss.example.com/avatars/expert_li.jpg",
                "bio": "从事农业技术推广20年，专注作物种植技术研究",
                "certifications": ["高级农艺师", "作物栽培专家"]
            },
            "publish_time": "2024-01-15T08:00:00Z",
            "update_time": "2024-01-15T08:30:00Z",
            "read_count": 1235,
            "like_count": 89,
            "comment_count": 23,
            "is_featured": true,
            "estimated_read_time": "3分钟",
            "content_quality_score": 9.2
        },
        "related_articles": [
            {
                "article_id": "ART20240115004",
                "title": "小麦病虫害防治指南",
                "cover_image": "https://oss.example.com/articles/wheat_pest.jpg",
                "similarity_score": 85
            }
        ],
        "user_interaction": {
            "has_liked": false,
            "has_bookmarked": false,
            "reading_progress": 0,
            "last_read_at": null
        }
    }
}
```

### 4. 获取文章分类
**接口路径**: `GET /api/content/categories`  
**认证要求**: 无需认证  
**功能描述**: 获取文章分类列表

#### 响应示例
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "categories": [
            {
                "category_id": "CAT_001",
                "category_name": "种植技术",
                "parent_id": "CAT_PARENT_001",
                "parent_name": "农业技术",
                "description": "作物种植相关技术和方法",
                "icon": "https://oss.example.com/icons/planting.png",
                "article_count": 234,
                "sort_order": 1,
                "children": [
                    {
                        "category_id": "CAT_001_001",
                        "category_name": "粮食作物",
                        "article_count": 123
                    },
                    {
                        "category_id": "CAT_001_002",
                        "category_name": "经济作物",
                        "article_count": 111
                    }
                ]
            },
            {
                "category_id": "CAT_002",
                "category_name": "设备维护",
                "parent_id": "CAT_PARENT_002",
                "parent_name": "农机技术",
                "description": "农机设备使用和维护",
                "icon": "https://oss.example.com/icons/machine.png",
                "article_count": 156,
                "sort_order": 2,
                "children": []
            }
        ],
        "category_tree": [
            {
                "name": "农业技术",
                "children": ["种植技术", "养殖技术", "植保技术"]
            },
            {
                "name": "农机技术", 
                "children": ["设备维护", "操作技巧", "故障排除"]
            }
        ]
    }
}
```

### 5. 获取专家列表
**接口路径**: `GET /api/content/experts`  
**认证要求**: 可选认证  
**功能描述**: 获取专家信息列表

#### 请求参数
```
?page={page}              # 页码，默认1
&limit={limit}            # 每页数量，默认20
&specialty={specialty}    # 专业领域筛选
&sort_by={field}          # 排序字段 (rating/experience/consultations)
&sort_order={desc|asc}    # 排序方向
```

#### 响应示例
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "experts": [
            {
                "expert_id": "EXP_001",
                "name": "李农业",
                "title": "高级农艺师",
                "avatar": "https://oss.example.com/avatars/expert_li.jpg",
                "specialties": ["作物栽培", "土壤改良", "病虫害防治"],
                "bio": "从事农业技术推广20年，专注作物种植技术研究",
                "education": "中国农业大学农学博士",
                "certifications": ["高级农艺师", "作物栽培专家"],
                "experience_years": 20,
                "consultation_count": 1234,
                "rating": 4.8,
                "response_rate": 95.6,
                "avg_response_time": "2小时",
                "online_status": "online",
                "consultation_fee": 0,
                "is_featured": true
            },
            {
                "expert_id": "EXP_002",
                "name": "王机械",
                "title": "农机专家",
                "avatar": "https://oss.example.com/avatars/expert_wang.jpg",
                "specialties": ["农机维修", "设备选型", "智能农机"],
                "bio": "专业农机技术服务15年，精通各类农机设备",
                "education": "农业工程硕士",
                "certifications": ["农机维修技师", "设备工程师"],
                "experience_years": 15,
                "consultation_count": 856,
                "rating": 4.7,
                "response_rate": 92.3,
                "avg_response_time": "3小时",
                "online_status": "offline",
                "consultation_fee": 50,
                "is_featured": false
            }
        ],
        "pagination": {
            "page": 1,
            "limit": 20,
            "total": 45,
            "pages": 3
        },
        "featured_experts": [
            {
                "expert_id": "EXP_001",
                "name": "李农业",
                "featured_reason": "本月咨询量最高"
            }
        ]
    }
}
```

### 6. 获取专家详情
**接口路径**: `GET /api/content/experts/{expert_id}`  
**认证要求**: 可选认证  
**功能描述**: 获取专家详细信息

#### 路径参数
- `expert_id`: 专家ID

#### 响应示例
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "expert_info": {
            "expert_id": "EXP_001",
            "name": "李农业",
            "title": "高级农艺师",
            "avatar": "https://oss.example.com/avatars/expert_li.jpg",
            "specialties": ["作物栽培", "土壤改良", "病虫害防治"],
            "bio": "从事农业技术推广20年，专注作物种植技术研究...",
            "detailed_intro": "李农业老师拥有20年丰富的农业实践经验...",
            "education": "中国农业大学农学博士",
            "certifications": [
                {
                    "name": "高级农艺师",
                    "issued_by": "农业部",
                    "issued_date": "2015-06-01"
                }
            ],
            "experience_years": 20,
            "work_history": [
                {
                    "company": "XX农业技术推广站",
                    "position": "技术主管",
                    "duration": "2010-2020"
                }
            ]
        },
        "consultation_stats": {
            "consultation_count": 1234,
            "rating": 4.8,
            "response_rate": 95.6,
            "avg_response_time": "2小时",
            "satisfaction_rate": 97.2,
            "repeat_consultation_rate": 68.5
        },
        "service_info": {
            "consultation_fee": 0,
            "online_status": "online",
            "available_time": "工作日 9:00-18:00",
            "consultation_methods": ["文字咨询", "语音咨询", "视频咨询"],
            "languages": ["中文"]
        },
        "recent_consultations": [
            {
                "consultation_id": "CONS_001",
                "question_title": "小麦叶片发黄怎么办？",
                "answered_at": "2024-01-15T10:30:00Z",
                "user_rating": 5,
                "is_public": true
            }
        ],
        "published_articles": [
            {
                "article_id": "ART20240115001",
                "title": "春季小麦种植技术要点",
                "publish_time": "2024-01-15T08:00:00Z",
                "read_count": 1234
            }
        ]
    }
}
```

---

## 🙋‍♂️ 用户咨询功能

### 7. 提交咨询问题
**接口路径**: `POST /api/user/consultations`  
**认证要求**: 需要认证 (惠农用户)  
**功能描述**: 用户提交专家咨询问题

#### 请求参数
```json
{
    "expert_id": "EXP_001",
    "title": "小麦叶片出现黄斑怎么办？",
    "content": "我家的小麦叶片最近出现了黄色斑点，不知道是什么原因，该如何处理？",
    "category": "病虫害防治",
    "urgency": "normal",
    "images": [
        "https://oss.example.com/consultations/wheat_problem_1.jpg",
        "https://oss.example.com/consultations/wheat_problem_2.jpg"
    ],
    "location": {
        "province": "山东省",
        "city": "济南市",
        "district": "历下区"
    },
    "crop_info": {
        "crop_type": "小麦",
        "planting_area": "50亩",
        "growth_stage": "拔节期",
        "planting_date": "2023-10-15"
    },
    "is_public": true
}
```

#### 响应示例
```json
{
    "code": 201,
    "message": "咨询提交成功",
    "data": {
        "consultation_id": "CONS20240115001",
        "consultation_number": "CN202401150001",
        "expert_id": "EXP_001",
        "expert_name": "李农业",
        "title": "小麦叶片出现黄斑怎么办？",
        "status": "pending",
        "submitted_at": "2024-01-15T11:00:00Z",
        "estimated_response_time": "2-4小时",
        "consultation_fee": 0,
        "tracking_info": {
            "current_stage": "expert_review",
            "progress_percentage": 25,
            "estimated_completion": "2024-01-15T15:00:00Z"
        }
    }
}
```

### 8. 获取咨询记录
**接口路径**: `GET /api/user/consultations`  
**认证要求**: 需要认证 (惠农用户)  
**功能描述**: 获取用户的咨询记录列表

#### 请求参数
```
?page={page}              # 页码，默认1
&limit={limit}            # 每页数量，默认10
&status={status}          # 状态筛选 (pending/answered/closed)
&expert_id={expert_id}    # 专家筛选
&category={category}      # 分类筛选
&date_from={date}         # 日期起始
&date_to={date}           # 日期结束
&sort_by={field}          # 排序字段
&sort_order={desc|asc}    # 排序方向
```

#### 响应示例
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "consultations": [
            {
                "consultation_id": "CONS20240115001",
                "consultation_number": "CN202401150001",
                "title": "小麦叶片出现黄斑怎么办？",
                "expert_info": {
                    "expert_id": "EXP_001",
                    "expert_name": "李农业",
                    "avatar": "https://oss.example.com/avatars/expert_li.jpg",
                    "title": "高级农艺师"
                },
                "category": "病虫害防治",
                "status": "answered",
                "status_text": "已回答",
                "urgency": "normal",
                "submitted_at": "2024-01-15T11:00:00Z",
                "answered_at": "2024-01-15T13:30:00Z",
                "response_time_hours": 2.5,
                "user_rating": 5,
                "has_new_message": false,
                "message_count": 3,
                "is_public": true
            }
        ],
        "pagination": {
            "page": 1,
            "limit": 10,
            "total": 12,
            "pages": 2
        },
        "summary": {
            "total_consultations": 12,
            "pending_consultations": 2,
            "answered_consultations": 8,
            "closed_consultations": 2,
            "average_rating": 4.6
        }
    }
}
```

---

## 🏢 OA后台内容管理

### 9. 创建文章
**接口路径**: `POST /api/oa/admin/content/articles`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 创建新文章

#### 请求参数
```json
{
    "title": "春季农业生产技术指导",
    "content": "春季是农业生产的关键时期...[文章内容]",
    "summary": "详细介绍春季农业生产的关键技术和注意事项",
    "cover_image": "https://oss.example.com/articles/spring_farming.jpg",
    "category_id": "CAT_001",
    "tags": ["春季生产", "农业技术", "种植指导"],
    "author_id": "AUTH_001",
    "is_featured": true,
    "publish_immediately": true,
    "scheduled_publish_time": null,
    "seo_keywords": "春季,农业,种植,技术指导",
    "meta_description": "春季农业生产技术指导，包含种植、管理等关键环节"
}
```

#### 响应示例
```json
{
    "code": 201,
    "message": "文章创建成功",
    "data": {
        "article_id": "ART20240115003",
        "title": "春季农业生产技术指导",
        "status": "published",
        "created_by": "ADMIN001",
        "created_at": "2024-01-15T11:30:00Z",
        "published_at": "2024-01-15T11:30:00Z",
        "article_url": "/articles/ART20240115003"
    }
}
```

### 10. 更新文章
**接口路径**: `PUT /api/oa/admin/content/articles/{article_id}`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 更新指定文章

#### 路径参数
- `article_id`: 文章ID

#### 请求参数
```json
{
    "title": "春季农业生产技术指导（更新版）",
    "content": "更新后的文章内容...",
    "summary": "更新后的摘要",
    "category_id": "CAT_001",
    "tags": ["春季生产", "农业技术", "种植指导", "更新"],
    "is_featured": true,
    "update_reason": "补充最新技术内容"
}
```

#### 响应示例
```json
{
    "code": 200,
    "message": "文章更新成功",
    "data": {
        "article_id": "ART20240115003",
        "title": "春季农业生产技术指导（更新版）",
        "status": "published",
        "updated_by": "ADMIN001",
        "updated_at": "2024-01-15T14:00:00Z",
        "version": 2
    }
}
```

### 11. 删除文章
**接口路径**: `DELETE /api/oa/admin/content/articles/{article_id}`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 删除指定文章

#### 路径参数
- `article_id`: 文章ID

#### 响应示例
```json
{
    "code": 200,
    "message": "文章删除成功",
    "data": {
        "article_id": "ART20240115003",
        "deleted_by": "ADMIN001",
        "deleted_at": "2024-01-15T14:30:00Z"
    }
}
```

### 12. 发布文章
**接口路径**: `POST /api/oa/admin/content/articles/{article_id}/publish`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 发布或重新发布文章

#### 路径参数
- `article_id`: 文章ID

#### 请求参数 (可选)
```json
{
    "publish_time": "2024-01-15T16:00:00Z",
    "publish_channels": ["website", "app", "wechat"],
    "notify_subscribers": true
}
```

#### 响应示例
```json
{
    "code": 200,
    "message": "文章发布成功",
    "data": {
        "article_id": "ART20240115003",
        "status": "published",
        "published_by": "ADMIN001",
        "published_at": "2024-01-15T16:00:00Z",
        "publish_channels": ["website", "app", "wechat"],
        "estimated_reach": 5000
    }
}
```

### 13. 获取公告列表
**接口路径**: `GET /api/oa/admin/content/announcements`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 获取系统公告列表

#### 请求参数
```
?limit={limit}            # 数量限制，默认10
&status={status}          # 状态筛选 (published/draft/expired)
```

#### 响应示例
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "announcements": [
            {
                "id": 1,
                "title": "系统维护通知",
                "content": "系统将于今晚22:00-24:00进行维护，期间可能影响部分功能使用。",
                "status": "published",
                "created_at": "2024-01-15T08:00:00Z",
                "updated_at": "2024-01-15T08:00:00Z"
            },
            {
                "id": 2,
                "title": "新功能上线公告",
                "content": "AI智能风险评估功能已正式上线，将大幅提升审批效率。",
                "status": "published",
                "created_at": "2024-01-14T10:00:00Z",
                "updated_at": "2024-01-14T10:00:00Z"
            },
            {
                "id": 3,
                "title": "节假日服务安排",
                "content": "春节期间客服时间调整为9:00-18:00，给您带来不便敬请谅解。",
                "status": "published",
                "created_at": "2024-01-13T15:30:00Z",
                "updated_at": "2024-01-13T15:30:00Z"
            }
        ],
        "total": 3,
        "limit": "10"
    }
}
```

### 14. 创建公告
**接口路径**: `POST /api/oa/admin/content/announcements`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 创建系统公告

#### 请求参数
```json
{
    "title": "重要系统升级公告",
    "content": "为提升系统性能和用户体验，系统将于本周末进行重大升级...",
    "type": "system",
    "priority": "high",
    "target_audience": "all",
    "effective_time": "2024-01-16T00:00:00Z",
    "expire_time": "2024-01-30T23:59:59Z",
    "auto_publish": true
}
```

#### 响应示例
```json
{
    "code": 200,
    "message": "公告创建成功",
    "data": {
        "id": 4,
        "title": "重要系统升级公告",
        "status": "published",
        "created_at": "2024-01-15T15:00:00Z"
    }
}
```

### 15. 更新公告
**接口路径**: `PUT /api/oa/admin/content/announcements/{id}`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 更新指定公告

#### 路径参数
- `id`: 公告ID

#### 请求参数
```json
{
    "title": "重要系统升级公告（更新）",
    "content": "更新后的公告内容...",
    "priority": "urgent"
}
```

#### 响应示例
```json
{
    "code": 200,
    "message": "公告更新成功",
    "data": {
        "id": "4",
        "updated_at": "2024-01-15T16:00:00Z"
    }
}
```

### 16. 删除公告
**接口路径**: `DELETE /api/oa/admin/content/announcements/{id}`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 删除指定公告

#### 路径参数
- `id`: 公告ID

#### 响应示例
```json
{
    "code": 200,
    "message": "公告删除成功",
    "data": {
        "id": "4",
        "deleted_at": "2024-01-15T16:30:00Z"
    }
}
```

### 17. 创建分类
**接口路径**: `POST /api/oa/admin/content/categories`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 创建内容分类

#### 请求参数
```json
{
    "name": "智慧农业",
    "parent_id": "CAT_PARENT_001",
    "description": "智慧农业技术和应用相关内容",
    "icon": "https://oss.example.com/icons/smart_farming.png",
    "sort_order": 10,
    "is_active": true
}
```

#### 响应示例
```json
{
    "code": 201,
    "message": "分类创建成功",
    "data": {
        "category_id": "CAT_004",
        "name": "智慧农业",
        "created_at": "2024-01-15T17:00:00Z"
    }
}
```

### 18. 创建专家
**接口路径**: `POST /api/oa/admin/content/experts`  
**认证要求**: 需要认证 (OA管理员)  
**功能描述**: 添加新专家

#### 请求参数
```json
{
    "name": "赵智能",
    "title": "智慧农业专家",
    "avatar": "https://oss.example.com/avatars/expert_zhao.jpg",
    "specialties": ["智慧农业", "物联网技术", "数据分析"],
    "bio": "专注智慧农业技术研发和应用推广",
    "education": "农业信息化博士",
    "certifications": ["智慧农业专家", "物联网工程师"],
    "experience_years": 12,
    "consultation_fee": 100,
    "available_time": "工作日 9:00-17:00",
    "consultation_methods": ["文字咨询", "视频咨询"],
    "is_featured": false,
    "is_active": true
}
```

#### 响应示例
```json
{
    "code": 201,
    "message": "专家添加成功",
    "data": {
        "expert_id": "EXP_003",
        "name": "赵智能",
        "status": "active",
        "created_at": "2024-01-15T17:30:00Z"
    }
}
```

---

## ⚠️ 错误码说明

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 5001 | 文章不存在 | 检查文章ID是否正确 |
| 5002 | 分类不存在 | 检查分类ID是否正确 |
| 5003 | 专家不存在 | 检查专家ID是否正确 |
| 5004 | 内容审核中 | 等待审核完成 |
| 5005 | 权限不足 | 检查用户权限 |
| 5006 | 文件上传失败 | 检查文件格式和大小 |
| 5007 | 咨询服务暂停 | 联系客服了解详情 |
| 5008 | 专家不在线 | 选择其他专家或稍后咨询 |
| 5009 | 内容违规 | 修改内容后重新提交 |
| 5010 | 操作频率过高 | 等待一段时间后重试 |

---

## 🔄 最佳实践

### 内容质量
1. **原创性**: 提供原创、高质量的农业技术内容
2. **实用性**: 确保内容贴近农户实际需求
3. **时效性**: 及时更新季节性和时效性内容
4. **专业性**: 确保技术内容的准确性和专业性

### 用户体验
1. **分类清晰**: 合理的内容分类和标签体系
2. **搜索优化**: 优化搜索功能，提高内容发现效率
3. **个性化推荐**: 基于用户行为提供个性化内容
4. **互动体验**: 鼓励用户参与评论和讨论

### 运营管理
1. **内容审核**: 建立完善的内容审核机制
2. **专家管理**: 严格专家资质审核和服务质量监控
3. **数据分析**: 定期分析内容数据，优化内容策略
4. **用户反馈**: 收集用户反馈，持续改进服务质量 