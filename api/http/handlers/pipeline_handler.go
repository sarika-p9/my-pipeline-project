package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
	"gorm.io/gorm"
)

type PipelineHandler struct {
	DB *gorm.DB
}

func NewPipelineHandler(db *gorm.DB) *PipelineHandler {
	return &PipelineHandler{DB: db}
}

func (h *PipelineHandler) CreatePipeline(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	pipeline := models.Pipeline{Name: req.Name}
	if err := h.DB.Create(&pipeline).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pipeline"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pipeline created successfully", "pipeline": pipeline})
}

func (h *PipelineHandler) ListPipelines(c *gin.Context) {
	var pipelines []models.Pipeline
	if err := h.DB.Find(&pipelines).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pipelines"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pipelines": pipelines})
}

func GetAPIStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "API is running"})
}

func ListWorkers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"workers": []string{}}) // Update logic later if needed
}
