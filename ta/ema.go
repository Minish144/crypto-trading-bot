package ta

import "github.com/shopspring/decimal"

func EMA(values []decimal.Decimal, period int) []decimal.Decimal {
	ema := make([]decimal.Decimal, len(values))

	// Calculate the initial EMA using a simple moving average
	sma := SMA(values[0:period])
	ema[period-1] = sma

	// Calculate subsequent EMAs using the previous EMA and the current value
	multiplier := decimal.NewFromFloat(2.0 / float64(period+1))

	for i := period; i < len(values); i++ {
		ema[i] = values[i].Div(ema[i-1]).Mul(multiplier).Add(ema[i-1])
	}

	return ema
}
