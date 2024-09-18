package config

import "flag"

type Config struct {
	Address string
	BaseURL string
}

func NewConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Address, "a", "localhost:8080", "HTTP-Server address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080",
		"Base URL for shortened links")
	flag.Parse()

	return cfg
}
