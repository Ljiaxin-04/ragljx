-- 添加 chunk_count 字段到 knowledge_documents 表
ALTER TABLE knowledge_documents ADD COLUMN IF NOT EXISTS chunk_count INTEGER DEFAULT 0;

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_doc_chunk_count ON knowledge_documents(chunk_count);

