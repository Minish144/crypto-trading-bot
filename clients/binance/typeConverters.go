package binance

import (
	"fmt"

	"github.com/Minish144/crypto-trading-bot/models"
	"github.com/Minish144/crypto-trading-bot/utils"
	gobinance "github.com/adshao/go-binance/v2"
)

func OrdersToModel(o *gobinance.Order) (*models.Order, error) {
	price, err := utils.StringToFloat64(o.Price)
	if err != nil {
		return nil, fmt.Errorf("failed to StringToFloat64 price: %w", err)
	}

	origQuantity, err := utils.StringToFloat64(o.OrigQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to StringToFloat64 origQuantity: %w", err)
	}

	executedQuantity, err := utils.StringToFloat64(o.ExecutedQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to StringToFloat64 executedQuantity: %w", err)
	}

	cummulativeQuoteQuantity, err := utils.StringToFloat64(o.CummulativeQuoteQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to StringToFloat64 cummulativeQuoteQuantity: %w", err)
	}

	status, ok := ostToModels[o.Status]
	if !ok {
		return nil, fmt.Errorf("failed to map status to model")
	}

	tif, ok := tiftToModels[o.TimeInForce]
	if !ok {
		return nil, fmt.Errorf("failed to map timeInForce to model")
	}

	ot, ok := otToModels[o.Type]
	if !ok {
		return nil, fmt.Errorf("failed to map orderType to model")
	}

	side, ok := stToModels[o.Side]
	if !ok {
		return nil, fmt.Errorf("failed to map sideType to model")
	}

	stopPrice, err := utils.StringToFloat64(o.StopPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to StringToFloat64 stopPrice: %w", err)
	}

	iceberg, err := utils.StringToFloat64(o.IcebergQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to StringToFloat64 icebergQuantity: %w", err)
	}

	oModel := models.Order{
		Symbol:                   o.Symbol,
		OrderID:                  o.OrderID,
		OrderListId:              o.OrderListId,
		ClientOrderID:            o.ClientOrderID,
		Price:                    price,
		OrigQuantity:             origQuantity,
		ExecutedQuantity:         executedQuantity,
		CummulativeQuoteQuantity: cummulativeQuoteQuantity,
		Status:                   status,
		TimeInForce:              tif,
		Type:                     ot,
		Side:                     side,
		StopPrice:                stopPrice,
		IcebergQuantity:          iceberg,
		Time:                     o.Time,
		UpdateTime:               o.UpdateTime,
		IsWorking:                o.IsWorking,
	}

	return &oModel, nil
}

func OrdersFromModel(o *models.Order) (*gobinance.Order, error) {
	status, ok := ostFromModels[o.Status]
	if !ok {
		return nil, fmt.Errorf("failed to map status from model")
	}

	tif, ok := tiftFromModels[o.TimeInForce]
	if !ok {
		return nil, fmt.Errorf("failed to map timeInForce from model")
	}

	ot, ok := otFromModels[o.Type]
	if !ok {
		return nil, fmt.Errorf("failed to map orderType from model")
	}

	side, ok := stFromModels[o.Side]
	if !ok {
		return nil, fmt.Errorf("failed to map sideType from model")
	}

	oBinanceModel := gobinance.Order{
		Symbol:                   o.Symbol,
		OrderID:                  o.OrderID,
		OrderListId:              o.OrderListId,
		ClientOrderID:            o.ClientOrderID,
		Price:                    utils.Float64ToString(o.Price, pricePrecision),
		OrigQuantity:             utils.Float64ToString(o.OrigQuantity, quantityPrecision),
		ExecutedQuantity:         utils.Float64ToString(o.ExecutedQuantity, quantityPrecision),
		CummulativeQuoteQuantity: utils.Float64ToString(o.CummulativeQuoteQuantity, quantityPrecision),
		Status:                   status,
		TimeInForce:              tif,
		Type:                     ot,
		Side:                     side,
		StopPrice:                utils.Float64ToString(o.StopPrice, pricePrecision),
		IcebergQuantity:          utils.Float64ToString(o.IcebergQuantity, quantityPrecision),
		Time:                     o.Time,
		UpdateTime:               o.UpdateTime,
		IsWorking:                o.IsWorking,
	}

	return &oBinanceModel, nil
}
