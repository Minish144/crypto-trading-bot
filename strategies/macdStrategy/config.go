package macdStrategy

import (
	"fmt"
	"time"

	"github.com/Minish144/crypto-trading-bot/utils"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Symbol            string `env:"STRATEGIES_MACD_SYMBOL"          envDefault:"BTCUSDT"` // trading pair symbol
	PricePrecision    uint   `env:"STRATEGIES_MACD_PRICE_PRECISION" envDefault:"3"`       // price float64 precision
	QuantityPrecision uint   `env:"STRATEGIES_MACD_QTY_PRECISION"   envDefault:"3"`       // quantity float64 precision
	Coins             struct {
		Quote string
		Base  string
	}
	Interval             time.Duration `env:"STRATEGIES_MACD_INTERVAL"                envDefault:"5m"`      // polling interval
	StopLossUpdatePeriod time.Duration `env:"STRATEGIES_MACD_STOP_LOSS_UPDATE_PERIOD" envDefault:"180m"`    // how often to update stop loss
	StopLossShare        float64       `env:"STRATEGIES_MACD_STOP_LOSS_SHARE"         envDefault:"0.85"`    // stop loss share of actual price
	BaseCoinForAmount    bool          `env:"STRATEGIES_MACD_BASE_COIN_FOR_AMOUNT"    envDefault:"false"`   // whether to use base coin for ORDER_AMOUNT
	OrderAmount          float64       `env:"STRATEGIES_MACD_ORDER_AMOUNT"            envDefault:"0.00005"` // quote coin amount for placing order
	MaxOrdersAmount      float64       `env:"STRATEGIES_MACD_MAX_ORDERS_AMOUNT"       envDefault:"100"`     // amount available for trading
	KlinesInterval       string        `env:"STRATEGIES_MACD_KLINES_INTERVAL"         envDefault:"15m"`     // klines interval
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

	cfg.OrderAmount = utils.RoundPrecision(cfg.OrderAmount, cfg.QuantityPrecision)

	return cfg, nil
}
