//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -source=contract.go -destination=./contract_mock_test.go -package=server
package server

import (
	"net"

	"pow/pkg/api"
)

type sessionFactory interface {
	NewSessionHandler() api.SessionHandler
}

type connectionSocket interface {
	net.Conn
}
