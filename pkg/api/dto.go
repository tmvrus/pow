package api

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
