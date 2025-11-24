package main

import (
	"fmt"
	"log"
	"ragljx/ioc"
	"ragljx/ioc/config/datasource"
	grpcConfig "ragljx/ioc/config/grpc"
	httpConfig "ragljx/ioc/config/http"
	kafkaConfig "ragljx/ioc/config/kafka"
	logConfig "ragljx/ioc/config/log"
	minioConfig "ragljx/ioc/config/minio"
	redisConfig "ragljx/ioc/config/redis"
)

func main() {
	// 配置加载请求
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "config/application.yaml"
	req.ConfigEnv.Enabled = false // 禁用环境变量，只测试文件加载
	req.ConfigEnv.Prefix = "RAGLJX"

	// 只加载配置文件，不初始化对象
	if err := ioc.DefaultStore.LoadFromFile(req.ConfigFile.Path); err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	fmt.Println("=== Configuration Test ===")
	fmt.Println()

	// 测试 PostgreSQL 配置
	fmt.Println("PostgreSQL Config:")
	pgObj := ioc.Config().Get("postgres")
	if pgObj != nil {
		pg := pgObj.(*datasource.PostgresDB)
		fmt.Printf("  Host: %s\n", pg.Host)
		fmt.Printf("  Port: %d\n", pg.Port)
		fmt.Printf("  Database: %s\n", pg.Database)
		fmt.Printf("  Username: %s\n", pg.Username)
		fmt.Printf("  Debug: %v\n", pg.Debug)
	}
	fmt.Println()

	// 测试 Redis 配置
	fmt.Println("Redis Config:")
	redisObj := ioc.Config().Get("redis")
	if redisObj != nil {
		redis := redisObj.(*redisConfig.Redis)
		fmt.Printf("  Endpoints: %v\n", redis.Endpoints)
		fmt.Printf("  DB: %d\n", redis.DB)
	}
	fmt.Println()

	// 测试 Kafka 配置
	fmt.Println("Kafka Config:")
	kafkaObj := ioc.Config().Get("kafka")
	if kafkaObj != nil {
		kafka := kafkaObj.(*kafkaConfig.Kafka)
		fmt.Printf("  Brokers: %v\n", kafka.Brokers)
	}
	fmt.Println()

	// 测试 MinIO 配置
	fmt.Println("MinIO Config:")
	minioObj := ioc.Config().Get("minio")
	if minioObj != nil {
		minio := minioObj.(*minioConfig.MinIO)
		fmt.Printf("  Endpoint: %s\n", minio.Endpoint)
		fmt.Printf("  AccessKeyID: %s\n", minio.AccessKeyID)
		fmt.Printf("  BucketName: %s\n", minio.BucketName)
		fmt.Printf("  UseSSL: %v\n", minio.UseSSL)
	}
	fmt.Println()

	// 测试 gRPC 配置
	fmt.Println("gRPC Config:")
	grpcObj := ioc.Config().Get("grpc")
	if grpcObj != nil {
		grpc := grpcObj.(*grpcConfig.GRPCClient)
		fmt.Printf("  PythonAddr: %s\n", grpc.PythonAddr)
	}
	fmt.Println()

	// 测试 HTTP 配置
	fmt.Println("HTTP Config:")
	httpObj := ioc.Config().Get("http")
	if httpObj != nil {
		http := httpObj.(*httpConfig.GinServer)
		fmt.Printf("  Host: %s\n", http.Host)
		fmt.Printf("  Port: %d\n", http.Port)
		fmt.Printf("  ReadTimeout: %d\n", http.ReadTimeout)
		fmt.Printf("  WriteTimeout: %d\n", http.WriteTimeout)
	}
	fmt.Println()

	// 测试 Log 配置
	fmt.Println("Log Config:")
	logObj := ioc.Config().Get("log")
	if logObj != nil {
		logger := logObj.(*logConfig.Logger)
		fmt.Printf("  Level: %s\n", logger.Level)
		fmt.Printf("  Output: %s\n", logger.Output)
		fmt.Printf("  Format: %s\n", logger.Format)
	}
	fmt.Println()

	fmt.Println("=== All configurations loaded successfully! ===")
}

