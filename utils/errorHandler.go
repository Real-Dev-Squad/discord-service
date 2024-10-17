package utils

import (
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

type errorPackage struct{}

var Errors *errorPackage = &errorPackage{}

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
	http.Error(response, `{"success": false, "message": "`+message+`", "status": `+strconv.Itoa(status)+`}`, status)
	return
}
