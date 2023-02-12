package macdStrategy

import (
	"context"
	"time"

	"github.com/MicahParks/go-ma"
	"github.com/Minish144/crypto-trading-bot/utils"
)

func (s *MACDStrategy) Start(ctx context.Context) error {
	go s.logic(ctx)

	for {
		select {
		case <-time.NewTicker(s.cfg.Interval).C:
			go s.logic(ctx)
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
		)
	}

	if l := len(klines); l < 36 {
		s.z.Warnw(
			"not enough klines received",
			"len", l,
		)

		return
	}

	signal, price, macd, signalEMA := s.getSignal(klines)
	price = utils.RoundPrecision(price, s.cfg.PricePrecision)

	s.z.Infow(
		"new signal",
		"signal", signal,
		"price", price,
		"MACD", macd,
		"signal_EMA", signalEMA,
	)

	if signal == signalBuy {
		if err := s.client.NewMarketBuyOrder(
			ctx,
			s.cfg.Symbol,
			s.cfg.OrderAmount,
		); err != nil {
			s.z.Warnw(
				"failed to place order",
				"side", "buy",
				"type", "market",
				"price", price,
				"quantity", s.cfg.OrderAmount,
				"error", err.Error(),
			)

			return
		}
	} else if signal == signalSell {
		if err := s.client.NewMarketSellOrder(ctx, s.cfg.Symbol, s.cfg.OrderAmount); err != nil {
			s.z.Warnw(
				"failed to place order",
				"side", "sell",
				"type", "market",
				"price", price,
				"quantity", s.cfg.OrderAmount,
				"error", err.Error(),
			)

			return
		}
	}
}

func (s *MACDStrategy) getSignal(klines []float64) (signal, float64, float64, float64) {
	macdSignal := ma.DefaultMACDSignal(klines[:ma.RequiredSamplesForDefaultMACDSignal])

	price := klines[len(klines)-1]
	signal := signalDoNothing

	results := macdSignal.Calculate(price)
	if results.BuySignal != nil && *results.BuySignal {
		signal = signalBuy
	} else {
		signal = signalSell
	}

	return signal, price, results.MACD.Result, results.SignalEMA

}

func (s *MACDStrategy) Stop(ctx context.Context) error {
	return nil
}
