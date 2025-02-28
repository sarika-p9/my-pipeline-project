package routes

import (
	"github.com/gorilla/mux"
	"github.com/sarika-p9/my-pipeline-project/internal/controllers" // Update with your actual module name
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/signup", controllers.SignUpHandler).Methods("POST")
}
