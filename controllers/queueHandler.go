package controllers

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func QueueHandler(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(response, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	logrus.Infof("QueueHandler: %s\n", string(body))
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

}
