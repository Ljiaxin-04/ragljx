import request from '@/utils/request'

/**
 * 获取知识库列表
 */
export function getKnowledgeBases(params) {
  return request({
    url: '/knowledge-bases',
    method: 'get',
    params
  })
}

/**
 * 获取知识库详情
 */
export function getKnowledgeBase(id) {
  return request({
    url: `/knowledge-bases/${id}`,
    method: 'get'
  })
}

/**
 * 创建知识库
 */
export function createKnowledgeBase(data) {
  return request({
    url: '/knowledge-bases',
    method: 'post',
    data
  })
}

/**
 * 更新知识库
 */
export function updateKnowledgeBase(id, data) {
  return request({
    url: `/knowledge-bases/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除知识库
 */
export function deleteKnowledgeBase(id) {
  return request({
    url: `/knowledge-bases/${id}`,
    method: 'delete'
  })
}

