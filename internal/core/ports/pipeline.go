package ports

import (
	"github.com/google/uuid"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
)

type PipelineRepository interface {
	SavePipelineExecution(execution *models.PipelineExecution) error
	UpdatePipelineExecution(execution *models.PipelineExecution) error
	SaveExecutionLog(logEntry *models.ExecutionLog) error
	GetPipelineStatus(pipelineID string) (string, error)

	GetUserByID(userID uuid.UUID) (*models.User, error)
	SaveUser(user *models.User) error
}
