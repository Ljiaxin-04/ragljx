<div align="center">

# 🧠 RAG 知识库系统

**企业级智能知识库解决方案**

基于 Go + Python + Vue 构建的 RAG（检索增强生成）知识库系统

让您的文档"活"起来，与知识对话

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![Python Version](https://img.shields.io/badge/Python-3.12+-3776AB?style=flat-square&logo=python)](https://python.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.x-4FC08D?style=flat-square&logo=vue.js)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

[功能特性](#-功能特性) • [快速开始](#-快速开始) • [系统架构](#-系统架构) • [技术栈](#-技术栈) • [文档](#-文档)

</div>

---

## 📖 项目简介

RAG 知识库系统是一个开箱即用的企业级智能文档管理与对话系统。通过先进的检索增强生成（RAG）技术，将您的私有文档转化为可交互的知识库，实现精准的语义检索和智能问答。

### 🎯 适用场景

- 📚 **企业知识库** - 构建内部知识管理系统，让员工快速获取信息
- 📋 **文档智能问答** - 上传产品手册、技术文档，实现智能客服
- 🎓 **教育培训** - 构建课程资料库，辅助学习和培训
- 🔬 **研究助手** - 管理论文、报告，快速检索相关内容
- 💼 **合同/法务管理** - 智能检索合同条款和法律文件

---

## ✨ 功能特性

<table>
<tr>
<td width="50%">

### 🔐 用户管理
- ✅ 用户注册、登录、JWT 认证
- ✅ 基于角色的权限控制（RBAC）
- ✅ 个人信息管理、密码修改
- ✅ 管理员用户管理功能

</td>
<td width="50%">

### 📚 知识库管理
- ✅ 创建、编辑、删除知识库
- ✅ 自定义嵌入模型
- ✅ 知识库统计信息展示
- ✅ 多知识库隔离管理

</td>
</tr>
<tr>
<td width="50%">

### 📄 文档管理
- ✅ 支持 12+ 种文档格式
- ✅ 自动文档解析与向量化
- ✅ 文档状态实时跟踪
- ✅ 批量上传、断点续传

</td>
<td width="50%">

### 💬 智能对话
- ✅ 基于 RAG 的精准问答
- ✅ 流式输出，实时响应
- ✅ 多知识库联合检索
- ✅ 来源文档追溯引用

</td>
</tr>
</table>

### 📁 支持的文档格式

| 类型 | 格式 |
|------|------|
| 📝 文本文件 | TXT, MD, RTF |
| 📊 Office 文档 | DOCX, XLSX, PPTX |
| 📕 PDF 文档 | PDF |
| 🌐 网页文件 | HTML |
| 📋 数据文件 | CSV, JSON, XML |

---

## 🏗 系统架构

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                                🖥 用户浏览器                                   │
│                         Vue 3  +  Element Plus                               │
└──────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼  HTTP / SSE
┌──────────────────────────────────────────────────────────────────────────────┐
│                           🚀 Go 后端服务 (8080)                                │
│                   Gin  ·  GORM  ·  JWT  ·  IOC 容器                           │
│               用户管理  │  知识库管理  │  文档管理  │  会话管理                    │
└──────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼  gRPC
┌──────────────────────────────────────────────────────────────────────────────┐
│                          🧠 Python AI 服务 (50051)                            │
│                    LlamaIndex  ·  OpenAI API  ·  gRPC                        │
│                 文档解析  │  向量化  │  语义检索  │  RAG 对话                    │
└──────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                                🔧 基础设施层                                   │
│                                                                              │
│   ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐   │
│   │PostgreSQL│   │  Qdrant  │   │  MinIO   │   │  Redis   │   │  Kafka   │   │
│   │ 业务数据  │   │ 向量存储   │   │ 文件存储  │   │   缓存    │   │ 消息队列  │   │
│   └──────────┘   └──────────┘   └──────────┘   └──────────┘   └──────────┘   │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## 🛠 技术栈

<table>
<tr>
<th align="center">层级</th>
<th align="center">技术</th>
<th align="center">说明</th>
</tr>
<tr>
<td><b>🖥 前端</b></td>
<td>Vue 3 + Vite + Element Plus + Pinia</td>
<td>现代化响应式用户界面</td>
</tr>
<tr>
<td><b>⚙️ 后端</b></td>
<td>Go + Gin + GORM + gRPC</td>
<td>高性能 RESTful API 服务</td>
</tr>
<tr>
<td><b>🤖 AI 服务</b></td>
<td>Python + LlamaIndex + gRPC</td>
<td>文档处理与 RAG 引擎</td>
</tr>
<tr>
<td><b>🗄 数据库</b></td>
<td>PostgreSQL</td>
<td>业务数据持久化存储</td>
</tr>
<tr>
<td><b>🔍 向量库</b></td>
<td>Qdrant</td>
<td>高性能向量相似度检索</td>
</tr>
<tr>
<td><b>📦 对象存储</b></td>
<td>MinIO</td>
<td>文档文件存储</td>
</tr>
<tr>
<td><b>💾 缓存</b></td>
<td>Redis</td>
<td>会话缓存（预留）</td>
</tr>
<tr>
<td><b>📨 消息队列</b></td>
<td>Kafka</td>
<td>异步任务处理（预留）</td>
</tr>
</table>

### 🤖 支持的 AI 模型

**嵌入模型（文档向量化）**
- OpenAI text-embedding-3-small / large
- 兼容 OpenAI API 的第三方服务

**对话模型（智能问答）**
- OpenAI GPT-4 / GPT-3.5
- DeepSeek (推荐，性价比高)
- 通义千问、智谱 GLM、月之暗面等

---

## 🚀 快速开始

### 📋 前置要求

- **Docker & Docker Compose** v2.0+
- **OpenAI 兼容的 API Key**（必需）

### ⚡ 一键启动

```bash
# 1. 克隆项目
git clone https://github.com/your-username/ragljx.git
cd ragljx

# 2. 配置 API Key（必需）
cp .env.example .env
# 编辑 .env 文件，填入您的 API Key

# 3. 启动所有服务
docker-compose up -d

# 4. 启动前端开发服务器
cd ragljx_web && npm install && npm run dev
```

### 🌐 访问系统

| 服务 | 地址 | 说明 |
|------|------|------|
| 🖥 前端界面 | http://localhost:5173 | 用户操作界面 |
| 📡 后端 API | http://localhost:8080 | RESTful API |
| 📦 MinIO 控制台 | http://localhost:9001 | 文件存储管理 |
| 🔍 Qdrant 面板 | http://localhost:6333/dashboard | 向量数据库管理 |

### 🔑 默认账号

```
用户名: admin
密码: 123456
```

⚠️ **首次登录后请立即修改密码！**

> 📖 详细的启动说明请参阅 [启动指南.md](启动指南.md)
>
> 🔑 API Key 配置请参阅 [API_KEY_配置指南.md](API_KEY_配置指南.md)

---

## 📖 使用指南

### 1️⃣ 创建知识库

1. 登录系统，进入「知识库管理」
2. 点击「创建知识库」
3. 填写名称、英文标识（用于向量存储）、描述
4. 选择嵌入模型，点击确定

### 2️⃣ 上传文档

1. 进入知识库，点击「上传文档」
2. 选择文档文件（支持批量上传）
3. 系统自动解析并向量化
4. 等待状态变为「已完成」

### 3️⃣ 智能对话

1. 进入「智能对话」页面
2. 创建新对话，选择关联的知识库
3. 输入问题，获取基于文档的智能回答
4. 查看回答来源，验证信息准确性

---

## 📊 核心流程

### 📄 文档处理流程

```
📤 用户上传文档
       ↓
⚙️ Go 后端接收
   ├── 计算文件 SHA256
   ├── 存储至 MinIO
   └── 创建数据库记录
       ↓
🔄 异步处理 (goroutine)
   ├── 从 MinIO 读取文件
   └── 调用 Python gRPC
       ↓
📖 文档解析 (ParseDocument)
   └── 提取纯文本内容
       ↓
🔢 向量化 (VectorizeDocument)
   ├── 文本分块 (chunk_size=512)
   ├── 调用 Embedding API
   └── 存入 Qdrant
       ↓
✅ 更新状态为 completed
```

### 💬 RAG 对话流程

```
❓ 用户提问
       ↓
🔐 JWT 认证
       ↓
🔍 语义检索
   ├── 问题向量化
   ├── Qdrant 相似度搜索
   └── 获取相关文档片段
       ↓
📝 构建 Prompt
   ├── 系统提示词
   ├── RAG 上下文
   └── 历史对话
       ↓
🤖 LLM 生成回答
       ↓
📚 返回答案 + 来源引用
```

---

## 📡 API 接口

### 🔐 认证相关
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/login` | 用户登录 |
| POST | `/api/v1/auth/register` | 用户注册 |
| POST | `/api/v1/auth/refresh` | 刷新 Token |

### 📚 知识库管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/knowledge-bases` | 获取知识库列表 |
| POST | `/api/v1/knowledge-bases` | 创建知识库 |
| PUT | `/api/v1/knowledge-bases/:id` | 更新知识库 |
| DELETE | `/api/v1/knowledge-bases/:id` | 删除知识库 |

### 📄 文档管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/knowledge-bases/:id/documents` | 获取文档列表 |
| POST | `/api/v1/knowledge-bases/:id/documents/upload` | 上传文档 |
| DELETE | `/api/v1/knowledge-bases/:id/documents/:docId` | 删除文档 |

### 💬 对话管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/chat/sessions` | 获取会话列表 |
| POST | `/api/v1/chat/sessions` | 创建会话 |
| GET | `/api/v1/chat/sessions/:id/messages/stream` | 流式对话 (SSE) |

---

## 📁 项目结构

```
ragljx/
├── 📂 ragljx_go/              # Go 后端服务
│   ├── cmd/server/           # 应用入口
│   ├── config/               # 配置文件
│   ├── internal/
│   │   ├── api/             # API 控制器
│   │   ├── middleware/      # 中间件
│   │   ├── model/           # 数据模型
│   │   ├── repository/      # 数据访问层
│   │   └── service/         # 业务逻辑层
│   ├── ioc/                  # IOC 容器
│   └── migrations/           # 数据库迁移
│
├── 📂 ragljx_py/              # Python AI 服务
│   ├── app/
│   │   ├── grpc_server/     # gRPC 服务
│   │   ├── services/        # 业务服务
│   │   └── proto/           # Proto 定义
│   ├── config.yaml           # 配置文件
│   └── main.py               # 入口文件
│
├── 📂 ragljx_web/             # Vue 前端
│   ├── src/
│   │   ├── api/             # API 封装
│   │   ├── components/      # 组件
│   │   ├── views/           # 页面
│   │   ├── router/          # 路由
│   │   └── stores/          # 状态管理
│   └── package.json
│
├── 📄 docker-compose.yml      # 完整部署配置
├── 📄 docker-compose.infra.yml # 基础设施配置
├── 📄 start_go.sh             # Go 服务启动脚本
├── 📄 start_python.sh         # Python 服务启动脚本
└── 📄 .env                    # 环境变量配置
```

---

## ❓ 常见问题

<details>
<summary><b>Q: 文档上传后一直显示"处理中"？</b></summary>

检查以下几点：
1. Python AI 服务是否正常运行
2. API Key 是否配置正确
3. 查看 Python 服务日志：`docker-compose logs ragljx_py`
</details>

<details>
<summary><b>Q: 对话无法获取回答？</b></summary>

1. 确认知识库中有已处理完成的文档
2. 检查对话模型 API Key 是否有效
3. 检查网络是否能访问 AI 服务
</details>

<details>
<summary><b>Q: 如何降低使用成本？</b></summary>

1. 对话模型使用 DeepSeek（约 OpenAI 的 1/10 价格）
2. 嵌入模型使用 text-embedding-3-small
3. 调整 RAG 参数减少 token 使用
</details>

<details>
<summary><b>Q: 支持私有化部署吗？</b></summary>

完全支持！系统设计为本地部署，所有数据存储在您自己的服务器上。
</details>

---

## 🛑 停止服务

```bash
# 停止所有 Docker 服务
docker-compose down

# 停止并删除数据（谨慎使用）
docker-compose down -v
```

---

## 🤝 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

---

## 📄 许可证

本项目采用 [MIT License](LICENSE) 开源许可证。

---


<div align="center">

**⭐ 如果这个项目对你有帮助，请给一个 Star！**

Made with ❤️ by RAG Knowledge Base Team

</div>