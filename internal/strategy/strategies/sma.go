package strategies

import (
	"context"
	"fmt"
	"time"

	"github.com/minish144/crypto-trading-bot/internal/exchange"
	"github.com/minish144/crypto-trading-bot/internal/strategy"
	"github.com/minish144/crypto-trading-bot/internal/ta"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var _ strategy.Strategy = (*SMAStrategy)(nil)

type SMAStrategy struct {
	ex exchange.Exchange
}

func NewSMAStrategy(ex exchange.Exchange) *SMAStrategy {
	return &SMAStrategy{ex: ex}
}

func (s *SMAStrategy) Init(ctx context.Context) {}

func (s *SMAStrategy) PriceCallback(ctx context.Context, symbol string, price decimal.Decimal, ts time.Time, interval time.Duration) {
	z := zap.S().With("context", "SMAStrategy.PriceCallback", "symbol", symbol)

	z.Infow("new price event", "price", price.String())

	ts50TicksBefore := ts.Add(-interval * time.Duration(50))
	ts200TicksBefore := ts.Add(-interval * time.Duration(200))

	klines50, err := s.ex.GetHistory(ctx, symbol, ts50TicksBefore, &ts, interval)
	if err != nil {
		z.Errorw("failed to get history", "error", err, "klines", 50)
		return
	}

	klines200, err := s.ex.GetHistory(ctx, symbol, ts200TicksBefore, &ts, interval)
	if err != nil {
		z.Errorw("failed to get history", "error", err, "klines", 200)
		return
	}

	sma50 := ta.SMA(klines50)
	sma200 := ta.SMA(klines200)

	if sma50.GreaterThan(sma200) {
		fmt.Println("BUY")
	}

	if sma50.LessThan(sma200) {
		fmt.Println("SELL")
	}
}
