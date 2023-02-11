package macdStrategy

import (
	"fmt"
	"time"

	"github.com/Minish144/crypto-trading-bot/utils"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Symbol string `env:"STRATEGIES_MACD_SYMBOL"   envDefault:"BTCUSDT"` // trading pair symbol
	Coins  struct {
		Quote string
		Base  string
	}
	Interval             time.Duration `env:"STRATEGIES_MACD_INTERVAL"                envDefault:"5m"`      // polling interval
	StopLossUpdatePeriod time.Duration `env:"STRATEGIES_MACD_STOP_LOSS_UPDATE_PERIOD" envDefault:"120m"`    // how often to update stop loss
	OrderAmount          float64       `env:"STRATEGIES_MACD_ORDER_AMOUNT"            envDefault:"0.00005"` // quote coin amount for placing order
}

func NewConfigFromEnv() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	sym, err := utils.ConvertSymbol(cfg.Symbol)
	if err != nil {
		return nil, fmt.Errorf("utils.ConvertSymbol: %w", err)
	}

	quote, err := utils.GetQuoteCoin(cfg.Symbol)
	if err != nil {
		return nil, fmt.Errorf("utils.ConvertSymbol: %w", err)
	}

	base, err := utils.GetBaseCoin(cfg.Symbol)
	if err != nil {
		return nil, fmt.Errorf("utils.ConvertSymbol: %w", err)
	}

	cfg.Symbol = sym
	cfg.Coins.Quote = quote
	cfg.Coins.Base = base

	return cfg, nil
}
