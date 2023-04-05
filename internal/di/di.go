package di

import (
	"context"
	"fmt"

	"github.com/minish144/crypto-trading-bot/pkg/config"
	"github.com/minish144/crypto-trading-bot/pkg/logger"
)

type DI struct {
	Config *config.Config
}

func NewDI() (*DI, error) {
	dic := &DI{}

	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	dic.Config = cfg

	if logger.NewLogger(cfg) != nil {
		return nil, fmt.Errorf("logger: %w", err)
	}

	return dic, nil
}

func (dic *DI) Start(ctx context.Context) context.Context {
	return ctx
}

func (dic *DI) Stop(ctx context.Context) {}
