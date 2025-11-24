package redis

import (
	"context"
	"log"
	"ragljx/ioc"

	"github.com/redis/go-redis/v9"
)

const (
	AppName = "redis"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Redis{
	Endpoints: []string{"localhost:6379"},
	Password:  "",
	DB:        0,
}

type Redis struct {
	ioc.ObjectImpl
	Endpoints []string `json:"endpoints" yaml:"endpoints" env:"ENDPOINTS" envSeparator:","`
	Password  string   `json:"password" yaml:"password" env:"PASSWORD"`
	DB        int      `json:"db" yaml:"db" env:"DB"`

	client redis.UniversalClient
}

func (r *Redis) Name() string {
	return AppName
}

func (r *Redis) Priority() int {
	return 698
}

func (r *Redis) Init() error {
	r.client = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    r.Endpoints,
		Password: r.Password,
		DB:       r.DB,
	})

	// 测试连接
	ctx := context.Background()
	if err := r.client.Ping(ctx).Err(); err != nil {
		return err
	}

	log.Printf("[redis] connected successfully to %v", r.Endpoints)
	return nil
}

func (r *Redis) Close(ctx context.Context) {
	if r.client != nil {
		r.client.Close()
		log.Printf("[redis] connection closed")
	}
}

func (r *Redis) Client() redis.UniversalClient {
	return r.client
}

// Get 全局获取方法
func Get() redis.UniversalClient {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig.Client()
	}
	return obj.(*Redis).Client()
}

