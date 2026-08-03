package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type DefaultErrorResponse struct {
	Detail  string `json:"detail"`
	Message string `json:"message"`
}

func (d *DefaultErrorResponse) LogErrorMessage(err ...error) {
	if len(err) > 0 {
		log.Print(errors.Join(err...).Error())
	}
}

func (d *DefaultErrorResponse) SendErrorMessage(w http.ResponseWriter, err ...error) {
	if len(err) > 0 {
		strErrors := errors.Join(err...).Error()
		d.Message = d.Message + " " + strErrors
	}
	JsonResponse(w, d, http.StatusInternalServerError)
}

// JsonResponse is a helper function to send JSON responses with a given status code.
func JsonResponse[T any](w http.ResponseWriter, data T, statusCode int) {
	// Encode into a buffer first so we know whether it succeeded before
	// committing to any status code or headers.
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)

	// Response is already committed at this point; a write failure here
	// (e.g. client disconnected) can only be logged, not turned into an
	// error response.
	if _, err := w.Write(buf.Bytes()); err != nil {
		// log.Printf("Error writing JSON response: %v", err)
	}
}
