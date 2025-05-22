package utils

import "encoding/json"

type jsonHandler struct{}

var Json *jsonHandler = &jsonHandler{}

func (Json *jsonHandler) ToJson(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
