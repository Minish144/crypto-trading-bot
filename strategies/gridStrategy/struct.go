package gridStrategy

import (
	"github.com/Minish144/crypto-trading-bot/clients"
	"go.uber.org/zap"
)

type GridStrategy struct {
	name   string
	cfg    *Config
	client clients.HttpClient
	test   bool
	z      *zap.SugaredLogger
}

func NewGridStrategy(c clients.HttpClient, cfg *Config) *GridStrategy {
	z := zap.S().With("context", "GridStrategy", "symbol", cfg.Symbol)

	return &GridStrategy{
		name:   "grid strategy",
		cfg:    cfg,
		client: c,
		z:      z,
	}
}

func (s *GridStrategy) Name() string {
	return s.name + ": " + s.cfg.Symbol
}
