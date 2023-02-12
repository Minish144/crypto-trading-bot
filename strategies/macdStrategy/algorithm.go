package macdStrategy

import (
	"context"
	"fmt"
	"time"

	"github.com/MicahParks/go-ma"
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

type signal int

const (
	signalBuy       signal = 1
	signalSell      signal = 2
	signalDoNothing signal = 3
)

func (s *MACDStrategy) logic(ctx context.Context) {
	t1 := time.Now()

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

	signal, price, macd, signalEMA := s.getSignal(klines)

	if signal == signalBuy {
		s.z.Infow(
			"new signal",
			"signal", "buy",
			"price", price,
			"MACD", macd,
			"signal_EMA", signalEMA,
		)
	} else if signal == signalSell {
		s.z.Infow(
			"new signal",
			"signal", "sell",
			"price", price,
			"price", price,
			"MACD", macd,
			"signal_EMA", signalEMA,
		)
	}

	s.z.Infow(
		"new signal",
		"signal", "do nothing",
		"price", price,
		"price", price,
		"MACD", macd,
		"signal_EMA", signalEMA,
	)

	t2 := time.Now()

	fmt.Println("execution time", t2.Sub(t1).Round(time.Millisecond))
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
