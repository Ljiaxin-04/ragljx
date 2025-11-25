package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"ragljx/ioc"
	"ragljx/ioc/config/datasource"
	httpConfig "ragljx/ioc/config/http"
	kafkaConfig "ragljx/ioc/config/kafka"
	minioConfig "ragljx/ioc/config/minio"
	redisConfig "ragljx/ioc/config/redis"
	grpcConfig "ragljx/ioc/config/grpc"
	"syscall"
	"time"

	// 导入 API 包以触发 init 注册
	_ "ragljx/internal/api/auth"
	_ "ragljx/internal/api/chat"
	_ "ragljx/internal/api/document"
	_ "ragljx/internal/api/knowledge_base"
	_ "ragljx/internal/api/user"
)

func main() {
	// 配置加载请求
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "config/application.yaml"
	req.ConfigEnv.Enabled = true
	req.ConfigEnv.Prefix = "RAGLJX"

	// 加载配置
	if err := ioc.ConfigIocObject(req); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 确保所有配置对象已加载
	_ = datasource.Get()
	_ = redisConfig.Get()
	_ = kafkaConfig.Get()
	_ = minioConfig.Get()
	_ = grpcConfig.Get()

	// 获取 HTTP 服务
	httpServer := httpConfig.Get()
	if httpServer == nil {
		log.Fatal("HTTP server not initialized")
	}

	// 打印已加载的对象
	log.Printf("Loaded configs: %v", ioc.Config().List())
	log.Printf("Loaded APIs: %v", ioc.Api().List())

	// 启动 HTTP 服务器
	go func() {
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	log.Println("Server started successfully")

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 关闭 HTTP 服务器
	httpServer.Close(ctx)

	log.Println("Server exited")
}

