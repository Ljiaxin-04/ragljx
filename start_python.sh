#!/bin/bash

# 启动 Python AI 服务

cd "$(dirname "$0")/ragljx_py"

echo "🚀 启动 Python AI 服务..."

# 设置环境变量
export RAGLJX_DB_HOST=localhost
export RAGLJX_DB_PORT=5432
export RAGLJX_DB_DATABASE=ragljx
export RAGLJX_DB_USERNAME=ragljx
export RAGLJX_DB_PASSWORD=ragljx_password
export RAGLJX_QDRANT_HOST=localhost
export RAGLJX_QDRANT_PORT=6333
export RAGLJX_MINIO_ENDPOINT=localhost:9000
export RAGLJX_MINIO_ACCESS_KEY=minioadmin
export RAGLJX_MINIO_SECRET_KEY=minioadmin
export RAGLJX_GRPC_HOST=0.0.0.0
export RAGLJX_GRPC_PORT=50051

# AI 配置（从 .env 文件读取，并导出为环境变量）
if [ -f "../.env" ]; then
    # 自动 export .env 中定义的变量
    set -a
    source "../.env"
    set +a
fi

echo "📋 配置信息："
echo "  - 数据库: ${RAGLJX_DB_HOST}:${RAGLJX_DB_PORT}"
echo "  - Qdrant: ${RAGLJX_QDRANT_HOST}:${RAGLJX_QDRANT_PORT}"
echo "  - MinIO: ${RAGLJX_MINIO_ENDPOINT}"
echo "  - gRPC: ${RAGLJX_GRPC_HOST}:${RAGLJX_GRPC_PORT}"
echo "  - 嵌入模型: ${EMBEDDING_MODEL}"
echo "  - 对话模型: ${CHAT_MODEL}"
echo ""

# 查找 Python3
PYTHON3=$(which python3)
if [ -z "$PYTHON3" ]; then
    echo "❌ 错误: 未找到 python3"
    exit 1
fi

echo "📍 使用 Python: $PYTHON3"
$PYTHON3 --version

# 检查虚拟环境
if [ ! -d "venv" ] || [ ! -f "venv/bin/activate" ]; then
    echo "📦 创建虚拟环境..."
    rm -rf venv
    $PYTHON3 -m venv venv
    if [ ! -f "venv/bin/activate" ]; then
        echo "❌ 错误: 虚拟环境创建失败"
        exit 1
    fi
fi

# 激活虚拟环境
echo "🔧 激活虚拟环境..."
source venv/bin/activate

# 验证激活成功
if [ -z "$VIRTUAL_ENV" ]; then
    echo "❌ 错误: 虚拟环境激活失败"
    exit 1
fi

echo "✅ 虚拟环境已激活: $VIRTUAL_ENV"

# 安装依赖
if [ ! -f "venv/.installed" ]; then
    echo "📦 安装依赖..."

    # 升级 pip
    python -m pip install --upgrade pip \
        -i https://pypi.tuna.tsinghua.edu.cn/simple \
        --trusted-host pypi.tuna.tsinghua.edu.cn

    # 使用国内镜像源安装依赖（避免 SSL 证书问题）
    pip install -r requirements.txt \
        -i https://pypi.tuna.tsinghua.edu.cn/simple \
        --trusted-host pypi.tuna.tsinghua.edu.cn

    if [ $? -eq 0 ]; then
        touch venv/.installed
        echo "✅ 依赖安装完成"
    else
        echo "❌ 依赖安装失败"
        exit 1
    fi
fi

# 启动服务
echo "✅ 启动服务..."
python main.py

