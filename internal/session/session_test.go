package session

import (
	"context"
	"testing"

	"pow/pkg/api"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSession_HappyPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	providerMock := NewMockwordProvider(ctrl)
	verifierMock := NewMocksolutionVerifier(ctrl)

	session := newSession(providerMock, verifierMock)

	request := api.NewDTO(api.InitialRequest)
	response := session.Handle(ctx, request)
	require.Equal(t, api.ChallengeResponse, response.State)
	require.NotEmpty(t, response.Payload)

	verifierMock.
		EXPECT().
		Verify(ctx, "sha256", session.verifierResource, "SOLVED").
		Return(nil)
	providerMock.
		EXPECT().
		Get(ctx).
		Return("HELLO", nil)

	request = api.NewDTO(api.SolveRequest)
	request.Payload = "SOLVED"
	response = session.Handle(ctx, request)
	require.Equal(t, api.GrantResponse, response.State)
	require.Equal(t, "HELLO", response.Payload)
}
