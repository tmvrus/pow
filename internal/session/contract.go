//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -source=contract.go -destination=./contract_mock_test.go -package=session
package session

import "context"

type hashCashVerifier interface {
	Verify(ctx context.Context, zeros int, alg, secret, verifiable string) error
}

type wordProvider interface {
	Get(ctx context.Context) (string, error)
}
