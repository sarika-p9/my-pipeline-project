package domain

import (
	"context"

	"github.com/google/uuid"
)

type Stage interface {
	GetID() uuid.UUID
	Execute(ctx context.Context, input interface{}) (interface{}, error)
	HandleError(ctx context.Context, err error) error
	Rollback(ctx context.Context, input interface{}) error
}
