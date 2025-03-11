package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDTO_Valid(t *testing.T) {
	t.Parallel()

	tt := []struct {
		in          DTO
		expectedErr string
	}{
		{
			expectedErr: "unsupported version",
		},
		{
			in:          DTO{Version: supportedVersion},
			expectedErr: "invalid DTO state",
		},
		{
			in: DTO{Version: supportedVersion, State: InitialRequest},
		},
		{
			in:          DTO{Version: supportedVersion, State: SolveRequest},
			expectedErr: "to be with payload",
		},
		{
			in: DTO{Version: supportedVersion, State: ErrorResponse, Payload: "Error"},
		},
	}

	for _, tc := range tt {
		err := tc.in.Valid()
		if tc.expectedErr != "" {
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErr)
		} else {
			require.NoError(t, err)
		}
	}
}
