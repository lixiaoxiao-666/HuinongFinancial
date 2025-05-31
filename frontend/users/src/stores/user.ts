import { defineStore } from 'pinia'
import { ref } from 'vue'

// 用户信息接口 - 根据API文档更新
export interface UserInfo {
  id: number
  uuid: string
  username?: string
  phone: string
  email?: string
  user_type: string
  status: string
  real_name: string
  id_card?: string // 脱敏显示
  avatar?: string
  gender?: string
  birthday?: string
  province?: string
  city?: string
  county?: string
  address?: string
  is_real_name_verified?: boolean
  is_bank_card_verified?: boolean
  is_credit_verified?: boolean
  last_login_time?: string
  created_at?: string
}

// 登录响应接口 - 保持向后兼容
export interface LoginResponse {
  user_id: string
  token: string
  expires_in: number
  user_type: string
}

// 新的登录API响应类型
export interface LoginApiResponse {
  user: {
    id: number
    uuid: string
    phone: string
    user_type: string
    real_name: string
    status: string
  }
  session: {
    access_token: string
    refresh_token: string
    expires_in: number
  }
}

export const useUserStore = defineStore('user', {
    // 状态
    state: () => ({
        token: ref(''),
        refreshToken: ref(''), // 添加刷新token
        userInfo: ref<UserInfo | null>(null),
        isLoggedIn: ref(false),
        loginExpireTime: ref<number>(0)
    }),
    // 操作
    actions: {
        // 设置token
        setToken(accessToken: string, refreshToken?: string, expiresIn?: number) {
            this.token = accessToken
            if (refreshToken) {
                this.refreshToken = refreshToken
            }
            this.isLoggedIn = true
            if (expiresIn) {
                this.loginExpireTime = Date.now() + expiresIn * 1000
            }
            // 存储到localStorage
            localStorage.setItem('access_token', accessToken)
            if (refreshToken) {
                localStorage.setItem('refresh_token', refreshToken)
            }
            localStorage.setItem('loginExpireTime', this.loginExpireTime.toString())
        },
        
        // 设置用户信息
        setUserInfo(userInfo: UserInfo) {
            this.userInfo = userInfo
            // 存储到localStorage
            localStorage.setItem('userInfo', JSON.stringify(userInfo))
        },
        
        // 登录 - 适配新的API响应格式
        login(loginData: LoginApiResponse | LoginResponse) {
            // 判断是新格式还是旧格式，保持向后兼容
            if ('session' in loginData) {
                // 新API格式
                this.setToken(
                    loginData.session.access_token, 
                    loginData.session.refresh_token, 
                    loginData.session.expires_in
                )
                // 设置基本用户信息
                if (loginData.user) {
                    const basicUserInfo: UserInfo = {
                        id: loginData.user.id,
                        uuid: loginData.user.uuid,
                        phone: loginData.user.phone,
                        user_type: loginData.user.user_type,
                        real_name: loginData.user.real_name,
                        status: loginData.user.status
                    }
                    this.setUserInfo(basicUserInfo)
                }
            } else {
                // 兼容旧格式
                this.setToken(loginData.token, undefined, loginData.expires_in)
                if (this.userInfo) {
                    this.userInfo.id = parseInt(loginData.user_id)
                }
            }
        },
        
        // 刷新token
        async refreshAccessToken(): Promise<boolean> {
            if (!this.refreshToken) {
                return false
            }
            
            try {
                // 这里需要导入API，为了避免循环依赖，我们在组件中处理
                // 或者将刷新逻辑移到API层
                return true
            } catch (error) {
                console.error('Token刷新失败:', error)
                this.logout()
                return false
            }
        },
        
        // 登出
        logout() {
            this.token = ''
            this.refreshToken = ''
            this.userInfo = null
            this.isLoggedIn = false
            this.loginExpireTime = 0
            // 清除localStorage
            localStorage.removeItem('access_token')
            localStorage.removeItem('refresh_token')
            localStorage.removeItem('userInfo')
            localStorage.removeItem('loginExpireTime')
            // 保持设备ID，用户下次登录时使用
            // localStorage.removeItem('device_id') // 不清除设备ID
        },
        
        // 从localStorage恢复状态
        restoreFromStorage() {
            const token = localStorage.getItem('access_token') || localStorage.getItem('token') // 兼容旧的token key
            const refreshToken = localStorage.getItem('refresh_token')
            const userInfo = localStorage.getItem('userInfo')
            const expireTime = localStorage.getItem('loginExpireTime')
            
            if (token && expireTime) {
                const expireTimeNum = parseInt(expireTime)
                // 检查token是否过期
                if (Date.now() < expireTimeNum) {
                    this.token = token
                    if (refreshToken) {
                        this.refreshToken = refreshToken
                    }
                    this.isLoggedIn = true
                    this.loginExpireTime = expireTimeNum
                } else {
                    // access token过期，但如果有refresh token，可以尝试刷新
                    if (refreshToken) {
                        this.refreshToken = refreshToken
                        // 这里可以触发自动刷新逻辑
                    } else {
                        // 完全过期，清除
                        this.logout()
                    }
                }
            }
            
            if (userInfo) {
                try {
                    this.userInfo = JSON.parse(userInfo)
                } catch (e) {
                    console.error('解析用户信息失败:', e)
                }
            }
        },
        
        // 检查token是否有效
        isTokenValid(): boolean {
            return this.isLoggedIn && this.token !== '' && Date.now() < this.loginExpireTime
        },

        // 检查是否需要刷新token（在过期前5分钟刷新）
        shouldRefreshToken(): boolean {
            const fiveMinutes = 5 * 60 * 1000
            return this.isLoggedIn && this.refreshToken !== '' && Date.now() > (this.loginExpireTime - fiveMinutes)
        },

        // 更新部分用户信息
        updateUserInfo(updates: Partial<UserInfo>) {
            if (this.userInfo) {
                this.userInfo = { ...this.userInfo, ...updates }
                localStorage.setItem('userInfo', JSON.stringify(this.userInfo))
            }
        }
    },
    // 计算
    getters: {
        getToken: (state) => state.token,
        getRefreshToken: (state) => state.refreshToken,
        getUserInfo: (state) => state.userInfo,
        getIsLoggedIn: (state) => state.isLoggedIn,
        // 获取遮盖后的手机号
        getMaskedPhone: (state) => {
            if (state.userInfo?.phone) {
                const phone = state.userInfo.phone
                return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
            }
            return ''
        },
        // 获取用户类型显示名称
        getUserTypeDisplay: (state) => {
            if (!state.userInfo?.user_type) return ''
            const userTypeMap: Record<string, string> = {
                'farmer': '农户',
                'farm_owner': '农场主', 
                'cooperative': '合作社',
                'enterprise': '企业'
            }
            return userTypeMap[state.userInfo.user_type] || state.userInfo.user_type
        },
        // 获取认证状态
        getAuthStatus: (state) => {
            if (!state.userInfo) return null
            return {
                realName: state.userInfo.is_real_name_verified || false,
                bankCard: state.userInfo.is_bank_card_verified || false,
                credit: state.userInfo.is_credit_verified || false
            }
        }
    },
})
