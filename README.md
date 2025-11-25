# RAG 知识库系统

基于 Go + Python + Vue 的企业级 RAG（检索增强生成）知识库系统，提供智能文档管理与对话功能。

## 系统架构

- **Go 后端服务** (`ragljx_go`): 使用 Gin + IOC 框架，提供 RESTful API，负责用户管理、知识库管理、文档管理、会话管理等
- **Python AI 服务** (`ragljx_py`): 使用 gRPC，负责文档解析、向量化、RAG 对话等 AI 功能
- **Vue 前端** (`ragljx_web`): 现代化的用户界面，提供知识库管理、文档上传、智能对话等功能

## 功能特性

✅ **用户管理**
- 用户注册、登录、JWT 认证
- 基于角色的权限控制（RBAC）
- 个人信息管理、密码修改

✅ **知识库管理**
- 创建、编辑、删除知识库
- 自定义嵌入模型和分块策略
- 知识库统计信息展示

✅ **文档管理**
- 支持多种文档格式（TXT, MD, PDF, DOCX, XLSX, PPTX, HTML, CSV, JSON, XML, RTF）
- 文档上传、自动解析、向量化
- 文档状态跟踪、重新处理

✅ **智能对话**
- 基于知识库的 RAG 对话
- 流式输出，实时响应
- 多知识库联合检索
- 对话历史管理
- 来源文档追溯

✅ **系统管理**
- 用户管理（管理员功能）
- 系统配置管理
- 日志记录

## 技术栈

### 后端 (Go)
- **Web 框架**: Gin
- **ORM**: GORM
- **依赖注入**: 自研 IOC 容器
- **RPC**: gRPC 客户端
- **数据库**: PostgreSQL
- **缓存**: Redis
- **消息队列**: Kafka (Redpanda)
- **对象存储**: MinIO
- **认证**: JWT

### AI 服务 (Python)
- **RPC 服务**: gRPC 服务器
- **RAG 框架**: LlamaIndex
- **向量数据库**: Qdrant
- **LLM**: OpenAI API (GPT-4)
- **嵌入模型**: OpenAI Embeddings
- **文档解析**: 支持 12+ 种文档格式

### 前端 (Vue)
- **框架**: Vue 3 (Composition API)
- **构建工具**: Vite
- **UI 组件库**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router
- **HTTP 客户端**: Axios
- **语言**: JavaScript

## 快速开始

### 前置要求

- **Docker & Docker Compose** (推荐使用 Docker 部署)
- **Go 1.21+** (本地开发)
- **Python 3.12+** (本地开发)
- **Node.js 20+** (前端开发)
- **OpenAI API Key** (必需，用于 AI 功能)

### 使用 Docker Compose 启动（推荐）

1. **克隆项目并进入目录**：
```bash
cd /Users/liang/projectljx/ragljx
```

2. **配置环境变量**（必需）：
```bash
# 创建 .env 文件
cat > .env << EOF
OPENAI_API_KEY=your_openai_api_key_here
OPENAI_API_BASE=https://api.openai.com/v1
EOF
```

3. **启动所有服务**：
```bash
docker-compose up -d
```

这将启动以下服务：
- PostgreSQL (数据库)
- Redis (缓存)
- Kafka/Redpanda (消息队列)
- MinIO (对象存储)
- Qdrant (向量数据库)
- Go 后端服务
- Python AI 服务

4. **查看服务状态**：
```bash
docker-compose ps
```

5. **查看日志**：
```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f ragljx_go
docker-compose logs -f ragljx_py
```

6. **访问应用**：
- 前端界面: http://localhost:5173 (开发模式)
- 后端 API: http://localhost:8080
- MinIO 控制台: http://localhost:9001
- Qdrant 控制台: http://localhost:6334

7. **默认登录账号**：
- 用户名: `admin`
- 密码: `123456`

⚠️ **首次登录后请立即修改密码！**

### 本地开发

#### 1. Go 后端服务

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

# 运行（确保 PostgreSQL、Redis 等服务已启动）
./bin/ragljx
```

#### 2. Python AI 服务

```bash
cd ragljx_py

# 创建虚拟环境
python3 -m venv venv
source venv/bin/activate  # Linux/Mac
# 或 venv\Scripts\activate  # Windows

# 安装依赖
pip install -r requirements.txt

# 生成 proto 文件
python3 -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. app/proto/rag_service.proto

# 配置环境变量
export OPENAI_API_KEY=your_api_key_here
export OPENAI_API_BASE=https://api.openai.com/v1

# 运行（确保 Qdrant 已启动）
python main.py
```

#### 3. Vue 前端

```bash
cd ragljx_web

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 预览生产版本
npm run preview
```

前端开发服务器将在 http://localhost:5173 启动

## 服务端口

| 服务 | 端口 | 说明 | 访问地址 |
|------|------|------|----------|
| Vue 前端 | 5173 | 用户界面（开发模式） | http://localhost:5173 |
| Go 后端 API | 8080 | RESTful API | http://localhost:8080 |
| Python gRPC | 50051 | AI 服务 gRPC 接口 | - |
| PostgreSQL | 5432 | 数据库 | - |
| Redis | 6379 | 缓存 | - |
| Kafka | 19092 | 消息队列 | - |
| MinIO API | 9000 | 对象存储 API | http://localhost:9000 |
| MinIO Console | 9001 | MinIO 管理界面 | http://localhost:9001 |
| Qdrant API | 6333 | 向量数据库 API | http://localhost:6333 |
| Qdrant Dashboard | 6334 | Qdrant 管理界面 | http://localhost:6334 |

## 系统架构详解

### 文档处理流程

```
用户上传文档
    ↓
Go 后端接收
    ├─ 保存到 MinIO 对象存储
    ├─ 创建数据库记录（状态：parsing）
    └─ 启动异步 goroutine
         ↓
    从 MinIO 下载文件
         ↓
    调用 Python gRPC 服务
         ├─ ParseDocument：解析文档内容
         └─ VectorizeDocument：向量化并存储
              ↓
         Qdrant 向量数据库
         ↓
    更新数据库状态（ready/failed）
```

### 对话流程

```
用户发送问题
    ↓
Go 后端接收
    ↓
调用 Python gRPC 服务
    ├─ 使用 Embedding 模型生成问题向量
    ├─ 在 Qdrant 中检索相关文档片段
    ├─ 构建 Prompt（问题 + 检索到的上下文）
    └─ 调用 LLM 生成回答
         ↓
    返回回答 + 来源文档
```

## API 文档

### 认证接口

- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/logout` - 用户登出
- `POST /api/v1/auth/refresh` - 刷新 Token
- `GET /api/v1/auth/me` - 获取当前用户信息

### 用户管理

- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/me` - 获取当前用户信息
- `GET /api/v1/users/:id` - 获取用户详情
- `POST /api/v1/users` - 创建用户
- `PUT /api/v1/users/:id` - 更新用户
- `PUT /api/v1/users/me/password` - 修改密码
- `DELETE /api/v1/users/:id` - 删除用户

### 知识库管理

- `GET /api/v1/knowledge-bases` - 获取知识库列表
- `GET /api/v1/knowledge-bases/:id` - 获取知识库详情
- `POST /api/v1/knowledge-bases` - 创建知识库（需要 `name` 和 `english_name` 字段）
- `PUT /api/v1/knowledge-bases/:id` - 更新知识库
- `DELETE /api/v1/knowledge-bases/:id` - 删除知识库

### 文档管理（嵌套路由）

- `GET /api/v1/knowledge-bases/:kb_id/documents` - 获取文档列表
- `GET /api/v1/knowledge-bases/:kb_id/documents/:doc_id` - 获取文档详情
- `POST /api/v1/knowledge-bases/:kb_id/documents/upload` - 上传文档（自动解析和向量化）
- `POST /api/v1/knowledge-bases/:kb_id/documents/:doc_id/vectorize` - 手动触发向量化
- `DELETE /api/v1/knowledge-bases/:kb_id/documents/:doc_id` - 删除文档

### 对话管理（嵌套路由）

- `GET /api/v1/chat/sessions` - 获取会话列表
- `GET /api/v1/chat/sessions/:id` - 获取会话详情
- `GET /api/v1/chat/sessions/:id/messages` - 获取会话消息列表
- `POST /api/v1/chat/sessions` - 创建会话
- `POST /api/v1/chat/sessions/:id/messages` - 发送消息（非流式）
- `GET /api/v1/chat/sessions/:id/messages/stream` - 发送消息（流式，SSE）
- `PUT /api/v1/chat/sessions/:id` - 更新会话
- `DELETE /api/v1/chat/sessions/:id` - 删除会话

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

## 使用指南

### 1. 创建知识库

1. 登录系统后，进入"知识库管理"页面
2. 点击"创建知识库"按钮
3. 填写知识库信息：
   - **名称**：知识库的中文名称（如：技术文档库）
   - **英文标识**：唯一的英文标识符（如：tech_docs），只能包含小写字母、数字和下划线
   - **描述**：知识库的描述信息（可选）
   - **嵌入模型**：选择向量化模型（默认 text-embedding-3-small）
   - **分块大小**：文档分块的大小（默认 500）
   - **分块重叠**：分块之间的重叠大小（默认 50）
4. 点击"确定"创建

### 2. 上传文档

1. 在知识库列表中，点击"查看文档"进入文档管理页面
2. 点击"上传文档"按钮
3. 选择要上传的文档（支持多文件上传）
4. 系统会**自动异步处理**：
   - 上传文件到 MinIO 对象存储
   - 调用 Python gRPC 服务解析文档内容
   - 调用 Python gRPC 服务进行向量化
   - 存储向量到 Qdrant 数据库
5. 等待文档状态变为"ready"（已完成）或"failed"（失败）

**文档状态说明**：
- `parsing` - 正在解析和向量化
- `ready` - 已完成，可用于检索
- `failed` - 处理失败，查看错误信息

### 3. 智能对话

1. 进入"智能对话"页面
2. 点击"新对话"创建一个对话会话
3. 选择要检索的知识库（可多选）
4. 在输入框中输入问题
5. 系统会基于选中的知识库进行检索并生成回答
6. 回答下方会显示参考来源文档

### 4. 用户管理（管理员）

1. 管理员用户可以进入"用户管理"页面
2. 可以创建、编辑、删除用户
3. 可以管理用户的角色和权限

## 数据库迁移

数据库表会在首次启动时自动创建（通过 `migrations/*.sql` 文件）。

### 默认管理员账号

系统会自动创建一个默认的超级管理员账号：

- **用户名**: `admin`
- **密码**: `123456`
- **邮箱**: `admin@ragljx.com`

⚠️ **重要提示**: 请在生产环境中立即修改默认密码！

### 工具脚本

#### 生成密码哈希

用于生成 bcrypt 密码哈希，可用于直接在数据库中更新密码：

```bash
cd ragljx_go
go run scripts/gen_password/main.go your_new_password
```

输出示例：
```
Password: your_new_password
Bcrypt Hash: $2a$10$xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
Verification: OK
```

#### 测试配置加载

用于测试配置文件是否正确加载：

```bash
cd ragljx_go
go run scripts/test_config/main.go
```

这将显示所有配置项的值，帮助您验证配置是否正确。

## 开发注意事项

### 1. OpenAI API 配置

**必须配置有效的 OpenAI API Key 才能使用 AI 功能**

- 通过环境变量设置：`export OPENAI_API_KEY=your_key_here`
- 或在 Python 服务配置文件中设置
- 支持自定义 API Base URL（如使用代理或第三方服务）

### 2. 模型选择

**嵌入模型**（用于向量化）：
- `text-embedding-3-small` - 推荐，性价比高
- `text-embedding-3-large` - 更高精度
- `text-embedding-ada-002` - 旧版模型

**对话模型**（用于生成回答）：
- `gpt-4` - 默认，质量最高
- `gpt-3.5-turbo` - 更快，成本更低
- 可在 Python 服务配置文件中修改

### 3. 文件限制

- **前端上传限制**：50MB
- **gRPC 消息大小限制**：100MB
- **支持的文档格式**：
  - 文本文件：TXT, MD
  - Office 文档：DOCX, XLSX, PPTX
  - PDF 文档：PDF
  - 网页文件：HTML
  - 数据文件：CSV, JSON, XML
  - 其他：RTF

### 4. 路由和认证

**RESTful API 设计**：
- 使用嵌套路由：`/knowledge-bases/:id/documents`
- 所有 API 需要 JWT 认证（除登录/注册）
- Token 存储在 `localStorage`

**JWT Token**：
- Access Token 有效期：24 小时
- Refresh Token 有效期：7 天
- 自动刷新机制

### 5. IOC 容器和路由注册

**初始化顺序**（按 Priority 值）：
1. Config 对象（priority 99）- 配置加载
2. HTTP Server（priority 800）- 创建 Gin 引擎，注册全局中间件
3. API 对象（priority -99）- 在 `Init()` 中调用 `Registry()` 注册路由

**中间件注册**：
- CORS、日志、恢复等中间件在 HTTP Server 的 `Init()` 中注册
- 必须在路由注册之前完成

### 6. 文档处理机制

**异步处理**：
- 上传后立即返回，后台异步处理
- 使用 goroutine 调用 Python gRPC 服务
- 实时更新数据库状态

**状态流转**：
```
pending → parsing → ready
                 ↘ failed
```

**错误处理**：
- 解析失败或向量化失败会更新状态为 `failed`
- 错误信息存储在数据库中
- 用户可以查看详细错误信息

## 故障排查

### 前端无法访问

**问题**: 访问 http://localhost:5173 无响应

**解决方案**:
```bash
cd ragljx_web
npm install
npm run dev
```

### Go 服务无法连接数据库

**问题**: Go 服务启动失败，提示数据库连接错误

**解决方案**:
- 检查 PostgreSQL 是否启动：`docker-compose ps postgres`
- 检查数据库配置是否正确（`ragljx_go/config/application.yaml`）
- 查看日志：`docker-compose logs postgres`
- 确保数据库端口 5432 未被占用

### Python 服务无法连接 Qdrant

**问题**: Python 服务启动失败，提示 Qdrant 连接错误

**解决方案**:
- 检查 Qdrant 是否启动：`docker-compose ps qdrant`
- 访问 Qdrant Dashboard: http://localhost:6334
- 查看日志：`docker-compose logs qdrant`

### gRPC 连接失败

**问题**: Go 服务提示无法连接 Python gRPC 服务

**解决方案**:
- 确保 Python 服务已启动：`docker-compose ps ragljx_py`
- 检查端口 50051 是否被占用：`lsof -i :50051`
- 查看 Python 服务日志：`docker-compose logs ragljx_py`

### OpenAI API 调用失败

**问题**: 对话功能无法使用，提示 API 错误

**解决方案**:
- 检查 OpenAI API Key 是否正确配置
- 检查 API Key 是否有效且有余额
- 检查网络是否能访问 OpenAI API
- 查看 Python 服务日志：`docker-compose logs ragljx_py`

### 文档上传失败

**问题**: 文档上传后状态一直是"parsing"或变成"failed"

**解决方案**:
1. **检查文件格式**：确保文件格式在支持列表中
2. **检查文件大小**：不超过 50MB
3. **查看 Go 服务日志**：
   ```bash
   docker-compose logs -f ragljx_go
   ```
4. **查看 Python 服务日志**：
   ```bash
   docker-compose logs -f ragljx_py
   ```
5. **检查 gRPC 连接**：
   - 确保 Python 服务在 50051 端口运行
   - 检查网络连接：`telnet localhost 50051`
6. **检查 MinIO**：
   ```bash
   docker-compose ps minio
   docker-compose logs minio
   ```
7. **检查 Qdrant**：
   ```bash
   docker-compose ps qdrant
   # 访问 Qdrant Dashboard
   open http://localhost:6334
   ```
8. **检查 OpenAI API**：
   - 确保 API Key 有效
   - 确保有足够余额
   - 检查网络能否访问 OpenAI API

### 前端登录后立即跳转到登录页

**问题**: 登录成功后又跳转回登录页

**解决方案**:
- 检查浏览器控制台是否有错误
- 检查后端 API 是否正常返回 Token
- 清除浏览器 localStorage：`localStorage.clear()`
- 检查后端 CORS 配置

## 停止服务

```bash
# 停止所有服务
docker-compose down

# 停止并删除数据卷（谨慎使用，会删除所有数据）
docker-compose down -v

# 停止特定服务
docker-compose stop ragljx_go
docker-compose stop ragljx_py
```

## 项目结构

```
ragljx/
├── ragljx_go/              # Go 后端服务
│   ├── cmd/                # 命令行入口
│   ├── config/             # 配置文件
│   ├── internal/           # 内部代码
│   │   ├── api/           # API 控制器
│   │   ├── middleware/    # 中间件
│   │   ├── model/         # 数据模型
│   │   ├── repository/    # 数据访问层
│   │   ├── service/       # 业务逻辑层
│   │   └── pkg/           # 工具包
│   ├── ioc/               # IOC 容器
│   ├── migrations/        # 数据库迁移
│   └── proto/             # Proto 文件
├── ragljx_py/             # Python AI 服务
│   ├── app/               # 应用代码
│   │   ├── grpc_server/  # gRPC 服务器
│   │   ├── services/     # 业务服务
│   │   ├── utils/        # 工具函数
│   │   └── proto/        # Proto 文件
│   ├── config.yaml        # 配置文件
│   ├── main.py           # 入口文件
│   └── requirements.txt   # Python 依赖
├── ragljx_web/            # Vue 前端
│   ├── src/
│   │   ├── api/          # API 封装
│   │   ├── components/   # 组件
│   │   ├── layouts/      # 布局
│   │   ├── router/       # 路由
│   │   ├── stores/       # 状态管理
│   │   ├── utils/        # 工具函数
│   │   ├── views/        # 页面
│   │   ├── App.vue       # 根组件
│   │   └── main.js       # 入口文件
│   ├── public/           # 静态资源
│   └── package.json      # 依赖配置
├── docker-compose.yml     # Docker 编排
└── README.md             # 项目文档
```

## 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

MIT License

## 联系方式

如有问题或建议，请提交 Issue。

