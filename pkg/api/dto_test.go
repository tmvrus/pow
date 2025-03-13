package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParsePayload(t *testing.T) {
	t.Parallel()

	bits, alg, resource, err := ParseChallengePayload("sha1:1:xxx")
	require.NoError(t, err)
	require.Equal(t, int32(1), bits)
	require.Equal(t, "sha1", alg)
	require.Equal(t, "xxx", resource)
}
