package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sarika-p9/my-pipeline-project/api/http/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/register", handlers.RegisterUser)
		api.POST("/login", handlers.Login)
	}

	return r
}
