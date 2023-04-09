package strategy

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Strategy interface {
	Init(ctx context.Context)
	PriceCallback(ctx context.Context, symbol string, price decimal.Decimal, ts time.Time, interval time.Duration)
}
