package utils

import (
	"encoding/json"
	"net/http"

	"github.com/Real-Dev-Squad/discord-service/errors"
	"github.com/sirupsen/logrus"
)

func WriteJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
    
    jsonData, err := json.Marshal(data)
    if err != nil {
        errors.HandleError(w, err)
        return
    }
    
    w.WriteHeader(statusCode)
    if _, err := w.Write(jsonData); err != nil {
        logrus.Errorf("failed to write response: %v", err)
    }
}