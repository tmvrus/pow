package server

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"time"

	"pow/pkg/api"
)

type handlerConfig struct {
	buffSize  int
	opTimeout time.Duration
}

type connectionHandler struct {
	log            *slog.Logger
	sessionHandler api.SessionHandler
	cfg            handlerConfig
	conn           connectionSocket
}

func newHandler(log *slog.Logger, h api.SessionHandler, cfg handlerConfig, conn connectionSocket) *connectionHandler {
	return &connectionHandler{
		log:            log,
		sessionHandler: h,
		cfg:            cfg,
		conn:           conn,
	}
}

func (h *connectionHandler) send(resp *api.DTO) error {
	deadline := time.Now().Add(h.cfg.opTimeout)
	err := h.conn.SetWriteDeadline(deadline)
	if err != nil {
		return fmt.Errorf("set write deadline: %w", err)
	}

	data, err := api.MarshalMessage(resp)
	if err != nil {
		return fmt.Errorf("marshal response: %w", err)
	}
	_, err = h.conn.Write(data)
	if err != nil {
		return fmt.Errorf("wirte marshaled response: %w", err)
	}

	return nil
}

func (h *connectionHandler) run(ctx context.Context) {
	defer func() {
		if e := recover(); e != nil {
			h.log.Error("got panic during handler run", "error", e)
		}
	}()

	input := bufio.NewReaderSize(h.conn, h.cfg.buffSize)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		deadline := time.Now().Add(h.cfg.opTimeout)
		err := h.conn.SetReadDeadline(deadline)
		if err != nil {
			h.log.Error("set read deadline", "error", err.Error())
			return
		}

		data, err := input.ReadSlice(api.EndFlag)
		if err != nil {
			h.log.Error("read request", "error", err.Error())
			return
		}

		req, err := api.UnmarshalMessage(data)
		if err != nil {
			h.log.Error("unmarshall request", "error", err.Error())
			return
		}

		resp := handleWithTimeout(ctx, h.sessionHandler, req, h.cfg.opTimeout)

		err = h.send(resp)
		if err != nil {
			h.log.Error("write response", "error", err.Error())
			return
		}

		if resp.State == api.ErrorResponse || resp.State == api.GrantResponse {
			h.log.Debug("got final api state, finish handling", "state", resp.State)
			return
		}
	}
}

func handleWithTimeout(ctx context.Context, h api.SessionHandler, req *api.DTO, t time.Duration) *api.DTO {
	ctx, cancel := context.WithTimeout(ctx, t)
	defer cancel()
	return h(ctx, req)
}
