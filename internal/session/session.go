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

type session struct {
	expectedState int
	provider      wordProvider

	verifier         solutionVerifier
	verifierResource string
}

func newSession(wp wordProvider, sv solutionVerifier) *session {
	s := &session{
		expectedState: api.InitialRequest,
		provider:      wp,
		verifier:      sv,
	}

	return s
}

// Handle entry point for each request.
func (s *session) Handle(ctx context.Context, req *api.DTO) *api.DTO {
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
		return s.handleInitialRequest()
	case api.SolveRequest:
		return s.handleSolveRequest(ctx, req)
	default:
		res := api.NewDTO(api.ErrorResponse)
		res.Payload = fmt.Sprintf("inconsistent session state: %d", req.State)
		return res
	}
}

func (s *session) handleSolveRequest(ctx context.Context, req *api.DTO) *api.DTO {
	err := s.verifier.Verify(ctx, defaultVerifierHash, s.verifierResource, req.Payload)
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

func (s *session) handleInitialRequest() *api.DTO {
	s.verifierResource = uuid.NewString()
	s.expectedState = api.SolveRequest
	return api.NewChallengeResponse(defaultVerifierComplexity, defaultVerifierHash, s.verifierResource)
}
