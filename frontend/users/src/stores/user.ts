import { defineStore } from 'pinia'
import { ref } from 'vue'

// 用户信息接口
export interface UserInfo {
  user_id: string
  phone: string
  nickname: string
  avatar_url: string
  real_name: string
  id_card_number: string
  address: string
}

// 登录响应接口
export interface LoginResponse {
  user_id: string
  token: string
  expires_in: number
  user_type: string
}

export const useUserStore = defineStore('user', {
    // 状态
    state: () => ({
        token: ref(''),
        userInfo: ref<UserInfo | null>(null),
        isLoggedIn: ref(false),
        loginExpireTime: ref<number>(0)
    }),
    // 操作
    actions: {
        // 设置token
        setToken(token: string, expiresIn?: number) {
            this.token = token
            this.isLoggedIn = true
            if (expiresIn) {
                this.loginExpireTime = Date.now() + expiresIn * 1000
            }
            // 存储到localStorage
            localStorage.setItem('token', token)
            localStorage.setItem('loginExpireTime', this.loginExpireTime.toString())
        },
        
        // 设置用户信息
        setUserInfo(userInfo: UserInfo) {
            this.userInfo = userInfo
            // 存储到localStorage
            localStorage.setItem('userInfo', JSON.stringify(userInfo))
        },
        
        // 登录
        login(loginData: LoginResponse) {
            this.setToken(loginData.token, loginData.expires_in)
            // 可以设置基本的用户ID信息
            if (this.userInfo) {
                this.userInfo.user_id = loginData.user_id
            }
        },
        
        // 登出
        logout() {
            this.token = ''
            this.userInfo = null
            this.isLoggedIn = false
            this.loginExpireTime = 0
            // 清除localStorage
            localStorage.removeItem('token')
            localStorage.removeItem('userInfo')
            localStorage.removeItem('loginExpireTime')
        },
        
        // 从localStorage恢复状态
        restoreFromStorage() {
            const token = localStorage.getItem('token')
            const userInfo = localStorage.getItem('userInfo')
            const expireTime = localStorage.getItem('loginExpireTime')
            
            if (token && expireTime) {
                const expireTimeNum = parseInt(expireTime)
                // 检查token是否过期
                if (Date.now() < expireTimeNum) {
                    this.token = token
                    this.isLoggedIn = true
                    this.loginExpireTime = expireTimeNum
                } else {
                    // token过期，清除
                    this.logout()
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
        }
    },
    // 计算
    getters: {
        getToken: (state) => state.token,
        getUserInfo: (state) => state.userInfo,
        getIsLoggedIn: (state) => state.isLoggedIn,
        // 获取遮盖后的手机号
        getMaskedPhone: (state) => {
            if (state.userInfo?.phone) {
                const phone = state.userInfo.phone
                return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
            }
            return ''
        }
    },
})
