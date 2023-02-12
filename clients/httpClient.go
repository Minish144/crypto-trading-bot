package clients

import (
	"context"

	"github.com/Minish144/crypto-trading-bot/models"
)

type HttpClient interface {
	Ping(ctx context.Context) error
	GetPrice(ctx context.Context, symbol string) (float64, error)
	GetBalance(ctx context.Context, coin string) (float64, float64, error)
	GetAssets(ctx context.Context) ([]models.Asset, error)
	GetKlines(ctx context.Context, symbol, interval string) ([]*models.Kline, error)
	GetKlinesCloses(ctx context.Context, symbol, interval string) ([]float64, error)

	// Orders
	NewOrder(
		ctx context.Context,
		symbol string,
		sideType models.SideType,
		orderType models.OrderType,
		tif models.TimeInForceType,
		price, quantity float64,
	) error
	NewLimitBuyOrder(ctx context.Context, symbol string, price, quantity float64) error
	NewLimitSellOrder(ctx context.Context, symbol string, price, quantity float64) error
	NewMarketBuyOrder(ctx context.Context, symbol string, quantity float64) error
	NewMarketSellOrder(ctx context.Context, symbol string, quantity float64) error
	GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error)
	CloseOrder(ctx context.Context, symbol string, orderId int64) error
}
