package minio

import (
	"context"
	"log"
	"ragljx/ioc"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	AppName = "minio"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &MinIO{
	Endpoint:        "localhost:9000",
	AccessKeyID:     "minioadmin",
	SecretAccessKey: "minioadmin",
	UseSSL:          false,
	BucketName:      "ragljx",
}

type MinIO struct {
	ioc.ObjectImpl
	Endpoint        string `json:"endpoint" yaml:"endpoint" env:"ENDPOINT"`
	AccessKeyID     string `json:"access_key_id" yaml:"access_key_id" env:"ACCESS_KEY_ID"`
	SecretAccessKey string `json:"secret_access_key" yaml:"secret_access_key" env:"SECRET_ACCESS_KEY"`
	UseSSL          bool   `json:"use_ssl" yaml:"use_ssl" env:"USE_SSL"`
	BucketName      string `json:"bucket_name" yaml:"bucket_name" env:"BUCKET_NAME"`

	client *minio.Client
}

func (m *MinIO) Name() string {
	return AppName
}

func (m *MinIO) Priority() int {
	return 695
}

func (m *MinIO) Init() error {
	client, err := minio.New(m.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.AccessKeyID, m.SecretAccessKey, ""),
		Secure: m.UseSSL,
	})
	if err != nil {
		return err
	}

	m.client = client

	// 确保 bucket 存在
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, m.BucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = client.MakeBucket(ctx, m.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		log.Printf("[minio] created bucket: %s", m.BucketName)
	}

	log.Printf("[minio] connected successfully to %s, bucket: %s", m.Endpoint, m.BucketName)
	return nil
}

func (m *MinIO) Close(ctx context.Context) {
	// MinIO client 不需要显式关闭
	log.Printf("[minio] closed")
}

func (m *MinIO) Client() *minio.Client {
	return m.client
}

func (m *MinIO) Bucket() string {
	return m.BucketName
}

// Get 全局获取方法
func Get() *MinIO {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*MinIO)
}

