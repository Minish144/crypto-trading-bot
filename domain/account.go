package domain

import "github.com/shopspring/decimal"

type InstrumentType string

const (
	InstrumentTypeStock    InstrumentType = "Stock"
	InstrumentTypeCurrency InstrumentType = "Currency"
	InstrumentTypeBond     InstrumentType = "Bond"
	InstrumentTypeEtf      InstrumentType = "Etf"
	InstrumentTypeCrypto   InstrumentType = "Crypto"
)

type Balance struct {
	Symbol string
	Type   InstrumentType
	Lots   decimal.Decimal
	Free   decimal.Decimal
	Locked decimal.Decimal
}

type Account struct {
	Balances []Balance
}
