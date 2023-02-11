package macdStrategy

func EMA12(candles [][]float64) []float64 {
	var ema12 []float64

	// Calculate EMA12
	for i, candle := range candles {
		if i == 0 {
			ema12 = append(ema12, candle[4])
		} else {
			ema12 = append(ema12, (candle[4]-ema12[i-1])*(2/13)+ema12[i-1])
		}
	}

	return ema12
}

func EMA26(candles [][]float64) []float64 {
	var ema26 []float64

	// Calculate EMA26
	for i, candle := range candles {
		if i == 0 {
			ema26 = append(ema26, candle[4])
		} else {
			ema26 = append(ema26, (candle[4]-ema26[i-1])*(2/27)+ema26[i-1])
		}
	}

	return ema26
}

func MACD(candles [][]float64) []float64 {
	var ema12, ema26, macdValues []float64

	ema12 = EMA12(candles)
	ema26 = EMA26(candles)

	// Calculate MACD
	for i := range candles {
		macdValues = append(macdValues, ema12[i]-ema26[i])
	}

	return macdValues
}
