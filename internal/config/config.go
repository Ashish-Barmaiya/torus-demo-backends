package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Service  string
	Instance string
	Port     int
}

func Load() (Config, error) {
	cfg := Config{
		Service:  os.Getenv("SERVICE"),
		Instance: os.Getenv("INSTANCE"),
		Port:     3000,
	}

	if port := os.Getenv("PORT"); port != "" {
		value, err := strconv.Atoi(port)
		if err != nil || value < 1 || value > 65535 {
			return Config{}, fmt.Errorf("invalid PORT %q", port)
		}
		cfg.Port = value
	}

	if cfg.Service == "" {
		return Config{}, fmt.Errorf("SERVICE must not be empty")
	}

	if cfg.Instance == "" {
		return Config{}, fmt.Errorf("INSTANCE must not be empty")
	}

	return cfg, nil
}
