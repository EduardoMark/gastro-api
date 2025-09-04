package jsonutils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func EncodeJson[T any](w http.ResponseWriter, statusCode int, data T) error {
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("failed to encode json: %v", err)
	}

	return nil
}
