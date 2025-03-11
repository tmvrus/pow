package api

const (
	supportedVersion = 1
)

// This `enum` represents reduced Challenge–response protocol states.
const (
	Invalid = iota
	ErrorResponse
	InitialRequest
	ChallengeResponse
	SolveRequest
	GrantResponse
)
