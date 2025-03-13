package client

import (
	"context"
	"testing"
	"time"

	"pow/internal/config"
	"pow/pkg/api"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_HandlerHappyPath(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	socket := NewMockconnectionSocket(ctrl)
	solver := NewMocksolver(ctrl)
	cfg := config.NewClientConfigWithDefaults()

	request, err := api.MarshalMessage(api.NewDTO(api.InitialRequest))
	require.NoError(t, err)
	prepareWrite(t, socket, request)

	challengeResp := api.NewDTO(api.ChallengeResponse)
	challengeResp.Payload = "sha256:10:resource"
	challengeRespByte, err := api.MarshalMessage(challengeResp)
	require.NoError(t, err)
	prepareRead(t, socket, challengeRespByte)

	solver.
		EXPECT().
		Solve(gomock.Any(), int32(10), "sha256", "resource").
		Return("solution", nil)

	solveRequest := api.NewDTO(api.SolveRequest)
	solveRequest.Payload = "solution"
	request, err = api.MarshalMessage(solveRequest)
	require.NoError(t, err)
	prepareWrite(t, socket, request)

	grantResp := api.NewDTO(api.GrantResponse)
	grantResp.Payload = "OK"
	grantRespByte, err := api.MarshalMessage(grantResp)
	require.NoError(t, err)
	prepareRead(t, socket, grantRespByte)

	result, err := newHandler(cfg, solver, socket).handle(context.Background())
	require.NoError(t, err)
	require.Equal(t, "OK", result)
}

func prepareWrite(t *testing.T, socket *MockconnectionSocket, data []byte) {
	t.Helper()

	socket.
		EXPECT().
		SetWriteDeadline(futureMatcher{}).
		Return(nil)
	socket.
		EXPECT().
		Write(byteMatcher{t: t, want: data}).
		Return(0, nil)
}

func prepareRead(t *testing.T, socket *MockconnectionSocket, data []byte) {
	t.Helper()

	socket.
		EXPECT().
		SetReadDeadline(futureMatcher{}).
		Return(nil)
	socket.
		EXPECT().
		Read(gomock.Any()).
		DoAndReturn(func(p []byte) (int, error) {
			return copy(p, data), nil
		})
}

type futureMatcher struct{}

func (m futureMatcher) Matches(x any) bool {
	t, ok := x.(time.Time)
	if !ok {
		return false
	}

	return t.After(time.Now())
}

func (m futureMatcher) String() string {
	return ""
}

type byteMatcher struct {
	t    *testing.T
	want []byte
}

func (m byteMatcher) Matches(x any) bool {
	got, ok := x.([]byte)
	if !ok {
		m.t.Errorf("expected []byte got %T", x)
		return false
	}

	gotS := string(got)
	wantS := string(m.want)
	if gotS != wantS {
		m.t.Errorf("expected %q, but got %q", wantS, gotS)
		return false
	}

	return true
}

func (m byteMatcher) String() string {
	return ""
}
