package di

import (
	"context"
	"fmt"

	"github.com/Minish144/crypto-trading-bot/config"
	"github.com/Minish144/crypto-trading-bot/logger"
	"go.uber.org/zap"
)

type DI struct {
	Config *config.Config
}

func NewDI() (*DI, error) {
	dic := &DI{}

	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("config.NewConfig: %w", err)
	}

	dic.Config = cfg

	if logger.NewLogger(cfg) != nil {
		return nil, fmt.Errorf("logger.NewLogger: %w", err)
	}

	return dic, nil
}

func (dic *DI) Start(ctx context.Context) context.Context {
	z := zap.S().With("context", "di.Start")

	z.Infow("starting service")

	return ctx
}

func (dic *DI) Stop() {}
