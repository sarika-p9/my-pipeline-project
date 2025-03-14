package domain

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sarika-p9/my-pipeline-project/internal/core/ports"
	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
)

type ParallelPipelineOrchestrator struct {
	PipelineID uuid.UUID
	UserID     uuid.UUID
	Stages     []Stage
	mu         sync.Mutex
	dbRepo     ports.PipelineRepository
}

func NewParallelPipelineOrchestrator(pipelineID uuid.UUID, dbRepo ports.PipelineRepository) *ParallelPipelineOrchestrator {
	return &ParallelPipelineOrchestrator{
		PipelineID: pipelineID,
		dbRepo:     dbRepo,
		Stages:     []Stage{},
	}
}

func (p *ParallelPipelineOrchestrator) AddStage(stage Stage) error {
	if stage == nil {
		return errors.New("stage cannot be nil")
	}
	p.mu.Lock()
	p.Stages = append(p.Stages, stage)
	p.mu.Unlock()
	return nil
}

func (p *ParallelPipelineOrchestrator) Execute(ctx context.Context, userID uuid.UUID, pipelineID uuid.UUID, input interface{}) (uuid.UUID, interface{}, error) {
	user, err := p.dbRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("Failed to validate user existence: %v", err)
		return pipelineID, nil, err
	}
	if user == nil {
		return pipelineID, nil, errors.New("user does not exist")
	}

	pipeline, err := p.dbRepo.GetPipelineByID(pipelineID)
	if err != nil {
		log.Printf("Failed to fetch pipeline details: %v", err)
		return pipelineID, nil, err
	}

	if err := p.dbRepo.UpdatePipelineExecution(&models.Pipelines{
		PipelineID:   pipelineID,
		PipelineName: pipeline.PipelineName,
		Status:       "Running",
		UpdatedAt:    time.Now(),
	}); err != nil {
		log.Printf("Failed to update pipeline execution status: %v", err)
		return pipelineID, nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]interface{}, 0, len(p.Stages))
	errorsSlice := make([]error, 0, len(p.Stages))

	for _, stage := range p.Stages {
		wg.Add(1)
		go func(stage Stage) {
			defer wg.Done()

			stageName := stage.GetName()

			infrastructure.WebSocket.SendMessage(pipeline.PipelineName, stageName, "Running")

			result, err := stage.Execute(ctx, pipeline.PipelineName, input)
			logEntry := &models.Stages{
				StageID:    stage.GetID(),
				StageName:  stageName,
				PipelineID: pipelineID,
				Status:     "Completed",
				Timestamp:  time.Now(),
			}

			if err != nil {
				logEntry.Status = "Failed"
				logEntry.ErrorMsg = err.Error()

				infrastructure.WebSocket.SendMessage(pipeline.PipelineName, stageName, "Failed")

				mu.Lock()
				errorsSlice = append(errorsSlice, err)
				mu.Unlock()
			} else {
				infrastructure.WebSocket.SendMessage(pipeline.PipelineName, stageName, "Completed")

				mu.Lock()
				results = append(results, result)
				mu.Unlock()
			}

			if err := p.dbRepo.SaveExecutionLog(logEntry); err != nil {
				log.Printf("Failed to save execution log: %v", err)
			}
		}(stage)
	}

	wg.Wait()

	finalStatus := "Completed"
	if len(errorsSlice) > 0 {
		finalStatus = "Failed"
	}

	if err := p.dbRepo.UpdatePipelineExecution(&models.Pipelines{
		PipelineID:   pipelineID,
		PipelineName: pipeline.PipelineName,
		Status:       finalStatus,
		UpdatedAt:    time.Now(),
	}); err != nil {
		log.Printf("Failed to update final pipeline execution status: %v", err)
	}

	return pipelineID, results, nil
}

func (p *ParallelPipelineOrchestrator) GetStatus(pipelineID uuid.UUID) (string, error) {
	return p.dbRepo.GetPipelineStatus(pipelineID.String())
}

func (p *ParallelPipelineOrchestrator) Cancel(pipelineID uuid.UUID, userID uuid.UUID) error {
	log.Printf("Cancelling pipeline: %s for user: %s", pipelineID, userID)

	status, err := p.dbRepo.GetPipelineStatus(pipelineID.String())
	if err != nil {
		log.Printf("Error fetching pipeline status: %v", err)
		return errors.New("pipeline not found")
	}

	if status == "Completed" {
		log.Printf("Pipeline %s is already completed, cannot cancel", pipelineID)
		return errors.New("cannot cancel a completed pipeline")
	}
	log.Printf("Cancelling pipeline %s...", pipelineID)

	err = p.dbRepo.UpdatePipelineExecution(&models.Pipelines{
		PipelineID: pipelineID,
		Status:     "Cancelled",
		UpdatedAt:  time.Now(),
	})

	if err != nil {
		log.Printf("Failed to update pipeline status: %v", err)
		return errors.New("failed to update pipeline status")
	}

	log.Printf("Pipeline %s successfully cancelled", pipelineID)
	return nil
}
