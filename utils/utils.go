package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, fmt.Sprintf("error encoding JSON: %v", err), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(data)
	return err

}

func WriteError(w http.ResponseWriter, status int, errorMessage string, err error) {
	var errMessage string
	if err != nil {
		errMessage = fmt.Sprintf("%s: %v", errorMessage, err)
	} else {
		errMessage = errorMessage
	}
	WriteJSON(w, status, map[string]string{"error": strings.TrimSpace(errMessage)})
}

func ParseJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	defer r.Body.Close()

	r.Body = http.MaxBytesReader(nil, r.Body, 1048576)

	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}
	return nil
}
