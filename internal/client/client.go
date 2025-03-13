package client

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"pow/internal/config"
)

type Client struct {
	log    *slog.Logger
	cfg    *config.ClientConfig
	solver solver
}

func NewClient(log *slog.Logger, cfg *config.ClientConfig, solver solver) *Client {
	return &Client{
		log:    log,
		cfg:    cfg,
		solver: solver,
	}
}

func (c *Client) Run(ctx context.Context) error {
	remoteAddr, err := net.ResolveTCPAddr("tcp", c.cfg.ServerAddress)
	if err != nil {
		return fmt.Errorf("solve address %q: %w", c.cfg.ServerAddress, err)
	}

	con, err := net.DialTCP("tcp", nil, remoteAddr)
	if err != nil {
		return fmt.Errorf("dial address %q: %w", remoteAddr.String(), err)
	}

	defer func() {
		if err := con.Close(); err != nil {
			c.log.Error("close connection", "error", err.Error())
		}
	}()

	wisdom, err := newHandler(c.cfg, c.solver, con).handle(ctx)
	if err != nil {
		return fmt.Errorf("connection handle: %w", err)
	}

	c.log.Debug("your piece of wisdom", "piece", wisdom)
	return nil
}
