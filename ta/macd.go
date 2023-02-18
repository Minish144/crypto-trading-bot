package ta

import "github.com/shopspring/decimal"

func MACD(closes []decimal.Decimal, fastPeriod int, slowPeriod int, signalPeriod int) decimal.Decimal {
	// Calculate the exponential moving averages (EMAs) of the closes
	fastEma := EMA(closes, fastPeriod)
	slowEma := EMA(closes, slowPeriod)

	// Calculate the MACD line by subtracting the slow EMA from the fast EMA
	macd := Subtract(fastEma, slowEma)

	// Calculate the signal line by taking the EMA of the MACD line
	signal := EMA(macd, signalPeriod)

	// Calculate the histogram by subtracting the signal line from the MACD line
	histogram := Subtract(macd, signal)

	// Return the last value of the histogram
	return histogram[len(histogram)-1]
}
