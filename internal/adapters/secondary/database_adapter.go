package secondary

import (
	"context"

	"github.com/google/uuid"
	"github.com/sarika-p9/my-pipeline-project/internal/core/ports"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
	"gorm.io/gorm"
)

type DatabaseAdapter struct {
	DB *gorm.DB
}

var _ ports.PipelineRepository = (*DatabaseAdapter)(nil)

func NewDatabaseAdapter(db *gorm.DB) *DatabaseAdapter {
	return &DatabaseAdapter{DB: db}
}

func (d *DatabaseAdapter) SaveUser(user *models.User) error {
	return d.DB.Create(user).Error
}

func (d *DatabaseAdapter) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := d.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

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

func (d *DatabaseAdapter) GetPipelineStatus(pipelineID string) (string, error) {
	parsedID, err := uuid.Parse(pipelineID)
	if err != nil {
		return "", err
	}

	var execution models.Pipelines
	if err := d.DB.Where("pipeline_id = ?", parsedID).First(&execution).Error; err != nil {
		return "", err
	}

	return execution.Status, nil
}

func (d *DatabaseAdapter) GetPipelinesByUser(userID string) ([]models.Pipelines, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	var pipelines []models.Pipelines
	err = d.DB.Select("pipeline_id, user_id, status, pipeline_name").
		Where("user_id = ?", parsedID).
		Find(&pipelines).Error
	return pipelines, err
}

func (d *DatabaseAdapter) SaveExecutionLog(logEntry *models.Stages) error {
	if logEntry.StageName == "" {
		logEntry.StageName = "Untitled Stage"
	}

	return d.DB.Create(logEntry).Error
}

func (r *DatabaseAdapter) UpdateStageStatus(stageID uuid.UUID, status string) error {
	return r.DB.Model(&models.Stages{}).
		Where("stage_id = ?", stageID).
		Update("status", status).
		Error
}

func (d *DatabaseAdapter) GetPipelineStages(pipelineID uuid.UUID) ([]models.Stages, error) {
	var stages []models.Stages
	if err := d.DB.Select("stage_id, pipeline_id, stage_name, status, error_msg, timestamp").
		Where("pipeline_id = ?", pipelineID).
		Find(&stages).Error; err != nil {
		return nil, err
	}
	return stages, nil
}

func (d *DatabaseAdapter) DeletePipeline(ctx context.Context, pipelineID string) error {
	parsedID, err := uuid.Parse(pipelineID)
	if err != nil {
		return err
	}

	return d.DB.WithContext(ctx).Where("pipeline_id = ?", parsedID).Delete(&models.Pipelines{}).Error
}

func (d *DatabaseAdapter) GetPipelineByID(pipelineID uuid.UUID) (*models.Pipelines, error) {
	var pipeline models.Pipelines
	if err := d.DB.Where("pipeline_id = ?", pipelineID).First(&pipeline).Error; err != nil {
		return nil, err
	}
	return &pipeline, nil
}
