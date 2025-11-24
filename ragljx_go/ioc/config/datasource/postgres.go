package datasource

import (
	"context"
	"fmt"
	"log"
	"ragljx/ioc"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	AppName = "postgres"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &PostgresDB{
	Host:         "localhost",
	Port:         5432,
	Database:     "ragljx",
	Username:     "ragljx",
	Password:     "ragljx_password",
	SSLMode:      "disable",
	Debug:        false,
	MaxIdleConns: 10,
	MaxOpenConns: 100,
}

type PostgresDB struct {
	ioc.ObjectImpl
	Host         string `json:"host" yaml:"host" env:"HOST"`
	Port         int    `json:"port" yaml:"port" env:"PORT"`
	Database     string `json:"database" yaml:"database" env:"DATABASE"`
	Username     string `json:"username" yaml:"username" env:"USERNAME"`
	Password     string `json:"password" yaml:"password" env:"PASSWORD"`
	SSLMode      string `json:"ssl_mode" yaml:"ssl_mode" env:"SSL_MODE"`
	Debug        bool   `json:"debug" yaml:"debug" env:"DEBUG"`
	MaxIdleConns int    `json:"max_idle_conns" yaml:"max_idle_conns" env:"MAX_IDLE_CONNS"`
	MaxOpenConns int    `json:"max_open_conns" yaml:"max_open_conns" env:"MAX_OPEN_CONNS"`

	db *gorm.DB
}

func (p *PostgresDB) Name() string {
	return AppName
}

func (p *PostgresDB) Priority() int {
	return 700
}

func (p *PostgresDB) Init() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
		p.Host, p.Username, p.Password, p.Database, p.Port, p.SSLMode,
	)

	logLevel := logger.Silent
	if p.Debug {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(p.MaxIdleConns)
	sqlDB.SetMaxOpenConns(p.MaxOpenConns)

	p.db = db
	log.Printf("[postgres] connected successfully to %s:%d/%s", p.Host, p.Port, p.Database)
	return nil
}

func (p *PostgresDB) Close(ctx context.Context) {
	if p.db != nil {
		sqlDB, _ := p.db.DB()
		if sqlDB != nil {
			sqlDB.Close()
			log.Printf("[postgres] connection closed")
		}
	}
}

func (p *PostgresDB) DB() *gorm.DB {
	return p.db
}

// Get 全局获取方法
func Get() *gorm.DB {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig.DB()
	}
	return obj.(*PostgresDB).DB()
}

