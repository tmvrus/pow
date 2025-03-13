package api

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"
)

// EndFlag indicates end of the message.
// Definitely a bad choice for production, but good enough for idea testing.
const EndFlag byte = 222

type DTO struct {
	Version int
	State   int
	Payload string
}

func NewDTO(s int) *DTO {
	return &DTO{
		Version: supportedVersion,
		State:   s,
	}
}

func MarshalMessage(d *DTO) ([]byte, error) {
	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(d)
	if err != nil {
		return nil, err
	}
	return append(buf.Bytes(), EndFlag), nil
}

func UnmarshalMessage(data []byte) (*DTO, error) {
	d := new(DTO)
	buf := bytes.NewBuffer(data)
	err := gob.NewDecoder(buf).Decode(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func NewChallengeResponse(bits int32, hashAlg, resource string) *DTO {
	d := NewDTO(ChallengeResponse)
	d.Payload = fmt.Sprintf("%s:%d:%s", hashAlg, bits, resource)
	return d
}

func ParseChallengePayload(payload string) (int32, string, string, error) {
	parts := strings.Split(payload, ":")
	if len(parts) != 3 {
		return 0, "", "", fmt.Errorf("invalid format: %q", payload)
	}

	bits, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return 0, "", "", fmt.Errorf("parse bits: %w", err)
	}

	return int32(bits), parts[0], parts[2], nil
}
