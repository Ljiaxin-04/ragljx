<template>
  <div class="documents-container">
    <div class="page-header">
      <div class="header-left">
        <el-button text @click="goBack">
          <el-icon>
            <ArrowLeft />
          </el-icon>
          返回
        </el-button>
        <h2>{{ knowledgeBaseName }}</h2>
      </div>
      <el-upload ref="uploadRef" :action="uploadUrl" :headers="uploadHeaders" :on-success="handleUploadSuccess"
        :on-error="handleUploadError" :before-upload="beforeUpload" :show-file-list="false" multiple>
        <el-button type="primary">
          <el-icon>
            <Upload />
          </el-icon>
          上传文档
        </el-button>
      </el-upload>
    </div>

    <el-card>
      <el-table v-loading="loading" :data="documents" style="width: 100%">
        <el-table-column prop="name" label="文档名称" min-width="200">
          <template #default="{ row }">
            <div class="doc-name">
              <el-icon>
                <Document />
              </el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="file_type" label="文件类型" width="140">
          <template #default="{ row }">
            <el-tag size="small">{{ formatFileType(row.file_type, row.name) }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="file_size" label="文件大小" width="120">
          <template #default="{ row }">
            {{ formatFileSize(row.file_size) }}
          </template>
        </el-table-column>

        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <div class="status-cell">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
              <el-progress
                v-if="row.status === 'processing' || row.status === 'pending'"
                :indeterminate="true"
                :stroke-width="4"
                :show-text="false"
                class="inline-progress"
              />
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="chunk_count" label="分块数" width="100" />

        <el-table-column prop="created_at" label="上传时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.status === 'failed'" type="primary" size="small" @click="handleReprocess(row)">
              重新处理
            </el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>

        <template #empty>
          <el-empty description="暂无文档，请上传文档" />
        </template>
      </el-table>

      <div class="pagination">
        <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :total="total"
          :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper" @size-change="fetchDocuments"
          @current-change="fetchDocuments" />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getDocuments,
  deleteDocument,
  reprocessDocument
} from '@/api/document'
import { getKnowledgeBase } from '@/api/knowledge'

const route = useRoute()
const router = useRouter()

const kbId = computed(() => route.params.id)
const knowledgeBaseName = ref('')
const loading = ref(false)
const documents = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const pollTimer = ref(null)

const uploadUrl = computed(() => {
  const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
  return `${baseURL}/knowledge-bases/${kbId.value}/documents/upload`
})

const uploadHeaders = computed(() => {
  const token = localStorage.getItem('access_token')
  return {
    Authorization: `Bearer ${token}`
  }
})

const fetchKnowledgeBase = async () => {
  try {
    const response = await getKnowledgeBase(kbId.value)
    knowledgeBaseName.value = response.data.name
  } catch (error) {
    console.error('Fetch knowledge base failed:', error)
  }
}

const fetchDocuments = async () => {
  loading.value = true
  try {
    const response = await getDocuments(kbId.value, {
      page: currentPage.value,
      page_size: pageSize.value
    })
    documents.value = response.data?.items || []
    total.value = response.data?.total || 0
  } catch (error) {
    console.error('Fetch documents failed:', error)
    ElMessage.error('获取文档列表失败')
  } finally {
    loading.value = false
    // 如果还有处理中/待处理的文档，保持轮询，否则清理
    const hasProcessing = documents.value.some(d => d.status === 'processing' || d.status === 'pending')
    if (!hasProcessing && pollTimer.value) {
      clearInterval(pollTimer.value)
      pollTimer.value = null
    }
  }
}

const beforeUpload = (file) => {
  const allowedTypes = [
    'text/plain',
    'text/markdown',
    'application/pdf',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'application/vnd.openxmlformats-officedocument.presentationml.presentation',
    'text/html',
    'text/csv',
    'application/json',
    'application/xml',
    'application/rtf'
  ]

  const isAllowedType = allowedTypes.includes(file.type) ||
    file.name.endsWith('.txt') ||
    file.name.endsWith('.md') ||
    file.name.endsWith('.pdf') ||
    file.name.endsWith('.docx') ||
    file.name.endsWith('.xlsx') ||
    file.name.endsWith('.pptx') ||
    file.name.endsWith('.html') ||
    file.name.endsWith('.csv') ||
    file.name.endsWith('.json') ||
    file.name.endsWith('.xml') ||
    file.name.endsWith('.rtf')

  if (!isAllowedType) {
    ElMessage.error('不支持的文件类型')
    return false
  }

  const isLt50M = file.size / 1024 / 1024 < 50
  if (!isLt50M) {
    ElMessage.error('文件大小不能超过 50MB')
    return false
  }

  return true
}

const handleUploadSuccess = () => {
  ElMessage.success('上传成功')
  fetchDocuments()
  startPolling()
}

const handleUploadError = (error) => {
  console.error('Upload failed:', error)
  ElMessage.error('上传失败')
}

const handleReprocess = async (doc) => {
  try {
    await reprocessDocument(kbId.value, doc.id)
    ElMessage.success('已提交重新处理请求')
    fetchDocuments()
    startPolling()
  } catch (error) {
    console.error('Reprocess failed:', error)
    ElMessage.error('重新处理失败')
  }
}

const handleDelete = (doc) => {
  ElMessageBox.confirm(`确定要删除文档"${doc.name}"吗？`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteDocument(kbId.value, doc.id)
      ElMessage.success('删除成功')
      fetchDocuments()
    } catch (error) {
      console.error('Delete failed:', error)
      ElMessage.error('删除失败')
    }
  }).catch(() => {
    // 取消操作
  })
}

const goBack = () => {
  router.push('/knowledge')
}

const formatFileSize = (bytes) => {
  if (!bytes) return '-'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  return (bytes / 1024 / 1024).toFixed(2) + ' MB'
}

const formatDateTime = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

const formatFileType = (mime, name) => {
  if (!mime && name) {
    const ext = name.split('.').pop() || ''
    return ext.toLowerCase()
  }
  // mime like application/pdf -> pdf
  if (mime && mime.includes('/')) {
    const part = mime.split('/').pop() || mime
    return part.toLowerCase()
  }
  return mime || '-'
}

const getStatusType = (status) => {
  const typeMap = {
    pending: 'info',
    processing: 'warning',
    completed: 'success',
    failed: 'danger'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status) => {
  const textMap = {
    pending: '待处理',
    processing: '处理中',
    completed: '已完成',
    failed: '失败'
  }
  return textMap[status] || status
}

const startPolling = () => {
  if (pollTimer.value) return
  pollTimer.value = setInterval(() => {
    fetchDocuments()
  }, 500)
}

onMounted(() => {
  fetchKnowledgeBase()
  fetchDocuments()
})

onUnmounted(() => {
  if (pollTimer.value) {
    clearInterval(pollTimer.value)
    pollTimer.value = null
  }
})
</script>

<style scoped>
.documents-container {
  height: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-left h2 {
  margin: 0;
  font-size: 24px;
  color: #333;
}

.doc-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.inline-progress {
  width: 68px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
