package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"myplants/internal/config"
)

var log *zap.Logger

// Init 初始化 zap 日志
func Init(cfg *config.LoggerConfig) error {
	if err := os.MkdirAll(cfg.Path, 0755); err != nil {
		return err
	}

	level := zapcore.InfoLevel
	_ = level.UnmarshalText([]byte(cfg.Level))

	var encoderCfg zapcore.EncoderConfig
	if cfg.Encoding == "json" {
		encoderCfg = zap.NewProductionEncoderConfig()
	} else {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	}
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	if cfg.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	logFile := filepath.Join(cfg.Path, cfg.Filename)
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	fileCore := zapcore.NewCore(encoder, zapcore.AddSync(f), level)
	consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)

	log = zap.New(zapcore.NewTee(fileCore, consoleCore), zap.AddCaller(), zap.AddCallerSkip(0))
	return nil
}

// L 返回 zap logger
func L() *zap.Logger {
	if log == nil {
		log, _ = zap.NewDevelopment()
	}
	return log
}

// S 返回 zap sugared logger
func S() *zap.SugaredLogger {
	return L().Sugar()
}
