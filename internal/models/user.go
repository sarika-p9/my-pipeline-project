package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey"`
	Name      string    `gorm:"type:varchar(100)";not null;default:'Sarika Gautam'`
	Email     string    `gorm:"type:varchar(100);unique;not null"`
	Role      string    `gorm:"type:varchar(20);not null;default:'worker';check:role IN ('super_admin', 'admin', 'manager', 'worker')"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	PipelineExecutions []Pipelines `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type Pipelines struct {
	PipelineID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Status        string    `gorm:"type:varchar(50);not null"`
	PipelineName  string    `gorm:"type:varchar(255);not null;default:'Untitled Pipeline'"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	ExecutionLogs []Stages  `gorm:"foreignKey:PipelineID;constraint:OnDelete:CASCADE;"`
}

type Stages struct {
	StageID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	PipelineID uuid.UUID `gorm:"type:uuid;not null;index"`
	StageName  string    `gorm:"type:varchar(255);not null;default:'Untitled Stage'"`
	Status     string    `gorm:"type:varchar(50);not null"`
	ErrorMsg   string    `gorm:"type:text"`
	Timestamp  time.Time `gorm:"autoCreateTime"`
}
