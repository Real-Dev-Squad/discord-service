package utils

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type errorPackage struct{}

var Errors *errorPackage = &errorPackage{}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (errorPackage) NewBadRequestError(response http.ResponseWriter, message string) {
	formatError(response, message, http.StatusBadRequest)
}

func (errorPackage) NewUnauthorisedError(response http.ResponseWriter, message ...string) {
	if len(message) == 0 {
		message = []string{"Unauthorized Access"}
	}
	formatError(response, message[0], http.StatusUnauthorized)
}

func (errorPackage) NewInternalError(response http.ResponseWriter) {
	formatError(response, "Internal Server Error", http.StatusInternalServerError)
}

func formatError(response http.ResponseWriter, message string, status int) {
	logrus.Error("Message : ", message)
	response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	error, _:= Json.ToJson(ErrorResponse{
		Success: false,
		Message: message,
		Status:  status,
	})
	http.Error(response, error, status)
}
