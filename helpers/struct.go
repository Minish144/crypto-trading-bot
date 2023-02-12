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

func (h *Helper) TotalHoldings(ctx context.Context) (float64, float64, error) {
	assets, err := h.c.GetAssets(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("c.GetAssets: %w", err)
	}

	balance, locked, err := h.c.GetBalance(ctx, h.baseCoin)
	if err != nil {
		return 0, 0, fmt.Errorf("c.GetBalance: %w", err)
	}

	for _, asset := range assets {
		pair := asset.Coin + h.baseCoin

		if asset.Free == 0 && asset.Locked == 0 {
			continue
		}

		price, err := h.c.GetPrice(ctx, pair)
		if err != nil {
			continue
		}

		balance += price * asset.Free
		locked += price * asset.Locked
	}

	return balance, locked, nil
}

func (h *Helper) StartLoggingHelpers(ctx context.Context) {
	z := zap.S().With("context", "Helper.LoggingHelpers")

	for {
		select {
		case <-time.NewTicker(5 * time.Second).C:
			balance, locked, err := h.TotalHoldings(ctx)
			if err != nil {
				z.Warnw(
					"failed to calculate total holdings",
					"error", err.Error(),
				)
			} else {
				z.Infow(
					"total holdings",
					"base_coin", h.baseCoin,
					"amount", balance,
					"amount_locked", locked,
					"total", balance+locked,
				)
			}
		case <-ctx.Done():
			return
		}
	}
}
