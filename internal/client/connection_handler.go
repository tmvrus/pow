package client

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"pow/internal/config"
	"pow/pkg/api"
)

type handler struct {
	con    connectionSocket
	solver solver
	cfg    *config.ClientConfig
}

func newHandler(cfg *config.ClientConfig, s solver, c connectionSocket) *handler {
	return &handler{
		cfg:    cfg,
		solver: s,
		con:    c,
	}
}

func (h *handler) handle(ctx context.Context) (string, error) {
	request := api.NewDTO(api.InitialRequest)
	err := h.send(ctx, request)
	if err != nil {
		return "", fmt.Errorf("send init request: %w", err)
	}

	response, err := h.receive(ctx)
	if err != nil {
		return "", fmt.Errorf("receive challenge response: %w", err)
	}
	if response.State != api.ChallengeResponse {
		return "", fmt.Errorf("api error response: %d", request.State)
	}

	solution, err := h.solve(ctx, response.Payload)
	if err != nil {
		return "", err
	}

	request = api.NewDTO(api.SolveRequest)
	request.Payload = solution
	err = h.send(ctx, request)
	if err != nil {
		return "", fmt.Errorf("send solve request: %w", err)
	}

	response, err = h.receive(ctx)
	if err != nil {
		return "", fmt.Errorf("receive grant response: %w", err)
	}

	if response.State != api.GrantResponse {
		return "", fmt.Errorf("api error response %q: %d", response.Payload, request.State)
	}

	return response.Payload, nil
}

func (h *handler) solve(ctx context.Context, payload string) (string, error) {
	bits, alg, resource, err := api.ParseChallengePayload(payload)
	if err != nil {
		return "", fmt.Errorf("parse challenge payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, h.cfg.OpTimeout)
	defer cancel()

	solution, err := h.solver.Solve(ctx, bits, alg, resource)
	if err != nil {
		return "", fmt.Errorf("solve: %w", err)
	}

	return solution, nil
}

func (h *handler) receive(ctx context.Context) (*api.DTO, error) {
	deadline := time.Now().Add(h.cfg.OpTimeout)
	err := h.con.SetReadDeadline(deadline)
	if err != nil {
		return nil, fmt.Errorf("set receive deadline: %w", err)
	}

	errCh := make(chan error, 1)
	dataCh := make(chan []byte, 1)

	go func() {
		data, err := bufio.NewReaderSize(h.con, h.cfg.MaxMessageSize).ReadSlice(api.EndFlag)
		if err != nil {
			errCh <- fmt.Errorf("receive response: %w", err)
		}
		dataCh <- data
	}()

	var data []byte
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	case data = <-dataCh:
	}

	response, err := api.UnmarshalMessage(data)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return response, nil
}

func (h *handler) send(ctx context.Context, r *api.DTO) error {
	deadline := time.Now().Add(h.cfg.OpTimeout)
	err := h.con.SetWriteDeadline(deadline)
	if err != nil {
		return fmt.Errorf("set deadline: %w", err)
	}

	data, err := api.MarshalMessage(r)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		_, err = h.con.Write(data)
		if err != nil {
			errCh <- fmt.Errorf("wite: %w", err)
		} else {
			close(errCh)
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err = <-errCh:
		if err != nil {
			return err
		}
	}

	return nil
}
