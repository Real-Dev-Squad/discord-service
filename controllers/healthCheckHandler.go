package controllers

import (
	"net/http"
	"time"

	"github.com/Real-Dev-Squad/discord-service/queue"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type HealthCheckResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Queue     struct {
		Connected bool   `json:"connected"`
		Error     string `json:"error,omitempty"`
	} `json:"queue"`
}

func HealthCheckHandler(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	response.Header().Set("Content-Type", "application/json")

	// Check queue connection
	queueInstance, err := queue.GetQueueInstance()
	queueStatus := struct {
		Connected bool   `json:"connected"`
		Error     string `json:"error,omitempty"`
	}{
		Connected: queueInstance != nil && queueInstance.Channel != nil,
	}

	if err != nil {
		logrus.Errorf("Queue health check failed: %v", err)
		queueStatus.Error = err.Error()
	}

	data := HealthCheckResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
		Queue:     queueStatus,
	}

	if err := utils.WriteResponse(data, response); err != nil {
		logrus.Errorf("Failed to write health check response: %v", err)
		http.Error(response, `{"status":"error","message":"Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
}
