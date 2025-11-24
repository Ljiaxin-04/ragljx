# RAG 知识库管理系统 - 详细设计文档

## 一、项目概述

### 1.1 项目背景
基于现有的gin_ljx(ioc)项目的技术积累，构建一个企业级的 RAG（检索增强生成）知识库管理系统。

### 1.2 技术栈选型

#### 后端技术栈
- **Go 服务（管理端）**
  - 框架：Gin + 自研 IOC 容器
  - 数据库：PostgreSQL（主库）
  - 缓存：Redis
  - 消息队列：Kafka
  - 对象存储：MinIO
  - RPC 通信：gRPC（与 Python 服务通信）
  - 认证：JWT Token
  
- **Python 服务（AI 处理端）**
  - 框架：Flask/FastAPI
  - 向量数据库：Qdrant
  - RAG 框架：LlamaIndex
  - 文档解析：PyMuPDF, python-docx, openpyxl 等
  - 向量化模型：OpenAI-compatible API
  - 对话模型：OpenAI-compatible API

#### 前端技术栈
- Vue 3 + TypeScript
- Element Plus / Ant Design Vue
- Axios
- Pinia（状态管理）

#### 部署方案
- Docker + Docker Compose
- 支持 Kubernetes 部署

---

## 二、系统架构设计

### 2.1 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                         前端层 (Vue3)                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  知识库管理  │  │  文档管理    │  │  对话管理    │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                              │ HTTP/HTTPS + JWT
┌─────────────────────────────▼─────────────────────────────────────┐
│                    Go 后端服务 (Gin + IOC)                        │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │  API 层 (Gin Router + JWT Middleware)                       │ │
│  └──────────────────────────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │  Service 层                                                  │ │
│  │  - 用户服务  - 知识库服务  - 文档服务  - 对话服务          │ │
│  └──────────────────────────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │  Repository 层 (GORM)                                        │ │
│  └──────────────────────────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │  IOC 容器 (配置管理)                                         │ │
│  │  - PostgreSQL  - Redis  - Kafka  - MinIO  - gRPC Client    │ │
│  └──────────────────────────────────────────────────────────────┘ │
└───────────────────────────┬───────────────────────────────────────┘
                            │ gRPC
┌───────────────────────────▼───────────────────────────────────────┐
│                  Python AI 服务 (Flask/FastAPI)                   │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │  gRPC Server                                                 │ │
│  └──────────────────────────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │  RAG 服务层                                                  │ │
│  │  - FileService (文档解析)                                   │ │
│  │  - VectorService (向量化、检索)                             │ │
│  │  - ChatService (对话生成)                                   │ │
│  └──────────────────────────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │  外部依赖                                                    │ │
│  │  - Qdrant (向量存储)  - LLM API  - Embedding API           │ │
│  └──────────────────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────────────────┘

┌───────────────────────────────────────────────────────────────────┐
│                        基础设施层                                  │
│  PostgreSQL  │  Redis  │  Kafka  │  MinIO  │  Qdrant             │
└───────────────────────────────────────────────────────────────────┘
```

### 2.2 服务职责划分

#### Go 服务职责
1. **用户认证与授权**：JWT Token 生成、验证、刷新
2. **知识库管理**：CRUD 操作、权限控制
3. **文档元数据管理**：文档上传、元数据存储、状态管理
4. **对话会话管理**：会话创建、历史记录存储
5. **任务调度**：通过 Kafka 异步调度文档处理任务
6. **API 网关**：统一对外接口、限流、日志

#### Python 服务职责
1. **文档解析**：支持 PDF、Word、Excel、Markdown 等格式
2. **文本向量化**：调用 Embedding API 生成向量
3. **向量存储与检索**：Qdrant 向量数据库操作
4. **RAG 对话**：基于检索结果生成回答
5. **gRPC 服务**：接收 Go 服务的调用请求

---

## 三、数据库设计

### 3.1 PostgreSQL 表结构设计

#### 3.1.1 用户相关表

**表名：users**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(128),
    phone VARCHAR(20),
    real_name VARCHAR(64),
    avatar_url VARCHAR(512),
    status VARCHAR(20) DEFAULT 'active', -- active, disabled
    is_admin BOOLEAN DEFAULT FALSE,
    last_login_at TIMESTAMP,
    password_updated_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
```

**表名：roles**
```sql
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL UNIQUE,
    description TEXT,
    is_builtin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**表名：user_roles**
```sql
CREATE TABLE user_roles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, role_id)
);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
```

**表名：permissions**
```sql
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    resource VARCHAR(128) NOT NULL, -- knowledge_base, document, chat
    action VARCHAR(64) NOT NULL,    -- create, read, update, delete
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**表名：role_permissions**
```sql
CREATE TABLE role_permissions (
    id SERIAL PRIMARY KEY,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role_id, permission_id)
);
```

#### 3.1.2 知识库相关表

**表名：knowledge_bases**
```sql
CREATE TABLE knowledge_bases (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    name VARCHAR(256) NOT NULL,
    english_name VARCHAR(256) NOT NULL UNIQUE,
    description TEXT,
    created_by_id INTEGER REFERENCES users(id),
    status VARCHAR(32) DEFAULT 'active', -- active, archived, deleted
    is_builtin BOOLEAN DEFAULT FALSE,
    embedding_model VARCHAR(128), -- 使用的向量化模型
    collection_name VARCHAR(256), -- Qdrant collection 名称
    document_count INTEGER DEFAULT 0,
    total_size BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_kb_english_name ON knowledge_bases(english_name);
CREATE INDEX idx_kb_created_by ON knowledge_bases(created_by_id);
CREATE INDEX idx_kb_status ON knowledge_bases(status);
```

**表名：knowledge_documents**
```sql
CREATE TABLE knowledge_documents (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    title TEXT NOT NULL,
    content TEXT, -- 解析后的文本内容
    object_key TEXT NOT NULL UNIQUE, -- MinIO 对象键
    mime VARCHAR(128),
    size BIGINT DEFAULT 0,
    checksum VARCHAR(128), -- SHA256
    knowledge_base_id VARCHAR(64) NOT NULL REFERENCES knowledge_bases(id) ON DELETE CASCADE,
    created_by_id INTEGER REFERENCES users(id),
    parsing_status VARCHAR(32) DEFAULT 'pending', -- pending, parsing, ready, vectorizing, success, failed
    error_message TEXT,
    is_builtin BOOLEAN DEFAULT FALSE,
    metadata JSONB, -- 额外元数据
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_doc_kb_id ON knowledge_documents(knowledge_base_id);
CREATE INDEX idx_doc_title ON knowledge_documents USING gin(to_tsvector('simple', title));
CREATE INDEX idx_doc_status ON knowledge_documents(parsing_status);
CREATE INDEX idx_doc_created_by ON knowledge_documents(created_by_id);
```

#### 3.1.3 对话相关表

**表名：chat_sessions**
```sql
CREATE TABLE chat_sessions (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    user_id INTEGER NOT NULL REFERENCES users(id),
    title VARCHAR(512),
    knowledge_base_ids TEXT[], -- 关联的知识库 ID 数组
    use_rag BOOLEAN DEFAULT TRUE,
    top_k INTEGER DEFAULT 3,
    similarity_threshold FLOAT,
    similarity_weight FLOAT,
    status VARCHAR(32) DEFAULT 'active', -- active, archived
    message_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_session_user_id ON chat_sessions(user_id);
CREATE INDEX idx_session_status ON chat_sessions(status);
```

**表名：chat_messages**
```sql
CREATE TABLE chat_messages (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    session_id VARCHAR(64) NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL, -- user, assistant, system
    content TEXT NOT NULL,
    rag_context TEXT, -- RAG 检索到的上下文
    rag_sources JSONB, -- 来源文档信息
    tokens_used INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_msg_session_id ON chat_messages(session_id);
CREATE INDEX idx_msg_created_at ON chat_messages(created_at);
```

#### 3.1.4 任务队列表

**表名：document_tasks**
```sql
CREATE TABLE document_tasks (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    document_id VARCHAR(64) NOT NULL REFERENCES knowledge_documents(id) ON DELETE CASCADE,
    task_type VARCHAR(32) NOT NULL, -- parse, vectorize, delete
    status VARCHAR(32) DEFAULT 'pending', -- pending, processing, success, failed
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
```

#### 3.1.5 系统配置表

**表名：system_configs**
```sql
CREATE TABLE system_configs (
    id SERIAL PRIMARY KEY,
    config_key VARCHAR(128) NOT NULL UNIQUE,
    config_value TEXT,
    config_type VARCHAR(32), -- string, int, float, bool, json
    description TEXT,
    is_encrypted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_config_key ON system_configs(config_key);
```

### 3.2 Redis 缓存设计

#### 缓存键命名规范
```
rag:user:{user_id}                    # 用户信息缓存，TTL: 1h
rag:kb:{kb_id}                        # 知识库信息缓存，TTL: 30m
rag:doc:{doc_id}                      # 文档元数据缓存，TTL: 30m
rag:session:{session_id}              # 会话信息缓存，TTL: 2h
rag:token:{token}                     # JWT Token 黑名单，TTL: token过期时间
rag:rate_limit:{user_id}:{api}       # API 限流，TTL: 1m
rag:task_queue:{task_type}            # 任务队列（List）
```

#### 缓存策略
1. **用户信息**：登录后缓存，修改时删除
2. **知识库/文档**：读多写少，采用 Cache-Aside 模式
3. **会话信息**：对话过程中频繁访问，缓存完整会话
4. **限流计数**：滑动窗口算法

### 3.3 Kafka 消息队列设计

#### Topic 设计
```
rag.document.parse       # 文档解析任务
rag.document.vectorize   # 文档向量化任务
rag.document.delete      # 文档删除任务
rag.chat.request         # 对话请求（可选异步）
rag.system.event         # 系统事件（审计日志等）
```

#### 消息格式（JSON）
```json
{
  "task_id": "uuid",
  "task_type": "parse|vectorize|delete",
  "document_id": "doc_uuid",
  "knowledge_base_id": "kb_uuid",
  "priority": 0,
  "created_at": "2024-01-01T00:00:00Z",
  "payload": {
    "object_key": "file.pdf",
    "collection_name": "rag_collection_kb_model"
  }
}
```

---

## 四、API 接口设计

### 4.1 认证接口

#### 4.1.1 用户登录
```
POST /api/v1/auth/login
Content-Type: application/json

Request:
{
  "username": "admin",
  "password": "password123"
}

Response:
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "refresh_token_string",
    "expires_in": 7200,
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "real_name": "管理员",
      "is_admin": true
    }
  }
}
```

#### 4.1.2 刷新 Token
```
POST /api/v1/auth/refresh
Authorization: Bearer {refresh_token}

Response:
{
  "code": 0,
  "data": {
    "token": "new_jwt_token",
    "expires_in": 7200
  }
}
```

#### 4.1.3 登出
```
POST /api/v1/auth/logout
Authorization: Bearer {token}

Response:
{
  "code": 0,
  "message": "logout success"
}
```

### 4.2 用户管理接口

#### 4.2.1 创建用户
```
POST /api/v1/users
Authorization: Bearer {token}
Content-Type: application/json

Request:
{
  "username": "user001",
  "password": "password123",
  "email": "user@example.com",
  "real_name": "张三",
  "phone": "13800138000",
  "role_ids": [2, 3]
}

Response:
{
  "code": 0,
  "data": {
    "id": 10,
    "username": "user001",
    "email": "user@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 4.2.2 获取用户列表
```
GET /api/v1/users?page=1&page_size=20&keyword=zhang
Authorization: Bearer {token}

Response:
{
  "code": 0,
  "data": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "items": [
      {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com",
        "real_name": "管理员",
        "status": "active",
        "roles": ["管理员"],
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

#### 4.2.3 更新用户
```
PUT /api/v1/users/{id}
Authorization: Bearer {token}

Request:
{
  "email": "newemail@example.com",
  "real_name": "新名字",
  "status": "active"
}

Response:
{
  "code": 0,
  "message": "user updated successfully"
}
```

#### 4.2.4 删除用户
```
DELETE /api/v1/users/{id}
Authorization: Bearer {token}

Response:
{
  "code": 0,
  "message": "user deleted successfully"
}
```

### 4.3 知识库管理接口

#### 4.3.1 创建知识库
```
POST /api/v1/knowledge-bases
Authorization: Bearer {token}

Request:
{
  "name": "技术文档库",
  "english_name": "tech_docs",
  "description": "存储所有技术文档",
  "embedding_model": "text-embedding-3-small"
}

Response:
{
  "code": 0,
  "data": {
    "id": "kb_uuid_123",
    "name": "技术文档库",
    "english_name": "tech_docs",
    "collection_name": "rag_collection_tech_docs_text-embedding-3-small",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 4.3.2 获取知识库列表
```
GET /api/v1/knowledge-bases?page=1&page_size=20&status=active
Authorization: Bearer {token}

Response:
{
  "code": 0,
  "data": {
    "total": 10,
    "items": [
      {
        "id": "kb_uuid_123",
        "name": "技术文档库",
        "english_name": "tech_docs",
        "description": "存储所有技术文档",
        "document_count": 150,
        "total_size": 52428800,
        "status": "active",
        "created_by": "admin",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

#### 4.3.3 获取知识库详情
```
GET /api/v1/knowledge-bases/{id}
Authorization: Bearer {token}

Response:
{
  "code": 0,
  "data": {
    "id": "kb_uuid_123",
    "name": "技术文档库",
    "english_name": "tech_docs",
    "description": "存储所有技术文档",
    "embedding_model": "text-embedding-3-small",
    "collection_name": "rag_collection_tech_docs_text-embedding-3-small",
    "document_count": 150,
    "total_size": 52428800,
    "status": "active",
    "created_by_id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 4.3.4 更新知识库
```
PUT /api/v1/knowledge-bases/{id}
Authorization: Bearer {token}

Request:
{
  "name": "新技术文档库",
  "description": "更新后的描述"
}

Response:
{
  "code": 0,
  "message": "knowledge base updated successfully"
}
```

#### 4.3.5 删除知识库
```
DELETE /api/v1/knowledge-bases/{id}
Authorization: Bearer {token}

Response:
{
  "code": 0,
  "message": "knowledge base deleted successfully"
}
```

### 4.4 文档管理接口

#### 4.4.1 上传文档
```
POST /api/v1/knowledge-bases/{kb_id}/documents
Authorization: Bearer {token}
Content-Type: multipart/form-data

Request:
- file: (binary)
- auto_vectorize: true/false (可选，默认 true)

Response:
{
  "code": 0,
  "data": {
    "id": "doc_uuid_456",
    "title": "技术文档.pdf",
    "object_key": "tech_docs/技术文档.pdf",
    "size": 1048576,
    "mime": "application/pdf",
    "parsing_status": "pending",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 4.4.2 批量上传文档
```
POST /api/v1/knowledge-bases/{kb_id}/documents/batch
Authorization: Bearer {token}
Content-Type: multipart/form-data

Request:
- files[]: (multiple files)
- auto_vectorize: true/false

Response:
{
  "code": 0,
  "data": {
    "success_count": 8,
    "failed_count": 2,
    "results": [
      {
        "filename": "doc1.pdf",
        "status": "success",
        "document_id": "doc_uuid_1"
      },
      {
        "filename": "doc2.pdf",
        "status": "failed",
        "error": "file too large"
      }
    ]
  }
}
```

#### 4.4.3 获取文档列表
```
GET /api/v1/knowledge-bases/{kb_id}/documents?page=1&page_size=20&status=success
Authorization: Bearer {token}

Response:
{
  "code": 0,
  "data": {
    "total": 150,
    "items": [
      {
        "id": "doc_uuid_456",
        "title": "技术文档.pdf",
        "size": 1048576,
        "mime": "application/pdf",
        "parsing_status": "success",
        "created_by": "admin",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

#### 4.4.4 获取文档详情
```
GET /api/v1/documents/{id}
Authorization: Bearer {token}

Response:
{
  "code": 0,
  "data": {
    "id": "doc_uuid_456",
    "title": "技术文档.pdf",
    "content": "文档解析后的文本内容...",
    "object_key": "tech_docs/技术文档.pdf",
    "size": 1048576,
    "mime": "application/pdf",
    "checksum": "sha256_hash",
    "knowledge_base_id": "kb_uuid_123",
    "parsing_status": "success",
    "created_by_id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 4.4.5 下载文档
```
GET /api/v1/documents/{id}/download
Authorization: Bearer {token}

Response:
Content-Type: application/pdf
Content-Disposition: attachment; filename="技术文档.pdf"

(binary file content)
```

#### 4.4.6 重新向量化文档
```
POST /api/v1/documents/{id}/vectorize
Authorization: Bearer {token}

Response:
{
  "code": 0,
  "message": "vectorization task submitted",
  "data": {
    "task_id": "task_uuid_789"
  }
}
```

#### 4.4.7 删除文档
```
DELETE /api/v1/documents/{id}
Authorization: Bearer {token}

Response:
{
  "code": 0,
  "message": "document deleted successfully"
}
```

#### 4.4.8 批量删除文档
```
POST /api/v1/documents/batch-delete
Authorization: Bearer {token}

Request:
{
  "document_ids": ["doc_uuid_1", "doc_uuid_2", "doc_uuid_3"]
}

Response:
{
  "code": 0,
  "data": {
    "success_count": 2,
    "failed_count": 1,
    "results": [
      {
        "document_id": "doc_uuid_1",
        "status": "success"
      },
      {
        "document_id": "doc_uuid_2",
        "status": "failed",
        "error": "document not found"
      }
    ]
  }
}
```

### 4.5 对话管理接口

#### 4.5.1 创建对话会话
```
POST /api/v1/chat/sessions
Authorization: Bearer {token}

Request:
{
  "title": "技术咨询",
  "knowledge_base_ids": ["kb_uuid_123", "kb_uuid_456"],
  "use_rag": true,
  "top_k": 3,
  "similarity_threshold": 0.7,
  "similarity_weight": 1.5
}

Response:
{
  "code": 0,
  "data": {
    "id": "session_uuid_111",
    "title": "技术咨询",
    "knowledge_base_ids": ["kb_uuid_123", "kb_uuid_456"],
    "use_rag": true,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 4.5.2 获取会话列表
```
GET /api/v1/chat/sessions?page=1&page_size=20
Authorization: Bearer {token}

Response:
{
  "code": 0,
  "data": {
    "total": 50,
    "items": [
      {
        "id": "session_uuid_111",
        "title": "技术咨询",
        "message_count": 15,
        "status": "active",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

#### 4.5.3 发送消息（同步）
```
POST /api/v1/chat/sessions/{session_id}/messages
Authorization: Bearer {token}

Request:
{
  "content": "什么是 Kubernetes？",
  "stream": false
}

Response:
{
  "code": 0,
  "data": {
    "id": "msg_uuid_222",
    "session_id": "session_uuid_111",
    "role": "assistant",
    "content": "Kubernetes 是一个开源的容器编排平台...",
    "rag_sources": [
      {
        "document_id": "doc_uuid_456",
        "title": "Kubernetes 入门.pdf",
        "score": 0.95,
        "snippet": "Kubernetes 是..."
      }
    ],
    "tokens_used": 150,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 4.5.4 发送消息（流式）
```
POST /api/v1/chat/sessions/{session_id}/messages
Authorization: Bearer {token}

Request:
{
  "content": "什么是 Kubernetes？",
  "stream": true
}

Response (Server-Sent Events):
Content-Type: text/event-stream

data: {"type":"start","message_id":"msg_uuid_222"}

data: {"type":"content","delta":"Kubernetes"}

data: {"type":"content","delta":" 是一个"}

data: {"type":"content","delta":"开源的容器编排平台"}

data: {"type":"sources","sources":[{"document_id":"doc_uuid_456","title":"Kubernetes 入门.pdf","score":0.95}]}

data: {"type":"done","tokens_used":150}
```

#### 4.5.5 获取会话消息历史
```
GET /api/v1/chat/sessions/{session_id}/messages?page=1&page_size=50
Authorization: Bearer {token}

Response:
{
  "code": 0,
  "data": {
    "total": 30,
    "items": [
      {
        "id": "msg_uuid_222",
        "role": "user",
        "content": "什么是 Kubernetes？",
        "created_at": "2024-01-01T00:00:00Z"
      },
      {
        "id": "msg_uuid_223",
        "role": "assistant",
        "content": "Kubernetes 是一个开源的容器编排平台...",
        "rag_sources": [...],
        "created_at": "2024-01-01T00:00:01Z"
      }
    ]
  }
}
```

#### 4.5.6 删除会话
```
DELETE /api/v1/chat/sessions/{session_id}
Authorization: Bearer {token}

Response:
{
  "code": 0,
  "message": "session deleted successfully"
}
```

#### 4.5.7 清空会话消息
```
POST /api/v1/chat/sessions/{session_id}/clear
Authorization: Bearer {token}

Response:
{
  "code": 0,
  "message": "session messages cleared"
}
```

---

## 五、gRPC 接口设计（Go ↔ Python）

### 5.1 Proto 定义

**文件：ragljx/proto/rag_service.proto**

```protobuf
syntax = "proto3";

package rag;

option go_package = "ragljx/proto/rag";

// RAG 服务定义
service RAGService {
  // 解析文档
  rpc ParseDocument(ParseDocumentRequest) returns (ParseDocumentResponse);

  // 向量化文档
  rpc VectorizeDocument(VectorizeDocumentRequest) returns (VectorizeDocumentResponse);

  // 删除文档向量
  rpc DeleteDocumentVectors(DeleteDocumentVectorsRequest) returns (DeleteDocumentVectorsResponse);

  // RAG 对话
  rpc Chat(ChatRequest) returns (ChatResponse);

  // RAG 对话（流式）
  rpc ChatStream(ChatRequest) returns (stream ChatStreamResponse);

  // 检索相关文档
  rpc RetrieveDocuments(RetrieveDocumentsRequest) returns (RetrieveDocumentsResponse);
}

// 解析文档请求
message ParseDocumentRequest {
  string document_id = 1;
  string object_key = 2;
  bytes file_content = 3;
  string mime_type = 4;
}

// 解析文档响应
message ParseDocumentResponse {
  bool success = 1;
  string content = 2;
  string error_message = 3;
}

// 向量化文档请求
message VectorizeDocumentRequest {
  string document_id = 1;
  string content = 2;
  string knowledge_base_id = 3;
  string collection_name = 4;
  string title = 5;
  string object_key = 6;
}

// 向量化文档响应
message VectorizeDocumentResponse {
  bool success = 1;
  string error_message = 2;
}

// 删除文档向量请求
message DeleteDocumentVectorsRequest {
  repeated string document_ids = 1;
  string collection_name = 2;
}

// 删除文档向量响应
message DeleteDocumentVectorsResponse {
  bool success = 1;
  int32 deleted_count = 2;
  string error_message = 3;
}

// 对话请求
message ChatRequest {
  string query = 1;
  bool use_rag = 2;
  repeated string knowledge_base_ids = 3;
  int32 top_k = 4;
  float similarity_threshold = 5;
  float similarity_weight = 6;
  repeated ChatMessage history = 7;
}

// 对话消息
message ChatMessage {
  string role = 1;  // user, assistant, system
  string content = 2;
}

// 对话响应
message ChatResponse {
  string content = 1;
  repeated RAGSource sources = 2;
  int32 tokens_used = 3;
}

// 流式对话响应
message ChatStreamResponse {
  string type = 1;  // start, content, sources, done
  string delta = 2;
  repeated RAGSource sources = 3;
  int32 tokens_used = 4;
}

// RAG 来源
message RAGSource {
  string document_id = 1;
  string title = 2;
  float score = 3;
  string snippet = 4;
}

// 检索文档请求
message RetrieveDocumentsRequest {
  string query = 1;
  repeated string knowledge_base_ids = 2;
  int32 top_k = 3;
  float similarity_threshold = 4;
}

// 检索文档响应
message RetrieveDocumentsResponse {
  repeated RAGSource documents = 1;
}
```

---

## 六、Go 服务详细设计

### 6.1 项目目录结构

```
ragljx_go/
├── cmd/
│   └── server/
│       └── main.go                 # 程序入口
├── internal/
│   ├── api/                        # API 层
│   │   ├── auth/                   # 认证相关接口
│   │   │   ├── login.go
│   │   │   ├── logout.go
│   │   │   └── refresh.go
│   │   ├── user/                   # 用户管理接口
│   │   │   ├── create.go
│   │   │   ├── list.go
│   │   │   ├── update.go
│   │   │   └── delete.go
│   │   ├── knowledge_base/         # 知识库接口
│   │   │   ├── create.go
│   │   │   ├── list.go
│   │   │   ├── get.go
│   │   │   ├── update.go
│   │   │   └── delete.go
│   │   ├── document/               # 文档接口
│   │   │   ├── upload.go
│   │   │   ├── list.go
│   │   │   ├── get.go
│   │   │   ├── download.go
│   │   │   ├── vectorize.go
│   │   │   └── delete.go
│   │   └── chat/                   # 对话接口
│   │       ├── session.go
│   │       ├── message.go
│   │       └── stream.go
│   ├── service/                    # 业务逻辑层
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── knowledge_base_service.go
│   │   ├── document_service.go
│   │   └── chat_service.go
│   ├── repository/                 # 数据访问层
│   │   ├── user_repo.go
│   │   ├── knowledge_base_repo.go
│   │   ├── document_repo.go
│   │   ├── chat_session_repo.go
│   │   └── chat_message_repo.go
│   ├── model/                      # 数据模型
│   │   ├── user.go
│   │   ├── role.go
│   │   ├── permission.go
│   │   ├── knowledge_base.go
│   │   ├── document.go
│   │   ├── chat_session.go
│   │   └── chat_message.go
│   ├── middleware/                 # 中间件
│   │   ├── auth.go                 # JWT 认证
│   │   ├── cors.go                 # 跨域
│   │   ├── logger.go               # 日志
│   │   ├── recovery.go             # 异常恢复
│   │   └── rate_limit.go           # 限流
│   ├── grpc_client/                # gRPC 客户端
│   │   └── rag_client.go
│   └── pkg/                        # 工具包
│       ├── response/               # 统一响应
│       ├── errors/                 # 错误处理
│       ├── jwt/                    # JWT 工具
│       ├── validator/              # 参数验证
│       └── utils/                  # 通用工具
├── ioc/                            # IOC 容器
│   ├── config/                     # 配置对象
│   │   ├── datasource/             # 数据库配置
│   │   │   └── postgres.go
│   │   ├── redis/                  # Redis 配置
│   │   │   └── redis.go
│   │   ├── kafka/                  # Kafka 配置
│   │   │   └── kafka.go
│   │   ├── minio/                  # MinIO 配置
│   │   │   └── minio.go
│   │   ├── grpc/                   # gRPC 配置
│   │   │   └── client.go
│   │   ├── http/                   # HTTP 服务配置
│   │   │   └── gin.go
│   │   └── log/                    # 日志配置
│   │       └── logger.go
│   ├── interface.go                # IOC 接口定义
│   ├── object.go                   # 对象基类
│   ├── store.go                    # 对象存储
│   └── load.go                     # 配置加载
├── proto/                          # gRPC proto 文件
│   └── rag/
│       ├── rag_service.proto
│       └── rag_service.pb.go       # 生成的代码
├── config/                         # 配置文件
│   ├── application.yaml            # 应用配置
│   └── application.example.yaml
├── migrations/                     # 数据库迁移
│   ├── 001_init_users.sql
│   ├── 002_init_knowledge_bases.sql
│   └── 003_init_chat_sessions.sql
├── scripts/                        # 脚本
│   ├── build.sh
│   ├── migrate.sh
│   └── proto_gen.sh
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### 6.2 IOC 容器设计

#### 6.2.1 PostgreSQL 配置对象

**文件：ioc/config/datasource/postgres.go**

```go
package datasource

import (
    "context"
    "fmt"
    "ragljx/ioc"
    "ragljx/ioc/config/log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func init() {
    ioc.Config().Registry(&PostgresDB{})
}

type PostgresDB struct {
    ioc.ObjectImpl
    Host     string `env:"POSTGRES_HOST" yaml:"host"`
    Port     int    `env:"POSTGRES_PORT" yaml:"port"`
    Database string `env:"POSTGRES_DB" yaml:"database"`
    Username string `env:"POSTGRES_USER" yaml:"username"`
    Password string `env:"POSTGRES_PASSWORD" yaml:"password"`
    SSLMode  string `env:"POSTGRES_SSLMODE" yaml:"ssl_mode"`
    Debug    bool   `env:"POSTGRES_DEBUG" yaml:"debug"`

    db  *gorm.DB
    log *log.Logger
}

func (p *PostgresDB) Name() string {
    return "postgres"
}

func (p *PostgresDB) Priority() int {
    return 700
}

func (p *PostgresDB) Init() error {
    p.log = log.Sub("postgres")

    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
        p.Host, p.Username, p.Password, p.Database, p.Port, p.SSLMode,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("failed to connect to postgres: %w", err)
    }

    if p.Debug {
        db = db.Debug()
    }

    p.db = db
    p.log.Info("postgres connected successfully")
    return nil
}

func (p *PostgresDB) Close(ctx context.Context) {
    if p.db != nil {
        sqlDB, _ := p.db.DB()
        if sqlDB != nil {
            sqlDB.Close()
        }
    }
}

func (p *PostgresDB) DB() *gorm.DB {
    return p.db
}

// 全局获取方法
func Get() *gorm.DB {
    return ioc.Config().Get("postgres").(*PostgresDB).DB()
}
```

#### 6.2.2 Redis 配置对象

**文件：ioc/config/redis/redis.go**

```go
package redis

import (
    "context"
    "ragljx/ioc"
    "ragljx/ioc/config/log"

    "github.com/redis/go-redis/v9"
)

func init() {
    ioc.Config().Registry(&Redis{})
}

type Redis struct {
    ioc.ObjectImpl
    Endpoints []string `env:"REDIS_ENDPOINTS" envSeparator:"," yaml:"endpoints"`
    Password  string   `env:"REDIS_PASSWORD" yaml:"password"`
    DB        int      `env:"REDIS_DB" yaml:"db"`

    client redis.UniversalClient
    log    *log.Logger
}

func (r *Redis) Name() string {
    return "redis"
}

func (r *Redis) Priority() int {
    return 698
}

func (r *Redis) Init() error {
    r.log = log.Sub("redis")

    r.client = redis.NewUniversalClient(&redis.UniversalOptions{
        Addrs:    r.Endpoints,
        Password: r.Password,
        DB:       r.DB,
    })

    // 测试连接
    ctx := context.Background()
    if err := r.client.Ping(ctx).Err(); err != nil {
        return err
    }

    r.log.Info("redis connected successfully")
    return nil
}

func (r *Redis) Close(ctx context.Context) {
    if r.client != nil {
        r.client.Close()
    }
}

func (r *Redis) Client() redis.UniversalClient {
    return r.client
}

func Get() redis.UniversalClient {
    return ioc.Config().Get("redis").(*Redis).Client()
}
```

#### 6.2.3 Kafka 配置对象

**文件：ioc/config/kafka/kafka.go**

```go
package kafka

import (
    "context"
    "ragljx/ioc"
    "ragljx/ioc/config/log"

    "github.com/segmentio/kafka-go"
)

func init() {
    ioc.Config().Registry(&Kafka{})
}

type Kafka struct {
    ioc.ObjectImpl
    Brokers  []string `env:"KAFKA_BROKERS" envSeparator:"," yaml:"brokers"`
    Username string   `env:"KAFKA_USERNAME" yaml:"username"`
    Password string   `env:"KAFKA_PASSWORD" yaml:"password"`

    log *log.Logger
}

func (k *Kafka) Name() string {
    return "kafka"
}

func (k *Kafka) Priority() int {
    return 696
}

func (k *Kafka) Init() error {
    k.log = log.Sub("kafka")
    k.log.Info("kafka initialized")
    return nil
}

func (k *Kafka) Close(ctx context.Context) {
    // Kafka writers/readers 由使用方管理
}

// 创建 Producer
func (k *Kafka) Producer(topic string) *kafka.Writer {
    return &kafka.Writer{
        Addr:                   kafka.TCP(k.Brokers...),
        Topic:                  topic,
        Balancer:               &kafka.LeastBytes{},
        AllowAutoTopicCreation: true,
    }
}

// 创建 Consumer
func (k *Kafka) Consumer(groupID string, topics []string) *kafka.Reader {
    return kafka.NewReader(kafka.ReaderConfig{
        Brokers: k.Brokers,
        GroupID: groupID,
        Topic:   topics[0], // 简化处理，实际可支持多 topic
    })
}

func Get() *Kafka {
    return ioc.Config().Get("kafka").(*Kafka)
}
```

#### 6.2.4 MinIO 配置对象

**文件：ioc/config/minio/minio.go**

```go
package minio

import (
    "context"
    "ragljx/ioc"
    "ragljx/ioc/config/log"

    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

func init() {
    ioc.Config().Registry(&MinIO{})
}

type MinIO struct {
    ioc.ObjectImpl
    Endpoint        string `env:"MINIO_ENDPOINT" yaml:"endpoint"`
    AccessKeyID     string `env:"MINIO_ACCESS_KEY" yaml:"access_key_id"`
    SecretAccessKey string `env:"MINIO_SECRET_KEY" yaml:"secret_access_key"`
    UseSSL          bool   `env:"MINIO_USE_SSL" yaml:"use_ssl"`
    BucketName      string `env:"MINIO_BUCKET" yaml:"bucket_name"`

    client *minio.Client
    log    *log.Logger
}

func (m *MinIO) Name() string {
    return "minio"
}

func (m *MinIO) Priority() int {
    return 695
}

func (m *MinIO) Init() error {
    m.log = log.Sub("minio")

    client, err := minio.New(m.Endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(m.AccessKeyID, m.SecretAccessKey, ""),
        Secure: m.UseSSL,
    })
    if err != nil {
        return err
    }

    m.client = client

    // 确保 bucket 存在
    ctx := context.Background()
    exists, err := client.BucketExists(ctx, m.BucketName)
    if err != nil {
        return err
    }

    if !exists {
        err = client.MakeBucket(ctx, m.BucketName, minio.MakeBucketOptions{})
        if err != nil {
            return err
        }
        m.log.Info("created bucket: %s", m.BucketName)
    }

    m.log.Info("minio connected successfully")
    return nil
}

func (m *MinIO) Close(ctx context.Context) {
    // MinIO client 不需要显式关闭
}

func (m *MinIO) Client() *minio.Client {
    return m.client
}

func (m *MinIO) Bucket() string {
    return m.BucketName
}

func Get() *MinIO {
    return ioc.Config().Get("minio").(*MinIO)
}
```

### 6.3 中间件设计

#### 6.3.1 JWT 认证中间件

**文件：internal/middleware/auth.go**

```go
package middleware

import (
    "net/http"
    "strings"
    "ragljx/internal/pkg/jwt"
    "ragljx/internal/pkg/response"

    "github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从 Header 获取 token
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            response.Error(c, http.StatusUnauthorized, "missing authorization header")
            c.Abort()
            return
        }

        // 解析 Bearer token
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            response.Error(c, http.StatusUnauthorized, "invalid authorization header")
            c.Abort()
            return
        }

        tokenString := parts[1]

        // 验证 token
        claims, err := jwt.ParseToken(tokenString)
        if err != nil {
            response.Error(c, http.StatusUnauthorized, "invalid token")
            c.Abort()
            return
        }

        // 将用户信息存入上下文
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("is_admin", claims.IsAdmin)

        c.Next()
    }
}

// 可选的认证中间件（允许未登录访问）
func OptionalAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader != "" {
            parts := strings.SplitN(authHeader, " ", 2)
            if len(parts) == 2 && parts[0] == "Bearer" {
                claims, err := jwt.ParseToken(parts[1])
                if err == nil {
                    c.Set("user_id", claims.UserID)
                    c.Set("username", claims.Username)
                    c.Set("is_admin", claims.IsAdmin)
                }
            }
        }
        c.Next()
    }
}
```

---

## 七、Python 服务详细设计

### 7.1 项目目录结构

```
ragljx_py/
├── app/
│   ├── __init__.py
│   ├── main.py                     # 程序入口
│   ├── config.py                   # 配置管理
│   ├── grpc_server/                # gRPC 服务
│   │   ├── __init__.py
│   │   ├── server.py               # gRPC 服务器
│   │   └── rag_service.py          # RAG 服务实现
│   ├── services/                   # 业务服务
│   │   ├── __init__.py
│   │   ├── file_service.py         # 文件解析服务
│   │   ├── vector_service.py       # 向量服务
│   │   └── chat_service.py         # 对话服务
│   ├── utils/                      # 工具类
│   │   ├── __init__.py
│   │   ├── minio_client.py         # MinIO 客户端
│   │   └── logger.py               # 日志工具
│   └── proto/                      # gRPC proto 生成文件
│       ├── __init__.py
│       ├── rag_service_pb2.py
│       └── rag_service_pb2_grpc.py
├── requirements.txt
├── Dockerfile
├── docker-compose.yml
└── README.md
```

### 7.2 配置管理

**文件：app/config.py**

```python
import os
from typing import Optional

class Config:
    """配置类"""

    # gRPC 服务配置
    GRPC_HOST: str = os.getenv("GRPC_HOST", "0.0.0.0")
    GRPC_PORT: int = int(os.getenv("GRPC_PORT", "50051"))

    # Qdrant 配置
    QDRANT_HOST: str = os.getenv("QDRANT_HOST", "localhost")
    QDRANT_PORT: int = int(os.getenv("QDRANT_PORT", "6333"))
    QDRANT_API_KEY: Optional[str] = os.getenv("QDRANT_API_KEY")
    QDRANT_COLLECTION_PREFIX: str = os.getenv("QDRANT_COLLECTION_PREFIX", "rag_collection")

    # Embedding 模型配置
    EMBEDDING_API_BASE: str = os.getenv("EMBEDDING_API_BASE", "https://api.openai.com/v1")
    EMBEDDING_API_KEY: str = os.getenv("EMBEDDING_API_KEY", "")
    EMBEDDING_MODEL: str = os.getenv("EMBEDDING_MODEL", "text-embedding-3-small")
    EMBEDDING_DIMENSION: int = int(os.getenv("EMBEDDING_DIMENSION", "1536"))

    # LLM 配置
    LLM_API_BASE: str = os.getenv("LLM_API_BASE", "https://api.openai.com/v1")
    LLM_API_KEY: str = os.getenv("LLM_API_KEY", "")
    LLM_MODEL: str = os.getenv("LLM_MODEL", "gpt-4")

    # MinIO 配置
    MINIO_ENDPOINT: str = os.getenv("MINIO_ENDPOINT", "localhost:9000")
    MINIO_ACCESS_KEY: str = os.getenv("MINIO_ACCESS_KEY", "minioadmin")
    MINIO_SECRET_KEY: str = os.getenv("MINIO_SECRET_KEY", "minioadmin")
    MINIO_BUCKET: str = os.getenv("MINIO_BUCKET", "ragljx")
    MINIO_USE_SSL: bool = os.getenv("MINIO_USE_SSL", "false").lower() == "true"

    # 日志配置
    LOG_LEVEL: str = os.getenv("LOG_LEVEL", "INFO")

config = Config()
```

### 7.3 gRPC 服务实现

**文件：app/grpc_server/rag_service.py**

```python
import grpc
from concurrent import futures
import logging
from app.proto import rag_service_pb2, rag_service_pb2_grpc
from app.services.file_service import RagFileService
from app.services.vector_service import RagVectorService
from app.services.chat_service import RagChatService

logger = logging.getLogger(__name__)

class RAGServiceImpl(rag_service_pb2_grpc.RAGServiceServicer):
    """RAG 服务实现"""

    def __init__(self):
        self.file_service = RagFileService()
        self.vector_service = RagVectorService()
        self.chat_service = RagChatService()

    def ParseDocument(self, request, context):
        """解析文档"""
        try:
            # 调用文件服务解析
            content = self.file_service.parse_file_content(
                request.file_content,
                request.mime_type
            )

            return rag_service_pb2.ParseDocumentResponse(
                success=True,
                content=content
            )
        except Exception as e:
            logger.error(f"Parse document failed: {e}")
            return rag_service_pb2.ParseDocumentResponse(
                success=False,
                error_message=str(e)
            )

    def VectorizeDocument(self, request, context):
        """向量化文档"""
        try:
            self.vector_service.upsert_document(
                doc_id=request.document_id,
                content=request.content,
                knowledge_base_id=request.knowledge_base_id,
                title=request.title,
                object_key=request.object_key,
                collection_name=request.collection_name
            )

            return rag_service_pb2.VectorizeDocumentResponse(success=True)
        except Exception as e:
            logger.error(f"Vectorize document failed: {e}")
            return rag_service_pb2.VectorizeDocumentResponse(
                success=False,
                error_message=str(e)
            )

    def DeleteDocumentVectors(self, request, context):
        """删除文档向量"""
        try:
            self.vector_service.delete_vectors(
                doc_ids=list(request.document_ids),
                collection_name=request.collection_name
            )

            return rag_service_pb2.DeleteDocumentVectorsResponse(
                success=True,
                deleted_count=len(request.document_ids)
            )
        except Exception as e:
            logger.error(f"Delete vectors failed: {e}")
            return rag_service_pb2.DeleteDocumentVectorsResponse(
                success=False,
                error_message=str(e)
            )

    def Chat(self, request, context):
        """RAG 对话"""
        try:
            # 转换历史消息
            history = [
                {"role": msg.role, "content": msg.content}
                for msg in request.history
            ]

            # 调用对话服务
            response_content = self.chat_service.chat(
                query=request.query,
                use_rag=request.use_rag,
                knowledge_base_ids=list(request.knowledge_base_ids) if request.knowledge_base_ids else None,
                top_k=request.top_k or 3,
                history=history
            )

            # 获取 RAG 来源
            sources = []
            if request.use_rag:
                selected = self.chat_service.get_rag_selectd(
                    query=request.query,
                    knowledge_base_ids=list(request.knowledge_base_ids),
                    top_k=request.top_k or 3
                )
                for node in selected:
                    sources.append(rag_service_pb2.RAGSource(
                        document_id=node.metadata.get('object_key', ''),
                        title=node.metadata.get('title', ''),
                        score=getattr(node, 'score', 0.0),
                        snippet=node.get_content()[:200]
                    ))

            return rag_service_pb2.ChatResponse(
                content=response_content,
                sources=sources,
                tokens_used=0  # 可以从 LLM 响应中获取
            )
        except Exception as e:
            logger.error(f"Chat failed: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return rag_service_pb2.ChatResponse()

    def ChatStream(self, request, context):
        """流式对话（简化实现）"""
        # 实际应该实现流式返回
        response = self.Chat(request, context)
        yield rag_service_pb2.ChatStreamResponse(
            type="content",
            delta=response.content
        )
        yield rag_service_pb2.ChatStreamResponse(
            type="done",
            tokens_used=response.tokens_used
        )
```

---

## 八、Docker 部署配置

### 8.1 Docker Compose 配置

**文件：docker-compose.yml**

```yaml
version: '3.8'

services:
  # PostgreSQL 数据库
  postgres:
    image: postgres:15-alpine
    container_name: ragljx_postgres
    environment:
      POSTGRES_DB: ragljx
      POSTGRES_USER: ragljx
      POSTGRES_PASSWORD: ragljx_password
      POSTGRES_INITDB_ARGS: "--encoding=UTF8"
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./ragljx_go/migrations:/docker-entrypoint-initdb.d
    networks:
      - ragljx_network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ragljx"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis 缓存
  redis:
    image: redis:7-alpine
    container_name: ragljx_redis
    command: redis-server --requirepass ragljx_redis_password
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - ragljx_network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Kafka
  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    container_name: ragljx_zookeeper
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    networks:
      - ragljx_network

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    container_name: ragljx_kafka
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    networks:
      - ragljx_network

  # MinIO 对象存储
  minio:
    image: minio/minio:latest
    container_name: ragljx_minio
    command: server /data --console-address ":9001"
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - minio_data:/data
    networks:
      - ragljx_network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3

  # Qdrant 向量数据库
  qdrant:
    image: qdrant/qdrant:latest
    container_name: ragljx_qdrant
    ports:
      - "6333:6333"
      - "6334:6334"
    volumes:
      - qdrant_data:/qdrant/storage
    networks:
      - ragljx_network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:6333/health"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Go 后端服务
  ragljx_go:
    build:
      context: ./ragljx_go
      dockerfile: Dockerfile
    container_name: ragljx_go
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      kafka:
        condition: service_started
      minio:
        condition: service_healthy
      ragljx_py:
        condition: service_started
    ports:
      - "8080:8080"
    environment:
      # PostgreSQL
      POSTGRES_HOST: postgres
      POSTGRES_PORT: 5432
      POSTGRES_DB: ragljx
      POSTGRES_USER: ragljx
      POSTGRES_PASSWORD: ragljx_password
      POSTGRES_SSLMODE: disable
      # Redis
      REDIS_ENDPOINTS: redis:6379
      REDIS_PASSWORD: ragljx_redis_password
      REDIS_DB: 0
      # Kafka
      KAFKA_BROKERS: kafka:29092
      # MinIO
      MINIO_ENDPOINT: minio:9000
      MINIO_ACCESS_KEY: minioadmin
      MINIO_SECRET_KEY: minioadmin
      MINIO_BUCKET: ragljx
      MINIO_USE_SSL: "false"
      # gRPC
      GRPC_PYTHON_ADDR: ragljx_py:50051
      # JWT
      JWT_SECRET_KEY: your-secret-key-change-in-production
      JWT_EXPIRE_HOURS: 2
    networks:
      - ragljx_network
    restart: unless-stopped

  # Python AI 服务
  ragljx_py:
    build:
      context: ./ragljx_py
      dockerfile: Dockerfile
    container_name: ragljx_py
    depends_on:
      qdrant:
        condition: service_healthy
      minio:
        condition: service_healthy
    ports:
      - "50051:50051"
    environment:
      # gRPC
      GRPC_HOST: 0.0.0.0
      GRPC_PORT: 50051
      # Qdrant
      QDRANT_HOST: qdrant
      QDRANT_PORT: 6333
      QDRANT_COLLECTION_PREFIX: rag_collection
      # Embedding
      EMBEDDING_API_BASE: ${EMBEDDING_API_BASE}
      EMBEDDING_API_KEY: ${EMBEDDING_API_KEY}
      EMBEDDING_MODEL: ${EMBEDDING_MODEL:-text-embedding-3-small}
      EMBEDDING_DIMENSION: ${EMBEDDING_DIMENSION:-1536}
      # LLM
      LLM_API_BASE: ${LLM_API_BASE}
      LLM_API_KEY: ${LLM_API_KEY}
      LLM_MODEL: ${LLM_MODEL:-gpt-4}
      # MinIO
      MINIO_ENDPOINT: minio:9000
      MINIO_ACCESS_KEY: minioadmin
      MINIO_SECRET_KEY: minioadmin
      MINIO_BUCKET: ragljx
      MINIO_USE_SSL: "false"
      # Log
      LOG_LEVEL: INFO
    networks:
      - ragljx_network
    restart: unless-stopped

  # Vue 前端服务
  ragljx_web:
    build:
      context: ./ragljx_web
      dockerfile: Dockerfile
    container_name: ragljx_web
    depends_on:
      - ragljx_go
    ports:
      - "80:80"
    networks:
      - ragljx_network
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
  minio_data:
  qdrant_data:

networks:
  ragljx_network:
    driver: bridge
```

### 8.2 Go 服务 Dockerfile

**文件：ragljx_go/Dockerfile**

```dockerfile
# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git make

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ragljx_server ./cmd/server

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/ragljx_server .
COPY --from=builder /app/config ./config

# 设置时区
ENV TZ=Asia/Shanghai

EXPOSE 8080

CMD ["./ragljx_server"]
```

### 8.3 Python 服务 Dockerfile

**文件：ragljx_py/Dockerfile**

```dockerfile
FROM python:3.11-slim

WORKDIR /app

# 安装系统依赖
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    && rm -rf /var/lib/apt/lists/*

# 复制依赖文件
COPY requirements.txt .

# 安装 Python 依赖
RUN pip install --no-cache-dir -r requirements.txt

# 复制源代码
COPY . .

# 设置环境变量
ENV PYTHONUNBUFFERED=1

EXPOSE 50051

CMD ["python", "-m", "app.main"]
```

### 8.4 前端 Dockerfile

**文件：ragljx_web/Dockerfile**

```dockerfile
# 构建阶段
FROM node:18-alpine AS builder

WORKDIR /app

# 复制 package 文件
COPY package*.json ./

# 安装依赖
RUN npm ci

# 复制源代码
COPY . .

# 构建
RUN npm run build

# 运行阶段
FROM nginx:alpine

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

---

## 九、开发流程与规范

### 9.1 开发环境搭建

#### 9.1.1 前置要求
- Go 1.21+
- Python 3.11+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

#### 9.1.2 本地开发步骤

**1. 克隆项目**
```bash
cd /Users/liang/projectljx/ragljx
```

**2. 启动基础设施**
```bash
docker-compose up -d postgres redis kafka minio qdrant
```

**3. 配置环境变量**
```bash
# Go 服务
cd ragljx_go
cp config/application.example.yaml config/application.yaml
# 编辑配置文件

# Python 服务
cd ../ragljx_py
cp .env.example .env
# 编辑环境变量
```

**4. 数据库迁移**
```bash
cd ragljx_go
make migrate
```

**5. 启动 Go 服务**
```bash
cd ragljx_go
go run cmd/server/main.go
```

**6. 启动 Python 服务**
```bash
cd ragljx_py
python -m app.main
```

**7. 启动前端服务**
```bash
cd ragljx_web
npm install
npm run dev
```

### 9.2 代码规范

#### 9.2.1 Go 代码规范
- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 使用 `golangci-lint` 进行静态检查
- 函数注释使用英文
- 错误处理必须显式处理

#### 9.2.2 Python 代码规范
- 遵循 PEP 8 规范
- 使用 `black` 格式化代码
- 使用 `flake8` 进行静态检查
- 使用类型注解（Type Hints）
- 文档字符串使用 Google 风格

#### 9.2.3 前端代码规范
- 遵循 Vue 3 官方风格指南
- 使用 ESLint + Prettier
- 组件命名使用 PascalCase
- 文件命名使用 kebab-case

### 9.3 Git 工作流

#### 分支策略
- `main`: 主分支，生产环境代码
- `develop`: 开发分支
- `feature/*`: 功能分支
- `bugfix/*`: 修复分支
- `release/*`: 发布分支

#### Commit 规范
```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type 类型：**
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具链相关

---

## 十、总结与后续规划

### 10.1 系统特点

1. **微服务架构**：Go 和 Python 服务分离，职责清晰
2. **IOC 容器**：统一的依赖注入和配置管理
3. **完整的认证授权**：基于 JWT 的用户认证和 RBAC 权限控制
4. **高性能**：Redis 缓存、Kafka 异步处理
5. **可扩展**：支持水平扩展，易于添加新功能
6. **容器化部署**：Docker Compose 一键部署

### 10.2 核心功能清单

- [x] 用户认证与授权（JWT + RBAC）
- [x] 知识库管理（CRUD）
- [x] 文档上传与解析（多格式支持）
- [x] 文档向量化（Qdrant）
- [x] RAG 对话（同步/流式）
- [x] 会话管理
- [x] 对象存储（MinIO）
- [x] 消息队列（Kafka）
- [x] 缓存机制（Redis）

### 10.3 后续开发建议

#### 第一阶段：核心功能实现（2-3周）
1. 搭建 Go 服务基础框架（IOC、路由、中间件）
2. 实现用户认证模块
3. 实现知识库和文档管理基础 CRUD
4. 搭建 Python gRPC 服务
5. 实现文档解析和向量化

#### 第二阶段：RAG 功能完善（2周）
1. 实现 RAG 对话功能
2. 优化向量检索性能
3. 实现流式对话
4. 添加对话历史管理

#### 第三阶段：前端开发（2-3周）
1. 搭建 Vue 3 项目框架
2. 实现登录注册页面
3. 实现知识库管理界面
4. 实现文档管理界面
5. 实现对话界面

#### 第四阶段：优化与测试（1-2周）
1. 性能优化
2. 单元测试和集成测试
3. 文档完善
4. Docker 部署测试

### 10.4 注意事项

1. **安全性**：
   - JWT Secret 必须使用强密码
   - 数据库密码不要使用默认值
   - 生产环境启用 HTTPS
   - 敏感配置使用环境变量

2. **性能优化**：
   - 合理使用 Redis 缓存
   - 数据库查询添加索引
   - 大文件上传使用分片
   - 向量检索结果缓存

3. **监控与日志**：
   - 添加 Prometheus 监控
   - 集成 ELK 日志系统
   - 添加链路追踪（Jaeger）

4. **扩展性**：
   - 预留接口扩展点
   - 配置项外部化
   - 支持多租户（可选）

---

## 附录

### A. 环境变量清单

**Go 服务环境变量**
```bash
# 数据库
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=ragljx
POSTGRES_USER=ragljx
POSTGRES_PASSWORD=ragljx_password
POSTGRES_SSLMODE=disable

# Redis
REDIS_ENDPOINTS=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Kafka
KAFKA_BROKERS=localhost:9092

# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=ragljx
MINIO_USE_SSL=false

# gRPC
GRPC_PYTHON_ADDR=localhost:50051

# JWT
JWT_SECRET_KEY=your-secret-key
JWT_EXPIRE_HOURS=2

# 服务
HTTP_PORT=8080
```

**Python 服务环境变量**
```bash
# gRPC
GRPC_HOST=0.0.0.0
GRPC_PORT=50051

# Qdrant
QDRANT_HOST=localhost
QDRANT_PORT=6333
QDRANT_API_KEY=
QDRANT_COLLECTION_PREFIX=rag_collection

# Embedding
EMBEDDING_API_BASE=https://api.openai.com/v1
EMBEDDING_API_KEY=sk-xxx
EMBEDDING_MODEL=text-embedding-3-small
EMBEDDING_DIMENSION=1536

# LLM
LLM_API_BASE=https://api.openai.com/v1
LLM_API_KEY=sk-xxx
LLM_MODEL=gpt-4

# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=ragljx
MINIO_USE_SSL=false

# Log
LOG_LEVEL=INFO
```

### B. 常用命令

**Go 服务**
```bash
# 运行
go run cmd/server/main.go

# 编译
go build -o bin/ragljx_server cmd/server/main.go

# 测试
go test ./...

# 生成 proto
protoc --go_out=. --go-grpc_out=. proto/rag_service.proto
```

**Python 服务**
```bash
# 运行
python -m app.main

# 安装依赖
pip install -r requirements.txt

# 生成 proto
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. proto/rag_service.proto
```

**Docker**
```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f ragljx_go

# 重启服务
docker-compose restart ragljx_go

# 停止所有服务
docker-compose down

# 清理数据
docker-compose down -v
```

---

**文档版本：v1.0**
**最后更新：2024-01-01**
**作者：RAG 系统开发团队**


