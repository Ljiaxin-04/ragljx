#!/bin/bash

# 启动 Go 后端服务

cd "$(dirname "$0")/ragljx_go"

echo "🚀 启动 Go 后端服务..."

# 设置环境变量
export RAGLJX_DB_HOST=localhost
export RAGLJX_DB_PORT=5432
export RAGLJX_DB_DATABASE=ragljx
export RAGLJX_DB_USERNAME=ragljx
export RAGLJX_DB_PASSWORD=ragljx_password
export RAGLJX_REDIS_HOST=localhost
export RAGLJX_REDIS_PORT=6379
export RAGLJX_KAFKA_BROKERS=localhost:19092
export RAGLJX_MINIO_ENDPOINT=localhost:9000
export RAGLJX_MINIO_ACCESS_KEY=minioadmin
export RAGLJX_MINIO_SECRET_KEY=minioadmin
export RAGLJX_GRPC_ADDRESS=localhost:50051
export RAGLJX_HTTP_PORT=8080

echo "📋 配置信息："
echo "  - 数据库: ${RAGLJX_DB_HOST}:${RAGLJX_DB_PORT}"
echo "  - Redis: ${RAGLJX_REDIS_HOST}:${RAGLJX_REDIS_PORT}"
echo "  - Kafka: ${RAGLJX_KAFKA_BROKERS}"
echo "  - MinIO: ${RAGLJX_MINIO_ENDPOINT}"
echo "  - gRPC: ${RAGLJX_GRPC_ADDRESS}"
echo "  - HTTP: 0.0.0.0:${RAGLJX_HTTP_PORT}"
echo ""

# 检查依赖
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到 go"
    exit 1
fi

# 启动服务
echo "✅ 启动服务..."
go run cmd/server/main.go

