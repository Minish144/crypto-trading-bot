package ta

import (
	"github.com/minish144/crypto-trading-bot/domain"
	"github.com/shopspring/decimal"
)

func SMA(klines []*domain.Kline) decimal.Decimal {
	n := len(klines)

	sum := decimal.Zero

	for _, kline := range klines {
		sum = sum.Add(kline.Close)
	}

	return sum.Div(decimal.NewFromInt(int64(n)))
}

func SMA_N(klines []*domain.Kline, n int) decimal.Decimal {
	if l := len(klines); l < n {
		n = l
	}

	sum := decimal.Zero

	for i := 0; i < n; i++ {
		sum = sum.Add(klines[i].Close)
	}

	return sum.Div(decimal.NewFromInt(int64(n)))
}
