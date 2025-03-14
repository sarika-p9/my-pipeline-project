package domain

import (
	"context"

	"github.com/google/uuid"
)

type PipelineOrchestrator interface {
	AddStage(stage Stage) error
	Execute(ctx context.Context, userID uuid.UUID, pipelineID uuid.UUID, input interface{}) (uuid.UUID, interface{}, error)
	GetStatus(pipelineID uuid.UUID) (string, error)
	Cancel(pipelineID uuid.UUID, userID uuid.UUID) error
}
