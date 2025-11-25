# 工具脚本

本目录包含用于开发和维护的工具脚本。

## 可用脚本

### 1. gen_password - 密码哈希生成器

生成 bcrypt 密码哈希，用于创建或更新用户密码。

**使用方法**:
```bash
go run gen_password/main.go <password>
```

**示例**:
```bash
go run gen_password/main.go mypassword123
```

**输出**:
```
Password: mypassword123
Bcrypt Hash: $2a$10$xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
Verification: OK
```

**用途**:
- 生成新用户的密码哈希
- 重置用户密码（直接在数据库中更新）
- 验证密码哈希是否正确

### 2. test_config - 配置测试工具

测试配置文件是否正确加载，显示所有配置项的值。

**使用方法**:
```bash
go run test_config/main.go
```

**输出示例**:
```
=== Configuration Test ===

PostgreSQL Config:
  Host: localhost
  Port: 5432
  Database: ragljx
  Username: ragljx
  Debug: true

Redis Config:
  Endpoints: [localhost:6379]
  DB: 0

Kafka Config:
  Brokers: [localhost:19092]

MinIO Config:
  Endpoint: localhost:9000
  AccessKeyID: minioadmin
  BucketName: ragljx
  UseSSL: false

gRPC Config:
  PythonAddr: localhost:50051

HTTP Config:
  Host: 0.0.0.0
  Port: 8080
  ReadTimeout: 60
  WriteTimeout: 60

Log Config:
  Level: debug
  Output: stdout
  Format: json

=== All configurations loaded successfully! ===
```

**用途**:
- 验证配置文件格式是否正确
- 检查配置项是否被正确读取
- 调试配置加载问题

### 3. init_db.sh - 数据库初始化脚本

初始化数据库，执行迁移脚本。

**使用方法**:
```bash
./init_db.sh
```

**用途**:
- 创建数据库表结构
- 执行数据库迁移
- 初始化默认数据

## 注意事项

1. 所有 Go 脚本都应该从 `ragljx_go` 目录运行
2. 确保配置文件 `config/application.yaml` 存在
3. 密码哈希生成器使用 bcrypt 算法，与系统保持一致
4. 配置测试工具不会初始化实际的数据库连接，只测试配置加载

## 添加新脚本

如果需要添加新的工具脚本：

1. 在 `scripts` 目录下创建新的子目录
2. 在子目录中创建 `main.go` 文件
3. 使用 `package main` 和 `func main()`
4. 更新本 README 文件

**示例结构**:
```
scripts/
├── gen_password/
│   └── main.go
├── test_config/
│   └── main.go
├── your_new_script/
│   └── main.go
└── README.md
```

这样可以避免多个 `main` 包冲突的问题。

