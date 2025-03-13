package hashcash

import (
	"context"
	"fmt"
	"testing"

	"github.com/PoW-HC/hashcash/pkg/hash"
	"github.com/PoW-HC/hashcash/pkg/pow"
	"github.com/stretchr/testify/require"
)

func Test_SolutionParser(t *testing.T) {
	t.Parallel()

	solution := "1:3:1741768305:resource::B/tAlsAyo+A=:MjRkZQ=="
	hc, err := parseSolution(solution)
	require.NoError(t, err)
	require.Equal(t, solution, hc.String())
}

func TestVerifier_HappyPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	hc, err := pow.InitHashcash(3, "resource", nil)
	require.NoError(t, err)

	hasher, err := hash.NewHasher("sha256")
	require.NoError(t, err)
	p := pow.New(hasher)
	hc, err = p.Compute(ctx, hc, 1<<20)
	require.NoError(t, err)

	err = NewVerifier().Verify(ctx, "sha256", "resource", hc.String())
	fmt.Println(hc.String())
	require.NoError(t, err)
}

func Test_SolverVerifier(t *testing.T) {
	t.Parallel()

	bits := int32(3)
	alg := "sha256"
	resource := "resource"

	solution, err := NewSolver(1<<30).Solve(context.Background(), bits, alg, resource)
	require.NoError(t, err)

	err = NewVerifier().Verify(context.Background(), alg, resource, solution)
	require.NoError(t, err)
}
