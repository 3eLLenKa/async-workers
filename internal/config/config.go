package config

import (
	"os"
	"strconv"
)

type Config struct {
	Workers   int
	QueueSize int
}

func Load() *Config {
	cfg := &Config{
		Workers:   4,
		QueueSize: 64,
	}

	if v, ok := os.LookupEnv("WORKERS"); ok {
		if val, err := strconv.Atoi(v); err == nil {
			cfg.Workers = val
		}
	}

	if v, ok := os.LookupEnv("QUEUE_SIZE"); ok {
		if val, err := strconv.Atoi(v); err == nil {
			cfg.QueueSize = val
		}
	}

	return cfg
}
