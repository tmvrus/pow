package config

import "time"

type ClientConfig struct {
	Address        string
	OpTimeout      time.Duration
	MaxMessageSize int
}

func NewClientConfigWithDefaults() *ClientConfig {
	return &ClientConfig{
		Address:        "127.0.0.1:22222",
		OpTimeout:      time.Minute,
		MaxMessageSize: 1024,
	}
}
