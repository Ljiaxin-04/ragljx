import request from '@/utils/request'

/**
 * 获取文档列表
 */
export function getDocuments(kbId, params) {
  return request({
    url: `/knowledge-bases/${kbId}/documents`,
    method: 'get',
    params
  })
}

/**
 * 获取文档详情
 */
export function getDocument(kbId, docId) {
  return request({
    url: `/knowledge-bases/${kbId}/documents/${docId}`,
    method: 'get'
  })
}

/**
 * 上传文档
 */
export function uploadDocument(kbId, formData, onUploadProgress) {
  return request({
    url: `/knowledge-bases/${kbId}/documents/upload`,
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    onUploadProgress
  })
}

/**
 * 删除文档
 */
export function deleteDocument(kbId, docId) {
  return request({
    url: `/knowledge-bases/${kbId}/documents/${docId}`,
    method: 'delete'
  })
}

/**
 * 重新处理文档
 */
export function reprocessDocument(kbId, docId) {
  return request({
    url: `/knowledge-bases/${kbId}/documents/${docId}/reprocess`,
    method: 'post'
  })
}

