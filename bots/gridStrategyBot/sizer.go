package gridStrategyBot

import (
	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/money"
	"github.com/thecolngroup/gou/dec"
	"github.com/thecolngroup/gou/num"
)

var _ money.Sizer = (*Sizer)(nil)

type Sizer struct {
	FixedCapital decimal.Decimal
	StepSize     float64
}

func NewSizer(fixedCapital decimal.Decimal, stepSize float64) *Sizer {
	return &Sizer{FixedCapital: fixedCapital, StepSize: stepSize}
}

func (s *Sizer) Size(price, _, _ decimal.Decimal) decimal.Decimal {
	size := num.RoundTo(s.FixedCapital.Div(price).InexactFloat64(), s.StepSize)
	return dec.New(num.NN(size, 0))
}
