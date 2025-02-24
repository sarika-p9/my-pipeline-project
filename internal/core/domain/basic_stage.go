package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type BasicStage struct {
	ID   uuid.UUID
	Name string
}

func NewBasicStage(name string) *BasicStage {
	return &BasicStage{
		ID:   uuid.New(),
		Name: name,
	}
}

func (s *BasicStage) GetID() uuid.UUID {
	return s.ID
}

func (s *BasicStage) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	fmt.Printf("Executing stage: %s with input: %v\n", s.Name, input)
	// Simulate success
	return fmt.Sprintf("%v_processed_by_%s", input, s.Name), nil
}

func (s *BasicStage) HandleError(ctx context.Context, err error) error {
	fmt.Printf("Error in stage %s: %v\n", s.Name, err)
	return nil
}

func (s *BasicStage) Rollback(ctx context.Context, input interface{}) error {
	fmt.Printf("Rolling back stage: %s for input: %v\n", s.Name, input)
	return nil
}
