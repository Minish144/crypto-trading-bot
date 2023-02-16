package config

import "github.com/caarlos0/env/v6"

type Config struct {
	Debug bool `env:"DEBUG" envDefault:"false"`

	LogLevel    string `env:"LOG_LEVEL"    envDefault:"info"`
	LogOutput   string `env:"LOG_OUTPUT"   envDefault:"stdout"`
	LogEncoding string `env:"LOG_ENCODING" envDefault:"json"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
