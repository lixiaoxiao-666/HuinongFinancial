# 数字惠农用户端 - API 使用指南

## 概述

本项目已完善了登录注册功能，集成了完整的用户认证、会话管理和Token刷新机制。

## 主要功能

### 1. 用户注册

```typescript
import { userApi } from '@/services/api'

// 用户注册
const registerUser = async () => {
  try {
    const response = await userApi.register({
      phone: '13800138000',
      password: 'password123',
      real_name: '张三'
      // 其他参数会自动填充（user_type, platform, device_info等）
    })
    
    console.log('注册成功:', response.data)
    // response.data 包含用户信息和会话信息
  } catch (error) {
    console.error('注册失败:', error)
  }
}
```

### 2. 用户登录

```typescript
import { userApi } from '@/services/api'
import { useUserStore } from '@/stores/user'

// 用户登录
const loginUser = async () => {
  try {
    const response = await userApi.login({
      phone: '13800138000',
      password: 'password123'
      // platform 和 device_info 会自动填充
    })
    
    const userStore = useUserStore()
    userStore.login(response.data) // 保存登录状态
    
    // 获取完整用户信息
    const userInfoResponse = await userApi.getUserInfo()
    userStore.setUserInfo(userInfoResponse.data)
    
  } catch (error) {
    console.error('登录失败:', error)
  }
}
```

### 3. Token 刷新

```typescript
import { authManager } from '@/utils/auth'

// 自动Token刷新（通常在应用启动时或请求拦截器中）
const checkAuth = async () => {
  const isAuthenticated = await authManager.checkAndRefreshToken()
  if (isAuthenticated) {
    console.log('用户已认证')
  } else {
    console.log('需要重新登录')
  }
}
```

### 4. 会话管理

```typescript
import { userApi } from '@/services/api'

// 获取用户所有会话
const getUserSessions = async () => {
  try {
    const response = await userApi.getUserSessions()
    console.log('用户会话列表:', response.data)
  } catch (error) {
    console.error('获取会话失败:', error)
  }
}

// 注销其他会话
const revokeOtherSessions = async () => {
  try {
    const response = await userApi.revokeOtherSessions()
    console.log('已注销其他会话，数量:', response.data.revoked_count)
  } catch (error) {
    console.error('注销失败:', error)
  }
}
```

### 5. 用户信息管理

```typescript
import { userApi } from '@/services/api'

// 更新用户信息
const updateUserProfile = async () => {
  try {
    await userApi.updateUserInfo({
      real_name: '新姓名',
      email: 'new@example.com',
      address: '新地址'
    })
    console.log('用户信息更新成功')
  } catch (error) {
    console.error('更新失败:', error)
  }
}

// 修改密码
const changePassword = async () => {
  try {
    await userApi.changePassword('旧密码', '新密码')
    console.log('密码修改成功')
  } catch (error) {
    console.error('密码修改失败:', error)
  }
}
```

### 6. 实名认证

```typescript
import { userApi, fileApi } from '@/services/api'

// 提交实名认证
const submitRealNameAuth = async () => {
  try {
    // 先上传身份证照片
    const frontFile = document.getElementById('idCardFront').files[0]
    const backFile = document.getElementById('idCardBack').files[0]
    const faceFile = document.getElementById('facePhoto').files[0]
    
    const frontResult = await fileApi.upload(frontFile, 'id_card_front')
    const backResult = await fileApi.upload(backFile, 'id_card_back')
    const faceResult = await fileApi.upload(faceFile, 'face_verify')
    
    // 提交认证申请
    const response = await userApi.submitRealNameAuth({
      id_card_number: '370123199001011234',
      real_name: '张三',
      id_card_front_img_url: frontResult.data.file_url,
      id_card_back_img_url: backResult.data.file_url,
      face_verify_img_url: faceResult.data.file_url
    })
    
    console.log('实名认证申请已提交:', response.data)
  } catch (error) {
    console.error('实名认证提交失败:', error)
  }
}
```

## 认证状态管理

使用 Pinia store 管理用户状态：

```typescript
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

// 检查登录状态
if (userStore.isLoggedIn) {
  console.log('用户已登录')
}

// 获取用户信息
const userInfo = userStore.getUserInfo
console.log('用户信息:', userInfo)

// 获取脱敏手机号
const maskedPhone = userStore.getMaskedPhone
console.log('手机号:', maskedPhone) // 138****8000

// 获取用户类型显示名称
const userTypeDisplay = userStore.getUserTypeDisplay
console.log('用户类型:', userTypeDisplay) // 农户

// 获取认证状态
const authStatus = userStore.getAuthStatus
console.log('认证状态:', authStatus)
// { realName: true, bankCard: false, credit: false }
```

## 路由守卫

应用已配置自动路由守卫，会自动检查用户认证状态：

```typescript
import { authGuard } from '@/utils/auth'

// 在 router/index.ts 中使用
router.beforeEach(authGuard)
```

## 设备信息

系统会自动检测并发送设备信息：

```typescript
import { getDeviceInfo } from '@/services/api'

const deviceInfo = getDeviceInfo()
console.log('设备信息:', deviceInfo)
// {
//   device_id: "web_abc123_1642345678",
//   device_type: "web",
//   device_name: "Chrome浏览器",
//   user_agent: "Mozilla/5.0...",
//   app_version: "1.0.0"
// }
```

## 错误处理

所有API调用都应包含错误处理：

```typescript
try {
  const response = await userApi.someMethod()
  // 处理成功响应
} catch (error) {
  if (error.message.includes('登录已过期')) {
    // Token过期，会自动跳转登录页
  } else {
    // 其他错误处理
    console.error('操作失败:', error.message)
  }
}
```

## 最佳实践

1. **自动Token刷新**：应用会在Token过期前5分钟自动刷新
2. **设备ID持久化**：设备ID存储在localStorage中，用户退出登录不会清除
3. **会话管理**：用户可以查看和管理所有活跃会话
4. **错误处理**：统一的错误处理机制
5. **类型安全**：完整的TypeScript类型定义

## 安全特性

- JWT Token认证
- Refresh Token轮换
- 设备信息记录
- 会话隔离
- 自动登出机制
- 请求拦截和响应处理
