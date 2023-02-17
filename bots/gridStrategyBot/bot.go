package gridStrategyBot

import (
	"context"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/money"
	"github.com/thecolngroup/alphakit/trader"
)

var (
	_                 trader.Bot = (*gridStrategyBot)(nil)
	ErrNotImplemented            = errors.New("not implemented")
)

type gridStrategyBot struct {
	ctx       context.Context
	asset     market.Asset
	sizer     money.Sizer
	predicter trader.Predicter

	dealer broker.Dealer
}

func MakeGridStrategyBotFromConfig(config map[string]any) (trader.Bot, error) {
	asset := config["asset"].(market.Asset)

	sizer := NewSizer(
		config["sizerFixedCapital"].(decimal.Decimal),
		config["sizerStepSize"].(float64),
	)

	dealer := config["dealer"].(broker.Dealer)

	ctx := config["ctx"].(context.Context)

	predicter := NewPredicter(
		ctx,
		config["maxKlinesAmount"].(int32),
	)

	bot := &gridStrategyBot{
		ctx:       ctx,
		asset:     asset,
		sizer:     sizer,
		predicter: predicter,
		dealer:    dealer,
	}

	return bot, nil
}

func (b *gridStrategyBot) Warmup(ctx context.Context, klines []market.Kline) error {
	if !b.predicter.Valid() {
		return fmt.Errorf("predicter valid: false")
	}

	for i := 0; i < len(klines); i++ {
		if err := b.predicter.ReceivePrice(ctx, klines[i]); err != nil {
			return fmt.Errorf("predicter receiver price: %w", err)
		}
	}

	return nil
}

func (b *gridStrategyBot) SetDealer(d broker.Dealer) {
	b.dealer = d
}

func (b *gridStrategyBot) SetAsset(a market.Asset) {
	b.asset = a
}

func (b *gridStrategyBot) Close(ctx context.Context) error {
	return ErrNotImplemented
}

func (b *gridStrategyBot) ReceivePrice(ctx context.Context, klines market.Kline) error {
	if !b.predicter.Valid() {
		return fmt.Errorf("predicter valid: false")
	}

	if err := b.predicter.ReceivePrice(ctx, klines); err != nil {
		return fmt.Errorf("predicter receiver price: %w", err)
	}

	return nil
}
