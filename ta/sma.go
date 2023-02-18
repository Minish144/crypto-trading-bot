package ta

import "github.com/shopspring/decimal"

func SMA(values []decimal.Decimal) decimal.Decimal {
	var sum decimal.Decimal

	for _, value := range values {
		sum = sum.Add(value)
	}

	l := decimal.NewFromInt(int64(len(values)))

	return sum.Div(l)
}
