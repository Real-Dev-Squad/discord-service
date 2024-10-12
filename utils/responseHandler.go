package utils

import (
	"encoding/json"
	"net/http"
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
