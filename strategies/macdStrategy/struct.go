package macdStrategy

import (
	"github.com/Minish144/crypto-trading-bot/clients"
	"go.uber.org/zap"
)

type MACDStrategy struct {
	name   string
	cfg    *Config
	client clients.HttpClient
	test   bool
	z      *zap.SugaredLogger
}

func NewMACDStrategy(c clients.HttpClient, cfg *Config) *MACDStrategy {
	z := zap.S().With("context", "MACDStrategy", "symbol", cfg.Symbol)

	return &MACDStrategy{
		name:   "MACD strategy",
		cfg:    cfg,
		client: c,
		z:      z,
	}
}

func (s *MACDStrategy) Name() string {
	return s.name + ": " + s.cfg.Symbol
}
