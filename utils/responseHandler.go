package utils

import (
	"net/http"
)

type successPackage struct{}

var Success *successPackage = &successPackage{}

func (Success *successPackage) NewDiscordResponse(response http.ResponseWriter, message string, data interface{}) {
	if data != nil {
		err := Encode(response, data)
		if err != nil {
			Errors.NewInternalError(response)
			return
		}
		response.WriteHeader(http.StatusOK)
	} else {
		response.Write([]byte(`{"success": true, "status": 200, "message": "` + message + `"}`))
	}
}
