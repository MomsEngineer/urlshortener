package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name        string
		Address     string
		BaseURL     string
		FilePath    string
		DataBaseDSN string
		args        []string
		expected    *Config
	}{
		{
			name: "config without env and flags",
			args: []string{"cmd"},
			expected: &Config{
				Address:     "localhost:8080",
				BaseURL:     "http://localhost:8080",
				FilePath:    "/tmp/short-url-db.json",
				DataBaseDSN: "",
			},
		},
		{
			name: "config without env and with flags",
			args: []string{
				"cmd", "-a", "localhost:9090", "-b", "http://localhost:7777",
				"-f", "test.json", "-d", "test:db",
			},
			expected: &Config{
				Address:     "localhost:9090",
				BaseURL:     "http://localhost:7777",
				FilePath:    "test.json",
				DataBaseDSN: "test:db",
			},
		},

		{
			name:        "config with env and flags",
			Address:     "localhost:9999",
			BaseURL:     "http://test",
			FilePath:    "test.json",
			DataBaseDSN: "test:db:config",
			args: []string{
				"cmd", "-a", "localhost:7070", "-b", "http://localhost:7777",
				"-f", "/tmp/test.json", "-d", "test:db:flag",
			},
			expected: &Config{
				Address:     "localhost:9999",
				BaseURL:     "http://test",
				FilePath:    "test.json",
				DataBaseDSN: "test:db:config",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.Address != "" {
				os.Setenv("SERVER_ADDRESS", tt.Address)
			}

			if tt.BaseURL != "" {
				os.Setenv("BASE_URL", tt.BaseURL)
			}

			if tt.FilePath != "" {
				os.Setenv("FILE_STORAGE_PATH", tt.FilePath)
			}

			if tt.DataBaseDSN != "" {
				os.Setenv("DATABASE_DSN", tt.DataBaseDSN)
			}

			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = tt.args

			cfg := NewConfig()
			assert.Equal(t, tt.expected, cfg)
		})
	}
}
