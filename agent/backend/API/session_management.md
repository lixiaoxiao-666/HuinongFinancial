# 会话管理系统 - API 使用指南

## 🎯 系统概述

数字惠农系统已升级为基于Redis的分布式会话管理，替代传统的JWT认证，提供更强大的会话控制和多后端支持。

### ✨ 核心特性
- 🔐 **分布式会话**: 多后端实例共享用户会话状态
- 📱 **多端登录**: 支持APP、Web、OA后台同时登录
- ⚡ **高性能**: Redis毫秒级会话验证
- 🔄 **实时同步**: 会话状态实时同步，支持强制下线
- 🛡️ **安全机制**: Token哈希存储、自动清理、设备绑定

## 📋 API 接口文档

### 1. 用户认证相关

#### 1.1 用户登录
```http
POST /api/auth/login
Content-Type: application/json

{
    "phone": "13800138000",
    "password": "password123",
    "platform": "app",
    "device_id": "iPhone_12_ABC123",
    "device_type": "ios",
    "device_name": "John's iPhone",
    "app_version": "1.0.0",
    "ip_address": "192.168.1.100",
    "location": "北京市朝阳区"
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "登录成功",
    "data": {
        "user": {
            "id": 1001,
            "phone": "13800138000",
            "real_name": "张三",
            "user_type": "farmer"
        },
        "session": {
            "session_id": "sess_abc123def456",
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "expires_in": 86400,
            "platform": "app"
        }
    }
}
```

#### 1.2 Token刷新
```http
POST /api/auth/refresh
Content-Type: application/x-www-form-urlencoded

refresh_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "刷新成功",
    "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expires_in": 86400
    }
}
```

### 2. 会话管理相关

#### 2.1 获取会话信息
```http
GET /api/user/session/info
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取成功",
    "data": [
        {
            "session_id": "sess_abc123def456",
            "platform": "app",
            "device_info": {
                "device_id": "iPhone_12_ABC123",
                "device_type": "ios",
                "device_name": "John's iPhone",
                "app_version": "1.0.0"
            },
            "network_info": {
                "ip_address": "192.168.1.100",
                "location": "北京市朝阳区"
            },
            "status": "active",
            "created_at": "2024-01-15T10:30:00Z",
            "last_active_at": "2024-01-15T14:25:30Z"
        }
    ]
}
```

#### 2.2 注销当前会话
```http
POST /api/user/logout
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "注销成功"
}
```

#### 2.3 注销其他设备会话
```http
POST /api/user/session/revoke-others
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "注销其他会话成功"
}
```

### 3. 管理员会话管理

#### 3.1 获取所有活跃会话
```http
GET /api/admin/sessions/active
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取活跃会话列表",
    "data": {
        "total": 150,
        "sessions": [
            {
                "session_id": "sess_abc123",
                "user_id": 1001,
                "platform": "app",
                "device_type": "ios",
                "ip_address": "192.168.1.100",
                "last_active_at": "2024-01-15T14:25:30Z"
            }
        ]
    }
}
```

#### 3.2 强制注销会话
```http
DELETE /api/admin/sessions/sess_abc123def456
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "会话已强制注销",
    "session_id": "sess_abc123def456"
}
```

#### 3.3 清理过期会话
```http
POST /api/admin/sessions/cleanup
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "会话清理完成",
    "cleaned": 23
}
```

## 🔧 客户端集成指南

### 1. 移动端 (APP)

#### 登录流程
```javascript
// 1. 用户登录
const loginResponse = await fetch('/api/auth/login', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({
        phone: '13800138000',
        password: 'password123',
        platform: 'app',
        device_id: getDeviceId(),
        device_type: 'ios',
        device_name: getDeviceName(),
        app_version: '1.0.0',
        ip_address: await getClientIP(),
        location: await getLocation()
    })
});

const { data } = await loginResponse.json();

// 2. 存储Token
localStorage.setItem('access_token', data.session.access_token);
localStorage.setItem('refresh_token', data.session.refresh_token);
localStorage.setItem('session_id', data.session.session_id);
```

#### 自动Token刷新
```javascript
// 拦截器自动刷新Token
axios.interceptors.response.use(
    response => response,
    async error => {
        if (error.response?.status === 401) {
            const refreshToken = localStorage.getItem('refresh_token');
            
            try {
                const refreshResponse = await fetch('/api/auth/refresh', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: `refresh_token=${refreshToken}`
                });
                
                const { data } = await refreshResponse.json();
                
                // 更新Token
                localStorage.setItem('access_token', data.access_token);
                localStorage.setItem('refresh_token', data.refresh_token);
                
                // 重试原请求
                return axios.request(error.config);
            } catch (refreshError) {
                // 刷新失败，跳转到登录页
                redirectToLogin();
            }
        }
        return Promise.reject(error);
    }
);
```

### 2. Web端

#### Axios配置
```javascript
// 创建axios实例
const api = axios.create({
    baseURL: '/api',
    timeout: 10000
});

// 请求拦截器
api.interceptors.request.use(config => {
    const token = localStorage.getItem('access_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// 响应拦截器
api.interceptors.response.use(
    response => response,
    async error => {
        if (error.response?.status === 401) {
            await refreshToken();
            return api.request(error.config);
        }
        return Promise.reject(error);
    }
);
```

#### 会话管理
```javascript
// 获取用户会话列表
const getUserSessions = async () => {
    const response = await api.get('/user/session/info');
    return response.data;
};

// 强制下线其他设备
const revokeOtherSessions = async () => {
    const response = await api.post('/user/session/revoke-others');
    return response.data;
};

// 安全登出
const logout = async () => {
    await api.post('/user/logout');
    localStorage.clear();
    window.location.href = '/login';
};
```

### 3. 管理员后台

#### 会话监控
```javascript
// 获取活跃会话统计
const getActiveSessionStats = async () => {
    const response = await api.get('/admin/sessions/active');
    return response.data;
};

// 强制注销用户会话
const forceLogoutSession = async (sessionId) => {
    const response = await api.delete(`/admin/sessions/${sessionId}`);
    return response.data;
};

// 清理过期会话
const cleanupExpiredSessions = async () => {
    const response = await api.post('/admin/sessions/cleanup');
    return response.data;
};
```

## 🔐 安全最佳实践

### 1. Token存储
- **移动端**: 使用系统安全存储 (Keychain/Keystore)
- **Web端**: 使用HttpOnly Cookie (推荐) 或 localStorage
- **永远不要**: 在URL参数中传递Token

### 2. 网络安全
- **HTTPS**: 生产环境强制使用HTTPS
- **Token传输**: 仅通过Authorization头传递
- **防重放**: 实现请求签名或nonce机制

### 3. 会话管理
- **超时设置**: 合理设置会话超时时间
- **设备限制**: 限制单用户最大会话数
- **异常检测**: 监控异常登录行为

### 4. 错误处理
```javascript
// 统一错误处理
const handleAuthError = (error) => {
    switch (error.response?.status) {
        case 401:
            // Token过期或无效
            redirectToLogin();
            break;
        case 403:
            // 权限不足
            showErrorMessage('权限不足');
            break;
        case 429:
            // 请求过频
            showErrorMessage('请求过于频繁，请稍后再试');
            break;
        default:
            showErrorMessage('网络错误，请稍后再试');
    }
};
```

## 📊 监控和调试

### 1. 开发调试
```javascript
// 启用调试模式
localStorage.setItem('debug_session', 'true');

// 查看当前会话信息
console.log('Current Session:', {
    sessionId: localStorage.getItem('session_id'),
    accessToken: localStorage.getItem('access_token'),
    refreshToken: localStorage.getItem('refresh_token')
});
```

### 2. 生产监控
- **会话创建/销毁QPS**: 监控登录登出频率
- **Token刷新成功率**: 监控刷新机制健康度
- **异常会话检测**: 监控可疑登录行为
- **Redis性能**: 监控缓存命中率和延迟

## 🚀 部署和配置

### 1. 配置文件
```yaml
# config.yaml
session:
  redis:
    db: 1
    pool_size: 20
  settings:
    access_token_ttl: 24h
    refresh_token_ttl: 168h
    max_sessions_per_user: 5
    cleanup_interval: 1h
  security:
    validate_ip: false
    max_concurrent_sessions: 5
    auto_refresh: true
```

### 2. 环境变量
```bash
# 生产环境配置
export HUINONG_SESSION_REDIS_HOST=redis.example.com
export HUINONG_SESSION_REDIS_PASSWORD=your-redis-password
export HUINONG_JWT_SECRET_KEY=your-super-secret-key
```

通过以上优化，您的系统现在拥有了企业级的分布式会话管理能力，支持多后端实例、多端登录、实时会话控制等高级功能！ 🎉 