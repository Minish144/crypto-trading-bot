package backtest

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/minish144/crypto-trading-bot/domain"
	"github.com/minish144/crypto-trading-bot/internal/backtest/cursor"
	"github.com/minish144/crypto-trading-bot/internal/exchange"
	"github.com/minish144/crypto-trading-bot/internal/strategy"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type Backtest interface {
	exchange.Exchange

	Start() error
}

var (
	_                 Backtest = (*backtest)(nil)
	errNotImplemented error    = errors.New("not implemented")
)

type wallet struct {
	InstrumentQuantity decimal.Decimal
}

type backtest struct {
	ctx         context.Context
	symbol      string
	exchange    exchange.Exchange
	strategy    strategy.Strategy
	cursor      cursor.Cursor
	wallet      *wallet
	limitOrders []*domain.Order
}

func NewBacktest(
	ctx context.Context,
	symbol string,
	exchange exchange.Exchange,
	strategy strategy.Strategy,
	startInstrumentQuantity decimal.Decimal,
	start time.Time,
	end *time.Time,
	interval time.Duration,
) *backtest {
	wallet := &wallet{InstrumentQuantity: startInstrumentQuantity}
	limitOrders := make([]*domain.Order, 0)

	klines, err := exchange.GetHistory(ctx, symbol, start, end, interval)
	if err != nil {
		log.Fatalf("exchange.GetHistory: %s\n", err.Error())
	}

	z := zap.S().With("context", "NewBacktest")

	z.Infow("klines loaded", "count", len(klines))

	cursor := cursor.NewFromKlines(klines)

	return &backtest{
		ctx:         ctx,
		symbol:      symbol,
		exchange:    exchange,
		strategy:    strategy,
		cursor:      cursor,
		wallet:      wallet,
		limitOrders: limitOrders,
	}
}

func (b *backtest) Start() error {
	return errNotImplemented
}

func (b *backtest) GetAccount(ctx context.Context) (domain.Account, error) {
	return domain.Account{}, errNotImplemented
}

func (b *backtest) GetBalance(ctx context.Context, symbol string) (decimal.Decimal, error) {
	return decimal.Zero, errNotImplemented
}

func (b *backtest) GetPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	return b.exchange.GetPrice(ctx, symbol)
}

func (b *backtest) GetOrder(ctx context.Context, orderId string) (domain.Order, error) {
	return b.exchange.GetOrder(ctx, orderId)
}

func (b *backtest) GetOpenOrders(ctx context.Context) ([]domain.Order, error) {
	return b.exchange.GetOpenOrders(ctx)
}

func (b *backtest) MakeOrder(ctx context.Context, o domain.Order) (domain.Order, error) {
	return b.exchange.MakeOrder(ctx, o)
}

func (b *backtest) MakeStopOrder(ctx context.Context, o domain.Order) (domain.Order, error) {
	return b.exchange.MakeStopOrder(ctx, o)
}

func (b *backtest) CancelOrder(ctx context.Context, orderId string) error {
	return b.exchange.CancelOrder(ctx, orderId)
}

func (b *backtest) CancelStopOrder(ctx context.Context, orderId string) error {
	return b.exchange.CancelOrder(ctx, orderId)
}

func (b *backtest) GetFees(ctx context.Context) ([]domain.Fee, error) {
	return b.exchange.GetFees(ctx)
}

func (b *backtest) GetHistory(ctx context.Context, symbol string, start time.Time, end *time.Time, interval time.Duration) ([]*domain.Kline, error) {
	return b.exchange.GetHistory(ctx, symbol, start, end, interval)
}

func (b *backtest) GetExchangeTime(ctx context.Context) (time.Time, error) {
	return b.exchange.GetExchangeTime(ctx)
}
