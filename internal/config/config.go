package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	DefaultPollSec int
}

func Load() Config {
	return Config{
		Port:           envOr("PORT", "8080"),
		DefaultPollSec: envIntOr("POLL_DEFAULT_SEC", 2),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envIntOr(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return fallback
}
