package helpers

import (
	"encoding/json"
	"io"
)

func ParseJSONBody[T any](body io.ReadCloser) (T, error) {
	var payload T
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&payload); err != nil {
		var zero T
		return zero, err
	}

	return payload, nil
}
