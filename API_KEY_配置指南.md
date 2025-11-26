<div align="center">

# 🔑 API Key 配置指南

**为 RAG 知识库系统配置 AI 模型**

</div>

---

## 📋 目录

- [概述](#-概述)
- [方式一：Docker 全容器启动配置](#-方式一docker-全容器启动配置)
- [方式二：本地开发启动配置](#-方式二本地开发启动配置)
- [支持的服务提供商](#-支持的服务提供商)
- [常见问题](#-常见问题)

---

## 📖 概述

本系统使用两个独立的 AI 模型：

| 模型类型 | 用途 | 推荐服务 |
|---------|------|---------|
| 🔢 **嵌入模型** | 文档向量化 | OpenAI text-embedding-3-small |
| 💬 **对话模型** | 智能问答 | DeepSeek / OpenAI GPT-4 |

> 💡 **提示**：两个模型可以使用不同的服务提供商，灵活搭配以优化成本和效果。

---

## 🐳 方式一：Docker 全容器启动配置

> 使用 `docker-compose.yml` 启动所有服务时的配置方法

### 步骤 1：创建 .env 文件

```bash
# 进入项目根目录
cd ragljx

# 复制模板文件
cp .env.example .env
```

### 步骤 2：编辑 .env 文件

```bash
# ========================================
# 嵌入模型配置（用于文档向量化）
# ========================================
EMBEDDING_API_KEY=your_embedding_api_key_here
EMBEDDING_API_BASE=https://api.openai.com/v1
EMBEDDING_MODEL=text-embedding-3-small

# ========================================
# 对话模型配置（用于 AI 对话）
# ========================================
CHAT_API_KEY=your_chat_api_key_here
CHAT_API_BASE=https://api.openai.com/v1
CHAT_MODEL=gpt-4
```

### 步骤 3：启动服务

```bash
# 启动所有容器
docker-compose up -d

# 验证配置是否生效
docker-compose logs ragljx_py | grep "initialized"
```

**预期输出**：

```
VectorService initialized with Qdrant at qdrant:6333
Using embedding model: text-embedding-3-small
ChatService initialized with model: gpt-4
```

### 🔧 配置原理

`docker-compose.yml` 中的 Python 服务会读取 `.env` 文件中的环境变量：

```yaml
ragljx_py:
  environment:
    EMBEDDING_API_KEY: ${EMBEDDING_API_KEY:-}
    EMBEDDING_API_BASE: ${EMBEDDING_API_BASE:-https://api.openai.com/v1}
    EMBEDDING_MODEL: ${EMBEDDING_MODEL:-text-embedding-3-small}
    CHAT_API_KEY: ${CHAT_API_KEY:-}
    CHAT_API_BASE: ${CHAT_API_BASE:-https://api.openai.com/v1}
    CHAT_MODEL: ${CHAT_MODEL:-gpt-4}
```

---

## 💻 方式二：本地开发启动配置

> 使用 `docker-compose.infra.yml` + 本地服务时的配置方法

### 步骤 1：创建 .env 文件

```bash
# 进入项目根目录
cd ragljx

# 复制模板文件
cp .env.example .env
```

### 步骤 2：编辑 .env 文件

```bash
# ========================================
# 嵌入模型配置（用于文档向量化）
# ========================================
EMBEDDING_API_KEY=your_embedding_api_key_here
EMBEDDING_API_BASE=https://api.openai.com/v1
EMBEDDING_MODEL=text-embedding-3-small

# ========================================
# 对话模型配置（用于 AI 对话）
# ========================================
CHAT_API_KEY=your_chat_api_key_here
CHAT_API_BASE=https://api.openai.com/v1
CHAT_MODEL=gpt-4
```

### 步骤 3：启动 Python 服务

```bash
# 运行启动脚本（会自动读取 .env 文件）
./start_python.sh
```

**预期输出**：

```
VectorService initialized with Qdrant at localhost:6333
Using embedding model: text-embedding-3-small
ChatService initialized with model: gpt-4
```

### 🔧 配置原理

`start_python.sh` 脚本会自动读取项目根目录的 `.env` 文件：

```bash
# start_python.sh 中的关键代码
if [ -f "../.env" ]; then
    set -a
    source "../.env"
    set +a
    echo "✅ 已加载 .env 配置文件"
fi
```

---

## 🌐 支持的服务提供商

### 方案 1：全部使用 OpenAI（推荐）

```bash
# 嵌入模型
EMBEDDING_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
EMBEDDING_API_BASE=https://api.openai.com/v1
EMBEDDING_MODEL=text-embedding-3-small

# 对话模型
CHAT_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
CHAT_API_BASE=https://api.openai.com/v1
CHAT_MODEL=gpt-4
```

**优点**：质量最高，兼容性最好
**缺点**：需要国际信用卡，国内访问可能需要代理

### 方案 2：嵌入用 OpenAI，对话用国内大模型

```bash
# 嵌入模型（OpenAI）
EMBEDDING_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
EMBEDDING_API_BASE=https://api.openai.com/v1
EMBEDDING_MODEL=text-embedding-3-small

# 对话模型（DeepSeek）
CHAT_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
CHAT_API_BASE=https://api.deepseek.com/v1
CHAT_MODEL=deepseek-chat
```

**优点**：嵌入质量高，对话成本低
**缺点**：需要两个账号

### 方案 3：使用 OpenAI 代理服务

如果您在国内无法直接访问 OpenAI，可以使用代理服务：

```bash
# 使用代理服务
EMBEDDING_API_KEY=your_api_key
EMBEDDING_API_BASE=https://api.openai-proxy.com/v1  # 替换为您的代理地址

CHAT_API_KEY=your_api_key
CHAT_API_BASE=https://api.openai-proxy.com/v1
```

## 🤖 国内大模型配置示例

### DeepSeek（推荐，性价比高）

```bash
CHAT_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
CHAT_API_BASE=https://api.deepseek.com/v1
CHAT_MODEL=deepseek-chat
```

- 官网：https://platform.deepseek.com/
- 价格：约 ¥0.001/1K tokens（输入），¥0.002/1K tokens（输出）
- 特点：兼容 OpenAI API，质量好，价格低

### 通义千问（阿里云）

```bash
CHAT_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
CHAT_API_BASE=https://dashscope.aliyuncs.com/compatible-mode/v1
CHAT_MODEL=qwen-turbo
```

- 官网：https://dashscope.aliyun.com/
- 可选模型：qwen-turbo, qwen-plus, qwen-max
- 特点：国内访问快，稳定性好

### 智谱 AI（GLM）

```bash
CHAT_API_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxxxxxxxx
CHAT_API_BASE=https://open.bigmodel.cn/api/paas/v4
CHAT_MODEL=glm-4
```

- 官网：https://open.bigmodel.cn/
- 可选模型：glm-4, glm-4-flash
- 特点：中文能力强

### 月之暗面（Moonshot）

```bash
CHAT_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
CHAT_API_BASE=https://api.moonshot.cn/v1
CHAT_MODEL=moonshot-v1-8k
```

- 官网：https://platform.moonshot.cn/
- 可选模型：moonshot-v1-8k, moonshot-v1-32k, moonshot-v1-128k
- 特点：长文本处理能力强

## 💰 成本估算

### OpenAI 价格（2024年）

**嵌入模型**：
- text-embedding-3-small: $0.00002/1K tokens
- text-embedding-3-large: $0.00013/1K tokens

**对话模型**：
- gpt-3.5-turbo: 输入 $0.0005/1K，输出 $0.0015/1K
- gpt-4: 输入 $0.03/1K，输出 $0.06/1K
- gpt-4-turbo: 输入 $0.01/1K，输出 $0.03/1K

### 使用示例成本

假设处理 100 个文档（每个 5000 字）+ 100 次对话（每次 500 字）：

**方案 1：OpenAI（gpt-4）**
- 嵌入：100 × 5000 / 1000 × $0.00002 = $0.01
- 对话：100 × 500 / 1000 × ($0.03 + $0.06) = $4.5
- **总计：约 $4.51**

**方案 2：OpenAI 嵌入 + DeepSeek 对话**
- 嵌入：$0.01
- 对话：100 × 500 / 1000 × (¥0.001 + ¥0.002) / 7 = $0.02
- **总计：约 $0.03**

---

## ✅ 验证配置

### 方式一验证（Docker 全容器）

```bash
# 查看 Python 服务日志
docker-compose logs ragljx_py | grep -E "(initialized|model)"
```

### 方式二验证（本地开发）

```bash
# 启动 Python 服务后查看输出
./start_python.sh
```

**成功标志**：

```
✅ VectorService initialized with Qdrant at localhost:6333
✅ Using embedding model: text-embedding-3-small
✅ ChatService initialized with model: deepseek-chat
✅ RAG gRPC server started on 0.0.0.0:50051
```

---

## ❓ 常见问题

<details>
<summary><b>Q1: 两个 API Key 必须不同吗？</b></summary>

不必须。如果使用同一个 OpenAI 账号，两个 Key 可以相同。
</details>

<details>
<summary><b>Q2: 可以只配置一个 Key 吗？</b></summary>

不可以。系统需要嵌入模型和对话模型都正常工作。但两个 Key 可以相同。
</details>

<details>
<summary><b>Q3: 嵌入模型可以用国内大模型吗？</b></summary>

理论上可以，但需要确保：
1. 该服务提供嵌入 API
2. API 格式兼容 OpenAI
3. 向量维度一致（默认 1536）

建议嵌入模型使用 OpenAI，对话模型可以灵活选择。
</details>

<details>
<summary><b>Q4: 如何降低成本？</b></summary>

1. 对话模型改用 DeepSeek（约 OpenAI 的 1/100 价格）
2. 嵌入模型使用 text-embedding-3-small
3. 减少 RAG 检索的 TopK 值
</details>

<details>
<summary><b>Q5: API Key 安全吗？</b></summary>

- `.env` 文件已在 `.gitignore` 中，不会提交到 Git
- Docker 容器中的环境变量是隔离的
- 建议定期轮换 API Key
- 在服务商控制台设置使用限额
</details>

---

## 📚 相关文档

| 服务商 | 文档链接 |
|--------|---------|
| OpenAI | [API 文档](https://platform.openai.com/docs/api-reference) |
| DeepSeek | [API 文档](https://platform.deepseek.com/api-docs/) |
| 通义千问 | [API 文档](https://help.aliyun.com/zh/dashscope/) |
| 智谱 AI | [API 文档](https://open.bigmodel.cn/dev/api) |
| 月之暗面 | [API 文档](https://platform.moonshot.cn/docs) |

---

<div align="center">

**🔑 配置完成后，请返回 [启动指南.md](启动指南.md) 继续启动服务**

</div>
