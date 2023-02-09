package gridStrategy

import (
	"context"
	"time"
)

func (s *GridStrategy) Start(ctx context.Context) error {
	for {
		select {
		case <-time.NewTicker(s.cfg.Interval).C:
			// get the current price of the symbol
			price, err := s.client.GetPrice(ctx, s.cfg.Symbol)
			if err != nil {
				s.z.Warnw("failed to get price", "error", err.Error())
				continue
			}

			s.z.Infow(
				"price",
				"amount", price,
			)

			// check the balance of the account
			balance, err := s.client.GetBalance(ctx, s.cfg.Coins.Base)
			if err != nil {
				s.z.Warnw("failed to get balance", "error", err.Error())
				continue
			}

			s.z.Infow(
				"balance",
				"base", s.cfg.Coins.Base,
				"amount", balance,
			)

			sellLevels, buyLevels := s.generateGrids(price)

			s.placeSellOrders(ctx, sellLevels, price, balance)
			s.placeBuyOrders(ctx, buyLevels, price, balance)
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *GridStrategy) generateGrids(price float64) ([]float64, []float64) {
	// calculate the price levels for each grid
	sellLevels := make([]float64, s.cfg.GridsAmount)
	buyLevels := make([]float64, s.cfg.GridsAmount)

	for i := 1; i <= s.cfg.GridsAmount; i++ {
		sellLevels[i-1] = price * (1 + (float64(i) * s.cfg.GridStep))
		buyLevels[i-1] = price * (1 - (float64(i) * s.cfg.GridStep))
	}

	return sellLevels, buyLevels
}

func (s *GridStrategy) placeSellOrders(ctx context.Context, levels []float64, price, balance float64) {
	// placing sell orders at each grid level
	for i, gridLevel := range levels {
		// calculate the size of the order
		orderSizeBase := balance * s.cfg.GridSize * float64(i+1)
		orderSizeQuote := orderSizeBase / price

		if err := s.client.NewLimitSellOrder(ctx, s.cfg.Symbol, gridLevel, orderSizeQuote); err != nil {
			s.z.Warnw(
				"failed to place order",
				"side", "sell",
				"type", "limit",
				"multiplier", (1 + (float64(i) * s.cfg.GridStep)),
				"price", gridLevel,
				"quantity", orderSizeQuote,
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
			"quantity", orderSizeQuote,
		)
	}
}

func (s *GridStrategy) placeBuyOrders(ctx context.Context, levels []float64, price, balance float64) {
	// placing buy orders at each grid level
	for i, gridLevel := range levels {
		// calculate the size of the order
		orderSizeBase := balance * s.cfg.GridSize * float64(i+1)
		orderSizeQuote := orderSizeBase / price

		if err := s.client.NewLimitBuyOrder(ctx, s.cfg.Symbol, gridLevel, orderSizeQuote); err != nil {
			s.z.Warnw(
				"failed to place order",
				"side", "buy",
				"type", "limit",
				"multiplier", (1 - (float64(i) * s.cfg.GridStep)),
				"price", gridLevel,
				"quantity", orderSizeQuote,
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
			"quantity", orderSizeQuote,
		)
	}
}

func (s *GridStrategy) Stop(ctx context.Context) error {
	return nil
}
