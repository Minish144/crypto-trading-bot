package tinkoffExchange

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/minish144/crypto-trading-bot/domain"
	"github.com/minish144/crypto-trading-bot/exchange"
	"github.com/minish144/crypto-trading-bot/pkg/tinkoff"
	"github.com/minish144/crypto-trading-bot/pkg/utils"
	"github.com/shopspring/decimal"
	investapi "github.com/tinkoff/invest-api-go-sdk"
)

var (
	ErrNotImplemented error             = errors.New("not implemented")
	_                 exchange.Exchange = (*TinkoffExchange)(nil)
)

const defaultAccount = ""

type TinkoffExchange struct {
	client tinkoff.TinkoffAPI
}

func NewTinkoffExchange(client tinkoff.TinkoffAPI) *TinkoffExchange {
	return &TinkoffExchange{client: client}
}

func (ex TinkoffExchange) GetAccount(ctx context.Context) (domain.Account, error) {
	return domain.Account{}, nil
}

func (ex TinkoffExchange) GetPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	respPrice, err := ex.client.MarketDataClient.GetLastPrices(ctx, &investapi.GetLastPricesRequest{Figi: []string{symbol}})
	if err != nil {
		return decimal.Zero, fmt.Errorf("MarketDataClient.GetLastPrices: %w", err)
	}

	if len(respPrice.LastPrices) == 0 {
		return decimal.Zero, fmt.Errorf("MarketDataClient.GetLastPrices: %s", "empty LastPrices response received")
	}

	pricePerLot := utils.IntFractToDecimal(respPrice.LastPrices[0].Price.Units, respPrice.LastPrices[0].Price.Nano)
	instrumentId := respPrice.LastPrices[0].InstrumentUid

	respInstrument, err := ex.client.InstrumentsClient.GetInstrumentBy(ctx, &investapi.InstrumentRequest{Id: instrumentId})
	if err != nil {
		return decimal.Zero, fmt.Errorf("InstrumentsClient.GetInstrumentBy: %w", err)
	}

	if respInstrument == nil || respInstrument.Instrument == nil {
		return decimal.Zero, fmt.Errorf("InstrumentsClient.GetInstrumentBy: %s", "instrument was not found")
	}

	lot := decimal.NewFromInt32(respInstrument.Instrument.Lot)
	iType := domain.InstrumentType(respInstrument.Instrument.InstrumentType)

	lotPrice := ex.lotPriceByType(pricePerLot, lot, iType)

	return lotPrice, ErrNotImplemented
}

// https://tinkoff.github.io/investAPI/faq_marketdata/
func (ex TinkoffExchange) lotPriceByType(price, lot decimal.Decimal, instrumentType domain.InstrumentType) decimal.Decimal {
	if lot.IsZero() {
		return price
	}

	switch instrumentType {
	case domain.InstrumentTypeBond:
		return price.Div(decimal.NewFromFloat(100)).Mul(lot)
	default:
		return price.Mul(lot)
	}
}

func (ex TinkoffExchange) GetOrder(ctx context.Context, orderId string) (domain.Order, error) {
	return domain.Order{}, ErrNotImplemented
}

func (ex TinkoffExchange) GetOpenOrders(ctx context.Context) ([]domain.Order, error) {
	return nil, ErrNotImplemented
}

func (ex TinkoffExchange) MarketOrder(ctx context.Context, o domain.Order) (domain.Order, error) {
	return domain.Order{}, ErrNotImplemented
}

func (ex TinkoffExchange) LimitOrder(ctx context.Context, o domain.Order) (domain.Order, error) {
	return domain.Order{}, ErrNotImplemented
}

func (ex TinkoffExchange) TakeProfitOrder(ctx context.Context, o domain.Order) (domain.Order, error) {
	return domain.Order{}, ErrNotImplemented
}

func (ex TinkoffExchange) StopLossOrder(ctx context.Context, o domain.Order) (domain.Order, error) {
	return domain.Order{}, ErrNotImplemented
}

func (ex TinkoffExchange) CancelOrder(ctx context.Context, o domain.Order) (domain.Order, error) {
	return domain.Order{}, ErrNotImplemented
}

func (ex TinkoffExchange) GetFees(ctx context.Context) ([]domain.Fee, error) {
	return nil, ErrNotImplemented
}

func (ex TinkoffExchange) GetHistory(
	ctx context.Context, symbol string, start time.Time, end *time.Time, interval time.Duration,
) ([]domain.Kline, error) {
	return nil, ErrNotImplemented
}

func (ex TinkoffExchange) GetExchangeTime(ctx context.Context) (time.Time, error) {
	return time.Time{}, ErrNotImplemented
}
