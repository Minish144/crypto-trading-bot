package gridStrategy

import (
	"context"
	"fmt"
	"time"

	"github.com/Minish144/crypto-trading-bot/models"
	"github.com/Minish144/crypto-trading-bot/utils"
	"go.uber.org/atomic"
)

func (s *GridStrategy) Start(ctx context.Context) error {
	stopLoss := atomic.NewFloat64(0)
	ordersChecksCounter := atomic.NewInt32(0)

	go s.stopLoss(ctx, stopLoss)
	// go s.checkBalances(ctx)
	go s.logic(ctx, stopLoss, ordersChecksCounter)

	for {
		select {
		case <-time.NewTicker(s.cfg.Interval).C:
			// go s.checkBalances(ctx)
			go s.logic(ctx, stopLoss, ordersChecksCounter)
		case <-time.NewTicker(s.cfg.StopLossUpdatePeriod).C:
			go s.stopLoss(ctx, stopLoss)
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *GridStrategy) logic(ctx context.Context, stopLoss *atomic.Float64, ordersChecksCounter *atomic.Int32) {
	// get the current price of the symbol
	price, err := s.client.GetPrice(ctx, s.cfg.Symbol)
	if err != nil {
		s.z.Warnw("failed to get price", "error", err.Error())
		return
	}

	amount := s.cfg.OrderAmount
	if s.cfg.BaseCoinForAmount {
		amount = utils.QuoteQtyFromBaseQty(price, s.cfg.OrderAmount)
	}

	// increment checks counter
	ordersChecksCounter.Inc()

	// close all orders if more checks were performed than expected
	if ordersChecksCounter.Load() >= int32(s.cfg.OrdersCheckRetriesMax) {
		orders, err := s.client.GetOpenOrders(ctx, s.cfg.Symbol)
		if err != nil {
			s.z.Warnw(
				"failed to get open orders for stop loss",
				"error", err.Error(),
			)

			return
		}

		s.closeOrders(ctx, orders)
	}

	// check whether all orders have been filled
	filled, err := s.allOrdersFilled(ctx)
	if err != nil {
		s.z.Warnw("failed to check open orders", "error", err.Error())
		return
	} else if !filled {
		return
	}

	// place buy and sell orders
	sellLevels, buyLevels := s.generateGridsAndStopLoss(price)
	s.placeSellOrders(ctx, sellLevels, price, amount)
	s.placeBuyOrders(ctx, buyLevels, price, amount)

	ordersChecksCounter.Store(0)
}

func (s *GridStrategy) generateGridsAndStopLoss(price float64) ([]float64, []float64) {
	// calculate the price levels for each grid
	sellLevels := make([]float64, s.cfg.GridsAmount)
	buyLevels := make([]float64, s.cfg.GridsAmount)

	for i := uint(1); i <= s.cfg.GridsAmount; i++ {
		sellLevels[i-1] = price * (1 + s.cfg.GridSize + (float64(i) * s.cfg.GridSize))
		buyLevels[i-1] = price * (1 - (float64(i) * s.cfg.GridSize))
	}

	return sellLevels, buyLevels
}

func (s *GridStrategy) allOrdersFilled(ctx context.Context) (bool, error) {
	orders, err := s.client.GetOpenOrders(ctx, s.cfg.Symbol)
	if err != nil {
		return false, fmt.Errorf("client.GetOpenOrders: %w", err)
	}

	return len(orders) == 0, nil
}

func (s *GridStrategy) closeOrders(ctx context.Context, orders []*models.Order) {
	for _, order := range orders {
		if err := s.client.CloseOrder(ctx, s.cfg.Symbol, order.OrderID); err != nil {
			s.z.Warnw(
				"failed to close order",
				"side", "buy",
				"type", "limit",
				"price", order.Price,
				"quantity", order.OrigQuantity,
				"error", err.Error(),
			)
		}
	}
}

// @TODO implement using stop-loss binance orders instead of market sell order
func (s *GridStrategy) stopLoss(ctx context.Context, stopLoss *atomic.Float64) {
	currentPrice, err := s.client.GetPrice(ctx, s.cfg.Symbol)
	if err != nil {
		s.z.Warnw(
			"failed to get price for stop loss",
			"error", err.Error(),
		)
		return
	}

	if stopLoss.Load() != 0 && stopLoss.Load() >= currentPrice {
		s.z.Infow(
			"stop loss triggered",
			"current_price", currentPrice,
			"stop_loss", stopLoss,
		)

		orders, err := s.client.GetOpenOrders(ctx, s.cfg.Symbol)
		if err != nil {
			s.z.Warnw(
				"failed to get open orders for stop loss",
				"error", err.Error(),
			)

			return
		}

		s.closeOrders(ctx, orders)

		go func() {
			balance, _, err := s.client.GetBalance(ctx, s.cfg.Coins.Quote)
			if err != nil {
				s.z.Warnw(
					"failed to get balance for stop loss",
					"coin", s.cfg.Coins.Quote,
					"error", err.Error(),
				)
			}

			s.client.NewLimitSellOrder(ctx, s.cfg.Symbol, currentPrice, balance)
		}()
	}

	stopLossActual := currentPrice * (1 - (1+float64(s.cfg.GridsAmount))*s.cfg.GridSize)
	s.z.Infow(
		"stop loss updated",
		"previous", stopLoss.Load(),
		"current", stopLossActual,
	)

	stopLoss.Store(stopLossActual)
}

func (s *GridStrategy) placeSellOrders(ctx context.Context, levels []float64, price, quantity float64) {
	// placing sell orders at each grid level
	for i, gridLevel := range levels {
		// calculate the size of the order
		if err := s.client.NewLimitSellOrder(
			ctx,
			s.cfg.Symbol,
			utils.RoundPrecision(gridLevel, s.cfg.PricePrecision),
			utils.RoundPrecision(quantity*(1+float64(i)*s.cfg.GridStep), s.cfg.QuantityPrecision),
		); err != nil {
			s.z.Warnw(
				"failed to place order",
				"side", "sell",
				"type", "limit",
				"multiplier", (1 + s.cfg.GridSize + (float64(i) * s.cfg.GridSize)),
				"price", utils.RoundPrecision(gridLevel, s.cfg.PricePrecision),
				"quantity", utils.RoundPrecision(quantity*(1+float64(i)*s.cfg.GridStep), s.cfg.QuantityPrecision),
				"error", err.Error(),
			)

			continue
		}

		s.z.Infow(
			"new order",
			"side", "sell",
			"type", "limit",
			"multiplier", (1 + s.cfg.GridSize + (float64(i) * s.cfg.GridSize)),
			"price", utils.RoundPrecision(gridLevel, s.cfg.PricePrecision),
			"quantity", utils.RoundPrecision(quantity*(1+float64(i)*s.cfg.GridStep), s.cfg.QuantityPrecision),
		)
	}
}

func (s *GridStrategy) placeBuyOrders(ctx context.Context, levels []float64, price, quantity float64) {
	// placing buy orders at each grid level
	for i, gridLevel := range levels {
		// calculate the size of the order
		if err := s.client.NewLimitBuyOrder(
			ctx,
			s.cfg.Symbol,
			utils.RoundPrecision(gridLevel, s.cfg.PricePrecision),
			quantity*(1+float64(i)*s.cfg.GridStep),
		); err != nil {
			s.z.Warnw(
				"failed to place order",
				"side", "buy",
				"type", "limit",
				"multiplier", (1 - (float64(i) * s.cfg.GridSize)),
				"price", utils.RoundPrecision(gridLevel, s.cfg.PricePrecision),
				"quantity", utils.RoundPrecision(quantity*(1+float64(i)*s.cfg.GridStep), s.cfg.QuantityPrecision),
				"error", err.Error(),
			)

			continue
		}

		s.z.Infow(
			"new order",
			"side", "buy",
			"type", "limit",
			"multiplier", (1 - (float64(i) * s.cfg.GridSize)),
			"price", utils.RoundPrecision(gridLevel, s.cfg.PricePrecision),
			"quantity", utils.RoundPrecision(quantity*(1+float64(i)*s.cfg.GridStep), s.cfg.QuantityPrecision),
		)
	}
}

func (s *GridStrategy) Stop(ctx context.Context) error {
	return nil
}
