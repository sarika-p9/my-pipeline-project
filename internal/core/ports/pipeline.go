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
	UpdateUser(userID uuid.UUID, updates map[string]interface{}) error
	GetPipelinesByUser(userID string) ([]models.PipelineExecution, error)
	GetPipelineStages(pipelineID uuid.UUID) ([]models.ExecutionLog, error)
}
