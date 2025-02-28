package primary

import (
	"context"
	"testing"

	pb "github.com/sarika-p9/my-pipeline-project/api/grpc/proto"
)

func TestCreatePipeline(t *testing.T) {
	server := &PipelineServer{}

	req := &pb.CreatePipelineRequest{Name: "TestPipeline"}
	res, err := server.CreatePipeline(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.Id == "" {
		t.Fatalf("Expected a valid ID, got empty string")
	}
}
