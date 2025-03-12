package api

import "context"

type SessionHandler func(ctx context.Context, req *DTO) *DTO
