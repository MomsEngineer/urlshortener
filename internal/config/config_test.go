package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "-a", "localhost:9090", "-b", "http://localhost:7777"}

	cfg := NewConfig()

	assert.Equal(t, "localhost:9090", cfg.Address)
	assert.Equal(t, "http://localhost:7777", cfg.BaseURL)
}
