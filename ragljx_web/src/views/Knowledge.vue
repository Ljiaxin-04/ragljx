<template>
  <div class="knowledge-container">
    <div class="page-header">
      <div>
        <div class="eyebrow">知识空间</div>
        <h2>知识库管理</h2>
        <p class="subtitle">管理你的业务文档与向量配置</p>
      </div>
      <el-button type="primary" round @click="showCreateDialog">
        <el-icon><Plus /></el-icon>
        创建知识库
      </el-button>
    </div>
    
    <div v-loading="loading" class="knowledge-list">
      <el-empty v-if="!loading && knowledgeBases.length === 0" description="暂无知识库，请创建一个" />
      
      <el-row :gutter="20">
        <el-col
          v-for="kb in knowledgeBases"
          :key="kb.id"
          :xs="24"
          :sm="12"
          :md="8"
          :lg="6"
        >
          <el-card class="kb-card" shadow="hover">
            <div class="kb-header">
              <div class="kb-icon-pill">
                <el-icon><Collection /></el-icon>
              </div>
              <div class="kb-title-block">
                <h3 class="kb-name">{{ kb.name }}</h3>
                <div class="kb-sub">{{ kb.english_name }}</div>
              </div>
              <el-tag type="success" effect="plain" size="small">运行中</el-tag>
            </div>
            
            <p class="kb-description">{{ kb.description || '暂无描述' }}</p>
            
            <div class="kb-stats">
              <div class="stat-item">
                <el-icon><Document /></el-icon>
                <span>{{ kb.document_count || 0 }} 个文档</span>
              </div>
              <div class="stat-item">
                <el-icon><Clock /></el-icon>
                <span>{{ formatDate(kb.created_at) }}</span>
              </div>
            </div>
            
            <div class="kb-actions">
              <el-button type="primary" size="small" round @click="goToDocuments(kb.id)">
                查看文档
              </el-button>
              <el-button size="small" round @click="showEditDialog(kb)">
                编辑
              </el-button>
              <el-button type="danger" size="small" round plain @click="handleDelete(kb)">
                删除
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
    
    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="520px"
      class="glass-dialog"
    >
      <div class="dialog-body">
        <el-form
          ref="formRef"
          :model="form"
          :rules="formRules"
          label-width="100px"
        >
          <el-form-item label="名称" prop="name">
            <el-input v-model="form.name" placeholder="请输入知识库名称" />
          </el-form-item>

          <el-form-item label="英文标识" prop="english_name">
            <el-input v-model="form.english_name" placeholder="请输入英文标识（如：my_kb）" />
          </el-form-item>

          <el-form-item label="描述" prop="description">
            <el-input
              v-model="form.description"
              type="textarea"
              :rows="4"
              placeholder="请输入知识库描述"
            />
          </el-form-item>
          
          <el-form-item label="嵌入模型" prop="embedding_model">
            <el-select v-model="form.embedding_model" placeholder="请选择嵌入模型" class="full-width">
              <el-option label="text-embedding-ada-002" value="text-embedding-ada-002" />
              <el-option label="text-embedding-3-small" value="text-embedding-3-small" />
              <el-option label="text-embedding-3-large" value="text-embedding-3-large" />
            </el-select>
          </el-form-item>
          
          <div class="inline-fields">
            <el-form-item label="分块大小" prop="chunk_size">
              <el-input-number v-model="form.chunk_size" :min="100" :max="2000" :step="100" />
            </el-form-item>
            <el-form-item label="分块重叠" prop="chunk_overlap">
              <el-input-number v-model="form.chunk_overlap" :min="0" :max="500" :step="50" />
            </el-form-item>
          </div>
        </el-form>
      </div>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useKnowledgeStore } from '@/stores/knowledge'
import {
  getKnowledgeBases,
  createKnowledgeBase,
  updateKnowledgeBase,
  deleteKnowledgeBase
} from '@/api/knowledge'

const router = useRouter()
const knowledgeStore = useKnowledgeStore()

const loading = ref(false)
const knowledgeBases = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('创建知识库')
const submitting = ref(false)
const formRef = ref(null)
const editingId = ref(null)

const form = reactive({
  name: '',
  english_name: '',
  description: '',
  embedding_model: 'text-embedding-3-small',
  chunk_size: 500,
  chunk_overlap: 50
})

const formRules = {
  name: [
    { required: true, message: '请输入知识库名称', trigger: 'blur' },
    { min: 2, max: 100, message: '名称长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  english_name: [
    { required: true, message: '请输入英文标识', trigger: 'blur' },
    { pattern: /^[a-z0-9_]+$/, message: '只能包含小写字母、数字和下划线', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  embedding_model: [
    { required: true, message: '请选择嵌入模型', trigger: 'change' }
  ]
}

const fetchKnowledgeBases = async () => {
  loading.value = true
  try {
    const response = await getKnowledgeBases({ page: 1, page_size: 100 })
    knowledgeBases.value = response.data?.items || []
    knowledgeStore.knowledgeBases = knowledgeBases.value
  } catch (error) {
    console.error('Fetch knowledge bases failed:', error)
    ElMessage.error('获取知识库列表失败')
  } finally {
    loading.value = false
  }
}

const showCreateDialog = () => {
  dialogTitle.value = '创建知识库'
  editingId.value = null
  resetForm()
  dialogVisible.value = true
}

const showEditDialog = (kb) => {
  dialogTitle.value = '编辑知识库'
  editingId.value = kb.id
  form.name = kb.name
  form.english_name = kb.english_name
  form.description = kb.description
  form.embedding_model = kb.embedding_model
  form.chunk_size = kb.chunk_size
  form.chunk_overlap = kb.chunk_overlap
  dialogVisible.value = true
}

const resetForm = () => {
  form.name = ''
  form.english_name = ''
  form.description = ''
  form.embedding_model = 'text-embedding-3-small'
  form.chunk_size = 500
  form.chunk_overlap = 50
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (editingId.value) {
          await updateKnowledgeBase(editingId.value, form)
          ElMessage.success('更新成功')
        } else {
          await createKnowledgeBase(form)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        fetchKnowledgeBases()
      } catch (error) {
        console.error('Submit failed:', error)
        ElMessage.error(error.response?.data?.message || '操作失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleDelete = (kb) => {
  ElMessageBox.confirm(`确定要删除知识库"${kb.name}"吗？`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteKnowledgeBase(kb.id)
      ElMessage.success('删除成功')
      fetchKnowledgeBases()
    } catch (error) {
      console.error('Delete failed:', error)
      ElMessage.error('删除失败')
    }
  }).catch(() => {
    // 取消操作
  })
}

const goToDocuments = (kbId) => {
  router.push(`/knowledge/${kbId}/documents`)
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

onMounted(() => {
  fetchKnowledgeBases()
})
</script>

<style scoped>
.knowledge-container {
  height: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 26px;
  color: var(--ui-text);
}

.subtitle {
  margin: 6px 0 0 0;
  color: var(--ui-subtext);
  font-size: 13px;
}

.eyebrow {
  font-size: 12px;
  color: var(--ui-subtext);
  letter-spacing: 0.6px;
}

.knowledge-list {
  min-height: 400px;
}

.kb-card {
  margin-bottom: 20px;
  transition: transform 0.3s, box-shadow 0.3s;
  border-radius: 16px;
  border: 1px solid var(--ui-border);
  box-shadow: var(--ui-shadow);
  background: linear-gradient(180deg, rgba(43, 104, 255, 0.03), rgba(17, 207, 161, 0.02) 60%, #fff 100%);
}

.kb-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 14px 34px rgba(19, 54, 109, 0.14);
}

.kb-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 15px;
}

.kb-icon-pill {
  width: 44px;
  height: 44px;
  border-radius: 14px;
  background: rgba(43, 104, 255, 0.1);
  display: grid;
  place-items: center;
  color: var(--ui-primary);
  font-size: 22px;
}

.kb-name {
  margin: 0;
  font-size: 18px;
  color: var(--ui-text);
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.kb-title-block {
  flex: 1;
  min-width: 0;
}

.kb-sub {
  font-size: 12px;
  color: var(--ui-subtext);
  margin-top: 4px;
}

.kb-description {
  color: #4a5568;
  font-size: 14px;
  margin: 0 0 15px 0;
  min-height: 40px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.kb-stats {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
  padding-top: 15px;
  border-top: 1px solid var(--ui-border);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 13px;
  color: var(--ui-subtext);
}

.kb-actions {
  display: flex;
  gap: 10px;
}

.kb-actions .el-button {
  flex: 1;
}

.glass-dialog :deep(.el-dialog__body) {
  background: rgba(255, 255, 255, 0.85);
  border-radius: 12px;
  backdrop-filter: blur(12px);
}

.glass-dialog :deep(.el-dialog__header) {
  border-bottom: none;
}

.dialog-body {
  background: linear-gradient(135deg, rgba(43, 104, 255, 0.05), rgba(17, 207, 161, 0.05));
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  padding: 12px;
}

.full-width {
  width: 100%;
}

.inline-fields {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
}
</style>
