package ta

import "github.com/shopspring/decimal"

func Subtract(values1 []decimal.Decimal, values2 []decimal.Decimal) []decimal.Decimal {
	result := make([]decimal.Decimal, len(values1))
	for i := 0; i < len(result); i++ {
		result[i] = values1[i].Sub(values2[i])
	}
	return result
}
