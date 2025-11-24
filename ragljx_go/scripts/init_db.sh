#!/bin/bash

# 数据库初始化脚本
# 用于在 Docker 容器启动时自动运行所有迁移

set -e

# 数据库连接参数
DB_HOST="${RAGLJX_DB_HOST:-localhost}"
DB_PORT="${RAGLJX_DB_PORT:-5432}"
DB_NAME="${RAGLJX_DB_DATABASE:-ragljx}"
DB_USER="${RAGLJX_DB_USERNAME:-ragljx}"
DB_PASSWORD="${RAGLJX_DB_PASSWORD:-ragljx_password}"

echo "========================================="
echo "RAG 系统数据库初始化"
echo "========================================="
echo "数据库主机: $DB_HOST:$DB_PORT"
echo "数据库名称: $DB_NAME"
echo "用户名: $DB_USER"
echo "========================================="

# 等待数据库就绪
echo "等待数据库就绪..."
until PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c '\q' 2>/dev/null; do
  echo "数据库未就绪，等待中..."
  sleep 2
done

echo "数据库已就绪！"
echo ""

# 执行所有迁移脚本
MIGRATION_DIR="./migrations"

if [ ! -d "$MIGRATION_DIR" ]; then
    echo "错误: 迁移目录不存在: $MIGRATION_DIR"
    exit 1
fi

echo "开始执行数据库迁移..."
echo ""

# 按文件名排序执行所有 .sql 文件
for sql_file in $(ls -1 $MIGRATION_DIR/*.sql | sort); do
    echo "执行迁移: $(basename $sql_file)"
    PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$sql_file"
    
    if [ $? -eq 0 ]; then
        echo "✓ $(basename $sql_file) 执行成功"
    else
        echo "✗ $(basename $sql_file) 执行失败"
        exit 1
    fi
    echo ""
done

echo "========================================="
echo "数据库迁移完成！"
echo "========================================="
echo ""
echo "默认管理员账号:"
echo "  用户名: admin"
echo "  密码: 123456"
echo "  邮箱: admin@ragljx.com"
echo ""
echo "⚠️  请在生产环境中立即修改默认密码！"
echo "========================================="

