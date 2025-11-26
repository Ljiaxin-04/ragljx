package grpc

import (
	"context"
	"log"
	"ragljx/ioc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	AppName = "grpc"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &GRPCClient{
	PythonAddr: "localhost:50051",
}

type GRPCClient struct {
	ioc.ObjectImpl
	PythonAddr string `json:"python_addr" yaml:"python_addr" env:"PYTHON_ADDR"`

	conn *grpc.ClientConn
}

func (g *GRPCClient) Name() string {
	return AppName
}

func (g *GRPCClient) Priority() int {
	return 694
}

func (g *GRPCClient) Init() error {
	// 为了稳定暂用Dial，暂时不改为NewClient
	conn, err := grpc.Dial(g.PythonAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	g.conn = conn
	log.Printf("[grpc_client] connected to python service at %s", g.PythonAddr)
	return nil
}

func (g *GRPCClient) Close(ctx context.Context) {
	if g.conn != nil {
		g.conn.Close()
		log.Printf("[grpc_client] connection closed")
	}
}

func (g *GRPCClient) Conn() *grpc.ClientConn {
	return g.conn
}

// Get 全局获取方法
func Get() *grpc.ClientConn {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig.Conn()
	}
	return obj.(*GRPCClient).Conn()
}

