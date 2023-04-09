package tinkoffExchange

import (
	"fmt"

	"github.com/minish144/crypto-trading-bot/domain"
	"github.com/minish144/crypto-trading-bot/pkg/tinkoff"
	"github.com/minish144/crypto-trading-bot/pkg/utils"
	"github.com/shopspring/decimal"
	investapi "github.com/tinkoff/invest-api-go-sdk"
)

type Instrument struct {
	FIGI              string
	ISIN              string
	Ticker            string
	Type              domain.InstrumentType
	Lot               decimal.Decimal
	MinPriceIncrement decimal.Decimal
	Name              string
}

func loadInstruments(client *tinkoff.TinkoffAPI) ([]*Instrument, error) {
	respCurrencies, err := client.InstrumentsClient.Currencies(
		client.Ctx,
		&investapi.InstrumentsRequest{InstrumentStatus: investapi.InstrumentStatus_INSTRUMENT_STATUS_BASE},
	)
	if err != nil {
		return nil, fmt.Errorf("InstrumentsClient.Currencies: %w", err)
	}

	respStocks, err := client.InstrumentsClient.Shares(
		client.Ctx,
		&investapi.InstrumentsRequest{InstrumentStatus: investapi.InstrumentStatus_INSTRUMENT_STATUS_BASE},
	)
	if err != nil {
		return nil, fmt.Errorf("InstrumentsClient.Shares: %w", err)
	}

	respBonds, err := client.InstrumentsClient.Bonds(
		client.Ctx,
		&investapi.InstrumentsRequest{InstrumentStatus: investapi.InstrumentStatus_INSTRUMENT_STATUS_BASE},
	)
	if err != nil {
		return nil, fmt.Errorf("InstrumentsClient.Bonds: %w", err)
	}

	respEtfs, err := client.InstrumentsClient.Etfs(
		client.Ctx,
		&investapi.InstrumentsRequest{InstrumentStatus: investapi.InstrumentStatus_INSTRUMENT_STATUS_BASE},
	)
	if err != nil {
		return nil, fmt.Errorf("InstrumentsClient.Etfs: %w", err)
	}

	instruments := make([]*Instrument, 0,
		len(respCurrencies.Instruments)+
			len(respStocks.Instruments)+
			len(respBonds.Instruments)+
			len(respEtfs.Instruments),
	)

	instruments = append(instruments, currenciesToInstruments(respCurrencies.Instruments)...)
	instruments = append(instruments, stocksToInstruments(respStocks.Instruments)...)
	instruments = append(instruments, bondsToInstruments(respBonds.Instruments)...)
	instruments = append(instruments, etfsToInstruments(respEtfs.Instruments)...)

	return instruments, nil
}

func currenciesToInstruments(xs []*investapi.Currency) []*Instrument {
	instruments := make([]*Instrument, 0, len(xs))

	for _, x := range xs {
		i := &Instrument{
			FIGI:   x.Figi,
			ISIN:   x.Isin,
			Ticker: x.Ticker,
			Type:   domain.InstrumentTypeCurrency,
			Lot:    decimal.NewFromInt32(x.Lot),
			Name:   x.Name,
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
			Ticker: x.Ticker,
			Type:   domain.InstrumentTypeStock,
			Lot:    decimal.NewFromInt32(x.Lot),
			Name:   x.Name,
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
			Ticker: x.Ticker,
			Type:   domain.InstrumentTypeBond,
			Lot:    decimal.NewFromInt32(x.Lot),
			Name:   x.Name,
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
			Ticker: x.Ticker,
			Type:   domain.InstrumentTypeEtf,
			Lot:    decimal.NewFromInt32(x.Lot),
			Name:   x.Name,
		}

		if x.MinPriceIncrement != nil {
			i.MinPriceIncrement = utils.IntFractToDecimal(x.MinPriceIncrement.Units, x.MinPriceIncrement.Nano)
		}

		instruments = append(instruments, i)
	}

	return instruments
}
