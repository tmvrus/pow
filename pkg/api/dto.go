package api

import (
	"bytes"
	"encoding/gob"
)

const EndFlag byte = '\n'

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

func Marshal(d *DTO) ([]byte, error) {
	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(d)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Unmarshal(data []byte) (*DTO, error) {
	d := new(DTO)
	buf := bytes.NewBuffer(data)
	err := gob.NewDecoder(buf).Decode(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}
