package domain

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, cmd CreateUserInfoCmd) (*UserInfo, error)
	GetLatest(ctx context.Context, userID uuid.UUID) (*UserInfo, error)
	List(ctx context.Context, userID uuid.UUID) ([]*UserInfo, error)
}

type CreateUserInfoCmd struct {
	UserID uuid.UUID
	Weight float64
	Height int64
	Age    int64
	Date   string
}
