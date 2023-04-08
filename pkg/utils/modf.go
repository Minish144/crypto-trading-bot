package utils

import (
	"github.com/shopspring/decimal"
)

func Modf(d decimal.Decimal) (int, int) {
	i := d.IntPart()

	p := decimal.NewFromInt32(d.Exponent())
	frac := d.Sub(decimal.NewFromInt(i)).Mul(decimal.NewFromInt(10).Pow(p)).IntPart()

	return int(i), int(frac)
}
