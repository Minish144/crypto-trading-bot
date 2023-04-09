package backtestExchange

import (
	"context"
	"errors"
	"time"

	"github.com/minish144/crypto-trading-bot/domain"
	"github.com/minish144/crypto-trading-bot/internal/exchange"
	"github.com/minish144/crypto-trading-bot/internal/strategy"
	"github.com/shopspring/decimal"
)

type BackTestExchange interface {
	exchange.Exchange
}

var (
	_                 BackTestExchange = (*backtestExchange)(nil)
	errNotImplemented error            = errors.New("not implemented")
)

type wallet struct {
	InstrumentQuantity decimal.Decimal
}

type backtestExchange struct {
	symbol      string
	exchange    exchange.Exchange
	strategy    strategy.Strategy
	wallet      *wallet
	limitOrders []*domain.Order
}

func NewBacktestExchange(
	symbol string,
	exchange exchange.Exchange,
	strategy strategy.Strategy,
	startInstrumentQuantity decimal.Decimal,
	start time.Time,
	end *time.Time,
	interval time.Duration,
) *backtestExchange {
	wallet := &wallet{InstrumentQuantity: startInstrumentQuantity}
	limitOrders := make([]*domain.Order, 0)

	return &backtestExchange{
		symbol:      symbol,
		exchange:    exchange,
		strategy:    strategy,
		wallet:      wallet,
		limitOrders: limitOrders,
	}
}

func (b *backtestExchange) GetAccount(ctx context.Context) (domain.Account, error) {
	return domain.Account{}, errNotImplemented
}

func (b *backtestExchange) GetBalance(ctx context.Context, symbol string) (decimal.Decimal, error) {
	return decimal.Zero, errNotImplemented
}

func (b *backtestExchange) GetPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	return b.exchange.GetPrice(ctx, symbol)
}

func (b *backtestExchange) GetOrder(ctx context.Context, orderId string) (domain.Order, error) {
	return b.exchange.GetOrder(ctx, orderId)
}

func (b *backtestExchange) GetOpenOrders(ctx context.Context) ([]domain.Order, error) {
	return b.exchange.GetOpenOrders(ctx)
}

func (b *backtestExchange) MakeOrder(ctx context.Context, o domain.Order) (domain.Order, error) {
	return b.exchange.MakeOrder(ctx, o)
}

func (b *backtestExchange) MakeStopOrder(ctx context.Context, o domain.Order) (domain.Order, error) {
	return b.exchange.MakeStopOrder(ctx, o)
}

func (b *backtestExchange) CancelOrder(ctx context.Context, orderId string) error {
	return b.exchange.CancelOrder(ctx, orderId)
}

func (b *backtestExchange) CancelStopOrder(ctx context.Context, orderId string) error {
	return b.exchange.CancelOrder(ctx, orderId)
}

func (b *backtestExchange) GetFees(ctx context.Context) ([]domain.Fee, error) {
	return b.exchange.GetFees(ctx)
}

func (b *backtestExchange) GetHistory(ctx context.Context, symbol string, start time.Time, end *time.Time, interval time.Duration) ([]*domain.Kline, error) {
	return b.exchange.GetHistory(ctx, symbol, start, end, interval)
}

func (b *backtestExchange) GetExchangeTime(ctx context.Context) (time.Time, error) {
	return b.exchange.GetExchangeTime(ctx)
}
