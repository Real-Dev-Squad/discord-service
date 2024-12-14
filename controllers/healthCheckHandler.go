package controllers

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthCheckResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func HealthCheckHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	data := HealthCheckResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if err := json.NewEncoder(response).Encode(data); err != nil {
		http.Error(response, `{"status":"error","message":"Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
}
