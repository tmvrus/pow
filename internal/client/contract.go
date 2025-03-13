//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -source=contract.go -destination=./contract_mock_test.go -package=client
package client

import (
	"context"
	"net"
)

type solver interface {
	Solve(ctx context.Context, bits int32, hashAlg string, resource string) (string, error)
}

type connectionSocket interface {
	net.Conn
}
