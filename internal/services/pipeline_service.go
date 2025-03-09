package services

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sarika-p9/my-pipeline-project/internal/core/domain"
	"github.com/sarika-p9/my-pipeline-project/internal/core/ports"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
)

type PipelineService struct {
	ParallelOrchestrators map[uuid.UUID]*domain.ParallelPipelineOrchestrator
	Repository            ports.PipelineRepository
	mu                    sync.RWMutex
}

func NewPipelineService(repo ports.PipelineRepository) *PipelineService {
	return &PipelineService{
		ParallelOrchestrators: make(map[uuid.UUID]*domain.ParallelPipelineOrchestrator),
		Repository:            repo,
	}
}

func (ps *PipelineService) CreatePipeline(userID uuid.UUID, name string, stageCount int) (uuid.UUID, error) {
	pipelineID := uuid.New()

	ps.mu.Lock()
	orchestrator := domain.NewParallelPipelineOrchestrator(pipelineID, ps.Repository)
	ps.ParallelOrchestrators[pipelineID] = orchestrator
	ps.mu.Unlock()

	for i := 0; i < stageCount; i++ {
		stage := domain.NewBaseStage()
		log.Printf("Adding Stage: %s to Pipeline: %s", stage.GetID(), pipelineID) // Debugging log
		if err := orchestrator.AddStage(stage); err != nil {
			return uuid.Nil, err
		}
	}

	err := ps.Repository.SavePipelineExecution(&models.Pipelines{
		PipelineID:   pipelineID,
		UserID:       userID,
		PipelineName: name, // ✅ Store pipeline name
		Status:       "Created",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		return uuid.Nil, err
	}

	return pipelineID, nil
}

// ✅ Start pipeline execution based on pipeline ID
func (ps *PipelineService) StartPipeline(ctx context.Context, userID uuid.UUID, pipelineID uuid.UUID, input interface{}) error {
	ps.mu.RLock()
	orchestrator, exists := ps.ParallelOrchestrators[pipelineID]
	ps.mu.RUnlock()

	if !exists {
		return errors.New("orchestrator not initialized for this pipeline")
	}

	status, err := ps.Repository.GetPipelineStatus(pipelineID.String())
	if err != nil {
		log.Printf("Failed to get pipeline status: %v", err)
		return err
	}
	if status != "Created" && status != "Paused" {
		return errors.New("invalid pipeline status: " + status)
	}

	// 🚀 Ensure stages are present before execution
	if len(orchestrator.Stages) == 0 {
		return errors.New("no stages found for this pipeline execution")
	}

	if err := ps.updatePipelineStatus(pipelineID, "Running"); err != nil {
		return err
	}

	stageID, _, err := orchestrator.Execute(ctx, userID, pipelineID, input)
	if err != nil {
		_ = ps.updatePipelineStatus(pipelineID, "Failed")
		ps.logExecutionError(pipelineID, stageID, err.Error())
		return err
	}

	return ps.updatePipelineStatus(pipelineID, "Completed")
}

// ✅ Retrieve pipeline status
func (ps *PipelineService) GetPipelineStatus(pipelineID uuid.UUID) (string, error) {
	ps.mu.RLock()
	orchestrator, exists := ps.ParallelOrchestrators[pipelineID]
	ps.mu.RUnlock()

	if !exists {
		return "", errors.New("orchestrator not found for pipeline")
	}

	return orchestrator.GetStatus(pipelineID)
}

// ✅ Cancel pipeline execution
func (ps *PipelineService) CancelPipeline(pipelineID uuid.UUID, userID uuid.UUID) error {
	ps.mu.RLock()
	orchestrator, exists := ps.ParallelOrchestrators[pipelineID]
	ps.mu.RUnlock()

	if !exists {
		log.Printf("Orchestrator not found for pipeline: %s", pipelineID)
		return errors.New("orchestrator not initialized for this pipeline")
	}

	log.Printf("Cancelling pipeline: %s by user: %s", pipelineID, userID)

	err := orchestrator.Cancel(pipelineID, userID)
	if err != nil {
		log.Printf("Failed to cancel pipeline: %v", err)
		_ = ps.updatePipelineStatus(pipelineID, "Failed to Cancel")
		return err
	}

	return ps.updatePipelineStatus(pipelineID, "Cancelled")
}

func (ps *PipelineService) updatePipelineStatus(pipelineID uuid.UUID, status string) error {
	return ps.Repository.UpdatePipelineExecution(&models.Pipelines{
		PipelineID: pipelineID,
		Status:     status,
		UpdatedAt:  time.Now(),
	})
}

func (ps *PipelineService) logExecutionError(pipelineID uuid.UUID, stageID uuid.UUID, errorMsg string) {
	logErr := ps.Repository.SaveExecutionLog(&models.Stages{
		StageID:    stageID,
		PipelineID: pipelineID,
		Status:     "Error",
		ErrorMsg:   errorMsg,
		Timestamp:  time.Now(),
	})
	if logErr != nil {
		log.Printf("Failed to log execution error: %v", logErr)
	}
}

func (ps *PipelineService) GetPipelinesByUser(userID string) ([]models.Pipelines, error) {
	return ps.Repository.GetPipelinesByUser(userID)
}

// GetPipelineStages fetches all stages for a given pipeline
func (ps *PipelineService) GetPipelineStages(pipelineID uuid.UUID) ([]models.Stages, error) {
	return ps.Repository.GetPipelineStages(pipelineID)
}
