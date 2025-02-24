package domain

import (
	"context"
	"fmt"
	"sync"
)

// Stage defines the interface for a pipeline stage

// PipelineOrchestrator handles the execution of stages
type PipelineOrchestrator struct {
	stages []Stage
}

// NewPipelineOrchestrator creates a new orchestrator
func NewPipelineOrchestrator() *PipelineOrchestrator {
	return &PipelineOrchestrator{
		stages: make([]Stage, 0),
	}
}

// AddStage adds a stage to the pipeline
func (p *PipelineOrchestrator) AddStage(stage Stage) {
	p.stages = append(p.stages, stage)
}

// Execute runs all stages in parallel
func (p *PipelineOrchestrator) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	var wg sync.WaitGroup
	errChan := make(chan error, len(p.stages))
	resultChan := make(chan interface{}, len(p.stages))

	for _, stage := range p.stages {
		wg.Add(1)
		go func(s Stage) {
			defer wg.Done()
			fmt.Printf("Executing stage: %s with input: %v\n", s.GetID(), input)
			result, err := s.Execute(ctx, input)
			if err != nil {
				errChan <- fmt.Errorf("stage %s failed: %w", s.GetID(), err)
				return
			}
			resultChan <- result
		}(stage)
	}

	wg.Wait()
	close(errChan)
	close(resultChan)

	if len(errChan) > 0 {
		for err := range errChan {
			return nil, err
		}
	}

	// Collect results (for simplicity, return the first one)
	for res := range resultChan {
		return res, nil
	}

	return input, nil
}
