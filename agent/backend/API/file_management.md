# 文件管理模块 - API 接口文档

## 📋 模块概述

文件管理模块为惠农APP/Web用户提供文件上传、下载、删除等功能。支持多种文件类型，包括身份认证材料、贷款申请附件、农机相关图片等。

### 核心功能
-   **文件上传**: 单文件和批量文件上传
-   **文件下载**: 安全的文件访问和下载
-   **文件管理**: 文件信息查询和删除
-   **权限控制**: 用户只能访问自己上传的文件

---

## 📁 惠农APP/Web - 文件管理接口

**接口路径前缀**: `/api/user/files`
**认证要求**: `RequireAuth` (惠农APP/Web用户)
**适用平台**: `app`, `web`

### 1.1 单文件上传

```http
POST /api/user/files/upload
Authorization: Bearer {access_token}
Content-Type: multipart/form-data

表单数据:
- file: 文件内容 (File)
- file_type: 文件类型 (string, 可选) - id_card_front, id_card_back, face_photo, bank_card, business_license, income_proof, loan_application, machine_photo, other
- category: 文件分类 (string, 可选) - auth_material, loan_document, machine_document, user_document
- description: 文件描述 (string, 可选)
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "文件上传成功",
    "data": {
        "file_id": "file_20240115_001",
        "original_filename": "身份证正面.jpg",
        "file_size": 2048576,
        "file_type": "id_card_front",
        "category": "auth_material",
        "content_type": "image/jpeg",
        "file_url": "https://storage.huinong.com/files/user_1001/file_20240115_001.jpg",
        "thumbnail_url": "https://storage.huinong.com/thumbnails/user_1001/file_20240115_001_thumb.jpg",
        "upload_time": "2024-01-15T10:30:00Z",
        "expires_at": "2024-04-15T10:30:00Z", // 如果有过期时间
        "is_verified": false, // 是否已通过审核验证
        "metadata": {
            "width": 1920,
            "height": 1080,
            "format": "JPEG"
        }
    }
}
```

### 1.2 批量文件上传

```http
POST /api/user/files/upload/batch
Authorization: Bearer {access_token}
Content-Type: multipart/form-data

表单数据:
- files: 多个文件 (File[])
- file_types: 对应的文件类型 (string[], 可选)
- category: 文件分类 (string, 可选)
- description: 批量描述 (string, 可选)
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "批量上传完成",
    "data": {
        "total_files": 3,
        "successful_uploads": 2,
        "failed_uploads": 1,
        "results": [
            {
                "index": 0,
                "status": "success",
                "file_id": "file_20240115_002",
                "original_filename": "身份证正面.jpg",
                "file_url": "https://storage.huinong.com/files/user_1001/file_20240115_002.jpg"
            },
            {
                "index": 1,
                "status": "success", 
                "file_id": "file_20240115_003",
                "original_filename": "身份证背面.jpg",
                "file_url": "https://storage.huinong.com/files/user_1001/file_20240115_003.jpg"
            },
            {
                "index": 2,
                "status": "failed",
                "error_code": 6002,
                "error_message": "文件格式不支持",
                "original_filename": "document.txt"
            }
        ]
    }
}
```

### 1.3 获取文件信息和下载

```http
GET /api/user/files/{file_id}
Authorization: Bearer {access_token}
```

**Query Parameters (可选)**:
-   `download` (boolean): `true` 表示直接下载文件，`false` 表示获取文件信息
-   `thumbnail` (boolean): `true` 表示获取缩略图（仅适用于图片）

**响应示例 (获取文件信息):**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "file_id": "file_20240115_001",
        "original_filename": "身份证正面.jpg",
        "file_size": 2048576,
        "file_type": "id_card_front",
        "category": "auth_material",
        "content_type": "image/jpeg",
        "file_url": "https://storage.huinong.com/files/user_1001/file_20240115_001.jpg",
        "thumbnail_url": "https://storage.huinong.com/thumbnails/user_1001/file_20240115_001_thumb.jpg",
        "download_url": "https://api.huinong.com/api/user/files/file_20240115_001?download=true&token=temp_download_token",
        "upload_time": "2024-01-15T10:30:00Z",
        "last_accessed": "2024-01-15T14:20:00Z",
        "download_count": 3,
        "is_verified": true,
        "verified_at": "2024-01-15T11:00:00Z",
        "status": "active", // active, deleted, expired
        "metadata": {
            "width": 1920,
            "height": 1080,
            "format": "JPEG",
            "file_hash": "sha256:abc123..."
        }
    }
}
```

**响应示例 (直接下载, download=true):**
返回文件二进制内容，并设置合适的HTTP头：
```http
Content-Type: image/jpeg
Content-Disposition: attachment; filename="身份证正面.jpg"
Content-Length: 2048576
Cache-Control: private, max-age=3600
```

### 1.4 删除文件

```http
DELETE /api/user/files/{file_id}
Authorization: Bearer {access_token}
```

**响应示例 (成功):**
```json
{
    "code": 200,
    "message": "文件删除成功",
    "data": {
        "file_id": "file_20240115_001",
        "deleted_at": "2024-01-15T16:00:00Z"
    }
}
```

### 1.5 获取用户文件列表

```http
GET /api/user/files?category=auth_material&file_type=id_card_front&page=1&limit=20
Authorization: Bearer {access_token}
```

**Query Parameters (可选)**:
-   `category` (string): 按分类筛选
-   `file_type` (string): 按文件类型筛选
-   `status` (string): 按状态筛选 (`active`, `deleted`, `expired`)
-   `is_verified` (boolean): 按验证状态筛选
-   `date_range_start`, `date_range_end` (string): 按上传时间筛选
-   `page`, `limit` (int): 分页参数

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total": 15,
        "files": [
            {
                "file_id": "file_20240115_001",
                "original_filename": "身份证正面.jpg",
                "file_size": 2048576,
                "file_type": "id_card_front",
                "category": "auth_material",
                "content_type": "image/jpeg",
                "thumbnail_url": "https://storage.huinong.com/thumbnails/user_1001/file_20240115_001_thumb.jpg",
                "upload_time": "2024-01-15T10:30:00Z",
                "is_verified": true,
                "status": "active"
            }
            // ... 更多文件
        ]
    }
}
```

---

## 📊 文件统计接口

### 2.1 获取用户文件统计

```http
GET /api/user/files/statistics
Authorization: Bearer {access_token}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total_files": 25,
        "total_size_bytes": 52428800,
        "total_size_readable": "50 MB",
        "by_category": {
            "auth_material": {
                "count": 6,
                "size_bytes": 12582912
            },
            "loan_document": {
                "count": 8,
                "size_bytes": 20971520
            },
            "machine_document": {
                "count": 5,
                "size_bytes": 10485760
            },
            "user_document": {
                "count": 6,
                "size_bytes": 8388608
            }
        },
        "by_file_type": {
            "id_card_front": {"count": 1, "size_bytes": 2048576},
            "id_card_back": {"count": 1, "size_bytes": 1985792},
            "bank_card": {"count": 1, "size_bytes": 1572864},
            "other": {"count": 22, "size_bytes": 46821568}
        },
        "verified_files": 18,
        "pending_verification": 7,
        "storage_quota": {
            "total_bytes": 1073741824, // 1GB
            "used_bytes": 52428800,
            "available_bytes": 1021312896,
            "usage_percentage": 4.9
        }
    }
}
```

---

## 🔧 错误码说明

| 错误码 | 说明 | 处理建议 |
|-------|------|---------|
| 6001 | 文件不存在 | 检查文件ID是否正确 |
| 6002 | 文件格式不支持 | 使用支持的文件格式 |
| 6003 | 文件大小超出限制 | 压缩文件或分片上传 |
| 6004 | 存储空间不足 | 清理无用文件或联系客服 |
| 6005 | 文件上传失败 | 检查网络连接后重试 |
| 6006 | 文件已被删除 | 重新上传文件 |
| 6007 | 文件类型参数无效 | 检查file_type参数 |
| 6008 | 文件访问权限不足 | 只能访问自己的文件 |
| 6009 | 文件已过期 | 重新上传文件 |
| 6010 | 批量上传部分失败 | 查看详细错误信息 |

---

## 📝 文件类型和限制说明

### 支持的文件类型
```javascript
const FILE_TYPES = {
    // 身份认证相关
    'id_card_front': {
        name: '身份证正面',
        extensions: ['.jpg', '.jpeg', '.png'],
        max_size: '5MB',
        required_for: 'real_name_auth'
    },
    'id_card_back': {
        name: '身份证背面', 
        extensions: ['.jpg', '.jpeg', '.png'],
        max_size: '5MB',
        required_for: 'real_name_auth'
    },
    'face_photo': {
        name: '人脸照片',
        extensions: ['.jpg', '.jpeg', '.png'],
        max_size: '3MB',
        required_for: 'real_name_auth'
    },
    'bank_card': {
        name: '银行卡照片',
        extensions: ['.jpg', '.jpeg', '.png'],
        max_size: '5MB',
        required_for: 'bank_card_auth'
    },
    
    // 贷款相关
    'business_license': {
        name: '营业执照',
        extensions: ['.jpg', '.jpeg', '.png', '.pdf'],
        max_size: '10MB',
        required_for: 'loan_application'
    },
    'income_proof': {
        name: '收入证明',
        extensions: ['.jpg', '.jpeg', '.png', '.pdf'],
        max_size: '10MB',
        required_for: 'loan_application'
    },
    'loan_application': {
        name: '贷款申请材料',
        extensions: ['.jpg', '.jpeg', '.png', '.pdf', '.doc', '.docx'],
        max_size: '20MB',
        required_for: 'loan_application'
    },
    
    // 农机相关
    'machine_photo': {
        name: '农机照片',
        extensions: ['.jpg', '.jpeg', '.png'],
        max_size: '10MB',
        required_for: 'machine_registration'
    },
    
    // 其他
    'other': {
        name: '其他文件',
        extensions: ['.jpg', '.jpeg', '.png', '.pdf', '.doc', '.docx', '.txt'],
        max_size: '50MB',
        required_for: null
    }
};
```

### 文件分类
```javascript
const FILE_CATEGORIES = {
    'auth_material': '认证材料',
    'loan_document': '贷款文档',
    'machine_document': '农机文档',
    'user_document': '用户文档'
};
```

---

## 📝 接口调用示例

### JavaScript示例
```javascript
// 单文件上传
const uploadFile = async (token, file, fileType) => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('file_type', fileType);
    formData.append('category', 'auth_material');
    
    const response = await fetch('/api/user/files/upload', {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`
        },
        body: formData
    });
    return response.json();
};

// 批量文件上传
const uploadMultipleFiles = async (token, files, fileTypes) => {
    const formData = new FormData();
    files.forEach((file, index) => {
        formData.append('files', file);
    });
    formData.append('file_types', JSON.stringify(fileTypes));
    formData.append('category', 'auth_material');
    
    const response = await fetch('/api/user/files/upload/batch', {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`
        },
        body: formData
    });
    return response.json();
};

// 获取文件信息
const getFileInfo = async (token, fileId) => {
    const response = await fetch(`/api/user/files/${fileId}`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};

// 下载文件
const downloadFile = async (token, fileId, filename) => {
    const response = await fetch(`/api/user/files/${fileId}?download=true`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    
    if (response.ok) {
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
    }
};

// 删除文件
const deleteFile = async (token, fileId) => {
    const response = await fetch(`/api/user/files/${fileId}`, {
        method: 'DELETE',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.json();
};
```

### Vue.js 文件上传组件示例
```vue
<template>
    <div class="file-upload">
        <el-upload
            :action="uploadUrl"
            :headers="uploadHeaders"
            :data="uploadData"
            :on-success="handleSuccess"
            :on-error="handleError"
            :before-upload="beforeUpload"
            drag
            multiple>
            <i class="el-icon-upload"></i>
            <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
        </el-upload>
    </div>
</template>

<script>
export default {
    data() {
        return {
            uploadUrl: '/api/user/files/upload/batch',
            uploadHeaders: {
                'Authorization': `Bearer ${this.$store.state.auth.token}`
            },
            uploadData: {
                category: 'auth_material'
            }
        };
    },
    methods: {
        beforeUpload(file) {
            const isValidType = ['image/jpeg', 'image/png', 'application/pdf'].includes(file.type);
            const isValidSize = file.size < 10 * 1024 * 1024; // 10MB
            
            if (!isValidType) {
                this.$message.error('文件格式不支持！');
                return false;
            }
            if (!isValidSize) {
                this.$message.error('文件大小不能超过10MB！');
                return false;
            }
            return true;
        },
        handleSuccess(response) {
            if (response.code === 200) {
                this.$message.success('上传成功');
                this.$emit('upload-success', response.data);
            } else {
                this.$message.error(response.message);
            }
        },
        handleError(error) {
            this.$message.error('上传失败，请重试');
        }
    }
};
</script>
```

### 安全注意事项
1. **文件类型验证**: 严格验证文件扩展名和MIME类型
2. **文件大小限制**: 根据文件类型设置合理的大小限制
3. **病毒扫描**: 上传的文件需要进行安全扫描
4. **权限控制**: 用户只能访问自己上传的文件
5. **存储安全**: 文件存储在安全的云存储服务中
6. **访问控制**: 使用临时URL或Token控制文件访问
7. **数据备份**: 重要文件需要定期备份 