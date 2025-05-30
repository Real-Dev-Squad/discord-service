package utils

import (
	"encoding/json"
	"net/http"
)

type successPackage struct{}

var Success *successPackage = &successPackage{}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (Success *successPackage) NewDiscordResponse(response http.ResponseWriter, message string, data interface{}) {
	response.WriteHeader(http.StatusOK)
	if data != nil {
		err := json.NewEncoder(response).Encode(data)
		if err != nil {
			Errors.NewInternalError(response)
			return
		}
	} else {
		res, _:= Json.ToJson(SuccessResponse{
			Success: true,
			Message: message,
			Status:  200,
		})
		response.Write([]byte(res))
	}
}
