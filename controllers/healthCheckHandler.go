package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HealthCheckHandler(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte(`{"status": "ok"}`))
}
