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
	errNotImplemented     error             = errors.New("not implemented")
	errPositionNotFound   error             = errors.New("position was not found")
	errInstrumentNotFound error             = errors.New("instrument was not found in cache")
	_                     exchange.Exchange = (*TinkoffExchange)(nil)
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

func (ex *TinkoffExchange) GetAccount(ctx context.Context) (domain.Account, error) {
	return domain.Account{}, errNotImplemented
}

func (ex *TinkoffExchange) GetBalance(ctx context.Context, symbol string) (decimal.Decimal, error) {
	resp, err := ex.client.OperationsClient.GetPortfolio(ctx, &investapi.PortfolioRequest{AccountId: defaultAccount})
	if err != nil {
		return decimal.Zero, nil
	}

	instrument := ex.cache.instruments[symbol]
	if instrument == nil {
		return decimal.Zero, errInstrumentNotFound
	}

	for _, p := range resp.Positions {
		if p.Figi == instrument.FIGI {
			q := p.Quantity
			if q == nil {
				return decimal.Zero, nil
			}

			return utils.IntFractToDecimal(q.Units, q.Nano), nil
		}
	}

	return decimal.Zero, errNotImplemented
}

func (ex *TinkoffExchange) GetPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	instrument := ex.cache.instruments[symbol]
	if instrument == nil {
		return decimal.Zero, errInstrumentNotFound
	}

	respPrice, err := ex.client.MarketDataClient.GetLastPrices(ex.client.Ctx, &investapi.GetLastPricesRequest{Figi: []string{instrument.FIGI}})
	if err != nil {
		return decimal.Zero, fmt.Errorf("MarketDataClient.GetLastPrices: %w", err)
	}

	if len(respPrice.LastPrices) == 0 {
		return decimal.Zero, fmt.Errorf("MarketDataClient.GetLastPrices: empty LastPrices response received")
	}

	pricePerLot := utils.IntFractToDecimal(respPrice.LastPrices[0].Price.Units, respPrice.LastPrices[0].Price.Nano)

	return pricePerLot, errNotImplemented
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
	return domain.Order{}, errNotImplemented
}

func (ex *TinkoffExchange) GetOpenOrders(ctx context.Context) ([]domain.Order, error) {
	return nil, errNotImplemented
}

func (ex *TinkoffExchange) MakeOrder(ctx context.Context, o domain.Order) (domain.Order, error) {
	orderReq, err := ex.orderReqFromDomain(o, defaultAccount)
	if err != nil {
		return domain.Order{}, err
	}

	resp, err := ex.client.OrdersClient.PostOrder(ex.client.Ctx, orderReq)
	if err != nil {
		return domain.Order{}, fmt.Errorf("OrdersClient.PostOrder: %w", err)
	}

	o.ID = resp.OrderId

	return o, nil
}

var (
	tinkoffOrderType = map[domain.OrderType]investapi.OrderType{
		domain.OrderTypeLimit:  investapi.OrderType_ORDER_TYPE_LIMIT,
		domain.OrderTypeMarket: investapi.OrderType_ORDER_TYPE_MARKET,
	}

	tinkoffOrderDirection = map[domain.OrderDirection]investapi.OrderDirection{
		domain.OrderDirectionBuy:  investapi.OrderDirection_ORDER_DIRECTION_BUY,
		domain.OrderDirectionSell: investapi.OrderDirection_ORDER_DIRECTION_SELL,
	}
)

func (ex *TinkoffExchange) orderReqFromDomain(o domain.Order, accountId string) (*investapi.PostOrderRequest, error) {
	instrument := ex.cache.instruments[o.Symbol]
	if instrument == nil {
		return nil, errInstrumentNotFound
	}

	tinkoffOrder := investapi.PostOrderRequest{
		Figi:      instrument.FIGI,
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

	if ot, ok := tinkoffOrderType[o.OrderType]; ok {
		tinkoffOrder.OrderType = ot
	}

	if od, ok := tinkoffOrderDirection[o.Direction]; ok {
		tinkoffOrder.Direction = od
	}

	return &tinkoffOrder, nil
}

// stop loss, stop limit, take profit orders are supposed to be made using this method
func (ex *TinkoffExchange) MakeStopOrder(ctx context.Context, o domain.Order) (domain.Order, error) {
	orderReq, err := ex.stopReqFromDomain(o, defaultAccount)
	if err != nil {
		return domain.Order{}, err
	}

	resp, err := ex.client.StopOrdersClient.PostStopOrder(ex.client.Ctx, orderReq)
	if err != nil {
		return domain.Order{}, fmt.Errorf("StopOrdersClient.PostStopOrder: %w", err)
	}

	o.ID = resp.StopOrderId

	return o, nil
}

var (
	tinkoffStopOrderTypes = map[domain.OrderType]investapi.StopOrderType{
		domain.OrderTypeStopLoss:   investapi.StopOrderType_STOP_ORDER_TYPE_STOP_LOSS,
		domain.OrderTypeStopLimit:  investapi.StopOrderType_STOP_ORDER_TYPE_STOP_LIMIT,
		domain.OrderTypeTakeProfit: investapi.StopOrderType_STOP_ORDER_TYPE_TAKE_PROFIT,
	}

	tinkoffStopOrderDirection = map[domain.OrderDirection]investapi.StopOrderDirection{
		domain.OrderDirectionBuy:  investapi.StopOrderDirection_STOP_ORDER_DIRECTION_BUY,
		domain.OrderDirectionSell: investapi.StopOrderDirection_STOP_ORDER_DIRECTION_SELL,
	}
)

func (ex *TinkoffExchange) stopReqFromDomain(o domain.Order, accountId string) (*investapi.PostStopOrderRequest, error) {
	instrument := ex.cache.instruments[o.Symbol]
	if instrument == nil {
		return nil, errInstrumentNotFound
	}

	tinkoffOrder := investapi.PostStopOrderRequest{
		Figi:      instrument.FIGI,
		Quantity:  int64(o.Quantity),
		AccountId: accountId,
	}

	if o.Price != nil {
		priceUnits, priceNano := utils.Modf(*o.Price)
		tinkoffOrder.Price = &investapi.Quotation{Units: int64(priceUnits), Nano: int32(priceNano)}
	}

	if ot, ok := tinkoffStopOrderTypes[o.OrderType]; ok {
		tinkoffOrder.StopOrderType = ot
	}

	if od, ok := tinkoffStopOrderDirection[o.Direction]; ok {
		tinkoffOrder.Direction = od
	}

	return &tinkoffOrder, nil
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
	return nil, errNotImplemented
}

// currently only 1 min, 5 min, 15 min, 1 hour, 1 day intervals are supported
func (ex *TinkoffExchange) GetHistory(
	ctx context.Context, symbol string, start time.Time, end *time.Time, interval time.Duration,
) ([]*domain.Kline, error) {
	instrument := ex.cache.instruments[symbol]
	if instrument == nil {
		return nil, errInstrumentNotFound
	}

	interv := ex.durationToTinkoffInterval(interval)

	tsStart := timestamppb.New(start)

	tsEnd := timestamppb.New(time.Now())
	if end != nil {
		tsEnd = timestamppb.New(*end)
	}

	resp, err := ex.client.MarketDataClient.GetCandles(
		ex.client.Ctx,
		&investapi.GetCandlesRequest{
			Figi:     instrument.FIGI,
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
		t := k.Time.AsTime()

		klines = append(klines, &domain.Kline{Low: l, High: h, Open: o, Close: c, Volume: v, Ts: t})
	}

	return klines, nil
}

var tinkoffIntervals = map[time.Duration]investapi.CandleInterval{
	domain.Interval1Min:  investapi.CandleInterval_CANDLE_INTERVAL_1_MIN,
	domain.Interval5Min:  investapi.CandleInterval_CANDLE_INTERVAL_5_MIN,
	domain.Interval15Min: investapi.CandleInterval_CANDLE_INTERVAL_15_MIN,
	domain.Interval1Hour: investapi.CandleInterval_CANDLE_INTERVAL_HOUR,
	domain.Interval1Day:  investapi.CandleInterval_CANDLE_INTERVAL_DAY,
}

func (ex *TinkoffExchange) durationToTinkoffInterval(duration time.Duration) investapi.CandleInterval {
	interval, ok := tinkoffIntervals[duration]
	if !ok {
		return investapi.CandleInterval_CANDLE_INTERVAL_UNSPECIFIED
	}

	return interval
}

func (ex *TinkoffExchange) GetExchangeTime(ctx context.Context) (time.Time, error) {
	return time.Now(), errNotImplemented
}
