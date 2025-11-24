-- 系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
    id SERIAL PRIMARY KEY,
    config_key VARCHAR(128) NOT NULL UNIQUE,
    config_value TEXT,
    config_type VARCHAR(32),
    description TEXT,
    is_encrypted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_config_key ON system_configs(config_key);

-- 插入默认系统配置
INSERT INTO system_configs (config_key, config_value, config_type, description) VALUES 
    ('max_file_size', '52428800', 'int', '最大文件上传大小（字节），默认50MB'),
    ('allowed_file_types', 'pdf,docx,xlsx,pptx,txt,md,html,csv,json,xml', 'string', '允许上传的文件类型'),
    ('default_top_k', '3', 'int', 'RAG检索默认返回文档数'),
    ('default_similarity_threshold', '0.7', 'float', 'RAG检索默认相似度阈值'),
    ('session_expire_hours', '24', 'int', '会话过期时间（小时）')
ON CONFLICT (config_key) DO NOTHING;

