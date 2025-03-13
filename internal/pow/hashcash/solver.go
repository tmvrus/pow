package hashcash

import (
	"context"
	"fmt"

	"github.com/PoW-HC/hashcash/pkg/hash"
	"github.com/PoW-HC/hashcash/pkg/pow"
)

type Solver struct {
	maxIterations int64
}

func NewSolver(maxIterations int64) Solver {
	return Solver{maxIterations: maxIterations}
}

func (s Solver) Solve(ctx context.Context, bits int32, hashAlg string, resource string) (string, error) {
	hasher, err := hash.NewHasher(hashAlg)
	if err != nil {
		return "", fmt.Errorf("create hasher %q: %w", hashAlg, err)
	}

	hashcash, err := pow.InitHashcash(bits, resource, nil)
	if err != nil {
		return "", fmt.Errorf("init hashcash: %w", err)
	}

	solution, err := pow.New(hasher).Compute(ctx, hashcash, s.maxIterations)
	if err != nil {
		return "", fmt.Errorf("compute solution: %w", err)
	}

	return solution.String(), nil
}
