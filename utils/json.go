package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func WriteJSON(w http.ResponseWriter, statusCode int, value any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(value)
}

func WriteErrorJSON(w http.ResponseWriter, statusCode int, err error) {
	WriteJSON(w, statusCode, map[string]string{"error": err.Error()})
}

func ParseJSON(r *http.Request, value any) error {
	if r.Body == nil {
		return fmt.Errorf("MISSING REQUEST BODY")
	}

	return json.NewDecoder(r.Body).Decode(value)
}
