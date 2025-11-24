-- 启用 UUID 扩展
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 知识库表
CREATE TABLE IF NOT EXISTS knowledge_bases (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    name VARCHAR(128) NOT NULL,
    english_name VARCHAR(128) NOT NULL UNIQUE,
    description TEXT,
    embedding_model VARCHAR(64) DEFAULT 'text-embedding-3-small',
    collection_name VARCHAR(255),
    document_count INTEGER DEFAULT 0,
    total_size BIGINT DEFAULT 0,
    status VARCHAR(32) DEFAULT 'active',
    is_builtin BOOLEAN DEFAULT FALSE,
    created_by_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_kb_english_name ON knowledge_bases(english_name);
CREATE INDEX idx_kb_status ON knowledge_bases(status);
CREATE INDEX idx_kb_created_by ON knowledge_bases(created_by_id);

-- 知识库文档表
CREATE TABLE IF NOT EXISTS knowledge_documents (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    knowledge_base_id VARCHAR(64) NOT NULL REFERENCES knowledge_bases(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    object_key VARCHAR(512),
    size BIGINT DEFAULT 0,
    mime VARCHAR(128),
    checksum VARCHAR(128),
    parsing_status VARCHAR(32) DEFAULT 'pending',
    error_message TEXT,
    is_builtin BOOLEAN DEFAULT FALSE,
    created_by_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_doc_kb_id ON knowledge_documents(knowledge_base_id);
CREATE INDEX idx_doc_status ON knowledge_documents(parsing_status);
CREATE INDEX idx_doc_created_by ON knowledge_documents(created_by_id);
CREATE INDEX idx_doc_checksum ON knowledge_documents(checksum);

-- 文档任务表
CREATE TABLE IF NOT EXISTS document_tasks (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    document_id VARCHAR(64) NOT NULL REFERENCES knowledge_documents(id) ON DELETE CASCADE,
    task_type VARCHAR(32) NOT NULL,
    status VARCHAR(32) DEFAULT 'pending',
    priority INTEGER DEFAULT 0,
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_task_doc_id ON document_tasks(document_id);
CREATE INDEX idx_task_status ON document_tasks(status);
CREATE INDEX idx_task_priority ON document_tasks(priority DESC);
CREATE INDEX idx_task_type ON document_tasks(task_type);

