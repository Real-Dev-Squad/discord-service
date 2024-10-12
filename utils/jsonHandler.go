package utils

import "encoding/json"

type jsonHandler struct{}

var Json *jsonHandler = &jsonHandler{}

func (Json *jsonHandler) ToJson(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(bytes)
}
