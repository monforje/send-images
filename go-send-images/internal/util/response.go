package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"send-images-backend/internal/logger"
)

// JSON отправляет JSON-ответ с указанным статусом и логирует.
func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logger.Error("Failed to encode JSON: %v", err)
	} else {
		logger.Debug("Response status %d, payload: %+v", status, payload)
	}
}

// JSONError — единый стиль для ошибок (и возврат как error).
func JSONError(w http.ResponseWriter, status int, message string) error {
	errObj := map[string]string{"error": message}
	JSON(w, status, errObj)
	return fmt.Errorf("http %d: %s", status, message)
}
