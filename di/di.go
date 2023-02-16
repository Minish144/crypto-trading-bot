package di

import (
	"context"
	"fmt"

	"github.com/Minish144/crypto-trading-bot/clients/binance"
	"github.com/Minish144/crypto-trading-bot/config"
	"github.com/Minish144/crypto-trading-bot/logger"
	"github.com/thecolngroup/alphakit/broker"
	"go.uber.org/zap"
)

type DI struct {
	Config *config.Config

	BinanceDealer broker.Dealer
}

func NewDI() (*DI, error) {
	dic := &DI{}

	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("config.NewConfig: %w", err)
	}

	dic.Config = cfg

	if logger.NewLogger(cfg) != nil {
		return nil, fmt.Errorf("logger.NewLogger: %w", err)
	}

	z := zap.S().With("context", "di.NewDI")

	dealer, err := binance.New(nil)
	if err != nil {
		z.Fatalw("binance new", "error", err.Error())
	}

	dic.BinanceDealer = dealer

	return dic, nil
}

func (dic *DI) Start(ctx context.Context) context.Context {
	z := zap.S().With("context", "di.Start")

	z.Infow("starting service")

	// fmt.Println(dic.BinanceDealer.PlaceOrder(ctx, broker.Order{
	// 	Asset:      market.Asset{Symbol: "BTCUSDT"},
	// 	Side:       broker.Sell,
	// 	Type:       broker.Limit,
	// 	Size:       decimal.NewFromFloat(0.01),
	// 	LimitPrice: decimal.NewFromFloat(24000.41),
	// }))

	fmt.Println(dic.BinanceDealer.CancelOrders(ctx))

	return ctx
}

func (dic *DI) Stop() {}
