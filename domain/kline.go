package domain

import "github.com/shopspring/decimal"

type Kline struct {
	Low    decimal.Decimal
	High   decimal.Decimal
	Open   decimal.Decimal
	Close  decimal.Decimal
	Volume decimal.Decimal
}
