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
	client := hirokisanBybit.NewClient().WithAuth(c.Key, c.Secret)

	return &BybitClient{client}
}

func (c *BybitClient) Ping(ctx context.Context) error {
	_, err := c.Future().USDTPerpetual().APIKeyInfo()

	return err
}

func (c *BybitClient) GetPrice(ctx context.Context, symbol string) (float64, error) {
	return 0, ErrNotImplemented
}

func (c *BybitClient) GetBalance(ctx context.Context, coin string) (float64, error) {
	return 0, ErrNotImplemented
}

func (c *BybitClient) NewLimitBuyOrder(ctx context.Context, symbol string, price, quantity float64) error {
	return ErrNotImplemented
}

func (c *BybitClient) NewLimitSellOrder(ctx context.Context, symbol string, price, quantity float64) error {
	return ErrNotImplemented
}

func (c *BybitClient) GetAssets(ctx context.Context) ([]models.Asset, error) {
	return nil, ErrNotImplemented
}
