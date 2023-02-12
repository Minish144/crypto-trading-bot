package binance

import (
	"context"
	"fmt"

	"github.com/Minish144/crypto-trading-bot/models"
	"github.com/Minish144/crypto-trading-bot/utils"
	gobinance "github.com/adshao/go-binance/v2"
)

const (
	pricePrecision    int = 6
	quantityPrecision     = 8
)

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
		return 0, fmt.Errorf("c.NewListPricesService.Do: %w", ParseError(err))
	} else if len(prices) == 0 {
		return 0, fmt.Errorf("c.NewListPricesService.Do: empty prices array received")
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
	tif models.TimeInForceType,
	price, quantity float64,
) error {
	side, ok := stFromModels[sideType]
	if !ok {
		return fmt.Errorf("failed to convert side to binance-type field")
	}

	oType, ok := otFromModels[orderType]
	if !ok {
		return fmt.Errorf("failed to convert orderType to binance-type field")
	}

	tifType, ok := tiftFromModels[tif]
	if !ok {
		return fmt.Errorf("failed to convert timeInForce to binance-type field")
	}

	return c.newBinanceOrder(
		ctx,
		symbol,
		side,
		oType,
		tifType,
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
		gobinance.TimeInForceTypeGTC,
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
	return c.newBinanceOrder(
		ctx,
		symbol,
		gobinance.SideTypeSell,
		gobinance.OrderTypeLimit,
		gobinance.TimeInForceTypeFOK,
		price,
		quantity,
	)
}

func (c *BinanceClient) NewMarketBuyOrder(ctx context.Context, symbol string, quantity float64) error {
	return c.newBinanceOrder(
		ctx,
		symbol,
		gobinance.SideTypeBuy,
		gobinance.OrderTypeMarket,
		gobinance.TimeInForceTypeFOK,
		0,
		quantity,
	)
}

func (c *BinanceClient) NewMarketSellOrder(ctx context.Context, symbol string, quantity float64) error {
	return c.newBinanceOrder(
		ctx,
		symbol,
		gobinance.SideTypeSell,
		gobinance.OrderTypeMarket,
		gobinance.TimeInForceTypeFOK,
		0,
		quantity,
	)
}

func (c *BinanceClient) newBinanceOrder(
	ctx context.Context,
	symbol string,
	sideType gobinance.SideType,
	orderType gobinance.OrderType,
	tif gobinance.TimeInForceType,
	price, quantity float64,
) error {
	request := c.NewCreateOrderService().
		Symbol(symbol).
		Side(sideType).
		Type(orderType).
		Quantity(utils.Float64ToString(quantity, quantityPrecision))

	if orderType != gobinance.OrderTypeMarket {
		request = request.Price(utils.Float64ToString(price, pricePrecision))
		request = request.TimeInForce(tif)
	}

	_, err := request.Do(ctx)
	if err != nil {
		return fmt.Errorf("c.NewCreateOrderService.Do: %v", ParseError(err))
	}

	return nil
}

func (c *BinanceClient) GetBalance(ctx context.Context, coin string) (float64, float64, error) {
	balances, err := c.NewGetAccountService().Do(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("c.NewGetAccountService.Do: %w", ParseError(err))
	}

	var (
		balance float64 = 0.0
		locked  float64 = 0.0
	)

	for _, asset := range balances.Balances {
		if asset.Asset == coin {
			balance, err = utils.StringToFloat64(asset.Free)
			if err != nil {
				return 0, 0, fmt.Errorf("utils.StringToFloat64: %w", err)
			}

			locked, err = utils.StringToFloat64(asset.Locked)
			if err != nil {
				return 0, 0, fmt.Errorf("utils.StringToFloat64: %w", err)
			}

			break
		}
	}

	return balance, locked, nil
}

func (c *BinanceClient) GetAssets(ctx context.Context) ([]models.Asset, error) {
	account, err := c.NewGetAccountService().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("c.NewGetAccountService.Do: %w", ParseError(err))
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
	if err != nil {
		return nil, fmt.Errorf("c.NewListOpenOrdersService.Do: %w", ParseError(err))
	}

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
		return fmt.Errorf("c.NewCancelOrderService.Do: %w", ParseError(err))
	}

	return nil
}

func (c *BinanceClient) GetKlines(ctx context.Context, symbol, interval string) ([]*models.Kline, error) {
	klines, err := c.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("c.NewKlinesService.Do: %w", ParseError(err))
	}

	klinesModels := make([]*models.Kline, len(klines))

	for i, kline := range klines {
		klinesModels[i], err = KlineToModel(kline)
		if err != nil {
			return nil, fmt.Errorf("KlineToModel: %w", err)
		}
	}

	return klinesModels, nil
}

func (c *BinanceClient) GetKlinesCloses(ctx context.Context, symbol, interval string) ([]float64, error) {
	klines, err := c.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("c.NewKlinesService.Do: %w", ParseError(err))
	}

	closes := make([]float64, len(klines))

	for i, klines := range klines {
		closeF64, err := utils.StringToFloat64(klines.Close)
		if err != nil {
			return nil, fmt.Errorf("utils.StringToFloat64: %w", err)
		}

		closes[i] = closeF64
	}

	return closes, nil
}
