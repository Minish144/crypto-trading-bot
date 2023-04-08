package exchange

import (
	"context"
	"time"

	"github.com/minish144/crypto-trading-bot/domain"
	"github.com/shopspring/decimal"
)

type Exchange interface {
	GetAccount(ctx context.Context) (domain.Account, error)
	GetPrice(ctx context.Context, symbol string) (decimal.Decimal, error)
	GetOrder(ctx context.Context, orderId string) (domain.Order, error)
	GetOpenOrders(ctx context.Context) ([]domain.Order, error)
	MakeOrder(ctx context.Context, o domain.Order) (domain.Order, error)
	MakeStopOrder(ctx context.Context, o domain.Order) (domain.Order, error)
	CancelOrder(ctx context.Context, o domain.Order) (domain.Order, error)
	GetFees(ctx context.Context) ([]domain.Fee, error)
	GetHistory(ctx context.Context, symbol string, start time.Time, end *time.Time, interval time.Duration) ([]*domain.Kline, error)
	GetExchangeTime(ctx context.Context) (time.Time, error)
}
