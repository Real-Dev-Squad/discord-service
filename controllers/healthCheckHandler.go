package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type HealthCheckResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

var WriteResponse = func(data interface{}, response http.ResponseWriter) error {
	return json.NewEncoder(response).Encode(data)
}

func HealthCheckHandler(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	response.Header().Set("Content-Type", "application/json")

	data := HealthCheckResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if err := WriteResponse(data, response); err != nil {
		http.Error(response, `{"status":"error","message":"Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
}
