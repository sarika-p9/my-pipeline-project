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

func (d *DatabaseAdapter) SavePipelineExecution(execution *models.Pipelines) error {
	return d.DB.Create(execution).Error
}

func (d *DatabaseAdapter) UpdatePipelineExecution(execution *models.Pipelines) error {
	return d.DB.Model(&models.Pipelines{}).
		Where("pipeline_id = ?", execution.PipelineID).
		Update("status", execution.Status).Error
}

func (d *DatabaseAdapter) SaveExecutionLog(logEntry *models.Stages) error {
	return d.DB.Create(logEntry).Error
}

// ✅ FIX: Accept pipelineID as string and convert to uuid.UUID
func (d *DatabaseAdapter) GetPipelineStatus(pipelineID string) (string, error) {
	parsedID, err := uuid.Parse(pipelineID)
	if err != nil {
		return "", err // Return an error if the pipelineID is invalid
	}

	var execution models.Pipelines
	if err := d.DB.Where("pipeline_id = ?", parsedID).First(&execution).Error; err != nil {
		return "", err
	}
	return execution.Status, nil
}

// ✅ FIX: Accept userID as string and convert to uuid.UUID
func (d *DatabaseAdapter) GetPipelinesByUser(userID string) ([]models.Pipelines, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	var pipelines []models.Pipelines
	err = d.DB.Where("user_id = ?", parsedID).Find(&pipelines).Error
	return pipelines, err
}

// GetPipelineStages fetches all stages associated with a pipeline
func (d *DatabaseAdapter) GetPipelineStages(pipelineID uuid.UUID) ([]models.Stages, error) {
	var stages []models.Stages
	if err := d.DB.Where("pipeline_id = ?", pipelineID).Find(&stages).Error; err != nil {
		return nil, err
	}
	return stages, nil
}
