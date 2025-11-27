package jwt

import (
	"context"
	"log"
	jwtpkg "ragljx/internal/pkg/jwt"
	"ragljx/ioc"
)

const (
	AppName = "jwt"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Config{
	SecretKey:          "your-secret-key-change-in-production",
	ExpireHours:        2,
	RefreshExpireHours: 168,
}

type Config struct {
	ioc.ObjectImpl
	SecretKey          string `json:"secret_key" yaml:"secret_key" env:"SECRET_KEY"`
	ExpireHours        int    `json:"expire_hours" yaml:"expire_hours" env:"EXPIRE_HOURS"`
	RefreshExpireHours int    `json:"refresh_expire_hours" yaml:"refresh_expire_hours" env:"REFRESH_EXPIRE_HOURS"`
}

func (c *Config) Name() string {
	return AppName
}

func (c *Config) Priority() int {
	return 750
}

func (c *Config) Init() error {
	jwtpkg.SetConfig(c.SecretKey, c.ExpireHours, c.RefreshExpireHours)
	log.Printf("[jwt] config loaded, expire=%dh, refresh_expire=%dh", c.ExpireHours, c.RefreshExpireHours)
	return nil
}

func (c *Config) Close(ctx context.Context) {}

// Get 全局获取方法
func Get() *Config {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*Config)
}
