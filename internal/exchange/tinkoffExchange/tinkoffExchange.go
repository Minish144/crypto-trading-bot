package tinkoffExchange

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/minish144/crypto-trading-bot/domain"
	"github.com/minish144/crypto-trading-bot/internal/exchange"
	"github.com/minish144/crypto-trading-bot/pkg/tinkoff"
	"github.com/minish144/crypto-trading-bot/pkg/utils"
	"github.com/shopspring/decimal"
	investapi "github.com/tinkoff/invest-api-go-sdk"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrNotImplemented error             = errors.New("not implemented")
	_                 exchange.Exchange = (*TinkoffExchange)(nil)
)

const defaultAccount = ""

type TinkoffExchange struct {
	client *tinkoff.TinkoffAPI
	cache  *TinkoffCache
}

func NewTinkoffExchange(ctx context.Context, token string, sandbox bool) *TinkoffExchange {
	ex := &TinkoffExchange{
		client: tinkoff.New(ctx, token, sandbox),
		cache:  NewCache(),
	}

	instruments, err := loadInstruments(ex.client)
	if err != nil {
		log.Fatalf("loadInstruments: %s", err.Error())
	}

	ex.cache.AddInstruments(instruments)

	return ex
}

func loadInstruments(client *tinkoff.TinkoffAPI) ([]*Instrument, error) {
	respCurrencies, err := client.InstrumentsClient.Currencies(
		client.Ctx,
		&investapi.InstrumentsRequest{InstrumentStatus: investapi.InstrumentStatus_INSTRUMENT_STATUS_BASE},
	)
	if err != nil {
		return nil, fmt.Errorf("InstrumentsClient.Currencies: %w", err)
	}

	respStocks, err := client.InstrumentsClient.Shares(
		client.Ctx,
		&investapi.InstrumentsRequest{InstrumentStatus: investapi.InstrumentStatus_INSTRUMENT_STATUS_BASE},
	)
	if err != nil {
		return nil, fmt.Errorf("InstrumentsClient.Shares: %w", err)
	}

	respBonds, err := client.InstrumentsClient.Bonds(
		client.Ctx,
		&investapi.InstrumentsRequest{InstrumentStatus: investapi.InstrumentStatus_INSTRUMENT_STATUS_BASE},
	)
	if err != nil {
		return nil, fmt.Errorf("InstrumentsClient.Bonds: %w", err)
	}

	respEtfs, err := client.InstrumentsClient.Etfs(
		client.Ctx,
		&investapi.InstrumentsRequest{InstrumentStatus: investapi.InstrumentStatus_INSTRUMENT_STATUS_BASE},
	)
	if err != nil {
		return nil, fmt.Errorf("InstrumentsClient.Etfs: %w", err)
	}

	instruments := make([]*Instrument, 0,
		len(respCurrencies.Instruments)+
			len(respStocks.Instruments)+
			len(respBonds.Instruments)+
			len(respEtfs.Instruments),
	)

	instruments = append(instruments, currenciesToInstruments(respCurrencies.Instruments)...)
	instruments = append(instruments, stocksToInstruments(respStocks.Instruments)...)
	instruments = append(instruments, bondsToInstruments(respBonds.Instruments)...)
	instruments = append(instruments, etfsToInstruments(respEtfs.Instruments)...)

	return instruments, nil
}

func (ex *TinkoffExchange) GetAccount(ctx context.Context) (domain.Account, error) {
	return domain.Account{}, nil
}

func (ex *TinkoffExchange) GetPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	respPrice, err := ex.client.MarketDataClient.GetLastPrices(ex.client.Ctx, &investapi.GetLastPricesRequest{Figi: []string{symbol}})
	if err != nil {
		return decimal.Zero, fmt.Errorf("MarketDataClient.GetLastPrices: %w", err)
	}

	if len(respPrice.LastPrices) == 0 {
		return decimal.Zero, fmt.Errorf("MarketDataClient.GetLastPrices: %s", "empty LastPrices response received")
	}

	pricePerLot := utils.IntFractToDecimal(respPrice.LastPrices[0].Price.Units, respPrice.LastPrices[0].Price.Nano)

	return pricePerLot, ErrNotImplemented
}

// https://tinkoff.github.io/investAPI/faq_marketdata/
func (ex *TinkoffExchange) lotPriceByType(price, lot decimal.Decimal, instrumentType domain.InstrumentType) decimal.Decimal {
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

func (ex *TinkoffExchange) GetOrder(ctx context.Context, orderId string) (domain.Order, error) {
	return domain.Order{}, ErrNotImplemented
}

func (ex *TinkoffExchange) GetOpenOrders(ctx context.Context) ([]domain.Order, error) {
	return nil, ErrNotImplemented
}

func (ex *TinkoffExchange) MakeOrder(ctx context.Context, o domain.Order) (domain.Order, error) {
	orderReq := ex.orderReqFromDomain(o, defaultAccount)

	resp, err := ex.client.OrdersClient.PostOrder(ex.client.Ctx, &orderReq)
	if err != nil {
		return domain.Order{}, fmt.Errorf("OrdersClient.PostOrder: %w", err)
	}

	o.ID = resp.OrderId

	return o, nil
}

func (ex *TinkoffExchange) orderReqFromDomain(o domain.Order, accountId string) investapi.PostOrderRequest {
	tinkoffOrder := investapi.PostOrderRequest{
		Figi:      o.Symbol,
		Quantity:  int64(o.Quantity),
		AccountId: accountId,
	}

	if o.ID != "" {
		tinkoffOrder.OrderId = o.ID
	}

	if o.Price != nil {
		priceUnits, priceNano := utils.Modf(*o.Price)
		tinkoffOrder.Price = &investapi.Quotation{Units: int64(priceUnits), Nano: int32(priceNano)}
	}

	switch o.OrderType {
	case domain.OrderTypeLimit:
		tinkoffOrder.OrderType = investapi.OrderType_ORDER_TYPE_LIMIT
	case domain.OrderTypeMarket:
		tinkoffOrder.OrderType = investapi.OrderType_ORDER_TYPE_MARKET
	}

	switch o.Direction {
	case domain.OrderDirectionBuy:
		tinkoffOrder.Direction = investapi.OrderDirection_ORDER_DIRECTION_BUY
	case domain.OrderDirectionSell:
		tinkoffOrder.Direction = investapi.OrderDirection_ORDER_DIRECTION_SELL
	}

	return tinkoffOrder
}

// stop loss, stop limit, take profit orders are supposed to be made using this method
func (ex *TinkoffExchange) MakeStopOrder(ctx context.Context, o domain.Order) (domain.Order, error) {
	orderReq := ex.stopReqFromDomain(o, defaultAccount)

	resp, err := ex.client.StopOrdersClient.PostStopOrder(ex.client.Ctx, &orderReq)
	if err != nil {
		return domain.Order{}, fmt.Errorf("StopOrdersClient.PostStopOrder: %w", err)
	}

	o.ID = resp.StopOrderId

	return o, nil
}

func (ex *TinkoffExchange) stopReqFromDomain(o domain.Order, accountId string) investapi.PostStopOrderRequest {
	tinkoffOrder := investapi.PostStopOrderRequest{
		Figi:      o.Symbol,
		Quantity:  int64(o.Quantity),
		AccountId: accountId,
	}

	if o.Price != nil {
		priceUnits, priceNano := utils.Modf(*o.Price)
		tinkoffOrder.Price = &investapi.Quotation{Units: int64(priceUnits), Nano: int32(priceNano)}
	}

	switch o.OrderType {
	case domain.OrderTypeStopLoss:
		tinkoffOrder.StopOrderType = investapi.StopOrderType_STOP_ORDER_TYPE_STOP_LOSS
	case domain.OrderTypeStopLimit:
		tinkoffOrder.StopOrderType = investapi.StopOrderType_STOP_ORDER_TYPE_STOP_LIMIT
	case domain.OrderTypeTakeProfit:
		tinkoffOrder.StopOrderType = investapi.StopOrderType_STOP_ORDER_TYPE_TAKE_PROFIT
	}

	switch o.Direction {
	case domain.OrderDirectionBuy:
		tinkoffOrder.Direction = investapi.StopOrderDirection_STOP_ORDER_DIRECTION_BUY
	case domain.OrderDirectionSell:
		tinkoffOrder.Direction = investapi.StopOrderDirection_STOP_ORDER_DIRECTION_SELL
	}

	return tinkoffOrder
}

func (ex *TinkoffExchange) CancelOrder(ctx context.Context, orderId string) error {
	_, err := ex.client.OrdersClient.CancelOrder(ex.client.Ctx, &investapi.CancelOrderRequest{AccountId: defaultAccount, OrderId: orderId})
	if err != nil {
		return fmt.Errorf("OrdersClient.CancelOrder: %w", err)
	}

	return nil
}

func (ex *TinkoffExchange) CancelStopOrder(ctx context.Context, orderId string) error {
	_, err := ex.client.StopOrdersClient.CancelStopOrder(
		ex.client.Ctx,
		&investapi.CancelStopOrderRequest{AccountId: defaultAccount, StopOrderId: orderId},
	)
	if err != nil {
		return fmt.Errorf("OrdersClient.CancelOrder: %w", err)
	}

	return nil
}

func (ex *TinkoffExchange) GetFees(ctx context.Context) ([]domain.Fee, error) {
	return nil, ErrNotImplemented
}

// currently only 1 min, 5 min, 15 min, 1 hour, 1 day intervals are supported
func (ex *TinkoffExchange) GetHistory(
	ctx context.Context, symbol string, start time.Time, end *time.Time, interval time.Duration,
) ([]*domain.Kline, error) {
	interv := ex.durationToTinkoffInterval(interval)

	tsStart := timestamppb.New(start)

	tsEnd := timestamppb.New(time.Now())
	if end != nil {
		tsEnd = timestamppb.New(*end)
	}

	fmt.Println(symbol, tsStart, tsEnd, interv)

	resp, err := ex.client.MarketDataClient.GetCandles(
		ex.client.Ctx,
		&investapi.GetCandlesRequest{
			Figi:     symbol,
			From:     tsStart,
			To:       tsEnd,
			Interval: interv,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("MarketDataClient.GetCandles: %w", err)
	}

	klines := make([]*domain.Kline, 0, len(resp.Candles))

	for _, k := range resp.Candles {
		l := utils.IntFractToDecimal(k.Low.Units, k.Low.Nano)
		h := utils.IntFractToDecimal(k.High.Units, k.High.Nano)
		o := utils.IntFractToDecimal(k.Open.Units, k.Open.Nano)
		c := utils.IntFractToDecimal(k.Close.Units, k.Close.Nano)
		v := decimal.NewFromInt(k.Volume)

		klines = append(klines, &domain.Kline{Low: l, High: h, Open: o, Close: c, Volume: v})
	}

	return klines, nil
}

func (ex *TinkoffExchange) durationToTinkoffInterval(duration time.Duration) investapi.CandleInterval {
	minutes := duration.Minutes()
	switch {
	case minutes == 1:
		return investapi.CandleInterval_CANDLE_INTERVAL_1_MIN
	case minutes == 5:
		return investapi.CandleInterval_CANDLE_INTERVAL_5_MIN
	case minutes == 15:
		return investapi.CandleInterval_CANDLE_INTERVAL_15_MIN
	case minutes == 60:
		return investapi.CandleInterval_CANDLE_INTERVAL_HOUR
	case minutes == 1440:
		return investapi.CandleInterval_CANDLE_INTERVAL_DAY
	default:
		return investapi.CandleInterval_CANDLE_INTERVAL_UNSPECIFIED
	}
}

func (ex *TinkoffExchange) GetExchangeTime(ctx context.Context) (time.Time, error) {
	return time.Now(), ErrNotImplemented
}
