# 会话管理系统 - API 使用指南

## 🎯 系统概述

数字惠农系统已升级为基于Redis的分布式会话管理，替代传统的JWT认证，提供更强大的会话控制和多后端支持。

### ✨ 核心特性
- 🔐 **分布式会话**: 多后端实例共享用户会话状态
- 📱 **多端登录**: 支持APP、Web、OA后台同时登录
- ⚡ **高性能**: Redis毫秒级会话验证
- 🔄 **实时同步**: 会话状态实时同步，支持强制下线
- 🛡️ **安全机制**: Token哈希存储、自动清理、设备绑定
- 📊 **智能设备识别**: 自动识别浏览器类型，优化设备名称显示

## 📋 API 接口文档

### 1. 用户认证相关

#### 1.1 用户登录
```http
POST /api/auth/login
Content-Type: application/json
X-Device-Info: {"device_type":"web","device_name":"Google Chrome","user_agent":"Mozilla/5.0...","platform":"app"}

{
    "phone": "13800138000",
    "password": "password123",
    "platform": "app",
    "device_id": "iPhone_12_ABC123",
    "device_type": "ios",
    "device_name": "张三的iPhone",
    "app_version": "1.0.0",
    "ip_address": "192.168.1.100",
    "location": "北京市朝阳区"
}
```

**请求头说明:**
- `X-Device-Info`: 自动设备信息检测（Web端自动生成）
  - `device_name`: 自动识别浏览器类型（如"Google Chrome"、"Microsoft Edge"）
  - `user_agent`: 完整的User-Agent字符串
  - `device_type`: 设备类型（web、ios、android等）
  - `platform`: 平台标识（app、web、oa等）

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
            "user_type": "farmer",
            "status": "active"
        },
        "session": {
            "session_id": "sess_abc123def456",
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "expires_in": 86400,
            "platform": "app",
            "created_at": "2024-01-15T10:30:00Z",
            "expires_at": "2024-01-16T10:30:00Z"
        }
    }
}
```

#### 1.2 OA管理员登录
```http
POST /api/oa/auth/login
Content-Type: application/json
X-Device-Info: {"device_type":"web","device_name":"Microsoft Edge","user_agent":"Mozilla/5.0...","platform":"oa"}

{
    "username": "admin",
    "password": "admin123",
    "platform": "oa",
    "device_type": "web",
    "device_name": "管理后台",
    "device_id": "oa_web_1640995200000",
    "app_version": "1.0.0"
}
```

**响应示例:**
```json
{
    "code": 200,
    "message": "登录成功",
    "data": {
        "user": {
            "id": 1,
            "username": "admin",
            "email": "admin@huinong.com",
            "real_name": "超级管理员",
            "role_id": 1,
            "department": "技术部",
            "position": "系统管理员",
            "status": "active",
            "permissions": ["*"]
        },
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expires_in": 86400,
        "session_id": "sess_oa_1640995200000_admin"
    }
}
```

**平台标识说明:**
- `app`: 移动端APP
- `web`: 普通Web端
- `oa`: OA管理后台
- `mini`: 微信小程序

#### 1.3 通用Token刷新
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

#### 1.4 OA专用Token刷新
```http
POST /api/oa/auth/refresh
Content-Type: application/json

{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
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

#### 1.5 会话验证（通用）
```http
GET /api/auth/validate
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "会话有效",
    "data": {
        "session_id": "sess_abc123def456",
        "user_id": 1001,
        "platform": "app",
        "status": "active",
        "expires_at": "2024-01-16T10:30:00Z",
        "last_active_at": "2024-01-15T14:25:30Z"
    }
}
```

#### 1.6 OA会话验证
```http
GET /api/oa/auth/validate
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "OA会话有效",
    "data": {
        "session_id": "sess_oa_1640995200000_admin",
        "user_id": 1,
        "username": "admin",
        "platform": "oa",
        "status": "active",
        "permissions": ["*"],
        "expires_at": "2024-01-16T10:30:00Z"
    }
}
```

### 2. 会话管理相关

#### 2.1 获取当前用户会话信息
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
                "device_name": "张三的iPhone",
                "app_version": "1.0.0",
                "user_agent": "HuinongApp/1.0.0 (iPhone; iOS 15.0)"
            },
            "network_info": {
                "ip_address": "192.168.1.100",
                "location": "北京市朝阳区"
            },
            "status": "active",
            "created_at": "2024-01-15T10:30:00Z",
            "last_active_at": "2024-01-15T14:25:30Z",
            "is_current": true
        },
        {
            "session_id": "sess_xyz789uvw012",
            "platform": "web",
            "device_info": {
                "device_id": "web_chrome_123456",
                "device_type": "web",
                "device_name": "Google Chrome",
                "app_version": "1.0.0",
                "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
            },
            "network_info": {
                "ip_address": "192.168.1.101",
                "location": "北京市朝阳区"
            },
            "status": "active",
            "created_at": "2024-01-15T08:15:00Z",
            "last_active_at": "2024-01-15T12:30:00Z",
            "is_current": false
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
    "message": "注销成功",
    "data": {
        "session_id": "sess_abc123def456",
        "logout_time": "2024-01-15T15:00:00Z"
    }
}
```

#### 2.3 OA管理员注销
```http
POST /api/oa/auth/logout
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "OA注销成功",
    "data": {
        "session_id": "sess_oa_1640995200000_admin",
        "logout_time": "2024-01-15T15:00:00Z"
    }
}
```

#### 2.4 注销其他设备会话
```http
POST /api/user/session/revoke-others
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "注销其他会话成功",
    "data": {
        "current_session": "sess_abc123def456",
        "revoked_sessions": [
            "sess_xyz789uvw012",
            "sess_def456ghi789"
        ],
        "revoked_count": 2
    }
}
```

#### 2.5 注销指定会话
```http
DELETE /api/user/session/{session_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "会话注销成功",
    "data": {
        "session_id": "sess_xyz789uvw012",
        "revoked_at": "2024-01-15T15:00:00Z"
    }
}
```

### 3. 管理员会话管理

#### 3.1 获取所有活跃会话
```http
GET /api/admin/sessions/active?page=1&limit=20&platform=all&status=active
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**查询参数:**
- `page`: 页码（默认1）
- `limit`: 每页数量（默认20，最大100）
- `platform`: 平台筛选（all、app、web、oa）
- `status`: 状态筛选（active、expired、revoked）
- `user_id`: 用户ID筛选
- `device_type`: 设备类型筛选

**响应示例:**
```json
{
    "code": 200,
    "message": "获取活跃会话列表成功",
    "data": {
        "total": 150,
        "page": 1,
        "limit": 20,
        "sessions": [
            {
                "session_id": "sess_abc123",
                "user_id": 1001,
                "username": "13800138000",
                "real_name": "张三",
                "platform": "app",
                "device_info": {
                    "device_type": "ios",
                    "device_name": "张三的iPhone",
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
}
```

#### 3.2 强制注销会话
```http
DELETE /api/admin/sessions/{session_id}
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "会话已强制注销",
    "data": {
        "session_id": "sess_abc123def456",
        "user_id": 1001,
        "revoked_at": "2024-01-15T15:00:00Z",
        "revoked_by": 1
    }
}
```

#### 3.3 强制注销用户所有会话
```http
DELETE /api/admin/users/{user_id}/sessions
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "用户所有会话已强制注销",
    "data": {
        "user_id": 1001,
        "revoked_sessions": [
            "sess_abc123def456",
            "sess_xyz789uvw012"
        ],
        "revoked_count": 2,
        "revoked_at": "2024-01-15T15:00:00Z"
    }
}
```

#### 3.4 清理过期会话
```http
POST /api/admin/sessions/cleanup
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "过期会话清理完成",
    "data": {
        "cleaned_sessions": 23,
        "cleanup_time": "2024-01-15T15:00:00Z",
        "next_cleanup": "2024-01-15T16:00:00Z"
    }
}
```

#### 3.5 会话统计信息
```http
GET /api/admin/sessions/statistics
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例:**
```json
{
    "code": 200,
    "message": "获取会话统计成功",
    "data": {
        "total_sessions": 1250,
        "active_sessions": 567,
        "platform_stats": {
            "app": 423,
            "web": 134,
            "oa": 10
        },
        "device_stats": {
            "ios": 234,
            "android": 189,
            "web": 144
        },
        "recent_24h": {
            "new_sessions": 89,
            "expired_sessions": 45,
            "revoked_sessions": 12
        }
    }
}
```

## 🔧 客户端集成指南

### 1. 移动端 (APP)

#### 智能设备信息获取
```javascript
// 获取设备信息
const getDeviceInfo = () => {
    const platform = Platform.OS; // 'ios' or 'android'
    const deviceName = `${DeviceInfo.getModel()} - ${DeviceInfo.getDeviceId()}`;
    
    return {
        platform: 'app',
        device_id: DeviceInfo.getUniqueId(),
        device_type: platform,
        device_name: deviceName,
        app_version: DeviceInfo.getVersion(),
        user_agent: DeviceInfo.getUserAgent()
    };
};

// 登录流程
const loginUser = async (phone, password) => {
    const deviceInfo = getDeviceInfo();
    const location = await getCurrentLocation();
    
    const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            phone,
            password,
            ...deviceInfo,
            ip_address: await getClientIP(),
            location: location.address
        })
    });

    const result = await response.json();
    
    if (result.code === 200) {
        // 存储会话信息
        await AsyncStorage.setItem('access_token', result.data.session.access_token);
        await AsyncStorage.setItem('refresh_token', result.data.session.refresh_token);
        await AsyncStorage.setItem('session_id', result.data.session.session_id);
        await AsyncStorage.setItem('user_info', JSON.stringify(result.data.user));
        
        return result.data;
    } else {
        throw new Error(result.message);
    }
};
```

#### 自动Token刷新机制
```javascript
import { createAuthRefreshInterceptor } from 'axios-auth-refresh';

// 创建axios实例
const apiClient = axios.create({
    baseURL: API_BASE_URL,
    timeout: 10000
});

// Token刷新函数
const refreshAuthLogic = async (failedRequest) => {
    const refreshToken = await AsyncStorage.getItem('refresh_token');
    
    if (!refreshToken) {
        throw new Error('No refresh token available');
    }

    try {
        const response = await fetch('/api/auth/refresh', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `refresh_token=${refreshToken}`
        });

        const result = await response.json();
        
        if (result.code === 200) {
            // 更新存储的token
            await AsyncStorage.setItem('access_token', result.data.access_token);
            await AsyncStorage.setItem('refresh_token', result.data.refresh_token);
            
            // 更新失败请求的header
            failedRequest.response.config.headers.Authorization = `Bearer ${result.data.access_token}`;
            
            return Promise.resolve();
        } else {
            throw new Error('Token refresh failed');
        }
    } catch (error) {
        // 刷新失败，跳转到登录页
        await clearAuthStorage();
        NavigationService.navigate('Login');
        throw error;
    }
};

// 设置刷新拦截器
createAuthRefreshInterceptor(apiClient, refreshAuthLogic);

// 请求拦截器
apiClient.interceptors.request.use(async (config) => {
    const token = await AsyncStorage.getItem('access_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});
```

### 2. Web端（Vue3 + Element Plus）

#### 优化的设备信息检测
```javascript
// utils/deviceDetection.js
export const getDeviceInfo = () => {
    const ua = navigator.userAgent;
    
    // 浏览器检测
    const getBrowserName = () => {
        if (ua.includes('Edg/')) return 'Microsoft Edge';
        if (ua.includes('Chrome/') && !ua.includes('Edg/')) return 'Google Chrome';
        if (ua.includes('Firefox/')) return 'Mozilla Firefox';
        if (ua.includes('Safari/') && !ua.includes('Chrome/') && !ua.includes('Edg/')) return 'Apple Safari';
        if (ua.includes('Opera/') || ua.includes('OPR/')) return 'Opera';
        return 'Unknown Browser';
    };

    // 操作系统检测
    const getOSName = () => {
        if (ua.includes('Windows NT')) return 'Windows';
        if (ua.includes('Mac OS X')) return 'macOS';
        if (ua.includes('Linux')) return 'Linux';
        if (ua.includes('Android')) return 'Android';
        if (ua.includes('iPhone') || ua.includes('iPad')) return 'iOS';
        return 'Unknown OS';
    };

    return {
        device_type: 'web',
        device_name: `${getBrowserName()} on ${getOSName()}`,
        user_agent: ua,
        platform: 'web',
        device_id: generateWebDeviceId(),
        app_version: '1.0.0'
    };
};

// 生成Web设备ID
const generateWebDeviceId = () => {
    let deviceId = localStorage.getItem('web_device_id');
    if (!deviceId) {
        deviceId = `web_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
        localStorage.setItem('web_device_id', deviceId);
    }
    return deviceId;
};
```

#### 增强的请求拦截器
```javascript
// utils/request.js
import axios from 'axios';
import { ElMessage, ElLoading } from 'element-plus';
import { useAuthStore } from '@/stores/auth';
import { getDeviceInfo } from '@/utils/deviceDetection';

const request = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '',
    timeout: 30000,
    headers: {
        'Content-Type': 'application/json'
    }
});

// 请求拦截器
request.interceptors.request.use((config) => {
    // 添加Token
    const authStore = useAuthStore();
    if (authStore.accessToken && !isPublicEndpoint(config.url)) {
        config.headers.Authorization = `Bearer ${authStore.accessToken}`;
    }

    // 添加设备信息
    const deviceInfo = getDeviceInfo();
    config.headers['X-Device-Info'] = JSON.stringify(deviceInfo);
    config.headers['X-Request-Time'] = Date.now().toString();

    return config;
});

// 响应拦截器 - 增强错误处理
request.interceptors.response.use(
    (response) => response,
    async (error) => {
        const { response } = error;
        
        if (response?.status === 401) {
            // 尝试刷新Token
            const refreshed = await handleTokenRefresh();
            if (refreshed) {
                // 重试原请求
                return request.request(error.config);
            } else {
                handleAuthError('登录已过期，请重新登录');
            }
        }
        
        // 其他错误处理
        handleAPIError(response);
        return Promise.reject(error);
    }
);

// Token刷新处理
const handleTokenRefresh = async () => {
    try {
        const authStore = useAuthStore();
        if (!authStore.refreshToken) return false;

        const response = await axios.post('/api/auth/refresh', 
            `refresh_token=${authStore.refreshToken}`,
            {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            }
        );

        if (response.data.code === 200) {
            authStore.updateTokens(
                response.data.data.access_token,
                response.data.data.refresh_token
            );
            return true;
        }
        return false;
    } catch (error) {
        console.error('Token刷新失败:', error);
        return false;
    }
};

// API错误处理
const handleAPIError = (response) => {
    const errorMessages = {
        400: '请求参数错误',
        401: '登录已过期',
        403: '权限不足',
        404: '请求的资源不存在',
        429: '请求过于频繁，请稍后再试',
        500: '服务器内部错误',
        502: '网关错误',
        503: '服务暂时不可用',
        504: '网关超时'
    };

    const message = response?.data?.message || errorMessages[response?.status] || '网络错误';
    ElMessage.error(message);
};
```

### 3. 管理员后台 (OA)

#### OA专用认证处理
```javascript
// stores/oaAuth.js
import { defineStore } from 'pinia';
import { oaLogin, oaRefreshToken, oaLogout, oaValidate } from '@/api/oa';

export const useOAAuthStore = defineStore('oaAuth', {
    state: () => ({
        user: null,
        accessToken: localStorage.getItem('oa_access_token'),
        refreshToken: localStorage.getItem('oa_refresh_token'),
        sessionId: localStorage.getItem('oa_session_id'),
        isAuthenticated: false
    }),

    actions: {
        async login(credentials) {
            try {
                const response = await oaLogin(credentials);
                
                this.user = response.data.user;
                this.accessToken = response.data.access_token;
                this.refreshToken = response.data.refresh_token;
                this.sessionId = response.data.session_id;
                this.isAuthenticated = true;

                // 持久化存储
                localStorage.setItem('oa_access_token', this.accessToken);
                localStorage.setItem('oa_refresh_token', this.refreshToken);
                localStorage.setItem('oa_session_id', this.sessionId);
                localStorage.setItem('oa_user', JSON.stringify(this.user));

                return response;
            } catch (error) {
                this.clearAuth();
                throw error;
            }
        },

        async refreshTokens() {
            try {
                const response = await oaRefreshToken({
                    refresh_token: this.refreshToken
                });

                this.accessToken = response.data.access_token;
                this.refreshToken = response.data.refresh_token;

                localStorage.setItem('oa_access_token', this.accessToken);
                localStorage.setItem('oa_refresh_token', this.refreshToken);

                return true;
            } catch (error) {
                this.clearAuth();
                return false;
            }
        },

        async logout() {
            try {
                await oaLogout();
            } catch (error) {
                console.error('OA注销失败:', error);
            } finally {
                this.clearAuth();
            }
        },

        async validateSession() {
            try {
                const response = await oaValidate();
                this.user = { ...this.user, ...response.data };
                return true;
            } catch (error) {
                this.clearAuth();
                return false;
            }
        },

        clearAuth() {
            this.user = null;
            this.accessToken = null;
            this.refreshToken = null;
            this.sessionId = null;
            this.isAuthenticated = false;

            localStorage.removeItem('oa_access_token');
            localStorage.removeItem('oa_refresh_token');
            localStorage.removeItem('oa_session_id');
            localStorage.removeItem('oa_user');
        },

        checkPermission(permission) {
            return this.user?.permissions?.includes('*') || 
                   this.user?.permissions?.includes(permission) ||
                   this.user?.role?.name === 'super_admin';
        }
    }
});
```

## 🔐 安全最佳实践

### 1. 强化的Token管理
```javascript
// 安全Token存储
class SecureTokenStorage {
    // 移动端：使用安全存储
    static async setTokenSecure(key, value) {
        if (Platform.OS === 'ios' || Platform.OS === 'android') {
            return await Keychain.setInternetCredentials(key, 'token', value);
        } else {
            // Web端：使用加密存储
            const encrypted = CryptoJS.AES.encrypt(value, 'secret-key').toString();
            localStorage.setItem(key, encrypted);
        }
    }

    static async getTokenSecure(key) {
        if (Platform.OS === 'ios' || Platform.OS === 'android') {
            const credentials = await Keychain.getInternetCredentials(key);
            return credentials ? credentials.password : null;
        } else {
            const encrypted = localStorage.getItem(key);
            if (encrypted) {
                const decrypted = CryptoJS.AES.decrypt(encrypted, 'secret-key');
                return decrypted.toString(CryptoJS.enc.Utf8);
            }
            return null;
        }
    }
}
```

### 2. 增强的会话监控
```javascript
// 会话监控和异常检测
class SessionMonitor {
    constructor() {
        this.lastActivity = Date.now();
        this.activityThreshold = 30 * 60 * 1000; // 30分钟
        this.init();
    }

    init() {
        // 监听用户活动
        ['mousedown', 'keydown', 'scroll', 'touchstart'].forEach(event => {
            document.addEventListener(event, () => {
                this.updateActivity();
            }, true);
        });

        // 定期检查会话状态
        setInterval(() => {
            this.checkSession();
        }, 5 * 60 * 1000); // 每5分钟检查一次
    }

    updateActivity() {
        this.lastActivity = Date.now();
    }

    async checkSession() {
        const inactive = Date.now() - this.lastActivity > this.activityThreshold;
        
        if (!inactive) {
            try {
                await api.get('/api/auth/validate');
            } catch (error) {
                if (error.response?.status === 401) {
                    this.handleSessionExpired();
                }
            }
        }
    }

    handleSessionExpired() {
        ElMessageBox.alert('会话已过期，请重新登录', '提示', {
            type: 'warning',
            callback: () => {
                this.redirectToLogin();
            }
        });
    }

    redirectToLogin() {
        const authStore = useAuthStore();
        authStore.logout();
        router.push('/login');
    }
}

// 启动会话监控
const sessionMonitor = new SessionMonitor();
```

### 3. 设备指纹增强
```javascript
// 设备指纹生成
class DeviceFingerprint {
    static async generate() {
        const canvas = document.createElement('canvas');
        const ctx = canvas.getContext('2d');
        ctx.textBaseline = 'top';
        ctx.font = '14px Arial';
        ctx.fillText('Device fingerprint', 2, 2);
        
        const fingerprint = {
            screen: `${screen.width}x${screen.height}x${screen.colorDepth}`,
            timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
            language: navigator.language,
            platform: navigator.platform,
            canvas: canvas.toDataURL(),
            webgl: this.getWebGLFingerprint(),
            plugins: Array.from(navigator.plugins).map(p => p.name).join(','),
            fonts: await this.detectFonts()
        };

        return btoa(JSON.stringify(fingerprint)).slice(0, 32);
    }

    static getWebGLFingerprint() {
        const gl = document.createElement('canvas').getContext('webgl');
        if (!gl) return 'no-webgl';
        
        return [
            gl.getParameter(gl.VENDOR),
            gl.getParameter(gl.RENDERER),
            gl.getParameter(gl.VERSION)
        ].join('|');
    }

    static async detectFonts() {
        const fonts = ['Arial', 'Helvetica', 'Times', 'Courier', 'Verdana', 'Georgia'];
        const detected = [];
        
        for (const font of fonts) {
            if (document.fonts.check(`12px ${font}`)) {
                detected.push(font);
            }
        }
        
        return detected.join(',');
    }
}
```

## 📊 监控和调试

### 1. 会话状态监控
```javascript
// 会话状态实时监控
class SessionStatsMonitor {
    constructor() {
        this.stats = {
            totalRequests: 0,
            successfulRequests: 0,
            failedRequests: 0,
            tokenRefreshes: 0,
            sessionErrors: 0
        };
    }

    trackRequest(success, error) {
        this.stats.totalRequests++;
        
        if (success) {
            this.stats.successfulRequests++;
        } else {
            this.stats.failedRequests++;
            
            if (error?.response?.status === 401) {
                this.stats.sessionErrors++;
            }
        }

        this.reportStats();
    }

    trackTokenRefresh() {
        this.stats.tokenRefreshes++;
        this.reportStats();
    }

    reportStats() {
        // 发送统计数据到监控系统
        if (this.stats.totalRequests % 100 === 0) {
            console.log('Session Stats:', this.stats);
            
            // 可以发送到分析服务
            // analytics.track('session_stats', this.stats);
        }
    }
}

const sessionStatsMonitor = new SessionStatsMonitor();
```

### 2. 开发调试工具
```javascript
// 开发环境调试工具
if (process.env.NODE_ENV === 'development') {
    window.SessionDebug = {
        // 查看当前会话信息
        getCurrentSession: async () => {
            try {
                const response = await api.get('/api/auth/validate');
                console.table(response.data);
                return response.data;
            } catch (error) {
                console.error('获取会话信息失败:', error);
            }
        },

        // 查看存储的Token
        getStoredTokens: () => {
            const tokens = {
                access_token: localStorage.getItem('access_token'),
                refresh_token: localStorage.getItem('refresh_token'),
                session_id: localStorage.getItem('session_id')
            };
            console.table(tokens);
            return tokens;
        },

        // 强制刷新Token
        forceRefreshToken: async () => {
            try {
                const refreshToken = localStorage.getItem('refresh_token');
                const response = await axios.post('/api/auth/refresh', 
                    `refresh_token=${refreshToken}`,
                    {
                        headers: {
                            'Content-Type': 'application/x-www-form-urlencoded'
                        }
                    }
                );
                console.log('Token刷新成功:', response.data);
                return response.data;
            } catch (error) {
                console.error('Token刷新失败:', error);
            }
        },

        // 清理所有会话数据
        clearAllSessions: () => {
            localStorage.clear();
            sessionStorage.clear();
            console.log('所有会话数据已清理');
        }
    };

    console.log('🔧 会话调试工具已加载，使用 SessionDebug.* 方法进行调试');
}
```

## 🚀 部署和配置

### 1. 生产环境配置
```yaml
# production.yaml
session:
  redis:
    # Redis集群配置
    cluster:
      enabled: true
      nodes:
        - "redis-node1:6379"
        - "redis-node2:6379"
        - "redis-node3:6379"
    password: "${REDIS_PASSWORD}"
    
    # 连接池优化
    pool_size: 50
    max_retries: 3
    dial_timeout: 5s
    read_timeout: 3s
    write_timeout: 3s

  settings:
    # 安全配置
    access_token_ttl: 2h          # 生产环境缩短访问令牌时间
    refresh_token_ttl: 24h        # 刷新令牌时间
    max_sessions_per_user: 3      # 限制并发会话数
    
    # 清理策略
    cleanup_interval: 30m
    batch_cleanup_size: 500
    
    # 安全特性
    validate_device_fingerprint: true
    enable_suspicious_activity_detection: true
    max_login_attempts: 5
    lockout_duration: 15m

  security:
    # IP验证
    validate_ip: true
    allow_ip_change: false
    
    # 设备绑定
    validate_device_id: true
    allow_device_change: false
    
    # 风险控制
    enable_risk_detection: true
    risk_threshold: 80
    
    # 自动续期
    auto_refresh: true
    refresh_threshold: 0.2  # 剩余20%时间时自动刷新
```

### 2. 监控和告警
```yaml
# 监控配置
monitoring:
  metrics:
    - name: session_creation_rate
      description: "会话创建速率"
      type: counter
      
    - name: session_validation_latency
      description: "会话验证延迟"
      type: histogram
      
    - name: token_refresh_success_rate
      description: "Token刷新成功率"
      type: gauge
      
    - name: redis_connection_pool_usage
      description: "Redis连接池使用率"
      type: gauge

  alerts:
    - name: high_session_creation_rate
      condition: "session_creation_rate > 1000/min"
      message: "会话创建速率异常"
      
    - name: low_token_refresh_success_rate
      condition: "token_refresh_success_rate < 0.95"
      message: "Token刷新成功率过低"
      
    - name: redis_connection_pool_exhausted
      condition: "redis_connection_pool_usage > 0.9"
      message: "Redis连接池即将耗尽"
```

通过以上完善的API文档和集成指南，您的Redis会话同步系统现在具备了企业级的功能和安全性，支持多端、多后端、智能设备识别和完善的监控体系！ 🎉 