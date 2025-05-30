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
	response.Header().Set("Content-Type", "application/json")

	data := HealthCheckResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if err := utils.WriteResponse(data, response); err != nil {
		res, _ := utils.Json.ToJson(utils.ErrorResponse{
			Success: false,
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
		http.Error(response, res, http.StatusInternalServerError)
		return
	}
}
