<template>
  <div class="chat-container">
    <el-container>
      <!-- 左侧会话列表 -->
      <el-aside width="280px" class="session-sidebar">
        <div class="sidebar-header">
          <h3>对话历史</h3>
          <el-button type="primary" size="small" @click="createNewSession">
            <el-icon>
              <Plus />
            </el-icon>
            新对话
          </el-button>
        </div>

        <el-scrollbar class="session-list">
          <div v-for="session in sessions" :key="session.id"
            :class="['session-item', { active: currentSessionId === session.id }]" @click="selectSession(session)">
            <div class="session-info">
              <div class="session-title">{{ session.title || '新对话' }}</div>
              <div class="session-time">{{ formatDateTime(session.created_at) }}</div>
            </div>
            <el-icon class="delete-icon" @click.stop="deleteSession(session)">
              <Delete />
            </el-icon>
          </div>

          <el-empty v-if="sessions.length === 0" description="暂无对话记录" />
        </el-scrollbar>
      </el-aside>

      <!-- 右侧对话区域 -->
      <el-main class="chat-main">
        <div v-if="!currentSessionId" class="empty-chat">
          <el-empty description="请选择或创建一个对话" />
        </div>

        <div v-else class="chat-content">
          <!-- 顶部标题和会话信息 -->
          <div class="chat-header">
            <div class="chat-header-left">
              <div class="chat-title">
                {{ currentSessionTitle }}
              </div>
              <div class="chat-subtitle">
                {{ selectedKnowledgeBases.length > 0 ? '已启用知识库问答（RAG）' : '纯模型对话' }}
              </div>
            </div>
          </div>

          <!-- 知识库选择 -->
          <div class="kb-selector">
            <span>选择知识库：</span>
            <el-select v-model="selectedKnowledgeBases" multiple placeholder="请选择知识库" style="width: 400px">
              <el-option v-for="kb in knowledgeBases" :key="kb.id" :label="kb.name" :value="kb.id" />
            </el-select>
          </div>

          <!-- 消息列表 -->
          <el-scrollbar ref="scrollbarRef" class="message-list">
            <div v-for="message in messages" :key="message.id" :class="['message-item', message.role]">
              <div class="message-avatar">
                <el-avatar v-if="message.role === 'user'" :size="36">
                  {{ userStore.nickname.charAt(0) }}
                </el-avatar>
                <el-icon v-else :size="36" color="#409EFF">
                  <ChatDotRound />
                </el-icon>
              </div>

              <div class="message-content">
                <div class="message-text" v-html="formatMessage(message.content)"></div>
                <div v-if="getSources(message).length > 0" class="message-sources">
                  <el-divider />
                  <div class="sources-title">📚 参考来源：</div>
                  <div v-for="(source, index) in getSources(message)" :key="index" class="source-item">
                    <div class="source-info">
                      <el-icon class="source-icon">
                        <Document />
                      </el-icon>
                      <div class="source-text">
                        <div class="source-name">
                          {{ source.file_name || source.document_name || source.title || '知识库文档' }}
                        </div>
                        <div class="source-score">
                          相似度: {{ (Number(source.score || 0) * 100).toFixed(1) }}%
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 加载中 -->
            <div v-if="isLoading" class="message-item assistant">
              <div class="message-avatar">
                <el-icon :size="36" color="#409EFF">
                  <ChatDotRound />
                </el-icon>
              </div>
              <div class="message-content">
                <div class="typing-indicator">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
            </div>
          </el-scrollbar>

          <!-- 输入框 -->
          <div class="input-area">
            <div class="input-wrapper">
              <el-input v-model="inputMessage" type="textarea" :rows="3"
                :placeholder="selectedKnowledgeBases.length === 0 ? '请先选择知识库...' : '输入您的问题，按 Enter 发送，Shift + Enter 换行...'"
                @keydown.enter.exact.prevent="sendMessage" :disabled="selectedKnowledgeBases.length === 0"
                class="message-input" />
              <div class="input-actions">
                <div class="input-hint">
                  <el-icon>
                    <InfoFilled />
                  </el-icon>
                  <span v-if="selectedKnowledgeBases.length === 0">请先选择知识库</span>
                  <span v-else>按 Enter 发送，Shift + Enter 换行</span>
                </div>
                <el-button type="primary" :loading="isLoading"
                  :disabled="!inputMessage.trim() || selectedKnowledgeBases.length === 0" @click="sendMessage"
                  size="large">
                  <el-icon>
                    <Promotion />
                  </el-icon>
                  {{ isLoading ? '发送中...' : '发送' }}
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, watch, computed, reactive } from 'vue'
import { onBeforeUnmount } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { useChatStore } from '@/stores/chat'
import { useKnowledgeStore } from '@/stores/knowledge'
import {
  getChatSessions,
  createChatSession,
  deleteChatSession,
  getChatMessages,
  sendMessageStream,
  updateChatSession
} from '@/api/chat'

const userStore = useUserStore()
const chatStore = useChatStore()
const knowledgeStore = useKnowledgeStore()

const sessions = ref([])
const currentSessionId = ref(null)
const messages = ref([])
const knowledgeBases = ref([])
const selectedKnowledgeBases = ref([])
const inputMessage = ref('')
const isLoading = ref(false)
const scrollbarRef = ref(null)
const eventSourceRef = ref(null)


const currentSession = computed(() =>
  sessions.value.find((s) => s.id === currentSessionId.value) || null
)

const currentSessionTitle = computed(() =>
  (currentSession.value && currentSession.value.title)
    ? currentSession.value.title
    : '新对话'
)

const fetchSessions = async () => {
  try {
    const response = await getChatSessions({ page: 1, page_size: 100 })
    const items = response.data?.items || []
    // 按创建时间倒序展示，最近的在上
    sessions.value = [...items].sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
  } catch (error) {
    console.error('Fetch sessions failed:', error)
  }
}

const fetchKnowledgeBases = async () => {
  try {
    await knowledgeStore.fetchKnowledgeBases({ page: 1, page_size: 100 })
    knowledgeBases.value = knowledgeStore.knowledgeBases
  } catch (error) {
    console.error('Fetch knowledge bases failed:', error)
  }
}

const getSources = (message) => {
  if (!message) return []

  const rawSources = message.sources || message.rag_sources || []

  const mapped = rawSources.map((s) => ({
    document_id: s.document_id || s.DocumentID || s.id || '',
    document_name: s.document_name || s.DocumentName || '',
    title: s.title || s.Title || '',
    file_name: s.file_name || s.FileName || '',
    score: typeof s.score === 'number'
      ? s.score
      : (typeof s.Score === 'number' ? s.Score : 0)
  }))

  // 按相似度排序，展示前3个来源
  if (!mapped.length) return mapped
  mapped.sort((a, b) => (b.score || 0) - (a.score || 0))
  return mapped.slice(0, 3)
}


const fetchMessages = async (sessionId) => {
  try {
    const response = await getChatMessages(sessionId, { page: 1, page_size: 100 })
    messages.value = response.data?.items || []
    scrollToBottom()
  } catch (error) {
    console.error('Fetch messages failed:', error)
  }
}

const createNewSession = async () => {
  try {
    const response = await createChatSession({
      title: '新对话',
      knowledge_base_ids: selectedKnowledgeBases.value,  // 使用当前选择的知识库
      use_rag: selectedKnowledgeBases.value.length > 0,  // 如果有知识库则启用 RAG
      top_k: 5,
      similarity_threshold: 0.7,
      similarity_weight: 1.5
    })
    const newSession = response.data
    sessions.value.unshift(newSession)
    selectSession(newSession)
    ElMessage.success('创建成功')
  } catch (error) {
    console.error('Create session failed:', error)
    ElMessage.error('创建失败')
  }
}

const selectSession = (session) => {
  currentSessionId.value = session.id
  // 同步当前会话绑定的知识库
  selectedKnowledgeBases.value = session.knowledge_base_ids || []
  chatStore.setCurrentSession(session)
  fetchMessages(session.id)
}

const deleteSession = (session) => {
  ElMessageBox.confirm('确定要删除这个对话吗？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteChatSession(session.id)
      sessions.value = sessions.value.filter(s => s.id !== session.id)
      if (currentSessionId.value === session.id) {
        currentSessionId.value = null
        messages.value = []
      }
      ElMessage.success('删除成功')
    } catch (error) {
      console.error('Delete session failed:', error)
      ElMessage.error('删除失败')
    }
  }).catch(() => {
    // 取消操作
  })
}

// 监听知识库选择变化，实时更新后端会话配置
watch(selectedKnowledgeBases, async (newVal) => {
  if (!currentSessionId.value) return

  try {
    await updateChatSession(currentSessionId.value, {
      knowledge_base_ids: newVal,
      use_rag: newVal.length > 0
    })

    const index = sessions.value.findIndex((s) => s.id === currentSessionId.value)
    if (index !== -1) {
      sessions.value[index] = {
        ...sessions.value[index],
        knowledge_base_ids: [...newVal],
        use_rag: newVal.length > 0
      }
    }
  } catch (error) {
    console.error('Update session knowledge bases failed:', error)
    ElMessage.error('更新知识库选择失败')
  }
})

const sendMessage = async () => {
  if (!inputMessage.value.trim() || selectedKnowledgeBases.value.length === 0) {
    return
  }

  if (!currentSessionId.value) {
    ElMessage.warning('请先创建或选择一个对话')
    return
  }

  const question = inputMessage.value.trim()
  inputMessage.value = ''

  // 添加用户消息
  const userMessage = {
    id: Date.now(),
    role: 'user',
    content: question,
    created_at: new Date().toISOString()
  }
  messages.value.push(userMessage)
  scrollToBottom()

  // 如果会话标题还是默认“新对话”，自动更新为首条提问的前30字符
  if (currentSession.value && (!currentSession.value.title || currentSession.value.title.startsWith('新对话'))) {
    const newTitle = question.slice(0, 30)
    updateChatSession(currentSession.value.id, { title: newTitle }).then(() => {
      // 同步前端列表与 store
      const idx = sessions.value.findIndex((s) => s.id === currentSession.value.id)
      if (idx !== -1) sessions.value[idx] = { ...sessions.value[idx], title: newTitle }
      chatStore.setCurrentSession({ ...currentSession.value, title: newTitle })
    }).catch((err) => {
      console.error('Update session title failed:', err)
    })
  }

  isLoading.value = true

  try {
    // 关闭上一个流，防止多个连接阻塞 UI
    if (eventSourceRef.value) {
      eventSourceRef.value.close()
      eventSourceRef.value = null
    }

    // 使用流式输出（不需要传递 knowledge_base_ids，会话已经包含了）
    const eventSource = sendMessageStream(currentSessionId.value, {
      question
    })
    eventSourceRef.value = eventSource

    let assistantMessage = null
    let buffer = ''
    let flushTimer = null
    let pendingSources = []

    const ensureAssistantMessage = () => {
      if (!assistantMessage) {
        // 使用 reactive，确保后续属性修改是可响应的
        assistantMessage = reactive({
          id: Date.now() + 1,
          role: 'assistant',
          content: '',
          sources: [],
          created_at: new Date().toISOString()
        })
        messages.value.push(assistantMessage)
      }
    }

    // 减少 DOM 抖动：缓冲内容，每 40ms 刷新一次
    const flushBuffer = () => {
      if (!assistantMessage || !buffer) return
      assistantMessage.content += buffer
      buffer = ''
      scrollToBottom()
    }

    const scheduleFlush = () => {
      if (flushTimer) return
      flushTimer = setTimeout(() => {
        flushTimer = null
        flushBuffer()
      }, 40)
    }

    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data)

      if (data.type === 'start') {
        ensureAssistantMessage()
        isLoading.value = false
      } else if (data.type === 'content') {
        ensureAssistantMessage()
        // 若后台未发送 start 事件，第一次内容也能关闭 loading
        isLoading.value = false
        buffer += data.content || ''
        scheduleFlush()
      } else if (data.type === 'sources') {
        pendingSources = data.sources || []
        isLoading.value = false
      } else if (data.type === 'error') {
        eventSource.close()
        eventSourceRef.value = null
        isLoading.value = false
        flushBuffer()
        ElMessage.error('对话失败: ' + (data.error || '未知错误'))
      } else if (data.type === 'done') {
        eventSource.close()
        eventSourceRef.value = null
        flushBuffer()
        ensureAssistantMessage()
        assistantMessage.sources = pendingSources
        isLoading.value = false
      }
    }

    eventSource.onerror = (error) => {
      console.error('Stream error:', error)
      eventSource.close()
      eventSourceRef.value = null
      isLoading.value = false
      flushBuffer()
      ElMessage.error('发送失败')
    }
  } catch (error) {
    console.error('Send message failed:', error)
    isLoading.value = false
    ElMessage.error('发送失败')
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (scrollbarRef.value) {
      const scrollElement = scrollbarRef.value.$el.querySelector('.el-scrollbar__wrap')
      if (scrollElement) {
        scrollElement.scrollTop = scrollElement.scrollHeight
      }
    }
  })
}

const formatMessage = (content) => {
  if (!content) return ''
  // 简单的 Markdown 转换
  return content
    .replace(/\n/g, '<br>')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
    .replace(/`(.*?)`/g, '<code>$1</code>')
}

const formatDateTime = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  if (isNaN(date.getTime())) return ''
  // 统一显示为 UTC+8（北京时间）
  return date.toLocaleString('zh-CN', {
    hour12: false,
    timeZone: 'Asia/Shanghai'
  })
}

onMounted(() => {
  fetchSessions()
  fetchKnowledgeBases()
})

onBeforeUnmount(() => {
  if (eventSourceRef.value) {
    eventSourceRef.value.close()
    eventSourceRef.value = null
  }
})
</script>

<style scoped>
.chat-container {
  height: calc(100vh - 120px);
  background: white;
  border-radius: 8px;
  overflow: hidden;
  /* 防止内容超出卡片区域 */
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.chat-container :deep(.el-container) {
  height: 100%;
  /* 确保容器占满父元素高度 */
}

.session-sidebar {
  border-right: 1px solid #e6e6e6;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sidebar-header h3 {
  margin: 0;
  font-size: 16px;
}

.session-list {
  flex: 1;
  padding: 10px;
}

.session-item {
  padding: 14px;
  margin-bottom: 8px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: all 0.3s;
  border: 1px solid transparent;
}

.session-item:hover {
  background-color: #f5f7fa;
  border-color: #e6e6e6;
  transform: translateX(4px);
}

.session-item.active {
  background-color: #ecf5ff;
  border-color: #409EFF;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2);
}

.session-info {
  flex: 1;
  overflow: hidden;
}

.session-title {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 4px;
}

.session-item.active .session-title {
  color: #409EFF;
}

.session-time {
  font-size: 12px;
  color: #999;
}

.delete-icon {
  color: #999;
  cursor: pointer;
  transition: all 0.3s;
  padding: 4px;
  border-radius: 4px;
}

.delete-icon:hover {
  color: #f56c6c;
  background-color: rgba(245, 108, 108, 0.1);
}

.chat-main {
  padding: 0;
  display: flex;
  flex-direction: column;
  height: 100%;
  /* 确保主区域占满高度 */
  overflow: hidden;
  /* 防止主区域本身滚动 */
}

.empty-chat {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chat-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  /* 确保内容不会超出容器 */
}

/* 顶部标题栏样式 */
.chat-header {
  padding: 16px 20px 8px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  flex-shrink: 0;
  /* 防止标题栏被压缩 */
}

.chat-header-left {
  display: flex;
  flex-direction: column;
}

.chat-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.chat-subtitle {
  margin-top: 4px;
  font-size: 12px;
  color: #909399;
}

/* 优化消息区域的留白和行距 - 已移至下方统一定义 */

.kb-selector {
  padding: 16px 20px;
  border-bottom: 1px solid #e6e6e6;
  background-color: #fafafa;
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
  /* 防止知识库选择器被压缩 */
}

.kb-selector span {
  font-size: 14px;
  font-weight: 500;
  color: #606266;
  white-space: nowrap;
}

.message-list {
  flex: 1;
  overflow-y: auto;
  /* 允许消息列表滚动 */
  min-height: 0;
  /* 确保 flex 子元素可以正确收缩 */
  padding: 20px;
}

.message-item {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message-item.user {
  flex-direction: row-reverse;
}

.message-avatar {
  flex-shrink: 0;
}

.message-content {
  max-width: 75%;
  max-height: 480px;
  overflow-y: auto;
  padding: 16px 18px;
  border-radius: 16px;
  background-color: #f6f8fc;
  box-shadow: 0 10px 24px rgba(23, 57, 107, 0.08);
  transition: all 0.25s ease;
}

.message-content:hover {
  box-shadow: 0 12px 30px rgba(23, 57, 107, 0.12);
}

.message-item.user .message-content {
  background: linear-gradient(135deg, #2b68ff 0%, #5f8bff 100%);
  color: white;
}

.message-text {
  line-height: 1.8;
  word-break: break-word;
  font-size: 15px;
  color: var(--ui-text);
}

.message-item.user .message-text {
  color: white;
}

.card {
  background: white;
  border: 1px solid var(--ui-border);
  border-radius: 14px;
  box-shadow: var(--ui-shadow);
}

.message-sources {
  margin-top: 14px;
  padding: 12px;
}

.sources-title {
  font-size: 13px;
  font-weight: 700;
  color: var(--ui-text);
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.source-grid {
  display: grid;
  gap: 10px;
}

.source-item {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 10px;
  align-items: center;
  padding: 10px 12px;
  border-radius: 10px;
  background: #f7f9ff;
  border: 1px solid #e4e9f5;
  transition: all 0.2s;
}

.source-item:hover {
  transform: translateY(-1px);
  box-shadow: 0 8px 18px rgba(24, 72, 160, 0.08);
}

.source-icon-wrap {
  width: 32px;
  height: 32px;
  border-radius: 10px;
  background: rgba(43, 104, 255, 0.08);
  display: grid;
  place-items: center;
  color: var(--ui-primary);
}

.source-text {
  flex: 1;
  min-width: 0;
}

.source-name {
  font-size: 13px;
  color: var(--ui-text);
  font-weight: 600;
  margin-bottom: 4px;
  word-break: break-all;
}

.source-score {
  font-size: 12px;
  color: var(--ui-accent);
  font-weight: 600;
}

.typing-indicator {
  display: flex;
  gap: 4px;
  margin-top: 6px;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #2b68ff;
  animation: typing 1.4s infinite;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {

  0%,
  60%,
  100% {
    transform: translateY(0);
  }

  30% {
    transform: translateY(-10px);
  }
}

.input-area {
  padding: 16px 18px;
  border-top: 1px solid var(--ui-border);
  background-color: transparent;
  flex-shrink: 0;
}

.input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.message-input {
  width: 100%;
}

.message-input :deep(.el-textarea__inner) {
  border-radius: 12px;
  border: 2px solid #e7ecf5;
  transition: all 0.3s;
  font-size: 14px;
  line-height: 1.6;
}

.message-input :deep(.el-textarea__inner):focus {
  border-color: var(--ui-primary);
  box-shadow: 0 0 0 3px rgba(43, 104, 255, 0.12);
}

.message-input :deep(.el-textarea__inner):disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

.input-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.input-hint {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #909399;
}

.input-hint .el-icon {
  font-size: 14px;
}
</style>
