package session

import (
	"context"
	"fmt"

	"pow/pkg/api"

	"github.com/google/uuid"
)

const (
	defaultVerifierComplexity = 5
	defaultVerifierHash       = "sha256"
)

type Session struct {
	expectedState int
	provider      wordProvider

	verifier       hashCashVerifier
	verifierSecret string
}

func StartNewSession(wp wordProvider, vr hashCashVerifier) *Session {
	s := &Session{
		expectedState: api.InitialRequest,
		provider:      wp,
		verifier:      vr,
	}

	return s
}

// Handle entry point for each request.
func (s *Session) Handle(ctx context.Context, req *api.DTO) *api.DTO {
	if err := req.Valid(); err != nil {
		res := api.NewDTO(api.ErrorResponse)
		res.Payload = err.Error()
		return res
	}

	if s.expectedState != req.State {
		res := api.NewDTO(api.ErrorResponse)
		res.Payload = fmt.Sprintf("got invalid state %d, expected %d", req.State, s.expectedState)
		return res
	}

	switch req.State {
	case api.InitialRequest:
		return s.handleInitialRequest(ctx)
	case api.SolveRequest:
		return s.handleSolveRequest(ctx, req)
	default:
		res := api.NewDTO(api.ErrorResponse)
		res.Payload = fmt.Sprintf("inconsistent session state: %d", req.State)
		return res
	}
}

func (s *Session) handleSolveRequest(ctx context.Context, req *api.DTO) *api.DTO {
	err := s.verifier.Verify(ctx, defaultVerifierComplexity, defaultVerifierHash, s.verifierSecret, req.Payload)
	if err != nil {
		r := api.NewDTO(api.ErrorResponse)
		r.Payload = fmt.Sprintf("verification failed: %s", err.Error())
		return r
	}

	wisdom, err := s.provider.Get(ctx)
	if err != nil {
		r := api.NewDTO(api.ErrorResponse)
		r.Payload = fmt.Sprintf("something wrong with wisdom source: %s", err.Error())
		return r
	}

	r := api.NewDTO(api.GrantResponse)
	r.Payload = wisdom
	return r
}

func (s *Session) handleInitialRequest(ctx context.Context) *api.DTO {
	s.verifierSecret = uuid.NewString()
	s.expectedState = api.SolveRequest

	r := api.NewDTO(api.ChallengeResponse)
	r.Payload = fmt.Sprintf("%s:%d:%s", defaultVerifierHash, defaultVerifierComplexity, s.verifierSecret)
	return r
}
