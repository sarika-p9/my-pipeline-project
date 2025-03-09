package routes

import (
	"github.com/sarika-p9/my-pipeline-project/internal/websocket"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/ws", func(c *gin.Context) {
		websocket.Manager.HandleConnections(c.Writer, c.Request)
	})
}
