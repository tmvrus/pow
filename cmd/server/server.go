package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"pow/internal/config"
	"pow/internal/pow/hashcash"
	"pow/internal/provider"
	"pow/internal/server"
	"pow/internal/session"
)

func main() {
	cfg := config.NewServerConfigWithDefaults()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	wisdomProvider := provider.NewProvider(cfg.ProviderHost)
	sessionFactory := session.NewFactory(wisdomProvider, hashcash.NewVerifier())

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	err := server.New(cfg, sessionFactory, log).Run(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Debug("application stopped")
			return
		}

		log.Error("failed to run application", "error", err.Error())
		return
	}
}
