package secondary

import (
	"github.com/google/uuid"
	"github.com/sarika-p9/my-pipeline-project/internal/core/ports"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
	"gorm.io/gorm"
)

type DatabaseAdapter struct {
	DB *gorm.DB
}

var _ ports.PipelineRepository = (*DatabaseAdapter)(nil)

// NewDatabaseAdapter initializes DatabaseAdapter with a gorm.DB instance
func NewDatabaseAdapter(db *gorm.DB) *DatabaseAdapter {
	return &DatabaseAdapter{DB: db}
}

// SaveUser inserts a new user into the database
func (d *DatabaseAdapter) SaveUser(user *models.User) error {
	return d.DB.Create(user).Error
}

// GetUserByID retrieves a user by their ID
func (d *DatabaseAdapter) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := d.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user details
func (d *DatabaseAdapter) UpdateUser(userID uuid.UUID, updates map[string]interface{}) error {
	return d.DB.Model(&models.User{}).Where("user_id = ?", userID).Updates(updates).Error
}

func (d *DatabaseAdapter) SavePipelineExecution(execution *models.PipelineExecution) error {
	return d.DB.Create(execution).Error
}

func (d *DatabaseAdapter) UpdatePipelineExecution(execution *models.PipelineExecution) error {
	return d.DB.Model(execution).Where("pipeline_id = ?", execution.PipelineID).Update("status", execution.Status).Error
}

func (d *DatabaseAdapter) SaveExecutionLog(logEntry *models.ExecutionLog) error {
	return d.DB.Create(logEntry).Error
}

func (d *DatabaseAdapter) GetPipelineStatus(pipelineID string) (string, error) {
	var execution models.PipelineExecution
	if err := d.DB.Where("pipeline_id = ?", pipelineID).First(&execution).Error; err != nil {
		return "", err
	}
	return execution.Status, nil
}

// GetPipelinesByUser retrieves all pipelines for a specific user
func (d *DatabaseAdapter) GetPipelinesByUser(userID string) ([]models.PipelineExecution, error) {
	var pipelines []models.PipelineExecution
	err := d.DB.Where("user_id = ?", userID).Find(&pipelines).Error
	return pipelines, err
}

// GetPipelineStages fetches all stages associated with a pipeline
func (d *DatabaseAdapter) GetPipelineStages(pipelineID uuid.UUID) ([]models.ExecutionLog, error) {
	var stages []models.ExecutionLog
	if err := d.DB.Where("pipeline_id = ?", pipelineID).Find(&stages).Error; err != nil {
		return nil, err
	}
	return stages, nil
}
