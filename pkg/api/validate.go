package api

import "fmt"

func (d DTO) Valid() error {
	if d.Version != supportedVersion {
		return fmt.Errorf("unsupported version %d", d.Version)
	}
	if d.State <= Invalid || d.State > GrantResponse {
		return fmt.Errorf("invalid DTO state %d", d.State)
	}

	if d.State == InitialRequest {
		return nil
	}

	if d.Payload == "" {
		return fmt.Errorf("state %d expected to be with payload", d.State)
	}

	return nil
}
