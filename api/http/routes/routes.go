package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sarikap9/my-pipeline-project/api/http/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
	}

	return r
}
