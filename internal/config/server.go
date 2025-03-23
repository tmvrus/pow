package config

import "time"

type ServerConfig struct {
	ListenAddress  string        `env:"LISTEN_ADDRESS"`
	ProviderHost   string        `env:"PROVIDER_HOST"`
	MaxConnections int           `env:"MAX_CONNECTIONS"`
	MaxMessageSize int           `env:"MAX_MESSAGE_SIZE"`
	MaxPerIPLimit  int           `env:"MAX_PER_IP_LIMIT"`
	OpTimeout      time.Duration `env:"OPERATION_TIMEOUT"`
}

func NewServerConfigWithDefaults() *ServerConfig {
	cfg := &ServerConfig{}
	cfg.ListenAddress = ":22222"
	cfg.MaxConnections = 20
	cfg.OpTimeout = time.Minute
	cfg.MaxMessageSize = 1024
	cfg.MaxPerIPLimit = 1
	cfg.ProviderHost = "http://programmingexcuses.com/"
	return cfg
}
