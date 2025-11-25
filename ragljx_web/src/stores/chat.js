import { defineStore } from 'pinia'
import { getChatSessions, getChatMessages } from '@/api/chat'

export const useChatStore = defineStore('chat', {
  state: () => ({
    sessions: [],
    currentSession: null,
    messages: [],
    loading: false
  }),

  getters: {
    sessionList: (state) => state.sessions,
    currentSessionId: (state) => state.currentSession?.id,
    messageList: (state) => state.messages
  },

  actions: {
    // 获取会话列表
    async fetchSessions(params = {}) {
      this.loading = true
      try {
        const response = await getChatSessions(params)
        this.sessions = response.data?.items || []
        return response
      } catch (error) {
        console.error('Fetch sessions failed:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    // 获取消息列表
    async fetchMessages(sessionId, params = {}) {
      this.loading = true
      try {
        const response = await getChatMessages(sessionId, params)
        this.messages = response.data?.items || []
        return response
      } catch (error) {
        console.error('Fetch messages failed:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    // 设置当前会话
    setCurrentSession(session) {
      this.currentSession = session
    },

    // 清除当前会话
    clearCurrentSession() {
      this.currentSession = null
      this.messages = []
    },

    // 添加会话
    addSession(session) {
      this.sessions.unshift(session)
    },

    // 删除会话
    removeSession(sessionId) {
      const index = this.sessions.findIndex(s => s.id === sessionId)
      if (index !== -1) {
        this.sessions.splice(index, 1)
      }
    },

    // 添加消息
    addMessage(message) {
      this.messages.push(message)
    },

    // 更新消息
    updateMessage(messageId, data) {
      const index = this.messages.findIndex(m => m.id === messageId)
      if (index !== -1) {
        this.messages[index] = { ...this.messages[index], ...data }
      }
    }
  }
})

