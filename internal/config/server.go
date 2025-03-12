package config

import "time"

type ServerConfig struct {
	Address        string
	ProviderHost   string
	MaxConnections int
	MaxMessageSize int
	OpTimeout      time.Duration
}

func NewServerConfigWithDefaults() *ServerConfig {
	cfg := &ServerConfig{}
	cfg.Address = "127.0.0.1:22222"
	cfg.MaxConnections = 20
	cfg.OpTimeout = time.Minute
	cfg.MaxMessageSize = 1024
	cfg.ProviderHost = "http://programmingexcuses.com/"
	return cfg
}
