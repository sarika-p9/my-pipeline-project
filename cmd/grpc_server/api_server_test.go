package main

// import (
// 	"bytes"
// 	"context"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang/mock/gomock"
// 	"github.com/sarika-p9/my-pipeline-project/api/grpc/proto"
// 	"github.com/stretchr/testify/assert"
// )

// func TestCreatePipelineAPI(t *testing.T) {
// 	// Setup Gin
// 	gin.SetMode(gin.TestMode)
// 	router := gin.Default()

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockClient := proto.NewMockPipelineServiceClient(ctrl)

// 	// Replace grpcClient with mock
// 	grpcClient = mockClient

// 	// Define API route
// 	router.POST("/pipelines", func(c *gin.Context) {
// 		// Use a local struct for JSON binding
// 		var newPipeline struct {
// 			Name string `json:"name"`
// 		}

// 		if err := c.BindJSON(&newPipeline); err != nil || newPipeline.Name == "" {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 			return
// 		}

// 		// Call the mocked gRPC method
// 		_, err := grpcClient.CreatePipeline(context.Background(), &proto.CreatePipelineRequest{
// 			Name: newPipeline.Name,
// 		})
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(http.StatusCreated, gin.H{"message": "Pipeline created!"})
// 	})

// 	// Mock gRPC response
// 	mockClient.EXPECT().
// 		CreatePipeline(gomock.Any(), &proto.CreatePipelineRequest{Name: "TestPipeline"}).
// 		Return(&proto.CreatePipelineResponse{Id: "1234", Message: "Pipeline created!"}, nil)

// 	// Create HTTP request
// 	payload := []byte(`{"name": "TestPipeline"}`)
// 	req, _ := http.NewRequest("POST", "/pipelines", bytes.NewBuffer(payload))
// 	req.Header.Set("Content-Type", "application/json")

// 	// Execute request
// 	resp := httptest.NewRecorder()
// 	router.ServeHTTP(resp, req)

// 	// Assertions
// 	assert.Equal(t, http.StatusCreated, resp.Code)
// 	assert.Contains(t, resp.Body.String(), "Pipeline created!")
// }
