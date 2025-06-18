package controllers

import (
	"net/http"
	"time"

	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/julienschmidt/httprouter"
)

type HealthCheckResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func HealthCheckHandler(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	data := HealthCheckResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	utils.WriteJSONResponse(response, http.StatusOK, data)
}
