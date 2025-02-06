package main

import (
	"net/http"

	"my-pipeline-project/pipeline" // Ensure this matches your module name from `go.mod`

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "API Server is running!"})
	})

	r.GET("/pipelines", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"pipelines": []string{}})
	})

	r.POST("/pipelines", func(c *gin.Context) {
		pipeline.CreatePipeline("TestPipeline") // Call the function from `pipeline` package
		c.JSON(http.StatusCreated, gin.H{"message": "Pipeline created!"})
	})

	r.GET("/workers", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"workers": []string{}})
	})

	r.Run(":8080")
}
