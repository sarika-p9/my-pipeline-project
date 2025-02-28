package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure" // Update with your actual module name
	"github.com/supabase-community/gotrue-go"
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := infrastructure.SupabaseAuth.Signup(r.Context(), gotrue.SignupRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		http.Error(w, "Failed to sign up user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User registered successfully. Please check your email for verification.",
		"user":    user,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
