# API Key 配置指南

本系统使用两个独立的 AI 模型：
1. **嵌入模型（Embedding Model）**：用于文档向量化
2. **对话模型（Chat Model）**：用于 AI 对话

您可以为这两个模型配置不同的 API Key 和服务提供商。

## 📝 配置步骤

### 1. 编辑 `.env` 文件

在项目根目录 `/Users/liang/projectljx/ragljx/.env` 中配置：

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

# ========================================
# 模型参数配置
# ========================================
CHAT_TEMPERATURE=0.7
CHAT_MAX_TOKENS=2000
```

### 2. 替换 API Key

将 `your_embedding_api_key_here` 和 `your_chat_api_key_here` 替换为您的真实 API Key。

**如果使用同一个 OpenAI 账号**，两个 Key 可以相同：
```bash
EMBEDDING_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
CHAT_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

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

## 🔧 测试配置

配置完成后，可以通过以下方式测试：

### 1. 测试 Python 服务启动

```bash
cd /Users/liang/projectljx/ragljx/ragljx_py

# 设置环境变量
export EMBEDDING_API_KEY="your_key"
export CHAT_API_KEY="your_key"

# 启动服务
python main.py
```

查看日志，应该看到：
```
VectorService initialized with Qdrant at localhost:6333
Using embedding model: text-embedding-3-small
ChatService initialized with model: gpt-4
```

### 2. 使用 Docker Compose

```bash
cd /Users/liang/projectljx/ragljx

# 启动服务
docker-compose up -d

# 查看 Python 服务日志
docker-compose logs -f ragljx_py
```

## ⚠️ 常见问题

### Q1: 两个 API Key 必须不同吗？

**A**: 不必须。如果使用同一个 OpenAI 账号，两个 Key 可以相同。

### Q2: 可以只配置一个 Key 吗？

**A**: 不可以。系统需要两个模型都正常工作。但两个 Key 可以相同。

### Q3: 嵌入模型可以用国内大模型吗？

**A**: 理论上可以，但需要确保：
1. 该服务提供嵌入 API
2. API 格式兼容 OpenAI
3. 向量维度一致（默认 1536）

建议嵌入模型使用 OpenAI，对话模型可以灵活选择。

### Q4: 如何降低成本？

**A**: 
1. 对话模型改用 gpt-3.5-turbo 或国内大模型
2. 减少 CHAT_MAX_TOKENS（默认 2000）
3. 提高 CHAT_TEMPERATURE 可能减少 token 使用

### Q5: API Key 安全吗？

**A**: 
- `.env` 文件已在 `.gitignore` 中，不会提交到 Git
- Docker 容器中的环境变量是隔离的
- 建议定期轮换 API Key
- 在 OpenAI 控制台设置使用限额

## 📚 相关文档

- [OpenAI API 文档](https://platform.openai.com/docs/api-reference)
- [DeepSeek API 文档](https://platform.deepseek.com/api-docs/)
- [通义千问 API 文档](https://help.aliyun.com/zh/dashscope/)
- [智谱 AI API 文档](https://open.bigmodel.cn/dev/api)

## 🆘 获取帮助

如果配置过程中遇到问题：

1. 检查 API Key 是否正确
2. 检查 API Base URL 是否正确
3. 查看服务日志：`docker-compose logs ragljx_py`
4. 测试 API 连接：`curl -H "Authorization: Bearer YOUR_KEY" https://api.openai.com/v1/models`

