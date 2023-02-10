package models

type Order struct {
	Symbol                   string
	OrderID                  int64
	OrderListId              int64
	ClientOrderID            string
	Price                    float64
	OrigQuantity             float64
	ExecutedQuantity         float64
	CummulativeQuoteQuantity float64
	Status                   OrderStatusType
	TimeInForce              TimeInForceType
	Type                     OrderType
	Side                     SideType
	StopPrice                float64
	IcebergQuantity          float64
	Time                     int64
	UpdateTime               int64
	IsWorking                bool
}
