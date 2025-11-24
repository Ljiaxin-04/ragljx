.PHONY: help build start stop restart logs clean init-db test-config gen-proto

# 默认目标
help:
	@echo "RAG 系统 - 可用命令:"
	@echo ""
	@echo "  make build        - 构建所有服务"
	@echo "  make start        - 启动所有服务"
	@echo "  make stop         - 停止所有服务"
	@echo "  make restart      - 重启所有服务"
	@echo "  make logs         - 查看所有服务日志"
	@echo "  make logs-go      - 查看 Go 服务日志"
	@echo "  make logs-py      - 查看 Python 服务日志"
	@echo "  make clean        - 清理所有容器和数据卷"
	@echo "  make init-db      - 初始化数据库"
	@echo "  make test-config  - 测试配置加载"
	@echo "  make gen-proto    - 生成 Proto 文件"
	@echo "  make gen-password - 生成密码哈希 (使用: make gen-password PWD=yourpassword)"
	@echo ""

# 构建所有服务
build:
	@echo "构建 Go 服务..."
	cd ragljx_go && go build -o bin/ragljx cmd/server/main.go
	@echo "构建 Docker 镜像..."
	docker-compose build

# 启动所有服务
start:
	@echo "启动所有服务..."
	docker-compose up -d
	@echo "等待服务启动..."
	sleep 5
	@echo "服务状态:"
	docker-compose ps

# 停止所有服务
stop:
	@echo "停止所有服务..."
	docker-compose down

# 重启所有服务
restart: stop start

# 查看所有服务日志
logs:
	docker-compose logs -f

# 查看 Go 服务日志
logs-go:
	docker-compose logs -f ragljx_go

# 查看 Python 服务日志
logs-py:
	docker-compose logs -f ragljx_py

# 清理所有容器和数据卷
clean:
	@echo "警告: 这将删除所有容器和数据卷！"
	@read -p "确认继续? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker-compose down -v; \
		echo "清理完成"; \
	else \
		echo "已取消"; \
	fi

# 初始化数据库
init-db:
	@echo "初始化数据库..."
	cd ragljx_go && ./scripts/init_db.sh

# 测试配置加载
test-config:
	@echo "测试配置加载..."
	cd ragljx_go && go run scripts/test_config.go

# 生成 Proto 文件
gen-proto:
	@echo "生成 Go Proto 文件..."
	cd ragljx_go && protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/rag/rag_service.proto
	@echo "生成 Python Proto 文件..."
	cd ragljx_py && python3 -m grpc_tools.protoc -I. \
		--python_out=. --grpc_python_out=. \
		app/proto/rag_service.proto
	@echo "Proto 文件生成完成"

# 生成密码哈希
gen-password:
	@if [ -z "$(PWD)" ]; then \
		echo "用法: make gen-password PWD=yourpassword"; \
		exit 1; \
	fi
	cd ragljx_go && go run scripts/gen_password.go $(PWD)

# 本地开发 - 启动 Go 服务
dev-go:
	cd ragljx_go && go run cmd/server/main.go

# 本地开发 - 启动 Python 服务
dev-py:
	cd ragljx_py && python main.py

# 安装 Go 依赖
deps-go:
	cd ragljx_go && go mod tidy

# 安装 Python 依赖
deps-py:
	cd ragljx_py && pip install -r requirements.txt

# 运行 Go 测试
test-go:
	cd ragljx_go && go test ./...

# 格式化 Go 代码
fmt-go:
	cd ragljx_go && go fmt ./...

# 检查 Go 代码
lint-go:
	cd ragljx_go && golangci-lint run

# 查看服务状态
status:
	docker-compose ps

# 进入 Go 容器
shell-go:
	docker-compose exec ragljx_go sh

# 进入 Python 容器
shell-py:
	docker-compose exec ragljx_py bash

# 进入 PostgreSQL 容器
shell-db:
	docker-compose exec postgres psql -U ragljx -d ragljx

# 备份数据库
backup-db:
	@echo "备份数据库..."
	docker-compose exec -T postgres pg_dump -U ragljx ragljx > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "备份完成"

# 恢复数据库
restore-db:
	@if [ -z "$(FILE)" ]; then \
		echo "用法: make restore-db FILE=backup.sql"; \
		exit 1; \
	fi
	@echo "恢复数据库..."
	docker-compose exec -T postgres psql -U ragljx -d ragljx < $(FILE)
	@echo "恢复完成"

