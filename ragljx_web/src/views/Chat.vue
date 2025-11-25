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
                <div v-if="message.sources && message.sources.length > 0" class="message-sources">
                  <el-divider />
                  <div class="sources-title">参考来源：</div>
                  <div v-for="(source, index) in message.sources" :key="index" class="source-item">
                    <el-tag size="small">{{ source.document_name }}</el-tag>
                    <span class="source-score">相似度: {{ (source.score * 100).toFixed(1) }}%</span>
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
            <el-input v-model="inputMessage" type="textarea" :rows="3" placeholder="请输入您的问题..."
              @keydown.enter.exact.prevent="sendMessage" />
            <el-button type="primary" :loading="isLoading"
              :disabled="!inputMessage.trim() || selectedKnowledgeBases.length === 0" @click="sendMessage">
              <el-icon>
                <Promotion />
              </el-icon>
              发送
            </el-button>
          </div>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { useChatStore } from '@/stores/chat'
import { useKnowledgeStore } from '@/stores/knowledge'
import {
  getChatSessions,
  createChatSession,
  deleteChatSession,
  getChatMessages,
  sendMessageStream
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

const fetchSessions = async () => {
  try {
    const response = await getChatSessions({ page: 1, page_size: 100 })
    sessions.value = response.data?.items || []
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
      knowledge_base_ids: []
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

  isLoading.value = true

  try {
    // 使用流式输出
    const eventSource = sendMessageStream(currentSessionId.value, {
      question,
      knowledge_base_ids: selectedKnowledgeBases.value
    })

    let assistantMessage = {
      id: Date.now() + 1,
      role: 'assistant',
      content: '',
      sources: [],
      created_at: new Date().toISOString()
    }
    messages.value.push(assistantMessage)

    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data)

      if (data.type === 'content') {
        assistantMessage.content += data.content
        scrollToBottom()
      } else if (data.type === 'sources') {
        assistantMessage.sources = data.sources
      } else if (data.type === 'done') {
        eventSource.close()
        isLoading.value = false
      }
    }

    eventSource.onerror = (error) => {
      console.error('Stream error:', error)
      eventSource.close()
      isLoading.value = false
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
  const now = new Date()
  const diff = now - date

  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + ' 分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + ' 小时前'
  return date.toLocaleDateString('zh-CN')
}

onMounted(() => {
  fetchSessions()
  fetchKnowledgeBases()
})
</script>

<style scoped>
.chat-container {
  height: calc(100vh - 120px);
  background: white;
  border-radius: 8px;
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
  padding: 12px;
  margin-bottom: 8px;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: background-color 0.3s;
}

.session-item:hover {
  background-color: #f5f7fa;
}

.session-item.active {
  background-color: #ecf5ff;
}

.session-info {
  flex: 1;
  overflow: hidden;
}

.session-title {
  font-size: 14px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.session-time {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.delete-icon {
  color: #999;
  cursor: pointer;
}

.delete-icon:hover {
  color: #f56c6c;
}

.chat-main {
  padding: 0;
  display: flex;
  flex-direction: column;
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
}

.kb-selector {
  padding: 15px 20px;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  gap: 10px;
}

.message-list {
  flex: 1;
  padding: 20px;
}

.message-item {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.message-item.user {
  flex-direction: row-reverse;
}

.message-avatar {
  flex-shrink: 0;
}

.message-content {
  max-width: 70%;
  padding: 12px 16px;
  border-radius: 8px;
  background-color: #f5f7fa;
}

.message-item.user .message-content {
  background-color: #409EFF;
  color: white;
}

.message-text {
  line-height: 1.6;
  word-break: break-word;
}

.message-sources {
  margin-top: 10px;
}

.sources-title {
  font-size: 12px;
  color: #666;
  margin-bottom: 8px;
}

.source-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
  font-size: 12px;
}

.source-score {
  color: #999;
}

.typing-indicator {
  display: flex;
  gap: 4px;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #409EFF;
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
  padding: 20px;
  border-top: 1px solid #e6e6e6;
  display: flex;
  gap: 10px;
}

.input-area .el-textarea {
  flex: 1;
}
</style>
