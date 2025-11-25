import { defineStore } from 'pinia'
import { login, logout, getCurrentUser } from '@/api/auth'
import router from '@/router'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('access_token') || '',
    refreshToken: localStorage.getItem('refresh_token') || '',
    userInfo: JSON.parse(localStorage.getItem('user_info') || 'null')
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
    username: (state) => state.userInfo?.username || '',
    nickname: (state) => state.userInfo?.nickname || state.userInfo?.username || '',
    avatar: (state) => state.userInfo?.avatar || '',
    roles: (state) => state.userInfo?.roles || [],
    isAdmin: (state) => state.userInfo?.roles?.some(role => role.name === 'admin') || false
  },

  actions: {
    // 登录
    async login(loginForm) {
      try {
        const response = await login(loginForm)
        const { access_token, refresh_token, user } = response.data
        
        this.token = access_token
        this.refreshToken = refresh_token
        this.userInfo = user
        
        // 保存到 localStorage
        localStorage.setItem('access_token', access_token)
        localStorage.setItem('refresh_token', refresh_token)
        localStorage.setItem('user_info', JSON.stringify(user))
        
        return response
      } catch (error) {
        console.error('Login failed:', error)
        throw error
      }
    },

    // 登出
    async logout() {
      try {
        await logout()
      } catch (error) {
        console.error('Logout failed:', error)
      } finally {
        this.token = ''
        this.refreshToken = ''
        this.userInfo = null
        
        // 清除 localStorage
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        localStorage.removeItem('user_info')
        
        // 跳转到登录页
        router.push('/login')
      }
    },

    // 获取用户信息
    async getUserInfo() {
      try {
        const response = await getCurrentUser()
        this.userInfo = response.data
        localStorage.setItem('user_info', JSON.stringify(response.data))
        return response
      } catch (error) {
        console.error('Get user info failed:', error)
        throw error
      }
    },

    // 清除用户信息
    clearUserInfo() {
      this.token = ''
      this.refreshToken = ''
      this.userInfo = null
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('user_info')
    }
  }
})

