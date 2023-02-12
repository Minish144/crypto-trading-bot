package strategies

import (
	"github.com/Minish144/crypto-trading-bot/strategies/gridStrategy"
	"github.com/Minish144/crypto-trading-bot/strategies/macdStrategy"
)

// interface validations
var _ Strategy = &gridStrategy.GridStrategy{}

// aliases
var (
	NewGridStrategy = gridStrategy.NewGridStrategy
	NewMACDStrategy = macdStrategy.NewMACDStrategy
)
