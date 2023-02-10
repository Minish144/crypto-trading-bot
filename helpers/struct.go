package helpers

import (
	"context"
	"fmt"
	"time"

	"github.com/Minish144/crypto-trading-bot/clients"
	"go.uber.org/zap"
)

type Helper struct {
	c        clients.HttpClient
	baseCoin string
}

func NewHelper(c clients.HttpClient, baseCoin string) *Helper {
	return &Helper{c: c, baseCoin: baseCoin}
}

func (h *Helper) TotalHoldings(ctx context.Context) (float64, error) {
	z := zap.S().With("context", "Helper.TotalHoldings")

	assets, err := h.c.GetAssets(ctx)
	if err != nil {
		return 0, fmt.Errorf("c.GetAssets: %w", err)
	}

	balance, err := h.c.GetBalance(ctx, h.baseCoin)
	if err != nil {
		return 0, fmt.Errorf("c.GetBalance: %w", err)
	}

	for _, asset := range assets {
		if asset.Coin == h.baseCoin || asset.Coin == "BUSD" {
			continue
		}

		pair := asset.Coin + h.baseCoin

		holdings, err := h.c.GetPrice(ctx, pair)
		if err != nil {
			z.Warnw("failed to get price", "symbol", pair)
		}

		balance += holdings
	}

	return balance, nil
}

func (h *Helper) StartLoggingHelpers(ctx context.Context) {
	z := zap.S().With("context", "Helper.LoggingHelpers")

	for {
		select {
		case <-time.NewTicker(15 * time.Second).C:
			balance, err := h.TotalHoldings(ctx)
			if err != nil {
				z.Warnw("failed to calculate total holdings", "error", err.Error())
			} else {
				z.Infow("total holdings", "base_coin", h.baseCoin, "amount", balance)
			}
		case <-ctx.Done():
			return
		}
	}
}
