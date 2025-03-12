package session

import (
	"pow/pkg/api"
)

type Factory struct {
	wp wordProvider
	sv solutionVerifier
}

func NewFactory(wp wordProvider, sv solutionVerifier) Factory {
	return Factory{wp: wp, sv: sv}
}

func (f Factory) NewSessionHandler() api.SessionHandler {
	return newSession(f.wp, f.sv).Handle
}
