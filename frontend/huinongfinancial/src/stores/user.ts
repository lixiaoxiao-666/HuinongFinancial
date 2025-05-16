import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', {
    // 状态
    state: () => ({
        token: ref(''),
        userInfo: ref({}),
    }),
    // 操作
    actions: {
        setToken(token: string) {
            this.token = token
        },
        setUserInfo(userInfo: any) {
            this.userInfo = userInfo
        },
    },
    // 计算
    getters: {
        getToken: (state) => state.token,
        getUserInfo: (state) => state.userInfo,
    },
})
