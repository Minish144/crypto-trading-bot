package macdStrategy

import (
	"context"
	"time"

	"github.com/MicahParks/go-ma"
	"github.com/Minish144/crypto-trading-bot/models"
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

func (s *MACDStrategy) logic(ctx context.Context) {
	klines, err := s.client.GetKlinesCloses(ctx, s.cfg.Symbol, s.cfg.KlinesInterval)
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

	signal := ma.DefaultMACDSignal(klines[:ma.RequiredSamplesForDefaultMACDSignal])

	// Iterate through the rest of the data and print the results.
	var results ma.MACDSignalResults
	for i, p := range klines[ma.RequiredSamplesForDefaultMACDSignal:] {
		results = signal.Calculate(p)

		// Interpret the buy signal.
		var buySignal string
		if results.BuySignal != nil {
			if results.BuySignal != nil && *results.BuySignal {
				buySignal = string(models.SideTypeBuy)
			} else {
				buySignal = string(models.SideTypeSell)
			}
		} else {
			continue
		}

		s.z.Infow(
			"new signal",
			"price", p,
			"price_index", i+ma.RequiredSamplesForDefaultMACDSignal,
			"MACD", results.MACD.Result,
			"signal_EMA", results.SignalEMA,
			"buy_signal", buySignal,
		)
	}
}

func (s *MACDStrategy) Stop(ctx context.Context) error {
	return nil
}
