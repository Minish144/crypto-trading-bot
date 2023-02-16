package binance

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	gobinance "github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/web"
)

var _ broker.Dealer = (*binanceDealer)(nil)

type binanceDealer struct {
	cfg *Config
	*futures.Client
}

const _newConfigFromEnvErr = "new config from env: %w"

func New(cfg *Config) (*binanceDealer, error) {
	if cfg == nil {
		cfgFromEnv, err := newBinanceConfigFromEnv()
		if err != nil {
			return nil, fmt.Errorf(_newConfigFromEnvErr, err)
		}

		cfg = cfgFromEnv
	}

	dealer := &binanceDealer{
		cfg: cfg,
	}

	if cfg.Test {
		gobinance.UseTestnet = true
		dealer.Client = futures.NewClient(cfg.TestKey, cfg.TestSecret)
	} else {
		dealer.Client = futures.NewClient(cfg.Key, cfg.Secret)
	}

	return dealer, nil
}

var ErrNotImplemented = errors.New("not implemented")

const (
	_accountServiceDoErrFmt  = "account service do: %w"
	_decimalFromStringErrFmt = "decimal from string: %w"
)

func (d *binanceDealer) GetBalance(ctx context.Context) (*broker.AccountBalance, *web.Response, error) {
	account, err := d.Client.
		NewGetAccountService().
		Do(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf(_accountServiceDoErrFmt, ParseError(err))
	}

	var (
		trade  decimal.Decimal
		equity decimal.Decimal
	)

	for _, balance := range account.Balances {
		_free, err := decimal.NewFromString(balance.Free)
		if err != nil {
			return nil, nil, fmt.Errorf(_decimalFromStringErrFmt, err)
		}

		_locked, err := decimal.NewFromString(balance.Locked)
		if err != nil {
			return nil, nil, fmt.Errorf(_decimalFromStringErrFmt, err)
		}

		if balance.Asset == d.cfg.Base {
			trade = trade.Add(_free)
			equity = equity.Add(_free).Add(_locked)

			continue
		}

		if _free.IsZero() && _locked.IsZero() {
			continue
		}

		assetPrice, err := d.getPrice(ctx, balance.Asset+d.cfg.Base)
		if err != nil {
			continue
		}

		_free = _free.Mul(assetPrice)
		_locked = _locked.Mul(assetPrice)

		trade = trade.Add(_free)
		equity = equity.Add(_free).Add(_locked)
	}

	return &broker.AccountBalance{
		Trade:  trade,
		Equity: equity,
	}, nil, nil
}

const (
	_listPriceServieDoErrFmt           = "list prices service do: %w"
	_listPriceServieDoEmptyArrayErrFmt = "list prices service do: empty prices array"
)

func (d *binanceDealer) getPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	prices, err := d.Client.NewListPricesService().
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		return decimal.Zero, fmt.Errorf(_listPriceServieDoErrFmt, ParseError(err))
	} else if len(prices) == 0 {
		return decimal.Zero, fmt.Errorf(_listPriceServieDoEmptyArrayErrFmt)
	}

	price, err := decimal.NewFromString(prices[len(prices)-1].Price)
	if err != nil {
		return decimal.Zero, fmt.Errorf(_decimalFromStringErrFmt, err)
	}

	return price, nil
}

var (
	_sideTypeMapTo = map[broker.OrderSide]gobinance.SideType{
		broker.Buy:  gobinance.SideTypeBuy,
		broker.Sell: gobinance.SideTypeSell,
	}

	_orderTypeMapTo = map[broker.OrderType]gobinance.OrderType{
		broker.Market: gobinance.OrderTypeMarket,
		broker.Limit:  gobinance.OrderTypeLimit,
	}
)

const (
	_orderServiceDoErrFmt      = "order service do: %w"
	_orderServiceNilRespErrFmt = "order service do: nil response"
	_alphakitOrderErrFmt       = "createOrderResponseToAlphakitOrder: %w"
)

var (
	_sideTypeMapFrom = map[gobinance.SideType]broker.OrderSide{
		gobinance.SideTypeBuy:  broker.Buy,
		gobinance.SideTypeSell: broker.Sell,
	}

	_orderTypeMapFrom = map[gobinance.OrderType]broker.OrderType{
		gobinance.OrderTypeMarket: broker.Market,
		gobinance.OrderTypeLimit:  broker.Limit,
	}
)

func (d *binanceDealer) PlaceOrder(ctx context.Context, o broker.Order) (*broker.Order, *web.Response, error) {
	svc := d.Client.NewCreateOrderService()

	if o.Type != broker.Market {
		svc.Price(o.LimitPrice.String())
		svc.TimeInForce(gobinance.TimeInForceTypeGTC)
	}

	if o.ID != "" {
		svc.NewClientOrderID(string(o.ID))
	}

	svc.Symbol(o.Asset.Symbol)
	svc.Side(_sideTypeMapTo[o.Side])
	svc.Type(_orderTypeMapTo[o.Type])
	svc.Quantity(o.Size.String())

	resp, err := svc.Do(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf(_orderServiceDoErrFmt, ParseError(err))
	} else if resp == nil {
		return nil, nil, fmt.Errorf(_orderServiceNilRespErrFmt)
	}

	oAlphakit, err := d.createOrderResponseToAlphakitOrder(resp)
	if err != nil {
		return nil, nil, fmt.Errorf(_alphakitOrderErrFmt, err)
	}

	oAlphakit.Size = o.Size

	return oAlphakit, nil, nil
}

const _getFillsErrFmt = "get fills: %w"

func (d *binanceDealer) createOrderResponseToAlphakitOrder(res *gobinance.CreateOrderResponse) (*broker.Order, error) {
	o := &broker.Order{
		ID:       broker.DealID(strconv.Itoa(int(res.OrderID))),
		OpenedAt: time.UnixMilli(res.TransactTime),
		Asset: market.Asset{
			Symbol: res.Symbol,
		},
		Side: _sideTypeMapFrom[res.Side],
		Type: _orderTypeMapFrom[res.Type],
	}

	size, err := decimal.NewFromString(res.OrigQuantity)
	if err != nil {
		return nil, fmt.Errorf(_decimalFromStringErrFmt, err)
	}
	o.Size = size

	price, err := decimal.NewFromString(res.Price)
	if err != nil {
		return nil, fmt.Errorf(_decimalFromStringErrFmt, err)
	}

	if o.Type == broker.Market {
		o.FilledAt = o.OpenedAt
		o.ClosedAt = o.OpenedAt

		if res.Status == gobinance.OrderStatusTypeFilled ||
			res.Status == gobinance.OrderStatusTypePartiallyFilled {
			fee, qty, price, err := d.getFills(res.Fills)
			if err != nil {
				return nil, fmt.Errorf(_getFillsErrFmt, err)
			}

			o.Fee = fee
			o.FilledPrice = price
			o.FilledSize = qty
		}
	} else {
		o.LimitPrice = price
	}

	return o, nil
}

func (d *binanceDealer) getFills(fills []*gobinance.Fill) (fee, qty, price decimal.Decimal, err error) {
	for _, fill := range fills {
		fillFee, err := decimal.NewFromString(fill.Commission)
		if err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, fmt.Errorf(_decimalFromStringErrFmt, err)
		}

		fillQty, err := decimal.NewFromString(fill.Quantity)
		if err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, fmt.Errorf(_decimalFromStringErrFmt, err)
		}

		fillPrice, err := decimal.NewFromString(fill.Price)
		if err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, fmt.Errorf(_decimalFromStringErrFmt, err)
		}

		fee = fee.Add(fillFee)
		qty = qty.Add(fillQty)
		price = price.Add(fillPrice)
	}

	return fee, qty, price, nil
}

const (
	_listOpenOrdersDoErrFmt = "list open orders do: %w"
	_cancelOrderDoErrFmt    = "cancel open order do: %w"
)

func (d *binanceDealer) CancelOrders(ctx context.Context) (*web.Response, error) {
	orders, err := d.Client.
		NewListOpenOrdersService().
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf(_listOpenOrdersDoErrFmt, ParseError(err))
	}

	for _, order := range orders {
		_, err := d.Client.
			NewCancelOrderService().
			Symbol(order.Symbol).
			OrderID(order.OrderID).
			Do(ctx)
		if err != nil {
			return nil, fmt.Errorf(_cancelOrderDoErrFmt, ParseError(err))
		}
	}

	return nil, nil
}

const _getAccountServiceDoErrFmt = "get account service do: %w"

func (d *binanceDealer) ListPositions(ctx context.Context, opts *web.ListOpts) ([]broker.Position, *web.Response, error) {
	return nil, nil, ErrNotImplemented
}

func (d *binanceDealer) ListRoundTurns(ctx context.Context, opts *web.ListOpts) ([]broker.RoundTurn, *web.Response, error) {
	return nil, nil, ErrNotImplemented
}
