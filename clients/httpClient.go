package clients

import (
	"context"

	"github.com/Minish144/crypto-trading-bot/models"
)

type HttpClient interface {
	Ping(ctx context.Context) error
	GetPrice(ctx context.Context, symbol string) (float64, error)
	GetBalance(ctx context.Context, coin string) (float64, error)
	GetAssets(ctx context.Context) ([]models.Asset, error)

	// Orders
	NewLimitBuyOrder(ctx context.Context, symbol string, price, quantity float64) error
	NewLimitSellOrder(ctx context.Context, symbol string, price, quantity float64) error
	GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error)
	CloseOrder(ctx context.Context, symbol string, orderId int64) error
}
