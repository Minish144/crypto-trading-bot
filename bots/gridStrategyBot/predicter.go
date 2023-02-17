package gridStrategyBot

import (
	"context"

	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/trader"
)

var _ trader.Predicter = (*Predcter)(nil)

type Predcter struct {
	klines          []*market.Kline
	maxKlinesAmount int32
	valid           bool
}

func NewPredicter(ctx context.Context, maxKlinesAmount int32) *Predcter {
	return &Predcter{
		klines:          make([]*market.Kline, 0),
		maxKlinesAmount: maxKlinesAmount,
		valid:           true,
	}
}

func (p *Predcter) ReceivePrice(ctx context.Context, k market.Kline) error {
	p.klines = append(p.klines, &k)

	return nil
}

func (p *Predcter) Predict() float64 {
	return 0
}

func (p *Predcter) Valid() bool {
	return p.valid
}
