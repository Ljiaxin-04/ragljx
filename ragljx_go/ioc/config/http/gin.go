package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"ragljx/ioc"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	AppName = "http"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &GinServer{
	Host:         "0.0.0.0",
	Port:         8080,
	ReadTimeout:  60,
	WriteTimeout: 60,
}

type GinServer struct {
	ioc.ObjectImpl
	Host         string `json:"host" yaml:"host" env:"HOST"`
	Port         int    `json:"port" yaml:"port" env:"PORT"`
	ReadTimeout  int    `json:"read_timeout" yaml:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout int    `json:"write_timeout" yaml:"write_timeout" env:"WRITE_TIMEOUT"`

	engine *gin.Engine
	server *http.Server
}

func (g *GinServer) Name() string {
	return AppName
}

func (g *GinServer) Priority() int {
	return -99
}

func (g *GinServer) Init() error {
	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	g.engine = gin.New()
	g.engine.Use(gin.Recovery())

	g.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", g.Host, g.Port),
		Handler:      g.engine,
		ReadTimeout:  time.Duration(g.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(g.WriteTimeout) * time.Second,
	}

	log.Printf("[http] server initialized on %s:%d", g.Host, g.Port)
	return nil
}

func (g *GinServer) Start() error {
	log.Printf("[http] server starting on %s:%d", g.Host, g.Port)
	return g.server.ListenAndServe()
}

func (g *GinServer) Close(ctx context.Context) {
	if g.server != nil {
		log.Printf("[http] shutting down server...")
		if err := g.server.Shutdown(ctx); err != nil {
			log.Printf("[http] server shutdown error: %v", err)
		} else {
			log.Printf("[http] server stopped")
		}
	}
}

func (g *GinServer) Engine() *gin.Engine {
	return g.engine
}

// Get 全局获取方法
func Get() *GinServer {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*GinServer)
}

