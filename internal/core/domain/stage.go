package domain

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sarika-p9/my-pipeline-project/internal/websocket"
)

type Stage interface {
	GetID() uuid.UUID
	Execute(ctx context.Context, pipelineName string, input interface{}) (interface{}, error)
	HandleError(ctx context.Context, err error) error
	Rollback(ctx context.Context, input interface{}) error
}

type BaseStage struct {
	ID     uuid.UUID
	Status string
}

func NewBaseStage() *BaseStage {
	return &BaseStage{ID: uuid.New(), Status: "Pending"}
}

func (s *BaseStage) GetID() uuid.UUID {
	return s.ID
}

func (s *BaseStage) Execute(ctx context.Context, pipelineName string, input interface{}) (interface{}, error) {
	log.Printf("Executing stage: %s for pipeline: %s", s.ID, pipelineName)

	// Update status to "Running" and send JSON message with pipelineName
	s.Status = "Running"
	websocket.Manager.BroadcastMessage(pipelineName+" - "+s.ID.String(), "Running") // ✅ Include pipelineName

	time.Sleep(5 * time.Second) // Simulating execution

	// Update status to "Completed" and send JSON message with pipelineName
	s.Status = "Completed"
	websocket.Manager.BroadcastMessage(pipelineName+" - "+s.ID.String(), "Completed") // ✅ Include pipelineName

	return input, nil
}

func (s *BaseStage) HandleError(ctx context.Context, err error) error {
	log.Printf("Error in stage %s execution: %v", s.ID, err)
	return errors.New("stage execution failed: " + err.Error())
}

func (s *BaseStage) Rollback(ctx context.Context, input interface{}) error {
	log.Printf("Rolling back stage %s due to failure. Input: %v", s.ID, input)
	return nil
}
