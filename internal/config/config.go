package config

import (
	"flag"

	"github.com/caarlos0/env"
	"github.com/gofiber/fiber/v2/log"
)

type Config struct {
	Address  string `env:"SERVER_ADDRESS"`
	BaseURL  string `env:"BASE_URL"`
	FilePath string `env:"FILE_STORAGE_PATH"`
}

func NewConfig() *Config {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		log.Error(err)
	}

	var a, b, f string
	flag.StringVar(&a, "a", "localhost:8080", "HTTP server address")
	flag.StringVar(&b, "b", "http://localhost:8080", "Base URL for shortened links")
	flag.StringVar(&f, "f", "/tmp/short-url-db.json", "The path to storage file")
	flag.Parse()

	if cfg.Address == "" {
		cfg.Address = a
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = b
	}

	if cfg.FilePath == "" {
		cfg.FilePath = f
	}

	return cfg
}
