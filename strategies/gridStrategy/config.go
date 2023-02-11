package gridStrategy

import (
	"fmt"
	"time"

	"github.com/Minish144/crypto-trading-bot/utils"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Symbol string `env:"STRATEGIES_GRID_SYMBOL"   envDefault:"BTCUSDT"` // trading pair symbol
	Coins  struct {
		Quote string
		Base  string
	}
	Interval    time.Duration `env:"STRATEGIES_GRID_INTERVAL"      envDefault:"5m"`      // polling interval
	GridSize    float64       `env:"STRATEGIES_GRID_SIZE"          envDefault:"0.01"`    // share of total funds to use for each grid level
	GridStep    float64       `env:"STRATEGIES_GRID_STEP"          envDefault:"0.005"`   // share increase/decrease of the price for each subsequent grid level
	GridsAmount int           `env:"STRATEGIES_GRIDS_AMOUNT"       envDefault:"3"`       // number of grids to create
	OrderAmount float64       `env:"STRATEGIES_GRIDS_ORDER_AMOUNT" envDefault:"0.00005"` // quote coin amount for placing order
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
