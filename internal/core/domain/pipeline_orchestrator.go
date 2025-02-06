package domain

import (
	"context"

	"github.com/google/uuid"
)

type PipelineOrchestrator interface {
	AddStage(stage Stage) error
	Execute(ctx context.Context, input interface{}) (interface{}, error)
	GetStatus(pipelineID uuid.UUID) (Status, error)
	Cancel(pipelineID uuid.UUID) error
}
