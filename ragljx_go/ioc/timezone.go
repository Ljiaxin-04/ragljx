package ioc

import (
	"context"
	"log"
	"time"
)

const (
	timezoneName = "timezone"
)

var defaultTimezone = &TimezoneConfig{
	Location: "Asia/Shanghai",
}

// TimezoneConfig 应用时区配置
type TimezoneConfig struct {
	ObjectImpl
	Location string `json:"location" yaml:"location" env:"LOCATION"`
}

func (t *TimezoneConfig) Name() string {
	return timezoneName
}

// 高优先级，尽早执行
func (t *TimezoneConfig) Priority() int {
	return 900
}

func (t *TimezoneConfig) Init() error {
	locValue := t.Location
	if locValue == "" {
		locValue = defaultTimezone.Location
	}

	loc, err := time.LoadLocation(locValue)
	if err != nil {
		return err
	}
	time.Local = loc
	log.Printf("[timezone] set time.Local to %s", loc.String())
	return nil
}

func (t *TimezoneConfig) Close(ctx context.Context) {}

// Get 全局获取方法
func GetTimezone() *TimezoneConfig {
	obj := Config().Get(timezoneName)
	if obj == nil {
		return defaultTimezone
	}
	return obj.(*TimezoneConfig)
}

func init() {
	// 注册到 Config 空间，随配置文件/环境变量一起加载
	Config().Registry(defaultTimezone)
}
