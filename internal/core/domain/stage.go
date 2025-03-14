package domain

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
)

type Stage interface {
	GetID() uuid.UUID
	GetName() string
	Execute(ctx context.Context, pipelineName string, input interface{}) (interface{}, error)
	HandleError(ctx context.Context, err error) error
	Rollback(ctx context.Context, input interface{}) error
}

type StageStatus struct {
	PipelineName string `json:"pipeline_name"`
	StageName    string `json:"stage_name"`
	Status       string `json:"status"`
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

	if s.Status != "Running" {
		s.Status = "Running"
		infrastructure.WebSocket.SendMessage(pipelineName, s.Name, "Running")
	}

	time.Sleep(4 * time.Second)

	if s.Status != "Completed" {
		s.Status = "Completed"
		infrastructure.WebSocket.SendMessage(pipelineName, s.Name, "Completed")
	}

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
