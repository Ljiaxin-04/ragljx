# RAG 系统

基于 Go + Python + Vue 的企业级 RAG（检索增强生成）系统。

## 系统架构

- **Go 后端服务** (`ragljx_go`): 使用 Gin + IOC 框架，提供 RESTful API，负责用户管理、知识库管理、文档管理等
- **Python AI 服务** (`ragljx_py`): 使用 gRPC，负责文档解析、向量化、RAG 对话等 AI 功能
- **Vue 前端** (`ragljx_web`): 用户界面（待实现）

## 技术栈

### 后端 (Go)
- Gin Web 框架
- GORM ORM
- IOC 依赖注入
- gRPC 客户端
- PostgreSQL
- Redis
- Kafka
- MinIO

### AI 服务 (Python)
- gRPC 服务器
- LlamaIndex RAG 框架
- Qdrant 向量数据库
- OpenAI API
- 多格式文档解析（PDF, DOCX, XLSX, PPTX, HTML, CSV, JSON, XML, RTF等）

### 前端 (Vue)
- Vue 3
- JavaScript
- Element Plus (待实现)

## 快速开始

### 前置要求

- Docker & Docker Compose
- Go 1.21+
- Python 3.12+
- Node.js 18+ (前端开发)

### 使用 Docker Compose 启动（推荐）

1. 克隆项目并进入目录：
```bash
cd /Users/liang/projectljx/ragljx
```

2. 配置环境变量（可选）：
```bash
# 创建 .env 文件
cat > .env << EOF
OPENAI_API_KEY=your_openai_api_key_here
OPENAI_API_BASE=https://api.openai.com/v1
EOF
```

3. 启动所有服务：
```bash
docker-compose up -d
```

4. 查看服务状态：
```bash
docker-compose ps
```

5. 查看日志：
```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f ragljx_go
docker-compose logs -f ragljx_py
```

### 本地开发

#### Go 后端服务

```bash
cd ragljx_go

# 安装依赖
go mod tidy

# 生成 proto 文件
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/rag/rag_service.proto

# 编译
go build -o bin/ragljx cmd/server/main.go

# 运行
./bin/ragljx
```

#### Python AI 服务

```bash
cd ragljx_py

# 创建虚拟环境
python3 -m venv venv
source venv/bin/activate  # Linux/Mac
# 或 venv\Scripts\activate  # Windows

# 安装依赖
pip install --trusted-host pypi.org --trusted-host files.pythonhosted.org -r requirements.txt

# 生成 proto 文件
python3 -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. app/proto/rag_service.proto

# 运行
python main.py
```

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| Go 后端 API | 8080 | RESTful API |
| Python gRPC | 50051 | AI 服务 gRPC 接口 |
| PostgreSQL | 5432 | 数据库 |
| Redis | 6379 | 缓存 |
| Kafka | 19092 | 消息队列 |
| MinIO | 9000 | 对象存储 API |
| MinIO Console | 9001 | MinIO 管理界面 |
| Qdrant | 6333 | 向量数据库 API |
| Qdrant Dashboard | 6334 | Qdrant 管理界面 |

## API 文档

### 认证接口

- `POST /api/auth/login` - 用户登录
- `POST /api/auth/logout` - 用户登出
- `POST /api/auth/refresh` - 刷新 Token

### 用户管理

- `GET /api/users` - 获取用户列表
- `GET /api/users/:id` - 获取用户详情
- `POST /api/users` - 创建用户
- `PUT /api/users/:id` - 更新用户
- `DELETE /api/users/:id` - 删除用户

### 知识库管理

- `GET /api/knowledge-bases` - 获取知识库列表
- `GET /api/knowledge-bases/:id` - 获取知识库详情
- `POST /api/knowledge-bases` - 创建知识库
- `PUT /api/knowledge-bases/:id` - 更新知识库
- `DELETE /api/knowledge-bases/:id` - 删除知识库

### 文档管理

- `GET /api/documents` - 获取文档列表
- `GET /api/documents/:id` - 获取文档详情
- `POST /api/documents/upload` - 上传文档
- `POST /api/documents/:id/vectorize` - 向量化文档
- `DELETE /api/documents/:id` - 删除文档

### 对话管理

- `GET /api/chat/sessions` - 获取会话列表
- `POST /api/chat/sessions` - 创建会话
- `POST /api/chat` - 发送消息（非流式）
- `POST /api/chat/stream` - 发送消息（流式）

## 配置说明

### Go 服务配置 (`ragljx_go/config/application.yaml`)

```yaml
http:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 60
  write_timeout: 60

postgres:
  host: "localhost"
  port: 5432
  database: "ragljx"
  username: "ragljx"
  password: "ragljx_password"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

# ... 其他配置
```

### Python 服务配置 (`ragljx_py/config.yaml`)

```yaml
grpc:
  host: "0.0.0.0"
  port: 50051

qdrant:
  host: "localhost"
  port: 6333

openai:
  api_key: ""  # 从环境变量读取
  api_base: "https://api.openai.com/v1"
  embedding_model: "text-embedding-3-small"
  chat_model: "gpt-4"

# ... 其他配置
```

## 数据库迁移

数据库表会在首次启动时自动创建（通过 `migrations/*.sql` 文件）。

### 默认管理员账号

系统会自动创建一个默认的超级管理员账号：

- **用户名**: `admin`
- **密码**: `123456`
- **邮箱**: `admin@ragljx.com`

⚠️ **重要提示**: 请在生产环境中立即修改默认密码！

您可以使用以下工具生成新的密码哈希：

```bash
cd ragljx_go
go run scripts/gen_password.go your_new_password
```

## 开发注意事项

1. **OpenAI API Key**: 需要配置有效的 OpenAI API Key 才能使用 AI 功能
2. **向量化模型**: 默认使用 `text-embedding-3-small`，可在配置文件中修改
3. **对话模型**: 默认使用 `gpt-4`，可在配置文件中修改
4. **文件大小限制**: gRPC 消息大小限制为 100MB

## 故障排查

### Go 服务无法连接数据库
- 检查 PostgreSQL 是否启动：`docker-compose ps postgres`
- 检查数据库配置是否正确
- 查看日志：`docker-compose logs postgres`

### Python 服务无法连接 Qdrant
- 检查 Qdrant 是否启动：`docker-compose ps qdrant`
- 访问 Qdrant Dashboard: http://localhost:6334

### gRPC 连接失败
- 确保 Python 服务已启动
- 检查端口 50051 是否被占用
- 查看 Python 服务日志：`docker-compose logs ragljx_py`

## 停止服务

```bash
# 停止所有服务
docker-compose down

# 停止并删除数据卷（谨慎使用）
docker-compose down -v
```

## 许可证

MIT License

