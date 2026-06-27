package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func ValidateRequestMediaType(r *http.Request, expectedType string) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != expectedType {
			return fmt.Errorf("Unsupported Media Type: %s", mediaType)
		}
	}
	return nil
}

func GetJsonBody(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("Failed to decode JSON body: %w", err)
	}
	return nil
}

func JsonResponse(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		err := json.NewEncoder(w).Encode(v)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode JSON response: %v", err), http.StatusInternalServerError)
		}
	}
}
