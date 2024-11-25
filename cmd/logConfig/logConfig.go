package logConfig

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type LoggerConfig struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"DEBUG"`
}

func ReadConfig() (*LoggerConfig, error) {
	config := LoggerConfig{}

	err := env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf("read logConfig error: %w", err)
	}

	return &config, err
}
