package utils

import "encoding/json"

type JSONHandler interface {
	ToJson(data interface{}) (string, error)
}

type jsonHandler struct{}

var Json JSONHandler = &jsonHandler{}

func (Json *jsonHandler) ToJson(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
