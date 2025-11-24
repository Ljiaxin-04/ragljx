package kafka

import (
	"context"
	"log"
	"ragljx/ioc"

	"github.com/segmentio/kafka-go"
)

const (
	AppName = "kafka"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Kafka{
	Brokers: []string{"localhost:9092"},
}

type Kafka struct {
	ioc.ObjectImpl
	Brokers  []string `json:"brokers" yaml:"brokers" env:"BROKERS" envSeparator:","`
	Username string   `json:"username" yaml:"username" env:"USERNAME"`
	Password string   `json:"password" yaml:"password" env:"PASSWORD"`
}

func (k *Kafka) Name() string {
	return AppName
}

func (k *Kafka) Priority() int {
	return 696
}

func (k *Kafka) Init() error {
	log.Printf("[kafka] initialized with brokers: %v", k.Brokers)
	return nil
}

func (k *Kafka) Close(ctx context.Context) {
	// Kafka writers/readers 由使用方管理
	log.Printf("[kafka] closed")
}

// Producer 创建 Producer
func (k *Kafka) Producer(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(k.Brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
}

// Consumer 创建 Consumer
func (k *Kafka) Consumer(groupID string, topics []string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: k.Brokers,
		GroupID: groupID,
		Topic:   topics[0], // 简化处理，实际可支持多 topic
	})
}

// Get 全局获取方法
func Get() *Kafka {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*Kafka)
}

