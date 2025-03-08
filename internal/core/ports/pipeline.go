package ports

import (
	"github.com/google/uuid"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
)

type PipelineRepository interface {
	SavePipelineExecution(execution *models.Pipelines) error
	UpdatePipelineExecution(execution *models.Pipelines) error
	SaveExecutionLog(logEntry *models.Stages) error
	GetPipelineStatus(pipelineID string) (string, error)
	GetUserByID(userID uuid.UUID) (*models.User, error)
	SaveUser(user *models.User) error
	UpdateUser(userID uuid.UUID, updates map[string]interface{}) error
	GetPipelinesByUser(userID string) ([]models.Pipelines, error)
	GetPipelineStages(pipelineID uuid.UUID) ([]models.Stages, error)
}
