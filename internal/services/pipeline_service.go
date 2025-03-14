package services

import (
	"context"
	"errors"
	"fmt"
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

func (ps *PipelineService) CreatePipeline(userID uuid.UUID, name string, stageCount int, stageNames []string) (uuid.UUID, error) {
	pipelineID := uuid.New()

	fmt.Printf("🚀 Creating Pipeline: %s\n", pipelineID)

	err := ps.Repository.SavePipelineExecution(&models.Pipelines{
		PipelineID:   pipelineID,
		UserID:       userID,
		PipelineName: name,
		Status:       "Created",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		return uuid.Nil, err
	}

	// ✅ Initialize orchestrator for this pipeline
	ps.mu.Lock()
	ps.ParallelOrchestrators[pipelineID] = domain.NewParallelPipelineOrchestrator(pipelineID, ps.Repository)
	ps.mu.Unlock()
	fmt.Printf("✅ Orchestrator initialized for pipeline: %s\n", pipelineID)

	// ✅ Insert pipeline stages
	if err := ps.InsertPipelineStages(pipelineID, stageNames); err != nil {
		fmt.Println("❌ Error inserting stages:", err)
		return uuid.Nil, err
	}

	return pipelineID, nil
}

// InsertPipelineStages inserts the stages into the database with "Pending" status
// InsertPipelineStages inserts the stages into the database with "Pending" status
func (ps *PipelineService) InsertPipelineStages(pipelineID uuid.UUID, stageNames []string) error {
	fmt.Printf("🔄 Inserting stages for pipeline: %s, Total stages: %d\n", pipelineID, len(stageNames))

	for _, stageName := range stageNames {
		fmt.Printf("🛠️ Inserting Stage: %s\n", stageName)

		stage := models.Stages{
			StageID:    uuid.New(),
			PipelineID: pipelineID,
			StageName:  stageName,
			Status:     "Pending",
		}

		if err := ps.Repository.SaveExecutionLog(&stage); err != nil {
			return err
		}
	}

	return nil
}

// ✅ Start pipeline execution based on pipeline ID
func (ps *PipelineService) StartPipeline(ctx context.Context, userID uuid.UUID, pipelineID uuid.UUID, input interface{}) error {
	fmt.Printf("🚀 Received request to start pipeline: %s\n", pipelineID)

	ps.mu.Lock()
	orchestrator, exists := ps.ParallelOrchestrators[pipelineID]
	if !exists {
		fmt.Println("⚠️ Orchestrator not found in memory, reinitializing...")
		orchestrator = domain.NewParallelPipelineOrchestrator(pipelineID, ps.Repository)
		ps.ParallelOrchestrators[pipelineID] = orchestrator
	}
	ps.mu.Unlock()

	fmt.Println("✅ Orchestrator initialized, updating pipeline status to Running...")
	if err := ps.updatePipelineStatus(pipelineID, "Running"); err != nil {
		return err
	}

	fmt.Println("🔄 Fetching pipeline stages...")
	stages, err := ps.Repository.GetPipelineStages(pipelineID)
	if err != nil {
		return err
	}

	for _, stage := range stages {
		fmt.Printf("🔄 Updating stage %s to Running\n", stage.StageName)
		if err := ps.Repository.UpdateStageStatus(stage.StageID, "Running"); err != nil {
			fmt.Println("❌ Failed to update stage status:", err)
			return err
		}

		// ✅ Create a BaseStage instance and call Execute
		baseStage := domain.NewBaseStage(stage.StageName) // Initialize BaseStage
		_, err := baseStage.Execute(ctx, pipelineID.String(), input)
		if err != nil {
			fmt.Println("❌ Error executing stage:", err)
			return err
		}

		// ✅ Mark stage as Completed
		fmt.Printf("✅ Stage %s Completed\n", stage.StageName)
		if err := ps.Repository.UpdateStageStatus(stage.StageID, "Completed"); err != nil {
			fmt.Println("❌ Failed to update stage to Completed:", err)
			return err
		}

		// ❌ Remove this sleep (already handled in `Execute`)
		// time.Sleep(5 * time.Second)
	}

	fmt.Printf("✅ Pipeline completed: %s\n", pipelineID)
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
