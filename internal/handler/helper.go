package handler

import (
	"encoding/json"
	"net/http"
)

func JSON(rw http.ResponseWriter, statusCode int, payload interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	if payload != nil {
		json.NewEncoder(rw).Encode(payload)
	}
}

func Error(rw http.ResponseWriter, err HTTPError) {
	JSON(rw, err.Code, map[string]string{"error": err.Message})
}
