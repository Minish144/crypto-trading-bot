package domain

import "github.com/shopspring/decimal"

type OrderDirection string

const (
	OrderDirectionUnspecified = ""
	OrderDirectionBuy         = "buy"
	OrderDirectionSell        = "sell"
)

type OrderType string

const (
	OrderTypeUnspecified = ""
	OrderTypeMarket      = "market"
	OrderTypeLimit       = "limit"
)

const (
	OrderTypeStopLoss   = "stopLoss"
	OrderTypeStopLimit  = "stopLimit"
	OrderTypeTakeProfit = "takeProfit"
)

var stops = map[OrderType]struct{}{
	OrderTypeStopLoss:   {},
	OrderTypeStopLimit:  {},
	OrderTypeTakeProfit: {},
}

func IsStopOrderType(ot OrderType) bool {
	_, ok := stops[ot]

	return ok
}

type Order struct {
	ID        string
	Symbol    string
	Quantity  int
	Price     *decimal.Decimal
	Direction OrderDirection
	OrderType OrderType
}
