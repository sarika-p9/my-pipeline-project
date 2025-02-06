package primary

import (
	"context"
	"log"

	pb "github.com/sarikap9/my-pipeline-project/api/grpc/proto"
)

type PipelineServer struct {
	pb.UnimplementedPipelineServiceServer
	pipelines []string
}

func (s *PipelineServer) CreatePipeline(ctx context.Context, req *pb.CreatePipelineRequest) (*pb.CreatePipelineResponse, error) {
	id := "pipeline-" + req.Name // Simulated ID
	s.pipelines = append(s.pipelines, id)
	log.Printf("Created pipeline: %s", id)
	return &pb.CreatePipelineResponse{Id: id, Message: "Pipeline Created!"}, nil
}

func (s *PipelineServer) ListPipelines(ctx context.Context, req *pb.ListPipelinesRequest) (*pb.ListPipelinesResponse, error) {
	return &pb.ListPipelinesResponse{Pipelines: s.pipelines}, nil
}
