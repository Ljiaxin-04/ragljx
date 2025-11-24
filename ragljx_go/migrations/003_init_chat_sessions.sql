-- 对话会话表
CREATE TABLE IF NOT EXISTS chat_sessions (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255),
    knowledge_base_ids TEXT[],
    use_rag BOOLEAN DEFAULT TRUE,
    top_k INTEGER DEFAULT 3,
    similarity_threshold FLOAT DEFAULT 0.7,
    similarity_weight FLOAT DEFAULT 1.5,
    status VARCHAR(32) DEFAULT 'active',
    message_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_session_user_id ON chat_sessions(user_id);
CREATE INDEX idx_session_status ON chat_sessions(status);
CREATE INDEX idx_session_created_at ON chat_sessions(created_at DESC);

-- 对话消息表
CREATE TABLE IF NOT EXISTS chat_messages (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    session_id VARCHAR(64) NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE,
    role VARCHAR(32) NOT NULL,
    content TEXT NOT NULL,
    rag_sources JSONB,
    tokens_used INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_msg_session_id ON chat_messages(session_id);
CREATE INDEX idx_msg_created_at ON chat_messages(created_at);
CREATE INDEX idx_msg_role ON chat_messages(role);

