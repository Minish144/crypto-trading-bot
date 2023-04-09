package di

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/minish144/crypto-trading-bot/domain"
	"github.com/minish144/crypto-trading-bot/internal/backtest"
	"github.com/minish144/crypto-trading-bot/internal/exchange/tinkoffExchange"
	"github.com/minish144/crypto-trading-bot/pkg/config"
	"github.com/minish144/crypto-trading-bot/pkg/logger"
	"github.com/shopspring/decimal"
)

type DI struct {
	Config *config.Config
}

func NewDI(ctx context.Context) (*DI, error) {
	dic := &DI{}

	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	dic.Config = cfg

	if logger.NewLogger(cfg) != nil {
		return nil, fmt.Errorf("logger: %w", err)
	}

	exchange := tinkoffExchange.NewTinkoffExchange(ctx, os.Getenv("TINKOFF_TOKEN"), false)
	_ = backtest.NewBacktest(ctx, "CNYRUB_TOM", exchange, nil, decimal.NewFromFloat(20000), time.Now().Add(-7*domain.Interval1Day), nil, domain.Interval1Hour)

	return dic, nil
}

func (dic *DI) Start(ctx context.Context) context.Context {
	return ctx
}

func (dic *DI) Stop(ctx context.Context) {}
