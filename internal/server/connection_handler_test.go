package server

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"pow/pkg/api"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_HandlerHappyPath(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	socket := NewMockconnectionSocket(ctrl)

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	ctx := context.Background()

	request := api.NewDTO(api.SolveRequest)
	request.Payload = "Payload"

	socket.
		EXPECT().
		SetReadDeadline(futureMatcher{}).
		Return(nil)
	socket.
		EXPECT().
		Read(gomock.Any()).
		DoAndReturn(func(p []byte) (int, error) {
			data, _ := api.MarshalMessage(request)
			data = append(data, api.EndFlag)
			return copy(p, data), nil
		}).Times(1)

	response := api.NewDTO(api.GrantResponse)
	response.Payload = "Payload"

	sessionHandler := api.SessionHandler(func(_ context.Context, req *api.DTO) *api.DTO {
		t.Helper()
		require.Equal(t, request, req)
		return response
	})

	socket.
		EXPECT().
		SetWriteDeadline(futureMatcher{}).
		Return(nil)
	data, _ := api.MarshalMessage(response)
	socket.
		EXPECT().
		Write(byteMatcher{t: t, want: data}).
		Return(0, nil)

	handler := newHandler(log, sessionHandler, handlerConfig{
		buffSize:  1024,
		opTimeout: time.Minute,
	}, socket)

	handler.run(ctx)
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
