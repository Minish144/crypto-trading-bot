package utils

func QuoteQtyFromBaseQty(basePrice float64, baseAmount float64) float64 {
	return baseAmount / basePrice
}
