package services

import (
	"context"
	"log"

	"github.com/sarika-p9/my-pipeline-project/api/grpc/proto"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
	"gorm.io/gorm"
)

// PipelineService implements the gRPC PipelineServiceServer interface
type PipelineService struct {
	proto.UnimplementedPipelineServiceServer
	DB *gorm.DB
}

// NewPipelineService initializes a new PipelineService
func NewPipelineService(db *gorm.DB) *PipelineService {
	return &PipelineService{DB: db}
}

// CreatePipeline handles pipeline creation
func (s *PipelineService) CreatePipeline(ctx context.Context, req *proto.CreatePipelineRequest) (*proto.CreatePipelineResponse, error) {
	pipeline := models.Pipeline{Name: req.Name}

	if err := s.DB.Create(&pipeline).Error; err != nil {
		log.Println("❌ Failed to create pipeline:", err)
		return nil, err
	}

	log.Println("✅ Pipeline created:", pipeline)
	return &proto.CreatePipelineResponse{Message: "Pipeline created successfully"}, nil
}

// ListPipelines retrieves all pipelines
func (s *PipelineService) ListPipelines(ctx context.Context, req *proto.ListPipelinesRequest) (*proto.ListPipelinesResponse, error) {
	var pipelines []models.Pipeline
	if err := s.DB.Find(&pipelines).Error; err != nil {
		log.Println("❌ Failed to fetch pipelines:", err)
		return nil, err
	}

	// ✅ Fix: Correctly initialize proto.Pipeline instances
	var protoPipelines []*proto.Pipeline
	for _, p := range pipelines {
		protoPipelines = append(protoPipelines, &proto.Pipeline{Name: p.Name}) // ✅ No more "undefined: proto.Pipeline"
	}

	log.Printf("✅ Retrieved %d pipelines\n", len(pipelines))
	return &proto.ListPipelinesResponse{Pipelines: protoPipelines}, nil
}
