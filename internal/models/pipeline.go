package models

import (
	"time"

	"github.com/google/uuid"
)

// Pipeline represents a manufacturing pipeline.
type Pipeline struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TableName sets the table name for the Pipeline model.
func (Pipeline) TableName() string {
	return "pipelines"
}
