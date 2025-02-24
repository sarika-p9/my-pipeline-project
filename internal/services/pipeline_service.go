package services
package models

import (
	"log"

	"gorm.io/gorm"
)

// CreatePipeline inserts a new pipeline into the database.
func CreatePipeline(db *gorm.DB, name string) error {
	pipeline := models.Pipeline{Name: name}

	if err := db.Create(&pipeline).Error; err != nil {
		return err
	}
	log.Println("✅ Pipeline created:", pipeline)
	return nil
}

// GetPipelines retrieves all pipelines.
func GetPipelines(db *gorm.DB) ([]models.Pipeline, error) {
	var pipelines []models.Pipeline
	if err := db.Find(&pipelines).Error; err != nil {
		return nil, err
	}
	log.Printf("✅ Retrieved %d pipelines\n", len(pipelines))
	return pipelines, nil
}
