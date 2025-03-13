package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"pow/internal/config"
)

type Server struct {
	log *slog.Logger
	cfg *config.ServerConfig

	sessionFactory     sessionFactory
	connectionLimiter  chan struct{}
	currentConnections *sync.WaitGroup
}

func New(cfg *config.ServerConfig, sf sessionFactory, l *slog.Logger) Server {
	return Server{
		log:                l,
		cfg:                cfg,
		sessionFactory:     sf,
		connectionLimiter:  make(chan struct{}, cfg.MaxConnections),
		currentConnections: &sync.WaitGroup{},
	}
}

func (s Server) Run(ctx context.Context) error {
	l, err := net.Listen("tcp", s.cfg.ListenAddress)
	if err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	s.log.Debug("ready to accept connections", "address", s.cfg.ListenAddress)

	go func() {
		<-ctx.Done()
		s.currentConnections.Wait()
		if err := l.Close(); err != nil {
			s.log.Error("failed to close listener", "error", err.Error())
		}
	}()

	if err := s.acceptLoop(ctx, l); err != nil {
		return err
	}
	return nil
}

func (s Server) acceptLoop(ctx context.Context, l net.Listener) error {
	for {
		select {
		case <-ctx.Done():
			s.log.Debug("got context done, stop application")
			return ctx.Err()
		default:
		}

		conn, err := l.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				s.log.Error("failed to accept connection", "error", err.Error())
			}
			continue
		}

		select {
		case s.connectionLimiter <- struct{}{}:
			s.currentConnections.Add(1)
			s.log.Debug("start session", "src", conn.RemoteAddr().String())

		default:
			s.log.Debug("drop session due the limit", "src", conn.RemoteAddr().String())
			if err := conn.Close(); err != nil {
				s.log.Error("failed to close connection", "error", err.Error())
			}
			continue
		}

		go func() {
			defer func() {
				if err := conn.Close(); err != nil {
					s.log.Error("failed to close connection", "error", err.Error())
				}
			}()

			start := time.Now()
			cfg := handlerConfig{
				opTimeout: s.cfg.OpTimeout,
				buffSize:  s.cfg.MaxMessageSize,
			}

			newHandler(s.log, s.sessionFactory.NewSessionHandler(), cfg, conn).run(ctx)

			<-s.connectionLimiter
			s.currentConnections.Done()

			s.log.Debug("session finished", "src", conn.RemoteAddr().String(), "duration", time.Since(start).String())
		}()
	}
}
