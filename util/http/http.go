package http

import (
	"encoding/json"
	"net/http"
)

// Response is a struct
type Response struct {
	Message interface{} `json:"message,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

// JSON is a func
func JSON(w http.ResponseWriter, code int, payload Response) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	w.Write(response)
}

// Text is a func
func Text(w http.ResponseWriter, code int, payload string) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(code)
	w.Write(response)
}
