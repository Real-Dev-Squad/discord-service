package utils

import (
	"encoding/json"
	"net/http"
)

var WriteResponse = func(data interface{}, response http.ResponseWriter) error {
	return json.NewEncoder(response).Encode(data)
}