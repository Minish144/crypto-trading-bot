package macdStrategy

import (
	"context"
	"fmt"
	"time"

	"github.com/Minish144/crypto-trading-bot/models"
	"github.com/Minish144/crypto-trading-bot/utils"
	"github.com/cinar/indicator"
	"go.uber.org/atomic"
)

func (s *MACDStrategy) Start(ctx context.Context) error {
	stopLoss := atomic.NewFloat64(0)

	go s.stopLoss(ctx, stopLoss)
	go s.logic(ctx)

	for {
		select {
		case <-time.NewTicker(s.cfg.Interval).C:
			go s.logic(ctx)
		case <-time.NewTicker(s.cfg.StopLossUpdatePeriod).C:
			go s.stopLoss(ctx, stopLoss)
		case <-ctx.Done():
			return nil
		}
	}
}

type signal string

const (
	signalBuy       signal = "buy"
	signalSell      signal = "sell"
	signalDoNothing signal = "do-nothing"
)

func (s *MACDStrategy) logic(ctx context.Context) {
	klines, err := s.client.GetKlinesCloses(
		ctx,
		s.cfg.Symbol,
		s.cfg.KlinesInterval,
	)
	if err != nil {
		s.z.Warnw(
			"failed to get klines",
			"interval", s.cfg.KlinesInterval,
			"error", err.Error(),
		)
	}

	if l := len(klines); l < 36 {
		s.z.Warnw(
			"not enough klines received",
			"len", l,
		)

		return
	}

	signal, macd := s.getSignal(klines)
	price := utils.RoundPrecision(klines[len(klines)-1], s.cfg.PricePrecision)

	amount := utils.RoundPrecision(s.cfg.OrderAmount, s.cfg.QuantityPrecision)
	if s.cfg.BaseCoinForAmount {
		amount = utils.RoundPrecision(utils.QuoteQtyFromBaseQty(price, s.cfg.OrderAmount), s.cfg.QuantityPrecision)
	}

	s.z.Infow(
		"new signal",
		"signal", signal,
		"price", price,
		"MACD", macd,
	)

	if signal == signalBuy {
		allowed, err := s.isBuyAllowedByLimit(ctx, price, amount)
		if err != nil {
			s.z.Warnw(
				"failed to check allowance",
				"side", "buy",
				"type", "market",
				"price", price,
				"quantity", amount,
				"error", err.Error(),
			)

			return
		}

		if !allowed {
			s.z.Warnw(
				"failed to place order",
				"side", "buy",
				"type", "market",
				"price", price,
				"quantity", amount,
				"error", "buy limit had been reached",
			)
		}

		if err := s.client.NewMarketBuyOrder(ctx, s.cfg.Symbol, amount); err != nil {
			s.z.Warnw(
				"failed to place order",
				"side", "buy",
				"type", "market",
				"price", price,
				"quantity", amount,
				"error", err.Error(),
			)

			return
		}
	} else if signal == signalSell {
		if err := s.client.NewMarketSellOrder(ctx, s.cfg.Symbol, amount); err != nil {
			s.z.Warnw(
				"failed to place order",
				"side", "sell",
				"type", "market",
				"price", price,
				"quantity", amount,
				"error", err.Error(),
			)

			return
		}
	}
}

func (s *MACDStrategy) getSignal(klines []float64) (signal, float64) {
	macdValues, signalValues := indicator.Macd(klines)

	latestIndex := len(macdValues) - 1

	macd := macdValues[latestIndex]
	signalLatest := signalValues[latestIndex]

	if macd > signalLatest {
		return signalBuy, macd
	} else if macd < signalLatest {
		return signalSell, macd
	} else {
		return signalDoNothing, macd
	}
}

func (s *MACDStrategy) isBuyAllowedByLimit(ctx context.Context, price, qty float64) (bool, error) {
	available := s.cfg.MaxOrdersAmount

	orders, err := s.client.GetOpenOrders(ctx, s.cfg.Symbol)
	if err != nil {
		return false, fmt.Errorf("client.GetOpenOrders: %w", err)
	}

	for _, order := range orders {
		if order.Side == models.SideTypeBuy {
			available -= order.Price * order.OrigQuantity
		}
	}

	available -= price * qty

	balance, locked, err := s.client.GetBalance(ctx, s.cfg.Coins.Quote)
	if err != nil {
		return false, fmt.Errorf("client.GetBalance: %w", err)
	}

	total := price * (balance + locked)

	return available >= 0 && total <= s.cfg.MaxOrdersAmount, nil
}

// @TODO implement using stop-loss binance orders instead of market sell order
func (s *MACDStrategy) stopLoss(ctx context.Context, stopLoss *atomic.Float64) {
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

	stopLossActual := currentPrice * s.cfg.StopLossShare

	s.z.Infow(
		"stop loss updated",
		"previous", stopLoss.Load(),
		"current", stopLossActual,
	)

	stopLoss.Store(stopLossActual)
}

func (s *MACDStrategy) closeOrders(ctx context.Context, orders []*models.Order) {
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

func (s *MACDStrategy) Stop(ctx context.Context) error {
	return nil
}
