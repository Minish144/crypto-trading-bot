package bybit

import "github.com/caarlos0/env/v6"

type Config struct {
	Key    string `env:"BYBIT_KEY"    envDefault:""`
	Secret string `env:"BYBIT_SECRET" envDefault:""`
}

func NewBybitConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
