package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sarika-p9/my-pipeline-project/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

// MockStage for testing
type MockStage struct {
	id      uuid.UUID
	process time.Duration
}

func (m *MockStage) GetID() uuid.UUID {
	return m.id
}

func (m *MockStage) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	time.Sleep(m.process) // Simulate work
	return input.(string) + "_processed_by_" + m.id.String(), nil
}

func (m *MockStage) HandleError(ctx context.Context, err error) error {
	return err
}

func (m *MockStage) Rollback(ctx context.Context, input interface{}) error {
	return nil
}

func TestPipelineParallelExecution(t *testing.T) {
	orchestrator := domain.NewPipelineOrchestrator()

	// Adding stages with simulated processing time
	stage1 := &MockStage{id: uuid.New(), process: 2 * time.Second}
	stage2 := &MockStage{id: uuid.New(), process: 2 * time.Second}
	stage3 := &MockStage{id: uuid.New(), process: 2 * time.Second}

	orchestrator.AddStage(stage1)
	orchestrator.AddStage(stage2)
	orchestrator.AddStage(stage3)

	start := time.Now()
	result, err := orchestrator.Execute(context.Background(), "input_data")
	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Since stages run in parallel, total time should be around 2 seconds, not 6
	if elapsed > 3*time.Second {
		t.Errorf("Expected parallel execution (~2s), but took %v", elapsed)
	}
}
