import request from '@/utils/request'

/**
 * 获取对话会话列表
 */
export function getChatSessions(params) {
  return request({
    url: '/chat/sessions',
    method: 'get',
    params
  })
}

/**
 * 创建对话会话
 */
export function createChatSession(data) {
  return request({
    url: '/chat/sessions',
    method: 'post',
    data
  })
}

/**
 * 获取会话详情
 */
export function getChatSession(sessionId) {
  return request({
    url: `/chat/sessions/${sessionId}`,
    method: 'get'
  })
}

/**
 * 删除会话
 */
export function deleteChatSession(sessionId) {
  return request({
    url: `/chat/sessions/${sessionId}`,
    method: 'delete'
  })
}

/**
 * 获取会话消息列表
 */
export function getChatMessages(sessionId, params) {
  return request({
    url: `/chat/sessions/${sessionId}/messages`,
    method: 'get',
    params
  })
}

/**
 * 发送消息（非流式）
 */
export function sendMessage(sessionId, data) {
  return request({
    url: `/chat/sessions/${sessionId}/messages`,
    method: 'post',
    data
  })
}

/**
 * 发送消息（流式）- 使用 EventSource
 */
export function sendMessageStream(sessionId, data) {
  const token = localStorage.getItem('access_token')
  const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
  
  // 构建查询参数
  const params = new URLSearchParams({
    question: data.question,
    knowledge_base_ids: JSON.stringify(data.knowledge_base_ids || []),
    stream: 'true'
  })
  
  const url = `${baseURL}/chat/sessions/${sessionId}/messages/stream?${params.toString()}`
  
  return new EventSource(url, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  })
}

