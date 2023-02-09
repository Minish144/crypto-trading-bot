package strategies

import "github.com/Minish144/crypto-trading-bot/strategies/gridStrategy"

// interface validations
var _ Strategy = &gridStrategy.GridStrategy{}

// aliases
var NewGridStrategy = gridStrategy.NewGridStrategy
