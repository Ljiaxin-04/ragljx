import { defineStore } from 'pinia'
import { getKnowledgeBases } from '@/api/knowledge'

export const useKnowledgeStore = defineStore('knowledge', {
  state: () => ({
    knowledgeBases: [],
    currentKnowledgeBase: null,
    loading: false
  }),

  getters: {
    knowledgeBaseList: (state) => state.knowledgeBases,
    currentKB: (state) => state.currentKnowledgeBase
  },

  actions: {
    // 获取知识库列表
    async fetchKnowledgeBases(params = {}) {
      this.loading = true
      try {
        const response = await getKnowledgeBases(params)
        this.knowledgeBases = response.data?.items || []
        return response
      } catch (error) {
        console.error('Fetch knowledge bases failed:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    // 设置当前知识库
    setCurrentKnowledgeBase(kb) {
      this.currentKnowledgeBase = kb
    },

    // 清除当前知识库
    clearCurrentKnowledgeBase() {
      this.currentKnowledgeBase = null
    },

    // 添加知识库到列表
    addKnowledgeBase(kb) {
      this.knowledgeBases.unshift(kb)
    },

    // 更新知识库
    updateKnowledgeBase(id, data) {
      const index = this.knowledgeBases.findIndex(kb => kb.id === id)
      if (index !== -1) {
        this.knowledgeBases[index] = { ...this.knowledgeBases[index], ...data }
      }
    },

    // 删除知识库
    removeKnowledgeBase(id) {
      const index = this.knowledgeBases.findIndex(kb => kb.id === id)
      if (index !== -1) {
        this.knowledgeBases.splice(index, 1)
      }
    }
  }
})

