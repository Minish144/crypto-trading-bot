package clients

import "context"

type HttpClient interface {
	Ping(ctx context.Context) error
	GetPrice(ctx context.Context, symbol string) (float64, error)
	GetBalance(ctx context.Context, coin string) (float64, error)

	// Orders
	NewLimitBuyOrder(ctx context.Context, symbol string, price, quantity float64) error
	NewLimitSellOrder(ctx context.Context, symbol string, price, quantity float64) error
}
