package utils

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type successPackage struct{}

var Success *successPackage = &successPackage{}

func (Success *successPackage) NewDiscordResponse(response http.ResponseWriter, message string, data interface{}) {
	response.WriteHeader(http.StatusOK)
	if data != nil {
		err := json.NewEncoder(response).Encode(data)
		if err != nil {
			Errors.NewInternalError(response)
			return
		}
	} else {
		response.Write([]byte(`{"success": true, "status": 200, "message": "` + message + `"}`))
	}
}

type responseHandler struct{}

var ResponseHandler *responseHandler = &responseHandler{}

func (r *responseHandler) WriteJSON(response http.ResponseWriter, data any, status int) {
	response.Header().Set("Content-Type", "application/json")

	bytes, err := Json.ToJson(data)
	if err != nil {
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response.WriteHeader(status)
	if _, err := response.Write([]byte(bytes)); err != nil {
		logrus.Errorf("Failed to write response: %v", err)
	}
}

