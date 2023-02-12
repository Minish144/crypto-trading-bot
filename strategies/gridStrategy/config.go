package gridStrategy

import (
	"fmt"
	"time"

	"github.com/Minish144/crypto-trading-bot/utils"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Symbol            string `env:"STRATEGIES_GRID_SYMBOL"          envDefault:"BTCUSDT"` // trading pair symbol
	PricePrecision    uint   `env:"STRATEGIES_GRID_PRICE_PRECISION" envDefault:"3"`       // price float64 precision
	QuantityPrecision uint   `env:"STRATEGIES_GRID_QTY_PRECISION"   envDefault:"3"`       // quantity float64 precision
	Coins             struct {
		Quote string
		Base  string
	}
	Interval              time.Duration `env:"STRATEGIES_GRID_INTERVAL"                 envDefault:"5m"`      // polling interval
	GridSize              float64       `env:"STRATEGIES_GRID_SIZE"                     envDefault:"0.01"`    // share of total funds to use for each grid level
	GridStep              float64       `env:"STRATEGIES_GRID_STEP"                     envDefault:"0.02"`    // share increase/decrease of the price for each subsequent grid level
	GridsAmount           uint          `env:"STRATEGIES_GRIDS_AMOUNT"                  envDefault:"3"`       // number of grids to create
	BaseCoinForAmount     bool          `env:"STRATEGIES_GRIDS_BASE_COIN_FOR_AMOUNT"    envDefault:"false"`   // whether to use base coin for ORDER_AMOUNT
	OrderAmount           float64       `env:"STRATEGIES_GRIDS_ORDER_AMOUNT"            envDefault:"0.00005"` // quote coin amount for placing order
	StopLossUpdatePeriod  time.Duration `env:"STRATEGIES_GRID_STOP_LOSS_UPDATE_PERIOD"  envDefault:"120m"`    // how often to update stop loss
	OrdersCheckRetriesMax uint          `env:"STRATEGIES_GRID_ORDERS_CHECK_RETRIES_MAX" envDefault:"3"`       // how many times try to make orders unit replacing
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
