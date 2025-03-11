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
	GetName() string
	Execute(ctx context.Context, pipelineName string, input interface{}) (interface{}, error)
	HandleError(ctx context.Context, err error) error
	Rollback(ctx context.Context, input interface{}) error
}

type BaseStage struct {
	ID     uuid.UUID
	Name   string
	Status string
}

func NewBaseStage(name string) *BaseStage {
	return &BaseStage{ID: uuid.New(), Name: name, Status: "Pending"}
}

func (s *BaseStage) GetID() uuid.UUID {
	return s.ID
}

func (s *BaseStage) GetName() string {
	return s.Name
}

func (s *BaseStage) Execute(ctx context.Context, pipelineName string, input interface{}) (interface{}, error) {
	log.Printf("Executing stage: %s (%s) for pipeline: %s", s.Name, s.ID, pipelineName)

	// Update status to "Running" and send JSON message with pipelineName and stage name
	s.Status = "Running"
	websocket.Manager.BroadcastMessage(pipelineName+" - "+s.Name, "Running")

	time.Sleep(5 * time.Second) // Simulating execution

	// Update status to "Completed" and send JSON message with pipelineName and stage name
	s.Status = "Completed"
	websocket.Manager.BroadcastMessage(pipelineName+" - "+s.Name, "Completed")

	return input, nil
}

func (s *BaseStage) HandleError(ctx context.Context, err error) error {
	log.Printf("Error in stage %s (%s) execution: %v", s.Name, s.ID, err)
	return errors.New("stage execution failed: " + err.Error())
}

func (s *BaseStage) Rollback(ctx context.Context, input interface{}) error {
	log.Printf("Rolling back stage %s (%s) due to failure. Input: %v", s.Name, s.ID, input)
	return nil
}
