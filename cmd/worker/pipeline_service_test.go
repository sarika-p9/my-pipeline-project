package main

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sarikap9/my-pipeline-project/api/grpc/proto"
	"github.com/stretchr/testify/assert"
)

func TestCreatePipeline(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := proto.NewMockPipelineServiceServer(ctrl)

	// Mock request and expected response
	req := &proto.CreatePipelineRequest{Name: "TestPipeline"}
	res := &proto.CreatePipelineResponse{
		Id:      "1234",
		Message: "Pipeline created successfully",
	}

	// Define expected behavior
	mockService.EXPECT().
		CreatePipeline(gomock.Any(), req).
		Return(res, nil)

	// Call the method
	resp, err := mockService.CreatePipeline(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "1234", resp.Id)
	assert.Equal(t, "Pipeline created successfully", resp.Message)
}
