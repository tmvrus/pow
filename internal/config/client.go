package config

import "time"

type ClientConfig struct {
	ServerAddress    string        `env:"SERVER_ADDRESS"`
	OpTimeout        time.Duration `env:"OPERATION_TIMEOUT"`
	MaxMessageSize   int           `env:"MAX_MESSAGE_SIZE"`
	MaxPOWIterations int64         `env:"MAX_ITERATIONS"`
}

func NewClientConfigWithDefaults() *ClientConfig {
	//nolint:mnd
	return &ClientConfig{
		ServerAddress:    "127.0.0.1:22222",
		OpTimeout:        time.Minute,
		MaxMessageSize:   1024,
		MaxPOWIterations: 1 << 30,
	}
}
