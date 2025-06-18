package utils

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err:= json.NewEncoder(w).Encode(data); err != nil {
		logrus.Errorf("failed to write response: %v", err)
	}
}