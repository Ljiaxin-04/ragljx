package log

import (
	"context"
	"log"
	"os"
	"ragljx/ioc"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	AppName = "log"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Logger{
	Level:  "info",
	Output: "stdout",
	Format: "json",
}

type Logger struct {
	ioc.ObjectImpl
	Level  string `json:"level" yaml:"level" env:"LEVEL"`
	Output string `json:"output" yaml:"output" env:"OUTPUT"`
	Format string `json:"format" yaml:"format" env:"FORMAT"`

	logger *zap.Logger
}

func (l *Logger) Name() string {
	return AppName
}

func (l *Logger) Priority() int {
	return 800
}

func (l *Logger) Init() error {
	// 解析日志级别
	level := zapcore.InfoLevel
	switch strings.ToLower(l.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}

	// 配置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var encoder zapcore.Encoder
	if l.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 配置输出
	var writer zapcore.WriteSyncer
	if l.Output == "stdout" {
		writer = zapcore.AddSync(os.Stdout)
	} else {
		file, err := os.OpenFile(l.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		writer = zapcore.AddSync(file)
	}

	core := zapcore.NewCore(encoder, writer, level)
	l.logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	log.Printf("[log] logger initialized with level: %s, format: %s", l.Level, l.Format)
	return nil
}

func (l *Logger) Close(ctx context.Context) {
	if l.logger != nil {
		l.logger.Sync()
	}
}

func (l *Logger) L() *zap.Logger {
	return l.logger
}

// Get 全局获取方法
func Get() *zap.Logger {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig.L()
	}
	return obj.(*Logger).L()
}

