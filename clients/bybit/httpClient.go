package bybit

import (
	hirokisanBybit "github.com/hirokisan/bybit/v2"
)

type BybitClient struct {
	*hirokisanBybit.Client
}

func NewBybitClient(c *Config) *BybitClient {
	client := hirokisanBybit.NewClient().WithAuth(c.Key, c.Secret)

	return &BybitClient{client}
}

func (c *BybitClient) Ping() error {
	_, err := c.Future().InverseFuture().APIKeyInfo()

	return err
}
