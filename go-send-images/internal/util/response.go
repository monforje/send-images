package util

import (
	"encoding/json"
	"net/http"
	"send-images-backend/internal/logger"
)

// JSON отправляет JSON-ответ с указанным статусом.
func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logger.Error("Failed to encode JSON: %v", err)
	}
}

// JSONError — единый стиль для ошибок.
func JSONError(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{
		"error": message,
	})
}
