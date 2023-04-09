package tinkoffExchange

import (
	"github.com/minish144/crypto-trading-bot/pkg/utils"
	"github.com/shopspring/decimal"
	investapi "github.com/tinkoff/invest-api-go-sdk"
)

type Instrument struct {
	FIGI              string
	ISIN              string
	Lot               decimal.Decimal
	MinPriceIncrement decimal.Decimal
	Name              string
	Ticker            string
}

func currenciesToInstruments(xs []*investapi.Currency) []*Instrument {
	instruments := make([]*Instrument, 0, len(xs))

	for _, x := range xs {
		i := &Instrument{
			FIGI:   x.Figi,
			ISIN:   x.Isin,
			Lot:    decimal.NewFromInt32(x.Lot),
			Name:   x.Name,
			Ticker: x.Ticker,
		}

		if x.MinPriceIncrement != nil {
			i.MinPriceIncrement = utils.IntFractToDecimal(x.MinPriceIncrement.Units, x.MinPriceIncrement.Nano)
		}

		instruments = append(instruments, i)
	}

	return instruments
}

func stocksToInstruments(xs []*investapi.Share) []*Instrument {
	instruments := make([]*Instrument, 0, len(xs))

	for _, x := range xs {
		i := &Instrument{
			FIGI:   x.Figi,
			ISIN:   x.Isin,
			Lot:    decimal.NewFromInt32(x.Lot),
			Name:   x.Name,
			Ticker: x.Ticker,
		}

		if x.MinPriceIncrement != nil {
			i.MinPriceIncrement = utils.IntFractToDecimal(x.MinPriceIncrement.Units, x.MinPriceIncrement.Nano)
		}

		instruments = append(instruments, i)
	}

	return instruments
}

func bondsToInstruments(xs []*investapi.Bond) []*Instrument {
	instruments := make([]*Instrument, 0, len(xs))

	for _, x := range xs {
		i := &Instrument{
			FIGI:   x.Figi,
			ISIN:   x.Isin,
			Lot:    decimal.NewFromInt32(x.Lot),
			Name:   x.Name,
			Ticker: x.Ticker,
		}

		if x.MinPriceIncrement != nil {
			i.MinPriceIncrement = utils.IntFractToDecimal(x.MinPriceIncrement.Units, x.MinPriceIncrement.Nano)
		}

		instruments = append(instruments, i)
	}

	return instruments
}

func etfsToInstruments(xs []*investapi.Etf) []*Instrument {
	instruments := make([]*Instrument, 0, len(xs))

	for _, x := range xs {
		i := &Instrument{
			FIGI:   x.Figi,
			ISIN:   x.Isin,
			Lot:    decimal.NewFromInt32(x.Lot),
			Name:   x.Name,
			Ticker: x.Ticker,
		}

		if x.MinPriceIncrement != nil {
			i.MinPriceIncrement = utils.IntFractToDecimal(x.MinPriceIncrement.Units, x.MinPriceIncrement.Nano)
		}

		instruments = append(instruments, i)
	}

	return instruments
}
