package di

import (
	"context"
	"fmt"

	"github.com/Minish144/bybit-bot/clients"
	"github.com/Minish144/bybit-bot/clients/bybit"
	"github.com/Minish144/bybit-bot/config"
	"github.com/Minish144/bybit-bot/logger"
	"go.uber.org/zap"
)

type DI struct {
	Config *config.Config

	Exchanges struct {
		Bybit struct {
			Config     *bybit.Config
			HttpClient clients.HttpClient
		}
	}
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

	if cfg.ExchangesEnables.Bybit {
		bbCfg, err := bybit.NewBybitConfig()
		if err != nil {
			return nil, fmt.Errorf("bybit.NewBybitConfig: %w", err)
		}

		dic.Exchanges.Bybit.Config = bbCfg
		dic.Exchanges.Bybit.HttpClient = bybit.NewBybitClient(bbCfg)
	}

	return dic, nil
}

func (dic *DI) Start(ctx context.Context) context.Context {
	z := zap.S().With("context", "di.Start")

	if dic.Config.ExchangesEnables.Bybit {
		if err := dic.Exchanges.Bybit.HttpClient.Ping(); err != nil {
			z.Fatalw(
				"failed to ping bybit",
				"error", err.Error(),
			)
		}
	}

	return ctx
}

func (dic *DI) Stop() {}
