package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
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
	ipLimiter          *bucketLimiter
	handlerConfig      handlerConfig
}

func New(cfg *config.ServerConfig, sf sessionFactory, l *slog.Logger) Server {
	return Server{
		log:                l,
		cfg:                cfg,
		sessionFactory:     sf,
		connectionLimiter:  make(chan struct{}, cfg.MaxConnections),
		currentConnections: &sync.WaitGroup{},
		ipLimiter:          newBucketLimiter(cfg.MaxConnections, cfg.MaxPerIPLimit),
		handlerConfig: handlerConfig{
			opTimeout: cfg.OpTimeout,
			buffSize:  cfg.MaxMessageSize,
		},
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

func (s Server) close(c io.Closer) {
	if err := c.Close(); err != nil {
		s.log.Error("failed to close connection", "error", err.Error())
	}
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

		remoteIP := parseIP(conn.RemoteAddr().String())
		if remoteIP == "" {
			s.log.Debug("drop session due malformed parseIP response", "src", conn.RemoteAddr().String())
			s.close(conn)
			continue
		}

		if !s.ipLimiter.acquire(remoteIP) {
			s.log.Debug("drop session due the ip limit", "src", conn.RemoteAddr().String())
			s.close(conn)
			continue
		}

		select {
		case s.connectionLimiter <- struct{}{}:
			s.currentConnections.Add(1)
			s.log.Debug("start session", "src", conn.RemoteAddr().String())
			go s.handleConnection(ctx, conn, s.handlerConfig, remoteIP)

		default:
			s.log.Debug("drop session due the connections limit", "src", conn.RemoteAddr().String())
			s.ipLimiter.release(remoteIP)
			s.close(conn)
			continue
		}

	}
}

func (s Server) handleConnection(ctx context.Context, conn connectionSocket, cfg handlerConfig, remoteIP string) {
	defer func() {
		s.close(conn)
	}()

	start := time.Now()

	newHandler(s.log, s.sessionFactory.NewSessionHandler(), cfg, conn).run(ctx)

	<-s.connectionLimiter
	s.currentConnections.Done()
	s.ipLimiter.release(remoteIP)

	s.log.Debug("session finished", "src", conn.RemoteAddr().String(), "duration", time.Since(start).String())
}

func parseIP(remoteAddr string) string {
	const portDelimiter = ":"
	portPosition := strings.LastIndex(remoteAddr, portDelimiter)
	if portPosition < 0 {
		return ""
	}

	return remoteAddr[:portPosition]
}
