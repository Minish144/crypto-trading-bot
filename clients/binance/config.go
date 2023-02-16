package binance

import "github.com/caarlos0/env/v6"

type Config struct {
	Base string `env:"BASE" envDefault:"USDT"`

	Key    string `env:"BINANCE_KEY"    envDefault:""`
	Secret string `env:"BINANCE_SECRET" envDefault:""`

	TestKey    string `env:"BINANCE_TEST_KEY"    envDefault:""`
	TestSecret string `env:"BINANCE_TEST_SECRET" envDefault:""`

	Test bool `env:"TEST" envDefault:"false"`
}

func newBinanceConfigFromEnv() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
