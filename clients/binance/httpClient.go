package binance

import (
	"context"
	"fmt"

	"github.com/Minish144/crypto-trading-bot/models"
	"github.com/Minish144/crypto-trading-bot/utils"
	gobinance "github.com/adshao/go-binance/v2"
)

const defaultRecWindow int64 = 5000

var defaultRecWindowOption = gobinance.WithRecvWindow(defaultRecWindow)

type BinanceClient struct {
	*gobinance.Client
}

func NewBinanceClient(c *Config, test bool) *BinanceClient {
	client := &BinanceClient{}

	if test {
		gobinance.UseTestnet = true
		client.Client = gobinance.NewClient(c.TestKey, c.TestSecret)
	} else {
		client.Client = gobinance.NewClient(c.Key, c.Secret)
	}

	return client
}

func (c *BinanceClient) Ping(ctx context.Context) error {
	return c.NewPingService().Do(ctx, gobinance.WithHeader("mock", "mock", false))
}

func (c *BinanceClient) GetPrice(ctx context.Context, symbol string) (float64, error) {
	prices, err := c.NewListPricesService().
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		return 0, fmt.Errorf("client.NewListPricesService.Do: %w", err)
	} else if len(prices) == 0 {
		return 0, fmt.Errorf("client.NewListPricesService.Do: empty prices array received")
	}

	pricef64, err := utils.StringToFloat64(prices[len(prices)-1].Price)
	if err != nil {
		return 0, fmt.Errorf("utils.StringToFloat64: %w", err)
	}

	return pricef64, nil
}

func (c *BinanceClient) NewOrder(
	ctx context.Context,
	symbol string,
	sideType models.SideType,
	orderType models.OrderType,
	price, quantity float64,
) error {
	side, ok := stFromModels[sideType]
	if !ok {
		return fmt.Errorf("failed to convert field to binance-type field")
	}

	oType, ok := otFromModels[orderType]
	if !ok {
		return fmt.Errorf("failed to convert field to binance-type field")
	}

	return c.newBinanceOrder(
		ctx,
		symbol,
		side,
		oType,
		price,
		quantity,
	)
}

func (c *BinanceClient) NewLimitBuyOrder(ctx context.Context, symbol string, price, quantity float64) error {
	return c.newBinanceOrder(
		ctx,
		symbol,
		gobinance.SideTypeBuy,
		gobinance.OrderTypeLimit,
		price,
		quantity,
	)
}

func (c *BinanceClient) NewLimitSellOrder(
	ctx context.Context,
	symbol string,
	price,
	quantity float64,
) error {
	return c.newBinanceOrder(ctx, symbol, gobinance.SideTypeSell, gobinance.OrderTypeLimit, price, quantity)
}

const (
	pricePrecision    int = 2
	quantityPrecision     = 6
)

func (c *BinanceClient) newBinanceOrder(
	ctx context.Context,
	symbol string,
	sideType gobinance.SideType,
	orderType gobinance.OrderType,
	price, quantity float64,
) error {
	_, err := c.NewCreateOrderService().
		Symbol(symbol).
		Side(sideType).
		Type(orderType).
		Price(utils.Float64ToString(price, pricePrecision)).
		Quantity(utils.Float64ToString(quantity, quantityPrecision)).
		TimeInForce(gobinance.TimeInForceTypeGTC).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("client.NewCreateOrderService.Do: %v", err)
	}

	return nil
}

func (c *BinanceClient) GetBalance(ctx context.Context, coin string) (float64, error) {
	balances, err := c.NewGetAccountService().Do(ctx)
	if err != nil {
		return 0, fmt.Errorf("client.NewGetAccountService.Do: %s", err.Error())
	}

	var balance float64 = 0.0

	for _, asset := range balances.Balances {
		if asset.Asset == coin {
			balance, err = utils.StringToFloat64(asset.Free)
			if err != nil {
				return 0, fmt.Errorf("utils.StringToFloat64: %w", err)
			}

			break
		}
	}

	return balance, nil
}

func (c *BinanceClient) GetAssets(ctx context.Context) ([]models.Asset, error) {
	account, err := c.NewGetAccountService().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("client.NewGetAccountService: %s", err.Error())
	}

	assets := make([]models.Asset, len(account.Balances))

	for i, asset := range account.Balances {
		free, err := utils.StringToFloat64(asset.Free)
		if err != nil {
			return nil, fmt.Errorf("utils.StringToFloat64: %w", err)
		}

		locked, err := utils.StringToFloat64(asset.Locked)
		if err != nil {
			return nil, fmt.Errorf("utils.StringToFloat64: %w", err)
		}

		assets[i] = models.Asset{
			Coin:   asset.Asset,
			Free:   free,
			Locked: locked,
		}
	}

	return assets, nil
}

func (c *BinanceClient) GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error) {
	binanceOrders, err := c.NewListOpenOrdersService().
		Symbol(symbol).
		Do(ctx)

	orders := make([]*models.Order, len(binanceOrders))

	for i, binanceOrder := range binanceOrders {
		orders[i], err = OrdersToModel(binanceOrder)
		if err != nil {
			return nil, fmt.Errorf("OrdersToModel: %w", err)
		}
	}

	return orders, nil
}

func (c *BinanceClient) CloseOrder(ctx context.Context, symbol string, orderId int64) error {
	if _, err := c.NewCancelOrderService().
		Symbol(symbol).
		OrderID(orderId).
		Do(ctx); err != nil {
		return fmt.Errorf("c.NewCancelOrderService.Do: %w", err)
	}

	return nil
}

func (c *BinanceClient) NewStopLoss(ctx context.Context, symbol string) error {

}
