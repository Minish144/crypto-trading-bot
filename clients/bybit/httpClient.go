package bybit

import (
	"context"
	"errors"

	"github.com/Minish144/crypto-trading-bot/models"
	hirokisanBybit "github.com/hirokisan/bybit/v2"
)

var ErrNotImplemented = errors.New("not implemented")

type BybitClient struct {
	*hirokisanBybit.Client
}

func NewBybitClient(c *Config, test bool) *BybitClient {
	return &BybitClient{hirokisanBybit.NewClient().WithAuth(c.Key, c.Secret)}
}

func (c *BybitClient) Ping(ctx context.Context) error {
	return ErrNotImplemented
}

func (c *BybitClient) GetPrice(ctx context.Context, symbol string) (float64, error) {
	return 0, ErrNotImplemented
}

func (c *BybitClient) GetBalance(ctx context.Context, coin string) (float64, error) {
	return 0, ErrNotImplemented
}

func (c *BybitClient) NewOrder(
	ctx context.Context,
	symbol string,
	sideType models.SideType,
	orderType models.OrderType,
	price, quantity float64,
) error {
	return ErrNotImplemented
}

func (c *BybitClient) NewLimitBuyOrder(ctx context.Context, symbol string, price, quantity float64) error {
	return ErrNotImplemented
}

func (c *BybitClient) NewLimitSellOrder(ctx context.Context, symbol string, price, quantity float64) error {
	return ErrNotImplemented
}

func (c *BybitClient) NewMarketBuyOrder(ctx context.Context, symbol string, price, quantity float64) error {
	return ErrNotImplemented
}

func (c *BybitClient) NewMarketSellOrder(ctx context.Context, symbol string, price, quantity float64) error {
	return ErrNotImplemented
}

func (c *BybitClient) GetAssets(ctx context.Context) ([]models.Asset, error) {
	return nil, ErrNotImplemented
}

func (c *BybitClient) GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error) {
	return nil, ErrNotImplemented
}

func (c *BybitClient) CloseOrder(ctx context.Context, symbol string, orderId int64) error {
	return ErrNotImplemented
}

func (c *BybitClient) GetKlines(ctx context.Context, symbol, interval string) ([]*models.Kline, error) {
	return nil, ErrNotImplemented
}

func (c *BybitClient) GetKlinesCloses(ctx context.Context, symbol, interval string) ([]float64, error) {
	return nil, ErrNotImplemented
}
