package logger

import (
	"github.com/minish144/crypto-trading-bot/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	debugLevel = "debug"
	infoLevel  = "info"
	warnLevel  = "warn"
	errorLevel = "error"
)

func NewLogger(cfg *config.Config) error {
	var level zapcore.Level

	switch cfg.LogLevel {
	case debugLevel:
		level = zapcore.DebugLevel
	case infoLevel:
		level = zapcore.InfoLevel
	case warnLevel:
		level = zapcore.WarnLevel
	case errorLevel:
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	zapCfg := zap.Config{
		Encoding:    cfg.LogEncoding,
		Level:       zap.NewAtomicLevelAt(level),
		OutputPaths: []string{cfg.LogOutput},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "msg",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
	}

	logger, err := zapCfg.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)

	return nil
}
