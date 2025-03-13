package hashcash

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PoW-HC/hashcash/pkg/hash"
	"github.com/PoW-HC/hashcash/pkg/pow"
)

type Verifier struct{}

func NewVerifier() *Verifier {
	return &Verifier{}
}

func (v *Verifier) Verify(_ context.Context, hashAlg, resource, solution string) error {
	hasher, err := hash.NewHasher(hashAlg)
	if err != nil {
		return fmt.Errorf("create new hasher for %q: %w", hashAlg, err)
	}

	parsed, err := parseSolution(solution)
	if err != nil {
		return fmt.Errorf("parse solution: %w", err)
	}

	err = pow.New(hasher).Verify(parsed, resource)
	if err != nil {
		return fmt.Errorf("verify solution: %w", err)
	}

	return nil
}

func parseSolution(s string) (*pow.Hashcach, error) {
	const validSolutionParts = 7
	parts := strings.Split(s, ":")
	if len(parts) != validSolutionParts {
		return nil, fmt.Errorf("invalid solution parts count: %d", len(parts))
	}
	version, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("parse solution version %q; %w", parts[0], err)
	}
	bits, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("parse solution bits %q; %w", parts[0], err)
	}
	ts, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse datetime: %w", err)
	}
	rand, err := base64.StdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, fmt.Errorf("decode random part: %w", err)
	}

	counterBytes, err := base64.StdEncoding.DecodeString(parts[6])
	if err != nil {
		return nil, fmt.Errorf("decode counter part: %w", err)
	}
	counter, err := strconv.ParseInt(string(counterBytes), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("parse counter: %w", err)
	}

	return pow.NewHashcach(
		int32(version),
		int32(bits),
		time.Unix(ts, 0),
		parts[3],
		parts[4],
		rand,
		counter,
	), nil
}
