package config

import "server/internal/app/storage"

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Storage  *storage.DbConfig
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
		Storage:  storage.NewConfig(),
	}
}
