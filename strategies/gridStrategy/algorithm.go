package gridStrategy

import (
	"context"
	"time"

	"go.uber.org/atomic"
)

func (s *GridStrategy) Start(ctx context.Context) error {
	stopLoss := atomic.NewFloat64(0)

	go func() {
		s.checkBalances(ctx)
		s.logic(ctx, stopLoss)

		for {
			select {
			case <-time.NewTicker(1 * time.Minute).C:
				s.checkBalances(ctx)
				s.logic(ctx, stopLoss)
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		select {
		case <-time.NewTicker(s.cfg.Interval).C:
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *GridStrategy) logic(ctx context.Context, stopLoss *atomic.Float64) {
	// get the current price of the symbol
	price, err := s.client.GetPrice(ctx, s.cfg.Symbol)
	if err != nil {
		s.z.Warnw("failed to get price", "error", err.Error())
		return
	}

	go s.stopLoss(ctx, stopLoss.Load(), price)

	sellLevels, buyLevels, stopLossValue := s.generateGridsAndStopLoss(price)
	s.placeSellOrders(ctx, sellLevels, price, s.cfg.OrderAmount)
	s.placeBuyOrders(ctx, buyLevels, price, s.cfg.OrderAmount)

	stopLoss.Store(stopLossValue)
	s.z.Infow(
		"stop loss updated",
		"amount", stopLoss,
	)
}

func (s *GridStrategy) generateGridsAndStopLoss(price float64) ([]float64, []float64, float64) {
	// calculate the price levels for each grid
	sellLevels := make([]float64, s.cfg.GridsAmount)
	buyLevels := make([]float64, s.cfg.GridsAmount)

	for i := 1; i <= s.cfg.GridsAmount; i++ {
		sellLevels[i-1] = price * (1 + (float64(i) * s.cfg.GridStep))
		buyLevels[i-1] = price * (1 - (float64(i) * s.cfg.GridStep))
	}

	stopLoss := price * (1 - (1+float64(s.cfg.GridsAmount))*s.cfg.GridSize)

	return sellLevels, buyLevels, stopLoss
}

func (s *GridStrategy) checkBalances(ctx context.Context) {
	// check the balance of the account
	balance, err := s.client.GetBalance(ctx, s.cfg.Coins.Base)
	if err != nil {
		s.z.Warnw("failed to get balance", "error", err.Error())
		return
	}

	s.z.Infow(
		"balance",
		"coin", s.cfg.Coins.Base,
		"amount", balance,
	)

	balanceQuote, err := s.client.GetBalance(ctx, s.cfg.Coins.Quote)
	if err != nil {
		s.z.Warnw("failed to get balance", "error", err.Error())
		return
	}

	price, err := s.client.GetPrice(ctx, s.cfg.Symbol)
	if err != nil {
		s.z.Warnw("failed to get price", "error", err.Error())
		return
	}

	s.z.Infow(
		"balance",
		"coin", s.cfg.Coins.Quote,
		"amount", balanceQuote,
		"base_equivalent", price*balanceQuote,
	)
}

func (s *GridStrategy) stopLoss(ctx context.Context, stopLoss float64, currentPrice float64) {
	if stopLoss != 0 && stopLoss >= currentPrice {
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

		go func() {
			balance, err := s.client.GetBalance(ctx, s.cfg.Coins.Quote)
			if err != nil {
				s.z.Warnw(
					"failed to get balance for stoploss",
					"coin", s.cfg.Coins.Quote,
					"error", err.Error(),
				)
			}

			s.client.NewLimitSellOrder(ctx, s.cfg.Symbol, currentPrice, balance)
		}()

		for _, order := range orders {
			if err := s.client.CloseOrder(ctx, s.cfg.Symbol, order.OrderID); err != nil {
				s.z.Warnw(
					"failed to close order after stop loss trigger",
					"side", "buy",
					"type", "limit",
					"price", order.Price,
					"quantity", order.OrigQuantity,
					"error", err.Error(),
				)
			}
		}
	}
}

func (s *GridStrategy) placeSellOrders(ctx context.Context, levels []float64, price, quantity float64) {
	// placing sell orders at each grid level
	for i, gridLevel := range levels {
		// calculate the size of the order
		if err := s.client.NewLimitSellOrder(ctx, s.cfg.Symbol, gridLevel, quantity); err != nil {
			s.z.Warnw(
				"failed to place order",
				"side", "sell",
				"type", "limit",
				"multiplier", (1 + (float64(i) * s.cfg.GridStep)),
				"price", gridLevel,
				"quantity", quantity,
				"error", err.Error(),
			)

			continue
		}

		s.z.Infow(
			"new order",
			"side", "sell",
			"type", "limit",
			"multiplier", (1 + (float64(i) * s.cfg.GridStep)),
			"price", gridLevel,
			"quantity", quantity,
		)
	}
}

func (s *GridStrategy) placeBuyOrders(ctx context.Context, levels []float64, price, quantity float64) {
	// placing buy orders at each grid level
	for i, gridLevel := range levels {
		// calculate the size of the order
		if err := s.client.NewLimitBuyOrder(ctx, s.cfg.Symbol, gridLevel, quantity); err != nil {
			s.z.Warnw(
				"failed to place order",
				"side", "buy",
				"type", "limit",
				"multiplier", (1 - (float64(i) * s.cfg.GridStep)),
				"price", gridLevel,
				"quantity", quantity,
				"error", err.Error(),
			)

			continue
		}

		s.z.Infow(
			"new order",
			"side", "buy",
			"type", "limit",
			"multiplier", (1 - (float64(i) * s.cfg.GridStep)),
			"price", gridLevel,
			"quantity", quantity,
		)
	}
}

func (s *GridStrategy) Stop(ctx context.Context) error {
	return nil
}
