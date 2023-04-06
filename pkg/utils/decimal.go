package utils

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func IntFractToDecimal(integer int64, fractional int32) decimal.Decimal {
	f, _ := decimal.NewFromString(fmt.Sprintf("%d.%d", integer, fractional))

	return f
}
