package di

import (
	"context"
	"fmt"

	"github.com/Minish144/crypto-trading-bot/clients"
	"github.com/Minish144/crypto-trading-bot/clients/binance"
	"github.com/Minish144/crypto-trading-bot/clients/bybit"
	"github.com/Minish144/crypto-trading-bot/config"
	"github.com/Minish144/crypto-trading-bot/helpers"
	"github.com/Minish144/crypto-trading-bot/logger"
	"github.com/Minish144/crypto-trading-bot/strategies"
	"github.com/Minish144/crypto-trading-bot/strategies/gridStrategy"
	"go.uber.org/zap"
)

type DI struct {
	Config *config.Config

	Exchanges struct {
		Bybit struct {
			Config     *bybit.Config
			HttpClient clients.HttpClient
		}

		Binance struct {
			Config     *binance.Config
			HttpClient clients.HttpClient
		}
	}

	Helpers struct {
		BinanceHelper *helpers.Helper
	}

	Strategies []strategies.Strategy
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
		dic.Exchanges.Bybit.HttpClient = bybit.NewBybitClient(bbCfg, cfg.Test)
	}

	if cfg.ExchangesEnables.Binance {
		bCfg, err := binance.NewBinanceConfig()
		if err != nil {
			return nil, fmt.Errorf("bybit.NewBinanceConfig: %w", err)
		}

		dic.Exchanges.Binance.Config = bCfg
		dic.Exchanges.Binance.HttpClient = binance.NewBinanceClient(bCfg, cfg.Test)
	}

	gridStrategyCfg, err := gridStrategy.NewConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("gridStrategy.NewConfigFromEnv: %w", err)
	}

	dic.Helpers.BinanceHelper = helpers.NewHelper(dic.Exchanges.Binance.HttpClient, dic.Config.BaseCoin)

	dic.Strategies = append(
		dic.Strategies,
		strategies.NewGridStrategy(
			dic.Exchanges.Binance.HttpClient,
			gridStrategyCfg,
		),
	)

	return dic, nil
}

func (dic *DI) Start(ctx context.Context) context.Context {
	z := zap.S().With("context", "di.Start")

	if dic.Config.ExchangesEnables.Bybit {
		if err := dic.Exchanges.Bybit.HttpClient.Ping(ctx); err != nil {
			z.Fatalw(
				"failed to ping bybit",
				"error", err.Error(),
			)
		}
	}

	if dic.Config.ExchangesEnables.Binance {
		if err := dic.Exchanges.Binance.HttpClient.Ping(ctx); err != nil {
			z.Fatalw(
				"failed to ping binance",
				"error", err.Error(),
			)
		}
	}

	go dic.Helpers.BinanceHelper.StartLoggingHelpers(ctx)

	for _, strategy := range dic.Strategies {
		z.Infow("starting strategy", "name", strategy.Name())
		go strategy.Start(ctx)
	}

	return ctx
}

func (dic *DI) Stop() {}
