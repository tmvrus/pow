package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseIP(t *testing.T) {
	t.Parallel()

	ip := parseIP("192.0.2.1:123")
	require.Equal(t, "192.0.2.1", ip)

	ip6 := parseIP("[2001:db8::1]:80")
	require.Equal(t, "[2001:db8::1]", ip6)

	empty := parseIP("malformed")
	require.Empty(t, empty)
}

func Test_BucketLimiter(t *testing.T) {
	t.Parallel()

	l := newBucketLimiter(10, 2)
	key1 := "xxx"
	key2 := "yyy"

	require.True(t, l.acquire(key1))
	require.True(t, l.acquire(key1))
	require.False(t, l.acquire(key1))

	require.True(t, l.acquire(key2))
	require.True(t, l.acquire(key2))

	l.release(key1)
	require.True(t, l.acquire(key1))
}
