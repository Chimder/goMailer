package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json") // Use Set for consistency and clarity
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func WriteError(w http.ResponseWriter, status int, errorMessage interface{}) {
	switch err := errorMessage.(type) {
	case error:
		WriteJSON(w, status, map[string]string{"error": err.Error()})
	case string:
		WriteJSON(w, status, map[string]string{"error": err})
	default:
		return
	}
}

func ParseJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}
	return nil
}
