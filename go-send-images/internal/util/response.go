package util

import (
	"encoding/json"
	"net/http"
)

// JSON sets the response header to "application/json", writes the given status code,
// and encodes the provided payload into JSON format for the http.ResponseWriter.

func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
