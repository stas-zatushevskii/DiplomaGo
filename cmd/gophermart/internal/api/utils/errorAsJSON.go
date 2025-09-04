package utils

import "encoding/json"

type errorResponse struct {
	Error string `json:"error"`
}

func ErrorAsJSON(err error) string {
	if err == nil {
		return ""
	}
	resp := errorResponse{Error: err.Error()}
	response, _ := json.Marshal(resp)
	return string(response)
}
