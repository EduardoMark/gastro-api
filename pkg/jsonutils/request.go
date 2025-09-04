package jsonutils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DecodeJson[T any](r *http.Request) (T, error) {
	var data T

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, fmt.Errorf("failed decode json: %v", err)
	}

	return data, nil
}
