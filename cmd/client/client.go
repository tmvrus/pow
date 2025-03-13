package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"pow/internal/client"
	"pow/internal/config"
	"pow/internal/pow/hashcash"

	"github.com/caarlos0/env/v11"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	cfg := config.NewClientConfigWithDefaults()
	err := env.Parse(cfg)
	if err != nil {
		log.Error("failed to parse env", "error", err.Error())
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	solver := hashcash.NewSolver(cfg.MaxPOWIterations)
	err = client.NewClient(log, cfg, solver).Run(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Error("failed to run client", "error", err.Error())
		return
	}

	log.Debug("client stopped")
}
