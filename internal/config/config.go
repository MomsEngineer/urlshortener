package config

import (
	"flag"

	"github.com/caarlos0/env"
	"github.com/gofiber/fiber/v2/log"
)

type Config struct {
	Address string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
}

func NewConfig() *Config {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		log.Error(err)
	}

	if cfg.Address == "" || cfg.BaseURL == "" {
		if cfg.Address == "" {
			flag.StringVar(&cfg.Address, "a", "localhost:8080", "HTTP server address")
		}
		if cfg.BaseURL == "" {
			flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080",
				"Base URL for shortened links")
		}
		flag.Parse()
	}

	return cfg
}
