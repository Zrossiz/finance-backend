package handler

import (
	"encoding/json"
	"net/http"
)

func parseJSONBody[T any](r *http.Request) (T, error) {
	var payload T
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&payload); err != nil {
		var zero T
		return zero, err
	}

	return payload, nil
}

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
